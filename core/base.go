package core

import (
	"database/sql"
	"log"
	"sync"

	"github.com/ctwj/urldb/plugin/hook"
	"github.com/ctwj/urldb/db/entity"
)

// BaseApp 实现 App 接口的基础结构
type BaseApp struct {
	mu            sync.RWMutex
	bootstrapped  bool
	db            *sql.DB
	router        RouterInterface
	config        ConfigInterface
	logger        LoggerInterface
	dataDir       string

	// 钩子实例
	onBootstrap      *hook.Hook[*BootstrapEvent]
	onServe          *hook.Hook[*ServeEvent]
	onTerminate      *hook.Hook[*TerminateEvent]

	onURLAdd         *hook.Hook[*URLEvent]
	onURLAccess      *hook.Hook[*URLAccessEvent]
	onURLUpdate      *hook.Hook[*URLEvent]
	onURLDelete      *hook.Hook[*URLEvent]

	onUserLogin      *hook.Hook[*UserEvent]
	onUserLogout     *hook.Hook[*UserEvent]
	onUserRegister   *hook.Hook[*UserEvent]

	onCategoryCreate *hook.Hook[*CategoryEvent]
	onTagAdd         *hook.Hook[*TagEvent]

	onAPIRequest     *hook.Hook[*APIEvent]
	onAPIResponse    *hook.Hook[*APIResponseEvent]

	onCustomEvent    *hook.Hook[*CustomEvent]
}

// NewBaseApp 创建新的基础应用实例
func NewBaseApp() *BaseApp {
	app := &BaseApp{
		onBootstrap:      &hook.Hook[*BootstrapEvent]{},
		onServe:          &hook.Hook[*ServeEvent]{},
		onTerminate:      &hook.Hook[*TerminateEvent]{},

		onURLAdd:         &hook.Hook[*URLEvent]{},
		onURLAccess:      &hook.Hook[*URLAccessEvent]{},
		onURLUpdate:      &hook.Hook[*URLEvent]{},
		onURLDelete:      &hook.Hook[*URLEvent]{},

		onUserLogin:      &hook.Hook[*UserEvent]{},
		onUserLogout:     &hook.Hook[*UserEvent]{},
		onUserRegister:   &hook.Hook[*UserEvent]{},

		onCategoryCreate: &hook.Hook[*CategoryEvent]{},
		onTagAdd:         &hook.Hook[*TagEvent]{},

		onAPIRequest:     &hook.Hook[*APIEvent]{},
		onAPIResponse:    &hook.Hook[*APIResponseEvent]{},

		onCustomEvent:    &hook.Hook[*CustomEvent]{},
	}

	return app
}

// --- App 接口实现 ---

func (app *BaseApp) Bootstrap() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.bootstrapped {
		return nil
	}

	// 触发启动钩子
	event := &BootstrapEvent{App: app}
	if err := app.onBootstrap.Trigger(event); err != nil {
		return err
	}

	app.bootstrapped = true
	return nil
}

func (app *BaseApp) IsBootstrapped() bool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.bootstrapped
}

func (app *BaseApp) Restart() error {
	// 实现应用重启逻辑
	// 这里可能需要重新初始化各个组件
	log.Println("Restarting application...")
	return nil
}

func (app *BaseApp) DataDir() string {
	return app.dataDir
}

func (app *BaseApp) DB() *sql.DB {
	return app.db
}

func (app *BaseApp) WithURLDB(fn func(interface{}) error) error {
	// 这里需要根据现有的数据库仓库结构来实现
	// 暂时返回一个简单的实现
	return nil
}

func (app *BaseApp) Router() RouterInterface {
	return app.router
}

func (app *BaseApp) Config() ConfigInterface {
	return app.config
}

func (app *BaseApp) Logger() LoggerInterface {
	return app.logger
}

// --- 钩子访问方法 ---

func (app *BaseApp) OnBootstrap() *hook.Hook[*BootstrapEvent] {
	return app.onBootstrap
}

func (app *BaseApp) OnServe() *hook.Hook[*ServeEvent] {
	return app.onServe
}

func (app *BaseApp) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.onTerminate
}

func (app *BaseApp) OnURLAdd() *hook.Hook[*URLEvent] {
	return app.onURLAdd
}

func (app *BaseApp) OnURLAccess() *hook.Hook[*URLAccessEvent] {
	return app.onURLAccess
}

func (app *BaseApp) OnURLUpdate() *hook.Hook[*URLEvent] {
	return app.onURLUpdate
}

func (app *BaseApp) OnURLDelete() *hook.Hook[*URLEvent] {
	return app.onURLDelete
}

func (app *BaseApp) OnUserLogin() *hook.Hook[*UserEvent] {
	return app.onUserLogin
}

func (app *BaseApp) OnUserLogout() *hook.Hook[*UserEvent] {
	return app.onUserLogout
}

func (app *BaseApp) OnUserRegister() *hook.Hook[*UserEvent] {
	return app.onUserRegister
}

func (app *BaseApp) OnCategoryCreate() *hook.Hook[*CategoryEvent] {
	return app.onCategoryCreate
}

func (app *BaseApp) OnTagAdd() *hook.Hook[*TagEvent] {
	return app.onTagAdd
}

func (app *BaseApp) OnAPIRequest() *hook.Hook[*APIEvent] {
	return app.onAPIRequest
}

func (app *BaseApp) OnAPIResponse() *hook.Hook[*APIResponseEvent] {
	return app.onAPIResponse
}

func (app *BaseApp) OnCustomEvent() *hook.Hook[*CustomEvent] {
	return app.onCustomEvent
}

// --- 设置方法 ---

func (app *BaseApp) SetDB(db *sql.DB) {
	app.db = db
}

func (app *BaseApp) SetRouter(router RouterInterface) {
	app.router = router
}

func (app *BaseApp) SetConfig(config ConfigInterface) {
	app.config = config
}

func (app *BaseApp) SetLogger(logger LoggerInterface) {
	app.logger = logger
}

func (app *BaseApp) SetDataDir(dir string) {
	app.dataDir = dir
}

// --- 便捷触发方法 ---

// TriggerURLAdd 触发 URL 添加事件
func (app *BaseApp) TriggerURLAdd(url *entity.Resource, data map[string]interface{}) error {
	event := &URLEvent{
		App:  app,
		URL:  url,
		Data: data,
	}
	return app.onURLAdd.Trigger(event)
}

// TriggerURLAccess 触发 URL 访问事件
func (app *BaseApp) TriggerURLAccess(url *entity.Resource, accessLog interface{}, request, response interface{}) error {
	event := &URLAccessEvent{
		App:       app,
		URL:       url,
		AccessLog: accessLog,
		Request:   request,
		Response:  response,
	}
	return app.onURLAccess.Trigger(event)
}

// TriggerUserLogin 触发用户登录事件
func (app *BaseApp) TriggerUserLogin(user *entity.User, data map[string]interface{}) error {
	event := &UserEvent{
		App:  app,
		User: user,
		Data: data,
	}
	return app.onUserLogin.Trigger(event)
}

// TriggerCustomEvent 触发自定义事件
func (app *BaseApp) TriggerCustomEvent(name string, data map[string]interface{}) error {
	event := &CustomEvent{
		App:  app,
		Name: name,
		Data: data,
	}
	return app.onCustomEvent.Trigger(event)
}