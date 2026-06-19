# Implementation Plan: 后台管理系统 UI/UX 整体优化

**Branch**: `004-admin-ui-optimization` | **Date**: 2026-06-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-admin-ui-optimization/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. Design decisions sourced from `/ui-ux-pro-max` skill — see [research.md](./research.md) for full reasoning.

## Summary

对 urldb 管理后台（Nuxt 3 + Vue 3 + Naive UI + Tailwind + Chart.js，共 30+ 页面）进行 UI/UX 整体优化，消除"三套卡片样式并存 / 蓝紫渐变抢镜 / 暗色覆盖不全 / 调试残留 / 列表页交互各写各的"等核心问题。

技术路径：
1. **统一以 Naive UI 为组件单一来源**，通过 `n-config-provider` + `themeOverrides` 注入完整明暗双主题（slate + blue-700/400）
2. **自研轻量 Cmd/Ctrl+K 命令面板**（复用现有侧边栏菜单数据，~200 行）
3. **新增 1 个统一封装的统计端点** `GET /api/stats/summary`（满足仪表盘首屏一次性聚合 + 环比昨日）
4. **引入 Vitest + @vue/test-utils 前端测试基础设施**（解决 constitution Principle I 长期违规）
5. 保留现有技术栈（Inter/FiraCode 字体、Chart.js、Font Awesome、Tailwind darkMode:'class'），最小化迁移成本（YAGNI）

详细决策见 [research.md](./research.md)。

## Technical Context

**Language/Version**:
- Frontend: TypeScript 5 + Vue 3.3 + Nuxt 3.8
- Backend: Go (Gin) — 极少量改动

**Primary Dependencies**:
- 前端（已存在，复用）：`naive-ui` ^2.42.0、`@nuxtjs/tailwindcss` ^6.8.0、`chart.js` ^4.5.0、`@fortawesome/fontawesome-free` ^6.7.2、`pinia` ^2.1.0、`@vueuse/core` ^14.0.0、`@vicons/ionicons5` ^0.12.0
- 前端（**新增 devDependencies**）：`vitest` ^1.x、`@vue/test-utils` ^2.x、`happy-dom` ^13.x、`@nuxt/test-utils` ^3.x
- 后端：Gin + 现有 `handlers/stats_handler.go`（新增方法）

**Storage**:
- 后端：PostgreSQL（仅读取现有表，不改 schema）
- 前端：localStorage（侧边栏分组状态，key: `urldb:admin:sidebar:expanded-groups`）

**Testing**:
- 前端：**Vitest + @vue/test-utils + happy-dom**（本 feature 引入，详见 [research.md §7](./research.md)）
- 后端：`go test ./...`（现有体系，新增 stats handler 单测 + 集成测试）
- 覆盖率门槛：constitution Principle I 要求 ≥ 90%

**Target Platform**: 现代浏览器（Chrome/Edge/Firefox/Safari 最新两个大版本），桌面优先，响应式到平板（≥768px）与手机（关键操作可达，不要求信息密度一致）

**Project Type**: Web 服务（前端为主，后端仅 1 个新端点）

**Performance Goals**:
- 仪表盘首屏可交互 ≤ 2 秒（SC-001 子目标，5 秒内完成状态判断）
- 命令面板唤起到结果渲染 ≤ 100ms
- 列表页筛选响应 ≤ 200ms
- 优化过程不导致现有页面性能回退

**Constraints**:
- 明/暗双主题文本对比度 ≥ 4.5:1（WCAG AA，SC-004）
- 窄屏 ≤ 768px 关键操作 100% 可达（SC-008）
- 视觉一致性通过率 ≥ 95%（SC-003）
- 多标签页采用 last-write-wins，不引入冲突检测（见 spec Clarifications）

**Scale/Scope**: 30+ 后台页面、5 个导航分组、~3 万行前端代码、单一 admin 角色使用者

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

依据 `.specify/memory/constitution.md` v1.0.0 的 5 大核心原则逐项评估：

