package entity

import (
	"time"
)

// PluginConfig 插件配置实体
type PluginConfig struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	PluginName string    `json:"plugin_name" gorm:"size:100;not null;uniqueIndex:idx_plugin_config_unique;comment:插件名称"`
	ConfigKey  string    `json:"config_key" gorm:"size:255;not null;uniqueIndex:idx_plugin_config_unique;comment:配置键"`
	ConfigValue string   `json:"config_value" gorm:"type:text;comment:配置值"`
	ConfigType string    `json:"config_type" gorm:"size:20;default:'string';comment:配置类型(string,int,bool,json)"`
	IsEncrypted bool     `json:"is_encrypted" gorm:"default:false;comment:是否加密"`
	Description string    `json:"description" gorm:"type:text;comment:配置描述"`
}

// TableName 指定表名
func (PluginConfig) TableName() string {
	return "plugin_configs"
}