package debug

import (
	"context"
	"fmt"
	"time"

	"github.com/ctwj/urldb/utils"
)

// Tracer 调试追踪器
type Tracer struct {
	collector *EventCollector
}

// NewTracer 创建新的追踪器
func NewTracer(collector *EventCollector) *Tracer {
	return &Tracer{
		collector: collector,
	}
}

// TraceFunction 跟踪函数执行
func (t *Tracer) TraceFunction(pluginName, functionName string, fn func() error) error {
	startTime := time.Now()
	correlationID := fmt.Sprintf("fn_%s_%s_%d", pluginName, functionName, startTime.UnixNano())

	// 记录函数开始
	t.collector.AddEventWithDetails(pluginName, EventTypeFunctionCall, LevelDebug,
		fmt.Sprintf("Calling function %s", functionName), map[string]string{
			"function": functionName,
			"correlation_id": correlationID,
		})

	// 执行函数
	err := fn()

	// 记录函数结束
	duration := time.Since(startTime)
	level := LevelDebug
	if err != nil {
		level = LevelError
	}

	t.collector.AddEventWithDetails(pluginName, EventTypeFunctionReturn, level,
		fmt.Sprintf("Function %s completed", functionName), map[string]string{
			"function": functionName,
			"duration": duration.String(),
			"correlation_id": correlationID,
			"error":  fmt.Sprintf("%v", err),
		})

	return err
}

// TraceFunctionWithResult 跟踪函数执行并返回结果
func (t *Tracer) TraceFunctionWithResult(pluginName, functionName string, fn func() (interface{}, error)) (interface{}, error) {
	startTime := time.Now()
	correlationID := fmt.Sprintf("fn_%s_%s_%d", pluginName, functionName, startTime.UnixNano())

	// 记录函数开始
	t.collector.AddEventWithDetails(pluginName, EventTypeFunctionCall, LevelDebug,
		fmt.Sprintf("Calling function %s", functionName), map[string]string{
			"function": functionName,
			"correlation_id": correlationID,
		})

	// 执行函数
	result, err := fn()

	// 记录函数结束
	duration := time.Since(startTime)
	level := LevelDebug
	if err != nil {
		level = LevelError
	}

	t.collector.AddEventWithDetails(pluginName, EventTypeFunctionReturn, level,
		fmt.Sprintf("Function %s completed", functionName), map[string]string{
			"function": functionName,
			"duration": duration.String(),
			"correlation_id": correlationID,
			"error":  fmt.Sprintf("%v", err),
		})

	return result, err
}

// WithTraceContext 创建带追踪上下文的函数
func (t *Tracer) WithTraceContext(ctx context.Context, pluginName, operation string) (context.Context, func(error)) {
	startTime := time.Now()
	correlationID := fmt.Sprintf("ctx_%s_%s_%d", pluginName, operation, startTime.UnixNano())

	// 记录操作开始
	t.collector.AddEventWithDetails(pluginName, EventTypeFunctionCall, LevelDebug,
		fmt.Sprintf("Starting operation %s", operation), map[string]string{
			"operation": operation,
			"correlation_id": correlationID,
		})

	// 将correlationID添加到上下文
	ctx = context.WithValue(ctx, "correlation_id", correlationID)

	return ctx, func(err error) {
		// 记录操作结束
		duration := time.Since(startTime)
		level := LevelDebug
		if err != nil {
			level = LevelError
		}

		t.collector.AddEventWithDetails(pluginName, EventTypeFunctionReturn, level,
			fmt.Sprintf("Operation %s completed", operation), map[string]string{
				"operation": operation,
				"duration": duration.String(),
				"correlation_id": correlationID,
				"error":  fmt.Sprintf("%v", err),
			})
	}
}

// TraceDataAccess 跟踪数据访问
func (t *Tracer) TraceDataAccess(pluginName, operation, dataType, key string, fn func() error) error {
	startTime := time.Now()
	correlationID := fmt.Sprintf("data_%s_%s_%s_%d", pluginName, operation, key, startTime.UnixNano())

	// 记录数据访问开始
	t.collector.AddEventWithDetails(pluginName, EventTypeDataAccess, LevelDebug,
		fmt.Sprintf("Accessing %s data: %s", dataType, key), map[string]string{
			"operation": operation,
			"data_type": dataType,
			"data_key":  key,
			"correlation_id": correlationID,
		})

	// 执行数据访问
	err := fn()

	// 记录数据访问结束
	duration := time.Since(startTime)
	level := LevelDebug
	if err != nil {
		level = LevelError
	}

	t.collector.AddEventWithDetails(pluginName, EventTypeDataAccess, level,
		fmt.Sprintf("%s data access completed: %s", operation, key), map[string]string{
			"operation": operation,
			"data_type": dataType,
			"data_key":  key,
			"duration": duration.String(),
			"correlation_id": correlationID,
			"error":  fmt.Sprintf("%v", err),
		})

	return err
}

