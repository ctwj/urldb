package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin/config"
	"github.com/ctwj/urldb/plugin/concurrency"
	"github.com/ctwj/urldb/plugin/monitor"
	"github.com/ctwj/urldb/plugin/registry"
	"github.com/ctwj/urldb/plugin/security"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// Manager is the plugin manager that handles plugin lifecycle
type Manager struct {
	plugins       map[string]types.Plugin
	instances     map[string]*types.PluginInstance
	registry      *registry.PluginRegistry
	depManager    *DependencyManager
	securityManager *security.SecurityManager
	monitor       *monitor.PluginMonitor
	configManager *config.ConfigManager
	lazyLoader    *LazyLoader
	concurrencyCtrl *concurrency.ConcurrencyController
	mutex         sync.RWMutex
	taskManager   interface{} // Reference to the existing task manager
	repoManager   *repo.RepositoryManager
	database      *gorm.DB
}


// NewManager creates a new plugin manager
func NewManager(taskManager interface{}, repoManager *repo.RepositoryManager, database *gorm.DB, pluginMonitor *monitor.PluginMonitor) *Manager {
	securityManager := security.NewSecurityManager(database)
	configManager := config.NewConfigManager()

	// 创建全局并发控制器，全局限制为50个并发任务
	concurrencyCtrl := concurrency.NewConcurrencyController(50)

	manager := &Manager{
		plugins:        make(map[string]types.Plugin),
		instances:      make(map[string]*types.PluginInstance),
		registry:       registry.NewPluginRegistry(),
		securityManager: securityManager,
		monitor:       pluginMonitor,
		configManager: configManager,
		concurrencyCtrl: concurrencyCtrl,
		taskManager:    taskManager,
		repoManager:    repoManager,
		database:       database,
	}
	manager.depManager = NewDependencyManager(manager)
	manager.lazyLoader = NewLazyLoader(manager)
	return manager
}

// GetLazyLoader returns the lazy loader
func (m *Manager) GetLazyLoader() *LazyLoader {
	return m.lazyLoader
}

// SetPluginConcurrencyLimit sets the concurrency limit for a specific plugin
func (m *Manager) SetPluginConcurrencyLimit(pluginName string, limit int) {
	if m.concurrencyCtrl != nil {
		m.concurrencyCtrl.SetPluginLimit(pluginName, limit)
	}
}

// GetConcurrencyStats returns concurrency control statistics
func (m *Manager) GetConcurrencyStats() map[string]interface{} {
	if m.concurrencyCtrl != nil {
		return m.concurrencyCtrl.GetStats()
	}
	return nil
}

// RegisterPlugin registers a plugin with the manager
func (m *Manager) RegisterPlugin(plugin types.Plugin) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	name := plugin.Name()
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	m.plugins[name] = plugin

	// 如果插件实现了ConfigurablePlugin接口，注册其配置模式和模板
	if configurablePlugin, ok := plugin.(types.ConfigurablePlugin); ok {
		// 注册配置模式
		schema := configurablePlugin.CreateConfigSchema()
		if schema != nil {
			if err := m.configManager.RegisterSchema(schema); err != nil {
				utils.Warn("Failed to register config schema for plugin %s: %v", name, err)
			} else {
				utils.Info("Config schema registered for plugin %s", name)
			}
		}

		// 注册配置模板
		template := configurablePlugin.CreateConfigTemplate()
		if template != nil {
			if err := m.configManager.RegisterTemplate(template); err != nil {
				utils.Warn("Failed to register config template for plugin %s: %v", name, err)
			} else {
				utils.Info("Config template registered for plugin %s", name)
			}
		}
	}

	utils.Info("Plugin registered: %s", name)
	return nil
}

// ValidateDependencies validates all plugin dependencies
func (m *Manager) ValidateDependencies() error {
	return m.depManager.ValidateDependencies()
}

// CheckPluginDependencies checks if all dependencies for a specific plugin are satisfied
func (m *Manager) CheckPluginDependencies(pluginName string) (bool, []string, error) {
	return m.depManager.CheckPluginDependencies(pluginName)
}

