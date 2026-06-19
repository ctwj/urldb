# Research: 链接检测优化 - 接入 PanCheck 服务

**Date**: 2026-06-13
**Status**: Complete — all deferred decisions resolved

本文档解决 spec 中推迟到规划阶段的实现决策。每项给出**决策**、**理由**、**备选**。所有决策均参照 `/Users/kerwin/Program/rust/tg_tool` 中已验证的 PanCheck 客户端实现 + 本项目技术栈（Go/Gin/GORM/PostgreSQL）。

---

## 1. 缓存存储选型

**Decision**: PostgreSQL（新增 `link_check_results` 表，via GORM），作为两处检测点共享的持久化缓存。

**Rationale**:
- 本项目主存储已是 PostgreSQL + GORM，缓存表走标准 `AutoMigrate` 即可，零新增基础设施。
- 持久化缓存意味着服务重启后结论不丢失，详情页打开即可命中（SC-003 要求二次检测请求下降 ≥80%，持久化缓存是达成该指标的前提）。
- tg_tool 同样用 DB 表（SQLite `link_check_results`）做缓存，行为可对齐。
- Redis 虽更轻量，但会引入本服务目前没有的额外依赖，过度工程化。

**Alternatives considered**:
- 进程内内存缓存（`sync.Map`/LRU）：重启即失，无法支撑 SC-003 跨重启/跨触发点的命中率。否决。
- Redis：需新增依赖与运维，对本规模过度。保留为未来横向扩展时的升级路径。

---

## 2. 缓存键与 URL 规范化

**Decision**: 缓存键 = 规范化后 URL 的 **SHA-256 十六进制**（`url_hash`），唯一索引。规范化规则（对齐 tg_tool `normalize_url`）：
1. `strings.TrimSpace`
2. 去掉 fragment（`#...`）
3. scheme + host 转小写（保留 path/query 原始大小写）
4. 去掉末尾 `/`

**Rationale**: 同一分享链接常以不同大小写/尾部斜杠/锚点形式出现，规范化保证命中；SHA-256 作为定长键便于唯一索引与批量查询。

**Alternatives considered**: 直接用原始 URL 做键——大小写/锚点差异导致漏命中，缓存效率低。否决。

---

## 3. 缓存写入策略

**Decision**: **仅缓存"有效/失效"二态结论**（与 spec FR-005、二态模型一致）。`pending`/`unknown` 不写缓存（异常批次视为未得出结论，保持 `is_valid` 原值，下次可重检）。缓存记录含 `status`、`platform`、`fail_reason`（失效时）、`checked_at`、`expires_at = checked_at + TTL`。

**Rationale**: 对齐 tg_tool（只缓存 valid/invalid）。异常瞬态若缓存会阻塞后续重检。

---

## 4. 运行参数默认值（可配置）

**Decision**: 全部纳入系统配置（`system_config` 表），默认值参照 tg_tool 成熟实践：

| 配置键 | 默认值 | 说明 |
|---|---|---|
| `pancheck_enabled` | `false` | PanCheck 总开关 |
| `pancheck_host` | `""` | 服务地址，如 `http://pancheck:6080` |
| `pancheck_timeout_seconds` | `60` | 单次 HTTP 请求超时 |
| `pancheck_batch_size` | `20` | 单批提交的 URL 数（tg_tool `URL_CHUNK_SIZE`） |
| `pancheck_concurrency` | `5` | 并发批次数（tg_tool，范围 1-20） |
| `pancheck_cache_ttl_hours` | `24` | 缓存有效期（小时） |

**Rationale**: tg_tool 已用这些值稳定运行；纳入系统配置使管理员可在不改代码/不重启下调整（SC-005）。`pancheck_enabled` 默认 `false`——未配置前不改变现有行为，避免误开。

**Alternatives considered**: 全部硬编码——无法运维调整，违反 SC-005。否决。

---

## 5. PanCheck 客户端调用契约

**Decision**: 调用 `POST {pancheck_host}/api/v1/links/check`，请求体：

```json
{
  "links": ["https://pan.quark.cn/s/xxx", "..."],
  "selected_platforms": ["quark","uc","baidu","tianyi","pan123","pan115","aliyun","xunlei","cmcc"]
}
```

响应解析（容错）：
- 优先读 `valid_links` / `invalid_links` / `pending_links`（字符串数组）；同时兼容对象数组 `[{url,platform,reason}]`。
- 按**规范化 URL** 匹配请求链接与返回分组；**invalid 优先于 valid**（防御性）。
- 既不在 valid 也不在 invalid 的链接 → 视为本次未得出结论（不翻转 `is_valid`、不写缓存）。
- 非 2xx / 网络错误 / JSON 解析失败 → 整批视为未得出结论（同上）。

**Rationale**: 对齐 tg_tool `PanCheckChecker::check` + `parse_pancheck_response` 的容错策略。

> **平台名映射注意**：PanCheck 服务端平台常量为 `pan123`/`pan115`/`cmcc`（见 `PanCheck-main/internal/model/platform.go`），而本项目历史用 `123pan`/`115`。由于本方案**总是发送全部 9 个平台**（不依赖按平台过滤），映射差异不影响检测；URL→平台的归属由 PanCheck 服务端识别。`selected_platforms` 统一使用 PanCheck 服务端常量集合。

详见 [contracts/pancheck-api.md](./contracts/pancheck-api.md)。

