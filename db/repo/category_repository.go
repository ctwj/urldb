package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// CategoryRepository Category的Repository接口
type CategoryRepository interface {
	BaseRepository[entity.Category]
	FindByName(name string) (*entity.Category, error)
	FindWithResources() ([]entity.Category, error)
	GetResourceCount(categoryID uint) (int64, error)
}

// CategoryRepositoryImpl Category的Repository实现
type CategoryRepositoryImpl struct {
	BaseRepositoryImpl[entity.Category]
}

// NewCategoryRepository 创建Category Repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Category]{db: db},
	}
}

// FindByName 根据名称查找
func (r *CategoryRepositoryImpl) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindWithResources 查找包含资源的分类
func (r *CategoryRepositoryImpl) FindWithResources() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Preload("Resources").Find(&categories).Error
	return categories, err
}

// GetResourceCount 获取分类下的资源数量
func (r *CategoryRepositoryImpl) GetResourceCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
