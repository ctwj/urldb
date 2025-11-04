package debug

import (
	"fmt"
	"strings"
	"time"

	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// Debugger 插件调试器
type Debugger struct {
	manager   *manager.Manager
	collector *EventCollector
	tracer    *Tracer
	config    DebugConfig
}

// NewDebugger 创建新的调试器
func NewDebugger(mgr *manager.Manager, config DebugConfig, db *gorm.DB) *Debugger {
	collector := NewEventCollector(config, db)
	tracer := NewTracer(collector)

	debugger := &Debugger{
		manager:   mgr,
		collector: collector,
		tracer:    tracer,
		config:    config,
	}

	// 注册插件生命周期事件监听器
	debugger.registerPluginEventListeners()

	return debugger
}

// Start 开始调试
func (d *Debugger) Start() {
	if !d.config.Enabled {
		utils.Info("Debugger is not enabled")
		return
	}

	utils.Info("Plugin debugger started")
	d.collector.AddEvent(DebugEvent{
		ID:         fmt.Sprintf("debug_start_%d", time.Now().UnixNano()),
		PluginName: "debugger",
		Timestamp:  time.Now(),
		EventType:  EventTypePluginStart,
		Level:      LevelInfo,
		Message:    "Debugger started",
	})
}

// Stop 停止调试
func (d *Debugger) Stop() {
	utils.Info("Plugin debugger stopped")
	d.collector.AddEvent(DebugEvent{
		ID:         fmt.Sprintf("debug_stop_%d", time.Now().UnixNano()),
		PluginName: "debugger",
		Timestamp:  time.Now(),
		EventType:  EventTypePluginStop,
		Level:      LevelInfo,
		Message:    "Debugger stopped",
	})

	if d.collector.fileWriter != nil {
		d.collector.fileWriter.Close()
	}
}

// TracePlugin 跟踪插件
func (d *Debugger) TracePlugin(pluginName string) error {
	if !d.config.Enabled {
		return fmt.Errorf("debugger is not enabled")
	}

	// 检查插件是否存在
	plugin, err := d.manager.GetPlugin(pluginName)
	if err != nil {
		return fmt.Errorf("plugin %s not found: %v", pluginName, err)
	}

	// 开始调试会话
	sessionID := fmt.Sprintf("trace_%s_%d", pluginName, time.Now().UnixNano())
	d.collector.StartSession(pluginName, sessionID)

	d.tracer.LogInfo(pluginName, "Started tracing plugin")
	return nil
}

// UntracePlugin 停止跟踪插件
func (d *Debugger) UntracePlugin(pluginName string) {
	d.tracer.LogInfo(pluginName, "Stopped tracing plugin")
}

// QueryEvents 查询调试事件
func (d *Debugger) QueryEvents(query DebugQuery) (*DebugQueryResponse, error) {
	return d.collector.QueryEvents(query)
}

// GetPluginStats 获取插件统计信息
func (d *Debugger) GetPluginStats(pluginName string) map[string]interface{} {
	return d.tracer.GetPluginStats(pluginName)
}

// GetStats 获取总体统计信息
func (d *Debugger) GetStats() *DebugStats {
	return d.collector.GetStats()
}

// GetPluginTraces 获取插件的跟踪信息
func (d *Debugger) GetPluginTraces(pluginName string) ([]DebugEvent, error) {
	query := DebugQuery{
		Filter: DebugFilter{
			PluginName: pluginName,
			Limit:      100, // 默认获取100个事件
		},
		SortBy:  "timestamp",
		SortDir: "desc",
	}

	result, err := d.collector.QueryEvents(query)
	if err != nil {
		return nil, err
	}

	return result.Events, nil
}

// ExportEvents 导出事件
func (d *Debugger) ExportEvents(filter DebugFilter, format string) ([]byte, error) {
	query := DebugQuery{
		Filter: filter,
		SortBy: "timestamp",
		SortDir: "asc",
	}

	result, err := d.collector.QueryEvents(query)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(format) {
	case "json":
		return d.exportToJSON(result)
	case "text", "txt":
		return d.exportToText(result)
	case "csv":
		return d.exportToCSV(result)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToJSON 导出为JSON格式
func (d *Debugger) exportToJSON(result *DebugQueryResponse) ([]byte, error) {
	// 使用Go的encoding/json包来序列化
	var builder strings.Builder
	builder.WriteString("[")

	for i, event := range result.Events {
		if i > 0 {
			builder.WriteString(",\n")
		}

		// 简单的JSON序列化实现
		builder.WriteString(fmt.Sprintf(`{"id":"%s","plugin_name":"%s","timestamp":"%s","event_type":"%s","level":"%s","message":"%s"`,
			event.ID, event.PluginName, event.Timestamp.Format(time.RFC3339), event.EventType, event.Level, event.Message))

		if len(event.Details) > 0 {
			builder.WriteString(`,"details":{`)
			first := true
			for k, v := range event.Details {
				if !first {
					builder.WriteString(",")
				}
				builder.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
				first = false
			}
			builder.WriteString("}")
		}

		builder.WriteString("}")
	}

	builder.WriteString("]")
	return []byte(builder.String()), nil
}

// exportToText 导出为文本格式
func (d *Debugger) exportToText(result *DebugQueryResponse) ([]byte, error) {
	var builder strings.Builder

	for _, event := range result.Events {
		builder.WriteString(d.formatEventAsText(event))
		builder.WriteString("\n")
	}

	return []byte(builder.String()), nil
}

// exportToCSV 导出为CSV格式
func (d *Debugger) exportToCSV(result *DebugQueryResponse) ([]byte, error) {
	var builder strings.Builder

	// CSV头部
	builder.WriteString("ID,PluginName,Timestamp,EventType,Level,Message\n")

	for _, event := range result.Events {
		builder.WriteString(fmt.Sprintf(`"%s","%s","%s","%s","%s","%s"`+"\n",
			event.ID, event.PluginName, event.Timestamp.Format(time.RFC3339),
			event.EventType, event.Level, event.Message))
	}

	return []byte(builder.String()), nil
}

// formatEventAsText 格式化事件为文本
func (d *Debugger) formatEventAsText(event DebugEvent) string {
	var builder strings.Builder

	builder.WriteString(event.Timestamp.Format("2006-01-02 15:04:05.000"))
	builder.WriteString(" [")
	builder.WriteString(string(event.Level))
	builder.WriteString("] [")
	builder.WriteString(event.PluginName)
	builder.WriteString("] ")
	builder.WriteString(string(event.EventType))
	builder.WriteString(": ")
	builder.WriteString(event.Message)

	if len(event.Details) > 0 {
		builder.WriteString(" {")
		first := true
		for k, v := range event.Details {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(v)
			first = false
		}
		builder.WriteString("}")
	}

	return builder.String()
}

// ClearEvents 清除事件
func (d *Debugger) ClearEvents() {
	d.collector.ClearEvents()
}

// UpdateConfig 更新配置
func (d *Debugger) UpdateConfig(config DebugConfig) {
	d.config = config
	d.collector.UpdateConfig(config)
}

// GetActiveSessions 获取活动会话
func (d *Debugger) GetActiveSessions() []*DebugSession {
	var sessions []*DebugSession
	for _, session := range d.collector.sessions {
		if session.Status == "active" {
			// 复制会话信息以避免并发问题
			sessionCopy := &DebugSession{
				ID:         session.ID,
				PluginName: session.PluginName,
				StartTime:  session.StartTime,
				EndTime:    session.EndTime,
				Status:     session.Status,
				Events:     make([]DebugEvent, len(session.Events)),
			}
			copy(sessionCopy.Events, session.Events)
			sessions = append(sessions, sessionCopy)
		}
	}
	return sessions
}

// registerPluginEventListeners 注册插件事件监听器
func (d *Debugger) registerPluginEventListeners() {
	// 这里可以注册钩子函数来监听插件事件
	// 由于Go的插件系统限制，可能需要在插件管理器中添加回调机制
	// 这里提供一个概念性的实现
}

// TraceFunction 跟踪函数执行
func (d *Debugger) TraceFunction(pluginName, functionName string, fn func() error) error {
	return d.tracer.TraceFunction(pluginName, functionName, fn)
}

// TraceFunctionWithResult 跟踪函数执行并返回结果
func (d *Debugger) TraceFunctionWithResult(pluginName, functionName string, fn func() (interface{}, error)) (interface{}, error) {
	return d.tracer.TraceFunctionWithResult(pluginName, functionName, fn)
}

// TraceDataAccess 跟踪数据访问
func (d *Debugger) TraceDataAccess(pluginName, operation, dataType, key string, fn func() error) error {
	return d.tracer.TraceDataAccess(pluginName, operation, dataType, key, fn)
}

// TraceConfigChange 跟踪配置变更
func (d *Debugger) TraceConfigChange(pluginName, key string, oldValue, newValue interface{}) {
	d.tracer.TraceConfigChange(pluginName, key, oldValue, newValue)
}

// TraceTaskExecution 跟踪任务执行
func (d *Debugger) TraceTaskExecution(pluginName, taskName string, fn func() error) error {
	return d.tracer.TraceTaskExecution(pluginName, taskName, fn)
}