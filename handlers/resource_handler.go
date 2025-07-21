package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"

	"github.com/gin-gonic/gin"
)

// GetResources 获取资源列表
func GetResources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	if search := c.Query("search"); search != "" {
		params["search"] = search
	}
	if panID := c.Query("pan_id"); panID != "" {
		if id, err := strconv.ParseUint(panID, 10, 32); err == nil {
			params["pan_id"] = uint(id)
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			params["category_id"] = uint(id)
		}
	}

	resources, total, err := repoManager.ResourceRepository.SearchWithFilters(params)

	// 搜索统计（仅非管理员）
	if search, ok := params["search"].(string); ok && search != "" {
		user, _ := c.Get("user")
		if user == nil || (user != nil && user.(entity.User).Role != "admin") {
			ip := c.ClientIP()
			userAgent := c.GetHeader("User-Agent")
			repoManager.SearchStatRepository.RecordSearch(search, ip, userAgent)
		}
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"data":      converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetResourceByID 根据ID获取资源
func GetResourceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	response := converter.ToResourceResponse(resource)
	SuccessResponse(c, response)
}

// CheckResourceExists 检查资源是否存在（测试FindExists函数）
func CheckResourceExists(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		ErrorResponse(c, "URL参数不能为空", http.StatusBadRequest)
		return
	}

	excludeIDStr := c.Query("exclude_id")
	var excludeID uint
	if excludeIDStr != "" {
		if id, err := strconv.ParseUint(excludeIDStr, 10, 32); err == nil {
			excludeID = uint(id)
		}
	}

	exists, err := repoManager.ResourceRepository.FindExists(url, excludeID)
	if err != nil {
		ErrorResponse(c, "检查失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"url":    url,
		"exists": exists,
	})
}

// CreateResource 创建资源
func CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource := &entity.Resource{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		PanID:       req.PanID,
		QuarkURL:    req.QuarkURL,
		FileSize:    req.FileSize,
		CategoryID:  req.CategoryID,
		IsValid:     req.IsValid,
		IsPublic:    req.IsPublic,
	}

	err := repoManager.ResourceRepository.Create(resource)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{
		"message":  "资源创建成功",
		"resource": converter.ToResourceResponse(resource),
	})
}

// UpdateResource 更新资源
func UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// 更新资源信息
	if req.Title != "" {
		resource.Title = req.Title
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.URL != "" {
		resource.URL = req.URL
	}
	if req.PanID != nil {
		resource.PanID = req.PanID
	}
	if req.QuarkURL != "" {
		resource.QuarkURL = req.QuarkURL
	}
	if req.FileSize != "" {
		resource.FileSize = req.FileSize
	}
	if req.CategoryID != nil {
		resource.CategoryID = req.CategoryID
	}
	resource.IsValid = req.IsValid
	resource.IsPublic = req.IsPublic

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = repoManager.ResourceRepository.Update(resource)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{"message": "资源更新成功"})
}

// DeleteResource 删除资源
func DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.ResourceRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "资源删除成功"})
}

// SearchResources 搜索资源
func SearchResources(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var resources []entity.Resource
	var total int64
	var err error

	if query == "" {
		// 搜索关键词为空时，返回最新记录（分页）
		resources, total, err = repoManager.ResourceRepository.FindWithRelationsPaginated(page, pageSize)
	} else {
		// 有搜索关键词时，执行搜索
		resources, total, err = repoManager.ResourceRepository.Search(query, nil, page, pageSize)
		// 新增：记录搜索关键词
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		repoManager.SearchStatRepository.RecordSearch(query, ip, userAgent)
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"resources": converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 增加资源浏览次数
func IncrementResourceViewCount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}
	err = repoManager.ResourceRepository.IncrementViewCount(uint(id))
	if err != nil {
		ErrorResponse(c, "增加浏览次数失败", http.StatusInternalServerError)
		return
	}
	SuccessResponse(c, gin.H{"message": "浏览次数+1"})
}

// BatchDeleteResources 批量删除资源
func BatchDeleteResources(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		ErrorResponse(c, "参数错误", 400)
		return
	}
	count := 0
	for _, id := range req.IDs {
		if err := repoManager.ResourceRepository.Delete(id); err == nil {
			count++
		}
	}
	SuccessResponse(c, gin.H{"deleted": count, "message": "批量删除成功"})
}
