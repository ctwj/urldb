package main

import (
	"fmt"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// DemoPlugin 演示插件
type DemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	config      map[string]interface{}
}

// NewDemoPlugin 创建演示插件
func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{
		name:        "demo-plugin",
		version:     "1.0.0",
		description: "这是一个演示插件，用于测试插件系统",
		author:      "Plugin System Developer",
		config:      make(map[string]interface{}),
	}
}

// Name 返回插件名称
func (p *DemoPlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *DemoPlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *DemoPlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *DemoPlugin) Author() string {
	return p.author
}

// Start 启动插件
func (p *DemoPlugin) Start() error {
	fmt.Printf("启动插件: %s v%s\n", p.name, p.version)
	// 模拟一些初始化工作
	time.Sleep(100 * time.Millisecond)
	fmt.Println("插件启动完成")
	return nil
}

// Stop 停止插件
func (p *DemoPlugin) Stop() error {
	fmt.Printf("停止插件: %s\n", p.name)
	// 模拟一些清理工作
	time.Sleep(50 * time.Millisecond)
	fmt.Println("插件停止完成")
	return nil
}

// GetConfig 获取插件配置
func (p *DemoPlugin) GetConfig() map[string]interface{} {
	return p.config
}

// UpdateConfig 更新插件配置
func (p *DemoPlugin) UpdateConfig(config map[string]interface{}) error {
	p.config = config
	fmt.Printf("更新插件配置: %v\n", config)
	return nil
}

// Dependencies 返回插件依赖
func (p *DemoPlugin) Dependencies() []string {
	// 这个演示插件没有依赖
	return []string{}
}

// Plugin 导出的插件实例（这是插件加载器查找的符号）
var Plugin types.Plugin = NewDemoPlugin()

func main() {
	// 这个main函数仅用于独立测试插件
	fmt.Println("演示插件测试")
	plugin := NewDemoPlugin()
	fmt.Printf("插件名称: %s\n", plugin.Name())
	fmt.Printf("插件版本: %s\n", plugin.Version())
	fmt.Printf("插件描述: %s\n", plugin.Description())

	// 测试启动和停止
	plugin.Start()
	plugin.Stop()
}