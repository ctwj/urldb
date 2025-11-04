# 插件系统开发指南

## 概述

urlDB插件系统提供了一个灵活的扩展机制，允许开发者创建自定义功能来增强系统能力。插件采用进程内加载模式，避免使用Go标准plugin包的限制。

## 插件生命周期

1. **注册 (Register)** - 插件被发现并注册到管理器
2. **初始化 (Initialize)** - 插件接收上下文并准备运行
3. **启动 (Start)** - 插件开始执行主要功能
4. **运行 (Running)** - 插件正常工作状态
5. **停止 (Stop)** - 插件停止运行
6. **清理 (Cleanup)** - 插件释放资源

## 创建插件

### 基本插件结构

### 可配置插件

插件可以实现 `ConfigurablePlugin` 接口来支持配置模式和模板：

```go
// ConfigurablePlugin is an optional interface for plugins that support configuration schemas
type ConfigurablePlugin interface {
    // CreateConfigSchema creates the configuration schema for the plugin
    CreateConfigSchema() *config.ConfigSchema

    // CreateConfigTemplate creates a default configuration template
    CreateConfigTemplate() *config.ConfigTemplate
}
```

```go
package myplugin

import (
    "github.com/ctwj/urldb/plugin/config"
    "github.com/ctwj/urldb/plugin/types"
)

// MyPlugin 实现插件接口
type MyPlugin struct {
    name        string
    version     string
    description string
    author      string
    context     types.PluginContext
}

// NewMyPlugin 创建新插件实例
func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        name:        "my-plugin",
        version:     "1.0.0",
        description: "我的自定义插件",
        author:      "开发者名称",
    }
}

// 实现必需的接口方法
func (p *MyPlugin) Name() string { return p.name }
func (p *MyPlugin) Version() string { return p.version }
func (p *MyPlugin) Description() string { return p.description }
func (p *MyPlugin) Author() string { return p.author }

func (p *MyPlugin) Initialize(ctx types.PluginContext) error {
    p.context = ctx
    p.context.LogInfo("插件初始化完成")
    return nil
}

func (p *MyPlugin) Start() error {
    p.context.LogInfo("插件启动")
    return nil
}

func (p *MyPlugin) Stop() error {
    p.context.LogInfo("插件停止")
    return nil
}

func (p *MyPlugin) Cleanup() error {
    p.context.LogInfo("插件清理")
    return nil
}

func (p *MyPlugin) Dependencies() []string {
    return []string{}
}

func (p *MyPlugin) CheckDependencies() map[string]bool {
    return make(map[string]bool)
}

// 实现可选的配置接口
func (p *MyPlugin) CreateConfigSchema() *config.ConfigSchema {
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

    return schema
}

func (p *MyPlugin) CreateConfigTemplate() *config.ConfigTemplate {
    config := map[string]interface{}{
        "interval": 30,
        "enabled":  true,
    }

    return &config.ConfigTemplate{
        Name:        "default-config",
        Description: "默认配置模板",
        Config:      config,
        Version:     p.version,
    }
}
```

## 插件上下文 API

### 日志功能

```go
// 记录不同级别的日志
p.context.LogDebug("调试信息: %s", "详细信息")
p.context.LogInfo("普通信息: %d", 42)
p.context.LogWarn("警告信息")
p.context.LogError("错误信息: %v", err)
```

### 配置管理

```go
// 设置配置
err := p.context.SetConfig("interval", 60)
err := p.context.SetConfig("enabled", true)
err := p.context.SetConfig("api_key", "secret-key")

// 获取配置
interval, err := p.context.GetConfig("interval")
enabled, err := p.context.GetConfig("enabled")
```

### 数据存储

```go
// 存储数据
data := map[string]interface{}{
    "last_update": time.Now().Unix(),
    "counter":     0,
    "status":      "active",
}
err := p.context.SetData("my_data", data, "app_state")

// 读取数据
retrievedData, err := p.context.GetData("my_data", "app_state")

// 删除数据
err := p.context.DeleteData("my_data", "app_state")
```

### 数据库访问

```go
// 获取数据库连接
db := p.context.GetDB()
if gormDB, ok := db.(*gorm.DB); ok {
    // 执行数据库操作
    var count int64
    err := gormDB.Model(&entity.Resource{}).Count(&count).Error
    if err != nil {
        p.context.LogError("查询失败: %v", err)
    } else {
        p.context.LogInfo("资源数量: %d", count)
    }
}
```

### 任务调度

```go
// 注册定时任务
err := p.context.RegisterTask("my-periodic-task", func() {
    p.context.LogInfo("执行定时任务于 %s", time.Now().Format(time.RFC3339))
    // 任务逻辑...
})

if err != nil {
    p.context.LogError("注册任务失败: %v", err)
}
```

## 自动注册插件