// GetLoadOrder returns the correct order to load plugins based on dependencies
func (m *Manager) GetLoadOrder() ([]string, error) {
	return m.depManager.GetLoadOrder()
}

// GetDependencyInfo returns dependency information for a plugin
func (m *Manager) GetDependencyInfo(pluginName string) (*types.PluginInfo, error) {
	return m.depManager.GetDependencyInfo(pluginName)
}

// CheckAllDependencies checks all plugin dependencies
func (m *Manager) CheckAllDependencies() map[string]map[string]bool {
	return m.depManager.CheckAllDependencies()
}

// UnregisterPlugin unregisters a plugin from the manager
func (m *Manager) UnregisterPlugin(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	delete(m.plugins, name)
	utils.Info("Plugin unregistered: %s", name)
	return nil
}

// GetPlugin returns a plugin by name
func (m *Manager) GetPlugin(name string) (types.Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// ListPlugins returns a list of all registered plugins
func (m *Manager) ListPlugins() []types.PluginInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var plugins []types.PluginInfo
	for name, plugin := range m.plugins {
		info := types.PluginInfo{
			Name:        name,
			Version:     plugin.Version(),
			Description: plugin.Description(),
			Author:      plugin.Author(),
		}

		// Get instance status if available
		if instance, exists := m.instances[name]; exists {
			info.Status = instance.Status
			info.LastError = instance.LastError
			info.StartTime = instance.StartTime
			info.StopTime = instance.StopTime
			info.RestartCount = instance.RestartCount
			info.HealthScore = instance.HealthScore
			// 添加监控相关字段
			info.TotalExecutionTime = instance.TotalExecutionTime
			info.TotalExecutions = instance.TotalExecutions
			info.TotalErrors = instance.TotalErrors
			info.LastExecutionTime = instance.LastExecutionTime
		} else {
			info.Status = types.StatusRegistered
		}

		plugins = append(plugins, info)
	}

	return plugins
}

// InitializePlugin initializes a plugin
func (m *Manager) InitializePlugin(name string, config map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Check if all dependencies are satisfied before initialization
	satisfied, unresolved, err := m.CheckPluginDependencies(name)
	if err != nil {
		return fmt.Errorf("failed to check dependencies for plugin %s: %v", name, err)
	}

	if !satisfied {
		return fmt.Errorf("plugin %s has unsatisfied dependencies: %v", name, unresolved)
	}

	// Validate plugin configuration if schema exists
	if err := m.validatePluginConfig(name, config); err != nil {
		return fmt.Errorf("plugin configuration validation failed: %v", err)
	}

	// Apply default values if schema exists
	if err := m.applyConfigDefaults(name, config); err != nil {
		utils.Warn("Failed to apply config defaults for plugin %s: %v", name, err)
	}

	// Create plugin context with repo manager and database
	context := NewPluginContext(name, m, config, m.repoManager, m.database)

	// 如果有监控器，创建增强的上下文
	var enhancedContext types.PluginContext = context
	if m.monitor != nil {
		// 创建插件实例（用于监控）
		instance := &types.PluginInstance{
			Plugin:  plugin,
			Context: context,
			Status:  types.StatusInitialized,
			Config:  config,
		}
		enhancedContext = monitor.NewEnhancedPluginContext(context, instance, m.monitor, name)
	}

	// Create plugin instance
	instance := &types.PluginInstance{
		Plugin:  plugin,
		Context: enhancedContext,
		Status:  types.StatusInitialized,
		Config:  config,
	}

	// Initialize the plugin
	if err := plugin.Initialize(enhancedContext); err != nil {
		instance.Status = types.StatusError
		instance.LastError = err
		m.instances[name] = instance
		return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
	}

	// Save configuration version
	if err := m.configManager.SaveVersion(name, plugin.Version(), "Initial configuration", "system", config); err != nil {
		utils.Warn("Failed to save initial config version for plugin %s: %v", name, err)
	}

	m.instances[name] = instance
	utils.Info("Plugin initialized: %s", name)
	return nil
}

