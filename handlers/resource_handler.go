package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	pan "github.com/ctwj/urldb/common"
	panutils "github.com/ctwj/urldb/common"
	commonutils "github.com/ctwj/urldb/common/utils"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetResources 获取资源列表
func GetResources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	utils.Info("资源列表请求 - page: %d, pageSize: %d, User-Agent: %s", page, pageSize, c.GetHeader("User-Agent"))

	params := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	if search := c.Query("search"); search != "" {
		params["search"] = search
	}
	if panID := c.Query("pan_id"); panID != "" {
		if id, err := strconv.ParseUint(panID, 10, 32); err == nil {
			params["pan_id"] = uint(id)
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		utils.Info("收到分类ID参数: %s", categoryID)
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			params["category_id"] = uint(id)
			utils.Info("解析分类ID成功: %d", uint(id))
		} else {
			utils.Error("解析分类ID失败: %v", err)
		}
	}
	if hasSaveURL := c.Query("has_save_url"); hasSaveURL != "" {
		if hasSaveURL == "true" {
			params["has_save_url"] = true
		} else if hasSaveURL == "false" {
			params["has_save_url"] = false
		}
	}
	if noSaveURL := c.Query("no_save_url"); noSaveURL != "" {
		if noSaveURL == "true" {
			params["no_save_url"] = true
		}
	}
	if panName := c.Query("pan_name"); panName != "" {
		params["pan_name"] = panName
	}

	// 获取违禁词配置（只获取一次）
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{} // 如果获取失败，使用空列表
	}

	var resources []entity.Resource
	var total int64

	// 如果有搜索关键词且启用了Meilisearch，优先使用Meilisearch搜索
	if search := c.Query("search"); search != "" && meilisearchManager != nil && meilisearchManager.IsEnabled() {
		// 构建Meilisearch过滤器
		filters := make(map[string]interface{})
		if panID := c.Query("pan_id"); panID != "" {
			if id, err := strconv.ParseUint(panID, 10, 32); err == nil {
				// 直接使用pan_id进行过滤
				filters["pan_id"] = id
			}
		}

		// 使用Meilisearch搜索
		service := meilisearchManager.GetService()
		if service != nil {
			docs, docTotal, err := service.Search(search, filters, page, pageSize)
			if err == nil {

				// 将Meilisearch文档转换为ResourceResponse（包含高亮信息）并处理违禁词
				var resourceResponses []dto.ResourceResponse
				for _, doc := range docs {
					resourceResponse := converter.ToResourceResponseFromMeilisearch(doc)

					// 处理违禁词（Meilisearch场景，需要处理高亮标记）
					if len(cleanWords) > 0 {
						forbiddenInfo := utils.CheckResourceForbiddenWords(resourceResponse.Title, resourceResponse.Description, cleanWords)
						if forbiddenInfo.HasForbiddenWords {
							resourceResponse.Title = forbiddenInfo.ProcessedTitle
							resourceResponse.Description = forbiddenInfo.ProcessedDesc
							resourceResponse.TitleHighlight = forbiddenInfo.ProcessedTitle
							resourceResponse.DescriptionHighlight = forbiddenInfo.ProcessedDesc
						}
						resourceResponse.HasForbiddenWords = forbiddenInfo.HasForbiddenWords
						resourceResponse.ForbiddenWords = forbiddenInfo.ForbiddenWords
					}

					resourceResponses = append(resourceResponses, resourceResponse)
				}

				// 返回Meilisearch搜索结果（包含高亮信息）
				SuccessResponse(c, gin.H{
					"data":      resourceResponses,
					"total":     docTotal,
					"page":      page,
					"page_size": pageSize,
					"source":    "meilisearch",
				})
				return
			} else {
				utils.Error("Meilisearch搜索失败，回退到数据库搜索: %v", err)
			}
		}
	}

	// 如果Meilisearch未启用、搜索失败或没有搜索关键词，使用数据库搜索
	if meilisearchManager == nil || !meilisearchManager.IsEnabled() || len(resources) == 0 {
		resources, total, err = repoManager.ResourceRepository.SearchWithFilters(params)
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 处理违禁词替换和标记
	var processedResources []entity.Resource
	if len(cleanWords) > 0 {
		processedResources = utils.ProcessResourcesForbiddenWords(resources, cleanWords)
		// 复制标签数据到处理后的资源
		for i := range processedResources {
			if i < len(resources) {
				processedResources[i].Tags = resources[i].Tags
			}
		}
	} else {
		processedResources = resources
	}

	// 转换为响应格式并添加违禁词标记
	var resourceResponses []gin.H
	for i, processedResource := range processedResources {
		// 使用原始资源进行检查违禁词（数据库搜索场景，使用普通处理）
		originalResource := resources[i]
		forbiddenInfo := utils.CheckResourceForbiddenWords(originalResource.Title, originalResource.Description, cleanWords)

		resourceResponse := gin.H{
			"id":          processedResource.ID,
			"key":         processedResource.Key, // 添加key字段
			"title":       forbiddenInfo.ProcessedTitle, // 使用处理后的标题
			"url":         processedResource.URL,
			"description": forbiddenInfo.ProcessedDesc, // 使用处理后的描述
			"pan_id":      processedResource.PanID,
			"view_count":  processedResource.ViewCount,
			"created_at":  processedResource.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  processedResource.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// 添加违禁词标记
		resourceResponse["has_forbidden_words"] = forbiddenInfo.HasForbiddenWords
		resourceResponse["forbidden_words"] = forbiddenInfo.ForbiddenWords

		// 添加标签信息（需要预加载）
		var tagResponses []gin.H
		if len(processedResource.Tags) > 0 {
			for _, tag := range processedResource.Tags {
				tagResponse := gin.H{
					"id":          tag.ID,
					"name":        tag.Name,
					"description": tag.Description,
				}
				tagResponses = append(tagResponses, tagResponse)
			}
		}
		resourceResponse["tags"] = tagResponses
		resourceResponse["cover"] = originalResource.Cover

		resourceResponses = append(resourceResponses, resourceResponse)
	}

	// 构建响应数据
	responseData := gin.H{
		"data":      resourceResponses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	SuccessResponse(c, responseData)
}

// GetResourceByID 根据ID获取资源
func GetResourceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	response := converter.ToResourceResponse(resource)
	SuccessResponse(c, response)
}

// GetResourcesByKey 根据Key获取资源组
func GetResourcesByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		ErrorResponse(c, "Key参数不能为空", http.StatusBadRequest)
		return
	}

	resources, err := repoManager.ResourceRepository.FindByKey(key)
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	if len(resources) == 0 {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// 转换为响应格式并处理违禁词
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{}
	}

	var responses []dto.ResourceResponse
	for _, resource := range resources {
		response := converter.ToResourceResponse(&resource)
		// 检查违禁词
		forbiddenInfo := utils.CheckResourceForbiddenWords(response.Title, response.Description, cleanWords)
		response.HasForbiddenWords = forbiddenInfo.HasForbiddenWords
		response.ForbiddenWords = forbiddenInfo.ForbiddenWords
		responses = append(responses, response)
	}

	SuccessResponse(c, gin.H{
		"resources": responses,
		"total":     len(responses),
		"key":       key,
	})
}

