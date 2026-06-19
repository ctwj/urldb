# Research: 后台管理系统 UI/UX 整体优化

**Branch**: `004-admin-ui-optimization` | **Date**: 2026-06-19
**Spec**: [spec.md](./spec.md)

> 本文件汇总 Phase 0 研究结论，解决 plan.md Technical Context 中所有 NEEDS CLARIFICATION，并记录由 `/ui-ux-pro-max` 产出的设计选型。

---

## 1. 设计 Pattern 与风格定位

**Decision**: 采用 **Data-Dense Dashboard** 风格（来自 ui-ux-pro-max `--design-system` 推荐）。

**Rationale**:
- 后台目标用户为少量内部管理员，需要"专业、高效、低干扰"，而非营销吸引力
- 现有 30+ 页面以数据表、KPI、筛选栏为主，数据密集型 pattern 与场景天然匹配
- ui-ux-pro-max 标注该风格 Performance ⚡ Excellent、Accessibility ✓ WCAG AA，与 constitution Principle I/IV 兼容

**Alternatives considered**:
- *Glassmorphism / Neumorphism*: 视觉花哨但对比度难保证，违反 SC-004（≥4.5:1），拒绝
- *Minimalist Landing*: 适合 C 页面而非内部后台，信息密度不足

**Key Effects**（来自推荐）: hover tooltips、行 hover 高亮、平滑筛选动画、加载 spinner、图表 hover zoom。

---

## 2. 配色方案（明暗双主题）

**Decision**: 沿用 Tailwind 默认 `slate` + `blue` 色系，主色从 `blue-600` 调整为 `blue-700`（明色）/ `blue-400`（暗色）以提升对比度。

**Rationale**:
- 项目 tailwind.config.js 已定义 `darkMode: 'class'` 与 primary blue 色阶，复用成本最低（YAGNI）
- ui-ux-pro-max 推荐 Primary `#1E40AF`（即 blue-800）+ Secondary `#3B82F6`（blue-500）+ CTA `#F59E0B`（amber-500），与现有蓝色主色一致
- 替换首页"蓝紫渐变"为"纯蓝主色 + 中性灰底"，消除视觉冲突（FR-003）

**完整色板（明色 / 暗色）**:

| 角色 | 明色 (Light) | 暗色 (Dark) | 对比度验证 |
|------|-------------|-------------|-----------|
| Background | `slate-50` #F8FAFC | `slate-900` #0F172A | — |
| Surface (卡片) | `white` #FFFFFF | `slate-800` #1E293B | — |
| Border | `slate-200` #E2E8F0 | `slate-700` #334155 | ✓ 可见 |
| Text Primary | `slate-900` #0F172A | `slate-50` #F8FAFC | ✓ 15:1 / 15:1 |
| Text Secondary | `slate-600` #475569 | `slate-400` #94A3B8 | ✓ 7:1 / 4.6:1 |
| Primary | `blue-700` #1D4ED8 | `blue-400` #60A5FA | ✓ 用于按钮/链接/焦点环 |
| Success | `emerald-600` #059669 | `emerald-400` #34D399 | ✓ |
| Warning | `amber-600` #D97706 | `amber-400` #FBBF24 | ✓ |
| Danger | `rose-600` #E11D48 | `rose-400` #FB7185 | ✓ |
| Info | `sky-600` #0284C7 | `sky-400` #38BDF8 | ✓ |

**Alternatives considered**:
- *引入全新色系（如 indigo/violet）*: 与系统既有蓝色品牌冲突，且 SC-003 视觉一致性要求统一，拒绝
- *自定义 CSS 变量主题引擎*: 违反 YAGNI，Tailwind `dark:` 变体已足够

---

## 3. 字体配对

**Decision**: 保留 **Inter**（正文/UI）+ **Fira Code**（数字、统计值、代码片段）。

