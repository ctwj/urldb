package repo

import (
	"fmt"

	"time"

	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// ResourceRepository Resource的Repository接口
type ResourceRepository interface {
	BaseRepository[entity.Resource]
	FindWithRelations() ([]entity.Resource, error)
	FindWithRelationsPaginated(page, limit int) ([]entity.Resource, int64, error)
	FindByCategoryID(categoryID uint) ([]entity.Resource, error)
	FindByCategoryIDPaginated(categoryID uint, page, limit int) ([]entity.Resource, int64, error)
	FindByPanID(panID uint) ([]entity.Resource, error)
	FindByPanIDPaginated(panID uint, page, limit int) ([]entity.Resource, int64, error)
	FindByIsValid(isValid bool) ([]entity.Resource, error)
	FindByIsPublic(isPublic bool) ([]entity.Resource, error)
	Search(query string, categoryID *uint, page, limit int) ([]entity.Resource, int64, error)
	SearchByPanID(query string, panID uint, page, limit int) ([]entity.Resource, int64, error)
	SearchWithFilters(params map[string]interface{}) ([]entity.Resource, int64, error)
	IncrementViewCount(id uint) error
	FindWithTags() ([]entity.Resource, error)
	UpdateWithTags(resource *entity.Resource, tagIDs []uint) error
	GetLatestResources(limit int) ([]entity.Resource, error)
	GetCachedLatestResources(limit int) ([]entity.Resource, error)
	InvalidateCache() error
	FindExists(url string, excludeID ...uint) (bool, error)
}

// ResourceRepositoryImpl Resource的Repository实现
type ResourceRepositoryImpl struct {
	BaseRepositoryImpl[entity.Resource]
	cache map[string]interface{}
}

// NewResourceRepository 创建Resource Repository
func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &ResourceRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Resource]{db: db},
		cache:              make(map[string]interface{}),
	}
}

// FindWithRelations 查找包含关联关系的资源
func (r *ResourceRepositoryImpl) FindWithRelations() ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Preload("Category").Preload("Pan").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindWithRelationsPaginated 分页查找包含关联关系的资源
func (r *ResourceRepositoryImpl) FindWithRelationsPaginated(page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit

	// 优化查询：只预加载必要的关联，并添加排序
	db := r.db.Model(&entity.Resource{}).
		Preload("Category").
		Preload("Pan").
		Order("updated_at DESC") // 按更新时间倒序，显示最新内容

	// 获取总数（使用缓存键）
	cacheKey := fmt.Sprintf("resources_total_%d_%d", page, limit)
	if cached, exists := r.cache[cacheKey]; exists {
		if totalCached, ok := cached.(int64); ok {
			total = totalCached
		}
	} else {
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		// 缓存总数（5分钟）
		r.cache[cacheKey] = total
		go func() {
			time.Sleep(5 * time.Minute)
			delete(r.cache, cacheKey)
		}()
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// FindByCategoryID 根据分类ID查找
func (r *ResourceRepositoryImpl) FindByCategoryID(categoryID uint) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("category_id = ?", categoryID).Preload("Category").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindByCategoryIDPaginated 分页根据分类ID查找
func (r *ResourceRepositoryImpl) FindByCategoryIDPaginated(categoryID uint, page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&entity.Resource{}).Where("category_id = ?", categoryID).Preload("Category").Preload("Tags").Order("updated_at DESC")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// FindByPanID 根据平台ID查找
func (r *ResourceRepositoryImpl) FindByPanID(panID uint) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("pan_id = ?", panID).Preload("Category").Preload("Tags").Find(&resources).Error
	return resources, err
}

// FindByPanIDPaginated 分页根据平台ID查找
func (r *ResourceRepositoryImpl) FindByPanIDPaginated(panID uint, page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&entity.Resource{}).Where("pan_id = ?", panID).Preload("Category").Preload("Tags").Order("updated_at DESC")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
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

	// 获取分页数据，按更新时间倒序
	err := db.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// SearchByPanID 在指定平台内搜索资源
func (r *ResourceRepositoryImpl) SearchByPanID(query string, panID uint, page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&entity.Resource{}).Preload("Category").Preload("Tags").Where("pan_id = ?", panID)

	// 构建查询条件
	if query != "" {
		db = db.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按更新时间倒序
	err := db.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// SearchWithFilters 根据参数进行搜索
func (r *ResourceRepositoryImpl) SearchWithFilters(params map[string]interface{}) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	db := r.db.Model(&entity.Resource{})

	// 处理参数
	for key, value := range params {
		switch key {
		case "query":
			if query, ok := value.(string); ok && query != "" {
				db = db.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
			}
		case "category_id":
			if categoryID, ok := value.(uint); ok {
				db = db.Where("category_id = ?", categoryID)
			}
		case "is_valid":
			if isValid, ok := value.(bool); ok {
				db = db.Where("is_valid = ?", isValid)
			}
		case "is_public":
			if isPublic, ok := value.(bool); ok {
				db = db.Where("is_public = ?", isPublic)
			}
		case "pan_id":
			if panID, ok := value.(uint); ok {
				db = db.Where("pan_id = ?", panID)
			}
		}
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按更新时间倒序
	err := db.Order("updated_at DESC").Find(&resources).Error
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

// GetLatestResources 获取最新资源
func (r *ResourceRepositoryImpl) GetLatestResources(limit int) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Order("created_at DESC").Limit(limit).Find(&resources).Error
	return resources, err
}

// GetCachedLatestResources 获取缓存的最新资源
func (r *ResourceRepositoryImpl) GetCachedLatestResources(limit int) ([]entity.Resource, error) {
	cacheKey := fmt.Sprintf("latest_resources_%d", limit)

	// 检查缓存
	if cached, exists := r.cache[cacheKey]; exists {
		if resources, ok := cached.([]entity.Resource); ok {
			return resources, nil
		}
	}

	// 从数据库获取
	resources, err := r.GetLatestResources(limit)
	if err != nil {
		return nil, err
	}

	// 缓存结果（5分钟过期）
	r.cache[cacheKey] = resources
	go func() {
		time.Sleep(5 * time.Minute)
		delete(r.cache, cacheKey)
	}()

	return resources, nil
}

// InvalidateCache 清除缓存
func (r *ResourceRepositoryImpl) InvalidateCache() error {
	r.cache = make(map[string]interface{})
	return nil
}

// FindExists 检查是否存在相同URL的资源
func (r *ResourceRepositoryImpl) FindExists(url string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&entity.Resource{}).Where("url = ?", url)

	// 如果有排除ID，则排除该记录（用于更新时排除自己）
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
