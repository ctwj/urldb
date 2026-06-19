# Implementation Plan: 转存文件定时自动清理

**Branch**: `002-auto-cleanup-transfer` | **Date**: 2026-06-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-auto-cleanup-transfer/spec.md`

## Summary

为现有 Quark 自动转存流程追加"转存文件定时自动清理"能力：在资源转存到 Quark 账号并经过可配置的保留时长后，由独立调度器定时扫描、调用既有的 `QuarkPanService.DeleteFiles` 删除账号中的转存文件，并清空资源上的转存字段（`fid` / `save_url`）以表达"已清理"状态。功能复用现有调度框架（`BaseScheduler` + `Manager`）、系统配置体系（`system_config` 表）与 Quark 网盘服务，不实现专门重试机制，不提供手动触发入口。

## Technical Context

**Language/Version**: Go 1.24.0（后端）/ Vue 3 + Nuxt 3 + Naive UI（前端，`web/`）
**Primary Dependencies**: Gin（HTTP 路由）、GORM（ORM）、`github.com/robfig/cron/v3` 未使用——现有调度器采用 `time.Ticker`；前端 `web/composables/useApi.ts` 统一封装 API
**Storage**: GORM 支持的数据库（配置在 `.env`），核心表 `resources` / `system_configs`
**Testing**: Go 标准 `testing` 包（项目暂无系统性测试套件）
**Target Platform**: Linux 服务器（Go 二进制 + Vue SSR 同进程）
**Project Type**: web-service（单体，前后端同仓库）
**Performance Goals**: 清理任务单轮处理数百级资源即可；删除调用受 Quark API 频率隐式约束，不需显式限流（MVP）
**Constraints**:
- 复用现有 `BaseScheduler` + `Manager` + `GlobalScheduler` 调度框架，不引入新依赖
- 复用 `QuarkPanService.DeleteFiles([]string)` 删除能力（`common/quark_pan.go:284`），不新增网盘接口
- 复用 `system_config` 表 + `SystemConfigRepository` 配置体系，新增配置键需同步 `config/config.go` 的迁移/校验逻辑
- 资源"已清理"通过 `fid` + `save_url` 是否为空判定，不新增独立清理状态字段（spec 澄清结论）
- 不实现重试机制、不提供手动触发入口（spec 澄清结论）
**Scale/Scope**: 单实例部署；预计单轮清理涉及资源量 < 1000；新增 1 个调度器 + 1 个清理服务 + 配置项 + 前端配置入口 + 资源列表状态展示

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

项目 `.specify/memory/constitution.md` 为模板状态（未填写实际约束原则），无约束性 Gate 需校验。

**结论**: PASS（无约束）。Phase 1 设计完成后将再次确认未引入与项目惯例冲突的复杂度。

## Project Structure

### Documentation (this feature)

```text
specs/002-auto-cleanup-transfer/
├── plan.md              # 本文件
├── spec.md              # /speckit-specify 产出
├── research.md          # Phase 0 研究产出
├── data-model.md        # Phase 1 数据模型
├── quickstart.md        # Phase 1 快速上手
├── contracts/
│   └── api-contract.md  # Phase 1 接口契约
├── checklists/
│   └── requirements.md  # /speckit-specify 质量清单
└── tasks.md             # /speckit-tasks 产出（本命令不创建）
```

### Source Code (repository root)

```text
# 后端（Go）—— 在现有分层中扩展
common/
└── quark_pan.go            # 复用 DeleteFiles / deleteSingleFile（已存在，不改）

db/
├── entity/
│   ├── resource.go                  # 新增字段：transferred_at / cleaned_at / clean_error_msg
│   └── system_config_constants.go   # 新增 ConfigKey：auto_cleanup_enabled / _retention_days / _interval_minutes
├── repo/
│   └── resource_repository.go       # 新增查询：FindDueForCleanup / MarkCleaned / MarkCleanError
├── connection.go                    # 新增配置项种子数据
└── dto/system_config.go             # 新增 DTO 字段
└── converter/system_config_converter.go  # 新增字段映射

scheduler/
├── cleanup_scheduler.go   # 【新增】CleanupScheduler，参照 ready_resource.go 模式
├── manager.go             # 注册 CleanupScheduler（参照 HotDrama/ReadyResource 注册方式）
└── global.go              # 暴露 Start/Stop + UpdateSchedulerStatusWithCleanup

services/
└── cleanup_service.go     # 【新增】封装清理业务逻辑（筛选+调用 DeleteFiles+更新字段）

handlers/
└── system_config_handler.go  # 配置更新时启停 CleanupScheduler

# 前端（Vue 3 / Nuxt）—— 配置入口与状态展示
web/
├── pages/admin/feature-config.vue      # 新增"自动清理"配置区块
├── pages/admin/resources.vue           # 展示已清理标记 / 失败原因
├── composables/useApi.ts               # 复用统一 API 封装
└── stores/systemConfig.ts              # 新增 cleanup 相关 state
```

**Structure Decision**: 采用既有 monorepo 分层扩展，不新增顶层目录。后端按 `entity → repo → service → scheduler → handler` 分层，前端按 `pages/admin + composables + stores` 组织。`cleanup_scheduler.go` 与 `cleanup_service.go` 为本期新增的两个核心文件，分别承担"调度"与"清理业务"职责，保持与 `ready_resource.go`（调度）+ `transfer_processor.go`（业务）的一致分层。

## Complexity Tracking

> Constitution Check 无违规，本表无需填写。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |
