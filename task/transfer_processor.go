package task

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// TransferProcessor 转存任务处理器
type TransferProcessor struct {
	repoMgr *repo.RepositoryManager
}

// NewTransferProcessor 创建转存任务处理器
func NewTransferProcessor(repoMgr *repo.RepositoryManager) *TransferProcessor {
	return &TransferProcessor{
		repoMgr: repoMgr,
	}
}

// GetTaskType 获取任务类型
func (tp *TransferProcessor) GetTaskType() string {
	return "transfer"
}

// TransferInput 转存任务输入数据结构
type TransferInput struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	CategoryID uint   `json:"category_id"`
	Tags       []uint `json:"tags"`
}

// TransferOutput 转存任务输出数据结构
type TransferOutput struct {
	ResourceID uint   `json:"resource_id,omitempty"`
	SaveURL    string `json:"save_url,omitempty"`
	Error      string `json:"error,omitempty"`
	Success    bool   `json:"success"`
	Time       string `json:"time"`
}

// Process 处理转存任务项
func (tp *TransferProcessor) Process(ctx context.Context, taskID uint, item *entity.TaskItem) error {
	utils.Info("开始处理转存任务项: %d", item.ID)

	// 解析输入数据
	var input TransferInput
	if err := json.Unmarshal([]byte(item.InputData), &input); err != nil {
		return fmt.Errorf("解析输入数据失败: %v", err)
	}

	// 验证输入数据
	if err := tp.validateInput(&input); err != nil {
		return fmt.Errorf("输入数据验证失败: %v", err)
	}

	// 检查资源是否已存在
	exists, existingResource, err := tp.checkResourceExists(input.URL)
	if err != nil {
		utils.Error("检查资源是否存在失败: %v", err)
	}

	if exists {
		// 资源已存在，更新输出数据
		output := TransferOutput{
			ResourceID: existingResource.ID,
			SaveURL:    existingResource.SaveURL,
			Success:    true,
			Time:       time.Now().Format("2006-01-02 15:04:05"),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Info("资源已存在，跳过转存: %s", input.Title)
		return nil
	}

	// 执行转存操作
	resourceID, saveURL, err := tp.performTransfer(ctx, &input)
	if err != nil {
		return fmt.Errorf("转存失败: %v", err)
	}

	// 更新输出数据
	output := TransferOutput{
		ResourceID: resourceID,
		SaveURL:    saveURL,
		Success:    true,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
	}

	outputJSON, _ := json.Marshal(output)
	item.OutputData = string(outputJSON)

	utils.Info("转存任务项处理完成: %d, 资源ID: %d", item.ID, resourceID)
	return nil
}

// validateInput 验证输入数据
func (tp *TransferProcessor) validateInput(input *TransferInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("标题不能为空")
	}

	if strings.TrimSpace(input.URL) == "" {
		return fmt.Errorf("链接不能为空")
	}

	// 验证URL格式
	if !tp.isValidURL(input.URL) {
		return fmt.Errorf("链接格式不正确")
	}

	return nil
}

// isValidURL 验证URL格式
func (tp *TransferProcessor) isValidURL(url string) bool {
	// 简单的URL验证，可以根据需要扩展
	quarkPattern := `https://pan\.quark\.cn/s/[a-zA-Z0-9]+`
	matched, _ := regexp.MatchString(quarkPattern, url)
	return matched
}

// checkResourceExists 检查资源是否已存在
func (tp *TransferProcessor) checkResourceExists(url string) (bool, *entity.Resource, error) {
	// 根据URL查找资源
	resource, err := tp.repoMgr.ResourceRepository.GetByURL(url)
	if err != nil {
		// 如果是未找到记录的错误，则表示资源不存在
		if strings.Contains(err.Error(), "record not found") {
			return false, nil, nil
		}
		return false, nil, err
	}

	return true, resource, nil
}

// performTransfer 执行转存操作
func (tp *TransferProcessor) performTransfer(ctx context.Context, input *TransferInput) (uint, string, error) {
	// 解析URL获取分享信息
	shareInfo, err := tp.parseShareURL(input.URL)
	if err != nil {
		return 0, "", fmt.Errorf("解析分享链接失败: %v", err)
	}

	// 创建资源记录
	var categoryID *uint
	if input.CategoryID != 0 {
		categoryID = &input.CategoryID
	}

	resource := &entity.Resource{
		Title:      input.Title,
		URL:        input.URL,
		CategoryID: categoryID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存资源到数据库
	err = tp.repoMgr.ResourceRepository.Create(resource)
	if err != nil {
		return 0, "", fmt.Errorf("保存资源失败: %v", err)
	}

	// 添加标签关联
	if len(input.Tags) > 0 {
		err = tp.addResourceTags(resource.ID, input.Tags)
		if err != nil {
			utils.Error("添加资源标签失败: %v", err)
		}
	}

	// 执行实际转存操作
	saveURL, err := tp.transferToCloud(ctx, shareInfo)
	if err != nil {
		utils.Error("云端转存失败: %v", err)
		// 转存失败但资源已创建，返回原始URL
		return resource.ID, input.URL, nil
	}

	// 更新资源的转存链接
	if saveURL != "" {
		err = tp.repoMgr.ResourceRepository.UpdateSaveURL(resource.ID, saveURL)
		if err != nil {
			utils.Error("更新转存链接失败: %v", err)
		}
	}

	return resource.ID, saveURL, nil
}

// ShareInfo 分享信息结构
type ShareInfo struct {
	PanType string
	ShareID string
	URL     string
}

// parseShareURL 解析分享链接
func (tp *TransferProcessor) parseShareURL(url string) (*ShareInfo, error) {
	// 解析夸克网盘链接
	quarkPattern := `https://pan\.quark\.cn/s/([a-zA-Z0-9]+)`
	re := regexp.MustCompile(quarkPattern)
	matches := re.FindStringSubmatch(url)

	if len(matches) >= 2 {
		return &ShareInfo{
			PanType: "quark",
			ShareID: matches[1],
			URL:     url,
		}, nil
	}

	return nil, fmt.Errorf("不支持的分享链接格式: %s", url)
}

// addResourceTags 添加资源标签
func (tp *TransferProcessor) addResourceTags(resourceID uint, tagIDs []uint) error {
	for _, tagID := range tagIDs {
		// 创建资源标签关联
		resourceTag := &entity.ResourceTag{
			ResourceID: resourceID,
			TagID:      tagID,
		}

		err := tp.repoMgr.ResourceRepository.CreateResourceTag(resourceTag)
		if err != nil {
			return fmt.Errorf("创建资源标签关联失败: %v", err)
		}
	}
	return nil
}

// transferToCloud 执行云端转存
func (tp *TransferProcessor) transferToCloud(ctx context.Context, shareInfo *ShareInfo) (string, error) {
	// 检查是否启用自动转存
	autoTransferEnabled, err := tp.repoMgr.SystemConfigRepository.GetConfigBool("auto_transfer")
	if err != nil || !autoTransferEnabled {
		utils.Info("自动转存未启用，跳过云端转存")
		return "", nil
	}

	// TODO: 实现云端转存逻辑
	utils.Info("云端转存功能暂未实现，跳过转存: %s", shareInfo.ShareID)
	return "", nil
}