// StartPlugin starts a plugin
func (m *Manager) StartPlugin(name string) error {
	startTime := time.Now()
	m.mutex.Lock()
	defer m.mutex.Unlock()

	instance, exists := m.instances[name]
	if !exists {
		return fmt.Errorf("plugin instance %s not found", name)
	}

	if instance.Status != types.StatusInitialized && instance.Status != types.StatusStopped {
		return fmt.Errorf("plugin %s is not in a startable state, current status: %s", name, instance.Status)
	}

	instance.Status = types.StatusStarting
	startErr := instance.Plugin.Start()

	// 记录执行时间和错误
	executionTime := time.Since(startTime)
	if m.monitor != nil {
		// 这里可以记录到监控器
	}

	if startErr != nil {
		instance.Status = types.StatusError
		instance.LastError = startErr
		instance.UpdateExecutionStats(executionTime, startErr)
		return fmt.Errorf("failed to start plugin %s: %v", name, startErr)
	}

	instance.Status = types.StatusRunning
	instance.StartTime = time.Now()
	instance.RestartCount++
	instance.UpdateExecutionStats(executionTime, nil)
	utils.Info("Plugin started: %s", name)
	return nil
}

// StopPlugin stops a plugin
func (m *Manager) StopPlugin(name string) error {
	startTime := time.Now()
	m.mutex.Lock()
	defer m.mutex.Unlock()

	instance, exists := m.instances[name]
	if !exists {
		return fmt.Errorf("plugin instance %s not found", name)
	}

	if instance.Status != types.StatusRunning {
		return fmt.Errorf("plugin %s is not running, current status: %s", name, instance.Status)
	}

	instance.Status = types.StatusStopping
	stopErr := instance.Plugin.Stop()

	// 记录执行时间和错误
	executionTime := time.Since(startTime)
	if m.monitor != nil {
		// 这里可以记录到监控器
	}

	if stopErr != nil {
		instance.Status = types.StatusError
		instance.LastError = stopErr
		instance.UpdateExecutionStats(executionTime, stopErr)
		return fmt.Errorf("failed to stop plugin %s: %v", name, stopErr)
	}

	now := time.Now()
	instance.Status = types.StatusStopped
	instance.StopTime = &now
	instance.UpdateExecutionStats(executionTime, nil)
	utils.Info("Plugin stopped: %s", name)
	return nil
}

// GetPluginStatus returns the status of a plugin
func (m *Manager) GetPluginStatus(name string) types.PluginStatus {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if instance, exists := m.instances[name]; exists {
		return instance.Status
	}

	if _, exists := m.plugins[name]; exists {
		return types.StatusRegistered
	}

	return types.StatusDisabled
}

// GetEnabledPlugins returns all enabled plugins
func (m *Manager) GetEnabledPlugins() []types.Plugin {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var enabled []types.Plugin
	for name, plugin := range m.plugins {
		if instance, exists := m.instances[name]; exists {
			if instance.Status == types.StatusRunning {
				enabled = append(enabled, plugin)
			}
		}
	}

	return enabled
}

// UninstallPlugin uninstalls a plugin, stopping it first and performing cleanup
func (m *Manager) UninstallPlugin(name string, force bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if plugin exists
	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Check if any other plugins depend on this plugin
	dependents := m.depManager.GetDependents(name)
	if len(dependents) > 0 && !force {
		return fmt.Errorf("plugin %s cannot be uninstalled because the following plugins depend on it: %v", name, dependents)
	}

	// Stop the plugin if it's running
	if instance, exists := m.instances[name]; exists {
		if instance.Status == types.StatusRunning {
			utils.Info("Stopping plugin %s before uninstall", name)
			if err := m.stopPluginInstance(instance); err != nil {
				if !force {
					return fmt.Errorf("failed to stop plugin %s before uninstall: %v", name, err)
				}
				utils.Warn("Forcing uninstall despite stop failure: %v", err)
			}
		}
	}

	// Get plugin data before removal for cleanup operations
	// var pluginData map[string]interface{}
	// if instance, exists := m.instances[name]; exists {
	// 	pluginData = instance.Config
	// }

	// Perform plugin cleanup
	if err := plugin.Cleanup(); err != nil {
		if !force {
			return fmt.Errorf("plugin %s cleanup failed: %v", name, err)
		}
		utils.Warn("Plugin cleanup failed but continuing with uninstall: %v", err)
	}

	// Clean up plugin data and configurations
	if err := m.cleanupPluginData(name); err != nil {
		utils.Warn("Error during plugin data cleanup: %v", err)
	}

	// Clean up any registered tasks by the plugin
	if err := m.cleanupPluginTasks(name); err != nil {
		utils.Warn("Error during plugin task cleanup: %v", err)
	}

	// Remove the plugin from instances and plugins map
	delete(m.instances, name)
	delete(m.plugins, name)

	// Update dependency graph
	m.depManager.RemovePlugin(name)

	utils.Info("Plugin uninstalled: %s", name)
	return nil
}

