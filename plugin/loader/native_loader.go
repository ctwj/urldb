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

// NativePluginLoader 直接加载 .so 文件的插件加载器
type NativePluginLoader struct {
	pluginDir string
}

// NewNativePluginLoader 创建原生插件加载器
func NewNativePluginLoader(pluginDir string) *NativePluginLoader {
	return &NativePluginLoader{
		pluginDir: pluginDir,
	}
}

// LoadPlugin 加载单个插件文件
func (l *NativePluginLoader) LoadPlugin(filename string) (types.Plugin, error) {
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

	utils.Info("成功加载插件: %s", filename)
	return pluginInstance, nil
}

// LoadAllPlugins 加载所有插件
func (l *NativePluginLoader) LoadAllPlugins() ([]types.Plugin, error) {
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
			plugin, err := l.LoadPlugin(file.Name())
			if err != nil {
				utils.Error("加载插件 %s 失败: %v", file.Name(), err)
				continue
			}
			plugins = append(plugins, plugin)
		}
	}

	utils.Info("总共加载了 %d 个插件", len(plugins))
	return plugins, nil
}

// isPluginFile 检查是否为插件文件
func (l *NativePluginLoader) isPluginFile(filename string) bool {
	// 根据平台检查文件扩展名
	return strings.HasSuffix(filename, l.getPluginExtension())
}

// getPluginExtension 获取当前平台的插件文件扩展名
func (l *NativePluginLoader) getPluginExtension() string {
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

// GetPluginInfo 获取插件文件信息
func (l *NativePluginLoader) GetPluginInfo(filename string) (map[string]interface{}, error) {
	pluginPath := filepath.Join(l.pluginDir, filename)

	info := make(map[string]interface{})
	info["filename"] = filename
	info["path"] = pluginPath

	// 获取文件信息
	fileInfo, err := os.Stat(pluginPath)
	if err != nil {
		return nil, err
	}

	info["size"] = fileInfo.Size()
	info["mod_time"] = fileInfo.ModTime()
	info["is_plugin"] = l.isPluginFile(filename)

	return info, nil
}

// ListPluginFiles 列出所有插件文件
func (l *NativePluginLoader) ListPluginFiles() ([]map[string]interface{}, error) {
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