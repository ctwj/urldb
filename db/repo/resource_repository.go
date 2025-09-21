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
	BatchFindByURLs(urls []string) ([]entity.Resource, error)
	GetResourcesForTransfer(panID uint, sinceTime time.Time, limit int) ([]*entity.Resource, error)
	GetByURL(url string) (*entity.Resource, error)
	UpdateSaveURL(id uint, saveURL string) error
	CreateResourceTag(resourceTag *entity.ResourceTag) error
	FindByIDs(ids []uint) ([]entity.Resource, error)
	FindUnsyncedToMeilisearch(page, limit int) ([]entity.Resource, int64, error)
	FindSyncedToMeilisearch(page, limit int) ([]entity.Resource, int64, error)
	CountUnsyncedToMeilisearch() (int64, error)
	CountSyncedToMeilisearch() (int64, error)
	MarkAsSyncedToMeilisearch(ids []uint) error
	MarkAllAsUnsyncedToMeilisearch() error
	FindAllWithPagination(page, limit int) ([]entity.Resource, int64, error)
	GetRandomResourceWithFilters(categoryFilter, tagFilter string, isPushSavedInfo bool) (*entity.Resource, error)
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

	db := r.db.Model(&entity.Resource{}).Preload("Category").Preload("Pan").Preload("Tags")

	// 处理参数
	for key, value := range params {
		switch key {
		case "search": // 添加search参数支持
			if query, ok := value.(string); ok && query != "" {
				db = db.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
			}
		case "category_id": // 添加category_id参数支持
			if categoryID, ok := value.(uint); ok {
				fmt.Printf("应用分类筛选: category_id = %d\n", categoryID)
				db = db.Where("category_id = ?", categoryID)
			} else {
				fmt.Printf("分类ID类型错误: %T, value: %v\n", value, value)
			}
		case "category": // 添加category参数支持（字符串形式）
			if category, ok := value.(string); ok && category != "" {
				// 根据分类名称查找分类ID
				var categoryEntity entity.Category
				if err := r.db.Where("name ILIKE ?", "%"+category+"%").First(&categoryEntity).Error; err == nil {
					db = db.Where("category_id = ?", categoryEntity.ID)
				}
			}
		case "tag": // 添加tag参数支持
			if tag, ok := value.(string); ok && tag != "" {
				// 根据标签名称查找相关资源
				var tagEntity entity.Tag
				if err := r.db.Where("name ILIKE ?", "%"+tag+"%").First(&tagEntity).Error; err == nil {
					// 通过中间表查找包含该标签的资源
					db = db.Joins("JOIN resource_tags ON resources.id = resource_tags.resource_id").
						Where("resource_tags.tag_id = ?", tagEntity.ID)
				}
			}
		case "pan_id": // 添加pan_id参数支持
			if panID, ok := value.(uint); ok {
				db = db.Where("pan_id = ?", panID)
			}
		case "is_valid":
			if isValid, ok := value.(bool); ok {
				db = db.Where("is_valid = ?", isValid)
			}
		case "is_public":
			if isPublic, ok := value.(bool); ok {
				db = db.Where("is_public = ?", isPublic)
			}
		case "has_save_url": // 添加has_save_url参数支持
			if hasSaveURL, ok := value.(bool); ok {
				fmt.Printf("处理 has_save_url 参数: %v\n", hasSaveURL)
				if hasSaveURL {
					// 有转存链接：save_url不为空且不为空格
					db = db.Where("save_url IS NOT NULL AND save_url != '' AND TRIM(save_url) != ''")
					fmt.Printf("应用 has_save_url=true 条件: save_url IS NOT NULL AND save_url != '' AND TRIM(save_url) != ''\n")
				} else {
					// 没有转存链接：save_url为空、NULL或只有空格
					db = db.Where("(save_url IS NULL OR save_url = '' OR TRIM(save_url) = '')")
					fmt.Printf("应用 has_save_url=false 条件: (save_url IS NULL OR save_url = '' OR TRIM(save_url) = '')\n")
				}
			}
		case "no_save_url": // 添加no_save_url参数支持（与has_save_url=false相同）
			if noSaveURL, ok := value.(bool); ok && noSaveURL {
				db = db.Where("(save_url IS NULL OR save_url = '' OR TRIM(save_url) = '')")
			}
		case "pan_name": // 添加pan_name参数支持
			if panName, ok := value.(string); ok && panName != "" {
				// 根据平台名称查找平台ID
				var panEntity entity.Pan
				if err := r.db.Where("name ILIKE ?", "%"+panName+"%").First(&panEntity).Error; err == nil {
					db = db.Where("pan_id = ?", panEntity.ID)
				}
			}
		}
	}

	// 管理后台显示所有资源，公开API才限制为有效的公开资源
	// 这里通过检查请求来源来判断是否为管理后台
	// 如果没有明确指定is_valid和is_public，则显示所有资源
	// 注意：这个逻辑可能需要根据实际需求调整
	if _, hasIsValid := params["is_valid"]; !hasIsValid {
		// 管理后台不限制is_valid
		// db = db.Where("is_valid = ?", true)
	}
	if _, hasIsPublic := params["is_public"]; !hasIsPublic {
		// 管理后台不限制is_public
		// db = db.Where("is_public = ?", true)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 处理分页参数
	page := 1
	pageSize := 20

	if pageVal, ok := params["page"].(int); ok && pageVal > 0 {
		page = pageVal
	}
	if pageSizeVal, ok := params["page_size"].(int); ok && pageSizeVal > 0 {
		pageSize = pageSizeVal
		fmt.Printf("原始pageSize: %d\n", pageSize)
		// 限制最大page_size为10000（管理后台需要更大的数据量）
		if pageSize > 10000 {
			pageSize = 10000
			fmt.Printf("pageSize超过10000，限制为: %d\n", pageSize)
		}
		fmt.Printf("最终pageSize: %d\n", pageSize)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据，按更新时间倒序
	err := db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&resources).Error
	fmt.Printf("查询结果: 总数=%d, 当前页数据量=%d, pageSize=%d\n", total, len(resources), pageSize)
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
	query := r.db.Model(&entity.Resource{}).Where("url = ? OR save_url = ?", url, url)

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

func (r *ResourceRepositoryImpl) BatchFindByURLs(urls []string) ([]entity.Resource, error) {
	var resources []entity.Resource
	if len(urls) == 0 {
		return resources, nil
	}
	err := r.db.Where("url IN ?", urls).Find(&resources).Error
	return resources, err
}

// GetResourcesForTransfer 获取需要转存的资源
func (r *ResourceRepositoryImpl) GetResourcesForTransfer(panID uint, sinceTime time.Time, limit int) ([]*entity.Resource, error) {
	var resources []*entity.Resource
	query := r.db.Where("pan_id = ? AND (save_url = '' OR save_url IS NULL) AND (error_msg = '' OR error_msg IS NULL)", panID)
	if !sinceTime.IsZero() {
		query = query.Where("created_at >= ?", sinceTime)
	}

	// 添加数量限制
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("created_at DESC").Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// GetByURL 根据URL获取资源
func (r *ResourceRepositoryImpl) GetByURL(url string) (*entity.Resource, error) {
	var resource entity.Resource
	err := r.db.Where("url = ?", url).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// FindByIDs 根据ID列表查找资源
func (r *ResourceRepositoryImpl) FindByIDs(ids []uint) ([]entity.Resource, error) {
	if len(ids) == 0 {
		return []entity.Resource{}, nil
	}

	var resources []entity.Resource
	err := r.db.Where("id IN ?", ids).Preload("Category").Preload("Pan").Preload("Tags").Find(&resources).Error
	return resources, err
}

// UpdateSaveURL 更新保存URL
func (r *ResourceRepositoryImpl) UpdateSaveURL(id uint, saveURL string) error {
	return r.db.Model(&entity.Resource{}).Where("id = ?", id).Update("save_url", saveURL).Error
}

// CreateResourceTag 创建资源与标签的关联
func (r *ResourceRepositoryImpl) CreateResourceTag(resourceTag *entity.ResourceTag) error {
	return r.db.Create(resourceTag).Error
}

// FindUnsyncedToMeilisearch 查找未同步到Meilisearch的资源
func (r *ResourceRepositoryImpl) FindUnsyncedToMeilisearch(page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit

	// 查询未同步的资源
	db := r.db.Model(&entity.Resource{}).
		Where("synced_to_meilisearch = ?", false).
		Preload("Category").
		Preload("Pan").
		Preload("Tags"). // 添加Tags预加载
		Order("updated_at DESC")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// CountUnsyncedToMeilisearch 统计未同步到Meilisearch的资源数量
func (r *ResourceRepositoryImpl) CountUnsyncedToMeilisearch() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).
		Where("synced_to_meilisearch = ?", false).
		Count(&count).Error
	return count, err
}

// MarkAsSyncedToMeilisearch 标记资源为已同步到Meilisearch
func (r *ResourceRepositoryImpl) MarkAsSyncedToMeilisearch(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	now := time.Now()
	return r.db.Model(&entity.Resource{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"synced_to_meilisearch": true,
			"synced_at":             now,
		}).Error
}

// MarkAllAsUnsyncedToMeilisearch 标记所有资源为未同步到Meilisearch
func (r *ResourceRepositoryImpl) MarkAllAsUnsyncedToMeilisearch() error {
	return r.db.Model(&entity.Resource{}).
		Where("1 = 1"). // 添加WHERE条件以更新所有记录
		Updates(map[string]interface{}{
			"synced_to_meilisearch": false,
			"synced_at":             nil,
		}).Error
}

// FindSyncedToMeilisearch 查找已同步到Meilisearch的资源
func (r *ResourceRepositoryImpl) FindSyncedToMeilisearch(page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit

	// 查询已同步的资源
	db := r.db.Model(&entity.Resource{}).
		Where("synced_to_meilisearch = ?", true).
		Preload("Category").
		Preload("Pan").
		Preload("Tags").
		Order("updated_at DESC")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// CountSyncedToMeilisearch 统计已同步到Meilisearch的资源数量
func (r *ResourceRepositoryImpl) CountSyncedToMeilisearch() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).
		Where("synced_to_meilisearch = ?", true).
		Count(&count).Error
	return count, err
}

// FindAllWithPagination 分页查找所有资源
func (r *ResourceRepositoryImpl) FindAllWithPagination(page, limit int) ([]entity.Resource, int64, error) {
	var resources []entity.Resource
	var total int64

	offset := (page - 1) * limit

	// 查询所有资源
	db := r.db.Model(&entity.Resource{}).
		Preload("Category").
		Preload("Pan").
		Preload("Tags").
		Order("updated_at DESC")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Find(&resources).Error
	return resources, total, err
}

// GetRandomResourceWithFilters 使用 PostgreSQL RANDOM() 功能随机获取一个符合条件的资源
func (r *ResourceRepositoryImpl) GetRandomResourceWithFilters(categoryFilter, tagFilter string, isPushSavedInfo bool) (*entity.Resource, error) {
	// 构建查询条件
	query := r.db.Model(&entity.Resource{}).Preload("Category").Preload("Pan").Preload("Tags")

	// 基础条件：有效且公开的资源
	query = query.Where("is_valid = ? AND is_public = ?", true, true)

	// 根据分类过滤
	if categoryFilter != "" {
		// 查找分类ID
		var categoryEntity entity.Category
		if err := r.db.Where("name ILIKE ?", "%"+categoryFilter+"%").First(&categoryEntity).Error; err == nil {
			query = query.Where("category_id = ?", categoryEntity.ID)
		}
	}

	// 根据标签过滤
	if tagFilter != "" {
		// 查找标签ID
		var tagEntity entity.Tag
		if err := r.db.Where("name ILIKE ?", "%"+tagFilter+"%").First(&tagEntity).Error; err == nil {
			// 通过中间表查找包含该标签的资源
			query = query.Joins("JOIN resource_tags ON resources.id = resource_tags.resource_id").
				Where("resource_tags.tag_id = ?", tagEntity.ID)
		}
	}

	// // 根据是否只推送已转存资源过滤
	// if isPushSavedInfo {
	// 	query = query.Where("save_url IS NOT NULL AND save_url != '' AND TRIM(save_url) != ''")
	// }

	// 使用 PostgreSQL 的 RANDOM() 进行随机排序，并限制为1个结果
	var resource entity.Resource
	err := query.Order("RANDOM()").Limit(1).First(&resource).Error

	if err != nil {
		return nil, err
	}

	return &resource, nil
}
