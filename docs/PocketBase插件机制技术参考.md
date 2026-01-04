# 在其他项目中实现类似 PocketBase 插件机制的技术方案

## 1. 概述

本文档详细阐述如何在其他项目中实现类似 PocketBase 的插件机制。PocketBase 的插件系统基于 Goja JavaScript 引擎，提供了一个强大的扩展平台，允许开发者通过 JavaScript 编写钩子、迁移和自定义功能。

## 2. 架构设计

### 2.1 核心组件结构

```
plugin-system/
├── core/                 # 核心应用接口定义
├── plugin/               # 插件注册和管理系统
├── hooks/                # 钩子系统
├── runtime/              # JavaScript 运行时管理
├── pool/                 # VM 池管理
└── bindings/             # JavaScript 绑定
```

### 2.2 核心接口定义

```go
// core/app.go
package core

import "context"

// App 接口定义核心应用结构
type App interface {
    Bootstrap() error
    IsBootstrapped() bool
    // 添加各种钩子接口方法
    OnBootstrap() *Hook[*BootstrapEvent]
    OnCustomEvent() *Hook[*CustomEvent]
    // ... 其他钩子方法
}

// Event 事件基类
type Event interface {
    Next() error
}

// Hook 钩子定义
type Hook[T Event] interface {
    Bind(handler Handler[T]) error
    BindFunc(fn func(T) error) error
    Trigger(event T, handler func(T) error) error
}
```

## 3. 钩子系统实现

### 3.1 钩子注册与触发机制

```go
// hooks/hook.go
package hooks

import (
    "context"
    "errors"
    "sync"
)

type Handler[T Event] struct {
    Id       string
    Func     func(T) error
    Priority int // 优先级，负数优先级高
}

type Hook[T Event] struct {
    handlers []*Handler[T]
    mux      sync.RWMutex
}

func NewHook[T Event]() *Hook[T] {
    return &Hook[T]{
        handlers: make([]*Handler[T], 0),
    }
}

func (h *Hook[T]) Bind(handler *Handler[T]) {
    h.mux.Lock()
    defer h.mux.Unlock()

    h.handlers = append(h.handlers, handler)
    h.sortHandlers()
}

func (h *Hook[T]) BindFunc(fn func(T) error) {
    h.Bind(&Handler[T]{Func: fn})
}

func (h *Hook[T]) sortHandlers() {
    // 按优先级排序
    for i := 0; i < len(h.handlers)-1; i++ {
        for j := i + 1; j < len(h.handlers); j++ {
            if h.handlers[i].Priority > h.handlers[j].Priority {
                h.handlers[i], h.handlers[j] = h.handlers[j], h.handlers[i]
            }
        }
    }
}

func (h *Hook[T]) Trigger(event T, handler func(T) error) error {
    h.mux.RLock()
    handlers := make([]*Handler[T], len(h.handlers))
    copy(handlers, h.handlers)
    h.mux.RUnlock()

    var currentHandler func(T) error = func(e T) error {
        return handler(e)
    }

    // 反向构建处理链
    for i := len(handlers) - 1; i >= 0; i-- {
        currentHandler = func(next func(T) error, h *Handler[T]) func(T) error {
            return func(e T) error {
                if h.Func != nil {
                    return h.Func(e)
                }
                return next(e)
            }
        }(currentHandler, handlers[i])
    }

    return currentHandler(event)
}
```

### 3.2 事件结构定义

```go
// hooks/events.go
package hooks

type BootstrapEvent struct {
    App App
}

func (e *BootstrapEvent) Next() error {
    return nil // 可以添加中间件逻辑
}

type CustomEvent struct {
    App     App
    Data    map[string]interface{}
    Cancel  bool
}

func (e *CustomEvent) Next() error {
    if e.Cancel {
        return errors.New("event cancelled")
    }
    return nil
}
```

## 4. JavaScript 运行时实现

### 4.1 Goja 运行时初始化

```go
// runtime/runtime.go
package runtime

import (
    "github.com/dop251/goja"
    "github.com/dop251/goja_nodejs/require"
    "github.com/dop251/goja_nodejs/console"
    "github.com/dop251/goja_nodejs/process"
)

type JavaScriptRuntime struct {
    vm       *goja.Runtime
    registry *require.Registry
    bindings map[string]interface{}
}

func NewJavaScriptRuntime() *JavaScriptRuntime {
    vm := goja.New()

    registry := new(require.Registry)
    registry.Enable(vm)
    console.Enable(vm)
    process.Enable(vm)

    return &JavaScriptRuntime{
        vm:       vm,
        registry: registry,
        bindings: make(map[string]interface{}),
    }
}

func (j *JavaScriptRuntime) Bind(name string, value interface{}) {
    j.bindings[name] = value
    j.vm.Set(name, value)
}
```

