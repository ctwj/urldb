package handlers

import (
	"runtime"
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
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

	// 获取今日数据
	today := utils.GetTodayString()

	// 今日新增资源数量
	var todayResources int64
	db.DB.Model(&entity.Resource{}).Where("DATE(created_at) = ?", today).Count(&todayResources)

	// 今日更新资源数量（包括新增和修改）
	var todayUpdates int64
	db.DB.Model(&entity.Resource{}).Where("DATE(updated_at) = ?", today).Count(&todayUpdates)

	// 今日浏览量 - 使用访问记录表统计今日访问量
	var todayViews int64
	todayViews, err := repoManager.ResourceViewRepository.GetTodayViews()
	if err != nil {
		utils.Error("获取今日访问量失败: %v", err)
		todayViews = 0
	}

	// 今日搜索量
	var todaySearches int64
	db.DB.Model(&entity.SearchStat{}).Where("DATE(date) = ?", today).Count(&todaySearches)

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

// 记录启动时间
var startTime = utils.GetCurrentTime()
