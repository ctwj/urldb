package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// TelegramLogEntry Telegram日志条目
type TelegramLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Category  string    `json:"category,omitempty"` // telegram, push, message等
}

// GetTelegramLogs 获取Telegram相关的日志
func GetTelegramLogs(startTime *time.Time, endTime *time.Time, limit int) ([]TelegramLogEntry, error) {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return []TelegramLogEntry{}, nil
	}

	files, err := filepath.Glob(filepath.Join(logDir, "app_*.log"))
	if err != nil {
		return nil, fmt.Errorf("查找日志文件失败: %v", err)
	}

	if len(files) == 0 {
		return []TelegramLogEntry{}, nil
	}

	// 按时间排序，最近的在前面
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	var allEntries []TelegramLogEntry

	// 编译Telegram相关的正则表达式
	telegramRegex := regexp.MustCompile(`(?i)(telegram|bot|推送|消息|频道|群组|\[.*telegram.*\])`)
	messageRegex := regexp.MustCompile(`\[(\w+)\]\s+(\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2})\s+(.*)`)

	for _, file := range files {
		entries, err := parseTelegramLogsFromFile(file, telegramRegex, messageRegex, startTime, endTime)
		if err != nil {
			continue // 跳过读取失败的文件
		}

		allEntries = append(allEntries, entries...)

		// 如果已经达到限制数量，退出
		if limit > 0 && len(allEntries) >= limit {
			break
		}
	}

	// 按时间排序，最新的在前面
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].Timestamp.After(allEntries[j].Timestamp)
	})

	// 限制返回数量
	if limit > 0 && len(allEntries) > limit {
		allEntries = allEntries[:limit]
	}

	// 分类日志
	for i := range allEntries {
		allEntries[i].Category = categorizeLog(allEntries[i].Message)
	}

	return allEntries, nil
}

// parseTelegramLogsFromFile 解析单个日志文件中的Telegram日志
func parseTelegramLogsFromFile(filePath string, telegramRegex, messageRegex *regexp.Regexp, startTime, endTime *time.Time) ([]TelegramLogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []TelegramLogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否是Telegram相关日志
		if !telegramRegex.MatchString(line) {
			continue
		}

		// 解析日志行
		entry, err := parseLogLine(line, messageRegex)
		if err != nil {
			continue
		}

		// 时间过滤
		if startTime != nil && entry.Timestamp.Before(*startTime) {
			continue
		}
		if endTime != nil && entry.Timestamp.After(*endTime) {
			continue
		}

		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// parseLogLine 解析单行日志
func parseLogLine(line string, messageRegex *regexp.Regexp) (TelegramLogEntry, error) {
	// 匹配日志格式: [LEVEL] 2006/01/02 15:04:05 message
	matches := messageRegex.FindStringSubmatch(line)
	if len(matches) < 4 {
		return TelegramLogEntry{}, fmt.Errorf("无法解析日志行")
	}

	level := matches[1]
	timeStr := matches[2]
	message := matches[3]

	// 解析时间
	timestamp, err := time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return TelegramLogEntry{}, fmt.Errorf("时间解析失败: %v", err)
	}

	return TelegramLogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
	}, nil
}

// categorizeLog 对日志进行分类
func categorizeLog(message string) string {
	message = strings.ToLower(message)

	switch {
	case strings.Contains(message, "推送") || strings.Contains(message, "push"):
		return "push"
	case strings.Contains(message, "消息") || strings.Contains(message, "message") || strings.Contains(message, "收到"):
		return "message"
	case strings.Contains(message, "频道") || strings.Contains(message, "群组") || strings.Contains(message, "register"):
		return "channel"
	case strings.Contains(message, "启动") || strings.Contains(message, "停止") || strings.Contains(message, "start") || strings.Contains(message, "stop"):
		return "service"
	default:
		return "general"
	}
}

// GetTelegramLogStats 获取Telegram日志统计
func GetTelegramLogStats(hours int) (map[string]interface{}, error) {
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	entries, err := GetTelegramLogs(&startTime, &endTime, 0)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_logs": len(entries),
		"categories": make(map[string]int),
		"levels":     make(map[string]int),
		"timeline":   make(map[string]int), // 按小时统计
	}

	categoryStats := stats["categories"].(map[string]int)
	levelStats := stats["levels"].(map[string]int)
	timelineStats := stats["timeline"].(map[string]int)

	for _, entry := range entries {
		// 分类统计
		categoryStats[entry.Category]++

		// 级别统计
		levelStats[entry.Level]++

		// 时间线统计（按小时）
		hourKey := entry.Timestamp.Format("2006-01-02 15:00")
		timelineStats[hourKey]++
	}

	return stats, nil
}

// ClearOldTelegramLogs 清理旧的Telegram日志（保留最近N天的日志）
func ClearOldTelegramLogs(daysToKeep int) error {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return nil // 日志目录不存在，无需清理
	}

	files, err := filepath.Glob(filepath.Join(logDir, "app_*.log"))
	if err != nil {
		return fmt.Errorf("查找日志文件失败: %v", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)

	deletedCount := 0
	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			if err := os.Remove(file); err == nil {
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		Info("清理了 %d 个旧的日志文件", deletedCount)
	}

	return nil
}
