package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	categories, err := repoManager.CategoryRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个分类的资源数量
	resourceCounts := make(map[uint]int64)
	for _, category := range categories {
		count, err := repoManager.CategoryRepository.GetResourceCount(category.ID)
		if err != nil {
			continue
		}
		resourceCounts[category.ID] = count
	}

	responses := converter.ToCategoryResponseList(categories, resourceCounts)
	SuccessResponse(c, responses)
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
		"category": converter.ToCategoryResponse(category, 0),
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
