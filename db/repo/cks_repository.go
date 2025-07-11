package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// CksRepository Cks的Repository接口
type CksRepository interface {
	BaseRepository[entity.Cks]
	FindByPanID(panID uint) ([]entity.Cks, error)
	FindByIsValid(isValid bool) ([]entity.Cks, error)
	UpdateSpace(id uint, space, leftSpace int64) error
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

// FindByIsValid 根据有效性查找
func (r *CksRepositoryImpl) FindByIsValid(isValid bool) ([]entity.Cks, error) {
	var cks []entity.Cks
	err := r.db.Where("is_valid = ?", isValid).Find(&cks).Error
	return cks, err
}

// UpdateSpace 更新空间信息
func (r *CksRepositoryImpl) UpdateSpace(id uint, space, leftSpace int64) error {
	return r.db.Model(&entity.Cks{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"space":      space,
			"left_space": leftSpace,
		}).Error
}

// DeleteByPanID 根据PanID删除
func (r *CksRepositoryImpl) DeleteByPanID(panID uint) error {
	return r.db.Where("pan_id = ?", panID).Delete(&entity.Cks{}).Error
}
