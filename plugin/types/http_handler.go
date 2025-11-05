package types

import "github.com/gin-gonic/gin"

// HTTPHandler is an optional interface for plugins that want to register HTTP routes
type HTTPHandler interface {
	// RegisterRoutes registers HTTP routes for the plugin
	RegisterRoutes(group *gin.RouterGroup)
}

// TaskHandler is an optional interface for plugins that want to handle background tasks
type TaskHandler interface {
	// RegisterTaskProcessor registers a task processor for the plugin
	RegisterTaskProcessor() TaskProcessor
}

// TaskProcessor handles plugin-specific tasks
type TaskProcessor interface {
	// Process processes a task
	Process(taskData map[string]interface{}) error
	// Validate validates task data
	Validate(taskData map[string]interface{}) error
	// GetName returns the name of the task processor
	GetName() string
}