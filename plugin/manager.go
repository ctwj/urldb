package plugin

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin/jsvm"
	"github.com/ctwj/urldb/utils"
)

// Manager 插件管理器
type Manager struct {
	app          core.App
	installer    *PluginInstaller
	jsvmConfig   jsvm.Config
	repoManager  *repo.RepositoryManager
	loadedPlugins map[string]bool
	mu           sync.RWMutex
}

// NewManager 创建插件管理器
func NewManager(app core.App) *Manager {
	return &Manager{
		app:           app,
		installer:     NewPluginInstaller("."),
		loadedPlugins: make(map[string]bool),
	}
}

// SetRepoManager 设置 RepositoryManager
func (m *Manager) SetRepoManager(repoManager *repo.RepositoryManager) {
	m.repoManager = repoManager
}

// RegisterJSVM 注册 JavaScript 虚拟机插件
func (m *Manager) RegisterJSVM(config jsvm.Config) error {
	m.jsvmConfig = config
	return jsvm.Register(m.app, config)
}

// RegisterJSVMWithRepo 注册 JavaScript 虚拟机插件（带RepositoryManager）
func (m *Manager) RegisterJSVMWithRepo(config jsvm.Config, repoManager *repo.RepositoryManager) error {
	m.jsvmConfig = config
	m.repoManager = repoManager
	return jsvm.RegisterWithRepo(m.app, config, repoManager)
}

// RegisterJSVMDefault 注册默认配置的 JSVM 插件
func (m *Manager) RegisterJSVMDefault() error {
	config := jsvm.Config{
		HooksWatch:      true,
		HooksPoolSize:   10,
		OnInit:          m.defaultOnInit,
		RouteRegister:   m.registerPluginRoute,
	}
	return m.RegisterJSVM(config)
}

// defaultOnInit 默认的 VM 初始化回调
func (m *Manager) defaultOnInit(vm *goja.Runtime) {
	// 可以在这里添加自定义的全局变量或函数
	// 例如：vm.Set("version", "1.0.0")
}

// InstallPlugin 安装插件
func (m *Manager) InstallPlugin(source string) error {
	// 判断是URL
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return m.installer.InstallFromURL(source)
	}

	// 判断是ZIP文件
	if len(source) > 4 && source[len(source)-4:] == ".zip" {
		return m.installer.InstallFromFile(source)
	}

	// 处理单文件插件 (.plugin.js)
	if len(source) > 10 && source[len(source)-10:] == ".plugin.js" {
		return m.installer.InstallSingleFile(source)
	}

	return fmt.Errorf("unsupported plugin source: %s (must be .zip, .plugin.js file or URL)", source)
}

// UninstallPlugin 卸载插件
func (m *Manager) UninstallPlugin(pluginName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查插件是否已加载
	if m.loadedPlugins[pluginName] {
		return fmt.Errorf("cannot uninstall plugin '%s' while it is loaded", pluginName)
	}

	return m.installer.Uninstall(pluginName)
}

// LoadPlugin 加载已安装的插件
func (m *Manager) LoadPlugin(pluginName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.loadedPlugins[pluginName] {
		return fmt.Errorf("plugin '%s' is already loaded", pluginName)
	}

	// 获取插件安装目录
	pluginDir := filepath.Join(m.installer.installedDir, pluginName)

	// 验证插件钩子目录
	pluginHooksDir := filepath.Join(pluginDir, "hooks")
	if _, err := os.Stat(pluginHooksDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin hooks directory not found: %s", pluginHooksDir)
	}

	// 使用增强的多目录扫描方式加载插件
	if err := m.loadPluginWithMultiDir(pluginName, pluginHooksDir); err != nil {
		return fmt.Errorf("failed to load plugin '%s': %w", pluginName, err)
	}

	// 标记插件为已加载
	m.loadedPlugins[pluginName] = true

	utils.Info("Plugin '%s' loaded successfully", pluginName)
	return nil
}

