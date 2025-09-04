package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	pan "github.com/ctwj/urldb/common"
	commonutils "github.com/ctwj/urldb/common/utils"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// GetResources 获取资源列表
func GetResources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	utils.Info("资源列表请求 - page: %d, pageSize: %d, User-Agent: %s", page, pageSize, c.GetHeader("User-Agent"))

	// 添加缓存控制头，优化 SSR 性能
	c.Header("Cache-Control", "public, max-age=30") // 30秒缓存，平衡性能和实时性
	c.Header("ETag", fmt.Sprintf("resources-%d-%d-%s-%s", page, pageSize, c.Query("search"), c.Query("pan_id")))

	params := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	if search := c.Query("search"); search != "" {
		params["search"] = search
	}
	if panID := c.Query("pan_id"); panID != "" {
		if id, err := strconv.ParseUint(panID, 10, 32); err == nil {
			params["pan_id"] = uint(id)
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		utils.Info("收到分类ID参数: %s", categoryID)
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			params["category_id"] = uint(id)
			utils.Info("解析分类ID成功: %d", uint(id))
		} else {
			utils.Error("解析分类ID失败: %v", err)
		}
	}
	if hasSaveURL := c.Query("has_save_url"); hasSaveURL != "" {
		if hasSaveURL == "true" {
			params["has_save_url"] = true
		} else if hasSaveURL == "false" {
			params["has_save_url"] = false
		}
	}
	if noSaveURL := c.Query("no_save_url"); noSaveURL != "" {
		if noSaveURL == "true" {
			params["no_save_url"] = true
		}
	}
	if panName := c.Query("pan_name"); panName != "" {
		params["pan_name"] = panName
	}

	resources, total, err := repoManager.ResourceRepository.SearchWithFilters(params)

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"data":      converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetResourceByID 根据ID获取资源
func GetResourceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	response := converter.ToResourceResponse(resource)
	SuccessResponse(c, response)
}

// CheckResourceExists 检查资源是否存在（测试FindExists函数）
func CheckResourceExists(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		ErrorResponse(c, "URL参数不能为空", http.StatusBadRequest)
		return
	}

	excludeIDStr := c.Query("exclude_id")
	var excludeID uint
	if excludeIDStr != "" {
		if id, err := strconv.ParseUint(excludeIDStr, 10, 32); err == nil {
			excludeID = uint(id)
		}
	}

	exists, err := repoManager.ResourceRepository.FindExists(url, excludeID)
	if err != nil {
		ErrorResponse(c, "检查失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"url":    url,
		"exists": exists,
	})
}

// CreateResource 创建资源
func CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource := &entity.Resource{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		PanID:       req.PanID,
		SaveURL:     req.SaveURL,
		FileSize:    req.FileSize,
		CategoryID:  req.CategoryID,
		IsValid:     req.IsValid,
		IsPublic:    req.IsPublic,
		Cover:       req.Cover,
		Author:      req.Author,
		ErrorMsg:    req.ErrorMsg,
	}

	err := repoManager.ResourceRepository.Create(resource)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{
		"message":  "资源创建成功",
		"resource": converter.ToResourceResponse(resource),
	})
}

// UpdateResource 更新资源
func UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// 更新资源信息
	if req.Title != "" {
		resource.Title = req.Title
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.URL != "" {
		resource.URL = req.URL
	}
	if req.PanID != nil {
		resource.PanID = req.PanID
	}
	if req.SaveURL != "" {
		resource.SaveURL = req.SaveURL
	}
	if req.FileSize != "" {
		resource.FileSize = req.FileSize
	}
	if req.CategoryID != nil {
		resource.CategoryID = req.CategoryID
	}
	resource.IsValid = req.IsValid
	resource.IsPublic = req.IsPublic
	if req.Cover != "" {
		resource.Cover = req.Cover
	}
	if req.Author != "" {
		resource.Author = req.Author
	}
	if req.ErrorMsg != "" {
		resource.ErrorMsg = req.ErrorMsg
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = repoManager.ResourceRepository.Update(resource)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{"message": "资源更新成功"})
}

// DeleteResource 删除资源
func DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.ResourceRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "资源删除成功"})
}

