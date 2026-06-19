# Quickstart: 转存文件定时自动清理

**Date**: 2026-06-14
**Feature**: 002-auto-cleanup-transfer

面向开发者的快速上手指南：本地拉起服务、开启清理功能、验证清理行为。

---

## 1. 环境准备

### 1.1 前置条件

- Go 1.24+
- Node.js 18+（前端，Nuxt 3）
- 数据库（配置在 `.env`，复用现有）
- 一个有效的 Quark 网盘账号 cookie（用于转存与删除）

### 1.2 拉起服务

```bash
# 后端（修改 golang 代码后需重启以保证编译生效）
go run main.go

# 前端（独立终端；前端改动自动热更新，无需重启）
cd web && npm run dev
```

后台访问：`http://localhost:<后端端口>/admin`（账号 `admin` / 密码 `password1`）

---

## 2. 开启自动清理（最小验证路径）

### 2.1 通过后台 UI

1. 进入"管理后台 → 功能配置"页面（`web/pages/admin/feature-config.vue`）
2. 找到"转存文件自动清理"区块（本期新增）
3. 开启开关，设置保留时长（天）与调度周期（分钟）
4. 保存

### 2.2 通过 API（等价路径）

```bash
# 开启清理，保留 7 天，每小时调度一次
curl -X PUT http://localhost:<port>/api/admin/system-config \
  -H "Content-Type: application/json" \
  -H "Cookie: <admin 会话 cookie>" \
  -d '{
    "auto_cleanup_enabled": true,
    "auto_cleanup_retention_days": 7,
    "auto_cleanup_interval_minutes": 60
  }'
```

---

## 3. 验证清理行为（短周期冒烟测试）

为快速验证，将保留时长设为最小值并制造一个"超期"资源：

1. 将 `auto_cleanup_retention_days` 暂设为 `1`、`auto_cleanup_interval_minutes` 设为 `1`（每分钟调度）
2. 确保至少一个资源满足：
   - `fid != ""` 且 `save_url != ""`
   - `transferred_at` 早于当前时间 1 天以上
   - （测试时可手动 SQL 更新该资源 `transferred_at = <1 天前>` 制造条件）
3. 等待下一个调度周期（≤ 1 分钟）
4. 观察日志：应出现"开始处理清理任务...""清理成功 - 资源ID: X"等条目
5. 查询该资源：`fid` 与 `save_url` 应已清空，`cleaned_at` 已写入

---

## 4. 失败场景验证

### 4.1 模拟"文件不存在"

- 手动在 Quark 网盘端删除某资源对应的文件
- 等待清理任务运行
- 预期：资源被标记为已清理（`fid`/`save_url` 清空），日志记录"文件不存在，视为清理成功"

### 4.2 模拟 cookie 失效

- 临时清空对应账号的 cookie
- 等待清理任务运行
- 预期：`clean_error_msg` 写入鉴权失败原因，`fid`/`save_url` 保留；下一轮任务（cookie 恢复后）自然再次尝试

---

## 5. 关闭清理

```bash
curl -X PUT http://localhost:<port>/api/admin/system-config \
  -H "Content-Type: application/json" \
  -H "Cookie: <admin 会话 cookie>" \
  -d '{"auto_cleanup_enabled": false}'
```

调度器收到停止信号后退出，不再执行任何删除。

---

## 6. 开发调试提示

- **日志位置**：`logs/` 目录；清理任务日志关键字：`清理任务`、`CleanupScheduler`
- **调度器状态**：`GlobalScheduler.IsCleanupRunning()`（实现后）可程序化查询
- **手动重置清理状态**（仅测试）：直接 SQL `UPDATE resources SET fid='', save_url='', cleaned_at=NOW(), clean_error_msg='' WHERE id=?` 模拟已清理
- **重新触发清理已清理资源**（仅测试）：`UPDATE resources SET fid='xxx', save_url='yyy', transferred_at='<过去>', cleaned_at=NULL WHERE id=?`
- **前端热更新**：修改 `web/` 下文件无需重启；修改 Go 代码需重启后端进程

---

## 7. 验收对照（与 spec 成功标准映射）

| SC | 快速验证方式 |
|----|------------|
| SC-001（1 周期内删除，成功率 95%） | 见第 3 节冒烟测试 |
| SC-002（1 分钟内完成配置且立即生效） | 第 2 节 UI/API 路径 + 不重启验证 |
| SC-003（账号空间下降） | 清理前后对比 Quark 账号已用容量 |
| SC-004（不影响其他功能） | 清理运行期间手动触发转存/分享/查询 |
| SC-005（失败可见 + 下轮自然重试） | 第 4 节失败场景验证 |
