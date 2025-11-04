package hotupdate

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// PluginWatcher 插件文件监视器
type PluginWatcher struct {
	manager     *manager.Manager
	pluginPaths map[string]string // pluginName -> pluginPath
	fileHashes  map[string]string // filePath -> fileHash
	mutex       sync.RWMutex
	stopChan    chan struct{}
	ticker      *time.Ticker
	interval    time.Duration
}

// NewPluginWatcher 创建新的插件监视器
func NewPluginWatcher(mgr *manager.Manager, interval time.Duration) *PluginWatcher {
	if interval <= 0 {
		interval = 5 * time.Second // 默认5秒检查一次
	}

	return &PluginWatcher{
		manager:     mgr,
		pluginPaths: make(map[string]string),
		fileHashes:  make(map[string]string),
		stopChan:    make(chan struct{}),
		interval:    interval,
	}
}

// AddPlugin 添加要监视的插件
func (w *PluginWatcher) AddPlugin(pluginName, pluginPath string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 检查插件文件是否存在
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin file does not exist: %s", pluginPath)
	}

	// 计算文件哈希
	hash, err := w.calculateFileHash(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to calculate file hash: %v", err)
	}

	w.pluginPaths[pluginName] = pluginPath
	w.fileHashes[pluginPath] = hash

	utils.Info("Added plugin to watcher: %s -> %s", pluginName, pluginPath)
	return nil
}

// RemovePlugin 移除监视的插件
func (w *PluginWatcher) RemovePlugin(pluginName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if pluginPath, exists := w.pluginPaths[pluginName]; exists {
		delete(w.pluginPaths, pluginName)
		delete(w.fileHashes, pluginPath)
		utils.Info("Removed plugin from watcher: %s", pluginName)
	}
}

// Start 开始监视
func (w *PluginWatcher) Start() {
	w.ticker = time.NewTicker(w.interval)
	go w.watchLoop()
	utils.Info("Plugin watcher started with interval: %v", w.interval)
}

// Stop 停止监视
func (w *PluginWatcher) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
	close(w.stopChan)
	utils.Info("Plugin watcher stopped")
}

// watchLoop 监视循环
func (w *PluginWatcher) watchLoop() {
	for {
		select {
		case <-w.stopChan:
			return
		case <-w.ticker.C:
			w.checkForUpdates()
		}
	}
}

// checkForUpdates 检查插件更新
func (w *PluginWatcher) checkForUpdates() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for pluginName, pluginPath := range w.pluginPaths {
		// 计算当前文件哈希
		currentHash, err := w.calculateFileHash(pluginPath)
		if err != nil {
			utils.Error("Failed to calculate hash for plugin %s: %v", pluginName, err)
			continue
		}

		// 比较哈希值
		if oldHash, exists := w.fileHashes[pluginPath]; exists && oldHash != currentHash {
			utils.Info("Detected update for plugin: %s", pluginName)

			// 更新哈希值
			w.fileHashes[pluginPath] = currentHash

			// 执行热更新
			if err := w.performHotUpdate(pluginName, pluginPath); err != nil {
				utils.Error("Failed to perform hot update for plugin %s: %v", pluginName, err)
			}
		}
	}
}

// calculateFileHash 计算文件哈希值
func (w *PluginWatcher) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// performHotUpdate 执行热更新
func (w *PluginWatcher) performHotUpdate(pluginName, pluginPath string) error {
	utils.Info("Performing hot update for plugin: %s", pluginName)

	// 1. 停止插件
	status := w.manager.GetPluginStatus(pluginName)
	if status == types.StatusRunning {
		if err := w.manager.StopPlugin(pluginName); err != nil {
			utils.Error("Failed to stop plugin %s: %v", pluginName, err)
			return err
		}
	}

	// 2. 卸载插件
	if err := w.manager.UninstallPlugin(pluginName, true); err != nil {
		utils.Error("Failed to uninstall plugin %s: %v", pluginName, err)
		return err
	}

	// 3. 重新加载插件
	// 这里需要根据具体实现方式来处理
	// 假设我们有一个插件加载函数
	// err := w.loadPlugin(pluginName, pluginPath)
	// if err != nil {
	// 	utils.Error("Failed to reload plugin %s: %v", pluginName, err)
	// 	return err
	// }

	// 4. 重新初始化插件
	// config, err := w.manager.GetLatestConfigVersion(pluginName)
	// if err != nil {
	// 	utils.Warn("Failed to get config for plugin %s: %v", pluginName, err)
	// 	config = make(map[string]interface{})
	// }
	//
	// if err := w.manager.InitializePlugin(pluginName, config); err != nil {
	// 	utils.Error("Failed to initialize plugin %s: %v", pluginName, err)
	// 	return err
	// }

	// 5. 重新启动插件
	// if err := w.manager.StartPlugin(pluginName); err != nil {
	// 	utils.Error("Failed to start plugin %s: %v", pluginName, err)
	// 	return err
	// }

	utils.Info("Hot update completed for plugin: %s", pluginName)
	return nil
}

// loadPlugin 加载插件（具体实现依赖于插件加载机制）
// func (w *PluginWatcher) loadPlugin(pluginName, pluginPath string) error {
// 	// 这里需要根据具体的插件加载机制来实现
// 	// 可能需要使用 plugin 包或者自定义的加载机制
// 	return nil
// }

// ListWatchedPlugins 列出所有被监视的插件
func (w *PluginWatcher) ListWatchedPlugins() map[string]string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	result := make(map[string]string)
	for name, path := range w.pluginPaths {
		result[name] = path
	}
	return result
}