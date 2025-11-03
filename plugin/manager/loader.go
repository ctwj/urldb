package manager

import (
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// PluginLoader handles loading and discovering plugins
type PluginLoader struct {
	manager *Manager
}

// NewPluginLoader creates a new plugin loader
func NewPluginLoader(manager *Manager) *PluginLoader {
	return &PluginLoader{
		manager: manager,
	}
}

// LoadPlugin loads a plugin by name
func (pl *PluginLoader) LoadPlugin(name string) error {
	// In a real implementation, this would load a plugin from a registry or configuration
	// For now, we'll just log the operation
	utils.Debug("Loading plugin: %s", name)
	return nil
}

// DiscoverPlugins discovers and registers all available plugins
func (pl *PluginLoader) DiscoverPlugins() error {
	// In a real implementation, this would discover plugins from a directory or configuration
	// For now, we'll just log the operation
	utils.Info("Discovering plugins...")
	return nil
}

// AutoRegisterPlugin automatically registers a plugin
func (pl *PluginLoader) AutoRegisterPlugin(plugin types.Plugin) error {
	return pl.manager.RegisterPlugin(plugin)
}

// LoadAllPlugins loads all registered plugins
func (pl *PluginLoader) LoadAllPlugins() error {
	// In a real implementation, this would load all plugins based on configuration
	// For now, we'll just log the operation
	utils.Info("Loading all plugins...")
	return nil
}