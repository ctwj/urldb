package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// ResourceRepository Resource的Repository接口
type ResourceRepository interface {
	BaseRepository[entity.Resource]
	FindWithRelations() ([]entity.Resource, error)
	FindByCategoryID(categoryID uint) ([]entity.Resource, error)
	FindByPanID(panID uint) ([]entity.Resource, error)
	FindByIsValid(isValid bool) ([]entity.Resource, error)
	FindByIsPublic(isPublic bool) ([]entity.Resource, error)
	Search(query string, categoryID *uint, page, limit int) ([]entity.Resource, int64, error)
	IncrementViewCount(id uint) error
	FindWithTags() ([]entity.Resource, error)
	UpdateWithTags(resource *entity.Resource, tagIDs []uint) error
}

// ResourceRepositoryImpl Resource的Repository实现
type ResourceRepositoryImpl struct {
	BaseRepositoryImpl[entity.Resource]
}

// NewResourceRepository 创建Resource Repository
func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &ResourceRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Resource]{db: db},
	}
}

// FindWithRelations 查找包含关联关系的资源
func (r *ResourceRepositoryImpl) FindWithRelations() ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Preload("Category").Preload("Pan").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindByCategoryID 根据分类ID查找
func (r *ResourceRepositoryImpl) FindByCategoryID(categoryID uint) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("category_id = ?", categoryID).Preload("Category").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindByPanID 根据平台ID查找
func (r *ResourceRepositoryImpl) FindByPanID(panID uint) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("pan_id = ?", panID).Preload("Category").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindByIsValid 根据有效性查找
func (r *ResourceRepositoryImpl) FindByIsValid(isValid bool) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("is_valid = ?", isValid).Preload("Category").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindByIsPublic 根据公开性查找
func (r *ResourceRepositoryImpl) FindByIsPublic(isPublic bool) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("is_public = ?", isPublic).Preload("Category").Preload("Tags").Find(&resources).Error
	return resources, err
}

// Search 搜索资源
func (r *ResourceRepositoryImpl) Search(query string, categoryID *uint, page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&entity.Resource{}).Preload("Category").Preload("Tags")

	// 构建查询条件
	if query != "" {
		db = db.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if categoryID != nil {
		db = db.Where("category_id = ?", *categoryID)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// IncrementViewCount 增加浏览次数
func (r *ResourceRepositoryImpl) IncrementViewCount(id uint) error {
	return r.db.Model(&entity.Resource{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// FindWithTags 查找包含标签的资源
func (r *ResourceRepositoryImpl) FindWithTags() ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Preload("Tags").Find(&resources).Error
	return resources, err
}

// UpdateWithTags 更新资源及其标签
func (r *ResourceRepositoryImpl) UpdateWithTags(resource *entity.Resource, tagIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新资源
		if err := tx.Save(resource).Error; err != nil {
			return err
		}

		// 删除旧的标签关联
		if err := tx.Where("resource_id = ?", resource.ID).Delete(&entity.ResourceTag{}).Error; err != nil {
			return err
		}

		// 创建新的标签关联
		for _, tagID := range tagIDs {
			resourceTag := &entity.ResourceTag{
				ResourceID: resource.ID,
				TagID:      tagID,
			}
			if err := tx.Create(resourceTag).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