// loadPluginWithMultiDir 使用多目录扫描方式加载插件
func (m *Manager) loadPluginWithMultiDir(pluginName, pluginHooksDir string) error {
	// 获取所有需要扫描的目录
	allHookDirs := []string{}

	// 添加原始 hooks 目录
	if m.jsvmConfig.HooksDir != "" {
		allHookDirs = append(allHookDirs, m.jsvmConfig.HooksDir)
	}

	// 添加已加载插件的所有钩子目录
	for loadedPluginName := range m.loadedPlugins {
		loadedPluginDir := filepath.Join(m.installer.installedDir, loadedPluginName)
		loadedHooksDir := filepath.Join(loadedPluginDir, "hooks")
		if _, err := os.Stat(loadedHooksDir); err == nil {
			allHookDirs = append(allHookDirs, loadedHooksDir)
		}
	}

	// 添加当前插件的钩子目录
	allHookDirs = append(allHookDirs, pluginHooksDir)

	// 创建合并的临时目录用于多目录扫描
	tempDir, err := os.MkdirTemp("", "hooks-merge-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 复制所有钩子文件到临时目录
	for i, hookDir := range allHookDirs {
		if err := m.copyHookFiles(hookDir, tempDir, i); err != nil {
			utils.Error("Failed to copy hook files from %s: %v", hookDir, err)
			continue
		}
	}

	// 临时修改配置使用合并目录
	originalHooksDir := m.jsvmConfig.HooksDir
	m.jsvmConfig.HooksDir = tempDir

	// 检查 app 是否为 nil
	if m.app == nil {
		return fmt.Errorf("plugin manager app instance is nil")
	}

	// 重新注册 JSVM
	if err := jsvm.RegisterWithRepo(m.app, m.jsvmConfig, m.repoManager); err != nil {
		m.jsvmConfig.HooksDir = originalHooksDir
		return err
	}

	// 恢复原始配置
	m.jsvmConfig.HooksDir = originalHooksDir

	return nil
}

// copyHookFiles 复制钩子文件到目标目录
func (m *Manager) copyHookFiles(srcDir string, destDir string, index int) error {
	// 扫描源目录中的所有 .plugin.js 文件
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 只复制 .plugin.js 文件
		if !strings.HasSuffix(file.Name(), ".plugin.js") {
			continue
		}

		srcPath := filepath.Join(srcDir, file.Name())

		// 为避免文件名冲突，添加索引前缀
		destName := fmt.Sprintf("%d_%s", index, file.Name())
		destPath := filepath.Join(destDir, destName)

		// 复制文件
		if err := m.copyFile(srcPath, destPath); err != nil {
			return err
		}
	}

	return nil
}

// copyFile 复制文件
func (m *Manager) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// UnloadPlugin 卸载已加载的插件
func (m *Manager) UnloadPlugin(pluginName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.loadedPlugins[pluginName] {
		return fmt.Errorf("plugin '%s' is not loaded", pluginName)
	}

	// 这里可以实现插件的卸载逻辑
	// 由于当前的 JSVM 架构限制，完全卸载比较复杂
	// 暂时只是标记为未加载
	delete(m.loadedPlugins, pluginName)

	utils.Info("Plugin '%s' unloaded", pluginName)
	return nil
}

// ReloadPlugin 重新加载插件
func (m *Manager) ReloadPlugin(pluginName string) error {
	if err := m.UnloadPlugin(pluginName); err != nil {
		return err
	}
	return m.LoadPlugin(pluginName)
}

// ListInstalledPlugins 列出已安装的插件
func (m *Manager) ListInstalledPlugins() ([]*PluginPackage, error) {
	return m.installer.ListInstalled()
}

// IsPluginInstalled 检查插件是否已安装
func (m *Manager) IsPluginInstalled(pluginName string) bool {
	return m.installer.IsInstalled(pluginName)
}

// IsPluginLoaded 检查插件是否已加载
func (m *Manager) IsPluginLoaded(pluginName string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loadedPlugins[pluginName]
}

// LoadAllInstalledPlugins 加载所有已安装的插件
func (m *Manager) LoadAllInstalledPlugins() error {
	plugins, err := m.installer.ListInstalled()
	if err != nil {
		return err
	}

	for _, pkg := range plugins {
		if err := m.LoadPlugin(pkg.Name); err != nil {
			utils.Error("Failed to load plugin '%s': %v", pkg.Name, err)
		}
	}

	return nil
}

// registerPluginRoute 注册插件路由（用于 JSVM 回调）
func (m *Manager) registerPluginRoute(method, path string, handler func() (interface{}, error)) error {
	// 这里可以实现动态路由注册
	// 由于架构限制，暂时只记录日志
	utils.Info("Plugin route registered: %s %s", method, path)
	return nil
}