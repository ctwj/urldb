package utils

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
	"gorm.io/gorm"
)

// Scheduler 定时任务管理器
type Scheduler struct {
	doubanService     *DoubanService
	hotDramaRepo      repo.HotDramaRepository
	readyResourceRepo repo.ReadyResourceRepository
	resourceRepo      repo.ResourceRepository
	systemConfigRepo  repo.SystemConfigRepository
	panRepo           repo.PanRepository
	cksRepo           repo.CksRepository
	// 新增
	tagRepo              repo.TagRepository
	categoryRepo         repo.CategoryRepository
	stopChan             chan bool
	isRunning            bool
	readyResourceRunning bool
	autoTransferRunning  bool
	processingMutex      sync.Mutex // 防止ready_resource任务重叠执行
	hotDramaMutex        sync.Mutex // 防止热播剧任务重叠执行
	autoTransferMutex    sync.Mutex // 防止自动转存任务重叠执行

	// 平台映射缓存
	panCache     map[string]*uint // serviceType -> panID
	panCacheOnce sync.Once
}

// NewScheduler 创建新的定时任务管理器
func NewScheduler(hotDramaRepo repo.HotDramaRepository, readyResourceRepo repo.ReadyResourceRepository, resourceRepo repo.ResourceRepository, systemConfigRepo repo.SystemConfigRepository, panRepo repo.PanRepository, cksRepo repo.CksRepository, tagRepo repo.TagRepository, categoryRepo repo.CategoryRepository) *Scheduler {
	return &Scheduler{
		doubanService:        NewDoubanService(),
		hotDramaRepo:         hotDramaRepo,
		readyResourceRepo:    readyResourceRepo,
		resourceRepo:         resourceRepo,
		systemConfigRepo:     systemConfigRepo,
		panRepo:              panRepo,
		cksRepo:              cksRepo,
		tagRepo:              tagRepo,
		categoryRepo:         categoryRepo,
		stopChan:             make(chan bool),
		isRunning:            false,
		readyResourceRunning: false,
		autoTransferRunning:  false,
		processingMutex:      sync.Mutex{},
		hotDramaMutex:        sync.Mutex{},
		autoTransferMutex:    sync.Mutex{},
		panCache:             make(map[string]*uint),
	}
}

// StartHotDramaScheduler 启动热播剧定时任务
func (s *Scheduler) StartHotDramaScheduler() {
	if s.isRunning {
		Info("热播剧定时任务已在运行中")
		return
	}

	s.isRunning = true
	Info("启动热播剧定时任务")

	go func() {
		ticker := time.NewTicker(12 * time.Hour) // 每12小时执行一次
		defer ticker.Stop()

		// 立即执行一次
		s.fetchHotDramaData()

		for {
			select {
			case <-ticker.C:
				// 使用TryLock防止任务重叠执行
				if s.hotDramaMutex.TryLock() {
					go func() {
						defer s.hotDramaMutex.Unlock()
						s.fetchHotDramaData()
					}()
				} else {
					Info("上一次热播剧任务还在执行中，跳过本次执行")
				}
			case <-s.stopChan:
				Info("停止热播剧定时任务")
				return
			}
		}
	}()
}

// StopHotDramaScheduler 停止热播剧定时任务
func (s *Scheduler) StopHotDramaScheduler() {
	if !s.isRunning {
		Info("热播剧定时任务未在运行")
		return
	}

	s.stopChan <- true
	s.isRunning = false
	Info("已发送停止信号给热播剧定时任务")
}

// fetchHotDramaData 获取热播剧数据
func (s *Scheduler) fetchHotDramaData() {
	Info("开始获取热播剧数据...")

	// 直接处理电影和电视剧数据，不再需要FetchHotDramaNames
	s.processHotDramaNames([]string{})
}

