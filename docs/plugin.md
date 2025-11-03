# 插件系统设计方案

## 目录
1. [概述](#概述)
2. [插件系统架构](#插件系统架构)
3. [插件接口定义](#插件接口定义)
4. [插件管理器](#插件管理器)
5. [插件配置管理](#插件配置管理)
6. [插件数据管理](#插件数据管理)
7. [插件日志管理](#插件日志管理)
8. [插件卸载机制](#插件卸载机制)
9. [UI管理界面](#ui管理界面)
10. [安全性和权限控制](#安全性和权限控制)
11. [实施建议](#实施建议)

## 概述

本方案设计了一个完整的插件系统，旨在为urlDB提供模块化扩展能力。插件系统支持动态加载、配置管理、数据存储、日志记录、任务调度等功能，使得新功能可以通过插件形式轻松集成到系统中。

## 实现方式说明

### 自研轻量级实现的原因

**1. 避免Go plugin包的限制**
- Go标准库的plugin包要求插件必须是.so文件，且仅在Linux/macOS平台支持
- 编译环境要求严格，插件和主程序必须使用完全相同的Go版本和编译环境
- 调试复杂，热更新支持有限，生产环境稳定性较差

**2. 保持系统架构一致性**
- urlDB采用分层架构设计，自研插件系统能更好地融入现有体系
- 复用现有组件：Repository管理器、配置管理器、任务管理器等
- 统一的错误处理、日志记录、监控体系

**3. 运维友好性**
- 单一二进制文件部署，无需管理额外的.so文件
- 插件配置统一存储在数据库，便于备份和迁移
- 插件状态可观测性强，集成到现有监控体系

### 自研轻量级实现的优势

**1. 跨平台兼容性**
- 支持Windows、Linux、macOS所有平台
- 无需特殊的编译环境和依赖
- 容器化部署友好

**2. 开发体验优化**
- 插件开发者只需关注业务逻辑，无需学习复杂的plugin API
- 支持热重载，开发调试效率高
- 完整的TypeScript类型支持（前端）和Go接口定义（后端）

**3. 运维管理简化**
- 插件版本管理统一，支持回滚操作
- 插件配置变更可追踪，支持审计
- 插件异常不影响主系统稳定性

**4. 性能优化**
- 进程内调用，避免IPC开销
- 内存共享，数据传递效率高
- 可针对具体场景优化加载策略

## 插件系统架构

### 架构图
```
┌─────────────────────────────────────────────────────────────┐
│                    主应用程序                               │
├─────────────────────────────────────────────────────────────┤
│  插件管理器  │  配置管理器  │  数据管理器  │  日志管理器      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  插件1 (内置)   插件2 (内置)   插件3 (内置)   ...           │
│  [网盘服务]    [通知服务]    [统计分析]                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件
1. **插件管理器**: 负责插件的加载、卸载、生命周期管理
2. **配置管理器**: 管理插件配置项的定义、存储和验证
3. **数据管理器**: 管理插件数据的存储、查询和清理
4. **日志管理器**: 管理插件日志的记录、查询和清理

## 插件加载机制

### 进程内加载设计

**设计理念**
- 采用进程内加载模式，插件作为Go包直接编译到主程序中
- 通过注册机制实现动态启用/禁用，而非动态加载.so文件
- 保持Go语言的类型安全和性能优势

**加载流程**
```go
// 1. 插件注册阶段（编译时）
func init() {
    pluginManager.Register(&NetworkDiskPlugin{})
    pluginManager.Register(&NotificationPlugin{})
}

// 2. 插件发现阶段（启动时）
func (pm *PluginManager) DiscoverPlugins() {
    for _, plugin := range pm.registeredPlugins {
        if pm.isPluginEnabled(plugin.Name()) {
            pm.loadedPlugins[plugin.Name()] = plugin
        }
    }
}

// 3. 插件初始化阶段（运行时）
func (pm *PluginManager) InitializePlugin(name string, config map[string]interface{}) error {
    plugin, exists := pm.loadedPlugins[name]
    if !exists {
        return fmt.Errorf("插件 %s 未找到", name)
    }

    context := pm.createPluginContext(name)
    return plugin.Initialize(context)
}
```

### 插件容器机制

**容器设计**
```go
type PluginContainer struct {
    plugins map[string]*PluginInstance
    mutex   sync.RWMutex
}

type PluginInstance struct {
    Plugin     Plugin
    Context    PluginContext
    Status     PluginStatus
    Config     map[string]interface{}
    StartTime  time.Time
    LastError  error
}
```

**隔离机制**
- 每个插件运行在独立的容器实例中
- 通过Context接口限制插件直接访问系统资源
- 插件间通信通过插件管理器进行调度

### 配置驱动的加载策略

**启用/禁用机制**
```go
type PluginLoadConfig struct {
    Name         string            `json:"name"`
    Enabled      bool              `json:"enabled"`
    AutoStart    bool              `json:"auto_start"`
    LoadOrder    int               `json:"load_order"`
    Dependencies []string          `json:"dependencies"`
    Config       map[string]interface{} `json:"config"`
}

// 从数据库加载插件配置
func (pm *PluginManager) loadPluginConfigs() error {
    configs, err := pm.repo.SystemConfigRepository.GetPluginConfigs()
    if err != nil {
        return err
    }

    for _, config := range configs {
        if config.Enabled {
            pm.enablePlugin(config.Name, config.Config)
        }
    }
    return nil
}
```

### 热重载支持

**配置热重载**
```go
func (pm *PluginManager) ReloadPlugin(name string) error {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    // 1. 停止插件
    if err := pm.stopPluginUnsafe(name); err != nil {
        return err
    }

    // 2. 重新加载配置
    config, err := pm.loadPluginConfig(name)
    if err != nil {
        return err
    }

    // 3. 重新初始化插件
    return pm.initializePluginUnsafe(name, config)
}
```

**状态保持机制**
- 热重载时保持插件数据的连续性
- 支持增量配置更新，避免全量重新初始化
- 提供回滚机制，重载失败时恢复原状态

## 插件注册和发现机制

### 自动注册机制

**编译时注册**
```go
// 插件包中的注册代码
package plugins

import (
    "github.com/ctwj/urldb/plugin"
)

func init() {
    plugin.GlobalManager.Register(&NetworkDiskPlugin{
        meta: PluginMeta{
            Name:        "network-disk",
            Version:     "1.0.0",
            Description: "网盘服务插件",
            Author:      "urlDB Team",
            Category:    "storage",
        },
    })
}
```

**插件元数据管理**
```go
type PluginMeta struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Author      string            `json:"author"`
    Category    string            `json:"category"`
    Tags        []string          `json:"tags"`
    Homepage    string            `json:"homepage"`
    License     string            `json:"license"`
    MinVersion  string            `json:"min_version"`    // 最低主程序版本
    Dependencies []string         `json:"dependencies"`  // 依赖的其他插件
    ConfigSchema *ConfigSchema    `json:"config_schema"` // 配置模式定义
}
```

### 插件发现策略

**静态发现**
```go
type PluginRegistry struct {
    plugins     map[string]Plugin
    metadata    map[string]*PluginMeta
    loadOrder   []string
    mutex       sync.RWMutex
}

// 扫描所有已注册的插件
func (pr *PluginRegistry) DiscoverPlugins() error {
    pr.mutex.Lock()
    defer pr.mutex.Unlock()

    // 按依赖关系排序
    sorted, err := pr.topologicalSort()
    if err != nil {
        return fmt.Errorf("插件依赖解析失败: %v", err)
    }

    pr.loadOrder = sorted
    utils.Info("发现 %d 个插件，加载顺序: %v", len(pr.plugins), sorted)
    return nil
}
```

**动态发现**
```go
// 支持运行时发现新的插件类型
func (pr *PluginRegistry) RegisterLate(plugin Plugin) error {
    pr.mutex.Lock()
    defer pr.mutex.Unlock()

    meta := plugin.GetMeta()
    if _, exists := pr.plugins[meta.Name]; exists {
        return fmt.Errorf("插件 %s 已存在", meta.Name)
    }

    pr.plugins[meta.Name] = plugin
    pr.metadata[meta.Name] = meta

    // 重新计算加载顺序
    _, err := pr.topologicalSort()
    return err
}
```

### 依赖解析机制

**依赖图构建**
```go
type DependencyGraph struct {
    nodes map[string]*PluginNode
    edges map[string][]string
}

type PluginNode struct {
    Name         string
    Plugin       Plugin
    Dependencies []string
    Dependents   []string
    Level        int // 依赖层级
}

func (pr *PluginRegistry) buildDependencyGraph() *DependencyGraph {
    graph := &DependencyGraph{
        nodes: make(map[string]*PluginNode),
        edges: make(map[string][]string),
    }

    // 构建节点
    for name, plugin := range pr.plugins {
        meta := plugin.GetMeta()
        graph.nodes[name] = &PluginNode{
            Name:         name,
            Plugin:       plugin,
            Dependencies: meta.Dependencies,
        }
    }

    // 构建边
    for name, node := range graph.nodes {
        for _, dep := range node.Dependencies {
            if _, exists := graph.nodes[dep]; !exists {
                utils.Warn("插件 %s 依赖的插件 %s 不存在", name, dep)
                continue
            }
            graph.edges[name] = append(graph.edges[name], dep)
            graph.nodes[dep].Dependents = append(graph.nodes[dep].Dependents, name)
        }
    }

    return graph
}
```

**拓扑排序**
```go
func (pr *PluginRegistry) topologicalSort() ([]string, error) {
    graph := pr.buildDependencyGraph()
    var result []string
    visited := make(map[string]bool)
    visiting := make(map[string]bool)

    var visit func(name string) error
    visit = func(name string) error {
        if visiting[name] {
            return fmt.Errorf("检测到循环依赖，涉及插件: %s", name)
        }
        if visited[name] {
            return nil
        }

        visiting[name] = true
        for _, dep := range graph.edges[name] {
            if err := visit(dep); err != nil {
                return err
            }
        }
        visiting[name] = false
        visited[name] = true

        result = append(result, name)
        return nil
    }

    for name := range graph.nodes {
        if !visited[name] {
            if err := visit(name); err != nil {
                return nil, err
            }
        }
    }

    return result, nil
}
```

### 插件验证机制

**兼容性检查**
```go
func (pr *PluginRegistry) validatePlugin(plugin Plugin) error {
    meta := plugin.GetMeta()

    // 1. 检查必需字段
    if meta.Name == "" || meta.Version == "" {
        return fmt.Errorf("插件名称和版本不能为空")
    }

    // 2. 检查版本兼容性
    if err := pr.checkVersionCompatibility(meta.MinVersion); err != nil {
        return fmt.Errorf("版本兼容性检查失败: %v", err)
    }

    // 3. 检查接口实现
    if err := pr.validateInterface(plugin); err != nil {
        return fmt.Errorf("接口验证失败: %v", err)
    }

    // 4. 检查配置模式
    if meta.ConfigSchema != nil {
        if err := pr.validateConfigSchema(meta.ConfigSchema); err != nil {
            return fmt.Errorf("配置模式验证失败: %v", err)
        }
    }

    return nil
}
```

**健康检查**
```go
func (pr *PluginRegistry) HealthCheck() map[string]error {
    pr.mutex.RLock()
    defer pr.mutex.RUnlock()

    results := make(map[string]error)

    for name, plugin := range pr.plugins {
        if healthChecker, ok := plugin.(HealthChecker); ok {
            if err := healthChecker.Check(); err != nil {
                results[name] = err
            }
        } else {
            // 默认健康检查：检查插件是否响应
            if plugin.GetMeta() == nil {
                results[name] = fmt.Errorf("插件元数据获取失败")
            }
        }
    }

    return results
}
```

## 插件接口定义

```go
// 插件基础接口
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

    // 配置管理
    ConfigSchema() *PluginConfigSchema
    ValidateConfig(config map[string]interface{}) error
}

// 插件上下文接口
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
    GetDB() *gorm.DB
}
```

## 插件管理器

### 功能职责
1. 插件注册和发现
2. 插件加载和卸载
3. 插件生命周期管理
4. 插件依赖解析
5. 插件状态监控

### 核心接口
```go
type PluginManager interface {
    // 插件注册
    RegisterPlugin(plugin Plugin) error
    UnregisterPlugin(name string) error

    // 插件发现
    DiscoverPlugins(path string) error
    LoadPlugin(name string) error
    UnloadPlugin(name string) error

    // 插件查询
    GetPlugin(name string) (Plugin, error)
    ListPlugins() []PluginInfo
    GetPluginStatus(name string) PluginStatus

    // 依赖管理
    ResolveDependencies(plugin Plugin) error
    CheckAllDependencies() map[string]bool

    // 生命周期管理
    InitializePlugin(name string, config map[string]interface{}) error
    StartPlugin(name string) error
    StopPlugin(name string) error
    RestartPlugin(name string) error
}
```

## 插件生命周期管理

### 生命周期状态

**状态定义**
```go
type PluginStatus string

const (
    StatusRegistered   PluginStatus = "registered"   // 已注册
    StatusInitialized  PluginStatus = "initialized"  // 已初始化
    StatusStarting     PluginStatus = "starting"     // 启动中
    StatusRunning      PluginStatus = "running"      // 运行中
    StatusStopping     PluginStatus = "stopping"     // 停止中
    StatusStopped      PluginStatus = "stopped"      // 已停止
    StatusError        PluginStatus = "error"        // 错误状态
    StatusDisabled     PluginStatus = "disabled"     // 已禁用
)

type PluginLifecycle struct {
    Status       PluginStatus    `json:"status"`
    StartTime    time.Time       `json:"start_time"`
    StopTime     *time.Time      `json:"stop_time,omitempty"`
    LastError    error           `json:"last_error,omitempty"`
    RestartCount int             `json:"restart_count"`
    HealthScore  float64         `json:"health_score"` // 0-100 健康评分
}
```

### 生命周期流程控制

**启动流程**
```go
func (pm *PluginManager) StartPlugin(name string) error {
    instance, exists := pm.instances[name]
    if !exists {
        return fmt.Errorf("插件 %s 不存在", name)
    }

    // 1. 状态检查
    if instance.Status != StatusInitialized && instance.Status != StatusStopped {
        return fmt.Errorf("插件状态不正确，当前状态: %s", instance.Status)
    }

    // 2. 更新状态为启动中
    pm.updateStatus(name, StatusStarting)

    // 3. 依赖检查
    if err := pm.checkDependencies(name); err != nil {
        pm.updateStatus(name, StatusError)
        return fmt.Errorf("依赖检查失败: %v", err)
    }

    // 4. 调用插件启动方法
    if err := instance.Plugin.Start(); err != nil {
        pm.updateStatus(name, StatusError)
        instance.LastError = err
        return fmt.Errorf("插件启动失败: %v", err)
    }

    // 5. 更新状态和统计信息
    pm.updateStatus(name, StatusRunning)
    instance.StartTime = time.Now()
    instance.RestartCount++

    // 6. 启动健康检查
    go pm.startHealthCheck(name)

    utils.Info("插件 %s 启动成功", name)
    return nil
}
```

**停止流程**
```go
func (pm *PluginManager) StopPlugin(name string, graceful bool) error {
    instance, exists := pm.instances[name]
    if !exists {
        return fmt.Errorf("插件 %s 不存在", name)
    }

    if instance.Status != StatusRunning {
        return fmt.Errorf("插件未在运行，当前状态: %s", instance.Status)
    }

    // 1. 更新状态为停止中
    pm.updateStatus(name, StatusStopping)

    // 2. 停止健康检查
    pm.stopHealthCheck(name)

    // 3. 优雅停止或强制停止
    var err error
    if graceful {
        // 优雅停止，等待正在执行的任务完成
        err = pm.gracefulStopPlugin(instance, 30*time.Second)
    } else {
        // 强制停止
        err = instance.Plugin.Stop()
    }

    // 4. 更新状态
    if err != nil {
        pm.updateStatus(name, StatusError)
        instance.LastError = err
        return fmt.Errorf("插件停止失败: %v", err)
    }

    now := time.Now()
    instance.StopTime = &now
    pm.updateStatus(name, StatusStopped)

    utils.Info("插件 %s 停止成功", name)
    return nil
}
```

### 健康检查机制

**健康检查策略**
```go
type HealthChecker interface {
    Check() error
    GetMetrics() map[string]interface{}
}

func (pm *PluginManager) startHealthCheck(name string) {
    instance := pm.instances[name]

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if instance.Status != StatusRunning {
                return
            }

            if err := pm.performHealthCheck(name); err != nil {
                utils.Error("插件 %s 健康检查失败: %v", name, err)
                pm.handleHealthCheckFailure(name, err)
            } else {
                instance.HealthScore = math.Min(100, instance.HealthScore+5)
            }

        case <-pm.ctx.Done():
            return
        }
    }
}

func (pm *PluginManager) performHealthCheck(name string) error {
    instance := pm.instances[name]

    if healthChecker, ok := instance.Plugin.(HealthChecker); ok {
        return healthChecker.Check()
    }

    // 默认健康检查：检查插件是否响应ping
    if pinger, ok := instance.Plugin.(Pinger); ok {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        return pinger.Ping(ctx)
    }

    return nil
}
```

### 故障恢复机制

**自动重启策略**
```go
type RestartPolicy struct {
    MaxRestarts    int           `json:"max_restarts"`     // 最大重启次数
    RestartDelay   time.Duration `json:"restart_delay"`    // 重启延迟
    BackoffFactor  float64       `json:"backoff_factor"`   // 退避因子
    MaxBackoff     time.Duration `json:"max_backoff"`      // 最大退避时间
}

func (pm *PluginManager) handleHealthCheckFailure(name string, err error) {
    instance := pm.instances[name]
    policy := pm.getRestartPolicy(name)

    // 1. 更新健康评分
    instance.HealthScore = math.Max(0, instance.HealthScore-10)
    instance.LastError = err

    // 2. 检查是否需要重启
    if instance.RestartCount >= policy.MaxRestarts {
        pm.updateStatus(name, StatusError)
        utils.Error("插件 %s 达到最大重启次数，停止重启", name)
        return
    }

    // 3. 计算延迟时间
    delay := time.Duration(float64(policy.RestartDelay) *
        math.Pow(policy.BackoffFactor, float64(instance.RestartCount)))
    delay = min(delay, policy.MaxBackoff)

    utils.Warn("插件 %s 健康检查失败，%v 后将重启", name, delay)

    // 4. 延迟重启
    time.AfterFunc(delay, func() {
        if pm.instances[name].Status == StatusRunning {
            pm.restartPlugin(name)
        }
    })
}
```

### 生命周期事件

**事件系统**
```go
type LifecycleEvent struct {
    PluginName string        `json:"plugin_name"`
    EventType  string        `json:"event_type"`
    Timestamp  time.Time     `json:"timestamp"`
    Data       interface{}   `json:"data,omitempty"`
    Error      error         `json:"error,omitempty"`
}

type EventHandler func(event LifecycleEvent)

func (pm *PluginManager) emitEvent(eventType, pluginName string, data interface{}, err error) {
    event := LifecycleEvent{
        PluginName: pluginName,
        EventType:  eventType,
        Timestamp:  time.Now(),
        Data:       data,
        Error:      err,
    }

    for _, handler := range pm.eventHandlers {
        go handler(event) // 异步处理事件
    }
}

// 注册事件处理器
func (pm *PluginManager) OnPluginStarted(handler EventHandler) {
    pm.eventHandlers = append(pm.eventHandlers, func(event LifecycleEvent) {
        if event.EventType == "plugin.started" {
            handler(event)
        }
    })
}
```

## 插件配置管理

### 配置定义格式
```go
type PluginConfigSchema struct {
    Name        string                 `json:"name"`
    Version     string                 `json:"version"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Properties  map[string]ConfigField `json:"properties"`
    Required    []string               `json:"required"`
    Category    string                 `json:"category"`
    Order       int                    `json:"order"`
}

type ConfigField struct {
    Type        string      `json:"type"`
    Title       string      `json:"title"`
    Description string      `json:"description"`
    Default     interface{} `json:"default"`
    Required    bool        `json:"required"`
    Min         *float64    `json:"min"`
    Max         *float64    `json:"max"`
    Pattern     string      `json:"pattern"`
    Enum        []string    `json:"enum"`
    Format      string      `json:"format"`
    IsSecret    bool        `json:"is_secret"`
    Group       string      `json:"group"`
    Order       int         `json:"order"`
    Placeholder string      `json:"placeholder"`
    Hint        string      `json:"hint"`
}
```

### 配置存储结构
```go
type SystemConfig struct {
    ID          uint      `gorm:"primaryKey"`
    ConfigKey   string    `gorm:"size:255;not null;index:idx_key"`
    ConfigValue string    `gorm:"type:text"`
    Description string    `gorm:"type:text"`
    DataType    string    `gorm:"size:50"`
    IsEncrypted bool      `gorm:"default:false"`
    IsPlugin    bool      `gorm:"default:false"`
    PluginName  string    `gorm:"size:100;index:idx_plugin"`
    Category    string    `gorm:"size:100;index:idx_category"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## 插件数据管理

### 数据存储结构
```go
type PluginData struct {
    ID          uint      `gorm:"primaryKey"`
    PluginName  string    `gorm:"size:100;index:idx_plugin"`
    DataType    string    `gorm:"size:100;index:idx_type"`
    DataKey     string    `gorm:"size:255;index:idx_data_key"`
    DataValue   string    `gorm:"type:text"`
    Metadata    string    `gorm:"type:json"`
    CreatedAt   time.Time `gorm:"index:idx_created"`
    UpdatedAt   time.Time
    ExpiresAt   *time.Time `gorm:"index:idx_expires"`
}
```

### 数据查询接口
```go
type DataQueryParams struct {
    PluginName  string    `json:"plugin_name"`
    DataType    string    `json:"data_type,omitempty"`
    DataKey     string    `json:"data_key,omitempty"`
    StartTime   time.Time `json:"start_time,omitempty"`
    EndTime     time.Time `json:"end_time,omitempty"`
    Page        int       `json:"page,omitempty"`
    PageSize    int       `json:"page_size,omitempty"`
    Search      string    `json:"search,omitempty"`
    SortBy      string    `json:"sort_by,omitempty"`
    SortOrder   string    `json:"sort_order,omitempty"`
}
```

## 插件日志管理

### 日志存储结构
```go
type PluginLog struct {
    ID         uint      `gorm:"primaryKey"`
    PluginName string    `gorm:"size:100;index:idx_plugin"`
    Level      string    `gorm:"size:20;index:idx_level"`
    Message    string    `gorm:"type:text"`
    Context    string    `gorm:"type:json"`
    TraceID    string    `gorm:"size:100;index:idx_trace"`
    CreatedAt  time.Time `gorm:"index:idx_created"`
}
```

### 日志查询接口
```go
type LogQueryParams struct {
    PluginName  string    `json:"plugin_name"`
    Level       string    `json:"level,omitempty"`
    StartTime   time.Time `json:"start_time,omitempty"`
    EndTime     time.Time `json:"end_time,omitempty"`
    Page        int       `json:"page,omitempty"`
    PageSize    int       `json:"page_size,omitempty"`
    Search      string    `json:"search,omitempty"`
    TraceID     string    `json:"trace_id,omitempty"`
}
```

## 与现有系统集成方式

### 集成架构设计

**统一入口模式**
urlDB插件系统采用统一入口模式与现有系统集成，通过插件管理器作为中央协调器：

```go
// 主程序启动时初始化插件系统
func initPluginSystem() {
    // 1. 创建插件管理器
    pluginManager := plugin.NewPluginManager(
        repoManager,      // 复用现有Repository管理器
        configManager,    // 复用现有配置管理器
        taskManager,      // 复用现有任务管理器
    )

    // 2. 插件发现和加载
    pluginManager.DiscoverPlugins()

    // 3. 注册到全局配置管理器
    config.SetPluginManager(pluginManager)

    // 4. 注册到全局任务管理器
    task.SetPluginTaskManager(pluginManager)

    // 5. 集成到现有HTTP路由系统
    setupPluginRoutes(router, pluginManager)
}
```

**依赖注入机制**
```go
type PluginContext struct {
    // 复用现有Repository
    RepoManager *repo.RepositoryManager

    // 复用现有配置管理器
    ConfigManager *config.ConfigManager

    // 复用现有任务管理器
    TaskManager *task.TaskManager

    // 复用现有日志系统
    Logger *utils.Logger

    // 复用现有数据库连接
    DB *gorm.DB

    // 复用现有监控系统
    Metrics *monitor.Metrics

    // 插件专用组件
    PluginDataStore  *PluginDataStore
    PluginConfigStore *PluginConfigStore
    PluginLoggerStore *PluginLoggerStore
}
```

### 与HTTP路由系统集成

**路由注册机制**
```go
// 插件路由自动注册
func setupPluginRoutes(router *gin.Engine, pm *plugin.Manager) {
    // 获取所有已启用的插件
    plugins := pm.GetEnabledPlugins()

    for _, plugin := range plugins {
        // 检查插件是否实现了HTTP路由接口
        if httpHandler, ok := plugin.(HTTPHandler); ok {
            // 为插件注册路由前缀
            pluginGroup := router.Group(fmt.Sprintf("/api/plugins/%s", plugin.Name()))

            // 插件注册自己的路由
            httpHandler.RegisterRoutes(pluginGroup)
        }
    }
}

// 插件路由接口定义
type HTTPHandler interface {
    RegisterRoutes(group *gin.RouterGroup)
}

// 插件路由实现示例
func (p *NetworkDiskPlugin) RegisterRoutes(group *gin.RouterGroup) {
    group.GET("/files", p.handleListFiles)
    group.POST("/files", p.handleUploadFile)
    group.GET("/files/:id", p.handleGetFile)
    group.DELETE("/files/:id", p.handleDeleteFile)
}
```

### 与任务调度系统集成

**任务处理器注册**
```go
// 插件任务处理器自动注册
func registerPluginTasks(pm *plugin.Manager, taskManager *task.TaskManager) {
    plugins := pm.GetEnabledPlugins()

    for _, plugin := range plugins {
        // 检查插件是否实现了任务处理器接口
        if taskHandler, ok := plugin.(TaskHandler); ok {
            // 注册插件任务处理器
            taskManager.RegisterProcessor(&PluginTaskProcessor{
                plugin: plugin,
                handler: taskHandler,
            })
        }
    }
}

// 插件任务处理器实现
type PluginTaskProcessor struct {
    plugin  Plugin
    handler TaskHandler
}

func (p *PluginTaskProcessor) Process(ctx context.Context, taskID uint, item *entity.TaskItem) error {
    // 调用插件的任务处理方法
    return p.handler.ProcessTask(ctx, taskID, item)
}

func (p *PluginTaskProcessor) GetTaskType() string {
    // 返回插件的任务类型标识
    return fmt.Sprintf("plugin.%s.task", p.plugin.Name())
}
```

### 与配置管理系统集成

**插件配置统一管理**
```go
// 插件配置与系统配置统一存储
type SystemConfig struct {
    ID          uint      `gorm:"primaryKey"`
    ConfigKey   string    `gorm:"size:255;not null;index:idx_key"`
    ConfigValue string    `gorm:"type:text"`
    Description string    `gorm:"type:text"`
    DataType    string    `gorm:"size:50"`
    IsEncrypted bool      `gorm:"default:false"`
    IsPlugin    bool      `gorm:"default:false"`     // 标识是否为插件配置
    PluginName  string    `gorm:"size:100;index:idx_plugin"` // 插件名称
    Category    string    `gorm:"size:100;index:idx_category"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 配置管理器扩展
func (cm *ConfigManager) GetPluginConfig(pluginName, key string) (*ConfigItem, error) {
    fullKey := fmt.Sprintf("plugin.%s.%s", pluginName, key)
    return cm.GetConfig(fullKey)
}

func (cm *ConfigManager) SetPluginConfig(pluginName, key, value string) error {
    fullKey := fmt.Sprintf("plugin.%s.%s", pluginName, key)
    return cm.SetConfig(fullKey, value)
}
```

### 与数据库系统集成

**数据存储统一管理**
```go
// 插件数据统一存储到系统表中
type PluginData struct {
    ID          uint      `gorm:"primaryKey"`
    PluginName  string    `gorm:"size:100;index:idx_plugin"`
    DataType    string    `gorm:"size:100;index:idx_type"`
    DataKey     string    `gorm:"size:255;index:idx_data_key"`
    DataValue   string    `gorm:"type:text"`
    Metadata    string    `gorm:"type:json"`
    CreatedAt   time.Time `gorm:"index:idx_created"`
    UpdatedAt   time.Time
    ExpiresAt   *time.Time `gorm:"index:idx_expires"`
}

// 插件数据访问接口
type PluginDataStore interface {
    Get(pluginName, dataType, key string) (*PluginData, error)
    Set(pluginName, dataType, key, value string, metadata map[string]interface{}) error
    Delete(pluginName, dataType, key string) error
    List(pluginName, dataType string, params *DataQueryParams) ([]*PluginData, error)
}

// 数据访问权限控制
func (pds *pluginDataStore) Get(pluginName, dataType, key string) (*PluginData, error) {
    // 确保插件只能访问自己的数据
    if !pds.context.IsPluginOwner(pluginName) {
        return nil, fmt.Errorf("权限不足，无法访问插件 %s 的数据", pluginName)
    }

    return pds.repo.PluginDataRepository.GetByPluginAndKey(pluginName, dataType, key)
}
```

### 与监控告警系统集成

**统一监控指标**
```go
// 插件监控指标集成
type PluginMetrics struct {
    // 复用现有监控系统
    metrics *monitor.Metrics

    // 插件专用指标
    pluginCallTotal     *prometheus.CounterVec
    pluginCallDuration  *prometheus.HistogramVec
    pluginErrorTotal    *prometheus.CounterVec
    pluginHealthScore   *prometheus.GaugeVec
}

func NewPluginMetrics(metrics *monitor.Metrics) *PluginMetrics {
    pm := &PluginMetrics{
        metrics: metrics,
        pluginCallTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: "urldb",
                Subsystem: "plugin",
                Name:      "call_total",
                Help:      "插件调用总次数",
            },
            []string{"plugin_name", "method"},
        ),
    }

    // 注册到现有监控系统
    metrics.RegisterCollector(pm.pluginCallTotal)
    return pm
}

// 插件中使用监控
func (p *NetworkDiskPlugin) ListFiles(ctx context.Context) ([]File, error) {
    timer := p.metrics.pluginCallDuration.WithLabelValues(p.Name(), "ListFiles").StartTimer()
    defer timer.ObserveDuration()
    defer p.metrics.pluginCallTotal.WithLabelValues(p.Name(), "ListFiles").Inc()

    files, err := p.listFiles(ctx)
    if err != nil {
        p.metrics.pluginErrorTotal.WithLabelValues(p.Name(), "ListFiles").Inc()
    }

    return files, err
}
```

### 与日志系统集成

**统一日志管理**
```go
// 插件日志复用现有日志系统
func (p *NetworkDiskPlugin) initLogger() {
    // 使用系统日志前缀
    p.logger = utils.NewPluginLogger(p.Name())
}

// 插件日志记录示例
func (p *NetworkDiskPlugin) ListFiles(ctx context.Context) ([]File, error) {
    p.logger.Debug("开始获取文件列表")

    files, err := p.listFiles(ctx)
    if err != nil {
        p.logger.Error("获取文件列表失败: %v", err)
        return nil, err
    }

    p.logger.Info("成功获取 %d 个文件", len(files))
    return files, nil
}
```

## 插件卸载机制

### 卸载策略
1. **完全卸载**: 删除所有相关数据和配置
2. **保留数据**: 保留数据备份，仅清理配置
3. **停用模式**: 仅停止插件，保留所有数据
4. **归档模式**: 将数据归档到专用表

### 卸载配置
```go
type UninstallConfig struct {
    RemoveTables     bool `json:"remove_tables"`     // 删除表结构
    BackupData       bool `json:"backup_data"`       // 备份数据
    RemoveConfig     bool `json:"remove_config"`     // 清理配置
    ArchiveTasks     bool `json:"archive_tasks"`     // 归档任务
}
```

### 回滚机制
```go
type UninstallRollback struct {
    PluginName    string    `json:"plugin_name"`
    BackupPath    string    `json:"backup_path"`
    RollbackSteps []string  `json:"rollback_steps"`
    CreatedAt     time.Time `json:"created_at"`
}
```

## UI管理界面

### 插件列表页面
- 插件基本信息展示
- 插件状态管理（启用/禁用）
- 插件配置入口
- 插件数据和日志查看入口

### 插件配置页面
- 基于JSON Schema的动态表单
- 配置项分组显示
- 实时配置验证
- 配置历史记录

### 插件数据页面
- 数据概览统计
- 数据查询和过滤
- 数据详情查看
- 数据导出功能

### 插件日志页面
- 日志统计图表
- 日志查询和过滤
- 日志详情查看
- 日志导出和清理

## 安全性和权限控制

### 权限管理
1. **插件加载权限**: 只有管理员可以加载新插件
2. **插件配置权限**: 根据角色控制配置项访问
3. **数据访问权限**: 插件只能访问自己的数据
4. **日志查看权限**: 根据角色控制日志访问级别

### 安全措施
1. **插件签名验证**: 验证插件来源和完整性
2. **沙箱隔离**: 限制插件系统资源访问
3. **数据加密**: 敏感配置项加密存储
4. **审计日志**: 记录所有插件管理操作

## 性能优化考虑

### 内存优化策略

**插件实例池化**
```go
// 插件实例复用机制
type PluginInstancePool struct {
    instances map[string]*sync.Pool  // 按插件类型分类的实例池
    mutex     sync.RWMutex
}

func (pip *PluginInstancePool) GetInstance(pluginType string) Plugin {
    pip.mutex.RLock()
    pool, exists := pip.instances[pluginType]
    pip.mutex.RUnlock()

    if !exists {
        return nil
    }

    if instance := pool.Get(); instance != nil {
        return instance.(Plugin)
    }

    return nil
}

func (pip *PluginInstancePool) PutInstance(pluginType string, instance Plugin) {
    pip.mutex.RLock()
    pool, exists := pip.instances[pluginType]
    pip.mutex.RUnlock()

    if !exists {
        return
    }

    // 重置实例状态
    if resetter, ok := instance.(Resetter); ok {
        resetter.Reset()
    }

    pool.Put(instance)
}
```

**内存泄漏防护**
```go
// 插件内存监控
type MemoryMonitor struct {
    pluginMemory map[string]uint64  // 插件内存使用量
    threshold    uint64             // 内存使用阈值
    alerts       chan MemoryAlert   // 内存告警通道
}

func (mm *MemoryMonitor) MonitorPlugin(pluginName string) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            usage := mm.getPluginMemoryUsage(pluginName)
            mm.pluginMemory[pluginName] = usage

            // 检查内存使用是否超过阈值
            if usage > mm.threshold {
                alert := MemoryAlert{
                    PluginName: pluginName,
                    Usage:      usage,
                    Threshold:  mm.threshold,
                    Timestamp:  time.Now(),
                }
                mm.alerts <- alert
            }

        case <-mm.ctx.Done():
            return
        }
    }
}
```

### 并发性能优化

**协程池管理**
```go
// 插件协程池
type PluginGoroutinePool struct {
    workers    chan chan Work  // 工作协程通道
    maxWorkers int            // 最大协程数
    jobQueue   chan Work      // 任务队列
    quit       chan bool      // 退出信号
}

// 工作单元定义
type Work struct {
    PluginName string
    Function   func() error
    Callback   func(error)
}

// 提交任务到协程池
func (pgp *PluginGoroutinePool) Submit(pluginName string, fn func() error, callback func(error)) {
    work := Work{
        PluginName: pluginName,
        Function:   fn,
        Callback:   callback,
    }

    select {
    case pgp.jobQueue <- work:
        // 任务提交成功
    default:
        // 任务队列满，执行降级策略
        go func() {
            err := fn()
            if callback != nil {
                callback(err)
            }
        }()
    }
}
```

**读写锁优化**
```go
// 分段锁机制
type SegmentedLock struct {
    segments []sync.RWMutex
    segmentCount int
}

func NewSegmentedLock(segmentCount int) *SegmentedLock {
    return &SegmentedLock{
        segments: make([]sync.RWMutex, segmentCount),
        segmentCount: segmentCount,
    }
}

func (sl *SegmentedLock) getSegment(key string) *sync.RWMutex {
    hash := fnv.New32a()
    hash.Write([]byte(key))
    index := hash.Sum32() % uint32(sl.segmentCount)
    return &sl.segments[index]
}

func (sl *SegmentedLock) Lock(key string) {
    sl.getSegment(key).Lock()
}

func (sl *SegmentedLock) RLock(key string) {
    sl.getSegment(key).RLock()
}
```

### 缓存优化策略

**多级缓存架构**
```go
// 插件缓存管理器
type PluginCacheManager struct {
    // L1: 本地内存缓存
    l1Cache *LocalCache

    // L2: 共享缓存（Redis等）
    l2Cache Cache

    // L3: 数据库
    database *gorm.DB
}

type CacheStrategy struct {
    PluginName     string        `json:"plugin_name"`
    CacheTTL       time.Duration `json:"cache_ttl"`
    CacheSize      int           `json:"cache_size"`
    EvictionPolicy string        `json:"eviction_policy"` // LRU, LFU, FIFO
    UseL2Cache     bool          `json:"use_l2_cache"`
}

// 智能缓存失效
func (pcm *PluginCacheManager) InvalidateCache(pluginName string, pattern string) {
    // 1. 失效本地缓存
    pcm.l1Cache.InvalidateByPattern(pattern)

    // 2. 失效共享缓存
    if pcm.l2Cache != nil {
        pcm.l2Cache.InvalidateByPattern(pattern)
    }

    // 3. 记录缓存失效事件
    utils.Info("插件 %s 缓存失效: %s", pluginName, pattern)
}
```

**缓存预热机制**
```go
// 缓存预热策略
type CacheWarmupStrategy struct {
    PluginName   string   `json:"plugin_name"`
    WarmupKeys   []string `json:"warmup_keys"`
    WarmupOnInit bool     `json:"warmup_on_init"`
    WarmupCron   string   `json:"warmup_cron"` // 定时预热表达式
}

func (pcm *PluginCacheManager) WarmupPluginCache(pluginName string) error {
    strategy := pcm.getWarmupStrategy(pluginName)
    if strategy == nil {
        return nil
    }

    var wg sync.WaitGroup
    errChan := make(chan error, len(strategy.WarmupKeys))

    for _, key := range strategy.WarmupKeys {
        wg.Add(1)
        go func(k string) {
            defer wg.Done()
            if err := pcm.preloadCacheKey(pluginName, k); err != nil {
                errChan <- fmt.Errorf("预热缓存键 %s 失败: %v", k, err)
            }
        }(key)
    }

    wg.Wait()
    close(errChan)

    // 收集错误
    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }

    if len(errors) > 0 {
        return fmt.Errorf("缓存预热完成，但有 %d 个错误: %v", len(errors), errors)
    }

    utils.Info("插件 %s 缓存预热完成", pluginName)
    return nil
}
```

### 数据库性能优化

**批量操作优化**
```go
// 插件数据批量操作
type PluginDataBatch struct {
    operations []DataOperation
    batchSize  int
    flushTimer *time.Timer
}

type DataOperation struct {
    Type      string      // INSERT, UPDATE, DELETE
    PluginName string
    DataType  string
    Key       string
    Value     string
    Metadata  map[string]interface{}
}

// 批量提交数据操作
func (pdb *PluginDataBatch) BatchCommit(operations []DataOperation) error {
    if len(operations) == 0 {
        return nil
    }

    // 按插件和数据类型分组
    groupedOps := make(map[string][]DataOperation)
    for _, op := range operations {
        key := fmt.Sprintf("%s:%s", op.PluginName, op.DataType)
        groupedOps[key] = append(groupedOps[key], op)
    }

    // 并发处理各组操作
    var wg sync.WaitGroup
    errChan := make(chan error, len(groupedOps))

    for _, ops := range groupedOps {
        wg.Add(1)
        go func(operations []DataOperation) {
            defer wg.Done()
            if err := pdb.processBatch(operations); err != nil {
                errChan <- err
            }
        }(ops)
    }

    wg.Wait()
    close(errChan)

    // 收集错误
    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }

    if len(errors) > 0 {
        return fmt.Errorf("批量操作完成，但有 %d 个错误: %v", len(errors), errors)
    }

    return nil
}
```

**查询优化**
```go
// 插件数据查询优化器
type QueryOptimizer struct {
    queryStats map[string]*QueryStats  // 查询统计
    slowQueryThreshold time.Duration    // 慢查询阈值
}

