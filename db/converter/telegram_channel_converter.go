package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// TelegramChannelToResponse 将TelegramChannel实体转换为响应DTO
func TelegramChannelToResponse(channel entity.TelegramChannel) dto.TelegramChannelResponse {
	return dto.TelegramChannelResponse{
		ID:                channel.ID,
		ChatID:            channel.ChatID,
		ChatName:          channel.ChatName,
		ChatType:          channel.ChatType,
		PushEnabled:       channel.PushEnabled,
		PushFrequency:     channel.PushFrequency,
		PushStartTime:     channel.PushStartTime,
		PushEndTime:       channel.PushEndTime,
		ContentCategories: channel.ContentCategories,
		ContentTags:       channel.ContentTags,
		IsActive:          channel.IsActive,
		LastPushAt:        channel.LastPushAt,
		RegisteredBy:      channel.RegisteredBy,
		RegisteredAt:      channel.RegisteredAt,
	}
}

// TelegramChannelsToResponse 将TelegramChannel实体列表转换为响应DTO列表
func TelegramChannelsToResponse(channels []entity.TelegramChannel) []dto.TelegramChannelResponse {
	var responses []dto.TelegramChannelResponse
	for _, channel := range channels {
		responses = append(responses, TelegramChannelToResponse(channel))
	}
	return responses
}

// RequestToTelegramChannel 将请求DTO转换为TelegramChannel实体
func RequestToTelegramChannel(req dto.TelegramChannelRequest, registeredBy string) entity.TelegramChannel {
	return entity.TelegramChannel{
		ChatID:            req.ChatID,
		ChatName:          req.ChatName,
		ChatType:          req.ChatType,
		PushEnabled:       req.PushEnabled,
		PushFrequency:     req.PushFrequency,
		PushStartTime:     req.PushStartTime,
		PushEndTime:       req.PushEndTime,
		ContentCategories: req.ContentCategories,
		ContentTags:       req.ContentTags,
		IsActive:          req.IsActive,
		RegisteredBy:      registeredBy,
		RegisteredAt:      time.Now(),
	}
}

// TelegramBotConfigToResponse 将Telegram bot配置转换为响应DTO
func TelegramBotConfigToResponse(
	botEnabled bool,
	botApiKey string,
	autoReplyEnabled bool,
	autoReplyTemplate string,
	autoDeleteEnabled bool,
	autoDeleteInterval int,
	proxyEnabled bool,
	proxyType string,
	proxyHost string,
	proxyPort int,
	proxyUsername string,
	proxyPassword string,
) dto.TelegramBotConfigResponse {
	return dto.TelegramBotConfigResponse{
		BotEnabled:         botEnabled,
		BotApiKey:          botApiKey,
		AutoReplyEnabled:   autoReplyEnabled,
		AutoReplyTemplate:  autoReplyTemplate,
		AutoDeleteEnabled:  autoDeleteEnabled,
		AutoDeleteInterval: autoDeleteInterval,
		ProxyEnabled:       proxyEnabled,
		ProxyType:          proxyType,
		ProxyHost:          proxyHost,
		ProxyPort:          proxyPort,
		ProxyUsername:      proxyUsername,
		ProxyPassword:      proxyPassword,
	}
}

// SystemConfigToTelegramBotConfig 将系统配置转换为Telegram bot配置响应
func SystemConfigToTelegramBotConfig(configs []entity.SystemConfig) dto.TelegramBotConfigResponse {
	botEnabled := false
	botApiKey := ""
	autoReplyEnabled := true
	autoReplyTemplate := "您好！我可以帮您搜索网盘资源，请输入您要搜索的内容。"
	autoDeleteEnabled := false
	autoDeleteInterval := 60
	proxyEnabled := false
	proxyType := "http"
	proxyHost := ""
	proxyPort := 8080
	proxyUsername := ""
	proxyPassword := ""

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyTelegramBotEnabled:
			botEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramBotApiKey:
			botApiKey = config.Value
		case entity.ConfigKeyTelegramAutoReplyEnabled:
			autoReplyEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramAutoReplyTemplate:
			autoReplyTemplate = config.Value
		case entity.ConfigKeyTelegramAutoDeleteEnabled:
			autoDeleteEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramAutoDeleteInterval:
			if config.Value != "" {
				// 简单解析整数，这里可以改进错误处理
				var val int
				if _, err := fmt.Sscanf(config.Value, "%d", &val); err == nil {
					autoDeleteInterval = val
				}
			}
		case entity.ConfigKeyTelegramProxyEnabled:
			proxyEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramProxyType:
			proxyType = config.Value
		case entity.ConfigKeyTelegramProxyHost:
			proxyHost = config.Value
		case entity.ConfigKeyTelegramProxyPort:
			if config.Value != "" {
				var val int
				if _, err := fmt.Sscanf(config.Value, "%d", &val); err == nil {
					proxyPort = val
				}
			}
		case entity.ConfigKeyTelegramProxyUsername:
			proxyUsername = config.Value
		case entity.ConfigKeyTelegramProxyPassword:
			proxyPassword = config.Value
		}
	}

	return TelegramBotConfigToResponse(
		botEnabled,
		botApiKey,
		autoReplyEnabled,
		autoReplyTemplate,
		autoDeleteEnabled,
		autoDeleteInterval,
		proxyEnabled,
		proxyType,
		proxyHost,
		proxyPort,
		proxyUsername,
		proxyPassword,
	)
}