// processHotDramaNames 处理热播剧名字
func (s *Scheduler) processHotDramaNames(dramaNames []string) {
	Info("开始处理热播剧数据，共 %d 个", len(dramaNames))

	// 收集所有数据
	var allDramas []*entity.HotDrama

	// 获取电影数据
	movieDramas := s.processMovieData()
	allDramas = append(allDramas, movieDramas...)

	// 获取电视剧数据
	tvDramas := s.processTvData()
	allDramas = append(allDramas, tvDramas...)

	// 清空数据库
	Info("准备清空数据库，当前共有 %d 条数据", len(allDramas))
	if err := s.hotDramaRepo.DeleteAll(); err != nil {
		Error("清空数据库失败: %v", err)
		return
	}
	Info("数据库清空完成")

	// 批量插入所有数据
	if len(allDramas) > 0 {
		Info("开始批量插入 %d 条数据", len(allDramas))
		if err := s.hotDramaRepo.BatchCreate(allDramas); err != nil {
			Error("批量插入数据失败: %v", err)
		} else {
			Info("成功批量插入 %d 条数据", len(allDramas))
		}
	} else {
		Info("没有数据需要插入")
	}

	Info("热播剧数据处理完成")
}

// processMovieData 处理电影数据
func (s *Scheduler) processMovieData() []*entity.HotDrama {
	Info("开始处理电影数据...")

	var movieDramas []*entity.HotDrama

	// 使用GetTypePage方法获取电影数据
	movieResult, err := s.doubanService.GetTypePage("热门", "全部")
	if err != nil {
		Error("获取电影榜单失败: %v", err)
		return movieDramas
	}

	if movieResult.Success && movieResult.Data != nil {
		Info("电影获取到 %d 个数据", len(movieResult.Data.Items))

		for _, item := range movieResult.Data.Items {
			drama := &entity.HotDrama{
				Title:        item.Title,
				CardSubtitle: item.CardSubtitle,
				EpisodesInfo: item.EpisodesInfo,
				IsNew:        item.IsNew,
				Rating:       item.Rating.Value,
				RatingCount:  item.Rating.Count,
				Year:         item.Year,
				Region:       item.Region,
				Genres:       strings.Join(item.Genres, ", "),
				Directors:    strings.Join(item.Directors, ", "),
				Actors:       strings.Join(item.Actors, ", "),
				PosterURL:    item.Pic.Normal,
				Category:     "电影",
				SubType:      "热门",
				Source:       "douban",
				DoubanID:     item.ID,
				DoubanURI:    item.URI,
			}

			movieDramas = append(movieDramas, drama)
			Info("收集电影: %s (评分: %.1f, 年份: %s, 地区: %s)",
				item.Title, item.Rating.Value, item.Year, item.Region)
		}
	} else {
		Warn("电影获取数据失败或为空")
	}

	Info("电影数据处理完成，共收集 %d 条数据", len(movieDramas))
	return movieDramas
}

// processTvData 处理电视剧数据
func (s *Scheduler) processTvData() []*entity.HotDrama {
	Info("开始处理电视剧数据...")

	var tvDramas []*entity.HotDrama

	// 获取所有tv类型
	tvTypes := s.doubanService.GetAllTvTypes()
	Info("获取到 %d 个tv类型: %v", len(tvTypes), tvTypes)

	// 遍历每个type，分别请求数据
	for _, tvType := range tvTypes {
		Info("正在处理tv类型: %s", tvType)

		// 使用GetTypePage方法请求数据
		tvResult, err := s.doubanService.GetTypePage("tv", tvType)
		if err != nil {
			Error("获取tv类型 %s 数据失败: %v", tvType, err)
			continue
		}

		if tvResult.Success && tvResult.Data != nil {
			Info("tv类型 %s 获取到 %d 个数据", tvType, len(tvResult.Data.Items))

			for _, item := range tvResult.Data.Items {
				drama := &entity.HotDrama{
					Title:        item.Title,
					CardSubtitle: item.CardSubtitle,
					EpisodesInfo: item.EpisodesInfo,
					IsNew:        item.IsNew,
					Rating:       item.Rating.Value,
					RatingCount:  item.Rating.Count,
					Year:         item.Year,
					Region:       item.Region,
					Genres:       strings.Join(item.Genres, ", "),
					Directors:    strings.Join(item.Directors, ", "),
					Actors:       strings.Join(item.Actors, ", "),
					PosterURL:    item.Pic.Normal,
					Category:     "电视剧",
					SubType:      tvType, // 使用具体的tv类型
					Source:       "douban",
					DoubanID:     item.ID,
					DoubanURI:    item.URI,
				}

				tvDramas = append(tvDramas, drama)
				Info("收集tv类型 %s: %s (评分: %.1f, 年份: %s, 地区: %s)",
					tvType, item.Title, item.Rating.Value, item.Year, item.Region)
			}
		} else {
			Warn("tv类型 %s 获取数据失败或为空", tvType)
		}

		// 每个type请求间隔1秒，避免请求过于频繁
		time.Sleep(1 * time.Second)
	}

	Info("电视剧数据处理完成，共收集 %d 条数据", len(tvDramas))
	return tvDramas
}

