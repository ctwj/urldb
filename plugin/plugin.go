package plugin

import (
	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
)

// GlobalManager is the global plugin manager instance
var GlobalManager *manager.Manager

// InitPluginSystem initializes the plugin system
func InitPluginSystem(taskManager *task.TaskManager) {
	// Initialize logger if not already initialized
	if utils.GetLogger() == nil {
		utils.InitLogger(nil)
	}

	// Create plugin manager
	GlobalManager = manager.NewManager(taskManager)

	// Log initialization
	utils.Info("Plugin system initialized")
}

// RegisterPlugin registers a plugin with the global manager
func RegisterPlugin(plugin types.Plugin) error {
	if GlobalManager == nil {
		return nil // Plugin system not initialized
	}

	return GlobalManager.RegisterPlugin(plugin)
}

// GetManager returns the global plugin manager
func GetManager() *manager.Manager {
	return GlobalManager
}