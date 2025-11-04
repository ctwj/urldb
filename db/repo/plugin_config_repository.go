package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// PluginConfigRepository 插件配置Repository接口
type PluginConfigRepository interface {
	BaseRepository[entity.PluginConfig]
	FindByPluginAndKey(pluginName, key string) (*entity.PluginConfig, error)
	FindByPlugin(pluginName string) ([]entity.PluginConfig, error)
	Upsert(pluginName, key, value, configType string, isEncrypted bool, description string) error
	DeleteByPluginAndKey(pluginName, key string) error
	DeleteByPlugin(pluginName string) error
}

// PluginConfigRepositoryImpl 插件配置Repository实现
type PluginConfigRepositoryImpl struct {
	BaseRepositoryImpl[entity.PluginConfig]
}

// NewPluginConfigRepository 创建插件配置Repository
func NewPluginConfigRepository(db *gorm.DB) PluginConfigRepository {
	return &PluginConfigRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.PluginConfig]{db: db},
	}
}

// FindByPluginAndKey 根据插件名称和键查找配置
func (r *PluginConfigRepositoryImpl) FindByPluginAndKey(pluginName, key string) (*entity.PluginConfig, error) {
	var config entity.PluginConfig
	err := r.db.Where("plugin_name = ? AND config_key = ?", pluginName, key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FindByPlugin 根据插件名称查找所有配置
func (r *PluginConfigRepositoryImpl) FindByPlugin(pluginName string) ([]entity.PluginConfig, error) {
	var configs []entity.PluginConfig
	err := r.db.Where("plugin_name = ?", pluginName).Find(&configs).Error
	return configs, err
}

// Upsert 创建或更新插件配置
func (r *PluginConfigRepositoryImpl) Upsert(pluginName, key, value, configType string, isEncrypted bool, description string) error {
	var existingConfig entity.PluginConfig
	err := r.db.Where("plugin_name = ? AND config_key = ?", pluginName, key).First(&existingConfig).Error

	if err != nil {
		// 如果不存在，则创建
		newConfig := entity.PluginConfig{
			PluginName:  pluginName,
			ConfigKey:   key,
			ConfigValue: value,
			ConfigType:  configType,
			IsEncrypted: isEncrypted,
			Description: description,
		}
		return r.db.Create(&newConfig).Error
	} else {
		// 如果存在，则更新
		existingConfig.ConfigValue = value
		existingConfig.ConfigType = configType
		existingConfig.IsEncrypted = isEncrypted
		existingConfig.Description = description
		return r.db.Save(&existingConfig).Error
	}
}

// DeleteByPluginAndKey 根据插件名称和键删除配置
func (r *PluginConfigRepositoryImpl) DeleteByPluginAndKey(pluginName, key string) error {
	return r.db.Where("plugin_name = ? AND config_key = ?", pluginName, key).Delete(&entity.PluginConfig{}).Error
}

// DeleteByPlugin 根据插件名称删除所有配置
func (r *PluginConfigRepositoryImpl) DeleteByPlugin(pluginName string) error {
	return r.db.Where("plugin_name = ?", pluginName).Delete(&entity.PluginConfig{}).Error
}