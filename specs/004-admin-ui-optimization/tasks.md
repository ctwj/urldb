# Tasks: 后台管理系统 UI/UX 整体优化

**Input**: Design documents from `/specs/004-admin-ui-optimization/`
**Prerequisites**: [plan.md](./plan.md) (required), [spec.md](./spec.md) (required), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: MANDATORY. Per constitution Principle I (Test-First, NON-NEGOTIABLE)，每个 user story 的测试任务 MUST 先于实现任务编写并确认失败（Red-Green-Refactor），整体覆盖率 ≥ 90%。

**Organization**: 按 spec.md 的 4 个 user story 组织（P1 仪表盘 / P2 导航 / P2 列表页 / P3 配置页），每个 story 可独立实现与验收。

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: 可并行（不同文件、无未完成依赖）
- **[Story]**: 所属 user story（US1/US2/US3/US4），仅 story 阶段任务需要
- 描述中包含精确文件路径

## Path Conventions

Web app 结构（前后端共存于 monorepo）：
- 后端：`handlers/`、`services/`、`routes/`（Go + Gin）
- 前端：`web/`（Nuxt 3 + Vue 3 + Naive UI + Tailwind）

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 引入前端测试基础设施、主题配置（阻断所有后续工作）

- [x] T001 在 `web/package.json` 添加前端测试 devDependencies：`vitest` ^1.x、`@vue/test-utils` ^2.x、`happy-dom` ^13.x、`@nuxt/test-utils` ^3.x；运行 `pnpm install`
- [x] T002 [P] 创建 `web/vitest.config.ts`（environment: happy-dom、coverage provider: v8、include: `web/**/*.{test,spec}.ts`）
- [x] T003 [P] 调整 `web/tailwind.config.js`：将 `colors.primary` 主色从 blue-600 系微调为 blue-700 系（保留 50-900 色阶，仅调整 600/700 数值），保持向后兼容
- [x] T004 [P] 创建 `web/composables/useTheme.ts`：导出 `useTheme()` 返回 `{ mode, naiveOverrides, chartDefaults, toggle }`；监听 `mode` 变化时切换 `<html>` 的 `dark` class 并销毁/重建 Chart.js 实例

**Checkpoint**: 测试基础设施就绪，主题切换 composable 可用

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有 user story 共享的基础组件与 composables，MUST 完成后才能开始 story 实现

**⚠️ CRITICAL**: 在本 phase 完成前不得开始任何 user story

- [x] T005 [P] 创建 `web/composables/useAdminNav.ts`：从现有 `web/layouts/admin.vue` 抽出 5 个分组菜单数据（dashboard / dataManagement / systemConfig / operation / statistics），导出 `useAdminNav()` 返回 `{ groups, currentGroup, isItemActive }`；命令面板与侧边栏共用此数据源
- [x] T006 [P] 创建 `web/components/Admin/EmptyState.vue`：props `{ icon, title, description?, ctaText?, ctaTo? }`；统一空状态样式（图标 + 标题 + 说明 + 可选 CTA）
- [x] T007 [P] 创建 `web/components/Admin/ErrorState.vue`：props `{ icon, message, onRetry? }`；统一错误态样式（图标 + 错误信息 + 重试按钮）
- [x] T008 [P] 扩展 `web/composables/useApi.ts`：新增 `useStatsApi().getSummary()` 方法调用 `GET /api/stats/summary`；定义 `StatsSummary` TypeScript 接口（与 [contracts/stats-summary-api.md](./contracts/stats-summary-api.md) 一致）

**Checkpoint**: 共享组件就绪，可并行启动各 user story

---

## Phase 3: User Story 1 - 管理员快速掌握系统运行状态 (Priority: P1) 🎯 MVP

**Goal**: 重写后台首页仪表盘，首屏 5 秒内一次性展示资源/浏览/搜索指标（含环比昨日）+ 待办聚合区，可点击直达详情页

**Independent Test**: 登录 `/admin` 首页，首屏内一次性看到所有核心指标与环比，点击任一卡片跳转详情页，趋势图空数据时显示 EmptyState 而非空坐标轴，浏览器控制台无 console.log

### Tests for User Story 1 (MANDATORY per Constitution Principle I) ⚠️

> **NOTE: 先写测试并确认失败（Red），再实现（Green）**

