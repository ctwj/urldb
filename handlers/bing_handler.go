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

	apiKeyValue := h.getConfigValue(entity.BingIndexConfigKeyAPIKey, "")

	fmt.Printf("[Bing] 获取配置 - enabled: %v, apiKey: %s\n", enabled, apiKeyValue)

	config := gin.H{
		"enabled": enabled,
		"apiKey":  apiKeyValue,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// UpdateBingIndexConfig 更新Bing索引配置
func (h *BingHandler) UpdateBingIndexConfig(c *gin.Context) {
	var request struct {
		Enabled bool   `json:"enabled"`
		APIKey  string `json:"apiKey"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	fmt.Printf("[Bing] 更新配置 - enabled: %v, apiKey: %s\n",
		request.Enabled, request.APIKey)

	// 准备要保存的配置
	configs := []entity.SystemConfig{
		{
			Key:   entity.BingIndexConfigKeyEnabled,
			Value: fmt.Sprintf("%t", request.Enabled),
			Type:  entity.ConfigTypeBool,
		},
	}

	// 如果提供了API密钥，则保存
	if request.APIKey != "" {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.BingIndexConfigKeyAPIKey,
			Value: request.APIKey,
			Type:  entity.ConfigTypeString,
		})
	}

	// 批量保存配置
	err := h.systemConfigRepo.UpsertConfigs(configs)
	if err != nil {
		fmt.Printf("[Bing] 保存配置失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	// 清除配置缓存，确保下次读取时从数据库获取最新值
	h.systemConfigRepo.ClearConfigCache()
	fmt.Printf("[Bing] 配置保存成功，缓存已清除\n")

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