type QueryStats struct {
    TotalCount    int64         `json:"total_count"`
    TotalDuration time.Duration `json:"total_duration"`
    AvgDuration   time.Duration `json:"avg_duration"`
    MaxDuration   time.Duration `json:"max_duration"`
    LastQueryTime time.Time     `json:"last_query_time"`
}

// 查询执行拦截器
func (qo *QueryOptimizer) InterceptQuery(pluginName, query string, fn func() (interface{}, error)) (interface{}, error) {
    startTime := time.Now()

    // 执行查询
    result, err := fn()

    duration := time.Since(startTime)
    queryKey := fmt.Sprintf("%s:%s", pluginName, utils.GetMD5Hash(query))

    // 更新统计信息
    stats := qo.getOrCreateStats(queryKey)
    stats.TotalCount++
    stats.TotalDuration += duration
    stats.AvgDuration = stats.TotalDuration / time.Duration(stats.TotalCount)
    if duration > stats.MaxDuration {
        stats.MaxDuration = duration
    }
    stats.LastQueryTime = time.Now()

    // 检查慢查询
    if duration > qo.slowQueryThreshold {
        utils.Warn("插件 %s 慢查询: %s, 耗时: %v", pluginName, query, duration)
        // 可以触发告警或自动优化
    }

    return result, err
}
```

### 网络性能优化

**连接池管理**
```go
// 插件网络连接池
type PluginConnectionPool struct {
    pools map[string]*ConnectionPool  // 按插件和目标地址分类
    mutex sync.RWMutex
}

