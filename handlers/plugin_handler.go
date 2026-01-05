package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin-system/manager/plugin"
	"github.com/gin-gonic/gin"
)

// PluginHandler 插件管理处理器
type PluginHandler struct {
	repoManager   *repo.RepositoryManager
	metadataParser *plugin.MetadataParser
	pluginManager *plugin.Manager
}

// NewPluginHandler 创建插件处理器
func NewPluginHandler(repoManager *repo.RepositoryManager, pluginManager *plugin.Manager) *PluginHandler {
	return &PluginHandler{
		repoManager:   repoManager,
		metadataParser: plugin.NewMetadataParser(),
		pluginManager: pluginManager,
	}
}

// getPluginConfigName 根据插件名称获取配置名称
// 新的命名系统：使用插件元数据中的@name字段，而不是文件名+.plugin
func (h *PluginHandler) getPluginConfigName(pluginName string) string {
	// 尝试从插件元数据获取真实的插件名称
	if metadata, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".plugin.js")); err == nil {
		return metadata.Name // 使用@name字段的值
	}

	// 如果解析失败，返回原始名称（兼容旧逻辑）
	return pluginName
}

// PluginListResponse 插件列表响应
type PluginListResponse struct {
	Success bool        `json:"success"`
	Data    []PluginInfo `json:"data"`
	Total   int         `json:"total"`
}

// PluginInfo 插件信息
type PluginInfo struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	DisplayName     string                 `json:"display_name"`
	Version         string                 `json:"version"`
	Description     string                 `json:"description"`
	Author          string                 `json:"author"`
	License         string                 `json:"license"`
	Category        string                 `json:"category"`
	Status          string                 `json:"status"`
	Enabled         bool                   `json:"enabled"`
	Config          map[string]interface{} `json:"config"`
	ConfigFields    map[string]interface{} `json:"config_fields"`
	ScheduledTasks  []ScheduledTaskInfo    `json:"scheduled_tasks"`    // 定时任务列表
	HasScheduledTask bool                  `json:"has_scheduled_task"` // 是否包含定时任务
	FileSize        int64                  `json:"file_size"`
	LastUpdated     time.Time              `json:"last_updated"`
	ExecutionStats  *ExecutionStats        `json:"execution_stats,omitempty"`
}

// ScheduledTaskInfo 定时任务信息
type ScheduledTaskInfo struct {
	Name      string                 `json:"name"`       // 任务名称
	Schedule  string                 `json:"schedule"`   // 调度表达式
	Line      int                    `json:"line"`       // 所在行号
	Frequency map[string]interface{} `json:"frequency"`  // 执行频率信息
}

// ExecutionStats 执行统计
type ExecutionStats struct {
	TotalExecutions int64   `json:"total_executions"`
	SuccessRate     float64 `json:"success_rate"`
	AverageTime     int64   `json:"average_time"`
	LastExecution   *time.Time `json:"last_execution,omitempty"`
}

