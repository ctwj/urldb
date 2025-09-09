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

	// 获取最近热门电影数据
	recentMovieDramas := h.processRecentMovies()
	allDramas = append(allDramas, recentMovieDramas...)

	// 获取最近热门剧集数据
	recentTVDramas := h.processRecentTVs()
	allDramas = append(allDramas, recentTVDramas...)

	// 获取最近热门综艺数据
	recentShowDramas := h.processRecentShows()
	allDramas = append(allDramas, recentShowDramas...)

	// 获取豆瓣电影Top250数据
	top250Dramas := h.processTop250Movies()
	allDramas = append(allDramas, top250Dramas...)

	// 设置排名顺序（保持豆瓣返回的顺序）
	for i, drama := range allDramas {
		drama.Rank = i
	}

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

// processRecentMovies 处理最近热门电影数据
func (h *HotDramaScheduler) processRecentMovies() []*entity.HotDrama {
	utils.Info("开始处理最近热门电影数据...")

	var recentMovies []*entity.HotDrama

	items, err := h.doubanService.GetRecentHotMovies()
	if err != nil {
		utils.Error(fmt.Sprintf("获取最近热门电影失败: %v", err))
		return recentMovies
	}

	utils.Info("最近热门电影获取到 %d 个数据", len(items))

	for _, item := range items {
		drama := h.convertDoubanItemToHotDrama(item, "电影", "热门")
		recentMovies = append(recentMovies, drama)
		utils.Info("收集最近热门电影: %s (评分: %.1f, 年份: %s, 地区: %s)",
			item.Title, item.Rating.Value, item.Year, item.Region)
	}

	utils.Info("最近热门电影数据处理完成，共收集 %d 条数据", len(recentMovies))
	return recentMovies
}

// processRecentTVs 处理最近热门剧集数据
func (h *HotDramaScheduler) processRecentTVs() []*entity.HotDrama {
	utils.Info("开始处理最近热门剧集数据...")

	var recentTVs []*entity.HotDrama

	items, err := h.doubanService.GetRecentHotTVs()
	if err != nil {
		utils.Error(fmt.Sprintf("获取最近热门剧集失败: %v", err))
		return recentTVs
	}

	utils.Info("最近热门剧集获取到 %d 个数据", len(items))

	for _, item := range items {
		drama := h.convertDoubanItemToHotDrama(item, "电视剧", "热门")
		recentTVs = append(recentTVs, drama)
		utils.Info("收集最近热门剧集: %s (评分: %.1f, 年份: %s, 地区: %s)",
			item.Title, item.Rating.Value, item.Year, item.Region)
	}

	utils.Info("最近热门剧集数据处理完成，共收集 %d 条数据", len(recentTVs))
	return recentTVs
}

// processRecentShows 处理最近热门综艺数据
func (h *HotDramaScheduler) processRecentShows() []*entity.HotDrama {
	utils.Info("开始处理最近热门综艺数据...")

	var recentShows []*entity.HotDrama

	items, err := h.doubanService.GetRecentHotShows()
	if err != nil {
		utils.Error(fmt.Sprintf("获取最近热门综艺失败: %v", err))
		return recentShows
	}

	utils.Info("最近热门综艺获取到 %d 个数据", len(items))

	for _, item := range items {
		drama := h.convertDoubanItemToHotDrama(item, "综艺", "热门")
		recentShows = append(recentShows, drama)
		utils.Info("收集最近热门综艺: %s (评分: %.1f, 年份: %s, 地区: %s)",
			item.Title, item.Rating.Value, item.Year, item.Region)
	}

	utils.Info("最近热门综艺数据处理完成，共收集 %d 条数据", len(recentShows))
	return recentShows
}

// processTop250Movies 处理豆瓣电影Top250数据
func (h *HotDramaScheduler) processTop250Movies() []*entity.HotDrama {
	utils.Info("开始处理豆瓣电影Top250数据...")

	var top250Movies []*entity.HotDrama

	items, err := h.doubanService.GetTop250Movies()
	if err != nil {
		utils.Error(fmt.Sprintf("获取豆瓣电影Top250失败: %v", err))
		return top250Movies
	}

	utils.Info("豆瓣电影Top250获取到 %d 个数据", len(items))

	for _, item := range items {
		drama := h.convertDoubanItemToHotDrama(item, "电影", "Top250")
		top250Movies = append(top250Movies, drama)
		utils.Info("收集豆瓣Top250电影: %s (评分: %.1f, 年份: %s, 地区: %s)",
			item.Title, item.Rating.Value, item.Year, item.Region)
	}

	utils.Info("豆瓣电影Top250数据处理完成，共收集 %d 条数据", len(top250Movies))
	return top250Movies
}

// convertDoubanItemToHotDrama 转换DoubanItem为HotDrama实体
func (h *HotDramaScheduler) convertDoubanItemToHotDrama(item utils.DoubanItem, category, subType string) *entity.HotDrama {
	return &entity.HotDrama{
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
		Category:     category,
		SubType:      subType,
		Source:       "douban",
		DoubanID:     item.ID,
		DoubanURI:    item.URI,
	}
}

// GetHotDramaNames 获取热播剧名称列表（公共方法）
func (h *HotDramaScheduler) GetHotDramaNames() ([]string, error) {
	// 由于删除了FetchHotDramaNames方法，返回空数组
	return []string{}, nil
}
