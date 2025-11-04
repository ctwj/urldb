# 插件测试框架文档

## 概述

本文档介绍urlDB插件测试框架的设计和使用方法。该框架提供了一套完整的工具来测试插件的功能、性能和稳定性。

## 框架组件

### 1. 基础测试框架 (`plugin/test/framework.go`)

基础测试框架提供了以下功能：

- `TestPluginContext`: 模拟插件上下文的实现，用于测试插件与系统核心的交互
- `TestPluginManager`: 插件生命周期管理器，用于测试插件的完整生命周期
- 日志记录和断言功能
- 配置和数据存储模拟
- 任务调度模拟
- 缓存系统模拟
- 安全权限模拟
- 并发控制模拟

### 2. 集成测试环境 (`plugin/test/integration.go`)

集成测试环境提供了以下功能：

- `IntegrationTestSuite`: 完整的集成测试套件，包括数据库、仓库管理器、任务管理器等
- `MockPlugin`: 模拟插件实现，用于测试插件管理器的功能
- 各种错误场景的模拟插件
- 依赖关系模拟
- 上下文操作模拟

### 3. 测试报告生成 (`plugin/test/reporting.go`)

测试报告生成器提供了以下功能：

- `TestReport`: 测试报告结构
- `TestReporter`: 测试报告生成器
- `TestingTWrapper`: 与Go测试框架集成的包装器
- `PluginTestHelper`: 插件测试助手，提供专门的插件测试功能

## 使用方法

### 1. 编写单元测试

要为插件编写单元测试，请参考以下示例：

```go
func TestMyPlugin(t *testing.T) {
    plugin := NewMyPlugin()

    // 创建测试上下文
    ctx := test.NewTestPluginContext()

    // 初始化插件
    if err := plugin.Initialize(ctx); err != nil {
        t.Fatalf("Failed to initialize plugin: %v", err)
    }

    // 验证初始化日志
    if !ctx.AssertLogContains(t, "INFO", "Plugin initialized") {
        t.Error("Expected initialization log")
    }

    // 测试其他功能...
}
```

### 2. 编写集成测试

要编写集成测试，请参考以下示例：

```go
func TestMyPluginIntegration(t *testing.T) {
    // 创建集成测试套件
    suite := test.NewIntegrationTestSuite()
    suite.Setup(t)
    defer suite.Teardown()

    // 注册插件
    plugin := NewMyPlugin()
    if err := suite.RegisterPlugin(plugin); err != nil {
        t.Fatalf("Failed to register plugin: %v", err)
    }

    // 运行集成测试
    config := map[string]interface{}{
        "setting1": "value1",
    }

    suite.RunPluginIntegrationTest(t, plugin.Name(), config)
}
```

### 3. 生成测试报告

测试报告会自动生成，也可以手动创建：

```go
func TestWithReporting(t *testing.T) {
    // 创建报告器
    reporter := test.NewTestReporter("MyTestSuite")
    wrapper := test.NewTestingTWrapper(t, reporter)

    // 使用包装器运行测试
    wrapper.Run("MyTest", func(t *testing.T) {
        // 测试代码...
    })

    // 生成报告
    textReport := reporter.GenerateTextReport()
    t.Logf("Test Report:\n%s", textReport)
}
```

## 运行测试

### 运行所有插件测试

```bash
go test ./plugin/...
```

### 运行特定测试

```bash
go test ./plugin/demo/ -v
```

### 生成测试覆盖率报告

```bash
go test ./plugin/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 最佳实践

### 1. 测试插件生命周期

确保测试插件的完整生命周期：

```go
func TestPluginLifecycle(t *testing.T) {
    manager := test.NewTestPluginManager()
    plugin := NewMyPlugin()

    // 注册插件
    manager.RegisterPlugin(plugin)

    // 测试完整生命周期
    config := map[string]interface{}{
        "config_key": "config_value",
    }

    if err := manager.RunPluginLifecycle(t, plugin.Name(), config); err != nil {
        t.Errorf("Plugin lifecycle failed: %v", err)
    }
}
```

### 2. 测试错误处理

确保测试插件在各种错误情况下的行为：

```go
func TestPluginErrorHandling(t *testing.T) {
    // 测试初始化错误
    pluginWithInitError := test.NewIntegrationTestSuite().
        CreateMockPlugin("error-plugin", "1.0.0").
        WithErrorOnInitialize()

    ctx := test.NewTestPluginContext()
    if err := pluginWithInitError.Initialize(ctx); err == nil {
        t.Error("Expected initialize error")
    }
}
```

### 3. 测试依赖关系

测试插件的依赖关系处理：

```go
func TestPluginDependencies(t *testing.T) {
    plugin := test.NewIntegrationTestSuite().
        CreateMockPlugin("dep-plugin", "1.0.0").
        WithDependencies([]string{"dep1", "dep2"})

    deps := plugin.Dependencies()
    if len(deps) != 2 {
        t.Errorf("Expected 2 dependencies, got %d", len(deps))
    }
}
```

### 4. 测试上下文操作

测试插件与系统上下文的交互：

```go
func TestPluginContextOperations(t *testing.T) {
    operations := []string{
        "log_info",
        "set_config",
        "get_data",
    }

    plugin := test.NewIntegrationTestSuite().
        CreateMockPlugin("context-plugin", "1.0.0").
        WithContextOperations(operations)

    ctx := test.NewTestPluginContext()
    plugin.Initialize(ctx)
    plugin.Start()

    // 验证操作结果
    if !ctx.AssertLogContains(t, "INFO", "Info message") {
        t.Error("Expected info log")
    }
}
```

## 扩展框架

### 添加新的测试功能

要扩展测试框架，可以：

1. 在`TestPluginContext`中添加新的模拟方法
2. 在`TestPluginManager`中添加新的测试辅助方法
3. 在`TestReporter`中添加新的报告功能

### 自定义报告格式

要创建自定义报告格式，可以：

1. 扩展`TestReport`结构
2. 创建新的报告生成方法
3. 实现特定的报告输出格式

## 故障排除

### 常见问题

1. **测试失败但没有错误信息**
   - 检查是否正确使用了测试断言
   - 确保测试上下文正确配置

2. **集成测试环境设置失败**
   - 检查数据库连接配置
   - 确保所有依赖服务可用

3. **测试报告不完整**
   - 确保正确使用了测试报告器
   - 检查测试是否正常完成

### 调试技巧

1. 使用`-v`标志运行测试以获取详细输出
2. 在测试中添加日志记录以跟踪执行流程
3. 使用测试报告来分析测试执行情况