package demo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/plugin/types"
	"gorm.io/gorm"
)

// FullDemoPlugin 是一个完整功能的演示插件
type FullDemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
}

// NewFullDemoPlugin 创建完整演示插件
func NewFullDemoPlugin() *FullDemoPlugin {
	return &FullDemoPlugin{
		name:        "full-demo-plugin",
		version:     "1.0.0",
		description: "A full-featured demo plugin demonstrating all plugin capabilities",
		author:      "urlDB Team",
	}
}

// Name returns the plugin name
func (p *FullDemoPlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *FullDemoPlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *FullDemoPlugin) Description() string {
	return p.description
}

// Author returns the plugin author
func (p *FullDemoPlugin) Author() string {
	return p.author
}

// Initialize initializes the plugin
func (p *FullDemoPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Full Demo plugin initialized")

	// 设置一些示例配置
	p.context.SetConfig("interval", 60)
	p.context.SetConfig("enabled", true)
	p.context.SetConfig("max_items", 100)

	// 存储一些示例数据
	data := map[string]interface{}{
		"last_updated": time.Now().Format(time.RFC3339),
		"counter":      0,
		"status":       "initialized",
	}
	p.context.SetData("demo_data", data, "demo_type")

	// 获取并验证配置
	interval, err := p.context.GetConfig("interval")
	if err != nil {
		p.context.LogError("Failed to get interval config: %v", err)
		return err
	}
	p.context.LogInfo("Configured interval: %v", interval)

	return nil
}

// Start starts the plugin
func (p *FullDemoPlugin) Start() error {
	p.context.LogInfo("Full Demo plugin started")

	// 注册定时任务
	err := p.context.RegisterTask("demo-periodic-task", p.executePeriodicTask)
	if err != nil {
		p.context.LogError("Failed to register periodic task: %v", err)
		return err
	}

	// 演示数据库访问
	p.demoDatabaseAccess()

	// 演示数据存储功能
	p.demoDataStorage()

	return nil
}

// Stop stops the plugin
func (p *FullDemoPlugin) Stop() error {
	p.context.LogInfo("Full Demo plugin stopped")
	return nil
}

// Cleanup cleans up the plugin
func (p *FullDemoPlugin) Cleanup() error {
	p.context.LogInfo("Full Demo plugin cleaned up")
	return nil
}

// Dependencies returns the plugin dependencies
func (p *FullDemoPlugin) Dependencies() []string {
	return []string{}
}

// CheckDependencies checks the plugin dependencies
func (p *FullDemoPlugin) CheckDependencies() map[string]bool {
	return make(map[string]bool)
}

// executePeriodicTask 执行周期性任务
func (p *FullDemoPlugin) executePeriodicTask() {
	p.context.LogInfo("Executing periodic task at %s", time.Now().Format(time.RFC3339))

	// 从数据库获取数据
	data, err := p.context.GetData("demo_data", "demo_type")
	if err != nil {
		p.context.LogError("Failed to get demo data: %v", err)
		return
	}

	p.context.LogInfo("Retrieved demo data: %v", data)

	// 更新数据计数器
	if dataMap, ok := data.(map[string]interface{}); ok {
		count, ok := dataMap["counter"].(float64) // json.Unmarshal converts numbers to float64
		if !ok {
			count = 0
		}
		count++

		// 更新数据
		dataMap["counter"] = count
		dataMap["last_updated"] = time.Now().Format(time.RFC3339)
		dataMap["status"] = "running"

		err = p.context.SetData("demo_data", dataMap, "demo_type")
		if err != nil {
			p.context.LogError("Failed to update demo data: %v", err)
		} else {
			p.context.LogInfo("Updated demo data, counter: %v", count)
		}
	}

	// 演示配置访问
	enabled, err := p.context.GetConfig("enabled")
	if err != nil {
		p.context.LogError("Failed to get enabled config: %v", err)
		return
	}

	p.context.LogInfo("Plugin enabled status: %v", enabled)
}

// demoDatabaseAccess 演示数据库访问
func (p *FullDemoPlugin) demoDatabaseAccess() {
	db := p.context.GetDB()
	if db == nil {
		p.context.LogError("Database connection not available")
		return
	}

	// 将db转换为*gorm.DB
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		p.context.LogError("Failed to cast database connection to *gorm.DB")
		return
	}

	// 尝试查询一些数据（如果存在的话）
	var count int64
	err := gormDB.Model(&entity.Resource{}).Count(&count).Error
	if err != nil {
		p.context.LogWarn("Failed to query resources: %v", err)
	} else {
		p.context.LogInfo("Database access demo: found %d resources", count)
	}
}

// demoDataStorage 演示数据存储功能
func (p *FullDemoPlugin) demoDataStorage() {
	// 存储一些复杂数据
	complexData := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice", "email": "alice@example.com"},
			{"id": 2, "name": "Bob", "email": "bob@example.com"},
		},
		"settings": map[string]interface{}{
			"theme":     "dark",
			"language":  "en",
			"notifications": true,
		},
		"timestamp": time.Now().Unix(),
	}

	err := p.context.SetData("complex_data", complexData, "user_settings")
	if err != nil {
		p.context.LogError("Failed to store complex data: %v", err)
	} else {
		p.context.LogInfo("Successfully stored complex data")
	}

	// 读取复杂数据
	retrievedData, err := p.context.GetData("complex_data", "user_settings")
	if err != nil {
		p.context.LogError("Failed to retrieve complex data: %v", err)
	} else {
		p.context.LogInfo("Successfully retrieved complex data: %v", retrievedData)
	}

	// 演示删除数据
	err = p.context.DeleteData("complex_data", "user_settings")
	if err != nil {
		p.context.LogError("Failed to delete data: %v", err)
	} else {
		p.context.LogInfo("Successfully deleted data")
	}
}