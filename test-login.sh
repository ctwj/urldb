#!/bin/bash

echo "测试登录功能..."

# 测试默认管理员账户登录
echo "测试默认管理员账户 (admin/password):"
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

echo -e "\n\n"

# 测试注册新用户
echo "测试注册新用户:"
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456","email":"test@example.com"}'

echo -e "\n\n"

# 测试新用户登录
echo "测试新用户登录:"
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'

echo -e "\n\n" 