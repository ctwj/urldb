package utils

import (
	"log"
	panutils "res_db/common"
	commonutils "res_db/common/utils"
	"res_db/db/entity"
	"res_db/db/repo"
	"strings"
	"sync"
	"time"
)

// Scheduler 定时任务管理器
type Scheduler struct {
	doubanService        *DoubanService
	hotDramaRepo         repo.HotDramaRepository
	readyResourceRepo    repo.ReadyResourceRepository
	resourceRepo         repo.ResourceRepository
	systemConfigRepo     repo.SystemConfigRepository
	panRepo              repo.PanRepository
	stopChan             chan bool
	isRunning            bool
	readyResourceRunning bool
	processingMutex      sync.Mutex // 防止ready_resource任务重叠执行
	hotDramaMutex        sync.Mutex // 防止热播剧任务重叠执行

	// 平台映射缓存
	panCache     map[string]*uint // serviceType -> panID
	panCacheOnce sync.Once
}

// NewScheduler 创建新的定时任务管理器
func NewScheduler(hotDramaRepo repo.HotDramaRepository, readyResourceRepo repo.ReadyResourceRepository, resourceRepo repo.ResourceRepository, systemConfigRepo repo.SystemConfigRepository, panRepo repo.PanRepository) *Scheduler {
	return &Scheduler{
		doubanService:        NewDoubanService(),
		hotDramaRepo:         hotDramaRepo,
		readyResourceRepo:    readyResourceRepo,
		resourceRepo:         resourceRepo,
		systemConfigRepo:     systemConfigRepo,
		panRepo:              panRepo,
		stopChan:             make(chan bool),
		isRunning:            false,
		readyResourceRunning: false,
		processingMutex:      sync.Mutex{},
		hotDramaMutex:        sync.Mutex{},
		panCache:             make(map[string]*uint),
	}
}

// StartHotDramaScheduler 启动热播剧定时任务
func (s *Scheduler) StartHotDramaScheduler() {
	if s.isRunning {
		log.Println("热播剧定时任务已在运行中")
		return
	}

	s.isRunning = true
	log.Println("启动热播剧定时任务")

	go func() {
		ticker := time.NewTicker(1 * time.Hour) // 每小时执行一次
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
					log.Println("上一次热播剧任务还在执行中，跳过本次执行")
				}
			case <-s.stopChan:
				log.Println("停止热播剧定时任务")
				return
			}
		}
	}()
}

// StopHotDramaScheduler 停止热播剧定时任务
func (s *Scheduler) StopHotDramaScheduler() {
	if !s.isRunning {
		log.Println("热播剧定时任务未在运行")
		return
	}

	s.stopChan <- true
	s.isRunning = false
	log.Println("已发送停止信号给热播剧定时任务")
}

// fetchHotDramaData 获取热播剧数据
func (s *Scheduler) fetchHotDramaData() {
	log.Println("开始获取热播剧数据...")

	dramaNames, err := s.doubanService.FetchHotDramaNames()
	if err != nil {
		log.Printf("获取热播剧数据失败: %v", err)
		return
	}

	log.Printf("成功获取到 %d 个热播剧: %v", len(dramaNames), dramaNames)

	// 处理获取到的热播剧数据
	s.processHotDramaNames(dramaNames)
}

// processHotDramaNames 处理热播剧名字
func (s *Scheduler) processHotDramaNames(dramaNames []string) {
	log.Printf("开始处理热播剧数据，共 %d 个", len(dramaNames))

	// 获取电影和电视剧的详细数据
	s.processMovieData()
	s.processTvData()

	log.Println("热播剧数据处理完成")
}

// processMovieData 处理电影数据
func (s *Scheduler) processMovieData() {
	log.Println("开始处理电影数据...")

	// 获取电影热门榜单（获取全部数据）
	movieResult, err := s.doubanService.GetMovieRanking("热门", "全部", 0, 0)
	if err != nil {
		log.Printf("获取电影榜单失败: %v", err)
		return
	}

	if movieResult.Success && movieResult.Data != nil {
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

			// 保存到数据库
			if err := s.hotDramaRepo.Upsert(drama); err != nil {
				log.Printf("保存电影数据失败: %v", err)
			} else {
				log.Printf("成功保存电影: %s (评分: %.1f, 年份: %s, 地区: %s)",
					item.Title, item.Rating.Value, item.Year, item.Region)
			}
		}
	}
}

