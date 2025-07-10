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
	GetResourceCount(tagID uint) (int64, error)
	FindByResourceID(resourceID uint) ([]entity.Tag, error)
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
