package entity

import (
	"time"

	"gorm.io/gorm"
)

// API 申请单状态（016-api-access-application）
const (
	ApiApplicationStatusPending  = "pending"  // 待审核
	ApiApplicationStatusApproved = "approved" // 已同意（已开通）
	ApiApplicationStatusRejected = "rejected" // 已拒绝
)

// API 申请单类型
const (
	ApiApplicationTypeApply   = "apply"   // 开通申请（默认）
	ApiApplicationTypeRenewal = "renewal" // 续期申请（1-6 个月；凭证未过期自动通过）
)

// ApiApplication API 开通申请单（用户中心「API 访问」）
type ApiApplication struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint           `json:"user_id" gorm:"not null;index;comment:申请人ID"`
	Email          string         `json:"email" gorm:"size:100;not null;comment:申请时邮箱快照"`
	Purpose        string         `json:"purpose" gorm:"type:text;not null;comment:用途说明"`
	Type           string         `json:"type" gorm:"size:20;not null;default:'apply';comment:类型:apply开通/renewal续期"`
	DurationMonths int            `json:"duration_months" gorm:"default:0;comment:续期月数(1-6,仅续期申请)"`
	Status         string         `json:"status" gorm:"size:20;not null;default:pending;index;comment:状态:pending/approved/rejected"`
	RejectReason   string         `json:"reject_reason" gorm:"type:text;comment:拒绝理由(可选)"`
	ReviewerID     *uint          `json:"reviewer_id" gorm:"comment:审核人ID"`
	ReviewerName   string         `json:"reviewer_name" gorm:"size:50;comment:审核人名快照"`
	ReviewedAt     *time.Time     `json:"reviewed_at" gorm:"comment:审核时间"`
	CreatedAt      time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (ApiApplication) TableName() string {
	return "api_applications"
}