// processTvData 处理电视剧数据
func (s *Scheduler) processTvData() {
	log.Println("开始处理电视剧数据...")

	// 获取电视剧热门榜单（获取全部数据）
	tvResult, err := s.doubanService.GetTvRanking("tv", "tv", 0, 0)
	if err != nil {
		log.Printf("获取电视剧榜单失败: %v", err)
		return
	}

	if tvResult.Success && tvResult.Data != nil {
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
				SubType:      "热门",
				Source:       "douban",
				DoubanID:     item.ID,
				DoubanURI:    item.URI,
			}

			// 保存到数据库
			if err := s.hotDramaRepo.Upsert(drama); err != nil {
				log.Printf("保存电视剧数据失败: %v", err)
			} else {
				log.Printf("成功保存电视剧: %s (评分: %.1f, 年份: %s, 地区: %s)",
					item.Title, item.Rating.Value, item.Year, item.Region)
			}
		}
	}
}

// IsRunning 检查定时任务是否在运行
func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}

// GetHotDramaNames 手动获取热播剧名字（用于测试或手动调用）
func (s *Scheduler) GetHotDramaNames() ([]string, error) {
	return s.doubanService.FetchHotDramaNames()
}

// StartReadyResourceScheduler 启动待处理资源自动处理任务
func (s *Scheduler) StartReadyResourceScheduler() {
	if s.readyResourceRunning {
		log.Println("待处理资源自动处理任务已在运行中")
		return
	}

	s.readyResourceRunning = true
	log.Println("启动待处理资源自动处理任务")

	go func() {
		// 获取系统配置中的间隔时间
		config, err := s.systemConfigRepo.GetOrCreateDefault()
		interval := 5 * time.Minute // 默认5分钟
		if err == nil && config.AutoProcessInterval > 0 {
			interval = time.Duration(config.AutoProcessInterval) * time.Minute
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Printf("待处理资源自动处理任务已启动，间隔时间: %v", interval)

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
					log.Println("上一次待处理资源任务还在执行中，跳过本次执行")
				}
			case <-s.stopChan:
				log.Println("停止待处理资源自动处理任务")
				return
			}
		}
	}()
}

// StopReadyResourceScheduler 停止待处理资源自动处理任务
func (s *Scheduler) StopReadyResourceScheduler() {
	if !s.readyResourceRunning {
		log.Println("待处理资源自动处理任务未在运行")
		return
	}

	s.stopChan <- true
	s.readyResourceRunning = false
	log.Println("已发送停止信号给待处理资源自动处理任务")
}

// processReadyResources 处理待处理资源
func (s *Scheduler) processReadyResources() {
	log.Println("开始处理待处理资源...")

	// 检查系统配置，确认是否启用自动处理
	config, err := s.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		log.Printf("获取系统配置失败: %v", err)
		return
	}

	if !config.AutoProcessReadyResources {
		log.Println("自动处理待处理资源功能已禁用")
		return
	}

	// 获取所有待处理资源
	readyResources, err := s.readyResourceRepo.FindAll()
	if err != nil {
		log.Printf("获取待处理资源失败: %v", err)
		return
	}

	if len(readyResources) == 0 {
		log.Println("没有待处理的资源")
		return
	}

	log.Printf("找到 %d 个待处理资源，开始处理...", len(readyResources))

	processedCount := 0
	factory := panutils.GetInstance() // 使用单例模式
	for _, readyResource := range readyResources {

		//readyResource.URL 是 查重
		exits, err := s.resourceRepo.FindExists(readyResource.URL)
		if err != nil {
			log.Printf("查重失败: %v", err)
			continue
		}
		if exits {
			log.Printf("资源已存在: %s", readyResource.URL)
			s.readyResourceRepo.Delete(readyResource.ID)
			continue
		}

		if err := s.convertReadyResourceToResource(readyResource, factory); err != nil {
			log.Printf("处理资源失败 (ID: %d): %v", readyResource.ID, err)
		}
		s.readyResourceRepo.Delete(readyResource.ID)
		processedCount++
		log.Printf("成功处理资源: %s", readyResource.URL)
	}

	log.Printf("待处理资源处理完成，共处理 %d 个资源", processedCount)
}

