package types

import (
	"context"
	"time"

	"github.com/ctwj/urldb/plugin/config"
)

// Plugin is the interface that all plugins must implement
type Plugin interface {
	// Basic information
	Name() string
	Version() string
	Description() string
	Author() string

	// Lifecycle
	Initialize(ctx PluginContext) error
	Start() error
	Stop() error
	Cleanup() error

	// Dependencies
	// Dependencies returns a list of plugin names that this plugin depends on
	Dependencies() []string
	// CheckDependencies checks the status of plugin dependencies and returns a map
	// where the key is the dependency name and the value indicates if it's satisfied
	CheckDependencies() map[string]bool
}

// ConfigurablePlugin is an optional interface for plugins that support configuration schemas
type ConfigurablePlugin interface {
	// CreateConfigSchema creates the configuration schema for the plugin
	CreateConfigSchema() *config.ConfigSchema

	// CreateConfigTemplate creates a default configuration template
	CreateConfigTemplate() *config.ConfigTemplate
}

// PluginInstance represents a running plugin instance
type PluginInstance struct {
	Plugin     Plugin
	Context    PluginContext
	Status     PluginStatus
	Config     map[string]interface{}
	StartTime  time.Time
	StopTime   *time.Time
	LastError  error
	RestartCount int
	HealthScore  float64
	// 监控相关字段
	TotalExecutionTime time.Duration
	TotalExecutions    int64
	TotalErrors        int64
	LastExecutionTime  time.Time
}

// UpdateExecutionStats 更新插件执行统计信息
func (instance *PluginInstance) UpdateExecutionStats(executionTime time.Duration, err error) {
	instance.TotalExecutionTime += executionTime
	instance.TotalExecutions++
	instance.LastExecutionTime = time.Now()

	if err != nil {
		instance.TotalErrors++
		instance.LastError = err
	}

	// 更新健康分数
	instance.updateHealthScore()
}

// updateHealthScore 更新插件健康分数
func (instance *PluginInstance) updateHealthScore() {
	// 简单的健康分数计算算法
	// 基于错误率和平均执行时间
	if instance.TotalExecutions == 0 {
		instance.HealthScore = 100.0
		return
	}

	// 计算错误率 (0-100)
	errorRate := float64(instance.TotalErrors) / float64(instance.TotalExecutions) * 100

	// 计算平均执行时间 (毫秒)
	avgExecutionTime := float64(instance.TotalExecutionTime.Milliseconds()) / float64(instance.TotalExecutions)

	// 基于错误率和执行时间计算健康分数
	healthScore := 100.0

	// 错误率影响 (权重 0.7)
	if errorRate > 0 {
		healthScore -= errorRate * 0.7
	}

	// 执行时间影响 (权重 0.3)
	// 假设超过1000ms为不健康
	if avgExecutionTime > 1000 {
		healthScore -= (avgExecutionTime / 1000) * 0.3
	}

	// 确保健康分数在0-100之间
	if healthScore < 0 {
		healthScore = 0
	}
	if healthScore > 100 {
		healthScore = 100
	}

	instance.HealthScore = healthScore
}

// GetInstanceInfo 获取插件实例信息
func (instance *PluginInstance) GetInstanceInfo() *PluginInstance {
	return &PluginInstance{
		Plugin:            instance.Plugin,
		Context:           instance.Context,
		Status:            instance.Status,
		Config:            instance.Config,
		StartTime:         instance.StartTime,
		StopTime:          instance.StopTime,
		LastError:         instance.LastError,
		RestartCount:      instance.RestartCount,
		HealthScore:       instance.HealthScore,
		TotalExecutionTime: instance.TotalExecutionTime,
		TotalExecutions:   instance.TotalExecutions,
		TotalErrors:       instance.TotalErrors,
		LastExecutionTime: instance.LastExecutionTime,
	}
}

// PluginContext provides the context for plugins to interact with the system
type PluginContext interface {
	// Logging
	LogDebug(msg string, args ...interface{})
	LogInfo(msg string, args ...interface{})
	LogWarn(msg string, args ...interface{})
	LogError(msg string, args ...interface{})

	// Configuration
	GetConfig(key string) (interface{}, error)
	SetConfig(key string, value interface{}) error

	// Data
	GetData(key string, dataType string) (interface{}, error)
	SetData(key string, value interface{}, dataType string) error
	DeleteData(key string, dataType string) error

	// Task scheduling
	RegisterTask(name string, task func()) error
	UnregisterTask(name string) error

	// Database access
	GetDB() interface{} // Returns *gorm.DB

	// Security
	// CheckPermission checks if the plugin has the specified permission
	CheckPermission(permissionType string, resource ...string) (bool, error)
	// RequestPermission requests a permission for the plugin
	RequestPermission(permissionType string, resource string) error
	// GetSecurityReport returns a security report for the plugin
	GetSecurityReport() (interface{}, error)

	// Cache
	// CacheSet sets a cache item
	CacheSet(key string, value interface{}, ttl time.Duration) error
	// CacheGet gets a cache item
	CacheGet(key string) (interface{}, error)
	// CacheDelete deletes a cache item
	CacheDelete(key string) error
	// CacheClear clears all cache items for the plugin
	CacheClear() error

	// Concurrency
	// ConcurrencyExecute executes a task with concurrency control
	ConcurrencyExecute(ctx context.Context, taskFunc func() error) error
	// SetConcurrencyLimit sets the concurrency limit for the plugin
	SetConcurrencyLimit(limit int) error
	// GetConcurrencyStats gets concurrency control statistics
	GetConcurrencyStats() (map[string]interface{}, error)
}

// PluginStatus represents the current status of a plugin
type PluginStatus string

const (
	StatusRegistered  PluginStatus = "registered"
	StatusInitialized PluginStatus = "initialized"
	StatusStarting    PluginStatus = "starting"
	StatusRunning     PluginStatus = "running"
	StatusStopping    PluginStatus = "stopping"
	StatusStopped     PluginStatus = "stopped"
	StatusError       PluginStatus = "error"
	StatusDisabled    PluginStatus = "disabled"
)

// PluginInfo contains information about a plugin
type PluginInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Status      PluginStatus      `json:"status"`
	Dependencies []string         `json:"dependencies"`
	LastError   error             `json:"last_error,omitempty"`
	StartTime   time.Time         `json:"start_time,omitempty"`
	StopTime    *time.Time        `json:"stop_time,omitempty"`
	RestartCount int              `json:"restart_count"`
	HealthScore  float64          `json:"health_score"`
	// 监控相关字段
	TotalExecutionTime time.Duration `json:"total_execution_time,omitempty"`
	TotalExecutions    int64         `json:"total_executions,omitempty"`
	TotalErrors        int64         `json:"total_errors,omitempty"`
	LastExecutionTime  time.Time     `json:"last_execution_time,omitempty"`
}