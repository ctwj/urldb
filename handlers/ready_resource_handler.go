package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetReadyResources 获取待处理资源列表
func GetReadyResources(c *gin.Context) {
	resources, err := repoManager.ReadyResourceRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := converter.ToReadyResourceResponseList(resources)
	c.JSON(http.StatusOK, responses)
}

// CreateReadyResource 创建待处理资源
func CreateReadyResource(c *gin.Context) {
	var req dto.CreateReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource := &entity.ReadyResource{
		Title: req.Title,
		URL:   req.URL,
		IP:    req.IP,
	}

	err := repoManager.ReadyResourceRepository.Create(resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      resource.ID,
		"message": "待处理资源创建成功",
	})
}

// BatchCreateReadyResources 批量创建待处理资源
func BatchCreateReadyResources(c *gin.Context) {
	var req dto.BatchCreateReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resources []entity.ReadyResource
	for _, reqResource := range req.Resources {
		resource := entity.ReadyResource{
			Title: reqResource.Title,
			URL:   reqResource.URL,
			IP:    reqResource.IP,
		}
		resources = append(resources, resource)
	}

	err := repoManager.ReadyResourceRepository.BatchCreate(resources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "批量创建成功",
		"count":   len(resources),
	})
}

// CreateReadyResourcesFromText 从文本创建待处理资源
func CreateReadyResourcesFromText(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文本内容不能为空"})
		return
	}

	lines := strings.Split(text, "\n")
	var resources []entity.ReadyResource

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 简单的URL提取逻辑
		if strings.Contains(line, "http") {
			resource := entity.ReadyResource{
				URL: line,
			}
			resources = append(resources, resource)
		}
	}

	if len(resources) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到有效的URL"})
		return
	}

	err := repoManager.ReadyResourceRepository.BatchCreate(resources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "从文本创建成功",
		"count":   len(resources),
	})
}

// DeleteReadyResource 删除待处理资源
func DeleteReadyResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = repoManager.ReadyResourceRepository.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "待处理资源删除成功"})
}

// ClearReadyResources 清空所有待处理资源
func ClearReadyResources(c *gin.Context) {
	resources, err := repoManager.ReadyResourceRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, resource := range resources {
		err = repoManager.ReadyResourceRepository.Delete(resource.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "所有待处理资源已清空",
		"count":   len(resources),
	})
}
