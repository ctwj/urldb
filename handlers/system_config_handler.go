package handlers

import (
	"net/http"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// SystemConfigHandler 系统配置处理器
type SystemConfigHandler struct {
	systemConfigRepo repo.SystemConfigRepository
}

// NewSystemConfigHandler 创建系统配置处理器
func NewSystemConfigHandler(systemConfigRepo repo.SystemConfigRepository) *SystemConfigHandler {
	return &SystemConfigHandler{
		systemConfigRepo: systemConfigRepo,
	}
}

// GetConfig 获取系统配置
func (h *SystemConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(config)
	SuccessResponse(c, configResponse)
}

// UpdateConfig 更新系统配置
func (h *SystemConfigHandler) UpdateConfig(c *gin.Context) {
	var req dto.SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.SiteTitle == "" {
		ErrorResponse(c, "网站标题不能为空", http.StatusBadRequest)
		return
	}

	if req.AutoProcessInterval < 1 || req.AutoProcessInterval > 1440 {
		ErrorResponse(c, "自动处理间隔必须在1-1440分钟之间", http.StatusBadRequest)
		return
	}

	if req.PageSize < 10 || req.PageSize > 500 {
		ErrorResponse(c, "每页显示数量必须在10-500之间", http.StatusBadRequest)
		return
	}

	// 验证自动转存配置
	if req.AutoTransferLimitDays < 0 || req.AutoTransferLimitDays > 365 {
		ErrorResponse(c, "自动转存限制天数必须在0-365之间", http.StatusBadRequest)
		return
	}

	if req.AutoTransferMinSpace < 100 || req.AutoTransferMinSpace > 1024 {
		ErrorResponse(c, "最小存储空间必须在100-1024GB之间", http.StatusBadRequest)
		return
	}

	// 转换为实体
	config := converter.RequestToSystemConfig(&req)
	if config == nil {
		ErrorResponse(c, "数据转换失败", http.StatusInternalServerError)
		return
	}

	// 保存配置
	err := h.systemConfigRepo.Upsert(config)
	if err != nil {
		ErrorResponse(c, "保存系统配置失败", http.StatusInternalServerError)
		return
	}

	// 返回更新后的配置
	updatedConfig, err := h.systemConfigRepo.FindFirst()
	if err != nil {
		ErrorResponse(c, "获取更新后的配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(updatedConfig)
	SuccessResponse(c, configResponse)
}

// GetSystemConfig 获取系统配置（使用全局repoManager）
func GetSystemConfig(c *gin.Context) {
	config, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(config)
	SuccessResponse(c, configResponse)
}

// UpdateSystemConfig 更新系统配置（使用全局repoManager）
func UpdateSystemConfig(c *gin.Context) {
	var req dto.SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.SiteTitle == "" {
		ErrorResponse(c, "网站标题不能为空", http.StatusBadRequest)
		return
	}

	if req.AutoProcessInterval < 1 || req.AutoProcessInterval > 1440 {
		ErrorResponse(c, "自动处理间隔必须在1-1440分钟之间", http.StatusBadRequest)
		return
	}

	if req.PageSize < 10 || req.PageSize > 500 {
		ErrorResponse(c, "每页显示数量必须在10-500之间", http.StatusBadRequest)
		return
	}

	// 验证自动转存配置
	if req.AutoTransferLimitDays < 0 || req.AutoTransferLimitDays > 365 {
		ErrorResponse(c, "自动转存限制天数必须在0-365之间", http.StatusBadRequest)
		return
	}

	if req.AutoTransferMinSpace < 100 || req.AutoTransferMinSpace > 1024 {
		ErrorResponse(c, "最小存储空间必须在100-1024GB之间", http.StatusBadRequest)
		return
	}

	// 转换为实体
	config := converter.RequestToSystemConfig(&req)
	if config == nil {
		ErrorResponse(c, "数据转换失败", http.StatusInternalServerError)
		return
	}

	// 保存配置
	err := repoManager.SystemConfigRepository.Upsert(config)
	if err != nil {
		ErrorResponse(c, "保存系统配置失败", http.StatusInternalServerError)
		return
	}

	// 根据配置更新定时任务状态（错误不影响配置保存）
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
		repoManager.TagRepository,
		repoManager.CategoryRepository,
	)
	if scheduler != nil {
		scheduler.UpdateSchedulerStatusWithAutoTransfer(req.AutoFetchHotDramaEnabled, req.AutoProcessReadyResources, req.AutoTransferEnabled)
	}

	// 返回更新后的配置
	updatedConfig, err := repoManager.SystemConfigRepository.FindFirst()
	if err != nil {
		ErrorResponse(c, "获取更新后的配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(updatedConfig)
	SuccessResponse(c, configResponse)
}

// 新增：公开获取系统配置（不含api_token）
func GetPublicSystemConfig(c *gin.Context) {
	config, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}
	configResponse := converter.SystemConfigToPublicResponse(config)
	SuccessResponse(c, configResponse)
}

// 新增：切换自动处理配置
func ToggleAutoProcess(c *gin.Context) {
	var req struct {
		AutoProcessReadyResources bool `json:"auto_process_ready_resources"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 获取当前配置
	config, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}

	// 只更新自动处理配置
	config.AutoProcessReadyResources = req.AutoProcessReadyResources

	// 保存配置
	err = repoManager.SystemConfigRepository.Upsert(config)
	if err != nil {
		ErrorResponse(c, "保存系统配置失败", http.StatusInternalServerError)
		return
	}

	// 更新定时任务状态
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
		repoManager.TagRepository,
		repoManager.CategoryRepository,
	)
	if scheduler != nil {
		scheduler.UpdateSchedulerStatusWithAutoTransfer(
			config.AutoFetchHotDramaEnabled,
			config.AutoProcessReadyResources,
			config.AutoTransferEnabled,
		)
	}

	// 返回更新后的配置
	configResponse := converter.SystemConfigToResponse(config)
	SuccessResponse(c, configResponse)
}
