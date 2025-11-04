package demo

import (
	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/types"
)

// init 函数会在包导入时自动调用
func init() {
	// 注册演示插件
	demoPlugin := NewDemoPlugin()
	RegisterPlugin(demoPlugin)

	// 注册完整演示插件
	fullDemoPlugin := NewFullDemoPlugin()
	RegisterPlugin(fullDemoPlugin)

	// 注册依赖插件
	dependentPlugin := NewDependentPlugin()
	RegisterPlugin(dependentPlugin)

	// 注册依赖演示插件
	databasePlugin := NewDatabasePlugin()
	RegisterPlugin(databasePlugin)

	authPlugin := NewAuthPlugin()
	RegisterPlugin(authPlugin)

	businessPlugin := NewBusinessPlugin()
	RegisterPlugin(businessPlugin)
}

// RegisterPlugin 注册插件到全局管理器
func RegisterPlugin(pluginInstance types.Plugin) {
	if plugin.GetManager() != nil {
		plugin.GetManager().RegisterPlugin(pluginInstance)
	}
}