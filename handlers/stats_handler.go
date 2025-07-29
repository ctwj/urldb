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

	// 获取今日更新数量
	var todayUpdates int64
	today := utils.GetCurrentTime().Format("2006-01-02")
	db.DB.Model(&entity.Resource{}).Where("DATE(updated_at) = ?", today).Count(&todayUpdates)

	SuccessResponse(c, gin.H{
		"total_resources":  totalResources,
		"total_categories": totalCategories,
		"total_tags":       totalTags,
		"total_views":      totalViews,
		"today_updates":    todayUpdates,
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
		"start_time": utils.FormatTime(startTime, "2006-01-02 15:04:05"),
		"version":    utils.Version,
		"environment": gin.H{
			"gin_mode": gin.Mode(),
		},
	})
}

// 记录启动时间
var startTime = utils.GetCurrentTime()
