# 资源管理系统

一个基于 Golang + Nuxt.js 的资源管理系统，参考网盘资源管理界面设计。

## 技术栈

### 后端
- **Golang** - 主要编程语言
- **Gin** - Web框架
- **PostgreSQL** - 数据库
- **lib/pq** - PostgreSQL驱动

### 前端
- **Nuxt.js 3** - Vue.js框架
- **Vue 3** - 前端框架
- **TypeScript** - 类型安全
- **Tailwind CSS** - 样式框架

## 项目结构

```
res_db/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块文件
├── env.example             # 环境变量示例
├── models/                 # 数据模型
│   ├── database.go         # 数据库连接
│   └── resource.go         # 资源模型
├── handlers/               # API处理器
│   ├── resource.go         # 资源相关API
│   └── category.go         # 分类相关API
├── web/                    # 前端项目
│   ├── nuxt.config.ts      # Nuxt配置
│   ├── package.json        # 前端依赖
│   ├── pages/              # 页面
│   ├── components/         # 组件
│   └── composables/        # 组合式函数
└── uploads/                # 文件上传目录
```

## 快速开始

### 1. 环境准备

确保已安装：
- Go 1.21+
- PostgreSQL 12+
- Node.js 18+

### 2. 数据库设置

```sql
CREATE DATABASE res_db;
```

### 3. 后端设置

```bash
# 复制环境变量文件
cp env.example .env

# 修改.env文件中的数据库配置

# 安装依赖
go mod tidy

# 运行后端
go run main.go
```

### 4. 前端设置

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

## API接口

### 资源管理
- `GET /api/resources` - 获取资源列表
- `POST /api/resources` - 创建资源
- `PUT /api/resources/:id` - 更新资源
- `DELETE /api/resources/:id` - 删除资源
- `GET /api/resources/:id` - 获取单个资源

### 分类管理
- `GET /api/categories` - 获取分类列表
- `POST /api/categories` - 创建分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

### 搜索和统计
- `GET /api/search` - 搜索资源
- `GET /api/stats` - 获取统计信息

## 功能特性

- 📁 资源分类管理
- 🔍 全文搜索
- 📊 统计信息
- 🏷️ 标签系统
- 📈 下载/浏览统计
- 🎨 现代化UI界面

## 开发

### 后端开发
```bash
# 热重载开发
go install github.com/cosmtrek/air@latest
air
```

### 前端开发
```bash
cd web
npm run dev
```

## 部署

### Docker部署
```bash
# 构建镜像
docker build -t res-db .

# 运行容器
docker run -p 8080:8080 res-db
```

## 许可证

MIT License 