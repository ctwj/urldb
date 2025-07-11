# 热播剧功能说明

## 功能概述

热播剧功能是一个自动获取和展示豆瓣热门电影、电视剧榜单的功能模块。系统会定时从豆瓣获取最新的热门影视作品信息，并保存到数据库中供用户浏览。

## 主要特性

### 1. 自动数据获取
- 每小时自动从豆瓣获取热门电影和电视剧数据
- 支持电影和电视剧两个分类
- 获取内容包括：剧名、评分、年份、导演、演员等详细信息

### 2. 数据存储
- 创建专门的热播剧数据表 `hot_dramas`
- 支持按豆瓣ID去重，避免重复数据
- 记录数据来源和获取时间

### 3. 前端展示
- 美观的卡片式布局展示热播剧信息
- 支持按分类筛选（全部/电影/电视剧）
- 分页显示，支持大量数据
- 响应式设计，适配各种设备

### 4. 管理功能
- 管理员可以手动启动/停止定时任务
- 支持手动获取热播剧数据
- 查看调度器运行状态

## 数据库结构

### hot_dramas 表
```sql
CREATE TABLE hot_dramas (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    rating DECIMAL(3,1) DEFAULT 0.0,
    year VARCHAR(10),
    directors VARCHAR(500),
    actors VARCHAR(1000),
    category VARCHAR(50),
    sub_type VARCHAR(50),
    source VARCHAR(50) DEFAULT 'douban',
    douban_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API 接口

### 热播剧管理
- `GET /api/hot-dramas` - 获取热播剧列表
- `GET /api/hot-dramas/:id` - 获取热播剧详情
- `POST /api/hot-dramas` - 创建热播剧记录（管理员）
- `PUT /api/hot-dramas/:id` - 更新热播剧记录（管理员）
- `DELETE /api/hot-dramas/:id` - 删除热播剧记录（管理员）

### 调度器管理
- `GET /api/scheduler/status` - 获取调度器状态
- `POST /api/scheduler/hot-drama/start` - 启动热播剧定时任务（管理员）
- `POST /api/scheduler/hot-drama/stop` - 停止热播剧定时任务（管理员）
- `GET /api/scheduler/hot-drama/names` - 手动获取热播剧名字（管理员）

## 配置说明

### 系统配置
在系统配置中有一个 `auto_fetch_hot_drama_enabled` 字段，用于控制是否启用自动获取热播剧功能：

- `true`: 启用自动获取，系统会根据配置的间隔时间自动获取数据
- `false`: 禁用自动获取，需要管理员手动启动

## 使用流程

### 1. 启用功能
1. 登录管理后台
2. 进入系统配置页面
3. 开启"自动拉取热播剧名字"选项
4. 保存配置

### 2. 查看热播剧
1. 在首页点击"热播剧"按钮
2. 进入热播剧页面
3. 可以按分类筛选查看
4. 支持分页浏览

### 3. 管理定时任务
1. 管理员可以手动启动/停止定时任务
2. 可以查看调度器运行状态
3. 可以手动触发数据获取

## 技术实现

### 后端架构
- **实体层**: `db/entity/hot_drama.go` - 定义热播剧数据结构
- **DTO层**: `db/dto/hot_drama.go` - 定义数据传输对象
- **转换器**: `db/converter/hot_drama_converter.go` - 实体与DTO转换
- **仓储层**: `db/repo/hot_drama_repository.go` - 数据库操作
- **处理器**: `handlers/hot_drama_handler.go` - API接口处理
- **调度器**: `utils/scheduler.go` - 定时任务管理
- **豆瓣服务**: `utils/douban_service.go` - 豆瓣API调用

### 前端实现
- **页面**: `web/pages/hot-dramas.vue` - 热播剧展示页面
- **导航**: 在首页添加热播剧入口
- **样式**: 使用Tailwind CSS实现响应式设计

## 注意事项

1. **数据来源**: 数据来源于豆瓣移动端API，如果API不可用会使用模拟数据
2. **频率限制**: 定时任务每小时执行一次，避免对豆瓣服务器造成压力
3. **数据去重**: 系统会根据豆瓣ID进行去重，避免重复数据
4. **权限控制**: 管理功能需要管理员权限
5. **错误处理**: 系统具备完善的错误处理机制，确保稳定性

## 扩展功能

未来可以考虑添加的功能：
1. 支持更多数据源（如IMDB、烂番茄等）
2. 添加用户收藏功能
3. 支持热播剧搜索
4. 添加数据统计和分析功能
5. 支持热播剧推荐算法 