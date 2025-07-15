# 待处理资源自动处理功能修复说明

## 问题描述

在管理后台开启"待处理资源自动处理"功能后，系统没有自动开始任务，并且出现外键约束错误：
```
ERROR: insert or update on table "resources" violates foreign key constraint "fk_resources_pan" (SQLSTATE 23503)
```

## 问题原因

1. **UpdateSchedulerStatus方法参数不完整**：`utils/global_scheduler.go` 中的 `UpdateSchedulerStatus` 方法只接收了 `autoFetchHotDramaEnabled` 参数，没有处理 `autoProcessReadyResources` 参数。

2. **UpdateSystemConfig调用参数不完整**：`handlers/system_config_handler.go` 中的 `UpdateSystemConfig` 函数只传递了热播剧相关的参数，没有传递待处理资源相关的参数。

3. **调度器间隔时间硬编码**：`utils/scheduler.go` 中的 `StartReadyResourceScheduler` 方法使用了硬编码的5分钟间隔，而不是使用系统配置中的 `AutoProcessInterval`。

4. **平台匹配机制优化**：使用 `serviceType` 来匹配平台，并添加平台映射缓存提高性能。

## 修复内容

### 1. 修复 UpdateSchedulerStatus 方法

**文件**: `utils/global_scheduler.go`

**修改前**:
```go
func (gs *GlobalScheduler) UpdateSchedulerStatus(autoFetchHotDramaEnabled bool) {
    // 只处理热播剧功能
}
```

**修改后**:
```go
func (gs *GlobalScheduler) UpdateSchedulerStatus(autoFetchHotDramaEnabled bool, autoProcessReadyResources bool) {
    // 处理热播剧自动拉取功能
    if autoFetchHotDramaEnabled {
        if !gs.scheduler.IsRunning() {
            log.Println("系统配置启用自动拉取热播剧，启动定时任务")
            gs.scheduler.StartHotDramaScheduler()
        }
    } else {
        if gs.scheduler.IsRunning() {
            log.Println("系统配置禁用自动拉取热播剧，停止定时任务")
            gs.scheduler.StopHotDramaScheduler()
        }
    }

    // 处理待处理资源自动处理功能
    if autoProcessReadyResources {
        if !gs.scheduler.IsReadyResourceRunning() {
            log.Println("系统配置启用自动处理待处理资源，启动定时任务")
            gs.scheduler.StartReadyResourceScheduler()
        }
    } else {
        if gs.scheduler.IsReadyResourceRunning() {
            log.Println("系统配置禁用自动处理待处理资源，停止定时任务")
            gs.scheduler.StopReadyResourceScheduler()
        }
    }
}
```

### 2. 修复 UpdateSystemConfig 函数

**文件**: `handlers/system_config_handler.go`

**修改前**:
```go
scheduler.UpdateSchedulerStatus(req.AutoFetchHotDramaEnabled)
```

**修改后**:
```go
scheduler.UpdateSchedulerStatus(req.AutoFetchHotDramaEnabled, req.AutoProcessReadyResources)
```

### 3. 修复调度器间隔时间配置

**文件**: `utils/scheduler.go`

**修改前**:
```go
ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
```

**修改后**:
```go
// 获取系统配置中的间隔时间
config, err := s.systemConfigRepo.GetOrCreateDefault()
interval := 5 * time.Minute // 默认5分钟
if err == nil && config.AutoProcessInterval > 0 {
    interval = time.Duration(config.AutoProcessInterval) * time.Minute
}

ticker := time.NewTicker(interval)
defer ticker.Stop()

log.Printf("待处理资源自动处理任务已启动，间隔时间: %v", interval)
```

### 4. 优化平台匹配机制

**文件**: `utils/scheduler.go`

**新增平台映射缓存**:
```go
type Scheduler struct {
    // ... 其他字段 ...
    
    // 平台映射缓存
    panCache     map[string]*uint // serviceType -> panID
    panCacheOnce sync.Once
}
```

