package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// TagRepository Tag的Repository接口
type TagRepository interface {
	BaseRepository[entity.Tag]
	FindByName(name string) (*entity.Tag, error)
	FindWithResources() ([]entity.Tag, error)
	FindByCategoryID(categoryID uint) ([]entity.Tag, error)
	FindByCategoryIDPaginated(categoryID uint, page, pageSize int) ([]entity.Tag, int64, error)
	GetResourceCount(tagID uint) (int64, error)
	FindByResourceID(resourceID uint) ([]entity.Tag, error)
	FindWithPagination(page, pageSize int) ([]entity.Tag, int64, error)
	Search(query string, page, pageSize int) ([]entity.Tag, int64, error)
	UpdateWithNulls(tag *entity.Tag) error
}

// TagRepositoryImpl Tag的Repository实现
type TagRepositoryImpl struct {
	BaseRepositoryImpl[entity.Tag]
}

// NewTagRepository 创建Tag Repository
func NewTagRepository(db *gorm.DB) TagRepository {
	return &TagRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Tag]{db: db},
	}
}

// FindByName 根据名称查找
func (r *TagRepositoryImpl) FindByName(name string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindWithResources 查找包含资源的标签
func (r *TagRepositoryImpl) FindWithResources() ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Preload("Resources").Find(&tags).Error
	return tags, err
}

// FindByCategoryID 根据分类ID查找标签
func (r *TagRepositoryImpl) FindByCategoryID(categoryID uint) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Where("category_id = ?", categoryID).Preload("Category").Find(&tags).Error
	return tags, err
}

// FindByCategoryIDPaginated 分页根据分类ID查找标签
func (r *TagRepositoryImpl) FindByCategoryIDPaginated(categoryID uint, page, pageSize int) ([]entity.Tag, int64, error) {
	var tags []entity.Tag
	var total int64

	// 获取总数
	err := r.db.Model(&entity.Tag{}).Where("category_id = ?", categoryID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Where("category_id = ?", categoryID).Preload("Category").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// GetResourceCount 获取标签下的资源数量
func (r *TagRepositoryImpl) GetResourceCount(tagID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.ResourceTag{}).Where("tag_id = ?", tagID).Count(&count).Error
	return count, err
}

// FindByResourceID 根据资源ID查找标签
func (r *TagRepositoryImpl) FindByResourceID(resourceID uint) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Joins("JOIN resource_tags ON tags.id = resource_tags.tag_id").
		Where("resource_tags.resource_id = ?", resourceID).Find(&tags).Error
	return tags, err
}

// FindWithPagination 分页查询标签
func (r *TagRepositoryImpl) FindWithPagination(page, pageSize int) ([]entity.Tag, int64, error) {
	var tags []entity.Tag
	var total int64

	// 获取总数
	err := r.db.Model(&entity.Tag{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Preload("Category").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// Search 搜索标签
func (r *TagRepositoryImpl) Search(query string, page, pageSize int) ([]entity.Tag, int64, error) {
	var tags []entity.Tag
	var total int64

	// 构建搜索条件
	searchQuery := "%" + query + "%"

	// 获取总数
	err := r.db.Model(&entity.Tag{}).Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页搜索
	offset := (page - 1) * pageSize
	err = r.db.Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).
		Preload("Category").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// UpdateWithNulls 更新标签，包括null值
func (r *TagRepositoryImpl) UpdateWithNulls(tag *entity.Tag) error {
	// 使用Select方法明确指定要更新的字段，包括null值
	return r.db.Model(tag).Select("name", "description", "category_id", "updated_at").Updates(tag).Error
}
