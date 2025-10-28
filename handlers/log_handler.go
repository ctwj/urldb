package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// GetSystemLogs 获取系统日志
func GetSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	level := c.Query("level")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	search := c.Query("search")

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

	// 使用日志查看器获取日志
	logViewer := utils.NewLogViewer("logs")

	// 获取日志文件列表
	logFiles, err := logViewer.GetLogFiles()
	if err != nil {
		ErrorResponse(c, "获取日志文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果指定了日期范围，只选择对应日期的日志文件
	if startDate != nil || endDate != nil {
		var filteredFiles []string
		for _, file := range logFiles {
			fileInfo, err := utils.GetFileInfo(file)
			if err != nil {
				continue
			}

			shouldInclude := true
			if startDate != nil {
				if fileInfo.ModTime().Before(*startDate) {
					shouldInclude = false
				}
			}
			if endDate != nil {
				if fileInfo.ModTime().After(*endDate) {
					shouldInclude = false
				}
			}

			if shouldInclude {
				filteredFiles = append(filteredFiles, file)
			}
		}
		logFiles = filteredFiles
	}

	// 限制读取的文件数量以提高性能
	if len(logFiles) > 10 {
		logFiles = logFiles[:10] // 只处理最近的10个文件
	}

	var allLogs []utils.LogEntry
	for _, file := range logFiles {
		// 读取日志文件
		fileLogs, err := logViewer.ParseLogEntriesFromFile(file, level, search)
		if err != nil {
			utils.Error("解析日志文件失败 %s: %v", file, err)
			continue
		}
		allLogs = append(allLogs, fileLogs...)
	}

	// 按时间排序（最新的在前）
	utils.SortLogEntriesByTime(allLogs, false)

	// 应用分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(allLogs) {
		start = len(allLogs)
	}
	if end > len(allLogs) {
		end = len(allLogs)
	}

	pagedLogs := allLogs[start:end]

	SuccessResponse(c, gin.H{
		"data":  pagedLogs,
		"total": len(allLogs),
		"page":  page,
		"limit": pageSize,
	})
}

// GetSystemLogFiles 获取系统日志文件列表
func GetSystemLogFiles(c *gin.Context) {
	logViewer := utils.NewLogViewer("logs")
	files, err := logViewer.GetLogFiles()
	if err != nil {
		ErrorResponse(c, "获取日志文件列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个文件的详细信息
	var fileInfos []gin.H
	for _, file := range files {
		info, err := utils.GetFileInfo(file)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, gin.H{
			"name":     info.Name(),
			"size":     info.Size(),
			"mod_time": info.ModTime(),
			"path":     file,
		})
	}

	SuccessResponse(c, gin.H{
		"data": fileInfos,
	})
}

// GetSystemLogSummary 获取系统日志统计摘要
func GetSystemLogSummary(c *gin.Context) {
	logViewer := utils.NewLogViewer("logs")
	files, err := logViewer.GetLogFiles()
	if err != nil {
		ErrorResponse(c, "获取日志文件列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取统计信息
	stats, err := logViewer.GetLogStats(files)
	if err != nil {
		ErrorResponse(c, "获取日志统计信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"summary":    stats,
		"files_count": len(files),
	})
}

// ClearSystemLogs 清理系统日志
func ClearSystemLogs(c *gin.Context) {
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

	logViewer := utils.NewLogViewer("logs")
	err = logViewer.CleanOldLogs(days)
	if err != nil {
		ErrorResponse(c, "清理系统日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "系统日志清理成功"})
}