**新增初始化缓存方法**:
```go
// initPanCache 初始化平台映射缓存
func (s *Scheduler) initPanCache() {
    s.panCacheOnce.Do(func() {
        // 获取所有平台数据
        pans, err := s.panRepo.FindAll()
        if err != nil {
            log.Printf("初始化平台缓存失败: %v", err)
            return
        }

        // 建立 ServiceType 到 PanID 的映射
        serviceTypeToPanName := map[string]string{
            "quark":   "quark",
            "alipan":  "aliyun", // 阿里云盘在数据库中的名称是 aliyun
            "baidu":   "baidu",
            "uc":      "uc",
            "unknown": "other",
        }

        // 创建平台名称到ID的映射
        panNameToID := make(map[string]*uint)
        for _, pan := range pans {
            panID := pan.ID
            panNameToID[pan.Name] = &panID
        }

        // 建立 ServiceType 到 PanID 的映射
        for serviceType, panName := range serviceTypeToPanName {
            if panID, exists := panNameToID[panName]; exists {
                s.panCache[serviceType] = panID
                log.Printf("平台映射缓存: %s -> %s (ID: %d)", serviceType, panName, *panID)
            } else {
                log.Printf("警告: 未找到平台 %s 对应的数据库记录", panName)
            }
        }

        // 确保有默认的 other 平台
        if otherID, exists := panNameToID["other"]; exists {
            s.panCache["unknown"] = otherID
        }

        log.Printf("平台映射缓存初始化完成，共 %d 个映射", len(s.panCache))
    })
}
```

**新增根据服务类型获取平台ID的方法**:
```go
// getPanIDByServiceType 根据服务类型获取平台ID
func (s *Scheduler) getPanIDByServiceType(serviceType panutils.ServiceType) *uint {
    s.initPanCache()
    
    serviceTypeStr := serviceType.String()
    if panID, exists := s.panCache[serviceTypeStr]; exists {
        return panID
    }
    
    // 如果找不到，返回 other 平台的ID
    if otherID, exists := s.panCache["other"]; exists {
        log.Printf("未找到服务类型 %s 的映射，使用默认平台 other", serviceTypeStr)
        return otherID
    }
    
    log.Printf("未找到服务类型 %s 的映射，且没有默认平台，返回nil", serviceTypeStr)
    return nil
}
```

**修改资源创建逻辑**:
```go
// 在 convertReadyResourceToResource 方法中
resource := &entity.Resource{
    Title:       title,
    Description: readyResource.Description,
    URL:         shareURL,
    PanID:       s.getPanIDByServiceType(serviceType), // 使用 serviceType 匹配
    IsValid:     true,
    IsPublic:    true,
}
```

### 5. 添加 PanRepository 依赖

**文件**: `utils/scheduler.go`

**修改前**:
```go
type Scheduler struct {
    doubanService        *DoubanService
    hotDramaRepo         repo.HotDramaRepository
    readyResourceRepo    repo.ReadyResourceRepository
    resourceRepo         repo.ResourceRepository
    systemConfigRepo     repo.SystemConfigRepository
    stopChan             chan bool
    isRunning            bool
    readyResourceRunning bool
    processingMutex      sync.Mutex
    hotDramaMutex        sync.Mutex
}
```

**修改后**:
```go
type Scheduler struct {
    doubanService        *DoubanService
    hotDramaRepo         repo.HotDramaRepository
    readyResourceRepo    repo.ReadyResourceRepository
    resourceRepo         repo.ResourceRepository
    systemConfigRepo     repo.SystemConfigRepository
    panRepo              repo.PanRepository
    stopChan             chan bool
    isRunning            bool
    readyResourceRunning bool
    processingMutex      sync.Mutex
    hotDramaMutex        sync.Mutex
    
    // 平台映射缓存
    panCache     map[string]*uint // serviceType -> panID
    panCacheOnce sync.Once
}
```

## 修复效果

现在当您在管理后台开启"待处理资源自动处理"功能时：

1. **系统会立即启动调度器** - 不再需要重启服务器
2. **使用配置的间隔时间** - 不再是固定的5分钟
3. **支持实时开关** - 可以随时开启或关闭功能
4. **正确的外键关联** - 不再出现外键约束错误
5. **智能平台识别** - 根据 serviceType 自动识别对应的平台
6. **高性能缓存** - 平台映射缓存，避免重复数据库查询

## 测试方法

运行测试脚本验证修复效果：
```bash
# 测试自动处理功能
chmod +x test-auto-process.sh
./test-auto-process.sh

# 测试平台匹配机制
chmod +x test-pan-mapping.sh
./test-pan-mapping.sh
```

## 注意事项

1. 确保数据库中有默认的平台数据
2. 确保系统配置中的 `auto_process_interval` 设置合理
3. 如果仍有问题，检查日志中的错误信息 