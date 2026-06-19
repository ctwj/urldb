# Data Model: 后台管理系统 UI/UX 整体优化

**Branch**: `004-admin-ui-optimization` | **Date**: 2026-06-19
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md) | **Research**: [research.md](./research.md)

> 本 feature 以后端最小改动 + 前端重构为主，核心数据模型分两部分：
> 1. **后端**：新增 1 个聚合统计响应（无新表、不改 schema）
> 2. **前端**：新增若干 UI 状态/配置类型（TypeScript interface）

---

## 后端数据模型

### 不涉及 schema 变更

本 feature **不新增数据库表、不修改现有表结构**。新增的 `GET /api/stats/summary` 端点仅对现有表（resources / reports / copyright_claims / tasks 等）执行聚合查询。

### 响应 Schema: `StatsSummary`

新端点返回的聚合数据结构（Go 结构体将位于 `services/stats_service.go`）：

```go
// StatsSummary 是 GET /api/stats/summary 的响应数据部分
type StatsSummary struct {
    Resources MetricWithTotal `json:"resources"`
    Views     MetricPair      `json:"views"`
    Searches  MetricPair      `json:"searches"`
    Todos     TodoSummary     `json:"todos"`
}

// MetricPair 用于环比对比（今日 vs 昨日）
type MetricPair struct {
    Today     int64 `json:"today"`
    Yesterday int64 `json:"yesterday"`
}

// MetricWithTotal 在环比基础上额外返回总量
type MetricWithTotal struct {
    MetricPair
    Total int64 `json:"total"`
}

// TodoSummary 待办/异常聚合
type TodoSummary struct {
    ReadyResources int `json:"ready_resources"` // 待处理资源数
    FailedTasks    int `json:"failed_tasks"`    // 失败任务数
    PendingReports int `json:"pending_reports"` // 待审核举报数
}
```

**计算规则**（services 层职责）：
- `resources.today`: 当日 0 点至今新增资源数
- `resources.yesterday`: 前一自然日新增资源数
- `resources.total`: 资源表全量计数
- `views.*`: 当日/昨日资源访问日志聚合（复用现有 `statsApi.getStats` 数据源）
- `searches.*`: 当日/昨日搜索日志聚合（复用现有 searches-trend 数据源）
- `todos.ready_resources`: `ready_resources` 表中状态为 pending 的计数
- `todos.failed_tasks`: `tasks` 表中状态为 failed 的计数
- `todos.pending_reports`: `reports` 表中状态为 pending 的计数

**约束**:
- 所有计数必须为非负整数；负值视为查询错误，返回 0 并记日志
- 时区：使用服务器本地时区（与现有 stats 端点一致），日期边界为 00:00:00
- 性能：单次请求 ≤ 500ms（5 次聚合查询，可并行）

---

## 前端数据模型

### 1. 指标卡片：`StatCardData`

```ts
// components/Admin/StatCard.vue 的 props 类型
interface StatCardData {
  /** 卡片唯一标识 */
  key: 'resources' | 'views' | 'searches'
  /** 卡片标题（如"今日资源/总资源数"） */
  label: string
  /** 图标类名（Font Awesome） */
  icon: string
  /** 今日数值 */
  today: number
  /** 昨日数值（用于环比） */
  yesterday: number
  /** 总量（可选，资源卡片有，浏览/搜索无） */
  total?: number
  /** 数值格式化方式 */
  format: 'integer' | 'compact'
  /** 点击跳转目标（FR-010） */
  to: string
}

// 计算属性：环比百分比
interface StatCardComparison {
  /** 变化绝对值（today - yesterday） */
  delta: number
  /** 变化百分比（yesterday 为 0 时返回 null，前端显示"新增"） */
  percent: number | null
  /** 方向标识 */
  direction: 'up' | 'down' | 'flat' | 'new'
}
```

**验证规则**:
- `today` / `yesterday` 必须 ≥ 0
- `to` 必须是合法的内部路由（以 `/admin/` 开头）
- `format='compact'` 时数值 ≥ 10000 才缩写（如 1.2 万）

### 2. 侧边栏导航：`NavItem` / `NavGroup`

```ts
// composables/useAdminNav.ts 导出
interface NavItem {
  /** 路由路径（唯一键） */
  to: string
  /** 显示文本 */
  label: string
  /** 图标类名 */
  icon: string
  /** 当前路由是否激活此项 */
  active: (route: { path: string }) => boolean
}

interface NavGroup {
  /** 分组标识 */
  key: 'dashboard' | 'dataManagement' | 'systemConfig' | 'operation' | 'statistics'
  /** 分组标题 */
  title: string
  /** 分组内菜单项 */
  items: NavItem[]
}

// 侧边栏状态（useStorage 持久化）
interface SidebarState {
  /** 用户手动调整的分组展开状态（dashboard 始终展开，不存储） */
  [groupKey: string]: boolean
}
```

**持久化规则**:
- localStorage key: `urldb:admin:sidebar:expanded-groups`
- 序列化为 JSON：`{"dataManagement":true,"systemConfig":false,"operation":false,"statistics":false}`
- 进入页面时：自动展开当前路由所属分组（覆盖存储值），其余按存储值
- 用户 toggle：写入 storage（跨会话保留）

### 3. 命令面板：`CommandPaletteItem`

```ts
// components/Admin/CommandPalette.vue 的索引项
interface CommandPaletteItem {
  /** 唯一键（路由路径） */
  id: string
  /** 显示标题（如"资源管理"） */
  title: string
  /** 所属分组标题（如"数据管理"，用于副标题） */
  group: string
  /** 图标类名 */
  icon: string
  /** 跳转目标 */
  to: string
  /** 搜索关键词（含中英文、拼音首字母） */
  keywords: string[]
}

// 过滤函数（纯函数，可单元测试）
type FilterFn = (items: CommandPaletteItem[], query: string) => CommandPaletteItem[]
```

