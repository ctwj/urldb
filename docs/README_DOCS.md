# 文档使用说明

## 概述

本项目使用 [docsify](https://docsify.js.org/) 生成文档网站。docsify 是一个轻量级的文档生成器，无需构建静态文件，只需要一个 `index.html` 文件即可。

## 文档结构

```
docs/
├── index.html              # 文档主页
├── docsify.config.js       # docsify 配置文件
├── README.md               # 首页内容
├── _sidebar.md             # 侧边栏导航
├── start-docs.sh           # 启动脚本
├── guide/                  # 使用指南
│   ├── quick-start.md      # 快速开始
│   ├── local-development.md # 本地开发
│   └── docker-deployment.md # Docker 部署
├── api/                    # API 文档
│   └── overview.md         # API 概览
├── architecture/           # 架构文档
│   └── overview.md         # 架构概览
├── faq.md                  # 常见问题
├── changelog.md            # 更新日志
└── license.md              # 许可证
```

## 快速启动

### 方法一：使用启动脚本（推荐）

```bash
# 进入文档目录
cd docs

# 运行启动脚本
./start-docs.sh
```

脚本会自动：
- 检查是否安装了 docsify-cli
- 如果没有安装，会自动安装
- 启动文档服务
- 在浏览器中打开文档

### 方法二：手动启动

```bash
# 安装 docsify-cli（如果未安装）
npm install -g docsify-cli

# 进入文档目录
cd docs

# 启动服务
docsify serve . --port 3000 --open
```

## 访问文档

启动成功后，文档将在以下地址可用：
- 本地访问：http://localhost:3000
- 局域网访问：http://[你的IP]:3000

## 文档特性

### 1. 搜索功能
- 支持全文搜索
- 搜索结果高亮显示
- 支持中文搜索

### 2. 代码高亮
支持多种编程语言的语法高亮：
- Go
- JavaScript/TypeScript
- SQL
- YAML
- JSON
- Bash

### 3. 代码复制
- 一键复制代码块
- 复制成功提示

### 4. 页面导航
- 侧边栏导航
- 页面间导航
- 自动回到顶部

### 5. 响应式设计
- 支持移动端访问
- 自适应屏幕尺寸

## 自定义配置

### 修改主题
在 `docsify.config.js` 中修改配置：

```javascript
window.$docsify = {
  name: '你的项目名称',
  repo: '你的仓库地址',
  // 其他配置...
}
```

### 添加新页面
1. 在相应目录下创建 `.md` 文件
2. 在 `_sidebar.md` 中添加导航链接
3. 刷新页面即可看到新页面

### 修改样式
可以通过添加自定义 CSS 来修改样式：

```html
<!-- 在 index.html 中添加 -->
<link rel="stylesheet" href="./custom.css">
```

## 部署到生产环境

### 静态部署
docsify 生成的文档可以部署到任何静态文件服务器：

```bash
# 构建静态文件（可选）
docsify generate docs docs/_site

# 部署到 GitHub Pages
git subtree push --prefix docs origin gh-pages
```

### Docker 部署
```bash
# 使用 nginx 镜像
docker run -d -p 80:80 -v $(pwd)/docs:/usr/share/nginx/html nginx
```

## 常见问题

### Q: 启动时提示端口被占用
A: 可以指定其他端口：
```bash
docsify serve . --port 3001
```

### Q: 搜索功能不工作
A: 确保在 `index.html` 中引入了搜索插件：
```html
<script src="//cdn.jsdelivr.net/npm/docsify@4/lib/plugins/search.min.js"></script>
```

### Q: 代码高亮不显示
A: 确保引入了相应的 Prism.js 组件：
```html
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-go.min.js"></script>
```

## 维护说明

### 更新文档
1. 修改相应的 `.md` 文件
2. 刷新浏览器即可看到更新

### 添加新功能
1. 在 `docsify.config.js` 中添加插件配置
2. 在 `index.html` 中引入相应的插件文件

### 版本控制
建议将文档与代码一起进行版本控制，确保文档与代码版本同步。

## 相关链接

- [docsify 官方文档](https://docsify.js.org/)
- [docsify 插件市场](https://docsify.js.org/#/plugins)
- [Markdown 语法指南](https://docsify.js.org/#/zh-cn/markdown) 