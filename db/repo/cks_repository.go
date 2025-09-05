package repo

import (
	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// CksRepository Cks的Repository接口
type CksRepository interface {
	BaseRepository[entity.Cks]
	FindByPanID(panID uint) ([]entity.Cks, error)
	FindByIds(ids []uint) ([]*entity.Cks, error)
	FindByIsValid(isValid bool) ([]entity.Cks, error)
	UpdateSpace(id uint, space, leftSpace int64) error
	DeleteByPanID(panID uint) error
	UpdateWithAllFields(cks *entity.Cks) error
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

// FindAll 查找所有Cks，预加载Pan关联数据
func (r *CksRepositoryImpl) FindAll() ([]entity.Cks, error) {
	var cks []entity.Cks
	err := r.db.Preload("Pan").Find(&cks).Error
	return cks, err
}

// FindByID 根据ID查找Cks，预加载Pan关联数据
func (r *CksRepositoryImpl) FindByID(id uint) (*entity.Cks, error) {
	var cks entity.Cks
	err := r.db.Preload("Pan").First(&cks, id).Error
	if err != nil {
		return nil, err
	}
	return &cks, nil
}

func (r *CksRepositoryImpl) FindByIds(ids []uint) ([]*entity.Cks, error) {
	var cks []*entity.Cks
	err := r.db.Preload("Pan").Where("id IN ?", ids).Find(&cks).Error
	if err != nil {
		return nil, err
	}
	return cks, nil
}

// UpdateWithAllFields 更新Cks，包括零值字段
func (r *CksRepositoryImpl) UpdateWithAllFields(cks *entity.Cks) error {
	return r.db.Save(cks).Error
}
