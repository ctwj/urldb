package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"time"

	panutils "github.com/ctwj/urldb/common"
	commonutils "github.com/ctwj/urldb/common/utils"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// ReadyResourceScheduler 待处理资源调度器
type ReadyResourceScheduler struct {
	*BaseScheduler
	readyResourceRunning bool
	processingMutex      sync.Mutex // 防止ready_resource任务重叠执行
}

// NewReadyResourceScheduler 创建待处理资源调度器
func NewReadyResourceScheduler(base *BaseScheduler) *ReadyResourceScheduler {
	return &ReadyResourceScheduler{
		BaseScheduler:        base,
		readyResourceRunning: false,
		processingMutex:      sync.Mutex{},
	}
}

// Start 启动待处理资源定时任务
func (r *ReadyResourceScheduler) Start() {
	if r.readyResourceRunning {
		utils.Info("待处理资源自动处理任务已在运行中")
		return
	}

	r.readyResourceRunning = true
	utils.Info("启动待处理资源自动处理任务")

	go func() {
		// 获取系统配置中的间隔时间
		interval := 3 * time.Minute // 默认3分钟
		if autoProcessInterval, err := r.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoProcessInterval); err == nil && autoProcessInterval > 0 {
			interval = time.Duration(autoProcessInterval) * time.Minute
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		utils.Info(fmt.Sprintf("待处理资源自动处理任务已启动，间隔时间: %v", interval))

		// 立即执行一次
		r.processReadyResources()

		for {
			select {
			case <-ticker.C:
				// 使用TryLock防止任务重叠执行
				if r.processingMutex.TryLock() {
					go func() {
						defer r.processingMutex.Unlock()
						r.processReadyResources()
					}()
				} else {
					utils.Info("上一次待处理资源任务还在执行中，跳过本次执行")
				}
			case <-r.GetStopChan():
				utils.Info("停止待处理资源自动处理任务")
				return
			}
		}
	}()
}

// Stop 停止待处理资源定时任务
func (r *ReadyResourceScheduler) Stop() {
	if !r.readyResourceRunning {
		utils.Info("待处理资源自动处理任务未在运行")
		return
	}

	r.GetStopChan() <- true
	r.readyResourceRunning = false
	utils.Info("已发送停止信号给待处理资源自动处理任务")
}

// IsReadyResourceRunning 检查待处理资源任务是否正在运行
func (r *ReadyResourceScheduler) IsReadyResourceRunning() bool {
	return r.readyResourceRunning
}

// processReadyResources 处理待处理资源
func (r *ReadyResourceScheduler) processReadyResources() {
	utils.Info("开始处理待处理资源...")

	// 检查系统配置，确认是否启用自动处理
	autoProcess, err := r.systemConfigRepo.GetConfigBool(entity.ConfigKeyAutoProcessReadyResources)
	if err != nil {
		utils.Error(fmt.Sprintf("获取系统配置失败: %v", err))
		return
	}

	if !autoProcess {
		utils.Info("自动处理待处理资源功能已禁用")
		return
	}

	// 获取所有没有错误的待处理资源
	readyResources, err := r.readyResourceRepo.FindAll()
	// readyResources, err := r.readyResourceRepo.FindWithoutErrors()
	if err != nil {
		utils.Error(fmt.Sprintf("获取待处理资源失败: %v", err))
		return
	}

	if len(readyResources) == 0 {
		utils.Info("没有待处理的资源")
		return
	}

	utils.Info(fmt.Sprintf("找到 %d 个待处理资源，开始处理...", len(readyResources)))

	processedCount := 0
	factory := panutils.GetInstance() // 使用单例模式
	for _, readyResource := range readyResources {

		//readyResource.URL 是 查重
		exits, err := r.resourceRepo.FindExists(readyResource.URL)
		if err != nil {
			utils.Error(fmt.Sprintf("查重失败: %v", err))
			continue
		}
		if exits {
			utils.Info(fmt.Sprintf("资源已存在: %s", readyResource.URL))
			r.readyResourceRepo.Delete(readyResource.ID)
			continue
		}

		if err := r.convertReadyResourceToResource(readyResource, factory); err != nil {
			utils.Error(fmt.Sprintf("处理资源失败 (ID: %d): %v", readyResource.ID, err))

			// 保存完整的错误信息
			readyResource.ErrorMsg = err.Error()

			if updateErr := r.readyResourceRepo.Update(&readyResource); updateErr != nil {
				utils.Error(fmt.Sprintf("更新错误信息失败 (ID: %d): %v", readyResource.ID, updateErr))
			} else {
				utils.Info(fmt.Sprintf("已保存错误信息到资源 (ID: %d): %s", readyResource.ID, err.Error()))
			}

			// 处理失败后删除资源，避免重复处理
			r.readyResourceRepo.Delete(readyResource.ID)
		} else {
			// 处理成功，删除readyResource
			r.readyResourceRepo.Delete(readyResource.ID)
			processedCount++
			utils.Info(fmt.Sprintf("成功处理资源: %s", readyResource.URL))
		}
	}

	utils.Info(fmt.Sprintf("待处理资源处理完成，共处理 %d 个资源", processedCount))
}

