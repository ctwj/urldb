#!/bin/bash

echo "=== 热播剧功能测试 ==="

# 测试获取热播剧列表
echo "1. 测试获取热播剧列表..."
curl -X GET "http://localhost:8080/api/hot-dramas" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n2. 测试按分类获取热播剧..."
curl -X GET "http://localhost:8080/api/hot-dramas?category=电影" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n3. 测试获取调度器状态..."
curl -X GET "http://localhost:8080/api/scheduler/status" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n4. 测试手动获取热播剧名字..."
curl -X GET "http://localhost:8080/api/scheduler/hot-drama/names" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  | jq '.'

echo -e "\n5. 测试启动热播剧定时任务..."
curl -X POST "http://localhost:8080/api/scheduler/hot-drama/start" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  | jq '.'

echo -e "\n=== 测试完成 ===" 