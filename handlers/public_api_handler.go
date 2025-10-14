package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"

	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
)

// PublicAPIHandler 公开API处理器
type PublicAPIHandler struct{}

// NewPublicAPIHandler 创建公开API处理器
func NewPublicAPIHandler() *PublicAPIHandler {
	return &PublicAPIHandler{}
}

// filterForbiddenWords 过滤包含违禁词的资源
func (h *PublicAPIHandler) filterForbiddenWords(resources []entity.Resource) ([]entity.Resource, []string) {
	// 获取违禁词配置
	forbiddenWords, err := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyForbiddenWords)
	if err != nil {
		// 如果获取失败，返回原资源列表
		return resources, nil
	}

	if forbiddenWords == "" {
		return resources, nil
	}

	// 分割违禁词
	words := strings.Split(forbiddenWords, ",")
	var filteredResources []entity.Resource
	var foundForbiddenWords []string

	for _, resource := range resources {
		shouldSkip := false
		title := strings.ToLower(resource.Title)
		description := strings.ToLower(resource.Description)

		for _, word := range words {
			word = strings.TrimSpace(word)
			if word != "" && (strings.Contains(title, strings.ToLower(word)) || strings.Contains(description, strings.ToLower(word))) {
				foundForbiddenWords = append(foundForbiddenWords, word)
				shouldSkip = true
				break
			}
		}

		if !shouldSkip {
			filteredResources = append(filteredResources, resource)
		}
	}

	// 去重违禁词
	uniqueForbiddenWords := make([]string, 0)
	wordMap := make(map[string]bool)
	for _, word := range foundForbiddenWords {
		if !wordMap[word] {
			wordMap[word] = true
			uniqueForbiddenWords = append(uniqueForbiddenWords, word)
		}
	}

	return filteredResources, uniqueForbiddenWords
}

// logAPIAccess 记录API访问日志
func (h *PublicAPIHandler) logAPIAccess(c *gin.Context, startTime time.Time, processCount int, responseData interface{}, errorMessage string) {
	endpoint := c.Request.URL.Path
	method := c.Request.Method
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 计算处理时间
	processingTime := time.Since(startTime).Milliseconds()

	// 获取查询参数
	var requestParams interface{}
	if method == "GET" {
		requestParams = c.Request.URL.Query()
	} else {
		// 对于POST请求，尝试从上下文中获取请求体（如果之前已解析）
		if req, exists := c.Get("request_body"); exists {
			requestParams = req
		}
	}

	// 异步记录日志，避免影响API响应时间
	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.Error("记录API访问日志时发生panic: %v", r)
			}
		}()

		err := repoManager.APIAccessLogRepository.RecordAccess(
			ip,
			userAgent,
			endpoint,
			method,
			requestParams,
			c.Writer.Status(),
			responseData,
			processCount,
			errorMessage,
			processingTime,
		)
		if err != nil {
			utils.Error("记录API访问日志失败: %v", err)
		}
	}()
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
	startTime := time.Now()

	var req dto.BatchReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logAPIAccess(c, startTime, 0, nil, "请求参数错误: "+err.Error())
		ErrorResponse(c, "请求参数错误: "+err.Error(), 400)
		return
	}

	// 存储请求体用于日志记录
	c.Set("request_body", req)

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
			h.logAPIAccess(c, startTime, len(createdResources), nil, "生成资源组标识失败: "+err.Error())
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

	responseData := gin.H{
		"created_count": len(createdResources),
		"created_ids":   createdResources,
	}
	h.logAPIAccess(c, startTime, len(createdResources), responseData, "")
	SuccessResponse(c, responseData)
}