// convertReadyResourceToResource 将待处理资源转换为正式资源
func (r *ReadyResourceScheduler) convertReadyResourceToResource(readyResource entity.ReadyResource, factory *panutils.PanFactory) error {
	utils.Debug(fmt.Sprintf("开始处理资源: %s", readyResource.URL))

	// 提取分享ID和服务类型
	shareID, serviceType := panutils.ExtractShareId(readyResource.URL)
	if serviceType == panutils.NotFound {
		utils.Warn(fmt.Sprintf("不支持的链接地址: %s", readyResource.URL))
		return fmt.Errorf("不支持的链接地址: %s", readyResource.URL)
	}

	utils.Debug(fmt.Sprintf("检测到服务类型: %s, 分享ID: %s", serviceType.String(), shareID))

	resource := &entity.Resource{
		Title:       derefString(readyResource.Title),
		Description: readyResource.Description,
		URL:         readyResource.URL,
		Cover:       readyResource.Img,
		IsValid:     true,
		IsPublic:    true,
		Key:         readyResource.Key,
		PanID:       r.getPanIDByServiceType(serviceType),
	}

	// 检查违禁词
	forbiddenWords, err := r.systemConfigRepo.GetConfigValue(entity.ConfigKeyForbiddenWords)
	if err == nil && forbiddenWords != "" {
		words := strings.Split(forbiddenWords, ",")
		var matchedWords []string
		title := strings.ToLower(resource.Title)
		description := strings.ToLower(resource.Description)

		for _, word := range words {
			word = strings.TrimSpace(word)
			if word != "" {
				wordLower := strings.ToLower(word)
				if strings.Contains(title, wordLower) || strings.Contains(description, wordLower) {
					matchedWords = append(matchedWords, word)
				}
			}
		}

		if len(matchedWords) > 0 {
			utils.Warn(fmt.Sprintf("资源包含违禁词: %s, 违禁词: %s", resource.Title, strings.Join(matchedWords, ", ")))
			return fmt.Errorf("存在违禁词: %s", strings.Join(matchedWords, ", "))
		}
	}

	// 不是夸克，直接保存
	if serviceType != panutils.Quark {
		// 检测是否有效
		checkResult, err := commonutils.CheckURL(readyResource.URL)
		if err != nil {
			utils.Error(fmt.Sprintf("链接检查失败: %v", err))
			return fmt.Errorf("链接检查失败: %v", err)
		}
		if !checkResult.Status {
			utils.Warn(fmt.Sprintf("链接无效: %s", readyResource.URL))
			return fmt.Errorf("链接无效: %s", readyResource.URL)
		}
	} else {
		// 获取夸克网盘账号的 cookie
		panID := r.getPanIDByServiceType(serviceType)
		if panID == nil {
			utils.Error("未找到对应的平台ID")
			return fmt.Errorf("未找到对应的平台ID")
		}

		accounts, err := r.cksRepo.FindByPanID(*panID)
		if err != nil {
			utils.Error(fmt.Sprintf("获取夸克网盘账号失败: %v", err))
			return fmt.Errorf("获取网盘账号失败: %v", err)
		}

		if len(accounts) == 0 {
			utils.Error("没有可用的夸克网盘账号")
			return fmt.Errorf("没有可用的夸克网盘账号")
		}

		// 选择第一个有效的账号
		var selectedAccount *entity.Cks
		for _, account := range accounts {
			if account.IsValid {
				selectedAccount = &account
				break
			}
		}

		if selectedAccount == nil {
			utils.Error("没有有效的夸克网盘账号")
			return fmt.Errorf("没有有效的夸克网盘账号")
		}

		utils.Debug(fmt.Sprintf("使用夸克网盘账号: %d, Cookie: %s", selectedAccount.ID, selectedAccount.Ck[:20]+"..."))

		// 准备配置
		config := &panutils.PanConfig{
			URL:         readyResource.URL,
			Code:        "", // 可以从readyResource中获取
			IsType:      1,  // 转存并分享后的资源信息  0 转存后分享， 1 只获取基本信息
			ExpiredType: 1,  // 永久分享
			AdFid:       "",
			Stoken:      "",
			Cookie:      selectedAccount.Ck, // 添加 cookie
		}

		// 通过工厂获取对应的网盘服务单例
		panService, err := factory.CreatePanService(readyResource.URL, config)
		if err != nil {
			utils.Error(fmt.Sprintf("获取网盘服务失败: %v", err))
			return fmt.Errorf("获取网盘服务失败: %v", err)
		}

		// 统一处理：尝试转存获取标题
		result, err := panService.Transfer(shareID)
		if err != nil {
			utils.Error(fmt.Sprintf("网盘信息获取失败: %v", err))
			return fmt.Errorf("网盘信息获取失败: %v", err)
		}

		if !result.Success {
			utils.Error(fmt.Sprintf("网盘信息获取失败: %s", result.Message))
			return fmt.Errorf("网盘信息获取失败: %s", result.Message)
		}

		// 从结果中提取标题等信息
		if result.Data != nil {
			if data, ok := result.Data.(map[string]interface{}); ok {
				if title, ok := data["title"].(string); ok && title != "" {
					resource.Title = title
				}
				if description, ok := data["description"].(string); ok && description != "" {
					resource.Description = description
				}
			}
		}
	}

	// 处理分类
	if readyResource.Category != "" {
		categoryID, err := r.resolveCategory(readyResource.Category, nil)
		if err != nil {
			utils.Error(fmt.Sprintf("解析分类失败: %v", err))
		} else {
			resource.CategoryID = categoryID
		}
	}

	// 处理标签
	if readyResource.Tags != "" {
		tagIDs, err := r.handleTags(readyResource.Tags)
		if err != nil {
			utils.Error(fmt.Sprintf("处理标签失败: %v", err))
		} else {
			// 保存资源
			err = r.resourceRepo.Create(resource)
			if err != nil {
				return fmt.Errorf("创建资源失败: %v", err)
			}

			// 创建资源标签关联
			for _, tagID := range tagIDs {
				err = r.resourceRepo.CreateResourceTag(resource.ID, tagID)
				if err != nil {
					utils.Error(fmt.Sprintf("创建资源标签关联失败: %v", err))
				}
			}
		}
	} else {
		// 保存资源
		err := r.resourceRepo.Create(resource)
		if err != nil {
			return fmt.Errorf("创建资源失败: %v", err)
		}
	}

	return nil
}

