package services

import (
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// StatsSummary 仪表盘首屏聚合数据
// 契约: specs/004-admin-ui-optimization/contracts/stats-summary-api.md
// 该结构序列化后与原 handlers.GetSummary 的 gin.H 嵌套结构等价
type StatsSummary struct {
	Resources ResourcesStats `json:"resources"`
	Views     ViewsStats     `json:"views"`
	Searches  SearchesStats  `json:"searches"`
	Todos     TodosStats     `json:"todos"`
}

// ResourcesStats 资源指标（今日/昨日/总量）
type ResourcesStats struct {
	Today     int64 `json:"today"`
	Yesterday int64 `json:"yesterday"`
	Total     int64 `json:"total"`
}

// ViewsStats 浏览量指标
type ViewsStats struct {
	Today     int64 `json:"today"`
	Yesterday int64 `json:"yesterday"`
}

// SearchesStats 搜索量指标
type SearchesStats struct {
	Today     int64 `json:"today"`
	Yesterday int64 `json:"yesterday"`
}

// TodosStats 待办事项聚合
type TodosStats struct {
	ReadyResources int64 `json:"ready_resources"`
	FailedTasks    int64 `json:"failed_tasks"`
	PendingReports int64 `json:"pending_reports"`
}

// StatsService 统计聚合服务
//
// 设计说明：
//   - 独立于 HTTP 层，便于集成测试直接调用 GetSummary()
//   - 通过构造函数注入 db 与 repoManager，避免对 package-level 全局状态的依赖
//   - 保持与现有 GetStats handler 直接操作 DB 的模式一致（无额外抽象层）
type StatsService struct {
	db   *gorm.DB
	repo *repo.RepositoryManager
}

// NewStatsService 创建统计聚合服务实例
func NewStatsService(database *gorm.DB, repoMgr *repo.RepositoryManager) *StatsService {
	return &StatsService{
		db:   database,
		repo: repoMgr,
	}
}

// GetSummary 仪表盘首屏聚合统计（含环比昨日与待办）
//
// 时区：按 Asia/Shanghai (UTC+8) 计算日期边界（00:00:00）
// 错误降级：浏览量查询失败时记 warning 并归零（与现有 GetStats 一致）
// 负值防护：理论不会出现负 COUNT，防御性 clamp 至 0
func (s *StatsService) GetSummary() (*StatsSummary, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := utils.GetCurrentTime()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfToday := startOfToday.Add(24 * time.Hour)
	startOfYesterday := startOfToday.AddDate(0, 0, -1)

	todayStr := startOfToday.Format(utils.TimeFormatDate)
	yesterdayStr := startOfYesterday.Format(utils.TimeFormatDate)

	// resources: 今日/昨日新增 + 总量
	var resourcesToday, resourcesYesterday, resourcesTotal int64
	s.db.Model(&entity.Resource{}).Where("created_at >= ? AND created_at < ?", startOfToday, endOfToday).Count(&resourcesToday)
	s.db.Model(&entity.Resource{}).Where("created_at >= ? AND created_at < ?", startOfYesterday, startOfToday).Count(&resourcesYesterday)
	s.db.Model(&entity.Resource{}).Count(&resourcesTotal)

	// views: 复用 repo GetViewsByDate（错误降级为 0）
	viewsToday, err := s.repo.ResourceViewRepository.GetViewsByDate(todayStr)
	if err != nil {
		utils.Error("GetSummary 获取今日浏览量失败: %v", err)
		viewsToday = 0
	}
	viewsYesterday, err := s.repo.ResourceViewRepository.GetViewsByDate(yesterdayStr)
	if err != nil {
		utils.Error("GetSummary 获取昨日浏览量失败: %v", err)
		viewsYesterday = 0
	}

	// searches: SearchStat 表按日期计数
	var searchesToday, searchesYesterday int64
	s.db.Model(&entity.SearchStat{}).Where("DATE(date) = ?", todayStr).Count(&searchesToday)
	s.db.Model(&entity.SearchStat{}).Where("DATE(date) = ?", yesterdayStr).Count(&searchesYesterday)

	// todos: 待处理资源(无错误)/失败任务/待审核举报
	var readyResources, failedTasks, pendingReports int64
	s.db.Model(&entity.ReadyResource{}).Where("error_msg = '' OR error_msg IS NULL").Count(&readyResources)
	s.db.Model(&entity.Task{}).Where("status = ?", entity.TaskStatusFailed).Count(&failedTasks)
	s.db.Model(&entity.Report{}).Where("status = ?", "pending").Count(&pendingReports)

	// 负值防护（防御性归零）
	clamp := func(v int64) int64 {
		if v < 0 {
			return 0
		}
		return v
	}

	return &StatsSummary{
		Resources: ResourcesStats{
			Today:     clamp(resourcesToday),
			Yesterday: clamp(resourcesYesterday),
			Total:     clamp(resourcesTotal),
		},
		Views: ViewsStats{
			Today:     clamp(viewsToday),
			Yesterday: clamp(viewsYesterday),
		},
		Searches: SearchesStats{
			Today:     clamp(searchesToday),
			Yesterday: clamp(searchesYesterday),
		},
		Todos: TodosStats{
			ReadyResources: clamp(readyResources),
			FailedTasks:    clamp(failedTasks),
			PendingReports: clamp(pendingReports),
		},
	}, nil
}

// 默认实例：在 main.go 中通过 SetDefaultStatsService 注入，
// 便于 handlers 包直接调用 GetSummary 而无需在每个 handler 中传参
var defaultStatsService *StatsService

// SetDefaultStatsService 设置默认统计服务实例（由 main.go 在初始化阶段调用）
func SetDefaultStatsService(s *StatsService) {
	defaultStatsService = s
}

// GetDefaultStatsService 获取默认统计服务实例（handler 与测试共用）
func GetDefaultStatsService() *StatsService {
	if defaultStatsService != nil {
		return defaultStatsService
	}
	// 兜底：未注入时从全局 db.DB 构造（用于开发/旧路径兼容）
	if db.DB != nil {
		defaultStatsService = NewStatsService(db.DB, repo.NewRepositoryManager(db.DB))
		return defaultStatsService
	}
	return nil
}
