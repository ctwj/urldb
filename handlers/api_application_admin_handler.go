package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// API 申请审核处理器（016-api-access-application，管理员端 B1-B4）

// GetApiApplicationStats 申请状态统计（管理端统计卡）
func GetApiApplicationStats(c *gin.Context) {
	counts, err := repoManager.ApiApplicationRepository.CountByStatus()
	if err != nil {
		utils.Error("GetApiApplicationStats - 统计查询失败 - error: %v", err)
		ErrorResponse(c, "获取统计失败", http.StatusInternalServerError)
		return
	}
	SuccessResponse(c, gin.H{
		"pending":  counts[entity.ApiApplicationStatusPending],
		"approved": counts[entity.ApiApplicationStatusApproved],
		"rejected": counts[entity.ApiApplicationStatusRejected],
	})
}

// ListApiApplications 申请列表（B1：status 过滤默认 pending，keyword 匹配邮箱/用户名/用途，分页）
func ListApiApplications(c *gin.Context) {
	status := c.DefaultQuery("status", "pending")
	keyword := strings.TrimSpace(c.Query("keyword"))

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	apps, total, err := repoManager.ApiApplicationRepository.ListByStatus(status, keyword, page, pageSize)
	if err != nil {
		utils.Error("ListApiApplications - 列表查询失败 - error: %v", err)
		ErrorResponse(c, "获取申请列表失败", http.StatusInternalServerError)
		return
	}

	// 联用户名（列表量级小，map 缓存避免重复查询）
	userCache := make(map[uint]string)
	credCache := make(map[uint]*entity.ApiCredential)
	items := make([]dto.ApiApplicationAdminItem, 0, len(apps))
	for _, app := range apps {
		username, ok := userCache[app.UserID]
		if !ok {
			if user, err := repoManager.UserRepository.FindByID(app.UserID); err == nil && user != nil {
				username = user.Username
			}
			userCache[app.UserID] = username
		}
		item := dto.ApiApplicationAdminItem{
			ID:           app.ID,
			UserID:       app.UserID,
			Username:     username,
			Email:        app.Email,
			Purpose:      app.Purpose,
			Status:       app.Status,
			RejectReason: app.RejectReason,
			ReviewerName: app.ReviewerName,
			ReviewedAt:   app.ReviewedAt,
			CreatedAt:    app.CreatedAt,
		}
		// 已同意申请附带凭证状态与到期时间（列表展示 + 有效期调整入口）
		if app.Status == entity.ApiApplicationStatusApproved {
			cred, ok := credCache[app.UserID]
			if !ok {
				cred, _ = repoManager.ApiCredentialRepository.FindByUserID(app.UserID)
				credCache[app.UserID] = cred
			}
			if cred != nil {
				item.CredentialStatus = cred.Status
				item.ExpiresAt = cred.ExpiresAt
				item.CredentialExpired = cred.IsExpired(utils.GetCurrentTime())
			}
		}
		items = append(items, item)
	}

	ListResponse(c, items, total)
}

// parseUintParam 解析路由 uint 参数
func parseUintParam(c *gin.Context, name string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ApproveApiApplication 同意申请（B2，FR-005：签发/重建凭证）
func ApproveApiApplication(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		ErrorResponse(c, "无效的申请ID", http.StatusBadRequest)
		return
	}
	reviewerID := c.GetUint("user_id")
	reviewerName, _ := c.Get("username")
	reviewerNameStr, _ := reviewerName.(string)

	app, expiresAt, err := services.ApproveApplication(id, reviewerID, reviewerNameStr)
	if err != nil {
		if errors.Is(err, services.ErrApplicationNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, "申请不存在", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrApplicationNotPending) {
			ErrorResponse(c, "该申请不在待审核状态", http.StatusBadRequest)
			return
		}
		utils.Error("ApproveApiApplication - 同意失败 - appID: %d, error: %v", id, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"application":           app,
		"credential_expires_at": expiresAt,
	})
}

// RejectApiApplication 拒绝申请（B3，理由可选）
func RejectApiApplication(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		ErrorResponse(c, "无效的申请ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 理由可选：允许空 body
		req.Reason = ""
	}

	reviewerID := c.GetUint("user_id")
	reviewerName, _ := c.Get("username")
	reviewerNameStr, _ := reviewerName.(string)

	app, err := services.RejectApplication(id, reviewerID, reviewerNameStr, req.Reason)
	if err != nil {
		if errors.Is(err, services.ErrApplicationNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, "申请不存在", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrApplicationNotPending) {
			ErrorResponse(c, "该申请不在待审核状态", http.StatusBadRequest)
			return
		}
		utils.Error("RejectApiApplication - 拒绝失败 - appID: %d, error: %v", id, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, app)
}

// DisableApiCredential 停用用户 API 权限（B4，FR-014）
func DisableApiCredential(c *gin.Context) {
	userID, err := parseUintParam(c, "userId")
	if err != nil {
		ErrorResponse(c, "无效的用户ID", http.StatusBadRequest)
		return
	}

	if err := services.DisableUserCredential(userID); err != nil {
		if errors.Is(err, services.ErrCredentialNotFound) {
			ErrorResponse(c, "该用户没有有效的 API 权限", http.StatusNotFound)
			return
		}
		utils.Error("DisableApiCredential - 停用失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	utils.Info("DisableApiCredential - 权限已停用 - userID: %d", userID)
	SuccessResponse(c, gin.H{"disabled": true})
}

// EnableApiCredential 启用用户 API 权限（B6，恢复被停用的凭证）
func EnableApiCredential(c *gin.Context) {
	userID, err := parseUintParam(c, "userId")
	if err != nil {
		ErrorResponse(c, "无效的用户ID", http.StatusBadRequest)
		return
	}

	if err := services.EnableUserCredential(userID); err != nil {
		if errors.Is(err, services.ErrCredentialNotFound) {
			ErrorResponse(c, "该用户没有已停用的 API 凭证", http.StatusNotFound)
			return
		}
		utils.Error("EnableApiCredential - 启用失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	adminName, _ := c.Get("username")
	utils.Info("EnableApiCredential - 权限已启用 - 管理员: %v, userID: %d", adminName, userID)
	SuccessResponse(c, gin.H{"enabled": true})
}

// UpdateApiCredentialExpiry 调整用户凭证到期时间（B5）
// body: { "expires_at": "RFC3339" }；空/缺省 = 永久
func UpdateApiCredentialExpiry(c *gin.Context) {
	userID, err := parseUintParam(c, "userId")
	if err != nil {
		ErrorResponse(c, "无效的用户ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ExpiresAt *string `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && strings.TrimSpace(*req.ExpiresAt) != "" {
		t, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpiresAt))
		if err != nil {
			ErrorResponse(c, "到期时间格式无效", http.StatusBadRequest)
			return
		}
		expiresAt = &t
	}

	cred, err := services.UpdateUserCredentialExpiry(userID, expiresAt)
	if err != nil {
		if errors.Is(err, services.ErrCredentialNotFound) {
			ErrorResponse(c, "该用户没有 API 凭证", http.StatusNotFound)
			return
		}
		utils.Error("UpdateApiCredentialExpiry - 调整失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	adminName, _ := c.Get("username")
	utils.Info("UpdateApiCredentialExpiry - 管理员调整到期时间 - 管理员: %v, userID: %d, expiresAt: %v", adminName, userID, cred.ExpiresAt)
	SuccessResponse(c, gin.H{
		"credential_status": cred.Status,
		"expires_at":        cred.ExpiresAt,
		"credential_expired": cred.IsExpired(utils.GetCurrentTime()),
	})
}
