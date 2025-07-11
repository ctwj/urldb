package utils

import (
	"log"
	"res_db/db/repo"
	"sync"
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
func GetGlobalScheduler(hotDramaRepo repo.HotDramaRepository) *GlobalScheduler {
	once.Do(func() {
		globalScheduler = &GlobalScheduler{
			scheduler: NewScheduler(hotDramaRepo),
		}
	})
	return globalScheduler
}

// StartHotDramaScheduler 启动热播剧定时任务
func (gs *GlobalScheduler) StartHotDramaScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.scheduler.IsRunning() {
		log.Println("热播剧定时任务已在运行中")
		return
	}

	gs.scheduler.StartHotDramaScheduler()
	log.Println("全局调度器已启动热播剧定时任务")
}

// StopHotDramaScheduler 停止热播剧定时任务
func (gs *GlobalScheduler) StopHotDramaScheduler() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.scheduler.IsRunning() {
		log.Println("热播剧定时任务未在运行")
		return
	}

	gs.scheduler.StopHotDramaScheduler()
	log.Println("全局调度器已停止热播剧定时任务")
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

// UpdateSchedulerStatus 根据系统配置更新调度器状态
func (gs *GlobalScheduler) UpdateSchedulerStatus(autoFetchHotDramaEnabled bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if autoFetchHotDramaEnabled {
		if !gs.scheduler.IsRunning() {
			log.Println("系统配置启用自动拉取热播剧，启动定时任务")
			gs.scheduler.StartHotDramaScheduler()
		}
	} else {
		if gs.scheduler.IsRunning() {
			log.Println("系统配置禁用自动拉取热播剧，停止定时任务")
			gs.scheduler.StopHotDramaScheduler()
		}
	}
}