type ConnectionPool struct {
    connections chan net.Conn
    factory     func() (net.Conn, error)
    maxSize     int
    currentSize int
    mutex       sync.Mutex
}

// 获取连接
func (pcp *PluginConnectionPool) GetConnection(pluginName, target string) (net.Conn, error) {
    poolKey := fmt.Sprintf("%s:%s", pluginName, target)

    pcp.mutex.RLock()
    pool, exists := pcp.pools[poolKey]
    pcp.mutex.RUnlock()

    if !exists {
        return nil, fmt.Errorf("连接池不存在: %s", poolKey)
    }

    select {
    case conn := <-pool.connections:
        // 检查连接是否有效
        if pcp.isConnectionValid(conn) {
            return conn, nil
        }
        // 连接无效，创建新连接
        fallthrough
    default:
        // 连接池为空，创建新连接
        return pool.factory()
    }
}

// 归还连接
func (pcp *PluginConnectionPool) PutConnection(pluginName, target string, conn net.Conn) {
    poolKey := fmt.Sprintf("%s:%s", pluginName, target)

    pcp.mutex.RLock()
    pool, exists := pcp.pools[poolKey]
    pcp.mutex.RUnlock()

    if !exists {
        conn.Close()
        return
    }

    select {
    case pool.connections <- conn:
        // 连接归还成功
    default:
        // 连接池已满，关闭连接
        conn.Close()
    }
}
```

## 实施建议

### 分阶段实施
1. **第一阶段**: 实现基础插件框架和管理器
2. **第二阶段**: 实现配置管理和UI界面
3. **第三阶段**: 实现数据和日志管理
4. **第四阶段**: 实现卸载机制和安全控制

### 最佳实践
1. **插件开发规范**: 制定插件开发标准和规范
2. **测试机制**: 建立插件测试和验证机制
3. **监控告警**: 实现插件运行状态监控
4. **文档完善**: 提供完整的插件开发和使用文档

## 部署和运维说明

### 部署架构

**容器化部署**
```dockerfile
# Dockerfile示例
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o urldb .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/urldb .
# 插件配置文件
COPY --from=builder /app/config/plugin-config.json ./config/
# 插件数据目录
RUN mkdir -p ./plugins/data

