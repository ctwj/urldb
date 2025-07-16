package entity

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"size:100;not null;unique;comment:标签名称"`
	Description string         `json:"description" gorm:"type:text;comment:标签描述"`
	CategoryID  *uint          `json:"category_id" gorm:"comment:分类ID，可以为空"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Category  Category   `json:"category" gorm:"foreignKey:CategoryID"`
	Resources []Resource `json:"resources" gorm:"many2many:resource_tags;"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
