package utils

import (
	"sync"
	"time"
)

// CacheData 缓存数据结构
type CacheData struct {
	Data      interface{}
	UpdatedAt time.Time
}

// CacheManager 通用缓存管理器
type CacheManager struct {
	cache map[string]*CacheData
	mutex sync.RWMutex
}

// NewCacheManager 创建缓存管理器
func NewCacheManager() *CacheManager {
	return &CacheManager{
		cache: make(map[string]*CacheData),
	}
}

// Set 设置缓存
func (cm *CacheManager) Set(key string, data interface{}) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.cache[key] = &CacheData{
		Data:      data,
		UpdatedAt: time.Now(),
	}
}

// Get 获取缓存
func (cm *CacheManager) Get(key string, ttl time.Duration) (interface{}, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if cachedData, exists := cm.cache[key]; exists {
		if time.Since(cachedData.UpdatedAt) < ttl {
			return cachedData.Data, true
		}
		// 缓存过期，删除
		delete(cm.cache, key)
	}
	return nil, false
}

// GetWithTTL 获取缓存并返回剩余TTL
func (cm *CacheManager) GetWithTTL(key string, ttl time.Duration) (interface{}, bool, time.Duration) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if cachedData, exists := cm.cache[key]; exists {
		elapsed := time.Since(cachedData.UpdatedAt)
		if elapsed < ttl {
			return cachedData.Data, true, ttl - elapsed
		}
		// 缓存过期，删除
		delete(cm.cache, key)
	}
	return nil, false, 0
}

// Delete 删除缓存
func (cm *CacheManager) Delete(key string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.cache, key)
}

// DeletePattern 删除匹配模式的缓存
func (cm *CacheManager) DeletePattern(pattern string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for key := range cm.cache {
		// 简单的字符串匹配，可以根据需要扩展为正则表达式
		if len(pattern) > 0 && (key == pattern || (len(key) >= len(pattern) && key[:len(pattern)] == pattern)) {
			delete(cm.cache, key)
		}
	}
}

// Clear 清空所有缓存
func (cm *CacheManager) Clear() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.cache = make(map[string]*CacheData)
}

// Size 获取缓存项数量
func (cm *CacheManager) Size() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.cache)
}

// CleanExpired 清理过期缓存
func (cm *CacheManager) CleanExpired(ttl time.Duration) int {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cleaned := 0
	now := time.Now()
	for key, cachedData := range cm.cache {
		if now.Sub(cachedData.UpdatedAt) >= ttl {
			delete(cm.cache, key)
			cleaned++
		}
	}
	return cleaned
}

// GetKeys 获取所有缓存键
func (cm *CacheManager) GetKeys() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	keys := make([]string, 0, len(cm.cache))
	for key := range cm.cache {
		keys = append(keys, key)
	}
	return keys
}

// 全局缓存管理器实例
var (
	// 热门资源缓存
	HotResourcesCache = NewCacheManager()

	// 相关资源缓存
	RelatedResourcesCache = NewCacheManager()

	// 系统配置缓存
	SystemConfigCache = NewCacheManager()

	// 分类缓存
	CategoriesCache = NewCacheManager()

	// 标签缓存
	TagsCache = NewCacheManager()

	// 资源有效性检测缓存
	ResourceValidityCache = NewCacheManager()
)

// GetHotResourcesCache 获取热门资源缓存管理器
func GetHotResourcesCache() *CacheManager {
	return HotResourcesCache
}

// GetRelatedResourcesCache 获取相关资源缓存管理器
func GetRelatedResourcesCache() *CacheManager {
	return RelatedResourcesCache
}

// GetSystemConfigCache 获取系统配置缓存管理器
func GetSystemConfigCache() *CacheManager {
	return SystemConfigCache
}

// GetCategoriesCache 获取分类缓存管理器
func GetCategoriesCache() *CacheManager {
	return CategoriesCache
}

// GetTagsCache 获取标签缓存管理器
func GetTagsCache() *CacheManager {
	return TagsCache
}

// GetResourceValidityCache 获取资源有效性检测缓存管理器
func GetResourceValidityCache() *CacheManager {
	return ResourceValidityCache
}

// ClearAllCaches 清空所有全局缓存
func ClearAllCaches() {
	HotResourcesCache.Clear()
	RelatedResourcesCache.Clear()
	SystemConfigCache.Clear()
	CategoriesCache.Clear()
	TagsCache.Clear()
	ResourceValidityCache.Clear()
}

// CleanAllExpiredCaches 清理所有过期缓存
func CleanAllExpiredCaches(ttl time.Duration) {
	totalCleaned := 0
	totalCleaned += HotResourcesCache.CleanExpired(ttl)
	totalCleaned += RelatedResourcesCache.CleanExpired(ttl)
	totalCleaned += SystemConfigCache.CleanExpired(ttl)
	totalCleaned += CategoriesCache.CleanExpired(ttl)
	totalCleaned += TagsCache.CleanExpired(ttl)
	totalCleaned += ResourceValidityCache.CleanExpired(ttl)

	if totalCleaned > 0 {
		Info("清理过期缓存完成，共清理 %d 个缓存项", totalCleaned)
	}
}