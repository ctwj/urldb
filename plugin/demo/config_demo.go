package demo

import (
	"time"

	"github.com/ctwj/urldb/plugin/config"
	"github.com/ctwj/urldb/plugin/types"
)

// ConfigDemoPlugin 演示插件配置功能
type ConfigDemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
}

// NewConfigDemoPlugin 创建新的配置演示插件
func NewConfigDemoPlugin() *ConfigDemoPlugin {
	return &ConfigDemoPlugin{
		name:        "config-demo-plugin",
		version:     "1.0.0",
		description: "插件配置功能演示插件",
		author:      "urlDB Team",
	}
}

// Name 返回插件名称
func (p *ConfigDemoPlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *ConfigDemoPlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *ConfigDemoPlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *ConfigDemoPlugin) Author() string {
	return p.author
}

// Initialize 初始化插件
func (p *ConfigDemoPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("配置演示插件初始化")

	// 演示读取配置
	interval, err := p.context.GetConfig("interval")
	if err != nil {
		p.context.LogWarn("无法获取interval配置: %v", err)
	} else {
		p.context.LogInfo("当前interval配置: %v", interval)
	}

	enabled, err := p.context.GetConfig("enabled")
	if err != nil {
		p.context.LogWarn("无法获取enabled配置: %v", err)
	} else {
		p.context.LogInfo("当前enabled配置: %v", enabled)
	}

	apiKey, err := p.context.GetConfig("api_key")
	if err != nil {
		p.context.LogWarn("无法获取api_key配置: %v", err)
	} else {
		p.context.LogInfo("当前api_key配置长度: %d", len(apiKey.(string)))
	}

	return nil
}

// Start 启动插件
func (p *ConfigDemoPlugin) Start() error {
	p.context.LogInfo("配置演示插件启动")

	// 注册定时任务
	err := p.context.RegisterTask("config-demo-task", p.demoTask)
	if err != nil {
		return err
	}

	return nil
}

// Stop 停止插件
func (p *ConfigDemoPlugin) Stop() error {
	p.context.LogInfo("配置演示插件停止")
	return nil
}

// Cleanup 清理插件
func (p *ConfigDemoPlugin) Cleanup() error {
	p.context.LogInfo("配置演示插件清理")
	return nil
}

// Dependencies 返回插件依赖
func (p *ConfigDemoPlugin) Dependencies() []string {
	return []string{}
}

// CheckDependencies 检查依赖
func (p *ConfigDemoPlugin) CheckDependencies() map[string]bool {
	return make(map[string]bool)
}

// CreateConfigSchema 创建插件配置模式
func (p *ConfigDemoPlugin) CreateConfigSchema() *config.ConfigSchema {
	schema := config.NewConfigSchema(p.name, p.version)

	// 添加配置字段
	intervalMin := 1.0
	intervalMax := 3600.0
	schema.AddField(config.ConfigField{
		Key:         "interval",
		Name:        "检查间隔",
		Description: "插件执行任务的时间间隔（秒）",
		Type:        "int",
		Required:    true,
		Default:     60,
		Min:         &intervalMin,
		Max:         &intervalMax,
	})

	schema.AddField(config.ConfigField{
		Key:         "enabled",
		Name:        "启用状态",
		Description: "插件是否启用",
		Type:        "bool",
		Required:    true,
		Default:     true,
	})

	schema.AddField(config.ConfigField{
		Key:         "api_key",
		Name:        "API密钥",
		Description: "访问外部服务的API密钥",
		Type:        "string",
		Required:    false,
		Encrypted:   true,
	})

	validProtocols := []string{"http", "https"}
	schema.AddField(config.ConfigField{
		Key:         "protocol",
		Name:        "协议",
		Description: "使用的网络协议",
		Type:        "string",
		Required:    false,
		Default:     "https",
		Enum:        validProtocols,
	})

	return schema
}

// demoTask 演示任务
func (p *ConfigDemoPlugin) demoTask() {
	p.context.LogInfo("执行配置演示任务于 %s", time.Now().Format(time.RFC3339))

	// 演示读取配置
	interval, err := p.context.GetConfig("interval")
	if err == nil {
		p.context.LogInfo("任务间隔: %v", interval)
	}

	enabled, err := p.context.GetConfig("enabled")
	if err == nil && enabled.(bool) {
		p.context.LogInfo("插件已启用，执行任务逻辑")
		// 在这里执行实际的任务逻辑
	} else {
		p.context.LogInfo("插件未启用，跳过任务执行")
	}
}

// CreateConfigTemplate 创建配置模板
func (p *ConfigDemoPlugin) CreateConfigTemplate() *config.ConfigTemplate {
	configData := map[string]interface{}{
		"interval": 30,
		"enabled":  true,
		"protocol": "https",
	}

	return &config.ConfigTemplate{
		Name:        "default-config",
		Description: "默认配置模板",
		Config:      configData,
		Version:     p.version,
	}
}