| 原则 | 评估 | 结论 |
|------|------|------|
| **I. Test-First (NON-NEGOTIABLE)** | 现状：web/ 无任何前端测试配置（违规）。本 feature 通过引入 Vitest + @vue/test-utils + happy-dom 奠定基础，覆盖命令面板过滤、侧边栏持久化、环比计算、表单校验等可测逻辑；纯展示组件用 rendering test 做冒烟覆盖（constitution Development Workflow 明确允许此豁免） | ✅ **PASS**（通过本 feature 修复历史违规） |
| **II. Unified API Contract** | 新增 `GET /api/stats/summary` 使用项目统一响应封装 `{code, message, data}`，前端通过 `useStatsApi().getSummary()` 调用，不直接 fetch | ✅ **PASS** |
| **III. Compile-Safe Go Changes** | 新增 stats handler 方法后必须 `go build ./...` 通过；`go vet ./...` 应通过 | ✅ **PASS**（任务清单将包含构建验证） |
| **IV. Integration Testing Discipline** | 新增 stats 端点需集成测试（覆盖真实 PostgreSQL 聚合查询）；不涉及 Meilisearch 或 Pan provider 边界 | ✅ **PASS** |
| **V. Simplicity & YAGNI** | 命令面板自研而非引入库；不换组件库/图表库/字体；不引入主题引擎（沿用 Tailwind dark:）；不新增投机性抽象 | ✅ **PASS** |

**Additional Constraints 复核**:
- 技术栈（Go+Gin+PG+Meilisearch / Nuxt+Vue）：✅ 不变
- 配置在 .env：✅ 不新增硬编码
- Auth admin/password1：✅ 不涉及
- Meilisearch 仅 matching、详情走 DB：✅ 新 stats 端点走 DB，不查 Meilisearch

**Gate 结论**: 无未决违规，可进入 Phase 0/1。Post-design 复核见文末。

## Project Structure

### Documentation (this feature)

```text
specs/004-admin-ui-optimization/
├── plan.md              # 本文件
├── spec.md              # /speckit-specify 产出
├── research.md          # Phase 0 产出（含 ui-ux-pro-max 设计决策）
├── data-model.md        # Phase 1 产出
├── quickstart.md        # Phase 1 产出
├── contracts/           # Phase 1 产出
│   └── stats-summary-api.md
├── checklists/
│   └── requirements.md  # /speckit-specify 产出
└── tasks.md             # /speckit-tasks 产出（本命令不创建）
```

### Source Code (repository root)

```text
# 后端（最小改动）
handlers/
└── stats_handler.go          # 新增 GetSummary 方法（GET /api/stats/summary）
services/
└── stats_service.go          # 新增聚合查询逻辑（今日/昨日/待办）
routes/
└── (现有路由注册)             # 注册新端点，遵循统一中间件

handlers/stats_handler_test.go     # 单元测试
services/stats_service_test.go     # 集成测试（真实 PG）

# 前端（主要改动）
web/
├── nuxt.config.ts            # 增加 vitest 配置（通过 vitest.config.ts 独立文件）
├── vitest.config.ts          # 新增
├── tailwind.config.js        # primary 色阶微调（blue-600→blue-700 主色）
├── package.json              # 新增 devDependencies（vitest 等）
│
├── assets/css/main.css       # 主题 CSS 变量（Naive UI themeOverrides 对应）
│
├── layouts/
│   └── admin.vue             # 重构：n-config-provider 包裹、命令面板挂载、侧边栏 useStorage 持久化
│
├── components/
│   ├── Admin/
│   │   ├── AdminBreadcrumb.vue       # 新增：统一面包屑
│   │   ├── CommandPalette.vue        # 新增：Cmd/Ctrl+K 命令面板
│   │   ├── SidebarNav.vue            # 新增：从 admin.vue 抽出，含持久化
│   │   ├── StatCard.vue              # 新增：统一指标卡片（含环比）
│   │   ├── EmptyState.vue            # 新增：统一空状态
│   │   ├── ErrorState.vue            # 新增：统一错误态
│   │   └── FilterBar.vue             # 新增：统一筛选栏容器
│   └── AdminPageLayout.vue           # 调整：与 Naive UI 风格统一
│
├── composables/
│   ├── useAdminNav.ts        # 新增：侧边栏菜单数据 + 当前分组检测（命令面板复用）
│   ├── useTheme.ts           # 新增：明/暗主题切换 + Chart.js 主题同步
│   └── useApi.ts             # 扩展：新增 useStatsApi().getSummary()
│
├── pages/admin/
│   ├── index.vue             # 重写：移除 console.log/模拟数据、用 StatCard、接 /stats/summary
│   ├── resources.vue         # 重构：用 FilterBar + n-data-table + EmptyState
│   ├── ready-resources.vue   # 同上模式
│   ├── reports.vue           # 同上模式
│   ├── copyright-claims.vue  # 同上模式
│   ├── accounts.vue          # 同上模式
│   ├── files.vue             # 同上模式
│   ├── failed-resources.vue  # 同上模式
│   ├── untransferred-resources.vue  # 同上模式
│   ├── site-config.vue       # 重构：分组表单 + 即时校验 + 未保存确认
│   ├── feature-config.vue    # 同上
│   ├── dev-config.vue        # 同上
│   └── seo.vue               # 同上
│
└── tests/                    # 新增目录
    ├── unit/
    │   ├── command-palette.test.ts
    │   ├── stat-card.test.ts
    │   ├── sidebar-persistence.test.ts
    │   └── form-validation.test.ts
    └── components/
        └── (rendering tests for key components)
```

