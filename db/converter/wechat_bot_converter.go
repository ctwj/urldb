package converter

import (
	"strconv"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// WechatBotConfigRequestToSystemConfigs 将微信机器人配置请求转换为系统配置实体
func WechatBotConfigRequestToSystemConfigs(req dto.WechatBotConfigRequest) []entity.SystemConfig {
	configs := []entity.SystemConfig{
		{Key: entity.ConfigKeyWechatBotEnabled, Value: wechatBoolToString(req.Enabled)},
		{Key: entity.ConfigKeyWechatAppId, Value: req.AppID},
		{Key: entity.ConfigKeyWechatAppSecret, Value: req.AppSecret},
		{Key: entity.ConfigKeyWechatToken, Value: req.Token},
		{Key: entity.ConfigKeyWechatEncodingAesKey, Value: req.EncodingAesKey},
		{Key: entity.ConfigKeyWechatWelcomeMessage, Value: req.WelcomeMessage},
		{Key: entity.ConfigKeyWechatAutoReplyEnabled, Value: wechatBoolToString(req.AutoReplyEnabled)},
		{Key: entity.ConfigKeyWechatSearchLimit, Value: wechatIntToString(req.SearchLimit)},
	}
	return configs
}

// SystemConfigToWechatBotConfig 将系统配置转换为微信机器人配置响应
func SystemConfigToWechatBotConfig(configs []entity.SystemConfig) dto.WechatBotConfigResponse {
	resp := dto.WechatBotConfigResponse{
		Enabled:          false,
		AppID:            "",
		AppSecret:        "",
		Token:            "",
		EncodingAesKey:   "",
		WelcomeMessage:   "欢迎关注老九网盘资源库！发送关键词即可搜索资源。",
		AutoReplyEnabled: true,
		SearchLimit:      5,
	}

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyWechatBotEnabled:
			resp.Enabled = config.Value == "true"
		case entity.ConfigKeyWechatAppId:
			resp.AppID = config.Value
		case entity.ConfigKeyWechatAppSecret:
			resp.AppSecret = config.Value
		case entity.ConfigKeyWechatToken:
			resp.Token = config.Value
		case entity.ConfigKeyWechatEncodingAesKey:
			resp.EncodingAesKey = config.Value
		case entity.ConfigKeyWechatWelcomeMessage:
			if config.Value != "" {
				resp.WelcomeMessage = config.Value
			}
		case entity.ConfigKeyWechatAutoReplyEnabled:
			resp.AutoReplyEnabled = config.Value == "true"
		case entity.ConfigKeyWechatSearchLimit:
			if config.Value != "" {
				resp.SearchLimit = wechatStringToInt(config.Value)
			}
		}
	}

	return resp
}

// 辅助函数 - 使用大写名称避免与其他文件中的函数冲突
func wechatBoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func wechatIntToString(i int) string {
	return strconv.Itoa(i)
}

func wechatStringToInt(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