EXPOSE 8080
CMD ["./urldb"]
```

**Kubernetes部署**
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: urldb-plugin-system
spec:
  replicas: 3
  selector:
    matchLabels:
      app: urldb
  template:
    metadata:
      labels:
        app: urldb
    spec:
      containers:
      - name: urldb
        image: urldb:latest
        ports:
        - containerPort: 8080
        env:
        - name: PLUGIN_ENABLED
          value: "true"
        - name: PLUGIN_CONFIG_PATH
          value: "/config/plugin-config.json"
        volumeMounts:
        - name: plugin-config
          mountPath: /config
        - name: plugin-data
          mountPath: /data/plugins
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
      volumes:
      - name: plugin-config
        configMap:
          name: urldb-plugin-config
      - name: plugin-data
        persistentVolumeClaim:
          claimName: urldb-plugin-data
```

### 运维监控

**健康检查端点**
```go
// 插件系统健康检查API
func setupPluginHealthRoutes(router *gin.Engine) {
    health := router.Group("/health/plugins")
    {
        // 插件系统整体健康状态
        health.GET("/", func(c *gin.Context) {
            status := pluginManager.GetSystemHealth()
            c.JSON(200, status)
        })

        // 单个插件健康状态
        health.GET("/:pluginName", func(c *gin.Context) {
            pluginName := c.Param("pluginName")
            status, err := pluginManager.GetPluginHealth(pluginName)
            if err != nil {
                c.JSON(404, gin.H{"error": err.Error()})
                return
            }
            c.JSON(200, status)
        })

        // 插件依赖关系图
        health.GET("/dependencies", func(c *gin.Context) {
            graph := pluginManager.GetDependencyGraph()
            c.JSON(200, graph)
        })
    }
}
```

