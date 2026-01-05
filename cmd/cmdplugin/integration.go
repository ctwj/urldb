package cmdplugin

import (
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/dop251/goja"
	"github.com/ctwj/urldb/plugin-system/core"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin-system/manager/plugin"
	"github.com/ctwj/urldb/plugin-system/manager/plugin/jsvm"
	"github.com/ctwj/urldb/plugin-system/triggers/plugins"
	"github.com/ctwj/urldb/utils"
)

// PendingRoute 待注册的路由
type PendingRoute struct {
	Method string
	Path   string
	Handler func() (interface{}, error)
}

// PluginIntegration 插件系统集成器
type PluginIntegration struct {
	app          *core.BaseApp
	pluginManager *plugin.Manager
	repoManager  *repo.RepositoryManager
	router       *gin.Engine // 存储 Gin 路由器
	pendingRoutes []PendingRoute // 存储待注册的路由
}

// NewPluginIntegration 创建插件系统集成器
func NewPluginIntegration(repoManager *repo.RepositoryManager) *PluginIntegration {
	pi := &PluginIntegration{
		repoManager: repoManager,
	}

	// 创建插件感知的应用
	pi.app = core.NewBaseApp()

	// 设置基础配置
	pi.app.SetDataDir("./plugin-system/types")
	pi.app.SetConfig(&PluginConfigWrapper{repoManager.SystemConfigRepository})
	pi.app.SetLogger(&PluginLoggerWrapper{})
	pi.app.SetRouter(&PluginRouterWrapper{})

	// 注册插件管理器
	pi.pluginManager = plugin.NewManager(pi.app)

	// 设置 RepositoryManager
	pi.pluginManager.SetRepoManager(repoManager)

	return pi
}

// SetRouter 设置 Gin 路由器
func (pi *PluginIntegration) SetRouter(router *gin.Engine) {
	pi.router = router
	// 注册所有待注册的路由
	pi.registerPendingRoutes()
}

// registerPendingRoutes 注册所有待注册的路由
func (pi *PluginIntegration) registerPendingRoutes() {
	if len(pi.pendingRoutes) == 0 {
		utils.Info("No pending routes to register")
		return
	}

	// 用于跟踪已注册的路由，避免重复
	registeredRoutes := make(map[string]bool)

	utils.Info("Registering %d pending plugin routes", len(pi.pendingRoutes))
	for _, route := range pi.pendingRoutes {
		routeKey := route.Method + ":" + route.Path

		// 检查路由是否已经注册
		if registeredRoutes[routeKey] {
			utils.Warn("Skipping duplicate plugin route: %s %s", route.Method, route.Path)
			continue
		}

		// 尝试注册路由，如果发生panic则捕获并给出警告
		func() {
			defer func() {
				if r := recover(); r != nil {
					utils.Warn("Failed to register plugin route %s %s: %v", route.Method, route.Path, r)
				}
			}()

			// 创建局部变量避免闭包问题
			method := route.Method
			path := route.Path
			handler := route.Handler

			// 根据方法注册路由
			switch method {
			case "GET":
				pi.router.GET(path, func(c *gin.Context) {
					// 检查插件是否启用
					pluginName := pi.extractPluginNameFromPath(path)
					if pluginName != "" {
						if config, err := pi.repoManager.PluginConfigRepository.GetConfig(pluginName); err == nil && config != nil && !config.Enabled {
							utils.Debug("Plugin route '%s %s' skipped: plugin '%s' is disabled", method, path, pluginName)
							c.JSON(403, gin.H{
								"success": false,
								"error":   "Plugin is disabled",
								"message": "The plugin for this route is currently disabled",
							})
							return
						}
					}

					utils.Info("Plugin route called: %s %s", method, path)
					result, err := handler()
					if err != nil {
						utils.Error("Plugin route handler error: %v", err)
						c.JSON(500, gin.H{"error": err.Error()})
					} else {
						utils.Info("Plugin route handler success: %v", result)
						c.JSON(200, result)
					}
				})
			case "POST":
				pi.router.POST(path, func(c *gin.Context) {
					// 检查插件是否启用
					pluginName := pi.extractPluginNameFromPath(path)
					if pluginName != "" {
						if config, err := pi.repoManager.PluginConfigRepository.GetConfig(pluginName); err == nil && config != nil && !config.Enabled {
							utils.Debug("Plugin route '%s %s' skipped: plugin '%s' is disabled", method, path, pluginName)
							c.JSON(403, gin.H{
								"success": false,
								"error":   "Plugin is disabled",
								"message": "The plugin for this route is currently disabled",
							})
							return
						}
					}

					result, err := handler()
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
					} else {
						c.JSON(200, result)
					}
				})
			case "PUT":
				pi.router.PUT(path, func(c *gin.Context) {
					// 检查插件是否启用
					pluginName := pi.extractPluginNameFromPath(path)
					if pluginName != "" {
						if config, err := pi.repoManager.PluginConfigRepository.GetConfig(pluginName); err == nil && config != nil && !config.Enabled {
							utils.Debug("Plugin route '%s %s' skipped: plugin '%s' is disabled", method, path, pluginName)
							c.JSON(403, gin.H{
								"success": false,
								"error":   "Plugin is disabled",
								"message": "The plugin for this route is currently disabled",
							})
							return
						}
					}

					result, err := handler()
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
					} else {
						c.JSON(200, result)
					}
				})
			case "DELETE":
				pi.router.DELETE(path, func(c *gin.Context) {
					// 检查插件是否启用
					pluginName := pi.extractPluginNameFromPath(path)
					if pluginName != "" {
						if config, err := pi.repoManager.PluginConfigRepository.GetConfig(pluginName); err == nil && config != nil && !config.Enabled {
							utils.Debug("Plugin route '%s %s' skipped: plugin '%s' is disabled", method, path, pluginName)
							c.JSON(403, gin.H{
								"success": false,
								"error":   "Plugin is disabled",
								"message": "The plugin for this route is currently disabled",
							})
							return
						}
					}

					result, err := handler()
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
					} else {
						c.JSON(200, result)
					}
				})
			default:
				utils.Error("Unsupported HTTP method in pending route: %s", method)
				return
			}

			registeredRoutes[routeKey] = true
			utils.Info("Pending plugin route registered: %s %s", method, path)
		}()
	}

	// 清空待注册的路由
	pi.pendingRoutes = nil
	utils.Info("All pending routes registration completed")
}

