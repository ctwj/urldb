# Data Model: 转存文件定时自动清理

**Date**: 2026-06-14
**Feature**: 002-auto-cleanup-transfer

## 实体变更总览

本期不新增独立实体表，仅在现有 `resources` 表新增字段、在 `system_configs` 表新增配置项种子数据。清理任务执行记录（spec 中标注"可选"）本期以日志输出为主，不单独建表——如后续需持久化可再扩展。

---

## 1. Resource 实体（`db/entity/resource.go`，扩展）

### 新增字段

| 字段 | Go 类型 | GORM 标签 | 语义 | 约束 |
|------|---------|-----------|------|------|
| `TransferredAt` | `*time.Time` | `gorm:"index;comment:转存完成时间"` | 转存成功写入 `fid`/`save_url` 时同步设置；nil 表示未转存 | 清理筛选主时间基准；新增索引加速筛选 |
| `CleanedAt` | `*time.Time` | `gorm:"comment:最近清理成功时间"` | 清理成功事务中写入；nil 表示从未成功清理 | 可空；用于 UI 展示与可观测性 |
| `CleanErrorMsg` | `string` | `gorm:"size:255;comment:清理失败原因"` | 最近一次清理失败原因；空字符串表示无失败或已成功 | 清理成功时同步清空 |
| `LastCleanErrorAt` | `*time.Time` | `gorm:"comment:最近清理失败时间"` | 最近一次清理失败时间 | 可空；用于 UI 展示失败时间 |

### 既有字段（复用，不改）

| 字段 | 在本特性中的角色 |
|------|------------------|
| `Fid` (`string`, size 128) | 删除目标；清理成功后清空；`fid != ""` 是进入清理候选集的必要条件 |
| `SaveURL` (`string`, size 500) | 转存后的分享链接；清理成功后清空；`save_url != ""` 同为必要条件 |
| `CkID` (`*uint`) | 账号 ID；清理时据此解析 cookie 实例，确保按正确账号删除（FR-012） |
| `PanID` (`*uint`) | 平台 ID；本期仅处理 `PanID` 指向 Quark 的资源 |
| `UpdatedAt` | 由 GORM 自动维护；不作为清理判定依据 |

### 状态判定（无独立状态字段，遵循 spec 澄清）

| 条件 | 语义 |
|------|------|
| `fid == "" && save_url == ""` | 已清理（或从未转存——区分依据是 `transferred_at` 是否为 nil） |
| `fid != "" && save_url != "" && transferred_at != nil` | 已转存，清理候选 |
| `transferred_at == nil` | 从未转存，不参与清理 |
| `clean_error_msg != ""` | 最近一次清理失败；下一轮自然重试（非专门重试机制） |

### 状态流转

```text
[未转存]                  [已转存待清理]              [已清理]
transferred_at=nil   ──转存成功──>   transferred_at set   ──清理成功──>   fid/save_url cleared
                                    fid/save_url set                   cleaned_at set
                                                                        clean_error_msg cleared
                                          │
                                          │ 清理失败
                                          ▼
                                  [已转存待清理 + 失败记录]
                                  clean_error_msg set
                                  last_clean_error_at set
                                  (fid/save_url 保留，下一轮自然再尝试)
```

注意：无"终态失败"——失败后下一轮筛选条件不变（`fid != "" && save_url != "" && transferred_at + retention ≤ now`），会自然再次纳入。

---

## 2. 系统配置（`system_configs` 表，新增种子）

### 新增配置项

| Key 常量 | Key 字符串 | 类型 | 默认值常量 | 默认值 | 说明 |
|----------|-----------|------|-----------|--------|------|
| `ConfigKeyAutoCleanupEnabled` | `auto_cleanup_enabled` | bool | `ConfigDefaultAutoCleanupEnabled` | `"false"` | 自动清理总开关 |
| `ConfigKeyAutoCleanupRetentionDays` | `auto_cleanup_retention_days` | int | `ConfigDefaultAutoCleanupRetentionDays` | `"7"` | 保留时长（天），须 > 0 |
| `ConfigKeyAutoCleanupIntervalMinutes` | `auto_cleanup_interval_minutes` | int | `ConfigDefaultAutoCleanupIntervalMinutes` | `"60"` | 调度周期（分钟），须 > 0 |

### 同步变更点

