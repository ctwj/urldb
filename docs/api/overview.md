# API 文档概览

## 概述

老九网盘资源数据库提供了一套完整的 RESTful API 接口，支持资源管理、搜索、热门剧获取等功能。所有 API 都需要进行认证，使用 API Token 进行身份验证。

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: API Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 认证说明

### 认证方式

所有 API 都需要提供 API Token 进行认证，支持两种方式：

1. **请求头方式**（推荐）
   ```
   X-API-Token: your_token_here
   ```

2. **查询参数方式**
   ```
   ?api_token=your_token_here
   ```

### 获取 Token

请联系管理员在系统配置中设置 API Token。

## API 接口列表

### 1. 单个添加资源

**接口描述**: 添加单个资源到待处理列表

**请求信息**:
- **方法**: `POST`
- **路径**: `/api/public/resources/add`
- **认证**: 必需

**请求参数**:
```json
{
  "title": "资源标题",
  "description": "资源描述",
  "url": "资源链接",
  "category": "分类名称",
  "tags": "标签1,标签2",
  "img": "封面图片链接",
  "source": "数据来源",
  "extra": "额外信息"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "资源添加成功，已进入待处理列表",
  "data": {
    "id": 123
  },
  "code": 200
}
```

### 2. 批量添加资源

**接口描述**: 批量添加多个资源到待处理列表

**请求信息**:
- **方法**: `POST`
- **路径**: `/api/public/resources/batch-add`
- **认证**: 必需

**请求参数**:
```json
{
  "resources": [
    {
      "title": "资源1",
      "url": "链接1",
      "description": "描述1"
    },
    {
      "title": "资源2", 
      "url": "链接2",
      "description": "描述2"
    }
  ]
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "批量添加成功，共添加 2 个资源",
  "data": {
    "created_count": 2,
    "created_ids": [123, 124]
  },
  "code": 200
}
```

### 3. 资源搜索

**接口描述**: 搜索资源，支持关键词、标签、分类过滤

**请求信息**:
- **方法**: `GET`
- **路径**: `/api/public/resources/search`
- **认证**: 必需

**查询参数**:
- `keyword` - 搜索关键词
- `tag` - 标签过滤
- `category` - 分类过滤
- `page` - 页码（默认1）
- `page_size` - 每页数量（默认20，最大100）

**响应示例**:
```json
{
  "success": true,
  "message": "搜索成功",
  "data": {
    "resources": [
      {
        "id": 1,
        "title": "资源标题",
        "url": "资源链接",
        "description": "资源描述",
        "view_count": 100,
        "created_at": "2024-12-19 10:00:00",
        "updated_at": "2024-12-19 10:00:00"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  },
  "code": 200
}
```

### 4. 热门剧列表

**接口描述**: 获取热门剧列表，支持分页

**请求信息**:
- **方法**: `GET`
- **路径**: `/api/public/hot-dramas`
- **认证**: 必需

**查询参数**:
- `page` - 页码（默认1）
- `page_size` - 每页数量（默认20，最大100）

**响应示例**:
```json
{
  "success": true,
  "message": "获取热门剧成功",
  "data": {
    "hot_dramas": [
      {
        "id": 1,
        "title": "剧名",
        "description": "剧集描述",
        "img": "封面图片",
        "url": "详情链接",
        "rating": 8.5,
        "year": "2024",
        "region": "中国大陆",
        "genres": "剧情,悬疑",
        "category": "电视剧",
        "created_at": "2024-12-19 10:00:00",
        "updated_at": "2024-12-19 10:00:00"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 20
  },
  "code": 200
}
```

## 错误码说明

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败（Token无效或缺失） |
| 500 | 服务器内部错误 |
| 503 | 系统维护中或API Token未配置 |

### 响应格式

所有 API 响应都遵循统一的格式：

```json
{
  "success": true/false,
  "message": "响应消息",
  "data": {}, // 响应数据
  "code": 200 // 状态码
}
```

## 使用示例

### cURL 示例