// initPanCache 初始化平台缓存
func (r *ReadyResourceScheduler) initPanCache() {
	r.panCacheOnce.Do(func() {
		// 获取所有平台数据
		pans, err := r.panRepo.FindAll()
		if err != nil {
			utils.Error(fmt.Sprintf("初始化平台缓存失败: %v", err))
			return
		}

		// 建立 ServiceType 到 PanID 的映射
		serviceTypeToPanName := map[string]string{
			"quark":   "quark",
			"alipan":  "aliyun", // 阿里云盘在数据库中的名称是 aliyun
			"baidu":   "baidu",
			"uc":      "uc",
			"xunlei":  "xunlei",
			"tianyi":  "tianyi",
			"123pan":  "123pan",
			"115":     "115",
			"unknown": "other",
		}

		// 创建平台名称到ID的映射
		panNameToID := make(map[string]*uint)
		for _, pan := range pans {
			panID := pan.ID
			panNameToID[pan.Name] = &panID
		}

		// 建立 ServiceType 到 PanID 的映射
		for serviceType, panName := range serviceTypeToPanName {
			if panID, exists := panNameToID[panName]; exists {
				r.panCache[serviceType] = panID
				utils.Info(fmt.Sprintf("平台映射缓存: %s -> %s (ID: %d)", serviceType, panName, *panID))
			} else {
				utils.Error(fmt.Sprintf("警告: 未找到平台 %s 对应的数据库记录", panName))
			}
		}

		// 确保有默认的 other 平台
		if otherID, exists := panNameToID["other"]; exists {
			r.panCache["unknown"] = otherID
		}

		utils.Info(fmt.Sprintf("平台映射缓存初始化完成，共 %d 个映射", len(r.panCache)))
	})
}