- `db/entity/system_config_constants.go`：新增 3 个 `ConfigKey` + 3 个 `ConfigDefault*` + 3 个 `ConfigResponseField*` 常量
- `db/connection.go` + `db/repo/system_config_repository.go`：种子初始化（参照现有 `ConfigKeyAutoProcessInterval` 模式，两处都需加）
- `config/config.go`：
  - 迁移 switch case（参照 line 328、365）：将 3 个新键纳入初始化逻辑
  - 校验 switch case（参照 line 373）：将新键纳入"重置为默认值"逻辑
- `db/dto/system_config.go`：
  - 请求 DTO 新增 `AutoCleanupEnabled *bool`、`AutoCleanupRetentionDays *int`、`AutoCleanupIntervalMinutes *int`
  - 响应 DTO 新增对应非指针字段
- `db/converter/system_config_converter.go`：双向映射（参照 `AutoTransferLimitDays` 模式）
- `handlers/system_config_handler.go`：校验逻辑（`retention_days` 须 > 0，`interval_minutes` 须 > 0）

### 配置校验规则

- `auto_cleanup_enabled`：接受 `true`/`false`
- `auto_cleanup_retention_days`：整数，范围 `[1, 365]`，否则拒绝并返回错误
- `auto_cleanup_interval_minutes`：整数，范围 `[1, 1440]`（最短 1 分钟，最长 24 小时），否则拒绝

---

## 3. Repository 层新增方法（`db/repo/resource_repository.go`）

| 方法 | 签名 | 用途 |
|------|------|------|
| `FindDueForCleanup` | `(retentionDays int, limit int) ([]*Resource, error)` | 筛选 `fid != '' AND save_url != '' AND transferred_at IS NOT NULL AND transferred_at + INTERVAL '? days' <= NOW()` 的资源，按 `transferred_at ASC` 排序，带 limit |
| `MarkCleaned` | `(id uint, cleanedAt time.Time) error` | 事务内：清空 `fid` + `save_url` + `clean_error_msg`，写入 `cleaned_at` |
| `MarkCleanError` | `(id uint, errMsg string, errAt time.Time) error` | 写入 `clean_error_msg` + `last_clean_error_at`，保留 `fid`/`save_url` |

**跨数据库兼容**：`transferred_at + INTERVAL` 在 SQLite 与 MySQL 语法不同，采用与 `FindDueForPush`（`telegram_channel_repository.go`）一致的"内存过滤"策略：查询时取出所有 `fid != '' AND save_url != '' AND transferred_at IS NOT NULL` 的资源，在 Go 代码中用 `time.Since(transferred_at) >= retentionDuration` 过滤。避免方言差异。

---

## 4. 调度器与服务层（新增文件，非数据模型但与查询紧密相关）

### `scheduler/cleanup_scheduler.go`（新增）

参照 `ReadyResourceScheduler` 结构：
- 嵌入 `*BaseScheduler`
- `processingMutex sync.Mutex`（防重叠）
- `running bool` 标志
- `Start()` / `Stop()` / `IsRunning()`
- 周期来源：`ConfigKeyAutoCleanupIntervalMinutes`（默认 60）
- 每轮：读 `ConfigKeyAutoCleanupEnabled`，若关闭则 return；否则调用 `services.CleanupService.Run()`

### `services/cleanup_service.go`（新增）

- 依赖注入：`repo.ResourceRepository`、`repo.SystemConfigRepository`、`repo.CksRepository`、`repo.PanRepository`
- `Run(ctx)` 流程：
  1. 读 `retention_days` 配置
  2. `resourceRepo.FindDueForCleanup(retention, limit=500)` 取候选
  3. 遍历每个资源：
     - 按 `ck_id` 解析账号实例（`CksRepository` 取 cookie）
     - 创建 `QuarkPanService`（或经 `pan.Factory`）传入 cookie
     - 调用 `DeleteFiles([]string{resource.Fid})`
     - 若成功或错误匹配"文件不存在" → `MarkCleaned`
     - 否则 → `MarkCleanError`
  4. 记录任务起止 + 处理统计到日志（Info 级别）

---

## 5. 迁移与兼容性

- `resources` 表新增 4 个字段：GORM `AutoMigrate` 在应用启动时自动添加列，既有行新字段为 NULL/空字符串，不影响存量数据
- `system_configs` 表新增 3 行种子：参照现有 `ConfigKeyAutoProcessInterval` 的初始化路径（`db/connection.go` + `db/repo/system_config_repository.go`），首次启动写入默认值
- 升级后默认 `auto_cleanup_enabled = false`：管理员必须显式开启才会执行清理，零意外删除风险
- `transferred_at` 对存量已转存资源为 nil：这些资源不会被清理（因筛选要求 `transferred_at IS NOT NULL`），需管理员知晓；新转存的资源会正确写入该字段
