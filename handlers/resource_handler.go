package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetResources 获取资源列表
func GetResources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryIDStr := c.Query("category_id")
	panIDStr := c.Query("pan_id")

	var resources []entity.Resource
	var err error

	if categoryIDStr != "" {
		categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
		resources, err = repoManager.ResourceRepository.FindByCategoryID(uint(categoryID))
	} else if panIDStr != "" {
		panID, _ := strconv.ParseUint(panIDStr, 10, 32)
		resources, err = repoManager.ResourceRepository.FindByPanID(uint(panID))
	} else {
		resources, err = repoManager.ResourceRepository.FindWithRelations()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回公开的资源
	var publicResources []entity.Resource
	for _, resource := range resources {
		if resource.IsPublic {
			publicResources = append(publicResources, resource)
		}
	}

	// 分页处理
	start := (page - 1) * limit
	end := start + limit
	if start >= len(publicResources) {
		start = len(publicResources)
	}
	if end > len(publicResources) {
		end = len(publicResources)
	}

	pagedResources := publicResources[start:end]
	responses := converter.ToResourceResponseList(pagedResources)

	c.JSON(http.StatusOK, gin.H{
		"resources": responses,
		"page":      page,
		"limit":     limit,
		"total":     len(publicResources),
	})
}

// GetResourceByID 根据ID获取资源
func GetResourceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	if !resource.IsPublic {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 增加浏览次数
	repoManager.ResourceRepository.IncrementViewCount(uint(id))

	response := converter.ToResourceResponse(resource)
	c.JSON(http.StatusOK, response)
}

// CreateResource 创建资源
func CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      resource.ID,
		"message": "资源创建成功",
	})
}

// UpdateResource 更新资源
func UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 更新字段
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

	err = repoManager.ResourceRepository.Update(resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理标签关联
	if req.TagIDs != nil {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "资源更新成功"})
}

// DeleteResource 删除资源
func DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = repoManager.ResourceRepository.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}

// SearchResources 搜索资源
func SearchResources(c *gin.Context) {
	query := c.Query("query")
	categoryIDStr := c.Query("category_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var categoryID *uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			temp := uint(id)
			categoryID = &temp
		}
	}

	// 记录搜索统计
	if query != "" {
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		repoManager.SearchStatRepository.RecordSearch(query, ip, userAgent)
	}

	resources, total, err := repoManager.ResourceRepository.Search(query, categoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回公开的资源
	var publicResources []entity.Resource
	for _, resource := range resources {
		if resource.IsPublic {
			publicResources = append(publicResources, resource)
		}
	}

	responses := converter.ToResourceResponseList(publicResources)

	c.JSON(http.StatusOK, dto.SearchResponse{
		Resources: responses,
		Total:     total,
		Page:      page,
		Limit:     limit,
	})
}
