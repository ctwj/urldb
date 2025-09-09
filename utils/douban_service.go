package utils

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// top250
// api: https://m.douban.com/rexxar/api/v2/subject_collection/movie_top250/items?start=0&count=10&items_only=1&type_tag=&for_mobile=1

// 最近热门电影 https://movie.douban.com/explore
// api: https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=20

// 最近热门剧集 https://movie.douban.com/tv/
// api: https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=20&limit=20

// 最近热门综艺
// api: https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?limit=50&category=show&type=show

// DoubanService 豆瓣服务
type DoubanService struct {
	baseURL string
	client  *resty.Client

	// 电影榜单配置 - 4个大类，每个大类下有5个小类
	MovieCategories map[string]map[string]map[string]string

	// 剧集榜单配置 - 2个大类
	TvCategories map[string]map[string]map[string]string
}

// DoubanItem 豆瓣项目
type DoubanItem struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	CardSubtitle string   `json:"card_subtitle"`
	EpisodesInfo string   `json:"episodes_info"`
	IsNew        bool     `json:"is_new"`
	Pic          PicInfo  `json:"pic"`
	Rating       Rating   `json:"rating"`
	Type         string   `json:"type"`
	URI          string   `json:"uri"`
	Year         string   `json:"year"`
	Directors    []string `json:"directors"`
	Actors       []string `json:"actors"`
	Region       string   `json:"region"`
	Genres       []string `json:"genres"`
}

// PicInfo 图片信息
type PicInfo struct {
	Large  string `json:"large"`
	Normal string `json:"normal"`
}

// Rating 评分
type Rating struct {
	Value     float64 `json:"value"`
	Count     int     `json:"count"`
	Max       int     `json:"max"`
	StarCount float64 `json:"star_count"`
}

// DoubanCategory 豆瓣分类
type DoubanCategory struct {
	Category string `json:"category"`
	Selected bool   `json:"selected"`
	Type     string `json:"type"`
	Title    string `json:"title"`
}

// DoubanResponse 豆瓣响应
type DoubanResponse struct {
	Items      []DoubanItem     `json:"items"`
	Categories []DoubanCategory `json:"categories"`
	Total      int              `json:"total"`
	IsMockData bool             `json:"is_mock_data,omitempty"`
	MockReason string           `json:"mock_reason,omitempty"`
	Notice     string           `json:"notice,omitempty"`
}

