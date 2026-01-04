# URLDB 插件系统 - 如何添加 JavaScript 变量和函数

本文档详细介绍如何在 URLDB 插件系统的 JavaScript 环境中添加新的 Go 变量和函数。

## 📋 目录

- [系统架构](#系统架构)
- [添加步骤](#添加步骤)
- [实践示例](#实践示例)
- [数据类型转换](#数据类型转换)
- [错误处理](#错误处理)
- [TypeScript 支持](#typescript-支持)
- [最佳实践](#最佳实践)

## 🏗️ 系统架构

URLDB 插件系统使用 **Goja** JavaScript 引擎将 Go 函数暴露给 JavaScript 环境。

### 核心文件结构

```
plugin/jsvm/
├── binds.go      # JavaScript 函数绑定定义
├── runtime.go    # 运行时初始化
├── pool.go       # VM 实例池管理
└── ...
```

### 绑定流程

1. **Go 函数定义** → 2. **注册到 VM** → 3. **JavaScript 调用** → 4. **Go 执行** → 5. **返回结果**

## 📝 添加步骤

### 步骤 1: 在 `binds.go` 中定义函数

在 `/Users/kerwin/Program/go/urldb/plugin/jsvm/binds.go` 中添加新的绑定函数：

```go
// 示例：添加配置相关绑定
func configBinds(vm *goja.Runtime, repoManager *repo.RepositoryManager) {
    // 获取插件配置函数
    vm.Set("getPluginConfig", func(pluginName string) goja.Value {
        // 从数据库查询插件配置
        config, err := repoManager.PluginConfigRepository.GetPluginConfig(pluginName)
        if err != nil {
            utils.Error("Failed to get plugin config: %v", err)
            return vm.ToValue(nil)
        }

        // 解析配置 JSON
        var configData interface{}
        if err := json.Unmarshal([]byte(config.Config), &configData); err != nil {
            utils.Error("Failed to parse config JSON: %v", err)
            return vm.ToValue(nil)
        }

        return vm.ToValue(configData)
    })

    // 设置插件配置函数
    vm.Set("setPluginConfig", func(pluginName string, configData goja.Value) error {
        // 将 JavaScript 数据转换为 JSON
        jsonData, err := json.Marshal(configData.Export())
        if err != nil {
            return fmt.Errorf("failed to marshal config: %v", err)
        }

        // 保存到数据库
        return repoManager.PluginConfigRepository.SetPluginConfig(pluginName, string(jsonData))
    })
}
```

### 步骤 2: 在 `runtime.go` 中注册绑定

在 `/Users/kerwin/Program/go/urldb/plugin/jsvm/runtime.go` 的 `sharedBinds()` 函数中调用新绑定：

```go
func sharedBinds(vm *goja.Runtime, app core.App, executors *vmsPool, repoManager *repo.RepositoryManager, routeRegister func(method, path string, handler func() (interface{}, error)) error) {
    // 现有绑定...
    baseBinds(vm)
    dbxBinds(vm)
    securityBinds(vm)
    osBinds(vm)

    // 添加新的配置绑定
    configBinds(vm, repoManager)

    // 需要传递应用实例的绑定...
    hooksBinds(app, vm, executors)
    cronBinds(app, vm, executors)
    routerBinds(app, vm, executors, routeRegister)
}
```

### 步骤 3: 添加 TypeScript 声明

在 `/Users/kerwin/Program/go/urldb/pb_data/types.d.ts` 中添加函数声明：

```typescript
// 配置相关函数声明
declare function getPluginConfig(pluginName: string): Record<string, any> | null;
declare function setPluginConfig(pluginName: string, config: Record<string, any>): void;

// 全局变量声明
declare const $app: App;
declare const __hooks: string;
```

### 步骤 4: 重新编译和测试

```bash
# 重新编译
go build -o urldb .

# 重启服务
./urldb
```

## 🎯 实践示例

### 示例 1: 简单函数绑定

**Go 代码：**
```go
// 在 binds.go 中添加
vm.Set("getSystemInfo", func() map[string]interface{} {
    return map[string]interface{}{
        "version": "1.0.0",
        "goVersion": runtime.Version(),
        "os": runtime.GOOS,
        "arch": runtime.GOARCH,
    }
})
```

**JavaScript 调用：**
```javascript
// 在插件中使用
const info = getSystemInfo();
console.log("系统信息:", JSON.stringify(info, null, 2));
```

### 示例 2: 带参数和错误处理的函数

**Go 代码：**
```go
vm.Set("readFileSafe", func(filename string) goja.Value {
    // 参数验证
    if filename == "" {
        utils.Error("Filename cannot be empty")
        return vm.ToValue(map[string]interface{}{
            "success": false,
            "error": "Filename cannot be empty",
        })
    }

    // 路径安全检查
    if strings.Contains(filename, "..") {
        utils.Error("Path traversal attempt blocked")
        return vm.ToValue(map[string]interface{}{
            "success": false,
            "error": "Invalid filename",
        })
    }

    // 读取文件
    content, err := os.ReadFile(filename)
    if err != nil {
        return vm.ToValue(map[string]interface{}{
            "success": false,
            "error": err.Error(),
        })
    }

    return vm.ToValue(map[string]interface{}{
        "success": true,
        "content": string(content),
    })
})
```

**JavaScript 调用：**
```javascript
const result = readFileSafe("test.txt");
if (result.success) {
    console.log("文件内容:", result.content);
} else {
    console.error("读取失败:", result.error);
}
```

### 示例 3: 异步回调函数

**Go 代码：**
```go
vm.Set("asyncOperation", func(callback goja.Value) {
    if fn, ok := goja.AssertFunction(callback); ok {
        // 在 goroutine 中执行异步操作
        go func() {
            time.Sleep(2 * time.Second) // 模拟耗时操作

            // 调用 JavaScript 回调
            _, err := fn(goja.Undefined(), vm.ToValue("异步操作完成"))
            if err != nil {
                utils.Error("Callback error: %v", err)
            }
        }()
    }
})
```

**JavaScript 调用：**
```javascript
asyncOperation(function(result) {
    console.log("回调结果:", result);
});
```

## 🔄 数据类型转换

### Go 到 JavaScript 的转换

| Go 类型 | JavaScript 类型 | 示例 |
|---------|----------------|------|
| `string` | `string` | `"hello"` |
| `int/int64` | `number` | `42` |
| `float64` | `number` | `3.14` |
| `bool` | `boolean` | `true` |
| `map[string]interface{}` | `Object` | `{key: "value"}` |
| `[]interface{}` | `Array` | `[1, 2, 3]` |
| `nil` | `null` | `null` |

### JavaScript 到 Go 的转换

使用 `goja.Value.Export()` 方法：

```go
vm.Set("processData", func(data goja.Value) {
    // 导出为 Go 类型
    exported := data.Export()

    switch v := exported.(type) {
    case string:
        fmt.Println("字符串:", v)
    case float64:
        fmt.Println("数字:", v)
    case map[string]interface{}:
        fmt.Println("对象:", v)
    case []interface{}:
        fmt.Println("数组:", v)
    default:
        fmt.Println("未知类型:", v)
    }
})
```

## ⚠️ 错误处理

### 1. 函数级错误处理

```go
vm.Set("safeDivide", func(a, b float64) goja.Value {
    if b == 0 {
        return vm.ToValue(map[string]interface{}{
            "error": "Division by zero",
            "success": false,
        })
    }

    result := a / b
    return vm.ToValue(map[string]interface{}{
        "result": result,
        "success": true,
    })
})
```

### 2. 捕获 JavaScript 异常

```go
vm.Set("executeJS", func(code string) goja.Value {
    value, err := vm.RunString(code)
    if err != nil {
        return vm.ToValue(map[string]interface{}{
            "error": err.Error(),
            "success": false,
        })
    }

    return vm.ToValue(map[string]interface{}{
        "result": value.Export(),
        "success": true,
    })
})
```

### 3. 使用 panic 恢复

```go
vm.Set("riskyOperation", func() goja.Value {
    defer func() {
        if r := recover(); r != nil {
            utils.Error("Panic recovered in riskyOperation: %v", r)
        }
    }()

    // 可能发生 panic 的代码
    // ...

    return vm.ToValue("操作成功")
})
```

## 🔧 TypeScript 支持

### 完整的类型声明示例

```typescript
// types.d.ts

declare global {
    // 系统信息函数
    declare function getSystemInfo(): {
        version: string;
        goVersion: string;
        os: string;
        arch: string;
    };

    // 文件操作函数
    declare function readFileSafe(filename: string): {
        success: boolean;
        content?: string;
        error?: string;
    };

    // 异步操作函数
    declare function asyncOperation(callback: (result: string) => void): void;

    // 数学工具函数
    declare function safeDivide(a: number, b: number): {
        success: boolean;
        result?: number;
        error?: string;
    };

    // 配置管理函数
    declare function getPluginConfig(pluginName: string): Record<string, any> | null;
    declare function setPluginConfig(pluginName: string, config: Record<string, any>): void;

    // 自定义对象类型
    interface PluginConfig {
        enabled: boolean;
        debug?: boolean;
        log_level?: 'debug' | 'info' | 'warn' | 'error';
        custom_data?: Record<string, any>;
    }
}

export {};
```

## 💡 最佳实践

### 1. 命名规范

```go
// ✅ 好的命名
vm.Set("getPluginConfig", ...)
vm.Set("readFileSafe", ...)
vm.Set("calculateHash", ...)

// ❌ 避免的命名
vm.Set("gpc", ...)           // 缩写不清晰
vm.Set("read_file", ...)     // 下划线不符合 JS 命名规范
vm.Set("internalFunc", ...)  // 不要暴露内部函数
```

### 2. 参数验证

```go
vm.Set("processUserInput", func(input string) goja.Value {
    // 输入验证
    if len(input) == 0 {
        return vm.ToValue(map[string]interface{}{
            "error": "Input cannot be empty",
            "code": "INVALID_INPUT",
        })
    }

    if len(input) > 1000 {
        return vm.ToValue(map[string]interface{}{
            "error": "Input too long (max 1000 chars)",
            "code": "INPUT_TOO_LONG",
        })
    }

    // 处理逻辑...

    return vm.ToValue(map[string]interface{}{
        "success": true,
        "result": processedInput,
    })
})
```

### 3. 性能优化

```go
// ✅ 使用对象池减少分配
vm.Set("processBulkData", func(data goja.Value) goja.Value {
    // 批量处理而不是逐个处理
    dataArray := data.Export().([]interface{})
    results := make([]interface{}, 0, len(dataArray))

    for _, item := range dataArray {
        // 处理每个项目
        result := processItem(item)
        results = append(results, result)
    }

    return vm.ToValue(results)
})

// ❌ 避免在循环中创建大量临时对象
vm.Set("badExample", func(items goja.Value) {
    for i := 0; i < 1000; i++ {
        // 每次循环都创建新对象，性能差
        temp := map[string]interface{}{
            "index": i,
            "data": expensiveOperation(),
        }
        // ...
    }
})
```

### 4. 安全考虑

```go
vm.Set("executeCommand", func(command string) goja.Value {
    // 安全检查：禁止危险命令
    dangerousCommands := []string{
        "rm", "format", "del", "shutdown", "reboot",
        "> /dev/", "mkfs", "fdisk", "mount", "umount",
    }

    for _, dangerous := range dangerousCommands {
        if strings.Contains(command, dangerous) {
            return vm.ToValue(map[string]interface{}{
                "error": "Dangerous command blocked",
                "success": false,
            })
        }
    }

    // 执行安全的命令...

    return vm.ToValue(map[string]interface{}{
        "success": true,
        "result": commandResult,
    })
})
```

### 5. 日志记录

```go
vm.Set("debugFunction", func(param string) goja.Value {
    utils.Debug("debugFunction called with param: %s", param)

    // 执行逻辑
    result := processData(param)

    utils.Info("debugFunction completed successfully")
    return vm.ToValue(result)
})
```

## 🚀 高级用法

### 1. 返回 JavaScript 函数

```go
vm.Set("createCounter", func() goja.Value {
    count := 0

    // 返回一个 JavaScript 函数
    counterFn := vm.ToValue(func() int {
        count++
        return count
    })

    return counterFn
})
```

### 2. 复杂数据结构

```go
type ComplexStruct struct {
    ID       int                    `json:"id"`
    Name     string                 `json:"name"`
    Metadata map[string]interface{} `json:"metadata"`
    Created  time.Time              `json:"created"`
}

vm.Set("getComplexData", func() goja.Value {
    data := ComplexStruct{
        ID:   1,
        Name: "Test Data",
        Metadata: map[string]interface{}{
            "tags":  []string{"tag1", "tag2"},
            "count": 42,
        },
        Created: time.Now(),
    }

    return vm.ToValue(data)
})
```

### 3. 流式处理

```go
vm.Set("streamData", func(callback goja.Value) {
    if fn, ok := goja.AssertFunction(callback); ok {
        go func() {
            for i := 0; i < 10; i++ {
                // 发送数据块
                _, err := fn(goja.Undefined(), vm.ToValue(map[string]interface{}{
                    "chunk": i,
                    "data":  fmt.Sprintf("Data chunk %d", i),
                }))

                if err != nil {
                    utils.Error("Stream callback error: %v", err)
                    break
                }

                time.Sleep(500 * time.Millisecond)
            }

            // 发送结束信号
            fn(goja.Undefined(), vm.ToValue(map[string]interface{}{
                "done": true,
            }))
        }()
    }
})
```

## 📚 总结

通过以上步骤和示例，你可以：

1. **定义 Go 函数**并绑定到 JavaScript 环境
2. **处理各种数据类型**的转换
3. **实现错误处理**和异常恢复
4. **添加 TypeScript 支持**以获得更好的开发体验
5. **遵循最佳实践**确保代码质量和安全性

记住关键原则：
- **安全第一**：验证所有输入
- **错误处理**：优雅处理所有异常情况
- **性能考虑**：避免不必要的资源消耗
- **类型安全**：使用 TypeScript 提供类型检查
- **文档完整**：为所有函数提供清晰的文档

这样就能构建一个强大、安全、易用的插件系统 JavaScript API！