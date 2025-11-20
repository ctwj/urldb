package converter

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// SystemConfigToResponse 将系统配置实体列表转换为响应DTO
func SystemConfigToResponse(configs []entity.SystemConfig) *dto.SystemConfigResponse {
	if len(configs) == 0 {
		return getDefaultConfigResponse()
	}

	response := getDefaultConfigResponse()

	// 将键值对转换为结构体
	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeySiteTitle:
			response.SiteTitle = config.Value
		case entity.ConfigKeySiteDescription:
			response.SiteDescription = config.Value
		case entity.ConfigKeyKeywords:
			response.Keywords = config.Value
		case entity.ConfigKeyAuthor:
			response.Author = config.Value
		case entity.ConfigKeyCopyright:
			response.Copyright = config.Value
		case entity.ConfigKeySiteLogo:
			response.SiteLogo = config.Value
		case entity.ConfigKeyAutoProcessReadyResources:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.AutoProcessReadyResources = val
			}
		case entity.ConfigKeyAutoProcessInterval:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response.AutoProcessInterval = val
			}
		case entity.ConfigKeyAutoTransferEnabled:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.AutoTransferEnabled = val
			}
		case entity.ConfigKeyAutoTransferLimitDays:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response.AutoTransferLimitDays = val
			}
		case entity.ConfigKeyAutoTransferMinSpace:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response.AutoTransferMinSpace = val
			}
		case entity.ConfigKeyAutoFetchHotDramaEnabled:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.AutoFetchHotDramaEnabled = val
			}
		case entity.ConfigKeyApiToken:
			response.ApiToken = config.Value
		case entity.ConfigKeyForbiddenWords:
			response.ForbiddenWords = config.Value
		case entity.ConfigKeyAdKeywords:
			response.AdKeywords = config.Value
		case entity.ConfigKeyAutoInsertAd:
			response.AutoInsertAd = config.Value
		case entity.ConfigKeyPageSize:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response.PageSize = val
			}
		case entity.ConfigKeyMaintenanceMode:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.MaintenanceMode = val
			}
		case entity.ConfigKeyEnableRegister:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.EnableRegister = val
			}
		case entity.ConfigKeyThirdPartyStatsCode:
			response.ThirdPartyStatsCode = config.Value
		case entity.ConfigKeyMeilisearchEnabled:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.MeilisearchEnabled = val
			}
		case entity.ConfigKeyMeilisearchHost:
			response.MeilisearchHost = config.Value
		case entity.ConfigKeyMeilisearchPort:
			response.MeilisearchPort = config.Value
		case entity.ConfigKeyMeilisearchMasterKey:
			response.MeilisearchMasterKey = config.Value
		case entity.ConfigKeyMeilisearchIndexName:
			response.MeilisearchIndexName = config.Value
		case entity.ConfigKeyEnableAnnouncements:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.EnableAnnouncements = val
			}
		case entity.ConfigKeyAnnouncements:
			if config.Value == "" || config.Value == "[]" {
				response.Announcements = ""
			} else {
				// 在响应时保持为字符串，后续由前端处理
				response.Announcements = config.Value
			}
		case entity.ConfigKeyEnableFloatButtons:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response.EnableFloatButtons = val
			}
		case entity.ConfigKeyWechatSearchImage:
			response.WechatSearchImage = config.Value
		case entity.ConfigKeyTelegramQrImage:
			response.TelegramQrImage = config.Value
		case entity.ConfigKeyQrCodeStyle:
			response.QrCodeStyle = config.Value
		}
	}

	// 设置时间戳（使用第一个配置的时间）
	if len(configs) > 0 {
		response.CreatedAt = configs[0].CreatedAt.Format(time.RFC3339)
		response.UpdatedAt = configs[0].UpdatedAt.Format(time.RFC3339)
	}

	return response
}