// stopPluginInstance stops a plugin instance (internal method)
func (m *Manager) stopPluginInstance(instance *types.PluginInstance) error {
	instance.Status = types.StatusStopping
	if err := instance.Plugin.Stop(); err != nil {
		instance.Status = types.StatusError
		instance.LastError = err
		return err
	}

	now := time.Now()
	instance.Status = types.StatusStopped
	instance.StopTime = &now
	utils.Info("Plugin stopped: %s", instance.Plugin.Name())
	return nil
}

// cleanupPluginData cleans up plugin-specific data
func (m *Manager) cleanupPluginData(name string) error {
	// Clean up plugin configurations from database
	if m.repoManager != nil && m.database != nil {
		// Get the plugin configuration repository from the manager
		pluginConfigRepo := m.repoManager.PluginConfigRepository
		if pluginConfigRepo != nil {
			// Remove plugin configurations
			if err := pluginConfigRepo.DeleteByPlugin(name); err != nil {
				utils.Error("Error deleting plugin config: %v", err)
				// Don't return error here as it's not critical for uninstallation
			}
		}

		// Clean up plugin data from database
		pluginDataRepo := m.repoManager.PluginDataRepository
		if pluginDataRepo != nil {
			// Remove plugin data
			// We'll need to handle this differently since the delete method requires dataType
			// Get all plugin data first, then delete by plugin name and data types
			var allPluginData []entity.PluginData
			if err := m.database.Where("plugin_name = ?", name).Find(&allPluginData).Error; err != nil {
				utils.Error("Error finding plugin data: %v", err)
			} else {
				// Delete each data type separately
				for _, data := range allPluginData {
					if err := pluginDataRepo.DeleteByPluginAndType(name, data.DataType); err != nil {
						utils.Error("Error deleting plugin data by type: %v", err)
					}
				}
			}
		}
	}

	// Remove any plugin-specific files or resources
	// (Implementation depends on plugin-specific needs)

	return nil
}

// cleanupPluginTasks cleans up plugin-registered tasks
func (m *Manager) cleanupPluginTasks(name string) error {
	// If we have access to a task manager, clean up plugin-registered tasks
	// This would require proper interface to task manager
	// For now, just a placeholder for future implementation
	if m.taskManager != nil {
		// Implementation depends on the actual task manager interface
		// The actual implementation would need to know how to unregister plugin tasks
	}

	return nil
}

// GetDependents returns plugins that depend on the specified plugin
func (m *Manager) GetDependents(name string) []string {
	return m.depManager.GetDependents(name)
}

// CanUninstall checks if a plugin can be safely uninstalled (no dependents or force option)
func (m *Manager) CanUninstall(name string) (bool, []string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if plugin exists
	if _, exists := m.plugins[name]; !exists {
		return false, nil, fmt.Errorf("plugin %s not found", name)
	}

	// Get dependents
	dependents := m.depManager.GetDependents(name)

	// Check if plugin is running
	var isRunning bool
	if instance, exists := m.instances[name]; exists {
		isRunning = instance.Status == types.StatusRunning
	}

	return len(dependents) == 0 && !isRunning, dependents, nil
}

