package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	panutils "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/scheduler"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// userResourceSubmitResult 单条提交/重检的内部结果
type userResourceSubmitResult struct {
	UserResource *entity.UserResource
	// Undetermined 检测通道启用但未得出结论（pending 的一种，提示语不同）
	Undetermined bool
}

// userResourceError 携带 HTTP 语义的提交错误（单条直接返回，批量逐条聚合）
type userResourceError struct {
	Code    int
	Message string
}

func (e *userResourceError) Error() string { return e.Message }

// SubmitUserResource 提交资源（单条，同步检测）
func SubmitUserResource(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		URL         string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "标题和链接不能为空", http.StatusBadRequest)
		return
	}

	userID := c.GetUint("user_id")

	result, err := submitSingleUserResource(userID, req.Title, req.Description, req.URL)
	if err != nil {
		if ure, ok := err.(*userResourceError); ok {
			ErrorResponse(c, ure.Message, ure.Code)
			return
		}
		utils.Error("SubmitUserResource - 内部错误: %v", err)
		ErrorResponse(c, "提交失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	ur := result.UserResource
	message := "提交成功"
	switch ur.Status {
	case entity.UserResourceStatusProcessing:
		message = "提交成功，检测有效，已进入处理队列"
	case entity.UserResourceStatusInvalid:
		message = "提交成功，但链接检测无效：" + ur.FailReason
	case entity.UserResourceStatusPending:
		if result.Undetermined {
			message = "提交成功，检测通道暂未得出结论，可稍后重新检测"
		} else {
			message = "提交成功，检测通道未启用，资源暂未检测"
		}
	}

	SuccessResponse(c, gin.H{
		"message":       message,
		"id":            ur.ID,
		"status":        ur.Status,
		"pan_id":        ur.PanID,
		"check_time":    ur.CheckTime,
		"check_method":  ur.CheckMethod,
		"fail_reason":   ur.FailReason,
		"will_publish":  ur.Status == entity.UserResourceStatusProcessing,
	})
}

// RecheckUserResource 手动重新检测（仅 pending/invalid 可重检，FR-011）
func RecheckUserResource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}
	userID := c.GetUint("user_id")

	ur, err := repoManager.UserResourceRepository.FindByIDAndUser(uint(id), userID)
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}
	if ur.Status != entity.UserResourceStatusPending && ur.Status != entity.UserResourceStatusInvalid {
		ErrorResponse(c, "当前状态不允许重新检测", http.StatusBadRequest)
		return
	}

	result, err := recheckUserResource(ur)
	if err != nil {
		if ure, ok := err.(*userResourceError); ok {
			ErrorResponse(c, ure.Message, ure.Code)
			return
		}
		utils.Error("RecheckUserResource - 内部错误: %v", err)
		ErrorResponse(c, "重新检测失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	updated := result.UserResource
	SuccessResponse(c, gin.H{
		"id":          updated.ID,
		"status":      updated.Status,
		"check_time":  updated.CheckTime,
		"fail_reason": updated.FailReason,
	})
}

// GetUserResources 我的资源列表（分页 + 状态/关键词过滤 + 五态统计，FR-005/FR-006）
func GetUserResources(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}
	status := c.Query("status")
	keyword := c.Query("keyword")

	list, total, err := repoManager.UserResourceRepository.ListByUser(userID, page, pageSize, status, keyword)
	if err != nil {
		utils.Error("GetUserResources - 查询失败: %v", err)
		ErrorResponse(c, "获取资源列表失败", http.StatusInternalServerError)
		return
	}

	stats, err := repoManager.UserResourceRepository.GetStatsByUser(userID)
	if err != nil {
		utils.Error("GetUserResources - 统计失败: %v", err)
		stats = map[string]int64{}
	}

	SuccessResponse(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"stats":     stats,
	})
}

