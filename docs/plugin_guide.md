# urlDB插件开发指南

## 目录
1. [插件系统概述](#插件系统概述)
2. [插件开发规范](#插件开发规范)
3. [插件实现步骤](#插件实现步骤)
4. [插件编译](#插件编译)
5. [插件加载和使用](#插件加载和使用)
6. [插件配置管理](#插件配置管理)
7. [插件数据管理](#插件数据管理)
8. [插件日志记录](#插件日志记录)
9. [插件任务调度](#插件任务调度)
10. [插件生命周期管理](#插件生命周期管理)
11. [插件依赖管理](#插件依赖管理)
12. [插件安全和权限](#插件安全和权限)
13. [插件测试](#插件测试)
14. [最佳实践](#最佳实践)

## 插件系统概述

urlDB插件系统是一个轻量级、高性能的插件框架，旨在为系统提供模块化扩展能力。插件系统采用进程内加载模式，与现有系统架构高度融合，提供完整的生命周期管理、配置管理、数据管理、日志记录和任务调度功能。

### 核心特性
- **轻量级实现**：避免复杂的.so文件管理
- **高性能**：进程内调用，无IPC开销
- **易集成**：与现有系统组件无缝集成
- **完整功能**：支持配置、数据、日志、任务等核心功能
- **安全可靠**：提供权限控制和资源隔离

## 插件开发规范

### 基础接口规范
所有插件必须实现`types.Plugin`接口：

```go
type Plugin interface {
    // 基本信息
    Name() string
    Version() string
    Description() string
    Author() string

    // 生命周期
    Initialize(ctx PluginContext) error
    Start() error
    Stop() error
    Cleanup() error

    // 依赖管理
    Dependencies() []string
    CheckDependencies() map[string]bool
}
```

### 上下文接口规范
插件通过`PluginContext`与系统交互：

```go
type PluginContext interface {
    // 日志功能
    LogDebug(msg string, args ...interface{})
    LogInfo(msg string, args ...interface{})
    LogWarn(msg string, args ...interface{})
    LogError(msg string, args ...interface{})

    // 配置功能
    GetConfig(key string) (interface{}, error)
    SetConfig(key string, value interface{}) error

    // 数据功能
    GetData(key string, dataType string) (interface{}, error)
    SetData(key string, value interface{}, dataType string) error
    DeleteData(key string, dataType string) error

    // 任务功能
    RegisterTask(name string, task func()) error
    UnregisterTask(name string) error

    // 数据库功能
    GetDB() interface{} // 返回 *gorm.DB
}
```

## 插件实现步骤

### 1. 创建插件结构体
```go
package myplugin

import (
    "github.com/ctwj/urldb/plugin/types"
)

type MyPlugin struct {
    name        string
    version     string
    description string
    author      string
    context     types.PluginContext
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        name:        "my-plugin",
        version:     "1.0.0",
        description: "My custom plugin",
        author:      "Developer Name",
    }
}
```

### 2. 实现基础接口方法
```go
// Name returns the plugin name
func (p *MyPlugin) Name() string {
    return p.name
}

// Version returns the plugin version
func (p *MyPlugin) Version() string {
    return p.version
}

// Description returns the plugin description
func (p *MyPlugin) Description() string {
    return p.description
}

// Author returns the plugin author
func (p *MyPlugin) Author() string {
    return p.author
}

// Dependencies returns plugin dependencies
func (p *MyPlugin) Dependencies() []string {
    return []string{}
}

// CheckDependencies checks plugin dependencies
func (p *MyPlugin) CheckDependencies() map[string]bool {
    return make(map[string]bool)
}
```

### 3. 实现生命周期方法
```go
// Initialize initializes the plugin
func (p *MyPlugin) Initialize(ctx types.PluginContext) error {
    p.context = ctx
    p.context.LogInfo("Plugin initialized")
    return nil
}

// Start starts the plugin
func (p *MyPlugin) Start() error {
    p.context.LogInfo("Plugin started")
    // 注册定时任务或其他初始化工作
    return nil
}

// Stop stops the plugin
func (p *MyPlugin) Stop() error {
    p.context.LogInfo("Plugin stopped")
    return nil
}

// Cleanup cleans up the plugin
func (p *MyPlugin) Cleanup() error {
    p.context.LogInfo("Plugin cleaned up")
    return nil
}
```

### 4. 实现插件功能
```go
// 示例：实现一个定时任务
func (p *MyPlugin) doWork() {
    p.context.LogInfo("Doing work...")
    // 实际业务逻辑
}
```

## 插件编译

### 1. 创建插件目录
```bash
mkdir -p /Users/kerwin/Program/go/urldb/plugin/myplugin
```

### 2. 编写插件代码
创建`myplugin.go`文件并实现插件逻辑。

### 3. 编译插件
由于采用进程内加载模式，插件需要在主程序中注册：

```go
// 在main.go中注册插件
import (
    "github.com/ctwj/urldb/plugin"
    "github.com/ctwj/urldb/plugin/myplugin"
)

func main() {
    // 初始化插件系统
    plugin.InitPluginSystem(taskManager)

    // 注册插件
    myPlugin := myplugin.NewMyPlugin()
    plugin.RegisterPlugin(myPlugin)

    // 其他初始化代码...
}
```

### 4. 构建主程序
```bash
cd /Users/kerwin/Program/go/urldb
go build -o urldb .
```

## 插件加载和使用

### 1. 插件自动注册
插件通过`init()`函数或在主程序中手动注册：

```go
// 方法1: 在插件包中使用init函数自动注册
func init() {
    plugin.RegisterPlugin(NewMyPlugin())
}

// 方法2: 在main.go中手动注册
func main() {
    plugin.InitPluginSystem(taskManager)
    plugin.RegisterPlugin(myplugin.NewMyPlugin())
    // ...
}
```

### 2. 插件配置
在系统启动时配置插件：

```go
// 初始化插件配置
config := map[string]interface{}{
    "setting1": "value1",
    "setting2": 42,
}

// 初始化插件
if err := pluginManager.InitializePlugin("my-plugin", config); err != nil {
    utils.Error("Failed to initialize plugin: %v", err)
    return
}

// 启动插件
if err := pluginManager.StartPlugin("my-plugin"); err != nil {
    utils.Error("Failed to start plugin: %v", err)
    return
}
```

### 3. 插件状态管理
```go
// 获取插件状态
status := pluginManager.GetPluginStatus("my-plugin")

// 停止插件
if err := pluginManager.StopPlugin("my-plugin"); err != nil {
    utils.Error("Failed to stop plugin: %v", err)
}
```

## 插件配置管理

### 1. 获取配置
```go
func (p *MyPlugin) doWork() {
    // 获取配置值
    setting1, err := p.context.GetConfig("setting1")
    if err != nil {
        p.context.LogError("Failed to get setting1: %v", err)
        return
    }

    p.context.LogInfo("Setting1 value: %v", setting1)
}
```

### 2. 设置配置
```go
func (p *MyPlugin) updateConfig() {
    // 设置配置值
    err := p.context.SetConfig("new_setting", "new_value")
    if err != nil {
        p.context.LogError("Failed to set config: %v", err)
        return
    }

    p.context.LogInfo("Config updated successfully")
}
```

### 3. 配置模式定义
```go
// 插件可以定义自己的配置模式
type ConfigSchema struct {
    Setting1 string `json:"setting1" validate:"required"`
    Setting2 int    `json:"setting2" validate:"min=1,max=100"`
    Enabled  bool   `json:"enabled"`
}
```

## 插件数据管理

### 1. 存储数据
```go
func (p *MyPlugin) saveData() {
    data := map[string]interface{}{
        "key1": "value1",
        "key2": 123,
    }

    // 存储数据
    err := p.context.SetData("my_data_key", data, "my_data_type")
    if err != nil {
        p.context.LogError("Failed to save data: %v", err)
        return
    }

    p.context.LogInfo("Data saved successfully")
}
```

### 2. 读取数据
```go
func (p *MyPlugin) loadData() {
    // 读取数据
    data, err := p.context.GetData("my_data_key", "my_data_type")
    if err != nil {
        p.context.LogError("Failed to load data: %v", err)
        return
    }

    p.context.LogInfo("Loaded data: %v", data)
}
```

### 3. 删除数据
```go
func (p *MyPlugin) deleteData() {
    // 删除数据
    err := p.context.DeleteData("my_data_key", "my_data_type")
    if err != nil {
        p.context.LogError("Failed to delete data: %v", err)
        return
    }

    p.context.LogInfo("Data deleted successfully")
}
```

## 插件日志记录

### 1. 日志级别
```go
func (p *MyPlugin) logExamples() {
    // 调试日志
    p.context.LogDebug("Debug message: %s", "debug info")

    // 信息日志
    p.context.LogInfo("Info message: %s", "process started")

    // 警告日志
    p.context.LogWarn("Warning message: %s", "deprecated feature used")

    // 错误日志
    p.context.LogError("Error message: %v", err)
}
```

### 2. 结构化日志
```go
func (p *MyPlugin) structuredLogging() {
    p.context.LogInfo("Processing resource: id=%d, title=%s, url=%s",
        resource.ID, resource.Title, resource.URL)
}
```

## 插件任务调度

### 1. 注册定时任务
```go
func (p *MyPlugin) Start() error {
    p.context.LogInfo("Plugin started")

    // 注册定时任务
    return p.context.RegisterTask("my-task", p.scheduledTask)
}

func (p *MyPlugin) scheduledTask() {
    p.context.LogInfo("Scheduled task executed")
    // 执行定时任务逻辑
}
```

### 2. 取消任务
```go
func (p *MyPlugin) Stop() error {
    // 取消定时任务
    err := p.context.UnregisterTask("my-task")
    if err != nil {
        p.context.LogError("Failed to unregister task: %v", err)
    }

    p.context.LogInfo("Plugin stopped")
    return nil
}
```

## 插件生命周期管理

### 1. 状态转换
```
Registered → Initialized → Starting → Running → Stopping → Stopped
     ↑                                            ↓
     └────────────── Cleanup ←───────────────────┘
```

### 2. 状态检查
```go
func (p *MyPlugin) checkStatus() {
    // 在插件方法中检查状态
    if p.context == nil {
        p.context.LogError("Plugin not initialized")
        return
    }

    // 执行业务逻辑
}
```

## 插件依赖管理

### 1. 声明依赖
```go
func (p *MyPlugin) Dependencies() []string {
    return []string{"database-plugin", "auth-plugin"}
}
```

### 2. 检查依赖
```go
func (p *MyPlugin) CheckDependencies() map[string]bool {
    deps := make(map[string]bool)
    deps["database-plugin"] = plugin.GetManager().GetPluginStatus("database-plugin") == types.StatusRunning
    deps["auth-plugin"] = plugin.GetManager().GetPluginStatus("auth-plugin") == types.StatusRunning
    return deps
}
```

## 插件安全和权限

### 1. 数据库访问控制
```go
func (p *MyPlugin) accessDatabase() {
    // 获取数据库连接
    db := p.context.GetDB().(*gorm.DB)

    // 执行安全的数据库操作
    var resources []Resource
    err := db.Where("is_public = ?", true).Find(&resources).Error
    if err != nil {
        p.context.LogError("Database query failed: %v", err)
        return
    }

    p.context.LogInfo("Found %d public resources", len(resources))
}
```

### 2. 权限检查
```go
func (p *MyPlugin) checkPermissions() bool {
    // 检查插件权限
    // 在实际实现中，可以从配置或上下文中获取权限信息
    return true
}
```

## 插件测试

### 1. 单元测试
```go
func TestMyPlugin_Name(t *testing.T) {
    plugin := NewMyPlugin()
    expected := "my-plugin"
    if plugin.Name() != expected {
        t.Errorf("Expected %s, got %s", expected, plugin.Name())
    }
}

func TestMyPlugin_Initialize(t *testing.T) {
    plugin := NewMyPlugin()
    mockContext := &MockPluginContext{}

    err := plugin.Initialize(mockContext)
    if err != nil {
        t.Errorf("Initialize failed: %v", err)
    }
}
```

### 2. 集成测试
```go
func TestMyPlugin_Lifecycle(t *testing.T) {
    plugin := NewMyPlugin()
    mockContext := &MockPluginContext{}

    // 测试完整生命周期
    if err := plugin.Initialize(mockContext); err != nil {
        t.Fatalf("Initialize failed: %v", err)
    }

    if err := plugin.Start(); err != nil {
        t.Fatalf("Start failed: %v", err)
    }

    if err := plugin.Stop(); err != nil {
        t.Fatalf("Stop failed: %v", err)
    }

    if err := plugin.Cleanup(); err != nil {
        t.Fatalf("Cleanup failed: %v", err)
    }
}
```

## 最佳实践

### 1. 错误处理
```go
func (p *MyPlugin) robustMethod() {
    defer func() {
        if r := recover(); r != nil {
            p.context.LogError("Panic recovered: %v", r)
        }
    }()

    // 业务逻辑
    result, err := someOperation()
    if err != nil {
        p.context.LogError("Operation failed: %v", err)
        return
    }

    p.context.LogInfo("Operation succeeded: %v", result)
}
```

### 2. 资源管理
```go
func (p *MyPlugin) manageResources() {
    // 确保资源正确释放
    defer func() {
        // 清理资源
        p.cleanup()
    }()

    // 使用资源
    p.useResources()
}
```

### 3. 配置验证
```go
func (p *MyPlugin) validateConfig() error {
    setting1, err := p.context.GetConfig("setting1")
    if err != nil {
        return fmt.Errorf("missing required config: setting1")
    }

    if setting1 == "" {
        return fmt.Errorf("setting1 cannot be empty")
    }

    return nil
}
```

### 4. 日志规范
```go
func (p *MyPlugin) logWithContext() {
    // 包含足够的上下文信息
    p.context.LogInfo("Processing user action: user_id=%d, action=%s, resource_id=%d",
        userID, action, resourceID)
}
```

### 5. 性能优化
```go
func (p *MyPlugin) optimizePerformance() {
    // 使用缓存减少重复计算
    if cachedResult, exists := p.getFromCache("key"); exists {
        p.context.LogInfo("Using cached result")
        return cachedResult
    }

    // 执行计算
    result := p.expensiveOperation()

    // 缓存结果
    p.setCache("key", result)

    return result
}
```

## 示例插件完整代码

以下是一个完整的示例插件实现：

```go
package example

import (
    "fmt"
    "time"

    "github.com/ctwj/urldb/plugin/types"
    "github.com/ctwj/urldb/utils"
)

// ExamplePlugin 示例插件
type ExamplePlugin struct {
    name        string
    version     string
    description string
    author      string
    context     types.PluginContext
}

// NewExamplePlugin 创建示例插件
func NewExamplePlugin() *ExamplePlugin {
    return &ExamplePlugin{
        name:        "example-plugin",
        version:     "1.0.0",
        description: "Example plugin for urlDB",
        author:      "urlDB Team",
    }
}

// Name 返回插件名称
func (p *ExamplePlugin) Name() string {
    return p.name
}

// Version 返回插件版本
func (p *ExamplePlugin) Version() string {
    return p.version
}

// Description 返回插件描述
func (p *ExamplePlugin) Description() string {
    return p.description
}

// Author 返回插件作者
func (p *ExamplePlugin) Author() string {
    return p.author
}

// Initialize 初始化插件
func (p *ExamplePlugin) Initialize(ctx types.PluginContext) error {
    p.context = ctx
    p.context.LogInfo("Example plugin initialized")
    return nil
}

// Start 启动插件
func (p *ExamplePlugin) Start() error {
    p.context.LogInfo("Example plugin started")

    // 注册定时任务
    return p.context.RegisterTask("example-task", p.scheduledTask)
}

// Stop 停止插件
func (p *ExamplePlugin) Stop() error {
    p.context.LogInfo("Example plugin stopped")
    return p.context.UnregisterTask("example-task")
}

// Cleanup 清理插件
func (p *ExamplePlugin) Cleanup() error {
    p.context.LogInfo("Example plugin cleaned up")
    return nil
}

// Dependencies 返回依赖列表
func (p *ExamplePlugin) Dependencies() []string {
    return []string{}
}

// CheckDependencies 检查依赖
func (p *ExamplePlugin) CheckDependencies() map[string]bool {
    return make(map[string]bool)
}

// scheduledTask 定时任务
func (p *ExamplePlugin) scheduledTask() {
    p.context.LogInfo("Executing scheduled task at %s", time.Now().Format(time.RFC3339))

    // 示例：获取配置
    interval, err := p.context.GetConfig("interval")
    if err != nil {
        p.context.LogDebug("Using default interval")
        interval = 60 // 默认60秒
    }

    p.context.LogInfo("Task interval: %v seconds", interval)

    // 示例：保存数据
    data := map[string]interface{}{
        "last_run": time.Now(),
        "status":   "success",
    }

    err = p.context.SetData("last_task", data, "task_history")
    if err != nil {
        p.context.LogError("Failed to save task data: %v", err)
    }
}
```

通过遵循本指南，您可以成功开发、编译、加载和使用urlDB插件系统中的插件。