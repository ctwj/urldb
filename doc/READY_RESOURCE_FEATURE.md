# 待处理资源功能说明

## 概述

新增了 `ready_resource` 表，用于批量添加和管理待处理的资源。支持两种输入格式，系统会自动识别标题和URL。

## 数据库表结构

```sql
CREATE TABLE ready_resource (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),           -- 标题（可选）
    url VARCHAR(500) NOT NULL,    -- URL（必填）
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip VARCHAR(45) DEFAULT NULL   -- 客户端IP
);
```

## 功能特性

### 🔧 **批量添加支持**

1. **单个添加**
   - 支持添加单个资源
   - 标题可选，URL必填

2. **JSON批量添加**
   - 支持JSON格式批量添加
   - 适合程序化操作

3. **文本批量添加**
   - 支持纯文本格式
   - 自动识别标题和URL
   - 支持两种格式

### 📝 **输入格式**

#### 格式1：标题和URL两行一组
```
电影标题1
https://pan.baidu.com/s/123456
电影标题2
https://pan.baidu.com/s/789012
```

#### 格式2：只有URL
```
https://pan.baidu.com/s/123456
https://pan.baidu.com/s/789012
https://pan.baidu.com/s/345678
```

### 🌐 **URL自动识别**

系统会自动识别以下类型的URL：
- HTTP/HTTPS链接
- FTP链接
- 磁力链接
- 百度网盘链接
- 阿里云盘链接
- 夸克网盘链接
- 天翼云盘链接
- 迅雷云盘链接
- 微云链接
- 蓝奏云链接
- 123云盘链接
- Google Drive链接
- Dropbox链接
- OneDrive链接
- 城通网盘链接
- 115网盘链接
- UC网盘链接

## API接口

### 1. 获取待处理资源列表
```http
GET /api/ready-resources
```

### 2. 创建单个待处理资源
```http
POST /api/ready-resources
Content-Type: application/json

{
  "title": "电影标题",
  "url": "https://pan.baidu.com/s/123456"
}
```

### 3. 批量创建待处理资源（JSON格式）
```http
POST /api/ready-resources/batch
Content-Type: application/json

{
  "resources": [
    {
      "title": "电影1",
      "url": "https://pan.baidu.com/s/111111"
    },
    {
      "title": "电影2",
      "url": "https://pan.baidu.com/s/222222"
    },
    {
      "url": "https://pan.baidu.com/s/333333"
    }
  ]
}
```

### 4. 从文本批量创建待处理资源
```http
POST /api/ready-resources/text
Content-Type: application/json

{
  "text": "电影标题1\nhttps://pan.baidu.com/s/444444\n电影标题2\nhttps://pan.baidu.com/s/555555"
}
```

### 5. 删除单个待处理资源
```http
DELETE /api/ready-resources/{id}
```

### 6. 清空所有待处理资源
```http
DELETE /api/ready-resources
```

## 前端页面

### 待处理资源管理页面 (`/ready-resources`)

功能特性：
- 📋 显示所有待处理资源
- ➕ 批量添加功能
- 🗑️ 删除单个资源
- 🗑️ 清空所有资源
- 🔄 刷新数据
- 📊 统计信息

### 管理员页面集成

在管理员页面 (`/admin`) 新增了待处理资源模块：
- 快速访问待处理资源管理
- 批量添加资源入口

## 使用示例

### 1. 通过前端页面添加

1. 访问 `/ready-resources` 页面
2. 点击"批量添加"按钮
3. 在文本框中输入资源内容
4. 点击"批量添加"提交

### 2. 通过API添加

```bash
# 添加单个资源
curl -X POST http://localhost:8080/api/ready-resources \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试电影",
    "url": "https://pan.baidu.com/s/123456"
  }'

# 批量添加（文本格式）
curl -X POST http://localhost:8080/api/ready-resources/text \
  -H "Content-Type: application/json" \
  -d '{
    "text": "电影标题1\nhttps://pan.baidu.com/s/444444\n电影标题2\nhttps://pan.baidu.com/s/555555"
  }'
```

## 工作流程

1. **批量添加** → 用户通过前端或API批量添加资源到 `ready_resource` 表
2. **自动处理** → 系统后续会自动处理这些资源（标题识别、平台判断等）
3. **正式资源** → 处理完成后移动到正式的 `resources` 表

## 技术实现

### 后端实现

1. **数据库模型** (`models/resource.go`)
   - `ReadyResource` 结构体
   - `CreateReadyResourceRequest` 请求结构
   - `BatchCreateReadyResourceRequest` 批量请求结构

2. **Handlers** (`handlers/ready_resource.go`)
   - `GetReadyResources` - 获取列表
   - `CreateReadyResource` - 创建单个
   - `BatchCreateReadyResources` - 批量创建（JSON）
   - `CreateReadyResourcesFromText` - 从文本批量创建
   - `DeleteReadyResource` - 删除单个
   - `ClearReadyResources` - 清空所有
   - `isURL` - URL识别函数

3. **路由配置** (`main.go`)
   - 添加所有待处理资源相关的API路由

### 前端实现

1. **API调用** (`composables/useApi.ts`)
   - `useReadyResourceApi` - 待处理资源API

2. **管理页面** (`pages/ready-resources.vue`)
   - 完整的待处理资源管理界面
   - 批量添加模态框
   - 资源列表显示
   - 操作功能

3. **管理员页面集成** (`pages/admin.vue`)
   - 添加待处理资源模块
   - 快速访问入口

## 测试

运行测试脚本验证功能：
```bash
./test-ready-resources.sh
```

## 注意事项

1. **URL识别**：系统会自动识别常见网盘和文件分享平台的URL
2. **标题处理**：如果只提供URL，标题字段为空，系统后续会自动处理
3. **IP记录**：自动记录添加资源的客户端IP
4. **事务处理**：批量操作使用数据库事务确保数据一致性
5. **错误处理**：完善的错误处理和用户反馈

## 后续扩展

1. **自动处理**：实现自动将待处理资源转换为正式资源
2. **标题识别**：通过URL内容自动识别资源标题
3. **平台分类**：自动识别资源所属平台
4. **重复检测**：检测重复URL避免重复添加
5. **批量操作**：支持批量删除、批量移动等功能

这个功能为网盘资源管理系统提供了灵活的批量添加机制，支持多种输入格式，为后续的自动化处理奠定了基础。 