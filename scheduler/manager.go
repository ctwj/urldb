package scheduler

import (
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// Manager 调度器管理器
type Manager struct {
	baseScheduler          *BaseScheduler
	hotDramaScheduler      *HotDramaScheduler
	readyResourceScheduler *ReadyResourceScheduler
}

// NewManager 创建调度器管理器
func NewManager(
	hotDramaRepo repo.HotDramaRepository,
	readyResourceRepo repo.ReadyResourceRepository,
	resourceRepo repo.ResourceRepository,
	systemConfigRepo repo.SystemConfigRepository,
	panRepo repo.PanRepository,
	cksRepo repo.CksRepository,
	tagRepo repo.TagRepository,
	categoryRepo repo.CategoryRepository,
) *Manager {
	// 创建基础调度器
	baseScheduler := NewBaseScheduler(
		hotDramaRepo,
		readyResourceRepo,
		resourceRepo,
		systemConfigRepo,
		panRepo,
		cksRepo,
		tagRepo,
		categoryRepo,
	)

	// 创建各个具体的调度器
	hotDramaScheduler := NewHotDramaScheduler(baseScheduler)
	readyResourceScheduler := NewReadyResourceScheduler(baseScheduler)

	return &Manager{
		baseScheduler:          baseScheduler,
		hotDramaScheduler:      hotDramaScheduler,
		readyResourceScheduler: readyResourceScheduler,
	}
}

// StartAll 启动所有调度任务
func (m *Manager) StartAll() {
	utils.Info("启动所有调度任务")

	// 启动热播剧调度任务
	m.hotDramaScheduler.Start()

	// 启动待处理资源调度任务
	m.readyResourceScheduler.Start()

	utils.Info("所有调度任务已启动")
}

// StopAll 停止所有调度任务
func (m *Manager) StopAll() {
	utils.Info("停止所有调度任务")

	// 停止热播剧调度任务
	m.hotDramaScheduler.Stop()

	// 停止待处理资源调度任务
	m.readyResourceScheduler.Stop()

	utils.Info("所有调度任务已停止")
}

// StartHotDramaScheduler 启动热播剧调度任务
func (m *Manager) StartHotDramaScheduler() {
	m.hotDramaScheduler.Start()
}

// StopHotDramaScheduler 停止热播剧调度任务
func (m *Manager) StopHotDramaScheduler() {
	m.hotDramaScheduler.Stop()
}

// IsHotDramaRunning 检查热播剧调度任务是否正在运行
func (m *Manager) IsHotDramaRunning() bool {
	return m.hotDramaScheduler.IsRunning()
}

// StartReadyResourceScheduler 启动待处理资源调度任务
func (m *Manager) StartReadyResourceScheduler() {
	m.readyResourceScheduler.Start()
}

// StopReadyResourceScheduler 停止待处理资源调度任务
func (m *Manager) StopReadyResourceScheduler() {
	m.readyResourceScheduler.Stop()
}

// IsReadyResourceRunning 检查待处理资源调度任务是否正在运行
func (m *Manager) IsReadyResourceRunning() bool {
	return m.readyResourceScheduler.IsReadyResourceRunning()
}

// GetHotDramaNames 获取热播剧名称列表
func (m *Manager) GetHotDramaNames() ([]string, error) {
	return m.hotDramaScheduler.GetHotDramaNames()
}

// GetStatus 获取所有调度任务的状态
func (m *Manager) GetStatus() map[string]bool {
	return map[string]bool{
		"hot_drama":      m.IsHotDramaRunning(),
		"ready_resource": m.IsReadyResourceRunning(),
	}
}
