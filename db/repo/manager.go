package repo

import (
	"gorm.io/gorm"
)

// RepositoryManager Repository管理器
type RepositoryManager struct {
	PanRepository           PanRepository
	CksRepository           CksRepository
	ResourceRepository      ResourceRepository
	CategoryRepository      CategoryRepository
	TagRepository           TagRepository
	ReadyResourceRepository ReadyResourceRepository
}

// NewRepositoryManager 创建Repository管理器
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		PanRepository:           NewPanRepository(db),
		CksRepository:           NewCksRepository(db),
		ResourceRepository:      NewResourceRepository(db),
		CategoryRepository:      NewCategoryRepository(db),
		TagRepository:           NewTagRepository(db),
		ReadyResourceRepository: NewReadyResourceRepository(db),
	}
}
