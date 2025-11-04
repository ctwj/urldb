package demo

import (
	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/types"
)

// 自动注册示例插件到全局管理器
func init() {
	// 注册配置演示插件
	configDemoPlugin := NewConfigDemoPlugin()
	RegisterConfigDemoPlugin(configDemoPlugin)
}

// RegisterConfigDemoPlugin 注册配置演示插件到全局管理器
func RegisterConfigDemoPlugin(pluginInstance types.Plugin) {
	if plugin.GetManager() != nil {
		plugin.GetManager().RegisterPlugin(pluginInstance)
	}
}