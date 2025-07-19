# 💻 本地开发

## 环境准备

### 1. 安装必需软件

#### Go 环境
```bash
# 下载并安装 Go 1.23+
# 访问 https://golang.org/dl/
# 或使用包管理器安装

# 验证安装
go version
```

#### Node.js 环境
```bash
# 下载并安装 Node.js 18+
# 访问 https://nodejs.org/
# 或使用 nvm 安装

# 验证安装
node --version
npm --version
```

#### PostgreSQL 数据库
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS (使用 Homebrew)
brew install postgresql

# 启动服务
sudo systemctl start postgresql  # Linux
brew services start postgresql   # macOS
```

#### pnpm (推荐)
```bash
# 安装 pnpm
npm install -g pnpm

# 验证安装
pnpm --version
```

### 2. 克隆项目

```bash
git clone https://github.com/ctwj/urldb.git
cd urldb
```

## 后端开发

### 1. 环境配置

```bash
# 复制环境变量文件
cp env.example .env

# 编辑环境变量
vim .env
```

### 2. 数据库设置

```sql
-- 登录 PostgreSQL
sudo -u postgres psql

-- 创建数据库
CREATE DATABASE url_db;

-- 创建用户（可选）
CREATE USER url_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE url_db TO url_user;

-- 退出
\q
```

### 3. 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 验证依赖
go mod verify
```

### 4. 启动后端服务

```bash
# 开发模式启动
go run main.go

# 或使用 air 热重载（推荐）
go install github.com/cosmtrek/air@latest
air
```

## 前端开发

### 1. 进入前端目录

```bash
cd web
```

### 2. 安装依赖

```bash
# 使用 pnpm (推荐)
pnpm install

# 或使用 npm
npm install
```

### 3. 启动开发服务器

```bash
# 开发模式
pnpm dev

# 或使用 npm
npm run dev
```

### 4. 访问前端

前端服务启动后，访问 http://localhost:3000

## 开发工具

### 推荐的 IDE 和插件

#### VS Code
- **Go** - Go 语言支持
- **Vetur** 或 **Volar** - Vue.js 支持
- **PostgreSQL** - 数据库支持
- **Docker** - Docker 支持
- **GitLens** - Git 增强

#### GoLand / IntelliJ IDEA
- 内置 Go 和 Vue.js 支持
- 数据库工具
- Docker 集成

### 代码格式化

```bash
# Go 代码格式化
go fmt ./...

# 前端代码格式化
cd web
pnpm format
```

### 代码检查

```bash
# Go 代码检查
go vet ./...

# 前端代码检查
cd web
pnpm lint
```

## 调试技巧

### 后端调试

```bash
# 使用 delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go

# 或使用 VS Code 调试配置
```

### 前端调试

```bash
# 启动开发服务器时开启调试
cd web
pnpm dev --inspect
```

### 数据库调试

```bash
# 连接数据库
psql -h localhost -U postgres -d url_db

# 查看表结构
\dt

# 查看数据
SELECT * FROM users LIMIT 5;
```

## 测试

### 后端测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./handlers

# 生成测试覆盖率报告
go test -cover ./...
```

### 前端测试

```bash
cd web

# 运行单元测试
pnpm test

# 运行 E2E 测试
pnpm test:e2e
```

## 构建

### 后端构建

```bash
# 构建二进制文件
go build -o urlDB main.go

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o urlDB-linux main.go
```

### 前端构建

```bash
cd web

# 构建生产版本
pnpm build

# 预览构建结果
pnpm preview
```

## 常见问题

### 1. 端口冲突

如果遇到端口被占用的问题：

```bash
# 查看端口占用
lsof -i :8080
lsof -i :3000

# 杀死进程
kill -9 <PID>
```

### 2. 数据库连接失败

检查 `.env` 文件中的数据库配置：

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=url_db
```

### 3. 前端依赖安装失败

```bash
# 清除缓存
pnpm store prune
rm -rf node_modules
pnpm install
```

## 下一步

- [了解项目架构](../architecture/overview.md)
- [查看 API 文档](../api/overview.md)
- [学习代码规范](../development/coding-standards.md) 