// GetPluginInfo returns detailed information about a plugin
func (m *Manager) GetPluginInfo(name string) (*types.PluginInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	info := &types.PluginInfo{
		Name:        plugin.Name(),
		Version:     plugin.Version(),
		Description: plugin.Description(),
		Author:      plugin.Author(),
		Dependencies: plugin.Dependencies(),
	}

	// Get instance status if available
	if instance, exists := m.instances[name]; exists {
		info.Status = instance.Status
		info.LastError = instance.LastError
		info.StartTime = instance.StartTime
		info.StopTime = instance.StopTime
		info.RestartCount = instance.RestartCount
		info.HealthScore = instance.HealthScore
		// 添加监控相关字段
		info.TotalExecutionTime = instance.TotalExecutionTime
		info.TotalExecutions = instance.TotalExecutions
		info.TotalErrors = instance.TotalErrors
		info.LastExecutionTime = instance.LastExecutionTime
	} else {
		info.Status = types.StatusRegistered
	}

	return info, nil
}

// GetSecurityManager returns the security manager
func (m *Manager) GetSecurityManager() *security.SecurityManager {
	return m.securityManager
}

// GetConfigManager returns the configuration manager
func (m *Manager) GetConfigManager() *config.ConfigManager {
	return m.configManager
}

// validatePluginConfig validates plugin configuration against its schema
func (m *Manager) validatePluginConfig(pluginName string, config map[string]interface{}) error {
	// Try to validate the configuration, but don't fail if no schema exists
	if err := m.configManager.ValidateConfig(pluginName, config); err != nil {
		// Check if it's a "schema not found" error
		if fmt.Sprintf("%v", err) == fmt.Sprintf("schema not found for plugin '%s'", pluginName) {
			// It's okay if no schema exists
			return nil
		}
		// Return other validation errors
		return err
	}
	return nil
}

// applyConfigDefaults applies default values from schema to plugin configuration
func (m *Manager) applyConfigDefaults(pluginName string, config map[string]interface{}) error {
	return m.configManager.ApplyDefaults(pluginName, config)
}

// RegisterConfigSchema registers a configuration schema for a plugin
func (m *Manager) RegisterConfigSchema(schema *config.ConfigSchema) error {
	return m.configManager.RegisterSchema(schema)
}

// GetConfigSchema returns the configuration schema for a plugin
func (m *Manager) GetConfigSchema(pluginName string) (*config.ConfigSchema, error) {
	return m.configManager.GetSchema(pluginName)
}

// RegisterConfigTemplate registers a configuration template
func (m *Manager) RegisterConfigTemplate(template *config.ConfigTemplate) error {
	return m.configManager.RegisterTemplate(template)
}

// GetConfigTemplate returns a configuration template
func (m *Manager) GetConfigTemplate(name string) (*config.ConfigTemplate, error) {
	return m.configManager.GetTemplate(name)
}

// ApplyConfigTemplate applies a configuration template to plugin config
func (m *Manager) ApplyConfigTemplate(pluginName, templateName string, config map[string]interface{}) error {
	return m.configManager.ApplyTemplate(pluginName, templateName, config)
}

// ListConfigTemplates lists all available configuration templates
func (m *Manager) ListConfigTemplates() []*config.ConfigTemplate {
	return m.configManager.ListTemplates()
}

// SaveConfigVersion saves a configuration version
func (m *Manager) SaveConfigVersion(pluginName, version, description, author string, config map[string]interface{}) error {
	return m.configManager.SaveVersion(pluginName, version, description, author, config)
}

// GetLatestConfigVersion gets the latest configuration version
func (m *Manager) GetLatestConfigVersion(pluginName string) (map[string]interface{}, error) {
	return m.configManager.GetLatestVersion(pluginName)
}

// RevertToConfigVersion reverts to a specific configuration version
func (m *Manager) RevertToConfigVersion(pluginName, version string) (map[string]interface{}, error) {
	return m.configManager.RevertToVersion(pluginName, version)
}

// ListConfigVersions lists all configuration versions for a plugin
func (m *Manager) ListConfigVersions(pluginName string) ([]*config.ConfigVersion, error) {
	return m.configManager.ListVersions(pluginName)
}