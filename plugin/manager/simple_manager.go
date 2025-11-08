package manager

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin/loader"
	"github.com/ctwj/urldb/plugin/monitor"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// SimpleManager 简化的插件管理器，使用.so文件直接加载
type SimpleManager struct {
	plugins        map[string]types.Plugin
	instances      map[string]*types.PluginInstance
	mutex          sync.RWMutex
	taskManager    *task.TaskManager
	repoManager    *repo.RepositoryManager
	db             *gorm.DB
	pluginMonitor  *monitor.PluginMonitor
	pluginLoader   *loader.SimplePluginLoader
}

// NewSimpleManager 创建简化版插件管理器
func NewSimpleManager(taskManager *task.TaskManager, repoManager *repo.RepositoryManager, database *gorm.DB, pluginMonitor *monitor.PluginMonitor) *SimpleManager {
	manager := &SimpleManager{
		plugins:       make(map[string]types.Plugin),
		instances:     make(map[string]*types.PluginInstance),
		taskManager:   taskManager,
		repoManager:   repoManager,
		db:            database,
		pluginMonitor: pluginMonitor,
		pluginLoader:  loader.NewSimplePluginLoader("./plugins"),
	}

	return manager
}

// RegisterPlugin 注册插件
func (sm *SimpleManager) RegisterPlugin(plugin types.Plugin) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	name := plugin.Name()
	if _, exists := sm.plugins[name]; exists {
		utils.Info("插件 %s 已存在，更新插件定义", name)
	}

	sm.plugins[name] = plugin
	utils.Info("成功注册插件: %s (版本: %s)", name, plugin.Version())
	return nil
}

// LoadPluginFromFile 从文件加载插件
func (sm *SimpleManager) LoadPluginFromFile(filepath string) error {
	plugin, err := sm.pluginLoader.LoadPlugin(filepath)
	if err != nil {
		return fmt.Errorf("加载插件失败: %v", err)
	}

	name := plugin.Name()
	if name == "" {
		return fmt.Errorf("插件名称不能为空")
	}

	sm.mutex.Lock()
	sm.plugins[name] = plugin
	sm.mutex.Unlock()

	utils.Info("成功从文件加载并注册插件: %s (版本: %s)", name, plugin.Version())
	return nil
}

// InitializePlugin 初始化插件
func (sm *SimpleManager) InitializePlugin(name string) error {
	sm.mutex.RLock()
	plugin, exists := sm.plugins[name]
	sm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("插件不存在: %s", name)
	}

	// 创建插件实例
	instance := &types.PluginInstance{
		Plugin: plugin,
		Status: types.StatusInitialized,
	}

	sm.mutex.Lock()
	sm.instances[name] = instance
	sm.mutex.Unlock()

	utils.Info("插件 %s 初始化完成", name)
	return nil
}

// StartPlugin 启动插件
func (sm *SimpleManager) StartPlugin(name string) error {
	sm.mutex.RLock()
	instance, exists := sm.instances[name]
	sm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("插件未初始化: %s", name)
	}

	if instance.Status == types.StatusRunning {
		return fmt.Errorf("插件已在运行状态: %s", name)
	}

	// 启动插件
	err := instance.Plugin.Start()
	if err != nil {
		return fmt.Errorf("启动插件失败: %v", err)
	}

	instance.Status = types.StatusRunning
	instance.StartTime = time.Now()
	utils.Info("插件 %s 启动成功", name)
	return nil
}

// StopPlugin 停止插件
func (sm *SimpleManager) StopPlugin(name string) error {
	sm.mutex.RLock()
	instance, exists := sm.instances[name]
	sm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("插件不存在: %s", name)
	}

	if instance.Status != types.StatusRunning {
		return fmt.Errorf("插件未在运行状态: %s", name)
	}

	// 停止插件
	err := instance.Plugin.Stop()
	if err != nil {
		return fmt.Errorf("停止插件失败: %v", err)
	}

	instance.Status = types.StatusStopped
	stopTime := time.Now()
	instance.StopTime = &stopTime
	utils.Info("插件 %s 停止成功", name)
	return nil
}

