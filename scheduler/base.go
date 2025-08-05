package scheduler

import (
	"sync"
	"time"

	"github.com/ctwj/urldb/db/repo"
)

// BaseScheduler 基础调度器结构
type BaseScheduler struct {
	// 共享的仓库
	hotDramaRepo      repo.HotDramaRepository
	readyResourceRepo repo.ReadyResourceRepository
	resourceRepo      repo.ResourceRepository
	systemConfigRepo  repo.SystemConfigRepository
	panRepo           repo.PanRepository
	cksRepo           repo.CksRepository
	tagRepo           repo.TagRepository
	categoryRepo      repo.CategoryRepository

	// 控制字段
	stopChan  chan bool
	isRunning bool

	// 平台映射缓存
	panCache     map[string]*uint // serviceType -> panID
	panCacheOnce sync.Once
}

// NewBaseScheduler 创建基础调度器
func NewBaseScheduler(
	hotDramaRepo repo.HotDramaRepository,
	readyResourceRepo repo.ReadyResourceRepository,
	resourceRepo repo.ResourceRepository,
	systemConfigRepo repo.SystemConfigRepository,
	panRepo repo.PanRepository,
	cksRepo repo.CksRepository,
	tagRepo repo.TagRepository,
	categoryRepo repo.CategoryRepository,
) *BaseScheduler {
	return &BaseScheduler{
		hotDramaRepo:      hotDramaRepo,
		readyResourceRepo: readyResourceRepo,
		resourceRepo:      resourceRepo,
		systemConfigRepo:  systemConfigRepo,
		panRepo:           panRepo,
		cksRepo:           cksRepo,
		tagRepo:           tagRepo,
		categoryRepo:      categoryRepo,
		stopChan:          make(chan bool),
		isRunning:         false,
		panCache:          make(map[string]*uint),
	}
}

// Stop 停止调度器
func (b *BaseScheduler) Stop() {
	if b.isRunning {
		b.stopChan <- true
		b.isRunning = false
	}
}

// IsRunning 检查是否正在运行
func (b *BaseScheduler) IsRunning() bool {
	return b.isRunning
}

// SetRunning 设置运行状态
func (b *BaseScheduler) SetRunning(running bool) {
	b.isRunning = running
}

// GetStopChan 获取停止通道
func (b *BaseScheduler) GetStopChan() chan bool {
	return b.stopChan
}

// SleepWithStopCheck 带停止检查的睡眠
func (b *BaseScheduler) SleepWithStopCheck(duration time.Duration) bool {
	select {
	case <-time.After(duration):
		return false
	case <-b.stopChan:
		return true
	}
}
