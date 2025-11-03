package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// Manager is the plugin manager that handles plugin lifecycle
type Manager struct {
	plugins    map[string]types.Plugin
	instances  map[string]*PluginInstance
	registry   *PluginRegistry
	mutex      sync.RWMutex
	taskManager interface{} // Reference to the existing task manager
}

// PluginInstance represents a running plugin instance
type PluginInstance struct {
	Plugin     types.Plugin
	Context    types.PluginContext
	Status     types.PluginStatus
	Config     map[string]interface{}
	StartTime  time.Time
	StopTime   *time.Time
	LastError  error
	RestartCount int
	HealthScore  float64
}

// NewManager creates a new plugin manager
func NewManager(taskManager interface{}) *Manager {
	return &Manager{
		plugins:    make(map[string]types.Plugin),
		instances:  make(map[string]*PluginInstance),
		registry:   NewPluginRegistry(),
		taskManager: taskManager,
	}
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
	utils.Info("Plugin registered: %s", name)
	return nil
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

	// Create plugin context
	context := NewPluginContext(name, m, config)

	// Create plugin instance
	instance := &PluginInstance{
		Plugin:  plugin,
		Context: context,
		Status:  types.StatusInitialized,
		Config:  config,
	}

	// Initialize the plugin
	if err := plugin.Initialize(context); err != nil {
		instance.Status = types.StatusError
		instance.LastError = err
		m.instances[name] = instance
		return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
	}

	m.instances[name] = instance
	utils.Info("Plugin initialized: %s", name)
	return nil
}

// StartPlugin starts a plugin
func (m *Manager) StartPlugin(name string) error {
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
	if err := instance.Plugin.Start(); err != nil {
		instance.Status = types.StatusError
		instance.LastError = err
		return fmt.Errorf("failed to start plugin %s: %v", name, err)
	}

	instance.Status = types.StatusRunning
	instance.StartTime = time.Now()
	instance.RestartCount++
	utils.Info("Plugin started: %s", name)
	return nil
}

// StopPlugin stops a plugin
func (m *Manager) StopPlugin(name string) error {
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
	if err := instance.Plugin.Stop(); err != nil {
		instance.Status = types.StatusError
		instance.LastError = err
		return fmt.Errorf("failed to stop plugin %s: %v", name, err)
	}

	now := time.Now()
	instance.Status = types.StatusStopped
	instance.StopTime = &now
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