package entity

import (
	"time"
)

// PluginData 插件数据实体
type PluginData struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	PluginName string    `json:"plugin_name" gorm:"size:100;not null;index:idx_plugin_data_plugin;comment:插件名称"`
	DataType   string    `json:"data_type" gorm:"size:100;not null;index:idx_plugin_data_type;comment:数据类型"`
	DataKey    string    `json:"data_key" gorm:"size:255;not null;index:idx_plugin_data_key;comment:数据键"`
	DataValue  string    `json:"data_value" gorm:"type:text;comment:数据值"`
	Metadata   string    `json:"metadata" gorm:"type:json;comment:元数据"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" gorm:"comment:过期时间"`
}

// TableName 指定表名
func (PluginData) TableName() string {
	return "plugin_data"
}