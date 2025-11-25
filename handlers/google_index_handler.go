package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/google" // 添加google包导入
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	goauth "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/searchconsole/v1"
)

// GoogleIndexHandler Google索引处理程序
type GoogleIndexHandler struct {
	repoMgr     *repo.RepositoryManager
	taskManager *task.TaskManager
}

// NewGoogleIndexHandler 创建Google索引处理程序
func NewGoogleIndexHandler(
	repoMgr *repo.RepositoryManager,
	taskManager *task.TaskManager,
) *GoogleIndexHandler {
	return &GoogleIndexHandler{
		repoMgr:     repoMgr,
		taskManager: taskManager,
	}
}

// GetAllConfig 获取所有Google索引配置（以分组形式返回）
func (h *GoogleIndexHandler) GetAllConfig(c *gin.Context) {
	// 获取通用配置
	enabledStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyEnabled)
	if err != nil {
		enabledStr = "false"
	}
	enabled := enabledStr == "true" || enabledStr == "1"

	siteURL, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.ConfigKeyWebsiteURL)
	if err != nil {
		siteURL = ""
	}

	siteName, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySiteName)
	if err != nil {
		siteName = ""
	}

	// 获取调度配置
	checkIntervalStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyCheckInterval)
	if err != nil {
		checkIntervalStr = "60"
	}
	checkInterval, _ := strconv.Atoi(checkIntervalStr)

	batchSizeStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyBatchSize)
	if err != nil {
		batchSizeStr = "100"
	}
	batchSize, _ := strconv.Atoi(batchSizeStr)

	concurrencyStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyConcurrency)
	if err != nil {
		concurrencyStr = "5"
	}
	concurrency, _ := strconv.Atoi(concurrencyStr)

	retryAttemptsStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyRetryAttempts)
	if err != nil {
		retryAttemptsStr = "3"
	}
	retryAttempts, _ := strconv.Atoi(retryAttemptsStr)

	retryDelayStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyRetryDelay)
	if err != nil {
		retryDelayStr = "2"
	}
	retryDelay, _ := strconv.Atoi(retryDelayStr)

	// 获取网站地图配置
	autoSitemapStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyAutoSitemap)
	if err != nil {
		autoSitemapStr = "false"
	}
	autoSitemap := autoSitemapStr == "true" || autoSitemapStr == "1"

	sitemapPath, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySitemapPath)
	if err != nil {
		sitemapPath = "/sitemap.xml"
	}

	sitemapSchedule, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySitemapSchedule)
	if err != nil {
		sitemapSchedule = "@daily"
	}

	// 获取认证配置
	credentialsFile, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyCredentialsFile)
	if err != nil {
		credentialsFile = ""
	}

	clientEmail, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyClientEmail)
	if err != nil {
		clientEmail = ""
	}

	clientID, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyClientID)
	if err != nil {
		clientID = ""
	}

	privateKey, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyPrivateKey)
	if err != nil {
		privateKey = ""
	}

	token, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyToken)
	if err != nil {
		token = ""
	}

	// 构建各组配置
	generalConfig := dto.GoogleIndexConfigGeneral{
		Enabled:  enabled,
		SiteURL:  siteURL,
		SiteName: siteName,
	}

	authConfig := dto.GoogleIndexConfigAuth{
		CredentialsFile: credentialsFile,
		ClientEmail:     clientEmail,
		ClientID:        clientID,
		PrivateKey:      privateKey,
		Token:           token,
	}

	scheduleConfig := dto.GoogleIndexConfigSchedule{
		CheckInterval: checkInterval,
		BatchSize:     batchSize,
		Concurrency:   concurrency,
		RetryAttempts: retryAttempts,
		RetryDelay:    retryDelay,
	}

	sitemapConfig := dto.GoogleIndexConfigSitemap{
		AutoSitemap:     autoSitemap,
		SitemapPath:     sitemapPath,
		SitemapSchedule: sitemapSchedule,
	}

	// 将配置对象转换为JSON字符串
	generalConfigJSON, _ := json.Marshal(generalConfig)
	authConfigJSON, _ := json.Marshal(authConfig)
	scheduleConfigJSON, _ := json.Marshal(scheduleConfig)
	sitemapConfigJSON, _ := json.Marshal(sitemapConfig)

	// 以数组格式返回所有配置组
	configs := []gin.H{
		{
			"group": "general",
			"key":   "general",
			"value": string(generalConfigJSON),
		},
		{
			"group": "auth",
			"key":   "credentials_file", // 使用具体的键名
			"value": string(authConfigJSON),
		},
		{
			"group": "schedule",
			"key":   "schedule",
			"value": string(scheduleConfigJSON),
		},
		{
			"group": "sitemap",
			"key":   "sitemap",
			"value": string(sitemapConfigJSON),
		},
	}

	SuccessResponse(c, configs)
}

