package converter

import (
	"time"

	"res_db/db/dto"
	"res_db/db/entity"
)

// SystemConfigToResponse 将系统配置实体转换为响应DTO
func SystemConfigToResponse(config *entity.SystemConfig) *dto.SystemConfigResponse {
	if config == nil {
		return nil
	}

	return &dto.SystemConfigResponse{
		ID:        config.ID,
		CreatedAt: config.CreatedAt.Format(time.RFC3339),
		UpdatedAt: config.UpdatedAt.Format(time.RFC3339),

		// SEO 配置
		SiteTitle:       config.SiteTitle,
		SiteDescription: config.SiteDescription,
		Keywords:        config.Keywords,
		Author:          config.Author,
		Copyright:       config.Copyright,

		// 自动处理配置组
		AutoProcessReadyResources: config.AutoProcessReadyResources,
		AutoProcessInterval:       config.AutoProcessInterval,
		AutoTransferEnabled:       config.AutoTransferEnabled,
		AutoFetchHotDramaEnabled:  config.AutoFetchHotDramaEnabled,

		// API配置
		ApiToken: config.ApiToken,

		// 其他配置
		PageSize:        config.PageSize,
		MaintenanceMode: config.MaintenanceMode,
	}
}

// RequestToSystemConfig 将请求DTO转换为系统配置实体
func RequestToSystemConfig(req *dto.SystemConfigRequest) *entity.SystemConfig {
	if req == nil {
		return nil
	}

	return &entity.SystemConfig{
		// SEO 配置
		SiteTitle:       req.SiteTitle,
		SiteDescription: req.SiteDescription,
		Keywords:        req.Keywords,
		Author:          req.Author,
		Copyright:       req.Copyright,

		// 自动处理配置组
		AutoProcessReadyResources: req.AutoProcessReadyResources,
		AutoProcessInterval:       req.AutoProcessInterval,
		AutoTransferEnabled:       req.AutoTransferEnabled,
		AutoFetchHotDramaEnabled:  req.AutoFetchHotDramaEnabled,

		// API配置
		ApiToken: req.ApiToken,

		// 其他配置
		PageSize:        req.PageSize,
		MaintenanceMode: req.MaintenanceMode,
	}
}
