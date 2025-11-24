package scheduler

import (
	"sync"

	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"
)

// GlobalScheduler 全局调度器管理器
type GlobalScheduler struct {
	manager *Manager
	mutex   sync.RWMutex
}

var (
	globalScheduler *GlobalScheduler
	once            sync.Once
	// 全局Meilisearch管理器
	globalMeilisearchManager *services.MeilisearchManager
)

// SetGlobalMeilisearchManager 设置全局Meilisearch管理器
func SetGlobalMeilisearchManager(manager *services.MeilisearchManager) {
	globalMeilisearchManager = manager
}

// GetGlobalMeilisearchManager 获取全局Meilisearch管理器
func GetGlobalMeilisearchManager() *services.MeilisearchManager {
	return globalMeilisearchManager
}

// GetGlobalScheduler 获取全局调度器实例（单例模式）
func GetGlobalScheduler(hotDramaRepo repo.HotDramaRepository, readyResourceRepo repo.ReadyResourceRepository, resourceRepo repo.ResourceRepository, systemConfigRepo repo.SystemConfigRepository, panRepo repo.PanRepository, cksRepo repo.CksRepository, tagRepo repo.TagRepository, categoryRepo repo.CategoryRepository) *GlobalScheduler {
	once.Do(func() {
		globalScheduler = &GlobalScheduler{
			manager: NewManager(hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo, panRepo, cksRepo, tagRepo, categoryRepo),
		}
	})
	return globalScheduler
}

// StartHotDramaScheduler 启动热播剧定时任务
func (gs *GlobalScheduler) StartHotDramaScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.manager.IsHotDramaRunning() {
		utils.Debug("热播剧定时任务已在运行中")
		return
	}

	gs.manager.StartHotDramaScheduler()
	utils.Debug("全局调度器已启动热播剧定时任务")
}

// StopHotDramaScheduler 停止热播剧定时任务
func (gs *GlobalScheduler) StopHotDramaScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.manager.IsHotDramaRunning() {
		utils.Debug("热播剧定时任务未在运行")
		return
	}

	gs.manager.StopHotDramaScheduler()
	utils.Debug("全局调度器已停止热播剧定时任务")
}

// IsHotDramaSchedulerRunning 检查热播剧定时任务是否在运行
func (gs *GlobalScheduler) IsHotDramaSchedulerRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.manager.IsHotDramaRunning()
}

// GetHotDramaNames 手动获取热播剧名字
func (gs *GlobalScheduler) GetHotDramaNames() ([]string, error) {
	return gs.manager.GetHotDramaNames()
}

// StartReadyResourceScheduler 启动待处理资源自动处理任务
func (gs *GlobalScheduler) StartReadyResourceScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.manager.IsReadyResourceRunning() {
		utils.Debug("待处理资源自动处理任务已在运行中")
		return
	}

	gs.manager.StartReadyResourceScheduler()
	utils.Debug("全局调度器已启动待处理资源自动处理任务")
}

// StopReadyResourceScheduler 停止待处理资源自动处理任务
func (gs *GlobalScheduler) StopReadyResourceScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.manager.IsReadyResourceRunning() {
		utils.Debug("待处理资源自动处理任务未在运行")
		return
	}

	gs.manager.StopReadyResourceScheduler()
	utils.Debug("全局调度器已停止待处理资源自动处理任务")
}

// IsReadyResourceRunning 检查待处理资源自动处理任务是否在运行
func (gs *GlobalScheduler) IsReadyResourceRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.manager.IsReadyResourceRunning()
}

// UpdateSchedulerStatusWithAutoTransfer 根据系统配置更新调度器状态（包含自动转存）
func (gs *GlobalScheduler) UpdateSchedulerStatusWithAutoTransfer(autoFetchHotDramaEnabled bool, autoProcessReadyResources bool, autoTransferEnabled bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// 处理热播剧自动拉取功能
	if autoFetchHotDramaEnabled {
		if !gs.manager.IsHotDramaRunning() {
			utils.Info("系统配置启用自动拉取热播剧，启动定时任务")
			gs.manager.StartHotDramaScheduler()
		}
	} else {
		if gs.manager.IsHotDramaRunning() {
			utils.Info("系统配置禁用自动拉取热播剧，停止定时任务")
			gs.manager.StopHotDramaScheduler()
		}
	}

	// 处理待处理资源自动处理功能
	if autoProcessReadyResources {
		if !gs.manager.IsReadyResourceRunning() {
			utils.Info("系统配置启用自动处理待处理资源，启动定时任务")
			gs.manager.StartReadyResourceScheduler()
		}
	} else {
		if gs.manager.IsReadyResourceRunning() {
			utils.Info("系统配置禁用自动处理待处理资源，停止定时任务")
			gs.manager.StopReadyResourceScheduler()
		}
	}

}

// StartSitemapScheduler 启动Sitemap调度任务
func (gs *GlobalScheduler) StartSitemapScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.manager.IsSitemapRunning() {
		utils.Debug("Sitemap定时任务已在运行中")
		return
	}

	gs.manager.StartSitemapScheduler()
	utils.Debug("全局调度器已启动Sitemap定时任务")
}

// StopSitemapScheduler 停止Sitemap调度任务
func (gs *GlobalScheduler) StopSitemapScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.manager.IsSitemapRunning() {
		utils.Debug("Sitemap定时任务未在运行")
		return
	}

	gs.manager.StopSitemapScheduler()
	utils.Debug("全局调度器已停止Sitemap定时任务")
}

// IsSitemapSchedulerRunning 检查Sitemap定时任务是否在运行
func (gs *GlobalScheduler) IsSitemapSchedulerRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.manager.IsSitemapRunning()
}

// UpdateSitemapConfig 更新Sitemap配置
func (gs *GlobalScheduler) UpdateSitemapConfig(enabled bool) error {
	return gs.manager.UpdateSitemapConfig(enabled)
}

// GetSitemapConfig 获取Sitemap配置
func (gs *GlobalScheduler) GetSitemapConfig() (bool, error) {
	return gs.manager.GetSitemapConfig()
}

// TriggerSitemapGeneration 手动触发sitemap生成
func (gs *GlobalScheduler) TriggerSitemapGeneration() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.manager.TriggerSitemapGeneration()
}

// StartGoogleIndexScheduler 启动Google索引调度任务
func (gs *GlobalScheduler) StartGoogleIndexScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.manager.IsGoogleIndexRunning() {
		utils.Debug("Google索引调度任务已在运行中")
		return
	}

	gs.manager.StartGoogleIndexScheduler()
	utils.Info("Google索引调度任务已启动")
}

// StopGoogleIndexScheduler 停止Google索引调度任务
func (gs *GlobalScheduler) StopGoogleIndexScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.manager.IsGoogleIndexRunning() {
		utils.Debug("Google索引调度任务未在运行")
		return
	}

	gs.manager.StopGoogleIndexScheduler()
	utils.Info("Google索引调度任务已停止")
}

// IsGoogleIndexSchedulerRunning 检查Google索引调度任务是否在运行
func (gs *GlobalScheduler) IsGoogleIndexSchedulerRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	return gs.manager.IsGoogleIndexRunning()
}
