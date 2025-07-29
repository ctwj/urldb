package middleware

import (
	"net/http"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/gin-gonic/gin"
)

var repoManager *repo.RepositoryManager

// SetRepositoryManager 设置Repository管理器
func SetRepositoryManager(rm *repo.RepositoryManager) {
	repoManager = rm
}

// PublicAPIAuth 公开API认证中间件
func PublicAPIAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取API Token
		apiToken := c.GetHeader("X-API-Token")
		if apiToken == "" {
			// 尝试从查询参数获取
			apiToken = c.Query("api_token")
		}

		if apiToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "缺少API Token",
				"code":    401,
			})
			c.Abort()
			return
		}

		// 验证API Token
		if repoManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "系统未初始化",
				"code":    500,
			})
			c.Abort()
			return
		}

		// 验证API Token
		apiTokenConfig, err := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyApiToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "系统配置获取失败",
				"code":    500,
			})
			c.Abort()
			return
		}

		if apiTokenConfig == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "API Token未配置",
				"code":    503,
			})
			c.Abort()
			return
		}

		if apiTokenConfig != apiToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "API Token无效",
				"code":    401,
			})
			c.Abort()
			return
		}

		// 检查维护模式
		maintenanceMode, err := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyMaintenanceMode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "系统配置获取失败",
				"code":    500,
			})
			c.Abort()
			return
		}

		if maintenanceMode {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "系统维护中，请稍后再试",
				"code":    503,
			})
			c.Abort()
			return
		}

		// 验证通过，继续处理
		c.Next()
	}
}
