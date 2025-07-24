package handlers

import (
	"strconv"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"

	"github.com/gin-gonic/gin"
)

// PublicAPIHandler 公开API处理器
type PublicAPIHandler struct{}

// NewPublicAPIHandler 创建公开API处理器
func NewPublicAPIHandler() *PublicAPIHandler {
	return &PublicAPIHandler{}
}

// AddBatchResources godoc
// @Summary 批量添加资源
// @Description 通过公开API批量添加多个资源到待处理列表
// @Tags PublicAPI
// @Accept json
// @Produce json
// @Param X-API-Token header string true "API访问令牌"
// @Param data body dto.BatchReadyResourceRequest true "批量资源信息"
// @Success 200 {object} map[string]interface{} "批量添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/public/resources/batch-add [post]
func (h *PublicAPIHandler) AddBatchResources(c *gin.Context) {
	var req dto.BatchReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	if len(req.Resources) == 0 {
		ErrorResponse(c, "资源列表不能为空", 400)
		return
	}

	// 收集所有待提交的URL，去重
	urlSet := make(map[string]struct{})
	for _, resource := range req.Resources {
		for _, u := range resource.Url {
			if u != "" {
				urlSet[u] = struct{}{}
			}
		}
	}
	uniqueUrls := make([]string, 0, len(urlSet))
	for url := range urlSet {
		uniqueUrls = append(uniqueUrls, url)
	}

	// 批量查重
	readyList, _ := repoManager.ReadyResourceRepository.BatchFindByURLs(uniqueUrls)
	existReadyUrls := make(map[string]struct{})
	for _, r := range readyList {
		existReadyUrls[r.URL] = struct{}{}
	}
	resourceList, _ := repoManager.ResourceRepository.BatchFindByURLs(uniqueUrls)
	existResourceUrls := make(map[string]struct{})
	for _, r := range resourceList {
		existResourceUrls[r.URL] = struct{}{}
	}

	var createdResources []uint
	for _, resourceReq := range req.Resources {
		// 生成 key（每组同一个 key）
		key, err := repoManager.ReadyResourceRepository.GenerateUniqueKey()
		if err != nil {
			ErrorResponse(c, "生成资源组标识失败: "+err.Error(), 500)
			return
		}
		for _, url := range resourceReq.Url {
			if url == "" {
				continue
			}
			if _, ok := existReadyUrls[url]; ok {
				continue
			}
			if _, ok := existResourceUrls[url]; ok {
				continue
			}
			readyResource := entity.ReadyResource{
				Title:       &resourceReq.Title,
				Description: resourceReq.Description,
				URL:         url,
				Category:    resourceReq.Category,
				Tags:        resourceReq.Tags,
				Img:         resourceReq.Img,
				Source:      "api",
				Extra:       resourceReq.Extra,
				Key:         key,
			}
			err := repoManager.ReadyResourceRepository.Create(&readyResource)
			if err == nil {
				createdResources = append(createdResources, readyResource.ID)
			}
		}
	}

	SuccessResponse(c, gin.H{
		"created_count": len(createdResources),
		"created_ids":   createdResources,
	})
}

// SearchResources godoc
// @Summary 资源搜索
// @Description 搜索资源，支持关键词、标签、分类过滤
// @Tags PublicAPI
// @Accept json
// @Produce json
// @Param X-API-Token header string true "API访问令牌"
// @Param keyword query string false "搜索关键词"
// @Param tag query string false "标签过滤"
// @Param category query string false "分类过滤"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20) maximum(100)
// @Success 200 {object} map[string]interface{} "搜索成功"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/public/resources/search [get]
func (h *PublicAPIHandler) SearchResources(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	tag := c.Query("tag")
	category := c.Query("category")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建搜索条件
	params := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	if keyword != "" {
		params["search"] = keyword
	}

	if tag != "" {
		params["tag"] = tag
	}

	if category != "" {
		params["category"] = category
	}

	// 执行搜索
	resources, total, err := repoManager.ResourceRepository.SearchWithFilters(params)
	if err != nil {
		ErrorResponse(c, "搜索失败: "+err.Error(), 500)
		return
	}

	// 转换为响应格式
	var resourceResponses []gin.H
	for _, resource := range resources {
		resourceResponses = append(resourceResponses, gin.H{
			"id":          resource.ID,
			"title":       resource.Title,
			"url":         resource.URL,
			"description": resource.Description,
			"view_count":  resource.ViewCount,
			"created_at":  resource.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  resource.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	SuccessResponse(c, gin.H{
		"list":  resourceResponses,
		"total": total,
		"page":  page,
		"limit": pageSize,
	})
}

// GetHotDramas godoc
// @Summary 获取热门剧列表
// @Description 获取热门剧列表，支持分页
// @Tags PublicAPI
// @Accept json
// @Produce json
// @Param X-API-Token header string true "API访问令牌"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20) maximum(100)
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/public/hot-dramas [get]
func (h *PublicAPIHandler) GetHotDramas(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取热门剧
	hotDramas, total, err := repoManager.HotDramaRepository.FindAll(page, pageSize)
	if err != nil {
		ErrorResponse(c, "获取热门剧失败: "+err.Error(), 500)
		return
	}

	// 转换为响应格式
	var hotDramaResponses []gin.H
	for _, drama := range hotDramas {
		hotDramaResponses = append(hotDramaResponses, gin.H{
			"id":          drama.ID,
			"title":       drama.Title,
			"description": drama.CardSubtitle, // 使用副标题作为描述
			"img":         drama.PosterURL,    // 使用海报URL作为图片
			"url":         drama.DoubanURI,    // 使用豆瓣链接作为URL
			"rating":      drama.Rating,
			"year":        drama.Year,
			"region":      drama.Region,
			"genres":      drama.Genres,
			"category":    drama.Category,
			"created_at":  drama.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  drama.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	SuccessResponse(c, gin.H{
		"hot_dramas": hotDramaResponses,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}
