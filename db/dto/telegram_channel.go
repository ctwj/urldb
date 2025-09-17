package dto

import "time"

// TelegramChannelRequest 创建/更新 Telegram 频道/群组请求
type TelegramChannelRequest struct {
	ChatID            int64  `json:"chat_id" binding:"required"`
	ChatName          string `json:"chat_name" binding:"required"`
	ChatType          string `json:"chat_type" binding:"required"` // channel 或 group
	PushEnabled       bool   `json:"push_enabled"`
	PushFrequency     int    `json:"push_frequency"`
	PushStartTime     string `json:"push_start_time"`
	PushEndTime       string `json:"push_end_time"`
	ContentCategories string `json:"content_categories"`
	ContentTags       string `json:"content_tags"`
	IsActive          bool   `json:"is_active"`
}

// TelegramChannelResponse Telegram 频道/群组响应
type TelegramChannelResponse struct {
	ID                uint       `json:"id"`
	ChatID            int64      `json:"chat_id"`
	ChatName          string     `json:"chat_name"`
	ChatType          string     `json:"chat_type"`
	PushEnabled       bool       `json:"push_enabled"`
	PushFrequency     int        `json:"push_frequency"`
	PushStartTime     string     `json:"push_start_time"`
	PushEndTime       string     `json:"push_end_time"`
	ContentCategories string     `json:"content_categories"`
	ContentTags       string     `json:"content_tags"`
	IsActive          bool       `json:"is_active"`
	LastPushAt        *time.Time `json:"last_push_at"`
	RegisteredBy      string     `json:"registered_by"`
	RegisteredAt      time.Time  `json:"registered_at"`
}

// TelegramBotConfigRequest Telegram 机器人配置请求
type TelegramBotConfigRequest struct {
	BotEnabled         *bool   `json:"bot_enabled"`
	BotApiKey          *string `json:"bot_api_key"`
	AutoReplyEnabled   *bool   `json:"auto_reply_enabled"`
	AutoReplyTemplate  *string `json:"auto_reply_template"`
	AutoDeleteEnabled  *bool   `json:"auto_delete_enabled"`
	AutoDeleteInterval *int    `json:"auto_delete_interval"`
}

// TelegramBotConfigResponse Telegram 机器人配置响应
type TelegramBotConfigResponse struct {
	BotEnabled         bool   `json:"bot_enabled"`
	BotApiKey          string `json:"bot_api_key"`
	AutoReplyEnabled   bool   `json:"auto_reply_enabled"`
	AutoReplyTemplate  string `json:"auto_reply_template"`
	AutoDeleteEnabled  bool   `json:"auto_delete_enabled"`
	AutoDeleteInterval int    `json:"auto_delete_interval"`
}

// ValidateTelegramApiKeyRequest 验证 Telegram API Key 请求
type ValidateTelegramApiKeyRequest struct {
	ApiKey        string `json:"api_key" binding:"required"`
	ProxyEnabled  bool   `json:"proxy_enabled"`
	ProxyType     string `json:"proxy_type"`
	ProxyHost     string `json:"proxy_host"`
	ProxyPort     int    `json:"proxy_port"`
	ProxyUsername string `json:"proxy_username"`
	ProxyPassword string `json:"proxy_password"`
}

// ValidateTelegramApiKeyResponse 验证 Telegram API Key 响应
type ValidateTelegramApiKeyResponse struct {
	Valid   bool                   `json:"valid"`
	Error   string                 `json:"error,omitempty"`
	BotInfo map[string]interface{} `json:"bot_info,omitempty"`
}
