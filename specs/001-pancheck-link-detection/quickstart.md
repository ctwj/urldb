# Quickstart: 链接检测优化 - 接入 PanCheck 服务

**Date**: 2026-06-13

面向开发者的端到端验证指南。前提：已部署一个可用的 PanCheck 服务实例（接口见 [contracts/pancheck-api.md](./contracts/pancheck-api.md)）。

## 1. 准备

```bash
# 确认在功能分支
git branch --show-current   # => 001-pancheck-link-detection

# 启动后端（开发环境会自动 AutoMigrate 新增 link_check_results 表 + PanCheck 默认配置）
# 数据库配置在 .env；账号 admin / 密码 password1
go build ./...              # CLAUDE.md 约定：修改 golang 代码须保证编译正常
```

## 2. 配置 PanCheck

管理后台 → 系统配置 → PanCheck 配置组：
- `pancheck_enabled` = `true`
- `pancheck_host` = `http://<你的 PanCheck 地址>:6080`
- 其余参数用默认值（超时 60s / 批量 20 / 并发 5 / TTL 24h）

保存后立即生效，无需重启（SC-005）。

## 3. 验证检测点 #1：添加资源后（scheduler）

1. 向"待处理资源"提交一个**已知有效**的夸克/阿里/百度等链接，触发 `ReadyResource → Resource` 转换。
2. 预期：有效链接转为正式资源（`is_valid=true`）。
3. 再提交一个**已知失效**的链接。
4. 预期：被拒绝，日志记录失效原因；夸克链接会先 PanCheck 校验、失效则不进入转存流程。

```bash
# 直查缓存表确认结论已写入
# psql: SELECT normalized_url, status, fail_reason, expires_at FROM link_check_results LIMIT 10;
```

## 4. 验证检测点 #2：前端详情页巡检

1. 浏览器打开公开资源页 `/r/<key>`。
2. 预期：页面加载后自动触发批量检测（5 分钟节流），各资源显示"有效/失效 + 失效原因"。
3. 对一个原本有效、手动改为失效的资源，刷新页面应看到状态翻转并写回 `is_valid` + 同步 Meilisearch。
4. 点"链接检测"按钮 = 强制重检（忽略 5 分钟节流，但仍走服务端缓存；手动重检见下）。

## 5. 验证缓存命中（SC-003）

1. 清空 `link_check_results`。
2. 对同一批链接做第一次检测 → 记录 PanCheck 实际被调用次数 N1。
3. 立即第二次检测同一批 → PanCheck 调用次数 N2 ≈ 0（全部命中缓存）。
4. 预期：N2 较 N1 下降 ≥80%。

## 6. 验证降级（SC-004 / FR-004）

- **关闭 PanCheck**（`pancheck_enabled=false`）：添加资源流程应直接放行、不报错；详情页接口返回 `detection_method: disabled`。
- **模拟不可达**（`pancheck_host` 指向不存在的地址）：检测应视为未得出结论，`is_valid` 保持原值，不产生"全部失效"误判；日志记录异常。

## 7. 验证夸克转存仍正常

提交一个有效的夸克链接：应先 PanCheck 校验通过，再执行既有"转存并分享"生成 `SaveURL`。确认 `SaveURL` 正常生成（转存业务不受影响）。

## 8. 验证旧代码已移除（FR-009）

```bash
# 应无输出（文件已删除）
ls common/utils/url_checker.go 2>/dev/null
# 编译应通过（无残留引用）
go build ./...
```

## 完成判据

- [ ] 两处检测点均经 PanCheck 得出结论
- [ ] 缓存命中使二次请求下降 ≥80%
- [ ] 降级不误判失效
- [ ] 夸克转存正常、有效性走 PanCheck
- [ ] 旧 `CheckURL` 代码已删除、编译通过
- [ ] 判定结果与 tg_tool 中 PanCheck 结论一致率 ≥95%（SC-001）
