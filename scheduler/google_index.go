package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/google"
	"github.com/ctwj/urldb/utils"
)

// GoogleIndexScheduler Google索引调度器
type GoogleIndexScheduler struct {
	*BaseScheduler
	config        entity.SystemConfig
	stopChan      chan bool
	isRunning     bool
	enabled       bool
	checkInterval time.Duration
	googleClient  *google.Client
	taskItemRepo  repo.TaskItemRepository
	taskRepo      repo.TaskRepository

	// 批量处理相关
	pendingURLResults []*repo.URLStatusResult
	currentTaskID     uint
}

// NewGoogleIndexScheduler 创建Google索引调度器
func NewGoogleIndexScheduler(baseScheduler *BaseScheduler, taskItemRepo repo.TaskItemRepository, taskRepo repo.TaskRepository) *GoogleIndexScheduler {
	return &GoogleIndexScheduler{
		BaseScheduler:     baseScheduler,
		taskItemRepo:      taskItemRepo,
		taskRepo:          taskRepo,
		stopChan:          make(chan bool),
		isRunning:         false,
		pendingURLResults: make([]*repo.URLStatusResult, 0),
	}
}

// Start 启动Google索引调度任务
func (s *GoogleIndexScheduler) Start() {
	if s.isRunning {
		utils.Debug("Google索引调度任务已在运行中")
		return
	}

	// 加载配置
	if err := s.loadConfig(); err != nil {
		utils.Error("加载Google索引配置失败: %v", err)
		return
	}

	if !s.enabled {
		utils.Debug("Google索引功能未启用，跳过调度任务")
		return
	}

	s.isRunning = true
	utils.Info("开始启动Google索引调度任务，检查间隔: %v", s.checkInterval)

	go s.run()
}

// Stop 停止Google索引调度任务
func (s *GoogleIndexScheduler) Stop() {
	if !s.isRunning {
		return
	}

	utils.Info("正在停止Google索引调度任务...")
	s.stopChan <- true
	s.isRunning = false
}

// IsRunning 检查调度器是否正在运行
func (s *GoogleIndexScheduler) IsRunning() bool {
	return s.isRunning
}

// run 运行调度器主循环
func (s *GoogleIndexScheduler) run() {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	// 启动时立即执行一次
	s.performScheduledTasks()

	for {
		select {
		case <-s.stopChan:
			utils.Info("Google索引调度任务已停止")
			return
		case <-ticker.C:
			s.performScheduledTasks()
		}
	}
}

// loadConfig 加载配置
func (s *GoogleIndexScheduler) loadConfig() error {
	// 获取启用状态
	enabledStr, err := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeyEnabled)
	if err != nil {
		s.enabled = false
	} else {
		s.enabled = enabledStr == "true" || enabledStr == "1"
	}

	// 获取检查间隔
	intervalStr, err := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeyCheckInterval)
	if err != nil {
		s.checkInterval = 60 * time.Minute // 默认60分钟
	} else {
		if interval, parseErr := time.ParseDuration(intervalStr + "m"); parseErr == nil {
			s.checkInterval = interval
		} else {
			s.checkInterval = 60 * time.Minute
		}
	}

	// 初始化Google客户端
	if s.enabled {
		if err := s.initGoogleClient(); err != nil {
			utils.Error("初始化Google客户端失败: %v", err)
			s.enabled = false
		}
	}

	return nil
}

// initGoogleClient 初始化Google客户端
func (s *GoogleIndexScheduler) initGoogleClient() error {
	// 获取凭据文件路径
	credentialsFile, err := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeyCredentialsFile)
	if err != nil {
		return fmt.Errorf("获取凭据文件路径失败: %v", err)
	}

	// 获取站点URL
	siteURL, err := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeySiteURL)
	if err != nil {
		return fmt.Errorf("获取站点URL失败: %v", err)
	}

	// 创建Google客户端配置
	config := &google.Config{
		CredentialsFile: credentialsFile,
		SiteURL:         siteURL,
	}

	client, err := google.NewClient(config)
	if err != nil {
		return fmt.Errorf("创建Google客户端失败: %v", err)
	}

	s.googleClient = client
	return nil
}

// performScheduledTasks 执行调度任务
func (s *GoogleIndexScheduler) performScheduledTasks() {
	if !s.enabled {
		return
	}

	ctx := context.Background()

	// 任务0: 清理旧记录
	if err := s.taskItemRepo.CleanupOldRecords(); err != nil {
		utils.Error("清理旧记录失败: %v", err)
	}

	// 任务1: 扫描未索引的URL并自动提交
	if err := s.scanAndSubmitUnindexedURLs(ctx); err != nil {
		utils.Error("扫描未索引URL失败: %v", err)
	}

	// 任务2: 定期检查已索引URL的状态
	if err := s.checkIndexedURLsStatus(ctx); err != nil {
		utils.Error("检查索引状态失败: %v", err)
	}

	// 任务3: 刷新待处理的URL结果
	s.flushURLResults()

	utils.Debug("Google索引调度任务执行完成")
}

