package monitor

import (
	"context"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// EnhancedPluginContext 增强的插件上下文，用于自动记录执行时间和错误
type EnhancedPluginContext struct {
	originalContext types.PluginContext
	pluginInstance  interface{} // 使用interface{}避免循环导入
	pluginMonitor   *PluginMonitor
	pluginName      string
}

// NewEnhancedPluginContext 创建新的增强插件上下文
func NewEnhancedPluginContext(original types.PluginContext, instance interface{}, monitor *PluginMonitor, pluginName string) *EnhancedPluginContext {
	return &EnhancedPluginContext{
		originalContext: original,
		pluginInstance:  instance,
		pluginMonitor:   monitor,
		pluginName:      pluginName,
	}
}

// ExecuteWithMonitoring 执行带监控的函数
func (e *EnhancedPluginContext) ExecuteWithMonitoring(operation string, fn func() error) error {
	startTime := time.Now()

	// 执行函数
	err := fn()

	// 计算执行时间
	executionTime := time.Since(startTime)

	// 更新插件实例统计信息
	if e.pluginInstance != nil {
		// 使用类型断言来调用方法
		if instance, ok := e.pluginInstance.(interface{ UpdateExecutionStats(time.Duration, error) }); ok {
			instance.UpdateExecutionStats(executionTime, err)
		}
	}

	// 记录活动到监控器
	if e.pluginMonitor != nil {
		// 这里需要获取插件对象，但在上下文中可能无法直接访问
		// 我们可以创建一个虚拟的插件对象用于记录
		// 在实际实现中，可能需要通过其他方式传递插件信息
	}

	return err
}

// LogDebug 记录调试日志
func (e *EnhancedPluginContext) LogDebug(msg string, args ...interface{}) {
	e.originalContext.LogDebug(msg, args...)
}

// LogInfo 记录信息日志
func (e *EnhancedPluginContext) LogInfo(msg string, args ...interface{}) {
	e.originalContext.LogInfo(msg, args...)
}

// LogWarn 记录警告日志
func (e *EnhancedPluginContext) LogWarn(msg string, args ...interface{}) {
	e.originalContext.LogWarn(msg, args...)
}

// LogError 记录错误日志
func (e *EnhancedPluginContext) LogError(msg string, args ...interface{}) {
	e.originalContext.LogError(msg, args...)
	// 记录错误到监控器
	if e.pluginMonitor != nil {
		// 这里可以记录错误到监控器
	}
}

// GetConfig 获取配置
func (e *EnhancedPluginContext) GetConfig(key string) (interface{}, error) {
	return e.originalContext.GetConfig(key)
}

// SetConfig 设置配置
func (e *EnhancedPluginContext) SetConfig(key string, value interface{}) error {
	return e.originalContext.SetConfig(key, value)
}

// GetData 获取数据
func (e *EnhancedPluginContext) GetData(key string, dataType string) (interface{}, error) {
	return e.originalContext.GetData(key, dataType)
}

// SetData 设置数据
func (e *EnhancedPluginContext) SetData(key string, value interface{}, dataType string) error {
	return e.originalContext.SetData(key, value, dataType)
}

// DeleteData 删除数据
func (e *EnhancedPluginContext) DeleteData(key string, dataType string) error {
	return e.originalContext.DeleteData(key, dataType)
}

// RegisterTask 注册任务
func (e *EnhancedPluginContext) RegisterTask(name string, task func()) error {
	return e.originalContext.RegisterTask(name, task)
}

// UnregisterTask 注销任务
func (e *EnhancedPluginContext) UnregisterTask(name string) error {
	return e.originalContext.UnregisterTask(name)
}

// GetDB 获取数据库连接
func (e *EnhancedPluginContext) GetDB() interface{} {
	return e.originalContext.GetDB()
}

// CheckPermission 检查权限
func (e *EnhancedPluginContext) CheckPermission(permissionType string, resource ...string) (bool, error) {
	return e.originalContext.CheckPermission(permissionType, resource...)
}

// RequestPermission 请求权限
func (e *EnhancedPluginContext) RequestPermission(permissionType string, resource string) error {
	return e.originalContext.RequestPermission(permissionType, resource)
}

// GetSecurityReport 获取安全报告
func (e *EnhancedPluginContext) GetSecurityReport() (interface{}, error) {
	return e.originalContext.GetSecurityReport()
}

// CacheSet 设置缓存项
func (e *EnhancedPluginContext) CacheSet(key string, value interface{}, ttl time.Duration) error {
	return e.originalContext.CacheSet(key, value, ttl)
}

// CacheGet 获取缓存项
func (e *EnhancedPluginContext) CacheGet(key string) (interface{}, error) {
	return e.originalContext.CacheGet(key)
}

// CacheDelete 删除缓存项
func (e *EnhancedPluginContext) CacheDelete(key string) error {
	return e.originalContext.CacheDelete(key)
}

// CacheClear 清除所有缓存项
func (e *EnhancedPluginContext) CacheClear() error {
	return e.originalContext.CacheClear()
}

// ConcurrencyExecute 执行带并发控制的任务
func (e *EnhancedPluginContext) ConcurrencyExecute(ctx context.Context, taskFunc func() error) error {
	return e.originalContext.ConcurrencyExecute(ctx, taskFunc)
}

// SetConcurrencyLimit 设置并发限制
func (e *EnhancedPluginContext) SetConcurrencyLimit(limit int) error {
	return e.originalContext.SetConcurrencyLimit(limit)
}

// GetConcurrencyStats 获取并发统计信息
func (e *EnhancedPluginContext) GetConcurrencyStats() (map[string]interface{}, error) {
	return e.originalContext.GetConcurrencyStats()
}