# 插件依赖管理功能说明

本文档详细说明了urlDB插件系统的依赖管理功能，包括依赖解析和验证机制以及依赖加载顺序管理。

## 1. 设计概述

插件依赖管理功能允许插件声明其依赖关系，系统会在插件初始化和启动时自动验证这些依赖关系，并确保插件按照正确的顺序加载。

### 核心组件

1. **DependencyManager**: 负责依赖解析和验证的核心组件
2. **DependencyGraph**: 依赖关系图，用于表示插件间的依赖关系
3. **PluginLoader**: 增强的插件加载器，支持依赖管理
4. **插件接口扩展**: 为插件添加依赖声明和检查方法

## 2. 依赖解析和验证机制

### 2.1 依赖声明

插件通过实现`Dependencies() []string`方法声明其依赖关系：

```go
// Dependencies returns the plugin dependencies
func (p *MyPlugin) Dependencies() []string {
    return []string{"database-plugin", "auth-plugin"}
}
```

### 2.2 依赖检查

插件可以通过实现`CheckDependencies() map[string]bool`方法来检查依赖状态：

```go
// CheckDependencies checks the plugin dependencies
func (p *MyPlugin) CheckDependencies() map[string]bool {
    dependencies := p.Dependencies()
    result := make(map[string]bool)

    for _, dep := range dependencies {
        // 检查依赖是否满足
        result[dep] = isDependencySatisfied(dep)
    }

    return result
}
```

### 2.3 系统级依赖验证

DependencyManager提供了以下验证功能：

1. **依赖存在性验证**: 确保所有声明的依赖都已注册
2. **循环依赖检测**: 检测并防止循环依赖
3. **依赖状态检查**: 验证依赖插件是否正在运行

```go
// ValidateDependencies validates all plugin dependencies
func (dm *DependencyManager) ValidateDependencies() error {
    // 检查所有依赖是否存在
    // 检测循环依赖
    // 返回验证结果
}
```

## 3. 依赖加载顺序管理

### 3.1 拓扑排序

系统使用拓扑排序算法确定插件的加载顺序，确保依赖项在依赖它们的插件之前加载。

```go
// GetLoadOrder returns the correct order to load plugins based on dependencies
func (dm *DependencyManager) GetLoadOrder() ([]string, error) {
    // 构建依赖图
    // 执行拓扑排序
    // 返回加载顺序
}
```

### 3.2 加载流程

1. 构建依赖图
2. 验证依赖关系
3. 执行拓扑排序确定加载顺序
4. 按顺序加载插件

## 4. 使用示例

### 4.1 声明依赖

```go
type MyPlugin struct {
    name        string
    version     string
    dependencies []string
}

func NewMyPlugin() *MyPlugin {
    return &MyPlugin{
        name:        "my-plugin",
        version:     "1.0.0",
        dependencies: []string{"demo-plugin"}, // 声明依赖
    }
}

func (p *MyPlugin) Dependencies() []string {
    return p.dependencies
}
```

### 4.2 检查依赖

```go
func (p *MyPlugin) Start() error {
    // 检查依赖状态
    satisfied, unresolved, err := pluginManager.CheckPluginDependencies(p.Name())
    if err != nil {
        return err
    }

    if !satisfied {
        return fmt.Errorf("unsatisfied dependencies: %v", unresolved)
    }

    // 依赖满足，继续启动
    return nil
}
```

### 4.3 系统级依赖管理

```go
// 验证所有依赖
if err := pluginManager.ValidateDependencies(); err != nil {
    log.Fatalf("Dependency validation failed: %v", err)
}

// 获取加载顺序
loadOrder, err := pluginManager.GetLoadOrder()
if err != nil {
    log.Fatalf("Failed to determine load order: %v", err)
}

// 按顺序加载插件
for _, pluginName := range loadOrder {
    if err := pluginManager.InitializePlugin(pluginName, config); err != nil {
        log.Fatalf("Failed to initialize plugin %s: %v", pluginName, err)
    }
}
```

## 5. API参考

### 5.1 DependencyManager方法

- `ValidateDependencies() error`: 验证所有插件依赖
- `CheckPluginDependencies(pluginName string) (bool, []string, error)`: 检查特定插件的依赖状态
- `GetLoadOrder() ([]string, error)`: 获取插件加载顺序
- `GetDependencyInfo(pluginName string) (*types.PluginInfo, error)`: 获取插件依赖信息
- `CheckAllDependencies() map[string]map[string]bool`: 检查所有插件的依赖状态

### 5.2 PluginLoader方法

- `LoadPluginWithDependencies(pluginName string) error`: 加载插件及其依赖
- `LoadAllPlugins() error`: 按依赖顺序加载所有插件

## 6. 最佳实践

1. **明确声明依赖**: 插件应明确声明所有必需的依赖
2. **避免循环依赖**: 设计时应避免插件间的循环依赖
3. **提供依赖检查**: 实现`CheckDependencies`方法以提供详细的依赖状态
4. **处理依赖失败**: 优雅地处理依赖不满足的情况
5. **测试依赖关系**: 编写测试确保依赖关系正确配置

## 7. 故障排除

### 7.1 常见错误

1. **依赖未找到**: 确保依赖插件已正确注册
2. **循环依赖**: 检查插件依赖关系图，消除循环依赖
3. **依赖未启动**: 确保依赖插件已正确启动并运行

### 7.2 调试工具

使用以下方法调试依赖问题：

```go
// 检查所有依赖状态
allDeps := pluginManager.CheckAllDependencies()
for plugin, deps := range allDeps {
    fmt.Printf("Plugin %s dependencies: %v\n", plugin, deps)
}

// 获取特定插件的依赖信息
info, err := pluginManager.GetDependencyInfo("my-plugin")
if err != nil {
    log.Printf("Failed to get dependency info: %v", err)
} else {
    fmt.Printf("Plugin info: %+v\n", info)
}
```