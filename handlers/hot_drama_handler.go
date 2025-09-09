package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/go-resty/resty/v2"

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

// GetPosterImage 获取海报图片代理
func GetPosterImage(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		ErrorResponse(c, "图片URL不能为空", http.StatusBadRequest)
		return
	}

	// 简单的URL验证
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		ErrorResponse(c, "无效的图片URL", http.StatusBadRequest)
		return
	}

	// 检查If-Modified-Since头，实现条件请求
	ifModifiedSince := c.GetHeader("If-Modified-Since")
	if ifModifiedSince != "" {
		// 如果存在，说明浏览器有缓存，检查是否过期
		ifLastModified, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", ifModifiedSince)
		if err == nil && time.Since(ifLastModified) < 86400*time.Second { // 24小时内
			c.Status(http.StatusNotModified)
			return
		}
	}

	// 检查ETag头 - 基于URL生成，保证相同URL有相同ETag
	ifNoneMatch := c.GetHeader("If-None-Match")
	if ifNoneMatch != "" {
		etag := fmt.Sprintf(`"%x"`, len(url)) // 简单的基于URL长度的ETag
		if ifNoneMatch == etag {
			c.Status(http.StatusNotModified)
			return
		}
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	resp, err := client.R().
		SetHeaders(map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Referer":         "https://m.douban.com/",
			"Accept":          "image/webp,image/apng,image/*,*/*;q=0.8",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		}).
		Get(url)

	if err != nil {
		ErrorResponse(c, "获取图片失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode() != 200 {
		ErrorResponse(c, fmt.Sprintf("获取图片失败，状态码: %d", resp.StatusCode()), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	contentType := resp.Header().Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}
	c.Header("Content-Type", contentType)

	// 增强缓存策略
	c.Header("Cache-Control", "public, max-age=604800, s-maxage=86400") // 客户端7天，代理1天
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	// 设置缓存验证头（基于URL长度生成的简单ETag）
	etag := fmt.Sprintf(`"%x"`, len(url))
	c.Header("ETag", etag)
	c.Header("Last-Modified", time.Now().Add(-86400*time.Second).Format("Mon, 02 Jan 2006 15:04:05 GMT")) // 设为1天前，避免立即过期

	// 返回图片数据
	c.Data(resp.StatusCode(), contentType, resp.Body())
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
	subType := c.Query("sub_type")

	var dramas []entity.HotDrama
	var total int64
	var err error

	// 如果page_size很大（比如>=1000），则获取所有数据
	if pageSize >= 1000 {
		if category != "" && subType != "" {
			dramas, total, err = repoManager.HotDramaRepository.FindByCategoryAndSubType(category, subType, 1, 10000)
		} else if category != "" {
			dramas, total, err = repoManager.HotDramaRepository.FindByCategory(category, 1, 10000)
		} else {
			dramas, total, err = repoManager.HotDramaRepository.FindAll(1, 10000)
		}
	} else {
		if category != "" && subType != "" {
			dramas, total, err = repoManager.HotDramaRepository.FindByCategoryAndSubType(category, subType, page, pageSize)
		} else if category != "" {
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
