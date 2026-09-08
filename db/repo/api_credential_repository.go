package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// ApiCredentialRepository API 凭证 Repository 接口（016-api-access-application）
type ApiCredentialRepository interface {
	BaseRepository[entity.ApiCredential]
	// FindByUserID 按归属用户查询（一用户一行；无则返回 nil,nil）
	FindByUserID(userID uint) (*entity.ApiCredential, error)
	// FindByKeyHash 按明文 key 哈希查询（开放接口鉴权）
	FindByKeyHash(keyHash string) (*entity.ApiCredential, error)
	// ExistsActiveUnexpiredByUser 是否存在有效且未过期凭证（重复申请拦截，FR-003）
	ExistsActiveUnexpiredByUser(userID uint, now time.Time) (bool, error)
	// SaveForUser 按归属用户保存（无则创建，有则按 ID 更新）
	SaveForUser(cred *entity.ApiCredential) error
	// FindExpiringBetween 查询到期时间落在 [start, end] 且未发过提醒的 active 凭证（FR-018 每日提醒）
	FindExpiringBetween(start, end time.Time) ([]*entity.ApiCredential, error)
}

// ApiCredentialRepositoryImpl API 凭证 Repository 实现
type ApiCredentialRepositoryImpl struct {
	BaseRepositoryImpl[entity.ApiCredential]
}

// NewApiCredentialRepository 创建 API 凭证 Repository
func NewApiCredentialRepository(db *gorm.DB) ApiCredentialRepository {
	return &ApiCredentialRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.ApiCredential]{db: db},
	}
}

// FindByUserID 按归属用户查询
func (r *ApiCredentialRepositoryImpl) FindByUserID(userID uint) (*entity.ApiCredential, error) {
	var cred entity.ApiCredential
	err := r.GetDB().Where("user_id = ?", userID).First(&cred).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cred, nil
}

// FindByKeyHash 按明文 key 哈希查询
func (r *ApiCredentialRepositoryImpl) FindByKeyHash(keyHash string) (*entity.ApiCredential, error) {
	var cred entity.ApiCredential
	err := r.GetDB().Where("key_hash = ?", keyHash).First(&cred).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cred, nil
}

// ExistsActiveUnexpiredByUser 是否存在有效且未过期凭证
func (r *ApiCredentialRepositoryImpl) ExistsActiveUnexpiredByUser(userID uint, now time.Time) (bool, error) {
	cred, err := r.FindByUserID(userID)
	if err != nil || cred == nil {
		return false, err
	}
	return cred.Status == entity.ApiCredentialStatusActive && !cred.IsExpired(now), nil
}

// SaveForUser 按归属用户保存
func (r *ApiCredentialRepositoryImpl) SaveForUser(cred *entity.ApiCredential) error {
	if cred.ID == 0 {
		var existing entity.ApiCredential
		err := r.GetDB().Where("user_id = ?", cred.UserID).First(&existing).Error
		if err == nil {
			cred.ID = existing.ID
			cred.CreatedAt = existing.CreatedAt
			return r.GetDB().Save(cred).Error
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		return r.GetDB().Create(cred).Error
	}
	return r.GetDB().Save(cred).Error
}

// FindExpiringBetween 查询到期时间落在 [start, end] 且未发过提醒的 active 凭证
func (r *ApiCredentialRepositoryImpl) FindExpiringBetween(start, end time.Time) ([]*entity.ApiCredential, error) {
	var list []*entity.ApiCredential
	err := r.GetDB().
		Where("status = ? AND expires_at IS NOT NULL AND expires_at >= ? AND expires_at <= ? AND reminder_sent_at IS NULL",
			entity.ApiCredentialStatusActive, start, end).
		Find(&list).Error
	return list, err
}
