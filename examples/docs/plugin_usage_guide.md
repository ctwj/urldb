# 插件系统的编译注册使用说明

## 1. 环境准备

### 1.1 系统要求
- Go 1.23.0 或更高版本
- 支持 Go plugin 系统的操作系统（Linux, macOS）
- PostgreSQL 数据库（用于插件数据存储）

### 1.2 项目依赖
```bash
go mod tidy
```

## 2. 创建插件

### 2.1 创建插件目录结构
```
my-plugin/
├── go.mod
├── plugin.go
└── main.go (可选，用于构建)
```

### 2.2 创建 go.mod 文件
```go
module github.com/your-org/your-project/plugins/my-plugin

go 1.23.0

replace github.com/ctwj/urldb => ../../../

require github.com/ctwj/urldb v0.0.0
```

### 2.3 实现插件接口

创建 `plugin.go` 文件：

```go
package main

import (
    "github.com/ctwj/urldb/plugin/types"
)

// MyPlugin 插件结构体
type MyPlugin struct {
    name        string
    version     string
    description string
    author      string
    dependencies []string
    context     types.PluginContext
}

// NewMyPlugin 创建新的插件实例
func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        name:        "my-plugin",
        version:     "1.0.0",
        description: "My custom plugin",
        author:      "Your Name",
        dependencies: []string{},
    }
}

// 实现插件接口方法
func (p *MyPlugin) Name() string {
    return p.name
}

func (p *MyPlugin) Version() string {
    return p.version
}

func (p *MyPlugin) Description() string {
    return p.description
}

func (p *MyPlugin) Author() string {
    return p.author
}

func (p *MyPlugin) Dependencies() []string {
    return p.dependencies
}

func (p *MyPlugin) CheckDependencies() map[string]bool {
    result := make(map[string]bool)
    for _, dep := range p.dependencies {
        result[dep] = true
    }
    return result
}

func (p *MyPlugin) Initialize(ctx types.PluginContext) error {
    p.context = ctx
    ctx.LogInfo("MyPlugin initialized")
    return nil
}

func (p *MyPlugin) Start() error {
    p.context.LogInfo("MyPlugin started")
    // 启动插件的主要功能
    return nil
}

func (p *MyPlugin) Stop() error {
    p.context.LogInfo("MyPlugin stopped")
    return nil
}

func (p *MyPlugin) Cleanup() error {
    p.context.LogInfo("MyPlugin cleaned up")
    return nil
}

// 导出插件实例（用于二进制插件）
var Plugin = NewMyPlugin()
```

## 3. 编译插件

### 3.1 编译为二进制插件（.so 文件）

```bash
# 在插件目录下编译
go build -buildmode=plugin -o my-plugin.so

# 或者，如果使用不同平台
# Linux: go build -buildmode=plugin -o my-plugin.so
# macOS: go build -buildmode=plugin -o my-plugin.so
# Windows: go build -buildmode=c-shared -o my-plugin.dll (不完全支持)
```

### 3.2 注意事项
- Windows 不完全支持 Go plugin 系统
- 需要确保插件与主程序使用相同版本的依赖
- 插件必须使用 `buildmode=plugin` 进行编译

## 4. 注册插件

### 4.1 源代码插件注册

将插件作为源代码集成到项目中：

```go
package main

import (
    "github.com/ctwj/urldb/plugin"
    "github.com/your-org/your-project/plugins/my-plugin"
)

func main() {
    // 初始化插件系统
    plugin.InitPluginSystem(taskManager, repoManager)

    // 创建并注册插件
    myPlugin := my_plugin.NewMyPlugin()  // 注意：下划线是包名的一部分
    if err := plugin.RegisterPlugin(myPlugin); err != nil {
        log.Fatal("Failed to register plugin:", err)
    }

    // 初始化插件
    plugin.GetManager().InitializePlugin("my-plugin", config)

    // 启动插件
    plugin.GetManager().StartPlugin("my-plugin")
}
```

### 4.2 二进制插件注册

使用插件加载器从文件系统加载二进制插件：

```go
package main

import (
    "github.com/ctwj/urldb/plugin"
    "github.com/ctwj/urldb/plugin/loader"
)

func main() {
    // 初始化插件系统
    plugin.InitPluginSystem(taskManager, repoManager)

    // 创建插件加载器
    pluginLoader := loader.NewSimplePluginLoader("./plugins")

    // 加载所有插件
    plugins, err := pluginLoader.LoadAllPlugins()
    if err != nil {
        log.Printf("Failed to load plugins: %v", err)
    } else {
        // 注册加载的插件
        for _, p := range plugins {
            if err := plugin.RegisterPlugin(p); err != nil {
                log.Printf("Failed to register plugin %s: %v", p.Name(), err)
            } else {
                log.Printf("Successfully registered plugin: %s", p.Name())
            }
        }
    }
}
```

## 5. 插件配置

### 5.1 创建插件配置
插件可以接收配置参数：

```go
config := map[string]interface{}{
    "interval": 30,     // 30秒间隔
    "enabled": true,    // 启用插件
    "custom_param": "value",
}

// 初始化插件时传入配置
plugin.GetManager().InitializePlugin("my-plugin", config)
```

### 5.2 在插件中使用配置

```go
func (p *MyPlugin) Initialize(ctx types.PluginContext) error {
    p.context = ctx
    ctx.LogInfo("MyPlugin initialized")

    // 获取配置
    if interval, err := ctx.GetConfig("interval"); err == nil {
        if intVal, ok := interval.(int); ok {
            p.context.LogInfo("Interval set to: %d seconds", intVal)
        }
    }

    return nil
}
```

