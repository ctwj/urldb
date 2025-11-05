package plugin

import (
	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/monitor"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
)

// GlobalManager is the global plugin manager instance
var GlobalManager *manager.SimpleManager

// GlobalPluginMonitor is the global plugin monitor instance
var GlobalPluginMonitor *monitor.PluginMonitor

// InitPluginSystem initializes the plugin system
func InitPluginSystem(taskManager *task.TaskManager, repoManager *repo.RepositoryManager) {
	// Initialize logger if not already initialized
	if utils.GetLogger() == nil {
		utils.InitLogger(nil)
	}

	// Create plugin manager with database and repo manager
	// 创建插件监控器
	GlobalPluginMonitor = monitor.NewPluginMonitor()
	GlobalManager = manager.NewSimpleManager(taskManager, repoManager, db.DB, GlobalPluginMonitor)

	// Load all plugins from filesystem
	if err := GlobalManager.LoadAllPluginsFromFilesystem(); err != nil {
		utils.Error("加载文件系统插件失败: %v", err)
	}

	// Log initialization
	utils.Info("Plugin system initialized")
}

// GetPluginMonitor returns the global plugin monitor
func GetPluginMonitor() *monitor.PluginMonitor {
	return GlobalPluginMonitor
}

// RegisterPlugin registers a plugin with the global manager
func RegisterPlugin(plugin types.Plugin) error {
	if GlobalManager == nil {
		return nil // Plugin system not initialized
	}

	return GlobalManager.RegisterPlugin(plugin)
}

// GetManager returns the global plugin manager
func GetManager() *manager.SimpleManager {
	return GlobalManager
}