// GetPlugins 获取插件列表
func (h *PluginHandler) GetPlugins(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	category := c.Query("category")

	// 扫描插件目录 - 包括 hooks 和已安装的插件
	hooksPlugins, err := h.metadataParser.ScanDirectory("./hooks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to scan hooks directory",
		})
		return
	}

	// 获取已安装的插件
	installedPlugins, err := h.pluginManager.ListInstalledPlugins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get installed plugins",
		})
		return
	}

	// 合并插件列表
	plugins := hooksPlugins
	for _, installedPkg := range installedPlugins {
		// 检查是否已经在 hooks 列表中（避免重复）
		found := false
		for _, hookPkg := range hooksPlugins {
			if hookPkg.Name == installedPkg.Name {
				found = true
				break
			}
		}
		if !found {
			// 将 PluginPackage 转换为 PluginMetadata
			metadata := &plugin.PluginMetadata{
				Name:        installedPkg.Name,
				Version:     installedPkg.Version,
				Description: installedPkg.Description,
				Author:      installedPkg.Author,
				DisplayName: installedPkg.Name, // 使用 name 作为显示名
				Category:    "utility",       // 默认分类
				Status:      "installed",     // 已安装状态
				Hooks:       installedPkg.Hooks,
				ConfigFields: make(map[string]*plugin.ConfigField),
				ScheduledTasks: []*plugin.ScheduledTask{},
			}
			plugins = append(plugins, metadata)
		}
	}

	// 转换为响应格式
	var pluginInfos []PluginInfo
	for _, metadata := range plugins {
		// 从数据库获取插件配置和状态（使用新的命名系统）
		configPluginName := metadata.Name
		pluginConfig, _ := h.repoManager.PluginConfigRepository.GetConfig(configPluginName)
		enabled := pluginConfig != nil && pluginConfig.Enabled

		// 应用过滤器
		if status != "" && metadata.Status != status {
			continue
		}
		if category != "" && metadata.Category != category {
			continue
		}

		// 获取执行统计
		stats := h.getExecutionStats(metadata.Name)

		// 转换定时任务信息
		var scheduledTasks []ScheduledTaskInfo
		for _, task := range metadata.ScheduledTasks {
			frequency := make(map[string]interface{})
			if task.Frequency != nil {
				frequency["expression"] = task.Frequency.Expression
				frequency["description"] = task.Frequency.Description
				frequency["interval"] = task.Frequency.Interval
				frequency["next_run"] = task.Frequency.NextRun
			}

			scheduledTasks = append(scheduledTasks, ScheduledTaskInfo{
				Name:      task.Name,
				Schedule:  task.Schedule,
				Line:      task.Line,
				Frequency: frequency,
			})
		}

		info := PluginInfo{
			ID:              metadata.Name,
			Name:            metadata.Name,
			DisplayName:     metadata.DisplayName,
			Version:         metadata.Version,
			Description:     metadata.Description,
			Author:          metadata.Author,
			License:         metadata.License,
			Category:        metadata.Category,
			Status:          metadata.Status,
			Enabled:         enabled,
			FileSize:        metadata.FileSize,
			LastUpdated:     metadata.LastUpdated,
			ExecutionStats:  stats,
			ConfigFields:    convertConfigFields(metadata.ConfigFields),
			ScheduledTasks:  scheduledTasks,
			HasScheduledTask: metadata.HasScheduledTask,
		}

		// 解析配置
		if pluginConfig != nil && pluginConfig.ConfigJSON != "" {
			// 解析JSON配置
			var configData map[string]interface{}
			if err := json.Unmarshal([]byte(pluginConfig.ConfigJSON), &configData); err == nil {
				info.Config = configData
			} else {
				info.Config = make(map[string]interface{})
			}
		}

		pluginInfos = append(pluginInfos, info)
	}

	// 分页处理
	total := len(pluginInfos)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	if start >= total {
		pluginInfos = []PluginInfo{}
	} else {
		pluginInfos = pluginInfos[start:end]
	}

	c.JSON(http.StatusOK, PluginListResponse{
		Success: true,
		Data:    pluginInfos,
		Total:   total,
	})
}

// GetPlugin 获取插件详情
func (h *PluginHandler) GetPlugin(c *gin.Context) {
	pluginName := c.Param("name")

	// 扫描插件文件
	metadata, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".plugin.js"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 获取插件配置（使用新的命名系统）
	configPluginName := h.getPluginConfigName(pluginName)
	pluginConfig, _ := h.repoManager.PluginConfigRepository.GetConfig(configPluginName)
	enabled := pluginConfig != nil && pluginConfig.Enabled

	// 获取执行统计
	stats := h.getExecutionStats(pluginName)

	// 获取最近的日志
	logs, _ := h.repoManager.PluginLogRepository.GetRecentLogs(pluginName, 10)

	// 转换定时任务信息
	var scheduledTasks []ScheduledTaskInfo
	for _, task := range metadata.ScheduledTasks {
		frequency := make(map[string]interface{})
		if task.Frequency != nil {
			frequency["expression"] = task.Frequency.Expression
			frequency["description"] = task.Frequency.Description
			frequency["interval"] = task.Frequency.Interval
			frequency["next_run"] = task.Frequency.NextRun
		}

		scheduledTasks = append(scheduledTasks, ScheduledTaskInfo{
			Name:      task.Name,
			Schedule:  task.Schedule,
			Line:      task.Line,
			Frequency: frequency,
		})
	}

	info := PluginInfo{
		ID:              metadata.Name,
		Name:            metadata.Name,
		DisplayName:     metadata.DisplayName,
		Version:         metadata.Version,
		Description:     metadata.Description,
		Author:          metadata.Author,
		License:         metadata.License,
		Category:        metadata.Category,
		Status:          metadata.Status,
		Enabled:         enabled,
		FileSize:        metadata.FileSize,
		LastUpdated:     metadata.LastUpdated,
		ExecutionStats:  stats,
		ConfigFields:    convertConfigFields(metadata.ConfigFields),
		ScheduledTasks:  scheduledTasks,
		HasScheduledTask: metadata.HasScheduledTask,
	}

	// 解析配置
	if pluginConfig != nil && pluginConfig.ConfigJSON != "" {
		// 解析JSON配置
		var configData map[string]interface{}
		if err := json.Unmarshal([]byte(pluginConfig.ConfigJSON), &configData); err == nil {
			info.Config = configData
		} else {
			info.Config = make(map[string]interface{})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"plugin": info,
			"logs":   logs,
		},
	})
}

