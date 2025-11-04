package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// HealthCheckResult 健康检查结果
type HealthCheckResult struct {
	PluginName    string        `json:"plugin_name"`
	Status        string        `json:"status"` // healthy, warning, critical, unknown
	HealthScore   float64       `json:"health_score"`
	ResponseTime  time.Duration `json:"response_time"`
	LastCheck     time.Time     `json:"last_check"`
	Details       interface{}   `json:"details,omitempty"`
	Error         string        `json:"error,omitempty"`
	Version       string        `json:"version"`
	Uptime        time.Duration `json:"uptime"`
	RestartCount  int           `json:"restart_count"`
}

// HealthChecker 健康检查器接口
type HealthChecker interface {
	Check(ctx context.Context, plugin types.Plugin) *HealthCheckResult
}

// DefaultHealthChecker 默认健康检查器
type DefaultHealthChecker struct {
	monitor *PluginMonitor
}

// NewDefaultHealthChecker 创建默认健康检查器
func NewDefaultHealthChecker(monitor *PluginMonitor) *DefaultHealthChecker {
	return &DefaultHealthChecker{
		monitor: monitor,
	}
}

// Check 执行健康检查
func (dhc *DefaultHealthChecker) Check(ctx context.Context, plugin types.Plugin) *HealthCheckResult {
	startTime := time.Now()
	result := &HealthCheckResult{
		PluginName:  plugin.Name(),
		Version:     plugin.Version(),
		LastCheck:   startTime,
		HealthScore: 100.0,
	}

	// 获取插件实例信息（如果可用）
	// 这里我们可以从插件管理器获取更多信息

	// 计算健康分数
	healthScore := dhc.monitor.GetPluginHealthScore(plugin.Name())
	result.HealthScore = healthScore

	// 根据健康分数设置状态
	if healthScore >= 90 {
		result.Status = "healthy"
	} else if healthScore >= 70 {
		result.Status = "warning"
	} else if healthScore >= 50 {
		result.Status = "critical"
	} else {
		result.Status = "unknown"
	}

	// 计算响应时间
	result.ResponseTime = time.Since(startTime)

	// 这里可以添加更具体的健康检查逻辑
	// 例如：检查插件是否响应、数据库连接是否正常等

	utils.Debug("Health check completed for plugin %s: status=%s, score=%.2f",
		plugin.Name(), result.Status, result.HealthScore)

	return result
}

// PluginHealthChecker 插件健康检查器
type PluginHealthChecker struct {
	checkers map[string]HealthChecker
	monitor  *PluginMonitor
}

// NewPluginHealthChecker 创建插件健康检查器
func NewPluginHealthChecker(monitor *PluginMonitor) *PluginHealthChecker {
	return &PluginHealthChecker{
		checkers: make(map[string]HealthChecker),
		monitor:  monitor,
	}
}

// RegisterChecker 注册健康检查器
func (phc *PluginHealthChecker) RegisterChecker(pluginName string, checker HealthChecker) {
	phc.checkers[pluginName] = checker
}

// Check 执行插件健康检查
func (phc *PluginHealthChecker) Check(ctx context.Context, plugin types.Plugin) *HealthCheckResult {
	// 检查是否有为该插件注册的特定检查器
	if checker, exists := phc.checkers[plugin.Name()]; exists {
		return checker.Check(ctx, plugin)
	}

	// 使用默认检查器
	defaultChecker := NewDefaultHealthChecker(phc.monitor)
	return defaultChecker.Check(ctx, plugin)
}

// BatchCheck 批量健康检查
func (phc *PluginHealthChecker) BatchCheck(ctx context.Context, plugins []types.Plugin) []*HealthCheckResult {
	results := make([]*HealthCheckResult, len(plugins))

	for i, plugin := range plugins {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			// 如果上下文已取消，标记剩余插件为未知状态
			for j := i; j < len(plugins); j++ {
				results[j] = &HealthCheckResult{
					PluginName: plugins[j].Name(),
					Status:     "unknown",
					Error:      "context cancelled",
					LastCheck:  time.Now(),
				}
			}
			return results
		default:
		}

		// 执行健康检查
		results[i] = phc.Check(ctx, plugin)
	}

	return results
}

// GetHealthStatusSummary 获取健康状态摘要
func (phc *PluginHealthChecker) GetHealthStatusSummary(results []*HealthCheckResult) map[string]int {
	summary := map[string]int{
		"healthy":   0,
		"warning":   0,
		"critical":  0,
		"unknown":   0,
		"total":     len(results),
	}

	for _, result := range results {
		summary[result.Status]++
	}

	return summary
}

// GetUnhealthyPlugins 获取不健康的插件
func (phc *PluginHealthChecker) GetUnhealthyPlugins(results []*HealthCheckResult) []*HealthCheckResult {
	var unhealthy []*HealthCheckResult
	for _, result := range results {
		if result.Status != "healthy" {
			unhealthy = append(unhealthy, result)
		}
	}
	return unhealthy
}

// HealthCheckReport 健康检查报告
type HealthCheckReport struct {
	Timestamp     time.Time              `json:"timestamp"`
	Summary       map[string]int         `json:"summary"`
	Results       []*HealthCheckResult   `json:"results"`
	Unhealthy     []*HealthCheckResult   `json:"unhealthy"`
	AverageScore  float64                `json:"average_score"`
	WorstPlugin   *HealthCheckResult     `json:"worst_plugin,omitempty"`
	BestPlugin    *HealthCheckResult     `json:"best_plugin,omitempty"`
}

// GenerateReport 生成健康检查报告
func (phc *PluginHealthChecker) GenerateReport(results []*HealthCheckResult) *HealthCheckReport {
	report := &HealthCheckReport{
		Timestamp: time.Now(),
		Summary:   phc.GetHealthStatusSummary(results),
		Results:   results,
		Unhealthy: phc.GetUnhealthyPlugins(results),
	}

	// 计算平均健康分数
	if len(results) > 0 {
		totalScore := 0.0
		for _, result := range results {
			totalScore += result.HealthScore
		}
		report.AverageScore = totalScore / float64(len(results))
	}

	// 找到最佳和最差的插件
	if len(results) > 0 {
		report.BestPlugin = results[0]
		report.WorstPlugin = results[0]
		for _, result := range results {
			if result.HealthScore > report.BestPlugin.HealthScore {
				report.BestPlugin = result
			}
			if result.HealthScore < report.WorstPlugin.HealthScore {
				report.WorstPlugin = result
			}
		}
	}

	return report
}

// CheckWithTimeout 带超时的健康检查
func (phc *PluginHealthChecker) CheckWithTimeout(plugin types.Plugin, timeout time.Duration) *HealthCheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan *HealthCheckResult, 1)

	go func() {
		result := phc.Check(ctx, plugin)
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		return &HealthCheckResult{
			PluginName: plugin.Name(),
			Status:     "unknown",
			Error:      fmt.Sprintf("health check timeout after %v", timeout),
			LastCheck:  time.Now(),
		}
	}
}