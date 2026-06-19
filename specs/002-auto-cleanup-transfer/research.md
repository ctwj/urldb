# Phase 0 Research: 转存文件定时自动清理

**Date**: 2026-06-14
**Feature**: 002-auto-cleanup-transfer
**Status**: All unknowns resolved

## 研究方法

本研究的依据为对现有代码库的静态分析（通过 CodeGraph AST 索引 + 关键文件直读），不涉及外部库文档查询——所有技术决策均落在"复用现有调度框架 / 配置体系 / 网盘服务"这一约束内，无新外部依赖。

---

## R1: 资源"转存完成时间"字段的来源

**Decision**: 新增 `transferred_at *time.Time` 字段到 `Resource` 实体；在转存成功写入 `SaveURL` / `Fid` 的同一事务中设置该字段。

**Rationale**:
- 现有 `Resource` 实体无精确的"转存完成时间"字段。`updated_at` 会被任意更新（包括管理员编辑、Meilisearch 同步标记等）触发，无法作为清理判定的时间基准。
- `created_at` 是资源创建时间，可能在转存之前就已存在（资源先入库后转存）。
- 清理判定公式 `transferred_at + retention_days ≤ now` 要求时间字段语义单一、只在转存成功时写入一次，新增独立字段是唯一可靠方案。
- 字段为可空指针 `*time.Time`：未转存的资源为 nil，天然不进入清理筛选。

**Alternatives considered**:
- 复用 `updated_at`：语义污染，任何更新都会推进"转存时间"，导致清理时机漂移。**否决**。
- 复用 `created_at`：资源创建与转存完成之间存在不确定延迟（人工录入、批量导入），会导致刚转存的资源被立即判定为超期。**否决**。
- 在 `task_items` 表推断：转存任务历史可作为时间来源，但跨表 JOIN 增加查询复杂度且 task 历史可能被清理。**否决**。

**影响范围**: `db/entity/resource.go`（新增字段 + GORM 标签）、`task/transfer_processor.go`（转存成功分支设置 `transferred_at`）、`db/repo/resource_repository.go`（新查询基于该字段）。

---

## R2: 清理调度器的集成方式

**Decision**: 仿照 `ReadyResourceScheduler` 模式新建 `CleanupScheduler`，通过 `Manager` 统一注册、`GlobalScheduler` 暴露启停 API，并在配置变更时由 `system_config_handler` 调用启停。

**Rationale**:
- `ReadyResourceScheduler`（`scheduler/ready_resource.go`）已示范了完整模式：嵌入 `*BaseScheduler`、`processingMutex sync.Mutex` 防重叠、`Start/Stop/IsRunning` 三件套、`time.Ticker` 周期触发、首启动立即执行一次。
- `Manager`（`scheduler/manager.go`）已聚合 `resourceRepo`、`systemConfigRepo`、`cksRepo` 等依赖，清理调度器所需依赖（resourceRepo + systemConfigRepo + 账号 cookie 来源）均已在 Manager 中可用。
- `GlobalScheduler.UpdateSchedulerStatusWithAutoTransfer` 已展示"根据配置布尔值启停调度器"的范式，新增 cleanup 开关只需扩展该方法或新增平行方法。
- 不引入 `robfig/cron`：现有调度器全部用 `time.Ticker`，保持一致性。

**Alternatives considered**:
- 复用 `ReadyResourceScheduler` 加分支：两类任务语义不同（待处理 vs 清理），耦合会导致单调度器职责膨胀、互斥锁竞争。**否决**。
- 接入 `task` 系统（像 TransferProcessor 那样）：task 系统面向"一次性任务项"，清理是"周期扫描"，模型不匹配。**否决**。
- 引入 cron 库：项目无 cron 依赖，新增仅为一个调度器不值。**否决**。

