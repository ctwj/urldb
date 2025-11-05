package handlers

import (
	"fmt"
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

	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	if err := manager.InitializePlugin(pluginName); err != nil {
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

	// 获取插件实例
	instance, err := manager.GetPluginInstance(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 如果插件支持配置接口，获取配置
	config := make(map[string]interface{})
	if configurablePlugin, ok := instance.Plugin.(interface{ GetConfig() map[string]interface{} }); ok {
		config = configurablePlugin.GetConfig()
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

	// 获取插件实例验证插件是否存在
	instance, err := manager.GetPluginInstance(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 检查插件是否支持配置更新
	configurablePlugin, ok := instance.Plugin.(interface{ UpdateConfig(map[string]interface{}) error })
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin does not support configuration updates"})
		return
	}

	// 如果插件正在运行，先停止
	instanceInfo, _ := manager.GetPluginInstance(pluginName)
	if instanceInfo.Status == types.StatusRunning {
		if err := manager.StopPlugin(pluginName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop plugin for configuration update: " + err.Error()})
			return
		}
	}

	// 更新配置
	if err := configurablePlugin.UpdateConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration: " + err.Error()})
		return
	}

	// 如果插件之前是运行状态，重新启动
	if instanceInfo.Status == types.StatusRunning {
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

	// 获取插件实例
	instance, err := manager.GetPluginInstance(pluginName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 获取依赖项列表（如果插件支持依赖接口）
	dependencies := make([]string, 0)
	dependents := make([]string, 0)

	if dependencyPlugin, ok := instance.Plugin.(interface{ Dependencies() []string }); ok {
		dependencies = dependencyPlugin.Dependencies()
	}

	c.JSON(http.StatusOK, gin.H{
		"plugin_info":  pluginInfo,
		"dependencies": dependencies,
		"dependents":   dependents,
	})
}

// GetPluginLoadOrder 获取插件加载顺序
func (ph *PluginHandler) GetPluginLoadOrder(c *gin.Context) {
	manager := plugin.GetManager()
	if manager == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Plugin manager not initialized"})
		return
	}

	// 简化版管理器直接返回所有插件名称
	plugins := manager.ListPlugins()
	loadOrder := make([]string, len(plugins))
	for i, plugin := range plugins {
		loadOrder[i] = plugin.Name
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

	// 检查是否有插件注册
	plugins := manager.ListPlugins()

	c.JSON(http.StatusOK, gin.H{
		"valid":   len(plugins) > 0, // 简单验证：如果有插件则认为有效
		"count":   len(plugins),
		"plugins": plugins,
		"message": fmt.Sprintf("Found %d plugins", len(plugins)),
	})
}