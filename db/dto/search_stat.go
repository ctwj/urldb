package dto

import "time"

// SearchStatRequest 搜索统计请求
type SearchStatRequest struct {
	Keyword string `json:"keyword" binding:"required"`
}

// SearchStatResponse 搜索统计响应
type SearchStatResponse struct {
	ID        uint      `json:"id"`
	Keyword   string    `json:"keyword"`
	Count     int       `json:"count"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// DailySearchStatResponse 每日搜索统计响应
type DailySearchStatResponse struct {
	Date           time.Time `json:"date"`
	TotalSearches  int       `json:"total_searches"`
	UniqueKeywords int       `json:"unique_keywords"`
}

// HotKeywordResponse 热门关键词响应
type HotKeywordResponse struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
	Rank    int    `json:"rank"`
}

// SearchTrendResponse 搜索趋势响应
type SearchTrendResponse struct {
	Days   []string `json:"days"`
	Values []int    `json:"values"`
}

// SearchStatsResponse 搜索统计总览响应
type SearchStatsResponse struct {
	TodaySearches int                       `json:"today_searches"`
	WeekSearches  int                       `json:"week_searches"`
	MonthSearches int                       `json:"month_searches"`
	HotKeywords   []HotKeywordResponse      `json:"hot_keywords"`
	DailyStats    []DailySearchStatResponse `json:"daily_stats"`
	SearchTrend   SearchTrendResponse       `json:"search_trend"`
}
