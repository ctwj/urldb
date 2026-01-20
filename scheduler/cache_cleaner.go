package scheduler

import (
	"time"

	"github.com/ctwj/urldb/utils"
)

// CacheCleaner 缓存清理调度器
type CacheCleaner struct {
	baseScheduler *BaseScheduler
	running       bool
	ticker        *time.Ticker
	stopChan      chan bool
}

// NewCacheCleaner 创建缓存清理调度器
func NewCacheCleaner(baseScheduler *BaseScheduler) *CacheCleaner {
	return &CacheCleaner{
		baseScheduler: baseScheduler,
		running:       false,
		ticker:        time.NewTicker(time.Hour), // 每小时执行一次
		stopChan:      make(chan bool),
	}
}

// Start 启动缓存清理任务
func (cc *CacheCleaner) Start() {
	if cc.running {
		utils.Warn("缓存清理任务已在运行中")
		return
	}

	cc.running = true
	utils.Info("启动缓存清理任务")

	go func() {
		for {
			select {
			case <-cc.ticker.C:
				cc.cleanCache()
			case <-cc.stopChan:
				cc.running = false
				utils.Info("缓存清理任务已停止")
				return
			}
		}
	}()
}

// Stop 停止缓存清理任务
func (cc *CacheCleaner) Stop() {
	if !cc.running {
		return
	}

	close(cc.stopChan)
	cc.ticker.Stop()
}

// cleanCache 执行缓存清理
func (cc *CacheCleaner) cleanCache() {
	// utils.Debug("开始清理过期缓存")

	// 清理过期缓存（1小时TTL）
	utils.CleanAllExpiredCaches(time.Hour)
	// utils.Debug("定期清理过期缓存完成")

	// 可以在这里添加其他缓存清理逻辑，比如：
	// - 清理特定模式的缓存
	// - 记录缓存统计信息
	cc.logCacheStats()
}

// logCacheStats 记录缓存统计信息
func (cc *CacheCleaner) logCacheStats() {
	hotCacheSize := utils.GetHotResourcesCache().Size()
	relatedCacheSize := utils.GetRelatedResourcesCache().Size()
	systemConfigSize := utils.GetSystemConfigCache().Size()
	categoriesSize := utils.GetCategoriesCache().Size()
	tagsSize := utils.GetTagsCache().Size()

	totalSize := hotCacheSize + relatedCacheSize + systemConfigSize + categoriesSize + tagsSize

	utils.Debug("缓存统计 - 热门资源: %d, 相关资源: %d, 系统配置: %d, 分类: %d, 标签: %d, 总计: %d",
		hotCacheSize, relatedCacheSize, systemConfigSize, categoriesSize, tagsSize, totalSize)

	// 如果缓存过多，可以记录警告
	if totalSize > 1000 {
		utils.Warn("缓存项数量过多: %d，建议检查缓存策略", totalSize)
	}
}

// IsRunning 检查是否正在运行
func (cc *CacheCleaner) IsRunning() bool {
	return cc.running
}
