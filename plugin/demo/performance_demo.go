package demo

import (
	"context"
	"fmt"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// PerformanceDemoPlugin 性能演示插件，展示新功能的使用
type PerformanceDemoPlugin struct {
	name        string
	version     string
	description string
	author      string
	dependencies []string
}

// NewPerformanceDemoPlugin 创建新的性能演示插件
func NewPerformanceDemoPlugin() *PerformanceDemoPlugin {
	return &PerformanceDemoPlugin{
		name:        "PerformanceDemoPlugin",
		version:     "1.0.0",
		description: "演示插件性能优化功能",
		author:      "Claude",
		dependencies: []string{},
	}
}

// Name 返回插件名称
func (p *PerformanceDemoPlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *PerformanceDemoPlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *PerformanceDemoPlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *PerformanceDemoPlugin) Author() string {
	return p.author
}

// Dependencies 返回插件依赖
func (p *PerformanceDemoPlugin) Dependencies() []string {
	return p.dependencies
}

// CheckDependencies 检查插件依赖
func (p *PerformanceDemoPlugin) CheckDependencies() map[string]bool {
	// 简单实现，实际插件中可能需要更复杂的依赖检查
	deps := make(map[string]bool)
	for _, dep := range p.dependencies {
		deps[dep] = true // 假设所有依赖都满足
	}
	return deps
}

// Initialize 初始化插件
func (p *PerformanceDemoPlugin) Initialize(ctx types.PluginContext) error {
	fmt.Printf("[%s] 初始化插件\n", p.name)

	// 设置并发限制
	if err := ctx.SetConcurrencyLimit(5); err != nil {
		return fmt.Errorf("设置并发限制失败: %v", err)
	}

	// 设置一些配置
	if err := ctx.SetConfig("demo_config_key", "demo_config_value"); err != nil {
		return fmt.Errorf("设置配置失败: %v", err)
	}

	// 设置一些数据
	if err := ctx.SetData("demo_data_key", "demo_data_value", "demo_type"); err != nil {
		return fmt.Errorf("设置数据失败: %v", err)
	}

	return nil
}

// Start 启动插件
func (p *PerformanceDemoPlugin) Start() error {
	fmt.Printf("[%s] 启动插件\n", p.name)
	return nil
}

// Stop 停止插件
func (p *PerformanceDemoPlugin) Stop() error {
	fmt.Printf("[%s] 停止插件\n", p.name)
	return nil
}

// Cleanup 清理插件
func (p *PerformanceDemoPlugin) Cleanup() error {
	fmt.Printf("[%s] 清理插件\n", p.name)
	return nil
}

// 演示缓存功能
func (p *PerformanceDemoPlugin) demoCache(ctx types.PluginContext) {
	fmt.Printf("[%s] 演示缓存功能\n", p.name)

	// 设置缓存项
	err := ctx.CacheSet("demo_cache_key", "demo_cache_value", 5*time.Minute)
	if err != nil {
		fmt.Printf("设置缓存失败: %v\n", err)
		return
	}

	// 获取缓存项
	value, err := ctx.CacheGet("demo_cache_key")
	if err != nil {
		fmt.Printf("获取缓存失败: %v\n", err)
		return
	}

	fmt.Printf("从缓存获取到值: %v\n", value)

	// 删除缓存项
	err = ctx.CacheDelete("demo_cache_key")
	if err != nil {
		fmt.Printf("删除缓存失败: %v\n", err)
		return
	}

	fmt.Printf("缓存项已删除\n")
}

// 演示并发控制功能
func (p *PerformanceDemoPlugin) demoConcurrency(ctx types.PluginContext) {
	fmt.Printf("[%s] 演示并发控制功能\n", p.name)

	// 创建多个并发任务
	tasks := 10
	results := make(chan string, tasks)

	// 启动多个并发任务
	for i := 0; i < tasks; i++ {
		go func(taskID int) {
			// 在并发控制下执行任务
			err := ctx.ConcurrencyExecute(context.Background(), func() error {
				// 模拟一些工作
				time.Sleep(100 * time.Millisecond)
				results <- fmt.Sprintf("任务 %d 完成", taskID)
				return nil
			})

			if err != nil {
				results <- fmt.Sprintf("任务 %d 失败: %v", taskID, err)
			}
		}(i)
	}

	// 收集结果
	for i := 0; i < tasks; i++ {
		result := <-results
		fmt.Println(result)
	}

	// 获取并发统计信息
	stats, err := ctx.GetConcurrencyStats()
	if err != nil {
		fmt.Printf("获取并发统计信息失败: %v\n", err)
	} else {
		fmt.Printf("并发统计信息: %+v\n", stats)
	}
}

// 演示数据存储优化功能
func (p *PerformanceDemoPlugin) demoDataStorage(ctx types.PluginContext) {
	fmt.Printf("[%s] 演示数据存储优化功能\n", p.name)

	// 设置数据
	err := ctx.SetData("optimized_key", "optimized_value", "demo_type")
	if err != nil {
		fmt.Printf("设置数据失败: %v\n", err)
		return
	}

	// 获取数据（第一次会从数据库获取并缓存）
	value1, err := ctx.GetData("optimized_key", "demo_type")
	if err != nil {
		fmt.Printf("获取数据失败: %v\n", err)
		return
	}
	fmt.Printf("第一次获取数据: %v\n", value1)

	// 再次获取数据（会从缓存获取）
	value2, err := ctx.GetData("optimized_key", "demo_type")
	if err != nil {
		fmt.Printf("获取数据失败: %v\n", err)
		return
	}
	fmt.Printf("第二次获取数据（来自缓存）: %v\n", value2)

	// 删除数据
	err = ctx.DeleteData("optimized_key", "demo_type")
	if err != nil {
		fmt.Printf("删除数据失败: %v\n", err)
		return
	}
	fmt.Printf("数据已删除\n")
}

// RunDemo 运行演示
func (p *PerformanceDemoPlugin) RunDemo(ctx types.PluginContext) {
	fmt.Printf("[%s] 运行性能优化演示\n", p.name)

	// 演示缓存功能
	p.demoCache(ctx)

	// 演示并发控制功能
	p.demoConcurrency(ctx)

	// 演示数据存储优化功能
	p.demoDataStorage(ctx)
}