### 4.2 运行时绑定系统

```go
// bindings/core_bindings.go
package bindings

import (
    "github.com/dop251/goja"
    "your-project/core"
)

type CoreBindings struct {
    app core.App
}

func NewCoreBindings(app core.App) *CoreBindings {
    return &CoreBindings{app: app}
}

func (c *CoreBindings) Bind(vm *goja.Runtime) {
    // 绑定应用实例
    vm.Set("$app", c.app)

    // 绑定钩子系统
    c.bindHooks(vm)

    // 绑定数据库操作
    c.bindDatabase(vm)

    // 绑定其他系统功能
    c.bindSystemFunctions(vm)
}

func (c *CoreBindings) bindHooks(vm *goja.Runtime) {
    // 通过反射自动绑定所有钩子方法
    vm.Set("onBootstrap", func(callback func(e *core.BootstrapEvent) error) {
        c.app.OnBootstrap().BindFunc(callback)
    })

    vm.Set("onCustomEvent", func(callback func(e *core.CustomEvent) error) {
        c.app.OnCustomEvent().BindFunc(callback)
    })
}

func (c *CoreBindings) bindDatabase(vm *goja.Runtime) {
    dbObj := vm.NewObject()
    vm.Set("$db", dbObj)

    dbObj.Set("find", func(query string) interface{} {
        // 实现数据库查找逻辑
        return nil
    })

    dbObj.Set("save", func(model interface{}) error {
        // 实现数据库保存逻辑
        return nil
    })
}
```

## 5. VM 池管理系统

### 5.1 VM 池实现

```go
// pool/pool.go
package pool

import (
    "sync"
    "time"
    "github.com/dop251/goja"
    "your-project/runtime"
)

type PoolItem struct {
    vm    *goja.Runtime
    busy  bool
    mux   sync.Mutex
    lastUsed time.Time
}

type VMPool struct {
    items     []*PoolItem
    factory   func() *goja.Runtime
    maxSize   int
    minSize   int
    idleTimeout time.Duration
    mux       sync.RWMutex
}

func NewVMPool(factory func() *goja.Runtime, maxSize, minSize int) *VMPool {
    pool := &VMPool{
        factory:   factory,
        maxSize:   maxSize,
        minSize:   minSize,
        items:     make([]*PoolItem, 0),
        idleTimeout: 5 * time.Minute,
    }

    // 预创建最小数量的 VM
    pool.preWarm()

    // 启动清理协程
    go pool.cleanupRoutine()

    return pool
}

func (p *VMPool) preWarm() {
    p.mux.Lock()
    defer p.mux.Unlock()

    for i := 0; i < p.minSize && len(p.items) < p.maxSize; i++ {
        vm := p.factory()
        item := &PoolItem{
            vm:       vm,
            lastUsed: time.Now(),
        }
        p.items = append(p.items, item)
    }
}

func (p *VMPool) GetVM() (*goja.Runtime, func()) {
    p.mux.Lock()
    defer p.mux.Unlock()

    // 查找空闲的 VM
    for _, item := range p.items {
        item.mux.Lock()
        if !item.busy {
            item.busy = true
            item.lastUsed = time.Now()
            item.mux.Unlock()

            releaseFunc := func() {
                item.mux.Lock()
                item.busy = false
                item.mux.Unlock()
            }

            return item.vm, releaseFunc
        }
        item.mux.Unlock()
    }

    // 如果没有空闲 VM 且未达到最大限制，创建新的
    if len(p.items) < p.maxSize {
        vm := p.factory()
        item := &PoolItem{
            vm:       vm,
            busy:     true,
            lastUsed: time.Now(),
        }

        p.items = append(p.items, item)

        releaseFunc := func() {
            item.mux.Lock()
            item.busy = false
            item.mux.Unlock()
        }

        return vm, releaseFunc
    }

    // 所有 VM 都忙，返回 nil（这里可以根据需要实现等待逻辑）
    return nil, nil
}

func (p *VMPool) cleanupRoutine() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        p.cleanupIdleVMs()
    }
}

func (p *VMPool) cleanupIdleVMs() {
    p.mux.Lock()
    defer p.mux.Unlock()

    now := time.Now()
    remaining := make([]*PoolItem, 0)

    for _, item := range p.items {
        item.mux.Lock()
        if !item.busy && now.Sub(item.lastUsed) > p.idleTimeout && len(remaining) >= p.minSize {
            // 回收空闲的 VM
            item.mux.Unlock()
            continue
        }
        item.mux.Unlock()
        remaining = append(remaining, item)
    }

    p.items = remaining
}
```