// CheckResourceExists 检查资源是否存在（测试FindExists函数）
func CheckResourceExists(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		ErrorResponse(c, "URL参数不能为空", http.StatusBadRequest)
		return
	}

	excludeIDStr := c.Query("exclude_id")
	var excludeID uint
	if excludeIDStr != "" {
		if id, err := strconv.ParseUint(excludeIDStr, 10, 32); err == nil {
			excludeID = uint(id)
		}
	}

	exists, err := repoManager.ResourceRepository.FindExists(url, excludeID)
	if err != nil {
		ErrorResponse(c, "检查失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"url":    url,
		"exists": exists,
	})
}

// CreateResource 创建资源
func CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource := &entity.Resource{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		PanID:       req.PanID,
		SaveURL:     req.SaveURL,
		FileSize:    req.FileSize,
		CategoryID:  req.CategoryID,
		IsValid:     req.IsValid,
		IsPublic:    req.IsPublic,
		Cover:       req.Cover,
		Author:      req.Author,
		ErrorMsg:    req.ErrorMsg,
	}

	err := repoManager.ResourceRepository.Create(resource)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 同步到Meilisearch
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		go func() {
			if err := meilisearchManager.SyncResourceToMeilisearch(resource); err != nil {
				utils.Error("同步资源到Meilisearch失败: %v", err)
			}
		}()
	}

	SuccessResponse(c, gin.H{
		"message":  "资源创建成功",
		"resource": converter.ToResourceResponse(resource),
	})
}

// UpdateResource 更新资源
func UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// 更新资源信息
	if req.Title != "" {
		resource.Title = req.Title
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.URL != "" {
		resource.URL = req.URL
	}
	if req.PanID != nil {
		resource.PanID = req.PanID
	}
	if req.SaveURL != "" {
		resource.SaveURL = req.SaveURL
	}
	if req.FileSize != "" {
		resource.FileSize = req.FileSize
	}
	if req.CategoryID != nil {
		resource.CategoryID = req.CategoryID
	}
	resource.IsValid = req.IsValid
	resource.IsPublic = req.IsPublic
	if req.Cover != "" {
		resource.Cover = req.Cover
	}
	if req.Author != "" {
		resource.Author = req.Author
	}
	if req.ErrorMsg != "" {
		resource.ErrorMsg = req.ErrorMsg
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		err = repoManager.ResourceRepository.UpdateWithTags(resource, req.TagIDs)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = repoManager.ResourceRepository.Update(resource)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 同步到Meilisearch
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		go func() {
			if err := meilisearchManager.SyncResourceToMeilisearch(resource); err != nil {
				utils.Error("同步资源到Meilisearch失败: %v", err)
			}
		}()
	}

	SuccessResponse(c, gin.H{"message": "资源更新成功"})
}

// DeleteResource 删除资源
func DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	// 使用事务确保删除操作的原子性
	err = repoManager.ResourceRepository.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. 先删除关联的访问记录（resource_views）
		if err := tx.Unscoped().Where("resource_id = ?", uint(id)).Delete(&entity.ResourceView{}).Error; err != nil {
			utils.Error("删除资源访问记录失败 (ID: %d): %v", uint(id), err)
			return err
		}

		// 2. 删除资源标签关联（resource_tags）
		if err := tx.Unscoped().Where("resource_id = ?", uint(id)).Delete(&entity.ResourceTag{}).Error; err != nil {
			utils.Error("删除资源标签关联失败 (ID: %d): %v", uint(id), err)
			return err
		}

		// 3. 最后删除资源本身
		if err := tx.Unscoped().Delete(&entity.Resource{}, uint(id)).Error; err != nil {
			utils.Error("删除资源失败 (ID: %d): %v", uint(id), err)
			return err
		}

		return nil
	})

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("成功从数据库物理删除资源及其关联数据 (ID: %d)", uint(id))

	// 如果启用了Meilisearch，尝试从Meilisearch中删除对应数据
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		go func() {
			service := meilisearchManager.GetService()
			if service != nil {
				if err := service.DeleteDocument(uint(id)); err != nil {
					utils.Error("从Meilisearch删除资源失败 (ID: %d): %v", uint(id), err)
				} else {
					utils.Info("成功从Meilisearch删除资源 (ID: %d)", uint(id))
				}
			} else {
				utils.Error("无法获取Meilisearch服务进行资源删除 (ID: %d)", uint(id))
			}
		}()
	}

	// 设置响应头，防止缓存
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	SuccessResponse(c, gin.H{"message": "资源删除成功"})
}

