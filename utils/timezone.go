package utils

import (
	"os"
	"time"
)

// InitTimezone 初始化时区设置
func InitTimezone() {
	// 从环境变量获取时区配置
	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		// 默认使用上海时间
		timezone = "Asia/Shanghai"
		Info("未配置时区，使用默认时区: %s", timezone)
	} else {
		Info("使用配置的时区: %s", timezone)
	}

	// 设置时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		Error("加载时区失败: %v，使用系统默认时区", err)
		return
	}

	// 设置全局时区
	time.Local = loc
	Info("时区设置成功: %s", timezone)
}

// GetCurrentTime 获取当前时间（使用配置的时区）
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeString 获取当前时间字符串（使用配置的时区）
func GetCurrentTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetCurrentTimeRFC3339 获取当前时间RFC3339格式（使用配置的时区）
func GetCurrentTimeRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

// ParseTime 解析时间字符串（使用配置的时区）
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}

// FormatTime 格式化时间（使用配置的时区）
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}
