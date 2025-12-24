package plugin

import (
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/plugin/jsvm"
)

// Manager 插件管理器
type Manager struct {
	app core.App
}

// NewManager 创建插件管理器
func NewManager(app core.App) *Manager {
	return &Manager{
		app: app,
	}
}

// RegisterJSVM 注册 JavaScript 虚拟机插件
func (m *Manager) RegisterJSVM(config jsvm.Config) error {
	return jsvm.Register(m.app, config)
}

// RegisterJSVMDefault 注册默认配置的 JSVM 插件
func (m *Manager) RegisterJSVMDefault() error {
	config := jsvm.Config{
		HooksWatch:      true,
		HooksPoolSize:   10,
		OnInit:          m.defaultOnInit,
	}
	return m.RegisterJSVM(config)
}

// defaultOnInit 默认的 VM 初始化回调
func (m *Manager) defaultOnInit(vm interface{}) {
	// 可以在这里添加自定义的全局变量或函数
	// 例如：vm.Set("version", "1.0.0")
}