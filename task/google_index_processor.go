package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
		config: &GoogleIndexProcessorConfig{
			RetryAttempts: 3,
			RetryDelay:    2 * time.Second,
		},
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
	case "url_submit":
		return gip.processURLSubmit(ctx, client, taskID, item, input)
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
	utils.Info("[GOOGLE-PROCESSOR] 开始网站地图提交任务: %s", input.SitemapURL)

	if input.SitemapURL == "" {
		errorMsg := "网站地图URL不能为空"
		utils.Error("[GOOGLE-PROCESSOR] %s", errorMsg)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 400, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// 验证网站地图URL格式
	if !strings.HasPrefix(input.SitemapURL, "http://") && !strings.HasPrefix(input.SitemapURL, "https://") {
		errorMsg := fmt.Sprintf("网站地图URL格式错误，必须以http://或https://开头: %s", input.SitemapURL)
		utils.Error("[GOOGLE-PROCESSOR] %s", errorMsg)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 400, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	utils.Info("[GOOGLE-PROCESSOR] 提交网站地图到Google...")
	// 提交网站地图
	err := client.SubmitSitemap(input.SitemapURL)
	if err != nil {
		utils.Error("[GOOGLE-PROCESSOR] 提交网站地图失败: %s, 错误: %v", input.SitemapURL, err)
		errorMessage := fmt.Sprintf("网站地图提交失败: %v", err)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, errorMessage)
		return fmt.Errorf("提交网站地图失败: %v", err)
	}

	// 更新任务项状态
	now := time.Now()
	gip.updateTaskItemStatus(item, entity.TaskItemStatusSuccess, "SUBMITTED", false, &now, 200, "")

	utils.Info("[GOOGLE-PROCESSOR] 网站地图提交任务完成: %s", input.SitemapURL)
	return nil
}

// processURLSubmit 处理URL提交到索引
func (gip *GoogleIndexProcessor) processURLSubmit(ctx context.Context, client *google.Client, taskID uint, item *entity.TaskItem, input GoogleIndexTaskInput) error {
	utils.Info("[GOOGLE-PROCESSOR] 开始URL提交到索引任务，URL数量: %d", len(input.URLs))

	submittedCount := 0
	failedCount := 0

	for i, url := range input.URLs {
		select {
		case <-ctx.Done():
			utils.Error("[GOOGLE-PROCESSOR] URL提交任务被取消")
			gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 0, "任务被取消")
			return ctx.Err()
		default:
			utils.Info("[GOOGLE-PROCESSOR] 处理URL %d/%d: %s", i+1, len(input.URLs), url)

			// 提交URL到索引
			err := client.PublishURL(url, "URL_UPDATED")
			if err != nil {
				utils.Error("[GOOGLE-PROCESSOR] 提交URL到索引失败: %s, 错误: %v", url, err)

				// 更新失败状态，但继续处理其他URL
				errorMessage := fmt.Sprintf("URL %s 提交失败: %v", url, err)
				gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, errorMessage)
				failedCount++

				// 即使失败也要等待，避免触发更多频率限制
				time.Sleep(2 * time.Second)
				continue
			}

			submittedCount++
			utils.Info("[GOOGLE-PROCESSOR] URL提交成功: %s", url)

			// Indexing API有严格的频率限制，增加延迟
			time.Sleep(2 * time.Second)
		}
	}

	// 更新任务项状态
	now := time.Now()
	statusMessage := fmt.Sprintf("成功提交: %d, 失败: %d", submittedCount, failedCount)

	// 根据提交结果确定任务项状态
	var finalStatus entity.TaskItemStatus
	var statusCode int
	if submittedCount > 0 && failedCount == 0 {
		// 全部成功
		finalStatus = entity.TaskItemStatusSuccess
		statusCode = 200
		utils.Info("[GOOGLE-PROCESSOR] URL提交任务全部成功: %s", statusMessage)
	} else if submittedCount == 0 && failedCount > 0 {
		// 全部失败
		finalStatus = entity.TaskItemStatusFailed
		statusCode = 500
		utils.Error("[GOOGLE-PROCESSOR] URL提交任务全部失败: %s", statusMessage)
	} else {
		// 部分成功
		finalStatus = entity.TaskItemStatusSuccess // 部分成功算作完成，但错误消息会显示失败数量
		statusCode = 206 // 206 Partial Content
		utils.Info("[GOOGLE-PROCESSOR] URL提交任务部分成功: %s", statusMessage)
	}

	gip.updateTaskItemStatus(item, finalStatus, "SUBMITTED", false, &now, statusCode, statusMessage)

	utils.Info("[GOOGLE-PROCESSOR] URL提交任务完成: %s", statusMessage)

	// 如果所有URL都提交失败，返回错误
	if submittedCount == 0 && failedCount > 0 {
		return fmt.Errorf("所有URL提交失败，失败数量: %d", failedCount)
	}

	return nil
}

