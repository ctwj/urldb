package scheduler

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	panutils "github.com/ctwj/urldb/common"
	commonutils "github.com/ctwj/urldb/common/utils"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// AutoTransferScheduler 自动转存调度器
type AutoTransferScheduler struct {
	*BaseScheduler
	autoTransferRunning bool
	autoTransferMutex   sync.Mutex // 防止自动转存任务重叠执行
}

// NewAutoTransferScheduler 创建自动转存调度器
func NewAutoTransferScheduler(base *BaseScheduler) *AutoTransferScheduler {
	return &AutoTransferScheduler{
		BaseScheduler:       base,
		autoTransferRunning: false,
		autoTransferMutex:   sync.Mutex{},
	}
}

// Start 启动自动转存定时任务
func (a *AutoTransferScheduler) Start() {

	// 自动转存已经放弃，不再自动缓存
	return

	if a.autoTransferRunning {
		utils.Info("自动转存定时任务已在运行中")
		return
	}

	a.autoTransferRunning = true
	utils.Info("启动自动转存定时任务")

	go func() {
		// 获取系统配置中的间隔时间
		interval := 5 * time.Minute // 默认5分钟
		if autoProcessInterval, err := a.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoProcessInterval); err == nil && autoProcessInterval > 0 {
			interval = time.Duration(autoProcessInterval) * time.Minute
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		utils.Info(fmt.Sprintf("自动转存定时任务已启动，间隔时间: %v", interval))

		// 立即执行一次
		a.processAutoTransfer()

		for {
			select {
			case <-ticker.C:
				// 使用TryLock防止任务重叠执行
				if a.autoTransferMutex.TryLock() {
					go func() {
						defer a.autoTransferMutex.Unlock()
						a.processAutoTransfer()
					}()
				} else {
					utils.Info("上一次自动转存任务还在执行中，跳过本次执行")
				}
			case <-a.GetStopChan():
				utils.Info("停止自动转存定时任务")
				return
			}
		}
	}()
}

// Stop 停止自动转存定时任务
func (a *AutoTransferScheduler) Stop() {
	if !a.autoTransferRunning {
		utils.Info("自动转存定时任务未在运行")
		return
	}

	a.GetStopChan() <- true
	a.autoTransferRunning = false
	utils.Info("已发送停止信号给自动转存定时任务")
}

// IsAutoTransferRunning 检查自动转存任务是否正在运行
func (a *AutoTransferScheduler) IsAutoTransferRunning() bool {
	return a.autoTransferRunning
}