// getPanIDByServiceType 根据服务类型获取平台ID
func (r *ReadyResourceScheduler) getPanIDByServiceType(serviceType panutils.ServiceType) *uint {
	r.initPanCache()

	serviceTypeStr := serviceType.String()
	if panID, exists := r.panCache[serviceTypeStr]; exists {
		return panID
	}

	// 如果找不到，返回 other 平台的ID
	if otherID, exists := r.panCache["other"]; exists {
		utils.Error(fmt.Sprintf("未找到服务类型 %s 的映射，使用默认平台 other", serviceTypeStr))
		return otherID
	}

	utils.Error(fmt.Sprintf("未找到服务类型 %s 的映射，且没有默认平台，返回nil", serviceTypeStr))
	return nil
}

// handleTags 处理标签
func (r *ReadyResourceScheduler) handleTags(tagStr string) ([]uint, error) {
	if tagStr == "" {
		return nil, nil
	}

	tagNames := splitTags(tagStr)
	var tagIDs []uint

	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// 查找或创建标签
		tag, err := r.tagRepo.FindByName(tagName)
		if err != nil {
			// 标签不存在，创建新标签
			tag = &entity.Tag{
				Name: tagName,
			}
			err = r.tagRepo.Create(tag)
			if err != nil {
				utils.Error(fmt.Sprintf("创建标签失败: %v", err))
				continue
			}
		}

		tagIDs = append(tagIDs, tag.ID)
	}

	return tagIDs, nil
}

// resolveCategory 解析分类
func (r *ReadyResourceScheduler) resolveCategory(categoryName string, tagIDs []uint) (*uint, error) {
	if categoryName == "" {
		return nil, nil
	}

	// 查找分类
	category, err := r.categoryRepo.FindByName(categoryName)
	if err != nil {
		// 分类不存在，创建新分类
		category = &entity.Category{
			Name: categoryName,
		}
		err = r.categoryRepo.Create(category)
		if err != nil {
			return nil, fmt.Errorf("创建分类失败: %v", err)
		}
	}

	return &category.ID, nil
}

// splitTags 分割标签字符串
func splitTags(tagStr string) []string {
	// 支持多种分隔符
	tagStr = strings.ReplaceAll(tagStr, "，", ",")
	tagStr = strings.ReplaceAll(tagStr, "；", ",")
	tagStr = strings.ReplaceAll(tagStr, ";", ",")
	tagStr = strings.ReplaceAll(tagStr, "、", ",")

	return strings.Split(tagStr, ",")
}

// derefString 解引用字符串指针
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