**性能监控指标**
```prometheus
# 插件系统监控指标示例
# 插件加载总数
urldb_plugin_load_total{status="success"} 42
urldb_plugin_load_total{status="failure"} 2

# 插件运行状态
urldb_plugin_status{plugin="network-disk",status="running"} 1
urldb_plugin_status{plugin="notification",status="stopped"} 1

# 插件调用延迟
urldb_plugin_call_duration_seconds_bucket{plugin="network-disk",le="0.005"} 123
urldb_plugin_call_duration_seconds_bucket{plugin="network-disk",le="0.01"} 456
urldb_plugin_call_duration_seconds_bucket{plugin="network-disk",le="0.025"} 789

# 插件内存使用
urldb_plugin_memory_bytes{plugin="network-disk"} 10485760
urldb_plugin_memory_bytes{plugin="notification"} 5242880
```

### 日志管理

**结构化日志**
```json
{
  "timestamp": "2023-12-01T10:30:45.123Z",
  "level": "INFO",
  "service": "urldb-plugin",
  "plugin": "network-disk",
  "component": "file-upload",
  "message": "文件上传成功",
  "context": {
    "file_id": "abc123",
    "file_size": 1024000,
    "user_id": "user001",
    "duration_ms": 150
  }
}
```

**日志轮转配置**
```yaml
# logrotate配置示例
/var/log/urldb/plugin/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 urldb urldb
    postrotate
        systemctl reload urldb > /dev/null 2>&1 || true
    endscript
}
```

