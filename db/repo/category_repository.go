package repo

import (
	"github.com/ctwj/panResManage/db/entity"

	"gorm.io/gorm"
)

// CategoryRepository Category的Repository接口
type CategoryRepository interface {
	BaseRepository[entity.Category]
	FindByName(name string) (*entity.Category, error)
	FindWithResources() ([]entity.Category, error)
	FindWithTags() ([]entity.Category, error)
	GetResourceCount(categoryID uint) (int64, error)
	GetTagCount(categoryID uint) (int64, error)
	GetTagNames(categoryID uint) ([]string, error)
	FindWithPagination(page, pageSize int) ([]entity.Category, int64, error)
	Search(query string, page, pageSize int) ([]entity.Category, int64, error)
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

// FindWithTags 查找包含标签的分类
func (r *CategoryRepositoryImpl) FindWithTags() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Preload("Tags").Find(&categories).Error
	return categories, err
}

// GetResourceCount 获取分类下的资源数量
func (r *CategoryRepositoryImpl) GetResourceCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// GetTagCount 获取分类下的标签数量
func (r *CategoryRepositoryImpl) GetTagCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Tag{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// GetTagNames 获取分类下的标签名称列表
func (r *CategoryRepositoryImpl) GetTagNames(categoryID uint) ([]string, error) {
	var tags []entity.Tag
	err := r.db.Model(&entity.Tag{}).Where("category_id = ?", categoryID).Select("name").Find(&tags).Error
	if err != nil {
		return nil, err
	}

	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return names, nil
}

// FindWithPagination 分页查询分类
func (r *CategoryRepositoryImpl) FindWithPagination(page, pageSize int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	// 获取总数
	err := r.db.Model(&entity.Category{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// Search 搜索分类
func (r *CategoryRepositoryImpl) Search(query string, page, pageSize int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	// 构建搜索条件
	searchQuery := "%" + query + "%"

	// 获取总数
	err := r.db.Model(&entity.Category{}).Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页搜索
	offset := (page - 1) * pageSize
	err = r.db.Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
