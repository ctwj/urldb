package converter

import (
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// ToAPIAccessLogResponse 将APIAccessLog实体转换为APIAccessLogResponse
func ToAPIAccessLogResponse(log *entity.APIAccessLog) dto.APIAccessLogResponse {
	return dto.APIAccessLogResponse{
		ID:             log.ID,
		IP:             log.IP,
		UserAgent:      log.UserAgent,
		Endpoint:       log.Endpoint,
		Method:         log.Method,
		RequestParams:  log.RequestParams,
		ResponseStatus: log.ResponseStatus,
		ResponseData:   log.ResponseData,
		ProcessCount:   log.ProcessCount,
		ErrorMessage:   log.ErrorMessage,
		ProcessingTime: log.ProcessingTime,
		CreatedAt:      log.CreatedAt,
	}
}

// ToAPIAccessLogResponseList 将APIAccessLog实体列表转换为APIAccessLogResponse列表
func ToAPIAccessLogResponseList(logs []entity.APIAccessLog) []dto.APIAccessLogResponse {
	responses := make([]dto.APIAccessLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = ToAPIAccessLogResponse(&log)
	}
	return responses
}

// ToAPIAccessLogSummaryResponse 将APIAccessLogSummary实体转换为APIAccessLogSummaryResponse
func ToAPIAccessLogSummaryResponse(summary *entity.APIAccessLogSummary) dto.APIAccessLogSummaryResponse {
	return dto.APIAccessLogSummaryResponse{
		TotalRequests: summary.TotalRequests,
		TodayRequests: summary.TodayRequests,
		WeekRequests:  summary.WeekRequests,
		MonthRequests: summary.MonthRequests,
		ErrorRequests: summary.ErrorRequests,
		UniqueIPs:     summary.UniqueIPs,
	}
}

// ToAPIAccessLogStatsResponse 将APIAccessLogStats实体转换为APIAccessLogStatsResponse
func ToAPIAccessLogStatsResponse(stat entity.APIAccessLogStats) dto.APIAccessLogStatsResponse {
	return dto.APIAccessLogStatsResponse{
		Endpoint:       stat.Endpoint,
		Method:         stat.Method,
		RequestCount:   stat.RequestCount,
		ErrorCount:     stat.ErrorCount,
		AvgProcessTime: stat.AvgProcessTime,
		LastAccess:     stat.LastAccess,
	}
}

// ToAPIAccessLogStatsResponseList 将APIAccessLogStats实体列表转换为APIAccessLogStatsResponse列表
func ToAPIAccessLogStatsResponseList(stats []entity.APIAccessLogStats) []dto.APIAccessLogStatsResponse {
	responses := make([]dto.APIAccessLogStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = ToAPIAccessLogStatsResponse(stat)
	}
	return responses
}
