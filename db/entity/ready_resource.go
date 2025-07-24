package entity

import (
	"time"

	"gorm.io/gorm"
)

// ReadyResource 待处理资源模型
type ReadyResource struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       *string        `json:"title" gorm:"size:255;comment:资源标题"`
	Description string         `json:"description" gorm:"type:text;comment:资源描述"`
	URL         string         `json:"url" gorm:"size:500;not null;comment:资源链接"`
	Category    string         `json:"category" gorm:"size:100;comment:资源分类"`
	Tags        string         `json:"tags" gorm:"size:500;comment:资源标签，多个标签用逗号分隔"`
	Img         string         `json:"img" gorm:"size:500;comment:封面链接"`
	Source      string         `json:"source" gorm:"size:100;comment:数据来源"`
	Extra       string         `json:"extra" gorm:"type:text;comment:额外附加数据"`
	Key         string         `json:"key" gorm:"size:64;index;comment:资源组标识，相同key表示同一组资源"`
	CreateTime  time.Time      `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
	IP          *string        `json:"ip" gorm:"size:45;comment:IP地址"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (ReadyResource) TableName() string {
	return "ready_resource"
}
