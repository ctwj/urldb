package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/google"
	"github.com/ctwj/urldb/utils"
)

// GoogleIndexProcessor Google索引任务处理器
type GoogleIndexProcessor struct {
	repoMgr *repo.RepositoryManager
	client  *google.Client
	config  *GoogleIndexProcessorConfig
}

// GoogleIndexProcessorConfig Google索引处理器配置
type GoogleIndexProcessorConfig struct {
	CredentialsFile string
	SiteURL         string
	TokenFile       string
	Concurrency     int
	RetryAttempts   int
	RetryDelay      time.Duration
}

// GoogleIndexTaskInput Google索引任务输入数据结构
type GoogleIndexTaskInput struct {
	URLs      []string `json:"urls"`
	Operation string   `json:"operation"` // indexing_check, sitemap_submit, batch_index
	SitemapURL string   `json:"sitemap_url,omitempty"`
}

// GoogleIndexTaskOutput Google索引任务输出数据结构
type GoogleIndexTaskOutput struct {
	URL         string                    `json:"url,omitempty"`
	IndexStatus string                    `json:"index_status,omitempty"`
	Error       string                    `json:"error,omitempty"`
	Success     bool                      `json:"success"`
	Message     string                    `json:"message"`
	Time        string                    `json:"time"`
	Result      *google.URLInspectionResult `json:"result,omitempty"`
}

// NewGoogleIndexProcessor 创建Google索引任务处理器
func NewGoogleIndexProcessor(repoMgr *repo.RepositoryManager) *GoogleIndexProcessor {
	return &GoogleIndexProcessor{
		repoMgr: repoMgr,
	}
}

// GetTaskType 获取任务类型
func (gip *GoogleIndexProcessor) GetTaskType() string {
	return "google_index"
}

// Process 处理Google索引任务项
func (gip *GoogleIndexProcessor) Process(ctx context.Context, taskID uint, item *entity.TaskItem) error {
	utils.Info("开始处理Google索引任务项: %d", item.ID)

	// 解析输入数据
	var input GoogleIndexTaskInput
	if err := json.Unmarshal([]byte(item.InputData), &input); err != nil {
		utils.Error("解析输入数据失败: %v", err)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 400, err.Error())
		return fmt.Errorf("解析输入数据失败: %v", err)
	}

	// 初始化Google客户端
	client, err := gip.initGoogleClient()
	if err != nil {
		utils.Error("初始化Google客户端失败: %v", err)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, err.Error())
		return fmt.Errorf("初始化Google客户端失败: %v", err)
	}

	// 根据操作类型执行不同任务
	switch input.Operation {
	case "url_indexing":
		return gip.processURLIndexing(ctx, client, taskID, item, input)
	case "sitemap_submit":
		return gip.processSitemapSubmit(ctx, client, taskID, item, input)
	case "status_check":
		return gip.processStatusCheck(ctx, client, taskID, item, input)
	default:
		errorMsg := fmt.Sprintf("不支持的操作类型: %s", input.Operation)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 400, errorMsg)
		return fmt.Errorf(errorMsg)
	}
}

// processURLIndexing 处理URL索引检查
func (gip *GoogleIndexProcessor) processURLIndexing(ctx context.Context, client *google.Client, taskID uint, item *entity.TaskItem, input GoogleIndexTaskInput) error {
	utils.Info("开始URL索引检查: %v", input.URLs)

	for _, url := range input.URLs {
		select {
		case <-ctx.Done():
			gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 0, "任务被取消")
			return ctx.Err()
		default:
			// 检查URL索引状态
			result, err := gip.inspectURL(client, url)
			if err != nil {
				utils.Error("检查URL索引状态失败: %s, 错误: %v", url, err)
				gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, err.Error())
				continue
			}

			// 更新任务项状态
			var lastCrawled *time.Time
			if result.IndexStatusResult.LastCrawled != "" {
				parsedTime, err := time.Parse(time.RFC3339, result.IndexStatusResult.LastCrawled)
				if err == nil {
					lastCrawled = &parsedTime
				}
			}

			gip.updateTaskItemStatus(item, entity.TaskItemStatusSuccess, result.IndexStatusResult.IndexingState, result.MobileUsabilityResult.MobileFriendly, lastCrawled, 200, "")

			// 更新URL状态记录
			gip.updateURLStatus(url, result.IndexStatusResult.IndexingState, lastCrawled)

			// 添加延迟避免API限制
			time.Sleep(100 * time.Millisecond)
		}
	}

	utils.Info("URL索引检查完成")
	return nil
}

