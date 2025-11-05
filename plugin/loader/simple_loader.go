package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// SimplePluginLoader 简单插件加载器，直接加载.so文件
type SimplePluginLoader struct {
	pluginDir string
}

// NewSimplePluginLoader 创建简单插件加载器
func NewSimplePluginLoader(pluginDir string) *SimplePluginLoader {
	return &SimplePluginLoader{
		pluginDir: pluginDir,
	}
}

// LoadPlugin 加载单个插件文件
func (l *SimplePluginLoader) LoadPlugin(filename string) (types.Plugin, error) {
	if !l.isPluginFile(filename) {
		return nil, fmt.Errorf("not a plugin file: %s", filename)
	}

	pluginPath := filepath.Join(l.pluginDir, filename)

	// 打开插件文件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %v", filename, err)
	}

	// 查找插件符号
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin symbol not found in %s: %v", filename, err)
	}

	// 类型断言
	pluginInstance, ok := sym.(types.Plugin)
	if !ok {
		return nil, fmt.Errorf("invalid plugin type in %s", filename)
	}

	utils.Info("成功加载插件: %s (名称: %s, 版本: %s)",
		filename, pluginInstance.Name(), pluginInstance.Version())
	return pluginInstance, nil
}

// LoadAllPlugins 加载所有插件
func (l *SimplePluginLoader) LoadAllPlugins() ([]types.Plugin, error) {
	var plugins []types.Plugin

	// 确保插件目录存在
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// 读取插件目录
	files, err := ioutil.ReadDir(l.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if l.isPluginFile(file.Name()) {
			pluginInstance, err := l.LoadPlugin(file.Name())
			if err != nil {
				utils.Error("加载插件 %s 失败: %v", file.Name(), err)
				continue
			}
			plugins = append(plugins, pluginInstance)
		}
	}

	utils.Info("总共加载了 %d 个插件", len(plugins))
	return plugins, nil
}

// isPluginFile 检查是否为插件文件
func (l *SimplePluginLoader) isPluginFile(filename string) bool {
	// 根据平台检查文件扩展名
	return strings.HasSuffix(filename, l.getPluginExtension())
}

// getPluginExtension 获取当前平台的插件文件扩展名
func (l *SimplePluginLoader) getPluginExtension() string {
	switch runtime.GOOS {
	case "linux":
		return ".so"
	case "windows":
		return ".dll"
	case "darwin":
		return ".dylib"
	default:
		return ""
	}
}

// GetPluginInfo 获取插件信息（通过调用插件方法）
func (l *SimplePluginLoader) GetPluginInfo(filename string) (map[string]interface{}, error) {
	pluginInstance, err := l.LoadPlugin(filename)
	if err != nil {
		return nil, err
	}

	info := make(map[string]interface{})
	info["name"] = pluginInstance.Name()
	info["version"] = pluginInstance.Version()
	info["description"] = pluginInstance.Description()
	info["author"] = pluginInstance.Author()
	info["filename"] = filename

	// 如果插件实现了ConfigurablePlugin接口，获取配置schema
	if configurablePlugin, ok := pluginInstance.(interface{ ConfigSchema() map[string]interface{} }); ok {
		info["config_schema"] = configurablePlugin.ConfigSchema()
	}

	return info, nil
}

// ListPluginFiles 列出所有插件文件及其信息
func (l *SimplePluginLoader) ListPluginFiles() ([]map[string]interface{}, error) {
	var pluginFiles []map[string]interface{}

	// 确保插件目录存在
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(l.pluginDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && l.isPluginFile(file.Name()) {
			info, err := l.GetPluginInfo(file.Name())
			if err != nil {
				utils.Error("获取插件信息失败: %v", err)
				continue
			}
			pluginFiles = append(pluginFiles, info)
		}
	}

	return pluginFiles, nil
}

// InstallPluginFromBytes 从字节数据安装插件
func (l *SimplePluginLoader) InstallPluginFromBytes(pluginName string, data []byte) error {
	// 确保插件目录存在
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return err
	}

	// 创建插件文件
	pluginPath := filepath.Join(l.pluginDir, pluginName+l.getPluginExtension())

	return ioutil.WriteFile(pluginPath, data, 0755)
}

// UninstallPlugin 卸载插件
func (l *SimplePluginLoader) UninstallPlugin(pluginName string) error {
	pluginPath := filepath.Join(l.pluginDir, pluginName+l.getPluginExtension())

	// 检查文件是否存在
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s not found", pluginName)
	}

	// 删除插件文件
	if err := os.Remove(pluginPath); err != nil {
		return fmt.Errorf("failed to remove plugin: %v", err)
	}

	utils.Info("成功卸载插件: %s", pluginName)
	return nil
}