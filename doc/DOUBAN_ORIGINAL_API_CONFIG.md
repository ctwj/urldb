# 原始豆瓣API Postman配置指南

## 基础信息

根据代码分析，原始豆瓣API的基础配置如下：

### 基础URL
```
https://m.douban.com/rexxar/api/v2
```

### 请求头配置
```
User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1
Referer: https://m.douban.com/
Accept: application/json, text/plain, */*
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
Accept-Encoding: gzip, deflate
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
```

## Postman配置步骤

### 1. 创建新的Collection

1. 打开Postman
2. 点击 "New" → "Collection"
3. 命名为 "豆瓣API"

### 2. 设置Collection级别的Headers

1. 选择刚创建的Collection
2. 点击 "Variables" 标签
3. 添加以下Headers：

| Key | Value |
|-----|-------|
| User-Agent | Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1 |
| Referer | https://m.douban.com/ |
| Accept | application/json, text/plain, */* |
| Accept-Language | zh-CN,zh;q=0.9,en;q=0.8 |
| Accept-Encoding | gzip, deflate |
| Connection | keep-alive |
| Sec-Fetch-Dest | empty |
| Sec-Fetch-Mode | cors |
| Sec-Fetch-Site | same-origin |

### 3. 创建电影榜单请求

#### 请求配置
- **Method**: GET
- **URL**: `https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie`

#### Query Parameters
| Key | Value | Description |
|-----|-------|-------------|
| start | 0 | 起始位置 |
| limit | 20 | 限制数量（0表示获取全部） |
| category | 热门 | 分类（热门/最新/豆瓣高分/冷门佳片） |
| type | 全部 | 类型（全部/华语/欧美/韩国/日本） |

#### 示例请求
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=20&category=热门&type=全部
```

### 4. 创建电视剧榜单请求

#### 请求配置
- **Method**: GET
- **URL**: `https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv`

#### Query Parameters
| Key | Value | Description |
|-----|-------|-------------|
| start | 0 | 起始位置 |
| limit | 20 | 限制数量（0表示获取全部） |
| category | tv | 分类（tv/show） |
| type | tv | 类型（tv/tv_domestic/tv_american/tv_japanese/tv_korean/tv_animation/tv_documentary） |

#### 示例请求
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=20&category=tv&type=tv
```

## 常用请求示例

### 1. 获取全部电影数据
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=0&category=热门&type=全部
```

### 2. 获取华语电影
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=20&category=热门&type=华语
```

### 3. 获取欧美电影
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=20&category=热门&type=欧美
```

### 4. 获取豆瓣高分电影
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=20&category=豆瓣高分&type=全部
```

### 5. 获取国产剧
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=20&category=tv&type=tv_domestic
```

### 6. 获取综艺节目
```
GET https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=20&category=show&type=show
```

## 响应格式

### 成功响应示例
```json
{
  "items": [
    {
      "id": "123456",
      "title": "电影标题",
      "card_subtitle": "导演 / 主演",
      "episodes_info": "",
      "is_new": false,
      "pic": {
        "large": "https://img1.doubanio.com/view/photo/s_ratio_poster/public/p123456.jpg",
        "normal": "https://img1.doubanio.com/view/photo/s_ratio_poster/public/p123456.jpg"
      },
      "rating": {
        "value": 8.5,
        "count": 12345,
        "max": 10,
        "star_count": 4.25
      },
      "type": "movie",
      "uri": "douban://douban.com/movie/123456",
      "year": "2024",
      "directors": ["导演名"],
      "actors": ["演员1", "演员2"],
      "region": "中国大陆",
      "genres": ["剧情", "动作"]
    }
  ],
  "total": 283,
  "categories": [
    {
      "category": "热门",
      "selected": true,
      "type": "全部",
      "title": "热门"
    }
  ]
}
```

## 注意事项

1. **请求频率**: 建议控制请求频率，避免被限制
2. **User-Agent**: 必须使用移动端User-Agent，否则可能返回错误
3. **Referer**: 必须设置正确的Referer头
4. **分页**: 使用start和limit参数进行分页
5. **全量获取**: 设置limit=0可以获取全部数据

## 错误处理

### 常见错误
- **403 Forbidden**: 请求头配置不正确
- **404 Not Found**: URL路径错误
- **429 Too Many Requests**: 请求频率过高

### 调试建议
1. 检查所有请求头是否正确设置
2. 确认URL路径和参数格式
3. 使用浏览器开发者工具对比请求
4. 查看响应状态码和错误信息 