package manager

import (
	"fmt"
	"sync"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// LazyLoader 实现插件懒加载机制
type LazyLoader struct {
	manager     *Manager
	loadedPlugins map[string]bool
	mutex       sync.RWMutex
}

// NewLazyLoader 创建新的懒加载器
func NewLazyLoader(manager *Manager) *LazyLoader {
	return &LazyLoader{
		manager:     manager,
		loadedPlugins: make(map[string]bool),
	}
}

// LoadPluginOnDemand 按需加载插件
func (ll *LazyLoader) LoadPluginOnDemand(name string) (types.Plugin, error) {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()

	// 检查插件是否已经加载
	if ll.loadedPlugins[name] {
		// 从管理器获取已加载的插件
		plugin, err := ll.manager.GetPlugin(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get loaded plugin %s: %v", name, err)
		}
		return plugin, nil
	}

	// 从注册表获取插件信息（不实际加载）
	plugin, err := ll.manager.registry.Get(name)
	if err != nil {
		return nil, fmt.Errorf("plugin %s not found in registry: %v", name, err)
	}

	// 标记插件为已加载
	ll.loadedPlugins[name] = true
	utils.Info("Plugin loaded on demand: %s", name)
	return plugin, nil
}

// UnloadPlugin 卸载插件（释放资源）
func (ll *LazyLoader) UnloadPlugin(name string) error {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()

	// 检查插件是否已加载
	if !ll.loadedPlugins[name] {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	// 停止插件实例（如果正在运行）
	if instance, exists := ll.manager.instances[name]; exists {
		if instance.Status == types.StatusRunning {
			if err := ll.manager.StopPlugin(name); err != nil {
				utils.Warn("Failed to stop plugin %s before unloading: %v", name, err)
			}
		}
	}

	// 从已加载插件列表中移除
	delete(ll.loadedPlugins, name)
	utils.Info("Plugin unloaded: %s", name)
	return nil
}

// IsPluginLoaded 检查插件是否已加载
func (ll *LazyLoader) IsPluginLoaded(name string) bool {
	ll.mutex.RLock()
	defer ll.mutex.RUnlock()
	return ll.loadedPlugins[name]
}

// GetLoadedPlugins 获取所有已加载的插件
func (ll *LazyLoader) GetLoadedPlugins() []string {
	ll.mutex.RLock()
	defer ll.mutex.RUnlock()

	var loaded []string
	for name := range ll.loadedPlugins {
		loaded = append(loaded, name)
	}
	return loaded
}

// PreloadPlugins 预加载指定的插件
func (ll *LazyLoader) PreloadPlugins(pluginNames []string) error {
	for _, name := range pluginNames {
		_, err := ll.LoadPluginOnDemand(name)
		if err != nil {
			return fmt.Errorf("failed to preload plugin %s: %v", name, err)
		}
	}
	return nil
}