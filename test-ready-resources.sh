#!/bin/bash

echo "测试待处理资源功能..."

# 测试1: 创建单个待处理资源
echo "1. 创建单个待处理资源..."
curl -X POST http://localhost:8080/api/ready-resources \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试电影",
    "url": "https://pan.baidu.com/s/123456"
  }'

echo -e "\n\n2. 批量创建待处理资源（JSON格式）..."
curl -X POST http://localhost:8080/api/ready-resources/batch \
  -H "Content-Type: application/json" \
  -d '{
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
  }'

echo -e "\n\n3. 从文本批量创建待处理资源（格式1：标题+URL）..."
curl -X POST http://localhost:8080/api/ready-resources/text \
  -H "Content-Type: application/json" \
  -d '{
    "text": "电影标题1\nhttps://pan.baidu.com/s/444444\n电影标题2\nhttps://pan.baidu.com/s/555555"
  }'

echo -e "\n\n4. 从文本批量创建待处理资源（格式2：只有URL）..."
curl -X POST http://localhost:8080/api/ready-resources/text \
  -H "Content-Type: application/json" \
  -d '{
    "text": "https://pan.baidu.com/s/666666\nhttps://pan.baidu.com/s/777777\nhttps://pan.baidu.com/s/888888"
  }'

echo -e "\n\n5. 获取所有待处理资源..."
curl http://localhost:8080/api/ready-resources

echo -e "\n\n6. 删除一个待处理资源（需要替换ID）..."
# 注意：这里需要替换实际的ID
curl -X DELETE http://localhost:8080/api/ready-resources/1

echo -e "\n\n7. 清空所有待处理资源..."
curl -X DELETE http://localhost:8080/api/ready-resources

echo -e "\n\n测试完成！" 