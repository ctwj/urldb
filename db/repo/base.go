package repo

import (
	"gorm.io/gorm"
)

// BaseRepository 基础Repository接口
type BaseRepository[T any] interface {
	Create(entity *T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]T, error)
	Update(entity *T) error
	Delete(id uint) error
	FindWithPagination(page, limit int) ([]T, int64, error)
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
