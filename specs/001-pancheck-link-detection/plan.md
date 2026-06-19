# Implementation Plan: 链接检测优化 - 接入 PanCheck 服务

**Branch**: `001-pancheck-link-detection` | **Date**: 2026-06-13 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-pancheck-link-detection/spec.md`

## Summary

用一个集中式的 **PanCheck 服务**替换 urldb 当前的两套并行链接检测逻辑（① `ReadyResource→Resource` 转换流程中的内联逐平台 `CheckURL`；② 前端详情页批量检测接口背后的 `performAdvancedValidityCheck`/`performQuarkValidityCheck`）。新增一个共享的链接检测服务层：负责 URL 规范化、按 URL 维度的检测结果缓存（PostgreSQL）、批量分组与并发提交、PanCheck HTTP 客户端调用与容错。两处检测触发点统一调用该服务；结论（有效/失效二态）写回 `Resource.is_valid` 并同步 Meilisearch。夸克的有效性检测改走 PanCheck（不再用 `Transfer()` 探测），但夸克"转存并分享"生成 `SaveURL` 的转存业务保留。旧的内联检测代码删除。

## Technical Context

**Language/Version**: Go 1.24.0（后端）；前端 Nuxt ^3.8.0 / Vue ^3.3.0 / TypeScript
**Primary Dependencies**: Gin 1.10.1（web）、GORM + PostgreSQL（存储）、`go-resty/resty/v2` v2.16.5（HTTP 客户端，已用于现有 url_checker）、`meilisearch/meilisearch-go` v0.33.1（检索同步）、cobra（CLI）
**Storage**: PostgreSQL（via GORM；`db/connection.go` 的 `AutoMigrate` 列表新增缓存表实体）
**Testing**: `go test`（现有 `common/utils/url_checker_test.go`、`common/pan_factory_test.go`）；新增 PanCheck 客户端与缓存服务的单元测试
**Target Platform**: Linux 服务器（自托管）
**Project Type**: web-service（Go 后端 + Nuxt 前端）
**Performance Goals**: 缓存命中使二次检测 PanCheck 请求数下降 ≥80%（SC-003）；详情页检测不阻塞页面渲染（依赖缓存/异步）；单次 PanCheck 请求超时 60s
**Constraints**: PanCheck 为外部服务，必须容错降级——不可达时保持 `is_valid` 原值、不误判失效（SC-004）；PanCheck 关闭/未配置时跳过检测放行（FR-004）
**Scale/Scope**: 自托管资源库；详情页检测接口当前无鉴权、由公开页 `/r/[key]` 触发，依赖共享缓存 + 前端 5 分钟节流控制负载

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

项目 constitution 文件（`.specify/memory/constitution.md`）当前为**未初始化模板**（全部为 `[PRINCIPLE_X_NAME]` / `[PLACEHOLDER]` 占位符），无已批准的核心原则可供评估。

**Gate 结果: PASS（无具体原则可违反，无违规）**。

> 建议：在后续迭代中通过 `/speckit-constitution` 建立项目原则（如"接口统一返回格式""前端 API 统一封装""修改 Golang 代码须保证编译通过"——这些已在 `CLAUDE.md` 项目约定中体现）。本计划已主动遵循 `CLAUDE.md` 中的项目约定。

## Project Structure

### Documentation (this feature)

```text
specs/001-pancheck-link-detection/
├── plan.md              # 本文件
├── research.md          # Phase 0：延迟决策研究结论
├── data-model.md        # Phase 1：实体与数据模型
├── quickstart.md        # Phase 1：快速验证指南
├── contracts/           # Phase 1：接口契约
│   └── pancheck-api.md  # PanCheck 服务接口契约（外部依赖）
└── tasks.md             # Phase 2（/speckit-tasks 生成，非本命令产物）
```

### Source Code (repository root)

```text
# —— 新增 ——
db/entity/
└── link_check_result.go                 # 缓存实体（按规范化 URL 维度）
db/repo/
└── link_check_result_repository.go      # 缓存仓储（读写/批量/过期清理）
services/
├── pancheck_client.go                   # PanCheck HTTP 客户端（调用 /api/v1/links/check、容错解析）
└── link_check_service.go                # 共享检测服务（URL 规范化、缓存命中、批量分组、并发、结论聚合、is_valid 写回 + Meilisearch 同步）
handlers/
└── (修改) resource_handler.go           # BatchCheckResourceValidity / CheckResourceValidity 改为调用 link_check_service
scheduler/
└── (修改) ready_resource.go             # convertReadyResourceToResource 的检测环节改为调用 link_check_service（含夸克前置校验）
db/entity/
└── (修改) system_config_constants.go    # 新增 PanCheck 配置键常量与默认值
db/connection.go                         # AutoMigrate 列表新增 LinkCheckResult；createIndexes 新增缓存表索引
db/connection.go (insertDefaultDataIfEmpty) # 新增 PanCheck 默认系统配置项
main.go                                  # 路由保持；DI 注入新服务

# —— 删除 ——
common/utils/url_checker.go              # 旧内联逐平台检测（CheckURL + 各 checkXxx）
common/utils/url_checker_test.go         # 对应测试

# —— 前端（小改） ——
web/composables/useApi.ts                # 已有 batchCheckResourceValidity 复用；新增 useSystemConfigApi 的 PanCheck 配置项读写
web/pages/admin/（系统配置页）            # 新增 PanCheck 服务地址/开关/参数的配置 UI（沿用现有系统配置页风格）
web/pages/r/[key].vue                    # 检测结果展示适配二态（有效/失效）+ 失效原因；触发逻辑不变
```

**Structure Decision**: 沿用现有"web-service"分层（handlers / services / db(entity+repo) / scheduler / common）。新增的检测能力收敛到 `services/` 层（`pancheck_client.go` + `link_check_service.go`），作为两处触发点的唯一入口，避免逻辑散落。缓存作为独立实体走标准 GORM AutoMigrate 流程。旧的内联检测代码整体删除（FR-009），不留回退。

## Complexity Tracking

> Constitution Check 无违规，无需填写。
