package services

import (
	"context"
	"sync"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// ResourceCheckResult 资源维度的检测结论
type ResourceCheckResult struct {
	// Status: "valid" / "invalid" / "undetermined"（未得出结论或未启用检测）
	Status string
	// FailReason 失效原因（仅 status=invalid 时有值）
	FailReason string
	// Platform 识别到的平台
	Platform string
	// DetectionMethod: "pancheck"（启用并得出结论或启用但未得出结论）/ "disabled"（未启用）
	DetectionMethod string
}

// LinkCheckService 链接检测共享服务（两处检测点的唯一入口）
// 注：客户端不再维护本地缓存，缓存由 PanCheck 服务端自行管理。
type LinkCheckService interface {
	// CheckURL 检测单个 URL（scheduler 用）
	CheckURL(ctx context.Context, rawURL string, ignoreCache bool) ResourceCheckResult
	// CheckResources 批量检测资源（详情页用），返回 map[resourceID] -> 结论
	CheckResources(ctx context.Context, resources []*entity.Resource, ignoreCache bool) map[uint]ResourceCheckResult
}

// linkCheckServiceImpl 服务实现
type linkCheckServiceImpl struct {
	client           *PanCheckClient
	systemConfigRepo repo.SystemConfigRepository
	resourceRepo     repo.ResourceRepository
	meiliManager     *MeilisearchManager
}

// NewLinkCheckService 创建链接检测服务
func NewLinkCheckService(
	client *PanCheckClient,
	systemConfigRepo repo.SystemConfigRepository,
	resourceRepo repo.ResourceRepository,
	meiliManager *MeilisearchManager,
) LinkCheckService {
	return &linkCheckServiceImpl{
		client:           client,
		systemConfigRepo: systemConfigRepo,
		resourceRepo:     resourceRepo,
		meiliManager:     meiliManager,
	}
}

// loadConfig 每次调用前读取配置（保证保存后立即生效，SC-005）
func (s *linkCheckServiceImpl) loadConfig() (enabled bool, host string, timeout, batch, concurrency int) {
	enabled, _ = s.systemConfigRepo.GetConfigBool(entity.ConfigKeyPanCheckEnabled)
	host, _ = s.systemConfigRepo.GetConfigValue(entity.ConfigKeyPanCheckHost)
	timeout, _ = s.systemConfigRepo.GetConfigInt(entity.ConfigKeyPanCheckTimeoutSeconds)
	batch, _ = s.systemConfigRepo.GetConfigInt(entity.ConfigKeyPanCheckBatchSize)
	concurrency, _ = s.systemConfigRepo.GetConfigInt(entity.ConfigKeyPanCheckConcurrency)

	if timeout <= 0 {
		timeout = 60
	}
	if batch <= 0 {
		batch = 20
	}
	if concurrency <= 0 {
		concurrency = 5
	}
	return
}

// CheckURL 检测单个 URL
func (s *linkCheckServiceImpl) CheckURL(ctx context.Context, rawURL string, ignoreCache bool) ResourceCheckResult {
	results := s.CheckResources(ctx, []*entity.Resource{{URL: rawURL}}, ignoreCache)
	for _, r := range results {
		return r
	}
	return ResourceCheckResult{Status: "undetermined", DetectionMethod: "disabled"}
}

// CheckResources 批量检测资源
func (s *linkCheckServiceImpl) CheckResources(ctx context.Context, resources []*entity.Resource, ignoreCache bool) map[uint]ResourceCheckResult {
	results := make(map[uint]ResourceCheckResult, len(resources))

	enabled, host, timeout, batch, concurrency := s.loadConfig()
	utils.Info("[PanCheck] CheckResources 入口: 资源数=%d, enabled=%v, host=%q, batch=%d, concurrency=%d, ignoreCache=%v", len(resources), enabled, host, batch, concurrency, ignoreCache)

	// 未启用或未配置 → 跳过检测，调用方按未检测处理
	if !enabled || host == "" {
		utils.Warn("[PanCheck] 检测被跳过(enabled=%v, host=%q) → 所有资源按未检测放行，不调用 PanCheck 服务", enabled, host)
		for _, res := range resources {
			results[res.ID] = ResourceCheckResult{
				Status:          "undetermined",
				DetectionMethod: "disabled",
			}
		}
		return results
	}

	// 构建待检测项：规范化 URL + 所属资源ID列表（一个规范化 URL 可能对应多个资源）
	type pendingItem struct {
		normalized  string
		resourceIDs []uint
	}
	pending := make(map[string]*pendingItem) // normalized -> item
	for _, res := range resources {
		// 先填充默认结论（启用但尚未得出结论）
		results[res.ID] = ResourceCheckResult{
			Status:          "undetermined",
			DetectionMethod: "pancheck",
		}
		n := NormalizeURL(res.URL)
		if n == "" {
			continue
		}
		if item, ok := pending[n]; ok {
			item.resourceIDs = append(item.resourceIDs, res.ID)
		} else {
			pending[n] = &pendingItem{
				normalized:  n,
				resourceIDs: []uint{res.ID},
			}
		}
	}

	if len(pending) == 0 {
		return results
	}

	// 分批并发调用 PanCheck（结果缓存由 PanCheck 服务端自行维护，客户端不再落库）
	urls := make([]string, 0, len(pending))
	for _, item := range pending {
		urls = append(urls, item.normalized)
	}
	utils.Info("[PanCheck] 待检测 URL 数=%d，将分 %d 批调用 PanCheck 服务 %s", len(urls), len(chunkStrings(urls, batch)), host)

	type batchOutcome struct {
		url     string
		outcome LinkCheckOutcome
		has     bool
	}
	outcomeCh := make(chan batchOutcome, len(urls))

	batches := chunkStrings(urls, batch)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, b := range batches {
		wg.Add(1)
		go func(batchURLs []string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			outcomes, err := s.client.Check(ctx, host, timeout, batchURLs)
			if err != nil {
				// 整批视为未得出结论（保持 is_valid 原值、可重试）
				utils.Warn("PanCheck 批次调用失败，视为未得出结论: %v", err)
				return
			}
			for _, u := range batchURLs {
				oc, ok := outcomes[u]
				outcomeCh <- batchOutcome{url: u, outcome: oc, has: ok}
			}
		}(b)
	}

	go func() {
		wg.Wait()
		close(outcomeCh)
	}()

	// 收集结论，填充结果
	for bo := range outcomeCh {
		item, ok := pending[bo.url]
		if !ok {
			continue
		}
		if !bo.has {
			// 该 URL 本次未得出结论（pending 或既不在 valid 也不在 invalid）
			continue
		}
		res := ResourceCheckResult{
			Status:          string(bo.outcome.Status),
			FailReason:      bo.outcome.FailReason,
			Platform:        bo.outcome.Platform,
			DetectionMethod: "pancheck",
		}
		for _, rid := range item.resourceIDs {
			results[rid] = res
		}
	}

	return results
}

// chunkStrings 将切片切分为固定大小的批次
func chunkStrings(in []string, size int) [][]string {
	if size <= 0 {
		size = 20
	}
	var out [][]string
	for i := 0; i < len(in); i += size {
		end := i + size
		if end > len(in) {
			end = len(in)
		}
		out = append(out, in[i:end])
	}
	return out
}

// ApplyValidityWriteback 仅在 is_valid 实际翻转时写回 DB + 同步 Meilisearch。
// 聚合规则：单资源多 URL 时任一 URL 失效即整体失效。
func ApplyValidityWriteback(
	resource *entity.Resource,
	result ResourceCheckResult,
	resourceRepo repo.ResourceRepository,
	meiliManager *MeilisearchManager,
) {
	if result.Status != "valid" && result.Status != "invalid" {
		return // 未得出结论，不翻转
	}
	newValid := result.Status == "valid"
	if newValid == resource.IsValid {
		return // 未翻转，不写回（避免写放大）
	}

	resource.IsValid = newValid
	if err := resourceRepo.UpdateIsValid(resource.ID, newValid); err != nil {
		utils.Error("写回资源有效性失败 - ID: %d, Error: %v", resource.ID, err)
		return
	}

	// 009-statistics-enhancement: 维护失效时间 invalidated_at，支撑每日失效统计（FR-007~009）
	if newValid {
		// 恢复有效：清空失效时间
		if err := resourceRepo.UpdateFields(resource.ID, map[string]interface{}{"invalidated_at": nil}); err != nil {
			utils.Error("清空 invalidated_at 失败 - ID: %d, Error: %v", resource.ID, err)
		}
		resource.InvalidatedAt = nil
	} else {
		// 失效：写入失效发生时间
		now := utils.GetCurrentTime()
		if err := resourceRepo.UpdateFields(resource.ID, map[string]interface{}{"invalidated_at": now}); err != nil {
			utils.Error("写入 invalidated_at 失败 - ID: %d, Error: %v", resource.ID, err)
		}
		resource.InvalidatedAt = &now
	}

	utils.Info("资源有效性翻转 - ID: %d, %v -> %v", resource.ID, !newValid, newValid)

	// 清除热门资源缓存，避免失效资源继续展示在热门列表
	utils.GetHotResourcesCache().DeletePattern("hot_resources_")

	// 同步 Meilisearch
	if meiliManager != nil && meiliManager.IsEnabled() {
		service := meiliManager.GetService()
		if service != nil {
			if err := service.UpdateResourceValidity(resource.ID, newValid); err != nil {
				utils.Error("同步 Meilisearch 有效性失败 - ID: %d, Error: %v", resource.ID, err)
			}
		}
	}
}