// IsRunning 检查定时任务是否在运行
func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}

// GetHotDramaNames 手动获取热播剧名字（用于测试或手动调用）
func (s *Scheduler) GetHotDramaNames() ([]string, error) {
	// 由于删除了FetchHotDramaNames方法，返回空数组
	return []string{}, nil
}

// StartReadyResourceScheduler 启动待处理资源自动处理任务
func (s *Scheduler) StartReadyResourceScheduler() {
	if s.readyResourceRunning {
		Info("待处理资源自动处理任务已在运行中")
		return
	}

	s.readyResourceRunning = true
	Info("启动待处理资源自动处理任务")

	go func() {
		// 获取系统配置中的间隔时间
		interval := 3 * time.Minute // 默认3分钟
		if autoProcessInterval, err := s.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoProcessInterval); err == nil && autoProcessInterval > 0 {
			interval = time.Duration(autoProcessInterval) * time.Minute
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		Info("待处理资源自动处理任务已启动，间隔时间: %v", interval)

		// 立即执行一次
		s.processReadyResources()

		for {
			select {
			case <-ticker.C:
				// 使用TryLock防止任务重叠执行
				if s.processingMutex.TryLock() {
					go func() {
						defer s.processingMutex.Unlock()
						s.processReadyResources()
					}()
				} else {
					Info("上一次待处理资源任务还在执行中，跳过本次执行")
				}
			case <-s.stopChan:
				Info("停止待处理资源自动处理任务")
				return
			}
		}
	}()
}

// StopReadyResourceScheduler 停止待处理资源自动处理任务
func (s *Scheduler) StopReadyResourceScheduler() {
	if !s.readyResourceRunning {
		Info("待处理资源自动处理任务未在运行")
		return
	}

	s.stopChan <- true
	s.readyResourceRunning = false
	Info("已发送停止信号给待处理资源自动处理任务")
}

// processReadyResources 处理待处理资源
func (s *Scheduler) processReadyResources() {
	Info("开始处理待处理资源...")

	// 检查系统配置，确认是否启用自动处理
	autoProcess, err := s.systemConfigRepo.GetConfigBool(entity.ConfigKeyAutoProcessReadyResources)
	if err != nil {
		Error("获取系统配置失败: %v", err)
		return
	}

	if !autoProcess {
		Info("自动处理待处理资源功能已禁用")
		return
	}

	// 获取所有没有错误的待处理资源
	readyResources, err := s.readyResourceRepo.FindAll()
	// readyResources, err := s.readyResourceRepo.FindWithoutErrors()
	if err != nil {
		Error("获取待处理资源失败: %v", err)
		return
	}

	if len(readyResources) == 0 {
		Info("没有待处理的资源")
		return
	}

	Info("找到 %d 个待处理资源，开始处理...", len(readyResources))

	processedCount := 0
	factory := panutils.GetInstance() // 使用单例模式
	for _, readyResource := range readyResources {

		//readyResource.URL 是 查重
		exits, err := s.resourceRepo.FindExists(readyResource.URL)
		if err != nil {
			Error("查重失败: %v", err)
			continue
		}
		if exits {
			Info("资源已存在: %s", readyResource.URL)
			s.readyResourceRepo.Delete(readyResource.ID)
			continue
		}

		if err := s.convertReadyResourceToResource(readyResource, factory); err != nil {
			Error("处理资源失败 (ID: %d): %v", readyResource.ID, err)

			// 保存完整的错误信息
			readyResource.ErrorMsg = err.Error()

			if updateErr := s.readyResourceRepo.Update(&readyResource); updateErr != nil {
				Error("更新错误信息失败 (ID: %d): %v", readyResource.ID, updateErr)
			} else {
				Info("已保存错误信息到资源 (ID: %d): %s", readyResource.ID, err.Error())
			}

			// 处理失败后删除资源，避免重复处理
			s.readyResourceRepo.Delete(readyResource.ID)
		} else {
			// 处理成功，删除readyResource
			s.readyResourceRepo.Delete(readyResource.ID)
			processedCount++
			Info("成功处理资源: %s", readyResource.URL)
		}
	}

	Info("待处理资源处理完成，共处理 %d 个资源", processedCount)
}

