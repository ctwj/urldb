package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PluginMonitor 插件监控器
type PluginMonitor struct {
	// 插件健康指标
	pluginHealthGauge *prometheus.GaugeVec
	pluginUpTimeGauge *prometheus.GaugeVec
	pluginErrorCounter *prometheus.CounterVec
	pluginResponseTimeHistogram *prometheus.HistogramVec

	// 插件活动监控
	pluginActivities map[string][]ActivityRecord
	activityMutex sync.RWMutex

	// 告警配置
	alertRules map[string]AlertRule
	alertMutex sync.RWMutex

	// 告警通道
	alertChannel chan Alert

	// 监控配置
	config MonitorConfig

	// 上下文
	ctx context.Context
	cancel context.CancelFunc
}

// ActivityRecord 活动记录
type ActivityRecord struct {
	Timestamp   time.Time   `json:"timestamp"`
	PluginName  string      `json:"plugin_name"`
	Operation   string      `json:"operation"`
	ExecutionTime time.Duration `json:"execution_time"`
	Error       error       `json:"error,omitempty"`
	Details     interface{} `json:"details,omitempty"`
}

// AlertRule 告警规则
type AlertRule struct {
	Name        string        `json:"name"`
	PluginName  string        `json:"plugin_name"`
	Condition   string        `json:"condition"` // threshold, spike, error_rate, etc.
	Metric      string        `json:"metric"`
	Threshold   float64       `json:"threshold"`
	Severity    string        `json:"severity"` // low, medium, high, critical
	Description string        `json:"description"`
	Enabled     bool          `json:"enabled"`
	LastFired   *time.Time    `json:"last_fired,omitempty"`
	CoolDown    time.Duration `json:"cool_down"` // 冷却时间
}

// Alert 告警信息
type Alert struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	PluginName  string    `json:"plugin_name"`
	RuleName    string    `json:"rule_name"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Metric      string    `json:"metric"`
	Resolved    bool      `json:"resolved"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	MaxActivitiesPerPlugin int           `json:"max_activities_per_plugin"`
	AlertChannelCapacity   int           `json:"alert_channel_capacity"`
	HealthCheckInterval    time.Duration `json:"health_check_interval"`
	MetricsCollectionInterval time.Duration `json:"metrics_collection_interval"`
}

// NewPluginMonitor 创建新的插件监控器
func NewPluginMonitor() *PluginMonitor {
	ctx, cancel := context.WithCancel(context.Background())

	monitor := &PluginMonitor{
		pluginHealthGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "urldb",
				Subsystem: "plugin",
				Name:      "health_score",
				Help:      "Plugin health score (0-100)",
			},
			[]string{"plugin_name", "plugin_version"},
		),
		pluginUpTimeGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "urldb",
				Subsystem: "plugin",
				Name:      "uptime_seconds",
				Help:      "Plugin uptime in seconds",
			},
			[]string{"plugin_name", "plugin_version"},
		),
		pluginErrorCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "urldb",
				Subsystem: "plugin",
				Name:      "errors_total",
				Help:      "Total number of plugin errors",
			},
			[]string{"plugin_name", "plugin_version", "error_type"},
		),
		pluginResponseTimeHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "urldb",
				Subsystem: "plugin",
				Name:      "response_time_seconds",
				Help:      "Plugin response time in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"plugin_name", "plugin_version", "operation"},
		),
		pluginActivities: make(map[string][]ActivityRecord),
		alertRules:       make(map[string]AlertRule),
		alertChannel:     make(chan Alert, 100),
		config: MonitorConfig{
			MaxActivitiesPerPlugin: 1000,
			AlertChannelCapacity:   100,
			HealthCheckInterval:    30 * time.Second,
			MetricsCollectionInterval: 15 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 启动监控协程
	go monitor.startMonitoring()

	return monitor
}

// startMonitoring 启动监控协程
func (pm *PluginMonitor) startMonitoring() {
	healthTicker := time.NewTicker(pm.config.HealthCheckInterval)
	defer healthTicker.Stop()

	metricsTicker := time.NewTicker(pm.config.MetricsCollectionInterval)
	defer metricsTicker.Stop()

	for {
		select {
		case <-pm.ctx.Done():
			return
		case <-healthTicker.C:
			pm.checkHealth()
		case <-metricsTicker.C:
			pm.collectMetrics()
		}
	}
}

