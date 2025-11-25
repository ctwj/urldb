package entity

import (
	"time"

	"gorm.io/gorm"
)

// TaskItemStatus 任务项状态
type TaskItemStatus string

const (
	TaskItemStatusPending    TaskItemStatus = "pending"    // 等待处理
	TaskItemStatusProcessing TaskItemStatus = "processing" // 处理中
	TaskItemStatusSuccess    TaskItemStatus = "success"    // 成功
	TaskItemStatusFailed     TaskItemStatus = "failed"     // 失败
	TaskItemStatusSkipped    TaskItemStatus = "skipped"    // 跳过
)

// TaskItem 任务项表（任务的详细记录）
type TaskItem struct {
	ID     uint `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID uint `json:"task_id" gorm:"not null;index;comment:任务ID"`

	// 通用任务项信息
	Status       TaskItemStatus `json:"status" gorm:"size:20;not null;default:pending;comment:处理状态"`
	ErrorMessage string         `json:"error_message" gorm:"type:text;comment:错误信息"`

	// 输入数据 (JSON格式存储，支持不同任务类型的不同数据结构)
	InputData string `json:"input_data" gorm:"type:text;not null;comment:输入数据(JSON格式)"`

	// 输出数据 (JSON格式存储，支持不同任务类型的不同结果数据)
	OutputData string `json:"output_data" gorm:"type:text;comment:输出数据(JSON格式)"`

	// 处理日志 (可选，用于记录详细的处理过程)
	ProcessLog string `json:"process_log" gorm:"type:text;comment:处理日志"`

	// Google索引特有字段 (当任务类型为google_index时使用)
	URL            string     `json:"url" gorm:"size:2048;comment:URL (Google索引专用)"`
	IndexStatus    string     `json:"index_status" gorm:"size:50;comment:索引状态 (Google索引专用)"`
	InspectResult  string     `json:"inspect_result" gorm:"type:text;comment:检查结果 (Google索引专用)"`
	MobileFriendly bool       `json:"mobile_friendly" gorm:"default:false;comment:是否移动友好 (Google索引专用)"`
	LastCrawled    *time.Time `json:"last_crawled" gorm:"comment:最后抓取时间 (Google索引专用)"`
	StatusCode     int        `json:"status_code" gorm:"default:0;comment:HTTP状态码 (Google索引专用)"`

	// 时间信息
	ProcessedAt *time.Time     `json:"processed_at" gorm:"comment:处理时间"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Task Task `json:"task" gorm:"foreignKey:TaskID"`
}

// TableName 指定表名
func (TaskItem) TableName() string {
	return "task_items"
}
