#!/bin/bash

echo "测试搜索统计功能..."

# 测试记录搜索
echo "1. 测试记录搜索..."
curl -X POST http://localhost:8080/api/search-stats/record \
  -H "Content-Type: application/json" \
  -d '{"keyword": "测试关键词"}'

echo -e "\n"

# 测试获取搜索统计
echo "2. 测试获取搜索统计..."
curl -X GET http://localhost:8080/api/search-stats

echo -e "\n"

# 测试获取热门关键词
echo "3. 测试获取热门关键词..."
curl -X GET "http://localhost:8080/api/search-stats/hot-keywords?days=30&limit=10"

echo -e "\n"

# 测试获取每日统计
echo "4. 测试获取每日统计..."
curl -X GET "http://localhost:8080/api/search-stats/daily?days=30"

echo -e "\n"

# 测试获取搜索趋势
echo "5. 测试获取搜索趋势..."
curl -X GET "http://localhost:8080/api/search-stats/trend?days=30"

echo -e "\n"

# 测试搜索资源（会触发搜索记录）
echo "6. 测试搜索资源..."
curl -X GET "http://localhost:8080/api/search?query=测试"

echo -e "\n"

echo "搜索统计功能测试完成！" 