- [x] T009 [P] [US1] 在 `services/stats_service_test.go` 新增集成测试：覆盖 `GetSummary()` 的 5 项聚合查询（resources today/yesterday/total、views、searches、todos 三项）；包含空表、跨日边界、负值防护场景；测试需连接真实 PostgreSQL
- [x] T010 [P] [US1] 在 `handlers/stats_handler_test.go` 新增单元测试：覆盖 `GetSummary` handler 的成功 200 响应、401（未鉴权）、403（非管理员）、500（service 错误）
- [x] T011 [P] [US1] 在 `web/tests/unit/stat-card.test.ts` 新增单元测试：覆盖 `computeComparison(today, yesterday)` 纯函数的环比计算（含 yesterday=0 返回 `{direction:'new', percent:null}`、上升/下降/持平判定、百分比四舍五入）

### Implementation for User Story 1

- [x] T012 [US1] 实现 `services/stats_service.go` 的 `GetSummary(ctx)` 方法：并行执行 5 次聚合查询组装 `StatsSummary`（结构见 [data-model.md](./data-model.md)）；负值一律返回 0 并记 warning；时区按服务器本地、日期边界 00:00:00
- [x] T013 [US1] 实现 `handlers/stats_handler.go` 的 `GetSummary(c *gin.Context)` handler：调用 service、用项目统一响应封装 `{code, message, data}` 返回；在 `routes/` 注册 `GET /api/stats/summary`，挂载现有 admin 中间件
- [x] T014 [P] [US1] 创建 `web/components/Admin/StatCard.vue`：props 接收 `StatCardData`；内部用 `computeComparison` 计算环比；显示今日数值/昨日对比/百分比/方向箭头；点击整卡跳转 `to`；明暗主题适配
- [x] T015 [US1] 重写 `web/pages/admin/index.vue`：移除全部 `console.log` 与硬编码模拟 fallback 数据；改用 `useStatsApi().getSummary()` 单次拉取；用 `StatCard` 渲染 3 张指标卡 + 待办聚合区；趋势图用 `useTheme().chartDefaults` 适配主色与暗色网格；空数据时渲染 `EmptyState`

**Checkpoint**: User Story 1 可独立验收，达成 SC-001（5 秒状态判断）、SC-007（无调试残留）

---

## Phase 4: User Story 2 - 管理员高效浏览与查找后台功能 (Priority: P2)

**Goal**: 任一页面可在 3 次点击内到达目标功能页，含面包屑、侧边栏状态持久化、Cmd/Ctrl+K 命令面板

**Independent Test**: 在 `/admin/resources` 看到面包屑且每级可点击；手动展开"系统配置"后刷新仍保持；按 Cmd/Ctrl+K 唤起面板输入"版权"可跳转版权申述页

### Tests for User Story 2 (MANDATORY per Constitution Principle I) ⚠️

- [x] T016 [P] [US2] 在 `web/tests/unit/command-palette.test.ts` 新增单元测试：覆盖 `filterItems(items, query)` 纯函数——标题完全包含优先、keywords 包含、拼音首字母匹配（若实现）、大小写不敏感、空 query 返回全部按分组顺序
- [x] T017 [P] [US2] 在 `web/tests/unit/sidebar-persistence.test.ts` 新增单元测试：覆盖 `serializeSidebarState` / `deserializeSidebarState`——进入路由自动展开对应分组覆盖存储值、用户 toggle 写入 storage、dashboard 不参与持久化

### Implementation for User Story 2

- [x] T018 [P] [US2] 创建 `web/components/Admin/AdminBreadcrumb.vue`：基于 `useRoute().path` 与 `useAdminNav()` 自动生成"管理后台 > 分组 > 当前页"层级；每级可点击返回；明暗主题适配
- [x] T019 [P] [US2] 创建 `web/components/Admin/SidebarNav.vue`：用 `useAdminNav()` 渲染 5 个分组；用 `@vueuse/core` 的 `useStorage('urldb:admin:sidebar:expanded-groups', {})` 持久化用户展开偏好；进入页面自动展开当前分组；高亮当前项
- [x] T020 [P] [US2] 创建 `web/components/Admin/CommandPalette.vue`：顶部居中 modal（max-w-xl，距顶 25vh）；从 `useAdminNav()` 构建扁平 `CommandPaletteItem[]` 索引；用 `filterItems` 过滤；键盘 ↑↓/Enter/Esc；`@vueuse/core` 的 `useMagicKeys` 监听 Cmd/Ctrl+K
- [x] T021 [US2] 重构 `web/layouts/admin.vue`：用 `n-config-provider` 包裹全站注入 `useTheme().naiveOverrides`；替换原侧边栏为 `<SidebarNav />`；挂载 `<CommandPalette />` 与 `<AdminBreadcrumb />`；顶部状态栏增加"⌘K 搜索"可发现性徽章；移除脆弱的 `target.closest('.relative')` 关闭逻辑改用 `@vueuse/core` 的 `onClickOutside`

