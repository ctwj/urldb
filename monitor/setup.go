package monitor

import (
	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
)

// SetupMonitoring 设置完整的监控系统
func SetupMonitoring(router *gin.Engine) {
	// 获取全局监控实例
	metrics := GetGlobalMetrics()

	// 设置健康检查端点
	metrics.SetupHealthCheck(router)

	// 设置指标端点
	router.GET("/metrics", metrics.MetricsHandler())

	utils.Info("监控系统已设置完成")
}

// SetGlobalErrorHandler 设置全局错误处理器
func SetGlobalErrorHandler(eh *ErrorHandler) {
	globalErrorHandler = eh
}