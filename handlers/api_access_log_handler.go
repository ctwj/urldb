package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/converter"

	"github.com/gin-gonic/gin"
)

// GetAPIAccessLogs 获取API访问日志
func GetAPIAccessLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	endpoint := c.Query("endpoint")
	ip := c.Query("ip")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// 设置为当天结束时间
			endOfDay := parsed.Add(24*time.Hour - time.Second)
			endDate = &endOfDay
		}
	}

	// 获取分页数据
	logs, total, err := repoManager.APIAccessLogRepository.FindWithFilters(page, pageSize, startDate, endDate, endpoint, ip)
	if err != nil {
		ErrorResponse(c, "获取API访问日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := converter.ToAPIAccessLogResponseList(logs)

	SuccessResponse(c, gin.H{
		"data":  response,
		"total": int(total),
		"page":  page,
		"limit": pageSize,
	})
}

// GetAPIAccessLogSummary 获取API访问日志汇总
func GetAPIAccessLogSummary(c *gin.Context) {
	summary, err := repoManager.APIAccessLogRepository.GetSummary()
	if err != nil {
		ErrorResponse(c, "获取API访问日志汇总失败: "+err.Error(), 500)
		return
	}

	response := converter.ToAPIAccessLogSummaryResponse(summary)
	SuccessResponse(c, response)
}

// GetAPIAccessLogStats 获取API访问日志统计
func GetAPIAccessLogStats(c *gin.Context) {
	stats, err := repoManager.APIAccessLogRepository.GetStatsByEndpoint()
	if err != nil {
		ErrorResponse(c, "获取API访问日志统计失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := converter.ToAPIAccessLogStatsResponseList(stats)
	SuccessResponse(c, response)
}

// ClearAPIAccessLogs 清理API访问日志
func ClearAPIAccessLogs(c *gin.Context) {
	daysStr := c.Query("days")
	if daysStr == "" {
		ErrorResponse(c, "请提供要清理的天数参数", http.StatusBadRequest)
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		ErrorResponse(c, "无效的天数参数", http.StatusBadRequest)
		return
	}

	err = repoManager.APIAccessLogRepository.ClearOldLogs(days)
	if err != nil {
		ErrorResponse(c, "清理API访问日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "API访问日志清理成功"})
}
