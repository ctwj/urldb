package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// PanRepository Pan的Repository接口
type PanRepository interface {
	BaseRepository[entity.Pan]
	FindByIsValid(isValid bool) ([]entity.Pan, error)
	FindWithCks() ([]entity.Pan, error)
	UpdateSpace(id uint, space, leftSpace int64) error
}

// PanRepositoryImpl Pan的Repository实现
type PanRepositoryImpl struct {
	BaseRepositoryImpl[entity.Pan]
}

// NewPanRepository 创建Pan Repository
func NewPanRepository(db *gorm.DB) PanRepository {
	return &PanRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Pan]{db: db},
	}
}

// FindByIsValid 根据有效性查找
func (r *PanRepositoryImpl) FindByIsValid(isValid bool) ([]entity.Pan, error) {
	var pans []entity.Pan
	err := r.db.Where("is_valid = ?", isValid).Find(&pans).Error
	return pans, err
}

// FindWithCks 查找包含Cks的Pan
func (r *PanRepositoryImpl) FindWithCks() ([]entity.Pan, error) {
	var pans []entity.Pan
	err := r.db.Preload("Cks").Find(&pans).Error
	return pans, err
}

// UpdateSpace 更新空间信息
func (r *PanRepositoryImpl) UpdateSpace(id uint, space, leftSpace int64) error {
	return r.db.Model(&entity.Pan{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"space":      space,
			"left_space": leftSpace,
		}).Error
}
