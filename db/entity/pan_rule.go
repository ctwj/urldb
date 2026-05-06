package entity

import (
	"time"

	"gorm.io/gorm"
)

// PanRule 网盘识别规则模型
type PanRule struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PanID       uint           `json:"pan_id" gorm:"not null;comment:关联网盘ID"`
	Name        string         `json:"name" gorm:"size:64;not null;comment:规则名称"`
	Domains     string         `json:"domains" gorm:"type:text;comment:域名列表，逗号分隔"`
	URLPatterns string         `json:"url_patterns" gorm:"type:text;comment:URL正则模式，逗号分隔"`
	Priority    int            `json:"priority" gorm:"default:1;comment:匹配优先级，数字越小优先级越高"`
	Enabled     bool           `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark      string         `json:"remark" gorm:"size:255;comment:备注说明"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Pan Pan `json:"pan" gorm:"foreignKey:PanID"`
}

// TableName 指定表名
func (PanRule) TableName() string {
	return "pan_rules"
}