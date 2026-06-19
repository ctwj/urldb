# API Contract: 转存文件定时自动清理

**Date**: 2026-06-14
**Feature**: 002-auto-cleanup-transfer

本项目为 Web 服务（Gin 后端 + Vue/Nuxt 前端），对外契约主要是管理后台 HTTP 接口与前端 API 封装。本期不新增独立端点，**复用现有"系统配置批量更新/查询"端点**承载清理配置，并在既有"资源列表"端点的响应中扩展清理相关字段。

---

## 1. 系统配置接口（复用，扩展字段）

### 1.1 获取系统配置

- **端点**: `GET /api/admin/system-config`
- **认证**: 复用现有 admin 会话/JWT 中间件（账号 `admin` / 密码 `password1`）
- **响应**: 沿用现有 `SystemConfigResponse` 结构，新增以下字段：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "auto_cleanup_enabled": false,
    "auto_cleanup_retention_days": 7,
    "auto_cleanup_interval_minutes": 60,
    "... 其他既有配置字段 ..."
  }
}
```

### 1.2 更新系统配置

- **端点**: `PUT /api/admin/system-config`
- **认证**: 同上
- **请求体**: 沿用现有 `SystemConfigRequest` 结构，新增可选字段（指针类型，未传则不修改）：

```json
{
  "auto_cleanup_enabled": true,
  "auto_cleanup_retention_days": 7,
  "auto_cleanup_interval_minutes": 60
}
```

- **校验规则**（失败返回 `code != 200` + `message` 说明原因）：

| 字段 | 规则 | 错误 message 示例 |
|------|------|-------------------|
| `auto_cleanup_retention_days` | 整数，`1 ≤ x ≤ 365` | "保留时长须在 1-365 天之间" |
| `auto_cleanup_interval_minutes` | 整数，`1 ≤ x ≤ 1440` | "调度周期须在 1-1440 分钟之间" |
| `auto_cleanup_enabled` | 布尔 | （类型不符由绑定层拦截） |

- **副作用**: 配置保存成功后，handler 调用 `GlobalScheduler.UpdateSchedulerStatusWithCleanup(...)`：
  - `enabled == true` 且调度器未运行 → 启动 `CleanupScheduler`
  - `enabled == false` 且调度器运行中 → 停止 `CleanupScheduler`
  - `interval_minutes` 变更 → 下一轮 ticker 重建时生效（重启调度器或下一周期自然应用，MVP 不强制立即重建）

- **统一响应格式**（遵循项目约定）:

```json
{
  "code": 200,
  "message": "配置更新成功",
  "data": null
}
```

---

## 2. 资源列表接口（复用，扩展字段）

### 2.1 资源列表

- **端点**: `GET /api/admin/resources`（既有，分页参数不变）
- **响应**: `Resource` 结构新增以下字段（既有字段保持不变）：

```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 123,
        "title": "...",
        "url": "...",
        "fid": "abc...|（空字符串表示已清理）",
        "save_url": "https://...|（空字符串表示已清理）",
        "transferred_at": "2026-06-10T12:00:00Z|null",
        "cleaned_at": "2026-06-14T13:00:00Z|null",
        "clean_error_msg": "|文件不存在以外的失败原因",
        "last_clean_error_at": "2026-06-14T12:55:00Z|null",
        "... 其他既有字段 ..."
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 200
  }
}
```

### 2.2 前端展示契约

- `fid == "" && save_url == ""` → 列表展示"已清理"标签 + `cleaned_at`（若 `transferred_at != nil`）
- `clean_error_msg != ""` → 鼠标悬浮/展开行展示失败原因 + `last_clean_error_at`
- `transferred_at == null && fid == ""` → 展示"未转存"（区别于"已清理"）

---

## 3. 调度器启停状态（可选查询，MVP 可不暴露独立端点）

调度器运行状态通过现有"系统配置详情"接口的 `auto_cleanup_enabled` 字段间接反映（配置为 true 即视为预期运行）。如后续需要"实际运行中"的精确状态，可扩展 `GET /api/admin/system-config/scheduler-status`，本期不实现。

---

## 4. 前端 API 封装契约（`web/composables/useApi.ts`）

遵循项目"前端所有 API 都需要统一封装"约定，本期**不新增独立 API 方法**——清理配置复用 `useSystemConfigApi().updateSystemConfig(payload)` 既有方法，仅扩展 payload 类型：

```typescript
// web/composables/useApi.ts 中 SystemConfigApi 既有方法签名不变
interface SystemConfigPayload {
  // ... 既有字段 ...
  auto_cleanup_enabled?: boolean
  auto_cleanup_retention_days?: number
  auto_cleanup_interval_minutes?: number
}
```

资源列表 API（`useResourceApi().getResources(params)`）方法签名不变，仅响应类型扩展清理字段。

---

## 5. 错误响应统一格式

所有校验失败/服务端错误遵循项目既有统一返回格式：

```json
{
  "code": 400,
  "message": "保留时长须在 1-365 天之间",
  "data": null
}
```

`code` 取值：`200` 成功、`400` 参数校验失败、`401` 未认证、`500` 服务端错误。

---

## 6. 契约不变性

- 本期不新增任何 HTTP 端点，所有变更均为既有端点的字段扩展
- 既有字段语义与结构保持不变，前端旧版本读取新响应不会破坏
- 新增字段在响应中始终存在（不为节省体积而省略），便于前端稳定解构
- 删除能力（`QuarkPanService.DeleteFiles`）不通过 HTTP 暴露，仅由内部调度器调用
