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

// LogEntry 日志条目
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	File      string
	Line      int
}

// LogViewer 日志查看器
type LogViewer struct {
	logDir string
}

// NewLogViewer 创建新的日志查看器
func NewLogViewer(logDir string) *LogViewer {
	return &LogViewer{
		logDir: logDir,
	}
}

// GetLogFiles 获取所有日志文件
func (lv *LogViewer) GetLogFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(lv.logDir, "*.log"))
	if err != nil {
		return nil, err
	}

	// 按修改时间排序，最新的在前
	sort.Slice(files, func(i, j int) bool {
		info1, _ := os.Stat(files[i])
		info2, _ := os.Stat(files[j])
		return info1.ModTime().After(info2.ModTime())
	})

	return files, nil
}

// ReadLogFile 读取日志文件
func (lv *LogViewer) ReadLogFile(filename string, lines int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []string
	scanner := bufio.NewScanner(file)

	// 如果指定了行数，先读取所有行到切片中
	if lines > 0 {
		var allLines []string
		for scanner.Scan() {
			allLines = append(allLines, scanner.Text())
		}

		// 返回最后N行
		start := len(allLines) - lines
		if start < 0 {
			start = 0
		}
		result = allLines[start:]
	} else {
		// 读取所有行
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}
	}

	return result, scanner.Err()
}

// SearchLogs 搜索日志
func (lv *LogViewer) SearchLogs(pattern string, files []string) ([]LogEntry, error) {
	var results []LogEntry
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("无效的正则表达式: %v", err)
	}

	for _, file := range files {
		entries, err := lv.searchFile(file, regex)
		if err != nil {
			Error("搜索文件 %s 失败: %v", file, err)
			continue
		}
		results = append(results, entries...)
	}

	// 按时间排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.Before(results[j].Timestamp)
	})

	return results, nil
}

// searchFile 搜索单个文件
func (lv *LogViewer) searchFile(filename string, regex *regexp.Regexp) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []LogEntry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if regex.MatchString(line) {
			entry := lv.parseLogLine(line)
			results = append(results, entry)
		}
	}

	return results, scanner.Err()
}

// parseLogLine 解析日志行
func (lv *LogViewer) parseLogLine(line string) LogEntry {
	entry := LogEntry{
		Message: line,
	}

	// 尝试解析日志格式 [DEBUG] [filename.go:123] message
	parts := strings.SplitN(line, " ", 3)
	if len(parts) >= 3 {
		// 解析级别
		if strings.HasPrefix(parts[0], "[") && strings.HasSuffix(parts[0], "]") {
			entry.Level = strings.Trim(parts[0], "[]")
		}

		// 解析文件位置
		if strings.HasPrefix(parts[1], "[") && strings.HasSuffix(parts[1], "]") {
			fileInfo := strings.Trim(parts[1], "[]")
			if strings.Contains(fileInfo, ":") {
				fileParts := strings.Split(fileInfo, ":")
				if len(fileParts) == 2 {
					entry.File = fileParts[0]
					fmt.Sscanf(fileParts[1], "%d", &entry.Line)
				}
			}
		}

		// 解析时间戳（如果存在）
		if len(parts) >= 3 {
			// 尝试解析时间戳格式
			timeStr := parts[2]
			if len(timeStr) >= 19 {
				if t, err := time.Parse("2006/01/02 15:04:05", timeStr[:19]); err == nil {
					entry.Timestamp = t
				}
			}
		}
	}

	return entry
}

// GetLogStats 获取日志统计信息
func (lv *LogViewer) GetLogStats(files []string) (map[string]int, error) {
	stats := map[string]int{
		"total":   0,
		"debug":   0,
		"info":    0,
		"warn":    0,
		"error":   0,
		"fatal":   0,
		"unknown": 0,
	}

	for _, file := range files {
		fileStats, err := lv.getFileStats(file)
		if err != nil {
			Error("获取文件 %s 统计失败: %v", file, err)
			continue
		}

		for level, count := range fileStats {
			stats[level] += count
		}
	}

	return stats, nil
}

// getFileStats 获取单个文件的统计信息
func (lv *LogViewer) getFileStats(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := map[string]int{
		"total":   0,
		"debug":   0,
		"info":    0,
		"warn":    0,
		"error":   0,
		"fatal":   0,
		"unknown": 0,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stats["total"]++
		line := scanner.Text()

		// 统计各级别日志数量
		if strings.Contains(line, "[DEBUG]") {
			stats["debug"]++
		} else if strings.Contains(line, "[INFO]") {
			stats["info"]++
		} else if strings.Contains(line, "[WARN]") {
			stats["warn"]++
		} else if strings.Contains(line, "[ERROR]") {
			stats["error"]++
		} else if strings.Contains(line, "[FATAL]") {
			stats["fatal"]++
		} else {
			stats["unknown"]++
		}
	}

	return stats, scanner.Err()
}

// TailLog 实时跟踪日志文件
func (lv *LogViewer) TailLog(filename string, callback func(string)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 移动到文件末尾
	file.Seek(0, 2)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		callback(scanner.Text())
	}

	return scanner.Err()
}

// CleanOldLogs 清理旧日志文件
func (lv *LogViewer) CleanOldLogs(days int) error {
	files, err := lv.GetLogFiles()
	if err != nil {
		return err
	}

	cutoffTime := GetCurrentTime().AddDate(0, 0, -days)
	deletedCount := 0

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			if err := os.Remove(file); err != nil {
				Error("删除旧日志文件失败 %s: %v", file, err)
			} else {
				deletedCount++
				Info("已删除旧日志文件: %s", file)
			}
		}
	}

	Info("清理完成，共删除 %d 个旧日志文件", deletedCount)
	return nil
}
