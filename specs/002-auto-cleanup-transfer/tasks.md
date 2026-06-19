---

description: "Task list for 转存文件定时自动清理 feature implementation"
---

# Tasks: 转存文件定时自动清理

**Input**: Design documents from `/specs/002-auto-cleanup-transfer/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/api-contract.md, quickstart.md

**Tests**: 项目无系统测试套件（plan.md 已述），不生成专门测试任务；以 `go build` 编译验证 + quickstart.md 冒烟路径作为每个故事的验收手段。

**Organization**: 按用户故事分组（US1=清理逻辑 MVP / US2=配置开关 UI / US3=可观测展示），支持独立实现与独立验收。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件、无未完成依赖）
- **[Story]**: 所属用户故事（US1/US2/US3）
- 描述中包含精确文件路径

## Path Conventions

本项目为 monorepo（Go 后端 + Vue/Nuxt 前端 `web/`）。后端路径相对仓库根，前端路径以 `web/` 开头。

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目已存在，无新增依赖、无新顶层目录。Setup 阶段无任务——直接进入 Foundational。

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有用户故事共享的数据层与配置基础设施。**必须先完成，US1/US2/US3 才能开始。**

**⚠️ CRITICAL**: 未完成此阶段不得开始任何用户故事。

- [X] T001 [P] 在 `db/entity/resource.go` 的 `Resource` 结构体新增 4 字段：`TransferredAt *time.Time`（`gorm:"index;comment:转存完成时间"`）、`CleanedAt *time.Time`（`gorm:"comment:最近清理成功时间"`）、`CleanErrorMsg string`（`gorm:"size:255;comment:清理失败原因"`）、`LastCleanErrorAt *time.Time`（`gorm:"comment:最近清理失败时间"`）
- [X] T002 [P] 在 `db/entity/system_config_constants.go` 新增 3 组常量：`ConfigKeyAutoCleanupEnabled`/`ConfigKeyAutoCleanupRetentionDays`/`ConfigKeyAutoCleanupIntervalMinutes`；对应 `ConfigDefaultAutoCleanupEnabled="false"`/`ConfigDefaultAutoCleanupRetentionDays="7"`/`ConfigDefaultAutoCleanupIntervalMinutes="60"`；对应 `ConfigResponseFieldAutoCleanup*` 字段名常量
- [X] T003 [P] 在 `db/connection.go` 与 `db/repo/system_config_repository.go` 的种子初始化数组中新增 3 条 cleanup 配置种子数据（参照现有 `ConfigKeyAutoProcessInterval` 模式，两处都要加）
- [X] T004 [P] 在 `config/config.go` 的初始化 switch（约 line 328、365）与校验 switch（约 line 373）中新增 3 个 cleanup 键的 case，使其参与配置初始化与默认值重置逻辑
- [X] T005 [P] 在 `db/dto/system_config.go` 的请求 DTO 新增 `AutoCleanupEnabled *bool`/`AutoCleanupRetentionDays *int`/`AutoCleanupIntervalMinutes *int`（指针类型，可选更新），响应 DTO 新增对应非指针字段
- [X] T006 [P] 在 `db/converter/system_config_converter.go` 新增 3 个 cleanup 字段的双向映射（参照 `AutoTransferLimitDays` 在 line 48、225、415、450 的模式：响应组装 + 请求转 entity + 默认值）
- [X] T007 [P] 在 `db/repo/resource_repository.go` 新增 3 个方法：`FindDueForCleanup(retentionDays int, limit int) ([]*entity.Resource, error)`（查询 `fid!='' AND save_url!='' AND transferred_at IS NOT NULL`，在 Go 内存中按 `time.Since(transferred_at) >= retention` 过滤，参照 `FindDueForPush` 跨数据库兼容模式）；`MarkCleaned(id uint, cleanedAt time.Time) error`（事务清空 `fid`+`save_url`+`clean_error_msg`，写入 `cleaned_at`）；`MarkCleanError(id uint, errMsg string, errAt time.Time) error`（写入 `clean_error_msg`+`last_clean_error_at`，保留 `fid`/`save_url`）
- [X] T008 [P] 在 `task/transfer_processor.go` 的转存成功分支（设置 `SaveURL`/`Fid` 的位置）同一事务中写入 `TransferredAt = utils.GetCurrentTime()`，确保后续清理有可信的时间基准
- [X] T009 运行 `go build ./...` 验证 Foundational 阶段编译通过

**Checkpoint**: 数据层与配置基础设施就绪，可开始用户故事实现。

---

## Phase 3: User Story 1 - 按保留时长自动清理转存文件 (Priority: P1) 🎯 MVP

**Goal**: 系统能在转存文件保存达到保留时长后，自动从 Quark 账号删除文件并清空资源转存字段。

**Independent Test**: 手动 SQL 将某资源 `transferred_at` 设为 1 天前、`auto_cleanup_retention_days=1`、`auto_cleanup_enabled=true`、`auto_cleanup_interval_minutes=1`，等待 ≤1 分钟，验证该资源 `fid`/`save_url` 被清空、`cleaned_at` 已写入、Quark 账号中文件已删除。

### Implementation for User Story 1

- [X] T010 [US1] 创建 `services/cleanup_service.go`：定义 `CleanupService` 结构（依赖 `repo.ResourceRepository`/`repo.SystemConfigRepository`/`repo.CksRepository`/`repo.PanRepository`）；实现 `Run(ctx)` 流程——读 `auto_cleanup_retention_days` 配置 → `FindDueForCleanup` 取候选 → 遍历资源按 `ck_id` 解析账号 cookie（防跨账号误删，FR-012）→ 经 `pan.Factory` 或 `QuarkPanService` 传入 cookie 调用 `DeleteFiles([]string{resource.Fid})` → 成功或错误匹配"文件不存在"/"not found"/"no such"时 `MarkCleaned`，否则 `MarkCleanError` → 日志记录任务起止 + 处理/成功/失败计数
- [X] T011 [US1] 创建 `scheduler/cleanup_scheduler.go`：参照 `scheduler/ready_resource.go` 模式实现 `CleanupScheduler`（嵌入 `*BaseScheduler`、`processingMutex sync.Mutex`、`running bool`、`Start/Stop/IsCleanupRunning`）；周期来源 `ConfigKeyAutoCleanupIntervalMinutes`（默认 60）；每轮先读 `ConfigKeyAutoCleanupEnabled`，关闭则 return，否则调用 `CleanupService.Run(ctx)`
- [X] T012 [US1] 在 `scheduler/manager.go` 注册 `CleanupScheduler`：新增 `cleanupScheduler` 字段、`StartCleanupScheduler`/`StopCleanupScheduler`/`IsCleanupRunning` 方法（参照 `HotDramaScheduler`/`ReadyResourceScheduler` 注册方式），`NewManager` 注入 `CleanupService` 依赖
- [X] T013 [US1] 在 `scheduler/global.go` 新增 `StartCleanupScheduler`/`StopCleanupScheduler`/`IsCleanupSchedulerRunning`/`UpdateSchedulerStatusWithCleanup(enabled bool)` 方法（参照 line 96-162 的 `ReadyResourceScheduler` 与 `UpdateSchedulerStatusWithAutoTransfer` 范式）
- [X] T014 [US1] 在 `handlers/system_config_handler.go` 的配置更新逻辑中新增 cleanup 三字段校验（`retention_days` 须 1-365、`interval_minutes` 须 1-1440，否则返回统一错误格式），并在配置保存成功后调用 `GlobalScheduler.UpdateSchedulerStatusWithCleanup(newEnabled)` 立即启停调度器（SC-002 立即生效）
- [X] T015 [US1] 在 `main.go` 启动流程中（参照现有 `globalScheduler.StartReadyResourceScheduler` 调用位置）根据 `auto_cleanup_enabled` 配置决定是否启动 `CleanupScheduler`
- [X] T016 [US1] 运行 `go build ./...` 编译验证，随后按 quickstart.md 第 3 节执行短周期冒烟测试（设 retention=1/interval=1，制造超期资源，验证清理成功路径）

**Checkpoint**: User Story 1 独立可用——启用配置后超期资源被自动清理。

---

## Phase 4: User Story 2 - 灵活配置保留时长与启用开关 (Priority: P2)

**Goal**: 管理员可在后台 UI 开启/关闭自动清理、设置保留时长与调度周期，配置立即生效。

**Independent Test**: 后台"功能配置"页修改 cleanup 三项配置并保存，验证后端调度器立即启停、非法值（如 0、负数、>365）被前端与后端双重拦截。

### Implementation for User Story 2

- [X] T017 [P] [US2] 在 `web/composables/useApi.ts` 扩展 `SystemConfigPayload`（请求）与响应类型，新增 `auto_cleanup_enabled?: boolean`/`auto_cleanup_retention_days?: number`/`auto_cleanup_interval_minutes?: number`（遵循项目"前端所有 API 统一封装"约定）
- [X] T018 [P] [US2] 在 `web/stores/systemConfig.ts` 新增 cleanup 相关 state 字段与 getter/setter（参照现有 `autoTransferEnabled` 等 state 模式）
- [X] T019 [US2] 在 `web/pages/admin/feature-config.vue` 新增"转存文件自动清理"配置区块：开关（`auto_cleanup_enabled`）、保留时长输入（`auto_cleanup_retention_days`，前端校验 1-365）、调度周期输入（`auto_cleanup_interval_minutes`，前端校验 1-1440）、保存按钮调用既有 `updateSystemConfig`；保存成功后刷新 store 并提示"配置已生效"
- [X] T020 [US2] 前端验证：开启/关闭开关 → 通过浏览器网络面板确认 PUT 请求与后端启停调度器联动；输入非法值 → 确认前端拦截 + 后端 400 错误提示

**Checkpoint**: User Story 2 独立可用——管理员可通过 UI 完整管理清理配置。

---

## Phase 5: User Story 3 - 清理过程可观测、可追溯 (Priority: P3)

**Goal**: 管理员可在资源列表查看已清理状态、清理时间、失败原因。

**Independent Test**: 触发一次清理任务后，资源列表正确展示"已清理"标记 + `cleaned_at`，失败资源展示 `clean_error_msg` + `last_clean_error_at`。

### Implementation for User Story 3

- [X] T021 [P] [US3] 在 `web/pages/admin/resources.vue` 的 `Resource` 接口（约 line 346）与表格列扩展：新增 `transferred_at`/`cleaned_at`/`clean_error_msg`/`last_clean_error_at` 字段；列表中以标签区分"已清理"（`fid==='' && save_url==='' && transferred_at`）/"未转存"（`transferred_at===null`）/"已转存"；失败原因以 tooltip 或展开行展示
- [X] T022 [US3] 确认后端资源列表响应（`GET /api/admin/resources`）已通过 GORM 自动序列化新字段（Foundational T001 已加 gorm 标签，通常无需额外 handler 改动）；如字段未暴露则在 `handlers/resource_handler.go` 补充响应组装
- [X] T023 [US3] 验证日志可观测性：触发清理任务后检查 `logs/` 中含任务起止时间、处理数、成功数、失败数（FR-007），不达标则补强 `cleanup_service.go` 的 Info 日志

**Checkpoint**: 全部用户故事独立可用，清理过程对管理员完全可见。

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 跨故事收敛与最终验证。

- [X] T024 [P] 运行 `go build ./...` 全量编译验证，确认无警告
- [X] T025 按 `quickstart.md` 第 1-7 节完整走查：本地拉起 → 开启清理 → 冒烟 → 失败场景 → 关闭清理 → 验收对照 SC-001~SC-005
- [X] T026 [P] 在 `ChangeLog.md` 追加本次功能版本记录（参照现有 changelog 风格）
- [X] T027 检查存量数据兼容性：确认 `resources` 表 AutoMigrate 正确新增 4 列、存量行新字段为 NULL/空、`transferred_at` 为 nil 的存量已转存资源不被误清理（筛选条件 `transferred_at IS NOT NULL` 保护）

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无任务，跳过
- **Foundational (Phase 2)**: 无外部依赖，可立即开始；**阻塞所有用户故事**
- **US1 (Phase 3)**: 依赖 Foundational 完成；提供调度器与启停联动
- **US2 (Phase 4)**: 依赖 Foundational 完成（配置项已就绪）；UI 与 US1 的后端启停联动协同（但可独立测试前端表单）
- **US3 (Phase 5)**: 依赖 Foundational 完成（字段已加）；与 US1 协同展示清理结果
- **Polish (Phase 6)**: 依赖所有目标用户故事完成

### User Story Dependencies

- **US1 (P1)**: Foundational 后即可开始，无跨故事依赖——后端配置可直接 SQL 写入测试
- **US2 (P2)**: Foundational 后即可开始；前端 UI 独立可测（表单校验、请求发送），与 US1 后端联调可在 US1 完成后进行
- **US3 (P3)**: Foundational 后即可开始；展示依赖清理结果数据，可与 US1 并行开发、US1 完成后联调

### Within Each User Story

- 数据层（Foundational）→ 服务层 → 调度器/handler → 前端
- 编译验证在每个故事末尾
- 冒烟验证在 US1 末尾（功能核心）

### Parallel Opportunities

- Foundational 中 T001-T008（除 T009 编译验证）全部 [P]，可并行（不同文件）
- US2 中 T017/T018 [P] 可并行
- US3 中 T021/T026 [P] 可并行
- US2 与 US3 可由不同开发者并行（前端不同页面）
- 不同用户故事可在 Foundational 完成后并行推进

---

## Parallel Example: Foundational

```bash
# 以下任务作用于不同文件，可并行：
Task T001: "Resource 新增 4 字段 in db/entity/resource.go"
Task T002: "系统配置常量 in db/entity/system_config_constants.go"
Task T003: "种子数据 in db/connection.go + db/repo/system_config_repository.go"
Task T004: "config.go 迁移/校验 in config/config.go"
Task T005: "DTO 字段 in db/dto/system_config.go"
Task T006: "converter 映射 in db/converter/system_config_converter.go"
Task T007: "Repository 方法 in db/repo/resource_repository.go"
Task T008: "设置 transferred_at in task/transfer_processor.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 2: Foundational（CRITICAL - 阻塞所有故事）
2. 完成 Phase 3: User Story 1（后端清理闭环）
3. **STOP and VALIDATE**: 通过 SQL 写入配置 + 制造超期资源，验证清理成功路径
4. 此时已具备"自动清理"核心能力（配置可通过 SQL/API 调整）

### Incremental Delivery

1. Foundational → 数据与配置基础就绪
2. +US1 → 后端清理 MVP（可先用 curl/API 操作配置）→ 验收
3. +US2 → 管理员可在 UI 完整管理配置 → 验收
4. +US3 → 清理状态对管理员可见 → 验收
5. Polish → 全量编译 + quickstart 走查 + changelog

### Parallel Team Strategy

单人推荐顺序执行（Foundational → US1 → US2 → US3 → Polish）。
多人可在 Foundational 完成后并行 US2（前端配置页）与 US3（前端资源页），US1 因涉及调度器/handler/manager 多处协同建议单人串行。

---

## Notes

- [P] 任务 = 不同文件、无未完成依赖
- [Story] 标签将任务映射到 spec.md 的用户故事
- 每个用户故事独立可完成、可验收
- 编译验证（`go build ./...`）是 golang 改动的硬性门槛（CLAUDE.md 约定）
- 前端改动无需重启服务（CLAUDE.md 约定）
- 提交粒度：每个任务或逻辑分组一个 commit
- 可在任何 Checkpoint 停下独立验收某故事
- 避免：模糊任务、同文件冲突、破坏独立性的跨故事依赖
