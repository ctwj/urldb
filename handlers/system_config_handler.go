package handlers

import (
	"net/http"

	pan "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/scheduler"
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
	// 先验证配置完整性
	if err := h.systemConfigRepo.ValidateConfigIntegrity(); err != nil {
		utils.Error("配置完整性检查失败: %v", err)
		// 如果配置不完整，尝试重新创建默认配置
		configs, err := h.systemConfigRepo.GetOrCreateDefault()
		if err != nil {
			ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
			return
		}
		configResponse := converter.SystemConfigToResponse(configs)
		SuccessResponse(c, configResponse)
		return
	}

	configs, err := h.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(configs)
	SuccessResponse(c, configResponse)
}

// UpdateConfig 更新系统配置
func (h *SystemConfigHandler) UpdateConfig(c *gin.Context) {
	var req dto.SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 验证参数 - 只验证提交的字段
	if req.SiteTitle != nil && (len(*req.SiteTitle) < 1 || len(*req.SiteTitle) > 100) {
		ErrorResponse(c, "网站标题长度必须在1-100字符之间", http.StatusBadRequest)
		return
	}

	if req.AutoProcessInterval != nil && (*req.AutoProcessInterval < 1 || *req.AutoProcessInterval > 1440) {
		ErrorResponse(c, "自动处理间隔必须在1-1440分钟之间", http.StatusBadRequest)
		return
	}

	if req.PageSize != nil && (*req.PageSize < 10 || *req.PageSize > 500) {
		ErrorResponse(c, "每页显示数量必须在10-500之间", http.StatusBadRequest)
		return
	}

	if req.AutoTransferMinSpace != nil && (*req.AutoTransferMinSpace < 100 || *req.AutoTransferMinSpace > 1024) {
		ErrorResponse(c, "最小存储空间必须在100-1024GB之间", http.StatusBadRequest)
		return
	}

	// 转换为实体
	configs := converter.RequestToSystemConfig(&req)
	if configs == nil {
		ErrorResponse(c, "数据转换失败", http.StatusInternalServerError)
		return
	}

	// 保存配置
	err := h.systemConfigRepo.UpsertConfigs(configs)
	if err != nil {
		ErrorResponse(c, "保存系统配置失败", http.StatusInternalServerError)
		return
	}

	// 刷新系统配置缓存
	pan.RefreshSystemConfigCache()

	// 返回更新后的配置
	updatedConfigs, err := h.systemConfigRepo.FindAll()
	if err != nil {
		ErrorResponse(c, "获取更新后的配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(updatedConfigs)
	SuccessResponse(c, configResponse)
}

// GetSystemConfig 获取系统配置（使用全局repoManager）
func GetSystemConfig(c *gin.Context) {
	configs, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}

	configResponse := converter.SystemConfigToResponse(configs)
	SuccessResponse(c, configResponse)
}

// UpdateSystemConfig 更新系统配置（使用全局repoManager）
func UpdateSystemConfig(c *gin.Context) {
	var req dto.SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 调试信息
	utils.Info("接收到的配置请求: %+v", req)

	// 获取当前配置作为备份
	currentConfigs, err := repoManager.SystemConfigRepository.FindAll()
	if err != nil {
		utils.Error("获取当前配置失败: %v", err)
	} else {
		utils.Info("当前配置数量: %d", len(currentConfigs))
	}

	// 验证参数 - 只验证提交的字段
	if req.SiteTitle != nil && (len(*req.SiteTitle) < 1 || len(*req.SiteTitle) > 100) {
		ErrorResponse(c, "网站标题长度必须在1-100字符之间", http.StatusBadRequest)
		return
	}

	if req.AutoProcessInterval != nil && (*req.AutoProcessInterval < 1 || *req.AutoProcessInterval > 1440) {
		ErrorResponse(c, "自动处理间隔必须在1-1440分钟之间", http.StatusBadRequest)
		return
	}

	if req.PageSize != nil && (*req.PageSize < 10 || *req.PageSize > 500) {
		ErrorResponse(c, "每页显示数量必须在10-500之间", http.StatusBadRequest)
		return
	}

	// 验证自动转存配置
	if req.AutoTransferLimitDays != nil && (*req.AutoTransferLimitDays < 0 || *req.AutoTransferLimitDays > 365) {
		ErrorResponse(c, "自动转存限制天数必须在0-365之间", http.StatusBadRequest)
		return
	}

	if req.AutoTransferMinSpace != nil && (*req.AutoTransferMinSpace < 100 || *req.AutoTransferMinSpace > 1024) {
		ErrorResponse(c, "最小存储空间必须在100-1024GB之间", http.StatusBadRequest)
		return
	}

	// 转换为实体
	configs := converter.RequestToSystemConfig(&req)
	if configs == nil {
		ErrorResponse(c, "数据转换失败", http.StatusInternalServerError)
		return
	}

	utils.Info("准备更新配置，配置项数量: %d", len(configs))

	// 保存配置
	err = repoManager.SystemConfigRepository.UpsertConfigs(configs)
	if err != nil {
		utils.Error("保存系统配置失败: %v", err)
		ErrorResponse(c, "保存系统配置失败", http.StatusInternalServerError)
		return
	}

	utils.Info("配置保存成功")

	// 安全刷新系统配置缓存
	if err := repoManager.SystemConfigRepository.SafeRefreshConfigCache(); err != nil {
		utils.Error("刷新配置缓存失败: %v", err)
		// 不返回错误，因为配置已经保存成功
	}

	// 刷新系统配置缓存
	pan.RefreshSystemConfigCache()

	// 重新加载Meilisearch配置（如果Meilisearch配置有变更）
	if req.MeilisearchEnabled != nil || req.MeilisearchHost != nil || req.MeilisearchPort != nil || req.MeilisearchMasterKey != nil || req.MeilisearchIndexName != nil {
		if meilisearchManager != nil {
			if err := meilisearchManager.ReloadConfig(); err != nil {
				utils.Error("重新加载Meilisearch配置失败: %v", err)
			} else {
				utils.Debug("Meilisearch配置重新加载成功")
			}
		}
	}

	// 根据配置更新定时任务状态（错误不影响配置保存）
	scheduler := scheduler.GetGlobalScheduler(
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
		// 只更新被设置的配置
		var autoFetchHotDrama, autoProcessReady, autoTransfer bool
		if req.AutoFetchHotDramaEnabled != nil {
			autoFetchHotDrama = *req.AutoFetchHotDramaEnabled
		}
		if req.AutoProcessReadyResources != nil {
			autoProcessReady = *req.AutoProcessReadyResources
		}
		if req.AutoTransferEnabled != nil {
			autoTransfer = *req.AutoTransferEnabled
		}
		scheduler.UpdateSchedulerStatusWithAutoTransfer(autoFetchHotDrama, autoProcessReady, autoTransfer)
	}

	// 返回更新后的配置
	updatedConfigs, err := repoManager.SystemConfigRepository.FindAll()
	if err != nil {
		utils.Error("获取更新后的配置失败: %v", err)
		ErrorResponse(c, "获取更新后的配置失败", http.StatusInternalServerError)
		return
	}

	utils.Info("配置更新完成，当前配置数量: %d", len(updatedConfigs))

	configResponse := converter.SystemConfigToResponse(updatedConfigs)
	SuccessResponse(c, configResponse)
}

// 新增：公开获取系统配置（不含api_token）
func GetPublicSystemConfig(c *gin.Context) {
	configs, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}
	configResponse := converter.SystemConfigToPublicResponse(configs)
	SuccessResponse(c, configResponse)
}

