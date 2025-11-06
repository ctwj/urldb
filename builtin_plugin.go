package main

import (
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// BuiltinPlugin 内置插件用于测试
type BuiltinPlugin struct {
	name        string
	version     string
	description string
	author      string
}

// NewBuiltinPlugin 创建内置插件实例
func NewBuiltinPlugin() *BuiltinPlugin {
	return &BuiltinPlugin{
		name:        "builtin-demo",
		version:     "1.0.0",
		description: "内置演示插件，用于测试插件系统功能",
		author:      "urlDB Team",
	}
}

// Name 返回插件名称
func (p *BuiltinPlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *BuiltinPlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *BuiltinPlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *BuiltinPlugin) Author() string {
	return p.author
}

// Initialize 初始化插件
func (p *BuiltinPlugin) Initialize(ctx types.PluginContext) error {
	utils.Info("Initializing builtin plugin: %s", p.name)
	ctx.LogInfo("Builtin plugin %s initialized successfully", p.name)
	return nil
}

// Start 启动插件
func (p *BuiltinPlugin) Start() error {
	utils.Info("Starting builtin plugin: %s", p.name)
	return nil
}

// Stop 停止插件
func (p *BuiltinPlugin) Stop() error {
	utils.Info("Stopping builtin plugin: %s", p.name)
	return nil
}

// Cleanup 清理插件资源
func (p *BuiltinPlugin) Cleanup() error {
	utils.Info("Cleaning up builtin plugin: %s", p.name)
	return nil
}

// Dependencies 返回插件依赖
func (p *BuiltinPlugin) Dependencies() []string {
	return []string{} // 无依赖
}

// CheckDependencies 检查插件依赖
func (p *BuiltinPlugin) CheckDependencies() map[string]bool {
	return map[string]bool{} // 无依赖需要检查
}