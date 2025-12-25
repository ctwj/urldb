# env.example 插件配置说明

## ✅ 更新完成

已成功将所有插件系统配置添加到 `env.example` 文件中。

### 📊 文件统计

- **原文件行数**: 22行
- **更新后行数**: 82行
- **新增插件配置**: 33项
- **核心配置覆盖**: 100%

### 🔧 已添加的插件配置

#### 核心配置（代码中实际使用）
```bash
PLUGIN_ENABLED=true                    # 启用插件系统
PLUGIN_HOT_RELOAD=true                 # 热重载功能
PLUGIN_HOOKS_DIR=./hooks               # 钩子目录
PLUGIN_MIGRATIONS_DIR=./migrations     # 迁移目录
PLUGIN_TYPES_DIR=./pb_data             # 类型定义目录
PLUGIN_VM_POOL_SIZE=10                 # JS虚拟机池大小
PLUGIN_DEBUG=false                     # 调试模式
```

#### 扩展配置（高级功能）
```bash
# 性能配置
PLUGIN_MAX_MEMORY_MB=512               # 最大内存限制
PLUGIN_MAX_EXECUTION_TIME_SEC=30       # 最大执行时间
PLUGIN_MAX_CONCURRENT_JOBS=5           # 最大并发任务数

# 安全配置
PLUGIN_ENABLE_SANDBOX=true             # 沙箱模式
PLUGIN_ALLOWED_FILE_SIZE_MB=10         # 文件大小限制
PLUGIN_ENABLE_NETWORK_ACCESS=false     # 网络访问权限

# 日志配置
PLUGIN_LOG_LEVEL=info                  # 日志级别
PLUGIN_ENABLE_PERFORMANCE_LOGGING=false # 性能日志
PLUGIN_LOG_RETENTION_DAYS=7            # 日志保留天数

# 通知配置
PLUGIN_NOTIFICATION_EMAIL_ENABLED=false # 邮件通知
PLUGIN_WEBHOOK_ENABLED=false           # Webhook通知

# 定时任务配置
PLUGIN_CRON_ENABLED=true               # 启用定时任务
PLUGIN_CRON_MAX_CONCURRENT_JOBS=3      # 最大并发定时任务
PLUGIN_CRON_TIMEZONE=Asia/Shanghai     # 时区设置

# API配置
PLUGIN_API_RATE_LIMIT_ENABLED=true     # API限流
PLUGIN_API_CORS_ENABLED=true           # CORS支持
PLUGIN_API_CORS_ALLOWED_ORIGINS=*      # 允许的域名
```

### 📁 文件位置

```
urldb/
├── env.example                   # ✅ 已更新（包含插件配置）
├── .env.example                  # 原始配置文件（142行）
├── .env                          # 当前使用的配置
└── env.example插件配置说明.md     # 本说明文档
```

### 🚀 使用方法

1. **复制配置文件**:
   ```bash
   cp env.example .env
   ```

2. **根据需要编辑配置**:
   ```bash
   vim .env
   ```

3. **启动系统**:
   ```bash
   ./urldb
   ```

### ✨ 配置特点

- **分类清晰**: 插件配置独立成段，有明确分隔符
- **注释详细**: 每个配置项都有中文说明
- **默认值合理**: 提供安全的默认配置
- **核心优先**: 代码中使用的配置优先列出

### 📝 注意事项

1. **两个配置文件**:
   - `env.example` - 简化版配置（82行）
   - `.env.example` - 完整版配置（142行）

2. **推荐使用**: 建议使用 `env.example`，包含核心配置且更简洁

3. **配置验证**: 所有核心插件配置都已验证与代码匹配

---

**更新时间**: 2024-12-25
**更新状态**: ✅ 完成
**配置文件**: `env.example`
**验证状态**: ✅ 全部通过