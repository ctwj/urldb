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
func GetGlobalScheduler(hotDramaRepo repo.HotDramaRepository, readyResourceRepo repo.ReadyResourceRepository, resourceRepo repo.ResourceRepository, systemConfigRepo repo.SystemConfigRepository, panRepo repo.PanRepository) *GlobalScheduler {
	once.Do(func() {
		globalScheduler = &GlobalScheduler{
			scheduler: NewScheduler(hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo, panRepo),
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
