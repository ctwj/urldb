package config

import (
	"fmt"
	"log"
)

// ExampleUsage 演示如何使用插件配置系统
func ExampleUsage() {
	// 创建插件管理器
	manager := NewConfigManager()

	// 1. 创建配置模式
	fmt.Println("1. 创建配置模式")
	schema := NewConfigSchema("example-plugin", "1.0.0")

	// 添加配置字段
	intervalMin := 1.0
	intervalMax := 3600.0
	schema.AddField(ConfigField{
		Key:         "interval",
		Name:        "检查间隔",
		Description: "插件执行任务的时间间隔（秒）",
		Type:        "int",
		Required:    true,
		Default:     60,
		Min:         &intervalMin,
		Max:         &intervalMax,
	})

	schema.AddField(ConfigField{
		Key:         "enabled",
		Name:        "启用状态",
		Description: "插件是否启用",
		Type:        "bool",
		Required:    true,
		Default:     true,
	})

	schema.AddField(ConfigField{
		Key:         "api_key",
		Name:        "API密钥",
		Description: "访问外部服务的API密钥",
		Type:        "string",
		Required:    false,
		Encrypted:   true,
	})

	// 注册模式
	if err := manager.RegisterSchema(schema); err != nil {
		log.Fatalf("注册模式失败: %v", err)
	}

	// 2. 创建配置模板
	fmt.Println("2. 创建配置模板")
	config := map[string]interface{}{
		"interval": 30,
		"enabled":  true,
		"protocol": "https",
	}

	template := &ConfigTemplate{
		Name:        "production-config",
		Description: "生产环境配置模板",
		Config:      config,
		Version:     "1.0.0",
	}

	if err := manager.RegisterTemplate(template); err != nil {
		log.Fatalf("注册模板失败: %v", err)
	}

	// 3. 验证配置
	fmt.Println("3. 验证配置")
	userConfig := map[string]interface{}{
		"interval": 120,
		"enabled":  true,
		"api_key":  "secret-key-12345",
	}

	if err := manager.ValidateConfig("example-plugin", userConfig); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	} else {
		fmt.Println("配置验证通过")
	}

	// 4. 保存配置版本
	fmt.Println("4. 保存配置版本")
	if err := manager.SaveVersion("example-plugin", "1.0.0", "初始生产配置", "admin", userConfig); err != nil {
		log.Fatalf("保存配置版本失败: %v", err)
	}

	// 5. 应用模板
	fmt.Println("5. 应用模板")
	newConfig := make(map[string]interface{})
	if err := manager.ApplyTemplate("example-plugin", "production-config", newConfig); err != nil {
		log.Fatalf("应用模板失败: %v", err)
	}

	fmt.Printf("应用模板后的配置: %+v\n", newConfig)

	// 6. 获取最新版本
	fmt.Println("6. 获取最新版本")
	latestConfig, err := manager.GetLatestVersion("example-plugin")
	if err != nil {
		log.Fatalf("获取最新版本失败: %v", err)
	}

	fmt.Printf("最新配置版本: %+v\n", latestConfig)

	// 7. 列出所有模板
	fmt.Println("7. 列出所有模板")
	templates := manager.ListTemplates()
	for _, tmpl := range templates {
		fmt.Printf("模板: %s - %s\n", tmpl.Name, tmpl.Description)
	}

	fmt.Println("配置系统演示完成")
}