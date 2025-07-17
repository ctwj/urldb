# 🚀 panResManage - 网盘资源管理系统

<div align="center>

![Go Version](https://img.shields.io/badge/Go-1230?logo=go&logoColor=white)
![Vue Version](https://img.shields.io/badge/Vue-334FC08D?logo=vue.js&logoColor=white)
![Nuxt Version](https://img.shields.io/badge/Nuxt-300.8+-00DC82?logo=nuxt.js&logoColor=white)
![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791go=postgresql&logoColor=white)

**一个现代化的网盘资源管理系统，支持多网盘自动化转存分享**

🌐 [在线演示](#) | 📖 [文档](#) | 🐛 问题反馈](#) | ⭐ [给个星标](#)

</div>

---

## 🔔 温馨提示

📌 **本项目仅供技术交流与学习使用**，自身不存储或提供任何资源文件及下载链接。

📌 **请勿将本项目用于任何违法用途**，否则后果自负。

📌 如有任何问题或建议，欢迎交流探讨！ 😊

> **免责声明**：本项目由 Trae AI 辅助编写。由于时间有限，仅在空闲时维护。如遇使用问题，请优先自行排查，感谢理解！

---

## ✨ 功能特性

### 🎯 核心功能
- **📁 多平台网盘支持** - 支持夸克网盘、阿里云盘、百度网盘、UC网盘
- **🔍 公开API** - 支持API数据录入，资源搜索
- **🏷️ 自动预处理** - 系统自动处理资源， 对数据进行有效性判断
- **📊 自动转存分享** - 有效资源，如果属于支持类型将自动转存分享
- **📱 多账号管理** - 同平台支持多账号管理

### 🛠️ 管理功能
- **📦 批量操作** - 批量添加、导入、管理资源
- **🔄 自动处理** - 待处理资源自动转存和分享
- **📈 热播剧管理** - 热门影视资源自动更新
- **⚙️ 系统配置** - 灵活的系统参数配置

### 🎨 用户体验
- **📱 响应式设计** - 支持桌面端和移动端
- **🌙 深色模式** - 支持明暗主题切换
- **⚡ 高性能** - 基于Go的高并发后端
- **🎯 现代化UI** - 基于Tailwind CSS的美观界面

---

## 🏗️ 技术架构

### 后端技术栈
- **🦀 Golang 10.23+** - 高性能后端语言
- **🌿 Gin** - 轻量级Web框架
- **🗄️ PostgreSQL** - 关系型数据库
- **🔧 GORM** - ORM框架
- **🔐 JWT** - 身份认证

### 前端技术栈
- **⚡ Nuxt.js 3** - Vue.js全栈框架
- **🎨 Vue 3** - 渐进式JavaScript框架
- **📝 TypeScript** - 类型安全的JavaScript
- **🎨 Tailwind CSS** - 实用优先的CSS框架
- **🔧 Pinia** - 状态管理

### 开发工具
- **🐳 Docker** - 容器化部署
- **📦 pnpm** - 快速包管理器
- **🔍 Air** - Go热重载工具

---

## 🚀 快速开始

### 环境要求

- **Docker** 和 **Docker Compose**
- 或者本地环境：
  - **Go** 1.23+
  - **Node.js** 18+
  - **PostgreSQL** 15+
  - **pnpm** (推荐) 或 npm

### 方式一：Docker 部署（推荐）

#### 使用启动脚本（最简单）
```bash
# 克隆项目
git clone https://github.com/ctwj/panResManage.git
cd panResManage

# 使用启动脚本
./docker-start.sh
```

#### 手动启动
```bash
# 克隆项目
git clone https://github.com/ctwj/panResManage.git
cd panResManage

# 使用 Docker Compose 启动
docker compose up --build -d

# 访问应用
# 前端: http://localhost:3000
# 后端API: http://localhost:8080
```

### 方式二：本地开发

#### 1. 克隆项目
```bash
git clone https://github.com/ctwj/panResManage.git
cd panResManage
```

#### 2. 后端设置
```bash
# 复制环境变量文件
cp env.example .env

# 编辑环境变量
vim .env

# 安装Go依赖
go mod tidy

# 启动后端服务
go run main.go
```

#### 3. 前端设置
```bash
# 进入前端目录
cd web

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

#### 4. 数据库设置
```sql
-- 创建数据库
CREATE DATABASE res_db;
```

---

## 📁 项目结构

```
l9pan/
├── 📁 common/                 # 通用功能模块
│   ├── 📄 pan_factory.go     # 网盘工厂模式
│   ├── 📄 alipan.go          # 阿里云盘实现
│   ├── 📄 baidu_pan.go       # 百度网盘实现
│   ├── 📄 quark_pan.go       # 夸克网盘实现
│   └── 📄 uc_pan.go          # UC网盘实现
├── 📁 db/                     # 数据库层
│   ├── 📁 entity/            # 数据实体
│   ├── 📁 repo/              # 数据仓库
│   ├── 📁 dto/               # 数据传输对象
│   └── 📁 converter/         # 数据转换器
├── 📁 handlers/               # API处理器
├── 📁 middleware/             # 中间件
├── 📁 utils/                  # 工具函数
├── 📁 web/                    # 前端项目
│   ├── 📁 pages/             # 页面组件
│   ├── 📁 components/        # 通用组件
│   ├── 📁 composables/       # 组合式函数
│   └── 📁 stores/            # 状态管理
├── 📁 docs/                   # 项目文档
├── 📄 main.go                # 主程序入口
├── 📄 Dockerfile             # Docker配置
├── 📄 docker-compose.yml     # Docker Compose配置
└── 📄 README.md              # 项目说明
```

---

## 🔧 配置说明

### 环境变量配置

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=res_db

# 服务器配置
PORT=8080
```

### Docker 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| frontend | 3000 | Nuxt.js 前端应用 |
| backend | 8080 | Go API 后端服务 |
| postgres | 5432 | PostgreSQL 数据库 |

### 支持的网盘平台

| 平台 | 状态 | 功能 |
|------|------|------|
| 夸克网盘 | ✅ 支持 | 转存、分享 |
| 阿里云盘 | 🚧 开发中 | 转存、分享 |
| 百度网盘 | 🚧 开发中 | 转存、分享 |
| UC网盘 | 🚧 开发中 | 转存、分享 |

---

## 📚 API 文档

### 公开统计

提供，批量入库和搜索api，通过 apiToken 授权

> 📖 完整API文档请访问：`http://p.l9.lc/doc.html`

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1**Fork** 本仓库2 **创建** 功能分支 (`git checkout -b feature/AmazingFeature`)
3** 更改 (`git commit -mAdd some AmazingFeature'`)
4. **推送** 到分支 (`git push origin feature/AmazingFeature`)
5. **创建** Pull Request


## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

```
MIT License

Copyright (c) 2024 L9Pan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 致谢

感谢以下开源项目的支持：

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Nuxt.js](https://nuxt.com/) - Vue.js全栈框架
- [Tailwind CSS](https://tailwindcss.com/) - CSS框架
- [GORM](https://gorm.io/) - Go ORM库

---

## 📞 联系我们

- **项目地址**: [https://github.com/ctwj/panResManage](https://github.com/ctwj/panResManage)
- **问题反馈**: [Issues](https://github.com/ctwj/panResManage/issues)
- **邮箱**: 510199617@qq.com

---

<div align="center>

**如果这个项目对您有帮助，请给我们一个 ⭐ Star！**

Made with ❤️ by [老九]

</div> 