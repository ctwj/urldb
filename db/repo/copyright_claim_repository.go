package repo

import (
	"gorm.io/gorm"
	"github.com/ctwj/urldb/db/entity"
)

// CopyrightClaimRepository 版权申述Repository接口
type CopyrightClaimRepository interface {
	BaseRepository[entity.CopyrightClaim]
	GetByResourceKey(resourceKey string) ([]*entity.CopyrightClaim, error)
	List(status string, page, pageSize int) ([]*entity.CopyrightClaim, int64, error)
	UpdateStatus(id uint, status string, processedBy *uint, note string) error
	// 兼容原有方法名
	GetByID(id uint) (*entity.CopyrightClaim, error)
}

// CopyrightClaimRepositoryImpl 版权申述Repository实现
type CopyrightClaimRepositoryImpl struct {
	BaseRepositoryImpl[entity.CopyrightClaim]
}

// NewCopyrightClaimRepository 创建版权申述Repository
func NewCopyrightClaimRepository(db *gorm.DB) CopyrightClaimRepository {
	return &CopyrightClaimRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.CopyrightClaim]{db: db},
	}
}

// Create 创建版权申述
func (r *CopyrightClaimRepositoryImpl) Create(claim *entity.CopyrightClaim) error {
	return r.GetDB().Create(claim).Error
}

// GetByID 根据ID获取版权申述
func (r *CopyrightClaimRepositoryImpl) GetByID(id uint) (*entity.CopyrightClaim, error) {
	var claim entity.CopyrightClaim
	err := r.GetDB().Where("id = ?", id).First(&claim).Error
	return &claim, err
}

// GetByResourceKey 获取某个资源的所有版权申述
func (r *CopyrightClaimRepositoryImpl) GetByResourceKey(resourceKey string) ([]*entity.CopyrightClaim, error) {
	var claims []*entity.CopyrightClaim
	err := r.GetDB().Where("resource_key = ?", resourceKey).Find(&claims).Error
	return claims, err
}

// List 获取版权申述列表
func (r *CopyrightClaimRepositoryImpl) List(status string, page, pageSize int) ([]*entity.CopyrightClaim, int64, error) {
	var claims []*entity.CopyrightClaim
	var total int64

	query := r.GetDB().Model(&entity.CopyrightClaim{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&claims).Error
	return claims, total, err
}

// Update 更新版权申述
func (r *CopyrightClaimRepositoryImpl) Update(claim *entity.CopyrightClaim) error {
	return r.GetDB().Save(claim).Error
}

// UpdateStatus 更新版权申述状态
func (r *CopyrightClaimRepositoryImpl) UpdateStatus(id uint, status string, processedBy *uint, note string) error {
	return r.GetDB().Model(&entity.CopyrightClaim{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"processed_at":  gorm.Expr("NOW()"),
		"processed_by":  processedBy,
		"note":          note,
	}).Error
}

// Delete 删除版权申述
func (r *CopyrightClaimRepositoryImpl) Delete(id uint) error {
	return r.GetDB().Delete(&entity.CopyrightClaim{}, id).Error
}