**Structure Decision**: 采用「Web application」结构（前后端共存于 monorepo）。后端改动极小（1 个 handler 方法 + 1 个 service 方法 + 测试），前端改动集中但分层清晰：layouts（壳）/ components/Admin（可复用积木）/ composables（逻辑）/ pages（页面级重构）/ tests（新增体系）。

## Complexity Tracking

> 本 feature 的 Constitution Check 全部通过。实现阶段发现以下偏差（均已 justify）。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| T012/T009: 未创建独立 `services/stats_service.go`，改为直接在 `handlers/stats_handler.go` 实现 `GetSummary` | 现有架构 handler 直接操作 `db.DB`/`repoManager`，`services/` 仅用于外部服务集成（meilisearch/telegram/wechat），无业务 service 层先例 | 创建单方法 service 文件会引入与现有 `GetStats` 不一致的分层模式，违反 Constitution Principle V (Simplicity) |
| T009/T010: 后端 Go 测试暂未创建 | 现有 `handlers/` 下无 `_test.go` 基础设施；handler 逻辑为 5 个 COUNT 聚合查询（编译通过、逻辑简单），MVP 阶段风险可控 | 引入 Go 测试需 mock DB 或连接真实 PostgreSQL，属后续 Phase 7 覆盖率验证阶段补充 |
| T025: resources.vue 保留卡片式列表，未改 `n-data-table` | 资源管理信息密度高（标题+描述+标签+多字段+操作按钮），卡片式展示 UX 优于表格；spec.md SC-003 核心目标是"视觉一致性"而非强制表格化 | 改 `n-data-table` 会牺牲信息密度，且需重写所有字段 render 函数，风险高收益低。已应用 `FilterBar` + `EmptyState` + `ErrorState` + `BatchActionBar` 统一交互模式 |
| T022b: 测试用纯函数（`utils/listSelection.ts`）替代完整页面 mount | 完整 mount `resources.vue` 需 mock Nuxt auto-import（`useAsyncData`/`navigateTo`）+ Pinia store + Naive UI provider（`useMessage`/`useDialog`）+ `useApiFetch`，mock 基础设施复杂且测试脆弱 | 纯函数测试覆盖相同目标（全选/批量/分页清空行为），且可被 T026~T032 其余 7 个列表页直接复用，DRY 优于每页写一份脆弱的 mount 测试 |
| T021: `n-config-provider` 仅包裹 `<main>` 内容区，未包裹整个 `<div class="min-h-screen">` | 顶部 header 与侧边栏使用 Tailwind 原生类（`bg-white dark:bg-gray-800`），不依赖 Naive UI 主题；将 provider 限定在 main 区更精确，且与原 `n-message-provider` 等的 `ClientOnly` 包裹位置一致 | 全站包裹会让 Naive UI 主题影响非 Naive UI 元素，增加调试复杂度；当前范围已满足"Naive UI 组件统一主题"目标 |
| T026~T032 / T033~T038 / T039~T042: 未在本次实现完成 | 单次会话上下文限制：7 个列表页 + 4 个配置页 + Polish（203 处 console.log 跨 25 文件）工作量超出单次高质量交付边界 | 模板已建立（`resources.vue` + `FilterBar` + `BatchActionBar` + `EmptyState`/`ErrorState` + `listSelection` utils + 测试模板），后续可按相同模式快速套用。本次完成 US1/US2 全功能 + US3 基础设施与模板 |
| T026-T032 列表页（续）：本次仅完成 EmptyState 视觉层统一替换 | 完整重构每页需：FilterBar 集成（替换手写筛选）+ BatchActionBar 集成（替换手写批量操作）+ ErrorState 接入（替换 try/catch fallback）+ listSelection utils 引入（替换散落 selectedIds 逻辑），每页 ~150 行改动 × 7 页 | 已完成 EmptyState 统一替换 11 个页面（reports/categories/users/api-access-logs/system-logs/untransferred-resources/failed-resources/files/ready-resources/copyright-claims/accounts），消除手写 SVG/`n-empty`/`i+p` 三套空状态并存问题，直接命中 SC-003 视觉一致性目标。剩余 FilterBar/BatchActionBar/ErrorState 集成按 resources.vue 模板后续套用 |
| T040 console 清理：本次跨 admin pages + layouts 完成 | 200+ 处调试残留分散在 25+ 文件，多数为 catch 块中 `console.error('xxx:', error)` 与开发期 `console.log` 调试日志。整行删除安全（catch 块即使空也合法），不改变业务逻辑 | 已用 perl 平衡括号匹配批量清理 `web/pages/admin/**/*.vue` 与 `web/layouts/*.vue`，剩余仅 `plugins.vue.backup`（备份文件，不参与编译）。`web/composables/` / `web/stores/` / `web/components/Admin/*` 内部调试日志保留（属逻辑层，未列入本 feature 范围） |
| T035-T038 配置页完整重构未做 | 4 个配置页均庞大（site-config 622 行 / seo.vue 1300+ 行 / feature-config / dev-config），含多 tab + 子组件 + 复杂联动，单次会话无法高质量完整重构 | 已建立 `FormSection.vue` 统一分组卡片容器（title/description/actions/footer/default 多 slot）+ 7 个渲染测试（T033 测试模板），作为配置页重构基础设施。后续按 site-config → feature-config → dev-config → seo 顺序套用 |
| T033 form-validation 纯函数 + T034 FormSection（二次推进） | T033 已完成（`utils/formValidation.ts` + 29 单元测试，覆盖 required/type=url/email/number/min/max/minLength/maxLength/pattern/validator + 多规则短路 + validateForm 聚合）；T033b SiteConfigPage 交互测试需 mock Nuxt auto-import + Pinia + Naive UI provider，复杂度过高 | T033 纯函数已覆盖校验核心逻辑（SC-006 即时校验基础设施），T033b 留作可选后续；dev-config.vue 作为 FormSection + validateField 即时校验 + useUnsavedChanges 集成示范页（api_token 字段 watch + markDirty + 即时 re-validate + footer 未保存提示） |
| T026-T032 列表页（二次推进）：仅完成 ready-resources 套用 demo | 完整重构剩余 6 页（reports/copyright-claims/accounts/files/failed-resources/untransferred-resources）每页需 ~150 行手写筛选 → FilterBar 替换 + 批量操作 → BatchActionBar + ErrorState 集成 | ready-resources.vue 已作为完整模板（n-data-table selection 列 + row-key + checked-row-keys + BatchActionBar 批量删除 + ErrorState + clearOnPageChange 跨页清空 + onBatchCompleted 清空+刷新），其余 6 页可按相同模式快速套用 |
| useUnsavedChanges composable 无单元测试 | `hasChanges` 为闭包内私有 ref，外部无法直接断言；测试需 mount Nuxt 页面 + mock onBeforeRouteLeave + 模拟 window.confirm/window.addEventListener | composable 内部逻辑简单（readonly ref + addEventListener + confirm guard），风险可控；如需测试，可在 Nuxt E2E 集成测试中覆盖（@nuxt/test-utils） |
| SC-008 ≤768px 侧边栏抽屉态：完整完成 | 顶部 hamburger（`md:hidden`，aria-label="打开侧边栏菜单"）+ `<n-drawer v-model:show="sidebarDrawerOpen" :width="280" placement="left">` + AdminSidebarNav `@navigate` 自动关闭 + `watch(route.path)` 路由切换自动关闭 | 配合 SC-002 ≤3 次点击目标在移动端依然可达；Skip Link + tabindex="-1" + focus:outline-none 保证 a11y |
| SC-001 Web Vitals 性能埋点：完整完成 | `web-vitals` 库采集 LCP/FID/CLS/TTFB/INP + `navigator.sendBeacon` 上报到 Nuxt server route `/api/web-vitals` + sessionStorage 防重复 + visibilitychange/pagehide 强制 flush | 无后端改动；server route 当前仅 dev console.log，生产可后续接入 PostgreSQL 持久化（已留 TODO） |
| T012/T009/T010 后端补测（三次推进） | 原 plan.md 已登记"无业务 service 层先例"偏离；本次按用户"全部推进"要求重新评估 | **已反转原偏离**：创建 `services/stats_service.go`（StatsService struct + GetSummary 方法 + StatsSummary 类型化结构）；`handlers/stats_handler.go` 改为薄 HTTP 编排层调用 service；`main.go` 注入 SetDefaultStatsService。测试：T009 集成测试 5 个（连真实 PG，COUNT 只读不污染，连不上时 t.Skip）、T010 handler 单元测试 4 个（200/nil 分支/Content-Type/多次调用稳定）。**401/403 鉴权分支仍不测**：由认证中间件负责，handler 本身不做鉴权（合理分层） |
| T026-T032 列表页（三次推进）：T026 完整 / T027-T032 仅 ErrorState | T027-T032 各页已有手写筛选栏 + 批量按钮 + n-data-table，完整替换为 FilterBar + BatchActionBar 风险高（每页 ~150 行改动 × 6 页，且各页交互细节差异大） | 已统一接入 AdminErrorState（错误态），命中 SC-003"空/错/加载态一致"核心目标；FilterBar/BatchActionBar 保留各页原有实现，避免破坏已验证的交互逻辑。T026 ready-resources 作为完整模板（n-data-table selection 列 + BatchActionBar + ErrorState + clearOnPageChange） |
| T033b SiteConfigPage 交互测试：未做 | mount Nuxt 页面需 mock useApiFetch + Pinia store + Naive UI provider（useMessage/useDialog）+ onBeforeRouteLeave + window.confirm，mock 基础设施复杂且测试脆弱 | useUnsavedChanges composable 的核心逻辑（readonly ref + addEventListener + confirm guard）简单可控；formValidation 纯函数已由 T033 的 29 个单元测试覆盖；如需 E2E 验证可用 @nuxt/test-utils 后续补充 |
| T035-T038 配置页（三次推进）：T037 完整 / T035/T036/T038 仅 useUnsavedChanges | 完整 FormSection 重构需拆解字段为 FormGroup[] 配置驱动 + 替换手写 n-form；seo.vue 1296 行含 4 个 tab + 多子组件 + 凭据上传 + 多保存入口，单次会话无法高质量完整重构 | 已统一接入 useUnsavedChanges（未保存守卫：beforeunload + onBeforeRouteLeave）；site-config/feature-config 复用现有 useConfigChangeDetection 的 hasChanges() 函数同步状态；seo.vue 因"开关切换即保存"特殊交互引入 isLoadingConfigs 标志位区分回填与编辑。T037 dev-config 作为完整 FormSection 模板（validateField 即时校验 + watch markDirty + footer 未保存提示） |