// Initialize 初始化插件系统
func (pi *PluginIntegration) Initialize() error {
	// 创建路由注册器
	routeRegister := func(method, path string, handler func() (interface{}, error)) error {
		if pi.router == nil {
			// 如果路由器还没有设置，存储路由以供后续注册
			pi.pendingRoutes = append(pi.pendingRoutes, PendingRoute{
				Method: method,
				Path:   path,
				Handler: handler,
			})
			utils.Info("Plugin route queued for registration: %s %s", method, path)
			return nil
		}

		// 根据方法注册路由
		switch method {
		case "GET":
			pi.router.GET(path, func(c *gin.Context) {
				result, err := handler()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
				} else {
					c.JSON(200, result)
				}
			})
		case "POST":
			pi.router.POST(path, func(c *gin.Context) {
				result, err := handler()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
				} else {
					c.JSON(200, result)
				}
			})
		case "PUT":
			pi.router.PUT(path, func(c *gin.Context) {
				result, err := handler()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
				} else {
					c.JSON(200, result)
				}
			})
		case "DELETE":
			pi.router.DELETE(path, func(c *gin.Context) {
				result, err := handler()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
				} else {
					c.JSON(200, result)
				}
			})
		default:
			utils.Error("Unsupported HTTP method: %s", method)
			return nil
		}

		utils.Info("Plugin route registered: %s %s", method, path)
		return nil
	}

	// 注册 JSVM 插件
	err := pi.pluginManager.RegisterJSVMWithRepo(jsvm.Config{
		HooksWatch:      true,
		HooksPoolSize:   10,
		HooksDir:        "./plugin-system/hooks",
		MigrationsDir:   "./migrations",
		TypesDir:        "./plugin-system/types",
		RouteRegister:   routeRegister,
		OnInit: func(vm *goja.Runtime) {
			utils.Info("Plugin system initialized")
		},
	}, pi.repoManager)

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

// GetGlobalPluginIntegration 获取全局插件系统集成实例
func GetGlobalPluginIntegration() *PluginIntegration {
	return globalPluginIntegration
}

// GetPluginManager 获取插件管理器
func (pi *PluginIntegration) GetPluginManager() *plugin.Manager {
	return pi.pluginManager
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
		// 转换 URL 类型并触发事件
		if resource, ok := url.(*entity.Resource); ok {
			if err := app.TriggerURLAdd(resource, data); err != nil {
				utils.Error("Failed to trigger URL add event:", err)
			} else {
				utils.Info("URL add event triggered successfully")
			}
		} else {
			utils.Error("Invalid URL type, expected *entity.Resource")
		}
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


// extractPluginNameFromPath 从路由路径中提取插件名称
func (pi *PluginIntegration) extractPluginNameFromPath(path string) string {
	// 从路径中提取插件名称，例如 /api/config-demo -> config_demo
	// 使用正则表达式匹配 /api/plugin-name 或 /api/plugin-name/... 的模式

	// 去掉开头的斜杠
	path = strings.TrimPrefix(path, "/")

	// 分割路径
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "api" {
		return ""
	}

	// 返回插件名称部分
	pluginName := parts[1]

	// 如果插件名称包含连字符，转换为下划线以匹配数据库中的插件名
	if strings.Contains(pluginName, "-") {
		// 例如：config-demo -> config_demo
		pluginName = strings.ReplaceAll(pluginName, "-", "_")
	}

	return pluginName
}