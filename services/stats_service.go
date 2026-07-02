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
	Resources              ResourcesStats         `json:"resources"`
	Views                  ViewsStats             `json:"views"`
	Searches               SearchesStats          `json:"searches"`
	Todos                  TodosStats             `json:"todos"`
	ViewPanDistribution    []ViewDistributionItem `json:"view_pan_distribution"`
	ViewSourceDistribution []ViewDistributionItem `json:"view_source_distribution"`
}

// ResourcesStats 资源指标（今日/昨日/总量 + 失效/同步）
type ResourcesStats struct {
	Today        int64 `json:"today"`
	Yesterday    int64 `json:"yesterday"`
	Total        int64 `json:"total"`
	InvalidTotal int64 `json:"invalid_total"` // 009: 失效资源总数（is_valid=false）
	SyncedTotal  int64 `json:"synced_total"`  // 009: 已同步搜索索引资源数
	TodayInvalid int64 `json:"today_invalid"` // 009: 今日失效资源数
	TodaySynced  int64 `json:"today_synced"`  // 009: 今日新同步资源数
}

// ViewsStats 浏览量指标
type ViewsStats struct {
	Today     int64 `json:"today"`
	Yesterday int64 `json:"yesterday"`
	Total     int64 `json:"total"` // 009: 访问（获取资源）总次数
}

// ViewDistributionItem 访问分布项（网盘 / 来源通用）。009-statistics-enhancement
type ViewDistributionItem struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Count   int64  `json:"count"`
	Percent int    `json:"percent"`
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

	// 009-statistics-enhancement: 失效/同步总量
	var invalidTotal, syncedTotal int64
	s.db.Model(&entity.Resource{}).Where("is_valid = ?", false).Count(&invalidTotal)
	syncedTotal, _ = s.repo.ResourceRepository.CountSyncedToMeilisearch()

	// 009: 今日失效 / 今日同步（仪表盘卡片"今日"小字）
	todayInvalid, _ := s.repo.ResourceRepository.CountInvalidByDate(todayStr)
	todaySynced, _ := s.repo.ResourceRepository.CountSyncedByDate(todayStr)

	// 009: 访问（获取资源）总次数 + 网盘/来源分布（均来自 resource_views）
	viewsTotal, vErr := s.repo.ResourceViewRepository.GetTotalViews()
	if vErr != nil {
		utils.Error("GetSummary 获取访问总次数失败: %v", vErr)
		viewsTotal = 0
	}
	panRows, pErr := s.repo.ResourceViewRepository.GetPanDistribution()
	if pErr != nil {
		utils.Error("GetSummary 获取网盘分布失败: %v", pErr)
		panRows = nil
	}
	sourceRows, sErr := s.repo.ResourceViewRepository.GetSourceDistribution()
	if sErr != nil {
		utils.Error("GetSummary 获取来源分布失败: %v", sErr)
		sourceRows = nil
	}

	// 负值防护（防御性归零）
	clamp := func(v int64) int64 {
		if v < 0 {
			return 0
		}
		return v
	}

	return &StatsSummary{
		Resources: ResourcesStats{
			Today:        clamp(resourcesToday),
			Yesterday:    clamp(resourcesYesterday),
			Total:        clamp(resourcesTotal),
			InvalidTotal: clamp(invalidTotal),
			SyncedTotal:  clamp(syncedTotal),
			TodayInvalid: clamp(todayInvalid),
			TodaySynced:  clamp(todaySynced),
		},
		Views: ViewsStats{
			Today:     clamp(viewsToday),
			Yesterday: clamp(viewsYesterday),
			Total:     clamp(viewsTotal),
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
		ViewPanDistribution:    toViewDistribution(panRows, "pan", viewsTotal),
		ViewSourceDistribution: toViewDistribution(sourceRows, "source", viewsTotal),
	}, nil
}

// toViewDistribution 将 repo 返回的分布行（map）转为前端分布项，并计算占比。
// keyField 为分布键列名（"pan" 或 "source"）；total 为该分布的总计数，用于求 percent。
func toViewDistribution(rows []map[string]interface{}, keyField string, total int64) []ViewDistributionItem {
	items := make([]ViewDistributionItem, 0, len(rows))
	for _, row := range rows {
		key := ""
		if v, ok := row[keyField]; ok && v != nil {
			key = interfaceToString(v)
		}
		var count int64
		if v, ok := row["count"]; ok {
			count = interfaceToInt64(v)
		}
		name := key
		if keyField == "source" {
			name = entity.SourceDisplayName(key)
		}
		percent := 0
		if total > 0 {
			percent = int(float64(count)/float64(total)*100 + 0.5)
		}
		items = append(items, ViewDistributionItem{Key: key, Name: name, Count: count, Percent: percent})
	}
	return items
}

// interfaceToString 从 interface{} 安全取字符串
func interfaceToString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	default:
		if v != nil {
			return ""
		}
		return ""
	}
}

// interfaceToInt64 从 interface{} 安全取 int64（兼容 PG/GORM 返回的 int64/float64 等）
func interfaceToInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	case float32:
		return int64(n)
	default:
		return 0
	}
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
