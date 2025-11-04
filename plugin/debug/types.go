package debug

import (
	"time"
)

// DebugEvent 调试事件
type DebugEvent struct {
	ID          string            `json:"id"`
	PluginName  string            `json:"plugin_name"`
	Timestamp   time.Time         `json:"timestamp"`
	EventType   DebugEventType    `json:"event_type"`
	Level       DebugLevel        `json:"level"`
	Message     string            `json:"message"`
	Details     map[string]string `json:"details,omitempty"`
	StackTrace  string            `json:"stack_trace,omitempty"`
	Duration    time.Duration     `json:"duration,omitempty"`
	Correlation string            `json:"correlation,omitempty"` // 用于关联相关事件
}

// DebugEventType 调试事件类型
type DebugEventType string

const (
	EventTypePluginLoad     DebugEventType = "plugin_load"
	EventTypePluginUnload   DebugEventType = "plugin_unload"
	EventTypePluginStart    DebugEventType = "plugin_start"
	EventTypePluginStop     DebugEventType = "plugin_stop"
	EventTypePluginError    DebugEventType = "plugin_error"
	EventTypeFunctionCall   DebugEventType = "function_call"
	EventTypeFunctionReturn DebugEventType = "function_return"
	EventTypeDataAccess     DebugEventType = "data_access"
	EventTypeConfigChange   DebugEventType = "config_change"
	EventTypeTaskExecute    DebugEventType = "task_execute"
	EventTypeTaskComplete   DebugEventType = "task_complete"
	EventTypeNetworkRequest DebugEventType = "network_request"
	EventTypeNetworkResponse DebugEventType = "network_response"
)

// DebugLevel 调试级别
type DebugLevel string

const (
	LevelTrace    DebugLevel = "trace"
	LevelDebug    DebugLevel = "debug"
	LevelInfo     DebugLevel = "info"
	LevelWarn     DebugLevel = "warn"
	LevelError    DebugLevel = "error"
	LevelCritical DebugLevel = "critical"
)

// DebugSession 调试会话
type DebugSession struct {
	ID         string    `json:"id"`
	PluginName string    `json:"plugin_name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time,omitempty"`
	Status     string    `json:"status"` // active, completed, error
	Events     []DebugEvent `json:"events"`
}

// DebugFilter 调试过滤器
type DebugFilter struct {
	PluginName  string         `json:"plugin_name,omitempty"`
	EventType   DebugEventType `json:"event_type,omitempty"`
	Level       DebugLevel     `json:"level,omitempty"`
	StartTime   time.Time      `json:"start_time,omitempty"`
	EndTime     time.Time      `json:"end_time,omitempty"`
	Contains    string         `json:"contains,omitempty"` // 消息包含的文本
	Correlation string         `json:"correlation,omitempty"`
	Limit       int            `json:"limit,omitempty"`
	Offset      int            `json:"offset,omitempty"`
}

// DebugConfig 调试配置
type DebugConfig struct {
	Enabled        bool          `json:"enabled"`
	Level          DebugLevel    `json:"level"`
	MaxEvents      int           `json:"max_events"`      // 最大事件数量
	MaxEventAge    time.Duration `json:"max_event_age"`   // 最大事件保存时间
	BufferSize     int           `json:"buffer_size"`     // 缓冲区大小
	FlushInterval  time.Duration `json:"flush_interval"`  // 刷新间隔
	OutputToFile   bool          `json:"output_to_file"`  // 是否输出到文件
	OutputFilePath string        `json:"output_file_path"` // 输出文件路径
	IncludeDetails bool          `json:"include_details"`  // 是否包含详细信息
	IncludeStack   bool          `json:"include_stack"`    // 是否包含堆栈信息
}

// DebugStats 调试统计信息
type DebugStats struct {
	TotalEvents     int64         `json:"total_events"`
	EventsByType    map[string]int64 `json:"events_by_type"`
	EventsByLevel   map[string]int64 `json:"events_by_level"`
	EventsByPlugin  map[string]int64 `json:"events_by_plugin"`
	StartTime       time.Time     `json:"start_time"`
	LastEventTime   time.Time     `json:"last_event_time"`
	AvgEventsPerMin float64       `json:"avg_events_per_min"`
	ErrorRate       float64       `json:"error_rate"`
}

// DebugQuery 查询调试信息的请求
type DebugQuery struct {
	Filter DebugFilter `json:"filter"`
	SortBy string      `json:"sort_by"` // timestamp, level, type
	SortDir string     `json:"sort_dir"` // asc, desc
}

// DebugQueryResponse 查询调试信息的响应
type DebugQueryResponse struct {
	Events     []DebugEvent `json:"events"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}