// RequestToSystemConfig 将请求DTO转换为系统配置实体列表
func RequestToSystemConfig(req *dto.SystemConfigRequest) []entity.SystemConfig {
	if req == nil {
		return nil
	}

	var configs []entity.SystemConfig
	var updatedKeys []string

	// 字符串字段 - 只处理被设置的字段
	if req.SiteTitle != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeySiteTitle, Value: *req.SiteTitle, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeySiteTitle)
	}
	if req.SiteDescription != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeySiteDescription, Value: *req.SiteDescription, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeySiteDescription)
	}
	if req.Keywords != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyKeywords, Value: *req.Keywords, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyKeywords)
	}
	if req.Author != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAuthor, Value: *req.Author, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAuthor)
	}
	if req.Copyright != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyCopyright, Value: *req.Copyright, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyCopyright)
	}
	if req.SiteLogo != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeySiteLogo, Value: *req.SiteLogo, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeySiteLogo)
	}
	if req.ApiToken != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyApiToken, Value: *req.ApiToken, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyApiToken)
	}
	if req.ForbiddenWords != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyForbiddenWords, Value: *req.ForbiddenWords, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyForbiddenWords)
	}
	if req.AdKeywords != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAdKeywords, Value: *req.AdKeywords, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAdKeywords)
	}
	if req.AutoInsertAd != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoInsertAd, Value: *req.AutoInsertAd, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoInsertAd)
	}

	// 布尔值字段 - 只处理被设置的字段
	if req.AutoProcessReadyResources != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoProcessReadyResources, Value: strconv.FormatBool(*req.AutoProcessReadyResources), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoProcessReadyResources)
	}
	if req.AutoTransferEnabled != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoTransferEnabled, Value: strconv.FormatBool(*req.AutoTransferEnabled), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoTransferEnabled)
	}
	if req.AutoFetchHotDramaEnabled != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoFetchHotDramaEnabled, Value: strconv.FormatBool(*req.AutoFetchHotDramaEnabled), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoFetchHotDramaEnabled)
	}
	if req.MaintenanceMode != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMaintenanceMode, Value: strconv.FormatBool(*req.MaintenanceMode), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyMaintenanceMode)
	}
	if req.EnableRegister != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyEnableRegister, Value: strconv.FormatBool(*req.EnableRegister), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyEnableRegister)
	}

	// 整数字段 - 只处理被设置的字段
	if req.AutoProcessInterval != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoProcessInterval, Value: strconv.Itoa(*req.AutoProcessInterval), Type: entity.ConfigTypeInt})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoProcessInterval)
	}
	if req.AutoTransferLimitDays != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoTransferLimitDays, Value: strconv.Itoa(*req.AutoTransferLimitDays), Type: entity.ConfigTypeInt})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoTransferLimitDays)
	}
	if req.AutoTransferMinSpace != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoTransferMinSpace, Value: strconv.Itoa(*req.AutoTransferMinSpace), Type: entity.ConfigTypeInt})
		updatedKeys = append(updatedKeys, entity.ConfigKeyAutoTransferMinSpace)
	}
	if req.PageSize != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyPageSize, Value: strconv.Itoa(*req.PageSize), Type: entity.ConfigTypeInt})
		updatedKeys = append(updatedKeys, entity.ConfigKeyPageSize)
	}

	// 三方统计配置 - 只处理被设置的字段
	if req.ThirdPartyStatsCode != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyThirdPartyStatsCode, Value: *req.ThirdPartyStatsCode, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyThirdPartyStatsCode)
	}

	// Meilisearch配置 - 只处理被设置的字段
	if req.MeilisearchEnabled != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMeilisearchEnabled, Value: strconv.FormatBool(*req.MeilisearchEnabled), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyMeilisearchEnabled)
	}
	if req.MeilisearchHost != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMeilisearchHost, Value: *req.MeilisearchHost, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyMeilisearchHost)
	}
	if req.MeilisearchPort != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMeilisearchPort, Value: *req.MeilisearchPort, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyMeilisearchPort)
	}
	if req.MeilisearchMasterKey != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMeilisearchMasterKey, Value: *req.MeilisearchMasterKey, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyMeilisearchMasterKey)
	}
	if req.MeilisearchIndexName != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMeilisearchIndexName, Value: *req.MeilisearchIndexName, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyMeilisearchIndexName)
	}

	// 界面配置处理
	if req.EnableAnnouncements != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyEnableAnnouncements, Value: strconv.FormatBool(*req.EnableAnnouncements), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyEnableAnnouncements)
	}
	if req.Announcements != nil {
		// 将数组转换为JSON字符串
		if jsonBytes, err := json.Marshal(*req.Announcements); err == nil {
			configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAnnouncements, Value: string(jsonBytes), Type: entity.ConfigTypeJSON})
			updatedKeys = append(updatedKeys, entity.ConfigKeyAnnouncements)
		}
	}
	if req.EnableFloatButtons != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyEnableFloatButtons, Value: strconv.FormatBool(*req.EnableFloatButtons), Type: entity.ConfigTypeBool})
		updatedKeys = append(updatedKeys, entity.ConfigKeyEnableFloatButtons)
	}
	if req.WechatSearchImage != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyWechatSearchImage, Value: *req.WechatSearchImage, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyWechatSearchImage)
	}
	if req.TelegramQrImage != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyTelegramQrImage, Value: *req.TelegramQrImage, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyTelegramQrImage)
	}
	if req.QrCodeStyle != nil {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyQrCodeStyle, Value: *req.QrCodeStyle, Type: entity.ConfigTypeString})
		updatedKeys = append(updatedKeys, entity.ConfigKeyQrCodeStyle)
	}

	// 记录更新的配置项
	if len(updatedKeys) > 0 {
		utils.Info("配置更新 - 被修改的配置项: %v", updatedKeys)
	}

	return configs
}