---

## 6. 失效原因（fail_reason）落位

**Decision**: 失效原因**存入缓存表** `link_check_results.fail_reason`；**不**写入 `Resource.ErrorMsg`（该字段语义为"转存失败原因"，属转存流程，不复用以免语义混淆）。资源维度展示失效原因时，按其 URL 查缓存表取最近一次 `fail_reason`。

**Rationale**: FR-007 要求"失效时记录失效原因"。缓存表天然按 URL 维度记录原因，资源展示时关联查询即可，无需在 `Resource` 上新增字段（与二态/is_valid 决策一致：不在资源上加多态字段）。

**Alternatives considered**:
- 复用 `Resource.ErrorMsg`：语义冲突（转存失败 vs 链接失效），且会覆盖既有转存错误信息。否决。
- 在 `Resource` 新增 `link_fail_reason` 字段：违反"不在资源上加多态字段"的澄清结论，且原因本就按 URL 缓存。否决。

---

## 7. 结论聚合与 is_valid 写回 + Meilisearch 同步

**Decision**: 单资源多 URL 时聚合规则——**任一 URL 失效即整体失效**（对齐 tg_tool `aggregate_link_status` 与 spec 边缘场景）。检测得出结论后：
- 仅当 `is_valid` **实际翻转**（true→false 或 false→true）时才写 DB + 同步 Meilisearch（避免每次巡检都产生写放大）。
- 写回逻辑复用现有 `ResourceRepository.Update` + 现有 Meilisearch 同步路径（`services/meilisearch_service.go`），与当前夸克写回路径一致，仅扩展到全平台。

**Rationale**: 写回仅在翻转时发生，配合共享缓存，使公开页匿名触发对 DB/Meilisearch 的压力可控（回应"公开触发"边缘场景）。

---

## 8. 两处触发点的接入方式

**Decision**: 新增 `services/link_check_service.go` 作为**唯一检测入口**，提供：
- `CheckResources(ctx, resources, ignoreCache) → map[resourceID]{status, failReason}` — 批量（详情页接口用）
- `CheckURL(ctx, url, ignoreCache) → {status, failReason}` — 单条（scheduler 用）

两处旧逻辑改为调用该服务：
- `handlers/resource_handler.go`：`BatchCheckResourceValidity`/`CheckResourceValidity` 内部把 `performAdvancedValidityCheck` 替换为 `linkCheckService.CheckResources(...)`；删除 `performAdvancedValidityCheck`/`performQuarkValidityCheck`/`performAlipanValidityCheck`。
- `scheduler/ready_resource.go`：`convertReadyResourceToResource` 中检测环节（当前非夸克走 `CheckURL`、夸克走转存）统一改为 `linkCheckService.CheckURL(...)`；夸克**先 PanCheck 校验**，通过后再执行既有的转存并分享（生成 `SaveURL`）流程。

**Rationale**: 单一入口消除两套并行逻辑（FR-009/FR-010），缓存与批量自然共享。

---

## 9. 降级与"PanCheck 关闭"语义

**Decision**: `link_check_service` 启动时/每次调用前读取 `pancheck_enabled` + `pancheck_host`：
- `enabled=false` 或 `host==""` → **跳过检测**，调用方按"未检测"处理：scheduler 直接放行资源（FR-004）；详情页接口返回"未启用检测"的明确状态（不报错）。
- `enabled=true` 但调用异常 → 异常批次视为未得出结论（保持 `is_valid`、不写缓存、可重试）。

**Rationale**: 对齐澄清结论"关闭即放行、移除旧代码、不回退"。

---

## 10. 旧代码删除清单

**Decision**: 删除以下文件/函数（FR-009，PanCheck 成为唯一主检测手段）：
- `common/utils/url_checker.go` 整文件（`CheckURL`、`extractShareID`、`checkUC`/`checkAliyun`/`checkQuark`/`check115`/`check123pan`/`checkTianyi`/`checkXunlei`/`checkBaidu`、`Test`、`CheckResult`）。
- `common/utils/url_checker_test.go` 整文件。
- `handlers/resource_handler.go` 中 `performAdvancedValidityCheck`、`performQuarkValidityCheck`、`performAlipanValidityCheck`。

**注意**: `common/pan_factory.go` 的 `ExtractShareId`/`ExtractServiceType`（转存流程与平台识别仍用）**保留**；只删检测相关。`extractShareID`（url_checker 内私有）随文件删除；平台识别如检测侧需要，由 PanCheck 服务端负责，客户端不再单独识别。

**Rationale**: 澄清明确"不保留旧代码作为回退"。

---

## 11. 前端改动范围

**Decision**: 前端改动最小化：
- 详情页 `/r/[key].vue`：检测结果展示从"夸克有效/无效/不支持检测"三态简化为"有效/失效 + 失效原因"；`onMounted` 自动触发 + 手动"链接检测"按钮逻辑**不变**（后端接口路径与请求体不变）；5 分钟节流保留。
- 系统配置页：新增 PanCheck 配置组（开关、host、超时、批量、并发、TTL），沿用现有系统配置页的统一封装风格。
- `useApi.ts`：`batchCheckResourceValidity` 复用；系统配置读写复用现有 `useSystemConfigApi`。

**Rationale**: CLAUDE.md 约定"前端 API 统一封装""前端改动不需重启"；后端接口契约不变，故前端触发逻辑零改动。
