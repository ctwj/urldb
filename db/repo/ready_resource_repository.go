package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// ReadyResourceRepository ReadyResource的Repository接口
type ReadyResourceRepository interface {
	BaseRepository[entity.ReadyResource]
	FindByURL(url string) (*entity.ReadyResource, error)
	FindByIP(ip string) ([]entity.ReadyResource, error)
	BatchCreate(resources []entity.ReadyResource) error
	DeleteByURL(url string) error
	FindAllWithinDays(days int) ([]entity.ReadyResource, error)
}

// ReadyResourceRepositoryImpl ReadyResource的Repository实现
type ReadyResourceRepositoryImpl struct {
	BaseRepositoryImpl[entity.ReadyResource]
}

// NewReadyResourceRepository 创建ReadyResource Repository
func NewReadyResourceRepository(db *gorm.DB) ReadyResourceRepository {
	return &ReadyResourceRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.ReadyResource]{db: db},
	}
}

// FindByURL 根据URL查找
func (r *ReadyResourceRepositoryImpl) FindByURL(url string) (*entity.ReadyResource, error) {
	var resource entity.ReadyResource
	err := r.db.Where("url = ?", url).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// FindByIP 根据IP查找
func (r *ReadyResourceRepositoryImpl) FindByIP(ip string) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	err := r.db.Where("ip = ?", ip).Find(&resources).Error
	return resources, err
}

// BatchCreate 批量创建
func (r *ReadyResourceRepositoryImpl) BatchCreate(resources []entity.ReadyResource) error {
	return r.db.Create(&resources).Error
}

// DeleteByURL 根据URL删除
func (r *ReadyResourceRepositoryImpl) DeleteByURL(url string) error {
	return r.db.Where("url = ?", url).Delete(&entity.ReadyResource{}).Error
}

// FindAllWithinDays 获取n天内的待处理资源，n=0时不限制
func (r *ReadyResourceRepositoryImpl) FindAllWithinDays(days int) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	db := r.db.Model(&entity.ReadyResource{})
	if days > 0 {
		since := time.Now().AddDate(0, 0, -days)
		db = db.Where("create_time >= ?", since)
	}
	err := db.Find(&resources).Error
	return resources, err
}
