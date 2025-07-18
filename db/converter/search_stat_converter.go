package converter

import (
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// ToSearchStatResponse 将SearchStat实体转换为SearchStatResponse
func ToSearchStatResponse(stat *entity.SearchStat) dto.SearchStatResponse {
	return dto.SearchStatResponse{
		ID:        stat.ID,
		Keyword:   stat.Keyword,
		Count:     stat.Count,
		Date:      stat.Date,
		CreatedAt: stat.CreatedAt,
	}
}

// ToSearchStatResponseList 将SearchStat实体列表转换为SearchStatResponse列表
func ToSearchStatResponseList(stats []entity.SearchStat) []dto.SearchStatResponse {
	responses := make([]dto.SearchStatResponse, len(stats))
	for i, stat := range stats {
		responses[i] = ToSearchStatResponse(&stat)
	}
	return responses
}

// ToDailySearchStatResponse 将DailySearchStat实体转换为DailySearchStatResponse
func ToDailySearchStatResponse(stat entity.DailySearchStat) dto.DailySearchStatResponse {
	return dto.DailySearchStatResponse{
		Date:           stat.Date,
		TotalSearches:  stat.TotalSearches,
		UniqueKeywords: stat.UniqueKeywords,
	}
}

// ToDailySearchStatResponseList 将DailySearchStat实体列表转换为DailySearchStatResponse列表
func ToDailySearchStatResponseList(stats []entity.DailySearchStat) []dto.DailySearchStatResponse {
	responses := make([]dto.DailySearchStatResponse, len(stats))
	for i, stat := range stats {
		responses[i] = ToDailySearchStatResponse(stat)
	}
	return responses
}

// ToHotKeywordResponse 将KeywordStat实体转换为HotKeywordResponse
func ToHotKeywordResponse(stat entity.KeywordStat) dto.HotKeywordResponse {
	return dto.HotKeywordResponse{
		Keyword: stat.Keyword,
		Count:   stat.Count,
		Rank:    stat.Rank,
	}
}

// ToHotKeywordResponseList 将KeywordStat实体列表转换为HotKeywordResponse列表
func ToHotKeywordResponseList(stats []entity.KeywordStat) []dto.HotKeywordResponse {
	responses := make([]dto.HotKeywordResponse, len(stats))
	for i, stat := range stats {
		responses[i] = ToHotKeywordResponse(stat)
	}
	return responses
}