// GetConfig 获取Google索引配置（原有接口，为保持兼容性）
func (h *GoogleIndexHandler) GetConfig(c *gin.Context) {
	// 获取通用配置
	enabledStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyEnabled)
	if err != nil {
		enabledStr = "false"
	}
	enabled := enabledStr == "true" || enabledStr == "1"

	siteURL, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.ConfigKeyWebsiteURL)
	if err != nil {
		siteURL = ""
	}

	siteName, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySiteName)
	if err != nil {
		siteName = ""
	}

	// 获取调度配置
	checkIntervalStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyCheckInterval)
	if err != nil {
		checkIntervalStr = "60"
	}
	checkInterval, _ := strconv.Atoi(checkIntervalStr)

	batchSizeStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyBatchSize)
	if err != nil {
		batchSizeStr = "100"
	}
	batchSize, _ := strconv.Atoi(batchSizeStr)

	concurrencyStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyConcurrency)
	if err != nil {
		concurrencyStr = "5"
	}
	concurrency, _ := strconv.Atoi(concurrencyStr)

	retryAttemptsStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyRetryAttempts)
	if err != nil {
		retryAttemptsStr = "3"
	}
	retryAttempts, _ := strconv.Atoi(retryAttemptsStr)

	retryDelayStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyRetryDelay)
	if err != nil {
		retryDelayStr = "2"
	}
	retryDelay, _ := strconv.Atoi(retryDelayStr)

	// 获取网站地图配置
	autoSitemapStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyAutoSitemap)
	if err != nil {
		autoSitemapStr = "false"
	}
	autoSitemap := autoSitemapStr == "true" || autoSitemapStr == "1"

	sitemapPath, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySitemapPath)
	if err != nil {
		sitemapPath = "/sitemap.xml"
	}

	sitemapSchedule, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySitemapSchedule)
	if err != nil {
		sitemapSchedule = "@daily"
	}

	config := dto.GoogleIndexConfigGeneral{
		Enabled:  enabled,
		SiteURL:  siteURL,
		SiteName: siteName,
	}

	scheduleConfig := dto.GoogleIndexConfigSchedule{
		CheckInterval: checkInterval,
		BatchSize:     batchSize,
		Concurrency:   concurrency,
		RetryAttempts: retryAttempts,
		RetryDelay:    retryDelay,
	}

	sitemapConfig := dto.GoogleIndexConfigSitemap{
		AutoSitemap:     autoSitemap,
		SitemapPath:     sitemapPath,
		SitemapSchedule: sitemapSchedule,
	}

	result := gin.H{
		"general":  config,
		"schedule": scheduleConfig,
		"sitemap":  sitemapConfig,
		"is_running": false, // 不再有独立的调度器，使用统一任务管理器
	}

	SuccessResponse(c, result)
}

