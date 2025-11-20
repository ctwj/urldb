package repo

import (
	"gorm.io/gorm"
	"github.com/ctwj/urldb/db/entity"
)

// ReportRepository 举报Repository接口
type ReportRepository interface {
	BaseRepository[entity.Report]
	GetByResourceKey(resourceKey string) ([]*entity.Report, error)
	List(status string, page, pageSize int) ([]*entity.Report, int64, error)
	UpdateStatus(id uint, status string, processedBy *uint, note string) error
	// 兼容原有方法名
	GetByID(id uint) (*entity.Report, error)
}

// ReportRepositoryImpl 举报Repository实现
type ReportRepositoryImpl struct {
	BaseRepositoryImpl[entity.Report]
}

// NewReportRepository 创建举报Repository
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Report]{db: db},
	}
}

// Create 创建举报
func (r *ReportRepositoryImpl) Create(report *entity.Report) error {
	return r.GetDB().Create(report).Error
}

// GetByID 根据ID获取举报
func (r *ReportRepositoryImpl) GetByID(id uint) (*entity.Report, error) {
	var report entity.Report
	err := r.GetDB().Where("id = ?", id).First(&report).Error
	return &report, err
}

// GetByResourceKey 获取某个资源的所有举报
func (r *ReportRepositoryImpl) GetByResourceKey(resourceKey string) ([]*entity.Report, error) {
	var reports []*entity.Report
	err := r.GetDB().Where("resource_key = ?", resourceKey).Find(&reports).Error
	return reports, err
}

// List 获取举报列表
func (r *ReportRepositoryImpl) List(status string, page, pageSize int) ([]*entity.Report, int64, error) {
	var reports []*entity.Report
	var total int64

	query := r.GetDB().Model(&entity.Report{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&reports).Error
	return reports, total, err
}

// Update 更新举报
func (r *ReportRepositoryImpl) Update(report *entity.Report) error {
	return r.GetDB().Save(report).Error
}

// UpdateStatus 更新举报状态
func (r *ReportRepositoryImpl) UpdateStatus(id uint, status string, processedBy *uint, note string) error {
	return r.GetDB().Model(&entity.Report{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"processed_at":  gorm.Expr("NOW()"),
		"processed_by":  processedBy,
		"note":          note,
	}).Error
}

// Delete 删除举报
func (r *ReportRepositoryImpl) Delete(id uint) error {
	return r.GetDB().Delete(&entity.Report{}, id).Error
}