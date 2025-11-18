package entity

import (
	"gorm.io/gorm"
	"time"
)

// Report 举报实体
type Report struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ResourceKey string         `gorm:"type:varchar(255);not null;index" json:"resource_key"` // 资源唯一标识
	Reason      string         `gorm:"type:varchar(100);not null" json:"reason"`             // 举报原因
	Description string         `gorm:"type:text" json:"description"`                         // 详细描述
	Contact     string         `gorm:"type:varchar(255)" json:"contact"`                     // 联系方式
	UserAgent   string         `gorm:"type:text" json:"user_agent"`                          // 用户代理
	IPAddress   string         `gorm:"type:varchar(45)" json:"ip_address"`                   // IP地址
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"`     // 处理状态: pending, approved, rejected
	ProcessedAt *time.Time     `json:"processed_at"`                                         // 处理时间
	ProcessedBy *uint          `json:"processed_by"`                                         // 处理人ID
	Note        string         `gorm:"type:text" json:"note"`                                // 处理备注
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

// TableName 表名
func (Report) TableName() string {
	return "reports"
}