// SystemConfigToPublicResponse 返回不含敏感配置的系统配置响应
func SystemConfigToPublicResponse(configs []entity.SystemConfig) map[string]interface{} {
	response := map[string]interface{}{
		entity.ConfigResponseFieldID:                  0,
		entity.ConfigResponseFieldCreatedAt:           utils.GetCurrentTimeString(),
		entity.ConfigResponseFieldUpdatedAt:           utils.GetCurrentTimeString(),
		entity.ConfigResponseFieldSiteTitle:           entity.ConfigDefaultSiteTitle,
		entity.ConfigResponseFieldSiteDescription:     entity.ConfigDefaultSiteDescription,
		entity.ConfigResponseFieldKeywords:            entity.ConfigDefaultKeywords,
		entity.ConfigResponseFieldAuthor:              entity.ConfigDefaultAuthor,
		entity.ConfigResponseFieldCopyright:           entity.ConfigDefaultCopyright,
		"site_logo":                                   "",
		entity.ConfigResponseFieldAdKeywords:          "",
		entity.ConfigResponseFieldAutoInsertAd:        "",
		entity.ConfigResponseFieldPageSize:            100,
		entity.ConfigResponseFieldMaintenanceMode:     false,
		entity.ConfigResponseFieldEnableRegister:      true, // 默认开启注册功能
		entity.ConfigResponseFieldThirdPartyStatsCode: "",
	}

	// 将键值对转换为map，过滤掉敏感配置
	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeySiteTitle:
			response[entity.ConfigResponseFieldSiteTitle] = config.Value
		case entity.ConfigKeySiteDescription:
			response[entity.ConfigResponseFieldSiteDescription] = config.Value
		case entity.ConfigKeyKeywords:
			response[entity.ConfigResponseFieldKeywords] = config.Value
		case entity.ConfigKeyAuthor:
			response[entity.ConfigResponseFieldAuthor] = config.Value
		case entity.ConfigKeyCopyright:
			response[entity.ConfigResponseFieldCopyright] = config.Value
		case entity.ConfigKeySiteLogo:
			response["site_logo"] = config.Value
		case entity.ConfigKeyAdKeywords:
			response[entity.ConfigResponseFieldAdKeywords] = config.Value
		case entity.ConfigKeyAutoInsertAd:
			response[entity.ConfigResponseFieldAutoInsertAd] = config.Value
		case entity.ConfigKeyPageSize:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response[entity.ConfigResponseFieldPageSize] = val
			}
		case entity.ConfigKeyMaintenanceMode:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response[entity.ConfigResponseFieldMaintenanceMode] = val
			}
		case entity.ConfigKeyEnableRegister:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response[entity.ConfigResponseFieldEnableRegister] = val
			}
		case entity.ConfigKeyThirdPartyStatsCode:
			response[entity.ConfigResponseFieldThirdPartyStatsCode] = config.Value
		case entity.ConfigKeyEnableAnnouncements:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response["enable_announcements"] = val
			}
		case entity.ConfigKeyAnnouncements:
			if config.Value == "" || config.Value == "[]" {
				response["announcements"] = ""
			} else {
				response["announcements"] = config.Value
			}
		case entity.ConfigKeyEnableFloatButtons:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response["enable_float_buttons"] = val
			}
		case entity.ConfigKeyWechatSearchImage:
			response["wechat_search_image"] = config.Value
		case entity.ConfigKeyTelegramQrImage:
			response["telegram_qr_image"] = config.Value
		case entity.ConfigKeyQrCodeStyle:
			response["qr_code_style"] = config.Value
		case entity.ConfigKeyAutoProcessReadyResources:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response["auto_process_ready_resources"] = val
			}
		case entity.ConfigKeyAutoTransferEnabled:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response["auto_transfer_enabled"] = val
			}
		// 跳过不需要返回给公众的配置
		case entity.ConfigKeyAutoProcessInterval:
		case entity.ConfigKeyAutoTransferLimitDays:
		case entity.ConfigKeyAutoTransferMinSpace:
		case entity.ConfigKeyAutoFetchHotDramaEnabled:
		case entity.ConfigKeyMeilisearchEnabled:
		case entity.ConfigKeyMeilisearchHost:
		case entity.ConfigKeyMeilisearchPort:
		case entity.ConfigKeyMeilisearchMasterKey:
		case entity.ConfigKeyMeilisearchIndexName:
		case entity.ConfigKeyForbiddenWords:
			// 这些配置不返回给公众
			continue
		}
	}

	// 设置时间戳（使用第一个配置的时间）
	if len(configs) > 0 {
		response[entity.ConfigResponseFieldCreatedAt] = configs[0].CreatedAt.Format(utils.TimeFormatDateTime)
		response[entity.ConfigResponseFieldUpdatedAt] = configs[0].UpdatedAt.Format(utils.TimeFormatDateTime)
	}

	return response
}

