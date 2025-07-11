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
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")
	search := c.Query("search")

	var resources []entity.Resource
	var total int64
	var err error

	// 设置响应头，启用缓存
	c.Header("Cache-Control", "public, max-age=300") // 5分钟缓存

	if search != "" {
		resources, total, err = repoManager.ResourceRepository.Search(search, nil, page, pageSize)
	} else if categoryID != "" {
		categoryIDUint, _ := strconv.ParseUint(categoryID, 10, 32)
		resources, total, err = repoManager.ResourceRepository.FindByCategoryIDPaginated(uint(categoryIDUint), page, pageSize)
	} else {
		// 使用分页查询，避免加载所有数据
		resources, total, err = repoManager.ResourceRepository.FindWithRelationsPaginated(page, pageSize)
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

// GetResourceByID 根据ID获取资源详情
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

	// 增加浏览次数
	if resource != nil {
		repoManager.ResourceRepository.IncrementViewCount(uint(id))
	}

	response := converter.ToResourceResponse(resource)
	SuccessResponse(c, response)
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