**影响范围**: `scheduler/cleanup_scheduler.go`（新增）、`scheduler/manager.go`（注册）、`scheduler/global.go`（启停方法）、`handlers/system_config_handler.go`（配置变更联动）。

---

## R3: Quark 删除能力的复用与"文件不存在"识别

**Decision**: 直接复用 `QuarkPanService.DeleteFiles([]string)` / `deleteSingleFile(fileID)`；将"文件不存在"类响应视为清理成功。

**Rationale**:
- `common/quark_pan.go:284` 的 `DeleteFiles` 已逐个调用 `deleteSingleFile`，后者 POST `https://drive-pc.quark.cn/1/clouddrive/file/delete` 并通过 `waitForTask` 等待完成——这正是清理所需的全部能力。
- `deleteSingleFile` 在 `response.Status != 200` 时返回 `fmt.Errorf("删除文件失败: %s", response.Message)`。Quark 对"文件不存在"返回的 status 与 message 需在实现阶段通过实测确认（典型表现为 status 非 200 + message 含"不存在"/"not found"），清理服务需在调用失败时检查错误字符串特征以区分"文件不存在"（视为成功）与"真实失败"（记录原因）。
- 清理时通过 `resource.fid` 直接调用 `DeleteFiles([]string{fid})`，单文件粒度便于失败定位与字段更新。

**风险与缓解**:
- Quark API 的"文件不存在"具体响应形态在 plan 阶段无法 100% 确定 → 实现阶段先以宽松匹配（message 含"不存在"/"not found"/"no such"）识别，并在日志中记录原始响应以便后续收紧匹配规则。
- `waitForTask` 是同步阻塞调用，大批量删除会拉长单轮任务 → MVP 阶段每轮处理量预计 < 1000，且 `processingMutex` 保证不重叠，可接受；后续如需可加并发上限。

**Alternatives considered**:
- 新写一个 `DeleteForCleanup` 方法：与现有 `DeleteFiles` 行为重复。**否决**。
- 批量传多 fid 一次调用：失败时无法定位是哪个 fid 出问题。**否决**（采用单文件循环，与现有 `DeleteFiles` 内部实现一致）。

**影响范围**: 无需修改 `common/quark_pan.go`；清理服务层（`services/cleanup_service.go`）封装错误识别逻辑。

---

## R4: 配置项命名与现有 `auto_transfer_*` 的区分

**Decision**: 新增独立配置键前缀 `auto_cleanup_*`，与现有 `auto_transfer_*`（转存）严格分离：

| 配置键 | 含义 | 类型 | 默认值 |
|--------|------|------|--------|
| `auto_cleanup_enabled` | 自动清理开关 | bool | `false` |
| `auto_cleanup_retention_days` | 保留时长（天） | int | `7` |
| `auto_cleanup_interval_minutes` | 调度周期（分钟） | int | `60` |

**Rationale**:
- 现有 `ConfigKeyAutoTransferLimitDays`（`auto_transfer_limit_days`）的语义经 `db/dto/system_config.go:17` 注释确认为"自动转存限制天数（0表示不限制）"——即转存时只取最近 N 天的新资源，**与"清理保留时长"完全无关**。复用该键会导致语义冲突与误操作。
- 采用分钟而非小时作为调度周期单位：与现有 `ConfigKeyAutoProcessInterval`（单位：分钟）保持一致，便于管理员形成统一心智模型；默认 60 分钟 = 1 小时符合 spec 澄清结论。
- 默认 `auto_cleanup_enabled = false`：与所有 `auto_*` 开关默认值一致，避免升级后静默启用。

**影响范围**:
- `db/entity/system_config_constants.go`：新增 3 个 `ConfigKey` + 3 个 `ConfigDefault` + 3 个 `ConfigResponseField` 常量。
- `db/connection.go` / `db/repo/system_config_repository.go`：种子数据初始化。
- `config/config.go`：迁移与校验 switch case 需新增 3 键。
- `db/dto/system_config.go` + `db/converter/system_config_converter.go`：DTO 字段 + 双向映射。