**Rationale**:
- 项目 `nuxt.config.ts` 已通过 `vfonts/Lato.css` + `vfonts/FiraCode.css` 引入，tailwind.config.js 已设 `fontFamily.sans: ['Inter', ...]`
- ui-ux-pro-max 推荐 Fira Sans + Fira Code（mood: dashboard/data/technical/precise），Fira Code 已在用
- Inter 是 dashboard/后台事实标准（GitHub、Linear、Vercel 都用），迁移成本为零
- 数字使用等宽字体（Fira Code）能避免 KPI 卡片数值跳动，提升"数据感"

**层级规范**:
- H1（页面标题）: Inter Semi Bold 24px / line-height 1.3
- H2（区块标题）: Inter Semi Bold 18px / line-height 1.4
- H3（卡片标题）: Inter Medium 16px / line-height 1.5
- Body: Inter Regular 14px / line-height 1.6（后台默认 14px，符合数据密集型）
- Caption/Muted: Inter Regular 12px / line-height 1.5
- Numeric/Code: Fira Code Regular 13px

**Alternatives considered**:
- *切换到 ui-ux-pro-max 推荐的 Fira Sans*: 收益边际，迁移成本（更换字体加载、回归测试）不划算，拒绝
- *引入衬线字体（如 Source Serif）*: 不符合"technical/precise" mood，拒绝

---

## 4. 组件风格规范

**Decision**: **统一以 Naive UI 组件为单一来源**，移除原生 `bg-white dark:bg-gray-800 rounded-lg shadow` 自定义卡片 div；通过 Naive UI 的 theme overrides 注入上述配色。

**Rationale**:
- 项目已通过 `unplugin-vue-components` + `NaiveUiResolver` 自动导入 Naive UI，且 `nuxt.config.ts` build.transpile 已配置 SSR 支持
- 当前"三套卡片样式并存"（n-card / 自定义 div / AdminPageLayout 插槽）正是 FR-001 要解决的核心问题
- Naive UI 原生支持 `darkTheme`，通过 `n-config-provider` + `themeOverrides` 可一次性统一全站，满足 SC-003（≥95% 一致性）

**统一规范**:

| 元素 | 规范 |
|------|------|
| 卡片 | `n-card` + `:bordered="false"` + 默认 padding；移除所有自定义 shadow div |
| 按钮 | `n-button`，type 语义化（primary/info/success/warning/error） |
| 输入/选择 | `n-input` / `n-select` / `n-checkbox`，统一 size="medium" |
| 筛选栏 | 统一容器：`n-card` 内 `grid grid-cols-1 md:grid-cols-{N} gap-4` |
| 表格 | 优先 `n-data-table`（内置排序/分页/选择），替代手写 v-for 列表 |
| 模态/抽屉 | `n-modal` / `n-drawer`，统一关闭行为与 esc 支持 |
| 反馈 | 已有 `useToast.ts`，统一用 `useMessage()` / `useNotification()` |
| 加载态 | 统一 `n-spin`；首屏用 `n-skeleton` 避免 layout shift |
| 空状态 | 统一 `<EmptyState>` 组件（图标 + 标题 + 说明 + 可选 CTA） |
| 错误态 | 统一 `<ErrorState>` 组件（图标 + 错误信息 + 重试按钮） |

**圆角与间距**:
- 圆角统一 `rounded-lg`（8px）；按钮/输入 `rounded-md`（6px）
- 间距基于 4px 网格：`gap-2`(8) / `gap-4`(16) / `gap-6`(24) / `gap-8`(32)
- 卡片间距：区块间 `space-y-6`，卡片内 padding `p-6`

**Alternatives considered**:
- *引入 shadcn-style 组件库*: 与现有 Naive UI 重复，违反 YAGNI，拒绝
- *纯 Tailwind 自研组件*: 失去 Naive UI 已有的无障碍/键盘支持，违反 FR-021

---

## 5. 命令面板（FR-007）

**Decision**: **自研轻量命令面板**（非引入第三方库），通过 `Cmd/Ctrl+K` 唤起，基于现有侧边栏菜单配置生成索引。

