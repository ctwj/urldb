package entity

import (
	"time"

	"gorm.io/gorm"
)

// 用户资源状态（015-user-resource-upload 五态机）
const (
	UserResourceStatusPending    = "pending"    // 未检测（检测通道未启用/超时/未得出结论）
	UserResourceStatusValid      = "valid"      // 有效（检测通过，待入队处理）
	UserResourceStatusInvalid    = "invalid"    // 无效（检测未通过或处理失败，可重检）
	UserResourceStatusProcessing = "processing" // 处理中（已入 ReadyResource 队列）
	UserResourceStatusPublished  = "published"  // 已公开（自动处理通道复检通过并入库）
)

// 检测方式（对齐 LinkCheckService.DetectionMethod）
const (
	UserResourceCheckMethodPanCheck = "pancheck"
	UserResourceCheckMethodDisabled = "disabled"
)

// UserResource 用户上传资源（用户中心「我的资源」）
type UserResource struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint      `json:"user_id" gorm:"not null;index:idx_user_resources_user_status,priority:1;index:idx_user_resources_user_url,priority:1;comment:提交者ID"`
	Title       string    `json:"title" gorm:"size:255;not null;comment:资源标题"`
	Description string    `json:"description" gorm:"type:text;comment:资源描述"`
	URL         string    `json:"url" gorm:"size:500;not null;index:idx_user_resources_user_url,priority:2;comment:网盘分享链接"`
	PanID       *uint     `json:"pan_id" gorm:"comment:平台ID"`
	Status      string    `json:"status" gorm:"size:20;not null;default:pending;index:idx_user_resources_user_status,priority:2;comment:状态:pending/valid/invalid/processing/published"`
	CheckTime   *time.Time `json:"check_time" gorm:"comment:最近检测时间"`
	CheckMethod string    `json:"check_method" gorm:"size:20;comment:检测方式:pancheck/disabled"`
	PublishResourceID *uint `json:"publish_resource_id" gorm:"index;comment:公开库资源ID(归属反查)"`
	FailReason  string    `json:"fail_reason" gorm:"type:text;comment:无效/处理失败原因"`
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Pan *Pan `json:"pan" gorm:"foreignKey:PanID"`
}

// TableName 指定表名
func (UserResource) TableName() string {
	return "user_resources"
}
