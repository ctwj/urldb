package entity

import (
	"time"
)

// SystemConfig 系统配置实体（键值对形式）
type SystemConfig struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 键值对配置
	Key   string `json:"key" gorm:"size:100;not null;unique;comment:配置键"`
	Value string `json:"value" gorm:"type:text"`
	Type  string `json:"type" gorm:"size:20;default:'string'"` // string, int, bool, json
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}
