package entity

import (
	"time"
)

// TelegramChannel Telegram 频道/群组实体
type TelegramChannel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Telegram 频道/群组信息
	ChatID   int64  `json:"chat_id" gorm:"not null;comment:Telegram 聊天ID"`
	ChatName string `json:"chat_name" gorm:"size:255;not null;comment:聊天名称"`
	ChatType string `json:"chat_type" gorm:"size:50;not null;comment:类型：channel/group"`

	// 推送配置
	PushEnabled       bool   `json:"push_enabled" gorm:"default:true;comment:是否启用推送"`
	PushFrequency     int    `json:"push_frequency" gorm:"default:5;comment:推送频率（分钟）"`
	PushStartTime     string `json:"push_start_time" gorm:"size:10;comment:推送开始时间，格式HH:mm"`
	PushEndTime       string `json:"push_end_time" gorm:"size:10;comment:推送结束时间，格式HH:mm"`
	ContentCategories string `json:"content_categories" gorm:"type:text;comment:推送的内容分类，用逗号分隔"`
	ContentTags       string `json:"content_tags" gorm:"type:text;comment:推送的标签，用逗号分隔"`

	// 频道状态
	IsActive   bool       `json:"is_active" gorm:"default:true;comment:是否活跃"`
	LastPushAt *time.Time `json:"last_push_at" gorm:"comment:最后推送时间"`

	// 注册信息
	RegisteredBy string    `json:"registered_by" gorm:"size:100;comment:注册者用户名"`
	RegisteredAt time.Time `json:"registered_at"`

	// API配置
	API             string `json:"api" gorm:"size:255;comment:API地址"`
	Token           string `json:"token" gorm:"size:255;comment:访问令牌"`
	ApiType         string `json:"api_type" gorm:"size:50;comment:API类型"`
	IsPushSavedInfo bool   `json:"is_push_saved_info" gorm:"default:false;comment:是否只推送已转存资源"`
}

// TableName 指定表名
func (TelegramChannel) TableName() string {
	return "telegram_channels"
}
