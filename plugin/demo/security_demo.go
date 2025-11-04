package demo

import (
	"fmt"
	"time"

	"github.com/ctwj/urldb/plugin/security"
	"github.com/ctwj/urldb/plugin/types"
)

// SecurityDemoPlugin is a demo plugin to demonstrate security features
type SecurityDemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	config      map[string]interface{}
}

// NewSecurityDemoPlugin creates a new security demo plugin
func NewSecurityDemoPlugin() *SecurityDemoPlugin {
	return &SecurityDemoPlugin{
		name:        "security_demo",
		version:     "1.0.0",
		description: "A demo plugin to demonstrate security features",
		author:      "urlDB",
		config:      make(map[string]interface{}),
	}
}

// Name returns the plugin name
func (p *SecurityDemoPlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *SecurityDemoPlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *SecurityDemoPlugin) Description() string {
	return p.description
}

// Author returns the plugin author
func (p *SecurityDemoPlugin) Author() string {
	return p.author
}

// Dependencies returns the plugin dependencies
func (p *SecurityDemoPlugin) Dependencies() []string {
	return []string{}
}

// CheckDependencies checks the plugin dependencies
func (p *SecurityDemoPlugin) CheckDependencies() map[string]bool {
	return make(map[string]bool)
}

// Initialize initializes the plugin
func (p *SecurityDemoPlugin) Initialize(ctx types.PluginContext) error {
	ctx.LogInfo("Initializing security demo plugin")

	// Request additional permissions
	ctx.RequestPermission(string(security.PermissionConfigWrite), p.name)
	ctx.RequestPermission(string(security.PermissionDataWrite), p.name)

	// Test permission
	hasPerm, err := ctx.CheckPermission(string(security.PermissionConfigRead))
	if err != nil {
		ctx.LogError("Error checking permission: %v", err)
		return err
	}

	if !hasPerm {
		ctx.LogWarn("Plugin does not have config read permission")
		return fmt.Errorf("plugin does not have required permissions")
	}

	// Set some config
	ctx.SetConfig("initialized", true)
	ctx.SetConfig("timestamp", time.Now().Unix())

	ctx.LogInfo("Security demo plugin initialized successfully")

	// Register a demo task
	ctx.RegisterTask("security_demo_task", func() {
		ctx.LogInfo("Security demo task executed")
		// Try to access config
		if initialized, err := ctx.GetConfig("initialized"); err == nil {
			ctx.LogInfo("Plugin initialized: %v", initialized)
		}

		// Try to write some data
		ctx.SetData("demo_key", "demo_value", "demo_type")
	})

	return nil
}

// Start starts the plugin
func (p *SecurityDemoPlugin) Start() error {
	return nil
}

// Stop stops the plugin
func (p *SecurityDemoPlugin) Stop() error {
	return nil
}

// Cleanup cleans up the plugin
func (p *SecurityDemoPlugin) Cleanup() error {
	return nil
}