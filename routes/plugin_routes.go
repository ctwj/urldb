package routes

import (
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/handlers"
	"github.com/gin-gonic/gin"
)

// SetupPluginRoutes 设置插件管理路由
func SetupPluginRoutes(router *gin.Engine, repoManager *repo.RepositoryManager) {
	pluginHandler := handlers.NewPluginHandler(repoManager)

	// 插件管理路由组
	pluginGroup := router.Group("/api/plugins")
	{
		// 插件列表和详情
		pluginGroup.GET("", pluginHandler.GetPlugins)               // 获取插件列表
		pluginGroup.GET("/stats", pluginHandler.GetPluginStats)     // 获取插件统计
		pluginGroup.GET("/:name", pluginHandler.GetPlugin)          // 获取插件详情
		pluginGroup.GET("/:name/logs", pluginHandler.GetPluginLogs) // 获取插件日志

		// 插件控制
		pluginGroup.POST("/:name/enable", pluginHandler.EnablePlugin)   // 启用插件
		pluginGroup.POST("/:name/disable", pluginHandler.DisablePlugin) // 禁用插件

		// 插件配置
		pluginGroup.PUT("/:name/config", pluginHandler.UpdatePluginConfig) // 更新插件配置

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

		// 插件安装（未来扩展）
		pluginGroup.POST("/install", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Plugin installation coming soon",
			})
		})

		// 插件卸载（未来扩展）
		pluginGroup.DELETE("/:name", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Plugin uninstallation coming soon",
			})
		})
	}
}