// SearchResources 搜索资源
func SearchResources(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var resources []entity.Resource
	var total int64
	var err error

	// 如果启用了Meilisearch，优先使用Meilisearch搜索
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		// 构建过滤器
		filters := make(map[string]interface{})
		if categoryID := c.Query("category_id"); categoryID != "" {
			if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
				filters["category"] = uint(id)
			}
		}

		// 使用Meilisearch搜索
		service := meilisearchManager.GetService()
		if service != nil {
			docs, docTotal, err := service.Search(query, filters, page, pageSize)
			if err == nil {
				// 将Meilisearch文档转换为Resource实体
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
		if query == "" {
			// 搜索关键词为空时，返回最新记录（分页）
			resources, total, err = repoManager.ResourceRepository.FindWithRelationsPaginated(page, pageSize)
		} else {
			// 有搜索关键词时，执行搜索
			resources, total, err = repoManager.ResourceRepository.Search(query, nil, page, pageSize)
		}
	}

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"resources": converter.ToResourceResponseList(resources),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 增加资源浏览次数
func IncrementResourceViewCount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}

	// 增加资源访问量
	err = repoManager.ResourceRepository.IncrementViewCount(uint(id))
	if err != nil {
		ErrorResponse(c, "增加浏览次数失败", http.StatusInternalServerError)
		return
	}

	// 记录访问记录
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	err = repoManager.ResourceViewRepository.RecordView(uint(id), ipAddress, userAgent)
	if err != nil {
		// 记录访问失败不影响主要功能，只记录日志
		utils.Error("记录资源访问失败: %v", err)
	}

	SuccessResponse(c, gin.H{"message": "浏览次数+1"})
}

// BatchDeleteResources 批量删除资源
func BatchDeleteResources(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		ErrorResponse(c, "参数错误", 400)
		return
	}

	var deletedCount int64

	// 使用事务确保批量删除操作的原子性
	err := repoManager.ResourceRepository.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. 先删除关联的访问记录（resource_views）
		if err := tx.Unscoped().Where("resource_id IN ?", req.IDs).Delete(&entity.ResourceView{}).Error; err != nil {
			utils.Error("批量删除资源访问记录失败: %v", err)
			return err
		}

		// 2. 删除资源标签关联（resource_tags）
		if err := tx.Unscoped().Where("resource_id IN ?", req.IDs).Delete(&entity.ResourceTag{}).Error; err != nil {
			utils.Error("批量删除资源标签关联失败: %v", err)
			return err
		}

		// 3. 最后删除资源本身
		result := tx.Unscoped().Delete(&entity.Resource{}, req.IDs)
		if result.Error != nil {
			utils.Error("批量删除资源失败: %v", result.Error)
			return result.Error
		}

		deletedCount = result.RowsAffected
		return nil
	})

	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("批量物理删除资源及其关联数据成功：删除 %d 个资源", deletedCount)

	// 如果启用了Meilisearch，异步删除对应的搜索数据
	if meilisearchManager != nil && meilisearchManager.IsEnabled() && len(req.IDs) > 0 {
		go func() {
			service := meilisearchManager.GetService()
			if service != nil {
				meilisearchDeletedCount := 0
				for _, id := range req.IDs {
					if err := service.DeleteDocument(id); err != nil {
						utils.Error("从Meilisearch批量删除资源失败 (ID: %d): %v", id, err)
					} else {
						meilisearchDeletedCount++
						utils.Info("成功从Meilisearch批量删除资源 (ID: %d)", id)
					}
				}
				utils.Info("Meilisearch批量删除完成：删除 %d 个资源", meilisearchDeletedCount)
			} else {
				utils.Error("批量删除时无法获取Meilisearch服务")
			}
		}()
	}

	// 设置响应头，防止缓存
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	SuccessResponse(c, gin.H{"deleted": deletedCount, "message": "批量删除成功"})
}