---

## R5: 资源字段清空的并发安全与状态判定

**Decision**: 清理成功时在同一 DB 事务中清空 `fid` / `save_url`、写入 `cleaned_at = now`、清空 `clean_error_msg`；失败时写入 `clean_error_msg` 与 `last_clean_error_at`，**不清空** `fid` / `save_url`。

**Rationale**:
- "已清理"判定逻辑（spec 澄清）：`fid == "" && save_url == ""` ⇒ 已清理。因此清理成功必须同时清空两者，且必须在同一事务中完成，避免中途失败导致"半清空"状态。
- `processingMutex`（调度器级）已保证同一轮任务不重叠，但资源级并发（如管理员手动编辑同一资源）仍可能发生 → 通过 GORM 行级更新（`UPDATE ... WHERE id = ?`）保证原子性，不依赖应用层锁。
- 清理失败保留 `fid` 是重试前提：下一轮筛选条件 `fid != "" && save_url != "" && transferred_at IS NOT NULL && transferred_at + retention ≤ now` 会自然再次纳入该资源。

**字段清单（Resource 新增）**:
| 字段 | 类型 | 语义 |
|------|------|------|
| `transferred_at` | `*time.Time` | 转存完成时间（R1） |
| `cleaned_at` | `*time.Time` | 最近一次清理成功时间；nil 表示从未清理成功 |
| `clean_error_msg` | `string` (size 255) | 最近一次清理失败原因；空表示无失败或已成功 |
| `last_clean_error_at` | `*time.Time` | 最近一次清理失败时间 |

**Alternatives considered**:
- 用单独 `clean_status` 枚举字段：spec 澄清已明确"不新增独立状态字段"，通过 `fid`/`save_url` 为空判定。**遵循澄清结论，否决**。
- 记录 `retry_count`：spec 澄清已明确"不实现重试机制"。**否决**。

---

## R6: "仅清理本系统转存文件"的边界保证（FR-012）

**Decision**: 清理筛选条件天然限定范围——只处理 `fid != "" && save_url != ""` 的资源，即只清理本系统写入过转存字段的记录；网盘账号中用户手动转存或原本存在的文件不在 `resources` 表中，不会被纳入。

**Rationale**:
- 系统从未获取过网盘账号的"完整文件列表"，`DeleteFiles` 仅按传入的 `fid` 删除，不存在"扫描账号全量文件后删除"的路径。
- 唯一风险点是 `fid` 与账号的对应关系：清理时必须通过 `resource.ck_id`（账号ID）获取对应 cookie 实例，而非使用任意账号 cookie 删除——否则可能跨账号误删。此约束在 `cleanup_service.go` 实现中强制：每个资源按其 `ck_id` 解析账号实例后再调用 `DeleteFiles`。

**影响范围**: `services/cleanup_service.go` 实现约束，无需额外表结构。

---

## 研究结论汇总

| 编号 | 决策 | 关键依据 |
|------|------|----------|
| R1 | 新增 `transferred_at` 字段 | `updated_at` 语义污染，需单一语义时间基准 |
| R2 | 仿 `ReadyResourceScheduler` 新建 `CleanupScheduler` | 现有调度框架成熟，保持一致性 |
| R3 | 复用 `QuarkPanService.DeleteFiles`，宽松匹配"文件不存在" | 删除能力已存在，无需重写 |
| R4 | 新增 `auto_cleanup_*` 配置前缀，与 `auto_transfer_*` 分离 | 语义不同，避免冲突 |
| R5 | 事务清空 `fid`/`save_url` + 新增 4 个字段 | 遵循澄清结论，原子性保证 |
| R6 | 按 `ck_id` 解析账号实例删除，天然限定范围 | 防跨账号误删 |

**所有 NEEDS CLARIFICATION 已解决**，可进入 Phase 1 设计。
