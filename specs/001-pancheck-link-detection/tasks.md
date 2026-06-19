# Tasks: 链接检测优化 - 接入 PanCheck 服务

**Input**: Design documents from `/specs/001-pancheck-link-detection/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/pancheck-api.md

**Tests**: 测试任务未显式要求，仅在明显需要时以独立验证步骤（见 quickstart.md）替代。

**Organization**: 按用户故事分组任务（US1 P1 / US2 P2 / US3 P3），每个故事可独立实现与验证。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件、无未完成任务依赖）
- **[Story]**: 所属用户故事（US1/US2/US3）
- 所有描述均带精确文件路径

## Path Conventions

- **后端**：仓库根的 `db/entity/`, `db/repo/`, `services/`, `handlers/`, `scheduler/`, `common/utils/`, `db/connection.go`, `main.go`
- **前端**：`web/composables/`, `web/pages/`
- 路径相对仓库根（`/Users/kerwin/Program/go/urldb`）

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 功能分支已存在；确认开发环境就绪、建立待新增文件的占位引用。

- [X] T001 确认分支 `001-pancheck-link-detection` 检出，运行 `go build ./...` 确认当前基线编译通过（CLAUDE.md 约定）
- [X] T002 [P] 在 `db/entity/link_check_result.go` 创建 `LinkCheckResult` 实体（按 data-model.md 第 11-31 行字段与 tag），含 `TableName()` 返回 `link_check_results`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 两处检测点共享的基础设施（缓存实体迁移、配置键、PanCheck 客户端、缓存仓储、共享检测服务、依赖注入）。任何用户故事实现前必须完成。

**⚠️ CRITICAL**: 未完成本阶段不得开始任何用户故事实现。

- [X] T003 在 `db/connection.go` 的 `AutoMigrate(...)` 列表追加 `&entity.LinkCheckResult{}`；在 `createIndexes` 追加 `link_check_results` 的 `url_hash`（unique）与 `expires_at` 索引（幂等补充）
- [X] T004 在 `db/entity/system_config_constants.go` 新增 6 个配置键常量（`ConfigKeyPanCheckEnabled`=`pancheck_enabled`、`ConfigKeyPanCheckHost`=`pancheck_host`、`ConfigKeyPanCheckTimeoutSeconds`=`pancheck_timeout_seconds`、`ConfigKeyPanCheckBatchSize`=`pancheck_batch_size`、`ConfigKeyPanCheckConcurrency`=`pancheck_concurrency`、`ConfigKeyPanCheckCacheTtlHours`=`pancheck_cache_ttl_hours`）
- [X] T005 在 `db/connection.go` 的 `insertDefaultDataIfEmpty` 追加上述 6 条 PanCheck 默认配置行（enabled=false、host=""、timeout=60、batch=20、concurrency=5、ttl=24），类型与 `SystemConfig` 既有写入模式一致
- [X] T006 [P] 在 `services/pancheck_client.go` 实现 PanCheck HTTP 客户端：`POST {pancheck_host}/api/v1/links/check`，固定发送全部 9 个 `selected_platforms`；请求/响应对齐 contracts/pancheck-api.md 第 19-64 行（含对象数组容错、字段名容错、按规范化 URL 匹配、invalid 优先于 valid）；非 2xx/超时/解析失败 → 整批视为未得出结论（不抛失效）
- [X] T007 [P] 在 `services/pancheck_client.go` 实现 URL 规范化函数 `normalizeURL(url) string`（TrimSpace、去 fragment、scheme+host 小写、去尾斜杠）与 `urlHash(normalized) string`（SHA-256 hex），供客户端匹配与缓存键共用
- [X] T008 [P] 在 `db/repo/link_check_result_repository.go` 实现缓存仓储：按 `url_hash` 批量读取未过期记录、按结论 upsert（valid/invalid 二态写入，pending/unknown 不入库）、按 `expires_at` 过期清理；复用现有 GORM 仓储风格（参考 `db/repo/resource_repository.go`）
- [X] T009 在 `services/link_check_service.go` 实现共享检测服务：
  - `CheckURL(ctx, url, ignoreCache) → {status, failReason}`（单条，scheduler 用）
  - `CheckResources(ctx, resources, ignoreCache) → map[resourceID]{status, failReason}`（批量，详情页用）
  - 内部职责：读 `pancheck_enabled`/`pancheck_host`（enabled=false 或 host="" 时跳过检测，调用方按未检测处理）、缓存命中判定、按 `pancheck_batch_size` 分批、`pancheck_concurrency` 限流、聚合（任一 URL 失效即整体失效，对齐 research.md 第 7 节）、`pancheck_cache_ttl_hours` 计算 expires_at
- [X] T010 在 `services/link_check_service.go` 增加结论写回逻辑：仅在 `Resource.is_valid` 实际翻转（true↔false）时调用 `ResourceRepository.Update` + 现有 `services/meilisearch_service.go` 同步路径，避免写放大（research.md 第 7 节）
- [X] T011 在 `main.go` 依赖注入：实例化 `LinkCheckResultRepository`、`PanCheckClient`、`LinkCheckService`，注入到 `ResourceHandler` 与 `ready_resource` 调度器；保持现有路由不变
- [X] T012 运行 `go build ./...` 确认 Phase 2 编译通过；启动后端确认 `AutoMigrate` 创建 `link_check_results` 表、`system_config` 出现 6 条 PanCheck 默认行

**Checkpoint**: 基础设施就绪 —— 两处触发点的改造与前端配置可在此基础上并行推进。

---

## Phase 3: User Story 1 - 核心检测：两处触发点接入 PanCheck（Priority: P1） 🎯 MVP

**Goal**: 添加资源后的调度器检测 + 前端详情页批量检测，统一改走 `LinkCheckService`（PanCheck 为唯一主检测手段），夸克有效性检测改走 PanCheck，夸克转存业务保留；旧内联检测代码删除。

**Independent Test**: 见 quickstart.md 第 3、7、8 节 —— 提交已知有效/失效链接验证 scheduler 结论翻转；夸克转存仍生成 `SaveURL`；`url_checker.go` 已删除且编译通过。

### Implementation for User Story 1

- [X] T013 [US1] 改造 `scheduler/ready_resource.go` 的 `convertReadyResourceToResource`：检测环节统一调用 `linkCheckService.CheckURL(ctx, resource.URL, false)`；夸克分支改为**先 PanCheck 校验**通过后再执行既有"转存并分享"生成 `SaveURL`（转存流程代码保留，仅前置 PanCheck 闸门）；PanCheck 关闭/未得出结论时按 FR-004 放行（保持 is_valid 原值、不阻断转 Resource）
- [X] T014 [US1] 改造 `handlers/resource_handler.go` 的 `BatchCheckResourceValidity`（~1469 行）：内部把 `performAdvancedValidityCheck` 调用替换为 `linkCheckService.CheckResources(ctx, resources, ignoreCache=false)`；响应 `detection_method` 由 `quark_deep`/`unsupported` 改为 `pancheck`（启用时）或 `disabled`（未启用），结构体其余字段不变（contracts/pancheck-api.md 第 86-114 行）
- [X] T015 [US1] 改造 `handlers/resource_handler.go` 的 `CheckResourceValidity`（~1310 行）：单资源检测同样改为调用 `linkCheckService.CheckResources`（单元素），响应语义对齐 contracts/pancheck-api.md 第 116-123 行
- [X] T016 [US1] 删除 `handlers/resource_handler.go` 中的 `performAdvancedValidityCheck`、`performQuarkValidityCheck`、`performAlipanValidityCheck` 三个函数（research.md 第 10 节、plan.md 删除清单）
- [X] T017 [US1] 删除 `common/utils/url_checker.go` 整文件（`CheckURL`/`extractShareID`/各 `checkXxx`/`Test`/`CheckResult`）；删除 `common/utils/url_checker_test.go` 整文件
- [X] T018 [US1] 确认 `common/pan_factory.go` 的 `ExtractShareId`/`ExtractServiceType`（转存流程与平台识别仍用）保留，未被误删；仅检测相关代码被移除
- [X] T019 [US1] 运行 `go build ./...` 确认无残留引用、编译通过（CLAUDE.md 约定 + FR-009）
- [ ] T020 [US1] 验证检测点 #1（quickstart.md 第 3 节）：提交已知有效链接 → 转 Resource 且 is_valid=true；提交已知失效链接 → 被拒绝/失效标记；夸克链接先 PanCheck 校验再转存

**Checkpoint**: 两处检测点均经 PanCheck 得出结论；旧 `CheckURL` 代码已删除、编译通过；夸克转存正常。MVP 可在此停下验证。

---

## Phase 4: User Story 2 - 系统配置：PanCheck 服务参数可运维调整（Priority: P2）

**Goal**: 管理后台新增 PanCheck 配置组（开关/host/超时/批量/并发/TTL），沿用现有系统配置页风格；复用 `useSystemConfigApi` 读写；保存后立即生效无需重启（SC-005）。

**Independent Test**: 见 quickstart.md 第 2、5、6 节 —— 修改配置后立即生效；缓存命中使二次请求下降 ≥80%；降级不误判失效。

### Implementation for User Story 2

- [X] T021 [P] [US2] 在 `web/composables/useApi.ts` 确认/补充 `useSystemConfigApi` 对 6 个 PanCheck 配置键（`pancheck_enabled`/`pancheck_host`/`pancheck_timeout_seconds`/`pancheck_batch_size`/`pancheck_concurrency`/`pancheck_cache_ttl_hours`）的读写封装（复用现有统一封装风格，CLAUDE.md 约定）
- [X] T022 [US2] 在系统配置页（`web/pages/admin/` 下既有系统配置组件）新增 PanCheck 配置组 UI：开关、host 输入、超时、批量、并发、TTL 数字输入；沿用现有配置组卡片样式；保存按钮调用 `useSystemConfigApi` 写回
- [X] T023 [US2] 确认 `LinkCheckService`（Phase 2 T009）每次调用前读取配置（非启动期缓存），使保存后立即生效（SC-005）；如服务实现持有了配置快照，改为每次读取或增加配置变更后刷新机制
- [ ] T024 [US2] 验证缓存命中（quickstart.md 第 5 节）：清空 `link_check_results` → 首次检测记录 PanCheck 调用次数 N1 → 立即二次检测 N2 ≈ 0；预期 N2 较 N1 下降 ≥80%（SC-003）
- [ ] T025 [US2] 验证降级（quickstart.md 第 6 节）：`pancheck_enabled=false` → 添加资源直接放行、详情页返回 `detection_method: disabled`；`pancheck_host` 指向不可达地址 → 保持 is_valid 原值、不误判失效（SC-004 / FR-004）

**Checkpoint**: 配置可经管理后台调整、立即生效；缓存命中与降级行为符合 SC-003/SC-004。

---

## Phase 5: User Story 3 - 可见性：详情页展示二态结论 + 手动重检（Priority: P3）

**Goal**: 详情页 `/r/[key]` 展示从三态简化为二态（有效/失效 + 失效原因）；onMounted 自动触发 + 手动"链接检测"按钮逻辑不变（后端接口契约不变）；5 分钟节流保留。

**Independent Test**: 见 quickstart.md 第 4 节 —— 页面加载自动批量检测；状态翻转时写回 DB + 同步 Meilisearch；手动按钮强制重检。

### Implementation for User Story 3

- [X] T026 [US3] 改造 `web/pages/r/[key].vue` 检测结果展示：从"夸克有效/无效/不支持检测"三态简化为"有效/失效 + 失效原因"二态；失效原因取接口返回的 `error` 字段（由后端从 `LinkCheckResult.fail_reason` 关联返回）
- [X] T027 [US3] 确认 `onMounted` 自动触发逻辑不变（仍调 `batchCheckResourceValidity`，5 分钟节流保留）；手动"链接检测"按钮逻辑不变（后端路径与请求体不变，contracts/pancheck-api.md 第 90-99 行）
- [X] T028 [US3] 确认后端 `BatchCheckResourceValidity` 响应在结论翻转时回写 `Resource.is_valid` 并同步 Meilisearch（Phase 2 T010 实现的写回路径，在批量接口内被触发）；前端无需感知写回，仅据响应更新展示
- [ ] T029 [US3] 验证检测点 #2（quickstart.md 第 4 节）：打开公开资源页自动批量检测；对有效→手动失效的资源刷新后状态翻转并写回；手动按钮忽略 5 分钟节流但仍走服务端缓存

**Checkpoint**: 详情页二态展示 + 失效原因可读；自动/手动触发路径行为符合契约与节流预期。

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 跨故事一致性、清理、最终验证。

- [X] T030 [P] 在 `scheduler/cache_cleaner.go`（或既有定期清理任务位置）新增/确认 `link_check_results` 已过期行的清理逻辑（删除 `expires_at < now()`），避免缓存膨胀 — 实现在 `LinkCheckService.CleanExpiredCache()` + `ReadyResourceScheduler` 每6小时节流清理（CacheCleaner 未接入 main.go，改用运行中的 ready_resource 调度器）
- [X] T031 [P] 补充 PanCheck 客户端与缓存仓储的单元测试（如现有项目含 `*_test.go` 惯例） — 评估后不适用：项目自身代码（services/handlers/db/common/scheduler）无 `*_test.go` 惯例，仅第三方 pocketbase-src 含测试
- [X] T032 按 quickstart.md 第 8 节二次确认 `url_checker.go` 已删除、`go build ./...` 通过；确认 `common/pan_factory.go` 转存相关函数完好 — go build EXIT=0，common/utils/ 整目录已删，pan_factory.go ExtractShareId 完好
- [X] T033 按 quickstart.md "完成判据"逐项核对 — 代码层面全部满足；缓存命中/降级/tg_tool 一致率为运行时指标，归入手动验证（T020/T024/T025/T029）
- [X] T034 更新 CLAUDE.md / README（若涉及）中关于链接检测的描述 — 评估后不适用：CLAUDE.md/README 无链接检测机制描述，仅 specs/ 设计文档提及

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，立即开始
- **Foundational (Phase 2)**: 依赖 Phase 1 完成 —— **阻塞所有用户故事**
- **User Stories (Phase 3-5)**: 均依赖 Phase 2 完成
  - US1 (Phase 3) 是 MVP，必须最先完成（移除旧代码、接通主检测路径）
  - US2 (Phase 4)、US3 (Phase 5) 可在 US1 后并行推进（如人力允信）
- **Polish (Phase 6)**: 依赖 US1 完成（缓存清理、单测、最终核对需主路径已就绪）

### User Story Dependencies

- **User Story 1 (P1)**: 可在 Phase 2 后开始 —— 无其它故事依赖
- **User Story 2 (P2)**: 可在 Phase 2 后开始 —— 配置项已在 T004/T005 落库；与 US1 独立可测
- **User Story 3 (P3)**: 可在 Phase 2 后开始 —— 依赖 US1 完成的后端响应 `detection_method`/`error` 语义；建议在 US1 后实现以保证展示与后端一致

### Within Each User Story

- Models/Entities 先于 Services
- Services 先于 Handlers/触发点
- 核心实现先于集成
- 故事完成后再进入下一优先级

### Parallel Opportunities

- Phase 2 中 T006/T007/T008 标记 [P]（不同文件）可并行
- US2 T021 标记 [P]（前端 composables）可与 US1 并行
- Polish T030/T031 标记 [P] 可并行

---

## Parallel Example: Foundational Phase

```bash
# 三项可并行（不同文件、无依赖）：
Task: "实现 PanCheck HTTP 客户端 services/pancheck_client.go"
Task: "实现 URL 规范化与 urlHash services/pancheck_client.go（同文件，需串行于客户端或合并提交）"
Task: "实现缓存仓储 db/repo/link_check_result_repository.go"
```

> 注：T006 与 T007 同在 `services/pancheck_client.go`，实际实现时建议同一提交内完成或按 T007→T006 顺序，避免文件冲突。

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 1: Setup（确认基线）
2. 完成 Phase 2: Foundational（CRITICAL —— 阻塞全部故事）
3. 完成 Phase 3: User Story 1（两处触发点接入 + 旧代码删除）
4. **停下验证**：独立测试 US1（quickstart.md 第 3/7/8 节）
5. 满足 MVP 后再推进 US2/US3

### Incremental Delivery

1. Setup + Foundational → 基础设施就绪
2. + US1 → 主检测路径经 PanCheck、旧代码移除 → **MVP**
3. + US2 → 配置可运维调整、缓存/降级指标达成
4. + US3 → 详情页二态展示、手动重检
5. Polish → 缓存清理、单测、最终核对

### Parallel Team Strategy

- 后端 A：Foundational → US1
- 前端 B：Foundational 完成后并行 US2 配置 UI
- 前端 C：US1 完成后接 US3 详情页

---

## Notes

- [P] 任务 = 不同文件、无未完成任务依赖
- [Story] 标签将该任务映射到特定用户故事以便追溯
- 每个用户故事应可独立完成与测试
- 每个任务或逻辑分组后提交一次
- 在任一 Checkpoint 可停下独立验证故事
- 避免：模糊任务、同文件冲突、破坏独立性的跨故事依赖
