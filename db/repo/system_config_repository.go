package repo

import (
	"fmt"
	"sync"

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
	GetCachedConfigs() map[string]string
	ClearConfigCache()
}

// SystemConfigRepositoryImpl 系统配置Repository实现
type SystemConfigRepositoryImpl struct {
	BaseRepositoryImpl[entity.SystemConfig]

	// 配置缓存
	configCache      map[string]string // key -> value
	configCacheOnce  sync.Once
	configCacheMutex sync.RWMutex
}

// NewSystemConfigRepository 创建系统配置Repository
func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &SystemConfigRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.SystemConfig]{db: db},
		configCache:        make(map[string]string),
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

	// 更新配置后刷新缓存
	r.refreshConfigCache()
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

// initConfigCache 初始化配置缓存
func (r *SystemConfigRepositoryImpl) initConfigCache() {
	r.configCacheOnce.Do(func() {
		// 获取所有配置
		configs, err := r.FindAll()
		if err != nil {
			// 如果获取失败，尝试创建默认配置
			configs, err = r.GetOrCreateDefault()
			if err != nil {
				return
			}
		}

		// 初始化缓存
		r.configCacheMutex.Lock()
		defer r.configCacheMutex.Unlock()

		for _, config := range configs {
			r.configCache[config.Key] = config.Value
		}
	})
}

// refreshConfigCache 刷新配置缓存
func (r *SystemConfigRepositoryImpl) refreshConfigCache() {
	// 重置Once，允许重新初始化
	r.configCacheOnce = sync.Once{}

	// 清空缓存
	r.configCacheMutex.Lock()
	r.configCache = make(map[string]string)
	r.configCacheMutex.Unlock()

	// 重新初始化缓存
	r.initConfigCache()
}

// GetConfigValue 获取配置值（字符串）
func (r *SystemConfigRepositoryImpl) GetConfigValue(key string) (string, error) {
	// 初始化缓存
	r.initConfigCache()

	// 从缓存中读取
	r.configCacheMutex.RLock()
	value, exists := r.configCache[key]
	r.configCacheMutex.RUnlock()

	if exists {
		return value, nil
	}

	// 如果缓存中没有，尝试从数据库获取（可能是新添加的配置）
	config, err := r.FindByKey(key)
	if err != nil {
		return "", err
	}

	// 更新缓存
	r.configCacheMutex.Lock()
	r.configCache[key] = config.Value
	r.configCacheMutex.Unlock()

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

// GetCachedConfigs 获取所有缓存的配置（用于调试）
func (r *SystemConfigRepositoryImpl) GetCachedConfigs() map[string]string {
	r.initConfigCache()

	r.configCacheMutex.RLock()
	defer r.configCacheMutex.RUnlock()

	// 返回缓存的副本
	result := make(map[string]string)
	for k, v := range r.configCache {
		result[k] = v
	}

	return result
}

// ClearConfigCache 清空配置缓存（用于测试或手动刷新）
func (r *SystemConfigRepositoryImpl) ClearConfigCache() {
	r.configCacheMutex.Lock()
	r.configCache = make(map[string]string)
	r.configCacheMutex.Unlock()

	// 重置Once，允许重新初始化
	r.configCacheOnce = sync.Once{}
}
