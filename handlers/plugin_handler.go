package handlers

import (
	"net/http"

	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/gin-gonic/gin"
)

// PluginHandler 插件管理处理器
type PluginHandler struct{}

// NewPluginHandler 创建插件管理处理器
func NewPluginHandler() *PluginHandler {
	return &PluginHandler{}
}

// GetPlugins 获取所有插件信息
func (ph *PluginHandler) GetPlugins(c *gin.Context) {
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	plugins := manager.ListPlugins()
	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
		"count":   len(plugins),
	})
}

// GetPlugin 获取指定插件信息
func (ph *PluginHandler) GetPlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	pluginInfo, err := manager.GetPluginInfo(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pluginInfo)
}

// InstallPlugin 安装插件（预留接口，具体实现依赖插件加载机制）
func (ph *PluginHandler) InstallPlugin(c *gin.Context) {
	// TODO: 实现插件安装功能
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Plugin installation is not implemented yet"})
}

// UninstallPlugin 卸载插件
func (ph *PluginHandler) UninstallPlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	force := false
	if c.Query("force") == "true" {
		force = true
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	// 检查是否可以安全卸载
	canUninstall, dependents, err := manager.CanUninstall(pluginName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !canUninstall && !force {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Plugin cannot be safely uninstalled",
			"dependents": dependents,
			"can_force":  true,
		})
		return
	}

	if err := manager.UninstallPlugin(pluginName, force); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin uninstalled successfully",
	})
}

// InitializePlugin 初始化插件
func (ph *PluginHandler) InitializePlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration format: " + err.Error()})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	if err := manager.InitializePlugin(pluginName, config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin initialized successfully",
	})
}

// StartPlugin 启动插件
func (ph *PluginHandler) StartPlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	if err := manager.StartPlugin(pluginName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin started successfully",
	})
}

// StopPlugin 停止插件
func (ph *PluginHandler) StopPlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	if err := manager.StopPlugin(pluginName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin stopped successfully",
	})
}

// GetPluginConfig 获取插件配置
func (ph *PluginHandler) GetPluginConfig(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	config, err := manager.GetLatestConfigVersion(pluginName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plugin_name": pluginName,
		"config":      config,
	})
}

// UpdatePluginConfig 更新插件配置
func (ph *PluginHandler) UpdatePluginConfig(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration format: " + err.Error()})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	// 获取插件信息验证插件是否存在
	_, err := manager.GetPluginInfo(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 如果插件正在运行，先停止
	status := manager.GetPluginStatus(pluginName)
	if status == types.StatusRunning {
		if err := manager.StopPlugin(pluginName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop plugin for configuration update: " + err.Error()})
			return
		}
	}

	// 保存新的配置
	if err := manager.SaveConfigVersion(pluginName, "latest", "Configuration updated via API", "system", config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save configuration: " + err.Error()})
		return
	}

	// 如果插件之前是运行状态，重新启动
	if status == types.StatusRunning {
		if err := manager.StartPlugin(pluginName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restart plugin after configuration update: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin configuration updated successfully",
		"config":  config,
	})
}

// GetPluginDependencies 获取插件依赖信息
func (ph *PluginHandler) GetPluginDependencies(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin name is required"})
		return
	}

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	pluginInfo, err := manager.GetPluginInfo(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 获取依赖检查结果
	dependenciesStatus := manager.CheckAllDependencies()

	// 获取依赖项列表
	dependents := manager.GetDependents(pluginName)

	c.JSON(http.StatusOK, gin.H{
		"plugin_info":    pluginInfo,
		"dependencies":   dependenciesStatus[pluginName],
		"dependents":     dependents,
		"dependencies_status": dependenciesStatus,
	})
}

// GetPluginLoadOrder 获取插件加载顺序
func (ph *PluginHandler) GetPluginLoadOrder(c *gin.Context) {
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	loadOrder, err := manager.GetLoadOrder()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"load_order": loadOrder,
		"count":      len(loadOrder),
	})
}

// ValidatePluginDependencies 验证插件依赖
func (ph *PluginHandler) ValidatePluginDependencies(c *gin.Context) {
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	err := manager.ValidateDependencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      err.Error(),
			"valid":      false,
			"message":    "Plugin dependencies validation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "All plugin dependencies are satisfied",
	})
}