// 新增：配置监控端点
func GetConfigStatus(c *gin.Context) {
	// 获取配置统计信息
	configs, err := repoManager.SystemConfigRepository.FindAll()
	if err != nil {
		ErrorResponse(c, "获取配置状态失败", http.StatusInternalServerError)
		return
	}

	// 验证配置完整性
	integrityErr := repoManager.SystemConfigRepository.ValidateConfigIntegrity()

	// 获取缓存状态
	cachedConfigs := repoManager.SystemConfigRepository.GetCachedConfigs()

	status := map[string]interface{}{
		"total_configs":   len(configs),
		"cached_configs":  len(cachedConfigs),
		"integrity_check": integrityErr == nil,
		"integrity_error": "",
		"last_check_time": utils.GetCurrentTimeString(),
	}

	if integrityErr != nil {
		status["integrity_error"] = integrityErr.Error()
	}

	SuccessResponse(c, status)
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
	configs, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取系统配置失败", http.StatusInternalServerError)
		return
	}

	// 更新自动处理配置
	for i, config := range configs {
		if config.Key == entity.ConfigKeyAutoProcessReadyResources {
			configs[i].Value = "true"
			if !req.AutoProcessReadyResources {
				configs[i].Value = "false"
			}
			break
		}
	}

	// 保存配置
	err = repoManager.SystemConfigRepository.UpsertConfigs(configs)
	if err != nil {
		ErrorResponse(c, "保存系统配置失败", http.StatusInternalServerError)
		return
	}

	// 更新定时任务状态
	scheduler := scheduler.GetGlobalScheduler(
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
		// 获取其他配置值
		autoFetchHotDrama, _ := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoFetchHotDramaEnabled)
		autoTransfer, _ := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)

		scheduler.UpdateSchedulerStatusWithAutoTransfer(
			autoFetchHotDrama,
			req.AutoProcessReadyResources,
			autoTransfer,
		)
	}

	// 返回更新后的配置
	configResponse := converter.SystemConfigToResponse(configs)
	SuccessResponse(c, configResponse)
}
