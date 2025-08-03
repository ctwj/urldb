package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// ReadyResourceRepository ReadyResource的Repository接口
type ReadyResourceRepository interface {
	BaseRepository[entity.ReadyResource]
	FindByURL(url string) (*entity.ReadyResource, error)
	FindByIP(ip string) ([]entity.ReadyResource, error)
	FindByKey(key string) ([]entity.ReadyResource, error)
	BatchCreate(resources []entity.ReadyResource) error
	DeleteByURL(url string) error
	DeleteByKey(key string) error
	FindAllWithinDays(days int) ([]entity.ReadyResource, error)
	BatchFindByURLs(urls []string) ([]entity.ReadyResource, error)
	GenerateUniqueKey() (string, error)
	FindWithErrors() ([]entity.ReadyResource, error)
	FindWithErrorsPaginated(page, limit int) ([]entity.ReadyResource, int64, error)
	FindWithErrorsIncludingDeleted() ([]entity.ReadyResource, error)
	FindWithErrorsPaginatedIncludingDeleted(page, limit int, errorFilter string) ([]entity.ReadyResource, int64, error)
	FindWithErrorsByQuery(errorFilter string) ([]entity.ReadyResource, error)
	FindWithoutErrors() ([]entity.ReadyResource, error)
	ClearErrorMsg(id uint) error
	ClearErrorMsgAndRestore(id uint) error
	ClearAllErrorsByQuery(errorFilter string) (int64, error) // 批量清除错误信息并真正删除资源
}

// ReadyResourceRepositoryImpl ReadyResource的Repository实现
type ReadyResourceRepositoryImpl struct {
	BaseRepositoryImpl[entity.ReadyResource]
}

// NewReadyResourceRepository 创建ReadyResource Repository
func NewReadyResourceRepository(db *gorm.DB) ReadyResourceRepository {
	return &ReadyResourceRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.ReadyResource]{db: db},
	}
}

// FindByURL 根据URL查找
func (r *ReadyResourceRepositoryImpl) FindByURL(url string) (*entity.ReadyResource, error) {
	var resource entity.ReadyResource
	err := r.db.Where("url = ?", url).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// FindByIP 根据IP查找
func (r *ReadyResourceRepositoryImpl) FindByIP(ip string) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	err := r.db.Where("ip = ?", ip).Find(&resources).Error
	return resources, err
}

// BatchCreate 批量创建
func (r *ReadyResourceRepositoryImpl) BatchCreate(resources []entity.ReadyResource) error {
	return r.db.Create(&resources).Error
}

// DeleteByURL 根据URL删除
func (r *ReadyResourceRepositoryImpl) DeleteByURL(url string) error {
	return r.db.Where("url = ?", url).Delete(&entity.ReadyResource{}).Error
}

// FindAllWithinDays 获取n天内的待处理资源，n=0时不限制
func (r *ReadyResourceRepositoryImpl) FindAllWithinDays(days int) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	db := r.db.Model(&entity.ReadyResource{})
	if days > 0 {
		since := time.Now().AddDate(0, 0, -days)
		db = db.Where("create_time >= ?", since)
	}
	err := db.Find(&resources).Error
	return resources, err
}

func (r *ReadyResourceRepositoryImpl) BatchFindByURLs(urls []string) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	if len(urls) == 0 {
		return resources, nil
	}
	err := r.db.Where("url IN ?", urls).Find(&resources).Error
	return resources, err
}

// FindByKey 根据Key查找
func (r *ReadyResourceRepositoryImpl) FindByKey(key string) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	err := r.db.Where("key = ?", key).Find(&resources).Error
	return resources, err
}

// DeleteByKey 根据Key删除
func (r *ReadyResourceRepositoryImpl) DeleteByKey(key string) error {
	return r.db.Where("key = ?", key).Delete(&entity.ReadyResource{}).Error
}

// GenerateUniqueKey 生成唯一的6位Base62 key
func (r *ReadyResourceRepositoryImpl) GenerateUniqueKey() (string, error) {
	for i := 0; i < 20; i++ {
		key, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 6)
		if err != nil {
			return "", err
		}
		var count int64
		err = r.db.Model(&entity.ReadyResource{}).Where("key = ?", key).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count == 0 {
			return key, nil
		}
	}
	return "", gorm.ErrInvalidData
}

