# 搜索统计功能

## 功能概述

搜索统计功能用于记录和分析用户的搜索行为，包括：
- 每日搜索量统计
- 热门关键词分析
- 搜索趋势图表
- 关键词热度排名

## 数据库设计

### 搜索统计表 (search_stats)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | SERIAL | 主键 |
| keyword | VARCHAR(255) | 搜索关键词 |
| count | INTEGER | 搜索次数 |
| date | DATE | 搜索日期 |
| ip | VARCHAR(45) | 用户IP |
| user_agent | VARCHAR(500) | 用户代理 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

## API接口

### 1. 记录搜索
```
POST /api/search-stats/record
Content-Type: application/json

{
  "keyword": "搜索关键词"
}
```

### 2. 获取搜索统计总览
```
GET /api/search-stats
```

返回数据：
```json
{
  "today_searches": 10,
  "week_searches": 50,
  "month_searches": 200,
  "hot_keywords": [
    {
      "keyword": "关键词",
      "count": 15,
      "rank": 1
    }
  ],
  "daily_stats": [...],
  "search_trend": {
    "days": ["01-01", "01-02"],
    "values": [10, 15]
  }
}
```

### 3. 获取热门关键词
```
GET /api/search-stats/hot-keywords?days=30&limit=10
```

参数：
- days: 统计天数（默认30）
- limit: 返回数量限制（默认10）

### 4. 获取每日统计
```
GET /api/search-stats/daily?days=30
```

### 5. 获取搜索趋势
```
GET /api/search-stats/trend?days=30
```

### 6. 获取关键词趋势
```
GET /api/search-stats/keyword/{keyword}/trend?days=30
```

## 前端页面

### 搜索统计页面 (/search-stats)

功能特性：
- 今日/本周/本月搜索量统计卡片
- 搜索趋势折线图
- 热门关键词排行榜
- 关键词热度可视化

## 自动记录

系统会在以下情况下自动记录搜索：
- 用户使用搜索功能时
- 记录用户IP和User-Agent信息
- 按日期聚合统计

## 使用说明

1. **管理员访问**：登录后可在管理员页面看到"搜索统计"模块
2. **查看统计**：点击"查看搜索统计"进入详细页面
3. **热门关键词**：查看最受欢迎的关键词排名
4. **趋势分析**：通过图表了解搜索量变化趋势

## 测试

运行测试脚本：
```bash
chmod +x test-search-stats.sh
./test-search-stats.sh
```

## 技术实现

### 后端
- 使用Repository模式管理数据访问
- 支持按日期聚合统计
- 提供多种统计维度API

### 前端
- 使用Chart.js绘制趋势图表
- 响应式设计适配不同设备
- 实时数据更新

## 注意事项

1. 搜索记录会保存用户IP，注意隐私保护
2. 大量搜索数据可能影响性能，建议定期清理
3. 关键词统计按天聚合，避免重复记录
4. 图表数据需要足够的历史数据才能显示趋势 