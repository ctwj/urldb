# 插件管理API文档

## ✅ 已实现的功能

### 📋 API接口列表

#### 1. **插件管理**
```http
# 获取插件列表
GET /api/plugins?page=1&limit=20&status=enabled&category=utility

# 获取插件详情
GET /api/plugins/{name}

# 获取插件统计信息
GET /api/plugins/stats
```

#### 2. **插件控制**
```http
# 启用插件
POST /api/plugins/{name}/enable

# 禁用插件
POST /api/plugins/{name}/disable
```

#### 3. **插件配置**
```http
# 更新插件配置
PUT /api/plugins/{name}/config
Content-Type: application/json

{
  "config": {
    "enabled": true,
    "log_level": "info",
    "max_retries": 3
  }
}
```

#### 4. **插件日志**
```http
# 获取插件日志
GET /api/plugins/{name}/logs?page=1&limit=50
```

#### 5. **插件市场（预留）**
```http
# 获取插件市场
GET /api/plugins/market

# 安装插件（预留）
POST /api/plugins/install

# 卸载插件（预留）
DELETE /api/plugins/{name}
```

## 📊 响应格式

### 成功响应
```json
{
  "success": true,
  "data": {
    // 具体数据
  }
}
```

### 错误响应
```json
{
  "success": false,
  "error": "错误信息"
}
```

## 🔧 插件信息结构

### PluginInfo
```json
{
  "id": "config_demo",
  "name": "config_demo",
  "version": "1.0.0",
  "description": "配置系统演示插件",
  "author": "URLDB Team",
  "license": "MIT",
  "category": "utility",
  "status": "installed",
  "enabled": true,
  "config": {},
  "file_size": 8697,
  "last_updated": "2024-12-25T08:04:59Z",
  "execution_stats": {
    "total_executions": 1000,
    "success_rate": 98.5,
    "average_time": 15,
    "last_execution": "2024-12-25T08:30:00Z"
  }
}
```

### ExecutionStats
```json
{
  "total_executions": 1000,
  "success_rate": 98.5,
  "average_time": 15,
  "last_execution": "2024-12-25T08:30:00Z"
}
```

## 🧪 测试示例

### 1. 获取插件列表
```bash
curl -X GET http://localhost:8080/api/plugins
```

**响应示例**:
```json
{
  "success": true,
  "data": [
    {
      "id": "config_demo",
      "name": "config_demo",
      "version": "1.0.0",
      "description": "配置系统演示插件",
      "author": "URLDB Team",
      "license": "MIT",
      "category": "utility",
      "status": "installed",
      "enabled": true,
      "file_size": 8697,
      "last_updated": "2024-12-25T08:04:59Z",
      "execution_stats": {
        "total_executions": 1000,
        "success_rate": 98.5,
        "average_time": 15
      }
    }
  ],
  "total": 1
}
```

### 2. 获取插件统计
```bash
curl -X GET http://localhost:8080/api/plugins/stats
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "total_plugins": 8,
    "enabled_plugins": 6,
    "disabled_plugins": 2,
    "total_executions": 12470,
    "success_rate": 98.2
  }
}
```

### 3. 启用插件
```bash
curl -X POST http://localhost:8080/api/plugins/config_demo/enable
```

**响应示例**:
```json
{
  "success": true,
  "message": "Plugin enabled successfully"
}
```

### 4. 更新插件配置
```bash
curl -X PUT http://localhost:8080/api/plugins/config_demo/config \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "enabled": true,
      "log_level": "debug",
      "max_retries": 5
    }
  }'
```

**响应示例**:
```json
{
  "success": true,
  "message": "Plugin config updated successfully"
}
```

## 🎯 插件元数据

插件文件开头应包含标准化的元数据：

```javascript
/// <reference path="../pb_data/types.d.ts" />

/**
 * @name config_demo
 * @version 1.0.0
 * @description 配置系统演示插件
 * @author URLDB Team
 * @license MIT
 * @category utility
 * @dependencies []
 * @permissions ["database:read", "config:manage"]
 * @hooks ["onURLAdd", "onUserLogin"]
 * @config_schema {
 *   "type": "object",
 *   "properties": {
 *     "enabled": {"type": "boolean", "default": true}
 *   }
 * }
 */
```

## 📁 数据库表结构

### plugin_configs
```sql
CREATE TABLE plugin_configs (
    id SERIAL PRIMARY KEY,
    plugin_name VARCHAR(255) UNIQUE NOT NULL,
    config_json TEXT NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### plugin_logs
```sql
CREATE TABLE plugin_logs (
    id SERIAL PRIMARY KEY,
    plugin_name VARCHAR(255) NOT NULL,
    hook_name VARCHAR(255) NOT NULL,
    execution_time INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🚀 下一步功能

1. **前端管理界面**: Vue.js组件开发
2. **插件市场**: 在线插件下载和安装
3. **版本管理**: 插件更新和回滚
4. **权限控制**: 细粒度插件权限管理
5. **实时监控**: WebSocket实时状态推送

---

**API版本**: v1.0
**最后更新**: 2024-12-25
**状态**: ✅ 核心功能已完成