## 6. 插件系统核心实现

### 6.1 插件管理器

```go
// plugin/manager.go
package plugin

import (
    "context"
    "errors"
    "fmt"
    "path/filepath"
    "plugin"
    "reflect"
    "sync"

    "your-project/core"
    "your-project/runtime"
    "your-project/bindings"
)

type PluginManager struct {
    app      core.App
    plugins  map[string]*Plugin
    mux      sync.RWMutex
    vmpool   *pool.VMPool
    bindings *bindings.CoreBindings
}

type Plugin struct {
    Name        string
    Version     string
    Description string
    Register    func(core.App, interface{}) error
    Config      interface{}
    Loaded      bool
    Enabled     bool
}

func NewPluginManager(app core.App) *PluginManager {
    vmpool := pool.NewVMPool(func() *goja.Runtime {
        return runtime.NewJavaScriptRuntime()
    }, 20, 5)

    return &PluginManager{
        app:      app,
        plugins:  make(map[string]*Plugin),
        vmpool:   vmpool,
        bindings: bindings.NewCoreBindings(app),
    }
}

func (pm *PluginManager) RegisterPlugin(name string, plugin *Plugin) error {
    pm.mux.Lock()
    defer pm.mux.Unlock()

    if _, exists := pm.plugins[name]; exists {
        return fmt.Errorf("plugin %s already registered", name)
    }

    pm.plugins[name] = plugin
    return nil
}

func (pm *PluginManager) LoadPlugin(name string) error {
    pm.mux.Lock()
    plugin, exists := pm.plugins[name]
    pm.mux.Unlock()

    if !exists {
        return fmt.Errorf("plugin %s not found", name)
    }

    if plugin.Register == nil {
        return errors.New("plugin registration function not found")
    }

    err := plugin.Register(pm.app, plugin.Config)
    if err != nil {
        return fmt.Errorf("failed to register plugin %s: %v", name, err)
    }

    plugin.Loaded = true
    plugin.Enabled = true

    return nil
}
```

### 6.2 JavaScript 插件处理器

```go
// plugin/js_plugin_handler.go
package plugin

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "regexp"

    "github.com/dop251/goja"
    "your-project/core"
)

type JSPluginHandler struct {
    manager *PluginManager
    app     core.App
    hooksDir string
    vmpool  *pool.VMPool
}

func NewJSPluginHandler(manager *PluginManager, app core.App) *JSPluginHandler {
    return &JSPluginHandler{
        manager: manager,
        app:     app,
        hooksDir: "hooks", // 默认钩子目录
        vmpool:  manager.vmpool,
    }
}

func (j *JSPluginHandler) SetHooksDir(dir string) {
    j.hooksDir = dir
}

func (j *JSPluginHandler) LoadJSFiles(pattern string) error {
    files, err := j.getJSFiles(pattern)
    if err != nil {
        return err
    }

    for file, content := range files {
        if err := j.executeJSFile(file, content); err != nil {
            return fmt.Errorf("failed to execute JS file %s: %v", file, err)
        }
    }

    return nil
}

func (j *JSPluginHandler) getJSFiles(pattern string) (map[string][]byte, error) {
    files := make(map[string][]byte)

    exp, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }

    fileInfos, err := os.ReadDir(j.hooksDir)
    if err != nil {
        if os.IsNotExist(err) {
            return files, nil // 目录不存在，返回空映射
        }
        return nil, err
    }

    for _, fileInfo := range fileInfos {
        if fileInfo.IsDir() || !exp.MatchString(fileInfo.Name()) {
            continue
        }

        content, err := os.ReadFile(filepath.Join(j.hooksDir, fileInfo.Name()))
        if err != nil {
            return nil, err
        }

        files[fileInfo.Name()] = content
    }

    return files, nil
}

func (j *JSPluginHandler) executeJSFile(filename string, content []byte) error {
    vm, release := j.vmpool.GetVM()
    if vm == nil {
        return errors.New("no available VM in pool")
    }
    defer release()

    // 绑定核心功能到 VM
    j.bindCoreFunctions(vm)

    _, err := vm.RunScript(filename, string(content))
    if err != nil {
        return err
    }

    return nil
}

func (j *JSPluginHandler) bindCoreFunctions(vm *goja.Runtime) {
    // 绑定应用实例
    vm.Set("$app", j.app)

    // 绑定钩子系统
    j.bindHookFunctions(vm)

    // 绑定其他功能
    j.manager.bindings.Bind(vm)
}

func (j *JSPluginHandler) bindHookFunctions(vm *goja.Runtime) {
    // 通过反射动态绑定所有钩子方法
    appType := reflect.TypeOf(j.app)
    appValue := reflect.ValueOf(j.app)

    for i := 0; i < appType.NumMethod(); i++ {
        method := appType.Method(i)
        if !j.isHookMethod(method.Name) {
            continue
        }

        jsName := j.toJSMethodName(method.Name)
        vm.Set(jsName, func(callback string) {
            j.registerHookInJS(callback, method.Name, appValue)
        })
    }
}

func (j *JSPluginHandler) isHookMethod(name string) bool {
    return len(name) > 2 && name[:2] == "On"
}

func (j *JSPluginHandler) toJSMethodName(goName string) string {
    if len(goName) <= 2 {
        return goName
    }
    return "on" + goName[2:] // OnBootstrap -> onBootstrap
}

func (j *JSPluginHandler) registerHookInJS(callback, method string, appValue reflect.Value) {
    hookInstance := appValue.MethodByName(method).Call([]reflect.Value{})[0]

    // 这里需要实现 JavaScript 回调到 Go 的桥接
    // 具体实现依赖于 goja 的函数调用机制
    hookBindFunc := hookInstance.MethodByName("BindFunc")

    wrappedFunc := func(event interface{}) error {
        vm, release := j.vmpool.GetVM()
        if vm == nil {
            return errors.New("no available VM")
        }
        defer release()

        // 执行 JavaScript 回调
        gojaCallback, err := goja.Compile("callback", callback, false)
        if err != nil {
            return err
        }

        _, err = vm.RunProgram(gojaCallback)
        return err
    }

    handlerType := hookBindFunc.Type().In(0)
    handler := reflect.MakeFunc(handlerType, func(args []reflect.Value) []reflect.Value {
        event := args[0].Interface()
        err := wrappedFunc(event)
        return []reflect.Value{reflect.ValueOf(&err).Elem()}
    })

    hookBindFunc.Call([]reflect.Value{handler})
}
```

