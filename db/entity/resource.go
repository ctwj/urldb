package entity

import (
	"time"

	"gorm.io/gorm"
)

// Resource 资源模型
type Resource struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"size:255;not null;comment:资源标题"`
	Description string         `json:"description" gorm:"type:text;comment:资源描述"`
	URL         string         `json:"url" gorm:"size:128;comment:资源链接"`
	PanID       *uint          `json:"pan_id" gorm:"comment:平台ID"`
	QuarkURL    string         `json:"quark_url" gorm:"size:500;comment:夸克链接"`
	FileSize    string         `json:"file_size" gorm:"size:100;comment:文件大小"`
	CategoryID  *uint          `json:"category_id" gorm:"comment:分类ID"`
	ViewCount   int            `json:"view_count" gorm:"default:0;comment:浏览次数"`
	IsValid     bool           `json:"is_valid" gorm:"default:true;comment:是否有效"`
	IsPublic    bool           `json:"is_public" gorm:"default:true;comment:是否公开"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	Pan      Pan      `json:"pan" gorm:"foreignKey:PanID"`
	Tags     []Tag    `json:"tags" gorm:"many2many:resource_tags;"`
}

// TableName 指定表名
func (Resource) TableName() string {
	return "resources"
}
