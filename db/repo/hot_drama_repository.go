package repo

import (
	"github.com/ctwj/panResManage/db/entity"

	"gorm.io/gorm"
)

// HotDramaRepository 热播剧仓储接口
type HotDramaRepository interface {
	Create(drama *entity.HotDrama) error
	FindByID(id uint) (*entity.HotDrama, error)
	FindAll(page, pageSize int) ([]entity.HotDrama, int64, error)
	FindByCategory(category string, page, pageSize int) ([]entity.HotDrama, int64, error)
	FindByDoubanID(doubanID string) (*entity.HotDrama, error)
	Upsert(drama *entity.HotDrama) error
	Delete(id uint) error
	DeleteByDoubanID(doubanID string) error
	DeleteOldRecords(days int) error
	DeleteAll() error
	BatchCreate(dramas []*entity.HotDrama) error
}

// hotDramaRepository 热播剧仓储实现
type hotDramaRepository struct {
	db *gorm.DB
}

// NewHotDramaRepository 创建热播剧仓储实例
func NewHotDramaRepository(db *gorm.DB) HotDramaRepository {
	return &hotDramaRepository{db: db}
}

// Create 创建热播剧记录
func (r *hotDramaRepository) Create(drama *entity.HotDrama) error {
	return r.db.Create(drama).Error
}

// FindByID 根据ID查找热播剧
func (r *hotDramaRepository) FindByID(id uint) (*entity.HotDrama, error) {
	var drama entity.HotDrama
	err := r.db.First(&drama, id).Error
	if err != nil {
		return nil, err
	}
	return &drama, nil
}

// FindAll 查找所有热播剧（分页）
func (r *hotDramaRepository) FindAll(page, pageSize int) ([]entity.HotDrama, int64, error) {
	var dramas []entity.HotDrama
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&entity.HotDrama{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&dramas).Error
	if err != nil {
		return nil, 0, err
	}

	return dramas, total, nil
}

// FindByCategory 根据分类查找热播剧（分页）
func (r *hotDramaRepository) FindByCategory(category string, page, pageSize int) ([]entity.HotDrama, int64, error) {
	var dramas []entity.HotDrama
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&entity.HotDrama{}).Where("category = ?", category).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("category = ?", category).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&dramas).Error
	if err != nil {
		return nil, 0, err
	}

	return dramas, total, nil
}

// FindByDoubanID 根据豆瓣ID查找热播剧
func (r *hotDramaRepository) FindByDoubanID(doubanID string) (*entity.HotDrama, error) {
	var drama entity.HotDrama
	err := r.db.Where("douban_id = ?", doubanID).First(&drama).Error
	if err != nil {
		return nil, err
	}
	return &drama, nil
}

// Upsert 插入或更新热播剧记录
func (r *hotDramaRepository) Upsert(drama *entity.HotDrama) error {
	if drama.DoubanID != "" {
		// 如果存在豆瓣ID，先尝试查找现有记录
		existing, err := r.FindByDoubanID(drama.DoubanID)
		if err == nil && existing != nil {
			// 更新现有记录
			drama.ID = existing.ID
			return r.db.Save(drama).Error
		}
	}

	// 创建新记录
	return r.Create(drama)
}

// Delete 删除热播剧记录
func (r *hotDramaRepository) Delete(id uint) error {
	return r.db.Delete(&entity.HotDrama{}, id).Error
}

// DeleteByDoubanID 根据豆瓣ID删除热播剧记录
func (r *hotDramaRepository) DeleteByDoubanID(doubanID string) error {
	return r.db.Where("douban_id = ?", doubanID).Delete(&entity.HotDrama{}).Error
}

// DeleteOldRecords 删除指定天数前的旧记录
func (r *hotDramaRepository) DeleteOldRecords(days int) error {
	return r.db.Where("created_at < NOW() - INTERVAL '? days'", days).Delete(&entity.HotDrama{}).Error
}

// DeleteAll 删除所有热播剧记录
func (r *hotDramaRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.HotDrama{}).Error
}

// BatchCreate 批量创建热播剧记录
func (r *hotDramaRepository) BatchCreate(dramas []*entity.HotDrama) error {
	if len(dramas) == 0 {
		return nil
	}
	return r.db.CreateInBatches(dramas, 100).Error
}
