package debug

import (
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// EventCollector 事件收集器
type EventCollector struct {
	config     DebugConfig
	events     []DebugEvent
	sessions   map[string]*DebugSession
	mutex      sync.RWMutex
	database   *gorm.DB
	fileWriter *DebugFileWriter
	stats      *DebugStats
}

// NewEventCollector 创建新的事件收集器
func NewEventCollector(config DebugConfig, db *gorm.DB) *EventCollector {
	collector := &EventCollector{
		config:   config,
		events:   make([]DebugEvent, 0, config.BufferSize),
		sessions: make(map[string]*DebugSession),
		database: db,
		stats: &DebugStats{
			EventsByType:   make(map[string]int64),
			EventsByLevel:  make(map[string]int64),
			EventsByPlugin: make(map[string]int64),
			StartTime:      time.Now(),
		},
	}

	// 初始化文件写入器
	if config.OutputToFile && config.OutputFilePath != "" {
		collector.fileWriter = NewDebugFileWriter(config.OutputFilePath)
	}

	// 启动定期清理和刷新
	if config.FlushInterval > 0 {
		go collector.startFlushTicker()
	}

	return collector
}

// AddEvent 添加调试事件
func (ec *EventCollector) AddEvent(event DebugEvent) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	// 检查是否启用调试
	if !ec.config.Enabled {
		return
	}

	// 检查级别过滤
	if !ec.isLevelEnabled(event.Level) {
		return
	}

	// 如果达到最大事件数，移除最旧的事件
	if ec.config.MaxEvents > 0 && len(ec.events) >= ec.config.MaxEvents {
		ec.events = ec.events[1:] // 移除第一个元素
	}

	// 添加事件
	ec.events = append(ec.events, event)

	// 更新统计信息
	ec.updateStats(event)

	// 输出到文件
	if ec.fileWriter != nil {
		ec.fileWriter.WriteEvent(event)
	}

	// 添加到会话
	if event.Correlation != "" {
		if session, exists := ec.sessions[event.Correlation]; exists {
			session.Events = append(session.Events, event)
		}
	}

	utils.Debug("Debug event added: %s - %s", event.EventType, event.Message)
}

// AddEventWithDetails 添加带详细信息的调试事件
func (ec *EventCollector) AddEventWithDetails(pluginName string, eventType DebugEventType, level DebugLevel, message string, details map[string]string) {
	event := DebugEvent{
		ID:         fmt.Sprintf("%s_%d", pluginName, time.Now().UnixNano()),
		PluginName: pluginName,
		Timestamp:  time.Now(),
		EventType:  eventType,
		Level:      level,
		Message:    message,
		Details:    details,
	}

	ec.AddEvent(event)
}

// StartSession 开始调试会话
func (ec *EventCollector) StartSession(pluginName, sessionID string) *DebugSession {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	session := &DebugSession{
		ID:         sessionID,
		PluginName: pluginName,
		StartTime:  time.Now(),
		Status:     "active",
		Events:     make([]DebugEvent, 0),
	}

	ec.sessions[sessionID] = session
	return session
}

// EndSession 结束调试会话
func (ec *EventCollector) EndSession(sessionID string, hasError bool) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	if session, exists := ec.sessions[sessionID]; exists {
		session.EndTime = time.Now()
		if hasError {
			session.Status = "error"
		} else {
			session.Status = "completed"
		}
	}
}

// GetSession 获取调试会话
func (ec *EventCollector) GetSession(sessionID string) *DebugSession {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	if session, exists := ec.sessions[sessionID]; exists {
		// 返回副本以避免并发问题
		sessionCopy := &DebugSession{
			ID:         session.ID,
			PluginName: session.PluginName,
			StartTime:  session.StartTime,
			EndTime:    session.EndTime,
			Status:     session.Status,
			Events:     make([]DebugEvent, len(session.Events)),
		}
		copy(sessionCopy.Events, session.Events)
		return sessionCopy
	}

	return nil
}

