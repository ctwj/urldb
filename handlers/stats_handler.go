package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
)

// GetStats 获取基础统计信息
func GetStats(c *gin.Context) {
	// 设置响应头，启用缓存
	c.Header("Cache-Control", "public, max-age=60") // 1分钟缓存

	// 获取数据库统计
	var totalResources, totalCategories, totalTags, totalViews int64
	db.DB.Model(&entity.Resource{}).Count(&totalResources)
	db.DB.Model(&entity.Category{}).Count(&totalCategories)
	db.DB.Model(&entity.Tag{}).Count(&totalTags)
	db.DB.Model(&entity.Resource{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)

	// 获取今日数据（在UTC+8时区的0点开始统计）
	now := utils.GetCurrentTime()
	// 使用UTC+8时区的今天0点
	loc, _ := time.LoadLocation("Asia/Shanghai") // UTC+8
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfToday := startOfToday.Add(24 * time.Hour)

	// 今日新增资源数量（从0点开始）
	var todayResources int64
	db.DB.Model(&entity.Resource{}).Where("created_at >= ? AND created_at < ?", startOfToday, endOfToday).Count(&todayResources)

	// 今日更新资源数量（包括新增和修改，从0点开始）
	var todayUpdates int64
	db.DB.Model(&entity.Resource{}).Where("updated_at >= ? AND updated_at < ?", startOfToday, endOfToday).Count(&todayUpdates)

	// 今日浏览量 - 使用访问记录表统计今日访问量
	var todayViews int64
	todayViews, err := repoManager.ResourceViewRepository.GetTodayViews()
	if err != nil {
		utils.Error("获取今日访问量失败: %v", err)
		todayViews = 0
	}

	// 今日搜索量（从0点开始）
	var todaySearches int64
	db.DB.Model(&entity.SearchStat{}).Where("date >= ? AND date < ?", startOfToday.Format(utils.TimeFormatDate), endOfToday.Format(utils.TimeFormatDate)).Count(&todaySearches)

	// 009: 今日失效资源数（按 invalidated_at 日期）
	todayInvalid, errInv := repoManager.ResourceRepository.CountInvalidByDate(startOfToday.Format(utils.TimeFormatDate))
	if errInv != nil {
		utils.Error("获取今日失效资源数失败: %v", errInv)
		todayInvalid = 0
	}

	// 添加调试日志
	utils.Info("统计数据 - 总资源: %d, 总分类: %d, 总标签: %d, 总浏览量: %d",
		totalResources, totalCategories, totalTags, totalViews)
	utils.Info("今日数据 - 新增资源: %d, 今日更新: %d, 今日浏览量: %d, 今日搜索: %d",
		todayResources, todayUpdates, todayViews, todaySearches)

	SuccessResponse(c, gin.H{
		"total_resources":  totalResources,
		"total_categories": totalCategories,
		"total_tags":       totalTags,
		"total_views":      totalViews,
		"today_resources":  todayResources,
		"today_updates":    todayUpdates,
		"today_views":      todayViews,
		"today_searches":   todaySearches,
		"today_invalid":    todayInvalid,
	})
}

// GetPerformanceStats 获取性能监控信息
func GetPerformanceStats(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 获取数据库连接池状态
	sqlDB, err := db.DB.DB()
	var dbStats gin.H
	if err == nil {
		stats := sqlDB.Stats()
		dbStats = gin.H{
			"max_open_connections": stats.MaxOpenConnections,
			"open_connections":     stats.OpenConnections,
			"in_use":               stats.InUse,
			"idle":                 stats.Idle,
			"wait_count":           stats.WaitCount,
			"wait_duration":        stats.WaitDuration,
		}
		// 添加调试日志
		utils.Info("数据库连接池状态 - MaxOpen: %d, Open: %d, InUse: %d, Idle: %d",
			stats.MaxOpenConnections, stats.OpenConnections, stats.InUse, stats.Idle)
	} else {
		dbStats = gin.H{
			"error": "无法获取数据库连接池状态: " + err.Error(),
		}
		utils.Error("获取数据库连接池状态失败: %v", err)
	}

	SuccessResponse(c, gin.H{
		"timestamp": utils.GetCurrentTime().Unix(),
		"memory": gin.H{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"num_gc":      m.NumGC,
			"heap_alloc":  m.HeapAlloc,
			"heap_sys":    m.HeapSys,
			"heap_idle":   m.HeapIdle,
			"heap_inuse":  m.HeapInuse,
		},
		"goroutines": runtime.NumGoroutine(),
		"database":   dbStats,
		"system": gin.H{
			"cpu_count":  runtime.NumCPU(),
			"go_version": runtime.Version(),
		},
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	SuccessResponse(c, gin.H{
		"uptime":     time.Since(startTime).String(),
		"start_time": utils.FormatTime(startTime, utils.TimeFormatDateTime),
		"version":    utils.Version,
		"environment": gin.H{
			"gin_mode": gin.Mode(),
		},
	})
}

// GetViewsTrend 获取访问量趋势数据
func GetViewsTrend(c *gin.Context) {
	// 使用访问记录表获取最近7天的访问量数据
	results, err := repoManager.ResourceViewRepository.GetViewsTrend(7)
	if err != nil {
		utils.Error("获取访问量趋势数据失败: %v", err)
		// 如果获取失败，返回空数据
		results = []map[string]interface{}{}
	}

	// 添加调试日志
	utils.Info("访问量趋势数据: %+v", results)
	for i, result := range results {
		utils.Info("第%d天: 日期=%s, 访问量=%d", i+1, result["date"], result["views"])
	}

	SuccessResponse(c, results)
}

// GetSearchesTrend 获取搜索量趋势数据
func GetSearchesTrend(c *gin.Context) {
	// 获取最近7天的搜索量数据
	var results []gin.H

	// 生成最近7天的日期
	for i := 6; i >= 0; i-- {
		date := utils.GetCurrentTime().AddDate(0, 0, -i)
		dateStr := date.Format(utils.TimeFormatDate)

		// 查询该日期的搜索量（从搜索统计表）
		var searches int64
		db.DB.Model(&entity.SearchStat{}).
			Where("DATE(date) = ?", dateStr).
			Count(&searches)

		// 如果没有搜索记录，返回0
		// 移除模拟数据生成逻辑，只返回真实数据

		results = append(results, gin.H{
			"date":     dateStr,
			"searches": searches,
		})
	}

	// 添加调试日志
	utils.Info("搜索量趋势数据: %+v", results)

	// 添加更详细的调试信息
	for i, result := range results {
		utils.Info("第%d天: 日期=%s, 搜索量=%d", i+1, result["date"], result["searches"])
	}

	SuccessResponse(c, results)
}

// GetSummary 仪表盘首屏聚合统计（含环比昨日与待办）
// 契约: contracts/stats-summary-api.md — 返回 {code,message,data} 统一封装
// 架构说明: 业务逻辑已抽至 services.StatsService.GetSummary()，handler 仅做 HTTP 编排
func GetSummary(c *gin.Context) {
	svc := services.GetDefaultStatsService()
	if svc == nil {
		ErrorResponse(c, "统计服务未初始化", http.StatusInternalServerError)
		return
	}
	summary, err := svc.GetSummary()
	if err != nil {
		utils.Error("GetSummary 失败: %v", err)
		ErrorResponse(c, "获取统计失败", http.StatusInternalServerError)
		return
	}
	SuccessResponse(c, summary)
}

// 记录启动时间
var startTime = utils.GetCurrentTime()