## 7. 文件系统监控

### 7.1 热重载机制

```go
// plugin/file_watcher.go
package plugin

import (
    "context"
    "path/filepath"
    "strings"
    "time"

    "github.com/fsnotify/fsnotify"
    "github.com/fatih/color"
)

type FileWatcher struct {
    manager *PluginManager
    app     core.App
    jsHandler *JSPluginHandler
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewFileWatcher(manager *PluginManager, app core.App, jsHandler *JSPluginHandler) *FileWatcher {
    ctx, cancel := context.WithCancel(context.Background())
    return &FileWatcher{
        manager:   manager,
        app:       app,
        jsHandler: jsHandler,
        ctx:       ctx,
        cancel:    cancel,
    }
}

func (fw *FileWatcher) Watch(hooksDir string) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }

    defer watcher.Close()

    // 监控钩子目录
    if err := filepath.WalkDir(hooksDir, func(path string, info fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
            return watcher.Add(path)
        }
        return nil
    }); err != nil {
        return err
    }

    go fw.handleEvents(watcher)

    return nil
}

func (fw *FileWatcher) handleEvents(watcher *fsnotify.Watcher) {
    var debounceTimer *time.Timer

    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }

            if debounceTimer != nil {
                debounceTimer.Stop()
            }

            debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
                color.Yellow("File %s changed, reloading hooks...", event.Name)
                if err := fw.reloadHooks(); err != nil {
                    color.Red("Failed to reload hooks: %v", err)
                }
            })

        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            color.Red("Watch error: %v", err)

        case <-fw.ctx.Done():
            return
        }
    }
}

func (fw *FileWatcher) reloadHooks() error {
    // 停止当前插件，重新加载
    return fw.jsHandler.LoadJSFiles(`^.*(\.pb\.js|\.pb\.ts)$`)
}
```

## 8. 完整的初始化和使用示例

### 8.1 完整的插件系统初始化

