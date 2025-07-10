#!/bin/bash

echo "测试新的表结构..."

# 测试创建平台
echo "1. 创建平台..."
curl -X POST http://localhost:8080/api/pans \
  -H "Content-Type: application/json" \
  -d '{
    "name": "百度网盘",
    "key": 1,
    "ck": "test_ck",
    "is_valid": true,
    "space": 2048,
    "left_space": 1024,
    "remark": "测试平台"
  }'

echo -e "\n\n2. 获取平台列表..."
curl http://localhost:8080/api/pans

echo -e "\n\n3. 创建资源（使用平台ID）..."
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

echo -e "\n\n4. 获取资源列表..."
curl http://localhost:8080/api/resources

echo -e "\n\n5. 按平台ID获取资源..."
curl "http://localhost:8080/api/resources?pan_id=1"

echo -e "\n\n测试完成！" 