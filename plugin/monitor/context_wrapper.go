package monitor

import (
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// MonitoredPluginContext 插件上下文装饰器，用于自动记录执行时间和错误
type MonitoredPluginContext struct {
	originalContext types.PluginContext
	pluginInstance  interface{} // 使用interface{}避免循环导入
	pluginMonitor   *PluginMonitor
}

// NewMonitoredPluginContext 创建新的监控插件上下文
func NewMonitoredPluginContext(original types.PluginContext, instance interface{}, monitor *PluginMonitor) *MonitoredPluginContext {
	return &MonitoredPluginContext{
		originalContext: original,
		pluginInstance:  instance,
		pluginMonitor:   monitor,
	}
}

// LogDebug 记录调试日志
func (m *MonitoredPluginContext) LogDebug(msg string, args ...interface{}) {
	m.originalContext.LogDebug(msg, args...)
}

// LogInfo 记录信息日志
func (m *MonitoredPluginContext) LogInfo(msg string, args ...interface{}) {
	m.originalContext.LogInfo(msg, args...)
}

// LogWarn 记录警告日志
func (m *MonitoredPluginContext) LogWarn(msg string, args ...interface{}) {
	m.originalContext.LogWarn(msg, args...)
}

// LogError 记录错误日志
func (m *MonitoredPluginContext) LogError(msg string, args ...interface{}) {
	m.originalContext.LogError(msg, args...)
}

// GetConfig 获取配置
func (m *MonitoredPluginContext) GetConfig(key string) (interface{}, error) {
	return m.originalContext.GetConfig(key)
}

// SetConfig 设置配置
func (m *MonitoredPluginContext) SetConfig(key string, value interface{}) error {
	return m.originalContext.SetConfig(key, value)
}

// GetData 获取数据
func (m *MonitoredPluginContext) GetData(key string, dataType string) (interface{}, error) {
	return m.originalContext.GetData(key, dataType)
}

// SetData 设置数据
func (m *MonitoredPluginContext) SetData(key string, value interface{}, dataType string) error {
	return m.originalContext.SetData(key, value, dataType)
}

// DeleteData 删除数据
func (m *MonitoredPluginContext) DeleteData(key string, dataType string) error {
	return m.originalContext.DeleteData(key, dataType)
}

// RegisterTask 注册任务
func (m *MonitoredPluginContext) RegisterTask(name string, task func()) error {
	return m.originalContext.RegisterTask(name, task)
}

// UnregisterTask 注销任务
func (m *MonitoredPluginContext) UnregisterTask(name string) error {
	return m.originalContext.UnregisterTask(name)
}

// GetDB 获取数据库连接
func (m *MonitoredPluginContext) GetDB() interface{} {
	return m.originalContext.GetDB()
}

// CheckPermission 检查权限
func (m *MonitoredPluginContext) CheckPermission(permissionType string, resource ...string) (bool, error) {
	return m.originalContext.CheckPermission(permissionType, resource...)
}

// RequestPermission 请求权限
func (m *MonitoredPluginContext) RequestPermission(permissionType string, resource string) error {
	return m.originalContext.RequestPermission(permissionType, resource)
}

// GetSecurityReport 获取安全报告
func (m *MonitoredPluginContext) GetSecurityReport() (interface{}, error) {
	return m.originalContext.GetSecurityReport()
}

// ExecuteWithMonitoring 执行带监控的函数
func (m *MonitoredPluginContext) ExecuteWithMonitoring(operation string, fn func() error) error {
	startTime := time.Now()

	// 执行函数
	err := fn()

	// 计算执行时间
	executionTime := time.Since(startTime)

	// 更新插件实例统计信息
	if m.pluginInstance != nil {
		// 使用类型断言来调用方法
		if instance, ok := m.pluginInstance.(interface{ UpdateExecutionStats(time.Duration, error) }); ok {
			instance.UpdateExecutionStats(executionTime, err)
		}
	}

	// 记录活动到监控器
	// 注意：这里我们假设插件实例包含插件信息
	// 在实际实现中，可能需要通过其他方式获取插件信息

	return err
}

// RecordActivity 记录插件活动
func (m *MonitoredPluginContext) RecordActivity(operation string, executionTime time.Duration, err error, details interface{}) {
	// 记录活动到监控器
	// 这里需要获取插件信息，可能需要通过其他方式传递

	// 更新插件实例统计信息
	if m.pluginInstance != nil {
		// 使用类型断言来调用方法
		if instance, ok := m.pluginInstance.(interface{ UpdateExecutionStats(time.Duration, error) }); ok {
			instance.UpdateExecutionStats(executionTime, err)
		}
	}
}