// convertReadyResourceToResource 将待处理资源转换为正式资源
func (s *Scheduler) convertReadyResourceToResource(readyResource entity.ReadyResource, factory *panutils.PanFactory) error {
	Debug("开始处理资源: %s", readyResource.URL)

	// 提取分享ID和服务类型
	shareID, serviceType := panutils.ExtractShareId(readyResource.URL)
	if serviceType == panutils.NotFound {
		Warn("不支持的链接地址: %s", readyResource.URL)
		return NewUnsupportedLinkError(readyResource.URL)
	}

	Debug("检测到服务类型: %s, 分享ID: %s", serviceType.String(), shareID)

	resource := &entity.Resource{
		Title:       derefString(readyResource.Title),
		Description: readyResource.Description,
		URL:         readyResource.URL,
		Cover:       readyResource.Img,
		IsValid:     true,
		IsPublic:    true,
		Key:         readyResource.Key,
		PanID:       s.getPanIDByServiceType(serviceType),
	}

	// 不是夸克，直接保存，
	if serviceType != panutils.Quark {
		// 检测是否有效
		checkResult, err := commonutils.CheckURL(readyResource.URL)
		if err != nil {
			Error("链接检查失败: %v", err)
			return NewLinkCheckError(readyResource.URL, err.Error())
		}
		if !checkResult.Status {
			Warn("链接无效: %s", readyResource.URL)
			return NewInvalidLinkError(readyResource.URL, "链接状态检查失败")
		}
	} else {
		// 获取夸克网盘账号的 cookie
		panID := s.getPanIDByServiceType(serviceType)
		if panID == nil {
			Error("未找到对应的平台ID")
			return NewPlatformNotFoundError(serviceType.String())
		}

		accounts, err := s.cksRepo.FindByPanID(*panID)
		if err != nil {
			Error("获取夸克网盘账号失败: %v", err)
			return NewServiceCreationError(readyResource.URL, fmt.Sprintf("获取网盘账号失败: %v", err))
		}

		if len(accounts) == 0 {
			Error("没有可用的夸克网盘账号")
			return NewNoAccountError(serviceType.String())
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
			Error("没有有效的夸克网盘账号")
			return NewNoValidAccountError(serviceType.String())
		}

		Debug("使用夸克网盘账号: %d, Cookie: %s", selectedAccount.ID, selectedAccount.Ck[:20]+"...")

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
			Error("获取网盘服务失败: %v", err)
			return NewServiceCreationError(readyResource.URL, err.Error())
		}

		// 统一处理：尝试转存获取标题
		result, err := panService.Transfer(shareID)
		if err != nil {
			Error("网盘信息获取失败: %v", err)
			return NewTransferFailedError(readyResource.URL, err.Error())
		}

		if !result.Success {
			Error("网盘信息获取失败: %s", result.Message)
			return NewTransferFailedError(readyResource.URL, result.Message)
		}

		// 如果获取到了标题，更新资源标题
		// if result.Title != "" {
		// 	resource.Title = result.Title
		// }
	}

	// 处理标签
	tagIDs, err := s.handleTags(readyResource.Tags)
	if err != nil {
		Error("处理标签失败: %v", err)
		return NewTagProcessingError(err.Error())
	}
	// 如果没有标签，tagIDs 可能为 nil，这是正常的
	if tagIDs == nil {
		tagIDs = []uint{} // 初始化为空数组
	}
	// 处理分类
	categoryID, err := s.resolveCategory(readyResource.Category, tagIDs)
	if err != nil {
		Error("处理分类失败: %v", err)
		return NewCategoryProcessingError(err.Error())
	}
	if categoryID != nil {
		resource.CategoryID = categoryID
	}
	// 保存资源
	err = s.resourceRepo.Create(resource)
	if err != nil {
		Error("资源保存失败: %v", err)
		return NewResourceSaveError(readyResource.URL, err.Error())
	}
	// 插入 resource_tags 关联
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			err := s.resourceRepo.CreateResourceTag(resource.ID, tagID)
			if err != nil {
				Error("插入资源标签关联失败: %v", err)
				// 这里不返回错误，因为资源已经保存成功，标签关联失败不影响主要功能
			}
		}
	} else {
		Debug("没有标签，跳过插入资源标签关联")
	}
	return nil
}