// getDefaultConfigResponse 获取默认配置响应
func getDefaultConfigResponse() *dto.SystemConfigResponse {
	return &dto.SystemConfigResponse{
		SiteTitle:                 entity.ConfigDefaultSiteTitle,
		SiteDescription:           entity.ConfigDefaultSiteDescription,
		Keywords:                  entity.ConfigDefaultKeywords,
		Author:                    entity.ConfigDefaultAuthor,
		Copyright:                 entity.ConfigDefaultCopyright,
		SiteLogo:                  "",
		AutoProcessReadyResources: false,
		AutoProcessInterval:       30,
		AutoTransferEnabled:       false,
		AutoTransferLimitDays:     0,
		AutoTransferMinSpace:      100,
		AutoFetchHotDramaEnabled:  false,
		ApiToken:                  entity.ConfigDefaultApiToken,
		ForbiddenWords:            entity.ConfigDefaultForbiddenWords,
		AdKeywords:                entity.ConfigDefaultAdKeywords,
		AutoInsertAd:              entity.ConfigDefaultAutoInsertAd,
		PageSize:                  100,
		MaintenanceMode:           false,
		EnableRegister:            true, // 默认开启注册功能
		ThirdPartyStatsCode:       entity.ConfigDefaultThirdPartyStatsCode,
		MeilisearchEnabled:        false,
		MeilisearchHost:           entity.ConfigDefaultMeilisearchHost,
		MeilisearchPort:           entity.ConfigDefaultMeilisearchPort,
		MeilisearchMasterKey:      entity.ConfigDefaultMeilisearchMasterKey,
		MeilisearchIndexName:      entity.ConfigDefaultMeilisearchIndexName,
		EnableAnnouncements:       false,
		Announcements:             "",
		EnableFloatButtons:        false,
		WechatSearchImage:         entity.ConfigDefaultWechatSearchImage,
		TelegramQrImage:           entity.ConfigDefaultTelegramQrImage,
		QrCodeStyle:               entity.ConfigDefaultQrCodeStyle,
	}
}
