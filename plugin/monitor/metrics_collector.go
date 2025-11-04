package monitor

import (
	"time"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
	"github.com/prometheus/client_golang/prometheus"
)

// PluginMetricsCollector 插件指标收集器
type PluginMetricsCollector struct {
	monitor *PluginMonitor
}

// NewPluginMetricsCollector 创建插件指标收集器
func NewPluginMetricsCollector(monitor *PluginMonitor) *PluginMetricsCollector {
	return &PluginMetricsCollector{
		monitor: monitor,
	}
}

// CollectPluginMetrics 收集插件指标
func (pmc *PluginMetricsCollector) CollectPluginMetrics(plugin types.Plugin, instance *types.PluginInstance) {
	if instance == nil {
		utils.Warn("Plugin instance is nil for plugin %s", plugin.Name())
		return
	}

	labels := prometheus.Labels{
		"plugin_name":    plugin.Name(),
		"plugin_version": plugin.Version(),
	}

	// 更新插件状态指标
	statusGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "urldb",
		Subsystem:   "plugin",
		Name:        "status",
		Help:        "Plugin status (1=running, 0=stopped)",
		ConstLabels: labels,
	})

	// 根据插件状态设置指标值
	if instance.Status == types.StatusRunning {
		statusGauge.Set(1)
	} else {
		statusGauge.Set(0)
	}

	// 更新重启次数指标
	restartCounter := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "urldb",
		Subsystem:   "plugin",
		Name:        "restart_count",
		Help:        "Plugin restart count",
		ConstLabels: labels,
	})
	restartCounter.Set(float64(instance.RestartCount))

	// 更新运行时间指标
	if instance.StartTime.IsZero() {
		// 插件未启动
		uptimeGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   "urldb",
			Subsystem:   "plugin",
			Name:        "uptime_seconds",
			Help:        "Plugin uptime in seconds",
			ConstLabels: labels,
		})
		uptimeGauge.Set(0)
	} else {
		uptime := time.Since(instance.StartTime).Seconds()
		pmc.monitor.pluginUpTimeGauge.With(labels).Set(uptime)
	}

	utils.Debug("Collected metrics for plugin %s", plugin.Name())
}

// CollectAllPluginsMetrics 收集所有插件指标
func (pmc *PluginMetricsCollector) CollectAllPluginsMetrics(plugins map[string]types.Plugin, instances map[string]*types.PluginInstance) {
	for name, plugin := range plugins {
		instance := instances[name]
		pmc.CollectPluginMetrics(plugin, instance)
	}
}

// GetPluginMetrics 获取插件指标摘要
func (pmc *PluginMetricsCollector) GetPluginMetrics(pluginName string) map[string]interface{} {
	// 由于Prometheus指标不能直接读取，我们只能返回一些基本的统计信息
	// 在实际应用中，这些指标会通过/metrics端点暴露给Prometheus

	return map[string]interface{}{
		"plugin_name": pluginName,
		"timestamp":   time.Now(),
		"info":        "Use /metrics endpoint to get detailed metrics",
	}
}

// GetPluginsMetricsSummary 获取所有插件指标摘要
func (pmc *PluginMetricsCollector) GetPluginsMetricsSummary(plugins []types.PluginInfo) map[string]interface{} {
	runningCount := 0
	stoppedCount := 0
	errorCount := 0
	totalCount := len(plugins)

	for _, plugin := range plugins {
		switch plugin.Status {
		case types.StatusRunning:
			runningCount++
		case types.StatusStopped, types.StatusDisabled:
			stoppedCount++
		case types.StatusError:
			errorCount++
		}
	}

	return map[string]interface{}{
		"total_plugins":  totalCount,
		"running":        runningCount,
		"stopped":        stoppedCount,
		"error":          errorCount,
		"timestamp":      time.Now(),
		"health_summary": "Use health check endpoint for detailed health information",
	}
}