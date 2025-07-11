package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetPans 获取平台列表
func GetPans(c *gin.Context) {
	pans, err := repoManager.PanRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := converter.ToPanResponseList(pans)
	c.JSON(http.StatusOK, responses)
}

// CreatePan 创建平台
func CreatePan(c *gin.Context) {
	var req dto.CreatePanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pan := &entity.Pan{
		Name:   req.Name,
		Key:    req.Key,
		Icon:   req.Icon,
		Remark: req.Remark,
	}

	err := repoManager.PanRepository.Create(pan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      pan.ID,
		"message": "平台创建成功",
	})
}

// UpdatePan 更新平台
func UpdatePan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req dto.UpdatePanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pan, err := repoManager.PanRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "平台不存在"})
		return
	}

	if req.Name != "" {
		pan.Name = req.Name
	}
	pan.Key = req.Key
	if req.Icon != "" {
		pan.Icon = req.Icon
	}
	if req.Remark != "" {
		pan.Remark = req.Remark
	}

	err = repoManager.PanRepository.Update(pan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "平台更新成功"})
}

// DeletePan 删除平台
func DeletePan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = repoManager.PanRepository.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "平台删除成功"})
}

// GetPan 根据ID获取平台
func GetPan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	pan, err := repoManager.PanRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "平台不存在"})
		return
	}

	response := converter.ToPanResponse(pan)
	c.JSON(http.StatusOK, response)
}
