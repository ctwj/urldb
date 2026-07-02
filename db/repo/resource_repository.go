package repo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
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
	UpdateIsValid(id uint, isValid bool) error
	CreateResourceTag(resourceTag *entity.ResourceTag) error
	FindByIDs(ids []uint) ([]entity.Resource, error)
	FindUnsyncedToMeilisearch(page, limit int) ([]entity.Resource, int64, error)
	FindSyncedToMeilisearch(page, limit int) ([]entity.Resource, int64, error)
	CountUnsyncedToMeilisearch() (int64, error)
	CountSyncedToMeilisearch() (int64, error)
	CountInvalid() (int64, error)
	CountInvalidByDate(date string) (int64, error)
	// FindFailedResourceMetaByKey 取某 key 下任一失效资源的基本信息（失效页面优化）。
	// 命中返回非 nil；该 key 下无失效资源时返回 (nil, nil)。
	FindFailedResourceMetaByKey(key string) (*FailedResourceMeta, error)
	MarkAsSyncedToMeilisearch(ids []uint) error
	MarkAllAsUnsyncedToMeilisearch() error
	FindAllWithPagination(page, limit int) ([]entity.Resource, int64, error)
	GetRandomResourceWithFilters(categoryFilter, tagFilter string, isPushSavedInfo bool) (*entity.Resource, error)
	DeleteRelatedResources(ckID uint) (int64, error)
	CountResourcesByCkID(ckID uint) (int64, error)
	FindByResourceKey(key string) ([]entity.Resource, error)
	FindByKey(key string) ([]entity.Resource, error)
	GetHotResources(limit int) ([]entity.Resource, error)
	GetTotalCount() (int64, error)
	GetAllValidResources() ([]entity.Resource, error)
	// 002-auto-cleanup-transfer 自动清理相关
	FindDueForCleanup(retentionDays int, limit int) ([]*entity.Resource, error)
	MarkCleaned(id uint, cleanedAt time.Time) error
	MarkCleanError(id uint, errMsg string, errAt time.Time) error
	// UpdateFields 按主键更新指定字段；用于重转等场景部分更新（避免 GORM Updates 跳过零值）
	UpdateFields(id uint, fields map[string]interface{}) error
	// GenerateUniqueKey 生成唯一的6位Base62资源Key（复用 BaseRepositoryImpl 实现）
	GenerateUniqueKey() (string, error)
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
	// 使用新的分页查询功能
	options := &PaginationOptions{
		Page:     page,
		PageSize: limit,
		OrderBy:  "updated_at",
		OrderDir: "desc",
		Preloads: []string{"Category", "Pan"},
	}

	result, err := PaginatedQuery[entity.Resource](r.db, options)
	if err != nil {
		return nil, 0, err
	}

	return result.Data, result.Total, nil
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

	// 管理后台Search方法不过滤is_valid，显示所有资源

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

	// 管理后台SearchByPanID方法不过滤is_valid，显示所有资源

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
	startTime := utils.GetCurrentTime()
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
		case "tag_ids": // 添加tag_ids参数支持（标签ID列表）
			if tagIdsStr, ok := value.(string); ok && tagIdsStr != "" {
				// 将逗号分隔的标签ID字符串转换为整数ID数组
				tagIdStrs := strings.Split(tagIdsStr, ",")
				var tagIds []uint
				for _, idStr := range tagIdStrs {
					idStr = strings.TrimSpace(idStr) // 去除空格
					if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
						tagIds = append(tagIds, uint(id))
					}
				}
				if len(tagIds) > 0 {
					// 通过中间表查找包含任一标签的资源
					db = db.Joins("JOIN resource_tags ON resources.id = resource_tags.resource_id").
						Where("resource_tags.tag_id IN ?", tagIds)
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
		case "exclude_ids": // 添加exclude_ids参数支持
			if excludeIDs, ok := value.([]uint); ok && len(excludeIDs) > 0 {
				// 限制排除ID的数量，避免SQL语句过长
				maxExcludeIDs := 5000 // 限制排除ID数量，避免SQL语句过长
				if len(excludeIDs) > maxExcludeIDs {
					// 只取最近的maxExcludeIDs个ID进行排除
					startIndex := len(excludeIDs) - maxExcludeIDs
					truncatedExcludeIDs := excludeIDs[startIndex:]
					db = db.Where("resources.id NOT IN ?", truncatedExcludeIDs)
					utils.Debug("SearchWithFilters: 排除ID数量过多，截取最近%d个ID", len(truncatedExcludeIDs))
				} else {
					db = db.Where("resources.id NOT IN ?", excludeIDs)
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

	// 处理排序参数
	orderBy := "updated_at"
	orderDir := "DESC"

	if orderByVal, ok := params["order_by"].(string); ok && orderByVal != "" {
		// 验证排序字段，防止SQL注入
		validOrderByFields := map[string]bool{
			"created_at":  true,
			"updated_at":  true,
			"view_count":  true,
			"title":       true,
			"id":          true,
		}
		if validOrderByFields[orderByVal] {
			orderBy = orderByVal
		}
	}

	if orderDirVal, ok := params["order_dir"].(string); ok && orderDirVal != "" {
		// 验证排序方向
		if orderDirVal == "ASC" || orderDirVal == "DESC" {
			orderDir = orderDirVal
		}
	}

	// 获取分页数据，应用排序
	queryStart := utils.GetCurrentTime()
	err := db.Order(fmt.Sprintf("%s %s", orderBy, orderDir)).Offset(offset).Limit(pageSize).Find(&resources).Error
	queryDuration := time.Since(queryStart)
	totalDuration := time.Since(startTime)
	utils.Debug("SearchWithFilters完成: 总数=%d, 当前页数据量=%d, 排序=%s %s, 查询耗时=%v, 总耗时=%v", total, len(resources), orderBy, orderDir, queryDuration, totalDuration)
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
	startTime := utils.GetCurrentTime()
	var resource entity.Resource
	err := r.db.Where("url = ?", url).First(&resource).Error
	queryDuration := time.Since(startTime)
	if err != nil {
		utils.Debug("GetByURL失败: URL=%s, 错误=%v, 查询耗时=%v", url, err, queryDuration)
		return nil, err
	}
	utils.Debug("GetByURL成功: URL=%s, 查询耗时=%v", url, queryDuration)
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

// UpdateIsValid 更新资源的 is_valid 字段。
// 必须使用列级 Update（而非 Updates(struct)），因为 GORM 的 Updates 会跳过 bool 零值（false），
// 导致 is_valid 从 true→false 的翻转无法落库。
func (r *ResourceRepositoryImpl) UpdateIsValid(id uint, isValid bool) error {
	return r.db.Model(&entity.Resource{}).Where("id = ?", id).Update("is_valid", isValid).Error
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

// CountInvalid 统计失效资源数量（is_valid=false）。009-statistics-enhancement FR-003
func (r *ResourceRepositoryImpl) CountInvalid() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).
		Where("is_valid = ?", false).
		Count(&count).Error
	return count, err
}

// CountInvalidByDate 统计指定日期发生失效的资源数（按 invalidated_at 的日期）。009 FR-007/008
func (r *ResourceRepositoryImpl) CountInvalidByDate(date string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).
		Where("DATE(invalidated_at) = ?", date).
		Count(&count).Error
	return count, err
}

// FailedResourceMeta 失效资源的基本信息（失效页面优化：标题/封面/key）
type FailedResourceMeta struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
	Key   string `json:"key"`
}

