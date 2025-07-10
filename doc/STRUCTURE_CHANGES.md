# 表结构变更说明

## 概述

根据需求，我们对 `resources` 表进行了重大结构调整，以支持一个资源对应一个 URL 的模式，并添加了平台标识功能。

## 主要变更

### 1. Resources 表结构变更

**变更前：**
- `url` 字段：JSON 格式，存储多个链接
- 一个资源记录包含多个链接

**变更后：**
- `url` 字段：VARCHAR(128)，存储单个链接
- 新增 `pan_id` 字段：INTEGER，关联到 `pan` 表，标识链接类型
- 一个资源对应一个链接，多个链接需要创建多条记录

### 2. 新增 Pan 表

```sql
CREATE TABLE pan (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) DEFAULT NULL,
    key INTEGER DEFAULT NULL,
    ck TEXT,
    is_valid BOOLEAN DEFAULT true,
    space BIGINT DEFAULT 0,
    left_space BIGINT DEFAULT 0,
    remark VARCHAR(64) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. API 变更

#### 资源相关 API

**新增查询参数：**
- `pan_id`：按平台ID筛选资源

**请求体变更：**
```json
{
  "title": "资源标题",
  "description": "资源描述",
  "url": "https://pan.baidu.com/s/123456",
  "pan_id": 1,  // 新增：平台ID
  "quark_url": "",
  "file_size": "100MB",
  "category_id": 1,
  "is_valid": true,
  "is_public": true,
  "tag_ids": [1, 2]
}
```

#### 平台相关 API

**新增 API 端点：**
- `GET /api/pans` - 获取平台列表
- `GET /api/pans/:id` - 获取单个平台
- `POST /api/pans` - 创建平台
- `PUT /api/pans/:id` - 更新平台
- `DELETE /api/pans/:id` - 删除平台

## 前端变更

### useApi.ts 新增方法

```typescript
// 按平台ID获取资源
const getResourcesByPan = async (panId: number, params?: any) => {
  return await $fetch('/resources', {
    baseURL: config.public.apiBase,
    params: { ...params, pan_id: panId }
  })
}
```

## 数据迁移建议

1. **备份现有数据**
   ```sql
   -- 备份现有资源表
   CREATE TABLE resources_backup AS SELECT * FROM resources;
   ```

2. **创建新的表结构**
   - 运行更新后的数据库初始化代码

3. **数据迁移策略**
   - 对于每个资源的多个链接，创建多条记录
   - 为每条记录分配适当的 `pan_id`
   - 保持原有的标签关联

## 使用示例

### 创建平台
```bash
curl -X POST http://localhost:8080/api/pans \
  -H "Content-Type: application/json" \
  -d '{
    "name": "百度网盘",
    "key": 1,
    "ck": "test_ck",
    "is_valid": true,
    "space": 2048,
    "left_space": 1024,
    "remark": "百度网盘平台"
  }'
```

### 创建资源
```bash
curl -X POST http://localhost:8080/api/resources \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试资源",
    "description": "这是一个测试资源",
    "url": "https://pan.baidu.com/s/123456",
    "pan_id": 1,
    "quark_url": "",
    "file_size": "100MB",
    "category_id": 1,
    "is_valid": true,
    "is_public": true,
    "tag_ids": [1]
  }'
```

### 按平台查询资源
```bash
curl "http://localhost:8080/api/resources?pan_id=1"
```

## 注意事项

1. **URL 长度限制**：URL 字段现在限制为 128 字符
2. **平台关联**：每个资源必须关联到一个平台（pan_id）
3. **数据完整性**：确保在创建资源时提供有效的 pan_id
4. **向后兼容**：API 响应格式保持兼容，只是新增了 pan_id 字段

## 测试

运行测试脚本验证新功能：
```bash
./test-setup.sh
``` 