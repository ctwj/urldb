# Contracts: PanCheck 服务接口（外部依赖）

**Date**: 2026-06-13

本文档定义 urldb 作为**客户端**调用外部 PanCheck 服务的接口契约，以及本系统对外（前端）保持不变的检测接口契约。

---

## 1. 外部依赖：PanCheck 服务

### 端点

```
POST {pancheck_host}/api/v1/links/check
```

`pancheck_host` 来自系统配置 `pancheck_host`（如 `http://pancheck:6080`）。无鉴权（沿用 PanCheck 默认公开访问，与 tg_tool 一致）。

### 请求

**Headers**: `Content-Type: application/json`
**Body**:

```json
{
  "links": [
    "https://pan.quark.cn/s/abc123",
    "https://www.alipan.com/s/xyz"
  ],
  "selected_platforms": ["quark","uc","baidu","tianyi","pan123","pan115","aliyun","xunlei","cmcc"]
}
```

| 字段 | 类型 | 说明 |
|---|---|---|
| `links` | string[] | 待检测 URL 列表 |
| `selected_platforms` | string[] | **固定发送全部 9 个平台**，不依赖按平台过滤 |

### 响应（成功）

```json
{
  "submission_id": 1,
  "valid_links": ["https://pan.quark.cn/s/abc123"],
  "invalid_links": ["https://www.alipan.com/s/xyz"],
  "pending_links": [],
  "total_duration": 2.5,
  "invalid_format_count": 0,
  "duplicate_count": 0
}
```

| 字段 | 类型 | 说明 |
|---|---|---|
| `valid_links` | string[] | 有效链接（规范化 URL 字符串） |
| `invalid_links` | string[] | 失效链接 |
| `pending_links` | string[] | 检测中（正常为空） |

**容错解析要求**（对齐 tg_tool）:
- 同时兼容对象数组形式：`valid_links: [{"url":"...","platform":"quark","reason":"..."}]`。
- 字段名容错：valid 组兼容 `valid_links`/`available`/`ok`/`valid`；invalid 组兼容 `invalid_links`/`unavailable`/`expired`/`dead`/`invalid`；pending 组兼容 `pending_links`/`pending`/`checking`。
- 按**规范化 URL** 匹配请求链接与返回分组（urldb 侧与 PanCheck 侧都做规范化，需一致）。
- **invalid 优先于 valid**（同一 URL 出现在两组时按失效处理）。
- 既不在 valid 也不在 invalid 的链接 → 本次未得出结论。

### 错误处理

| 情况 | urldb 客户端行为 |
|---|---|
| 网络超时（> `pancheck_timeout_seconds`） | 整批视为未得出结论：不翻转 `is_valid`、不写缓存、可重试 |
| 非 2xx 响应 | 同上 |
| JSON 解析失败 / 结构不符 | 同上 |
| `enabled=false` 或 `host==""` | 不发起请求，调用方按"未启用检测"处理 |

任何异常都**不得**把链接判为失效（SC-004）。

### 平台常量映射

urldb 客户端**固定发送** PanCheck 服务端平台常量集合：
`["quark","uc","baidu","tianyi","pan123","pan115","aliyun","xunlei","cmcc"]`

> 注意：urldb 历史用 `123pan`/`115`，PanCheck 服务端用 `pan123`/`pan115`/`cmcc`。因本方案总发送全集、不依赖过滤，差异不影响检测；URL→平台归属由 PanCheck 服务端识别。

---

## 2. 本系统对外接口（前端调用，契约不变）

### 详情页批量检测（保持现有契约）

```
POST /api/resources/validity/batch
```

**请求**（不变）:
```json
{ "ids": [1, 2, 3] }
```
（最多 20 个 ID）

**响应**（结构不变，语义扩展到全平台）:
```json
{
  "code": 200,
  "data": {
    "results": [
      { "resource_id": 1, "is_valid": true,  "last_checked": "2026-06-13T16:00:00Z", "detection_method": "pancheck", "error": "" },
      { "resource_id": 2, "is_valid": false, "last_checked": "2026-06-13T16:00:00Z", "detection_method": "pancheck", "error": "链接已失效" }
    ],
    "total": 2
  }
}
```

> 统一返回格式遵循 CLAUDE.md 约定（`code`/`data`）。`detection_method` 由旧的 `quark_deep`/`unsupported` 统一改为 `pancheck`（启用时）或 `disabled`（未启用）。`is_valid` 仅在结论翻转时回写 DB + 同步 Meilisearch。

### 单资源检测（保持现有契约）

```
GET /api/resources/:id/validity
```

返回单个资源的检测结论，语义同上。

---

## 3. 配置接口（前端管理后台）

复用现有系统配置读写接口（`useSystemConfigApi`）。新增配置键见 [data-model.md](../data-model.md#systemconfig新增配置键)。前端系统配置页新增 PanCheck 配置组。
