package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"
	"res_db/db/entity"
	"res_db/db/repo"

	"github.com/gin-gonic/gin"
)

// HotDramaHandler 热播剧处理器
type HotDramaHandler struct {
	hotDramaRepo repo.HotDramaRepository
}

// NewHotDramaHandler 创建热播剧处理器
func NewHotDramaHandler(hotDramaRepo repo.HotDramaRepository) *HotDramaHandler {
	return &HotDramaHandler{
		hotDramaRepo: hotDramaRepo,
	}
}

// GetHotDramaList 获取热播剧列表
func (h *HotDramaHandler) GetHotDramaList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	var dramas []entity.HotDrama
	var total int64
	var err error

	if category != "" {
		dramas, total, err = h.hotDramaRepo.FindByCategory(category, page, pageSize)
	} else {
		dramas, total, err = h.hotDramaRepo.FindAll(page, pageSize)
	}

	if err != nil {
		ErrorResponse(c, "获取热播剧列表失败", http.StatusInternalServerError)
		return
	}

	response := converter.HotDramaListToResponse(dramas)
	response.Total = int(total)

	SuccessResponse(c, response)
}

// GetHotDramaByID 根据ID获取热播剧详情
func (h *HotDramaHandler) GetHotDramaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	drama, err := h.hotDramaRepo.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "热播剧不存在", http.StatusNotFound)
		return
	}

	response := converter.HotDramaToResponse(drama)
	SuccessResponse(c, response)
}

// CreateHotDrama 创建热播剧记录
func CreateHotDrama(c *gin.Context) {
	var req dto.HotDramaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	drama := converter.RequestToHotDrama(&req)
	if drama == nil {
		ErrorResponse(c, "数据转换失败", http.StatusInternalServerError)
		return
	}

	err := repoManager.HotDramaRepository.Create(drama)
	if err != nil {
		ErrorResponse(c, "创建热播剧记录失败", http.StatusInternalServerError)
		return
	}

	response := converter.HotDramaToResponse(drama)
	SuccessResponse(c, response)
}

// UpdateHotDrama 更新热播剧记录
func UpdateHotDrama(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.HotDramaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	drama := converter.RequestToHotDrama(&req)
	if drama == nil {
		ErrorResponse(c, "数据转换失败", http.StatusInternalServerError)
		return
	}
	drama.ID = uint(id)

	err = repoManager.HotDramaRepository.Upsert(drama)
	if err != nil {
		ErrorResponse(c, "更新热播剧记录失败", http.StatusInternalServerError)
		return
	}

	response := converter.HotDramaToResponse(drama)
	SuccessResponse(c, response)
}

// DeleteHotDrama 删除热播剧记录
func DeleteHotDrama(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.HotDramaRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, "删除热播剧记录失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "删除热播剧记录成功"})
}

// GetHotDramaList 获取热播剧列表（使用全局repoManager）
func GetHotDramaList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	var dramas []entity.HotDrama
	var total int64
	var err error

	// 如果page_size很大（比如>=1000），则获取所有数据
	if pageSize >= 1000 {
		if category != "" {
			dramas, total, err = repoManager.HotDramaRepository.FindByCategory(category, 1, 10000)
		} else {
			dramas, total, err = repoManager.HotDramaRepository.FindAll(1, 10000)
		}
	} else {
		if category != "" {
			dramas, total, err = repoManager.HotDramaRepository.FindByCategory(category, page, pageSize)
		} else {
			dramas, total, err = repoManager.HotDramaRepository.FindAll(page, pageSize)
		}
	}

	if err != nil {
		ErrorResponse(c, "获取热播剧列表失败", http.StatusInternalServerError)
		return
	}

	response := converter.HotDramaListToResponse(dramas)
	response.Total = int(total)

	SuccessResponse(c, response)
}

// GetHotDramaByID 根据ID获取热播剧详情（使用全局repoManager）
func GetHotDramaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	drama, err := repoManager.HotDramaRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "热播剧不存在", http.StatusNotFound)
		return
	}

	response := converter.HotDramaToResponse(drama)
	SuccessResponse(c, response)
}
