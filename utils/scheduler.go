package utils

import (
	"log"
	"res_db/db/entity"
	"res_db/db/repo"
	panutils "res_db/utils/pan"
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
	stopChan             chan bool
	isRunning            bool
	readyResourceRunning bool
	processingMutex      sync.Mutex // 防止ready_resource任务重叠执行
	hotDramaMutex        sync.Mutex // 防止热播剧任务重叠执行
}

// NewScheduler 创建新的定时任务管理器
func NewScheduler(hotDramaRepo repo.HotDramaRepository, readyResourceRepo repo.ReadyResourceRepository, resourceRepo repo.ResourceRepository, systemConfigRepo repo.SystemConfigRepository) *Scheduler {
	return &Scheduler{
		doubanService:        NewDoubanService(),
		hotDramaRepo:         hotDramaRepo,
		readyResourceRepo:    readyResourceRepo,
		resourceRepo:         resourceRepo,
		systemConfigRepo:     systemConfigRepo,
		stopChan:             make(chan bool),
		isRunning:            false,
		readyResourceRunning: false,
		processingMutex:      sync.Mutex{},
		hotDramaMutex:        sync.Mutex{},
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

	// 获取电影热门榜单
	movieResult, err := s.doubanService.GetMovieRanking("热门", "全部", 0, 20)
	if err != nil {
		log.Printf("获取电影榜单失败: %v", err)
		return
	}

	if movieResult.Success && movieResult.Data != nil {
		for _, item := range movieResult.Data.Items {
			drama := &entity.HotDrama{
				Title:     item.Title,
				Rating:    item.Rating.Value,
				Year:      item.Year,
				Directors: strings.Join(item.Directors, ", "),
				Actors:    strings.Join(item.Actors, ", "),
				Category:  "电影",
				SubType:   "热门",
				Source:    "douban",
				DoubanID:  item.ID,
			}

			// 保存到数据库
			if err := s.hotDramaRepo.Upsert(drama); err != nil {
				log.Printf("保存电影数据失败: %v", err)
			} else {
				log.Printf("成功保存电影: %s", item.Title)
			}
		}
	}
}

// processTvData 处理电视剧数据
func (s *Scheduler) processTvData() {
	log.Println("开始处理电视剧数据...")

	// 获取电视剧热门榜单
	tvResult, err := s.doubanService.GetTvRanking("tv", "tv", 0, 20)
	if err != nil {
		log.Printf("获取电视剧榜单失败: %v", err)
		return
	}

	if tvResult.Success && tvResult.Data != nil {
		for _, item := range tvResult.Data.Items {
			drama := &entity.HotDrama{
				Title:     item.Title,
				Rating:    item.Rating.Value,
				Year:      item.Year,
				Directors: strings.Join(item.Directors, ", "),
				Actors:    strings.Join(item.Actors, ", "),
				Category:  "电视剧",
				SubType:   "热门",
				Source:    "douban",
				DoubanID:  item.ID,
			}

			// 保存到数据库
			if err := s.hotDramaRepo.Upsert(drama); err != nil {
				log.Printf("保存电视剧数据失败: %v", err)
			} else {
				log.Printf("成功保存电视剧: %s", item.Title)
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
		ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
		defer ticker.Stop()

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
	factory := panutils.NewPanFactory()
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

	// 使用工厂模式创建对应的网盘服务
	config := &panutils.PanConfig{
		URL:         readyResource.URL,
		Code:        "", // 可以从readyResource中获取
		IsType:      0,  // 转存并分享后的资源信息
		ExpiredType: 1,  // 永久分享
		AdFid:       "",
		Stoken:      "",
	}

	panService, err := factory.CreatePanService(readyResource.URL, config)
	if err != nil {
		log.Printf("创建网盘服务失败: %v", err)
		return err
	}

	// 提取分享ID
	shareID, serviceType := panutils.ExtractShareId(readyResource.URL)
	if serviceType == panutils.NotFound {
		log.Printf("不支持的链接地址: %s", readyResource.URL)
		return nil
	}

	log.Printf("检测到服务类型: %s, 分享ID: %s", serviceType.String(), shareID)

	// 根据服务类型进行不同处理
	switch serviceType {
	case panutils.Quark:
		// 夸克网盘：直接转存
		result, err := panService.Transfer(shareID)
		if err != nil {
			log.Printf("夸克网盘转存失败: %v", err)
			return err
		}

		if !result.Success {
			log.Printf("夸克网盘转存失败: %s", result.Message)
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
				PanID:       s.determinePanID(readyResource.URL),
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

	case panutils.Alipan:
		// 阿里云盘：需要检查URL有效性
		checkResult, _ := CheckURL(readyResource.URL)
		if !checkResult.Status {
			log.Printf("阿里云盘链接无效: %s", readyResource.URL)
			return nil
		}

		// 如果有标题，直接创建资源
		if readyResource.Title != nil && *readyResource.Title != "" {
			resource := &entity.Resource{
				Title:       *readyResource.Title,
				Description: readyResource.Description,
				URL:         readyResource.URL,
				PanID:       s.determinePanID(readyResource.URL),
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

		// 尝试转存获取标题
		result, err := panService.Transfer(shareID)
		if err != nil {
			log.Printf("阿里云盘转存失败: %v", err)
			return err
		}

		if result.Success {
			if resultData, ok := result.Data.(map[string]interface{}); ok {
				title := resultData["title"].(string)
				shareURL := resultData["shareUrl"].(string)

				resource := &entity.Resource{
					Title:       title,
					Description: readyResource.Description,
					URL:         shareURL,
					PanID:       s.determinePanID(readyResource.URL),
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
		}

	default:
		log.Printf("暂不支持的服务类型: %s", serviceType.String())
		return nil
	}

	return nil
}

// determinePanID 根据URL确定平台ID
func (s *Scheduler) determinePanID(url string) *uint {
	url = strings.ToLower(url)

	// 这里可以根据你的平台配置来判断
	// 示例逻辑，你需要根据实际情况调整
	if strings.Contains(url, "pan.baidu.com") {
		panID := uint(1) // 百度网盘
		return &panID
	} else if strings.Contains(url, "www.aliyundrive.com") {
		panID := uint(2) // 阿里云盘
		return &panID
	} else if strings.Contains(url, "pan.quark.cn") {
		panID := uint(3) // 夸克网盘
		return &panID
	}

	return nil
}

// getOrCreateCategory 获取或创建分类
func (s *Scheduler) getOrCreateCategory(categoryName string) (uint, error) {
	// 这里需要实现分类的查找和创建逻辑
	// 由于没有CategoryRepository的注入，这里先返回0
	// 你可以根据需要添加CategoryRepository的依赖
	return 0, nil
}

// IsReadyResourceRunning 检查待处理资源自动处理任务是否在运行
func (s *Scheduler) IsReadyResourceRunning() bool {
	return s.readyResourceRunning
}
