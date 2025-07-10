package handlers

import (
	"net/http"
	"strconv"

	"res_db/db/converter"
	"res_db/db/dto"

	"github.com/gin-gonic/gin"
)

// RecordSearch 记录搜索
func RecordSearch(c *gin.Context) {
	var req dto.SearchStatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取客户端IP和User-Agent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 记录搜索
	err := repoManager.SearchStatRepository.RecordSearch(req.Keyword, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录搜索失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "搜索记录成功"})
}

// GetSearchStats 获取搜索统计总览
func GetSearchStats(c *gin.Context) {
	// 获取今日搜索量
	todayStats, err := repoManager.SearchStatRepository.GetDailyStats(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取今日统计失败"})
		return
	}

	// 获取本周搜索量
	weekStats, err := repoManager.SearchStatRepository.GetDailyStats(7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取本周统计失败"})
		return
	}

	// 获取本月搜索量
	monthStats, err := repoManager.SearchStatRepository.GetDailyStats(30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取本月统计失败"})
		return
	}

	// 获取热门关键词
	hotKeywords, err := repoManager.SearchStatRepository.GetHotKeywords(30, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门关键词失败"})
		return
	}

	// 获取搜索趋势
	searchTrend, err := repoManager.SearchStatRepository.GetSearchTrend(30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取搜索趋势失败"})
		return
	}

	// 计算总搜索量
	var todaySearches, weekSearches, monthSearches int
	if len(todayStats) > 0 {
		todaySearches = todayStats[0].TotalSearches
	}
	for _, stat := range weekStats {
		weekSearches += stat.TotalSearches
	}
	for _, stat := range monthStats {
		monthSearches += stat.TotalSearches
	}

	// 构建趋势数据
	var trendDays []string
	var trendValues []int
	for _, stat := range searchTrend {
		trendDays = append(trendDays, stat.Date.Format("01-02"))
		trendValues = append(trendValues, stat.TotalSearches)
	}

	response := dto.SearchStatsResponse{
		TodaySearches: todaySearches,
		WeekSearches:  weekSearches,
		MonthSearches: monthSearches,
		HotKeywords:   converter.ToHotKeywordResponseList(hotKeywords),
		DailyStats:    converter.ToDailySearchStatResponseList(searchTrend),
		SearchTrend: dto.SearchTrendResponse{
			Days:   trendDays,
			Values: trendValues,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetHotKeywords 获取热门关键词
func GetHotKeywords(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	limitStr := c.DefaultQuery("limit", "10")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的限制参数"})
		return
	}

	keywords, err := repoManager.SearchStatRepository.GetHotKeywords(days, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门关键词失败"})
		return
	}

	response := converter.ToHotKeywordResponseList(keywords)
	c.JSON(http.StatusOK, response)
}

// GetDailyStats 获取每日统计
func GetDailyStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数"})
		return
	}

	stats, err := repoManager.SearchStatRepository.GetDailyStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取每日统计失败"})
		return
	}

	response := converter.ToDailySearchStatResponseList(stats)
	c.JSON(http.StatusOK, response)
}

// GetSearchTrend 获取搜索趋势
func GetSearchTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数"})
		return
	}

	trend, err := repoManager.SearchStatRepository.GetSearchTrend(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取搜索趋势失败"})
		return
	}

	response := converter.ToDailySearchStatResponseList(trend)
	c.JSON(http.StatusOK, response)
}

// GetKeywordTrend 获取关键词趋势
func GetKeywordTrend(c *gin.Context) {
	keyword := c.Param("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "关键词不能为空"})
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的天数参数"})
		return
	}

	trend, err := repoManager.SearchStatRepository.GetKeywordTrend(keyword, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取关键词趋势失败"})
		return
	}

	response := converter.ToDailySearchStatResponseList(trend)
	c.JSON(http.StatusOK, response)
}