// SearchResources 搜索资源
func SearchResources(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var resources []entity.Resource
	var total int64
	var err error

	if query == "" {
		// 搜索关键词为空时，返回最新记录（分页）
		resources, total, err = repoManager.ResourceRepository.FindWithRelationsPaginated(page, pageSize)
	} else {
		// 有搜索关键词时，执行搜索
		resources, total, err = repoManager.ResourceRepository.Search(query, nil, page, pageSize)
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"resources": converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 增加资源浏览次数
func IncrementResourceViewCount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}

	// 增加资源访问量
	err = repoManager.ResourceRepository.IncrementViewCount(uint(id))
	if err != nil {
		ErrorResponse(c, "增加浏览次数失败", http.StatusInternalServerError)
		return
	}

	// 记录访问记录
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	err = repoManager.ResourceViewRepository.RecordView(uint(id), ipAddress, userAgent)
	if err != nil {
		// 记录访问失败不影响主要功能，只记录日志
		utils.Error("记录资源访问失败: %v", err)
	}

	SuccessResponse(c, gin.H{"message": "浏览次数+1"})
}

// BatchDeleteResources 批量删除资源
func BatchDeleteResources(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		ErrorResponse(c, "参数错误", 400)
		return
	}
	count := 0
	for _, id := range req.IDs {
		if err := repoManager.ResourceRepository.Delete(id); err == nil {
			count++
		}
	}
	SuccessResponse(c, gin.H{"deleted": count, "message": "批量删除成功"})
}

// GetResourceLink 获取资源链接（智能转存）
func GetResourceLink(c *gin.Context) {
	// 获取资源ID
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}

	utils.Info("获取资源链接请求 - resourceID: %d", resourceID)

	// 查询资源信息
	resource, err := repoManager.ResourceRepository.FindByID(uint(resourceID))
	if err != nil {
		utils.Error("查询资源失败: %v", err)
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// 查询平台信息
	var panInfo entity.Pan
	if resource.PanID != nil {
		panPtr, err := repoManager.PanRepository.FindByID(*resource.PanID)
		if err != nil {
			utils.Error("查询平台信息失败: %v", err)
		} else if panPtr != nil {
			panInfo = *panPtr
		}
	}

	utils.Info("资源信息 - 平台: %s, 原始链接: %s, 转存链接: %s", panInfo.Name, resource.URL, resource.SaveURL)

	// 统计访问次数
	err = repoManager.ResourceRepository.IncrementViewCount(uint(resourceID))
	if err != nil {
		utils.Error("增加资源访问量失败: %v", err)
	}

	// 记录访问记录
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	err = repoManager.ResourceViewRepository.RecordView(uint(resourceID), ipAddress, userAgent)
	if err != nil {
		utils.Error("记录资源访问失败: %v", err)
	}

	// 如果不是夸克网盘，直接返回原链接
	if panInfo.Name != "quark" && panInfo.Name != "xunlei" {
		utils.Info("非夸克和迅雷资源，直接返回原链接")
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
		})
		return
	}

	// 如果已存在转存链接，直接返回
	if resource.SaveURL != "" {
		utils.Info("已存在转存链接，直接返回: %s", resource.SaveURL)
		SuccessResponse(c, gin.H{
			"url":         resource.SaveURL,
			"type":        "transferred",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
		})
		return
	}

	// 检查是否开启自动转存
	autoTransferEnabled, err := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
	if err != nil {
		utils.Error("获取自动转存配置失败: %v", err)
		// 配置获取失败，返回原链接
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "",
		})
		return
	}

	if !autoTransferEnabled {
		utils.Info("自动转存功能未开启，返回原链接")
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "",
		})
		return
	}

	// 执行自动转存
	utils.Info("开始执行自动转存")
	transferResult := performAutoTransfer(resource)

	if transferResult.Success {
		utils.Info("自动转存成功，返回转存链接: %s", transferResult.SaveURL)
		SuccessResponse(c, gin.H{
			"url":         transferResult.SaveURL,
			"type":        "transferred",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "资源易和谐，请及时用手机夸克扫码转存",
		})
	} else {
		utils.Error("自动转存失败: %s", transferResult.ErrorMsg)
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "",
		})
	}
}

// TransferResult 转存结果
type TransferResult struct {
	Success  bool   `json:"success"`
	Fid      string `json:"fid"`
	SaveURL  string `json:"save_url"`
	ErrorMsg string `json:"error_msg"`
}

