package main

import (
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/demo"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
)

func main() {
	// Initialize logger
	utils.InitLogger(nil)

	// Initialize database
	if err := db.InitDB(); err != nil {
		utils.Fatal("Failed to initialize database: %v", err)
	}

	// Create repository manager
	repoManager := repo.NewRepositoryManager(db.DB)

	// Create task manager
	taskManager := task.NewTaskManager(repoManager)

	// Initialize plugin system
	plugin.InitPluginSystem(taskManager, repoManager)

	// Register demo plugins
	demoPlugin := demo.NewFullDemoPlugin()
	if err := plugin.RegisterPlugin(demoPlugin); err != nil {
		utils.Error("Failed to register demo plugin: %v", err)
		return
	}

	securityDemoPlugin := demo.NewSecurityDemoPlugin()
	if err := plugin.RegisterPlugin(securityDemoPlugin); err != nil {
		utils.Error("Failed to register security demo plugin: %v", err)
		return
	}

	// Initialize plugin
	config := map[string]interface{}{
		"interval": 30, // 30 seconds
		"enabled":  true,
	}
	if err := plugin.GetManager().InitializePlugin("full-demo-plugin", config); err != nil {
		utils.Error("Failed to initialize demo plugin: %v", err)
		return
	}

	// Start plugin
	if err := plugin.GetManager().StartPlugin("full-demo-plugin"); err != nil {
		utils.Error("Failed to start demo plugin: %v", err)
		return
	}

	// Keep the application running
	utils.Info("Plugin system test started. Running for 2 minutes...")
	time.Sleep(2 * time.Minute)

	// Stop plugin
	if err := plugin.GetManager().StopPlugin("full-demo-plugin"); err != nil {
		utils.Error("Failed to stop demo plugin: %v", err)
		return
	}

	utils.Info("Plugin system test completed successfully.")
}