// TraceConfigChange 跟踪配置变更
func (t *Tracer) TraceConfigChange(pluginName, key string, oldValue, newValue interface{}) {
	correlationID := fmt.Sprintf("config_%s_%s_%d", pluginName, key, time.Now().UnixNano())

	t.collector.AddEventWithDetails(pluginName, EventTypeConfigChange, LevelInfo,
		fmt.Sprintf("Configuration changed: %s", key), map[string]string{
			"key": key,
			"old_value": fmt.Sprintf("%v", oldValue),
			"new_value": fmt.Sprintf("%v", newValue),
			"correlation_id": correlationID,
		})
}

// TraceTaskExecution 跟踪任务执行
func (t *Tracer) TraceTaskExecution(pluginName, taskName string, fn func() error) error {
	startTime := time.Now()
	correlationID := fmt.Sprintf("task_%s_%s_%d", pluginName, taskName, startTime.UnixNano())

	// 记录任务开始
	t.collector.AddEventWithDetails(pluginName, EventTypeTaskExecute, LevelInfo,
		fmt.Sprintf("Starting task: %s", taskName), map[string]string{
			"task_name": taskName,
			"correlation_id": correlationID,
		})

	// 执行任务
	err := fn()

	// 记录任务结束
	duration := time.Since(startTime)
	level := LevelInfo
	if err != nil {
		level = LevelError
	}

	t.collector.AddEventWithDetails(pluginName, EventTypeTaskComplete, level,
		fmt.Sprintf("Task completed: %s", taskName), map[string]string{
			"task_name": taskName,
			"duration": duration.String(),
			"correlation_id": correlationID,
			"error":  fmt.Sprintf("%v", err),
		})

	return err
}

// LogPluginLifecycle 记录插件生命周期事件
func (t *Tracer) LogPluginLifecycle(pluginName string, eventType DebugEventType, message string, details map[string]string) {
	if details == nil {
		details = make(map[string]string)
	}
	details["plugin_name"] = pluginName

	level := LevelInfo
	if eventType == EventTypePluginError {
		level = LevelError
	}

	t.collector.AddEventWithDetails(pluginName, eventType, level, message, details)
}

// LogError 记录错误
func (t *Tracer) LogError(pluginName, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	t.collector.AddEventWithDetails(pluginName, EventTypePluginError, LevelError, formattedMessage, nil)
}

// LogWarn 记录警告
func (t *Tracer) LogWarn(pluginName, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	t.collector.AddEventWithDetails(pluginName, EventTypePluginError, LevelWarn, formattedMessage, nil)
}

// LogInfo 记录信息
func (t *Tracer) LogInfo(pluginName, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	t.collector.AddEventWithDetails(pluginName, EventTypePluginError, LevelInfo, formattedMessage, nil)
}

// LogDebug 记录调试信息
func (t *Tracer) LogDebug(pluginName, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	t.collector.AddEventWithDetails(pluginName, EventTypePluginError, LevelDebug, formattedMessage, nil)
}

// GetPluginStats 获取插件统计信息
func (t *Tracer) GetPluginStats(pluginName string) map[string]interface{} {
	stats := t.collector.GetStats()

	return map[string]interface{}{
		"total_events":   stats.EventsByPlugin[pluginName],
		"error_count":    stats.EventsByLevel[string(LevelError)] + stats.EventsByLevel[string(LevelCritical)],
		"last_event_time": stats.LastEventTime,
		"recent_events":  t.getRecentPluginEvents(pluginName, 10),
	}
}

// getRecentPluginEvents 获取插件最近的事件
func (t *Tracer) getRecentPluginEvents(pluginName string, count int) []DebugEvent {
	query := DebugQuery{
		Filter: DebugFilter{
			PluginName: pluginName,
			Limit:      count,
		},
		SortBy:  "timestamp",
		SortDir: "desc",
	}

	result, err := t.collector.QueryEvents(query)
	if err != nil {
		utils.Error("Failed to get recent plugin events: %v", err)
		return nil
	}

	return result.Events
}