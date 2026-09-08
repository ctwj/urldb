package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// 开放资源查询接口（016-api-access-application，US3，contracts/api.md C1）
// 鉴权由 middleware.ApiKeyAuth() 完成；限流由 middleware.ApiRateLimit() 完成

// openMaxPageSize 单次返回数量上限（FR-010：最大 100 条）
const openMaxPageSize = 100

// OpenSearchResources 公开资源搜索（X-API-Key）
func OpenSearchResources(c *gin.Context) {
	startTime := time.Now()

	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		ErrorResponse(c, "keyword 不能为空", http.StatusBadRequest)
		return
	}
	if len([]rune(keyword)) > 100 {
		ErrorResponse(c, "keyword 不能超过 100 字", http.StatusBadRequest)
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "20"))
	if err != nil || size <= 0 {
		ErrorResponse(c, "size 参数非法", http.StatusBadRequest)
		return
	}
	// FR-010：超过上限按 100 处理；接口无翻页，固定取第一页
	if size > openMaxPageSize {
		size = openMaxPageSize
	}

	var items []gin.H

	// 优先 Meilisearch（is_valid 过滤在索引侧；is_public 回查 DB 校验，FR-011 双保险）
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		service := meilisearchManager.GetService()
		if service != nil {
			docs, _, searchErr := service.Search(keyword, map[string]interface{}{"is_valid": true}, 1, size)
			if searchErr == nil {
				ids := make([]uint, 0, len(docs))
				for _, doc := range docs {
					ids = append(ids, doc.ID)
				}
				publicSet := make(map[uint]bool, len(ids))
				if len(ids) > 0 {
					if resources, err := repoManager.ResourceRepository.FindByIDs(ids); err == nil {
						for i := range resources {
							publicSet[resources[i].ID] = resources[i].IsPublic
						}
					} else {
						utils.Error("OpenSearchResources - 公开状态回查失败: %v", err)
					}
				}
				for _, doc := range docs {
					if !publicSet[doc.ID] {
						continue
					}
					// 契约：仅返回 url（save_url 优先）、cover，不返回 id/save_url/category
					itemURL := doc.SaveURL
					if itemURL == "" {
						itemURL = doc.URL
					}
					items = append(items, gin.H{
						"title":       doc.Title,
						"description": doc.Description,
						"url":         itemURL,
						"key":         doc.Key,
						"cover":       doc.Cover,
						"tags":        doc.Tags,
						"created_at":  doc.CreatedAt,
					})
				}
			} else {
				utils.Error("OpenSearchResources - Meilisearch 搜索失败，回退数据库搜索: %v", searchErr)
			}
		}
	}

	// Meilisearch 未启用/失败：数据库回退（SQL 层同条件，FR-011）
	if items == nil {
		resources, _, err := repoManager.ResourceRepository.SearchWithFilters(map[string]interface{}{
			"search":    keyword,
			"is_valid":  true,
			"is_public": true,
			"page":      1,
			"page_size": size,
		})
		if err != nil {
			utils.Error("OpenSearchResources - 数据库搜索失败: %v", err)
			ErrorResponse(c, "搜索失败，请稍后重试", http.StatusInternalServerError)
			return
		}
		items = make([]gin.H, 0, len(resources))
		for i := range resources {
			r := resources[i]
			// 契约：仅返回 url（save_url 优先）、cover，不返回 id/save_url/category/pan_name/file_size
			itemURL := r.SaveURL
			if itemURL == "" {
				itemURL = r.URL
			}
			tagNames := make([]string, 0, len(r.Tags))
			for _, t := range r.Tags {
				tagNames = append(tagNames, t.Name)
			}
			items = append(items, gin.H{
				"title":       r.Title,
				"description": r.Description,
				"url":         itemURL,
				"key":         r.Key,
				"cover":       r.Cover,
				"tags":        tagNames,
				"created_at":  r.CreatedAt,
			})
		}
	}

	// 访问留痕（复用 api_access_logs，异步避免阻塞响应）
	go func(statusCode, processCount int, params string) {
		log := &entity.APIAccessLog{
			IP:             c.ClientIP(),
			UserAgent:      c.GetHeader("User-Agent"),
			Endpoint:       c.Request.URL.Path,
			Method:         c.Request.Method,
			RequestParams:  params,
			ResponseStatus: statusCode,
			ProcessCount:   processCount,
			ProcessingTime: time.Since(startTime).Milliseconds(),
		}
		if err := repoManager.APIAccessLogRepository.Create(log); err != nil {
			utils.Error("OpenSearchResources - 访问日志写入失败: %v", err)
		}
	}(http.StatusOK, len(items), c.Request.URL.RawQuery)

	SuccessResponse(c, gin.H{
		"list": items,
		"size": size,
	})
}