// EnablePlugin 启用插件
func (h *PluginHandler) EnablePlugin(c *gin.Context) {
	pluginName := c.Param("name")

	// 检查插件是否存在
	if _, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".plugin.js")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 更新数据库中的插件状态（使用新的命名系统）
	configPluginName := h.getPluginConfigName(pluginName)
	err := h.repoManager.PluginConfigRepository.SetEnabled(configPluginName, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to enable plugin",
		})
		return
	}

	// 更新元数据状态
	plugin.UpdatePluginStatus(pluginName, "enabled")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin enabled successfully",
	})
}

// DisablePlugin 禁用插件
func (h *PluginHandler) DisablePlugin(c *gin.Context) {
	pluginName := c.Param("name")

	// 检查插件是否存在
	if _, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".plugin.js")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 更新数据库中的插件状态（使用新的命名系统）
	configPluginName := h.getPluginConfigName(pluginName)
	err := h.repoManager.PluginConfigRepository.SetEnabled(configPluginName, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to disable plugin",
		})
		return
	}

	// 更新元数据状态
	plugin.UpdatePluginStatus(pluginName, "disabled")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin disabled successfully",
	})
}

// UpdatePluginConfig 更新插件配置
func (h *PluginHandler) UpdatePluginConfig(c *gin.Context) {
	pluginName := c.Param("name")

	var request struct {
		Config map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// 检查插件是否存在
	if _, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".plugin.js")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 这里应该验证配置格式
	// TODO: 使用JSON Schema验证配置

	// 更新数据库中的配置（使用新的命名系统）
	configPluginName := h.getPluginConfigName(pluginName)
	err := h.repoManager.PluginConfigRepository.SetConfig(configPluginName, request.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update plugin config",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin config updated successfully",
	})
}

// GetPluginStats 获取插件统计信息
func (h *PluginHandler) GetPluginStats(c *gin.Context) {
	// 扫描所有插件
	plugins, err := h.metadataParser.ScanDirectory("./hooks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to scan plugins",
		})
		return
	}

	totalPlugins := len(plugins)
	enabledPlugins := 0
	disabledPlugins := 0
	totalExecutions := int64(0)
	totalErrors := int64(0)

	for _, metadata := range plugins {
		// 获取插件状态（使用新的命名系统）
		configPluginName := metadata.Name
		pluginConfig, _ := h.repoManager.PluginConfigRepository.GetConfig(configPluginName)
		if pluginConfig != nil && pluginConfig.Enabled {
			enabledPlugins++
		} else {
			disabledPlugins++
		}

		// 获取执行统计
		stats := h.getExecutionStats(metadata.Name)
		totalExecutions += stats.TotalExecutions
		// 这里应该统计错误数量
	}

	successRate := float64(100)
	if totalExecutions > 0 {
		successRate = float64(totalExecutions-totalErrors) / float64(totalExecutions) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_plugins":    totalPlugins,
			"enabled_plugins":  enabledPlugins,
			"disabled_plugins": disabledPlugins,
			"total_executions": totalExecutions,
			"success_rate":     successRate,
		},
	})
}

// GetPluginLogs 获取插件日志
func (h *PluginHandler) GetPluginLogs(c *gin.Context) {
	pluginName := c.Param("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// 获取插件日志
	logs, total, err := h.repoManager.PluginLogRepository.GetLogs(pluginName, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get plugin logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":  logs,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// getExecutionStats 获取插件执行统计
func (h *PluginHandler) getExecutionStats(pluginName string) *ExecutionStats {
	// 这里应该从数据库获取实际的统计数据
	// 暂时返回模拟数据
	stats := &ExecutionStats{
		TotalExecutions: 1000,
		SuccessRate:     98.5,
		AverageTime:     15,
	}

	// 设置最近执行时间
	now := time.Now()
	stats.LastExecution = &now

	return stats
}

// InstallPlugin 安装插件
func (h *PluginHandler) InstallPlugin(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	// 检查是否是文件上传请求
	if strings.Contains(contentType, "multipart/form-data") {
		// 文件上传请求
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("File upload failed: %v", err),
			})
			return
		}
		defer file.Close()

		// 检查文件名
		filename := header.Filename
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "File name is required",
			})
			return
		}

		// 保存上传的文件到临时位置
		tempDir := "./temp"
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Failed to create temp directory: %v", err),
			})
			return
		}

		tempPath := filepath.Join(tempDir, filename)
		if err := c.SaveUploadedFile(header, tempPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Failed to save uploaded file: %v", err),
			})
			return
		}

		// 使用插件管理器安装插件
		if err := h.pluginManager.InstallPlugin(tempPath); err != nil {
			// 清理临时文件
			os.Remove(tempPath)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Failed to install plugin: %v", err),
			})
			return
		}

		// 清理临时文件
		os.Remove(tempPath)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Plugin installed successfully",
		})
		return
	}

	// JSON格式请求（URL安装）
	var jsonRequest struct {
		Source string `json:"source"` // 文件路径或URL
	}
	if err := c.ShouldBindJSON(&jsonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: expected JSON with 'source' field or multipart/form-data file upload",
		})
		return
	}

	if jsonRequest.Source == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Source is required",
		})
		return
	}

	// 使用插件管理器安装插件
	if err := h.pluginManager.InstallPlugin(jsonRequest.Source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to install plugin: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin installed successfully",
	})
}