**索引构建**:
- 数据源：`useAdminNav()` 返回的 `NavGroup[]`
- 扁平化为 `CommandPaletteItem[]`，每项的 `keywords` 包含：标题、分组标题、英文翻译（可选）、拼音首字母（可选）

**过滤算法**（优先级降序）:
1. 标题完全包含 query → 最高优先级
2. keywords 数组中任一包含 query
3. 拼音首字母包含 query（若引入 pinyin-pro）
- 大小写不敏感
- query 为空时返回全部（按分组顺序）

### 4. 列表页通用模式：`FilterConfig` / `BatchAction`

```ts
// components/Admin/FilterBar.vue 的配置
interface FilterConfig {
  /** 搜索框配置（可选） */
  search?: { placeholder: string; key: string }
  /** 下拉筛选配置数组 */
  selects: Array<{
    key: string
    placeholder: string
    options: Array<{ label: string; value: string | number }>
  }>
}

// 批量操作配置
interface BatchAction {
  /** 操作标识 */
  key: string
  /** 按钮文本 */
  label: string
  /** Naive UI button type */
  type: 'primary' | 'info' | 'warning' | 'error'
  /** 图标 */
  icon: string
  /** 是否需要二次确认 */
  confirm?: { title: string; content: string }
  /** 执行函数，返回 Promise<BatchActionResult> */
  handler: (ids: (string | number)[]) => Promise<BatchActionResult>
}

interface BatchActionResult {
  success: boolean
  affected?: number
  message?: string
}
```

### 5. 表单分组：`FormGroup`

```ts
// 配置类页面（site-config / feature-config 等）的字段分组
interface FormGroup {
  /** 分组标题 */
  title: string
  /** 分组说明（展示在标题下方） */
  description?: string
  /** 是否默认展开 */
  defaultExpanded?: boolean
  /** 字段数组 */
  fields: FormField[]
}

interface FormField {
  /** 字段 key（对应后端字段名） */
  key: string
  /** 标签 */
  label: string
  /** 字段类型 */
  type: 'text' | 'number' | 'select' | 'switch' | 'textarea' | 'url'
  /** 是否必填 */
  required?: boolean
  /** 默认值 */
  default?: any
  /** 即时校验规则 */
  rules?: Array<{ validate: (v: any) => boolean; message: string }>
  /** 帮助文本 */
  help?: string
}
```

### 6. 主题配置：`ThemeConfig`

```ts
// composables/useTheme.ts
type ThemeMode = 'light' | 'dark'

interface ThemeConfig {
  mode: ThemeMode
  /** Naive UI GlobalThemeOverrides */
  naiveOverrides: import('naive-ui').GlobalThemeOverrides
  /** Chart.js 配置 */
  chartDefaults: {
    fontFamily: string
    gridColor: string
    primaryColor: string
  }
}
```

**主题切换副作用**:
- `<html>` 添加/移除 `dark` class（Tailwind darkMode:'class' 依赖）
- `n-config-provider` 的 `theme` 切换 `darkTheme` / `null`
- Chart.js 实例需销毁重建（监听 mode 变化）

---

## 数据流（关键场景）

### 仪表盘首屏加载

```text
admin/index.vue onMounted
  └─> useStatsApi().getSummary()  (单次 HTTP)
        └─> GET /api/stats/summary
              └─> stats_service.GetSummary()
                    ├─> 并行查询：resources/views/searches/todos
                    └─> 组装 StatsSummary 返回
  └─> 前端将 response 映射为 StatCardData[] 渲染卡片
  └─> StatCard 内部计算 StatCardComparison（环比）
```

### 命令面板触发跳转

```text
layouts/admin.vue 监听 Cmd/Ctrl+K (useMagicKeys from @vueuse/core)
  └─> CommandPalette 打开
        └─> useAdminNav() 构建 CommandPaletteItem[] 索引
        └─> 用户输入 → filterFn 过滤
        └─> Enter → router.push(item.to)
```

### 侧边栏状态持久化

```text
SidebarNav 初始化:
  const stored = useStorage('urldb:admin:sidebar:expanded-groups', {})
  watchEffect:
    - 当前路由变化 → 自动展开对应分组（覆盖 stored）
    - 用户 toggle → 写入 stored（跨会话保留）
```

---

## Post-Design Constitution Re-check

| 原则 | 设计落地点 | 结论 |
|------|-----------|------|
| **I. Test-First** | `StatCardComparison` 计算、`filterFn` 过滤、`SidebarState` 序列化、`stats_service.GetSummary` 均为纯逻辑/纯查询，可独立测试；新增 `tests/unit/` 覆盖 | ✅ |
| **II. Unified API Contract** | `StatsSummary` 经 `{code, message, data}` 封装，前端走 `useStatsApi` | ✅ |
| **III. Compile-Safe** | 后端改动仅在 `stats_handler.go` / `stats_service.go`，遵循现有模式 | ✅ |
| **IV. Integration Testing** | `stats_service_test.go` 覆盖真实 PG 聚合查询（不计其数边界：空表、跨日、负值防护） | ✅ |
| **V. Simplicity** | 无新表、无新依赖（前端用现有 @vueuse/core）、无投机抽象（CommandPalette 仅做页面跳转） | ✅ |

无未决问题，可进入 `/speckit-tasks`。
