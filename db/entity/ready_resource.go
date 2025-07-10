package entity

import (
	"time"

	"gorm.io/gorm"
)

// ReadyResource 待处理资源模型
type ReadyResource struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      *string        `json:"title" gorm:"size:255;comment:资源标题"`
	URL        string         `json:"url" gorm:"size:500;not null;comment:资源链接"`
	CreateTime time.Time      `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
	IP         *string        `json:"ip" gorm:"size:45;comment:IP地址"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (ReadyResource) TableName() string {
	return "ready_resource"
}
