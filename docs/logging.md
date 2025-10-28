# 日志系统说明

## 概述

本项目使用自定义的日志系统，支持多种日志级别、环境差异化配置和结构化日志记录。

## 日志级别

日志系统支持以下级别（按严重程度递增）：

1. **DEBUG** - 调试信息，用于开发和故障排除
2. **INFO** - 一般信息，记录系统正常运行状态
3. **WARN** - 警告信息，表示可能的问题但不影响系统运行
4. **ERROR** - 错误信息，表示系统错误但可以继续运行
5. **FATAL** - 致命错误，系统将退出

## 环境配置

### 日志级别配置

可以通过环境变量配置日志级别：

```bash
# 设置日志级别（DEBUG, INFO, WARN, ERROR, FATAL）
LOG_LEVEL=DEBUG

# 或者启用调试模式（等同于DEBUG级别）
DEBUG=true
```

默认情况下，开发环境使用DEBUG级别，生产环境使用INFO级别。

### 结构化日志

可以通过环境变量启用结构化日志（JSON格式）：

```bash
# 启用结构化日志
STRUCTURED_LOG=true
```

## 使用方法

### 基本日志记录

```go
import "github.com/ctwj/urldb/utils"

// 基本日志记录
utils.Debug("调试信息: %s", debugInfo)
utils.Info("一般信息: %s", info)
utils.Warn("警告信息: %s", warning)
utils.Error("错误信息: %s", err)
utils.Fatal("致命错误: %s", fatalErr) // 程序将退出
```

### 结构化日志记录

结构化日志允许添加额外的字段信息，便于日志分析：

```go
// 带字段的结构化日志
utils.DebugWithFields(map[string]interface{}{
    "user_id": 123,
    "action": "login",
    "ip": "192.168.1.1",
}, "用户登录调试信息")

utils.InfoWithFields(map[string]interface{}{
    "task_id": 456,
    "status": "completed",
    "duration_ms": 1250,
}, "任务处理完成")

utils.ErrorWithFields(map[string]interface{}{
    "error_code": 500,
    "error": "database connection failed",
    "component": "database",
}, "数据库连接失败: %v", err)
```

## 日志输出

日志默认输出到：
- 控制台（标准输出）
- 文件（logs目录下的app_日期.log文件）

日志文件支持轮转，单个文件最大100MB，最多保留5个备份文件，日志文件最长保留30天。

## 最佳实践

1. **选择合适的日志级别**：
   - DEBUG：详细的调试信息，仅在开发和故障排除时使用
   - INFO：重要的业务流程和状态变更
   - WARN：可预期的问题和异常情况
   - ERROR：系统错误和异常
   - FATAL：系统无法继续运行的致命错误

2. **使用结构化日志**：
   - 对于需要后续分析的日志，使用结构化日志
   - 添加有意义的字段，如用户ID、任务ID、请求ID等
   - 避免在字段中包含敏感信息

3. **性能监控**：
   - 记录关键操作的执行时间
   - 使用duration_ms字段记录毫秒级耗时

4. **安全日志**：
   - 记录所有认证和授权相关的操作
   - 包含客户端IP和用户信息
   - 记录失败的访问尝试

## 示例

```go
// 性能监控示例
startTime := time.Now()
// 执行操作...
duration := time.Since(startTime)
utils.DebugWithFields(map[string]interface{}{
    "operation": "database_query",
    "duration_ms": duration.Milliseconds(),
}, "数据库查询完成，耗时: %v", duration)

// 安全日志示例
utils.InfoWithFields(map[string]interface{}{
    "user_id": userID,
    "ip": clientIP,
    "action": "login",
    "status": "success",
}, "用户登录成功 - 用户ID: %d, IP: %s", userID, clientIP)
```