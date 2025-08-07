package entity

import (
	"time"
	"gorm.io/gorm"
)

// ResourceView 资源访问记录
type ResourceView struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceID uint           `json:"resource_id" gorm:"not null;index;comment:资源ID"`
	IPAddress  string         `json:"ip_address" gorm:"size:45;comment:访问者IP地址"`
	UserAgent  string         `json:"user_agent" gorm:"type:text;comment:用户代理"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime;comment:访问时间"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// 关联关系
	Resource Resource `json:"resource" gorm:"foreignKey:ResourceID"`
}

// TableName 指定表名
func (ResourceView) TableName() string {
	return "resource_views"
} 