package converter

import (
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

	// 只添加有值的字段
	if req.SiteTitle != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeySiteTitle, Value: req.SiteTitle, Type: entity.ConfigTypeString})
	}
	if req.SiteDescription != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeySiteDescription, Value: req.SiteDescription, Type: entity.ConfigTypeString})
	}
	if req.Keywords != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyKeywords, Value: req.Keywords, Type: entity.ConfigTypeString})
	}
	if req.Author != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAuthor, Value: req.Author, Type: entity.ConfigTypeString})
	}
	if req.Copyright != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyCopyright, Value: req.Copyright, Type: entity.ConfigTypeString})
	}
	if req.ApiToken != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyApiToken, Value: req.ApiToken, Type: entity.ConfigTypeString})
	}
	if req.ForbiddenWords != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyForbiddenWords, Value: req.ForbiddenWords, Type: entity.ConfigTypeString})
	}

	// 布尔值字段 - 只处理实际提交的字段
	// 注意：由于 Go 的零值机制，我们需要通过其他方式判断字段是否被提交
	// 这里暂时保持原样，但建议前端只提交有变化的字段
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoProcessReadyResources, Value: strconv.FormatBool(req.AutoProcessReadyResources), Type: entity.ConfigTypeBool})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoTransferEnabled, Value: strconv.FormatBool(req.AutoTransferEnabled), Type: entity.ConfigTypeBool})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoFetchHotDramaEnabled, Value: strconv.FormatBool(req.AutoFetchHotDramaEnabled), Type: entity.ConfigTypeBool})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyMaintenanceMode, Value: strconv.FormatBool(req.MaintenanceMode), Type: entity.ConfigTypeBool})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyEnableRegister, Value: strconv.FormatBool(req.EnableRegister), Type: entity.ConfigTypeBool})

	// 整数字段 - 添加所有提交的字段，包括0值
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoProcessInterval, Value: strconv.Itoa(req.AutoProcessInterval), Type: entity.ConfigTypeInt})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoTransferLimitDays, Value: strconv.Itoa(req.AutoTransferLimitDays), Type: entity.ConfigTypeInt})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyAutoTransferMinSpace, Value: strconv.Itoa(req.AutoTransferMinSpace), Type: entity.ConfigTypeInt})
	configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyPageSize, Value: strconv.Itoa(req.PageSize), Type: entity.ConfigTypeInt})

	// 三方统计配置
	if req.ThirdPartyStatsCode != "" {
		configs = append(configs, entity.SystemConfig{Key: entity.ConfigKeyThirdPartyStatsCode, Value: req.ThirdPartyStatsCode, Type: entity.ConfigTypeString})
	}

	return configs
}

// SystemConfigToPublicResponse 返回不含 api_token 的系统配置响应
func SystemConfigToPublicResponse(configs []entity.SystemConfig) map[string]interface{} {
	response := map[string]interface{}{
		entity.ConfigResponseFieldID:                        0,
		entity.ConfigResponseFieldCreatedAt:                 utils.GetCurrentTimeString(),
		entity.ConfigResponseFieldUpdatedAt:                 utils.GetCurrentTimeString(),
		entity.ConfigResponseFieldSiteTitle:                 entity.ConfigDefaultSiteTitle,
		entity.ConfigResponseFieldSiteDescription:           entity.ConfigDefaultSiteDescription,
		entity.ConfigResponseFieldKeywords:                  entity.ConfigDefaultKeywords,
		entity.ConfigResponseFieldAuthor:                    entity.ConfigDefaultAuthor,
		entity.ConfigResponseFieldCopyright:                 entity.ConfigDefaultCopyright,
		entity.ConfigResponseFieldAutoProcessReadyResources: false,
		entity.ConfigResponseFieldAutoProcessInterval:       30,
		entity.ConfigResponseFieldAutoTransferEnabled:       false,
		entity.ConfigResponseFieldAutoTransferLimitDays:     0,
		entity.ConfigResponseFieldAutoTransferMinSpace:      100,
		entity.ConfigResponseFieldAutoFetchHotDramaEnabled:  false,
		entity.ConfigResponseFieldForbiddenWords:            "",
		entity.ConfigResponseFieldPageSize:                  100,
		entity.ConfigResponseFieldMaintenanceMode:           false,
		entity.ConfigResponseFieldEnableRegister:            true, // 默认开启注册功能
	}

	// 将键值对转换为map
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
		case entity.ConfigKeyAutoProcessReadyResources:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response[entity.ConfigResponseFieldAutoProcessReadyResources] = val
			}
		case entity.ConfigKeyAutoProcessInterval:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response[entity.ConfigResponseFieldAutoProcessInterval] = val
			}
		case entity.ConfigKeyAutoTransferEnabled:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response[entity.ConfigResponseFieldAutoTransferEnabled] = val
			}
		case entity.ConfigKeyAutoTransferLimitDays:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response[entity.ConfigResponseFieldAutoTransferLimitDays] = val
			}
		case entity.ConfigKeyAutoTransferMinSpace:
			if val, err := strconv.Atoi(config.Value); err == nil {
				response[entity.ConfigResponseFieldAutoTransferMinSpace] = val
			}
		case entity.ConfigKeyAutoFetchHotDramaEnabled:
			if val, err := strconv.ParseBool(config.Value); err == nil {
				response[entity.ConfigResponseFieldAutoFetchHotDramaEnabled] = val
			}
		case entity.ConfigKeyForbiddenWords:
			response[entity.ConfigResponseFieldForbiddenWords] = config.Value
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
		AutoProcessReadyResources: false,
		AutoProcessInterval:       30,
		AutoTransferEnabled:       false,
		AutoTransferLimitDays:     0,
		AutoTransferMinSpace:      100,
		AutoFetchHotDramaEnabled:  false,
		ApiToken:                  entity.ConfigDefaultApiToken,
		ForbiddenWords:            entity.ConfigDefaultForbiddenWords,
		PageSize:                  100,
		MaintenanceMode:           false,
		EnableRegister:            true, // 默认开启注册功能
		ThirdPartyStatsCode:       entity.ConfigDefaultThirdPartyStatsCode,
	}
}
