package main

import (
	"fmt"
	"github.com/ctwj/urldb/plugin/types"
)

// Plugin2 示例插件2
type Plugin2 struct{}

// Name 返回插件名称
func (p *Plugin2) Name() string {
	return "demo-plugin-2"
}

// Version 返回插件版本
func (p *Plugin2) Version() string {
	return "1.0.0"
}

// Description 返回插件描述
func (p *Plugin2) Description() string {
	return "这是一个简单的示例插件2"
}

// Author 返回插件作者
func (p *Plugin2) Author() string {
	return "Demo Author"
}

// Initialize 初始化插件
func (p *Plugin2) Initialize(ctx types.PluginContext) error {
	ctx.LogInfo("示例插件2初始化")
	return nil
}

// Start 启动插件
func (p *Plugin2) Start() error {
	fmt.Println("示例插件2启动")
	return nil
}

// Stop 停止插件
func (p *Plugin2) Stop() error {
	fmt.Println("示例插件2停止")
	return nil
}

// Cleanup 清理插件
func (p *Plugin2) Cleanup() error {
	fmt.Println("示例插件2清理")
	return nil
}

// Dependencies 返回插件依赖
func (p *Plugin2) Dependencies() []string {
	return []string{}
}

// CheckDependencies 检查插件依赖
func (p *Plugin2) CheckDependencies() map[string]bool {
	return map[string]bool{}
}

// 导出插件实例
var Plugin = &Plugin2{}

func main() {
	// 编译为 .so 文件时，此函数不会被使用
}