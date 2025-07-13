# 网盘工厂模式使用说明

## 概述

本项目实现了一个优雅的网盘服务工厂模式，用于处理不同网盘平台的资源转存和文件管理。该设计采用了工厂模式和策略模式，使得代码结构清晰、易于扩展。

## 架构设计

### 核心组件

1. **PanFactory** - 网盘工厂，负责创建对应的网盘服务实例
2. **PanService** - 网盘服务接口，定义所有网盘服务必须实现的方法
3. **BasePanService** - 基础网盘服务，提供通用的HTTP请求和配置管理
4. **具体实现类** - 各网盘平台的具体实现

### 支持的服务类型

- **Quark** - 夸克网盘
- **Alipan** - 阿里云盘
- **BaiduPan** - 百度网盘（基础框架）
- **UC** - UC网盘（基础框架）

## 使用方法

### 基本使用

```go
package main

import (
    "log"
    "res_db/utils"
)

func main() {
    // 创建工厂实例
    factory := utils.NewPanFactory()
    
    // 准备配置
    config := &utils.PanConfig{
        URL:         "https://pan.quark.cn/s/123456789",
        Code:        "",
        IsType:      0, // 0: 转存并分享后的资源信息, 1: 直接获取资源信息
        ExpiredType: 1, // 1: 分享永久, 2: 临时
        AdFid:       "",
        Stoken:      "",
    }
    
    // 创建对应的网盘服务
    panService, err := factory.CreatePanService(config.URL, config)
    if err != nil {
        log.Fatalf("创建网盘服务失败: %v", err)
    }
    
    // 提取分享ID
    shareID, serviceType := utils.ExtractShareId(config.URL)
    if serviceType == utils.NotFound {
        log.Fatal("不支持的链接格式")
    }
    
    // 执行转存
    result, err := panService.Transfer(shareID)
    if err != nil {
        log.Fatalf("转存失败: %v", err)
    }
    
    if result.Success {
        log.Printf("转存成功: %s", result.Message)
        // 处理转存结果
        if resultData, ok := result.Data.(map[string]interface{}); ok {
            title := resultData["title"].(string)
            shareURL := resultData["shareUrl"].(string)
            log.Printf("标题: %s", title)
            log.Printf("分享链接: %s", shareURL)
        }
    } else {
        log.Printf("转存失败: %s", result.Message)
    }
}
```

### 在定时任务中使用

在 `utils/scheduler.go` 的 `convertReadyResourceToResource` 函数中已经集成了工厂模式：

```go
func (s *Scheduler) convertReadyResourceToResource(readyResource entity.ReadyResource) error {
    // 使用工厂模式创建对应的网盘服务
    factory := NewPanFactory()
    config := &PanConfig{
        URL:         readyResource.URL,
        Code:        "",
        IsType:      0,
        ExpiredType: 1,
        AdFid:       "",
        Stoken:      "",
    }
    
    panService, err := factory.CreatePanService(readyResource.URL, config)
    if err != nil {
        return err
    }
    
    // 根据服务类型进行不同处理
    shareID, serviceType := ExtractShareId(readyResource.URL)
    switch serviceType {
    case Quark:
        // 夸克网盘处理逻辑
    case Alipan:
        // 阿里云盘处理逻辑
    // ... 其他服务类型
    }
    
    return nil
}
```

## 扩展新的网盘服务

### 1. 创建新的服务实现

