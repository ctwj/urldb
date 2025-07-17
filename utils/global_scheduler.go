package utils

import (
	"sync"

	"github.com/ctwj/panResManage/db/repo"
)

// GlobalScheduler 全局调度器管理器
type GlobalScheduler struct {
	scheduler *Scheduler
	mutex     sync.RWMutex
}

var (
	globalScheduler *GlobalScheduler
	once            sync.Once
)

// GetGlobalScheduler 获取全局调度器实例（单例模式）
func GetGlobalScheduler(hotDramaRepo repo.HotDramaRepository, readyResourceRepo repo.ReadyResourceRepository, resourceRepo repo.ResourceRepository, systemConfigRepo repo.SystemConfigRepository, panRepo repo.PanRepository, cksRepo repo.CksRepository) *GlobalScheduler {
	once.Do(func() {
		globalScheduler = &GlobalScheduler{
			scheduler: NewScheduler(hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo, panRepo, cksRepo),
		}
	})
	return globalScheduler
}

// StartHotDramaScheduler 启动热播剧定时任务
func (gs *GlobalScheduler) StartHotDramaScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.scheduler.IsRunning() {
		Info("热播剧定时任务已在运行中")
		return
	}

	gs.scheduler.StartHotDramaScheduler()
	Info("全局调度器已启动热播剧定时任务")
}

// StopHotDramaScheduler 停止热播剧定时任务
func (gs *GlobalScheduler) StopHotDramaScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.scheduler.IsRunning() {
		Info("热播剧定时任务未在运行")
		return
	}

	gs.scheduler.StopHotDramaScheduler()
	Info("全局调度器已停止热播剧定时任务")
}

// IsHotDramaSchedulerRunning 检查热播剧定时任务是否在运行
func (gs *GlobalScheduler) IsHotDramaSchedulerRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.scheduler.IsRunning()
}

// GetHotDramaNames 手动获取热播剧名字
func (gs *GlobalScheduler) GetHotDramaNames() ([]string, error) {
	return gs.scheduler.GetHotDramaNames()
}

// StartReadyResourceScheduler 启动待处理资源自动处理任务
func (gs *GlobalScheduler) StartReadyResourceScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.scheduler.IsReadyResourceRunning() {
		Info("待处理资源自动处理任务已在运行中")
		return
	}

	gs.scheduler.StartReadyResourceScheduler()
	Info("全局调度器已启动待处理资源自动处理任务")
}

// StopReadyResourceScheduler 停止待处理资源自动处理任务
func (gs *GlobalScheduler) StopReadyResourceScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.scheduler.IsReadyResourceRunning() {
		Info("待处理资源自动处理任务未在运行")
		return
	}

	gs.scheduler.StopReadyResourceScheduler()
	Info("全局调度器已停止待处理资源自动处理任务")
}

// IsReadyResourceRunning 检查待处理资源自动处理任务是否在运行
func (gs *GlobalScheduler) IsReadyResourceRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.scheduler.IsReadyResourceRunning()
}

// ProcessReadyResources 手动触发待处理资源处理
func (gs *GlobalScheduler) ProcessReadyResources() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.scheduler.processReadyResources()
}

// UpdateSchedulerStatus 根据系统配置更新调度器状态
func (gs *GlobalScheduler) UpdateSchedulerStatus(autoFetchHotDramaEnabled bool, autoProcessReadyResources bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// 处理热播剧自动拉取功能
	if autoFetchHotDramaEnabled {
		if !gs.scheduler.IsRunning() {
			Info("系统配置启用自动拉取热播剧，启动定时任务")
			gs.scheduler.StartHotDramaScheduler()
		}
	} else {
		if gs.scheduler.IsRunning() {
			Info("系统配置禁用自动拉取热播剧，停止定时任务")
			gs.scheduler.StopHotDramaScheduler()
		}
	}

	// 处理待处理资源自动处理功能
	if autoProcessReadyResources {
		if !gs.scheduler.IsReadyResourceRunning() {
			Info("系统配置启用自动处理待处理资源，启动定时任务")
			gs.scheduler.StartReadyResourceScheduler()
		}
	} else {
		if gs.scheduler.IsReadyResourceRunning() {
			Info("系统配置禁用自动处理待处理资源，停止定时任务")
			gs.scheduler.StopReadyResourceScheduler()
		}
	}
}

// StartAutoTransferScheduler 启动自动转存定时任务
func (gs *GlobalScheduler) StartAutoTransferScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.scheduler.IsAutoTransferRunning() {
		Info("自动转存定时任务已在运行中")
		return
	}

	gs.scheduler.StartAutoTransferScheduler()
	Info("全局调度器已启动自动转存定时任务")
}

// StopAutoTransferScheduler 停止自动转存定时任务
func (gs *GlobalScheduler) StopAutoTransferScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.scheduler.IsAutoTransferRunning() {
		Info("自动转存定时任务未在运行")
		return
	}

	gs.scheduler.StopAutoTransferScheduler()
	Info("全局调度器已停止自动转存定时任务")
}

// IsAutoTransferRunning 检查自动转存定时任务是否在运行
func (gs *GlobalScheduler) IsAutoTransferRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.scheduler.IsAutoTransferRunning()
}

// ProcessAutoTransfer 手动触发自动转存处理
func (gs *GlobalScheduler) ProcessAutoTransfer() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.scheduler.processAutoTransfer()
}

// UpdateSchedulerStatusWithAutoTransfer 根据系统配置更新调度器状态（包含自动转存）
func (gs *GlobalScheduler) UpdateSchedulerStatusWithAutoTransfer(autoFetchHotDramaEnabled bool, autoProcessReadyResources bool, autoTransferEnabled bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// 处理热播剧自动拉取功能
	if autoFetchHotDramaEnabled {
		if !gs.scheduler.IsRunning() {
			Info("系统配置启用自动拉取热播剧，启动定时任务")
			gs.scheduler.StartHotDramaScheduler()
		}
	} else {
		if gs.scheduler.IsRunning() {
			Info("系统配置禁用自动拉取热播剧，停止定时任务")
			gs.scheduler.StopHotDramaScheduler()
		}
	}

	// 处理待处理资源自动处理功能
	if autoProcessReadyResources {
		if !gs.scheduler.IsReadyResourceRunning() {
			Info("系统配置启用自动处理待处理资源，启动定时任务")
			gs.scheduler.StartReadyResourceScheduler()
		}
	} else {
		if gs.scheduler.IsReadyResourceRunning() {
			Info("系统配置禁用自动处理待处理资源，停止定时任务")
			gs.scheduler.StopReadyResourceScheduler()
		}
	}

	// 处理自动转存功能
	if autoTransferEnabled {
		if !gs.scheduler.IsAutoTransferRunning() {
			Info("系统配置启用自动转存，启动定时任务")
			gs.scheduler.StartAutoTransferScheduler()
		}
	} else {
		if gs.scheduler.IsAutoTransferRunning() {
			Info("系统配置禁用自动转存，停止定时任务")
			gs.scheduler.StopAutoTransferScheduler()
		}
	}
}
