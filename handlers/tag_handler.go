package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	var tags []entity.Tag
	var total int64
	var err error

	if search != "" {
		// 搜索标签
		tags, total, err = repoManager.TagRepository.Search(search, page, pageSize)
	} else {
		// 分页查询
		tags, total, err = repoManager.TagRepository.FindWithPagination(page, pageSize)
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个标签的资源数量
	resourceCounts := make(map[uint]int64)
	for _, tag := range tags {
		count, err := repoManager.TagRepository.GetResourceCount(tag.ID)
		if err != nil {
			continue
		}
		resourceCounts[tag.ID] = count
	}

	responses := converter.ToTagResponseList(tags, resourceCounts)

	// 返回分页格式的响应
	SuccessResponse(c, gin.H{
		"items":     responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	tag := &entity.Tag{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	}

	err := repoManager.TagRepository.Create(tag)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "标签创建成功",
		"tag":     converter.ToTagResponse(tag, 0),
	})
}

// GetTagByID 根据ID获取标签详情
func GetTagByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	tag, err := repoManager.TagRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "标签不存在", http.StatusNotFound)
		return
	}

	// 获取资源数量
	resourceCount, _ := repoManager.TagRepository.GetResourceCount(tag.ID)
	response := converter.ToTagResponse(tag, resourceCount)
	SuccessResponse(c, response)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	tag, err := repoManager.TagRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "标签不存在", http.StatusNotFound)
		return
	}

	// 添加调试信息
	fmt.Printf("更新标签 - ID: %d, 请求CategoryID: %v, 当前CategoryID: %v\n", id, req.CategoryID, tag.CategoryID)

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Description != "" {
		tag.Description = req.Description
	}
	// 处理CategoryID，包括设置为null的情况
	tag.CategoryID = req.CategoryID

	// 添加调试信息
	fmt.Printf("更新后CategoryID: %v\n", tag.CategoryID)

	// 使用专门的更新方法，确保能更新null值
	err = repoManager.TagRepository.UpdateWithNulls(tag)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "标签更新成功"})
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.TagRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "标签删除成功"})
}

// GetTagByID 根据ID获取标签详情（使用全局repoManager）
func GetTagByIDGlobal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	tag, err := repoManager.TagRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "标签不存在", http.StatusNotFound)
		return
	}

	// 获取资源数量
	resourceCount, _ := repoManager.TagRepository.GetResourceCount(tag.ID)
	response := converter.ToTagResponse(tag, resourceCount)
	SuccessResponse(c, response)
}

// GetTags 获取标签列表（使用全局repoManager）
func GetTagsGlobal(c *gin.Context) {
	tags, err := repoManager.TagRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个标签的资源数量
	resourceCounts := make(map[uint]int64)
	for _, tag := range tags {
		count, err := repoManager.TagRepository.GetResourceCount(tag.ID)
		if err != nil {
			continue
		}
		resourceCounts[tag.ID] = count
	}

	responses := converter.ToTagResponseList(tags, resourceCounts)
	SuccessResponse(c, responses)
}

// GetTagsByCategory 根据分类ID获取标签列表
func GetTagsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的分类ID", http.StatusBadRequest)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tags, total, err := repoManager.TagRepository.FindByCategoryIDPaginated(uint(categoryID), page, pageSize)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个标签的资源数量
	resourceCounts := make(map[uint]int64)
	for _, tag := range tags {
		count, err := repoManager.TagRepository.GetResourceCount(tag.ID)
		if err != nil {
			continue
		}
		resourceCounts[tag.ID] = count
	}

	responses := converter.ToTagResponseList(tags, resourceCounts)

	// 返回分页格式的响应
	SuccessResponse(c, gin.H{
		"items":     responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
