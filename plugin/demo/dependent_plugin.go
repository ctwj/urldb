package demo

import (
	"github.com/ctwj/urldb/plugin/types"
)

// DependentPlugin is a plugin that depends on other plugins
type DependentPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
	dependencies []string
}

// NewDependentPlugin creates a new dependent plugin
func NewDependentPlugin() *DependentPlugin {
	return &DependentPlugin{
		name:        "dependent-plugin",
		version:     "1.0.0",
		description: "A plugin that demonstrates dependency management features",
		author:      "urlDB Team",
		dependencies: []string{"demo-plugin"}, // This plugin depends on demo-plugin
	}
}

// Name returns the plugin name
func (p *DependentPlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *DependentPlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *DependentPlugin) Description() string {
	return p.description
}

// Author returns the plugin author
func (p *DependentPlugin) Author() string {
	return p.author
}

// Initialize initializes the plugin
func (p *DependentPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Dependent plugin initialized")
	return nil
}

// Start starts the plugin
func (p *DependentPlugin) Start() error {
	p.context.LogInfo("Dependent plugin started")

	// Check dependencies
	depStatus := p.CheckDependencies()
	for dep, satisfied := range depStatus {
		if satisfied {
			p.context.LogInfo("Dependency %s is satisfied", dep)
		} else {
			p.context.LogWarn("Dependency %s is NOT satisfied", dep)
		}
	}

	return nil
}

// Stop stops the plugin
func (p *DependentPlugin) Stop() error {
	p.context.LogInfo("Dependent plugin stopped")
	return nil
}

// Cleanup cleans up the plugin
func (p *DependentPlugin) Cleanup() error {
	p.context.LogInfo("Dependent plugin cleaned up")
	return nil
}

// Dependencies returns the plugin dependencies
func (p *DependentPlugin) Dependencies() []string {
	return p.dependencies
}

// CheckDependencies checks the plugin dependencies
func (p *DependentPlugin) CheckDependencies() map[string]bool {
	dependencies := p.Dependencies()
	result := make(map[string]bool)

	// In a real implementation, this would check with the plugin manager
	// For demo purposes, we'll just return a mock result
	for _, dep := range dependencies {
		// Assume all dependencies are satisfied for demo
		result[dep] = true
	}

	return result
}