### 故障排查

**常见问题诊断**
```bash
# 1. 检查插件系统状态
curl -s http://localhost:8080/health/plugins | jq '.'

# 2. 查看插件日志
tail -f /var/log/urldb/plugin/system.log | grep "ERROR"

# 3. 检查插件配置
cat /config/plugin-config.json | jq '.'

# 4. 监控插件性能
curl -s http://localhost:9090/metrics | grep "urldb_plugin"

# 5. 插件热重载
curl -X POST http://localhost:8080/api/plugins/reload?name=network-disk
```

**性能分析工具**
```go
// 插件性能分析
func (pm *PluginManager) ProfilePlugin(pluginName string) *PluginProfile {
    profile := &PluginProfile{
        PluginName: pluginName,
        StartTime:  time.Now(),
    }

    // 启用pprof分析
    if pm.profilerEnabled {
        profile.ProfileData = pprof.StartCPUProfile()
    }

    return profile
}

// 分析报告生成
func (pp *PluginProfile) GenerateReport() *ProfileReport {
    report := &ProfileReport{
        PluginName:    pp.PluginName,
        Duration:      time.Since(pp.StartTime),
        MemoryUsage:   pp.getMemoryUsage(),
        GoroutineCount: runtime.NumGoroutine(),
        CallStatistics: pp.getCallStats(),
    }

    if pp.ProfileData != nil {
        report.CPUProfile = pprof.StopCPUProfile(pp.ProfileData)
    }

    return report
}
```

