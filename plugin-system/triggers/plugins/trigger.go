package plugins

import (
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// PluginApp 插件应用接口
type PluginApp interface {
	TriggerURLAdd(url *entity.Resource, data map[string]interface{}) error
	TriggerUserLogin(user *entity.User, data map[string]interface{}) error
	TriggerURLAccess(url *entity.Resource, accessLog interface{}, request, response interface{}) error
	TriggerReadyResourceAdd(readyResource *entity.ReadyResource, data map[string]interface{}) error
}

var (
	// 全局插件应用实例
	pluginApp PluginApp
)

// SetPluginApp 设置插件应用实例
func SetPluginApp(app PluginApp) {
	pluginApp = app
}

// GetPluginApp 获取插件应用实例
func GetPluginApp() PluginApp {
	return pluginApp
}

// TriggerURLAdd 触发 URL 添加事件
func TriggerURLAdd(url *entity.Resource, data map[string]interface{}) {
	if pluginApp != nil {
		if err := pluginApp.TriggerURLAdd(url, data); err != nil {
			utils.Error("Failed to trigger URL add event: %v", err)
		} else {
			utils.Info("URL add event triggered for: %s", url.URL)
		}
	}
}

// TriggerUserLogin 触发用户登录事件
func TriggerUserLogin(user *entity.User, data map[string]interface{}) {
	if pluginApp != nil {
		if err := pluginApp.TriggerUserLogin(user, data); err != nil {
			utils.Error("Failed to trigger user login event: %v", err)
		} else {
			utils.Info("User login event triggered for: %s", user.Username)
		}
	}
}

// TriggerURLAccess 触发 URL 访问事件
func TriggerURLAccess(url *entity.Resource, accessLog interface{}, request, response interface{}) {
	if pluginApp != nil {
		if err := pluginApp.TriggerURLAccess(url, accessLog, request, response); err != nil {
			utils.Error("Failed to trigger URL access event: %v", err)
		} else {
			utils.Info("URL access event triggered for: %s", url.URL)
		}
	}
}

// TriggerReadyResourceAdd 触发待处理资源添加事件
func TriggerReadyResourceAdd(readyResource *entity.ReadyResource, data map[string]interface{}) {
	if pluginApp != nil {
		if err := pluginApp.TriggerReadyResourceAdd(readyResource, data); err != nil {
			utils.Error("Failed to trigger ready resource add event: %v", err)
		} else {
			utils.Info("Ready resource add event triggered for: %s", readyResource.URL)
		}
	}
}