// processSitemapSubmit 处理网站地图提交
func (gip *GoogleIndexProcessor) processSitemapSubmit(ctx context.Context, client *google.Client, taskID uint, item *entity.TaskItem, input GoogleIndexTaskInput) error {
	utils.Info("开始网站地图提交: %s", input.SitemapURL)

	if input.SitemapURL == "" {
		errorMsg := "网站地图URL不能为空"
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 400, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// 提交网站地图
	err := client.SubmitSitemap(input.SitemapURL)
	if err != nil {
		utils.Error("提交网站地图失败: %s, 错误: %v", input.SitemapURL, err)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, err.Error())
		return fmt.Errorf("提交网站地图失败: %v", err)
	}

	// 更新任务项状态
	now := time.Now()
	gip.updateTaskItemStatus(item, entity.TaskItemStatusSuccess, "SUBMITTED", false, &now, 200, "")

	utils.Info("网站地图提交完成: %s", input.SitemapURL)
	return nil
}

// processStatusCheck 处理状态检查
func (gip *GoogleIndexProcessor) processStatusCheck(ctx context.Context, client *google.Client, taskID uint, item *entity.TaskItem, input GoogleIndexTaskInput) error {
	utils.Info("开始状态检查: %v", input.URLs)

	for _, url := range input.URLs {
		select {
		case <-ctx.Done():
			gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 0, "任务被取消")
			return ctx.Err()
		default:
			// 检查URL状态
			result, err := gip.inspectURL(client, url)
			if err != nil {
				utils.Error("检查URL状态失败: %s, 错误: %v", url, err)
				gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, err.Error())
				continue
			}

			// 更新任务项状态
			var lastCrawled *time.Time
			if result.IndexStatusResult.LastCrawled != "" {
				parsedTime, err := time.Parse(time.RFC3339, result.IndexStatusResult.LastCrawled)
				if err == nil {
					lastCrawled = &parsedTime
				}
			}

			gip.updateTaskItemStatus(item, entity.TaskItemStatusSuccess, result.IndexStatusResult.IndexingState, result.MobileUsabilityResult.MobileFriendly, lastCrawled, 200, "")

			utils.Info("URL状态检查完成: %s, 状态: %s", url, result.IndexStatusResult.IndexingState)
		}
	}

	return nil
}

// initGoogleClient 初始化Google客户端
func (gip *GoogleIndexProcessor) initGoogleClient() (*google.Client, error) {
	// 从配置中获取Google认证信息
	credentialsFile, err := gip.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeyCredentialsFile)
	if err != nil || credentialsFile == "" {
		return nil, fmt.Errorf("未配置Google认证文件: %v", err)
	}

	siteURL, err := gip.repoMgr.SystemConfigRepository.GetConfigValue(entity.GoogleIndexConfigKeySiteURL)
	if err != nil || siteURL == "" {
		return nil, fmt.Errorf("未配置网站URL: %v", err)
	}

	config := &google.Config{
		CredentialsFile: credentialsFile,
		SiteURL:         siteURL,
		TokenFile:       "google_token.json", // 使用固定token文件名
	}

	client, err := google.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建Google客户端失败: %v", err)
	}

	return client, nil
}

