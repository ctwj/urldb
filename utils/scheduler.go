package utils

import (
	"log"
	"res_db/db/entity"
	"res_db/db/repo"
	"strings"
	"time"
)

// Scheduler 定时任务管理器
type Scheduler struct {
	doubanService *DoubanService
	hotDramaRepo  repo.HotDramaRepository
	stopChan      chan bool
	isRunning     bool
}

// NewScheduler 创建新的定时任务管理器
func NewScheduler(hotDramaRepo repo.HotDramaRepository) *Scheduler {
	return &Scheduler{
		doubanService: NewDoubanService(),
		hotDramaRepo:  hotDramaRepo,
		stopChan:      make(chan bool),
		isRunning:     false,
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
				s.fetchHotDramaData()
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
