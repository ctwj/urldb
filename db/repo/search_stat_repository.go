package repo

import (
	"time"

	"res_db/db/entity"

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

// RecordSearch 记录搜索
func (r *SearchStatRepositoryImpl) RecordSearch(keyword, ip, userAgent string) error {
	today := time.Now().Truncate(24 * time.Hour)

	// 查找今天是否已有该关键词的记录
	var stat entity.SearchStat
	err := r.db.Where("keyword = ? AND date = ?", keyword, today).First(&stat).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		stat = entity.SearchStat{
			Keyword:   keyword,
			Count:     1,
			Date:      today,
			IP:        ip,
			UserAgent: userAgent,
		}
		return r.db.Create(&stat).Error
	} else if err != nil {
		return err
	}

	// 更新现有记录
	stat.Count++
	stat.IP = ip
	stat.UserAgent = userAgent
	return r.db.Save(&stat).Error
}

// GetDailyStats 获取每日统计
func (r *SearchStatRepositoryImpl) GetDailyStats(days int) ([]entity.DailySearchStat, error) {
	var stats []entity.DailySearchStat

	query := `
		SELECT 
			date,
			SUM(count) as total_searches,
			COUNT(DISTINCT keyword) as unique_keywords
		FROM search_stats 
		WHERE date >= CURRENT_DATE - INTERVAL '? days'
		GROUP BY date 
		ORDER BY date DESC
	`

	err := r.db.Raw(query, days).Scan(&stats).Error
	return stats, err
}

// GetHotKeywords 获取热门关键词
func (r *SearchStatRepositoryImpl) GetHotKeywords(days int, limit int) ([]entity.KeywordStat, error) {
	var keywords []entity.KeywordStat

	query := `
		SELECT 
			keyword,
			SUM(count) as count,
			RANK() OVER (ORDER BY SUM(count) DESC) as rank
		FROM search_stats 
		WHERE date >= CURRENT_DATE - INTERVAL '? days'
		GROUP BY keyword 
		ORDER BY count DESC 
		LIMIT ?
	`

	err := r.db.Raw(query, days, limit).Scan(&keywords).Error
	return keywords, err
}

// GetSearchTrend 获取搜索趋势
func (r *SearchStatRepositoryImpl) GetSearchTrend(days int) ([]entity.DailySearchStat, error) {
	var stats []entity.DailySearchStat

	query := `
		SELECT 
			date,
			SUM(count) as total_searches,
			COUNT(DISTINCT keyword) as unique_keywords
		FROM search_stats 
		WHERE date >= CURRENT_DATE - INTERVAL '? days'
		GROUP BY date 
		ORDER BY date ASC
	`

	err := r.db.Raw(query, days).Scan(&stats).Error
	return stats, err
}

// GetKeywordTrend 获取关键词趋势
func (r *SearchStatRepositoryImpl) GetKeywordTrend(keyword string, days int) ([]entity.DailySearchStat, error) {
	var stats []entity.DailySearchStat

	query := `
		SELECT 
			date,
			SUM(count) as total_searches,
			COUNT(DISTINCT keyword) as unique_keywords
		FROM search_stats 
		WHERE keyword = ? AND date >= CURRENT_DATE - INTERVAL '? days'
		GROUP BY date 
		ORDER BY date ASC
	`

	err := r.db.Raw(query, keyword, days).Scan(&stats).Error
	return stats, err
}
