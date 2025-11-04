package manager

import (
	"fmt"

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
	// Check if plugin is already loaded
	if _, exists := pl.manager.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	utils.Info("Plugin loaded: %s", name)
	return nil
}

// LoadPluginWithDependencies loads a plugin and its dependencies
func (pl *PluginLoader) LoadPluginWithDependencies(pluginName string) error {
	// Check if all dependencies are satisfied
	satisfied, unresolved, err := pl.manager.CheckPluginDependencies(pluginName)
	if err != nil {
		return fmt.Errorf("failed to check dependencies for plugin %s: %v", pluginName, err)
	}

	if !satisfied {
		return fmt.Errorf("plugin %s has unsatisfied dependencies: %v", pluginName, unresolved)
	}

	utils.Info("All dependencies satisfied for plugin: %s", pluginName)
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

// LoadAllPlugins loads all registered plugins in the correct dependency order
func (pl *PluginLoader) LoadAllPlugins() error {
	// Get the correct load order based on dependencies
	loadOrder, err := pl.manager.GetLoadOrder()
	if err != nil {
		return fmt.Errorf("failed to determine plugin load order: %v", err)
	}

	utils.Info("Plugin load order: %v", loadOrder)

	// Validate all dependencies before loading
	if err := pl.manager.ValidateDependencies(); err != nil {
		return fmt.Errorf("dependency validation failed: %v", err)
	}

	// Load plugins in the determined order
	for _, pluginName := range loadOrder {
		utils.Info("Loading plugin: %s", pluginName)
		if err := pl.LoadPlugin(pluginName); err != nil {
			utils.Error("Failed to load plugin %s: %v", pluginName, err)
			return err
		}
	}

	utils.Info("All plugins loaded successfully")
	return nil
}