// checkHealth 检查插件健康状态
func (pm *PluginMonitor) checkHealth() {
	// 这里可以实现具体的健康检查逻辑
	// 检查插件响应时间、错误率等指标
}

// collectMetrics 收集监控指标
func (pm *PluginMonitor) collectMetrics() {
	// 定期检查告警规则
	pm.evaluateAlertRules()
}

// RecordActivity 记录插件活动
func (pm *PluginMonitor) RecordActivity(plugin types.Plugin, operation string, executionTime time.Duration, err error, details interface{}) {
	record := ActivityRecord{
		Timestamp:   time.Now(),
		PluginName:  plugin.Name(),
		Operation:   operation,
		ExecutionTime: executionTime,
		Error:       err,
		Details:     details,
	}

	pm.activityMutex.Lock()
	defer pm.activityMutex.Unlock()

	activities := pm.pluginActivities[plugin.Name()]
	if len(activities) >= pm.config.MaxActivitiesPerPlugin {
		// 移除最旧的活动记录
		activities = activities[1:]
	}
	activities = append(activities, record)
	pm.pluginActivities[plugin.Name()] = activities

	// 更新监控指标
	labels := prometheus.Labels{
		"plugin_name": plugin.Name(),
		"plugin_version": plugin.Version(),
	}

	pm.pluginResponseTimeHistogram.With(labels).Observe(executionTime.Seconds())

	// 如果有错误，增加错误计数器
	if err != nil {
		pm.pluginErrorCounter.WithLabelValues(plugin.Name(), plugin.Version(), "execution_error").Inc()
	}

	// 更新插件健康分数
	pm.updateHealthScore(plugin)
}

// updateHealthScore 更新插件健康分数
func (pm *PluginMonitor) updateHealthScore(plugin types.Plugin) {
	labels := prometheus.Labels{
		"plugin_name": plugin.Name(),
		"plugin_version": plugin.Version(),
	}

	// 计算健康分数（示例：基于错误率和响应时间）
	healthScore := pm.calculateHealthScore(plugin)
	pm.pluginHealthGauge.With(labels).Set(healthScore)
}

// calculateHealthScore 计算插件健康分数
func (pm *PluginMonitor) calculateHealthScore(plugin types.Plugin) float64 {
	pm.activityMutex.RLock()
	defer pm.activityMutex.RUnlock()

	activities := pm.pluginActivities[plugin.Name()]
	if len(activities) == 0 {
		return 100.0 // 没有活动记录时默认健康
	}

	// 计算最近一段时间内的错误率和平均响应时间
	recentStart := time.Now().Add(-5 * time.Minute)
	errorCount := 0
	totalCount := 0
	totalResponseTime := time.Duration(0)

	for _, activity := range activities {
		if activity.Timestamp.After(recentStart) {
			totalCount++
			if activity.Error != nil {
				errorCount++
			}
			totalResponseTime += activity.ExecutionTime
		}
	}

	if totalCount == 0 {
		return 100.0
	}

	// 计算错误率（0-100）
	errorRate := float64(errorCount) / float64(totalCount) * 100

	// 计算平均响应时间（毫秒）
	avgResponseTime := float64(totalResponseTime.Milliseconds()) / float64(totalCount)

	// 基于错误率和响应时间计算健康分数
	healthScore := 100.0
	if errorRate > 0 {
		// 错误率越高，健康分数越低
		healthScore -= errorRate * 0.8
	}
	if avgResponseTime > 1000 { // 超过1秒
		// 响应时间越长，健康分数越低
		healthScore -= (avgResponseTime / 1000) * 0.2
	}

	// 确保健康分数在0-100之间
	if healthScore < 0 {
		healthScore = 0
	}
	if healthScore > 100 {
		healthScore = 100
	}

	return healthScore
}

