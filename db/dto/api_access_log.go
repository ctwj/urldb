package dto

import "time"

// APIAccessLogResponse API访问日志响应
type APIAccessLogResponse struct {
	ID             uint      `json:"id"`
	IP             string    `json:"ip"`
	UserAgent      string    `json:"user_agent"`
	Endpoint       string    `json:"endpoint"`
	Method         string    `json:"method"`
	RequestParams  string    `json:"request_params"`
	ResponseStatus int       `json:"response_status"`
	ResponseData   string    `json:"response_data"`
	ProcessCount   int       `json:"process_count"`
	ErrorMessage   string    `json:"error_message"`
	ProcessingTime int64     `json:"processing_time"`
	CreatedAt      time.Time `json:"created_at"`
}

// APIAccessLogSummaryResponse API访问日志汇总响应
type APIAccessLogSummaryResponse struct {
	TotalRequests int64 `json:"total_requests"`
	TodayRequests int64 `json:"today_requests"`
	WeekRequests  int64 `json:"week_requests"`
	MonthRequests int64 `json:"month_requests"`
	ErrorRequests int64 `json:"error_requests"`
	UniqueIPs     int64 `json:"unique_ips"`
}

// APIAccessLogStatsResponse 按端点统计响应
type APIAccessLogStatsResponse struct {
	Endpoint       string    `json:"endpoint"`
	Method         string    `json:"method"`
	RequestCount   int64     `json:"request_count"`
	ErrorCount     int64     `json:"error_count"`
	AvgProcessTime int64     `json:"avg_process_time"`
	LastAccess     time.Time `json:"last_access"`
}

// APIAccessLogListResponse API访问日志列表响应
type APIAccessLogListResponse struct {
	Data  []APIAccessLogResponse `json:"data"`
	Total int64                  `json:"total"`
}

// APIAccessLogFilterRequest API访问日志过滤请求
type APIAccessLogFilterRequest struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	IP        string `json:"ip,omitempty"`
	Page      int    `json:"page,omitempty" default:"1"`
	PageSize  int    `json:"page_size,omitempty" default:"20"`
}
