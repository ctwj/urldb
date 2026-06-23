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
	PanID      uint   `json:"pan_id"`
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
	startTime := utils.GetCurrentTime()
	utils.InfoWithFields(map[string]interface{}{
		"task_item_id": item.ID,
		"task_id":      taskID,
	}, "开始处理转存任务项: %d", item.ID)

	// 解析输入数据
	parseStart := utils.GetCurrentTime()
	var input TransferInput
	if err := json.Unmarshal([]byte(item.InputData), &input); err != nil {
		parseDuration := time.Since(parseStart)
		utils.ErrorWithFields(map[string]interface{}{
			"error":       err.Error(),
			"duration_ms": parseDuration.Milliseconds(),
		}, "解析输入数据失败: %v，耗时: %v", err, parseDuration)
		return fmt.Errorf("解析输入数据失败: %v", err)
	}
	parseDuration := time.Since(parseStart)
	utils.DebugWithFields(map[string]interface{}{
		"duration_ms": parseDuration.Milliseconds(),
	}, "解析输入数据完成，耗时: %v", parseDuration)

	// 验证输入数据
	validateStart := utils.GetCurrentTime()
	if err := tp.validateInput(&input); err != nil {
		validateDuration := time.Since(validateStart)
		utils.Error("输入数据验证失败: %v，耗时: %v", err, validateDuration)
		return fmt.Errorf("输入数据验证失败: %v", err)
	}
	validateDuration := time.Since(validateStart)
	utils.DebugWithFields(map[string]interface{}{
		"duration_ms": validateDuration.Milliseconds(),
	}, "输入数据验证完成，耗时: %v", validateDuration)

	// 获取任务配置中的账号信息
	configStart := utils.GetCurrentTime()
	var selectedAccounts []uint
	task, err := tp.repoMgr.TaskRepository.GetByID(taskID)
	if err == nil && task.Config != "" {
		var taskConfig map[string]interface{}
		if err := json.Unmarshal([]byte(task.Config), &taskConfig); err == nil {
			if accounts, ok := taskConfig["selected_accounts"].([]interface{}); ok {
				for _, acc := range accounts {
					if accID, ok := acc.(float64); ok {
						selectedAccounts = append(selectedAccounts, uint(accID))
					}
				}
			}
		}
	}
	configDuration := time.Since(configStart)
	utils.Debug("获取任务配置完成，耗时: %v", configDuration)

	if len(selectedAccounts) == 0 {
		utils.Error("失败: %v", "没有指定转存账号")
	}

	// 检查资源是否已存在
	checkStart := utils.GetCurrentTime()
	exists, existingResource, err := tp.checkResourceExists(input.URL)
	checkDuration := time.Since(checkStart)
	if err != nil {
		utils.Error("检查资源是否存在失败: %v，耗时: %v", err, checkDuration)
	} else {
		utils.Debug("检查资源是否存在完成，耗时: %v", checkDuration)
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

			elapsedTime := time.Since(startTime)
			utils.Info("资源已存在且有转存链接，跳过转存: %s，总耗时: %v", input.Title, elapsedTime)
			return nil
		}
	}

	// 查询出 账号列表
	cksStart := utils.GetCurrentTime()
	cks, err := tp.repoMgr.CksRepository.FindByIds(selectedAccounts)
	cksDuration := time.Since(cksStart)
	if err != nil {
		utils.Error("读取账号失败: %v，耗时: %v", err, cksDuration)
	} else {
		utils.Debug("读取账号完成，账号数量: %d，耗时: %v", len(cks), cksDuration)
	}

	// 执行转存操作
	transferStart := utils.GetCurrentTime()
	resourceID, saveURL, err := tp.performTransfer(ctx, &input, cks, existingResource)
	transferDuration := time.Since(transferStart)
	if err != nil {
		// 转存失败，更新输出数据
		output := TransferOutput{
			Error:   err.Error(),
			Success: false,
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		elapsedTime := time.Since(startTime)
		utils.ErrorWithFields(map[string]interface{}{
			"task_item_id": item.ID,
			"error":        err.Error(),
			"duration_ms":  transferDuration.Milliseconds(),
			"total_ms":     elapsedTime.Milliseconds(),
		}, "转存任务项处理失败: %d, 错误: %v，转存耗时: %v，总耗时: %v", item.ID, err, transferDuration, elapsedTime)
		// performTransfer / transferToCloud 已在错误中带"转存失败:"前缀，不再重复包装，
		// 否则日志里会出现三层"转存失败: 转存失败: 转存失败:" 的丑陋嵌套。
		return err
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

		elapsedTime := time.Since(startTime)
		utils.Error("转存任务项处理失败: %d, 未获取到分享链接，总耗时: %v", item.ID, elapsedTime)
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

	elapsedTime := time.Since(startTime)
	utils.InfoWithFields(map[string]interface{}{
		"task_item_id":         item.ID,
		"resource_id":          resourceID,
		"save_url":             saveURL,
		"transfer_duration_ms": transferDuration.Milliseconds(),
		"total_duration_ms":    elapsedTime.Milliseconds(),
	}, "转存任务项处理完成: %d, 资源ID: %d, 转存链接: %s，转存耗时: %v，总耗时: %v", item.ID, resourceID, saveURL, transferDuration, elapsedTime)
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
	patterns := []string{
		`https://pan\.quark\.cn/s/[a-zA-Z0-9]+`,        // 夸克网盘
		`https://pan\.xunlei\.com/s/.+`,                // 迅雷网盘
		`https?://pan\.baidu\.com/s/[a-zA-Z0-9_-]+`,    // 百度网盘 /s/ 格式
		`https?://pan\.baidu\.com/share/init\?surl=.+`, // 百度网盘 /share/init?surl= 格式
	}
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}
	return false
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
// existing 非空表示系统中已存在该 URL 的 resource（重转场景）—— 此时走 Update 路径并写入
// transferred_at，允许自动清理调度回收网盘文件；existing 为空表示全新 URL，走 Create 路径，
// 不写 transferred_at，避免清理服务误删用户主动入库的资源。
func (tp *TransferProcessor) performTransfer(ctx context.Context, input *TransferInput, cks []*entity.Cks, existing *entity.Resource) (uint, string, error) {
	// 从 cks 中，挑选出，能够转存的账号，
	urlType := pan.ExtractServiceType(input.URL)
	if urlType == pan.NotFound {
		return 0, "", fmt.Errorf("未识别资源类型: %v", input.URL)
	}

	serviceType := ""
	switch urlType {
	case pan.Quark:
		serviceType = "quark"
	case pan.Xunlei:
		serviceType = "xunlei"
	case pan.BaiduPan:
		serviceType = "baidu"
	default:
		serviceType = ""
	}

	var account *entity.Cks
	for _, ck := range cks {
		if ck.ServiceType == serviceType {
			account = ck
		}
	}
	if account == nil {
		return 0, "", fmt.Errorf("为找到匹配的账号: %v", serviceType)
	}

	// 先执行转存操作
	saveData, err := tp.transferToCloud(ctx, input.URL, account)
	if err != nil {
		utils.Error("云端转存失败: %v", err)
		// transferToCloud 已带"转存失败:"前缀，直接透传避免多层嵌套
		return 0, "", err
	}

	// 验证转存链接是否有效
	if saveData.SaveURL == "" {
		utils.Error("转存成功但未获取到分享链接")
		return 0, "", fmt.Errorf("转存成功但未获取到分享链接")
	}

	// 转存成功，创建资源记录
	var categoryID *uint
	if input.CategoryID != 0 {
		categoryID = &input.CategoryID
	}

	// 确定平台ID  根据 serviceType 确认 panId
	panID, _ := tp.repoMgr.PanRepository.FindIdByServiceType(serviceType)
	panIdInt := uint(panID)

	// 绑定转存账号 ID，供 cleanup_service 解析删除文件所需的 cookie（FR-012 防跨账号误删）
	accountID := account.ID

	now := time.Now()

	// 分两种情况：
	// (1) existing != nil：系统中已存在该 URL（来自 telegram 抓取等），转存为中转备份，需要写入
	//     transferred_at，让清理调度器到期回收；走 Update 避免重复 Create。
	// (2) existing == nil：用户在数据转存管理界面主动粘贴的全新链接，转存结果就是最终落地，
	//     不应被自动清理 → 不写 transferred_at；走 Create。
	var resourceID uint
	if existing != nil {
		// 重转：更新现有 resource。只更新转存相关字段，避免覆盖 Title/Category/Tags 等已有信息。
		if err := tp.repoMgr.ResourceRepository.UpdateFields(existing.ID, map[string]interface{}{
			"save_url":      saveData.SaveURL,
			"fid":           saveData.Fid,
			"ck_id":         accountID,
			"transferred_at": now,
			"error_msg":     "",
			"updated_at":    now,
		}); err != nil {
			utils.Error("更新资源转存信息失败: %v", err)
			return 0, "", fmt.Errorf("更新资源失败: %v", err)
		}
		resourceID = existing.ID
		utils.Info("转存成功，资源已更新 - 资源ID: %d, 转存链接: %s", resourceID, saveData.SaveURL)
	} else {
		// 生成 6 位 Base62 唯一 key（供 /r/:key 短链访问）。失败时回退到无 key 继续，
		// 避免转存成功但资源落库失败的尴尬，运维可后续手动补 key。
		key, keyErr := tp.repoMgr.ResourceRepository.GenerateUniqueKey()
		if keyErr != nil {
			utils.Error("生成资源 key 失败（继续落库但无 key）: %v", keyErr)
		}

		resource := &entity.Resource{
			Title:      input.Title,
			URL:        input.URL,
			CategoryID: categoryID,
			PanID:      &panIdInt,        // 设置平台ID
			CkID:       &accountID,       // 绑定转存账号（cleanup 删除时据此解析 cookie）
			SaveURL:    saveData.SaveURL, // 直接设置转存链接
			Fid:        saveData.Fid,     // 记录转存文件ID（清理时依据）
			Key:        key,
			// 注意：不写 TransferredAt —— 新建资源不应被自动清理
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tp.repoMgr.ResourceRepository.Create(resource); err != nil {
			utils.Error("保存转存成功的资源失败: %v", err)
			return 0, "", fmt.Errorf("保存资源失败: %v", err)
		}
		resourceID = resource.ID

		// 添加标签关联（仅新建场景；重转时原 resource 已有 tags）
		if len(input.Tags) > 0 {
			if err := tp.addResourceTags(resourceID, input.Tags); err != nil {
				utils.Error("添加资源标签失败: %v", err)
				// 标签添加失败不影响资源创建，只记录错误
			}
		}

		utils.Info("转存成功，资源已创建（不纳入自动清理）- 资源ID: %d, 转存链接: %s", resourceID, saveData.SaveURL)
	}

	return resourceID, saveData.SaveURL, nil
}

// ShareInfo 分享信息结构
type ShareInfo struct {
	PanType string
	ShareID string
	URL     string
}

// // parseShareURL 解析分享链接
// func (tp *TransferProcessor) parseShareURL(url string) (*ShareInfo, error) {
// 	// 解析夸克网盘链接
// 	quarkPattern := `https://pan\.quark\.cn/s/([a-zA-Z0-9]+)`
// 	re := regexp.MustCompile(quarkPattern)
// 	matches := re.FindStringSubmatch(url)

// 	if len(matches) >= 2 {
// 		return &ShareInfo{
// 			PanType: "quark",
// 			ShareID: matches[1],
// 			URL:     url,
// 		}, nil
// 	}

// 	return nil, fmt.Errorf("不支持的分享链接格式: %s", url)
// }

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
func (tp *TransferProcessor) transferToCloud(ctx context.Context, url string, account *entity.Cks) (*TransferResult, error) {

	// 创建网盘服务工厂
	factory := pan.NewPanFactory()

	service, err := factory.CreatePanService(url, &pan.PanConfig{
		URL:         url,
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	service.SetCKSRepository(tp.repoMgr.CksRepository, *account)

	// 提取分享ID
	shareID, _ := pan.ExtractShareId(url)

	// 执行转存
	transferResult, err := service.Transfer(shareID) // 有些链接还需要其他信息从 url 中自行解析
	if err != nil {
		utils.Error("转存失败: %v", err)
		return nil, fmt.Errorf("转存失败: %v", err)
	}

	if transferResult == nil || !transferResult.Success {
		errMsg := "转存失败"
		if transferResult != nil && transferResult.Message != "" {
			errMsg = transferResult.Message
		}
		return nil, fmt.Errorf("转存失败: %v", errMsg)
	}

	// 提取转存链接
	var saveURL string
	var fid string

	if data, ok := transferResult.Data.(map[string]interface{}); ok {
		if v, ok := data["shareUrl"]; ok {
			saveURL, _ = v.(string)
		}
		if v, ok := data["fid"]; ok {
			fid, _ = v.(string)
		}
	}
	if saveURL == "" {
		saveURL = transferResult.ShareURL
	}

	if saveURL == "" {
		return nil, fmt.Errorf("转存失败: %v", "转存成功但未获取到分享链接")
	}

	utils.Info("转存成功 - 原分享fid: %s, 转存链接: %s", fid, saveURL)

	return &TransferResult{
		Success: true,
		SaveURL: saveURL,
		Fid:     fid,
	}, nil

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
	Fid      string `json:"fid`
	ErrorMsg string `json:"error_msg"`
}