// UpdateConfig 更新Google索引配置
func (h *GoogleIndexHandler) UpdateConfig(c *gin.Context) {
	var req struct {
		General  dto.GoogleIndexConfigGeneral  `json:"general"`
		Schedule dto.GoogleIndexConfigSchedule `json:"schedule"`
		Sitemap  dto.GoogleIndexConfigSitemap  `json:"sitemap"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GoogleIndexHandler.UpdateConfig - 用户更新Google索引配置 - 用户: %s, IP: %s", username, clientIP)

	// 准备要更新的配置项
	configs := []entity.SystemConfig{
		{
			Key:   entity.GoogleIndexConfigKeyEnabled,
			Value: strconv.FormatBool(req.General.Enabled),
			Type:  entity.ConfigTypeBool,
		},
				{
			Key:   entity.GoogleIndexConfigKeySiteName,
			Value: req.General.SiteName,
			Type:  entity.ConfigTypeString,
		},
		{
			Key:   entity.GoogleIndexConfigKeyCheckInterval,
			Value: strconv.Itoa(req.Schedule.CheckInterval),
			Type:  entity.ConfigTypeInt,
		},
		{
			Key:   entity.GoogleIndexConfigKeyBatchSize,
			Value: strconv.Itoa(req.Schedule.BatchSize),
			Type:  entity.ConfigTypeInt,
		},
		{
			Key:   entity.GoogleIndexConfigKeyConcurrency,
			Value: strconv.Itoa(req.Schedule.Concurrency),
			Type:  entity.ConfigTypeInt,
		},
		{
			Key:   entity.GoogleIndexConfigKeyRetryAttempts,
			Value: strconv.Itoa(req.Schedule.RetryAttempts),
			Type:  entity.ConfigTypeInt,
		},
		{
			Key:   entity.GoogleIndexConfigKeyRetryDelay,
			Value: strconv.Itoa(req.Schedule.RetryDelay),
			Type:  entity.ConfigTypeInt,
		},
		{
			Key:   entity.GoogleIndexConfigKeyAutoSitemap,
			Value: strconv.FormatBool(req.Sitemap.AutoSitemap),
			Type:  entity.ConfigTypeBool,
		},
		{
			Key:   entity.GoogleIndexConfigKeySitemapPath,
			Value: req.Sitemap.SitemapPath,
			Type:  entity.ConfigTypeString,
		},
		{
			Key:   entity.GoogleIndexConfigKeySitemapSchedule,
			Value: req.Sitemap.SitemapSchedule,
			Type:  entity.ConfigTypeString,
		},
	}

	// 批量更新配置
	err := h.repoMgr.SystemConfigRepository.UpsertConfigs(configs)
	if err != nil {
		utils.Error("更新系统配置失败: %v", err)
		ErrorResponse(c, "更新配置失败", http.StatusInternalServerError)
		return
	}

	utils.Info("Google索引配置更新成功 - 用户: %s, IP: %s", username, clientIP)
	SuccessResponse(c, gin.H{
		"message": "配置更新成功",
	})
}

// CreateTask 创建Google索引任务
func (h *GoogleIndexHandler) CreateTask(c *gin.Context) {
	var req dto.GoogleIndexTaskInput
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GoogleIndexHandler.CreateTask - 用户创建Google索引任务 - 用户: %s, 任务类型: %s, 任务标题: %s, IP: %s", username, req.Type, req.Title, clientIP)

	// 创建通用任务
	task, err := h.taskManager.CreateTask(string(entity.TaskTypeGoogleIndex), req.Title, req.Description, req.ConfigID)
	if err != nil {
		utils.Error("创建Google索引任务失败: %v", err)
		ErrorResponse(c, "创建任务失败", http.StatusInternalServerError)
		return
	}

	// 根据任务类型创建任务项
	var taskItems []*entity.TaskItem

	switch req.Type {
	case "url_indexing", "status_check", "batch_index", "manual_check", "url_submit":
		// 为每个URL创建任务项
		for _, url := range req.URLs {
			itemData := map[string]interface{}{
				"urls":      []string{url},
				"operation": req.Type,
			}
			itemDataJSON, _ := json.Marshal(itemData)

			taskItem := &entity.TaskItem{
				URL:       url,
				InputData: string(itemDataJSON),
			}
			taskItems = append(taskItems, taskItem)
		}

	case "sitemap_submit":
		// 为网站地图创建任务项
		itemData := map[string]interface{}{
			"sitemap_url": req.SitemapURL,
			"operation":   "sitemap_submit",
		}
		itemDataJSON, _ := json.Marshal(itemData)

		taskItem := &entity.TaskItem{
			URL:       req.SitemapURL,
			InputData: string(itemDataJSON),
		}
		taskItems = append(taskItems, taskItem)
	}

	// 批量创建任务项
	err = h.taskManager.CreateTaskItems(task.ID, taskItems)
	if err != nil {
		utils.Error("创建任务项失败: %v", err)
		ErrorResponse(c, "创建任务项失败", http.StatusInternalServerError)
		return
	}

	// 更新任务的总项目数
	err = h.repoMgr.TaskRepository.UpdateTotalItems(task.ID, len(taskItems))
	if err != nil {
		utils.Error("更新任务总项目数失败: %v", err)
		// 不返回错误，因为任务和任务项已经创建成功
	}

	utils.Info("Google索引任务创建完成: %d, 任务类型: %s, 总项目数: %d", task.ID, req.Type, len(taskItems))
	SuccessResponse(c, gin.H{
		"task_id":     task.ID,
		"total_items": len(taskItems),
		"message":     "任务创建成功",
	})
}

// StartTask 启动Google索引任务
func (h *GoogleIndexHandler) StartTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID", http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GoogleIndexHandler.StartTask - 用户启动Google索引任务 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)

	// 使用任务管理器启动任务
	err = h.taskManager.StartTask(uint(taskID))
	if err != nil {
		utils.Error("启动Google索引任务失败: %v", err)
		ErrorResponse(c, "启动任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("Google索引任务启动成功: %d", taskID)
	SuccessResponse(c, gin.H{
		"message": "任务启动成功",
	})
}


// GetTaskStatus 获取任务状态
func (h *GoogleIndexHandler) GetTaskStatus(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID", http.StatusBadRequest)
		return
	}

	task, err := h.taskManager.GetTask(uint(taskID))
	if err != nil {
		ErrorResponse(c, "获取任务失败", http.StatusInternalServerError)
		return
	}

	if task == nil {
		ErrorResponse(c, "任务不存在", http.StatusNotFound)
		return
	}

	// 获取任务项统计
	stats, err := h.taskManager.GetTaskItemStats(task.ID)
	if err != nil {
		utils.Error("获取任务项统计失败: %v", err)
		stats = make(map[string]int)
	}

	taskOutput := converter.TaskToGoogleIndexTaskOutput(task, stats)

	result := gin.H{
		"id":                  taskOutput.ID,
		"name":                taskOutput.Name,
		"type":                taskOutput.Type,
		"status":              taskOutput.Status,
		"description":         taskOutput.Description,
		"progress":            taskOutput.Progress,
		"total_items":         taskOutput.TotalItems,
		"processed_items":     taskOutput.ProcessedItems,
		"successful_items":    taskOutput.SuccessfulItems,
		"failed_items":        taskOutput.FailedItems,
		"pending_items":       taskOutput.PendingItems,
		"processing_items":    taskOutput.ProcessingItems,
		"indexed_urls":        taskOutput.IndexedURLs,
		"failed_urls":         taskOutput.FailedURLs,
		"started_at":          taskOutput.StartedAt,
		"completed_at":        taskOutput.CompletedAt,
		"created_at":          taskOutput.CreatedAt,
		"updated_at":          taskOutput.UpdatedAt,
		"progress_data":       taskOutput.ProgressData,
		"stats":               stats,
	}

	SuccessResponse(c, result)
}

// GetTasks 获取任务列表
func (h *GoogleIndexHandler) GetTasks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	taskTypeStr := c.Query("type")
	statusStr := c.Query("status")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 根据参数筛选任务类型，如果有指定则使用，否则默认为Google索引类型
	taskType := string(entity.TaskTypeGoogleIndex)
	if taskTypeStr != "" {
		taskType = taskTypeStr
	}

	// 获取指定状态的任务，默认查找所有状态
	status := statusStr

	// 获取任务列表 - 目前我们没有Query方法，直接获取所有任务然后做筛选
	tasksList, total, err := h.repoMgr.TaskRepository.GetList(page, pageSize, taskType, status)
	if err != nil {
		ErrorResponse(c, "获取任务列表失败", http.StatusInternalServerError)
		return
	}

	taskOutputs := make([]dto.GoogleIndexTaskOutput, len(tasksList))
	for i, task := range tasksList {
		// 获取任务统计信息
		stats, err := h.taskManager.GetTaskItemStats(task.ID)
		if err != nil {
			stats = make(map[string]int)
		}
		taskOutputs[i] = converter.TaskToGoogleIndexTaskOutput(task, stats)
	}

	result := dto.GoogleIndexTaskListResponse{
		Tasks:      taskOutputs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	SuccessResponse(c, result)
}

// GetTaskItems 获取任务项列表
func (h *GoogleIndexHandler) GetTaskItems(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID", http.StatusBadRequest)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")
	statusStr := c.Query("status")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 50
	}

	// 获取任务项列表
	items, total, err := h.taskManager.QueryTaskItems(uint(taskID), page, pageSize, statusStr)
	if err != nil {
		ErrorResponse(c, "获取任务项列表失败", http.StatusInternalServerError)
		return
	}

	// 注意：我们还没有TaskItemToGoogleIndexTaskItemOutput转换器，需要创建一个
	itemOutputs := make([]dto.GoogleIndexTaskItemOutput, len(items))
	for i, item := range items {
		// 手动构建输出结构
		itemOutputs[i] = dto.GoogleIndexTaskItemOutput{
			ID:           item.ID,
			TaskID:       item.TaskID,
			URL:          item.URL,
			Status:       string(item.Status),
			IndexStatus:  item.IndexStatus,
			ErrorMessage: item.ErrorMessage,
			InspectResult: item.InspectResult,
			MobileFriendly: item.MobileFriendly,
			LastCrawled:  item.LastCrawled,
			StatusCode:   item.StatusCode,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
			StartedAt:    item.ProcessedAt, // 任务项处理完成时间
			CompletedAt:  item.ProcessedAt,
		}
	}

	result := dto.GoogleIndexTaskItemPageResponse{
		Items: itemOutputs,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}

	SuccessResponse(c, result)
}

// UploadCredentials 上传Google索引凭据
func (h *GoogleIndexHandler) UploadCredentials(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		ErrorResponse(c, "未提供凭据文件", http.StatusBadRequest)
		return
	}

	// 验证文件扩展名必须是.json
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".json" {
		ErrorResponse(c, "仅支持JSON格式的凭据文件", http.StatusBadRequest)
		return
	}

	// 验证文件大小（限制5MB）
	if file.Size > 5*1024*1024 {
		ErrorResponse(c, "文件大小不能超过5MB", http.StatusBadRequest)
		return
	}

	// 确保data目录存在
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		ErrorResponse(c, "创建数据目录失败", http.StatusInternalServerError)
		return
	}

	// 使用固定的文件名保存凭据
	fixedFileName := "google_credentials.json"
	filePath := filepath.Join(dataDir, fixedFileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		ErrorResponse(c, "保存凭据文件失败", http.StatusInternalServerError)
		return
	}

	// 设置文件权限
	if err := os.Chmod(filePath, 0600); err != nil {
		utils.Warn("设置凭据文件权限失败: %v", err)
	}

	// 返回成功响应
	accessPath := filepath.Join("data", fixedFileName)
	response := map[string]interface{}{
		"success":    true,
		"message":    "凭据文件上传成功",
		"file_name":  fixedFileName,
		"file_path":  accessPath,
		"full_path":  filePath,
	}

	SuccessResponse(c, response)
}

// makeSafeFileName 生成安全的文件名，移除危险字符
func (h *GoogleIndexHandler) makeSafeFileName(filename string) string {
	// 移除路径分隔符和特殊字符
	safeName := strings.ReplaceAll(filename, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, "..", "_")

	// 限制文件名长度
	if len(safeName) > 100 {
		ext := filepath.Ext(safeName)
		name := safeName[:100-len(ext)]
		safeName = name + ext
	}

	return safeName
}

// ValidateCredentials 验证Google索引凭据
func (h *GoogleIndexHandler) ValidateCredentials(c *gin.Context) {
	// 使用固定的凭据文件路径
	credentialsFile := "data/google_credentials.json"

	// 检查凭据文件是否存在
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		ErrorResponse(c, "凭据文件不存在", http.StatusBadRequest)
		return
	}

	// 尝试创建Google客户端并验证凭据
	config, err := h.loadCredentials(credentialsFile)
	if err != nil {
		ErrorResponse(c, "凭据格式错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证凭据是否有效（尝试获取token）
	err = h.getValidToken(config)
	if err != nil {
		ErrorResponse(c, "凭据验证失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "凭据验证成功",
		"valid":   true,
	}

	SuccessResponse(c, response)
}

// GetStatus 获取Google索引状态
func (h *GoogleIndexHandler) GetStatus(c *gin.Context) {
	// 获取通用配置
	enabledStr, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyEnabled)
	if err != nil {
		enabledStr = "false"
	}
	enabled := enabledStr == "true" || enabledStr == "1"

	siteURL, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.ConfigKeyWebsiteURL)
	if err != nil {
		siteURL = ""
	}

	// 获取认证配置
	credentialsFile, err := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyCredentialsFile)
	if err != nil {
		credentialsFile = ""
	}

	// 验证凭据是否有效
	authValid := false
	if credentialsFile != "" {
		if _, err := os.Stat(credentialsFile); !os.IsNotExist(err) {
			// 检查凭据文件是否存在且可读
			authValid = true
		}
	}

	// 获取统计信息（从数据库查询实际的资源数量）
	totalURLs := 0
	indexedURLs := 0
	notIndexedURLs := 0
	errorURLs := 0

	// 查询resources表获取总URL数
	totalResources, err := h.repoMgr.ResourceRepository.GetTotalCount()
	if err == nil {
		totalURLs = int(totalResources)
	}

	// 查询任务项统计获取索引状态
	taskStats, err := h.repoMgr.TaskItemRepository.GetIndexStats()
	if err == nil {
		indexedURLs = taskStats["indexed"]
		notIndexedURLs = taskStats["not_indexed"]
		errorURLs = taskStats["error"]
	}

	statusResponse := dto.GoogleIndexStatusResponse{
		Enabled:           enabled,
		SiteURL:           siteURL,
		LastCheckTime:     time.Now(),
		TotalURLs:         totalURLs,
		IndexedURLs:       indexedURLs,
		NotIndexedURLs:    notIndexedURLs,
		ErrorURLs:         errorURLs,
		LastSitemapSubmit: time.Time{},
		AuthValid:         authValid,
	}

	SuccessResponse(c, statusResponse)
}

// loadCredentials 从文件加载凭据
func (h *GoogleIndexHandler) loadCredentials(credentialsFile string) (*google.Config, error) {
	// 读取凭据文件
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("无法读取凭据文件: %v", err)
	}

	// 验证是否为有效的JSON
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, fmt.Errorf("凭据文件格式不是有效的JSON: %v", err)
	}

	// 检查必需字段
	requiredFields := []string{"type", "project_id", "private_key_id", "private_key", "client_email", "client_id"}
	for _, field := range requiredFields {
		if _, exists := temp[field]; !exists {
			return nil, fmt.Errorf("凭据文件缺少必需字段: %s", field)
		}
	}

	// 检查凭据类型
	if temp["type"] != "service_account" {
		return nil, fmt.Errorf("仅支持服务账号类型的凭据")
	}

	// 尝试从JSON数据加载凭据，需要指定作用域
	scopes := []string{
		"https://www.googleapis.com/auth/webmasters",
		"https://www.googleapis.com/auth/indexing",
	}
	jwtConfig, err := goauth.JWTConfigFromJSON(data, scopes...)
	if err != nil {
		return nil, fmt.Errorf("创建JWT配置失败: %v", err)
	}

	// 创建一个简单的配置对象，暂时只存储JWT配置
	config := &google.Config{
		CredentialsFile: credentialsFile,
	}

	// 为了验证凭据，我们尝试获取token源
	ctx := context.Background()
	tokenSource := jwtConfig.TokenSource(ctx)
	_ = tokenSource // 实际验证在getValidToken中进行

	return config, nil
}

// getValidToken 获取有效的token


// UpdateGoogleIndexConfig 更新Google索引配置（支持分组配置）
func (h *GoogleIndexHandler) UpdateGoogleIndexConfig(c *gin.Context) {
	var req dto.GoogleIndexConfigInput
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GoogleIndexHandler.UpdateGoogleIndexConfig - 用户更新Google索引分组配置 - 用户: %s, 组: %s, 键: %s, IP: %s", username, req.Group, req.Key, clientIP)

	// 处理不同的配置组

	switch req.Group {
	case "general":
		switch req.Key {
		case "general":
			// 解析general配置
			var generalConfig dto.GoogleIndexConfigGeneral
			if err := json.Unmarshal([]byte(req.Value), &generalConfig); err != nil {
				ErrorResponse(c, "通用配置格式错误: "+err.Error(), http.StatusBadRequest)
				return
			}
			// 存储各个字段
			generalConfigs := []entity.SystemConfig{
				{Key: entity.GoogleIndexConfigKeyEnabled, Value: strconv.FormatBool(generalConfig.Enabled), Type: entity.ConfigTypeBool},
								{Key: entity.GoogleIndexConfigKeySiteName, Value: generalConfig.SiteName, Type: entity.ConfigTypeString},
			}
			err := h.repoMgr.SystemConfigRepository.UpsertConfigs(generalConfigs)
			if err != nil {
				ErrorResponse(c, "保存通用配置失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			ErrorResponse(c, "未知的通用配置键: "+req.Key, http.StatusBadRequest)
			return
		}
	case "auth":
		switch req.Key {
		case "credentials_file":
			// 解析认证配置
			var authConfig dto.GoogleIndexConfigAuth
			if err := json.Unmarshal([]byte(req.Value), &authConfig); err != nil {
				ErrorResponse(c, "认证配置格式错误: "+err.Error(), http.StatusBadRequest)
				return
			}
			// 存储认证相关配置
			authConfigs := []entity.SystemConfig{
				{Key: entity.GoogleIndexConfigKeyCredentialsFile, Value: authConfig.CredentialsFile, Type: entity.ConfigTypeString},
				{Key: entity.GoogleIndexConfigKeyClientEmail, Value: authConfig.ClientEmail, Type: entity.ConfigTypeString},
				{Key: entity.GoogleIndexConfigKeyClientID, Value: authConfig.ClientID, Type: entity.ConfigTypeString},
				{Key: entity.GoogleIndexConfigKeyPrivateKey, Value: authConfig.PrivateKey, Type: entity.ConfigTypeString},
			}
			err := h.repoMgr.SystemConfigRepository.UpsertConfigs(authConfigs)
			if err != nil {
				ErrorResponse(c, "保存认证配置失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			ErrorResponse(c, "未知的认证配置键: "+req.Key, http.StatusBadRequest)
			return
		}
	case "schedule":
		switch req.Key {
		case "schedule":
			// 解析调度配置
			var scheduleConfig dto.GoogleIndexConfigSchedule
			if err := json.Unmarshal([]byte(req.Value), &scheduleConfig); err != nil {
				ErrorResponse(c, "调度配置格式错误: "+err.Error(), http.StatusBadRequest)
				return
			}
			// 存储调度相关配置
			scheduleConfigs := []entity.SystemConfig{
				{Key: entity.GoogleIndexConfigKeyCheckInterval, Value: strconv.Itoa(scheduleConfig.CheckInterval), Type: entity.ConfigTypeInt},
				{Key: entity.GoogleIndexConfigKeyBatchSize, Value: strconv.Itoa(scheduleConfig.BatchSize), Type: entity.ConfigTypeInt},
				{Key: entity.GoogleIndexConfigKeyConcurrency, Value: strconv.Itoa(scheduleConfig.Concurrency), Type: entity.ConfigTypeInt},
				{Key: entity.GoogleIndexConfigKeyRetryAttempts, Value: strconv.Itoa(scheduleConfig.RetryAttempts), Type: entity.ConfigTypeInt},
				{Key: entity.GoogleIndexConfigKeyRetryDelay, Value: strconv.Itoa(scheduleConfig.RetryDelay), Type: entity.ConfigTypeInt},
			}
			err := h.repoMgr.SystemConfigRepository.UpsertConfigs(scheduleConfigs)
			if err != nil {
				ErrorResponse(c, "保存调度配置失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			ErrorResponse(c, "未知的调度配置键: "+req.Key, http.StatusBadRequest)
			return
		}
	case "sitemap":
		switch req.Key {
		case "sitemap":
			// 解析网站地图配置
			var sitemapConfig dto.GoogleIndexConfigSitemap
			if err := json.Unmarshal([]byte(req.Value), &sitemapConfig); err != nil {
				ErrorResponse(c, "网站地图配置格式错误: "+err.Error(), http.StatusBadRequest)
				return
			}
			// 存储网站地图相关配置
			sitemapConfigs := []entity.SystemConfig{
				{Key: entity.GoogleIndexConfigKeyAutoSitemap, Value: strconv.FormatBool(sitemapConfig.AutoSitemap), Type: entity.ConfigTypeBool},
				{Key: entity.GoogleIndexConfigKeySitemapPath, Value: sitemapConfig.SitemapPath, Type: entity.ConfigTypeString},
				{Key: entity.GoogleIndexConfigKeySitemapSchedule, Value: sitemapConfig.SitemapSchedule, Type: entity.ConfigTypeString},
			}
			err := h.repoMgr.SystemConfigRepository.UpsertConfigs(sitemapConfigs)
			if err != nil {
				ErrorResponse(c, "保存网站地图配置失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			ErrorResponse(c, "未知的网站地图配置键: "+req.Key, http.StatusBadRequest)
			return
		}
		default:
		ErrorResponse(c, "未知的配置组: "+req.Group, http.StatusBadRequest)
		return
	}

	utils.Info("Google索引分组配置更新成功 - 组: %s, 键: %s - 用户: %s, IP: %s", req.Group, req.Key, username, clientIP)
	SuccessResponse(c, gin.H{
		"message": "配置更新成功",
		"group":   req.Group,
		"key":     req.Key,
	})
}

func (h *GoogleIndexHandler) getValidToken(config *google.Config) error {
	// 重新读取凭据文件进行验证
	data, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("无法读取凭据文件: %v", err)
	}

	// 尝试从JSON数据加载凭据，需要指定作用域
	scopes := []string{
		"https://www.googleapis.com/auth/webmasters",
		"https://www.googleapis.com/auth/indexing",
	}

	// 检查凭据类型
	var credentialsMap map[string]interface{}
	if err := json.Unmarshal(data, &credentialsMap); err != nil {
		return fmt.Errorf("解析凭据失败: %v", err)
	}

	credType, ok := credentialsMap["type"].(string)
	if !ok {
		return fmt.Errorf("未知的凭据类型")
	}

	ctx := context.Background()
	var client *http.Client

	if credType == "service_account" {
		// 服务账号凭据
		jwtConfig, err := goauth.JWTConfigFromJSON(data, scopes...)
		if err != nil {
			return fmt.Errorf("创建JWT配置失败: %v", err)
		}
		client = jwtConfig.Client(ctx)

		// 尝试获取一个测试URL来验证凭据
		siteURL, _ := h.repoMgr.SystemConfigRepository.GetConfigValue(entity.ConfigKeyWebsiteURL)
		if siteURL == "" {
			siteURL = "https://example.com" // 使用默认URL进行测试
		}

		// 创建Search Console服务进行测试
		searchService, err := searchconsole.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return fmt.Errorf("创建Search Console服务失败: %v", err)
		}

		// 尝试获取站点列表来验证凭据
		_, err = searchService.Sites.List().Do()
		if err != nil {
			return fmt.Errorf("凭据验证失败 - 无法访问Google Search Console API: %v", err)
		}

		utils.Info("Google服务账号凭据验证成功")

	} else {
		// OAuth2客户端凭据
		oauthConfig, err := goauth.ConfigFromJSON(data, scopes...)
		if err != nil {
			return fmt.Errorf("创建OAuth配置失败: %v", err)
		}

		// 尝试从文件读取token
		tokenFile := "data/google_token.json"
		token, err := h.tokenFromFile(tokenFile)
		if err != nil {
			return fmt.Errorf("未找到有效的token文件，请先完成OAuth认证流程: %v", err)
		}

		// 验证token是否过期
		if !token.Valid() {
			return fmt.Errorf("Token已过期，请重新认证")
		}

		client = oauthConfig.Client(ctx, token)

		// 测试API访问
		searchService, err := searchconsole.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return fmt.Errorf("创建Search Console服务失败: %v", err)
		}

		_, err = searchService.Sites.List().Do()
		if err != nil {
			return fmt.Errorf("凭据验证失败 - 无法访问Google Search Console API: %v", err)
		}

		utils.Info("Google OAuth2凭据验证成功")
	}

	return nil
}

// tokenFromFile 从文件读取token
func (h *GoogleIndexHandler) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// SubmitURLsToIndex 提交URL到Google索引
func (h *GoogleIndexHandler) SubmitURLsToIndex(c *gin.Context) {
	var req struct {
		URLs []string `json:"urls" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GoogleIndexHandler.SubmitURLsToIndex - 用户提交URL到索引 - 用户: %s, URL数量: %d, IP: %s", username, len(req.URLs), clientIP)

	// 创建通用任务
	var configID *uint
	task, err := h.taskManager.CreateTask(string(entity.TaskTypeGoogleIndex), fmt.Sprintf("手动URL提交任务 - %d个URL", len(req.URLs)), fmt.Sprintf("手动提交 %d 个URL到Google索引", len(req.URLs)), configID)
	if err != nil {
		utils.Error("创建Google索引任务失败: %v", err)
		ErrorResponse(c, "创建任务失败", http.StatusInternalServerError)
		return
	}

	// 为每个URL创建任务项
	var taskItems []*entity.TaskItem
	for _, url := range req.URLs {
		itemData := map[string]interface{}{
			"urls":      []string{url},
			"operation": "url_submit",
		}
		itemDataJSON, _ := json.Marshal(itemData)

		taskItem := &entity.TaskItem{
			URL:       url,
			InputData: string(itemDataJSON),
		}
		taskItems = append(taskItems, taskItem)
	}

	// 批量创建任务项
	err = h.taskManager.CreateTaskItems(task.ID, taskItems)
	if err != nil {
		utils.Error("创建任务项失败: %v", err)
		ErrorResponse(c, "创建任务项失败", http.StatusInternalServerError)
		return
	}

	// 更新任务的总项目数
	err = h.repoMgr.TaskRepository.UpdateTotalItems(task.ID, len(taskItems))
	if err != nil {
		utils.Error("更新任务总项目数失败: %v", err)
	}

	// 自动启动任务
	err = h.taskManager.StartTask(task.ID)
	if err != nil {
		utils.Error("启动Google索引任务失败: %v", err)
		ErrorResponse(c, "启动任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("Google索引URL提交任务创建完成: %d, 总项目数: %d", task.ID, len(taskItems))
	SuccessResponse(c, gin.H{
		"task_id":     task.ID,
		"total_items": len(taskItems),
		"message":     "URL提交任务已创建并启动",
	})
}

// DiagnosePermissions 诊断Google API权限
func (h *GoogleIndexHandler) DiagnosePermissions(c *gin.Context) {
	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("GoogleIndexHandler.DiagnosePermissions - 用户执行权限诊断 - 用户: %s, IP: %s", username, clientIP)

	// 凭据文件路径
	credentialsFile := "data/google_credentials.json"

	// 检查凭据文件是否存在
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		ErrorResponse(c, "凭据文件不存在: "+credentialsFile, http.StatusBadRequest)
		return
	}

	// 读取凭据文件
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		ErrorResponse(c, "读取凭据文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 解析凭据
	var creds struct {
		Type        string `json:"type"`
		ProjectID   string `json:"project_id"`
		ClientEmail string `json:"client_email"`
	}
	if err := json.Unmarshal(data, &creds); err != nil {
		ErrorResponse(c, "解析凭据文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建Google客户端
	ctx := context.Background()
	config, err := goauth.JWTConfigFromJSON(data, searchconsole.WebmastersScope)
	if err != nil {
		ErrorResponse(c, "创建JWT配置失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := config.Client(ctx)

	// 创建Search Console服务
	service, err := searchconsole.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		ErrorResponse(c, "创建Search Console服务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 诊断结果结构
	diagnosis := map[string]interface{}{
		"credentials": map[string]interface{}{
			"file_exists":    true,
			"service_account": creds.ClientEmail,
			"project_id":     creds.ProjectID,
			"type":          creds.Type,
		},
		"api_access": map[string]interface{}{
			"search_console_enabled": true,
			"sites_count":           0,
			"sites":                []interface{}{},
		},
		"site_tests": []interface{}{},
		"recommendations": []string{},
	}

	// 测试基本API访问 - 获取站点列表
	sites, err := service.Sites.List().Do()
	if err != nil {
		diagnosis["api_access"].(map[string]interface{})["sites_error"] = err.Error()
		diagnosis["recommendations"] = append(diagnosis["recommendations"].([]string),
			"无法访问Search Console API，请检查服务账号权限")
	} else {
		diagnosis["api_access"].(map[string]interface{})["sites_count"] = len(sites.SiteEntry)

		// 添加站点列表
		for i, site := range sites.SiteEntry {
			if i < 5 { // 只显示前5个站点
				diagnosis["api_access"].(map[string]interface{})["sites"] = append(
					diagnosis["api_access"].(map[string]interface{})["sites"].([]interface{}),
					map[string]interface{}{
						"url":             site.SiteUrl,
						"permission_level": site.PermissionLevel,
					})
			}
		}
	}

	// 测试特定站点访问
	targetSite := "https://pan.l9.lc"
	siteFormats := []string{
		targetSite,
		targetSite + "/",
		"sc-domain:pan.l9.lc",
	}

	for _, siteURL := range siteFormats {
		siteTest := map[string]interface{}{
			"site_format": siteURL,
			"site_access": false,
			"url_inspect": false,
		}

		// 测试站点访问
		_, err = service.Sites.Get(siteURL).Do()
		if err != nil {
			siteTest["site_error"] = err.Error()
		} else {
			siteTest["site_access"] = true
		}

		// 测试URL检查
		testURL := targetSite + "/test"
		inspectionRequest := &searchconsole.InspectUrlIndexRequest{
			InspectionUrl: testURL,
			SiteUrl:       siteURL,
			LanguageCode:  "zh-CN",
		}

		_, err = service.UrlInspection.Index.Inspect(inspectionRequest).Do()
		if err != nil {
			siteTest["inspect_error"] = err.Error()
		} else {
			siteTest["url_inspect"] = true
		}

		diagnosis["site_tests"] = append(diagnosis["site_tests"].([]interface{}), siteTest)
	}

	// 生成建议
	sitesCount := diagnosis["api_access"].(map[string]interface{})["sites_count"].(int)
	if sitesCount == 0 {
		diagnosis["recommendations"] = append(diagnosis["recommendations"].([]string),
			"服务账号未被授权访问任何Search Console站点",
			"请在Google Search Console中添加服务账号为用户")
	}

	// 检查是否有任何站点测试成功
	hasSuccessfulSiteTest := false
	for _, test := range diagnosis["site_tests"].([]interface{}) {
		if testMap, ok := test.(map[string]interface{}); ok {
			if siteAccess, exists := testMap["site_access"].(bool); exists && siteAccess {
				hasSuccessfulSiteTest = true
				break
			}
		}
	}

	if !hasSuccessfulSiteTest {
		diagnosis["recommendations"] = append(diagnosis["recommendations"].([]string),
			"所有站点访问测试失败，请检查站点所有权验证和服务账号权限")
	}

	// 添加具体操作步骤
	diagnosis["recommendations"] = append(diagnosis["recommendations"].([]string),
		"1. 登录Google Search Console: https://search.google.com/search-console",
		"2. 选择站点 https://pan.l9.lc",
		"3. 进入 设置 → 用户和权限",
		"4. 添加用户: "+creds.ClientEmail,
		"5. 授予 '所有者' 或 '完整' 权限",
		"6. 等待权限生效（可能需要几分钟）",
		"7. 确保Indexing API已启用: https://console.cloud.google.com/apis/library/indexing.googleapis.com")

	utils.Info("Google权限诊断完成，站点数量: %d", sitesCount)
	SuccessResponse(c, gin.H{
		"diagnosis": diagnosis,
		"message":   "权限诊断完成",
	})
}