// TelegramBotConfigRequestToSystemConfigs 将Telegram bot配置请求转换为系统配置实体列表
func TelegramBotConfigRequestToSystemConfigs(req dto.TelegramBotConfigRequest) []entity.SystemConfig {
	configs := []entity.SystemConfig{}

	// 添加调试日志
	utils.Debug("[TELEGRAM:CONVERTER] 转换请求: %+v", req)

	if req.BotEnabled != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramBotEnabled,
			Value: boolToString(*req.BotEnabled),
			Type:  entity.ConfigTypeBool,
		})
	}

	if req.BotApiKey != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramBotApiKey,
			Value: *req.BotApiKey,
			Type:  entity.ConfigTypeString,
		})
	}

	if req.AutoReplyEnabled != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramAutoReplyEnabled,
			Value: boolToString(*req.AutoReplyEnabled),
			Type:  entity.ConfigTypeBool,
		})
	}

	if req.AutoReplyTemplate != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramAutoReplyTemplate,
			Value: *req.AutoReplyTemplate,
			Type:  entity.ConfigTypeString,
		})
	}

	if req.AutoDeleteEnabled != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramAutoDeleteEnabled,
			Value: boolToString(*req.AutoDeleteEnabled),
			Type:  entity.ConfigTypeBool,
		})
	}

	if req.AutoDeleteInterval != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramAutoDeleteInterval,
			Value: intToString(*req.AutoDeleteInterval),
			Type:  entity.ConfigTypeInt,
		})
	}

	if req.ProxyEnabled != nil {
		utils.Debug("[TELEGRAM:CONVERTER] 添加代理启用配置: %v", *req.ProxyEnabled)
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramProxyEnabled,
			Value: boolToString(*req.ProxyEnabled),
			Type:  entity.ConfigTypeBool,
		})
	}

	if req.ProxyType != nil {
		utils.Debug("[TELEGRAM:CONVERTER] 添加代理类型配置: %s", *req.ProxyType)
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramProxyType,
			Value: *req.ProxyType,
			Type:  entity.ConfigTypeString,
		})
	}

	if req.ProxyHost != nil {
		utils.Debug("[TELEGRAM:CONVERTER] 添加代理主机配置: %s", *req.ProxyHost)
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramProxyHost,
			Value: *req.ProxyHost,
			Type:  entity.ConfigTypeString,
		})
	}

	if req.ProxyPort != nil {
		utils.Debug("[TELEGRAM:CONVERTER] 添加代理端口配置: %d", *req.ProxyPort)
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramProxyPort,
			Value: intToString(*req.ProxyPort),
			Type:  entity.ConfigTypeInt,
		})
	}

	if req.ProxyUsername != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramProxyUsername,
			Value: *req.ProxyUsername,
			Type:  entity.ConfigTypeString,
		})
	}

	if req.ProxyPassword != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyTelegramProxyPassword,
			Value: *req.ProxyPassword,
			Type:  entity.ConfigTypeString,
		})
	}

	utils.Debug("[TELEGRAM:CONVERTER] 转换完成，共生成 %d 个配置项", len(configs))
	for i, config := range configs {
		if strings.Contains(config.Key, "proxy") {
			utils.Debug("[TELEGRAM:CONVERTER] 配置项 %d: %s = %s", i+1, config.Key, config.Value)
		}
	}

	return configs
}

// 辅助函数：布尔转换为字符串
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// 辅助函数：整数转换为字符串
func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}
