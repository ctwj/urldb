# API 响应格式标准化

## 概述

为了统一API响应格式，提高前后端协作效率，所有API接口都使用标准化的响应格式。

## 标准响应格式

### 基础响应结构

```json
{
  "success": true,
  "message": "操作成功",
  "data": {},
  "error": "",
  "pagination": {}
}
```

### 字段说明

- `success`: 布尔值，表示操作是否成功
- `message`: 字符串，成功时的提示信息（可选）
- `data`: 对象/数组，返回的数据内容（可选）
- `error`: 字符串，错误信息（仅在失败时返回）
- `pagination`: 对象，分页信息（仅在分页接口时返回）

## 分页响应格式

### 分页信息结构

```json
{
  "success": true,
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 100,
    "total": 1002,
    "total_pages": 11
  }
}
```

### 分页参数

- `page`: 当前页码（从1开始）
- `page_size`: 每页条数
- `total`: 总记录数
- `total_pages`: 总页数

## 响应类型

### 1. 成功响应

```go
// 普通成功响应
SuccessResponse(c, data, "操作成功")

// 简单成功响应（无数据）
SimpleSuccessResponse(c, "操作成功")

// 创建成功响应
CreatedResponse(c, data, "创建成功")
```

### 2. 错误响应

```go
// 错误响应
ErrorResponse(c, http.StatusBadRequest, "参数错误")
```

### 3. 分页响应

```go
// 分页响应
PaginatedResponse(c, data, page, pageSize, total)
```

## 接口示例

### 获取待处理资源列表

**请求：**
```
GET /api/ready-resources?page=1&page_size=100
```

**响应：**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "示例资源",
      "url": "https://example.com",
      "create_time": "2024-01-01T00:00:00Z",
      "ip": "127.0.0.1"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 100,
    "total": 1002,
    "total_pages": 11
  }
}
```

### 创建待处理资源

**请求：**
```
POST /api/ready-resources
{
  "title": "新资源",
  "url": "https://example.com"
}
```

**响应：**
```json
{
  "success": true,
  "message": "待处理资源创建成功",
  "data": {
    "id": 1003
  }
}
```

### 错误响应示例

**响应：**
```json
{
  "success": false,
  "error": "参数错误：标题不能为空"
}
```

## 前端调用示例

### 获取分页数据

```typescript
const response = await api.getReadyResources({
  page: 1,
  page_size: 100
})

if (response.success) {
  const resources = response.data
  const pagination = response.pagination
  // 处理数据
}
```

### 处理错误

```typescript
try {
  const response = await api.createResource(data)
  if (response.success) {
    // 成功处理
  }
} catch (error) {
  // 网络错误等
}
```

## 实施规范

1. **所有新接口**必须使用标准化响应格式
2. **现有接口**逐步迁移到标准化格式
3. **错误处理**统一使用ErrorResponse
4. **分页接口**必须使用PaginatedResponse
5. **前端调用**统一处理success字段

## 迁移计划

1. ✅ 待处理资源接口（ready-resources）
2. 🔄 资源管理接口（resources）
3. 🔄 分类管理接口（categories）
4. 🔄 用户管理接口（users）
5. 🔄 统计接口（stats） 