**Rationale**:
- 需求范围明确且窄：仅页面跳转（30+ 页面），无需 fuzzy 文件搜索、命令执行等重型功能
- 现有 `layouts/admin.vue` 已维护完整的 5 个分组菜单数据结构（dashboardItems / dataManagementItems / systemConfigItems / operationItems / statisticsItems），可直接复用为索引源
- 自研约 200 行 Vue 代码（modal + 输入 + 过滤 + 键盘导航），少于引入 `vue-command-palette` 后的定制成本
- ui-ux-pro-max UX 规则要求：键盘导航 Tab 顺序匹配视觉、可见 focus ring、Skip Links——自研可完全掌控

**UI 设计**:
- 顶部居中 modal（max-w-xl），距顶部 25vh
- 输入框 + 结果列表（显示页面标题、所属分组、图标）
- 中英文模糊匹配（标题/拼音首字母）
- 键盘：↑↓ 选择、Enter 跳转、Esc 关闭
- 可发现性：顶部状态栏显示 `⌘K 搜索` 提示徽章

**Alternatives considered**:
- *vue-command-palette / @vueuse/core 的 useMagicKeys*: 引入依赖但定制需求少，性价比低
- *顶部固定搜索框*: 占用状态栏空间，与现有"自动处理/转存状态"指示冲突，拒绝

---

## 6. 图表样式

**Decision**: **保留 Chart.js**，统一主色与暗色适配，移除 `admin/index.vue` 中的 console.log 与硬编码模拟 fallback。

**Rationale**:
- 项目已用 Chart.js 4.5.0，ui-ux-pro-max chart 推荐首项即 Chart.js（趋势图场景）
- 当前问题不在库选型，而在：调试残留、模拟数据、暗色网格未适配、主色与系统不一致

**统一规范**:
- 折线图主色：明色 `blue-700` #1D4ED8 / 暗色 `blue-400` #60A5FA
- 多系列区分：blue / emerald / amber（与状态色一致）
- 填充：20% opacity（ui-ux-pro-max 推荐）
- 网格线：明色 `rgba(0,0,0,0.08)` / 暗色 `rgba(255,255,255,0.08)`
- 字体：Chart.js 全局 `Chart.defaults.font.family = 'Inter, sans-serif'`
- 空数据态：检测 dataset 全 0 时销毁图表并显示 `<EmptyState>`，不再用硬编码模拟数据

**Alternatives considered**:
- *迁移到 ApexCharts / Recharts*: 无功能收益，迁移成本高，违反 YAGNI

---

## 7. 前端测试基础设施（Constitution Principle I 关键）

**Decision**: **为本 feature 引入最小前端测试基础设施**：Vitest + @vue/test-utils + happy-dom，仅覆盖新增的可单元化逻辑（命令面板过滤、侧边栏状态持久化、指标卡片环比计算、表单校验）。

**Rationale / Constitution Gate 处理**:
- Constitution Principle I (Test-First NON-NEGOTIABLE) 要求所有生产代码先有测试，覆盖率 ≥ 90%
- 现状：web/ 目录无任何测试配置（无 vitest.config / playwright.config，package.json 无测试依赖）→ **现状即违规**
- 本 feature 主要是 UI 重构，纯展示组件 TDD 不切实际（constitution Development Workflow 已豁免"presentational components where TDD is impractical"，但要求"minimum a rendering or interaction test MUST accompany the change"）
- 折中：引入 Vitest 覆盖**可测试逻辑**（命令面板过滤算法、侧边栏状态序列化、环比计算、API mock），UI 视觉部分用 component rendering test 做冒烟测试

**引入依赖**（devDependencies）:
- `vitest` ^1.x
- `@vue/test-utils` ^2.x
- `happy-dom` ^13.x（轻量 DOM 环境，快于 jsdom）
- `@nuxt/test-utils` ^3.x（Nuxt 集成测试场景）

**覆盖率目标**: 新增逻辑代码 ≥ 90%（constitution gate）；UI 组件 rendering test 覆盖关键交互（命令面板、侧边栏、指标卡点击、批量操作确认）。

**Alternatives considered**:
- *跳过前端测试*: 直接违反 constitution NON-NEGOTIABLE 原则，不可行
- *Playwright E2E 全覆盖*: 成本高、CI 慢，留待后续迭代，本 feature 先用 Vitest 奠定基础
- *仅靠手动验证*: 违反 FR-013~FR-020 的可验收性，且无法防止回归