// // getOrCreateCategory 获取或创建分类
// func (s *Scheduler) getOrCreateCategory(categoryName string) (uint, error) {
// 	// 这里需要实现分类的查找和创建逻辑
// 	// 由于没有CategoryRepository的注入，这里先返回0
// 	// 你可以根据需要添加CategoryRepository的依赖
// 	return 0, nil
// }

// initPanCache 初始化平台映射缓存
func (s *Scheduler) initPanCache() {
	s.panCacheOnce.Do(func() {
		// 获取所有平台数据
		pans, err := s.panRepo.FindAll()
		if err != nil {
			Error("初始化平台缓存失败: %v", err)
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
				s.panCache[serviceType] = panID
				Debug("平台映射缓存: %s -> %s (ID: %d)", serviceType, panName, *panID)
			} else {
				Warn("警告: 未找到平台 %s 对应的数据库记录", panName)
			}
		}

		// 确保有默认的 other 平台
		if otherID, exists := panNameToID["other"]; exists {
			s.panCache["unknown"] = otherID
		}

		Info("平台映射缓存初始化完成，共 %d 个映射", len(s.panCache))
	})
}

// getPanIDByServiceType 根据服务类型获取平台ID
func (s *Scheduler) getPanIDByServiceType(serviceType panutils.ServiceType) *uint {
	s.initPanCache()

	serviceTypeStr := serviceType.String()
	if panID, exists := s.panCache[serviceTypeStr]; exists {
		return panID
	}

	// 如果找不到，返回 other 平台的ID
	if otherID, exists := s.panCache["other"]; exists {
		Warn("未找到服务类型 %s 的映射，使用默认平台 other", serviceTypeStr)
		return otherID
	}

	Warn("未找到服务类型 %s 的映射，且没有默认平台，返回nil", serviceTypeStr)
	return nil
}

// IsReadyResourceRunning 检查待处理资源自动处理任务是否在运行
func (s *Scheduler) IsReadyResourceRunning() bool {
	return s.readyResourceRunning
}

// StartAutoTransferScheduler 启动自动转存定时任务
func (s *Scheduler) StartAutoTransferScheduler() {
	if s.autoTransferRunning {
		Info("自动转存定时任务已在运行中")
		return
	}

	s.autoTransferRunning = true
	Info("启动自动转存定时任务")

	go func() {
		// 获取系统配置中的间隔时间
		interval := 5 * time.Minute // 默认5分钟
		if autoProcessInterval, err := s.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoProcessInterval); err == nil && autoProcessInterval > 0 {
			interval = time.Duration(autoProcessInterval) * time.Minute
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		Info("自动转存定时任务已启动，间隔时间: %v", interval)

		// 立即执行一次
		s.processAutoTransfer()

		for {
			select {
			case <-ticker.C:
				// 使用TryLock防止任务重叠执行
				if s.autoTransferMutex.TryLock() {
					go func() {
						defer s.autoTransferMutex.Unlock()
						s.processAutoTransfer()
					}()
				} else {
					Info("上一次自动转存任务还在执行中，跳过本次执行")
				}
			case <-s.stopChan:
				Info("停止自动转存定时任务")
				return
			}
		}
	}()
}

// StopAutoTransferScheduler 停止自动转存定时任务
func (s *Scheduler) StopAutoTransferScheduler() {
	if !s.autoTransferRunning {
		Info("自动转存定时任务未在运行")
		return
	}

	s.stopChan <- true
	s.autoTransferRunning = false
	Info("已发送停止信号给自动转存定时任务")
}

// IsAutoTransferRunning 检查自动转存定时任务是否在运行
func (s *Scheduler) IsAutoTransferRunning() bool {
	return s.autoTransferRunning
}

