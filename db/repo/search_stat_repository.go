package repo

import (
	"fmt"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"gorm.io/gorm"
)

// SearchStatRepository 搜索统计Repository接口
type SearchStatRepository interface {
	BaseRepository[entity.SearchStat]
	RecordSearch(keyword, ip, userAgent string) error
	GetDailyStats(days int) ([]entity.DailySearchStat, error)
	GetHotKeywords(days int, limit int) ([]entity.KeywordStat, error)
	GetSearchTrend(days int) ([]entity.DailySearchStat, error)
	GetKeywordTrend(keyword string, days int) ([]entity.DailySearchStat, error)
	GetSummary() (map[string]int64, error)
}

// SearchStatRepositoryImpl 搜索统计Repository实现
type SearchStatRepositoryImpl struct {
	BaseRepositoryImpl[entity.SearchStat]
}

// NewSearchStatRepository 创建搜索统计Repository
func NewSearchStatRepository(db *gorm.DB) SearchStatRepository {
	return &SearchStatRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.SearchStat]{db: db},
	}
}

// RecordSearch 记录搜索（每次都插入新记录）
func (r *SearchStatRepositoryImpl) RecordSearch(keyword, ip, userAgent string) error {
	stat := entity.SearchStat{
		Keyword:   keyword,
		Count:     1,
		Date:      utils.GetCurrentTime(), // 可保留 date 字段，实际用 created_at 统计
		IP:        ip,
		UserAgent: userAgent,
	}
	return r.db.Create(&stat).Error
}

// GetDailyStats 获取每日统计
func (r *SearchStatRepositoryImpl) GetDailyStats(days int) ([]entity.DailySearchStat, error) {
	var stats []entity.DailySearchStat

	query := fmt.Sprintf(`
		SELECT 
			date,
			SUM(count) as total_searches,
			COUNT(DISTINCT keyword) as unique_keywords
		FROM search_stats 
		WHERE date >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY date 
		ORDER BY date DESC
	`, days)

	err := r.db.Raw(query).Scan(&stats).Error
	return stats, err
}

// GetHotKeywords 获取热门关键词
func (r *SearchStatRepositoryImpl) GetHotKeywords(days int, limit int) ([]entity.KeywordStat, error) {
	var keywords []entity.KeywordStat

	query := fmt.Sprintf(`
		SELECT 
			keyword,
			SUM(count) as count,
			RANK() OVER (ORDER BY SUM(count) DESC) as rank
		FROM search_stats 
		WHERE date >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY keyword 
		ORDER BY count DESC 
		LIMIT ?
	`, days)

	err := r.db.Raw(query, limit).Scan(&keywords).Error
	return keywords, err
}

// GetSearchTrend 获取搜索趋势
func (r *SearchStatRepositoryImpl) GetSearchTrend(days int) ([]entity.DailySearchStat, error) {
	var stats []entity.DailySearchStat

	query := fmt.Sprintf(`
		SELECT 
			date,
			SUM(count) as total_searches,
			COUNT(DISTINCT keyword) as unique_keywords
		FROM search_stats 
		WHERE date >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY date 
		ORDER BY date ASC
	`, days)

	err := r.db.Raw(query).Scan(&stats).Error
	return stats, err
}

// GetKeywordTrend 获取关键词趋势
func (r *SearchStatRepositoryImpl) GetKeywordTrend(keyword string, days int) ([]entity.DailySearchStat, error) {
	var stats []entity.DailySearchStat

	query := fmt.Sprintf(`
		SELECT 
			date,
			SUM(count) as total_searches,
			COUNT(DISTINCT keyword) as unique_keywords
		FROM search_stats 
		WHERE keyword = ? AND date >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY date 
		ORDER BY date ASC
	`, days)

	err := r.db.Raw(query, keyword).Scan(&stats).Error
	return stats, err
}

// GetSummary 获取搜索统计汇总
func (r *SearchStatRepositoryImpl) GetSummary() (map[string]int64, error) {
	var total, today, week, month, keywords int64
	now := utils.GetCurrentTime()
	todayStr := now.Format(utils.TimeFormatDate)
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1).Format(utils.TimeFormatDate) // 周一
	monthStart := now.Format("2006-01") + "-01"

	// 总搜索次数
	if err := r.db.Model(&entity.SearchStat{}).Count(&total).Error; err != nil {
		return nil, err
	}
	// 今日搜索次数
	if err := r.db.Model(&entity.SearchStat{}).Where("DATE(created_at) = ?", todayStr).Count(&today).Error; err != nil {
		return nil, err
	}
	// 本周搜索次数
	if err := r.db.Model(&entity.SearchStat{}).Where("created_at >= ?", weekStart).Count(&week).Error; err != nil {
		return nil, err
	}
	// 本月搜索次数
	if err := r.db.Model(&entity.SearchStat{}).Where("created_at >= ?", monthStart).Count(&month).Error; err != nil {
		return nil, err
	}
	// 总关键词数
	if err := r.db.Model(&entity.SearchStat{}).Distinct("keyword").Count(&keywords).Error; err != nil {
		return nil, err
	}
	return map[string]int64{
		"total":    total,
		"today":    today,
		"week":     week,
		"month":    month,
		"keywords": keywords,
	}, nil
}
