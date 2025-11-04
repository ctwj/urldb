package main

import (
	"fmt"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// UninstallDemoPlugin is a demo plugin to demonstrate uninstall functionality
type UninstallDemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	dependencies []string
	context     types.PluginContext
}

// NewUninstallDemoPlugin creates a new uninstall demo plugin
func NewUninstallDemoPlugin() *UninstallDemoPlugin {
	return &UninstallDemoPlugin{
		name:        "uninstall-demo",
		version:     "1.0.0",
		description: "A demo plugin to demonstrate uninstall functionality",
		author:      "Plugin Developer",
		dependencies: []string{},
	}
}

// Name returns the plugin name
func (p *UninstallDemoPlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *UninstallDemoPlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *UninstallDemoPlugin) Description() string {
	return p.description
}

// Author returns the plugin author
func (p *UninstallDemoPlugin) Author() string {
	return p.author
}

// Dependencies returns the plugin dependencies
func (p *UninstallDemoPlugin) Dependencies() []string {
	return p.dependencies
}

// CheckDependencies checks the status of plugin dependencies
func (p *UninstallDemoPlugin) CheckDependencies() map[string]bool {
	// For demo purposes, all dependencies are satisfied
	result := make(map[string]bool)
	for _, dep := range p.dependencies {
		result[dep] = true
	}
	return result
}

// Initialize initializes the plugin
func (p *UninstallDemoPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Initializing uninstall demo plugin")

	// Set some demo configuration
	if err := p.context.SetConfig("demo_setting", "uninstall_demo_value"); err != nil {
		return err
	}

	// Set some demo data
	if err := p.context.SetData("demo_key", "uninstall_demo_data", "demo_type"); err != nil {
		return err
	}

	p.context.LogInfo("Uninstall demo plugin initialized successfully")
	return nil
}

// Start starts the plugin
func (p *UninstallDemoPlugin) Start() error {
	p.context.LogInfo("Starting uninstall demo plugin")

	// Simulate some work
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			p.context.LogInfo("Uninstall demo plugin working... %d", i+1)
		}
		p.context.LogInfo("Uninstall demo plugin work completed")
	}()

	p.context.LogInfo("Uninstall demo plugin started successfully")
	return nil
}

// Stop stops the plugin
func (p *UninstallDemoPlugin) Stop() error {
	p.context.LogInfo("Stopping uninstall demo plugin")

	// Perform any necessary cleanup before stopping
	p.context.LogInfo("Uninstall demo plugin stopped successfully")
	return nil
}

// Cleanup performs final cleanup when uninstalling the plugin
func (p *UninstallDemoPlugin) Cleanup() error {
	p.context.LogInfo("Cleaning up uninstall demo plugin")

	// Perform any final cleanup operations
	// This might include:
	// - Removing temporary files
	// - Cleaning up external resources
	// - Notifying external services

	p.context.LogInfo("Uninstall demo plugin cleaned up successfully")
	return nil
}

// GetContext returns the plugin context (for testing purposes)
func (p *UninstallDemoPlugin) GetContext() types.PluginContext {
	return p.context
}

func main() {
	// This is just for demonstration purposes
	// In a real scenario, the plugin would be loaded by the plugin manager
	plugin := NewUninstallDemoPlugin()
	fmt.Printf("Plugin: %s v%s\n", plugin.Name(), plugin.Version())
	fmt.Printf("Description: %s\n", plugin.Description())
	fmt.Printf("Author: %s\n", plugin.Author())
}