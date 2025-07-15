#!/bin/bash

echo "=== 测试原始豆瓣API ==="
echo "时间: $(date)"
echo

# 设置请求头
HEADERS=(
  "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"
  "Referer: https://m.douban.com/"
  "Accept: application/json, text/plain, */*"
  "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8"
  "Accept-Encoding: gzip, deflate"
  "Connection: keep-alive"
  "Sec-Fetch-Dest: empty"
  "Sec-Fetch-Mode: cors"
  "Sec-Fetch-Site: same-origin"
)

# 构建请求头字符串
HEADER_STR=""
for header in "${HEADERS[@]}"; do
  HEADER_STR="$HEADER_STR -H \"$header\""
done

# 测试电影榜单API
echo "1. 测试电影榜单API（获取前10条）..."
MOVIE_URL="https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=10&category=热门&type=全部"
echo "请求URL: $MOVIE_URL"
echo

eval "curl -s $HEADER_STR \"$MOVIE_URL\" | jq '.'"
echo

# 测试电视剧榜单API
echo "2. 测试电视剧榜单API（获取前10条）..."
TV_URL="https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=10&category=tv&type=tv"
echo "请求URL: $TV_URL"
echo

eval "curl -s $HEADER_STR \"$TV_URL\" | jq '.'"
echo

# 测试获取总数
echo "3. 测试获取电影总数..."
TOTAL_URL="https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=1&category=热门&type=全部"
echo "请求URL: $TOTAL_URL"
echo

eval "curl -s $HEADER_STR \"$TOTAL_URL\" | jq '.total'"
echo

echo "=== 测试完成 ==="
echo "如果看到JSON响应数据，说明API配置正确"
echo "如果返回403错误，请检查请求头配置" 