// SetAlertRule 设置告警规则
func (pm *PluginMonitor) SetAlertRule(rule AlertRule) error {
	pm.alertMutex.Lock()
	defer pm.alertMutex.Unlock()

	if rule.Name == "" {
		return fmt.Errorf("alert rule name cannot be empty")
	}

	pm.alertRules[rule.Name] = rule
	utils.Info("Alert rule set: %s for plugin %s", rule.Name, rule.PluginName)
	return nil
}

// RemoveAlertRule 移除告警规则
func (pm *PluginMonitor) RemoveAlertRule(ruleName string) {
	pm.alertMutex.Lock()
	defer pm.alertMutex.Unlock()

	delete(pm.alertRules, ruleName)
	utils.Info("Alert rule removed: %s", ruleName)
}

// evaluateAlertRules 评估告警规则
func (pm *PluginMonitor) evaluateAlertRules() {
	pm.alertMutex.RLock()
	defer pm.alertMutex.RUnlock()

	now := time.Now()

	for _, rule := range pm.alertRules {
		if !rule.Enabled {
			continue
		}

		// 检查是否在冷却期内
		if rule.LastFired != nil && now.Sub(*rule.LastFired) < rule.CoolDown {
			continue
		}

		// 评估规则
		triggered, value, threshold := pm.evaluateRule(rule)
		if triggered {
			alert := Alert{
				ID:          fmt.Sprintf("alert_%s_%d", rule.Name, now.UnixNano()),
				Timestamp:   now,
				PluginName:  rule.PluginName,
				RuleName:    rule.Name,
				Severity:    rule.Severity,
				Message:     rule.Description,
				Value:       value,
				Threshold:   threshold,
				Metric:      rule.Metric,
				Resolved:    false,
			}

			// 更新最后触发时间
			rule.LastFired = &now
			pm.alertRules[rule.Name] = rule

			// 发送告警
			select {
			case pm.alertChannel <- alert:
				utils.Warn("Alert triggered: %s for plugin %s (value: %.2f, threshold: %.2f)",
					rule.Name, rule.PluginName, value, threshold)
			default:
				utils.Warn("Alert channel is full, alert dropped: %s", rule.Name)
			}
		}
	}
}

// evaluateRule 评估单个规则
func (pm *PluginMonitor) evaluateRule(rule AlertRule) (bool, float64, float64) {
	switch rule.Metric {
	case "error_rate":
		return pm.evaluateErrorRate(rule)
	case "response_time":
		return pm.evaluateResponseTime(rule)
	case "health_score":
		return pm.evaluateHealthScore(rule)
	default:
		return false, 0, 0
	}
}

// evaluateErrorRate 评估错误率
func (pm *PluginMonitor) evaluateErrorRate(rule AlertRule) (bool, float64, float64) {
	pm.activityMutex.RLock()
	defer pm.activityMutex.RUnlock()

	activities := pm.pluginActivities[rule.PluginName]
	if len(activities) == 0 {
		return false, 0, rule.Threshold
	}

	recentStart := time.Now().Add(-1 * time.Minute) // 最近1分钟
	errorCount := 0
	totalCount := 0

	for _, activity := range activities {
		if activity.Timestamp.After(recentStart) {
			totalCount++
			if activity.Error != nil {
				errorCount++
			}
		}
	}

	if totalCount == 0 {
		return false, 0, rule.Threshold
	}

	errorRate := float64(errorCount) / float64(totalCount) * 100
	return errorRate > rule.Threshold, errorRate, rule.Threshold
}

// evaluateResponseTime 评估响应时间
func (pm *PluginMonitor) evaluateResponseTime(rule AlertRule) (bool, float64, float64) {
	pm.activityMutex.RLock()
	defer pm.activityMutex.RUnlock()

	activities := pm.pluginActivities[rule.PluginName]
	if len(activities) == 0 {
		return false, 0, rule.Threshold
	}

	recentStart := time.Now().Add(-1 * time.Minute) // 最近1分钟
	totalResponseTime := time.Duration(0)
	count := 0

	for _, activity := range activities {
		if activity.Timestamp.After(recentStart) {
			totalResponseTime += activity.ExecutionTime
			count++
		}
	}

	if count == 0 {
		return false, 0, rule.Threshold
	}

	avgResponseTime := float64(totalResponseTime.Milliseconds()) / float64(count)
	return avgResponseTime > rule.Threshold, avgResponseTime, rule.Threshold
}

