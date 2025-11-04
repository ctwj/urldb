package config

import (
	"testing"
)

func TestConfigSchema(t *testing.T) {
	// 创建配置模式
	schema := NewConfigSchema("test-plugin", "1.0.0")

	// 添加字段
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

	// 验证配置
	config := map[string]interface{}{
		"interval": 30,
		"enabled":  true,
	}

	validator := NewConfigValidator(schema)
	if err := validator.Validate(config); err != nil {
		t.Errorf("配置验证失败: %v", err)
	}

	// 测试无效配置
	invalidConfig := map[string]interface{}{
		"interval": 5000, // 超出最大值
		"enabled":  true,
	}

	if err := validator.Validate(invalidConfig); err == nil {
		t.Error("应该验证失败，但没有失败")
	}
}

func TestConfigTemplate(t *testing.T) {
	// 创建模板管理器
	manager := NewConfigTemplateManager()

	// 创建模板
	config := map[string]interface{}{
		"interval": 30,
		"enabled":  true,
		"protocol": "https",
	}

	template := &ConfigTemplate{
		Name:        "default-config",
		Description: "默认配置模板",
		Config:      config,
		Version:     "1.0.0",
	}

	// 注册模板
	if err := manager.RegisterTemplate(template); err != nil {
		t.Errorf("注册模板失败: %v", err)
	}

	// 获取模板
	retrievedTemplate, err := manager.GetTemplate("default-config")
	if err != nil {
		t.Errorf("获取模板失败: %v", err)
	}

	if retrievedTemplate.Name != "default-config" {
		t.Error("模板名称不匹配")
	}

	// 应用模板到配置
	targetConfig := make(map[string]interface{})
	if err := manager.ApplyTemplate("default-config", targetConfig); err != nil {
		t.Errorf("应用模板失败: %v", err)
	}

	if targetConfig["interval"] != 30 {
		t.Error("模板应用不正确")
	}
}

func TestConfigVersion(t *testing.T) {
	// 创建版本管理器
	manager := NewConfigVersionManager(3)

	// 保存版本
	config1 := map[string]interface{}{
		"interval": 30,
		"enabled":  true,
	}

	if err := manager.SaveVersion("test-plugin", "1.0.0", "初始版本", "tester", config1); err != nil {
		t.Errorf("保存版本失败: %v", err)
	}

	// 获取最新版本
	latest, err := manager.GetLatestVersion("test-plugin")
	if err != nil {
		t.Errorf("获取最新版本失败: %v", err)
	}

	if latest.Version != "1.0.0" {
		t.Error("版本不匹配")
	}

	// 保存更多版本以测试限制
	config2 := map[string]interface{}{
		"interval": 60,
		"enabled":  true,
	}

	config3 := map[string]interface{}{
		"interval": 90,
		"enabled":  false,
	}

	config4 := map[string]interface{}{
		"interval": 120,
		"enabled":  true,
	}

	manager.SaveVersion("test-plugin", "1.1.0", "第二版本", "tester", config2)
	manager.SaveVersion("test-plugin", "1.2.0", "第三版本", "tester", config3)
	manager.SaveVersion("test-plugin", "1.3.0", "第四版本", "tester", config4)

	// 检查版本数量限制
	versions, _ := manager.ListVersions("test-plugin")
	if len(versions) != 3 {
		t.Errorf("版本数量不正确，期望3个，实际%d个", len(versions))
	}

	// 最新版本应该是1.3.0
	latest, _ = manager.GetLatestVersion("test-plugin")
	if latest.Version != "1.3.0" {
		t.Error("最新版本不正确")
	}
}