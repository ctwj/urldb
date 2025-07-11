package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"

	"github.com/gin-gonic/gin"
)

// GetCks 获取Cookie列表
func GetCks(c *gin.Context) {
	cks, err := repoManager.CksRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := converter.ToCksResponseList(cks)
	c.JSON(http.StatusOK, responses)
}

// CreateCks 创建Cookie
func CreateCks(c *gin.Context) {
	var req dto.CreateCksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cks := &entity.Cks{
		PanID:     req.PanID,
		Idx:       req.Idx,
		Ck:        req.Ck,
		IsValid:   req.IsValid,
		Space:     req.Space,
		LeftSpace: req.LeftSpace,
		Remark:    req.Remark,
	}

	err := repoManager.CksRepository.Create(cks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      cks.ID,
		"message": "Cookie创建成功",
	})
}

// UpdateCks 更新Cookie
func UpdateCks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req dto.UpdateCksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cookie不存在"})
		return
	}

	if req.PanID != 0 {
		cks.PanID = req.PanID
	}
	cks.Idx = req.Idx
	if req.Ck != "" {
		cks.Ck = req.Ck
	}
	cks.IsValid = req.IsValid
	cks.Space = req.Space
	cks.LeftSpace = req.LeftSpace
	if req.Remark != "" {
		cks.Remark = req.Remark
	}

	err = repoManager.CksRepository.Update(cks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cookie更新成功"})
}

// DeleteCks 删除Cookie
func DeleteCks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = repoManager.CksRepository.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cookie删除成功"})
}

// GetCksByID 根据ID获取Cookie
func GetCksByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cookie不存在"})
		return
	}

	response := converter.ToCksResponse(cks)
	c.JSON(http.StatusOK, response)
}