// evaluateHealthScore 评估健康分数
func (pm *PluginMonitor) evaluateHealthScore(rule AlertRule) (bool, float64, float64) {
	healthScore := pm.calculateHealthScoreForPlugin(rule.PluginName)
	return healthScore < rule.Threshold, healthScore, rule.Threshold
}

// calculateHealthScoreForPlugin 为特定插件计算健康分数
func (pm *PluginMonitor) calculateHealthScoreForPlugin(pluginName string) float64 {
	pm.activityMutex.RLock()
	defer pm.activityMutex.RUnlock()

	activities := pm.pluginActivities[pluginName]
	if len(activities) == 0 {
		return 100.0
	}

	// 计算最近一段时间内的错误率和平均响应时间
	recentStart := time.Now().Add(-5 * time.Minute)
	errorCount := 0
	totalCount := 0
	totalResponseTime := time.Duration(0)

	for _, activity := range activities {
		if activity.Timestamp.After(recentStart) {
			totalCount++
			if activity.Error != nil {
				errorCount++
			}
			totalResponseTime += activity.ExecutionTime
		}
	}

	if totalCount == 0 {
		return 100.0
	}

	// 计算错误率（0-100）
	errorRate := float64(errorCount) / float64(totalCount) * 100

	// 计算平均响应时间（毫秒）
	avgResponseTime := float64(totalResponseTime.Milliseconds()) / float64(totalCount)

	// 基于错误率和响应时间计算健康分数
	healthScore := 100.0
	if errorRate > 0 {
		// 错误率越高，健康分数越低
		healthScore -= errorRate * 0.8
	}
	if avgResponseTime > 1000 { // 超过1秒
		// 响应时间越长，健康分数越低
		healthScore -= (avgResponseTime / 1000) * 0.2
	}

	// 确保健康分数在0-100之间
	if healthScore < 0 {
		healthScore = 0
	}
	if healthScore > 100 {
		healthScore = 100
	}

	return healthScore
}

// GetActivities 获取插件活动记录
func (pm *PluginMonitor) GetActivities(pluginName string, limit int) []ActivityRecord {
	pm.activityMutex.RLock()
	defer pm.activityMutex.RUnlock()

	activities := pm.pluginActivities[pluginName]
	if limit <= 0 || limit >= len(activities) {
		result := make([]ActivityRecord, len(activities))
		copy(result, activities)
		return result
	}

	start := len(activities) - limit
	result := make([]ActivityRecord, limit)
	copy(result, activities[start:])
	return result
}

// GetAlerts 获取告警通道
func (pm *PluginMonitor) GetAlerts() <-chan Alert {
	return pm.alertChannel
}

// GetAlertRules 获取告警规则
func (pm *PluginMonitor) GetAlertRules() map[string]AlertRule {
	pm.alertMutex.RLock()
	defer pm.alertMutex.RUnlock()

	result := make(map[string]AlertRule)
	for k, v := range pm.alertRules {
		result[k] = v
	}
	return result
}

// GetPluginHealthScore 获取插件健康分数
func (pm *PluginMonitor) GetPluginHealthScore(pluginName string) float64 {
	return pm.calculateHealthScoreForPlugin(pluginName)
}

// Stop 停止监控器
func (pm *PluginMonitor) Stop() {
	pm.cancel()
}

// GetMonitorStats 获取监控统计信息
func (pm *PluginMonitor) GetMonitorStats() map[string]interface{} {
	pm.activityMutex.RLock()
	defer pm.activityMutex.RUnlock()

	totalActivities := 0
	for _, activities := range pm.pluginActivities {
		totalActivities += len(activities)
	}

	return map[string]interface{}{
		"total_plugins_monitored": len(pm.pluginActivities),
		"total_activities":        totalActivities,
		"alert_rules_count":       len(pm.alertRules),
		"config":                  pm.config,
		"timestamp":               time.Now(),
	}
}