// FindFailedResourceMetaByKey 取某 key 下任一失效资源的基本信息（失效页面优化）。
// 走 key 索引 + LIMIT 1，只读一行；未命中返回 (nil, nil)。
func (r *ResourceRepositoryImpl) FindFailedResourceMetaByKey(key string) (*FailedResourceMeta, error) {
	var meta FailedResourceMeta
	err := r.db.Model(&entity.Resource{}).
		Select("title, cover, key").
		Where("key = ? AND is_valid = ?", key, false).
		Limit(1).
		Scan(&meta).Error
	if err != nil {
		return nil, err
	}
	// Scan 无记录时返回零值 struct（不报错）；据此判断是否命中
	if meta.Key == "" && meta.Title == "" {
		return nil, nil
	}
	return &meta, nil
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

// DeleteRelatedResources 删除关联资源，清空 fid、ck_id 和 save_url 三个字段
func (r *ResourceRepositoryImpl) DeleteRelatedResources(ckID uint) (int64, error) {
	result := r.db.Model(&entity.Resource{}).
		Where("ck_id = ?", ckID).
		Updates(map[string]interface{}{
			"fid":      nil, // 清空 fid 字段
			"ck_id":    0,   // 清空 ck_id 字段
			"save_url": "",  // 清空 save_url 字段
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// CountResourcesByCkID 统计指定账号ID的资源数量
func (r *ResourceRepositoryImpl) CountResourcesByCkID(ckID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Resource{}).
		Where("ck_id = ?", ckID).
		Count(&count).Error
	return count, err
}

// FindByKey 根据Key查找资源（同一组资源）
func (r *ResourceRepositoryImpl) FindByKey(key string) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.db.Where("key = ? AND is_valid = ?", key, true).
		Preload("Category").
		Preload("Pan").
		Preload("Tags").
		Order("pan_id ASC").
		Find(&resources).Error
	return resources, err
}

// GetHotResources 获取热门资源（按查看次数排序，去重，限制数量）
func (r *ResourceRepositoryImpl) GetHotResources(limit int) ([]entity.Resource, error) {
	var resources []entity.Resource

	// 按key分组，获取每个key中查看次数最高的资源，然后按查看次数排序
	err := r.db.Table("resources").
		Select(`
			resources.*,
			ROW_NUMBER() OVER (PARTITION BY key ORDER BY view_count DESC) as rn
		`).
		Where("is_public = ? AND is_valid = ? AND view_count > 0", true, true).
		Preload("Category").
		Preload("Pan").
		Preload("Tags").
		Order("view_count DESC").
		Limit(limit * 2). // 获取更多数据以确保去重后有足够的结果
		Find(&resources).Error

	if err != nil {
		return nil, err
	}

	// 按key去重，保留每个key的第一个（即查看次数最高的）
	seenKeys := make(map[string]bool)
	var hotResources []entity.Resource
	for _, resource := range resources {
		if !seenKeys[resource.Key] {
			seenKeys[resource.Key] = true
			hotResources = append(hotResources, resource)
			if len(hotResources) >= limit {
				break
			}
		}
	}

	return hotResources, nil
}

// FindByResourceKey 根据资源Key查找资源
func (r *ResourceRepositoryImpl) FindByResourceKey(key string) ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.GetDB().Where("key = ?", key).Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// GetTotalCount 获取资源总数
func (r *ResourceRepositoryImpl) GetTotalCount() (int64, error) {
	var count int64
	err := r.GetDB().Model(&entity.Resource{}).Count(&count).Error
	return count, err
}

// GetAllValidResources 获取所有有效的公开资源
func (r *ResourceRepositoryImpl) GetAllValidResources() ([]entity.Resource, error) {
	var resources []entity.Resource
	err := r.GetDB().Where("is_valid = ? AND is_public = ?", true, true).
		Find(&resources).Error
	return resources, err
}

// FindDueForCleanup 查询需要清理的已转存资源
// 条件：已转存(transferred_at不为空)、未清理(cleaned_at为空)、转存时间超过保留期、有有效 fid
// 使用内存过滤模式以保证跨数据库兼容（参照 telegram_channel_repository.go 的 FindDueForPush）
func (r *ResourceRepositoryImpl) FindDueForCleanup(retentionDays int, limit int) ([]*entity.Resource, error) {
	if retentionDays <= 0 {
		return nil, fmt.Errorf("保留天数必须大于0")
	}
	threshold := time.Now().AddDate(0, 0, -retentionDays)

	// 由于不同数据库对指针时间字段的查询行为不一致，
	// 先查询候选集合（已转存且 fid 不为空），再在内存中过滤
	var candidates []*entity.Resource
	query := r.db.Where("transferred_at IS NOT NULL").
		Where("fid IS NOT NULL AND fid <> ''").
		Where("cleaned_at IS NULL")
	if limit > 0 {
		query = query.Limit(limit * 2) // 多取一些以应对过滤后的数量
	}
	err := query.Find(&candidates).Error
	if err != nil {
		return nil, err
	}

	// 内存过滤：转存时间超过保留期
	result := make([]*entity.Resource, 0, limit)
	for _, res := range candidates {
		if res.TransferredAt != nil && res.TransferredAt.Before(threshold) {
			result = append(result, res)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}

	utils.Debug("FindDueForCleanup: 候选=%d, 命中=%d, 保留期=%d天", len(candidates), len(result), retentionDays)
	return result, nil
}

// MarkCleaned 标记资源清理成功：清空 fid/save_url 等转存写入字段，并设置 cleaned_at
func (r *ResourceRepositoryImpl) MarkCleaned(id uint, cleanedAt time.Time) error {
	return r.db.Model(&entity.Resource{}).Where("id = ?", id).Updates(map[string]interface{}{
		"fid":              "",
		"save_url":         "",
		"cleaned_at":       cleanedAt,
		"clean_error_msg":  "",
		"last_clean_error_at": nil,
	}).Error
}

// MarkCleanError 标记资源清理失败：记录失败原因和时间，不清空转存字段（等待下一轮重试）
func (r *ResourceRepositoryImpl) MarkCleanError(id uint, errMsg string, errAt time.Time) error {
	// 截断错误信息避免超出字段长度
	if len(errMsg) > 255 {
		errMsg = errMsg[:255]
	}
	return r.db.Model(&entity.Resource{}).Where("id = ?", id).Updates(map[string]interface{}{
		"clean_error_msg":    errMsg,
		"last_clean_error_at": errAt,
	}).Error
}

// UpdateFields 按主键更新指定字段；用于重转等场景的部分字段更新，
// 避免使用 GORM Updates(struct) 时跳过零值字段。
func (r *ResourceRepositoryImpl) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&entity.Resource{}).Where("id = ?", id).Updates(fields).Error
}

// GenerateUniqueKey 生成唯一的6位Base62资源Key，用于 /r/:key 短链访问。
// 委托给 BaseRepositoryImpl.GenerateUniqueKey（泛型实现，按 resource.key 字段查重）。
func (r *ResourceRepositoryImpl) GenerateUniqueKey() (string, error) {
	return r.BaseRepositoryImpl.GenerateUniqueKey("key")
}
