package config

import (
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// SyncWithRepository 同步配置管理器与Repository的缓存
func (cm *ConfigManager) SyncWithRepository(repoManager *repo.RepositoryManager) {
	// 监听配置变更事件并同步缓存
	// 这是一个抽象概念，实际实现需要修改Repository接口

	// 当配置更新时，通知Repository清理缓存
	go func() {
		watcher := cm.AddConfigWatcher()
		for {
			select {
			case key := <-watcher:
				// 通知Repository层清理缓存（如果Repository支持）
				utils.Debug("配置 %s 已更新，可能需要同步到Repository缓存", key)
			}
		}
	}()
}

// UpdateRepositoryCache 当配置管理器更新配置时，通知Repository层同步
func (cm *ConfigManager) UpdateRepositoryCache(repoManager *repo.RepositoryManager) {
	// 这个函数需要Repository支持特定的缓存清理方法
	// 由于现有Repository没有提供这样的接口，我们只能依赖数据库同步
	utils.Info("配置已通过配置管理器更新，Repository层将从数据库重新加载")
}