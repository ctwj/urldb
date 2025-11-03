package manager

import (
	"fmt"
	"sync"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// PluginRegistry manages plugin registration
type PluginRegistry struct {
	plugins  map[string]types.Plugin
	mutex    sync.RWMutex
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]types.Plugin),
	}
}

// Register registers a plugin
func (pr *PluginRegistry) Register(plugin types.Plugin) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	name := plugin.Name()
	if _, exists := pr.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	pr.plugins[name] = plugin
	utils.Debug("Plugin registered in registry: %s", name)
	return nil
}

// Unregister unregisters a plugin
func (pr *PluginRegistry) Unregister(name string) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if _, exists := pr.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	delete(pr.plugins, name)
	utils.Debug("Plugin unregistered from registry: %s", name)
	return nil
}

// Get returns a plugin by name
func (pr *PluginRegistry) Get(name string) (types.Plugin, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	plugin, exists := pr.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List returns all registered plugins
func (pr *PluginRegistry) List() []types.Plugin {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	var plugins []types.Plugin
	for _, plugin := range pr.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}