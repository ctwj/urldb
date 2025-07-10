package entity

import (
	"time"
)

// ResourceTag 资源标签关联表
type ResourceTag struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceID uint      `json:"resource_id" gorm:"not null;comment:资源ID"`
	TagID      uint      `json:"tag_id" gorm:"not null;comment:标签ID"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (ResourceTag) TableName() string {
	return "resource_tags"
}
