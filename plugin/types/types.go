package types

import (
	"context"
	"time"
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
	Dependencies() []string
	CheckDependencies() map[string]bool
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
}