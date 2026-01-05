package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/plugin-system/triggers/plugins"

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

	// 4. 在过滤之前触发插件系统待处理资源添加事件（包括被过滤的URL）
	for _, reqResource := range req.Resources {
		if len(reqResource.URL) == 0 {
			continue
		}
		for _, url := range reqResource.URL {
			if url == "" {
				continue
			}

			// 检查URL是否已存在，用于标识过滤状态
			isFiltered := false
			filterReason := ""
			if _, ok := existReadyUrls[url]; ok {
				isFiltered = true
				filterReason = "exists_in_ready_table"
			} else if _, ok := existResourceUrls[url]; ok {
				isFiltered = true
				filterReason = "exists_in_resource_table"
			}

			// 创建临时资源对象用于事件触发
			tempResource := entity.ReadyResource{
				Title:       reqResource.Title,
				Description: reqResource.Description,
				URL:         url,
				Category:    reqResource.Category,
				Tags:        reqResource.Tags,
				Img:         reqResource.Img,
				Source:      reqResource.Source,
				Extra:       reqResource.Extra,
				IP:          reqResource.IP,
				Key:         "",
			}

			// 触发插件事件
			plugins.TriggerReadyResourceAdd(&tempResource, map[string]interface{}{
				"request_id":    c.GetString("request_id"),
				"user_agent":    c.GetHeader("User-Agent"),
				"ip":            c.ClientIP(),
				"batch":         true,
				"is_filtered":   isFiltered,
				"filter_reason": filterReason,
			})
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

	// 触发插件系统待处理资源添加事件
	for _, resource := range resources {
		plugins.TriggerReadyResourceAdd(&resource, map[string]interface{}{
			"request_id": c.GetString("request_id"),
			"user_agent": c.GetHeader("User-Agent"),
			"ip":         c.ClientIP(),
			"source":     "text",
		})
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

// getRetryableErrorCount 统计可重试的错误数量
func getRetryableErrorCount(resources []entity.ReadyResource) int {
	count := 0

	for _, resource := range resources {
		if resource.ErrorMsg != "" {
			errorMsg := strings.ToUpper(resource.ErrorMsg)
			// 检查错误类型标记
			if strings.Contains(resource.ErrorMsg, "[NO_ACCOUNT]") ||
				strings.Contains(resource.ErrorMsg, "[NO_VALID_ACCOUNT]") ||
				strings.Contains(resource.ErrorMsg, "[TRANSFER_FAILED]") ||
				strings.Contains(resource.ErrorMsg, "[LINK_CHECK_FAILED]") {
				count++
			} else if strings.Contains(errorMsg, "没有可用的网盘账号") ||
				strings.Contains(errorMsg, "没有有效的网盘账号") ||
				strings.Contains(errorMsg, "网盘信息获取失败") ||
				strings.Contains(errorMsg, "链接检查失败") {
				count++
			}
		}
	}
	return count
}

// GetReadyResourcesWithErrors 获取有错误信息的待处理资源
func GetReadyResourcesWithErrors(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "100")
	errorFilter := c.Query("error_filter")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 1000 {
		pageSize = 100
	}

	// 获取有错误的资源（分页，包括软删除的）
	resources, total, err := repoManager.ReadyResourceRepository.FindWithErrorsPaginatedIncludingDeleted(page, pageSize, errorFilter)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToReadyResourceResponseList(resources)

	// 统计错误类型
	errorTypeStats := make(map[string]int)
	for _, resource := range resources {
		if resource.ErrorMsg != "" {
			// 尝试从错误信息中提取错误类型
			if len(resource.ErrorMsg) > 0 && resource.ErrorMsg[0] == '[' {
				endIndex := strings.Index(resource.ErrorMsg, "]")
				if endIndex > 0 {
					errorType := resource.ErrorMsg[1:endIndex]
					errorTypeStats[errorType]++
				} else {
					errorTypeStats["UNKNOWN"]++
				}
			} else {
				// 如果没有错误类型标记，尝试从错误信息中推断
				errorMsg := strings.ToUpper(resource.ErrorMsg)
				if strings.Contains(errorMsg, "不支持的链接") {
					errorTypeStats["UNSUPPORTED_LINK"]++
				} else if strings.Contains(errorMsg, "链接无效") {
					errorTypeStats["INVALID_LINK"]++
				} else if strings.Contains(errorMsg, "没有可用的网盘账号") {
					errorTypeStats["NO_ACCOUNT"]++
				} else if strings.Contains(errorMsg, "没有有效的网盘账号") {
					errorTypeStats["NO_VALID_ACCOUNT"]++
				} else if strings.Contains(errorMsg, "网盘信息获取失败") {
					errorTypeStats["TRANSFER_FAILED"]++
				} else if strings.Contains(errorMsg, "创建网盘服务失败") {
					errorTypeStats["SERVICE_CREATION_FAILED"]++
				} else if strings.Contains(errorMsg, "处理标签失败") {
					errorTypeStats["TAG_PROCESSING_FAILED"]++
				} else if strings.Contains(errorMsg, "处理分类失败") {
					errorTypeStats["CATEGORY_PROCESSING_FAILED"]++
				} else if strings.Contains(errorMsg, "资源保存失败") {
					errorTypeStats["RESOURCE_SAVE_FAILED"]++
				} else if strings.Contains(errorMsg, "未找到对应的平台ID") {
					errorTypeStats["PLATFORM_NOT_FOUND"]++
				} else if strings.Contains(errorMsg, "链接检查失败") {
					errorTypeStats["LINK_CHECK_FAILED"]++
				} else {
					errorTypeStats["UNKNOWN"]++
				}
			}
		}
	}

	SuccessResponse(c, gin.H{
		"data":            responses,
		"page":            page,
		"page_size":       pageSize,
		"total":           total,
		"count":           len(resources),
		"error_stats":     errorTypeStats,
		"retryable_count": getRetryableErrorCount(resources),
	})
}

// ClearErrorMsg 清除指定资源的错误信息
func ClearErrorMsg(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.ReadyResourceRepository.ClearErrorMsg(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "错误信息已清除"})
}

// RetryFailedResources 重试失败的资源
func RetryFailedResources(c *gin.Context) {
	// 获取有错误的资源
	resources, err := repoManager.ReadyResourceRepository.FindWithErrors()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(resources) == 0 {
		SuccessResponse(c, gin.H{
			"message": "没有需要重试的资源",
			"count":   0,
		})
		return
	}

	// 只重试可重试的错误
	clearedCount := 0
	skippedCount := 0

	for _, resource := range resources {
		isRetryable := false
		errorMsg := strings.ToUpper(resource.ErrorMsg)

		// 检查错误类型标记
		if strings.Contains(resource.ErrorMsg, "[NO_ACCOUNT]") ||
			strings.Contains(resource.ErrorMsg, "[NO_VALID_ACCOUNT]") ||
			strings.Contains(resource.ErrorMsg, "[TRANSFER_FAILED]") ||
			strings.Contains(resource.ErrorMsg, "[LINK_CHECK_FAILED]") {
			isRetryable = true
		} else if strings.Contains(errorMsg, "没有可用的网盘账号") ||
			strings.Contains(errorMsg, "没有有效的网盘账号") ||
			strings.Contains(errorMsg, "网盘信息获取失败") ||
			strings.Contains(errorMsg, "链接检查失败") {
			isRetryable = true
		}

		if isRetryable {
			if err := repoManager.ReadyResourceRepository.ClearErrorMsg(resource.ID); err == nil {
				clearedCount++
			}
		} else {
			skippedCount++
		}
	}

	SuccessResponse(c, gin.H{
		"message":         "已清除可重试资源的错误信息，资源将在下次调度时重新处理",
		"total_count":     len(resources),
		"cleared_count":   clearedCount,
		"skipped_count":   skippedCount,
		"retryable_count": getRetryableErrorCount(resources),
	})
}

// BatchRestoreToReadyPool 批量将失败资源重新放入待处理池
func BatchRestoreToReadyPool(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		ErrorResponse(c, "资源ID列表不能为空", http.StatusBadRequest)
		return
	}

	successCount := 0
	failedCount := 0

	for _, id := range req.IDs {
		// 清除错误信息并恢复软删除的资源
		err := repoManager.ReadyResourceRepository.ClearErrorMsgAndRestore(id)
		if err != nil {
			failedCount++
			continue
		}
		successCount++
	}

	SuccessResponse(c, gin.H{
		"message":       "批量重新放入待处理池操作完成",
		"total_count":   len(req.IDs),
		"success_count": successCount,
		"failed_count":  failedCount,
	})
}

