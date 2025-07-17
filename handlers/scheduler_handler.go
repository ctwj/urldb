package handlers

import (
	"net/http"

	"github.com/ctwj/panResManage/utils"
	"github.com/gin-gonic/gin"
)

// GetSchedulerStatus 获取调度器状态
func GetSchedulerStatus(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)

	status := gin.H{
		"hot_drama_scheduler_running":      scheduler.IsHotDramaSchedulerRunning(),
		"ready_resource_scheduler_running": scheduler.IsReadyResourceRunning(),
		"auto_transfer_scheduler_running":  scheduler.IsAutoTransferRunning(),
	}

	SuccessResponse(c, status)
}

// 启动热播剧定时任务
func StartHotDramaScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	if scheduler.IsHotDramaSchedulerRunning() {
		ErrorResponse(c, "热播剧定时任务已在运行中", http.StatusBadRequest)
		return
	}
	scheduler.StartHotDramaScheduler()
	SuccessResponse(c, gin.H{"message": "热播剧定时任务已启动"})
}

// 停止热播剧定时任务
func StopHotDramaScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	if !scheduler.IsHotDramaSchedulerRunning() {
		ErrorResponse(c, "热播剧定时任务未在运行", http.StatusBadRequest)
		return
	}
	scheduler.StopHotDramaScheduler()
	SuccessResponse(c, gin.H{"message": "热播剧定时任务已停止"})
}

// 手动触发热播剧定时任务
func TriggerHotDramaScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	scheduler.StartHotDramaScheduler() // 直接启动一次
	SuccessResponse(c, gin.H{"message": "手动触发热播剧定时任务成功"})
}

// 手动获取热播剧名字
func FetchHotDramaNames(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	names, err := scheduler.GetHotDramaNames()
	if err != nil {
		ErrorResponse(c, "获取热播剧名字失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	SuccessResponse(c, gin.H{"names": names})
}

// 启动待处理资源自动处理任务
func StartReadyResourceScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	if scheduler.IsReadyResourceRunning() {
		ErrorResponse(c, "待处理资源自动处理任务已在运行中", http.StatusBadRequest)
		return
	}
	scheduler.StartReadyResourceScheduler()
	SuccessResponse(c, gin.H{"message": "待处理资源自动处理任务已启动"})
}

// 停止待处理资源自动处理任务
func StopReadyResourceScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	if !scheduler.IsReadyResourceRunning() {
		ErrorResponse(c, "待处理资源自动处理任务未在运行", http.StatusBadRequest)
		return
	}
	scheduler.StopReadyResourceScheduler()
	SuccessResponse(c, gin.H{"message": "待处理资源自动处理任务已停止"})
}

// 手动触发待处理资源自动处理任务
func TriggerReadyResourceScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	// 手动触发一次处理
	scheduler.ProcessReadyResources()
	SuccessResponse(c, gin.H{"message": "手动触发待处理资源自动处理任务成功"})
}

// 启动自动转存定时任务
func StartAutoTransferScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	if scheduler.IsAutoTransferRunning() {
		ErrorResponse(c, "自动转存定时任务已在运行中", http.StatusBadRequest)
		return
	}
	scheduler.StartAutoTransferScheduler()
	SuccessResponse(c, gin.H{"message": "自动转存定时任务已启动"})
}

// 停止自动转存定时任务
func StopAutoTransferScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	if !scheduler.IsAutoTransferRunning() {
		ErrorResponse(c, "自动转存定时任务未在运行", http.StatusBadRequest)
		return
	}
	scheduler.StopAutoTransferScheduler()
	SuccessResponse(c, gin.H{"message": "自动转存定时任务已停止"})
}

// 手动触发自动转存定时任务
func TriggerAutoTransferScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
	)
	// 手动触发一次处理
	scheduler.ProcessAutoTransfer()
	SuccessResponse(c, gin.H{"message": "手动触发自动转存定时任务成功"})
}