// QueryEvents 查询事件
func (ec *EventCollector) QueryEvents(query DebugQuery) (*DebugQueryResponse, error) {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	// 应用过滤器
	filteredEvents := ec.applyFilter(query.Filter)

	// 排序
	sortedEvents := ec.sortEvents(filteredEvents, query.SortBy, query.SortDir)

	// 分页
	total := int64(len(sortedEvents))
	pageSize := query.Filter.Limit
	if pageSize <= 0 {
		pageSize = 50 // 默认页面大小
	}

	page := 1
	if query.Filter.Offset > 0 {
		page = query.Filter.Offset/pageSize + 1
	}

	start := query.Filter.Offset
	if start >= len(sortedEvents) {
		start = len(sortedEvents)
	}

	end := start + pageSize
	if end > len(sortedEvents) {
		end = len(sortedEvents)
	}

	pagedEvents := sortedEvents[start:end]

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &DebugQueryResponse{
		Events:     pagedEvents,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetStats 获取统计信息
func (ec *EventCollector) GetStats() *DebugStats {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	// 创建统计信息副本
	statsCopy := &DebugStats{
		TotalEvents:     ec.stats.TotalEvents,
		EventsByType:    make(map[string]int64),
		EventsByLevel:   make(map[string]int64),
		EventsByPlugin:  make(map[string]int64),
		StartTime:       ec.stats.StartTime,
		LastEventTime:   ec.stats.LastEventTime,
		AvgEventsPerMin: ec.stats.AvgEventsPerMin,
		ErrorRate:       ec.stats.ErrorRate,
	}

	// 复制映射
	for k, v := range ec.stats.EventsByType {
		statsCopy.EventsByType[k] = v
	}
	for k, v := range ec.stats.EventsByLevel {
		statsCopy.EventsByLevel[k] = v
	}
	for k, v := range ec.stats.EventsByPlugin {
		statsCopy.EventsByPlugin[k] = v
	}

	return statsCopy
}

// ClearEvents 清除所有事件
func (ec *EventCollector) ClearEvents() {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	ec.events = make([]DebugEvent, 0, ec.config.BufferSize)
	ec.stats = &DebugStats{
		EventsByType:   make(map[string]int64),
		EventsByLevel:  make(map[string]int64),
		EventsByPlugin: make(map[string]int64),
		StartTime:      time.Now(),
	}
}

// UpdateConfig 更新配置
func (ec *EventCollector) UpdateConfig(config DebugConfig) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	ec.config = config

	// 更新文件写入器
	if config.OutputToFile && config.OutputFilePath != "" {
		if ec.fileWriter == nil || ec.fileWriter.filePath != config.OutputFilePath {
			ec.fileWriter = NewDebugFileWriter(config.OutputFilePath)
		}
	} else {
		ec.fileWriter = nil
	}
}

// isLevelEnabled 检查级别是否启用
func (ec *EventCollector) isLevelEnabled(level DebugLevel) bool {
	levelOrder := map[DebugLevel]int{
		LevelTrace:    0,
		LevelDebug:    1,
		LevelInfo:     2,
		LevelWarn:     3,
		LevelError:    4,
		LevelCritical: 5,
	}

	currentLevelOrder := levelOrder[ec.config.Level]
	eventLevelOrder := levelOrder[level]

	return eventLevelOrder >= currentLevelOrder
}

// applyFilter 应用过滤器
func (ec *EventCollector) applyFilter(filter DebugFilter) []DebugEvent {
	var result []DebugEvent

	for _, event := range ec.events {
		// 应用过滤器条件
		if filter.PluginName != "" && event.PluginName != filter.PluginName {
			continue
		}

		if filter.EventType != "" && event.EventType != filter.EventType {
			continue
		}

		if filter.Level != "" && event.Level != filter.Level {
			continue
		}

		if !filter.StartTime.IsZero() && event.Timestamp.Before(filter.StartTime) {
			continue
		}

		if !filter.EndTime.IsZero() && event.Timestamp.After(filter.EndTime) {
			continue
		}

		if filter.Contains != "" && !contains(event.Message, filter.Contains) {
			continue
		}

		if filter.Correlation != "" && event.Correlation != filter.Correlation {
			continue
		}

		result = append(result, event)
	}

	return result
}

// sortEvents 排序事件
func (ec *EventCollector) sortEvents(events []DebugEvent, sortBy, sortDir string) []DebugEvent {
	// 简化实现，实际可以根据sortBy参数进行不同字段的排序
	// 这里默认按时间戳排序
	// 排序逻辑可以根据需要扩展
	return events
}

// updateStats 更新统计信息
func (ec *EventCollector) updateStats(event DebugEvent) {
	ec.stats.TotalEvents++
	ec.stats.EventsByType[string(event.EventType)]++
	ec.stats.EventsByLevel[string(event.Level)]++
	ec.stats.EventsByPlugin[event.PluginName]++
	ec.stats.LastEventTime = event.Timestamp

	// 计算平均事件数/分钟
	duration := time.Since(ec.stats.StartTime)
	if duration > 0 {
		ec.stats.AvgEventsPerMin = float64(ec.stats.TotalEvents) / duration.Minutes()
	}

	// 计算错误率
	if ec.stats.TotalEvents > 0 {
		errorCount := ec.stats.EventsByLevel[string(LevelError)] + ec.stats.EventsByLevel[string(LevelCritical)]
		ec.stats.ErrorRate = float64(errorCount) / float64(ec.stats.TotalEvents) * 100
	}
}

// startFlushTicker 启动刷新定时器
func (ec *EventCollector) startFlushTicker() {
	ticker := time.NewTicker(ec.config.FlushInterval)
	defer ticker.Stop()

	for range ticker.C {
		ec.mutex.Lock()
		// 清理过期事件
		ec.cleanupExpiredEvents()
		ec.mutex.Unlock()
	}
}

// cleanupExpiredEvents 清理过期事件
func (ec *EventCollector) cleanupExpiredEvents() {
	if ec.config.MaxEventAge <= 0 {
		return
	}

	expirationTime := time.Now().Add(-ec.config.MaxEventAge)
	var newEvents []DebugEvent

	for _, event := range ec.events {
		if event.Timestamp.After(expirationTime) {
			newEvents = append(newEvents, event)
		}
	}

	ec.events = newEvents
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(s) == len(substr) && s == substr ||
		 len(s) > len(substr) && (s == substr ||
			strings.Contains(strings.ToLower(s), strings.ToLower(substr))))
}