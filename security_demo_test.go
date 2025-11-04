package main

import (
	"fmt"
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

	// Register security demo plugin
	securityDemoPlugin := demo.NewSecurityDemoPlugin()
	if err := plugin.RegisterPlugin(securityDemoPlugin); err != nil {
		utils.Error("Failed to register security demo plugin: %v", err)
		return
	}

	// Initialize plugin
	config := map[string]interface{}{
		"test_mode": true,
		"duration":  30, // 30 seconds
	}
	if err := plugin.GetManager().InitializePlugin("security_demo", config); err != nil {
		utils.Error("Failed to initialize security demo plugin: %v", err)
		return
	}

	// Start plugin
	if err := plugin.GetManager().StartPlugin("security_demo"); err != nil {
		utils.Error("Failed to start security demo plugin: %v", err)
		return
	}

	// Keep the application running
	utils.Info("Security demo plugin test started. Running for 1 minute...")
	time.Sleep(1 * time.Minute)

	// Get security report
	manager := plugin.GetManager()
	if manager != nil && manager.GetSecurityManager() != nil {
		securityMgr := manager.GetSecurityManager()
		report := securityMgr.CreateSecurityReport("security_demo")
		fmt.Printf("Security Report for security_demo:\n")
		fmt.Printf("  Security Score: %.2f\n", report.SecurityScore)
		fmt.Printf("  Permissions: %d\n", len(report.Permissions))
		fmt.Printf("  Activities: %d\n", len(report.Activities))
		fmt.Printf("  Alerts: %d\n", len(report.Alerts))
		fmt.Printf("  Issues: %d\n", len(report.Issues))
		fmt.Printf("  Recommendations: %d\n", len(report.Recommendations))
	}

	// Stop plugin
	if err := plugin.GetManager().StopPlugin("security_demo"); err != nil {
		utils.Error("Failed to stop security demo plugin: %v", err)
		return
	}

	utils.Info("Security demo plugin test completed successfully.")
}