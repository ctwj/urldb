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
	Ck        string         `json:"ck" gorm:"type:text;comment:cookie"`
	IsValid   bool           `json:"is_valid" gorm:"default:true;comment:是否有效"`
	Space     int64          `json:"space" gorm:"default:0;comment:总空间"`
	LeftSpace int64          `json:"left_space" gorm:"default:0;comment:剩余空间"`
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
