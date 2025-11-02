package dto

// WechatBotConfigRequest 微信公众号机器人配置请求
type WechatBotConfigRequest struct {
	Enabled         bool   `json:"enabled"`
	AppID           string `json:"app_id"`
	AppSecret       string `json:"app_secret"`
	Token           string `json:"token"`
	EncodingAesKey  string `json:"encoding_aes_key"`
	WelcomeMessage  string `json:"welcome_message"`
	AutoReplyEnabled bool   `json:"auto_reply_enabled"`
	SearchLimit     int    `json:"search_limit"`
}

// WechatBotConfigResponse 微信公众号机器人配置响应
type WechatBotConfigResponse struct {
	Enabled         bool   `json:"enabled"`
	AppID           string `json:"app_id"`
	AppSecret       string `json:"app_secret"`
	Token           string `json:"token"`
	EncodingAesKey  string `json:"encoding_aes_key"`
	WelcomeMessage  string `json:"welcome_message"`
	AutoReplyEnabled bool   `json:"auto_reply_enabled"`
	SearchLimit     int    `json:"search_limit"`
}