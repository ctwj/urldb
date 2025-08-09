package entity

import (
	"time"

	"gorm.io/gorm"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   // 等待中
	TaskStatusRunning   TaskStatus = "running"   // 运行中
	TaskStatusPaused    TaskStatus = "paused"    // 已暂停
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 失败
	TaskStatusCancelled TaskStatus = "cancelled" // 已取消
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeBatchTransfer TaskType = "batch_transfer" // 批量转存
)

// Task 任务表
type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"size:255;not null;comment:任务标题"`
	Type        TaskType   `json:"type" gorm:"size:50;not null;comment:任务类型"`
	Status      TaskStatus `json:"status" gorm:"size:20;not null;default:pending;comment:任务状态"`
	Description string     `json:"description" gorm:"type:text;comment:任务描述"`

	// 进度信息
	TotalItems     int `json:"total_items" gorm:"not null;default:0;comment:总项目数"`
	ProcessedItems int `json:"processed_items" gorm:"not null;default:0;comment:已处理项目数"`
	SuccessItems   int `json:"success_items" gorm:"not null;default:0;comment:成功项目数"`
	FailedItems    int `json:"failed_items" gorm:"not null;default:0;comment:失败项目数"`

	// 任务配置 (JSON格式存储)
	Config string `json:"config" gorm:"type:text;comment:任务配置"`

	// 时间信息
	StartedAt   *time.Time     `json:"started_at" gorm:"comment:开始时间"`
	CompletedAt *time.Time     `json:"completed_at" gorm:"comment:完成时间"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	TaskItems []TaskItem `json:"task_items" gorm:"foreignKey:TaskID"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}