// BatchRestoreToReadyPoolByQuery 根据查询条件批量将失败资源重新放入待处理池
func BatchRestoreToReadyPoolByQuery(c *gin.Context) {
	var req struct {
		ErrorFilter string `json:"error_filter"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 根据查询条件获取所有符合条件的资源
	resources, err := repoManager.ReadyResourceRepository.FindWithErrorsByQuery(req.ErrorFilter)
	if err != nil {
		ErrorResponse(c, "查询资源失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(resources) == 0 {
		SuccessResponse(c, gin.H{
			"message":       "没有找到符合条件的资源",
			"total_count":   0,
			"success_count": 0,
			"failed_count":  0,
		})
		return
	}

	successCount := 0
	failedCount := 0
	for _, resource := range resources {
		err := repoManager.ReadyResourceRepository.ClearErrorMsgAndRestore(resource.ID)
		if err != nil {
			failedCount++
			continue
		}
		successCount++
	}

	SuccessResponse(c, gin.H{
		"message":       "批量重新放入待处理池操作完成",
		"total_count":   len(resources),
		"success_count": successCount,
		"failed_count":  failedCount,
	})
}

// ClearAllErrorsByQuery 根据查询条件批量清除错误信息并删除资源
func ClearAllErrorsByQuery(c *gin.Context) {
	var req struct {
		ErrorFilter string `json:"error_filter"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 根据查询条件批量删除失败资源
	affectedRows, err := repoManager.ReadyResourceRepository.ClearAllErrorsByQuery(req.ErrorFilter)
	if err != nil {
		ErrorResponse(c, "批量删除失败资源失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message":       "批量删除失败资源操作完成",
		"affected_rows": affectedRows,
	})
}