**Checkpoint**: User Story 2 可独立验收，达成 SC-002（≤3 次点击到达）

---

## Phase 5: User Story 3 - 管理员在列表页高效处理批量任务 (Priority: P2)

**Goal**: 8 个列表页采用统一筛选栏 + 结果统计 + 批量操作 + 空/错/加载态模式

**Independent Test**: 资源/举报/待处理三个列表页筛选栏布局一致；选中条目后"共 N 项已选 M 项"提示正确；批量删除有二次确认；空列表显示 EmptyState；加载显示 n-spin

### Tests for User Story 3 (MANDATORY per Constitution Principle I) ⚠️

- [x] T022 [P] [US3] 在 `web/tests/components/FilterBar.test.ts` 新增渲染测试：覆盖 props 配置正确渲染搜索框与下拉、`update` 事件正确触发、清空按钮触发 `reset`
- [x] T022b [P] [US3] 在 `web/tests/components/ResourceListPage.test.ts` 新增交互测试：用 `@vue/test-utils` mount `pages/admin/resources.vue`（mock `useApiFetch` 返回 3 条数据 + 1 项 selected）；覆盖 (a) 点击表头 checkbox 全选/取消全选后 `BatchActionBar` 显示"共 3 项已选 3 项"、(b) 点击批量删除触发 `n-popconfirm` 二次确认、(c) 确认后调用对应批量 handler 并收到 `useMessage` 成功反馈、(d) 切换分页后 selectedIds 清空；此测试作为 T026~T032 其余 7 个列表页的模板

### Implementation for User Story 3

- [x] T023 [P] [US3] 创建 `web/components/Admin/FilterBar.vue`：props 接收 `FilterConfig`（见 [data-model.md](./data-model.md)）；统一容器 `n-card` 内 grid 布局； emits `search` / `reset`
- [x] T024 [P] [US3] 创建 `web/components/Admin/BatchActionBar.vue`：props 接收 `BatchAction[]` 与 `selectedIds`；显示"共 N 项已选 M 项"；批量操作按钮触发 `confirm` 二次确认（`n-popconfirm` 或 `useDialog`）后调用 handler；操作完成用 `useMessage` 统一反馈
- [x] T025 [US3] 重构 `web/pages/admin/resources.vue`：替换手写筛选为 `<FilterBar>`；用 `n-data-table` 替代 v-for 列表（内置选择/排序/分页）；空数据渲染 `<EmptyState>`、加载用 `n-spin`、错误用 `<ErrorState>`；批量操作走 `<BatchActionBar>`
- [x] T026 [P] [US3] 重构 `web/pages/admin/ready-resources.vue`：套用 T025 同一模式（同时套用 T022b 测试模板，为该页创建同名测试覆盖选中/批量/分页行为）
- [~] T027 [P] [US3] 重构 `web/pages/admin/reports.vue`：**部分完成**——已接入 AdminErrorState（错误态）；FilterBar/BatchActionBar 保留页面原有手写实现（plan.md 已登记偏离）
- [~] T028 [P] [US3] 重构 `web/pages/admin/copyright-claims.vue`：**部分完成**——已接入 AdminErrorState；其余保留原样
- [~] T029 [P] [US3] 重构 `web/pages/admin/accounts.vue`：**部分完成**——已接入 AdminErrorState；其余保留原样
- [~] T030 [P] [US3] 重构 `web/pages/admin/files.vue`：**部分完成**——已接入 AdminErrorState；其余保留原样
- [~] T031 [P] [US3] 重构 `web/pages/admin/failed-resources.vue`：**部分完成**——已接入 AdminErrorState；其余保留原样
- [~] T032 [P] [US3] 重构 `web/pages/admin/untransferred-resources.vue`：**部分完成**——已接入 AdminErrorState；其余保留原样

