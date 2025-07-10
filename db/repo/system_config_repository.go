package repo

import (
	"res_db/db/entity"

	"gorm.io/gorm"
)

// SystemConfigRepository 系统配置Repository接口
type SystemConfigRepository interface {
	BaseRepository[entity.SystemConfig]
	FindFirst() (*entity.SystemConfig, error)
	GetOrCreateDefault() (*entity.SystemConfig, error)
	Upsert(config *entity.SystemConfig) error
}

// SystemConfigRepositoryImpl 系统配置Repository实现
type SystemConfigRepositoryImpl struct {
	BaseRepositoryImpl[entity.SystemConfig]
}

// NewSystemConfigRepository 创建系统配置Repository
func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &SystemConfigRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.SystemConfig]{db: db},
	}
}

// FindFirst 获取第一个配置（通常只有一个配置）
func (r *SystemConfigRepositoryImpl) FindFirst() (*entity.SystemConfig, error) {
	var config entity.SystemConfig
	err := r.db.First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Upsert 创建或更新系统配置
func (r *SystemConfigRepositoryImpl) Upsert(config *entity.SystemConfig) error {
	var existingConfig entity.SystemConfig
	err := r.db.First(&existingConfig).Error

	if err != nil {
		// 如果不存在，则创建
		return r.db.Create(config).Error
	} else {
		// 如果存在，则更新
		config.ID = existingConfig.ID
		return r.db.Save(config).Error
	}
}

// GetOrCreateDefault 获取配置或创建默认配置
func (r *SystemConfigRepositoryImpl) GetOrCreateDefault() (*entity.SystemConfig, error) {
	config, err := r.FindFirst()
	if err != nil {
		// 创建默认配置
		defaultConfig := &entity.SystemConfig{
			SiteTitle:                 "网盘资源管理系统",
			SiteDescription:           "专业的网盘资源管理系统",
			Keywords:                  "网盘,资源管理,文件分享",
			Author:                    "系统管理员",
			Copyright:                 "© 2024 网盘资源管理系统",
			AutoProcessReadyResources: false,
			AutoProcessInterval:       30,
			PageSize:                  100,
			MaintenanceMode:           false,
		}

		err = r.db.Create(defaultConfig).Error
		if err != nil {
			return nil, err
		}

		return defaultConfig, nil
	}

	return config, nil
}
