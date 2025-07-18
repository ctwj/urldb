package repo

import (
	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// PanRepository Pan的Repository接口
type PanRepository interface {
	BaseRepository[entity.Pan]
	FindWithCks() ([]entity.Pan, error)
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

// FindWithCks 查找包含Cks的Pan
func (r *PanRepositoryImpl) FindWithCks() ([]entity.Pan, error) {
	var pans []entity.Pan
	err := r.db.Preload("Cks").Find(&pans).Error
	return pans, err
}