// processAutoTransfer 处理自动转存
func (s *Scheduler) processAutoTransfer() {
	Info("开始处理自动转存...")

	// 检查系统配置，确认是否启用自动转存
	autoTransferEnabled, err := s.systemConfigRepo.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
	if err != nil {
		Error("获取系统配置失败: %v", err)
		return
	}

	if !autoTransferEnabled {
		Info("自动转存功能已禁用")
		return
	}

	// 获取quark平台ID
	panRepoImpl, ok := s.panRepo.(interface{ GetDB() *gorm.DB })
	if !ok {
		Error("panRepo不支持GetDB方法")
		return
	}
	var quarkPan entity.Pan
	err = panRepoImpl.GetDB().Where("name = ?", "quark").First(&quarkPan).Error
	if err != nil {
		Error("未找到quark平台: %v", err)
		return
	}
	quarkPanID := quarkPan.ID

	// 获取所有账号
	accounts, err := s.cksRepo.FindAll()
	if err != nil {
		Error("获取网盘账号失败: %v", err)
		return
	}

	// 获取最小存储空间配置
	autoTransferMinSpace, err := s.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		Error("获取最小存储空间配置失败: %v", err)
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
		Info("没有可用的quark网盘账号")
		return
	}

	Info("找到 %d 个可用quark网盘账号，开始自动转存处理...", len(validAccounts))

	// 获取需要转存的资源
	resources, err := s.getResourcesForTransfer(quarkPanID)
	if err != nil {
		Error("获取需要转存的资源失败: %v", err)
		return
	}

	if len(resources) == 0 {
		Info("没有需要转存的资源")
		return
	}

	Info("找到 %d 个需要转存的资源", len(resources))

	// 获取违禁词配置
	forbiddenWords, err := s.systemConfigRepo.GetConfigValue(entity.ConfigKeyForbiddenWords)
	if err != nil {
		Error("获取违禁词配置失败: %v", err)
		forbiddenWords = "" // 如果获取失败，使用空字符串
	}

	// 过滤包含违禁词的资源
	var filteredResources []*entity.Resource
	if forbiddenWords != "" {
		words := strings.Split(forbiddenWords, ",")
		for _, resource := range resources {
			shouldSkip := false
			title := strings.ToLower(resource.Title)
			description := strings.ToLower(resource.Description)

			for _, word := range words {
				word = strings.TrimSpace(word)
				if word != "" && (strings.Contains(title, strings.ToLower(word)) || strings.Contains(description, strings.ToLower(word))) {
					Info("跳过包含违禁词 '%s' 的资源: %s", word, resource.Title)
					shouldSkip = true
					break
				}
			}

			if !shouldSkip {
				filteredResources = append(filteredResources, resource)
			}
		}
		Info("违禁词过滤后，剩余 %d 个资源需要转存", len(filteredResources))
	} else {
		filteredResources = resources
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
				if err := s.transferResource(res, []entity.Cks{acc}, factory); err != nil {
					Error("转存资源失败 (ID: %d): %v", res.ID, err)
				} else {
					Info("成功转存资源: %s", res.Title)
					rand.Seed(GetCurrentTime().UnixNano())
					sleepSec := rand.Intn(3) + 1 // 1,2,3
					time.Sleep(time.Duration(sleepSec) * time.Second)
				}
			}
		}(account)
	}
	wg.Wait()
	Info("自动转存处理完成，账号数: %d，资源数: %d", len(validAccounts), len(filteredResources))
}

// getResourcesForTransfer 获取需要转存的资源
func (s *Scheduler) getResourcesForTransfer(quarkPanID uint) ([]*entity.Resource, error) {
	// 获取自动转存限制天数配置
	autoTransferLimitDays, err := s.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferLimitDays)
	if err != nil {
		Error("获取自动转存限制天数配置失败: %v", err)
		return nil, err
	}

	days := autoTransferLimitDays
	var sinceTime time.Time
	if days > 0 {
		sinceTime = GetCurrentTime().AddDate(0, 0, -days)
	} else {
		sinceTime = time.Time{}
	}

	repoImpl, ok := s.resourceRepo.(*repo.ResourceRepositoryImpl)
	if !ok {
		return nil, fmt.Errorf("resourceRepo不是ResourceRepositoryImpl类型")
	}
	return repoImpl.GetResourcesForTransfer(quarkPanID, sinceTime)
}

var resourceUpdateMutex sync.Mutex // 全局互斥锁，保证多协程安全