```bash
# 设置API Token
API_TOKEN="your_api_token_here"

# 单个添加资源
curl -X POST "http://localhost:8080/api/public/resources/add" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: $API_TOKEN" \
  -d '{
    "title": "测试资源",
    "url": "https://example.com/resource",
    "description": "测试描述"
  }'

# 搜索资源
curl -X GET "http://localhost:8080/api/public/resources/search?keyword=测试" \
  -H "X-API-Token: $API_TOKEN"

# 获取热门剧
curl -X GET "http://localhost:8080/api/public/hot-dramas?page=1&page_size=5" \
  -H "X-API-Token: $API_TOKEN"
```

### JavaScript 示例

```javascript
const API_TOKEN = 'your_api_token_here';
const BASE_URL = 'http://localhost:8080/api';

// 添加资源
async function addResource(resourceData) {
  const response = await fetch(`${BASE_URL}/public/resources/add`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Token': API_TOKEN
    },
    body: JSON.stringify(resourceData)
  });
  return await response.json();
}

// 搜索资源
async function searchResources(keyword, page = 1) {
  const response = await fetch(
    `${BASE_URL}/public/resources/search?keyword=${encodeURIComponent(keyword)}&page=${page}`,
    {
      headers: {
        'X-API-Token': API_TOKEN
      }
    }
  );
  return await response.json();
}

// 获取热门剧
async function getHotDramas(page = 1, pageSize = 20) {
  const response = await fetch(
    `${BASE_URL}/public/hot-dramas?page=${page}&page_size=${pageSize}`,
    {
      headers: {
        'X-API-Token': API_TOKEN
      }
    }
  );
  return await response.json();
}
```

### Python 示例

```python
import requests

API_TOKEN = 'your_api_token_here'
BASE_URL = 'http://localhost:8080/api'

headers = {
    'X-API-Token': API_TOKEN,
    'Content-Type': 'application/json'
}

# 添加资源
def add_resource(resource_data):
    response = requests.post(
        f'{BASE_URL}/public/resources/add',
        headers=headers,
        json=resource_data
    )
    return response.json()

# 搜索资源
def search_resources(keyword, page=1):
    params = {
        'keyword': keyword,
        'page': page
    }
    response = requests.get(
        f'{BASE_URL}/public/resources/search',
        headers={'X-API-Token': API_TOKEN},
        params=params
    )
    return response.json()

# 获取热门剧
def get_hot_dramas(page=1, page_size=20):
    params = {
        'page': page,
        'page_size': page_size
    }
    response = requests.get(
        f'{BASE_URL}/public/hot-dramas',
        headers={'X-API-Token': API_TOKEN},
        params=params
    )
    return response.json()
```

## 最佳实践

### 1. 错误处理

始终检查响应的 `success` 字段和 HTTP 状态码：

```javascript
const response = await fetch(url, options);
const data = await response.json();

if (!response.ok || !data.success) {
  console.error('API调用失败:', data.message);
  // 处理错误
}
```

### 2. 分页处理

对于支持分页的接口，建议实现分页逻辑：

```javascript
async function getAllResources(keyword) {
  let allResources = [];
  let page = 1;
  let hasMore = true;
  
  while (hasMore) {
    const response = await searchResources(keyword, page);
    if (response.success) {
      allResources.push(...response.data.resources);
      hasMore = response.data.resources.length > 0;
      page++;
    } else {
      break;
    }
  }
  
  return allResources;
}
```

### 3. 请求频率限制

避免过于频繁的 API 调用，建议实现请求间隔：

```javascript
function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function searchWithDelay(keyword) {
  const result = await searchResources(keyword);
  await delay(1000); // 等待1秒
  return result;
}
```

## 注意事项

1. **Token 安全**: 请妥善保管您的 API Token，不要泄露给他人
2. **请求限制**: 避免过于频繁的请求，以免影响系统性能
3. **数据格式**: 确保请求数据格式正确，特别是 JSON 格式
4. **错误处理**: 始终实现适当的错误处理机制
5. **版本兼容**: API 可能会进行版本更新，请关注更新通知

## 技术支持

如果您在使用 API 过程中遇到问题，请：

1. 检查 API Token 是否正确
2. 确认请求格式是否符合要求
3. 查看错误响应中的详细信息
4. 联系技术支持团队

---

**注意**: 本站内容由网络爬虫自动抓取。本站不储存、复制、传播任何文件，仅作个人公益学习，请在获取后24小时内删除！ 