// convertReadyResourceToResource 将待处理资源转换为正式资源
func (s *Scheduler) convertReadyResourceToResource(readyResource entity.ReadyResource, factory *panutils.PanFactory) error {
	log.Printf("开始处理资源: %s", readyResource.URL)

	// 提取分享ID和服务类型
	shareID, serviceType := panutils.ExtractShareId(readyResource.URL)
	if serviceType == panutils.NotFound {
		log.Printf("不支持的链接地址: %s", readyResource.URL)
		return nil
	}

	log.Printf("检测到服务类型: %s, 分享ID: %s", serviceType.String(), shareID)

	// 不是夸克，直接保存，
	if serviceType != panutils.Quark {
		// 检测是否有效
		checkResult, _ := commonutils.CheckURL(readyResource.URL)
		if !checkResult.Status {
			log.Printf("链接无效: %s", readyResource.URL)
			return nil
		}

		// 入库
	}

	// 准备配置
	config := &panutils.PanConfig{
		URL:         readyResource.URL,
		Code:        "", // 可以从readyResource中获取
		IsType:      1,  // 转存并分享后的资源信息  0 转存后分享， 1 只获取基本信息
		ExpiredType: 1,  // 永久分享
		AdFid:       "",
		Stoken:      "",
	}

	// 通过工厂获取对应的网盘服务单例
	panService, err := factory.CreatePanService(readyResource.URL, config)
	if err != nil {
		log.Printf("获取网盘服务失败: %v", err)
		return err
	}

	// 阿里云盘特殊处理：检查URL有效性
	// if serviceType == panutils.Alipan {
	// 	checkResult, _ := CheckURL(readyResource.URL)
	// 	if !checkResult.Status {
	// 		log.Printf("阿里云盘链接无效: %s", readyResource.URL)
	// 		return nil
	// 	}

	// 	// 如果有标题，直接创建资源
	// 	if readyResource.Title != nil && *readyResource.Title != "" {
	// 		resource := &entity.Resource{
	// 			Title:       *readyResource.Title,
	// 			Description: readyResource.Description,
	// 			URL:         readyResource.URL,
	// 			PanID:       s.determinePanID(readyResource.URL),
	// 			IsValid:     true,
	// 			IsPublic:    true,
	// 		}

	// 		// 如果有分类信息，尝试查找或创建分类
	// 		if readyResource.Category != "" {
	// 			categoryID, err := s.getOrCreateCategory(readyResource.Category)
	// 			if err == nil {
	// 				resource.CategoryID = &categoryID
	// 			}
	// 		}

	// 		return s.resourceRepo.Create(resource)
	// 	}
	// }

	// 统一处理：尝试转存获取标题
	result, err := panService.Transfer(shareID)
	if err != nil {
		log.Printf("网盘信息获取失败: %v", err)
		return err
	}

	if !result.Success {
		log.Printf("网盘信息获取失败: %s", result.Message)
		return nil
	}

	// 提取转存结果
	if resultData, ok := result.Data.(map[string]interface{}); ok {
		title := resultData["title"].(string)
		shareURL := resultData["shareUrl"].(string)
		// fid := resultData["fid"].(string) // 暂时未使用

		// 创建资源记录
		resource := &entity.Resource{
			Title:       title,
			Description: readyResource.Description,
			URL:         shareURL,
			PanID:       s.getPanIDByServiceType(serviceType),
			IsValid:     true,
			IsPublic:    true,
		}

		// 如果有分类信息，尝试查找或创建分类
		if readyResource.Category != "" {
			categoryID, err := s.getOrCreateCategory(readyResource.Category)
			if err == nil {
				resource.CategoryID = &categoryID
			}
		}

		return s.resourceRepo.Create(resource)
	}

	log.Printf("转存结果格式异常")
	return nil
}

// getOrCreateCategory 获取或创建分类
func (s *Scheduler) getOrCreateCategory(categoryName string) (uint, error) {
	// 这里需要实现分类的查找和创建逻辑
	// 由于没有CategoryRepository的注入，这里先返回0
	// 你可以根据需要添加CategoryRepository的依赖
	return 0, nil
}

// initPanCache 初始化平台映射缓存
func (s *Scheduler) initPanCache() {
	s.panCacheOnce.Do(func() {
		// 获取所有平台数据
		pans, err := s.panRepo.FindAll()
		if err != nil {
			log.Printf("初始化平台缓存失败: %v", err)
			return
		}

		// 建立 ServiceType 到 PanID 的映射
		serviceTypeToPanName := map[string]string{
			"quark":   "quark",
			"alipan":  "aliyun", // 阿里云盘在数据库中的名称是 aliyun
			"baidu":   "baidu",
			"uc":      "uc",
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
				log.Printf("平台映射缓存: %s -> %s (ID: %d)", serviceType, panName, *panID)
			} else {
				log.Printf("警告: 未找到平台 %s 对应的数据库记录", panName)
			}
		}

		// 确保有默认的 other 平台
		if otherID, exists := panNameToID["other"]; exists {
			s.panCache["unknown"] = otherID
		}

		log.Printf("平台映射缓存初始化完成，共 %d 个映射", len(s.panCache))
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
		log.Printf("未找到服务类型 %s 的映射，使用默认平台 other", serviceTypeStr)
		return otherID
	}

	log.Printf("未找到服务类型 %s 的映射，且没有默认平台，返回nil", serviceTypeStr)
	return nil
}

// IsReadyResourceRunning 检查待处理资源自动处理任务是否在运行
func (s *Scheduler) IsReadyResourceRunning() bool {
	return s.readyResourceRunning
}
