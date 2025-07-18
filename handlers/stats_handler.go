package handlers

import (
	"runtime"
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
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

	SuccessResponse(c, gin.H{
		"total_resources":  totalResources,
		"total_categories": totalCategories,
		"total_tags":       totalTags,
		"total_views":      totalViews,
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
		dbStats = gin.H{
			"max_open_connections": sqlDB.Stats().MaxOpenConnections,
			"open_connections":     sqlDB.Stats().OpenConnections,
			"in_use":               sqlDB.Stats().InUse,
			"idle":                 sqlDB.Stats().Idle,
		}
	} else {
		dbStats = gin.H{
			"error": "无法获取数据库连接池状态",
		}
	}

	SuccessResponse(c, gin.H{
		"timestamp": time.Now().Unix(),
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
		"start_time": startTime.Format("2006-01-02 15:04:05"),
		"version":    "1.0.0",
		"environment": gin.H{
			"gin_mode": gin.Mode(),
		},
	})
}

// 记录启动时间
var startTime = time.Now()
