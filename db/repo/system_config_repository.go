package repo

import (
	"fmt"

	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// SystemConfigRepository 系统配置Repository接口
type SystemConfigRepository interface {
	BaseRepository[entity.SystemConfig]
	FindAll() ([]entity.SystemConfig, error)
	FindByKey(key string) (*entity.SystemConfig, error)
	GetOrCreateDefault() ([]entity.SystemConfig, error)
	UpsertConfigs(configs []entity.SystemConfig) error
	GetConfigValue(key string) (string, error)
	GetConfigBool(key string) (bool, error)
	GetConfigInt(key string) (int, error)
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

// FindAll 获取所有配置
func (r *SystemConfigRepositoryImpl) FindAll() ([]entity.SystemConfig, error) {
	var configs []entity.SystemConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// FindByKey 根据键查找配置
func (r *SystemConfigRepositoryImpl) FindByKey(key string) (*entity.SystemConfig, error) {
	var config entity.SystemConfig
	err := r.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpsertConfigs 批量创建或更新配置
func (r *SystemConfigRepositoryImpl) UpsertConfigs(configs []entity.SystemConfig) error {
	for _, config := range configs {
		var existingConfig entity.SystemConfig
		err := r.db.Where("key = ?", config.Key).First(&existingConfig).Error

		if err != nil {
			// 如果不存在，则创建
			if err := r.db.Create(&config).Error; err != nil {
				return err
			}
		} else {
			// 如果存在，则更新
			config.ID = existingConfig.ID
			if err := r.db.Save(&config).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// GetOrCreateDefault 获取配置或创建默认配置
func (r *SystemConfigRepositoryImpl) GetOrCreateDefault() ([]entity.SystemConfig, error) {
	configs, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// 如果没有配置，创建默认配置
	if len(configs) == 0 {
		defaultConfigs := []entity.SystemConfig{
			{Key: entity.ConfigKeySiteTitle, Value: entity.ConfigDefaultSiteTitle, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeySiteDescription, Value: entity.ConfigDefaultSiteDescription, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyKeywords, Value: entity.ConfigDefaultKeywords, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyAuthor, Value: entity.ConfigDefaultAuthor, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyCopyright, Value: entity.ConfigDefaultCopyright, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyAutoProcessReadyResources, Value: entity.ConfigDefaultAutoProcessReadyResources, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyAutoProcessInterval, Value: entity.ConfigDefaultAutoProcessInterval, Type: entity.ConfigTypeInt},
			{Key: entity.ConfigKeyAutoTransferEnabled, Value: entity.ConfigDefaultAutoTransferEnabled, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyAutoTransferLimitDays, Value: entity.ConfigDefaultAutoTransferLimitDays, Type: entity.ConfigTypeInt},
			{Key: entity.ConfigKeyAutoTransferMinSpace, Value: entity.ConfigDefaultAutoTransferMinSpace, Type: entity.ConfigTypeInt},
			{Key: entity.ConfigKeyAutoFetchHotDramaEnabled, Value: entity.ConfigDefaultAutoFetchHotDramaEnabled, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyApiToken, Value: entity.ConfigDefaultApiToken, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyPageSize, Value: entity.ConfigDefaultPageSize, Type: entity.ConfigTypeInt},
			{Key: entity.ConfigKeyMaintenanceMode, Value: entity.ConfigDefaultMaintenanceMode, Type: entity.ConfigTypeBool},
		}

		err = r.UpsertConfigs(defaultConfigs)
		if err != nil {
			return nil, err
		}

		return defaultConfigs, nil
	}

	return configs, nil
}

// GetConfigValue 获取配置值（字符串）
func (r *SystemConfigRepositoryImpl) GetConfigValue(key string) (string, error) {
	config, err := r.FindByKey(key)
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

// GetConfigBool 获取配置值（布尔）
func (r *SystemConfigRepositoryImpl) GetConfigBool(key string) (bool, error) {
	value, err := r.GetConfigValue(key)
	if err != nil {
		return false, err
	}

	switch value {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, nil
	}
}

// GetConfigInt 获取配置值（整数）
func (r *SystemConfigRepositoryImpl) GetConfigInt(key string) (int, error) {
	value, err := r.GetConfigValue(key)
	if err != nil {
		return 0, err
	}

	// 这里需要导入 strconv 包，但为了避免循环导入，我们使用简单的转换
	var result int
	_, err = fmt.Sscanf(value, "%d", &result)
	return result, err
}
