package entity

import (
	"time"

	"gorm.io/gorm"
)

// Pan 第三方平台表
type Pan struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"size:64;comment:平台名称"`
	Key       int            `json:"key" gorm:"comment:平台标识"`
	Icon      string         `json:"icon" gorm:"size:128;comment:图标文字"`
	Remark    string         `json:"remark" gorm:"size:64;not null;comment:备注"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Cks []Cks `json:"cks" gorm:"foreignKey:PanID"`
}

// TableName 指定表名
func (Pan) TableName() string {
	return "pan"
}
