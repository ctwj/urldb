package handlers

import (
	"net/http"
	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/repo"

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
		ErrorResponse(c, http.StatusInternalServerError, "获取系统配置失败")
		return
	}

	configResponse := converter.SystemConfigToResponse(config)
	SuccessResponse(c, configResponse, "获取系统配置成功")
}

// UpdateConfig 更新系统配置
func (h *SystemConfigHandler) UpdateConfig(c *gin.Context) {
	var req dto.SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 验证参数
	if req.SiteTitle == "" {
		ErrorResponse(c, http.StatusBadRequest, "网站标题不能为空")
		return
	}

	if req.AutoProcessInterval < 1 || req.AutoProcessInterval > 1440 {
		ErrorResponse(c, http.StatusBadRequest, "自动处理间隔必须在1-1440分钟之间")
		return
	}

	if req.PageSize < 10 || req.PageSize > 500 {
		ErrorResponse(c, http.StatusBadRequest, "每页显示数量必须在10-500之间")
		return
	}

	// 转换为实体
	config := converter.RequestToSystemConfig(&req)
	if config == nil {
		ErrorResponse(c, http.StatusInternalServerError, "数据转换失败")
		return
	}

	// 保存配置
	err := h.systemConfigRepo.Upsert(config)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "保存系统配置失败")
		return
	}

	// 返回更新后的配置
	updatedConfig, err := h.systemConfigRepo.FindFirst()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取更新后的配置失败")
		return
	}

	configResponse := converter.SystemConfigToResponse(updatedConfig)
	SuccessResponse(c, configResponse, "系统配置保存成功")
}

// GetSystemConfig 获取系统配置（使用全局repoManager）
func GetSystemConfig(c *gin.Context) {
	config, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取系统配置失败")
		return
	}

	configResponse := converter.SystemConfigToResponse(config)
	SuccessResponse(c, configResponse, "获取系统配置成功")
}

// UpdateSystemConfig 更新系统配置（使用全局repoManager）
func UpdateSystemConfig(c *gin.Context) {
	var req dto.SystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 验证参数
	if req.SiteTitle == "" {
		ErrorResponse(c, http.StatusBadRequest, "网站标题不能为空")
		return
	}

	if req.AutoProcessInterval < 1 || req.AutoProcessInterval > 1440 {
		ErrorResponse(c, http.StatusBadRequest, "自动处理间隔必须在1-1440分钟之间")
		return
	}

	if req.PageSize < 10 || req.PageSize > 500 {
		ErrorResponse(c, http.StatusBadRequest, "每页显示数量必须在10-500之间")
		return
	}

	// 转换为实体
	config := converter.RequestToSystemConfig(&req)
	if config == nil {
		ErrorResponse(c, http.StatusInternalServerError, "数据转换失败")
		return
	}

	// 保存配置
	err := repoManager.SystemConfigRepository.Upsert(config)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "保存系统配置失败")
		return
	}

	// 返回更新后的配置
	updatedConfig, err := repoManager.SystemConfigRepository.FindFirst()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取更新后的配置失败")
		return
	}

	configResponse := converter.SystemConfigToResponse(updatedConfig)
	SuccessResponse(c, configResponse, "系统配置保存成功")
}
