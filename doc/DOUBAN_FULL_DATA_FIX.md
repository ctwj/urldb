# 豆瓣服务获取全部数据功能优化

## 问题描述

用户反馈豆瓣接口返回了 `total: 283`，但当前代码只获取了部分数据，希望能够一次性获取全部数据而不是分页获取。

## 解决方案

### 1. 修改豆瓣服务逻辑

在 `utils/douban_service.go` 中：

- 为 `GetMovieRanking` 和 `GetTvRanking` 方法添加了 `limit=0` 的特殊处理
- 当 `limit=0` 时，先获取第一页来确定总数，然后一次性获取全部数据
- 重构了代码，将实际的API调用逻辑提取到 `getMovieRankingPage` 和 `getTvRankingPage` 方法中

### 2. 更新调度器配置

在 `utils/scheduler.go` 中：

- 将电影数据处理从 `limit=20` 改为 `limit=0`
- 将电视剧数据处理从 `limit=20` 改为 `limit=0`
- 这样调度器会自动获取并处理全部数据

### 3. 功能特点

- **智能检测**: 当传入 `limit=0` 时，自动检测总数并获取全部数据
- **向后兼容**: 原有的分页功能保持不变
- **日志记录**: 详细记录获取过程，便于调试
- **错误处理**: 如果获取总数失败，会回退到默认行为

## 使用方法

### API调用

```bash
# 获取全部电影数据
curl "http://localhost:8080/api/hot-dramas/movies?category=热门&type=全部&start=0&limit=0"

# 获取全部电视剧数据
curl "http://localhost:8080/api/hot-dramas/tv?category=tv&type=tv&start=0&limit=0"
```

### 调度器自动处理

调度器现在会自动获取全部数据：

```go
// 电影数据处理
movieResult, err := s.doubanService.GetMovieRanking("热门", "全部", 0, 0)

// 电视剧数据处理
tvResult, err := s.doubanService.GetTvRanking("tv", "tv", 0, 0)
```

## 日志输出

当使用 `limit=0` 时，会看到类似以下的日志：

```
=== 开始获取电影榜单 ===
参数: category=热门, rankingType=全部, start=0, limit=0
检测到limit=0，将尝试获取全部数据
检测到总数为: 283，将一次性获取全部数据
```

## 测试验证

使用提供的测试脚本：

```bash
chmod +x test-douban-full-data.sh
./test-douban-full-data.sh
```

## 优势

1. **数据完整性**: 确保获取到所有可用数据
2. **性能优化**: 减少API调用次数
3. **灵活性**: 支持分页和全量获取两种模式
4. **可维护性**: 代码结构清晰，易于理解和维护

## 注意事项

1. 获取全部数据可能会增加单次请求的响应时间
2. 建议在调度器中使用此功能，避免影响用户界面的响应速度
3. 如果API返回的数据量很大，需要考虑内存使用情况 