package plugin

import (
	"fmt"
	"strconv"
	"strings"
)

// CronFrequency Cron 频率信息
type CronFrequency struct {
	Expression string `json:"expression"` // 原始 Cron 表达式
	Description string `json:"description"` // 友好的描述
	Interval    string `json:"interval"`    // 执行间隔
	NextRun     string `json:"next_run"`     // 下次执行时间（简化版）
}

// ParseCronExpression 解析 Cron 表达式并返回友好的描述
func ParseCronExpression(expression string) *CronFrequency {
	if expression == "" {
		return &CronFrequency{
			Expression: expression,
			Description: "无效的 Cron 表达式",
			Interval:    "未知",
		}
	}

	parts := strings.Fields(expression)
	if len(parts) != 5 {
		return &CronFrequency{
			Expression: expression,
			Description: "无效的 Cron 表达式",
			Interval:    "未知",
		}
	}

	minute, hour, day, month, weekday := parts[0], parts[1], parts[2], parts[3], parts[4]

	// 解析并生成描述
	description := generateCronDescription(minute, hour, day, month, weekday)
	interval := generateCronInterval(minute, hour, day, month, weekday)

	return &CronFrequency{
		Expression: expression,
		Description: description,
		Interval:    interval,
		NextRun:     "每" + interval, // 简化版本
	}
}

// generateCronDescription 生成 Cron 表达式的中文描述
func generateCronDescription(minute, hour, day, month, weekday string) string {
	var desc strings.Builder

	// 处理特殊情况
	if minute == "*" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "每分钟执行一次"
	}

	if minute == "0" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "每小时执行一次"
	}

	if minute == "0" && hour == "0" && day == "*" && month == "*" && weekday == "*" {
		return "每天午夜执行一次"
	}

	// 处理分钟
	if minute != "*" {
		if minute == "0" {
			desc.WriteString("整点")
		} else if minute == "*/1" {
			desc.WriteString("每分钟")
		} else if strings.HasPrefix(minute, "*/") {
			interval := strings.TrimPrefix(minute, "*/")
			desc.WriteString(fmt.Sprintf("每%d分钟", parseInterval(interval)))
		} else {
			desc.WriteString(fmt.Sprintf("第%s分钟", minute))
		}
	}

	// 处理小时
	if hour != "*" {
		if desc.Len() > 0 {
			desc.WriteString("的")
		}
		if hour == "0" {
			desc.WriteString("午夜")
		} else if hour == "*/1" {
			desc.WriteString("每小时")
		} else if strings.HasPrefix(hour, "*/") {
			interval := strings.TrimPrefix(hour, "*/")
			desc.WriteString(fmt.Sprintf("每%d小时", parseInterval(interval)))
		} else {
			desc.WriteString(fmt.Sprintf("%s点", hour))
		}
	}

	// 处理天
	if day != "*" {
		if desc.Len() > 0 {
			desc.WriteString("的")
		}
		if day == "*/1" {
			desc.WriteString("每天")
		} else if strings.HasPrefix(day, "*/") {
			interval := strings.TrimPrefix(day, "*/")
			desc.WriteString(fmt.Sprintf("每%d天", parseInterval(interval)))
		} else {
			desc.WriteString(fmt.Sprintf("每月%s号", day))
		}
	}

	// 处理月份
	if month != "*" {
		if desc.Len() > 0 {
			desc.WriteString("的")
		}
		if month == "*/1" {
			desc.WriteString("每月")
		} else if strings.HasPrefix(month, "*/") {
			interval := strings.TrimPrefix(month, "*/")
			desc.WriteString(fmt.Sprintf("每%d月", parseInterval(interval)))
		} else {
			desc.WriteString(fmt.Sprintf("%s月", month))
		}
	}

	// 处理星期
	if weekday != "*" {
		if desc.Len() > 0 {
			desc.WriteString("的")
		}
		weekdayNames := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
		if weekday == "*/1" {
			desc.WriteString("每天")
		} else if strings.HasPrefix(weekday, "*/") {
			interval := strings.TrimPrefix(weekday, "*/")
			desc.WriteString(fmt.Sprintf("每%d天", parseInterval(interval)))
		} else {
			if dayIndex, err := strconv.Atoi(weekday); err == nil && dayIndex >= 0 && dayIndex <= 6 {
				desc.WriteString(weekdayNames[dayIndex])
			} else {
				desc.WriteString(weekday)
			}
		}
	}

	desc.WriteString("执行")
	return desc.String()
}

// generateCronInterval 生成执行间隔的简化描述
func generateCronInterval(minute, hour, day, month, weekday string) string {
	// 检查最常见的模式
	if minute == "*" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "分钟"
	}

	if minute == "0" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "小时"
	}

	if minute == "0" && hour == "0" && day == "*" && month == "*" && weekday == "*" {
		return "天"
	}

	if minute == "0" && hour == "0" && day == "1" && month == "*" && weekday == "*" {
		return "月"
	}

	if minute == "0" && hour == "0" && day == "*" && month == "*" && weekday == "1" {
		return "周"
	}

	// 处理带间隔的表达式
	if strings.HasPrefix(minute, "*/") {
		interval := strings.TrimPrefix(minute, "*/")
		return fmt.Sprintf("%d分钟", parseInterval(interval))
	}

	if strings.HasPrefix(hour, "*/") {
		interval := strings.TrimPrefix(hour, "*/")
		return fmt.Sprintf("%d小时", parseInterval(interval))
	}

	if strings.HasPrefix(day, "*/") {
		interval := strings.TrimPrefix(day, "*/")
		return fmt.Sprintf("%d天", parseInterval(interval))
	}

	// 默认返回自定义
	return "自定义"
}

// parseInterval 解析间隔字符串
func parseInterval(intervalStr string) int {
	if interval, err := strconv.Atoi(intervalStr); err == nil {
		return interval
	}
	return 1
}

// GetCronFrequencyColor 获取频率对应的颜色类型
func GetCronFrequencyColor(interval string) string {
	switch {
	case strings.Contains(interval, "分钟"):
		return "warning" // 黄色，表示频繁
	case strings.Contains(interval, "小时"):
		return "info" // 蓝色，表示中等频率
	case strings.Contains(interval, "天"):
		return "success" // 绿色，表示低频率
	case strings.Contains(interval, "周"):
		return "success" // 绿色，表示低频率
	case strings.Contains(interval, "月"):
		return "default" // 灰色，表示很低频率
	default:
		return "default" // 灰色，自定义
	}
}