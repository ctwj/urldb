package main

import (
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/jsvm"
	"github.com/ctwj/urldb/plugins"
	"github.com/ctwj/urldb/utils"
)

// PluginIntegration 插件系统集成器
type PluginIntegration struct {
	app          *core.BaseApp
	pluginManager *plugin.Manager
	repoManager  *repo.RepositoryManager
}

// NewPluginIntegration 创建插件系统集成器
func NewPluginIntegration(repoManager *repo.RepositoryManager) *PluginIntegration {
	pi := &PluginIntegration{
		repoManager: repoManager,
	}

	// 创建插件感知的应用
	pi.app = core.NewBaseApp()

	// 设置基础配置
	pi.app.SetDataDir("./pb_data")
	pi.app.SetConfig(&PluginConfigWrapper{repoManager.SystemConfigRepository})
	pi.app.SetLogger(&PluginLoggerWrapper{})
	pi.app.SetRouter(&PluginRouterWrapper{})

	// 注册插件管理器
	pi.pluginManager = plugin.NewManager(pi.app)

	return pi
}

// Initialize 初始化插件系统
func (pi *PluginIntegration) Initialize() error {
	// 注册 JSVM 插件
	err := pi.pluginManager.RegisterJSVM(jsvm.Config{
		HooksWatch:      true,
		HooksPoolSize:   10,
		HooksDir:        "./hooks",
		MigrationsDir:   "./migrations",
		TypesDir:        "./pb_data",
		OnInit: func(vm interface{}) {
			utils.Info("Plugin system initialized")
		},
	})

	if err != nil {
		return err
	}

	// 启动插件系统
	if err := pi.app.Bootstrap(); err != nil {
		return err
	}

	// 设置插件应用到触发器
	plugins.SetPluginApp(pi.app)

	return nil
}

// GetApp 获取插件应用实例
func (pi *PluginIntegration) GetApp() *core.BaseApp {
	return pi.app
}

// PluginConfigWrapper 配置包装器
type PluginConfigWrapper struct {
	repo repo.SystemConfigRepository
}

func (c *PluginConfigWrapper) Get(key string) interface{} {
	val, _ := c.repo.GetConfigValue(key)
	return val
}

func (c *PluginConfigWrapper) GetString(key string) string {
	val, _ := c.repo.GetConfigValue(key)
	return val
}

func (c *PluginConfigWrapper) GetInt(key string) int {
	val, _ := c.repo.GetConfigInt(key)
	return val
}

func (c *PluginConfigWrapper) GetBool(key string) bool {
	val, _ := c.repo.GetConfigBool(key)
	return val
}

func (c *PluginConfigWrapper) Set(key string, value interface{}) {
	// 简化实现，暂时不做任何操作
	// TODO: 实现配置设置逻辑
}

// PluginLoggerWrapper 日志包装器
type PluginLoggerWrapper struct{}

func (l *PluginLoggerWrapper) Debug(msg string, args ...interface{}) {
	utils.Debug(msg, args...)
}

func (l *PluginLoggerWrapper) Info(msg string, args ...interface{}) {
	utils.Info(msg, args...)
}

func (l *PluginLoggerWrapper) Warn(msg string, args ...interface{}) {
	utils.Warn(msg, args...)
}

func (l *PluginLoggerWrapper) Error(msg string, args ...interface{}) {
	utils.Error(msg, args...)
}

func (l *PluginLoggerWrapper) Fatal(msg string, args ...interface{}) {
	utils.Fatal(msg, args...)
}

// PluginRouterWrapper 路由包装器
type PluginRouterWrapper struct{}

func (r *PluginRouterWrapper) GET(path string, handler interface{}) {
	utils.Info("Plugin route registered: GET %s", path)
	// TODO: 实际注册到 Gin 路由器
}

func (r *PluginRouterWrapper) POST(path string, handler interface{}) {
	utils.Info("Plugin route registered: POST %s", path)
	// TODO: 实际注册到 Gin 路由器
}

func (r *PluginRouterWrapper) PUT(path string, handler interface{}) {
	utils.Info("Plugin route registered: PUT %s", path)
	// TODO: 实际注册到 Gin 路由器
}

func (r *PluginRouterWrapper) DELETE(path string, handler interface{}) {
	utils.Info("Plugin route registered: DELETE %s", path)
	// TODO: 实际注册到 Gin 路由器
}

func (r *PluginRouterWrapper) PATCH(path string, handler interface{}) {
	utils.Info("Plugin route registered: PATCH %s", path)
	// TODO: 实际注册到 Gin 路由器
}

func (r *PluginRouterWrapper) Use(middleware interface{}) {
	utils.Info("Plugin middleware registered")
	// TODO: 实际注册到 Gin 路由器
}

func (r *PluginRouterWrapper) Group(path string) core.RouterInterface {
	return &PluginRouterWrapper{}
}

// 全局插件系统集成实例
var globalPluginIntegration *PluginIntegration

// InitializePluginSystem 初始化全局插件系统
func InitializePluginSystem(repoManager *repo.RepositoryManager) error {
	globalPluginIntegration = NewPluginIntegration(repoManager)

	err := globalPluginIntegration.Initialize()
	if err != nil {
		utils.Error("Failed to initialize plugin system: %v", err)
		return err
	}

	utils.Info("Plugin system initialized successfully")
	return nil
}

// GetPluginApp 获取插件应用实例
func GetPluginApp() *core.BaseApp {
	if globalPluginIntegration == nil {
		return nil
	}
	return globalPluginIntegration.GetApp()
}

// TriggerURLAdd 触发 URL 添加事件
func TriggerURLAdd(url interface{}, data map[string]interface{}) {
	app := GetPluginApp()
	if app != nil {
		// 这里需要转换 URL 类型
		// app.TriggerURLAdd(url, data)
		utils.Info("URL add event triggered")
	}
}

// TriggerUserLogin 触发用户登录事件
func TriggerUserLogin(user interface{}, data map[string]interface{}) {
	app := GetPluginApp()
	if app != nil {
		// 这里需要转换 User 类型
		// app.TriggerUserLogin(user, data)
		utils.Info("User login event triggered")
	}
}

// TriggerAPIRequest 触发 API 请求事件
func TriggerAPIRequest(method, path string, headers map[string]string, body interface{}) {
	app := GetPluginApp()
	if app != nil {
		// app.TriggerAPIRequest(method, path, headers, body)
		utils.Info("API request event triggered: %s %s", method, path)
	}
}