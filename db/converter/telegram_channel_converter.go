package converter

import (
	"fmt"
	"time"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
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
) dto.TelegramBotConfigResponse {
	return dto.TelegramBotConfigResponse{
		BotEnabled:         botEnabled,
		BotApiKey:          botApiKey,
		AutoReplyEnabled:   autoReplyEnabled,
		AutoReplyTemplate:  autoReplyTemplate,
		AutoDeleteEnabled:  autoDeleteEnabled,
		AutoDeleteInterval: autoDeleteInterval,
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
		}
	}

	return TelegramBotConfigToResponse(
		botEnabled,
		botApiKey,
		autoReplyEnabled,
		autoReplyTemplate,
		autoDeleteEnabled,
		autoDeleteInterval,
	)
}

// TelegramBotConfigRequestToSystemConfigs 将Telegram bot配置请求转换为系统配置实体列表
func TelegramBotConfigRequestToSystemConfigs(req dto.TelegramBotConfigRequest) []entity.SystemConfig {
	configs := []entity.SystemConfig{}

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