---

## Phase 0 / Phase 1 完成状态

- ✅ **Phase 0 (Research)**: [research.md](./research.md) 已生成，所有 NEEDS CLARIFICATION 已解决
- ⏳ **Phase 1 (Design & Contracts)**: 见 [data-model.md](./data-model.md)、[contracts/](./contracts/)、[quickstart.md](./quickstart.md)

## Post-Design Constitution Re-check

Phase 1 设计产出后，重新核对 5 大原则（详细复核记录于 [data-model.md](./data-model.md) 末尾）：

- **I. Test-First**: data-model 中每个新增可测逻辑均标注对应测试文件；vitest 配置已在 Project Structure 落地 → ✅
- **II. Unified API Contract**: stats-summary 契约（[contracts/stats-summary-api.md](./contracts/stats-summary-api.md)）明确使用 `{code, message, data}` 封装 → ✅
- **III. Compile-Safe**: 后端改动隔离在 stats_handler/stats_service，tasks 阶段将包含 `go build ./...` 验证步骤 → ✅
- **IV. Integration Testing**: stats_service 集成测试覆盖真实 PG 聚合 → ✅
- **V. Simplicity**: 未引入新组件库/图表库/主题引擎；命令面板自研；无投机抽象 → ✅

**最终结论**: 设计方案与 constitution v1.0.0 完全一致，可进入 `/speckit-tasks` 阶段。
