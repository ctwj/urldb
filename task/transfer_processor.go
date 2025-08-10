package task

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	pan "github.com/ctwj/urldb/common"
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
		// 检查已存在的资源是否有有效的转存链接
		if existingResource.SaveURL == "" {
			// 资源存在但没有转存链接，需要重新转存
			utils.Info("资源已存在但无转存链接，重新转存: %s", input.Title)
		} else {
			// 资源已存在且有转存链接，跳过转存
			output := TransferOutput{
				ResourceID: existingResource.ID,
				SaveURL:    existingResource.SaveURL,
				Success:    true,
				Time:       utils.GetCurrentTimeString(),
			}

			outputJSON, _ := json.Marshal(output)
			item.OutputData = string(outputJSON)

			utils.Info("资源已存在且有转存链接，跳过转存: %s", input.Title)
			return nil
		}
	}

	// 执行转存操作
	resourceID, saveURL, err := tp.performTransfer(ctx, &input)
	if err != nil {
		// 转存失败，更新输出数据
		output := TransferOutput{
			Error:   err.Error(),
			Success: false,
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("转存任务项处理失败: %d, 错误: %v", item.ID, err)
		return fmt.Errorf("转存失败: %v", err)
	}

	// 验证转存结果
	if saveURL == "" {
		output := TransferOutput{
			Error:   "转存成功但未获取到分享链接",
			Success: false,
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("转存任务项处理失败: %d, 未获取到分享链接", item.ID)
		return fmt.Errorf("转存成功但未获取到分享链接")
	}

	// 转存成功，更新输出数据
	output := TransferOutput{
		ResourceID: resourceID,
		SaveURL:    saveURL,
		Success:    true,
		Time:       utils.GetCurrentTimeString(),
	}

	outputJSON, _ := json.Marshal(output)
	item.OutputData = string(outputJSON)

	utils.Info("转存任务项处理完成: %d, 资源ID: %d, 转存链接: %s", item.ID, resourceID, saveURL)
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

	// 先执行转存操作
	saveURL, err := tp.transferToCloud(ctx, shareInfo)
	if err != nil {
		utils.Error("云端转存失败: %v", err)
		return 0, "", fmt.Errorf("转存失败: %v", err)
	}

	// 验证转存链接是否有效
	if saveURL == "" {
		utils.Error("转存成功但未获取到分享链接")
		return 0, "", fmt.Errorf("转存成功但未获取到分享链接")
	}

	// 转存成功，创建资源记录
	var categoryID *uint
	if input.CategoryID != 0 {
		categoryID = &input.CategoryID
	}

	resource := &entity.Resource{
		Title:      input.Title,
		URL:        input.URL,
		CategoryID: categoryID,
		SaveURL:    saveURL, // 直接设置转存链接
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存资源到数据库
	err = tp.repoMgr.ResourceRepository.Create(resource)
	if err != nil {
		utils.Error("保存转存成功的资源失败: %v", err)
		return 0, "", fmt.Errorf("保存资源失败: %v", err)
	}

	// 添加标签关联
	if len(input.Tags) > 0 {
		err = tp.addResourceTags(resource.ID, input.Tags)
		if err != nil {
			utils.Error("添加资源标签失败: %v", err)
			// 标签添加失败不影响资源创建，只记录错误
		}
	}

	utils.Info("转存成功，资源已创建 - 资源ID: %d, 转存链接: %s", resource.ID, saveURL)
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
	// 转存任务独立于自动转存开关，直接执行转存逻辑
	// 获取转存相关的配置（如最小存储空间等），但不检查自动转存开关

	// 获取夸克平台ID
	quarkPanID, err := tp.getQuarkPanID()
	if err != nil {
		return "", fmt.Errorf("获取夸克平台ID失败: %v", err)
	}

	// 获取可用的夸克账号
	accounts, err := tp.repoMgr.CksRepository.FindAll()
	if err != nil {
		return "", fmt.Errorf("获取网盘账号失败: %v", err)
	}

	// 获取最小存储空间配置（转存任务需要关注此配置）
	autoTransferMinSpace, err := tp.repoMgr.SystemConfigRepository.GetConfigInt("auto_transfer_min_space")
	if err != nil {
		utils.Error("获取最小存储空间配置失败: %v", err)
		autoTransferMinSpace = 5 // 默认5GB
	}

	// 过滤：只保留已激活、夸克平台、剩余空间足够的账号
	minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024
	var validAccounts []entity.Cks
	for _, acc := range accounts {
		if acc.IsValid && acc.PanID == quarkPanID && acc.LeftSpace >= minSpaceBytes {
			validAccounts = append(validAccounts, acc)
		}
	}

	if len(validAccounts) == 0 {
		return "", fmt.Errorf("没有可用的夸克网盘账号（需要剩余空间 >= %d GB）", autoTransferMinSpace)
	}

	utils.Info("找到 %d 个可用夸克网盘账号，开始转存处理...", len(validAccounts))

	// 使用第一个可用账号进行转存
	account := validAccounts[0]

	// 创建网盘服务工厂
	factory := pan.NewPanFactory()

	// 执行转存
	result := tp.transferSingleResource(shareInfo, account, factory)
	if !result.Success {
		return "", fmt.Errorf("转存失败: %s", result.ErrorMsg)
	}

	return result.SaveURL, nil
}

// getQuarkPanID 获取夸克网盘ID
func (tp *TransferProcessor) getQuarkPanID() (uint, error) {
	// 通过FindAll方法查找所有平台，然后过滤出quark平台
	pans, err := tp.repoMgr.PanRepository.FindAll()
	if err != nil {
		return 0, fmt.Errorf("查询平台信息失败: %v", err)
	}

	for _, p := range pans {
		if p.Name == "quark" {
			return p.ID, nil
		}
	}

	return 0, fmt.Errorf("未找到quark平台")
}

// TransferResult 转存结果
type TransferResult struct {
	Success  bool   `json:"success"`
	SaveURL  string `json:"save_url"`
	ErrorMsg string `json:"error_msg"`
}

// transferSingleResource 转存单个资源
func (tp *TransferProcessor) transferSingleResource(shareInfo *ShareInfo, account entity.Cks, factory *pan.PanFactory) TransferResult {
	utils.Info("开始转存资源 - 分享ID: %s, 账号: %s", shareInfo.ShareID, account.Username)

	service, err := factory.CreatePanService(shareInfo.URL, &pan.PanConfig{
		URL:         shareInfo.URL,
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	if err != nil {
		utils.Error("创建网盘服务失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("创建网盘服务失败: %v", err),
		}
	}

	// 执行转存
	transferResult, err := service.Transfer(shareInfo.ShareID)
	if err != nil {
		utils.Error("转存失败: %v", err)
		return TransferResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("转存失败: %v", err),
		}
	}

	if transferResult == nil || !transferResult.Success {
		errMsg := "转存失败"
		if transferResult != nil && transferResult.Message != "" {
			errMsg = transferResult.Message
		}
		utils.Error("转存失败: %s", errMsg)
		return TransferResult{
			Success:  false,
			ErrorMsg: errMsg,
		}
	}

	// 提取转存链接
	var saveURL string
	if data, ok := transferResult.Data.(map[string]interface{}); ok {
		if v, ok := data["shareUrl"]; ok {
			saveURL, _ = v.(string)
		}
	}
	if saveURL == "" {
		saveURL = transferResult.ShareURL
	}

	// 验证转存链接是否有效
	if saveURL == "" {
		utils.Error("转存成功但未获取到分享链接 - 分享ID: %s", shareInfo.ShareID)
		return TransferResult{
			Success:  false,
			ErrorMsg: "转存成功但未获取到分享链接",
		}
	}

	// 验证链接格式
	if !strings.HasPrefix(saveURL, "http") {
		utils.Error("转存链接格式无效 - 分享ID: %s, 链接: %s", shareInfo.ShareID, saveURL)
		return TransferResult{
			Success:  false,
			ErrorMsg: "转存链接格式无效",
		}
	}

	utils.Info("转存成功 - 分享ID: %s, 转存链接: %s", shareInfo.ShareID, saveURL)

	return TransferResult{
		Success: true,
		SaveURL: saveURL,
	}
}