// inspectURL 检查URL索引状态
func (gip *GoogleIndexProcessor) inspectURL(client *google.Client, url string) (*google.URLInspectionResult, error) {
	// 重试机制
	var result *google.URLInspectionResult
	var err error

	for attempt := 0; attempt <= gip.config.RetryAttempts; attempt++ {
		result, err = client.InspectURL(url)
		if err == nil {
			break // 成功则退出重试循环
		}

		if attempt < gip.config.RetryAttempts {
			utils.Info("URL检查失败，第%d次重试: %s, 错误: %v", attempt+1, url, err)
			time.Sleep(gip.config.RetryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("检查URL失败: %v", err)
	}

	return result, nil
}

// updateTaskItemStatus 更新任务项状态
func (gip *GoogleIndexProcessor) updateTaskItemStatus(item *entity.TaskItem, status entity.TaskItemStatus, indexStatus string, mobileFriendly bool, lastCrawled *time.Time, statusCode int, errorMessage string) {
	item.Status = status
	item.ErrorMessage = errorMessage

	// 更新Google索引特有字段
	item.IndexStatus = indexStatus
	item.MobileFriendly = mobileFriendly
	item.LastCrawled = lastCrawled
	item.StatusCode = statusCode

	now := time.Now()
	item.ProcessedAt = &now

	// 保存更新
	if err := gip.repoMgr.TaskItemRepository.Update(item); err != nil {
		utils.Error("更新任务项状态失败: %v", err)
	}
}

// updateURLStatus 更新URL状态记录（使用任务项存储）
func (gip *GoogleIndexProcessor) updateURLStatus(url string, indexStatus string, lastCrawled *time.Time) {
	// 在任务项中记录URL状态，而不是使用专门的URL状态表
	// 此功能现在通过任务系统中的TaskItem记录来跟踪
	utils.Debug("URL状态已更新: %s, 状态: %s", url, indexStatus)
}

// BatchProcessURLs 批量处理URLs
func (gip *GoogleIndexProcessor) BatchProcessURLs(ctx context.Context, urls []string, operation string, taskID uint) error {
	utils.Info("开始批量处理URLs，数量: %d, 操作: %s", len(urls), operation)

	// 根据并发数创建工作池
	semaphore := make(chan struct{}, gip.config.Concurrency)
	errChan := make(chan error, len(urls))

	for _, url := range urls {
		go func(u string) {
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			// 处理单个URL
			client, err := gip.initGoogleClient()
			if err != nil {
				errChan <- fmt.Errorf("初始化客户端失败: %v", err)
				return
			}

			result, err := gip.inspectURL(client, u)
			if err != nil {
				utils.Error("处理URL失败: %s, 错误: %v", u, err)
				errChan <- err
				return
			}

			// 更新状态
			var lastCrawled *time.Time
			if result.IndexStatusResult.LastCrawled != "" {
				parsedTime, err := time.Parse(time.RFC3339, result.IndexStatusResult.LastCrawled)
				if err == nil {
					lastCrawled = &parsedTime
				}
			}

			// 创建任务项记录
			now := time.Now()
			inputData := map[string]interface{}{
				"urls":      []string{u},
				"operation": "url_indexing",
			}
			inputDataJSON, _ := json.Marshal(inputData)

			taskItem := &entity.TaskItem{
				TaskID:       taskID,
				Status:       entity.TaskItemStatusSuccess,
				InputData:    string(inputDataJSON),
				URL:          u,
				IndexStatus:  result.IndexStatusResult.IndexingState,
				MobileFriendly: result.MobileUsabilityResult.MobileFriendly,
				LastCrawled:  lastCrawled,
				StatusCode:   200,
				ProcessedAt:  &now,
			}

			if err := gip.repoMgr.TaskItemRepository.Create(taskItem); err != nil {
				utils.Error("创建任务项失败: %v", err)
			}

			// 更新URL状态
			gip.updateURLStatus(u, result.IndexStatusResult.IndexingState, lastCrawled)

			errChan <- nil
		}(url)
	}

	// 等待所有goroutine完成
	for i := 0; i < len(urls); i++ {
		err := <-errChan
		if err != nil {
			utils.Error("批量处理URL时出错: %v", err)
		}
	}

	utils.Info("批量处理URLs完成")
	return nil
}

// SubmitSitemap 提交网站地图
func (gip *GoogleIndexProcessor) SubmitSitemap(ctx context.Context, sitemapURL string, taskID uint) error {
	utils.Info("开始提交网站地图: %s", sitemapURL)

	client, err := gip.initGoogleClient()
	if err != nil {
		return fmt.Errorf("初始化Google客户端失败: %v", err)
	}

	err = client.SubmitSitemap(sitemapURL)
	if err != nil {
		return fmt.Errorf("提交网站地图失败: %v", err)
	}

	// 创建任务项记录
	now := time.Now()
	inputData := map[string]interface{}{
		"sitemap_url": sitemapURL,
		"operation":   "sitemap_submit",
	}
	inputDataJSON, _ := json.Marshal(inputData)

	taskItem := &entity.TaskItem{
		TaskID:       taskID,
		Status:       entity.TaskItemStatusSuccess,
		InputData:    string(inputDataJSON),
		URL:          sitemapURL,
		IndexStatus:  "SUBMITTED",
		StatusCode:   200,
		ProcessedAt:  &now,
	}

	if err := gip.repoMgr.TaskItemRepository.Create(taskItem); err != nil {
		utils.Error("创建任务项失败: %v", err)
	}

	utils.Info("网站地图提交完成: %s", sitemapURL)
	return nil
}