// processAutoTransfer 处理自动转存
func (a *AutoTransferScheduler) processAutoTransfer() {
	utils.Info("开始处理自动转存...")

	// 检查系统配置，确认是否启用自动转存
	autoTransferEnabled, err := a.systemConfigRepo.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
	if err != nil {
		utils.Error(fmt.Sprintf("获取系统配置失败: %v", err))
		return
	}

	if !autoTransferEnabled {
		utils.Info("自动转存功能已禁用")
		return
	}

	// 获取quark平台ID
	quarkPanID, err := a.getQuarkPanID()
	if err != nil {
		utils.Error(fmt.Sprintf("获取夸克网盘ID失败: %v", err))
		return
	}

	// 获取所有账号
	accounts, err := a.cksRepo.FindAll()
	if err != nil {
		utils.Error(fmt.Sprintf("获取网盘账号失败: %v", err))
		return
	}

	// 获取最小存储空间配置
	autoTransferMinSpace, err := a.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		utils.Error(fmt.Sprintf("获取最小存储空间配置失败: %v", err))
		return
	}

	// 过滤：只保留已激活、quark平台、剩余空间足够的账号
	minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024
	var validAccounts []entity.Cks
	for _, acc := range accounts {
		if acc.IsValid && acc.PanID == quarkPanID && acc.LeftSpace >= minSpaceBytes {
			validAccounts = append(validAccounts, acc)
		}
	}

	if len(validAccounts) == 0 {
		utils.Info("没有可用的quark网盘账号")
		return
	}

	utils.Info(fmt.Sprintf("找到 %d 个可用quark网盘账号，开始自动转存处理...", len(validAccounts)))

	// 计算处理数量限制
	// 假设每5秒转存一个资源，每分钟20个，5分钟100个
	// 根据时间间隔和账号数量计算大致的处理数量
	interval := 5 * time.Minute // 默认5分钟
	if autoProcessInterval, err := a.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoProcessInterval); err == nil && autoProcessInterval > 0 {
		interval = time.Duration(autoProcessInterval) * time.Minute
	}

	// 计算每分钟能处理的资源数量：账号数 * 12（每分钟12个，即每5秒一个）
	resourcesPerMinute := len(validAccounts) * 12
	// 根据时间间隔计算总处理数量
	maxProcessCount := int(float64(resourcesPerMinute) * interval.Minutes())

	utils.Info(fmt.Sprintf("时间间隔: %v, 账号数: %d, 每分钟处理能力: %d, 最大处理数量: %d",
		interval, len(validAccounts), resourcesPerMinute, maxProcessCount))

	// 获取需要转存的资源（限制数量）
	resources, err := a.getResourcesForTransfer(quarkPanID, maxProcessCount)
	if err != nil {
		utils.Error(fmt.Sprintf("获取需要转存的资源失败: %v", err))
		return
	}

	if len(resources) == 0 {
		utils.Info("没有需要转存的资源")
		return
	}

	utils.Info(fmt.Sprintf("找到 %d 个需要转存的资源", len(resources)))

	// 获取违禁词配置
	forbiddenWords, err := a.systemConfigRepo.GetConfigValue(entity.ConfigKeyForbiddenWords)
	if err != nil {
		utils.Error(fmt.Sprintf("获取违禁词配置失败: %v", err))
		forbiddenWords = "" // 如果获取失败，使用空字符串
	}

	// 过滤包含违禁词的资源，并标记违禁词错误
	var filteredResources []*entity.Resource
	var forbiddenResources []*entity.Resource

	if forbiddenWords != "" {
		words := strings.Split(forbiddenWords, ",")
		// 清理违禁词数组，去除空格
		var cleanWords []string
		for _, word := range words {
			word = strings.TrimSpace(word)
			if word != "" {
				cleanWords = append(cleanWords, word)
			}
		}

		for _, resource := range resources {
			shouldSkip := false
			var matchedWords []string
			title := strings.ToLower(resource.Title)
			description := strings.ToLower(resource.Description)

			for _, word := range cleanWords {
				wordLower := strings.ToLower(word)
				if strings.Contains(title, wordLower) || strings.Contains(description, wordLower) {
					matchedWords = append(matchedWords, word)
					shouldSkip = true
				}
			}

			if shouldSkip {
				// 标记为违禁词错误
				resource.ErrorMsg = fmt.Sprintf("存在违禁词 (共 %d 个)", len(matchedWords))
				forbiddenResources = append(forbiddenResources, resource)
				utils.Info(fmt.Sprintf("标记违禁词资源: %s (包含 %d 个违禁词)", resource.Title, len(matchedWords)))
			} else {
				filteredResources = append(filteredResources, resource)
			}
		}
		utils.Info(fmt.Sprintf("违禁词过滤后，剩余 %d 个资源需要转存，违禁词资源 %d 个", len(filteredResources), len(forbiddenResources)))
	} else {
		filteredResources = resources
	}

	// 注意：资源数量已在数据库查询时限制，无需再次限制

	// 保存违禁词资源的错误信息
	for _, resource := range forbiddenResources {
		if err := a.resourceRepo.Update(resource); err != nil {
			utils.Error(fmt.Sprintf("保存违禁词错误信息失败 (ID: %d): %v", resource.ID, err))
		}
	}

	// 并发自动转存
	resourceCh := make(chan *entity.Resource, len(filteredResources))
	for _, res := range filteredResources {
		resourceCh <- res
	}
	close(resourceCh)

	var wg sync.WaitGroup
	for _, account := range validAccounts {
		wg.Add(1)
		go func(acc entity.Cks) {
			defer wg.Done()
			factory := panutils.GetInstance() // 使用单例模式
			for res := range resourceCh {
				if err := a.transferResource(res, []entity.Cks{acc}, factory); err != nil {
					utils.Error(fmt.Sprintf("转存资源失败 (ID: %d): %v", res.ID, err))
				} else {
					utils.Info(fmt.Sprintf("成功转存资源: %s", res.Title))
					rand.Seed(utils.GetCurrentTime().UnixNano())
					sleepSec := rand.Intn(3) + 1 // 1,2,3
					time.Sleep(time.Duration(sleepSec) * time.Second)
				}
			}
		}(account)
	}
	wg.Wait()
	utils.Info(fmt.Sprintf("自动转存处理完成，账号数: %d，处理资源数: %d，违禁词资源数: %d",
		len(validAccounts), len(filteredResources), len(forbiddenResources)))
}

// getQuarkPanID 获取夸克网盘ID
func (a *AutoTransferScheduler) getQuarkPanID() (uint, error) {
	// 获取panRepo的实现，以便访问数据库
	panRepoImpl, ok := a.panRepo.(interface{ GetDB() *gorm.DB })
	if !ok {
		return 0, fmt.Errorf("panRepo不支持GetDB方法")
	}

	var quarkPan entity.Pan
	err := panRepoImpl.GetDB().Where("name = ?", "quark").First(&quarkPan).Error
	if err != nil {
		return 0, fmt.Errorf("未找到quark平台: %v", err)
	}

	return quarkPan.ID, nil
}

