# 插件卸载机制

本文档详细说明了urlDB插件系统的卸载机制，包括如何安全地卸载插件、清理相关数据和处理依赖关系。

## 插件卸载流程

插件卸载是一个多步骤的过程，确保插件被安全地停止并清理所有相关资源。

### 1. 依赖检查

在卸载插件之前，系统会检查是否有其他插件依赖于该插件：

```go
dependents := pluginManager.GetDependents(pluginName)
if len(dependents) > 0 {
    // 有依赖插件，不能安全卸载（除非使用强制模式）
}
```

### 2. 停止插件

如果插件正在运行，系统会先停止它：

```go
if pluginStatus == types.StatusRunning {
    err := pluginManager.StopPlugin(pluginName)
    if err != nil {
        // 处理停止错误
    }
}
```

### 3. 执行插件清理

调用插件的`Cleanup()`方法，让插件执行自定义清理逻辑：

```go
err := plugin.Cleanup()
if err != nil {
    // 处理清理错误
}
```

### 4. 清理数据和配置

系统会自动清理插件的配置和数据：

- 删除插件配置（plugin_config表中的相关记录）
- 删除插件数据（plugin_data表中的相关记录）

### 5. 清理任务

清理插件注册的任何后台任务。

### 6. 更新依赖图

从依赖图中移除插件信息。

## API 使用示例

### 基本卸载

```go
// 非强制卸载（推荐）
err := pluginManager.UninstallPlugin("plugin-name", false)
if err != nil {
    log.Printf("卸载失败: %v", err)
}
```

### 强制卸载

```go
// 强制卸载（即使有错误也继续）
err := pluginManager.UninstallPlugin("plugin-name", true)
if err != nil {
    log.Printf("强制卸载完成，但存在错误: %v", err)
}
```

### 检查卸载安全性

```go
// 检查插件是否可以安全卸载
canUninstall, dependents, err := pluginManager.CanUninstall("plugin-name")
if err != nil {
    log.Printf("检查失败: %v", err)
    return
}

if !canUninstall {
    log.Printf("插件不能安全卸载，依赖插件: %v", dependents)
}
```

## 插件开发者的责任

插件开发者需要实现以下接口方法以支持卸载：

### Cleanup 方法

```go
func (p *MyPlugin) Cleanup() error {
    // 执行清理操作
    // 例如：关闭外部连接、删除临时文件等
    return nil
}
```

### 最佳实践

1. 在`Cleanup`方法中释放所有外部资源
2. 不要在`Cleanup`中删除插件核心文件
3. 处理可能的错误情况，尽可能完成清理工作

## 错误处理

卸载过程中可能出现的错误：

1. 插件不存在
2. 存在依赖插件
3. 停止插件失败
4. 插件清理失败
5. 数据清理失败

使用非强制模式时，任何错误都会导致卸载失败。使用强制模式时，即使出现错误也会继续执行卸载过程。

## 示例代码

完整的卸载示例：

```go
// 检查是否可以安全卸载
canUninstall, dependents, err := pluginManager.CanUninstall("my-plugin")
if err != nil {
    log.Printf("检查失败: %v", err)
    return
}

if !canUninstall {
    log.Printf("插件不能安全卸载，依赖插件: %v", dependents)
    return
}

// 执行卸载
err = pluginManager.UninstallPlugin("my-plugin", false)
if err != nil {
    log.Printf("卸载失败: %v", err)
    return
}

log.Println("插件卸载成功")
```

## 注意事项

1. 卸载是不可逆操作，执行后插件数据将被永久删除
2. 建议在卸载前备份重要数据
3. 强制卸载可能会导致数据不一致，应谨慎使用
4. 卸载后需要重启相关服务才能重新安装同名插件