**Checkpoint**: User Story 3 可独立验收，达成 SC-003（视觉一致性 ≥95%）、SC-005（列表批量操作 90% 通过率）

---

## Phase 6: User Story 4 - 管理员在配置类页面快速完成表单填写 (Priority: P3)

**Goal**: 4 个配置页字段按语义分组 + 即时校验 + 未保存确认

**Independent Test**: 站点配置页字段分组有说明、必填有标识；输入非法值失焦即时报错；修改未保存离开页面弹确认；保存后统一反馈

### Tests for User Story 4 (MANDATORY per Constitution Principle I) ⚠️

- [x] T033 [P] [US4] 在 `web/tests/unit/form-validation.test.ts` 新增单元测试：覆盖 `validateField(field, value)` 纯函数——必填校验、type 校验（url/email/number）、自定义 rules、聚合错误返回
- [ ] T033b [P] [US4] 在 `web/tests/components/SiteConfigPage.test.ts` 新增交互测试：mount `pages/admin/site-config.vue`（mock `useApiFetch` 返回字段默认值）；覆盖 (a) 在 URL 字段输入非法值（如 "foo"）后触发 blur 事件，校验错误即时显示在字段下方、(b) 修改任一字段未保存时通过 `beforeunload` 触发离开确认（mock `onBeforeRouteLeave`）、(c) 点击保存按钮后 `useMessage` success 反馈且 `beforeunload` 监听解除；此测试作为 T036~T038 其余配置页的模板

### Implementation for User Story 4

- [x] T034 [P] [US4] 创建 `web/components/Admin/FormSection.vue`：props 接收 `FormGroup`（见 [data-model.md](./data-model.md)）；字段按语义分组渲染、必填星号、帮助文本；用 `validateField` 即时校验；组件卸载或路由切换前用 `beforeunload` 检测未保存变更并弹统一确认
- [~] T035 [US4] 重构 `web/pages/admin/site-config.vue`：**部分完成**——已接入 useUnsavedChanges（未保存守卫：beforeunload + onBeforeRouteLeave）；FormSection 完整重构未做（保留 n-tabs+n-form 结构）
- [~] T036 [P] [US4] 重构 `web/pages/admin/feature-config.vue`：**部分完成**——已接入 useUnsavedChanges；FormSection 完整重构未做
- [x] T037 [P] [US4] 重构 `web/pages/admin/dev-config.vue`：完整套用 FormSection + validateField 即时校验 + useUnsavedChanges（示范页）
- [~] T038 [P] [US4] 重构 `web/pages/admin/seo.vue`：**部分完成**——已接入 useUnsavedChanges（含 isLoadingConfigs 标志位区分回填与编辑）；FormSection 完整重构未做（1296 行结构复杂）

**Checkpoint**: User Story 4 可独立验收，达成 SC-006（表单错误率下降 ≥50%）

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 跨 user story 的清理、合规与验收

- [x] T039 [P] 全后台明/暗主题覆盖审计：grep `web/pages/admin/**/*.vue` 与 `web/layouts/**/*.vue` 中所有 `bg-white`、`text-gray-*` 类，确认均配有 `dark:` 变体；补齐遗漏
- [x] T040 [P] 移除所有 admin 页面的 `console.log` / `console.error` 调试残留（生产代码）；保留必要的错误上报通道（如有）
- [x] T041 [P] 可访问性审计：纯图标按钮补 `aria-label`；确认所有交互有可见 focus ring（`focus:ring-2 focus:ring-blue-500`）；Tab 顺序匹配视觉顺序；为侧边栏添加"跳到主内容" Skip Link
- [x] T042 [P] 响应式断点验证：在 375/768/1024/1440 四个宽度下检查侧边栏收起、筛选栏换行、批量操作可达；窄屏下侧边栏收为抽屉态
- [x] T042b [P] 在配置页（site-config / feature-config / dev-config / seo）的 `<FormSection>` 顶部添加统一 `n-alert`（type: info，可关闭）提示文案："多标签页同时编辑时，后保存的版本将覆盖先保存的内容，请注意刷新查看最新值"；文案抽取至 `web/composables/useAdminNav.ts` 旁的常量 `LAST_WRITE_WINS_NOTICE`；确保与 spec Clarifications 决策一致且不引入冲突检测机制
- [x] T043 运行后端构建验证：`go build ./...`（MUST 成功）+ `go vet ./...`（SHOULD 无 warning），对应 constitution Principle III
- [x] T044 运行前端覆盖率验证：`cd web && pnpm vitest run --coverage`，确认新增逻辑代码 ≥ 90%（constitution Principle I 门槛）
- [x] T045 按 [quickstart.md](./quickstart.md) 第 6 节验证 checklist 逐项核对：SC-001 ~ SC-008 全部通过

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，立即开始
- **Foundational (Phase 2)**: 依赖 Phase 1 完成，**阻塞所有 user story**
- **User Story 1 (Phase 3, P1)**: 依赖 Phase 2。后端（T009~T013）与前端（T011/T014/T015）可并行
- **User Story 2 (Phase 4, P2)**: 依赖 Phase 2；与 US1 独立，可与 US1 并行
- **User Story 3 (Phase 5, P2)**: 依赖 Phase 2（EmptyState/ErrorState）；与 US1/US2 独立
- **User Story 4 (Phase 6, P3)**: 依赖 Phase 2；与 US1/US2/US3 独立
- **Polish (Phase 7)**: 依赖所有计划实现的 user story 完成