// performAutoTransfer 执行自动转存
func performAutoTransfer(resource *entity.Resource) TransferResult {
	utils.Info("开始执行资源转存 - ID: %d, URL: %s", resource.ID, resource.URL)

	// 平台ID
	panID := resource.PanID

	// 获取可用的夸克账号
	accounts, err := repoManager.CksRepository.FindByPanID(*panID)
	if err != nil {
		utils.Error("获取网盘账号失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("获取网盘账号失败: %v", err),
		}
	}

	// 测试阶段，移除最小限制
	// 获取最小存储空间配置
	// autoTransferMinSpace, err := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	// if err != nil {
	// 	utils.Error("获取最小存储空间配置失败: %v", err)
	// 	autoTransferMinSpace = 5 // 默认5GB
	// }

	// // 过滤：只保留已激活、夸克平台、剩余空间足够的账号
	// minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024
	// var validAccounts []entity.Cks
	// for _, acc := range accounts {
	// 	if acc.IsValid && acc.PanID == *panID && acc.LeftSpace >= minSpaceBytes {
	// 		validAccounts = append(validAccounts, acc)
	// 	}
	// }

	// if len(validAccounts) == 0 {
	// 	utils.Info("没有可用的网盘账号")
	// 	return TransferResult{
	// 		Success:  false,
	// 		ErrorMsg: "没有可用的网盘账号",
	// 	}
	// }

	// utils.Info("找到 %d 个可用网盘账号，开始转存处理...", len(validAccounts))

	// 使用第一个可用账号进行转存
	// account := validAccounts[0]
	account := accounts[0]

	// 创建网盘服务工厂
	factory := pan.NewPanFactory()

	// 执行转存
	result := transferSingleResource(resource, account, factory)

	if result.Success {
		// 更新资源的转存信息
		resource.SaveURL = result.SaveURL
		resource.Fid = result.Fid
		resource.CkID = &account.ID
		resource.ErrorMsg = ""
		if err := repoManager.ResourceRepository.Update(resource); err != nil {
			utils.Error("更新资源转存信息失败: %v", err)
		}
	} else {
		// 更新错误信息
		resource.ErrorMsg = result.ErrorMsg
		if err := repoManager.ResourceRepository.Update(resource); err != nil {
			utils.Error("更新资源错误信息失败: %v", err)
		}
	}

	return result
}

// transferSingleResource 转存单个资源
func transferSingleResource(resource *entity.Resource, account entity.Cks, factory *pan.PanFactory) TransferResult {
	utils.Info("开始转存资源 - 资源ID: %d, 账号: %s", resource.ID, account.Username)

	service, err := factory.CreatePanService(resource.URL, &pan.PanConfig{
		URL:         resource.URL,
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	if err != nil {
		utils.Error("创建网盘服务失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("创建网盘服务失败: %v", err),
		}
	}

	// 设置账号信息
	service.SetCKSRepository(repoManager.CksRepository, account)

	// 提取分享ID
	shareID, _ := commonutils.ExtractShareIdString(resource.URL)
	if shareID == "" {
		return TransferResult{
			Success:  false,
			ErrorMsg: "无效的分享链接",
		}
	}

	// 执行转存
	transferResult, err := service.Transfer(shareID) // 有些链接还需要其他信息从 url 中自行解析
	if err != nil {
		utils.Error("转存失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("转存失败: %v", err),
		}
	}

	if transferResult == nil || !transferResult.Success {
		errMsg := "转存失败"
		if transferResult != nil && transferResult.Message != "" {
			errMsg = transferResult.Message
		}
		utils.Error("转存失败: %s", errMsg)
		return TransferResult{
			Success:  false,
			ErrorMsg: errMsg,
		}
	}

	// 提取转存链接
	var saveURL string
	var fid string

	if data, ok := transferResult.Data.(map[string]interface{}); ok {
		if v, ok := data["shareUrl"]; ok {
			saveURL, _ = v.(string)
		}
		if v, ok := data["fid"]; ok {
			fid, _ = v.(string)
		}
	}
	if saveURL == "" {
		saveURL = transferResult.ShareURL
	}

	if saveURL == "" {
		return TransferResult{
			Success:  false,
			ErrorMsg: "转存成功但未获取到分享链接",
		}
	}

	utils.Info("转存成功 - 资源ID: %d, 转存链接: %s", resource.ID, saveURL)

	return TransferResult{
		Success: true,
		SaveURL: saveURL,
		Fid:     fid,
	}
}

// getQuarkPanID 获取夸克网盘ID
func getQuarkPanID() (uint, error) {
	// 通过FindAll方法查找所有平台，然后过滤出quark平台
	pans, err := repoManager.PanRepository.FindAll()
	if err != nil {
		return 0, fmt.Errorf("查询平台信息失败: %v", err)
	}

	for _, p := range pans {
		if p.Name == "quark" {
			return p.ID, nil
		}
	}

	return 0, fmt.Errorf("未找到quark平台")
}
