

<div align="center">
  
# 🚀 urlDB - 老九网盘资源数据库
  
![Go Version](https://img.shields.io/badge/Go-1230?logo=go&logoColor=white)
![Vue Version](https://img.shields.io/badge/Vue-334FC08D?logo=vue.js&logoColor=white)
![Nuxt Version](https://img.shields.io/badge/Nuxt-300.8+-00DC82?logo=nuxt.js&logoColor=white)
![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791go=postgresql&logoColor=white)

</div>
<div>

**一个现代化的网盘资源数据库，支持多网盘自动化转存分享，支持百度网盘，阿里云盘，夸克网盘， 天翼云盘，迅雷云盘，123云盘，115网盘，UC网盘 **

免费电报资源机器人:  [@L9ResBot](https://t.me/L9ResBot)  发送 搜索 + 名字 可搜索资源  





</div>

<div align="center">

  🌐 [在线演示](https://pan.l9.lc) | 📖 [文档](https://ecn5khs4t956.feishu.cn/wiki/PsnDwtxghiP0mLkTiruczKtxnwd?from=from_copylink) | 🐛 [问题反馈](https://github.com/ctwj/urldb/issues) | ⭐ [给个星标](https://github.com/ctwj/urldb)
  
### 支持的网盘平台

| 平台 | 录入 | 转存 | 分享 |
|------|-------|-----|------|
| 百度网盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |
| 阿里云盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |
| 夸克网盘 | ✅ 支持 | ✅ 支持 | ✅ 支持 |
| 天翼云盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |
| 迅雷云盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |
| UC网盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |
| 123云盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |
| 115网盘 | ✅ 支持 | 🚧 开发中 | 🚧 开发中 |

</div>

---

## 🔔 版本改动


- [文档说明](https://ecn5khs4t956.feishu.cn/wiki/PsnDwtxghiP0mLkTiruczKtxnwd?from=from_copylink)
- [服务器要求](https://ecn5khs4t956.feishu.cn/wiki/W8YBww1Mmiu4Cdkp5W4c8pFNnMf?from=from_copylink) 
- [QQ机器人](https://github.com/ctwj/astrbot_plugin_urldb) 
- [Telegram机器人](https://ecn5khs4t956.feishu.cn/wiki/SwkQw6AzRiFes7kxJXac3pd2ncb?from=from_copylink)
- [微信公众号自动回复](https://ecn5khs4t956.feishu.cn/wiki/APOEwOyDYicKGHk7gTzcQKpynkf?from=from_copylink)

### v1.3.6
1. 插件
2. bug fix

[详细改动记录](https://github.com/ctwj/urldb/blob/main/ChangeLog.md) 

当前特性
1. 支持API，手动批量录入资源
2. 支持，自动判断资源有效性
3. 支持自动转存
4. 支持平台多账号管理
5. 支持简单的数据统计
6. 支持Meilisearch


---

## 📸 项目截图


### 🏠 首页
![首页](https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/github/index.webp)

### 🔧 后台管理
![后台管理](https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/github/admin.webp)

### ⚙️ 系统配置
![系统配置](https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/github/config.webp)

### 🔍 批量转存
![资源搜索](https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/github/save.webp)

### 👤 多账号管理
![账号管理](https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/github/account.webp)

---

## ✨ 功能特性

### 🎯 核心功能
- **📁 多平台网盘支持** - 支持夸克网盘、阿里云盘、百度网盘、UC网盘
- **🔍 公开API** - 支持API数据录入，资源搜索
- **🏷️ 自动预处理** - 系统自动处理资源， 对数据进行有效性判断
- **📊 自动转存分享** - 有效资源，如果属于支持类型将自动转存分享
- **📱 多账号管理** - 同平台支持多账号管理

## 🏗️ 技术架构

### 后端技术栈
- **🦀 Golang 10.23+** - 高性能后端语言
- **🌿 Gin** - 轻量级Web框架
- **🗄️ PostgreSQL** - 关系型数据库
- **🔧 GORM** - ORM框架
- **🔐 JWT** - 身份认证
- **🔌 插件系统** - 基于事件驱动的可扩展架构

### 前端技术栈
- **⚡ Nuxt.js 3** - Vue.js全栈框架
- **🎨 Vue 3** - 渐进式JavaScript框架
- **📝 TypeScript** - 类型安全的JavaScript
- **🎨 Tailwind CSS** - 实用优先的CSS框架
- **🔧 Pinia** - 状态管理

---

## 🔧 配置说明

### 环境变量配置

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=url_db

# 服务器配置
PORT=8080

# 时区配置
TIMEZONE=Asia/Shanghai

# 日志配置
LOG_LEVEL=INFO          # 日志级别 (DEBUG, INFO, WARN, ERROR, FATAL)
DEBUG=false             # 调试模式开关
STRUCTURED_LOG=false    # 结构化日志开关 (JSON格式)
```
---

## 📄 许可证

本项目采用 [GPL License](LICENSE) 许可证。

---

## 🙏 致谢

感谢以下开源项目的支持：

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Nuxt.js](https://nuxt.com/) - Vue.js全栈框架
- [Tailwind CSS](https://tailwindcss.com/) - CSS框架
- [GORM](https://gorm.io/) - Go ORM库

---

## 📞 交流群
- **TG**: [Telegram 技术交流群](https://t.me/+QF9OMpOv-PBjZGEx)

---

<div align="center">

**如果这个项目对您有帮助，请给我们一个 ⭐ Star！**

Made with ❤️ by [老九]

</div> 
