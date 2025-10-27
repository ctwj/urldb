package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// ConfigManager 统一配置管理器
type ConfigManager struct {
	repo *repo.RepositoryManager

	// 内存缓存
	cache      map[string]*ConfigItem
	cacheMutex sync.RWMutex
	cacheOnce  sync.Once

	// 配置更新通知
	configUpdateCh chan string
	watchers       []chan string
	watcherMutex   sync.Mutex

	// 加载时间
	lastLoadTime time.Time
}

// ConfigItem 配置项结构
type ConfigItem struct {
	Key        string      `json:"key"`
	Value      string      `json:"value"`
	Type       string      `json:"type"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Group      string      `json:"group"`    // 配置分组
	Category   string      `json:"category"` // 配置分类
	IsSensitive bool      `json:"is_sensitive"` // 是否是敏感信息
}

// ConfigGroup 配置分组
type ConfigGroup string

const (
	GroupDatabase   ConfigGroup = "database"
	GroupServer     ConfigGroup = "server"
	GroupSecurity   ConfigGroup = "security"
	GroupSearch     ConfigGroup = "search"
	GroupTelegram   ConfigGroup = "telegram"
	GroupCache      ConfigGroup = "cache"
	GroupMeilisearch ConfigGroup = "meilisearch"
	GroupSEO        ConfigGroup = "seo"
	GroupAutoProcess ConfigGroup = "auto_process"
	GroupOther      ConfigGroup = "other"
)

// NewConfigManager 创建配置管理器
func NewConfigManager(repoManager *repo.RepositoryManager) *ConfigManager {
	cm := &ConfigManager{
		repo:           repoManager,
		cache:          make(map[string]*ConfigItem),
		configUpdateCh: make(chan string, 100), // 缓冲通道防止阻塞
	}

	// 启动配置更新监听器
	go cm.startConfigUpdateListener()

	return cm
}

// startConfigUpdateListener 启动配置更新监听器
func (cm *ConfigManager) startConfigUpdateListener() {
	for key := range cm.configUpdateCh {
		cm.notifyWatchers(key)
	}
}

// notifyWatchers 通知所有监听器配置已更新
func (cm *ConfigManager) notifyWatchers(key string) {
	cm.watcherMutex.Lock()
	defer cm.watcherMutex.Unlock()

	for _, watcher := range cm.watchers {
		select {
		case watcher <- key:
		default:
			// 如果通道阻塞，跳过该监听器
			utils.Warn("配置监听器通道阻塞，跳过通知: %s", key)
		}
	}
}

// AddConfigWatcher 添加配置变更监听器
func (cm *ConfigManager) AddConfigWatcher() chan string {
	cm.watcherMutex.Lock()
	defer cm.watcherMutex.Unlock()

	watcher := make(chan string, 10) // 为每个监听器创建缓冲通道
	cm.watchers = append(cm.watchers, watcher)
	return watcher
}

// GetConfig 获取配置项
func (cm *ConfigManager) GetConfig(key string) (*ConfigItem, error) {
	// 先尝试从内存缓存获取
	item, exists := cm.getCachedConfig(key)
	if exists {
		return item, nil
	}

	// 如果缓存中没有，从数据库获取
	config, err := cm.repo.SystemConfigRepository.FindByKey(key)
	if err != nil {
		return nil, err
	}

	// 将数据库配置转换为ConfigItem并缓存
	item = &ConfigItem{
		Key:       config.Key,
		Value:     config.Value,
		Type:      config.Type,
		UpdatedAt: time.Now(),
	}

	if group := cm.getGroupByConfigKey(key); group != "" {
		item.Group = string(group)
	}

	if category := cm.getCategoryByConfigKey(key); category != "" {
		item.Category = category
	}

	item.IsSensitive = cm.isSensitiveConfig(key)

	// 缓存配置
	cm.setCachedConfig(key, item)

	return item, nil
}

// GetConfigValue 获取配置值
func (cm *ConfigManager) GetConfigValue(key string) (string, error) {
	item, err := cm.GetConfig(key)
	if err != nil {
		return "", err
	}
	return item.Value, nil
}

// GetConfigBool 获取布尔值配置
func (cm *ConfigManager) GetConfigBool(key string) (bool, error) {
	value, err := cm.GetConfigValue(key)
	if err != nil {
		return false, err
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off", "":
		return false, nil
	default:
		return false, fmt.Errorf("无法将配置值 '%s' 转换为布尔值", value)
	}
}

// GetConfigInt 获取整数值配置
func (cm *ConfigManager) GetConfigInt(key string) (int, error) {
	value, err := cm.GetConfigValue(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

// GetConfigInt64 获取64位整数值配置
func (cm *ConfigManager) GetConfigInt64(key string) (int64, error) {
	value, err := cm.GetConfigValue(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(value, 10, 64)
}

// GetConfigFloat64 获取浮点数配置
func (cm *ConfigManager) GetConfigFloat64(key string) (float64, error) {
	value, err := cm.GetConfigValue(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(value, 64)
}

// SetConfig 设置配置值
func (cm *ConfigManager) SetConfig(key, value string) error {
	// 更新数据库
	config := &entity.SystemConfig{
		Key:   key,
		Value: value,
		Type:  "string", // 默认类型，实际类型应该从现有配置中获取
	}

	// 获取现有配置以确定类型
	existing, err := cm.repo.SystemConfigRepository.FindByKey(key)
	if err == nil {
		config.Type = existing.Type
	} else {
		// 如果配置不存在，尝试从默认配置中获取类型
		config.Type = cm.getDefaultConfigType(key)
	}

	// 保存到数据库
	err = cm.repo.SystemConfigRepository.UpsertConfigs([]entity.SystemConfig{*config})
	if err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	// 更新缓存
	item := &ConfigItem{
		Key:       config.Key,
		Value:     config.Value,
		Type:      config.Type,
		UpdatedAt: time.Now(),
	}

	if group := cm.getGroupByConfigKey(key); group != "" {
		item.Group = string(group)
	}

	if category := cm.getCategoryByConfigKey(key); category != "" {
		item.Category = category
	}

	item.IsSensitive = cm.isSensitiveConfig(key)

	cm.setCachedConfig(key, item)

	// 发送更新通知
	cm.configUpdateCh <- key

	utils.Info("配置已更新: %s = %s", key, value)

	return nil
}

// SetConfigWithType 设置配置值（指定类型）
func (cm *ConfigManager) SetConfigWithType(key, value, configType string) error {
	config := &entity.SystemConfig{
		Key:   key,
		Value: value,
		Type:  configType,
	}

	err := cm.repo.SystemConfigRepository.UpsertConfigs([]entity.SystemConfig{*config})
	if err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	// 更新缓存
	item := &ConfigItem{
		Key:       config.Key,
		Value:     config.Value,
		Type:      config.Type,
		UpdatedAt: time.Now(),
	}

	if group := cm.getGroupByConfigKey(key); group != "" {
		item.Group = string(group)
	}

	if category := cm.getCategoryByConfigKey(key); category != "" {
		item.Category = category
	}

	item.IsSensitive = cm.isSensitiveConfig(key)

	cm.setCachedConfig(key, item)

	// 发送更新通知
	cm.configUpdateCh <- key

	utils.Info("配置已更新: %s = %s (type: %s)", key, value, configType)

	return nil
}

// getGroupByConfigKey 根据配置键获取分组
func (cm *ConfigManager) getGroupByConfigKey(key string) ConfigGroup {
	switch {
	case strings.HasPrefix(key, "database_"), strings.HasPrefix(key, "db_"):
		return GroupDatabase
	case strings.HasPrefix(key, "server_"), strings.HasPrefix(key, "port"), strings.HasPrefix(key, "host"):
		return GroupServer
	case strings.HasPrefix(key, "api_"), strings.HasPrefix(key, "jwt_"), strings.HasPrefix(key, "password"):
		return GroupSecurity
	case strings.Contains(key, "meilisearch"):
		return GroupMeilisearch
	case strings.Contains(key, "telegram"):
		return GroupTelegram
	case strings.Contains(key, "cache"), strings.Contains(key, "redis"):
		return GroupCache
	case strings.Contains(key, "seo"), strings.Contains(key, "title"), strings.Contains(key, "keyword"):
		return GroupSEO
	case strings.Contains(key, "auto_"):
		return GroupAutoProcess
	case strings.Contains(key, "forbidden"), strings.Contains(key, "ad_"):
		return GroupOther
	default:
		return GroupOther
	}
}

// getCategoryByConfigKey 根据配置键获取分类
func (cm *ConfigManager) getCategoryByConfigKey(key string) string {
	switch {
	case key == entity.ConfigKeySiteTitle || key == entity.ConfigKeySiteDescription:
		return "basic_info"
	case key == entity.ConfigKeyKeywords || key == entity.ConfigKeyAuthor:
		return "seo"
	case key == entity.ConfigKeyAutoProcessReadyResources || key == entity.ConfigKeyAutoProcessInterval:
		return "auto_process"
	case key == entity.ConfigKeyAutoTransferEnabled || key == entity.ConfigKeyAutoTransferLimitDays:
		return "auto_transfer"
	case key == entity.ConfigKeyMeilisearchEnabled || key == entity.ConfigKeyMeilisearchHost:
		return "search"
	case key == entity.ConfigKeyTelegramBotEnabled || key == entity.ConfigKeyTelegramBotApiKey:
		return "telegram"
	case key == entity.ConfigKeyMaintenanceMode || key == entity.ConfigKeyEnableRegister:
		return "system"
	case key == entity.ConfigKeyForbiddenWords || key == entity.ConfigKeyAdKeywords:
		return "filtering"
	default:
		return "other"
	}
}

// isSensitiveConfig 判断是否是敏感配置
func (cm *ConfigManager) isSensitiveConfig(key string) bool {
	switch key {
	case entity.ConfigKeyApiToken,
		entity.ConfigKeyMeilisearchMasterKey,
		entity.ConfigKeyTelegramBotApiKey,
		entity.ConfigKeyTelegramProxyUsername,
		entity.ConfigKeyTelegramProxyPassword:
		return true
	default:
		return strings.Contains(strings.ToLower(key), "password") ||
			strings.Contains(strings.ToLower(key), "secret") ||
			strings.Contains(strings.ToLower(key), "key") ||
			strings.Contains(strings.ToLower(key), "token")
	}
}

// getDefaultConfigType 获取默认配置类型
func (cm *ConfigManager) getDefaultConfigType(key string) string {
	switch key {
	case entity.ConfigKeyAutoProcessReadyResources,
		entity.ConfigKeyAutoTransferEnabled,
		entity.ConfigKeyAutoFetchHotDramaEnabled,
		entity.ConfigKeyMaintenanceMode,
		entity.ConfigKeyEnableRegister,
		entity.ConfigKeyMeilisearchEnabled,
		entity.ConfigKeyTelegramBotEnabled:
		return entity.ConfigTypeBool
	case entity.ConfigKeyAutoProcessInterval,
		entity.ConfigKeyAutoTransferLimitDays,
		entity.ConfigKeyAutoTransferMinSpace,
		entity.ConfigKeyPageSize:
		return entity.ConfigTypeInt
	case entity.ConfigKeyAnnouncements:
		return entity.ConfigTypeJSON
	default:
		return entity.ConfigTypeString
	}
}

// LoadAllConfigs 加载所有配置到缓存
func (cm *ConfigManager) LoadAllConfigs() error {
	configs, err := cm.repo.SystemConfigRepository.FindAll()
	if err != nil {
		return fmt.Errorf("加载所有配置失败: %v", err)
	}

	cm.cacheMutex.Lock()
	defer cm.cacheMutex.Unlock()

	// 清空现有缓存
	cm.cache = make(map[string]*ConfigItem)

	// 更新缓存
	for _, config := range configs {
		item := &ConfigItem{
			Key:       config.Key,
			Value:     config.Value,
			Type:      config.Type,
			UpdatedAt: time.Now(), // 实际应该从数据库获取
		}

		if group := cm.getGroupByConfigKey(config.Key); group != "" {
			item.Group = string(group)
		}

		if category := cm.getCategoryByConfigKey(config.Key); category != "" {
			item.Category = category
		}

		item.IsSensitive = cm.isSensitiveConfig(config.Key)

		cm.cache[config.Key] = item
	}

	cm.lastLoadTime = time.Now()

	utils.Info("已加载 %d 个配置项到缓存", len(configs))

	return nil
}

// RefreshConfigCache 刷新配置缓存
func (cm *ConfigManager) RefreshConfigCache() error {
	return cm.LoadAllConfigs()
}

// GetCachedConfig 获取缓存的配置
func (cm *ConfigManager) getCachedConfig(key string) (*ConfigItem, bool) {
	cm.cacheMutex.RLock()
	defer cm.cacheMutex.RUnlock()

	item, exists := cm.cache[key]
	return item, exists
}

// setCachedConfig 设置缓存的配置
func (cm *ConfigManager) setCachedConfig(key string, item *ConfigItem) {
	cm.cacheMutex.Lock()
	defer cm.cacheMutex.Unlock()

	cm.cache[key] = item
}

// GetConfigByGroup 按分组获取配置
func (cm *ConfigManager) GetConfigByGroup(group ConfigGroup) (map[string]*ConfigItem, error) {
	cm.cacheMutex.RLock()
	defer cm.cacheMutex.RUnlock()

	result := make(map[string]*ConfigItem)

	for key, item := range cm.cache {
		if ConfigGroup(item.Group) == group {
			result[key] = item
		}
	}

	return result, nil
}

// GetConfigByCategory 按分类获取配置
func (cm *ConfigManager) GetConfigByCategory(category string) (map[string]*ConfigItem, error) {
	cm.cacheMutex.RLock()
	defer cm.cacheMutex.RUnlock()

	result := make(map[string]*ConfigItem)

	for key, item := range cm.cache {
		if item.Category == category {
			result[key] = item
		}
	}

	return result, nil
}

// DeleteConfig 删除配置
func (cm *ConfigManager) DeleteConfig(key string) error {
	// 先查找配置获取ID
	config, err := cm.repo.SystemConfigRepository.FindByKey(key)
	if err != nil {
		return fmt.Errorf("查找配置失败: %v", err)
	}

	// 从数据库删除
	err = cm.repo.SystemConfigRepository.Delete(config.ID)
	if err != nil {
		return fmt.Errorf("删除配置失败: %v", err)
	}

	// 从缓存中移除
	cm.cacheMutex.Lock()
	delete(cm.cache, key)
	cm.cacheMutex.Unlock()

	utils.Info("配置已删除: %s", key)

	return nil
}

// GetSensitiveConfigKeys 获取所有敏感配置键
func (cm *ConfigManager) GetSensitiveConfigKeys() []string {
	cm.cacheMutex.RLock()
	defer cm.cacheMutex.RUnlock()

	var sensitiveKeys []string
	for key, item := range cm.cache {
		if item.IsSensitive {
			sensitiveKeys = append(sensitiveKeys, key)
		}
	}

	return sensitiveKeys
}

// GetConfigWithMask 获取配置值（敏感配置会被遮蔽）
func (cm *ConfigManager) GetConfigWithMask(key string) (*ConfigItem, error) {
	item, err := cm.GetConfig(key)
	if err != nil {
		return nil, err
	}

	if item.IsSensitive {
		// 创建副本并遮蔽敏感值
		maskedItem := *item
		maskedItem.Value = cm.maskSensitiveValue(item.Value)
		return &maskedItem, nil
	}

	return item, nil
}

// maskSensitiveValue 遮蔽敏感值
func (cm *ConfigManager) maskSensitiveValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}

	// 保留前2个和后2个字符，中间用****替代
	return value[:2] + "****" + value[len(value)-2:]
}

// GetConfigAsJSON 获取配置为JSON格式
func (cm *ConfigManager) GetConfigAsJSON() ([]byte, error) {
	cm.cacheMutex.RLock()
	defer cm.cacheMutex.RUnlock()

	// 创建副本，敏感配置使用遮蔽值
	configMap := make(map[string]*ConfigItem)
	for key, item := range cm.cache {
		if item.IsSensitive {
			maskedItem := *item
			maskedItem.Value = cm.maskSensitiveValue(item.Value)
			configMap[key] = &maskedItem
		} else {
			configMap[key] = item
		}
	}

	return json.MarshalIndent(configMap, "", "  ")
}

// GetConfigStatistics 获取配置统计信息
func (cm *ConfigManager) GetConfigStatistics() map[string]interface{} {
	cm.cacheMutex.RLock()
	defer cm.cacheMutex.RUnlock()

	stats := map[string]interface{}{
		"total_configs":          len(cm.cache),
		"last_load_time":         cm.lastLoadTime,
		"cache_size_bytes":       len(cm.cache) * 100, // 估算每个配置约100字节
		"groups":                 make(map[string]int),
		"types":                  make(map[string]int),
		"categories":             make(map[string]int),
		"sensitive_configs":      0,
		"config_keys":            make([]string, 0),
	}

	groups := make(map[string]int)
	types := make(map[string]int)
	categories := make(map[string]int)

	for key, item := range cm.cache {
		// 统计分组
		groups[item.Group]++

		// 统计类型
		types[item.Type]++

		// 统计分类
		categories[item.Category]++

		// 统计敏感配置
		if item.IsSensitive {
			stats["sensitive_configs"] = stats["sensitive_configs"].(int) + 1
		}

		// 添加配置键到列表
		keys := stats["config_keys"].([]string)
		keys = append(keys, key)
		stats["config_keys"] = keys
	}

	stats["groups"] = groups
	stats["types"] = types
	stats["categories"] = categories

	return stats
}

// GetEnvironmentConfig 从环境变量获取配置
func (cm *ConfigManager) GetEnvironmentConfig(key string) (string, bool) {
	value := os.Getenv(key)
	if value != "" {
		return value, true
	}

	// 尝试使用大写版本的键
	value = os.Getenv(strings.ToUpper(key))
	if value != "" {
		return value, true
	}

	// 尝试使用大写带下划线的格式
	upperKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	value = os.Getenv(upperKey)
	if value != "" {
		return value, true
	}

	return "", false
}

// GetConfigWithEnvFallback 获取配置，环境变量优先
func (cm *ConfigManager) GetConfigWithEnvFallback(configKey, envKey string) (string, error) {
	// 优先从环境变量获取
	if envValue, exists := cm.GetEnvironmentConfig(envKey); exists {
		return envValue, nil
	}

	// 如果环境变量不存在，从数据库获取
	return cm.GetConfigValue(configKey)
}

// GetConfigIntWithEnvFallback 获取整数配置，环境变量优先
func (cm *ConfigManager) GetConfigIntWithEnvFallback(configKey, envKey string) (int, error) {
	// 优先从环境变量获取
	if envValue, exists := cm.GetEnvironmentConfig(envKey); exists {
		return strconv.Atoi(envValue)
	}

	// 如果环境变量不存在，从数据库获取
	return cm.GetConfigInt(configKey)
}

// GetConfigBoolWithEnvFallback 获取布尔配置，环境变量优先
func (cm *ConfigManager) GetConfigBoolWithEnvFallback(configKey, envKey string) (bool, error) {
	// 优先从环境变量获取
	if envValue, exists := cm.GetEnvironmentConfig(envKey); exists {
		switch strings.ToLower(envValue) {
		case "true", "1", "yes", "on":
			return true, nil
		case "false", "0", "no", "off", "":
			return false, nil
		default:
			return false, fmt.Errorf("无法将环境变量值 '%s' 转换为布尔值", envValue)
		}
	}

	// 如果环境变量不存在，从数据库获取
	return cm.GetConfigBool(configKey)
}