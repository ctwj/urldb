# 插件示例目录说明

本目录包含了不同类型的插件示例，用于演示和测试插件系统的功能。

## 目录结构

### binary-plugins/
包含编译后的二进制插件示例：
- `demo-plugin-1.dylib` 和 `demo-plugin-2.dylib` - 预编译的动态库插件
- `plugin1/` 和 `plugin2/` - 可编译的 Go 插件源代码和编译产物

### go-plugins/
包含 Go 源代码形式的插件示例：
- `demo/` - 各种功能演示插件的源代码（与 plugin/demo 目录内容相同）

### plugin_demo/
包含插件系统的使用示例。

## 插件类型说明

### 二进制插件 (Binary Plugins)
这些插件已经编译为动态库文件，可以直接加载使用：
- `.dylib` 文件 (macOS)
- `.so` 文件 (Linux)
- `.dll` 文件 (Windows)

### Go 源码插件 (Go Source Plugins)
这些插件以 Go 源代码形式提供，需要编译后才能使用：
- 完整功能演示插件
- 配置管理演示插件
- 依赖管理演示插件
- 安全功能演示插件
- 性能测试演示插件

## 使用说明

### 编译二进制插件
```bash
cd binary-plugins/plugin1
go build -buildmode=plugin -o plugin1.so main.go
```

### 运行插件示例
```bash
cd plugin_demo
go run plugin_demo.go
```

## 注意事项

- `plugin/demo/` 目录包含原始的插件源代码
- `go-plugins/demo/` 目录是 `plugin/demo/` 的副本，用于示例展示
- 两者内容相同，但位于不同位置以满足不同的使用需求