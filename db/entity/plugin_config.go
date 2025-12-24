package entity

import (
	"time"
)

// PluginConfig 插件配置表
type PluginConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PluginName  string    `gorm:"uniqueIndex;not null" json:"plugin_name"`
	ConfigJSON  string    `gorm:"type:text;not null" json:"config_json"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// PluginLog 插件执行日志表
type PluginLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PluginName   string    `gorm:"not null;index" json:"plugin_name"`
	HookName     string    `gorm:"not null" json:"hook_name"`
	ExecutionTime int      `gorm:"not null" json:"execution_time"` // 毫秒
	Success      bool      `gorm:"not null" json:"success"`
	ErrorMessage *string   `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// CustomEvent 自定义事件表
type CustomEvent struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	EventName string     `gorm:"not null;index" json:"event_name"`
	EventData string     `gorm:"type:text;not null" json:"event_data"`
	Processed bool       `gorm:"default:false;index" json:"processed"`
	CreatedAt time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`
}

// CronJob 定时任务表
type CronJob struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	JobName      string     `gorm:"uniqueIndex;not null" json:"job_name"`
	Schedule     string     `gorm:"not null" json:"schedule"`
	HandlerScript string    `gorm:"type:text;not null" json:"handler_script"`
	Enabled      bool       `gorm:"default:true" json:"enabled"`
	LastRun      *time.Time `json:"last_run"`
	NextRun      *time.Time `gorm:"index" json:"next_run"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserPreference 用户偏好设置表
type UserPreference struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	UserID            string     `gorm:"not null;index" json:"user_id"`
	Category          string     `gorm:"not null" json:"category"`
	PreferencesJSON   string     `gorm:"type:text" json:"preferences_json"`
	NotificationEmail bool       `gorm:"default:true" json:"notification_email"`
	NotificationPush  bool       `gorm:"default:true" json:"notification_push"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserInterest 用户兴趣标签表
type UserInterest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"not null;index" json:"user_id"`
	Category    string    `gorm:"not null" json:"category"`
	Score       int       `gorm:"default:1" json:"score"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`
}

// URLStats URL统计表
type URLStats struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URLID     string    `gorm:"not null;index" json:"url_id"`
	Domain    string    `gorm:"not null;index" json:"domain"`
	Category  string    `gorm:"index" json:"category"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// AccessStats 访问统计表
type AccessStats struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	URLID       string     `gorm:"not null;index" json:"url_id"`
	UserID      *string    `gorm:"index" json:"user_id"`
	AccessTime  time.Time  `gorm:"not null;index" json:"access_time"`
	Referrer    string     `gorm:"type:text" json:"referrer"`
	UserAgent   string     `gorm:"type:text" json:"user_agent"`
	IPAddress   string     `gorm:"index" json:"ip_address"`
}

// PopularResources 热门资源表
type PopularResources struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	URLID        string     `gorm:"uniqueIndex;not null" json:"url_id"`
	AccessCount  int        `gorm:"default:1" json:"access_count"`
	LastAccessed time.Time  `gorm:"autoUpdateTime" json:"last_accessed"`
}

// DomainPattern 域名模式表
type DomainPattern struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Domain    string    `gorm:"uniqueIndex;not null" json:"domain"`
	Count     int       `gorm:"default:1" json:"count"`
	LastSeen  time.Time `gorm:"autoUpdateTime" json:"last_seen"`
}

// PathPattern 路径模式表
type PathPattern struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Pattern   string    `gorm:"uniqueIndex;not null" json:"pattern"`
	Count     int       `gorm:"default:1" json:"count"`
	LastSeen  time.Time `gorm:"autoUpdateTime" json:"last_seen"`
}

// RealTimeMetrics 实时指标表
type RealTimeMetrics struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Metric    string    `gorm:"not null;index" json:"metric"`
	Value     float64   `gorm:"not null" json:"value"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
}

// ClassificationStats 分类统计表
type ClassificationStats struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Category    string     `gorm:"uniqueIndex;not null" json:"category"`
	TagCount    int        `gorm:"default:0" json:"tag_count"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// DailyReport 每日报告表
type DailyReport struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `gorm:"not null;index" json:"type"`
	ReportData string    `gorm:"type:text;not null" json:"report_data"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// SystemHealth 系统健康检查表
type SystemHealth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    string    `gorm:"not null" json:"status"`
	Metrics   string    `gorm:"type:text" json:"metrics"`
	CheckedAt time.Time `gorm:"autoCreateTime" json:"checked_at"`
}