// DeleteUserResource 删除我的资源（归属校验统一 404 防探测，FR-006/FR-007）
func DeleteUserResource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}
	userID := c.GetUint("user_id")

	ur, err := repoManager.UserResourceRepository.FindByIDAndUser(uint(id), userID)
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// processing 状态：尽力移除仍在队列中的来源记录（已被消费则仅删用户记录）
	if ur.Status == entity.UserResourceStatusProcessing {
		if _, err := repoManager.ReadyResourceRepository.DeleteByUserResourceID(ur.ID); err != nil {
			utils.Warn("DeleteUserResource - 移除队列记录失败（忽略）: id=%d, %v", ur.ID, err)
		}
	}

	if err := repoManager.UserResourceRepository.Delete(ur.ID); err != nil {
		utils.Error("DeleteUserResource - 删除失败: %v", err)
		ErrorResponse(c, "删除失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, nil)
}

// BatchSubmitUserResources 批量提交（逐条处理、逐条返回结果，FR-010；单次上限与整体拒绝语义见契约 §2）
func BatchSubmitUserResources(c *gin.Context) {
	const maxBatchItems = 50

	var req struct {
		Items []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "items 不能为空", http.StatusBadRequest)
		return
	}
	if len(req.Items) == 0 {
		ErrorResponse(c, "items 不能为空", http.StatusBadRequest)
		return
	}
	if len(req.Items) > maxBatchItems {
		ErrorResponse(c, fmt.Sprintf("单次最多提交 %d 条", maxBatchItems), http.StatusBadRequest)
		return
	}

	userID := c.GetUint("user_id")

	// 整体拒绝：当日余量不足且一条都无法处理（FR-009）
	limit := readUserUploadDailyLimit()
	if limit > 0 {
		count, err := repoManager.UserResourceRepository.CountTodayByUser(userID)
		if err != nil {
			ErrorResponse(c, "提交失败，请稍后重试", http.StatusInternalServerError)
			return
		}
		if count >= int64(limit) {
			ErrorResponse(c, fmt.Sprintf("今日提交已达上限（%d 条），请明天再试", limit), http.StatusTooManyRequests)
			return
		}
	}

	type batchResult struct {
		Index   int    `json:"index"`
		Title   string `json:"title"`
		OK      bool   `json:"ok"`
		ID      uint   `json:"id,omitempty"`
		Status  string `json:"status,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	results := make([]batchResult, 0, len(req.Items))
	successCount, failCount := 0, 0

	for i, item := range req.Items {
		result, err := submitSingleUserResource(userID, item.Title, item.Description, item.URL)
		br := batchResult{Index: i, Title: item.Title}
		if err != nil {
			if ure, ok := err.(*userResourceError); ok {
				br.Reason = ure.Message
			} else {
				utils.Error("BatchSubmitUserResources - 第 %d 条内部错误: %v", i, err)
				br.Reason = "提交失败，请稍后重试"
			}
			failCount++
		} else {
			br.OK = true
			br.ID = result.UserResource.ID
			br.Status = result.UserResource.Status
			successCount++
		}
		results = append(results, br)
	}

	SuccessResponse(c, gin.H{
		"results":       results,
		"success_count": successCount,
		"fail_count":    failCount,
	})
}

// submitSingleUserResource 单条提交共享逻辑（提交与批量复用，US3）
func submitSingleUserResource(userID uint, title, description, rawURL string) (*userResourceSubmitResult, error) {
	// 平台识别（D2：与自动处理通道同一识别入口）
	shareID, serviceType := panutils.ExtractShareId(rawURL)
	if serviceType == panutils.NotFound || shareID == "" {
		return nil, &userResourceError{Code: http.StatusBadRequest, Message: "不支持的平台或链接格式"}
	}

	// 公开库查重（FR-008）
	exists, err := repoManager.ResourceRepository.FindExists(rawURL, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &userResourceError{Code: http.StatusConflict, Message: "该资源已收录"}
	}

	// 本人查重（FR-008）
	dup, err := repoManager.UserResourceRepository.FindByUserAndURL(userID, rawURL)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if dup != nil {
		return nil, &userResourceError{Code: http.StatusConflict, Message: "您已提交过该资源"}
	}

	// 频率限制（FR-009，0 表示不限制；配置缺失时回退默认值，兼容存量库未插入新配置行的场景）
	limit := readUserUploadDailyLimit()
	if limit > 0 {
		count, err := repoManager.UserResourceRepository.CountTodayByUser(userID)
		if err != nil {
			return nil, err
		}
		if count >= int64(limit) {
			return nil, &userResourceError{Code: http.StatusTooManyRequests, Message: fmt.Sprintf("今日提交已达上限（%d 条），请明天再试", limit)}
		}
	}

	// 创建用户资源（pending 起步）
	ur := &entity.UserResource{
		UserID:      userID,
		Title:       title,
		Description: description,
		URL:         rawURL,
		Status:      entity.UserResourceStatusPending,
	}
	if err := repoManager.UserResourceRepository.Create(ur); err != nil {
		return nil, err
	}

	// 同步检测并推进状态（D3）
	if err := checkAndAdvanceUserResource(ur); err != nil {
		return nil, err
	}

	undetermined := ur.Status == entity.UserResourceStatusPending &&
		ur.CheckMethod == entity.UserResourceCheckMethodPanCheck
	return &userResourceSubmitResult{UserResource: ur, Undetermined: undetermined}, nil
}

// recheckUserResource 重检已有记录（复用检测推进逻辑）
func recheckUserResource(ur *entity.UserResource) (*userResourceSubmitResult, error) {
	ur.FailReason = ""
	if err := repoManager.UserResourceRepository.Update(ur); err != nil {
		return nil, err
	}
	if err := checkAndAdvanceUserResource(ur); err != nil {
		return nil, err
	}
	undetermined := ur.Status == entity.UserResourceStatusPending &&
		ur.CheckMethod == entity.UserResourceCheckMethodPanCheck
	return &userResourceSubmitResult{UserResource: ur, Undetermined: undetermined}, nil
}

// checkAndAdvanceUserResource 同步调用检测通道并推进状态：
// valid → 入 ReadyResource 队列（processing）；invalid → invalid+原因；其余保持 pending（FR-004/FR-007/D3）
func checkAndAdvanceUserResource(ur *entity.UserResource) error {
	now := utils.GetCurrentTime()
	ur.CheckTime = &now

	svc := scheduler.GetGlobalLinkCheckService()
	if svc == nil {
		ur.CheckMethod = entity.UserResourceCheckMethodDisabled
		return repoManager.UserResourceRepository.Update(ur)
	}

	result := svc.CheckURL(context.Background(), ur.URL, false)
	switch {
	case result.Status == "valid":
		ur.CheckMethod = entity.UserResourceCheckMethodPanCheck
		if err := enqueueUserResource(ur); err != nil {
			return err
		}
	case result.Status == "invalid":
		ur.CheckMethod = entity.UserResourceCheckMethodPanCheck
		ur.Status = entity.UserResourceStatusInvalid
		ur.FailReason = result.FailReason
	default:
		// undetermined：区分通道未启用与未得出结论（Edge Cases）
		if result.DetectionMethod == entity.UserResourceCheckMethodDisabled {
			ur.CheckMethod = entity.UserResourceCheckMethodDisabled
		} else {
			ur.CheckMethod = entity.UserResourceCheckMethodPanCheck
		}
	}

	return repoManager.UserResourceRepository.Update(ur)
}

// enqueueUserResource 检测有效后写入 ReadyResource 队列（D4），
// Source=user_upload 标记来源，Extra 携带 user_resource_id 供 scheduler 回写
func enqueueUserResource(ur *entity.UserResource) error {
	panID := resolvePanID(ur.URL)
	ur.PanID = panID

	ready := &entity.ReadyResource{
		Title:       &ur.Title,
		Description: ur.Description,
		URL:         ur.URL,
		Source:      "user_upload",
		Extra:       fmt.Sprintf("%d", ur.ID),
	}
	if err := repoManager.ReadyResourceRepository.Create(ready); err != nil {
		return err
	}

	ur.Status = entity.UserResourceStatusProcessing
	return nil
}

// readUserUploadDailyLimit 读取每用户每日提交上限（FR-009）。
// 配置行缺失或读取失败时回退默认值 50（存量生产库 MIGRATE=false 不会自动插入新配置行）。
func readUserUploadDailyLimit() int {
	limit, err := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyUserUploadDailyLimit)
	if err != nil || limit < 0 {
		var fallback int
		fmt.Sscanf(entity.ConfigDefaultUserUploadDailyLimit, "%d", &fallback)
		return fallback
	}
	return limit
}

// resolvePanID 按 URL 识别平台并解析 pans.id（复用 FindIdByServiceType，与处理通道一致的回退语义）
func resolvePanID(rawURL string) *uint {
	_, serviceType := panutils.ExtractShareId(rawURL)
	if serviceType == panutils.NotFound {
		return nil
	}
	id, err := repoManager.PanRepository.FindIdByServiceType(serviceType.String())
	if err != nil || id <= 0 {
		utils.Warn("resolvePanID - 未找到服务类型 %s 的平台映射", serviceType.String())
		return nil
	}
	panID := uint(id)
	return &panID
}
