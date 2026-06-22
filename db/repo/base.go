package repo

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// keyAlphabet Base62 字符集，用于生成 6 位短 key
const keyAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// keyLength 资源 key 长度（6 位 Base62，约 568 亿组合，碰撞概率可忽略）
const keyLength = 6

// keyMaxRetry 生成 key 的最大重试次数（防极端情况下死循环）
const keyMaxRetry = 20

// BaseRepository 基础Repository接口
type BaseRepository[T any] interface {
	Create(entity *T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]T, error)
	Update(entity *T) error
	Delete(id uint) error
	FindWithPagination(page, limit int) ([]T, int64, error)
	GetDB() *gorm.DB
}

// BaseRepositoryImpl 基础Repository实现
type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础Repository
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{db: db}
}

// Create 创建实体
func (r *BaseRepositoryImpl[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// FindByID 根据ID查找实体
func (r *BaseRepositoryImpl[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll 查找所有实体
func (r *BaseRepositoryImpl[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// Update 更新实体
func (r *BaseRepositoryImpl[T]) Update(entity *T) error {
	return r.db.Model(entity).Updates(entity).Error
}

// Delete 删除实体
func (r *BaseRepositoryImpl[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// FindWithPagination 分页查找
func (r *BaseRepositoryImpl[T]) FindWithPagination(page, limit int) ([]T, int64, error) {
	var entities []T
	var total int64

	offset := (page - 1) * limit

	// 获取总数
	if err := r.db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Offset(offset).Limit(limit).Find(&entities).Error
	return entities, total, err
}

func (r *BaseRepositoryImpl[T]) GetDB() *gorm.DB {
	return r.db
}

// GenerateUniqueKey 生成唯一的 6 位 Base62 key，并针对当前 repo 的实体类型查重。
// column 指定查重字段名（如 "key"）；最多重试 keyMaxRetry 次，冲突即返回 gorm.ErrInvalidData。
// 所有嵌入 BaseRepositoryImpl 的子 repo 可直接复用，无需各自实现。
func (r *BaseRepositoryImpl[T]) GenerateUniqueKey(column string) (string, error) {
	for i := 0; i < keyMaxRetry; i++ {
		key, err := gonanoid.Generate(keyAlphabet, keyLength)
		if err != nil {
			return "", err
		}
		var count int64
		if err := r.db.Model(new(T)).Where(column+" = ?", key).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return key, nil
		}
	}
	return "", gorm.ErrInvalidData
}
