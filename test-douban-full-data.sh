#!/bin/bash

echo "=== 测试豆瓣服务获取全部数据功能 ==="
echo "时间: $(date)"
echo

# 测试电影榜单获取全部数据
echo "1. 测试电影榜单获取全部数据..."
curl -s "http://localhost:8080/api/hot-dramas/movies?category=热门&type=全部&start=0&limit=0" | jq '.'
echo

# 测试电视剧榜单获取全部数据
echo "2. 测试电视剧榜单获取全部数据..."
curl -s "http://localhost:8080/api/hot-dramas/tv?category=tv&type=tv&start=0&limit=0" | jq '.'
echo

# 测试调度器处理全部数据
echo "3. 测试调度器处理全部数据..."
echo "检查日志中的数据处理情况..."
echo "应该看到类似 '检测到limit=0，将尝试获取全部数据' 的日志"
echo

echo "=== 测试完成 ==="
echo "如果看到 '检测到总数为: XXX，将一次性获取全部数据' 的日志，说明功能正常" 