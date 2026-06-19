# Quickstart: 后台管理系统 UI/UX 整体优化

**Branch**: `004-admin-ui-optimization` | **Date**: 2026-06-19

> 本指南说明如何在本地启动、验证、测试本 feature 的改动。前端 dev server 热重载，无需重启；后端改动需重新编译运行。

---

## 1. 前置条件

| 工具 | 版本 | 用途 |
|------|------|------|
| Go | ≥ 1.21 | 后端编译运行 |
| Node.js | ≥ 18 | 前端构建 |
| pnpm | 9.x（项目锁定） | 前端包管理 |
| PostgreSQL | ≥ 12 | 后端数据源 |
| Meilisearch | 现有版本 | 搜索（本 feature 不改） |

## 2. 环境配置

`.env` 已存在并配置好（CLAUDE.md 约束）。如需调整：

```bash
# 关键变量（示例）
DB_HOST=localhost
DB_PORT=5432
DB_USER=...
DB_PASSWORD=...
DB_NAME=urldb
MEILISEARCH_URL=http://localhost:7700
```

## 3. 启动开发环境

### 后端（Go + Gin）

```bash
# 从仓库根目录
go build -o urldb .
./urldb
# 默认监听 :8080
```

修改 Go 代码后**必须重新编译运行**（CLAUDE.md: "如果修改了 golang 代码，修改完成后，保证编译正常"）。

### 前端（Nuxt 3 dev server）

```bash
cd web
pnpm install          # 首次或依赖变更后
pnpm dev              # 启动 dev server，默认 :3000
```

前端改动**自动热重载，无需重启**（CLAUDE.md 约束）。

浏览器访问 http://localhost:3000/admin，使用 `admin` / `password1` 登录。

## 4. 新增前端测试依赖

本 feature 引入 Vitest 测试基础设施（constitution Principle I 要求）：

```bash
cd web
pnpm add -D vitest @vue/test-utils happy-dom @nuxt/test-utils
```

## 5. 运行测试

### 前端单元测试

```bash
cd web
pnpm vitest run              # 单次运行
pnpm vitest                  # watch 模式
pnpm vitest run --coverage   # 覆盖率（需另装 @vitest/coverage-v8）
```

覆盖率门槛：新增逻辑代码 ≥ 90%（constitution Principle I）。

### 后端测试

```bash
# 单元测试
go test ./handlers/... ./services/...

# 集成测试（需启动 PostgreSQL）
go test ./services/... -tags=integration

# 覆盖率
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## 6. 验证 Checklist

完成实现后，按以下清单逐项验证（对应 spec 的成功指标 SC-001 ~ SC-008）：

### 视觉与主题

- [ ] 后台任一页面切换明/暗主题，所有卡片/图表/文字/图标正确适配，无对比度不足
- [ ] 用 [axe DevTools](https://www.deque.com/axe/devtools/) 或浏览器内置检查器验证正文对比度 ≥ 4.5:1（SC-004）
- [ ] 30+ 页面随机抽查 5 个，卡片/按钮/筛选栏/空状态/加载态样式一致（SC-003）

### 仪表盘（P1，SC-001）

- [ ] 登录后首屏 5 秒内可见：资源/浏览/搜索指标 + 待办聚合区
- [ ] 每张指标卡片显示环比昨日百分比与方向箭头
- [ ] 点击任一卡片跳转到对应详情页
- [ ] 趋势图为空（断开 API）时显示"暂无数据"而非空坐标轴
- [ ] 浏览器控制台**无** `console.log` 输出（FR-012）

### 导航（P2，SC-002）

- [ ] 任一二级页面顶部显示面包屑，每级可点击返回
- [ ] 进入 `/admin/resources`，侧边栏"数据管理"分组自动展开并高亮
- [ ] 手动展开"系统配置"分组，刷新页面或重开浏览器后仍保持展开（localStorage）
- [ ] 按 `Cmd+K`（macOS）或 `Ctrl+K`（Windows/Linux）唤起命令面板
- [ ] 在命令面板输入"版权"，能找到并跳转到"版权申述"页

### 列表页（P2，SC-005）

- [ ] 资源管理、举报管理、待处理资源三个列表页的筛选栏布局一致
- [ ] 任一列表页选中若干条目，"共 N 项，已选 M 项"提示正确
- [ ] 批量删除时弹出二次确认
- [ ] 操作完成后出现统一格式的成功/失败 toast
- [ ] 列表为空时显示统一空状态（图标 + 说明 + 可选 CTA）

### 配置页（P3，SC-006）

- [ ] 站点配置页字段按语义分组，每组有说明
- [ ] 必填项有统一标识
- [ ] 在字段输入非法值并离开焦点，显示即时校验错误
- [ ] 修改未保存时尝试离开页面，弹出"未保存确认"

### 响应式（SC-008）

- [ ] 浏览器宽度调至 ≤ 768px，侧边栏收起为图标态或抽屉态
- [ ] 窄屏下导航开关、筛选、批量操作、保存按钮均可达且不错位

### 多标签页（Clarifications 决策）

- [ ] 开两个标签页分别打开同一配置页，A 改值保存，B 改不同值保存 → 后保存者覆盖（last-write-wins，无冲突提示）

## 7. 性能验证

```bash
# 浏览器 DevTools Network 标签
# 1. 进入 /admin 首页，确认仅 1 次 /api/stats/summary 请求（非 3 次分散请求）
# 2. 该请求 p95 ≤ 500ms

# Lighthouse（可选）
# Performance ≥ 90（桌面端 /admin 页面）
# Accessibility ≥ 95
```

## 8. Constitution 合规自检

| 原则 | 自检命令 | 期望 |
|------|---------|------|
| I. Test-First | `cd web && pnpm vitest run --coverage` | 新增逻辑 ≥ 90% |
| I. Test-First | `go test ./... -cover` | 新增 handler/service ≥ 90% |
| II. Unified API | grep 确认无 `fetch('/api/stats/summary')` 直接调用 | 全部走 `useStatsApi` |
| III. Compile-Safe | `go build ./...` | 成功，无报错 |
| III. Compile-Safe | `go vet ./...` | 无 warning |
| IV. Integration | `go test ./services/... -tags=integration` | 通过 |
| V. Simplicity | 检查未引入新组件库/图表库 | package.json 无投机依赖 |

## 9. 故障排查

| 现象 | 排查方向 |
|------|---------|
| 前端改动不生效 | 确认 dev server 运行中；清除 `.nuxt` 缓存：`rm -rf web/.nuxt` |
| 暗色模式部分元素白底 | 检查该元素是否用了 `bg-white` 而未加 `dark:bg-slate-800` |
| Chart.js 暗色下网格不显示 | 确认 `useTheme` 监听主题变化并销毁重建实例 |
| 命令面板无法唤起 | 检查 `useMagicKeys` 是否在 `onMounted` 中注册；浏览器是否拦截了 Cmd+K |
| 侧边栏状态不持久化 | 浏览器 DevTools → Application → LocalStorage 查看 `urldb:admin:sidebar:expanded-groups` |
| `/api/stats/summary` 401 | 确认已登录 admin；session 未过期 |

## 10. 下一步

实现完成后运行：

```
/speckit-tasks      # 基于本 plan + data-model + contracts 生成任务清单
/speckit-analyze    # 跨产物一致性检查
/speckit-implement  # 执行实现
```