## 6. 插件管理

### 6.1 启动插件

```go
// 启动单个插件
if err := plugin.GetManager().StartPlugin("my-plugin"); err != nil {
    log.Printf("Failed to start plugin: %v", err)
}

// 启动所有插件
plugins, _ := plugin.GetManager().GetAllPlugins()
for _, name := range plugins {
    plugin.GetManager().StartPlugin(name)
}
```

### 6.2 停止插件

```go
// 停止单个插件
if err := plugin.GetManager().StopPlugin("my-plugin"); err != nil {
    log.Printf("Failed to stop plugin: %v", err)
}

// 停止所有插件
plugins, _ := plugin.GetManager().GetAllPlugins()
for _, name := range plugins {
    plugin.GetManager().StopPlugin(name)
}
```

### 6.3 检查插件状态

```go
// 检查插件是否正在运行
if plugin.GetManager().IsPluginRunning("my-plugin") {
    log.Println("Plugin is running")
} else {
    log.Println("Plugin is not running")
}

// 获取所有插件信息
status := plugin.GetManager().GetAllPluginStatus()
for name, info := range status {
    log.Printf("Plugin %s: %s", name, info.Status)
}
```

## 7. 插件依赖管理

### 7.1 定义插件依赖

```go
type MyPlugin struct {
    dependencies []string
}

func (p *MyPlugin) Dependencies() []string {
    return []string{"dependency-plugin-1", "dependency-plugin-2"}
}

func (p *MyPlugin) CheckDependencies() map[string]bool {
    result := make(map[string]bool)

    // 检查依赖是否满足
    for _, dep := range p.dependencies {
        // 检查依赖插件是否存在且已启动
        result[dep] = plugin.GetManager().IsPluginRunning(dep)
    }

    return result
}
```

### 7.2 验证依赖

```go
// 验证所有插件依赖
if err := plugin.GetManager().ValidateDependencies(); err != nil {
    log.Printf("Dependency validation failed: %v", err)
}

// 检查特定插件依赖
ok, unresolved, err := plugin.GetManager().CheckPluginDependencies("my-plugin")
if !ok {
    log.Printf("Unresolved dependencies: %v", unresolved)
}
```

## 8. 实际使用示例

### 8.1 完整的插件使用示例

```go
package main

import (
    "log"
    "time"

    "github.com/ctwj/urldb/db"
    "github.com/ctwj/urldb/db/repo"
    "github.com/ctwj/urldb/plugin"
    "github.com/ctwj/urldb/plugin/loader"
    "github.com/ctwj/urldb/task"
    "github.com/ctwj/urldb/utils"
)

func main() {
    // 初始化日志
    utils.InitLogger(nil)

    // 初始化数据库
    if err := db.InitDB(); err != nil {
        utils.Fatal("Failed to initialize database: %v", err)
    }

    // 创建管理器
    repoManager := repo.NewRepositoryManager(db.DB)
    taskManager := task.NewTaskManager(repoManager)

    // 初始化插件系统
    plugin.InitPluginSystem(taskManager, repoManager)

    // 加载二进制插件
    loadBinaryPlugins()

    // 注册源代码插件
    registerSourcePlugins()

    // 等待运行
    log.Println("Plugin system ready. Running for 2 minutes...")
    time.Sleep(2 * time.Minute)

    // 停止插件
    stopAllPlugins()
}

func loadBinaryPlugins() {
    pluginLoader := loader.NewSimplePluginLoader("./plugins")

    if plugins, err := pluginLoader.LoadAllPlugins(); err == nil {
        for _, p := range plugins {
            if err := plugin.RegisterPlugin(p); err != nil {
                log.Printf("Failed to register binary plugin %s: %v", p.Name(), err)
            }
        }
    }
}

func registerSourcePlugins() {
    // 注册源代码插件
    // myPlugin := my_plugin.NewMyPlugin()
    // plugin.RegisterPlugin(myPlugin)
}

func stopAllPlugins() {
    plugins, _ := plugin.GetManager().GetAllPlugins()
    for _, name := range plugins {
        if err := plugin.GetManager().StopPlugin(name); err != nil {
            log.Printf("Failed to stop plugin %s: %v", name, err)
        }
    }
}
```

## 9. 构建和部署

### 9.1 编译主程序

```bash
cd examples/plugin_demo
go build -o plugin_demo main.go
./plugin_demo
```

### 9.2 部署插件

1. 将编译好的 `.so` 插件文件复制到插件目录
2. 确保主程序有读取插件文件的权限
3. 根据需要配置插件参数

```bash
# 创建插件目录
mkdir -p plugins

# 复制插件
cp my-plugin.so plugins/

# 运行主程序
./plugin_demo
```

## 10. 常见问题和解决方案

### 10.1 插件加载失败
- 检查 `.so` 文件是否与主程序使用相同版本的 Go 编译
- 确保插件依赖的库与主程序兼容
- 检查插件文件权限

### 10.2 依赖关系问题
- 确保插件依赖的其他插件已正确注册
- 检查依赖关系循环

### 10.3 运行时错误
- 确保插件实现符合接口要求
- 检查插件初始化参数

## 11. 最佳实践

1. **插件命名**：使用有意义的插件名称
2. **版本管理**：维护插件版本
3. **错误处理**：在插件中添加适当的错误处理
4. **日志记录**：使用插件上下文记录日志
5. **资源清理**：确保在 `Cleanup` 方法中释放资源
6. **安全考虑**：验证输入参数，限制资源使用
7. **文档**：为插件提供使用文档