// GetResourceLink 获取资源链接（智能转存）
func GetResourceLink(c *gin.Context) {
	// 获取资源ID
	resourceIDStr := c.Param("id")
	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}

	utils.Info("获取资源链接请求 - resourceID: %d", resourceID)

	// 查询资源信息
	resource, err := repoManager.ResourceRepository.FindByID(uint(resourceID))
	if err != nil {
		utils.Error("查询资源失败: %v", err)
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	// 查询平台信息
	var panInfo entity.Pan
	if resource.PanID != nil {
		panPtr, err := repoManager.PanRepository.FindByID(*resource.PanID)
		if err != nil {
			utils.Error("查询平台信息失败: %v", err)
		} else if panPtr != nil {
			panInfo = *panPtr
		}
	}

	utils.Info("资源信息 - 平台: %s, 原始链接: %s, 转存链接: %s", panInfo.Name, resource.URL, resource.SaveURL)

	// 统计访问次数
	err = repoManager.ResourceRepository.IncrementViewCount(uint(resourceID))
	if err != nil {
		utils.Error("增加资源访问量失败: %v", err)
	}

	// 记录访问记录
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	err = repoManager.ResourceViewRepository.RecordView(uint(resourceID), ipAddress, userAgent)
	if err != nil {
		utils.Error("记录资源访问失败: %v", err)
	}

	// 如果不是夸克网盘，直接返回原链接
	if panInfo.Name != "quark" && panInfo.Name != "xunlei" {
		utils.Info("非夸克和迅雷资源，直接返回原链接")
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
		})
		return
	}

	// 如果已存在转存链接，直接返回
	if resource.SaveURL != "" {
		utils.Info("已存在转存链接，直接返回: %s", resource.SaveURL)
		SuccessResponse(c, gin.H{
			"url":         resource.SaveURL,
			"type":        "transferred",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
		})
		return
	}

	// 检查是否开启自动转存
	autoTransferEnabled, err := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
	if err != nil {
		utils.Error("获取自动转存配置失败: %v", err)
		// 配置获取失败，返回原链接
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "",
		})
		return
	}

	if !autoTransferEnabled {
		utils.Info("自动转存功能未开启，返回原链接")
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "",
		})
		return
	}

	// 执行自动转存
	utils.Info("开始执行自动转存")
	transferResult := performAutoTransfer(resource)

	if transferResult.Success {
		utils.Info("自动转存成功，返回转存链接: %s", transferResult.SaveURL)
		SuccessResponse(c, gin.H{
			"url":         transferResult.SaveURL,
			"type":        "transferred",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "资源易和谐，请及时用手机夸克扫码转存",
		})
	} else {
		utils.Error("自动转存失败: %s", transferResult.ErrorMsg)
		SuccessResponse(c, gin.H{
			"url":         resource.URL,
			"type":        "original",
			"platform":    panInfo.Remark,
			"resource_id": resource.ID,
			"message":     "",
		})
	}
}

// TransferResult 转存结果
type TransferResult struct {
	Success  bool   `json:"success"`
	Fid      string `json:"fid"`
	SaveURL  string `json:"save_url"`
	ErrorMsg string `json:"error_msg"`
}

// performAutoTransfer 执行自动转存
func performAutoTransfer(resource *entity.Resource) TransferResult {
	utils.Info("开始执行资源转存 - ID: %d, URL: %s", resource.ID, resource.URL)

	// 平台ID
	panID := resource.PanID

	// 获取可用的夸克账号
	accounts, err := repoManager.CksRepository.FindByPanID(*panID)
	if err != nil {
		utils.Error("获取网盘账号失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("获取网盘账号失败: %v", err),
		}
	}

	// 测试阶段，移除最小限制
	// 获取最小存储空间配置
	autoTransferMinSpace, err := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		utils.Error("获取最小存储空间配置失败: %v", err)
		autoTransferMinSpace = 5 // 默认5GB
	}

	// 过滤：只保留已激活、夸克平台、剩余空间足够的账号
	minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024
	var validAccounts []entity.Cks
	for _, acc := range accounts {
		if acc.IsValid && acc.PanID == *panID && acc.LeftSpace >= minSpaceBytes {
			validAccounts = append(validAccounts, acc)
		}
	}

	if len(validAccounts) == 0 {
		utils.Info("没有可用的网盘账号")
		return TransferResult{
			Success:  false,
			ErrorMsg: "没有可用的网盘账号",
		}
	}

	utils.Info("找到 %d 个可用网盘账号，开始转存处理...", len(validAccounts))

	// 使用第一个可用账号进行转存
	account := validAccounts[0]
	// account := accounts[0]

	// 创建网盘服务工厂
	factory := pan.NewPanFactory()

	// 执行转存
	result := transferSingleResource(resource, account, factory)

	if result.Success {
		// 更新资源的转存信息
		resource.SaveURL = result.SaveURL
		resource.Fid = result.Fid
		resource.CkID = &account.ID
		resource.ErrorMsg = ""
		if err := repoManager.ResourceRepository.Update(resource); err != nil {
			utils.Error("更新资源转存信息失败: %v", err)
		}
	} else {
		// 更新错误信息
		resource.ErrorMsg = result.ErrorMsg
		if err := repoManager.ResourceRepository.Update(resource); err != nil {
			utils.Error("更新资源错误信息失败: %v", err)
		}
	}

	return result
}

