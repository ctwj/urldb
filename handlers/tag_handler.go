package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	tags, err := repoManager.TagRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToTagResponseList(tags)
	SuccessResponse(c, responses)
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
	}

	err := repoManager.TagRepository.Create(tag)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "标签创建成功",
		"tag":     converter.ToTagResponse(tag),
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

	response := converter.ToTagResponse(tag)
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

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Description != "" {
		tag.Description = req.Description
	}

	err = repoManager.TagRepository.Update(tag)
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

	response := converter.ToTagResponse(tag)
	SuccessResponse(c, response)
}

// GetTags 获取标签列表（使用全局repoManager）
func GetTagsGlobal(c *gin.Context) {
	tags, err := repoManager.TagRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToTagResponseList(tags)
	SuccessResponse(c, responses)
}
