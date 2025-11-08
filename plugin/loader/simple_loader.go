package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
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

	// 尝试直接断言
	pluginInstance, ok := sym.(types.Plugin)
	if !ok {
		// 如果直接断言失败，尝试使用反射来创建一个兼容的包装器
		utils.Info("类型断言失败: %s，尝试使用反射创建包装器", filename)

		// 使用反射来动态创建一个实现了types.Plugin接口的包装器
		wrapper, err := l.createPluginWrapper(sym)
		if err != nil {
			utils.Error("创建插件包装器失败: %v", err)
			return nil, fmt.Errorf("invalid plugin type in %s: %v", filename, err)
		}

		utils.Info("使用反射包装器加载插件: %s (名称: %s, 版本: %s)",
			filename, wrapper.Name(), wrapper.Version())
		return wrapper, nil
	}

	utils.Info("成功加载插件: %s (名称: %s, 版本: %s)",
		filename, pluginInstance.Name(), pluginInstance.Version())
	return pluginInstance, nil
}

// LoadAllPlugins 加载所有插件
func (l *SimplePluginLoader) LoadAllPlugins() ([]types.Plugin, error) {
	var plugins []types.Plugin

	utils.Info("开始从目录加载插件: %s", l.pluginDir)

	// 确保插件目录存在
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// 读取插件目录
	files, err := ioutil.ReadDir(l.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %v", err)
	}

	utils.Info("插件目录中的文件数量: %d", len(files))

	for _, file := range files {
		utils.Info("检查文件: %s (是否为目录: %t)", file.Name(), file.IsDir())
		if file.IsDir() {
			continue
		}

		if l.isPluginFile(file.Name()) {
			utils.Info("尝试加载插件文件: %s", file.Name())
			pluginInstance, err := l.LoadPlugin(file.Name())
			if err != nil {
				utils.Error("加载插件 %s 失败: %v", file.Name(), err)
				continue
			}
			plugins = append(plugins, pluginInstance)
			utils.Info("成功加载插件: %s", pluginInstance.Name())
		} else {
			utils.Info("文件 %s 不是插件文件", file.Name())
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

// pluginWrapper 是一个使用反射来包装插件实例的结构体
type pluginWrapper struct {
	original interface{}
	value    reflect.Value
	methods  map[string]reflect.Value
}

// createPluginWrapper 创建一个插件包装器
func (l *SimplePluginLoader) createPluginWrapper(plugin interface{}) (types.Plugin, error) {
	pluginValue := reflect.ValueOf(plugin)

	if pluginValue.Kind() != reflect.Ptr && pluginValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("plugin must be a pointer or struct")
	}

	// 获取所有方法
	methods := make(map[string]reflect.Value)

	// 获取方法列表
	methodNames := []string{
		"Name", "Version", "Description", "Author",
		"Initialize", "Start", "Stop", "Cleanup",
		"Dependencies", "CheckDependencies",
	}

	for _, methodName := range methodNames {
		method := pluginValue.MethodByName(methodName)
		if method.IsValid() {
			methods[methodName] = method
		} else if pluginValue.Kind() == reflect.Ptr && pluginValue.Elem().IsValid() {
			// 如果是结构体指针，检查结构体本身是否有方法
			method = pluginValue.Elem().MethodByName(methodName)
			if method.IsValid() {
				methods[methodName] = method
			}
		}
	}

	// 至少要有一些方法
	if len(methods) == 0 {
		return nil, fmt.Errorf("plugin has no valid methods")
	}

	wrapper := &pluginWrapper{
		original: plugin,
		value:    pluginValue,
		methods:  methods,
	}

	// 验证包装器是否实现了所有必要的方法
	if err := l.validatePluginWrapper(wrapper); err != nil {
		return nil, fmt.Errorf("plugin does not implement required methods: %v", err)
	}

	return wrapper, nil
}

// validatePluginWrapper 验证包装器是否实现了所有必要的方法
func (l *SimplePluginLoader) validatePluginWrapper(wrapper *pluginWrapper) error {
	// 检查是否有所需的方法
	methods := []string{
		"Name", "Version", "Description", "Author",
		"Initialize", "Start", "Stop", "Cleanup",
		"Dependencies", "CheckDependencies",
	}

	for _, method := range methods {
		// 优先检查缓存的方法
		if _, exists := wrapper.methods[method]; exists {
			continue
		}

		// 检查作为方法
		if methodValue := wrapper.value.MethodByName(method); methodValue.IsValid() {
			continue
		}

		// 如果作为方法不存在，检查是否可以直接调用（通过接口）
		// 首先检查是否为指针类型
		if wrapper.value.Kind() == reflect.Ptr && wrapper.value.Elem().IsValid() {
			if field := wrapper.value.Elem().FieldByName(method); field.IsValid() && field.Kind() == reflect.Func {
				continue
			}
		}

		return fmt.Errorf("missing method: %s", method)
	}

	// 尝试调用Name方法，验证基本功能
	name := wrapper.Name()
	if name == "" || name == "unknown" {
		return fmt.Errorf("Name method returned empty string or unknown")
	}

	return nil
}

// 为pluginWrapper实现types.Plugin接口的所有方法
func (w *pluginWrapper) Name() string {
	if method, exists := w.methods["Name"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return "unknown"
}

func (w *pluginWrapper) Version() string {
	if method, exists := w.methods["Version"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return "unknown"
}

func (w *pluginWrapper) Description() string {
	if method, exists := w.methods["Description"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return "No description"
}

func (w *pluginWrapper) Author() string {
	if method, exists := w.methods["Author"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return "Unknown"
}

func (w *pluginWrapper) Initialize(ctx types.PluginContext) error {
	if method, exists := w.methods["Initialize"]; exists && method.IsValid() {
		// 检查参数类型是否兼容
		if method.Type().NumIn() > 0 {
			results := method.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if len(results) > 0 {
				// 假设返回的是error类型
				if errVal := results[0].Interface(); errVal != nil {
					if err, ok := errVal.(error); ok {
						return err
					}
					return fmt.Errorf("initialization failed: %v", errVal)
				}
			}
			return nil
		}
	}

	return nil // 如果方法不存在，返回nil表示成功
}

func (w *pluginWrapper) Start() error {
	if method, exists := w.methods["Start"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			// 假设返回的是error类型
			if errVal := results[0].Interface(); errVal != nil {
				if err, ok := errVal.(error); ok {
					return err
				}
				return fmt.Errorf("start failed: %v", errVal)
			}
		}
		return nil
	}

	return nil // 如果方法不存在，返回nil表示成功
}

func (w *pluginWrapper) Stop() error {
	if method, exists := w.methods["Stop"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			// 假设返回的是error类型
			if errVal := results[0].Interface(); errVal != nil {
				if err, ok := errVal.(error); ok {
					return err
				}
				return fmt.Errorf("stop failed: %v", errVal)
			}
		}
		return nil
	}

	return nil // 如果方法不存在，返回nil表示成功
}

func (w *pluginWrapper) Cleanup() error {
	if method, exists := w.methods["Cleanup"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			// 假设返回的是error类型
			if errVal := results[0].Interface(); errVal != nil {
				if err, ok := errVal.(error); ok {
					return err
				}
				return fmt.Errorf("cleanup failed: %v", errVal)
			}
		}
		return nil
	}

	return nil // 如果方法不存在，返回nil表示成功
}

func (w *pluginWrapper) Dependencies() []string {
	if method, exists := w.methods["Dependencies"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			if deps, ok := results[0].Interface().([]string); ok {
				return deps
			}
		}
	}
	return []string{}
}

func (w *pluginWrapper) CheckDependencies() map[string]bool {
	if method, exists := w.methods["CheckDependencies"]; exists && method.IsValid() {
		results := method.Call([]reflect.Value{})
		if len(results) > 0 {
			if checks, ok := results[0].Interface().(map[string]bool); ok {
				return checks
			}
		}
	}
	return map[string]bool{}
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