package demo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// DemoPlugin is a demo plugin that fetches a random resource from the database every minute
type DemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
}

// NewDemoPlugin creates a new demo plugin
func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{
		name:        "demo-plugin",
		version:     "1.0.0",
		description: "A demo plugin that fetches a random resource from the database every minute and logs it",
		author:      "urlDB Team",
	}
}

// Name returns the plugin name
func (p *DemoPlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *DemoPlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *DemoPlugin) Description() string {
	return p.description
}

// Author returns the plugin author
func (p *DemoPlugin) Author() string {
	return p.author
}

// Initialize initializes the plugin
func (p *DemoPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Demo plugin initialized")
	return nil
}

// Start starts the plugin
func (p *DemoPlugin) Start() error {
	p.context.LogInfo("Demo plugin started")

	// Register a task to run every minute
	return p.context.RegisterTask("demo-task", p.fetchAndLogResource)
}

// Stop stops the plugin
func (p *DemoPlugin) Stop() error {
	p.context.LogInfo("Demo plugin stopped")
	return nil
}

// Cleanup cleans up the plugin
func (p *DemoPlugin) Cleanup() error {
	p.context.LogInfo("Demo plugin cleaned up")
	return nil
}

// Dependencies returns the plugin dependencies
func (p *DemoPlugin) Dependencies() []string {
	return []string{}
}

// CheckDependencies checks the plugin dependencies
func (p *DemoPlugin) CheckDependencies() map[string]bool {
	return make(map[string]bool)
}

// fetchAndLogResource fetches a random resource and logs it
func (p *DemoPlugin) fetchAndLogResource() {
	// Simulate fetching a resource from the database
	resource := p.fetchRandomResource()

	if resource != nil {
		p.context.LogInfo("Fetched resource: ID=%d, Title=%s, URL=%s",
			resource.ID, resource.Title, resource.URL)
	} else {
		p.context.LogWarn("No resources found in database")
	}
}

// Resource represents a resource entity
type Resource struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// fetchRandomResource simulates fetching a random resource from the database
func (p *DemoPlugin) fetchRandomResource() *Resource {
	// In a real implementation, this would query the actual database
	// For demo purposes, we'll generate a random resource

	// Simulate some resources
	resources := []Resource{
		{ID: 1, Title: "Go语言编程指南", Description: "学习Go语言的完整指南", URL: "https://example.com/go-guide", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "微服务架构设计", Description: "构建可扩展的微服务系统", URL: "https://example.com/microservices", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Title: "容器化部署实践", Description: "Docker和Kubernetes实战", URL: "https://example.com/container-deployment", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 4, Title: "数据库优化技巧", Description: "提升数据库性能的方法", URL: "https://example.com/db-optimization", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 5, Title: "前端框架比较", Description: "React vs Vue vs Angular", URL: "https://example.com/frontend-frameworks", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Return a random resource
	if len(resources) > 0 {
		return &resources[rand.Intn(len(resources))]
	}

	return nil
}