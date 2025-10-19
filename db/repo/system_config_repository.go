package repo

import (
	"fmt"
	"sync"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

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
	SafeRefreshConfigCache() error
	ValidateConfigIntegrity() error
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
	// 使用事务确保数据一致性
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 在更新前备份当前配置
		var existingConfigs []entity.SystemConfig
		if err := tx.Find(&existingConfigs).Error; err != nil {
			utils.Error("备份配置失败: %v", err)
			// 不返回错误，继续执行更新
		}

		for _, config := range configs {
			var existingConfig entity.SystemConfig
			err := tx.Where("key = ?", config.Key).First(&existingConfig).Error

			if err != nil {
				// 如果不存在，则创建
				if err := tx.Create(&config).Error; err != nil {
					utils.Error("创建配置失败 [%s]: %v", config.Key, err)
					return fmt.Errorf("创建配置失败 [%s]: %v", config.Key, err)
				}
			} else {
				// 如果存在，则更新
				config.ID = existingConfig.ID
				if err := tx.Save(&config).Error; err != nil {
					utils.Error("更新配置失败 [%s]: %v", config.Key, err)
					return fmt.Errorf("更新配置失败 [%s]: %v", config.Key, err)
				}
			}
		}

		// 更新成功后刷新缓存
		r.refreshConfigCache()
		return nil
	})
}