// transferResource 转存单个资源
func (s *Scheduler) transferResource(resource *entity.Resource, accounts []entity.Cks, factory *panutils.PanFactory) error {
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
	autoTransferMinSpace, err := s.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		Error("获取最小存储空间配置失败: %v", err)
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
		resourceUpdateMutex.Lock()
		defer resourceUpdateMutex.Unlock()
		s.resourceRepo.Update(&entity.Resource{
			ID:       resource.ID,
			ErrorMsg: err.Error(),
		})
		return fmt.Errorf("转存失败: %v", err)
	}

	if result == nil || !result.Success {
		errMsg := "转存失败"
		if result != nil && result.Message != "" {
			errMsg = result.Message
		}
		resourceUpdateMutex.Lock()
		defer resourceUpdateMutex.Unlock()
		s.resourceRepo.Update(&entity.Resource{
			ID:       resource.ID,
			ErrorMsg: errMsg,
		})
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

	resourceUpdateMutex.Lock()
	defer resourceUpdateMutex.Unlock()
	err = s.resourceRepo.Update(&entity.Resource{
		ID:       resource.ID,
		SaveURL:  saveURL,
		CkID:     &account.ID,
		Fid:      fid,
		ErrorMsg: "",
	})
	if err != nil {
		return fmt.Errorf("保存转存结果失败: %v", err)
	}

	return nil
}

// selectBestAccount 选择最佳账号
func (s *Scheduler) selectBestAccount(accounts []entity.Cks) *entity.Cks {
	if len(accounts) == 0 {
		return nil
	}

	// 获取最小存储空间配置
	autoTransferMinSpace, err := s.systemConfigRepo.GetConfigInt(entity.ConfigKeyAutoTransferMinSpace)
	if err != nil {
		Error("获取最小存储空间配置失败: %v", err)
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

		score := s.calculateAccountScore(account)
		if score > bestScore {
			bestScore = score
			bestAccount = account
		}
	}

	return bestAccount
}

// calculateAccountScore 计算账号评分
func (s *Scheduler) calculateAccountScore(account *entity.Cks) int64 {
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

// 分割标签，支持中英文逗号
func splitTags(tagStr string) []string {
	tagStr = strings.ReplaceAll(tagStr, "，", ",")
	return strings.Split(tagStr, ",")
}

// 处理标签，返回所有标签ID
func (s *Scheduler) handleTags(tagStr string) ([]uint, error) {
	if tagStr == "" {
		Debug("标签字符串为空，返回空数组")
		return []uint{}, nil // 返回空数组而不是 nil
	}

	Debug("开始处理标签字符串: %s", tagStr)
	tagNames := splitTags(tagStr)
	Debug("分割后的标签名称: %v", tagNames)

	var tagIDs []uint
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			Debug("跳过空标签名称")
			continue
		}

		Debug("查找标签: %s", name)
		tag, err := s.tagRepo.FindByName(name)
		if err != nil {
			Debug("标签 %s 不存在，创建新标签", name)
			// 不存在则新建
			tag = &entity.Tag{Name: name}
			err = s.tagRepo.Create(tag)
			if err != nil {
				Error("创建标签 %s 失败: %v", name, err)
				return nil, fmt.Errorf("创建标签 %s 失败: %v", name, err)
			}
			Debug("成功创建标签: %s (ID: %d)", name, tag.ID)
		} else {
			Debug("找到已存在的标签: %s (ID: %d)", name, tag.ID)
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	Debug("处理完成，标签ID列表: %v", tagIDs)
	return tagIDs, nil
}

// 分类处理逻辑
func (s *Scheduler) resolveCategory(categoryName string, tagIDs []uint) (*uint, error) {
	Debug("开始处理分类，分类名称: %s, 标签ID列表: %v", categoryName, tagIDs)

	if categoryName != "" {
		Debug("查找分类: %s", categoryName)
		cat, err := s.categoryRepo.FindByName(categoryName)
		if err != nil {
			Debug("分类 %s 不存在: %v", categoryName, err)
		} else if cat != nil {
			Debug("找到分类: %s (ID: %d)", categoryName, cat.ID)
			return &cat.ID, nil
		}
	}

	// 没有分类，尝试用标签反查
	if len(tagIDs) == 0 {
		Debug("没有标签，无法通过标签反查分类")
		return nil, nil
	}

	Debug("尝试通过标签反查分类")
	for _, tagID := range tagIDs {
		Debug("查找标签ID: %d", tagID)
		tag, err := s.tagRepo.GetByID(tagID)
		if err != nil {
			Debug("查找标签ID %d 失败: %v", tagID, err)
			continue
		}
		if tag != nil && tag.CategoryID != nil {
			Debug("通过标签 %s (ID: %d) 找到分类ID: %d", tag.Name, tagID, *tag.CategoryID)
			return tag.CategoryID, nil
		}
	}

	Debug("未找到分类，返回nil")
	return nil, nil
}

// 工具函数，解引用string指针
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