// processStatusCheck 处理状态检查
func (gip *GoogleIndexProcessor) processStatusCheck(ctx context.Context, client *google.Client, taskID uint, item *entity.TaskItem, input GoogleIndexTaskInput) error {
	utils.Info("开始状态检查: %v", input.URLs)

	successCount := 0
	failedCount := 0

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
				failedCount++
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
			successCount++

			utils.Info("URL状态检查完成: %s, 状态: %s", url, result.IndexStatusResult.IndexingState)
		}
	}

	// 如果所有URL都检查失败，返回错误
	if successCount == 0 && failedCount > 0 {
		errorMsg := fmt.Sprintf("所有URL状态检查失败，失败数量: %d", failedCount)
		gip.updateTaskItemStatus(item, entity.TaskItemStatusFailed, "", false, nil, 500, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	return nil
}

// initGoogleClient 初始化Google客户端
func (gip *GoogleIndexProcessor) initGoogleClient() (*google.Client, error) {
	utils.Info("[GOOGLE-PROCESSOR] 开始初始化Google客户端")

	// 使用固定的凭据文件路径，与验证逻辑保持一致
	credentialsFile := "data/google_credentials.json"
	utils.Info("[GOOGLE-PROCESSOR] 检查凭据文件: %s", credentialsFile)

	// 检查凭据文件是否存在
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		utils.Error("[GOOGLE-PROCESSOR] Google凭据文件不存在: %s", credentialsFile)
		return nil, fmt.Errorf("Google凭据文件不存在: %s", credentialsFile)
	}
	utils.Info("[GOOGLE-PROCESSOR] 凭据文件存在")

	// 从配置中获取网站URL
	siteURL, err := gip.repoMgr.SystemConfigRepository.GetConfigValue(entity.ConfigKeyWebsiteURL)
	if err != nil {
		utils.Error("[GOOGLE-PROCESSOR] 获取网站URL配置失败: %v", err)
		return nil, fmt.Errorf("获取网站URL配置失败: %v", err)
	}

	if siteURL == "" {
		utils.Error("[GOOGLE-PROCESSOR] 网站URL配置为空")
		return nil, fmt.Errorf("网站URL配置为空，请在站点配置中设置正确的网站URL")
	}

	if siteURL == "https://example.com" {
		utils.Error("[GOOGLE-PROCESSOR] 网站URL仍为默认值，请更新为实际网站URL")
		return nil, fmt.Errorf("网站URL仍为默认值，请在站点配置中设置正确的网站URL")
	}

	utils.Info("[GOOGLE-PROCESSOR] 使用网站URL: %s", siteURL)

	config := &google.Config{
		CredentialsFile: credentialsFile,
		SiteURL:         siteURL,
		TokenFile:       "data/google_token.json", // 使用固定token文件名，放在data目录下
	}

	utils.Info("[GOOGLE-PROCESSOR] 开始创建Google客户端...")
	client, err := google.NewClient(config)
	if err != nil {
		utils.Error("[GOOGLE-PROCESSOR] 创建Google客户端失败: %v", err)
		return nil, fmt.Errorf("创建Google客户端失败: %v", err)
	}

	utils.Info("[GOOGLE-PROCESSOR] Google客户端初始化成功")
	return client, nil
}

// inspectURL 检查URL索引状态
func (gip *GoogleIndexProcessor) inspectURL(client *google.Client, url string) (*google.URLInspectionResult, error) {
	utils.Info("[GOOGLE-PROCESSOR] 开始检查URL索引状态: %s", url)

	// 重试机制
	var result *google.URLInspectionResult
	var err error

	for attempt := 0; attempt <= gip.config.RetryAttempts; attempt++ {
		utils.Info("[GOOGLE-PROCESSOR] URL检查尝试 %d/%d: %s", attempt+1, gip.config.RetryAttempts+1, url)
		result, err = client.InspectURL(url)
		if err == nil {
			utils.Info("[GOOGLE-PROCESSOR] URL检查成功: %s", url)
			break // 成功则退出重试循环
		}

		if attempt < gip.config.RetryAttempts {
			utils.Info("[GOOGLE-PROCESSOR] URL检查失败，第%d次重试: %s, 错误: %v", attempt+1, url, err)
			time.Sleep(gip.config.RetryDelay)
		}
	}

	if err != nil {
		utils.Error("[GOOGLE-PROCESSOR] URL检查最终失败: %s, 错误: %v", url, err)
		return nil, fmt.Errorf("检查URL失败: %v", err)
	}

	if result != nil {
		utils.Info("[GOOGLE-PROCESSOR] URL检查结果: %s - 索引状态: %s", url, result.IndexStatusResult.IndexingState)
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