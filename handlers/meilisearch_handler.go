package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
)

// MeilisearchHandler Meilisearch处理器
type MeilisearchHandler struct {
	meilisearchManager *services.MeilisearchManager
}

// NewMeilisearchHandler 创建Meilisearch处理器
func NewMeilisearchHandler(meilisearchManager *services.MeilisearchManager) *MeilisearchHandler {
	return &MeilisearchHandler{
		meilisearchManager: meilisearchManager,
	}
}

// TestConnection 测试Meilisearch连接
func (h *MeilisearchHandler) TestConnection(c *gin.Context) {
	var req struct {
		Host      string      `json:"host"`
		Port      interface{} `json:"port"` // 支持字符串或数字
		MasterKey string      `json:"masterKey"`
		IndexName string      `json:"indexName"` // 可选字段
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 验证必要字段
	if req.Host == "" {
		ErrorResponse(c, "主机地址不能为空", http.StatusBadRequest)
		return
	}

	// 转换port为字符串
	var portStr string
	switch v := req.Port.(type) {
	case string:
		portStr = v
	case float64:
		portStr = strconv.Itoa(int(v))
	case int:
		portStr = strconv.Itoa(v)
	default:
		portStr = "7700" // 默认端口
	}

	// 如果没有提供索引名称，使用默认值
	indexName := req.IndexName
	if indexName == "" {
		indexName = "resources"
	}

	// 创建临时服务进行测试
	service := services.NewMeilisearchService(req.Host, portStr, req.MasterKey, indexName, true)

	if err := service.HealthCheck(); err != nil {
		ErrorResponse(c, "连接测试失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	SuccessResponse(c, gin.H{"message": "连接测试成功"})
}

// GetStatus 获取Meilisearch状态
func (h *MeilisearchHandler) GetStatus(c *gin.Context) {
	if h.meilisearchManager == nil {
		SuccessResponse(c, gin.H{
			"enabled": false,
			"healthy": false,
			"message": "Meilisearch未初始化",
		})
		return
	}

	status, err := h.meilisearchManager.GetStatusWithHealthCheck()
	if err != nil {
		ErrorResponse(c, "获取状态失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, status)
}

// GetUnsyncedCount 获取未同步资源数量
func (h *MeilisearchHandler) GetUnsyncedCount(c *gin.Context) {
	if h.meilisearchManager == nil {
		SuccessResponse(c, gin.H{"count": 0})
		return
	}

	count, err := h.meilisearchManager.GetUnsyncedCount()
	if err != nil {
		ErrorResponse(c, "获取未同步数量失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"count": count})
}

// GetUnsyncedResources 获取未同步的资源
func (h *MeilisearchHandler) GetUnsyncedResources(c *gin.Context) {
	if h.meilisearchManager == nil {
		SuccessResponse(c, gin.H{
			"resources": []interface{}{},
			"total":     0,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resources, total, err := h.meilisearchManager.GetUnsyncedResources(page, pageSize)
	if err != nil {
		ErrorResponse(c, "获取未同步资源失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"resources": converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetSyncedResources 获取已同步的资源
func (h *MeilisearchHandler) GetSyncedResources(c *gin.Context) {
	if h.meilisearchManager == nil {
		SuccessResponse(c, gin.H{
			"resources": []interface{}{},
			"total":     0,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resources, total, err := h.meilisearchManager.GetSyncedResources(page, pageSize)
	if err != nil {
		ErrorResponse(c, "获取已同步资源失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"resources": converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAllResources 获取所有资源
func (h *MeilisearchHandler) GetAllResources(c *gin.Context) {
	if h.meilisearchManager == nil {
		SuccessResponse(c, gin.H{
			"resources": []interface{}{},
			"total":     0,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resources, total, err := h.meilisearchManager.GetAllResources(page, pageSize)
	if err != nil {
		ErrorResponse(c, "获取所有资源失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"resources": converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// SyncAllResources 同步所有资源
func (h *MeilisearchHandler) SyncAllResources(c *gin.Context) {
	if h.meilisearchManager == nil {
		ErrorResponse(c, "Meilisearch未初始化", http.StatusInternalServerError)
		return
	}

	utils.Info("开始同步所有资源到Meilisearch...")

	_, err := h.meilisearchManager.SyncAllResources()
	if err != nil {
		ErrorResponse(c, "同步失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "同步已开始，请查看进度",
	})
}

// GetSyncProgress 获取同步进度
func (h *MeilisearchHandler) GetSyncProgress(c *gin.Context) {
	if h.meilisearchManager == nil {
		ErrorResponse(c, "Meilisearch未初始化", http.StatusInternalServerError)
		return
	}

	progress := h.meilisearchManager.GetSyncProgress()
	SuccessResponse(c, progress)
}

// StopSync 停止同步
func (h *MeilisearchHandler) StopSync(c *gin.Context) {
	if h.meilisearchManager == nil {
		ErrorResponse(c, "Meilisearch未初始化", http.StatusInternalServerError)
		return
	}

	h.meilisearchManager.StopSync()
	SuccessResponse(c, gin.H{
		"message": "同步已停止",
	})
}

// ClearIndex 清空索引
func (h *MeilisearchHandler) ClearIndex(c *gin.Context) {
	if h.meilisearchManager == nil {
		ErrorResponse(c, "Meilisearch未初始化", http.StatusInternalServerError)
		return
	}

	if err := h.meilisearchManager.ClearIndex(); err != nil {
		ErrorResponse(c, "清空索引失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "清空索引成功"})
}

// UpdateIndexSettings 更新索引设置
func (h *MeilisearchHandler) UpdateIndexSettings(c *gin.Context) {
	if h.meilisearchManager == nil {
		ErrorResponse(c, "Meilisearch未初始化", http.StatusInternalServerError)
		return
	}

	service := h.meilisearchManager.GetService()
	if service == nil {
		ErrorResponse(c, "Meilisearch服务未初始化", http.StatusInternalServerError)
		return
	}

	if err := service.UpdateIndexSettings(); err != nil {
		ErrorResponse(c, "更新索引设置失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "索引设置更新成功"})
}