// FindWithErrors 查找有错误信息的资源（包括软删除的）
func (r *ReadyResourceRepositoryImpl) FindWithErrors() ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	err := r.db.Unscoped().Where("error_msg != '' AND error_msg IS NOT NULL").Find(&resources).Error
	return resources, err
}

// FindWithErrorsPaginated 分页查找有错误信息的资源（包括软删除的）
func (r *ReadyResourceRepositoryImpl) FindWithErrorsPaginated(page, limit int) ([]entity.ReadyResource, int64, error) {
	var resources []entity.ReadyResource
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&entity.ReadyResource{}).Unscoped().Where("error_msg != '' AND error_msg IS NOT NULL")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&resources).Error
	return resources, total, err
}

// FindWithoutErrors 查找没有错误信息的资源
func (r *ReadyResourceRepositoryImpl) FindWithoutErrors() ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	err := r.db.Where("error_msg = '' OR error_msg IS NULL").Find(&resources).Error
	return resources, err
}

// FindWithErrorsIncludingDeleted 查找有错误信息的资源（包括软删除的，用于管理页面）
func (r *ReadyResourceRepositoryImpl) FindWithErrorsIncludingDeleted() ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource
	err := r.db.Unscoped().Where("error_msg != '' AND error_msg IS NOT NULL").Find(&resources).Error
	return resources, err
}

// FindWithErrorsPaginatedIncludingDeleted 分页查找有错误信息的资源（包括软删除的，用于管理页面）
func (r *ReadyResourceRepositoryImpl) FindWithErrorsPaginatedIncludingDeleted(page, limit int, errorFilter string) ([]entity.ReadyResource, int64, error) {
	var resources []entity.ReadyResource
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&entity.ReadyResource{}).Unscoped().Where("error_msg != '' AND error_msg IS NOT NULL")

	// 如果有错误过滤条件，添加到查询中
	if errorFilter != "" {
		db = db.Where("error_msg ILIKE ?", "%"+errorFilter+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&resources).Error
	return resources, total, err
}

// ClearErrorMsg 清除指定资源的错误信息
func (r *ReadyResourceRepositoryImpl) ClearErrorMsg(id uint) error {
	return r.db.Model(&entity.ReadyResource{}).Where("id = ?", id).Update("error_msg", "").Update("deleted_at", nil).Error
}

// ClearErrorMsgAndRestore 清除错误信息并恢复软删除的资源
func (r *ReadyResourceRepositoryImpl) ClearErrorMsgAndRestore(id uint) error {
	return r.db.Unscoped().Model(&entity.ReadyResource{}).Where("id = ?", id).Updates(map[string]interface{}{
		"error_msg":  "",
		"deleted_at": nil,
	}).Error
}

// FindWithErrorsByQuery 根据查询条件查找有错误信息的资源（不分页，用于批量操作）
func (r *ReadyResourceRepositoryImpl) FindWithErrorsByQuery(errorFilter string) ([]entity.ReadyResource, error) {
	var resources []entity.ReadyResource

	db := r.db.Model(&entity.ReadyResource{}).Unscoped().Where("error_msg != '' AND error_msg IS NOT NULL")

	// 如果有错误过滤条件，添加到查询中
	if errorFilter != "" {
		db = db.Where("error_msg ILIKE ?", "%"+errorFilter+"%")
	}

	err := db.Order("created_at DESC").Find(&resources).Error
	return resources, err
}

// ClearAllErrorsByQuery 根据查询条件批量清除错误信息并真正删除资源
func (r *ReadyResourceRepositoryImpl) ClearAllErrorsByQuery(errorFilter string) (int64, error) {
	db := r.db.Unscoped().Model(&entity.ReadyResource{}).Where("error_msg != '' AND error_msg IS NOT NULL")

	// 如果有错误过滤条件，添加到查询中
	if errorFilter != "" {
		db = db.Where("error_msg ILIKE ?", "%"+errorFilter+"%")
	}

	// 真正删除资源（物理删除）
	result := db.Delete(&entity.ReadyResource{})

	return result.RowsAffected, result.Error
}
