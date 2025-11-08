package main

import (
	"fmt"
	"github.com/ctwj/urldb/plugin/types"
)

// 示例插件实现
type ExamplePlugin struct{}

func (p *ExamplePlugin) Name() string {
	return "example-plugin"
}

func (p *ExamplePlugin) Version() string {
	return "1.0.0"
}

func (p *ExamplePlugin) Description() string {
	return "这是一个示例插件"
}

func (p *ExamplePlugin) Author() string {
	return "Test Author"
}

func (p *ExamplePlugin) Initialize(ctx types.PluginContext) error {
	ctx.LogInfo("示例插件初始化")
	return nil
}

func (p *ExamplePlugin) Start() error {
	fmt.Println("示例插件启动")
	return nil
}

func (p *ExamplePlugin) Stop() error {
	fmt.Println("示例插件停止")
	return nil
}

func (p *ExamplePlugin) Cleanup() error {
	fmt.Println("示例插件清理")
	return nil
}

func (p *ExamplePlugin) Dependencies() []string {
	return []string{}
}

func (p *ExamplePlugin) CheckDependencies() map[string]bool {
	return map[string]bool{}
}

// 导出插件实例
var Plugin = &ExamplePlugin{}

func main() {
	// 编译为 .so 文件时，此函数不会被使用
}