# 用户认证系统

## 概述

本项目已成功集成了完整的用户认证系统，包括用户注册、登录、权限管理等功能。

## 功能特性

### 1. 用户管理
- **用户注册**: 支持新用户注册
- **用户登录**: JWT令牌认证
- **用户管理**: 管理员可以创建、编辑、删除用户
- **角色管理**: 支持用户(user)和管理员(admin)角色
- **状态管理**: 支持用户激活/禁用状态

### 2. 认证机制
- **JWT令牌**: 使用JWT进行身份验证
- **密码加密**: 使用bcrypt进行密码哈希
- **中间件保护**: 路由级别的权限控制

### 3. 权限控制
- **公开接口**: 资源查看、搜索等
- **用户接口**: 个人资料查看
- **管理员接口**: 所有管理功能需要管理员权限

## 数据库结构

### users表
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## API接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `GET /api/auth/profile` - 获取用户信息

### 用户管理接口（管理员）
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

## 前端页面

### 登录页面 (`/login`)
- 用户名/密码登录
- 默认管理员账户: admin / password
- 登录成功后跳转到管理页面

### 注册页面 (`/register`)
- 新用户注册
- 注册成功后跳转到登录页面

### 管理页面 (`/admin`)
- 需要登录才能访问
- 显示用户信息和退出登录按钮
- 各种管理功能的入口

### 用户管理页面 (`/users`)
- 仅管理员可访问
- 用户列表展示
- 创建、编辑、删除用户功能

## 中间件

### AuthMiddleware
- 验证JWT令牌
- 将用户信息存储到上下文中

### AdminMiddleware
- 检查用户角色是否为管理员
- 保护管理员专用接口

## 默认数据

系统启动时会自动创建默认管理员账户：
- 用户名: admin
- 密码: password
- 角色: admin
- 邮箱: admin@example.com

## 安全特性

1. **密码加密**: 使用bcrypt进行密码哈希
2. **JWT令牌**: 24小时有效期的JWT令牌
3. **角色权限**: 基于角色的访问控制
4. **输入验证**: 服务器端数据验证
5. **SQL注入防护**: 使用GORM进行参数化查询

## 使用说明

1. 启动服务器后，访问 `/login` 页面
2. 使用默认管理员账户登录: admin / password
3. 登录成功后可以访问管理功能
4. 在用户管理页面可以创建新用户或修改现有用户

## 技术栈

- **后端**: Go + Gin + GORM + JWT
- **前端**: Nuxt.js + Tailwind CSS
- **数据库**: PostgreSQL
- **认证**: JWT + bcrypt 