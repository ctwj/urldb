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
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 1000 {
		pageSize = 100
	}

	// 获取分页数据
	resources, total, err := repoManager.ReadyResourceRepository.FindWithPagination(page, pageSize)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToReadyResourceResponseList(resources)

	// 使用标准化的分页响应格式
	SuccessResponse(c, gin.H{
		"data":      responses,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// CreateReadyResource 创建待处理资源
func CreateReadyResource(c *gin.Context) {
	var req dto.CreateReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource := &entity.ReadyResource{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    req.Category,
		Tags:        req.Tags,
		Img:         req.Img,
		Source:      req.Source,
		Extra:       req.Extra,
		IP:          req.IP,
	}

	err := repoManager.ReadyResourceRepository.Create(resource)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"id":      resource.ID,
		"message": "待处理资源创建成功",
	})
}

// BatchCreateReadyResources 批量创建待处理资源
func BatchCreateReadyResources(c *gin.Context) {
	var req dto.BatchCreateReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	var resources []entity.ReadyResource
	for _, reqResource := range req.Resources {
		resource := entity.ReadyResource{
			Title:       reqResource.Title,
			Description: reqResource.Description,
			URL:         reqResource.URL,
			Category:    reqResource.Category,
			Tags:        reqResource.Tags,
			Img:         reqResource.Img,
			Source:      reqResource.Source,
			Extra:       reqResource.Extra,
			IP:          reqResource.IP,
		}
		resources = append(resources, resource)
	}

	err := repoManager.ReadyResourceRepository.BatchCreate(resources)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"count":   len(resources),
		"message": "批量创建成功",
	})
}

// CreateReadyResourcesFromText 从文本创建待处理资源
func CreateReadyResourcesFromText(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		ErrorResponse(c, "文本内容不能为空", http.StatusBadRequest)
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
		ErrorResponse(c, "未找到有效的URL", http.StatusBadRequest)
		return
	}

	err := repoManager.ReadyResourceRepository.BatchCreate(resources)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"count":   len(resources),
		"message": "从文本创建成功",
	})
}

// DeleteReadyResource 删除待处理资源
func DeleteReadyResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.ReadyResourceRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "待处理资源删除成功"})
}

// ClearReadyResources 清空所有待处理资源
func ClearReadyResources(c *gin.Context) {
	resources, err := repoManager.ReadyResourceRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, resource := range resources {
		err = repoManager.ReadyResourceRepository.Delete(resource.ID)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{
		"deleted_count": len(resources),
		"message":       "所有待处理资源已清空",
	})
}
