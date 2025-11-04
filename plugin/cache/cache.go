package cache

import (
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// CacheManager 缓存管理器
type CacheManager struct {
	cache    map[string]*CacheItem
	mutex    sync.RWMutex
	defaultTTL time.Duration
}

// NewCacheManager 创建新的缓存管理器
func NewCacheManager(defaultTTL time.Duration) *CacheManager {
	cm := &CacheManager{
		cache:      make(map[string]*CacheItem),
		defaultTTL: defaultTTL,
	}

	// 启动定期清理过期缓存的goroutine
	go cm.cleanupExpired()

	return cm
}

// Set 设置缓存项
func (cm *CacheManager) Set(key string, value interface{}, ttl time.Duration) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	expiration := time.Now().Add(ttl)
	cm.cache[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get 获取缓存项
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	item, exists := cm.cache[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Delete 删除缓存项
func (cm *CacheManager) Delete(key string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.cache, key)
}

// Clear 清空所有缓存
func (cm *CacheManager) Clear() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.cache = make(map[string]*CacheItem)
}

// cleanupExpired 定期清理过期缓存
func (cm *CacheManager) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		cm.mutex.Lock()
		now := time.Now()
		for key, item := range cm.cache {
			if now.After(item.Expiration) {
				delete(cm.cache, key)
				utils.Debug("Cleaned up expired cache item: %s", key)
			}
		}
		cm.mutex.Unlock()
	}
}

// GetStats 获取缓存统计信息
func (cm *CacheManager) GetStats() map[string]interface{} {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_items"] = len(cm.cache)

	// 计算过期项数量
	now := time.Now()
	expiredCount := 0
	for _, item := range cm.cache {
		if now.After(item.Expiration) {
			expiredCount++
		}
	}
	stats["expired_items"] = expiredCount

	return stats
}