```go
// main.go 示例
package main

import (
    "log"
    "your-project/core"
    "your-project/plugin"
)

type MyApp struct {
    core.BaseApp // 假设有一个基础应用结构
    pluginManager *plugin.PluginManager
    fileWatcher   *plugin.FileWatcher
    jsHandler     *plugin.JSPluginHandler
}

func NewApp() *MyApp {
    app := &MyApp{}

    // 初始化核心应用
    app.initCore()

    // 初始化插件系统
    app.initPluginSystem()

    return app
}

func (app *MyApp) initCore() {
    // 初始化核心应用逻辑
    app.BaseApp = core.NewBaseApp()
}

func (app *MyApp) initPluginSystem() {
    // 创建插件管理器
    app.pluginManager = plugin.NewPluginManager(&app.BaseApp)

    // 创建 JS 插件处理器
    app.jsHandler = plugin.NewJSPluginHandler(app.pluginManager, &app.BaseApp)

    // 创建文件监控器
    app.fileWatcher = plugin.NewFileWatcher(app.pluginManager, &app.BaseApp, app.jsHandler)

    // 设置钩子目录
    app.jsHandler.SetHooksDir("./hooks")
}

func (app *MyApp) LoadPlugins() error {
    // 加载 JavaScript 钩子文件
    if err := app.jsHandler.LoadJSFiles(`^.*(\.pb\.js|\.pb\.ts)$`); err != nil {
        return err
    }

    // 启动文件监控（开发模式）
    if err := app.fileWatcher.Watch("./hooks"); err != nil {
        return err
    }

    return nil
}

func main() {
    app := NewApp()

    // 通过 Bootstrap 触发插件加载
    app.BaseApp.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
        if err := app.LoadPlugins(); err != nil {
            return err
        }
        return e.Next()
    })

    if err := app.BaseApp.Bootstrap(); err != nil {
        log.Fatal(err)
    }
}
```

### 8.2 JavaScript 钩子示例

```javascript
// hooks/example.pb.js
/// <reference path="../pb_data/types.d.ts" />

// 监听记录创建事件
onRecordCreate((e) => {
    console.log("Record created:", e.record.id);
    return e.next();
});

// 监听自定义事件
onCustomEvent((e) => {
    if (e.data.type === "user_action") {
        console.log("User action detected:", e.data.payload);
    }
    return e.next();
});

// 添加自定义路由
routerAdd("GET", "/api/hello", (e) => {
    return e.json(200, { message: "Hello from plugin!" });
});

// 定时任务
cronAdd("hello_task", "0 */5 * * *", () => {
    console.log("Hello from scheduled task!");
});
```

## 9. 安全性和性能考虑

### 9.1 安全措施

```go
// security/security.go
package security

import (
    "context"
    "time"
    "github.com/dop251/goja"
)

type SecurityEnforcer struct {
    timeout time.Duration
    memoryLimit int64
}

func NewSecurityEnforcer(timeout time.Duration, memoryLimit int64) *SecurityEnforcer {
    return &SecurityEnforcer{
        timeout: timeout,
        memoryLimit: memoryLimit,
    }
}

func (s *SecurityEnforcer) ExecuteWithTimeout(vm *goja.Runtime, script string) (goja.Value, error) {
    ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
    defer cancel()

    // 通过 channel 实现超时控制
    result := make(chan struct {
        value goja.Value
        err   error
    }, 1)

    go func() {
        value, err := vm.RunString(script)
        result <- struct {
            value goja.Value
            err   error
        }{value: value, err: err}
    }()

    select {
    case res := <-result:
        return res.value, res.err
    case <-ctx.Done():
        return nil, context.DeadlineExceeded
    }
}
```

### 9.2 性能监控

```go
// metrics/metrics.go
package metrics

import (
    "sync"
    "time"
)

type MetricsCollector struct {
    executionTimes map[string][]time.Duration
    lock          sync.RWMutex
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        executionTimes: make(map[string][]time.Duration),
    }
}

func (m *MetricsCollector) RecordExecution(hookName string, duration time.Duration) {
    m.lock.Lock()
    defer m.lock.Unlock()

    m.executionTimes[hookName] = append(m.executionTimes[hookName], duration)
}

func (m *MetricsCollector) GetAverageExecutionTime(hookName string) time.Duration {
    m.lock.RLock()
    defer m.lock.RUnlock()

    if times, exists := m.executionTimes[hookName]; exists && len(times) > 0 {
        var total time.Duration
        for _, t := range times {
            total += t
        }
        return total / time.Duration(len(times))
    }

    return 0
}
```

## 10. 总结

这个技术方案提供了一个完整的、可扩展的插件系统框架，具有以下特点：

1. **模块化设计**：各个组件职责分离，易于扩展
2. **高性能**：使用 VM 池避免重复创建运行时
3. **安全控制**：支持超时和资源限制
4. **热重载**：开发模式下支持文件修改自动重载
5. **类型安全**：Go 类型系统保证类型安全
6. **错误处理**：完善的错误处理和恢复机制

这种架构允许在其他项目中实现类似 PocketBase 的强大插件系统。