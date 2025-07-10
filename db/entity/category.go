package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"size:100;not null;unique;comment:分类名称"`
	Description string         `json:"description" gorm:"type:text;comment:分类描述"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Resources []Resource `json:"resources" gorm:"foreignKey:CategoryID"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