---

## 8. API 端点策略（Constitution Principle II）

**Decision**: **新增 1 个统一封装的后端统计端点** `GET /api/stats/summary`，返回仪表盘所需聚合数据（今日 + 昨日的资源/浏览/搜索 + 待办聚合），前端不再发多次请求拼装。

**Rationale**:
- FR-009 要求首屏一次性展示核心指标 + 环比昨日 + 待办聚合；当前 `admin/index.vue` 发了 3 次请求（stats / views-trend / searches-trend），且趋势端点返回的是 7 日序列，不适合"环比昨日"的卡片对比
- 新端点遵循 constitution Principle II（Unified API Contract）：使用项目统一响应封装 `{code, message, data}`，前端通过 `useStatsApi().getSummary()` 调用
- 后端实现位置：复用现有 `handlers/stats_handler.go`（若存在）或新建，编译遵循 Principle III

**响应契约（草案，详见 contracts/）**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "resources": { "today": 12, "yesterday": 8, "total": 1234 },
    "views": { "today": 450, "yesterday": 380 },
    "searches": { "today": 89, "yesterday": 72 },
    "todos": {
      "ready_resources": 3,
      "failed_tasks": 1,
      "pending_reports": 2
    }
  }
}
```

**Alternatives considered**:
- *前端组合多个现有端点*: 多次 RTT、无法原子读取、环比计算散落前端，违反 FR-009 首屏 5 秒目标
- *新增多个细分端点*: 过度拆分，YAGNI 违规

---

## 9. 侧边栏状态持久化（FR-006）

**Decision**: 用 `@vueuse/core` 的 `useStorage`（项目已依赖 `@vueuse/core` ^14.0.0）持久化用户展开的分组状态。

**Rationale**:
- 项目已引入 `@vueuse/core`，无需新增依赖
- `useStorage` 自动处理 JSON 序列化、跨标签页同步（localStorage 事件）——但本 feature 已决定多标签页采用 last-write-wins，跨标签页同步属于副作用但不影响正确性
- key 规范：`urldb:admin:sidebar:expanded-groups` → `{ dataManagement: true, systemConfig: false, ... }`

**行为**:
- 进入页面时：自动展开当前路由所属分组（覆盖存储值）+ 读取用户其他分组偏好
- 用户手动 toggle 分组：写入 storage
- 重置入口：用户菜单可加"重置布局"（可选，P2+）

---

## 10. 风险与未决项（Deferred to tasks 阶段）

- **图表暗色主题切换时机**: Chart.js 实例不会自动响应 Tailwind `dark:` class 切换，需监听主题变化并重建实例——实现细节留待 tasks 阶段
- **拼音匹配**: 命令面板的中文拼音首字母匹配需引入 `pinyin-pro` 或简化为"包含子串即匹配"——tasks 阶段决定
- **现有 30+ 页面重构顺序**: 应按 P1（仪表盘）→ P2（侧边栏/导航/列表页统一）→ P3（配置页表单）分批，避免大爆炸式 PR——tasks 阶段拆分

---

## 总结：所有 NEEDS CLARIFICATION 已解决

| 原未知项 | 决策 | 来源 |
|---------|------|------|
| 配色方案 | Tailwind slate+blue，主色 blue-700/400 | ui-ux-pro-max + 现状 |
| 字体配对 | Inter + Fira Code（保留） | ui-ux-pro-max + 现状 |
| 组件库 | Naive UI 单一来源 + theme overrides | 现状 + FR-001 |
| 命令面板 | 自研轻量（Cmd/Ctrl+K） | YAGNI |
| 图表 | 保留 Chart.js，统一主色 | ui-ux-pro-max + 现状 |
| 前端测试 | 引入 Vitest + @vue/test-utils | constitution Principle I |
| 统计端点 | 新增 GET /api/stats/summary | FR-009 + Principle II |
| 状态持久化 | @vueuse/core useStorage | 现状依赖 |
