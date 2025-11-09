package main

import (
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/examples/plugin_demo/full_demo_plugin"
	"github.com/ctwj/urldb/examples/plugin_demo/security_demo_plugin"
	"github.com/ctwj/urldb/examples/plugin_demo/uninstall_demo_plugin"
	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/loader"
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

	// Load binary plugins from filesystem
	// Create a simple plugin loader to load the binary plugins
	pluginLoader1 := loader.NewSimplePluginLoader("./binary_plugin1")
	pluginLoader2 := loader.NewSimplePluginLoader("./binary_plugin2")

	// Load first binary plugin
	if plugins1, err := pluginLoader1.LoadAllPlugins(); err != nil {
		utils.Error("Failed to load binary plugin1: %v", err)
		// Continue execution even if plugin loading fails
	} else {
		// Register the loaded plugins
		for _, p := range plugins1 {
			if err := plugin.RegisterPlugin(p); err != nil {
				utils.Error("Failed to register binary plugin1: %v", err)
			}
		}
	}

	// Load second binary plugin
	if plugins2, err := pluginLoader2.LoadAllPlugins(); err != nil {
		utils.Error("Failed to load binary plugin2: %v", err)
		// Continue execution even if plugin loading fails
	} else {
		// Register the loaded plugins
		for _, p := range plugins2 {
			if err := plugin.RegisterPlugin(p); err != nil {
				utils.Error("Failed to register binary plugin2: %v", err)
			}
		}
	}

	// Register demo plugins from new structure
	fullDemoPlugin := full_demo_plugin.NewFullDemoPlugin()
	if err := plugin.RegisterPlugin(fullDemoPlugin); err != nil {
		utils.Error("Failed to register full demo plugin: %v", err)
		return
	}

	securityDemoPlugin := security_demo_plugin.NewSecurityDemoPlugin()
	if err := plugin.RegisterPlugin(securityDemoPlugin); err != nil {
		utils.Error("Failed to register security demo plugin: %v", err)
		return
	}

	uninstallDemoPlugin := uninstall_demo_plugin.NewUninstallDemoPlugin()
	if err := plugin.RegisterPlugin(uninstallDemoPlugin); err != nil {
		utils.Error("Failed to register uninstall demo plugin: %v", err)
		return
	}

	// Initialize plugins
	config := map[string]interface{}{
		"interval": 30, // 30 seconds
		"enabled":  true,
	}
	if err := plugin.GetManager().InitializePlugin("full-demo-plugin", config); err != nil {
		utils.Error("Failed to initialize full demo plugin: %v", err)
		return
	}

	if err := plugin.GetManager().InitializePluginForHandler("security_demo"); err != nil {
		utils.Error("Failed to initialize security demo plugin: %v", err)
		return
	}

	if err := plugin.GetManager().InitializePluginForHandler("uninstall-demo"); err != nil {
		utils.Error("Failed to initialize uninstall demo plugin: %v", err)
		return
	}

	// Start plugins
	if err := plugin.GetManager().StartPlugin("full-demo-plugin"); err != nil {
		utils.Error("Failed to start full demo plugin: %v", err)
		return
	}

	if err := plugin.GetManager().StartPlugin("security_demo"); err != nil {
		utils.Error("Failed to start security demo plugin: %v", err)
		return
	}

	if err := plugin.GetManager().StartPlugin("uninstall-demo"); err != nil {
		utils.Error("Failed to start uninstall demo plugin: %v", err)
		return
	}

	// Keep the application running
	utils.Info("Plugin system test started. Running for 2 minutes...")
	time.Sleep(2 * time.Minute)

	// Stop plugins
	if err := plugin.GetManager().StopPlugin("full-demo-plugin"); err != nil {
		utils.Error("Failed to stop full demo plugin: %v", err)
		return
	}

	if err := plugin.GetManager().StopPlugin("security_demo"); err != nil {
		utils.Error("Failed to stop security demo plugin: %v", err)
		return
	}

	if err := plugin.GetManager().StopPlugin("uninstall-demo"); err != nil {
		utils.Error("Failed to stop uninstall demo plugin: %v", err)
		return
	}

	utils.Info("Plugin system test completed successfully.")
}