// transferSingleResource 转存单个资源
func transferSingleResource(resource *entity.Resource, account entity.Cks, factory *pan.PanFactory) TransferResult {
	utils.Info("开始转存资源 - 资源ID: %d, 账号: %s", resource.ID, account.Username)

	service, err := factory.CreatePanService(resource.URL, &pan.PanConfig{
		URL:         resource.URL,
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	if err != nil {
		utils.Error("创建网盘服务失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("创建网盘服务失败: %v", err),
		}
	}

	// 设置账号信息
	service.SetCKSRepository(repoManager.CksRepository, account)

	// 提取分享ID
	shareID, _ := commonutils.ExtractShareIdString(resource.URL)
	if shareID == "" {
		return TransferResult{
			Success:  false,
			ErrorMsg: "无效的分享链接",
		}
	}

	// 执行转存
	transferResult, err := service.Transfer(shareID) // 有些链接还需要其他信息从 url 中自行解析
	if err != nil {
		utils.Error("转存失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("转存失败: %v", err),
		}
	}

	if transferResult == nil || !transferResult.Success {
		errMsg := "转存失败"
		if transferResult != nil && transferResult.Message != "" {
			errMsg = transferResult.Message
		}
		utils.Error("转存失败: %s", errMsg)
		return TransferResult{
			Success:  false,
			ErrorMsg: errMsg,
		}
	}

	// 提取转存链接
	var saveURL string
	var fid string

	if data, ok := transferResult.Data.(map[string]interface{}); ok {
		if v, ok := data["shareUrl"]; ok {
			saveURL, _ = v.(string)
		}
		if v, ok := data["fid"]; ok {
			fid, _ = v.(string)
		}
	}
	if saveURL == "" {
		saveURL = transferResult.ShareURL
	}

	if saveURL == "" {
		return TransferResult{
			Success:  false,
			ErrorMsg: "转存成功但未获取到分享链接",
		}
	}

	utils.Info("转存成功 - 资源ID: %d, 转存链接: %s", resource.ID, saveURL)

	return TransferResult{
		Success: true,
		SaveURL: saveURL,
		Fid:     fid,
	}
}

// GetHotResources 获取热门资源
func GetHotResources(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	utils.Info("获取热门资源请求 - limit: %d", limit)

	// 限制最大请求数量
	if limit > 20 {
		limit = 20
	}
	if limit <= 0 {
		limit = 10
	}

	// 使用公共缓存机制
	cacheKey := fmt.Sprintf("hot_resources_%d", limit)
	ttl := time.Hour // 1小时缓存
	cacheManager := utils.GetHotResourcesCache()

	// 尝试从缓存获取
	if cachedData, found := cacheManager.Get(cacheKey, ttl); found {
		utils.Info("使用热门资源缓存 - key: %s", cacheKey)
		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("ETag", fmt.Sprintf("hot-resources-%d", len(cachedData.([]gin.H))))

		// 转换为正确的类型
		if data, ok := cachedData.([]gin.H); ok {
			SuccessResponse(c, gin.H{
				"data":    data,
				"total":   len(data),
				"limit":   limit,
				"cached":  true,
			})
		}
		return
	}

	// 缓存未命中，从数据库获取
	resources, err := repoManager.ResourceRepository.GetHotResources(limit)
	if err != nil {
		utils.Error("获取热门资源失败: %v", err)
		ErrorResponse(c, "获取热门资源失败", http.StatusInternalServerError)
		return
	}

	// 获取违禁词配置
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{}
	}

	// 处理违禁词并转换为响应格式
	var resourceResponses []gin.H
	for _, resource := range resources {
		// 检查违禁词
		forbiddenInfo := utils.CheckResourceForbiddenWords(resource.Title, resource.Description, cleanWords)

		resourceResponse := gin.H{
			"id":          resource.ID,
			"key":         resource.Key,
			"title":       forbiddenInfo.ProcessedTitle,
			"url":         resource.URL,
			"description": forbiddenInfo.ProcessedDesc,
			"pan_id":      resource.PanID,
			"view_count":  resource.ViewCount,
			"created_at":  resource.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  resource.UpdatedAt.Format("2006-01-02 15:04:05"),
			"cover":       resource.Cover,
			"author":      resource.Author,
			"file_size":   resource.FileSize,
		}

		// 添加违禁词标记
		resourceResponse["has_forbidden_words"] = forbiddenInfo.HasForbiddenWords
		resourceResponse["forbidden_words"] = forbiddenInfo.ForbiddenWords

		// 添加标签信息
		var tagResponses []gin.H
		if len(resource.Tags) > 0 {
			for _, tag := range resource.Tags {
				tagResponse := gin.H{
					"id":          tag.ID,
					"name":        tag.Name,
					"description": tag.Description,
				}
				tagResponses = append(tagResponses, tagResponse)
			}
		}
		resourceResponse["tags"] = tagResponses

		resourceResponses = append(resourceResponses, resourceResponse)
	}

	// 存储到缓存
	cacheManager.Set(cacheKey, resourceResponses)
	utils.Info("热门资源已缓存 - key: %s, count: %d", cacheKey, len(resourceResponses))

	// 设置缓存头
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("ETag", fmt.Sprintf("hot-resources-%d", len(resourceResponses)))

	SuccessResponse(c, gin.H{
		"data":    resourceResponses,
		"total":   len(resourceResponses),
		"limit":   limit,
		"cached":  false,
	})
}

// GetRelatedResources 获取相关资源
func GetRelatedResources(c *gin.Context) {
	// 获取查询参数
	key := c.Query("key")                   // 当前资源的key
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "8"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	utils.Info("获取相关资源请求 - key: %s, limit: %d", key, limit)

	if key == "" {
		ErrorResponse(c, "缺少资源key参数", http.StatusBadRequest)
		return
	}

	// 首先通过key获取当前资源信息
	currentResources, err := repoManager.ResourceRepository.FindByKey(key)
	if err != nil {
		utils.Error("获取当前资源失败: %v", err)
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	if len(currentResources) == 0 {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	currentResource := &currentResources[0] // 取第一个资源作为当前资源

	var resources []entity.Resource
	var total int64

	// 获取当前资源的标签ID列表
	var tagIDsList []string
	if currentResource.Tags != nil {
		for _, tag := range currentResource.Tags {
			tagIDsList = append(tagIDsList, strconv.Itoa(int(tag.ID)))
		}
	}

	utils.Info("当前资源标签: %v", tagIDsList)

	// 1. 优先使用Meilisearch进行标签搜索
	if meilisearchManager != nil && meilisearchManager.IsEnabled() && len(tagIDsList) > 0 {
		service := meilisearchManager.GetService()
		if service != nil {
			// 使用标签进行搜索
			filters := make(map[string]interface{})
			filters["tag_ids"] = tagIDsList

			// 使用当前资源的标题作为搜索关键词，提高相关性
			searchQuery := currentResource.Title
			if searchQuery == "" {
				searchQuery = strings.Join(tagIDsList, " ") // 如果没有标题，使用标签作为搜索词
			}

			docs, docTotal, err := service.Search(searchQuery, filters, page, limit)
			if err == nil && len(docs) > 0 {
				// 转换为Resource实体
				for _, doc := range docs {
					// 排除当前资源
					if doc.Key == key {
						continue
					}
					resource := entity.Resource{
						ID:          doc.ID,
						Title:       doc.Title,
						Description: doc.Description,
						URL:         doc.URL,
						SaveURL:     doc.SaveURL,
						FileSize:    doc.FileSize,
						Key:         doc.Key,
						PanID:       doc.PanID,
						ViewCount:   0, // Meilisearch文档中没有ViewCount字段，设为默认值
						CreatedAt:   doc.CreatedAt,
						UpdatedAt:   doc.UpdatedAt,
						Cover:       doc.Cover,
						Author:      doc.Author,
					}
					resources = append(resources, resource)
				}
				total = docTotal
				utils.Info("Meilisearch搜索到 %d 个相关资源", len(resources))
			} else {
				utils.Error("Meilisearch搜索失败，回退到标签搜索: %v", err)
			}
		}
	}

	// 2. 如果Meilisearch未启用、搜索失败或没有结果，使用数据库标签搜索
	if len(resources) == 0 {
		params := map[string]interface{}{
			"page":      page,
			"page_size": limit,
			"is_public": true,
			"order_by":  "updated_at",
			"order_dir": "desc",
		}

		// 使用当前资源的标签进行搜索
		if len(tagIDsList) > 0 {
			params["tag_ids"] = strings.Join(tagIDsList, ",")
		} else {
			// 如果没有标签，使用当前资源的分类作为搜索条件
			if currentResource.CategoryID != nil && *currentResource.CategoryID > 0 {
				params["category_id"] = *currentResource.CategoryID
			}
		}

		var err error
		resources, total, err = repoManager.ResourceRepository.SearchWithFilters(params)
		if err != nil {
			utils.Error("搜索相关资源失败: %v", err)
			ErrorResponse(c, "搜索相关资源失败", http.StatusInternalServerError)
			return
		}

		// 排除当前资源
		var filteredResources []entity.Resource
		for _, resource := range resources {
			if resource.Key != key {
				filteredResources = append(filteredResources, resource)
			}
		}
		resources = filteredResources
		total = int64(len(filteredResources))
	}

	utils.Info("标签搜索到 %d 个相关资源", len(resources))

	// 获取违禁词配置
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{}
	}

	// 处理违禁词并转换为响应格式
	var resourceResponses []gin.H
	for _, resource := range resources {
		// 检查违禁词
		forbiddenInfo := utils.CheckResourceForbiddenWords(resource.Title, resource.Description, cleanWords)

		resourceResponse := gin.H{
			"id":          resource.ID,
			"key":         resource.Key,
			"title":       forbiddenInfo.ProcessedTitle,
			"url":         resource.URL,
			"description": forbiddenInfo.ProcessedDesc,
			"pan_id":      resource.PanID,
			"view_count":  resource.ViewCount,
			"created_at":  resource.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  resource.UpdatedAt.Format("2006-01-02 15:04:05"),
			"cover":       resource.Cover,
			"author":      resource.Author,
			"file_size":   resource.FileSize,
		}

		// 添加违禁词标记
		resourceResponse["has_forbidden_words"] = forbiddenInfo.HasForbiddenWords
		resourceResponse["forbidden_words"] = forbiddenInfo.ForbiddenWords

		// 添加标签信息
		var tagResponses []gin.H
		if len(resource.Tags) > 0 {
			for _, tag := range resource.Tags {
				tagResponse := gin.H{
					"id":          tag.ID,
					"name":        tag.Name,
					"description": tag.Description,
				}
				tagResponses = append(tagResponses, tagResponse)
			}
		}
		resourceResponse["tags"] = tagResponses

		resourceResponses = append(resourceResponses, resourceResponse)
	}

	// 构建响应数据
	responseData := gin.H{
		"data":      resourceResponses,
		"total":     total,
		"page":      page,
		"page_size": limit,
		"source":    "database",
	}

	if meilisearchManager != nil && meilisearchManager.IsEnabled() && len(tagIDsList) > 0 {
		responseData["source"] = "meilisearch"
	}

	SuccessResponse(c, responseData)
}

// CheckResourceValidity 检查资源链接有效性
func CheckResourceValidity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的资源ID", http.StatusBadRequest)
		return
	}

	// 查询资源信息
	resource, err := repoManager.ResourceRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "资源不存在", http.StatusNotFound)
		return
	}

	utils.Info("开始检测资源有效性 - ID: %d, URL: %s", resource.ID, resource.URL)

	// 检查缓存
	cacheKey := fmt.Sprintf("resource_validity_%d", resource.ID)
	cacheManager := utils.GetResourceValidityCache()
	ttl := 5 * time.Minute // 5分钟缓存

	if cachedData, found := cacheManager.Get(cacheKey, ttl); found {
		if result, ok := cachedData.(gin.H); ok {
			utils.Info("使用资源有效性缓存 - ID: %d", resource.ID)
			result["cached"] = true
			SuccessResponse(c, result)
			return
		}
	}

	// 执行检测：只使用深度检测实现
	isValid, detectionMethod, err := performAdvancedValidityCheck(resource)

	if err != nil {
		utils.Error("深度检测资源链接失败 - ID: %d, Error: %v", resource.ID, err)

		// 深度检测失败，但不标记为无效（用户可自行验证）
		result := gin.H{
			"resource_id":      resource.ID,
			"url":             resource.URL,
			"is_valid":        resource.IsValid, // 保持原始状态
			"last_checked":    time.Now().Format(time.RFC3339),
			"error":           err.Error(),
			"detection_method": detectionMethod,
			"cached":          false,
			"note":            "当前网盘暂不支持自动检测，建议用户自行验证",
		}
		cacheManager.Set(cacheKey, result)
		SuccessResponse(c, result)
		return
	}

	// 只有明确检测出无效的资源才更新数据库状态
	// 如果检测成功且结果与数据库状态不同，则更新
	if (detectionMethod == "quark_deep" || detectionMethod == "baidu_deep") && isValid != resource.IsValid {
		resource.IsValid = isValid
		updateErr := repoManager.ResourceRepository.Update(resource)
		if updateErr != nil {
			utils.Error("更新资源有效性状态失败 - ID: %d, Error: %v", resource.ID, updateErr)
		} else {
			utils.Info("更新资源有效性状态 - ID: %d, Status: %v, Method: %s", resource.ID, isValid, detectionMethod)
		}
	}

	// 构建检测结果
	result := gin.H{
		"resource_id":      resource.ID,
		"url":             resource.URL,
		"is_valid":        isValid,
		"last_checked":    time.Now().Format(time.RFC3339),
		"detection_method": detectionMethod,
		"cached":          false,
	}

	// 缓存检测结果
	cacheManager.Set(cacheKey, result)

	utils.Info("资源有效性检测完成 - ID: %d, Valid: %v, Method: %s", resource.ID, isValid, detectionMethod)
	SuccessResponse(c, result)
}

