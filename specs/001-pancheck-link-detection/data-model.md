# Data Model: 链接检测优化 - 接入 PanCheck 服务

**Date**: 2026-06-13

## 新增实体

### LinkCheckResult（链接检测结果缓存）

按规范化 URL 维度的共享检测缓存，供两处检测点读写。仅缓存"有效/失效"二态结论。

```go
// db/entity/link_check_result.go
package entity

import "time"

type LinkCheckResult struct {
    ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    URLHash       string    `json:"url_hash" gorm:"size:64;uniqueIndex;not null;comment:规范化URL的SHA-256十六进制"`
    NormalizedURL string    `json:"normalized_url" gorm:"size:512;not null;comment:规范化后的URL"`
    Platform      string    `json:"platform" gorm:"size:32;comment:网盘平台(quark/uc/baidu/tianyi/pan123/pan115/aliyun/xunlei/cmcc)"`
    Status        string    `json:"status" gorm:"size:16;not null;comment:有效(valid)/失效(invalid)"`
    FailReason    string    `json:"fail_reason" gorm:"size:255;comment:失效原因(仅失效时)"`
    CheckedAt     time.Time `json:"checked_at" gorm:"not null;comment:检测时间"`
    ExpiresAt     time.Time `json:"expires_at" gorm:"not null;index;comment:缓存过期时间(checked_at+TTL)"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

func (LinkCheckResult) TableName() string { return "link_check_results" }
```

**字段说明**:
| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `url_hash` | string(64) | UNIQUE, NOT NULL | 缓存键；规范化 URL 的 SHA-256 hex |
| `normalized_url` | string(512) | NOT NULL | 规范化后 URL（trim、去 fragment、scheme+host 小写、去尾斜杠） |
| `platform` | string(32) | nullable | PanCheck 识别的平台 |
| `status` | string(16) | NOT NULL | `valid` 或 `invalid`（仅此二态入缓存） |
| `fail_reason` | string(255) | nullable | 失效原因，仅 `status=invalid` 时有值 |
| `checked_at` | timestamp | NOT NULL | 本次检测时间 |
| `expires_at` | timestamp | NOT NULL, INDEX | `checked_at + TTL`；过期后允许重检 |

**状态机**:
- 缓存行只存在两种 `status`：`valid` ↔ `invalid`。重新检测得出新结论时直接 `UPDATE`（含 upsert by `url_hash`），无需历史版本。
- `pending`/`unknown` **不入库**（异常批次跳过写缓存）。

**唯一性规则**: `url_hash` 唯一（一个规范化 URL 一行缓存）。

**索引**（在 `db/connection.go` `createIndexes` 新增）:
- `idx_link_check_results_url_hash` — 已由 `uniqueIndex` tag 覆盖
- `idx_link_check_results_expires_at` — 已由 `index` tag 覆盖，供过期清理任务使用

**生命周期**:
- 写入：检测得出 valid/invalid 结论时 upsert。
- 失效/刷新：`expires_at` 过期后视为未缓存，重检覆盖；可由缓存清理任务（`scheduler/cache_cleaner.go`）定期删除已过期行。

---

## 现有实体变更

### Resource（仅使用方式变更，无结构改动）

`Resource.is_valid`（bool）继续作为业务有效性标记，由新的检测路径写入。**无新增字段**（澄清结论：二态，is_valid 足够）。失效原因不存于 `Resource`，而存于 `LinkCheckResult.fail_reason`，资源展示时按其 URL 关联查询。

### SystemConfig（新增配置键）

在 `db/entity/system_config_constants.go` 新增常量；在 `db/connection.go` `insertDefaultDataIfEmpty` 新增默认行：

| 配置键常量 | 键值 | 类型 | 默认值 |
|---|---|---|---|
| `ConfigKeyPanCheckEnabled` | `pancheck_enabled` | bool | `false` |
| `ConfigKeyPanCheckHost` | `pancheck_host` | string | `""` |
| `ConfigKeyPanCheckTimeoutSeconds` | `pancheck_timeout_seconds` | int | `60` |
| `ConfigKeyPanCheckBatchSize` | `pancheck_batch_size` | int | `20` |
| `ConfigKeyPanCheckConcurrency` | `pancheck_concurrency` | int | `5` |
| `ConfigKeyPanCheckCacheTtlHours` | `pancheck_cache_ttl_hours` | int | `24` |

---

## 迁移

- `db/connection.go` 的 `AutoMigrate(...)` 列表追加 `&entity.LinkCheckResult{}`。
- `createIndexes` 追加 `link_check_results` 的索引（`url_hash` 唯一索引、`expires_at` 索引；GORM tag 已声明，可由 AutoMigrate 自动创建，`createIndexes` 作幂等补充）。
- `insertDefaultDataIfEmpty` 追加上述 6 条 PanCheck 默认配置。
- 无需手写 SQL 迁移脚本；遵循项目既有 GORM AutoMigrate 模式。

## 实体关系

```text
Resource (1) ──url──▶ (0..1) LinkCheckResult
        以 Resource.URL 规范化后匹配 LinkCheckResult.normalized_url / url_hash
        （弱关联，无外键；查询时按 url_hash join/lookup）
```
