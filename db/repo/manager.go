package repo

import (
	"gorm.io/gorm"
)

// RepositoryManager Repository管理器
type RepositoryManager struct {
	DB                        *gorm.DB
	PanRepository             PanRepository
	CksRepository             CksRepository
	ResourceRepository        ResourceRepository
	CategoryRepository        CategoryRepository
	TagRepository             TagRepository
	ReadyResourceRepository   ReadyResourceRepository
	UserRepository            UserRepository
	SearchStatRepository      SearchStatRepository
	SystemConfigRepository    SystemConfigRepository
	HotDramaRepository        HotDramaRepository
	ResourceViewRepository    ResourceViewRepository
	TaskRepository            TaskRepository
	TaskItemRepository        TaskItemRepository
	FileRepository            FileRepository
	TelegramChannelRepository TelegramChannelRepository
	APIAccessLogRepository    APIAccessLogRepository
	ReportRepository          ReportRepository
	CopyrightClaimRepository  CopyrightClaimRepository
	PluginConfigRepository    *PluginConfigRepository
	PluginLogRepository       *PluginLogRepository
	CronJobRepository         *CronJobRepository
}

// NewRepositoryManager 创建Repository管理器
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		DB:                        db,
		PanRepository:             NewPanRepository(db),
		CksRepository:             NewCksRepository(db),
		ResourceRepository:        NewResourceRepository(db),
		CategoryRepository:        NewCategoryRepository(db),
		TagRepository:             NewTagRepository(db),
		ReadyResourceRepository:   NewReadyResourceRepository(db),
		UserRepository:            NewUserRepository(db),
		SearchStatRepository:      NewSearchStatRepository(db),
		SystemConfigRepository:    NewSystemConfigRepository(db),
		HotDramaRepository:        NewHotDramaRepository(db),
		ResourceViewRepository:    NewResourceViewRepository(db),
		TaskRepository:            NewTaskRepository(db),
		TaskItemRepository:        NewTaskItemRepository(db),
		FileRepository:            NewFileRepository(db),
		TelegramChannelRepository: NewTelegramChannelRepository(db),
		APIAccessLogRepository:    NewAPIAccessLogRepository(db),
		ReportRepository:          NewReportRepository(db),
		CopyrightClaimRepository:  NewCopyrightClaimRepository(db),
		PluginConfigRepository:    NewPluginConfigRepository(db),
		PluginLogRepository:       NewPluginLogRepository(db),
		CronJobRepository:         NewCronJobRepository(db),
	}
}

// GetDB 获取数据库连接
func (rm *RepositoryManager) GetDB() *gorm.DB {
	return rm.DB
}

// GetPluginLogRepository 获取插件日志仓库
func (rm *RepositoryManager) GetPluginLogRepository() *PluginLogRepository {
	return rm.PluginLogRepository
}
