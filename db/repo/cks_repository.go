package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// CksRepository Cks的Repository接口
type CksRepository interface {
	BaseRepository[entity.Cks]
	FindByPanID(panID uint) ([]entity.Cks, error)
	FindByPanIDAndType(panID uint, ckType string) ([]entity.Cks, error)
	DeleteByPanID(panID uint) error
}

// CksRepositoryImpl Cks的Repository实现
type CksRepositoryImpl struct {
	BaseRepositoryImpl[entity.Cks]
}

// NewCksRepository 创建Cks Repository
func NewCksRepository(db *gorm.DB) CksRepository {
	return &CksRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Cks]{db: db},
	}
}

// FindByPanID 根据PanID查找
func (r *CksRepositoryImpl) FindByPanID(panID uint) ([]entity.Cks, error) {
	var cks []entity.Cks
	err := r.db.Where("pan_id = ?", panID).Find(&cks).Error
	return cks, err
}

// FindByPanIDAndType 根据PanID和类型查找
func (r *CksRepositoryImpl) FindByPanIDAndType(panID uint, ckType string) ([]entity.Cks, error) {
	var cks []entity.Cks
	err := r.db.Where("pan_id = ? AND t = ?", panID, ckType).Find(&cks).Error
	return cks, err
}

// DeleteByPanID 根据PanID删除
func (r *CksRepositoryImpl) DeleteByPanID(panID uint) error {
	return r.db.Where("pan_id = ?", panID).Delete(&entity.Cks{}).Error
}