// SearchResources godoc
// @Summary 资源搜索
// @Description 搜索资源，支持关键词、标签、分类过滤，自动过滤包含违禁词的资源
// @Tags PublicAPI
// @Accept json
// @Produce json
// @Param X-API-Token header string true "API访问令牌"
// @Param keyword query string false "搜索关键词"
// @Param tag query string false "标签过滤"
// @Param category query string false "分类过滤"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20) maximum(100)
// @Success 200 {object} map[string]interface{} "搜索成功，如果存在违禁词过滤会返回forbidden_words_filtered字段"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/public/resources/search [get]
func (h *PublicAPIHandler) SearchResources(c *gin.Context) {
	startTime := time.Now()

	// 获取查询参数
	keyword := c.Query("keyword")
	tag := c.Query("tag")
	category := c.Query("category")
	panID := c.Query("pan_id")
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

	var resources []entity.Resource
	var total int64

	// 如果启用了Meilisearch，优先使用Meilisearch搜索
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		// 构建过滤器
		filters := make(map[string]interface{})
		if category != "" {
			filters["category"] = category
		}
		if tag != "" {
			filters["tags"] = tag
		}
		if panID != "" {
			if id, err := strconv.ParseUint(panID, 10, 32); err == nil {
				// 根据pan_id获取pan_name
				pan, err := repoManager.PanRepository.FindByID(uint(id))
				if err == nil && pan != nil {
					filters["pan_name"] = pan.Name
				}
			}
		}

		// 使用Meilisearch搜索
		service := meilisearchManager.GetService()
		if service != nil {
			docs, docTotal, err := service.Search(keyword, filters, page, pageSize)
			if err == nil {
				// 将Meilisearch文档转换为Resource实体（保持兼容性）
				for _, doc := range docs {
					resource := entity.Resource{
						ID:          doc.ID,
						Title:       doc.Title,
						Description: doc.Description,
						URL:         doc.URL,
						SaveURL:     doc.SaveURL,
						FileSize:    doc.FileSize,
						Key:         doc.Key,
						PanID:       doc.PanID,
						Cover:       doc.Cover,
						CreatedAt:   doc.CreatedAt,
						UpdatedAt:   doc.UpdatedAt,
					}
					resources = append(resources, resource)
				}
				total = docTotal
			} else {
				utils.Error("Meilisearch搜索失败，回退到数据库搜索: %v", err)
			}
		}
	}

	// 如果Meilisearch未启用或搜索失败，使用数据库搜索
	if meilisearchManager == nil || !meilisearchManager.IsEnabled() || err != nil {
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
		if panID != "" {
			if id, err := strconv.ParseUint(panID, 10, 32); err == nil {
				params["pan_id"] = uint(id)
			}
		}

		// 执行数据库搜索
		resources, total, err = repoManager.ResourceRepository.SearchWithFilters(params)
		if err != nil {
			h.logAPIAccess(c, startTime, 0, nil, "搜索失败: "+err.Error())
			ErrorResponse(c, "搜索失败: "+err.Error(), 500)
			return
		}
	}

	// 获取违禁词配置（只获取一次）
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{} // 如果获取失败，使用空列表
	}

	// 转换为响应格式并添加违禁词标记
	var resourceResponses []gin.H
	for i, processedResource := range resources {
		originalResource := resources[i]
		forbiddenInfo := utils.CheckResourceForbiddenWords(originalResource.Title, originalResource.Description, cleanWords)

		resourceResponse := gin.H{
			"id":          processedResource.ID,
			"title":       forbiddenInfo.ProcessedTitle, // 使用处理后的标题
			"url":         processedResource.URL,
			"description": forbiddenInfo.ProcessedDesc, // 使用处理后的描述
			"view_count":  processedResource.ViewCount,
			"created_at":  processedResource.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  processedResource.UpdatedAt.Format("2006-01-02 15:04:05"),
			"cover":       processedResource.Cover, // 添加封面字段
		}

		// 添加违禁词标记
		resourceResponse["has_forbidden_words"] = forbiddenInfo.HasForbiddenWords
		resourceResponse["forbidden_words"] = forbiddenInfo.ForbiddenWords
		resourceResponses = append(resourceResponses, resourceResponse)
	}

	// 构建响应数据
	responseData := gin.H{
		"list":  resourceResponses,
		"total": total,
		"page":  page,
		"limit": pageSize,
	}

	h.logAPIAccess(c, startTime, len(resourceResponses), responseData, "")
	SuccessResponse(c, responseData)
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
	startTime := time.Now()

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
		h.logAPIAccess(c, startTime, 0, nil, "获取热门剧失败: "+err.Error())
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

	responseData := gin.H{
		"hot_dramas": hotDramaResponses,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}
	h.logAPIAccess(c, startTime, len(hotDramaResponses), responseData, "")
	SuccessResponse(c, responseData)
}
