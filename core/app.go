package core

import (
	"database/sql"

	"github.com/ctwj/urldb/plugin/hook"
	"github.com/ctwj/urldb/db/entity"
)

// App 定义 URLDB 应用的核心接口
type App interface {
	// 基础应用方法
	Bootstrap() error
	IsBootstrapped() bool
	Restart() error
	DataDir() string

	// 数据库访问
	DB() *sql.DB
	WithURLDB(fn func(interface{}) error) error

	// 路由和HTTP服务（假设使用 Gin）
	Router() RouterInterface

	// 配置和日志
	Config() ConfigInterface
	Logger() LoggerInterface

	// --- 插件系统钩子 ---

	// 应用生命周期钩子
	OnBootstrap() *hook.Hook[*BootstrapEvent]
	OnServe() *hook.Hook[*ServeEvent]
	OnTerminate() *hook.Hook[*TerminateEvent]

	// URL 相关钩子
	OnURLAdd() *hook.Hook[*URLEvent]
	OnURLAccess() *hook.Hook[*URLAccessEvent]
	OnURLUpdate() *hook.Hook[*URLEvent]
	OnURLDelete() *hook.Hook[*URLEvent]

	// 用户相关钩子
	OnUserLogin() *hook.Hook[*UserEvent]
	OnUserLogout() *hook.Hook[*UserEvent]
	OnUserRegister() *hook.Hook[*UserEvent]

	// 分类和标签钩子
	OnCategoryCreate() *hook.Hook[*CategoryEvent]
	OnTagAdd() *hook.Hook[*TagEvent]

	// API 钩子
	OnAPIRequest() *hook.Hook[*APIEvent]
	OnAPIResponse() *hook.Hook[*APIResponseEvent]

	// 自定义事件钩子
	OnCustomEvent() *hook.Hook[*CustomEvent]

	// 待处理资源钩子
	OnReadyResourceAdd() *hook.Hook[*ReadyResourceEvent]
}

// RouterInterface 路由接口（适配你的路由框架）
type RouterInterface interface {
	GET(path string, handler interface{})
	POST(path string, handler interface{})
	PUT(path string, handler interface{})
	DELETE(path string, handler interface{})
	PATCH(path string, handler interface{})

	// 中间件支持
	Use(middleware interface{})
	Group(path string) RouterInterface
}

// ConfigInterface 配置接口
type ConfigInterface interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{})
}

// LoggerInterface 日志接口
type LoggerInterface interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// BootstrapEvent 应用启动事件
type BootstrapEvent struct {
	hook.Event
	App App
}

// ServeEvent 服务启动事件
type ServeEvent struct {
	hook.Event
	App    App
	Router RouterInterface
}

// TerminateEvent 应用终止事件
type TerminateEvent struct {
	hook.Event
	App App
}

// URLEvent URL 相关事件
type URLEvent struct {
	hook.Event
	App  App
	URL  *entity.Resource
	Data map[string]interface{} // 额外数据
}

// URLAccessEvent URL 访问事件
type URLAccessEvent struct {
	hook.Event
	App       App
	URL       *entity.Resource
	AccessLog interface{} // 可以用 APIAccessLog
	Request   interface{} // HTTP 请求对象
	Response  interface{} // HTTP 响应对象
}

// UserEvent 用户相关事件
type UserEvent struct {
	hook.Event
	App  App
	User *entity.User
	Data map[string]interface{} // 额外数据，如登录信息等
}

// CategoryEvent 分类相关事件
type CategoryEvent struct {
	hook.Event
	App      App
	Category *entity.Category
}

// TagEvent 标签相关事件
type TagEvent struct {
	hook.Event
	App App
	Tag *entity.Tag
	URL *entity.Resource
}

// APIEvent API 请求事件
type APIEvent struct {
	hook.Event
	App     App
	Request interface{} // HTTP 请求对象
	Path    string
	Method  string
	Headers map[string]string
	Body    interface{}
}

// APIResponseEvent API 响应事件
type APIResponseEvent struct {
	hook.Event
	App      App
	Request  interface{} // HTTP 请求对象
	Response interface{} // HTTP 响应对象
	Status   int
	Body     interface{}
}

// CustomEvent 自定义事件
type CustomEvent struct {
	hook.Event
	App  App
	Name string
	Data map[string]interface{}
}

// ReadyResourceEvent 待处理资源事件
type ReadyResourceEvent struct {
	hook.Event
	App            App
	ReadyResource  *entity.ReadyResource
	Data           map[string]interface{} // 额外数据
}