package routes

import (
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/handlers"
	"github.com/ctwj/urldb/plugin"
	"github.com/gin-gonic/gin"
)

// SetupPluginRoutes 设置插件管理路由
func SetupPluginRoutes(router *gin.Engine, repoManager *repo.RepositoryManager, pluginManager *plugin.Manager) {
	pluginHandler := handlers.NewPluginHandler(repoManager, pluginManager)

	// 插件管理路由组
	pluginGroup := router.Group("/api/plugins")
	{
		// 插件列表和详情
		pluginGroup.GET("", pluginHandler.GetPlugins)                    // 获取插件列表（钩子插件）
		pluginGroup.GET("/installed", pluginHandler.GetInstalledPlugins)  // 获取已安装的插件
		pluginGroup.GET("/stats", pluginHandler.GetPluginStats)           // 获取插件统计
		pluginGroup.GET("/:name", pluginHandler.GetPlugin)                // 获取插件详情
		pluginGroup.GET("/:name/logs", pluginHandler.GetPluginLogs)       // 获取插件日志

		// 插件安装和卸载
		pluginGroup.POST("/install", pluginHandler.InstallPlugin)          // 安装插件
		pluginGroup.DELETE("/:name", pluginHandler.UninstallPlugin)        // 卸载插件

		// 插件加载和卸载
		pluginGroup.POST("/:name/load", pluginHandler.LoadPlugin)          // 加载插件
		pluginGroup.POST("/:name/unload", pluginHandler.UnloadPlugin)      // 卸载已加载的插件
		pluginGroup.POST("/:name/reload", pluginHandler.ReloadPlugin)       // 重新加载插件

		// 插件控制
		pluginGroup.POST("/:name/enable", pluginHandler.EnablePlugin)      // 启用插件
		pluginGroup.POST("/:name/disable", pluginHandler.DisablePlugin)    // 禁用插件

		// 插件配置
		pluginGroup.PUT("/:name/config", pluginHandler.UpdatePluginConfig)  // 更新插件配置

		// 插件市场（未来扩展）
		pluginGroup.GET("/market", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"data": gin.H{
					"message": "Plugin market coming soon",
					"plugins": []interface{}{},
				},
			})
		})
	}
}
