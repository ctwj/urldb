package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"

	"github.com/gin-gonic/gin"
)

// RecordSearch 记录搜索
func RecordSearch(c *gin.Context) {
	var req dto.SearchStatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 获取客户端IP和User-Agent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 记录搜索
	err := repoManager.SearchStatRepository.RecordSearch(req.Keyword, ip, userAgent)
	if err != nil {
		ErrorResponse(c, "记录搜索失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "搜索记录成功"})
}

// GetSearchStats 获取搜索统计（使用全局repoManager）
func GetSearchStats(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	stats, total, err := repoManager.SearchStatRepository.FindWithPagination(page, pageSize)
	if err != nil {
		ErrorResponse(c, "获取搜索统计失败", http.StatusInternalServerError)
		return
	}

	response := converter.ToSearchStatResponseList(stats)

	SuccessResponse(c, gin.H{
		"data":  response,
		"total": int(total),
	})
}

// GetHotKeywords 获取热门关键词（使用全局repoManager）
func GetHotKeywords(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	limitStr := c.DefaultQuery("limit", "10")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		ErrorResponse(c, "无效的天数参数", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ErrorResponse(c, "无效的限制参数", http.StatusBadRequest)
		return
	}

	keywords, err := repoManager.SearchStatRepository.GetHotKeywords(days, limit)
	if err != nil {
		ErrorResponse(c, "获取热门关键词失败", http.StatusInternalServerError)
		return
	}

	response := converter.ToHotKeywordResponseList(keywords)
	SuccessResponse(c, response)
}

// GetDailyStats 获取每日统计（使用全局repoManager）
func GetDailyStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		ErrorResponse(c, "无效的天数参数", http.StatusBadRequest)
		return
	}

	stats, err := repoManager.SearchStatRepository.GetDailyStats(days)
	if err != nil {
		ErrorResponse(c, "获取每日统计失败", http.StatusInternalServerError)
		return
	}

	response := converter.ToDailySearchStatResponseList(stats)
	SuccessResponse(c, response)
}

// GetSearchTrend 获取搜索趋势（使用全局repoManager）
func GetSearchTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		ErrorResponse(c, "无效的天数参数", http.StatusBadRequest)
		return
	}

	trend, err := repoManager.SearchStatRepository.GetSearchTrend(days)
	if err != nil {
		ErrorResponse(c, "获取搜索趋势失败", http.StatusInternalServerError)
		return
	}

	response := converter.ToDailySearchStatResponseList(trend)
	SuccessResponse(c, response)
}

// GetKeywordTrend 获取关键词趋势（使用全局repoManager）
func GetKeywordTrend(c *gin.Context) {
	keyword := c.Param("keyword")
	if keyword == "" {
		ErrorResponse(c, "关键词不能为空", http.StatusBadRequest)
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		ErrorResponse(c, "无效的天数参数", http.StatusBadRequest)
		return
	}

	trend, err := repoManager.SearchStatRepository.GetKeywordTrend(keyword, days)
	if err != nil {
		ErrorResponse(c, "获取关键词趋势失败", http.StatusInternalServerError)
		return
	}

	response := converter.ToDailySearchStatResponseList(trend)
	SuccessResponse(c, response)
}
