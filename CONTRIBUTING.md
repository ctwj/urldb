# 🤝 贡献指南

感谢您对 L9Pan 项目的关注！我们欢迎所有形式的贡献。

## 📋 贡献类型

- 🐛 Bug 修复
- ✨ 新功能开发
- 📚 文档改进
- 🎨 UI/UX 优化
- ⚡ 性能优化
- 🧪 测试用例
- 🔧 工具改进

## 🚀 开始贡献

### 1. Fork 项目

1. 访问 [GitHub 仓库](https://github.com/ctwj/urldb)
2. 点击右上角的 Fork 按钮
3. 选择您的 GitHub 账户

### 2. 克隆您的 Fork

```bash
git clone https://github.com/YOUR_USERNAME/resManage.git
cd resManage
```

### 3. 添加上游仓库

```bash
git remote add upstream https://github.com/ctwj/urldb.git
```

### 4. 创建功能分支

```bash
git checkout -b feature/your-feature-name
# 或者
git checkout -b fix/your-bug-fix
```

## 💻 开发环境设置

### 后端开发

```bash
# 安装 Go 依赖
go mod tidy

# 复制环境变量文件
cp env.example .env

# 编辑环境变量
vim .env

# 启动开发服务器
go run .
```

### 前端开发

```bash
# 进入前端目录
cd web

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

### 数据库设置

```sql
CREATE DATABASE url_db;
```

## 📝 开发规范

### 代码风格

#### Go 代码规范
- 遵循 [Go 官方代码规范](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代码
- 函数名使用驼峰命名法
- 包名使用小写字母
- 错误处理要明确

#### TypeScript/JavaScript 规范
- 使用 TypeScript 编写代码
- 遵循 ESLint 规则
- 使用 Prettier 格式化代码
- 组件名使用 PascalCase
- 变量名使用 camelCase

### 提交信息规范

使用中文描述提交信息，格式如下：

```
类型(范围): 简短描述

详细描述（可选）

相关Issue: #123```

类型包括：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(资源管理): 添加批量导入功能

- 支持Excel文件导入
- 添加导入进度显示
- 优化错误处理

相关Issue: #45
```

### 分支命名规范

- `feature/功能名称` - 新功能开发
- `fix/问题描述` - Bug修复
- `docs/文档类型` - 文档更新
- `refactor/重构内容` - 代码重构
- `test/测试内容` - 测试相关

## 🧪 测试

### 后端测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./handlers

# 运行测试并显示覆盖率
go test -cover ./...
```

### 前端测试

```bash
cd web

# 运行单元测试
pnpm test

# 运行E2E测试
pnpm test:e2## 📋 Pull Request 流程

###1. 确保代码质量

-  ] 代码通过所有测试
-  ] 代码符合项目规范
-  ] 添加了必要的文档
-  ] 更新了相关文档

### 2 提交 Pull Request

1推送您的分支到您的 Fork
```bash
git push origin feature/your-feature-name
```

2. 在 GitHub 上创建 Pull Request3. 填写 PR 模板4 等待代码审查

###3PR 模板

```markdown
## 📝 描述

简要描述您的更改

## 🎯 类型

- ] Bug 修复
-  ] 新功能
- ] 文档更新
- ] 代码重构
- ] 性能优化
- 相关

## 🔗 相关 Issue

关闭 #123

## ✅ 检查清单

-  ] 我的代码遵循项目的代码规范
- ] 我已经自测了我的更改
- ] 我已经添加了必要的测试
- ] 我已经更新了相关文档
- ] 我的更改不会产生新的警告
-  我添加了示例来证明我的更改是有效的

## 📸 截图（如果适用）

## 🔧 测试

请描述您如何测试您的更改

## 📋 其他信息

任何其他信息或上下文
```

## 🐛 报告 Bug

### Bug 报告模板

```markdown
## 🐛 Bug 描述

清晰简洁地描述 bug

## 🔄 重现步骤
1. 进入...2. 点击 '....3. 滚动到 ....4 看到错误

## ✅ 预期行为

清晰简洁地描述您期望发生的事情

## 📸 截图

如果适用，添加截图以帮助解释您的问题

## 💻 环境信息

- 操作系统: [例如 Windows 10]
- 浏览器: [例如 Chrome, Safari]
- 版本:例如 22]

## 📱 移动设备（如果适用）

- 设备: [例如 iPhone 6]
- 操作系统: 例如 iOS 8.1]
- 浏览器:例如 Safari]
- 版本: [例如 22]

## 🔧 其他上下文

在此处添加有关问题的任何其他上下文
```

## 💡 功能请求

### 功能请求模板

```markdown
## 💡 功能描述

清晰简洁地描述您想要的功能

## 🎯 问题描述

清晰简洁地描述这个功能要解决的问题

## 💭 建议的解决方案

清晰简洁地描述您希望发生的事情

## 🔄 替代方案

清晰简洁地描述您考虑过的任何替代解决方案或功能

## 📋 其他上下文

在此处添加有关功能请求的任何其他上下文或截图
```

## 📞 联系我们

如果您有任何问题或需要帮助，请通过以下方式联系我们：

- 📧 邮箱: 510199617@qq.com
- 💬 讨论区: [GitHub Discussions](https://github.com/your-username/l9pan/discussions)
- 🐛 问题反馈: [GitHub Issues](https://github.com/your-username/l9pan/issues)

## 🙏 致谢

感谢所有为 L9项目做出贡献的开发者！

---

**记住：每一个贡献，无论大小，都是宝贵的！** 🎉 