// UninstallPlugin 卸载插件
func (h *PluginHandler) UninstallPlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Plugin name is required",
		})
		return
	}

	// 检查插件是否已安装（包括已安装目录和hooks目录）
	isInstalled := h.pluginManager.IsPluginInstalled(pluginName)
	isHooksPlugin := false

	// 如果不在已安装目录中，检查hooks目录
	if !isInstalled {
		hooksFile := filepath.Join("./hooks", pluginName+".plugin.js")
		if _, err := os.Stat(hooksFile); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Plugin is not installed",
			})
			return
		}
		// 标记为hooks目录中的插件
		isInstalled = true
		isHooksPlugin = true
	}

	// 卸载插件
	var err error
	if isHooksPlugin {
		// 对于hooks目录中的插件，直接删除文件
		hooksFile := filepath.Join("./hooks", pluginName+".plugin.js")
		err = os.Remove(hooksFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Failed to delete hooks plugin file: %v", err),
			})
			return
		}
	} else {
		// 对于已安装目录中的插件，使用标准卸载方法
		err = h.pluginManager.UninstallPlugin(pluginName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Failed to uninstall plugin: %v", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin uninstalled successfully",
	})
}

// LoadPlugin 加载插件
func (h *PluginHandler) LoadPlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Plugin name is required",
		})
		return
	}

	// 检查插件是否已安装
	if !h.pluginManager.IsPluginInstalled(pluginName) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin is not installed",
		})
		return
	}

	// 检查插件是否已加载
	if h.pluginManager.IsPluginLoaded(pluginName) {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Plugin is already loaded",
		})
		return
	}

	// 加载插件
	if err := h.pluginManager.LoadPlugin(pluginName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to load plugin: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin loaded successfully",
	})
}

// UnloadPlugin 卸载已加载的插件
func (h *PluginHandler) UnloadPlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Plugin name is required",
		})
		return
	}

	// 检查插件是否已加载
	if !h.pluginManager.IsPluginLoaded(pluginName) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin is not loaded",
		})
		return
	}

	// 卸载插件
	if err := h.pluginManager.UnloadPlugin(pluginName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to unload plugin: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin unloaded successfully",
	})
}

// ReloadPlugin 重新加载插件
func (h *PluginHandler) ReloadPlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Plugin name is required",
		})
		return
	}

	// 检查插件是否已安装
	if !h.pluginManager.IsPluginInstalled(pluginName) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin is not installed",
		})
		return
	}

	// 重新加载插件
	if err := h.pluginManager.ReloadPlugin(pluginName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to reload plugin: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Plugin reloaded successfully",
	})
}

// GetInstalledPlugins 获取已安装的插件列表
func (h *PluginHandler) GetInstalledPlugins(c *gin.Context) {
	plugins, err := h.pluginManager.ListInstalledPlugins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get installed plugins",
		})
		return
	}

	// 转换为响应格式
	var pluginInfos []PluginInfo
	for _, pkg := range plugins {
		pluginInfo := PluginInfo{
			ID:           pkg.Name,
			Name:         pkg.Name,
			Version:      pkg.Version,
			Description:  pkg.Description,
			Author:       pkg.Author,
			Status:       "installed",
			Enabled:      h.pluginManager.IsPluginLoaded(pkg.Name),
			ConfigFields: convertConfigFields(nil), // 插件包配置字段
		}

		pluginInfos = append(pluginInfos, pluginInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pluginInfos,
		"total":   len(pluginInfos),
	})
}

// convertConfigFields 转换配置字段为JSON格式
func convertConfigFields(fields map[string]*plugin.ConfigField) map[string]interface{} {
	if fields == nil {
		return nil
	}

	result := make(map[string]interface{})
	for name, field := range fields {
		result[name] = map[string]interface{}{
			"type":        field.Type,
			"name":        field.Name,
			"label":       field.Label,
			"description": field.Description,
			"required":    field.Required,
			"default":     field.Default,
			"options":     field.Options,
			"validation":  field.Validation,
		}
	}
	return result
}