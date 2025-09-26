package task

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	pan "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// ExpansionProcessor 扩容任务处理器
type ExpansionProcessor struct {
	repoMgr *repo.RepositoryManager
}

// NewExpansionProcessor 创建扩容任务处理器
func NewExpansionProcessor(repoMgr *repo.RepositoryManager) *ExpansionProcessor {
	return &ExpansionProcessor{
		repoMgr: repoMgr,
	}
}

// GetTaskType 获取任务类型
func (ep *ExpansionProcessor) GetTaskType() string {
	return "expansion"
}

// ExpansionInput 扩容任务输入数据结构
type ExpansionInput struct {
	PanAccountID uint                   `json:"pan_account_id"`
	DataSource   map[string]interface{} `json:"data_source,omitempty"`
}

// TransferredResource 转存成功的资源信息
type TransferredResource struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// ExpansionOutput 扩容任务输出数据结构
type ExpansionOutput struct {
	Success              bool                  `json:"success"`
	Message              string                `json:"message"`
	Error                string                `json:"error,omitempty"`
	Time                 string                `json:"time"`
	TransferredResources []TransferredResource `json:"transferred_resources,omitempty"`
}

// Process 处理扩容任务项
func (ep *ExpansionProcessor) Process(ctx context.Context, taskID uint, item *entity.TaskItem) error {
	utils.Info("开始处理扩容任务项: %d", item.ID)

	// 解析输入数据
	var input ExpansionInput
	if err := json.Unmarshal([]byte(item.InputData), &input); err != nil {
		return fmt.Errorf("解析输入数据失败: %v", err)
	}

	// 验证输入数据
	if err := ep.validateInput(&input); err != nil {
		return fmt.Errorf("输入数据验证失败: %v", err)
	}

	// 检查账号是否已经扩容过
	exists, err := ep.checkExpansionExists(input.PanAccountID)
	if err != nil {
		utils.Error("检查扩容记录失败: %v", err)
		return fmt.Errorf("检查扩容记录失败: %v", err)
	}

	if exists {
		output := ExpansionOutput{
			Success: false,
			Message: "账号已扩容过",
			Error:   "每个账号只能扩容一次",
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Info("账号已扩容过，跳过扩容: 账号ID %d", input.PanAccountID)
		return fmt.Errorf("账号已扩容过")
	}

	// 检查账号类型（只支持quark账号）
	if err := ep.checkAccountType(input.PanAccountID); err != nil {
		output := ExpansionOutput{
			Success: false,
			Message: "账号类型不支持扩容",
			Error:   err.Error(),
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("账号类型不支持扩容: %v", err)
		return err
	}

	// 执行扩容操作（传入数据源）
	transferred, err := ep.performExpansion(ctx, input.PanAccountID, input.DataSource)
	if err != nil {
		output := ExpansionOutput{
			Success: false,
			Message: "扩容失败",
			Error:   err.Error(),
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("扩容任务项处理失败: %d, 错误: %v", item.ID, err)
		return fmt.Errorf("扩容失败: %v", err)
	}

	// 扩容成功
	output := ExpansionOutput{
		Success:              true,
		Message:              "扩容成功",
		Time:                 utils.GetCurrentTimeString(),
		TransferredResources: transferred,
	}

	outputJSON, _ := json.Marshal(output)
	item.OutputData = string(outputJSON)

	utils.Info("扩容任务项处理完成: %d, 账号ID: %d", item.ID, input.PanAccountID)
	return nil
}

// validateInput 验证输入数据
func (ep *ExpansionProcessor) validateInput(input *ExpansionInput) error {
	if input.PanAccountID == 0 {
		return fmt.Errorf("账号ID不能为空")
	}
	return nil
}

// checkExpansionExists 检查账号是否已经扩容过
func (ep *ExpansionProcessor) checkExpansionExists(panAccountID uint) (bool, error) {
	// 查询所有expansion类型的任务
	tasks, _, err := ep.repoMgr.TaskRepository.GetList(1, 1000, "expansion", "completed")
	if err != nil {
		return false, fmt.Errorf("获取扩容任务列表失败: %v", err)
	}

	// 检查每个任务的配置中是否包含该账号ID
	for _, task := range tasks {
		if task.Config != "" {
			var taskConfig map[string]interface{}
			if err := json.Unmarshal([]byte(task.Config), &taskConfig); err == nil {
				if configAccountID, ok := taskConfig["pan_account_id"].(float64); ok {
					if uint(configAccountID) == panAccountID {
						// 找到了该账号的扩容任务，检查任务状态
						if task.Status == "completed" {
							// 如果任务已完成，说明已经扩容过
							return true, nil
						}
					}
				}
			}
		}
	}

	return false, nil
}

// checkAccountType 检查账号类型（只支持quark账号）
func (ep *ExpansionProcessor) checkAccountType(panAccountID uint) error {
	// 获取账号信息
	cks, err := ep.repoMgr.CksRepository.FindByID(panAccountID)
	if err != nil {
		return fmt.Errorf("获取账号信息失败: %v", err)
	}

	// 检查是否为quark账号
	if cks.ServiceType != "quark" {
		return fmt.Errorf("只支持quark账号扩容，当前账号类型: %s", cks.ServiceType)
	}

	return nil
}

// performExpansion 执行扩容操作
func (ep *ExpansionProcessor) performExpansion(ctx context.Context, panAccountID uint, dataSource map[string]interface{}) ([]TransferredResource, error) {
	rand.Seed(time.Now().UnixNano())
	utils.Info("执行扩容操作，账号ID: %d, 数据源: %v", panAccountID, dataSource)

	transferred := []TransferredResource{}

	// 获取账号信息
	account, err := ep.repoMgr.CksRepository.FindByID(panAccountID)
	if err != nil {
		return nil, fmt.Errorf("获取账号信息失败: %v", err)
	}

	// 创建网盘服务工厂
	factory := pan.NewPanFactory()
	service, err := factory.CreatePanServiceByType(pan.Quark, &pan.PanConfig{
		URL:         "",
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	if err != nil {
		return nil, fmt.Errorf("创建网盘服务失败: %v", err)
	}
	service.SetCKSRepository(ep.repoMgr.CksRepository, *account)

	// 定义扩容分类列表（按优先级排序）
	categories := []string{
		"情色", "喜剧", "动作", "科幻", "动画", "悬疑", "犯罪", "惊悚",
		"冒险", "恐怖", "战争", "传记", "剧情", "爱情", "家庭", "儿童",
		"音乐", "历史", "奇幻", "歌舞", "武侠", "灾难", "西部", "古装", "运动",
	}

	// 获取数据源类型
	dataSourceType := "internal"
	var thirdPartyURL string
	if dataSource != nil {
		if dsType, ok := dataSource["type"].(string); ok {
			dataSourceType = dsType
			if dsType == "third-party" {
				if url, ok := dataSource["url"].(string); ok {
					thirdPartyURL = url
				}
			}
		}
	}

	utils.Info("使用数据源类型: %s", dataSourceType)

	totalTransferred := 0
	totalFailed := 0

	// 逐个处理分类
	for _, category := range categories {
		utils.Info("开始处理分类: %s", category)

		// 获取该分类的资源
		resources, err := ep.getHotResources(category)
		if err != nil {
			utils.Error("获取分类 %s 的资源失败: %v", category, err)
			continue
		}

		if len(resources) == 0 {
			utils.Info("分类 %s 没有可用资源，跳过", category)
			continue
		}

		utils.Info("分类 %s 获取到 %d 个资源", category, len(resources))

		// 转存该分类的资源（限制每个分类最多转存20个）
		maxPerCategory := 20
		transferredCount := 0

		for _, resource := range resources {
			if transferredCount >= maxPerCategory {
				break
			}

			// 检查是否还有存储空间
			hasSpace, err := ep.checkStorageSpace(service, &account.Ck)
			if err != nil {
				utils.Error("检查存储空间失败: %v", err)
				return transferred, fmt.Errorf("检查存储空间失败: %v", err)
			}

			if !hasSpace {
				utils.Info("存储空间不足，停止扩容，但保存已转存的资源")
				// 存储空间不足时，停止继续转存，但返回已转存的资源作为成功结果
				break
			}

			// 获取资源 , dataSourceType, thirdPartyURL
			resource, err := ep.getResourcesByHot(resource, dataSourceType, thirdPartyURL, *account, service)
			if resource == nil || err != nil {
				if resource != nil {
					utils.Error("获取资源失败: %s, 错误: %v", resource.Title, err)
				} else {
					utils.Error("获取资源失败, 错误: %v", err)
				}
				totalFailed++
				continue
			}

			// 执行转存
			saveURL, err := ep.transferResource(ctx, service, resource)
			if err != nil {
				utils.Error("转存资源失败: %s, 错误: %v", resource.Title, err)
				totalFailed++
				continue
			}

			// 随机休眠1-3秒，避免请求过于频繁
			sleepDuration := time.Duration(rand.Intn(3)+1) * time.Second
			time.Sleep(sleepDuration)

			// 保存转存结果到任务输出
			transferred = append(transferred, TransferredResource{
				Title: resource.Title,
				URL:   saveURL,
			})

			totalTransferred++
			transferredCount++
			utils.Info("成功转存资源: %s -> %s", resource.Title, saveURL)

			// 每转存5个资源检查一次存储空间
			if totalTransferred%5 == 0 {
				utils.Info("已转存 %d 个资源，检查存储空间", totalTransferred)
			}
		}

		utils.Info("分类 %s 处理完成，转存 %d 个资源", category, transferredCount)
	}

	utils.Info("扩容完成，总共转存: %d 个资源，失败: %d 个资源", totalTransferred, totalFailed)
	return transferred, nil
}

// getResourcesForCategory 获取指定分类的资源
func (ep *ExpansionProcessor) getResourcesByHot(
	resource *entity.HotDrama, dataSourceType,
	thirdPartyURL string,
	entity entity.Cks,
	service pan.PanService,
) (*entity.Resource, error) {
	if dataSourceType == "third-party" && thirdPartyURL != "" {
		// 从第三方API获取资源
		return ep.getResourcesFromThirdPartyAPI(resource, thirdPartyURL)
	}

	// 从内部数据库获取资源
	return ep.getResourcesFromInternalDB(resource, entity, service)
}

// getResourcesFromInternalDB 根据 HotDrama 的title 获取数据库中资源，并且资源的类型和 account 的资源类型一致
func (ep *ExpansionProcessor) getResourcesFromInternalDB(HotDrama *entity.HotDrama, account entity.Cks, service pan.PanService) (*entity.Resource, error) {
	// 获取账号对应的平台ID
	panIDInt, err := ep.repoMgr.PanRepository.FindIdByServiceType(account.ServiceType)
	if err != nil {
		return nil, fmt.Errorf("获取平台ID失败: %v", err)
	}
	panID := uint(panIDInt)

	// 1. 搜索标题
	params := map[string]interface{}{
		"search":    HotDrama.Title,
		"pan_id":    panID,
		"is_valid":  true,
		"page":      1,
		"page_size": 10,
	}
	resources, _, err := ep.repoMgr.ResourceRepository.SearchWithFilters(params)
	if err != nil {
		return nil, fmt.Errorf("搜索资源失败: %v", err)
	}

	// 检查结果是否有效，通过服务验证
	for _, res := range resources {
		if res.IsValid && res.URL != "" {
			// 使用服务验证资源是否可转存
			shareID, _ := pan.ExtractShareId(res.URL)
			if shareID != "" {
				result, err := service.Transfer(shareID)
				if err == nil && result != nil && result.Success {
					return &res, nil
				}
			}
		}
	}

	// 3. 没有有效资源，返回错误信息
	return nil, fmt.Errorf("未找到有效的资源")
}

// getResourcesFromInternalDB 从内部数据库获取资源
func (ep *ExpansionProcessor) getHotResources(category string) ([]*entity.HotDrama, error) {
	// 获取该分类下sub_type为"排行"的资源
	dramas, _, err := ep.repoMgr.HotDramaRepository.FindByCategoryAndSubType(category, "排行", 1, 20)
	if err != nil {
		return nil, fmt.Errorf("获取分类 %s 的资源失败: %v", category, err)
	}

	// 如果没有找到"排行"类型的资源，尝试获取该分类下的所有资源
	if len(dramas) == 0 {
		dramas, _, err = ep.repoMgr.HotDramaRepository.FindByCategory(category, 1, 20)
		if err != nil {
			return nil, fmt.Errorf("获取分类 %s 的资源失败: %v", category, err)
		}
	}

	// 转换为指针数组
	result := make([]*entity.HotDrama, len(dramas))
	for i := range dramas {
		result[i] = &dramas[i]
	}

	return result, nil
}

// getResourcesFromThirdPartyAPI 从第三方API获取资源
func (ep *ExpansionProcessor) getResourcesFromThirdPartyAPI(resource *entity.HotDrama, apiURL string) (*entity.Resource, error) {
	// 构建API请求URL，添加分类参数
	requestURL := fmt.Sprintf("%s?category=%s&limit=20", apiURL, resource)

	// 发送HTTP请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("请求第三方API失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("第三方API返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应数据（假设API返回JSON格式的资源列表）
	var apiResponse struct {
		Data []*entity.HotDrama `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, fmt.Errorf("解析第三方API响应失败: %v", err)
	}

	return nil, nil
}

// checkStorageSpace 检查存储空间是否足够
func (ep *ExpansionProcessor) checkStorageSpace(service pan.PanService, ck *string) (bool, error) {
	userInfo, err := service.GetUserInfo(ck)
	if err != nil {
		utils.Error("获取用户信息失败: %v", err)
		// 如果无法获取用户信息，假设还有空间继续
		return true, nil
	}

	// 检查是否还有足够的空间（保留至少10GB空间）
	const reservedSpaceGB = 100
	reservedSpaceBytes := int64(reservedSpaceGB * 1024 * 1024 * 1024)

	if userInfo.TotalSpace-userInfo.UsedSpace <= reservedSpaceBytes {
		utils.Info("存储空间不足，已使用: %d bytes，总容量: %d bytes",
			userInfo.UsedSpace, userInfo.TotalSpace)
		return false, nil
	}

	return true, nil
}

// transferResource 执行单个资源的转存
func (ep *ExpansionProcessor) transferResource(ctx context.Context, service pan.PanService, res *entity.Resource) (string, error) {
	// 如果没有URL，跳过转存
	if res.URL == "" {
		return "", fmt.Errorf("资源 %s 没有有效的URL", res.URL)
	}

	// 提取分享ID
	shareID, _ := pan.ExtractShareId(res.URL)
	if shareID == "" {
		return "", fmt.Errorf("无法从URL %s 提取分享ID", res.URL)
	}

	// 执行转存
	result, err := service.Transfer(shareID)
	if err != nil {
		return "", fmt.Errorf("转存失败: %v", err)
	}

	if result == nil || !result.Success {
		errorMsg := "转存失败"
		if result != nil {
			errorMsg = result.Message
		}
		return "", fmt.Errorf("转存失败: %s", errorMsg)
	}

	// 提取转存链接
	var saveURL string
	if result.Data != nil {
		if data, ok := result.Data.(map[string]interface{}); ok {
			if v, ok := data["shareUrl"]; ok {
				saveURL, _ = v.(string)
			}
		}
	}
	if saveURL == "" {
		saveURL = result.ShareURL
	}

	if saveURL == "" {
		return "", fmt.Errorf("转存成功但未获取到分享链接")
	}

	return saveURL, nil
}

// recordTransferredResource 记录转存成功的资源
func (ep *ExpansionProcessor) recordTransferredResource(drama *entity.HotDrama, accountID uint, saveURL string) error {
	// 获取夸克网盘的平台ID
	panIDInt, err := ep.repoMgr.PanRepository.FindIdByServiceType("quark")
	if err != nil {
		utils.Error("获取夸克网盘平台ID失败: %v", err)
		return err
	}

	// 转换为uint
	panID := uint(panIDInt)

	// 创建资源记录
	resource := &entity.Resource{
		Title:     drama.Title,
		URL:       drama.PosterURL,
		SaveURL:   saveURL,
		PanID:     &panID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsValid:   true,
		IsPublic:  false, // 扩容资源默认不公开
	}

	// 保存到数据库
	err = ep.repoMgr.ResourceRepository.Create(resource)
	if err != nil {
		return fmt.Errorf("保存资源记录失败: %v", err)
	}

	utils.Info("成功记录转存资源: %s (ID: %d)", drama.Title, resource.ID)
	return nil
}
