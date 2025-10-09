package entity

import (
	"time"

	"gorm.io/gorm"
)

// APIAccessLog API访问日志模型
type APIAccessLog struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	IP             string         `json:"ip" gorm:"size:45;not null;comment:客户端IP地址"`
	UserAgent      string         `json:"user_agent" gorm:"size:500;comment:用户代理"`
	Endpoint       string         `json:"endpoint" gorm:"size:255;not null;comment:访问的接口路径"`
	Method         string         `json:"method" gorm:"size:10;not null;comment:HTTP方法"`
	RequestParams  string         `json:"request_params" gorm:"type:text;comment:查询参数(JSON格式)"`
	ResponseStatus int            `json:"response_status" gorm:"default:200;comment:响应状态码"`
	ResponseData   string         `json:"response_data" gorm:"type:text;comment:响应数据(JSON格式)"`
	ProcessCount   int            `json:"process_count" gorm:"default:0;comment:处理数量(查询结果数或添加的数量)"`
	ErrorMessage   string         `json:"error_message" gorm:"size:500;comment:错误消息"`
	ProcessingTime int64          `json:"processing_time" gorm:"comment:处理时间(毫秒)"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (APIAccessLog) TableName() string {
	return "api_access_logs"
}

// APIAccessLogSummary API访问日志汇总统计
type APIAccessLogSummary struct {
	TotalRequests int64 `json:"total_requests"`
	TodayRequests int64 `json:"today_requests"`
	WeekRequests  int64 `json:"week_requests"`
	MonthRequests int64 `json:"month_requests"`
	ErrorRequests int64 `json:"error_requests"`
	UniqueIPs     int64 `json:"unique_ips"`
}

// APIAccessLogStats 按端点统计
type APIAccessLogStats struct {
	Endpoint       string    `json:"endpoint"`
	Method         string    `json:"method"`
	RequestCount   int64     `json:"request_count"`
	ErrorCount     int64     `json:"error_count"`
	AvgProcessTime int64     `json:"avg_process_time"`
	LastAccess     time.Time `json:"last_access"`
}