```go
// 在 utils/ 目录下创建新文件，例如 new_pan.go
package utils

// NewPanService 新网盘服务
type NewPanService struct {
    *BasePanService
}

// NewNewPanService 创建新网盘服务
func NewNewPanService(config *PanConfig) *NewPanService {
    service := &NewPanService{
        BasePanService: NewBasePanService(config),
    }
    
    // 设置请求头
    service.SetHeaders(map[string]string{
        "User-Agent": "Mozilla/5.0 ...",
        // 其他必要的请求头
    })
    
    return service
}

// GetServiceType 获取服务类型
func (n *NewPanService) GetServiceType() ServiceType {
    return NewPan // 需要在 ServiceType 中定义新的类型
}

// Transfer 转存分享链接
func (n *NewPanService) Transfer(shareID string) (*TransferResult, error) {
    // 实现转存逻辑
    return SuccessResult("转存成功", data), nil
}

// GetFiles 获取文件列表
func (n *NewPanService) GetFiles(pdirFid string) (*TransferResult, error) {
    // 实现文件列表获取逻辑
    return SuccessResult("获取成功", files), nil
}

// DeleteFiles 删除文件
func (n *NewPanService) DeleteFiles(fileList []string) (*TransferResult, error) {
    // 实现文件删除逻辑
    return SuccessResult("删除成功", nil), nil
}
```

### 2. 更新工厂

在 `utils/pan_factory.go` 中添加新的服务类型：

```go
const (
    Quark ServiceType = iota
    Alipan
    BaiduPan
    UC
    NewPan // 添加新的服务类型
    NotFound
)

// 在 String 方法中添加
func (s ServiceType) String() string {
    switch s {
    // ... 现有case
    case NewPan:
        return "newpan"
    default:
        return "unknown"
    }
}

// 在 CreatePanService 方法中添加
func (f *PanFactory) CreatePanService(url string, config *PanConfig) (PanService, error) {
    serviceType := ExtractServiceType(url)
    
    switch serviceType {
    // ... 现有case
    case NewPan:
        return NewNewPanService(config), nil
    default:
        return nil, fmt.Errorf("不支持的服务类型: %s", url)
    }
}

// 在 ExtractServiceType 方法中添加URL模式
func ExtractServiceType(url string) ServiceType {
    url = strings.ToLower(url)
    
    patterns := map[string]ServiceType{
        // ... 现有模式
        "newpan.example.com": NewPan,
    }
    
    for pattern, serviceType := range patterns {
        if strings.Contains(url, pattern) {
            return serviceType
        }
    }
    
    return NotFound
}
```

## 配置说明

### PanConfig 配置项

- **URL** - 分享链接地址
- **Code** - 分享码（如果需要）
- **IsType** - 处理类型：0=转存并分享后的资源信息，1=直接获取资源信息
- **ExpiredType** - 有效期类型：1=分享永久，2=临时
- **AdFid** - 夸克专用，分享时带上这个文件的fid
- **Stoken** - 夸克专用，分享token

### 环境配置

某些网盘服务可能需要特定的配置，如：

- **夸克网盘** - 需要配置cookie
- **阿里云盘** - 需要配置access_token
- **百度网盘** - 需要配置相关认证信息

## 错误处理

所有服务方法都返回 `*TransferResult` 和 `error`：

```go
type TransferResult struct {
    Success  bool        `json:"success"`
    Message  string      `json:"message"`
    Data     interface{} `json:"data,omitempty"`
    ShareURL string      `json:"shareUrl,omitempty"`
    Title    string      `json:"title,omitempty"`
    Fid      string      `json:"fid,omitempty"`
}
```

## 测试

运行测试：

```bash
cd utils
go test -v
```

## 注意事项

1. **并发安全** - 每个服务实例都是独立的，可以安全地在多个goroutine中使用
2. **错误重试** - 基础服务提供了带重试的HTTP请求方法
3. **资源管理** - 记得及时清理不需要的资源
4. **配置管理** - 敏感配置（如token）应该通过环境变量或配置文件管理

## 性能优化

1. **连接池** - 基础服务使用HTTP客户端连接池
2. **超时控制** - 设置了合理的请求超时时间
3. **重试机制** - 提供了可配置的重试机制
4. **缓存** - 可以实现token缓存等优化

这个工厂模式设计使得代码具有良好的可维护性和可扩展性，可以轻松添加新的网盘服务支持。 