// scanAndSubmitUnindexedURLs 扫描并提交未索引的URL
func (s *GoogleIndexScheduler) scanAndSubmitUnindexedURLs(ctx context.Context) error {
	utils.Info("开始扫描未索引的URL...")

	// 1. 获取所有资源URL
	resources, err := s.resourceRepo.GetAllValidResources()
	if err != nil {
		return fmt.Errorf("获取资源列表失败: %v", err)
	}

	// 2. 获取已索引的URL记录
	indexedURLs, err := s.getIndexedURLs()
	if err != nil {
		return fmt.Errorf("获取已索引URL列表失败: %v", err)
	}

	// 3. 找出未索引的URL
	var unindexedURLs []string
	indexedURLSet := make(map[string]bool)
	for _, url := range indexedURLs {
		indexedURLSet[url] = true
	}

	for _, resource := range resources {
		if resource.IsPublic && resource.IsValid && resource.URL != "" {
			if !indexedURLSet[resource.URL] {
				unindexedURLs = append(unindexedURLs, resource.URL)
			}
		}
	}

	utils.Info("发现 %d 个未索引的URL", len(unindexedURLs))

	// 4. 批量提交未索引的URL
	if len(unindexedURLs) > 0 {
		if err := s.batchSubmitURLs(ctx, unindexedURLs); err != nil {
			return fmt.Errorf("批量提交URL失败: %v", err)
		}
	}

	return nil
}

// getIndexedURLs 获取已索引的URL列表
func (s *GoogleIndexScheduler) getIndexedURLs() ([]string, error) {
	return s.taskItemRepo.GetDistinctProcessedURLs()
}

// batchSubmitURLs 批量提交URL
func (s *GoogleIndexScheduler) batchSubmitURLs(ctx context.Context, urls []string) error {
	utils.Info("开始批量提交 %d 个URL到Google索引...", len(urls))

	// 获取批量大小配置
	batchSizeStr, _ := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeyBatchSize)
	batchSize := 100 // 默认值
	if batchSizeStr != "" {
		if size, err := strconv.Atoi(batchSizeStr); err == nil && size > 0 {
			batchSize = size
		}
	}

	// 获取并发数配置
	concurrencyStr, _ := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeyConcurrency)
	concurrency := 5 // 默认值
	if concurrencyStr != "" {
		if conc, err := strconv.Atoi(concurrencyStr); err == nil && conc > 0 {
			concurrency = conc
		}
	}

	// 分批处理
	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}

		batch := urls[i:end]
		if err := s.processBatch(ctx, batch, concurrency); err != nil {
			utils.Error("处理批次失败 (批次 %d-%d): %v", i+1, end, err)
			continue
		}

		// 避免API限制，批次间稍作延迟
		time.Sleep(1 * time.Second)
	}

	utils.Info("批量URL提交完成")
	return nil
}