// GetOrCreateDefault 获取配置或创建默认配置
func (r *SystemConfigRepositoryImpl) GetOrCreateDefault() ([]entity.SystemConfig, error) {
	configs, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// 如果没有配置，创建默认配置
	if len(configs) == 0 {
		utils.Info("未找到任何配置，创建默认配置")
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
			{Key: entity.ConfigKeyForbiddenWords, Value: entity.ConfigDefaultForbiddenWords, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyAdKeywords, Value: entity.ConfigDefaultAdKeywords, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyAutoInsertAd, Value: entity.ConfigDefaultAutoInsertAd, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyPageSize, Value: entity.ConfigDefaultPageSize, Type: entity.ConfigTypeInt},
			{Key: entity.ConfigKeyMaintenanceMode, Value: entity.ConfigDefaultMaintenanceMode, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyEnableRegister, Value: entity.ConfigDefaultEnableRegister, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyThirdPartyStatsCode, Value: entity.ConfigDefaultThirdPartyStatsCode, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyMeilisearchEnabled, Value: entity.ConfigDefaultMeilisearchEnabled, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyMeilisearchHost, Value: entity.ConfigDefaultMeilisearchHost, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyMeilisearchPort, Value: entity.ConfigDefaultMeilisearchPort, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyMeilisearchMasterKey, Value: entity.ConfigDefaultMeilisearchMasterKey, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyMeilisearchIndexName, Value: entity.ConfigDefaultMeilisearchIndexName, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyEnableAnnouncements, Value: entity.ConfigDefaultEnableAnnouncements, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyAnnouncements, Value: entity.ConfigDefaultAnnouncements, Type: entity.ConfigTypeJSON},
			{Key: entity.ConfigKeyEnableFloatButtons, Value: entity.ConfigDefaultEnableFloatButtons, Type: entity.ConfigTypeBool},
			{Key: entity.ConfigKeyWechatSearchImage, Value: entity.ConfigDefaultWechatSearchImage, Type: entity.ConfigTypeString},
			{Key: entity.ConfigKeyTelegramQrImage, Value: entity.ConfigDefaultTelegramQrImage, Type: entity.ConfigTypeString},
		}

		err = r.UpsertConfigs(defaultConfigs)
		if err != nil {
			return nil, err
		}

		return defaultConfigs, nil
	}

	// 检查是否有缺失的配置项，如果有则添加
	requiredConfigs := map[string]entity.SystemConfig{
		entity.ConfigKeySiteTitle:                 {Key: entity.ConfigKeySiteTitle, Value: entity.ConfigDefaultSiteTitle, Type: entity.ConfigTypeString},
		entity.ConfigKeySiteDescription:           {Key: entity.ConfigKeySiteDescription, Value: entity.ConfigDefaultSiteDescription, Type: entity.ConfigTypeString},
		entity.ConfigKeyKeywords:                  {Key: entity.ConfigKeyKeywords, Value: entity.ConfigDefaultKeywords, Type: entity.ConfigTypeString},
		entity.ConfigKeyAuthor:                    {Key: entity.ConfigKeyAuthor, Value: entity.ConfigDefaultAuthor, Type: entity.ConfigTypeString},
		entity.ConfigKeyCopyright:                 {Key: entity.ConfigKeyCopyright, Value: entity.ConfigDefaultCopyright, Type: entity.ConfigTypeString},
		entity.ConfigKeyAutoProcessReadyResources: {Key: entity.ConfigKeyAutoProcessReadyResources, Value: entity.ConfigDefaultAutoProcessReadyResources, Type: entity.ConfigTypeBool},
		entity.ConfigKeyAutoProcessInterval:       {Key: entity.ConfigKeyAutoProcessInterval, Value: entity.ConfigDefaultAutoProcessInterval, Type: entity.ConfigTypeInt},
		entity.ConfigKeyAutoTransferEnabled:       {Key: entity.ConfigKeyAutoTransferEnabled, Value: entity.ConfigDefaultAutoTransferEnabled, Type: entity.ConfigTypeBool},
		entity.ConfigKeyAutoTransferLimitDays:     {Key: entity.ConfigKeyAutoTransferLimitDays, Value: entity.ConfigDefaultAutoTransferLimitDays, Type: entity.ConfigTypeInt},
		entity.ConfigKeyAutoTransferMinSpace:      {Key: entity.ConfigKeyAutoTransferMinSpace, Value: entity.ConfigDefaultAutoTransferMinSpace, Type: entity.ConfigTypeInt},
		entity.ConfigKeyAutoFetchHotDramaEnabled:  {Key: entity.ConfigKeyAutoFetchHotDramaEnabled, Value: entity.ConfigDefaultAutoFetchHotDramaEnabled, Type: entity.ConfigTypeBool},
		entity.ConfigKeyApiToken:                  {Key: entity.ConfigKeyApiToken, Value: entity.ConfigDefaultApiToken, Type: entity.ConfigTypeString},
		entity.ConfigKeyForbiddenWords:            {Key: entity.ConfigKeyForbiddenWords, Value: entity.ConfigDefaultForbiddenWords, Type: entity.ConfigTypeString},
		entity.ConfigKeyAdKeywords:                {Key: entity.ConfigKeyAdKeywords, Value: entity.ConfigDefaultAdKeywords, Type: entity.ConfigTypeString},
		entity.ConfigKeyAutoInsertAd:              {Key: entity.ConfigKeyAutoInsertAd, Value: entity.ConfigDefaultAutoInsertAd, Type: entity.ConfigTypeString},
		entity.ConfigKeyPageSize:                  {Key: entity.ConfigKeyPageSize, Value: entity.ConfigDefaultPageSize, Type: entity.ConfigTypeInt},
		entity.ConfigKeyMaintenanceMode:           {Key: entity.ConfigKeyMaintenanceMode, Value: entity.ConfigDefaultMaintenanceMode, Type: entity.ConfigTypeBool},
		entity.ConfigKeyEnableRegister:            {Key: entity.ConfigKeyEnableRegister, Value: entity.ConfigDefaultEnableRegister, Type: entity.ConfigTypeBool},
		entity.ConfigKeyThirdPartyStatsCode:       {Key: entity.ConfigKeyThirdPartyStatsCode, Value: entity.ConfigDefaultThirdPartyStatsCode, Type: entity.ConfigTypeString},
		entity.ConfigKeyMeilisearchEnabled:        {Key: entity.ConfigKeyMeilisearchEnabled, Value: entity.ConfigDefaultMeilisearchEnabled, Type: entity.ConfigTypeBool},
		entity.ConfigKeyMeilisearchHost:           {Key: entity.ConfigKeyMeilisearchHost, Value: entity.ConfigDefaultMeilisearchHost, Type: entity.ConfigTypeString},
		entity.ConfigKeyMeilisearchPort:           {Key: entity.ConfigKeyMeilisearchPort, Value: entity.ConfigDefaultMeilisearchPort, Type: entity.ConfigTypeString},
		entity.ConfigKeyMeilisearchMasterKey:      {Key: entity.ConfigKeyMeilisearchMasterKey, Value: entity.ConfigDefaultMeilisearchMasterKey, Type: entity.ConfigTypeString},
		entity.ConfigKeyMeilisearchIndexName:      {Key: entity.ConfigKeyMeilisearchIndexName, Value: entity.ConfigDefaultMeilisearchIndexName, Type: entity.ConfigTypeString},
		entity.ConfigKeyEnableAnnouncements:       {Key: entity.ConfigKeyEnableAnnouncements, Value: entity.ConfigDefaultEnableAnnouncements, Type: entity.ConfigTypeBool},
		entity.ConfigKeyAnnouncements:             {Key: entity.ConfigKeyAnnouncements, Value: entity.ConfigDefaultAnnouncements, Type: entity.ConfigTypeJSON},
		entity.ConfigKeyEnableFloatButtons:        {Key: entity.ConfigKeyEnableFloatButtons, Value: entity.ConfigDefaultEnableFloatButtons, Type: entity.ConfigTypeBool},
		entity.ConfigKeyWechatSearchImage:         {Key: entity.ConfigKeyWechatSearchImage, Value: entity.ConfigDefaultWechatSearchImage, Type: entity.ConfigTypeString},
		entity.ConfigKeyTelegramQrImage:           {Key: entity.ConfigKeyTelegramQrImage, Value: entity.ConfigDefaultTelegramQrImage, Type: entity.ConfigTypeString},
	}

	// 检查现有配置中是否有缺失的配置项
	existingKeys := make(map[string]bool)
	for _, config := range configs {
		existingKeys[config.Key] = true
	}

	// 找出缺失的配置项
	var missingConfigs []entity.SystemConfig
	for key, requiredConfig := range requiredConfigs {
		if !existingKeys[key] {
			missingConfigs = append(missingConfigs, requiredConfig)
		}
	}

	// 如果有缺失的配置项，则添加它们
	if len(missingConfigs) > 0 {
		err = r.UpsertConfigs(missingConfigs)
		if err != nil {
			return nil, err
		}
		// 重新获取所有配置
		configs, err = r.FindAll()
		if err != nil {
			return nil, err
		}
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

// SafeRefreshConfigCache 安全的刷新配置缓存（带错误处理）
func (r *SystemConfigRepositoryImpl) SafeRefreshConfigCache() error {
	defer func() {
		if r := recover(); r != nil {
			utils.Error("配置缓存刷新时发生panic: %v", r)
		}
	}()

	r.refreshConfigCache()
	return nil
}

// ValidateConfigIntegrity 验证配置完整性
func (r *SystemConfigRepositoryImpl) ValidateConfigIntegrity() error {
	configs, err := r.FindAll()
	if err != nil {
		return fmt.Errorf("获取配置失败: %v", err)
	}

	// 检查关键配置是否存在
	requiredKeys := []string{
		entity.ConfigKeySiteTitle,
		entity.ConfigKeySiteDescription,
		entity.ConfigKeyKeywords,
		entity.ConfigKeyAuthor,
		entity.ConfigKeyCopyright,
		entity.ConfigKeyAutoProcessReadyResources,
		entity.ConfigKeyAutoProcessInterval,
		entity.ConfigKeyAutoTransferEnabled,
		entity.ConfigKeyAutoTransferLimitDays,
		entity.ConfigKeyAutoTransferMinSpace,
		entity.ConfigKeyAutoFetchHotDramaEnabled,
		entity.ConfigKeyApiToken,
		entity.ConfigKeyPageSize,
		entity.ConfigKeyMaintenanceMode,
		entity.ConfigKeyEnableRegister,
		entity.ConfigKeyThirdPartyStatsCode,
	}

	existingKeys := make(map[string]bool)
	for _, config := range configs {
		existingKeys[config.Key] = true
	}

	var missingKeys []string
	for _, key := range requiredKeys {
		if !existingKeys[key] {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		utils.Error("发现缺失的配置项: %v", missingKeys)
		return fmt.Errorf("配置不完整，缺失: %v", missingKeys)
	}

	utils.Info("配置完整性检查通过")
	return nil
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