```go
package myplugin

import (
    "github.com/ctwj/urldb/plugin"
    "github.com/ctwj/urldb/plugin/types"
)

// init 函数在包导入时自动调用
func init() {
    plugin := NewMyPlugin()
    RegisterPlugin(plugin)
}

// RegisterPlugin 注册插件到全局管理器
func RegisterPlugin(pluginInstance types.Plugin) {
    if plugin.GetManager() != nil {
        plugin.GetManager().RegisterPlugin(pluginInstance)
    }
}
```

## 完整示例插件

```go
package demo

import (
    "time"

    "github.com/ctwj/urldb/db/entity"
    "github.com/ctwj/urldb/plugin/types"
    "gorm.io/gorm"
)

// FullDemoPlugin 完整功能演示插件
type FullDemoPlugin struct {
    name        string
    version     string
    description string
    author      string
    context     types.PluginContext
}

func NewFullDemoPlugin() *FullDemoPlugin {
    return &FullDemoPlugin{
        name:        "full-demo-plugin",
        version:     "1.0.0",
        description: "完整功能演示插件",
        author:      "urlDB Team",
    }
}

// ... 实现接口方法

func (p *FullDemoPlugin) Initialize(ctx types.PluginContext) error {
    p.context = ctx
    p.context.LogInfo("演示插件初始化")

    // 设置配置
    p.context.SetConfig("interval", 60)
    p.context.SetConfig("enabled", true)

    // 存储初始数据
    data := map[string]interface{}{
        "start_time": time.Now().Format(time.RFC3339),
        "counter":    0,
    }
    p.context.SetData("demo_stats", data, "monitoring")

    return nil
}

func (p *FullDemoPlugin) Start() error {
    p.context.LogInfo("演示插件启动")

    // 注册定时任务
    err := p.context.RegisterTask("demo-task", p.demoTask)
    if err != nil {
        return err
    }

    // 演示数据库访问
    p.demoDatabaseAccess()

    return nil
}

func (p *FullDemoPlugin) demoTask() {
    p.context.LogInfo("执行演示任务")

    // 更新计数器
    data, err := p.context.GetData("demo_stats", "monitoring")
    if err == nil {
        if stats, ok := data.(map[string]interface{}); ok {
            counter, _ := stats["counter"].(float64)
            stats["counter"] = counter + 1
            stats["last_update"] = time.Now().Format(time.RFC3339)
            p.context.SetData("demo_stats", stats, "monitoring")
        }
    }
}
```

## 最佳实践

1. **错误处理**：始终检查错误并适当处理
2. **资源清理**：在Cleanup方法中释放所有资源
3. **配置验证**：在初始化时验证配置有效性
4. **日志记录**：使用适当的日志级别记录重要事件
5. **性能考虑**：避免在插件中执行阻塞操作

## 常见问题

### Q: 插件如何访问数据库？
A: 通过 `p.context.GetDB()` 获取数据库连接，转换为 `*gorm.DB` 后使用。

### Q: 插件如何存储持久化数据？
A: 使用 `SetData()`、`GetData()`、`DeleteData()` 方法。

### Q: 插件如何注册定时任务？
A: 使用 `RegisterTask()` 方法注册任务函数。

### Q: 插件如何记录日志？
A: 使用 `LogDebug()`、`LogInfo()`、`LogWarn()`、`LogError()` 方法。

## 部署流程

1. 将插件代码放在 `plugin/` 目录下
2. 确保包含自动注册的 `init()` 函数
3. 构建主应用程序：`go build -o main .`
4. 启动应用程序，插件将自动注册
5. 通过API或管理界面启用插件

## API参考

### PluginContext 接口

- `LogDebug(msg string, args ...interface{})` - 调试日志
- `LogInfo(msg string, args ...interface{})` - 信息日志
- `LogWarn(msg string, args ...interface{})` - 警告日志
- `LogError(msg string, args ...interface{})` - 错误日志
- `GetConfig(key string) (interface{}, error)` - 获取配置
- `SetConfig(key string, value interface{}) error` - 设置配置
- `GetData(key string, dataType string) (interface{}, error)` - 获取数据
- `SetData(key string, value interface{}, dataType string) error` - 设置数据
- `DeleteData(key string, dataType string) error` - 删除数据
- `RegisterTask(name string, taskFunc func()) error` - 注册任务
- `GetDB() interface{}` - 获取数据库连接

### Plugin 接口

- `Name() string` - 插件名称
- `Version() string` - 插件版本
- `Description() string` - 插件描述
- `Author() string` - 插件作者
- `Initialize(ctx PluginContext) error` - 初始化插件
- `Start() error` - 启动插件
- `Stop() error` - 停止插件
- `Cleanup() error` - 清理插件
- `Dependencies() []string` - 获取依赖
- `CheckDependencies() map[string]bool` - 检查依赖