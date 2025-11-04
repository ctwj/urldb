package hotupdate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// PluginUpdater 插件更新器
type PluginUpdater struct {
	manager   *manager.Manager
	watcher   *PluginWatcher
	pluginDir string
}

// NewPluginUpdater 创建新的插件更新器
func NewPluginUpdater(mgr *manager.Manager, pluginDir string) *PluginUpdater {
	return &PluginUpdater{
		manager:   mgr,
		pluginDir: pluginDir,
	}
}

// StartUpdaterWithWatcher 启动更新器并设置文件监视
func (u *PluginUpdater) StartUpdaterWithWatcher(interval time.Duration) error {
	// 创建监视器
	u.watcher = NewPluginWatcher(u.manager, interval)
	u.watcher.Start()

	// 扫描并添加所有插件文件到监视器
	if err := u.scanAndAddPlugins(); err != nil {
		return fmt.Errorf("failed to scan plugins: %v", err)
	}

	utils.Info("Plugin updater with watcher started")
	return nil
}

// scanAndAddPlugins 扫描并添加插件
func (u *PluginUpdater) scanAndAddPlugins() error {
	if u.watcher == nil {
		return fmt.Errorf("watcher not initialized")
	}

	// 遍历插件目录
	err := filepath.Walk(u.pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为插件文件（通常是 .so 文件）
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".so") {
			// 从文件名提取插件名称
			pluginName := strings.TrimSuffix(filepath.Base(path), ".so")

			// 如果插件已注册，添加到监视器
			_, err := u.manager.GetPlugin(pluginName)
			if err == nil { // 插件已存在
				if err := u.watcher.AddPlugin(pluginName, path); err != nil {
					utils.Warn("Failed to add plugin to watcher %s: %v", pluginName, err)
				}
			}
		}

		return nil
	})

	return err
}

// UpdatePlugin 手动更新插件
func (u *PluginUpdater) UpdatePlugin(pluginName string, pluginPath string) error {
	utils.Info("Manually updating plugin: %s", pluginName)

	// 检查插件是否存在
	_, err := u.manager.GetPlugin(pluginName)
	if err != nil {
		return fmt.Errorf("plugin %s not found: %v", pluginName, err)
	}

	// 检查插件文件是否存在
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin file does not exist: %s", pluginPath)
	}

	// 停止插件
	status := u.manager.GetPluginStatus(pluginName)
	if status == types.StatusRunning {
		if err := u.manager.StopPlugin(pluginName); err != nil {
			utils.Error("Failed to stop plugin %s: %v", pluginName, err)
			return err
		}
	}

	// 备份当前插件文件
	pluginDir := filepath.Dir(pluginPath)
	backupPath := filepath.Join(pluginDir, fmt.Sprintf("%s.backup.%d.so", pluginName, time.Now().Unix()))
	if err := copyFile(pluginPath, backupPath); err != nil {
		utils.Warn("Failed to backup plugin file: %v", err)
	}

	// 移除旧插件
	if err := u.manager.UninstallPlugin(pluginName, true); err != nil {
		utils.Error("Failed to uninstall plugin %s: %v", pluginName, err)
		return err
	}

	// 更新插件文件（模拟更新过程）
	// 实际实现中可能需要从远程下载或从其他位置复制新版本
	if err := copyFile(pluginPath, pluginPath); err != nil {
		utils.Error("Failed to update plugin file: %v", err)
		// 尝试恢复备份
		if _, err := os.Stat(backupPath); err == nil {
			os.Rename(backupPath, pluginPath)
		}
		return err
	}

	// 加载新插件
	// 这里需要根据具体的插件加载机制来实现
	// newPlugin, err := u.loadPlugin(pluginPath)
	// if err != nil {
	// 	utils.Error("Failed to load new plugin: %v", err)
	// 	return err
	// }

	// 注册新插件
	// if err := u.manager.RegisterPlugin(newPlugin); err != nil {
	// 	utils.Error("Failed to register new plugin: %v", err)
	// 	return err
	// }

	// 重新初始化插件
	config, err := u.manager.GetLatestConfigVersion(pluginName)
	if err != nil {
		utils.Warn("Failed to get config for plugin %s: %v", pluginName, err)
		config = make(map[string]interface{})
	}

	if err := u.manager.InitializePlugin(pluginName, config); err != nil {
		utils.Error("Failed to initialize plugin %s: %v", pluginName, err)
		return err
	}

	// 重新启动插件
	if err := u.manager.StartPlugin(pluginName); err != nil {
		utils.Error("Failed to start plugin %s: %v", pluginName, err)
		return err
	}

	// 更新监视器中的插件路径
	if u.watcher != nil {
		u.watcher.RemovePlugin(pluginName)
		if err := u.watcher.AddPlugin(pluginName, pluginPath); err != nil {
			utils.Warn("Failed to add updated plugin to watcher: %v", err)
		}
	}

	// 清理备份文件
	if _, err := os.Stat(backupPath); err == nil {
		os.Remove(backupPath)
	}

	utils.Info("Plugin updated successfully: %s", pluginName)
	return nil
}

// UpdatePluginFromURL 从URL更新插件
func (u *PluginUpdater) UpdatePluginFromURL(pluginName, url string) error {
	// 这里需要实现从URL下载插件的逻辑
	// 简化实现，返回错误表示未实现
	utils.Info("Updating plugin %s from URL: %s", pluginName, url)
	return fmt.Errorf("not implemented: UpdatePluginFromURL")
}

// GetWatcher 返回监视器
func (u *PluginUpdater) GetWatcher() *PluginWatcher {
	return u.watcher
}

// Stop 停止更新器
func (u *PluginUpdater) Stop() {
	if u.watcher != nil {
		u.watcher.Stop()
	}
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// loadPlugin 加载插件（简化实现，实际需要根据具体插件加载机制）
func (u *PluginUpdater) loadPlugin(pluginPath string) (types.Plugin, error) {
	// 这里需要实现实际的插件加载逻辑
	// 可能使用 Go 的 plugin 包或其他自定义机制
	return nil, fmt.Errorf("not implemented: loadPlugin")
}