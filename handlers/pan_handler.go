package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctwj/panResManage/db/converter"
	"github.com/ctwj/panResManage/db/dto"
	"github.com/ctwj/panResManage/db/entity"

	"github.com/gin-gonic/gin"
)

// GetPans 获取平台列表
func GetPans(c *gin.Context) {
	pans, err := repoManager.PanRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToPanResponseList(pans)
	ListResponse(c, responses, int64(len(responses)))
}

// CreatePan 创建平台
func CreatePan(c *gin.Context) {
	var req dto.CreatePanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"id":      pan.ID,
		"message": "平台创建成功",
	})
}

// UpdatePan 更新平台
func UpdatePan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdatePanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	pan, err := repoManager.PanRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "平台不存在", http.StatusNotFound)
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
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "平台更新成功"})
}

// DeletePan 删除平台
func DeletePan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.PanRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "平台删除成功"})
}

// GetPan 根据ID获取平台
func GetPan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	pan, err := repoManager.PanRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "平台不存在", http.StatusNotFound)
		return
	}

	response := converter.ToPanResponse(pan)
	SuccessResponse(c, response)
}