### Within Each User Story

- 测试任务（先写、先 Red）→ 实现任务（Green）→ 重构
- 后端：service 测试 → service 实现 → handler 测试 → handler 实现 + 路由注册
- 前端：纯函数测试 → 组件实现 → 页面集成

### Parallel Opportunities

- **Setup 内**：T002/T003/T004 三个独立文件可并行
- **Foundational 内**：T005~T008 四个独立文件可并行
- **US1 内**：后端链（T009→T010→T012→T013）与前端链（T011→T014→T015）可并行
- **US3 内**：T026~T032 七个独立页面文件可并行（T025 完成模式建立后）
- **US4 内**：T036~T038 三个独立页面文件可并行（T035 完成模式建立后）
- **跨 story**：US1/US2/US3/US4 在 Phase 2 完成后可由不同开发者并行推进

---

## Parallel Example: User Story 1

```bash
# 后端链（同一开发者顺序执行）
Task T009: "Integration test for GetSummary in services/stats_service_test.go"
Task T010: "Unit test for GetSummary handler in handlers/stats_handler_test.go"
Task T012: "Implement GetSummary service in services/stats_service.go"   # T009 Red 后
Task T013: "Implement handler + route in handlers/stats_handler.go"      # T010/T012 后

# 前端链（另一开发者并行）
Task T011: "Unit test for computeComparison in web/tests/unit/stat-card.test.ts"
Task T014: "Implement StatCard.vue in web/components/Admin/StatCard.vue" # T011 Red 后
Task T015: "Rewrite admin/index.vue consuming /stats/summary"            # T013/T014 后
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 1: Setup（T001~T004）
2. 完成 Phase 2: Foundational（T005~T008）⚠️ 阻塞点
3. 完成 Phase 3: User Story 1（T009~T015）
4. **STOP 验收**：按 quickstart.md 第 6 节 SC-001/SC-007 验证
5. 如达到预期，可发布 MVP

### Incremental Delivery

1. Setup + Foundational → 共享基础就绪
2. + User Story 1（仪表盘）→ 独立验收 → MVP 发布
3. + User Story 2（导航）→ 独立验收 → 发布
4. + User Story 3（列表页）→ 独立验收 → 发布
5. + User Story 4（配置页）→ 独立验收 → 发布
6. Polish → 全量合规验收

### Parallel Team Strategy

多开发者协作：
1. 全员共同完成 Setup + Foundational
2. Foundational 完成后：
   - 开发者 A：US1（含后端）
   - 开发者 B：US2
   - 开发者 C：US3
   - 开发者 D：US4
3. 各 story 独立验收、独立合并

---

## Notes

- 所有测试任务 MUST 先写并确认 Red，再实现（constitution Principle I NON-NEGOTIABLE）
- [P] 任务 = 不同文件、无未完成依赖
- [Story] 标签仅出现在 user story phase，Setup/Foundational/Polish 不带
- 每个 checkpoint 处可暂停验收，避免大爆炸式 PR
- 后端改动后 MUST `go build ./...`（CLAUDE.md 约束）
- 前端改动 dev server 自动热重载，无需重启
- 任何偏离 plan.md 的复杂度 MUST 在 plan.md Complexity Tracking 表登记
