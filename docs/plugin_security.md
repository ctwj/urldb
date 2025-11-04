# 插件安全机制

urlDB插件系统现在包含了一套完整的安全机制，包括权限控制和行为监控，以确保插件在安全的环境中运行。

## 权限控制系统

### 权限类型

urlDB插件系统支持以下权限类型：

1. **系统权限**
   - `system:read` - 系统读取权限
   - `system:write` - 系统写入权限
   - `system:execute` - 系统执行权限

2. **数据库权限**
   - `database:read` - 数据库读取权限
   - `database:write` - 数据库写入权限
   - `database:exec` - 数据库执行权限

3. **网络权限**
   - `network:connect` - 网络连接权限
   - `network:listen` - 网络监听权限

4. **文件权限**
   - `file:read` - 文件读取权限
   - `file:write` - 文件写入权限
   - `file:exec` - 文件执行权限

5. **任务权限**
   - `task:schedule` - 任务调度权限
   - `task:control` - 任务控制权限

6. **配置权限**
   - `config:read` - 配置读取权限
   - `config:write` - 配置写入权限

7. **数据权限**
   - `data:read` - 数据读取权限
   - `data:write` - 数据写入权限

### 权限管理

插件默认具有以下基本权限：
- 读取自身配置和数据
- 写入自身数据
- 调度任务

插件可以通过`RequestPermission`方法请求额外权限，但需要管理员手动批准。

### 权限检查

插件可以通过`CheckPermission`方法检查是否具有特定权限：

```go
hasPerm, err := ctx.CheckPermission(string(security.PermissionConfigWrite))
if err != nil {
    // 处理错误
}
if !hasPerm {
    // 没有权限
}
```

## 行为监控系统

### 活动日志

系统会自动记录插件的以下活动：
- 日志输出（info, warn, error）
- 配置读写
- 数据读写
- 任务注册和执行
- 权限请求和拒绝

### 执行时间监控

系统会监控插件任务的执行时间，如果超过阈值会生成警报。

### 安全警报

当检测到以下行为时，系统会生成安全警报：
- 执行时间过长
- 数据库查询过多
- 文件操作过多
- 连接到可疑主机

### 安全报告

插件可以通过`GetSecurityReport`方法获取安全报告，报告包含：
- 插件权限列表
- 最近活动记录
- 安全警报
- 安全评分
- 安全问题和建议

## 使用示例

### 检查权限

```go
hasPerm, err := ctx.CheckPermission(string(security.PermissionConfigWrite))
if err != nil {
    ctx.LogError("Error checking permission: %v", err)
    return err
}

if !hasPerm {
    ctx.LogWarn("Plugin does not have config write permission")
    return fmt.Errorf("insufficient permissions")
}
```

### 请求权限

```go
// 请求写入配置的权限
err := ctx.RequestPermission(string(security.PermissionConfigWrite), pluginName)
if err != nil {
    ctx.LogError("Error requesting permission: %v", err)
}
```

### 获取安全报告

```go
report, err := ctx.GetSecurityReport()
if err != nil {
    ctx.LogError("Error getting security report: %v", err)
    return err
}

// 使用报告数据
```

## 安全最佳实践

1. **最小权限原则**: 只请求必需的权限
2. **输入验证**: 验证所有输入数据
3. **错误处理**: 妥善处理所有错误情况
4. **资源清理**: 及时释放使用的资源
5. **日志记录**: 记录重要的操作和事件

## 监控和审计

系统管理员可以通过以下方式监控插件活动：
- 查看插件活动日志
- 检查安全警报
- 定期审查插件权限
- 分析安全报告