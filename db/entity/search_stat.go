package entity

import (
	"time"

	"gorm.io/gorm"
)

// SearchStat 搜索统计模型
type SearchStat struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Keyword   string         `json:"keyword" gorm:"size:255;not null;comment:搜索关键词"`
	Count     int            `json:"count" gorm:"default:1;comment:搜索次数"`
	Date      time.Time      `json:"date" gorm:"type:date;not null;comment:搜索日期"`
	IP        string         `json:"ip" gorm:"size:45;comment:用户IP"`
	UserAgent string         `json:"user_agent" gorm:"size:500;comment:用户代理"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (SearchStat) TableName() string {
	return "search_stats"
}

// DailySearchStat 每日搜索统计
type DailySearchStat struct {
	Date           time.Time `json:"date"`
	TotalSearches  int       `json:"total_searches"`
	UniqueKeywords int       `json:"unique_keywords"`
}

// KeywordStat 关键词统计
type KeywordStat struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
	Rank    int    `json:"rank"`
}