// UninstallPlugin 卸载插件
func (sm *SimpleManager) UninstallPlugin(name string, force bool) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// 先停止插件
	if instance, exists := sm.instances[name]; exists && instance.Status == types.StatusRunning {
		if err := instance.Plugin.Stop(); err != nil {
			if !force {
				return fmt.Errorf("停止插件失败: %v", err)
			}
			utils.Warn("强制停止插件: %s", name)
		}
		instance.Status = types.StatusStopped
		stopTime := time.Now()
		instance.StopTime = &stopTime
	}

	// 从管理器中移除插件
	delete(sm.plugins, name)
	delete(sm.instances, name)

	// 尝试从文件系统卸载
	if err := sm.pluginLoader.UninstallPlugin(name); err != nil {
		if !force {
			return fmt.Errorf("卸载插件文件失败: %v", err)
		}
		utils.Warn("插件文件可能已不存在: %s", name)
	}

	utils.Info("插件 %s 卸载完成", name)
	return nil
}

// GetPlugin 获取插件
func (sm *SimpleManager) GetPlugin(name string) (types.Plugin, error) {
	sm.mutex.RLock()
	plugin, exists := sm.plugins[name]
	sm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("插件不存在: %s", name)
	}

	return plugin, nil
}

// GetPluginInstance 获取插件实例
func (sm *SimpleManager) GetPluginInstance(name string) (*types.PluginInstance, error) {
	sm.mutex.RLock()
	instance, exists := sm.instances[name]
	sm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("插件实例不存在: %s", name)
	}

	return instance, nil
}

// GetPluginInfo 获取插件信息
func (sm *SimpleManager) GetPluginInfo(name string) (*types.PluginInfo, error) {
	sm.mutex.RLock()
	plugin, exists := sm.plugins[name]
	instance, instanceExists := sm.instances[name]
	sm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("插件不存在: %s", name)
	}

	info := &types.PluginInfo{
		Name:        plugin.Name(),
		Version:     plugin.Version(),
		Description: plugin.Description(),
		Author:      plugin.Author(),
	}

	if instanceExists {
		info.Status = instance.Status
	}

	return info, nil
}

// ListPlugins 列出所有插件
func (sm *SimpleManager) ListPlugins() []types.PluginInfo {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var plugins []types.PluginInfo
	for name, plugin := range sm.plugins {
		instance, exists := sm.instances[name]
		status := types.StatusRegistered
		if exists {
			status = instance.Status
		}

		info := types.PluginInfo{
			Name:        plugin.Name(),
			Version:     plugin.Version(),
			Description: plugin.Description(),
			Author:      plugin.Author(),
			Status:      status,
		}

		// 尝试获取更多统计信息（如果插件支持监控）
		if sm.pluginMonitor != nil {
			// 简化处理，不直接调用不存在的方法
			// 这里可以留空或添加其他监控逻辑
		}

		plugins = append(plugins, info)
	}

	return plugins
}

// GetEnabledPlugins 获取所有启用的插件
func (sm *SimpleManager) GetEnabledPlugins() []types.Plugin {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var plugins []types.Plugin
	for _, plugin := range sm.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

// InstallPluginFromFile 从文件安装插件
func (sm *SimpleManager) InstallPluginFromFile(pluginFilepath string) error {
	data, err := ioutil.ReadFile(pluginFilepath)
	if err != nil {
		return fmt.Errorf("failed to read plugin file: %v", err)
	}

	filename := filepath.Base(pluginFilepath)
	// 提取不带扩展名的文件名作为插件名
	pluginName := strings.TrimSuffix(filename, filepath.Ext(filename))

	return sm.pluginLoader.InstallPluginFromBytes(pluginName, data)
}

// LoadAllPluginsFromFilesystem 从文件系统加载所有.so文件
func (sm *SimpleManager) LoadAllPluginsFromFilesystem() error {
	plugins, err := sm.pluginLoader.LoadAllPlugins()
	if err != nil {
		return fmt.Errorf("加载插件失败: %v", err)
	}

	for _, plugin := range plugins {
		name := plugin.Name()
		if name == "" {
			utils.Error("发现插件名称为空，跳过: %v", plugin)
			continue
		}

		sm.mutex.Lock()
		sm.plugins[name] = plugin
		sm.mutex.Unlock()

		utils.Info("从文件系统加载插件: %s (版本: %s)", name, plugin.Version())
	}

	return nil
}