// DoubanResult 豆瓣结果
type DoubanResult struct {
	Success bool            `json:"success"`
	Data    *DoubanResponse `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

// NewDoubanService 创建新的豆瓣服务
func NewDoubanService() *DoubanService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetHeaders(map[string]string{
		"User-Agent":       "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
		"Referer":          "https://m.douban.com/",
		"Accept":           "application/json, text/plain, */*",
		"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept-Encoding":  "gzip, deflate",
		"Connection":       "keep-alive",
		"Sec-Fetch-Dest":   "empty",
		"Sec-Fetch-Mode":   "cors",
		"Sec-Fetch-Site":   "same-origin",
		"Cache-Control":    "no-cache",
		"Pragma":           "no-cache",
		"X-Requested-With": "XMLHttpRequest",
		"Origin":           "https://m.douban.com",
	})

	// 启用自动解压缩
	client.SetDisableWarn(true)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	return &DoubanService{
		baseURL: "https://m.douban.com/rexxar/api/v2",
		client:  client,
	}
}

// GetRecentHotMovies fetches recent hot movies
func (ds *DoubanService) GetRecentHotMovies() ([]DoubanItem, error) {
	url := "https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie"
	params := map[string]string{
		"start": "0",
		"limit": "20",
	}
	items := []DoubanItem{}
	for {
		pageItems, total, err := ds.fetchPage(url, params)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if len(items) >= total {
			break
		}
		start := len(items)
		params["start"] = strconv.Itoa(start)
	}
	return items, nil
}

// GetRecentHotTVs fetches recent hot TV shows
func (ds *DoubanService) GetRecentHotTVs() ([]DoubanItem, error) {
	url := "https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv"
	params := map[string]string{
		"start": "0",
		"limit": "300",
	}
	items := []DoubanItem{}
	for {
		pageItems, total, err := ds.fetchPage(url, params)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if len(items) >= total {
			break
		}
		start := len(items)
		params["start"] = strconv.Itoa(start)
	}
	return items, nil
}

// GetRecentHotShows fetches recent hot shows
func (ds *DoubanService) GetRecentHotShows() ([]DoubanItem, error) {
	url := "https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv"
	params := map[string]string{
		"limit":    "300",
		"category": "show",
		"type":     "show",
		"start":    "0",
	}
	items := []DoubanItem{}
	for {
		pageItems, total, err := ds.fetchPage(url, params)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if len(items) >= total {
			break
		}
		start := len(items)
		params["start"] = strconv.Itoa(start)
	}
	return items, nil
}

// GetTop250Movies fetches top 250 movies
func (ds *DoubanService) GetTop250Movies() ([]DoubanItem, error) {
	url := "https://m.douban.com/rexxar/api/v2/subject_collection/movie_top250/items"
	params := map[string]string{
		"start":      "0",
		"count":      "250",
		"items_only": "1",
		"type_tag":   "",
		"for_mobile": "1",
	}
	items, _, err := ds.fetchPage(url, params)
	return items, err
}

// fetchPage fetches a page of items from a given URL and parameters
func (ds *DoubanService) fetchPage(url string, params map[string]string) ([]DoubanItem, int, error) {
	var response *resty.Response
	var err error

	response, err = ds.client.R().
		SetQueryParams(params).
		Get(url)

	if err != nil {
		return nil, 0, err
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(response.Body(), &apiResponse); err != nil {
		return nil, 0, err
	}

	items := ds.extractItems(apiResponse)
	total := ds.extractTotal(apiResponse)

	return items, total, nil
}

// extractTotal extracts the total number of items from the API response
func (ds *DoubanService) extractTotal(response map[string]interface{}) int {
	if totalData, ok := response["total"]; ok {
		if totalFloat, ok := totalData.(float64); ok {
			return int(totalFloat)
		}
	}
	return 0
}

// extractItems 从API响应中提取项目列表
func (ds *DoubanService) extractItems(response map[string]interface{}) []DoubanItem {
	var items []DoubanItem

	// 根据实际接口返回格式，数据在 subject_collection_items 字段中
	if itemsData, ok := response["subject_collection_items"]; ok {
		if itemsBytes, err := json.Marshal(itemsData); err == nil {
			if err := json.Unmarshal(itemsBytes, &items); err != nil {
				log.Printf("解析subject_collection_items字段失败: %v", err)
			}
		}
	} else if itemsData, ok := response["items"]; ok {
		// 兼容旧的items字段
		if itemsBytes, err := json.Marshal(itemsData); err == nil {
			if err := json.Unmarshal(itemsBytes, &items); err != nil {
				log.Printf("解析items字段失败: %v", err)
			}
		}
	} else if subjectsData, ok := response["subjects"]; ok {
		// 兼容subjects字段
		if subjectsBytes, err := json.Marshal(subjectsData); err == nil {
			if err := json.Unmarshal(subjectsBytes, &items); err != nil {
				log.Printf("解析subjects字段失败: %v", err)
			}
		}
	}

	log.Printf("从API响应中提取到 %d 个项目", len(items))

	// 解析每个项目的card_subtitle，提取年份、地区、类型、导演、演员信息
	for i := range items {
		ds.parseCardSubtitle(&items[i])
	}

	return items
}

// parseCardSubtitle 解析card_subtitle字段
func (ds *DoubanService) parseCardSubtitle(item *DoubanItem) {
	if item.CardSubtitle == "" {
		return
	}

	// card_subtitle格式: "2025 / 中国大陆 / 剧情 爱情 / 丁梓光 / 杨紫 李现"
	parts := strings.Split(item.CardSubtitle, " / ")
	if len(parts) >= 4 {
		// 年份
		if len(parts) > 0 {
			item.Year = strings.TrimSpace(parts[0])
		}

		// 地区
		if len(parts) > 1 {
			item.Region = strings.TrimSpace(parts[1])
		}

		// 类型（可能有多个，用空格分隔）
		if len(parts) > 2 {
			genresStr := strings.TrimSpace(parts[2])
			item.Genres = strings.Fields(genresStr)
		}

		// 导演（可能有多个，用空格分隔）
		if len(parts) > 3 {
			directorsStr := strings.TrimSpace(parts[3])
			item.Directors = strings.Fields(directorsStr)
		}

		// 演员（可能有多个，用空格分隔）
		if len(parts) > 4 {
			actorsStr := strings.TrimSpace(parts[4])
			item.Actors = strings.Fields(actorsStr)
		}
	}
}
