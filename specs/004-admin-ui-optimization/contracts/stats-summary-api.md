# Contract: GET /api/stats/summary

**Feature**: 004-admin-ui-optimization
**Date**: 2026-06-19
**Constitution ref**: Principle II (Unified API Contract)

> 本契约定义后台仪表盘首屏聚合统计端点。前端必须通过 `useStatsApi().getSummary()` 调用，禁止直接 fetch。

## 端点

```
GET /api/stats/summary
```

- **认证**: 必须携带管理员 session/token（admin / password1）
- **权限**: 仅管理员角色可访问（沿用现有 admin 中间件）
- **缓存**: 不缓存（数据实时性要求高，仪表盘每次进入都拉取最新值）

## 请求

### Query Parameters

无。

### Headers

```
Cookie: <admin session>     # 现有鉴权方式
# 或
Authorization: Bearer <token>
```

## 响应

### 成功响应（HTTP 200）

遵循项目统一封装 `{code, message, data}`：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "resources": {
      "today": 12,
      "yesterday": 8,
      "total": 12345
    },
    "views": {
      "today": 450,
      "yesterday": 380
    },
    "searches": {
      "today": 89,
      "yesterday": 72
    },
    "todos": {
      "ready_resources": 3,
      "failed_tasks": 1,
      "pending_reports": 2
    }
  }
}
```

### 字段定义

| 路径 | 类型 | 说明 | 约束 |
|------|------|------|------|
| `code` | int | 业务状态码，200 表示成功 | 固定 200 |
| `message` | string | 状态描述 | 成功固定 "success" |
| `data.resources.today` | int | 今日新增资源数 | ≥ 0 |
| `data.resources.yesterday` | int | 昨日新增资源数 | ≥ 0 |
| `data.resources.total` | int | 资源总量 | ≥ 0 |
| `data.views.today` | int | 今日浏览量 | ≥ 0 |
| `data.views.yesterday` | int | 昨日浏览量 | ≥ 0 |
| `data.searches.today` | int | 今日搜索量 | ≥ 0 |
| `data.searches.yesterday` | int | 昨日搜索量 | ≥ 0 |
| `data.todos.ready_resources` | int | 待处理资源数 | ≥ 0 |
| `data.todos.failed_tasks` | int | 失败任务数 | ≥ 0 |
| `data.todos.pending_reports` | int | 待审核举报数 | ≥ 0 |

### 错误响应

遵循项目统一错误封装：

| HTTP | code | 场景 | message 示例 |
|------|------|------|-------------|
| 401 | 401 | 未鉴权 / session 失效 | "未登录或登录已过期" |
| 403 | 403 | 已登录但非管理员 | "权限不足" |
| 500 | 500 | 聚合查询失败 | "获取统计数据失败" |

错误响应体：

```json
{
  "code": 500,
  "message": "获取统计数据失败",
  "data": null
}
```

## 业务规则

1. **时区**: 使用服务器本地时区，日期边界为 `00:00:00`（与现有 stats 端点一致）
2. **数据源**:
   - `resources.*`: `resources` 表的 `created_at` 字段聚合
   - `views.*` / `searches.*`: 复用现有 `views-trend` / `searches-trend` 端点的数据源（不查 Meilisearch，仅 DB）
   - `todos.*`: `ready_resources` / `tasks` / `reports` 表的状态字段计数
3. **性能预算**: 单次响应 ≤ 500ms（p95），5 次聚合查询应并行执行
4. **空表降级**: 查询返回 0 而非 null（前端无需处理 null）
5. **负值防护**: 任何聚合结果为负（理论不应发生）一律返回 0 并记 warning 日志

## 前端调用契约

```ts
// composables/useApi.ts 扩展
export const useStatsApi = () => {
  // ...现有方法

  /** 获取仪表盘首屏聚合统计（含环比昨日与待办） */
  const getSummary = (): Promise<StatsSummary> =>
    useApiFetch('/stats/summary').then(parseApiResponse)

  return { /* ... */, getSummary }
}

// 类型定义（与后端 StatsSummary 对应）
interface StatsSummary {
  resources: { today: number; yesterday: number; total: number }
  views: { today: number; yesterday: number }
  searches: { today: number; yesterday: number }
  todos: { ready_resources: number; failed_tasks: number; pending_reports: number }
}
```

## 验收场景

| # | 场景 | 期望 |
|---|------|------|
| 1 | 管理员已登录，调用端点 | 返回 200 + 完整 data |
| 2 | 未登录调用 | 返回 401 |
| 3 | 非管理员登录调用 | 返回 403 |
| 4 | 系统刚部署（空表） | 返回 200 + 所有数值为 0 |
| 5 | 跨日时刻（00:00:01）调用 | today/yesterday 边界正确切换 |
| 6 | 单次响应耗时 | p95 ≤ 500ms |

## Constitution Compliance

- ✅ **Principle II**: 使用 `{code, message, data}` 统一封装；前端走 `useStatsApi`，无直接 fetch
- ✅ **Principle IV**: 集成测试覆盖空表、跨日、聚合正确性
- ✅ **Principle III**: 后端实现不破坏 `go build ./...`

## 不在范围内

- 不提供分页/筛选参数（仪表盘固定展示当日聚合）
- 不提供历史时间范围查询（趋势走现有 `/stats/views-trend` / `/stats/searches-trend`）
- 不写缓存层（实时性优先）
