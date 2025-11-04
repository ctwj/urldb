package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/monitor"
	"github.com/gin-gonic/gin"
)

// PluginMonitorHandler 插件监控处理器
type PluginMonitorHandler struct {
	monitor *monitor.PluginMonitor
	checker *monitor.PluginHealthChecker
}

// NewPluginMonitorHandler 创建插件监控处理器
func NewPluginMonitorHandler() *PluginMonitorHandler {
	pluginMonitor := plugin.GetPluginMonitor()
	healthChecker := monitor.NewPluginHealthChecker(pluginMonitor)

	return &PluginMonitorHandler{
		monitor: pluginMonitor,
		checker: healthChecker,
	}
}

// GetPluginHealth 获取插件健康状态
func (pmh *PluginMonitorHandler) GetPluginHealth(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	// 获取插件管理器
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	// 获取插件
	pluginInstance, err := manager.GetPlugin(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	// 执行健康检查
	result := pmh.checker.Check(c.Request.Context(), pluginInstance)

	c.JSON(http.StatusOK, result)
}

// GetAllPluginsHealth 获取所有插件健康状态
func (pmh *PluginMonitorHandler) GetAllPluginsHealth(c *gin.Context) {
	// 获取插件管理器
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	// 获取所有插件
	plugins := manager.GetEnabledPlugins()

	// 执行批量健康检查
	results := pmh.checker.BatchCheck(c.Request.Context(), plugins)

	// 生成报告
	report := pmh.checker.GenerateReport(results)

	c.JSON(http.StatusOK, report)
}

// GetPluginActivities 获取插件活动记录
func (pmh *PluginMonitorHandler) GetPluginActivities(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	limit := 100 // 默认限制100条记录
	if limitParam := c.Query("limit"); limitParam != "" {
		// 这里可以解析limit参数
	}

	activities := pmh.monitor.GetActivities(pluginName, limit)

	c.JSON(http.StatusOK, gin.H{
		"plugin_name": pluginName,
		"activities":  activities,
		"count":       len(activities),
	})
}

// GetPluginMetrics 获取插件指标
func (pmh *PluginMonitorHandler) GetPluginMetrics(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	metrics := pmh.monitor.GetMonitorStats()

	c.JSON(http.StatusOK, metrics)
}

// GetAlertRules 获取告警规则
func (pmh *PluginMonitorHandler) GetAlertRules(c *gin.Context) {
	rules := pmh.monitor.GetAlertRules()

	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
		"count": len(rules),
	})
}

// CreateAlertRule 创建告警规则
func (pmh *PluginMonitorHandler) CreateAlertRule(c *gin.Context) {
	var rule monitor.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pmh.monitor.SetAlertRule(rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alert rule created successfully",
		"rule":    rule,
	})
}

// DeleteAlertRule 删除告警规则
func (pmh *PluginMonitorHandler) DeleteAlertRule(c *gin.Context) {
	ruleName := c.Param("name")
	if ruleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rule name is required"})
		return
	}

	pmh.monitor.RemoveAlertRule(ruleName)

	c.JSON(http.StatusOK, gin.H{
		"message": "Alert rule deleted successfully",
	})
}

// GetPluginMonitorStats 获取插件监控统计信息
func (pmh *PluginMonitorHandler) GetPluginMonitorStats(c *gin.Context) {
	stats := pmh.monitor.GetMonitorStats()

	c.JSON(http.StatusOK, stats)
}

// StreamAlerts 流式获取告警信息
func (pmh *PluginMonitorHandler) StreamAlerts(c *gin.Context) {
	// 设置SSE头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 获取告警通道
	alerts := pmh.monitor.GetAlerts()

	// 客户端断开连接时的处理
	clientGone := c.Request.Context().Done()

	// 发送初始连接消息
	c.Stream(func(w io.Writer) bool {
		select {
		case alert := <-alerts:
			// 发送告警信息
			c.SSEvent("alert", alert)
			return true
		case <-clientGone:
			// 客户端断开连接
			return false
		case <-time.After(30 * time.Second):
			// 发送心跳消息保持连接
			c.SSEvent("ping", time.Now().Unix())
			return true
		}
	})
}

// GetPluginHealthHistory 获取插件健康历史
func (pmh *PluginMonitorHandler) GetPluginHealthHistory(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	// 获取插件管理器
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	// 获取插件信息
	pluginInfo, err := manager.GetPluginInfo(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	// 获取活动记录作为健康历史的替代
	activities := pmh.monitor.GetActivities(pluginName, 50)

	// 构建健康历史数据
	history := make([]map[string]interface{}, 0)
	for _, activity := range activities {
		if activity.Error != nil {
			history = append(history, map[string]interface{}{
				"timestamp": activity.Timestamp,
				"status":    "error",
				"message":   activity.Error.Error(),
				"operation": activity.Operation,
				"duration":  activity.ExecutionTime,
			})
		} else {
			history = append(history, map[string]interface{}{
				"timestamp": activity.Timestamp,
				"status":    "success",
				"operation": activity.Operation,
				"duration":  activity.ExecutionTime,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"plugin":  pluginInfo,
		"history": history,
		"count":   len(history),
	})
}