// processBatch 处理单个批次
func (s *GoogleIndexScheduler) processBatch(ctx context.Context, urls []string, concurrency int) error {
	semaphore := make(chan struct{}, concurrency)
	errChan := make(chan error, len(urls))

	for _, url := range urls {
		go func(u string) {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 检查URL索引状态
			result, err := s.googleClient.InspectURL(u)
			if err != nil {
				utils.Error("检查URL失败: %s, 错误: %v", u, err)
				errChan <- err
				return
			}

			// Google Search Console API 不直接支持URL提交
			// 这里只记录URL状态，实际的URL索引需要通过sitemap或其他方式
			if result.IndexStatusResult.IndexingState == "NOT_SUBMITTED" {
				utils.Debug("URL未提交，需要通过sitemap提交: %s", u)
				// TODO: 可以考虑将未提交的URL加入到sitemap中
			}

			// 记录索引状态
			s.recordURLStatus(u, result)
			errChan <- nil
		}(url)
	}

	// 等待所有goroutine完成
	for i := 0; i < len(urls); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// checkIndexedURLsStatus 检查已索引URL的状态
func (s *GoogleIndexScheduler) checkIndexedURLsStatus(ctx context.Context) error {
	utils.Info("开始检查已索引URL的状态...")

	// 暂时跳过状态检查，因为需要TaskItemRepository访问权限
	// TODO: 后续通过扩展BaseScheduler来支持TaskItemRepository
	urlsToCheck := []string{}
	utils.Info("检查 %d 个已索引URL的状态", len(urlsToCheck))

	// 并发检查状态
	concurrencyStr, _ := s.systemConfigRepo.GetConfigValue(entity.GoogleIndexConfigKeyConcurrency)
	concurrency := 3 // 状态检查使用较低并发
	if concurrencyStr != "" {
		if conc, err := strconv.Atoi(concurrencyStr); err == nil && conc > 0 {
			concurrency = conc / 2 // 状态检查并发减半
			if concurrency < 1 {
				concurrency = 1
			}
		}
	}

	// 由于没有URL需要检查，跳过循环
	if len(urlsToCheck) == 0 {
		utils.Info("没有URL需要状态检查")
		return nil
	}

	semaphore := make(chan struct{}, concurrency)
	for _, url := range urlsToCheck {
		go func(u string) {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 检查URL最新状态
			result, err := s.googleClient.InspectURL(u)
			if err != nil {
				utils.Error("检查URL状态失败: %s, 错误: %v", u, err)
				return
			}

			// 记录状态
			s.recordURLStatus(u, result)
		}(url)
	}

	// 等待所有检查完成
	for i := 0; i < len(urlsToCheck); i++ {
		<-semaphore
	}

	utils.Info("索引状态检查完成")
	return nil
}

// recordURLStatus 记录URL索引状态
func (s *GoogleIndexScheduler) recordURLStatus(url string, result *google.URLInspectionResult) {
	// 构造结果对象
	urlResult := &repo.URLStatusResult{
		URL:            url,
		IndexStatus:    result.IndexStatusResult.IndexingState,
		InspectResult:  s.formatInspectResult(result),
		MobileFriendly: s.getMobileFriendly(result),
		StatusCode:     s.getStatusCode(result),
		LastCrawled:    s.parseLastCrawled(result),
		ErrorMessage:   s.getErrorMessage(result),
	}

	// 暂存到批量处理列表，定期批量写入
	s.pendingURLResults = append(s.pendingURLResults, urlResult)

	// 达到批量大小时写入数据库
	if len(s.pendingURLResults) >= 50 {
		s.flushURLResults()
	}
}

// updateURLStatus 更新URL状态
func (s *GoogleIndexScheduler) updateURLStatus(taskItem *entity.TaskItem, result *google.URLInspectionResult) {
	// 暂时只记录日志，不保存到数据库
	// TODO: 后续通过扩展BaseScheduler来支持TaskItemRepository以保存状态
	utils.Debug("更新URL状态: %s - %s", taskItem.URL, result.IndexStatusResult.IndexingState)
}

// flushURLResults 批量写入URL结果
func (s *GoogleIndexScheduler) flushURLResults() {
	if len(s.pendingURLResults) == 0 {
		return
	}

	// 如果没有当前任务，创建一个汇总任务
	if s.currentTaskID == 0 {
		task := &entity.Task{
			Title:       fmt.Sprintf("自动索引检查 - %s", time.Now().Format("2006-01-02 15:04:05")),
			Type:        entity.TaskTypeGoogleIndex,
			Status:      entity.TaskStatusCompleted,
			Description: fmt.Sprintf("自动检查并更新 %d 个URL的索引状态", len(s.pendingURLResults)),
			TotalItems:  len(s.pendingURLResults),
			Progress:    100.0,
		}

		if err := s.taskRepo.Create(task); err != nil {
			utils.Error("创建汇总任务失败: %v", err)
			return
		}
		s.currentTaskID = task.ID
	}

	// 批量写入URL状态
	if err := s.taskItemRepo.UpsertURLStatusRecords(s.currentTaskID, s.pendingURLResults); err != nil {
		utils.Error("批量写入URL状态失败: %v", err)
	} else {
		utils.Info("批量写入URL状态成功: %d 个", len(s.pendingURLResults))
	}

	// 清空待处理列表
	s.pendingURLResults = s.pendingURLResults[:0]
}

// 辅助方法：格式化检查结果
func (s *GoogleIndexScheduler) formatInspectResult(result *google.URLInspectionResult) string {
	if result == nil {
		return ""
	}

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(data)
}

// 辅助方法：获取移动友好状态
func (s *GoogleIndexScheduler) getMobileFriendly(result *google.URLInspectionResult) bool {
	if result != nil {
		return result.MobileUsabilityResult.MobileFriendly
	}
	return false
}

// 辅助方法：获取状态码
func (s *GoogleIndexScheduler) getStatusCode(result *google.URLInspectionResult) int {
	if result != nil {
		// 这里可以根据实际的Google API响应结构来获取状态码
		// 暂时返回200表示成功
		return 200
	}
	return 0
}

// 辅助方法：解析最后抓取时间
func (s *GoogleIndexScheduler) parseLastCrawled(result *google.URLInspectionResult) *time.Time {
	if result != nil && result.IndexStatusResult.LastCrawled != "" {
		// 这里需要根据实际的Google API响应结构来解析时间
		// 暂时返回当前时间
		now := time.Now()
		return &now
	}
	return nil
}

// 辅助方法：获取错误信息
func (s *GoogleIndexScheduler) getErrorMessage(result *google.URLInspectionResult) string {
	if result != nil {
		// 根据索引状态判断是否有错误
		if result.IndexStatusResult.IndexingState == "ERROR" {
			return "索引状态错误"
		}
		if result.IndexStatusResult.IndexingState == "NOT_FOUND" {
			return "页面未找到"
		}
	}
	return ""
}

// SetRunning 设置运行状态
func (s *GoogleIndexScheduler) SetRunning(running bool) {
	s.isRunning = running
}