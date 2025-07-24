package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"

	"github.com/gin-gonic/gin"
)

// GetReadyResources 获取待处理资源列表
func GetReadyResources(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 1000 {
		pageSize = 100
	}

	// 获取分页数据
	resources, total, err := repoManager.ReadyResourceRepository.FindWithPagination(page, pageSize)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToReadyResourceResponseList(resources)

	// 使用标准化的分页响应格式
	SuccessResponse(c, gin.H{
		"data":      responses,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// BatchCreateReadyResources 批量创建待处理资源
func BatchCreateReadyResources(c *gin.Context) {
	var req dto.BatchCreateReadyResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. 先收集所有待提交的URL，去重
	urlSet := make(map[string]struct{})
	for _, reqResource := range req.Resources {
		if len(reqResource.URL) == 0 {
			continue
		}
		for _, u := range reqResource.URL {
			if u != "" {
				urlSet[u] = struct{}{}
			}
		}
	}
	uniqueUrls := make([]string, 0, len(urlSet))
	for url := range urlSet {
		uniqueUrls = append(uniqueUrls, url)
	}

	// 2. 批量查询待处理资源表中已存在的URL
	existReadyUrls := make(map[string]struct{})
	if len(uniqueUrls) > 0 {
		readyList, _ := repoManager.ReadyResourceRepository.BatchFindByURLs(uniqueUrls)
		for _, r := range readyList {
			existReadyUrls[r.URL] = struct{}{}
		}
	}

	// 3. 批量查询资源表中已存在的URL
	existResourceUrls := make(map[string]struct{})
	if len(uniqueUrls) > 0 {
		resourceList, _ := repoManager.ResourceRepository.BatchFindByURLs(uniqueUrls)
		for _, r := range resourceList {
			existResourceUrls[r.URL] = struct{}{}
		}
	}

	// 5. 过滤掉已存在的URL
	var resources []entity.ReadyResource
	for _, reqResource := range req.Resources {
		if len(reqResource.URL) == 0 {
			continue
		}
		key, err := repoManager.ReadyResourceRepository.GenerateUniqueKey()
		if err != nil {
			ErrorResponse(c, "生成批量资源组标识失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for _, url := range reqResource.URL {
			if url == "" {
				continue
			}
			if _, ok := existReadyUrls[url]; ok {
				continue
			}
			if _, ok := existResourceUrls[url]; ok {
				continue
			}

			resource := entity.ReadyResource{
				Title:       reqResource.Title,
				Description: reqResource.Description,
				URL:         url,
				Category:    reqResource.Category,
				Tags:        reqResource.Tags,
				Img:         reqResource.Img,
				Source:      reqResource.Source,
				Extra:       reqResource.Extra,
				IP:          reqResource.IP,
				Key:         key,
			}
			resources = append(resources, resource)
		}
	}

	if len(resources) == 0 {
		SuccessResponse(c, gin.H{
			"count":   0,
			"message": "无新增资源，所有URL均已存在",
		})
		return
	}

	err := repoManager.ReadyResourceRepository.BatchCreate(resources)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"count":   len(resources),
		"message": "批量创建成功",
	})
}

// CreateReadyResourcesFromText 从文本创建待处理资源
func CreateReadyResourcesFromText(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		ErrorResponse(c, "文本内容不能为空", http.StatusBadRequest)
		return
	}

	lines := strings.Split(text, "\n")
	var resources []entity.ReadyResource

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 简单的URL提取逻辑
		if strings.Contains(line, "http") {
			resource := entity.ReadyResource{
				URL: line,
			}
			resources = append(resources, resource)
		}
	}

	if len(resources) == 0 {
		ErrorResponse(c, "未找到有效的URL", http.StatusBadRequest)
		return
	}

	err := repoManager.ReadyResourceRepository.BatchCreate(resources)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"count":   len(resources),
		"message": "从文本创建成功",
	})
}

// DeleteReadyResource 删除待处理资源
func DeleteReadyResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.ReadyResourceRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "待处理资源删除成功"})
}

// ClearReadyResources 清空所有待处理资源
func ClearReadyResources(c *gin.Context) {
	resources, err := repoManager.ReadyResourceRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, resource := range resources {
		err = repoManager.ReadyResourceRepository.Delete(resource.ID)
		if err != nil {
			ErrorResponse(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{
		"deleted_count": len(resources),
		"message":       "所有待处理资源已清空",
	})
}

// GetReadyResourcesByKey 根据key获取待处理资源
func GetReadyResourcesByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		ErrorResponse(c, "key参数不能为空", http.StatusBadRequest)
		return
	}

	resources, err := repoManager.ReadyResourceRepository.FindByKey(key)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToReadyResourceResponseList(resources)

	SuccessResponse(c, gin.H{
		"data":  responses,
		"key":   key,
		"count": len(resources),
	})
}

// DeleteReadyResourcesByKey 根据key删除待处理资源
func DeleteReadyResourcesByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		ErrorResponse(c, "key参数不能为空", http.StatusBadRequest)
		return
	}

	// 先查询要删除的资源数量
	resources, err := repoManager.ReadyResourceRepository.FindByKey(key)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(resources) == 0 {
		ErrorResponse(c, "未找到指定key的资源", http.StatusNotFound)
		return
	}

	// 删除所有具有相同key的资源
	err = repoManager.ReadyResourceRepository.DeleteByKey(key)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"deleted_count": len(resources),
		"key":           key,
		"message":       "资源组删除成功",
	})
}
