package main

import (
	"fmt"
	"github.com/ctwj/urldb/plugin/types"
)

// Plugin1 示例插件1
type Plugin1 struct{}

// Name 返回插件名称
func (p *Plugin1) Name() string {
	return "demo-plugin-1"
}

// Version 返回插件版本
func (p *Plugin1) Version() string {
	return "1.0.0"
}

// Description 返回插件描述
func (p *Plugin1) Description() string {
	return "这是一个简单的示例插件1"
}

// Author 返回插件作者
func (p *Plugin1) Author() string {
	return "Demo Author"
}

// Initialize 初始化插件
func (p *Plugin1) Initialize(ctx types.PluginContext) error {
	ctx.LogInfo("示例插件1初始化")
	return nil
}

// Start 启动插件
func (p *Plugin1) Start() error {
	fmt.Println("示例插件1启动")
	return nil
}

// Stop 停止插件
func (p *Plugin1) Stop() error {
	fmt.Println("示例插件1停止")
	return nil
}

// Cleanup 清理插件
func (p *Plugin1) Cleanup() error {
	fmt.Println("示例插件1清理")
	return nil
}

// Dependencies 返回插件依赖
func (p *Plugin1) Dependencies() []string {
	return []string{}
}

// CheckDependencies 检查插件依赖
func (p *Plugin1) CheckDependencies() map[string]bool {
	return map[string]bool{}
}

// 导出插件实例
var Plugin = &Plugin1{}

func main() {
	// 编译为 .so 文件时，此函数不会被使用
}