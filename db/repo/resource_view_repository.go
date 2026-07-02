package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// ResourceViewRepository 资源访问记录仓库接口
type ResourceViewRepository interface {
	BaseRepository[entity.ResourceView]
	RecordView(resourceID uint, ipAddress, userAgent, source string) error
	GetTotalViews() (int64, error)
	GetPanDistribution() ([]map[string]interface{}, error)
	GetSourceDistribution() ([]map[string]interface{}, error)
	GetTodayViews() (int64, error)
	GetViewsByDate(date string) (int64, error)
	GetViewsTrend(days int) ([]map[string]interface{}, error)
	GetResourceViews(resourceID uint, limit int) ([]entity.ResourceView, error)
}

// ResourceViewRepositoryImpl 资源访问记录仓库实现
type ResourceViewRepositoryImpl struct {
	BaseRepositoryImpl[entity.ResourceView]
}

// NewResourceViewRepository 创建资源访问记录仓库
func NewResourceViewRepository(db *gorm.DB) ResourceViewRepository {
	return &ResourceViewRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.ResourceView]{db: db},
	}
}

// RecordView 记录资源访问
func (r *ResourceViewRepositoryImpl) RecordView(resourceID uint, ipAddress, userAgent, source string) error {
	view := &entity.ResourceView{
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Source:     source,
	}
	return r.db.Create(view).Error
}

// GetTotalViews 获取获取资源（访问）总次数。009-statistics-enhancement FR-001
func (r *ResourceViewRepositoryImpl) GetTotalViews() (int64, error) {
	var count int64
	err := r.db.Model(&entity.ResourceView{}).Count(&count).Error
	return count, err
}

// GetPanDistribution 获取访问（获取资源）的网盘分布：resource_views join resources+pan，按网盘聚合。009 FR-004
func (r *ResourceViewRepositoryImpl) GetPanDistribution() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("resource_views").
		Select("COALESCE(NULLIF(pan.remark, ''), pan.name, '未知') as pan, COUNT(*) as count").
		Joins("LEFT JOIN resources ON resources.id = resource_views.resource_id").
		Joins("LEFT JOIN pan ON pan.id = resources.pan_id").
		Group("pan.remark, pan.name").
		Order("count DESC").
		Scan(&results).Error
	return results, err
}

// GetSourceDistribution 获取访问（获取资源）的来源渠道分布。009-statistics-enhancement FR-025
func (r *ResourceViewRepositoryImpl) GetSourceDistribution() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("resource_views").
		Select("source, COUNT(*) as count").
		Group("source").
		Order("count DESC").
		Scan(&results).Error
	return results, err
}

// GetTodayViews 获取今日访问量
func (r *ResourceViewRepositoryImpl) GetTodayViews() (int64, error) {
	today := utils.GetTodayString()
	var count int64
	err := r.db.Model(&entity.ResourceView{}).
		Where("DATE(created_at) = ?", today).
		Count(&count).Error
	return count, err
}

// GetViewsByDate 获取指定日期的访问量
func (r *ResourceViewRepositoryImpl) GetViewsByDate(date string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.ResourceView{}).
		Where("DATE(created_at) = ?", date).
		Count(&count).Error
	return count, err
}

// GetViewsTrend 获取访问量趋势数据
func (r *ResourceViewRepositoryImpl) GetViewsTrend(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	for i := days - 1; i >= 0; i-- {
		date := utils.GetCurrentTime().AddDate(0, 0, -i)
		dateStr := date.Format(utils.TimeFormatDate)

		count, err := r.GetViewsByDate(dateStr)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"date":  dateStr,
			"views": count,
		})
	}

	return results, nil
}

// GetResourceViews 获取指定资源的访问记录
func (r *ResourceViewRepositoryImpl) GetResourceViews(resourceID uint, limit int) ([]entity.ResourceView, error) {
	var views []entity.ResourceView
	err := r.db.Where("resource_id = ?", resourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&views).Error
	return views, err
}
