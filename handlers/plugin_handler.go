package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin"
	"github.com/gin-gonic/gin"
)

// PluginHandler 插件管理处理器
type PluginHandler struct {
	repoManager   *repo.RepositoryManager
	metadataParser *plugin.MetadataParser
}

// NewPluginHandler 创建插件处理器
func NewPluginHandler(repoManager *repo.RepositoryManager) *PluginHandler {
	return &PluginHandler{
		repoManager:   repoManager,
		metadataParser: plugin.NewMetadataParser(),
	}
}

// PluginListResponse 插件列表响应
type PluginListResponse struct {
	Success bool        `json:"success"`
	Data    []PluginInfo `json:"data"`
	Total   int         `json:"total"`
}

// PluginInfo 插件信息
type PluginInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	License     string                 `json:"license"`
	Category    string                 `json:"category"`
	Status      string                 `json:"status"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
	FileSize    int64                  `json:"file_size"`
	LastUpdated time.Time              `json:"last_updated"`
	ExecutionStats *ExecutionStats     `json:"execution_stats,omitempty"`
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

	// 扫描插件目录
	plugins, err := h.metadataParser.ScanDirectory("./hooks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to scan plugins directory",
		})
		return
	}

	// 转换为响应格式
	var pluginInfos []PluginInfo
	for _, metadata := range plugins {
		// 从数据库获取插件配置和状态
		pluginConfig, _ := h.repoManager.PluginConfigRepository.GetConfig(metadata.Name)
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

		info := PluginInfo{
			ID:          metadata.Name,
			Name:        metadata.Name,
			Version:     metadata.Version,
			Description: metadata.Description,
			Author:      metadata.Author,
			License:     metadata.License,
			Category:    metadata.Category,
			Status:      metadata.Status,
			Enabled:     enabled,
			FileSize:    metadata.FileSize,
			LastUpdated: metadata.LastUpdated,
			ExecutionStats: stats,
		}

		// 解析配置
		if pluginConfig != nil {
			// 这里应该解析JSON配置
			info.Config = make(map[string]interface{})
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
	metadata, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".pb.js"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 获取插件配置（使用带.pb后缀的名称保持一致性）
	configPluginName := pluginName + ".pb"
	pluginConfig, _ := h.repoManager.PluginConfigRepository.GetConfig(configPluginName)
	enabled := pluginConfig != nil && pluginConfig.Enabled

	// 获取执行统计
	stats := h.getExecutionStats(pluginName)

	// 获取最近的日志
	logs, _ := h.repoManager.PluginLogRepository.GetRecentLogs(pluginName, 10)

	info := PluginInfo{
		ID:          metadata.Name,
		Name:        metadata.Name,
		Version:     metadata.Version,
		Description: metadata.Description,
		Author:      metadata.Author,
		License:     metadata.License,
		Category:    metadata.Category,
		Status:      metadata.Status,
		Enabled:     enabled,
		FileSize:    metadata.FileSize,
		LastUpdated: metadata.LastUpdated,
		ExecutionStats: stats,
	}

	// 解析配置
	if pluginConfig != nil {
		// 这里应该解析JSON配置
		info.Config = make(map[string]interface{})
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
	if _, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".pb.js")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 更新数据库中的插件状态（使用带.pb后缀的名称保持一致性）
	configPluginName := pluginName + ".pb"
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
	if _, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".pb.js")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 更新数据库中的插件状态（使用带.pb后缀的名称保持一致性）
	configPluginName := pluginName + ".pb"
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
	if _, err := h.metadataParser.ParseFile(filepath.Join("./hooks", pluginName+".pb.js")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Plugin not found",
		})
		return
	}

	// 这里应该验证配置格式
	// TODO: 使用JSON Schema验证配置

	// 更新数据库中的配置（使用带.pb后缀的名称保持一致性）
	configPluginName := pluginName + ".pb"
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
		// 获取插件状态
		pluginConfig, _ := h.repoManager.PluginConfigRepository.GetConfig(metadata.Name)
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