// performAdvancedValidityCheck 执行深度检测（只使用具体网盘服务）
func performAdvancedValidityCheck(resource *entity.Resource) (bool, string, error) {
	// 提取分享ID和服务类型
	shareID, serviceType := panutils.ExtractShareId(resource.URL)
	if serviceType == panutils.NotFound {
		return false, "unsupported", fmt.Errorf("不支持的网盘服务: %s", resource.URL)
	}

	utils.Info("开始深度检测 - Service: %s, ShareID: %s", serviceType.String(), shareID)

	// 根据服务类型选择检测策略
	switch serviceType {
	case panutils.Quark:
		return performQuarkValidityCheck(resource, shareID)
	case panutils.BaiduPan:
		return performBaiduValidityCheck(resource, shareID)
	case panutils.Alipan:
		return performAlipanValidityCheck(resource, shareID)
	case panutils.UC, panutils.Xunlei, panutils.Tianyi, panutils.Pan123, panutils.Pan115:
		// 这些网盘暂未实现深度检测，返回不支持提示
		return false, "unsupported", fmt.Errorf("当前网盘类型 %s 暂不支持深度检测，请等待后续更新", serviceType.String())
	default:
		return false, "unsupported", fmt.Errorf("未知的网盘服务类型: %s", serviceType.String())
	}
}

