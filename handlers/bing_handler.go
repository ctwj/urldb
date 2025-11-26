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
	config := gin.H{
		"enabled": h.getConfigValue(entity.BingIndexConfigKeyEnabled, "false") == "true",
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

	// 保存配置
	err := h.setConfigValue(entity.BingIndexConfigKeyEnabled, fmt.Sprintf("%t", request.Enabled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

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
	config := entity.SystemConfig{
		Key:   key,
		Value: value,
		Type:  entity.ConfigTypeString,
	}
	return h.systemConfigRepo.UpsertConfigs([]entity.SystemConfig{config})
}