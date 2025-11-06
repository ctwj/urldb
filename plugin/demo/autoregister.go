package demo

import (
	"sync"

	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/types"
)

var (
	pendingPlugins []types.Plugin
	once           sync.Once
)

// init 函数会在包导入时自动调用
func init() {
	// 准备要注册的插件，但不立即注册
	pendingPlugins = append(pendingPlugins,
		NewDemoPlugin(),
		NewFullDemoPlugin(),
		NewDependentPlugin(),
		NewDatabasePlugin(),
		NewAuthPlugin(),
		NewBusinessPlugin(),
	)
}

// RegisterPendingPlugins 注册所有待注册的插件
func RegisterPendingPlugins() {
	once.Do(func() {
		if plugin.GetManager() != nil {
			for _, pluginInstance := range pendingPlugins {
				if err := plugin.GetManager().RegisterPlugin(pluginInstance); err != nil {
					plugin.GetLogger().Error("Failed to register plugin %s: %v", pluginInstance.Name(), err)
				} else {
					plugin.GetLogger().Info("Successfully registered plugin: %s", pluginInstance.Name())
				}
			}
		}
	})
}

// RegisterPlugin 注册插件到全局管理器
func RegisterPlugin(pluginInstance types.Plugin) {
	if plugin.GetManager() != nil {
		if err := plugin.GetManager().RegisterPlugin(pluginInstance); err != nil {
			plugin.GetLogger().Error("Failed to register plugin %s: %v", pluginInstance.Name(), err)
		} else {
			plugin.GetLogger().Info("Successfully registered plugin: %s", pluginInstance.Name())
		}
	} else {
		// 管理器还未初始化，添加到待注册列表
		pendingPlugins = append(pendingPlugins, pluginInstance)
	}
}