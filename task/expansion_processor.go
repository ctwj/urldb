package task

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
	startTime := utils.GetCurrentTime()
	utils.InfoWithFields(map[string]interface{}{
		"task_item_id": item.ID,
		"task_id":      taskID,
	}, "开始处理扩容任务项: %d", item.ID)

	// 解析输入数据
	parseStart := utils.GetCurrentTime()
	var input ExpansionInput
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
	if err := ep.validateInput(&input); err != nil {
		validateDuration := time.Since(validateStart)
		utils.Error("输入数据验证失败: %v，耗时: %v", err, validateDuration)
		return fmt.Errorf("输入数据验证失败: %v", err)
	}
	validateDuration := time.Since(validateStart)
	utils.Debug("输入数据验证完成，耗时: %v", validateDuration)

	// 检查账号是否已经扩容过
	checkExpansionStart := utils.GetCurrentTime()
	exists, err := ep.checkExpansionExists(input.PanAccountID)
	checkExpansionDuration := time.Since(checkExpansionStart)
	if err != nil {
		utils.Error("检查扩容记录失败: %v，耗时: %v", err, checkExpansionDuration)
		return fmt.Errorf("检查扩容记录失败: %v", err)
	}
	utils.Debug("检查扩容记录完成，耗时: %v", checkExpansionDuration)

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
	checkAccountTypeStart := utils.GetCurrentTime()
	if err := ep.checkAccountType(input.PanAccountID); err != nil {
		checkAccountTypeDuration := time.Since(checkAccountTypeStart)
		output := ExpansionOutput{
			Success: false,
			Message: "账号类型不支持扩容",
			Error:   err.Error(),
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("账号类型不支持扩容: %v，耗时: %v", err, checkAccountTypeDuration)
		return err
	}
	checkAccountTypeDuration := time.Since(checkAccountTypeStart)
	utils.Debug("检查账号类型完成，耗时: %v", checkAccountTypeDuration)

	// 执行扩容操作（传入数据源）
	expansionStart := utils.GetCurrentTime()
	transferred, err := ep.performExpansion(ctx, input.PanAccountID, input.DataSource)
	expansionDuration := time.Since(expansionStart)
	if err != nil {
		output := ExpansionOutput{
			Success: false,
			Message: "扩容失败",
			Error:   err.Error(),
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("扩容任务项处理失败: %d, 错误: %v，总耗时: %v", item.ID, err, expansionDuration)
		return fmt.Errorf("扩容失败: %v", err)
	}
	utils.Debug("扩容操作完成，耗时: %v", expansionDuration)

	// 扩容成功
	output := ExpansionOutput{
		Success:              true,
		Message:              "扩容成功",
		Time:                 utils.GetCurrentTimeString(),
		TransferredResources: transferred,
	}

	outputJSON, _ := json.Marshal(output)
	item.OutputData = string(outputJSON)

	elapsedTime := time.Since(startTime)
	utils.InfoWithFields(map[string]interface{}{
		"task_item_id": item.ID,
		"account_id":   input.PanAccountID,
		"duration_ms":  elapsedTime.Milliseconds(),
	}, "扩容任务项处理完成: %d, 账号ID: %d, 总耗时: %v", item.ID, input.PanAccountID, elapsedTime)
	return nil
}

// validateInput 验证输入数据
func (ep *ExpansionProcessor) validateInput(input *ExpansionInput) error {
	startTime := utils.GetCurrentTime()

	if input.PanAccountID == 0 {
		utils.Error("账号ID验证失败，账号ID不能为空，耗时: %v", time.Since(startTime))
		return fmt.Errorf("账号ID不能为空")
	}

	utils.Debug("输入数据验证完成，耗时: %v", time.Since(startTime))
	return nil
}

// checkExpansionExists 检查账号是否已经扩容过
func (ep *ExpansionProcessor) checkExpansionExists(panAccountID uint) (bool, error) {
	startTime := utils.GetCurrentTime()

	// 查询所有expansion类型的任务
	tasksStart := utils.GetCurrentTime()
	tasks, _, err := ep.repoMgr.TaskRepository.GetList(1, 1000, "expansion", "completed")
	tasksDuration := time.Since(tasksStart)
	if err != nil {
		utils.Error("获取扩容任务列表失败: %v，耗时: %v", err, tasksDuration)
		return false, fmt.Errorf("获取扩容任务列表失败: %v", err)
	}
	utils.Debug("获取扩容任务列表完成，找到 %d 个任务，耗时: %v", len(tasks), tasksDuration)

	// 检查每个任务的配置中是否包含该账号ID
	checkStart := utils.GetCurrentTime()
	for _, task := range tasks {
		if task.Config != "" {
			var taskConfig map[string]interface{}
			if err := json.Unmarshal([]byte(task.Config), &taskConfig); err == nil {
				if configAccountID, ok := taskConfig["pan_account_id"].(float64); ok {
					if uint(configAccountID) == panAccountID {
						// 找到了该账号的扩容任务，检查任务状态
						if task.Status == "completed" {
							// 如果任务已完成，说明已经扩容过
							checkDuration := time.Since(checkStart)
							utils.Debug("检查扩容记录完成，账号已扩容，耗时: %v", checkDuration)
							return true, nil
						}
					}
				}
			}
		}
	}
	checkDuration := time.Since(checkStart)
	utils.Debug("检查扩容记录完成，账号未扩容，耗时: %v", checkDuration)

	totalDuration := time.Since(startTime)
	utils.Debug("检查扩容记录完成，账号未扩容，总耗时: %v", totalDuration)
	return false, nil
}

// checkAccountType 检查账号类型（只支持quark账号）
func (ep *ExpansionProcessor) checkAccountType(panAccountID uint) error {
	startTime := utils.GetCurrentTime()

	// 获取账号信息
	accountStart := utils.GetCurrentTime()
	cks, err := ep.repoMgr.CksRepository.FindByID(panAccountID)
	accountDuration := time.Since(accountStart)
	if err != nil {
		utils.Error("获取账号信息失败: %v，耗时: %v", err, accountDuration)
		return fmt.Errorf("获取账号信息失败: %v", err)
	}
	utils.Debug("获取账号信息完成，耗时: %v", accountDuration)

	// 检查是否为quark账号
	serviceCheckStart := utils.GetCurrentTime()
	if cks.ServiceType != "quark" {
		serviceCheckDuration := time.Since(serviceCheckStart)
		utils.Error("账号类型检查失败，当前账号类型: %s，耗时: %v", cks.ServiceType, serviceCheckDuration)
		return fmt.Errorf("只支持quark账号扩容，当前账号类型: %s", cks.ServiceType)
	}
	serviceCheckDuration := time.Since(serviceCheckStart)
	utils.Debug("账号类型检查完成，为quark账号，耗时: %v", serviceCheckDuration)

	totalDuration := time.Since(startTime)
	utils.Debug("账号类型检查完成，总耗时: %v", totalDuration)
	return nil
}

// performExpansion 执行扩容操作
func (ep *ExpansionProcessor) performExpansion(ctx context.Context, panAccountID uint, dataSource map[string]interface{}) ([]TransferredResource, error) {
	rand.Seed(time.Now().UnixNano())
	startTime := utils.GetCurrentTime()
	utils.Info("执行扩容操作，账号ID: %d, 数据源: %v", panAccountID, dataSource)

	transferred := []TransferredResource{}

	// 获取账号信息
	accountStart := utils.GetCurrentTime()
	account, err := ep.repoMgr.CksRepository.FindByID(panAccountID)
	accountDuration := time.Since(accountStart)
	if err != nil {
		utils.Error("获取账号信息失败: %v，耗时: %v", err, accountDuration)
		return nil, fmt.Errorf("获取账号信息失败: %v", err)
	}
	utils.Debug("获取账号信息完成，耗时: %v", accountDuration)

	// 创建网盘服务工厂
	serviceStart := utils.GetCurrentTime()
	factory := pan.NewPanFactory()
	service, err := factory.CreatePanServiceByType(pan.Quark, &pan.PanConfig{
		URL:         "",
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	serviceDuration := time.Since(serviceStart)
	if err != nil {
		utils.Error("创建网盘服务失败: %v，耗时: %v", err, serviceDuration)
		return nil, fmt.Errorf("创建网盘服务失败: %v", err)
	}
	service.SetCKSRepository(ep.repoMgr.CksRepository, *account)
	utils.Debug("创建网盘服务完成，耗时: %v", serviceDuration)

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
		resourcesStart := utils.GetCurrentTime()
		resources, err := ep.getHotResources(category)
		resourcesDuration := time.Since(resourcesStart)
		if err != nil {
			utils.Error("获取分类 %s 的资源失败: %v，耗时: %v", category, err, resourcesDuration)
			continue
		}
		utils.Debug("获取分类 %s 的资源完成，耗时: %v", category, resourcesDuration)

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
			storageCheckStart := utils.GetCurrentTime()
			hasSpace, err := ep.checkStorageSpace(service, &account.Ck)
			storageCheckDuration := time.Since(storageCheckStart)
			if err != nil {
				utils.Error("检查存储空间失败: %v，耗时: %v", err, storageCheckDuration)
				return transferred, fmt.Errorf("检查存储空间失败: %v", err)
			}
			utils.Debug("检查存储空间完成，耗时: %v", storageCheckDuration)

			if !hasSpace {
				utils.Info("存储空间不足，停止扩容，但保存已转存的资源")
				// 存储空间不足时，停止继续转存，但返回已转存的资源作为成功结果
				break
			}

			// 获取资源 , dataSourceType, thirdPartyURL
			resourceGetStart := utils.GetCurrentTime()
			resource, err := ep.getResourcesByHot(resource, dataSourceType, thirdPartyURL, *account, service)
			resourceGetDuration := time.Since(resourceGetStart)
			if resource == nil || err != nil {
				if resource != nil {
					utils.Error("获取资源失败: %s, 错误: %v，耗时: %v", resource.Title, err, resourceGetDuration)
				} else {
					utils.Error("获取资源失败, 错误: %v，耗时: %v", err, resourceGetDuration)
				}
				totalFailed++
				continue
			}
			utils.Debug("获取资源完成，耗时: %v", resourceGetDuration)

			// 执行转存
			transferStart := utils.GetCurrentTime()
			saveURL, err := ep.transferResource(ctx, service, resource, *account)
			transferDuration := time.Since(transferStart)
			if err != nil {
				utils.Error("转存资源失败: %s, 错误: %v，耗时: %v", resource.Title, err, transferDuration)
				totalFailed++
				continue
			}
			utils.Debug("转存资源完成，耗时: %v", transferDuration)

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

	elapsedTime := time.Since(startTime)
	utils.Info("扩容完成，总共转存: %d 个资源，失败: %d 个资源，总耗时: %v", totalTransferred, totalFailed, elapsedTime)
	return transferred, nil
}

// getResourcesForCategory 获取指定分类的资源
func (ep *ExpansionProcessor) getResourcesByHot(
	resource *entity.HotDrama, dataSourceType,
	thirdPartyURL string,
	entity entity.Cks,
	service pan.PanService,
) (*entity.Resource, error) {
	startTime := utils.GetCurrentTime()

	if dataSourceType == "third-party" && thirdPartyURL != "" {
		// 从第三方API获取资源
		thirdPartyStart := utils.GetCurrentTime()
		result, err := ep.getResourcesFromThirdPartyAPI(resource, thirdPartyURL)
		thirdPartyDuration := time.Since(thirdPartyStart)
		utils.Debug("从第三方API获取资源完成，耗时: %v", thirdPartyDuration)
		return result, err
	}

	// 从内部数据库获取资源
	internalStart := utils.GetCurrentTime()
	result, err := ep.getResourcesFromInternalDB(resource, entity, service)
	internalDuration := time.Since(internalStart)
	utils.Debug("从内部数据库获取资源完成，耗时: %v", internalDuration)

	totalDuration := time.Since(startTime)
	utils.Debug("获取资源完成: %s，总耗时: %v", resource.Title, totalDuration)
	return result, err
}

// getResourcesFromInternalDB 根据 HotDrama 的title 获取数据库中资源，并且资源的类型和 account 的资源类型一致
func (ep *ExpansionProcessor) getResourcesFromInternalDB(HotDrama *entity.HotDrama, account entity.Cks, service pan.PanService) (*entity.Resource, error) {
	startTime := utils.GetCurrentTime()

	// 修改配置 isType = 1 只检测，不转存
	configStart := utils.GetCurrentTime()
	service.UpdateConfig(&pan.PanConfig{
		URL:         "",
		ExpiredType: 0,
		IsType:      1,
		Cookie:      account.Ck,
	})
	utils.Debug("更新服务配置完成，耗时: %v", time.Since(configStart))
	panID := account.PanID

	// 1. 搜索标题
	searchStart := utils.GetCurrentTime()
	params := map[string]interface{}{
		"search":    HotDrama.Title,
		"pan_id":    panID,
		"is_valid":  true,
		"page":      1,
		"page_size": 10,
	}
	resources, _, err := ep.repoMgr.ResourceRepository.SearchWithFilters(params)
	searchDuration := time.Since(searchStart)
	if err != nil {
		utils.Error("搜索资源失败: %v，耗时: %v", err, searchDuration)
		return nil, fmt.Errorf("搜索资源失败: %v", err)
	}
	utils.Debug("搜索资源完成，找到 %d 个资源，耗时: %v", len(resources), searchDuration)

	// 检查结果是否有效，通过服务验证
	validateStart := utils.GetCurrentTime()
	for _, res := range resources {
		if res.IsValid && res.URL != "" {
			// 使用服务验证资源是否可转存
			shareID, _ := pan.ExtractShareId(res.URL)
			if shareID != "" {
				result, err := service.Transfer(shareID)
				if err == nil && result != nil && result.Success {
					validateDuration := time.Since(validateStart)
					utils.Debug("验证资源成功: %s，耗时: %v", res.Title, validateDuration)
					return &res, nil
				}
			}
		}
	}
	validateDuration := time.Since(validateStart)
	utils.Debug("验证资源完成，未找到有效资源，耗时: %v", validateDuration)

	totalDuration := time.Since(startTime)
	utils.Debug("从内部数据库获取资源完成: %s，总耗时: %v", HotDrama.Title, totalDuration)
	// 3. 没有有效资源，返回错误信息
	return nil, fmt.Errorf("未找到有效的资源")
}

// getResourcesFromInternalDB 从内部数据库获取资源
func (ep *ExpansionProcessor) getHotResources(category string) ([]*entity.HotDrama, error) {
	startTime := utils.GetCurrentTime()

	// 获取该分类下sub_type为"排行"的资源
	rankedStart := utils.GetCurrentTime()
	dramas, _, err := ep.repoMgr.HotDramaRepository.FindByCategoryAndSubType(category, "排行", 1, 20)
	rankedDuration := time.Since(rankedStart)
	if err != nil {
		utils.Error("获取分类 %s 的排行资源失败: %v，耗时: %v", category, err, rankedDuration)
		return nil, fmt.Errorf("获取分类 %s 的资源失败: %v", category, err)
	}
	utils.Debug("获取分类 %s 的排行资源完成，找到 %d 个资源，耗时: %v", category, len(dramas), rankedDuration)

	// 如果没有找到"排行"类型的资源，尝试获取该分类下的所有资源
	if len(dramas) == 0 {
		allStart := utils.GetCurrentTime()
		dramas, _, err = ep.repoMgr.HotDramaRepository.FindByCategory(category, 1, 20)
		allDuration := time.Since(allStart)
		if err != nil {
			utils.Error("获取分类 %s 的所有资源失败: %v，耗时: %v", category, err, allDuration)
			return nil, fmt.Errorf("获取分类 %s 的资源失败: %v", category, err)
		}
		utils.Debug("获取分类 %s 的所有资源完成，找到 %d 个资源，耗时: %v", category, len(dramas), allDuration)
	}

	// 转换为指针数组
	convertStart := utils.GetCurrentTime()
	result := make([]*entity.HotDrama, len(dramas))
	for i := range dramas {
		result[i] = &dramas[i]
	}
	convertDuration := time.Since(convertStart)
	utils.Debug("转换资源数组完成，耗时: %v", convertDuration)

	totalDuration := time.Since(startTime)
	utils.Debug("获取热门资源完成: 分类 %s，总数 %d，总耗时: %v", category, len(result), totalDuration)
	return result, nil
}

// getResourcesFromThirdPartyAPI 从第三方API获取资源
func (ep *ExpansionProcessor) getResourcesFromThirdPartyAPI(resource *entity.HotDrama, apiURL string) (*entity.Resource, error) {
	// 构建API请求URL，添加分类参数
	// requestURL := fmt.Sprintf("%s?category=%s&limit=20", apiURL, resource)

	// TODO 使用第三方API接口，请求资源

	return nil, nil
}

// checkStorageSpace 检查存储空间是否足够
func (ep *ExpansionProcessor) checkStorageSpace(service pan.PanService, ck *string) (bool, error) {
	startTime := utils.GetCurrentTime()

	userInfoStart := utils.GetCurrentTime()
	userInfo, err := service.GetUserInfo(ck)
	userInfoDuration := time.Since(userInfoStart)
	if err != nil {
		utils.Error("获取用户信息失败: %v，耗时: %v", err, userInfoDuration)
		// 如果无法获取用户信息，假设还有空间继续
		return true, nil
	}
	utils.Debug("获取用户信息完成，耗时: %v", userInfoDuration)

	// 检查是否还有足够的空间（保留至少10GB空间）
	const reservedSpaceGB = 100
	reservedSpaceBytes := int64(reservedSpaceGB * 1024 * 1024 * 1024)

	if userInfo.TotalSpace-userInfo.UsedSpace <= reservedSpaceBytes {
		utils.Info("存储空间不足，已使用: %d bytes，总容量: %d bytes，检查耗时: %v",
			userInfo.UsedSpace, userInfo.TotalSpace, time.Since(startTime))
		return false, nil
	}

	totalDuration := time.Since(startTime)
	utils.Debug("存储空间检查完成，有足够空间，耗时: %v", totalDuration)
	return true, nil
}

// transferResource 执行单个资源的转存
func (ep *ExpansionProcessor) transferResource(ctx context.Context, service pan.PanService, res *entity.Resource, account entity.Cks) (string, error) {
	startTime := utils.GetCurrentTime()

	// 修改配置 isType = 0 转存
	configStart := utils.GetCurrentTime()
	service.UpdateConfig(&pan.PanConfig{
		URL:         "",
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	utils.Debug("更新服务配置完成，耗时: %v", time.Since(configStart))

	// 如果没有URL，跳过转存
	if res.URL == "" {
		utils.Error("资源 %s 没有有效的URL", res.URL)
		return "", fmt.Errorf("资源 %s 没有有效的URL", res.URL)
	}

	// 提取分享ID
	extractStart := utils.GetCurrentTime()
	shareID, _ := pan.ExtractShareId(res.URL)
	extractDuration := time.Since(extractStart)
	if shareID == "" {
		utils.Error("无法从URL %s 提取分享ID，耗时: %v", res.URL, extractDuration)
		return "", fmt.Errorf("无法从URL %s 提取分享ID", res.URL)
	}
	utils.Debug("提取分享ID完成: %s，耗时: %v", shareID, extractDuration)

	// 执行转存
	transferStart := utils.GetCurrentTime()
	result, err := service.Transfer(shareID)
	transferDuration := time.Since(transferStart)
	if err != nil {
		utils.Error("转存失败: %v，耗时: %v", err, transferDuration)
		return "", fmt.Errorf("转存失败: %v", err)
	}

	if result == nil || !result.Success {
		errorMsg := "转存失败"
		if result != nil {
			errorMsg = result.Message
		}
		utils.Error("转存结果失败: %s，耗时: %v", errorMsg, time.Since(transferStart))
		return "", fmt.Errorf("转存失败: %s", errorMsg)
	}

	// 提取转存链接
	extractURLStart := utils.GetCurrentTime()
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
		extractURLDuration := time.Since(extractURLStart)
		utils.Error("转存成功但未获取到分享链接，耗时: %v", extractURLDuration)
		return "", fmt.Errorf("转存成功但未获取到分享链接")
	}

	totalDuration := time.Since(startTime)
	utils.Debug("转存资源完成: %s -> %s，总耗时: %v", res.Title, saveURL, totalDuration)
	return saveURL, nil
}

// recordTransferredResource 记录转存成功的资源
// func (ep *ExpansionProcessor) recordTransferredResource(drama *entity.HotDrama, accountID uint, saveURL string) error {
// 	// 获取夸克网盘的平台ID
// 	panIDInt, err := ep.repoMgr.PanRepository.FindIdByServiceType("quark")
// 	if err != nil {
// 		utils.Error("获取夸克网盘平台ID失败: %v", err)
// 		return err
// 	}

// 	// 转换为uint
// 	panID := uint(panIDInt)

// 	// 创建资源记录
// 	resource := &entity.Resource{
// 		Title:     drama.Title,
// 		URL:       drama.PosterURL,
// 		SaveURL:   saveURL,
// 		PanID:     &panID,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 		IsValid:   true,
// 		IsPublic:  false, // 扩容资源默认不公开
// 	}

// 	// 保存到数据库
// 	err = ep.repoMgr.ResourceRepository.Create(resource)
// 	if err != nil {
// 		return fmt.Errorf("保存资源记录失败: %v", err)
// 	}

// 	utils.Info("成功记录转存资源: %s (ID: %d)", drama.Title, resource.ID)
// 	return nil
// }