### 备份与恢复

**配置备份策略**
```bash
#!/bin/bash
# plugin-backup.sh

BACKUP_DIR="/backup/plugins"
DATE=$(date +%Y%m%d_%H%M%S)

# 备份插件配置
mkdir -p ${BACKUP_DIR}/${DATE}
cp /config/plugin-config.json ${BACKUP_DIR}/${DATE}/
cp -r /data/plugins/* ${BACKUP_DIR}/${DATE}/plugins-data/

# 创建备份清单
cat > ${BACKUP_DIR}/${DATE}/manifest.json << EOF
{
  "backup_time": "$(date -Iseconds)",
  "plugin_configs": ["plugin-config.json"],
  "plugin_data_dirs": ["plugins-data"],
  "version": "1.0.0"
}
EOF

# 保留最近7天的备份
find ${BACKUP_DIR} -type d -mtime +7 -exec rm -rf {} \;
```

**灾难恢复流程**
```yaml
# disaster-recovery.yaml
recovery_plan:
  version: "1.0"
  steps:
    # 1. 系统状态检查
    - name: "system_check"
      command: "curl -f http://localhost:8080/health"
      timeout: "30s"

    # 2. 数据库连接检查
    - name: "database_check"
      command: "mysqladmin ping"
      timeout: "10s"

    # 3. 插件系统恢复
    - name: "plugin_recovery"
      command: "./scripts/recover-plugins.sh"
      timeout: "300s"

    # 4. 服务重启
    - name: "service_restart"
      command: "systemctl restart urldb"
      timeout: "60s"

    # 5. 健康验证
    - name: "health_check"
      command: "curl -f http://localhost:8080/health/plugins"
      timeout: "30s"

  notifications:
    failure:
      email: "admin@urldb.com"
      webhook: ""
    success:
      email: "ops@urldb.com"
```