// getResourcesForTransfer 获取需要转存的资源
func (a *AutoTransferScheduler) getResourcesForTransfer(quarkPanID uint, limit int) ([]*entity.Resource, error) {
	// 获取最近24小时内的资源
	sinceTime := utils.GetCurrentTime().Add(-24 * time.Hour)

	// 使用资源仓库的方法获取需要转存的资源
	repoImpl, ok := a.resourceRepo.(*repo.ResourceRepositoryImpl)
	if !ok {
		return nil, fmt.Errorf("资源仓库类型错误")
	}

	return repoImpl.GetResourcesForTransfer(quarkPanID, sinceTime, limit)
}

// transferResource 转存单个资源
func (a *AutoTransferScheduler) transferResource(resource *entity.Resource, accounts []entity.Cks, factory *panutils.PanFactory) error {
	if len(accounts) == 0 {
		return fmt.Errorf("没有可用的网盘账号")
	}
	account := accounts[0]

	service, err := factory.CreatePanService(resource.URL, &panutils.PanConfig{
		URL:         resource.URL,
		ExpiredType: 0,
		IsType:      0,
		Cookie:      account.Ck,
	})
	if err != nil {
		return fmt.Errorf("创建网盘服务失败: %v", err)
	}

	// 获取最小存储空间配置
	autoTransferMinSpace, err := a.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		utils.Error(fmt.Sprintf("获取最小存储空间配置失败: %v", err))
		return err
	}

	// 检查账号剩余空间
	minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024
	if account.LeftSpace < minSpaceBytes {
		return fmt.Errorf("账号剩余空间不足，需要 %d GB，当前剩余 %d GB", autoTransferMinSpace, account.LeftSpace/1024/1024/1024)
	}

	// 提取分享ID
	shareID, _ := commonutils.ExtractShareIdString(resource.URL)

	// 转存资源
	result, err := service.Transfer(shareID)
	if err != nil {
		// 更新错误信息
		resource.ErrorMsg = err.Error()
		a.resourceRepo.Update(resource)
		return fmt.Errorf("转存失败: %v", err)
	}

	if result == nil || !result.Success {
		errMsg := "转存失败"
		if result != nil && result.Message != "" {
			errMsg = result.Message
		}
		// 更新错误信息
		resource.ErrorMsg = errMsg
		a.resourceRepo.Update(resource)
		return fmt.Errorf("转存失败: %s", errMsg)
	}

	// 提取转存链接、fid等
	var saveURL, fid string
	if data, ok := result.Data.(map[string]interface{}); ok {
		if v, ok := data["shareUrl"]; ok {
			saveURL, _ = v.(string)
		}
		if v, ok := data["fid"]; ok {
			fid, _ = v.(string)
		}
	}
	if saveURL == "" {
		saveURL = result.ShareURL
	}

	// 更新资源信息
	resource.SaveURL = saveURL
	resource.CkID = &account.ID
	resource.Fid = fid
	resource.ErrorMsg = ""

	// 保存更新
	err = a.resourceRepo.Update(resource)
	if err != nil {
		return fmt.Errorf("保存转存结果失败: %v", err)
	}

	return nil
}

// selectBestAccount 选择最佳账号
func (a *AutoTransferScheduler) selectBestAccount(accounts []entity.Cks) *entity.Cks {
	if len(accounts) == 0 {
		return nil
	}

	// 获取最小存储空间配置
	autoTransferMinSpace, err := a.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		utils.Error(fmt.Sprintf("获取最小存储空间配置失败: %v", err))
		return &accounts[0] // 返回第一个账号
	}

	minSpaceBytes := int64(autoTransferMinSpace) * 1024 * 1024 * 1024

	var bestAccount *entity.Cks
	var bestScore int64 = -1

	for i := range accounts {
		account := &accounts[i]
		if account.LeftSpace < minSpaceBytes {
			continue // 跳过空间不足的账号
		}

		score := a.calculateAccountScore(account)
		if score > bestScore {
			bestScore = score
			bestAccount = account
		}
	}

	return bestAccount
}

// calculateAccountScore 计算账号评分
func (a *AutoTransferScheduler) calculateAccountScore(account *entity.Cks) int64 {
	// TODO: 实现账号评分算法
	// 1. VIP账号加分
	// 2. 剩余空间大的账号加分
	// 3. 使用率低的账号加分
	// 4. 可以根据历史使用情况调整评分

	score := int64(0)

	// VIP账号加分
	if account.VipStatus {
		score += 1000
	}

	// 剩余空间加分（每GB加1分）
	score += account.LeftSpace / (1024 * 1024 * 1024)

	// 使用率加分（使用率越低分数越高）
	if account.Space > 0 {
		usageRate := float64(account.UsedSpace) / float64(account.Space)
		score += int64((1 - usageRate) * 500) // 使用率越低，加分越多
	}

	return score
}
