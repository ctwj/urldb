package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// HotDramaScheduler 热播剧调度器
type HotDramaScheduler struct {
	*BaseScheduler
	doubanService *utils.DoubanService
	hotDramaMutex sync.Mutex // 防止热播剧任务重叠执行
}

// NewHotDramaScheduler 创建热播剧调度器
func NewHotDramaScheduler(base *BaseScheduler) *HotDramaScheduler {
	return &HotDramaScheduler{
		BaseScheduler: base,
		doubanService: utils.NewDoubanService(),
		hotDramaMutex: sync.Mutex{},
	}
}

// Start 启动热播剧定时任务
func (h *HotDramaScheduler) Start() {
	if h.IsRunning() {
		utils.Info("热播剧定时任务已在运行中")
		return
	}

	h.SetRunning(true)
	utils.Info("启动热播剧定时任务")

	go func() {
		ticker := time.NewTicker(12 * time.Hour) // 每12小时执行一次
		defer ticker.Stop()

		// 立即执行一次
		h.fetchHotDramaData()

		for {
			select {
			case <-ticker.C:
				// 使用TryLock防止任务重叠执行
				if h.hotDramaMutex.TryLock() {
					go func() {
						defer h.hotDramaMutex.Unlock()
						h.fetchHotDramaData()
					}()
				} else {
					utils.Info("上一次热播剧任务还在执行中，跳过本次执行")
				}
			case <-h.GetStopChan():
				utils.Info("停止热播剧定时任务")
				return
			}
		}
	}()
}

// Stop 停止热播剧定时任务
func (h *HotDramaScheduler) Stop() {
	if !h.IsRunning() {
		utils.Info("热播剧定时任务未在运行")
		return
	}

	h.GetStopChan() <- true
	h.SetRunning(false)
	utils.Info("已发送停止信号给热播剧定时任务")
}

// fetchHotDramaData 获取热播剧数据
func (h *HotDramaScheduler) fetchHotDramaData() {
	utils.Info("开始获取热播剧数据...")

	// 直接处理电影和电视剧数据，不再需要FetchHotDramaNames
	h.processHotDramaNames([]string{})
}

// processHotDramaNames 处理热播剧名称
func (h *HotDramaScheduler) processHotDramaNames(dramaNames []string) {
	utils.Info("开始处理热播剧数据，共 %d 个", len(dramaNames))

	// 收集所有数据
	var allDramas []*entity.HotDrama

	// 获取电影数据
	movieDramas := h.processMovieData()
	allDramas = append(allDramas, movieDramas...)

	// 获取电视剧数据
	tvDramas := h.processTvData()
	allDramas = append(allDramas, tvDramas...)

	// 清空数据库
	utils.Info("准备清空数据库，当前共有 %d 条数据", len(allDramas))
	if err := h.hotDramaRepo.DeleteAll(); err != nil {
		utils.Error(fmt.Sprintf("清空数据库失败: %v", err))
		return
	}
	utils.Info("数据库清空完成")

	// 批量插入所有数据
	if len(allDramas) > 0 {
		utils.Info("开始批量插入 %d 条数据", len(allDramas))
		if err := h.hotDramaRepo.BatchCreate(allDramas); err != nil {
			utils.Error(fmt.Sprintf("批量插入数据失败: %v", err))
		} else {
			utils.Info("成功批量插入 %d 条数据", len(allDramas))
		}
	} else {
		utils.Info("没有数据需要插入")
	}

	utils.Info("热播剧数据处理完成")
}

// processMovieData 处理电影数据
func (h *HotDramaScheduler) processMovieData() []*entity.HotDrama {
	utils.Info("开始处理电影数据...")

	var movieDramas []*entity.HotDrama

	// 使用GetTypePage方法获取电影数据
	movieResult, err := h.doubanService.GetTypePage("movie_top250", "全部")
	if err != nil {
		utils.Error(fmt.Sprintf("获取电影榜单失败: %v", err))
		return movieDramas
	}

	if movieResult.Success && movieResult.Data != nil {
		utils.Info("电影获取到 %d 个数据", len(movieResult.Data.Items))

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
			utils.Info("收集电影: %s (评分: %.1f, 年份: %s, 地区: %s)",
				item.Title, item.Rating.Value, item.Year, item.Region)
		}
	} else {
		utils.Warn("电影获取数据失败或为空")
	}

	utils.Info("电影数据处理完成，共收集 %d 条数据", len(movieDramas))
	return movieDramas
}

// processTvData 处理电视剧数据
func (h *HotDramaScheduler) processTvData() []*entity.HotDrama {
	utils.Info("开始处理电视剧数据...")

	var tvDramas []*entity.HotDrama

	// 获取所有tv类型
	tvTypes := h.doubanService.GetAllTvTypes()
	utils.Info("获取到 %d 个tv类型: %v", len(tvTypes), tvTypes)

	// 遍历每个type，分别请求数据
	for _, tvType := range tvTypes {
		utils.Info("正在处理tv类型: %s", tvType)

		// 使用GetTypePage方法请求数据
		tvResult, err := h.doubanService.GetTypePage("tv", tvType)
		if err != nil {
			utils.Error(fmt.Sprintf("获取tv类型 %s 数据失败: %v", tvType, err))
			continue
		}

		if tvResult.Success && tvResult.Data != nil {
			utils.Info("tv类型 %s 获取到 %d 个数据", tvType, len(tvResult.Data.Items))

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
				utils.Info("收集tv类型 %s: %s (评分: %.1f, 年份: %s, 地区: %s)",
					tvType, item.Title, item.Rating.Value, item.Year, item.Region)
			}
		} else {
			utils.Warn("tv类型 %s 获取数据失败或为空", tvType)
		}
	}

	utils.Info("电视剧数据处理完成，共收集 %d 条数据", len(tvDramas))
	return tvDramas
}

// GetHotDramaNames 获取热播剧名称列表（公共方法）
func (h *HotDramaScheduler) GetHotDramaNames() ([]string, error) {
	// 由于删除了FetchHotDramaNames方法，返回空数组
	return []string{}, nil
}
