package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// BingHandler Bing相关处理器
type BingHandler struct {
	systemConfigRepo  repo.SystemConfigRepository
}

// NewBingHandler 创建Bing处理器
func NewBingHandler(siteURL string, repoManager *repo.RepositoryManager) *BingHandler {
	return &BingHandler{
		systemConfigRepo: repoManager.SystemConfigRepository,
	}
}

// GetBingIndexConfig 获取Bing索引配置
func (h *BingHandler) GetBingIndexConfig(c *gin.Context) {
	enabledValue := h.getConfigValue(entity.BingIndexConfigKeyEnabled, "false")
	enabled := enabledValue == "true"
	
	fmt.Printf("[Bing] 获取配置 - 原始值: %s, 转换后: %v\n", enabledValue, enabled)
	
	config := gin.H{
		"enabled": enabled,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// UpdateBingIndexConfig 更新Bing索引配置
func (h *BingHandler) UpdateBingIndexConfig(c *gin.Context) {
	var request struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	fmt.Printf("[Bing] 更新配置 - 接收到的值: %v\n", request.Enabled)

	// 保存配置
	valueToSave := fmt.Sprintf("%t", request.Enabled)
	fmt.Printf("[Bing] 保存配置 - 键: %s, 值: %s\n", entity.BingIndexConfigKeyEnabled, valueToSave)
	
	err := h.setConfigValue(entity.BingIndexConfigKeyEnabled, valueToSave)
	if err != nil {
		fmt.Printf("[Bing] 保存配置失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	fmt.Printf("[Bing] 配置保存成功\n")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "配置更新成功",
	})
}

// getConfigValue 获取配置值，返回默认值如果配置不存在
func (h *BingHandler) getConfigValue(key string, defaultValue string) string {
	value, err := h.systemConfigRepo.GetConfigValue(key)
	if err != nil || value == "" {
		return defaultValue
	}
	return value
}

// setConfigValue 保存配置值
func (h *BingHandler) setConfigValue(key string, value string) error {
	// 根据key确定配置类型
	configType := entity.ConfigTypeString
	if key == entity.BingIndexConfigKeyEnabled {
		configType = entity.ConfigTypeBool
	}
	
	config := entity.SystemConfig{
		Key:   key,
		Value: value,
		Type:  configType,
	}
	return h.systemConfigRepo.UpsertConfigs([]entity.SystemConfig{config})
}