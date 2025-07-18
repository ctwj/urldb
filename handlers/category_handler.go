package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	var categories []entity.Category
	var total int64
	var err error

	if search != "" {
		// 搜索分类
		categories, total, err = repoManager.CategoryRepository.Search(search, page, pageSize)
	} else {
		// 分页查询
		categories, total, err = repoManager.CategoryRepository.FindWithPagination(page, pageSize)
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个分类的资源数量和标签名称
	resourceCounts := make(map[uint]int64)
	tagNamesMap := make(map[uint][]string)
	for _, category := range categories {
		// 获取资源数量
		resourceCount, err := repoManager.CategoryRepository.GetResourceCount(category.ID)
		if err != nil {
			continue
		}
		resourceCounts[category.ID] = resourceCount

		// 获取标签名称
		tagNames, err := repoManager.CategoryRepository.GetTagNames(category.ID)
		if err != nil {
			continue
		}
		tagNamesMap[category.ID] = tagNames
	}

	responses := converter.ToCategoryResponseList(categories, resourceCounts, tagNamesMap)

	// 返回分页格式的响应
	SuccessResponse(c, gin.H{
		"items":     responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	err := repoManager.CategoryRepository.Create(category)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message":  "分类创建成功",
		"category": converter.ToCategoryResponse(category, 0, []string{}),
	})
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := repoManager.CategoryRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "分类不存在", http.StatusNotFound)
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	err = repoManager.CategoryRepository.Update(category)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "分类更新成功"})
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.CategoryRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "分类删除成功"})
}
