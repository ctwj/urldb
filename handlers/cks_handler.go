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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToCksResponseList(cks)
	SuccessResponse(c, responses)
}

// CreateCks 创建Cookie
func CreateCks(c *gin.Context) {
	var req dto.CreateCksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "Cookie创建成功",
		"cks":     converter.ToCksResponse(cks),
	})
}

// GetCksByID 根据ID获取Cookie详情
func GetCksByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "Cookie不存在", http.StatusNotFound)
		return
	}

	response := converter.ToCksResponse(cks)
	SuccessResponse(c, response)
}

// UpdateCks 更新Cookie
func UpdateCks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "Cookie不存在", http.StatusNotFound)
		return
	}

	if req.PanID != 0 {
		cks.PanID = req.PanID
	}
	if req.Idx != 0 {
		cks.Idx = req.Idx
	}
	if req.Ck != "" {
		cks.Ck = req.Ck
	}
	cks.IsValid = req.IsValid
	if req.Space != 0 {
		cks.Space = req.Space
	}
	if req.LeftSpace != 0 {
		cks.LeftSpace = req.LeftSpace
	}
	if req.Remark != "" {
		cks.Remark = req.Remark
	}

	err = repoManager.CksRepository.Update(cks)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "Cookie更新成功"})
}

// DeleteCks 删除Cookie
func DeleteCks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.CksRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "Cookie删除成功"})
}

// GetCksByID 根据ID获取Cookie详情（使用全局repoManager）
func GetCksByIDGlobal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "Cookie不存在", http.StatusNotFound)
		return
	}

	response := converter.ToCksResponse(cks)
	SuccessResponse(c, response)
}