### 升级与维护

**灰度发布策略**
```go
// 插件灰度发布管理
type PluginRollout struct {
    PluginName    string            `json:"plugin_name"`
    Version       string            `json:"version"`
    TargetPercent float64           `json:"target_percent"` // 目标发布百分比
    CurrentPercent float64          `json:"current_percent"` // 当前发布百分比
    Strategy      RolloutStrategy   `json:"strategy"`       // 发布策略
    Status        RolloutStatus     `json:"status"`
    StartTime     time.Time         `json:"start_time"`
    EndTime       *time.Time        `json:"end_time,omitempty"`
}

type RolloutStrategy struct {
    Gradual bool              `json:"gradual"`     // 是否渐进式发布
    Steps   []RolloutStep     `json:"steps"`       // 发布步骤
    Canary  *CanaryConfig     `json:"canary,omitempty"` // 金丝雀配置
}

// 渐进式发布控制器
func (pr *PluginRollout) Execute() error {
    switch pr.Strategy.Type {
    case "gradual":
        return pr.executeGradualRollout()
    case "canary":
        return pr.executeCanaryRollout()
    case "blue-green":
        return pr.executeBlueGreenRollout()
    default:
        return fmt.Errorf("不支持的发布策略: %s", pr.Strategy.Type)
    }
}
```

**版本兼容性管理**
```go
// 插件版本兼容性检查
type CompatibilityMatrix struct {
    PluginName     string                   `json:"plugin_name"`
    CoreVersion    string                   `json:"core_version"`
    CompatibleVersions []PluginVersionInfo `json:"compatible_versions"`
}

type PluginVersionInfo struct {
    Version    string   `json:"version"`
    MinCore    string   `json:"min_core_version"`
    MaxCore    string   `json:"max_core_version"`
    Deprecated bool     `json:"deprecated"`
    ReleaseNotes string `json:"release_notes"`
}

// 兼容性验证
func (cm *CompatibilityMatrix) CheckCompatibility(pluginVersion, coreVersion string) error {
    pluginInfo := cm.getPluginVersionInfo(pluginVersion)
    if pluginInfo == nil {
        return fmt.Errorf("插件版本 %s 未找到", pluginVersion)
    }

    // 检查是否已弃用
    if pluginInfo.Deprecated {
        return fmt.Errorf("插件版本 %s 已弃用", pluginVersion)
    }

    // 检查核心版本兼容性
    if !cm.isCoreVersionCompatible(coreVersion, pluginInfo.MinCore, pluginInfo.MaxCore) {
        return fmt.Errorf("插件版本 %s 与核心版本 %s 不兼容", pluginVersion, coreVersion)
    }

    return nil
}
```