// performQuarkValidityCheck 夸克网盘深度检测
func performQuarkValidityCheck(resource *entity.Resource, shareID string) (bool, string, error) {
	// 获取夸克网盘账号
	panID, err := getQuarkPanID()
	if err != nil {
		return false, "quark_failed", fmt.Errorf("获取夸克平台ID失败: %v", err)
	}

	accounts, err := repoManager.CksRepository.FindByPanID(panID)
	if err != nil {
		return false, "quark_failed", fmt.Errorf("获取夸克网盘账号失败: %v", err)
	}

	if len(accounts) == 0 {
		return false, "quark_failed", fmt.Errorf("没有可用的夸克网盘账号")
	}

	// 选择第一个有效账号
	var selectedAccount *entity.Cks
	for _, account := range accounts {
		if account.IsValid {
			selectedAccount = &account
			break
		}
	}

	if selectedAccount == nil {
		return false, "quark_failed", fmt.Errorf("没有有效的夸克网盘账号")
	}

	// 创建网盘服务配置
	config := &pan.PanConfig{
		URL:         resource.URL,
		Code:        "",
		IsType:      1, // 只获取基本信息，不转存
		ExpiredType: 1,
		AdFid:       "",
		Stoken:      "",
		Cookie:      selectedAccount.Ck,
	}

	// 创建夸克网盘服务
	factory := pan.NewPanFactory()
	panService, err := factory.CreatePanService(resource.URL, config)
	if err != nil {
		return false, "quark_failed", fmt.Errorf("创建夸克网盘服务失败: %v", err)
	}

	// 执行深度检测（Transfer方法）
	utils.Info("执行夸克网盘深度检测 - ShareID: %s", shareID)
	result, err := panService.Transfer(shareID)
	if err != nil {
		return false, "quark_failed", fmt.Errorf("夸克网盘检测失败: %v", err)
	}

	if !result.Success {
		return false, "quark_failed", fmt.Errorf("夸克网盘链接无效: %s", result.Message)
	}

	utils.Info("夸克网盘深度检测成功 - ShareID: %s", shareID)
	return true, "quark_deep", nil
}

// performAlipanValidityCheck 阿里云盘深度检测
func performAlipanValidityCheck(resource *entity.Resource, shareID string) (bool, string, error) {
	// 阿里云盘深度检测暂未实现
	utils.Info("阿里云盘暂不支持深度检测 - ShareID: %s", shareID)
	return false, "unsupported", fmt.Errorf("阿里云盘暂不支持深度检测，请等待后续更新")
}

