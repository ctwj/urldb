package main

import (
	"time"

	"github.com/ctwj/urldb/plugin/demo"
	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/utils"
)

func main() {
	// Initialize logger
	utils.InitLogger(nil)

	// Create plugin manager
	pluginManager := manager.NewManager(nil) // In a real implementation, pass the actual task manager

	// Create and register demo plugin
	demoPlugin := demo.NewDemoPlugin()
	if err := pluginManager.RegisterPlugin(demoPlugin); err != nil {
		utils.Error("Failed to register demo plugin: %v", err)
		return
	}

	// Initialize plugin
	config := map[string]interface{}{
		"interval": 60, // 60 seconds
	}
	if err := pluginManager.InitializePlugin("demo-plugin", config); err != nil {
		utils.Error("Failed to initialize demo plugin: %v", err)
		return
	}

	// Start plugin
	if err := pluginManager.StartPlugin("demo-plugin"); err != nil {
		utils.Error("Failed to start demo plugin: %v", err)
		return
	}

	// Keep the application running
	utils.Info("Demo plugin system started. Press Ctrl+C to exit.")
	time.Sleep(5 * time.Minute) // Run for 5 minutes for demonstration

	// Stop plugin
	if err := pluginManager.StopPlugin("demo-plugin"); err != nil {
		utils.Error("Failed to stop demo plugin: %v", err)
		return
	}

	utils.Info("Demo plugin system stopped.")
}