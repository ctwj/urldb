package handlers

import (
	"net/http"
	"res_db/utils"

	"github.com/gin-gonic/gin"
)

// GetSchedulerStatus 获取调度器状态
func GetSchedulerStatus(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(repoManager.HotDramaRepository)

	status := gin.H{
		"hot_drama_scheduler_running": scheduler.IsHotDramaSchedulerRunning(),
	}

	SuccessResponse(c, status)
}

// 启动热播剧定时任务
func StartHotDramaScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(repoManager.HotDramaRepository)
	if scheduler.IsHotDramaSchedulerRunning() {
		ErrorResponse(c, "热播剧定时任务已在运行中", http.StatusBadRequest)
		return
	}
	scheduler.StartHotDramaScheduler()
	SuccessResponse(c, gin.H{"message": "热播剧定时任务已启动"})
}

// 停止热播剧定时任务
func StopHotDramaScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(repoManager.HotDramaRepository)
	if !scheduler.IsHotDramaSchedulerRunning() {
		ErrorResponse(c, "热播剧定时任务未在运行", http.StatusBadRequest)
		return
	}
	scheduler.StopHotDramaScheduler()
	SuccessResponse(c, gin.H{"message": "热播剧定时任务已停止"})
}

// 手动触发热播剧定时任务
func TriggerHotDramaScheduler(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(repoManager.HotDramaRepository)
	scheduler.StartHotDramaScheduler() // 直接启动一次
	SuccessResponse(c, gin.H{"message": "手动触发热播剧定时任务成功"})
}

// 手动获取热播剧名字
func FetchHotDramaNames(c *gin.Context) {
	scheduler := utils.GetGlobalScheduler(repoManager.HotDramaRepository)
	names, err := scheduler.GetHotDramaNames()
	if err != nil {
		ErrorResponse(c, "获取热播剧名字失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	SuccessResponse(c, gin.H{"names": names})
}