// performBaiduValidityCheck 百度网盘深度检测
func performBaiduValidityCheck(resource *entity.Resource, shareID string) (bool, string, error) {
	// 获取百度网盘账号
	panID, err := getBaiduPanID()
	if err != nil {
		return false, "baidu_failed", fmt.Errorf("获取百度平台ID失败: %v", err)
	}

	accounts, err := repoManager.CksRepository.FindByPanID(panID)
	if err != nil {
		return false, "baidu_failed", fmt.Errorf("获取百度网盘账号失败: %v", err)
	}

	if len(accounts) == 0 {
		return false, "baidu_failed", fmt.Errorf("没有可用的百度网盘账号")
	}

	// 选择第一个有效的账号
	var selectedAccount *entity.Cks
	for _, account := range accounts {
		if account.IsValid {
			selectedAccount = &account
			break
		}
	}

	if selectedAccount == nil {
		return false, "baidu_failed", fmt.Errorf("没有有效的百度网盘账号")
	}

	utils.Info("使用百度网盘账号进行深度检测 - AccountID: %d", selectedAccount.ID)

	// 准备配置
	config := &panutils.PanConfig{
		URL:         resource.URL,
		Code:        "", // 可以从URL中提取
		IsType:      1,  // 检验模式：只获取基本信息，不实际转存
		ExpiredType: 1,  // 永久分享
		AdFid:       "",
		Stoken:      "",
		Cookie:      selectedAccount.Ck,
	}

	// 创建网盘服务工厂
	factory := pan.NewPanFactory()
	// 通过工厂获取对应的网盘服务
	panService, err := factory.CreatePanService(resource.URL, config)
	if err != nil {
		return false, "baidu_failed", fmt.Errorf("创建百度网盘服务失败: %v", err)
	}

	// 执行深度检测（Transfer方法，检验模式）
	utils.Info("执行百度网盘深度检测 - ShareID: %s", shareID)
	result, err := panService.Transfer(shareID)
	if err != nil {
		return false, "baidu_failed", fmt.Errorf("百度网盘检测失败: %v", err)
	}

	if !result.Success {
		// 检查是否是因为没有Cookie或其他配置问题
		if strings.Contains(result.Message, "cookie") || strings.Contains(result.Message, "登录") {
			return false, "baidu_failed", fmt.Errorf("百度网盘账号配置问题: %s", result.Message)
		}
		return false, "baidu_deep", fmt.Errorf("百度网盘链接无效: %s", result.Message)
	}

	utils.Info("百度网盘深度检测成功 - ShareID: %s", shareID)
	return true, "baidu_deep", nil
}


// BatchCheckResourceValidity 批量检查资源链接有效性
func BatchCheckResourceValidity(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		ErrorResponse(c, "ID列表不能为空", http.StatusBadRequest)
		return
	}

	if len(req.IDs) > 20 {
		ErrorResponse(c, "单次最多检测20个资源", http.StatusBadRequest)
		return
	}

	utils.Info("开始批量检测资源有效性 - Count: %d", len(req.IDs))

	cacheManager := utils.GetResourceValidityCache()
	ttl := 5 * time.Minute
	results := make([]gin.H, 0, len(req.IDs))

	for _, id := range req.IDs {
		// 查询资源信息
		resource, err := repoManager.ResourceRepository.FindByID(id)
		if err != nil {
			results = append(results, gin.H{
				"resource_id": id,
				"is_valid":    false,
				"error":       "资源不存在",
				"cached":      false,
			})
			continue
		}

		// 检查缓存
		cacheKey := fmt.Sprintf("resource_validity_%d", id)
		if cachedData, found := cacheManager.Get(cacheKey, ttl); found {
			if result, ok := cachedData.(gin.H); ok {
				result["cached"] = true
				results = append(results, result)
				continue
			}
		}

		// 执行深度检测
		isValid, detectionMethod, err := performAdvancedValidityCheck(resource)

		if err != nil {
			// 深度检测失败，但不标记为无效（用户可自行验证）
			result := gin.H{
				"resource_id":      id,
				"url":             resource.URL,
				"is_valid":        resource.IsValid, // 保持原始状态
				"last_checked":    time.Now().Format(time.RFC3339),
				"error":           err.Error(),
				"detection_method": detectionMethod,
				"cached":          false,
				"note":            "当前网盘暂不支持自动检测，建议用户自行验证",
			}
			cacheManager.Set(cacheKey, result)
			results = append(results, result)
			continue
		}

		// 只有明确检测出无效的资源才更新数据库状态
		if (detectionMethod == "quark_deep" || detectionMethod == "baidu_deep") && isValid != resource.IsValid {
			resource.IsValid = isValid
			updateErr := repoManager.ResourceRepository.Update(resource)
			if updateErr != nil {
				utils.Error("更新资源有效性状态失败 - ID: %d, Error: %v", id, updateErr)
			}
		}

		result := gin.H{
			"resource_id":      id,
			"url":             resource.URL,
			"is_valid":        isValid,
			"last_checked":    time.Now().Format(time.RFC3339),
			"detection_method": detectionMethod,
			"cached":          false,
		}

		cacheManager.Set(cacheKey, result)
		results = append(results, result)
	}

	utils.Info("批量检测资源有效性完成 - Count: %d", len(results))
	SuccessResponse(c, gin.H{
		"results": results,
		"total":   len(results),
	})
}

// getQuarkPanID 获取夸克网盘ID
func getQuarkPanID() (uint, error) {
	// 通过FindAll方法查找所有平台，然后过滤出quark平台
	pans, err := repoManager.PanRepository.FindAll()
	if err != nil {
		return 0, fmt.Errorf("查询平台信息失败: %v", err)
	}

	for _, p := range pans {
		if p.Name == "quark" {
			return p.ID, nil
		}
	}

	return 0, fmt.Errorf("未找到quark平台")
}

// getBaiduPanID 获取百度网盘ID
func getBaiduPanID() (uint, error) {
	// 通过FindAll方法查找所有平台，然后过滤出baidu平台
	pans, err := repoManager.PanRepository.FindAll()
	if err != nil {
		return 0, fmt.Errorf("查询平台信息失败: %v", err)
	}

	for _, p := range pans {
		if p.Name == "baidu" {
			return p.ID, nil
		}
	}

	return 0, fmt.Errorf("未找到baidu平台")
}

