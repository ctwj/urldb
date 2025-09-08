package utils

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

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

	// 初始化剧集榜单配置
	tvCategories := map[string]map[string]map[string]string{
		"最近热门剧集": {
			// "国产剧": {"category": "tv", "type": "tv_domestic"},
			"欧美剧": {"category": "tv", "type": "tv_american"},
			"日剧":  {"category": "tv", "type": "tv_japanese"},
			"韩剧":  {"category": "tv", "type": "tv_korean"},
			"动画":  {"category": "tv", "type": "tv_animation"},
			"纪录片": {"category": "tv", "type": "tv_documentary"},
		},
		"最近热门综艺": {
			// "综合": {"category": "show", "type": "show"},
			"国内": {"category": "show", "type": "show_domestic"},
			"国外": {"category": "show", "type": "show_foreign"},
		},
	}

	return &DoubanService{
		baseURL:      "https://m.douban.com/rexxar/api/v2",
		client:       client,
		TvCategories: tvCategories,
	}
}

// GetTypePage 获取指定类型的数据
func (ds *DoubanService) GetTypePage(category, rankingType string) (*DoubanResult, error) {
	// 构建请求参数
	params := map[string]string{
		"start":  "0",
		"limit":  "50",
		"os":     "window",
		"_":      "0",
		"loc_id": "108288",
	}

	Debug("请求参数: %+v", params)
	Debug("请求URL: %s/subject_collection/%s/items", ds.baseURL, rankingType)

	var response *resty.Response
	var err error

	// 尝试调用豆瓣API
	Debug("开始发送HTTP请求...")
	response, err = ds.client.R().
		SetQueryParams(params).
		Get(ds.baseURL + "/subject_collection/" + rankingType + "/items")

	if err != nil {
		Error("=== 豆瓣API调用失败 ===")
		Error("错误详情: %v", err)
		return &DoubanResult{
			Success: false,
			Message: "API调用失败: " + err.Error(),
		}, nil
	}

	Debug("=== HTTP请求成功 ===")
	Debug("响应状态码: %d", response.StatusCode())
	Debug("响应体长度: %d bytes", len(response.Body()))

	// 记录响应体的前500个字符用于调试
	responseBody := string(response.Body())
	Debug("响应体原始长度: %d 字符", len(responseBody))

	if len(responseBody) > 500 {
		Debug("响应体前500字符: %s...", responseBody[:500])
	} else {
		Debug("完整响应体: %s", responseBody)
	}

	// 检查响应体是否包含有效JSON
	if len(responseBody) == 0 {
		Warn("=== 响应体为空 ===")
		return &DoubanResult{
			Success: false,
			Message: "响应体为空",
		}, nil
	}

	// 尝试解析JSON
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(response.Body(), &apiResponse); err != nil {
		Error("=== 解析API响应失败 ===")
		Error("JSON解析错误: %v", err)
		Debug("响应体内容: %s", string(response.Body()))

		// 尝试检查是否是HTML错误页面
		if len(responseBody) > 100 && (strings.Contains(responseBody, "<html>") || strings.Contains(responseBody, "<!DOCTYPE")) {
			Warn("检测到HTML响应，可能是错误页面")
			return &DoubanResult{
				Success: false,
				Message: "返回HTML错误页面",
			}, nil
		}

		return &DoubanResult{
			Success: false,
			Message: "解析API响应失败: " + err.Error(),
		}, nil
	}

	log.Printf("=== JSON解析成功 ===")
	log.Printf("解析后的数据结构: %+v", apiResponse)

	// 打印完整的API响应JSON
	log.Printf("=== 完整API响应JSON ===")
	if responseBytes, err := json.MarshalIndent(apiResponse, "", "  "); err == nil {
		log.Printf("完整响应:\n%s", string(responseBytes))
	} else {
		log.Printf("序列化响应失败: %v", err)
	}

	// 处理豆瓣移动端API的响应格式
	items := ds.extractItems(apiResponse)
	categories := ds.extractCategories(apiResponse)

	log.Printf("提取到的数据数量: %d", len(items))
	log.Printf("提取到的分类数量: %d", len(categories))

	// 如果没有获取到真实数据，返回空结果
	if len(items) == 0 {
		log.Printf("=== API返回空数据 ===")
		return &DoubanResult{
			Success: true,
			Data: &DoubanResponse{
				Items:      []DoubanItem{},
				Total:      0,
				Categories: []DoubanCategory{},
			},
		}, nil
	}

	// 如果没有获取到categories，使用默认分类
	if len(categories) == 0 {
		log.Printf("=== 使用默认分类 ===")
		categories = []DoubanCategory{
			{Category: category, Selected: true, Type: rankingType, Title: rankingType},
		}
	}

	// 根据请求的category和type更新selected状态
	for i := range categories {
		categories[i].Selected = categories[i].Category == category && categories[i].Type == rankingType
	}

	// 限制返回数量（最多50条）
	limit := 50
	if len(items) > limit {
		log.Printf("限制返回数量从 %d 到 %d", len(items), limit)
		items = items[:limit]
	}

	// 获取总数，优先使用API返回的total字段
	total := len(items)
	if totalData, ok := apiResponse["total"]; ok {
		if totalFloat, ok := totalData.(float64); ok {
			total = int(totalFloat)
		}
	}

	result := &DoubanResponse{
		Items:      items,
		Total:      total,
		Categories: categories,
		IsMockData: false,
		MockReason: "",
	}

	log.Printf("=== 数据获取完成 ===")
	log.Printf("最终返回数据数量: %d", len(items))

	return &DoubanResult{
		Success: true,
		Data:    result,
	}, nil
}

// GetTvByType 获取指定type的全部剧集数据
func (ds *DoubanService) GetTvByType(tvType string) ([]map[string]interface{}, error) {
	url := ds.baseURL + "/subject_collection/" + tvType + "/items"
	params := map[string]string{
		"start": "0",
		"limit": "1000", // 假设不会超过1000条
	}

	resp, err := ds.client.R().
		SetQueryParams(params).
		Get(url)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	items, ok := result["subject_collection_items"].([]interface{})
	if !ok {
		return nil, nil // 没有数据
	}

	// 转换为[]map[string]interface{}
	var out []map[string]interface{}
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// GetAllTvTypes 获取所有tv类型（type列表）
func (ds *DoubanService) GetAllTvTypes() []string {
	types := []string{}
	for _, sub := range ds.TvCategories {
		for _, v := range sub {
			if t, ok := v["type"]; ok {
				types = append(types, t)
			}
		}
	}
	return types
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

// extractCategories 从API响应中提取分类列表
func (ds *DoubanService) extractCategories(response map[string]interface{}) []DoubanCategory {
	var categories []DoubanCategory

	if categoriesData, ok := response["categories"]; ok {
		if categoriesBytes, err := json.Marshal(categoriesData); err == nil {
			json.Unmarshal(categoriesBytes, &categories)
		}
	}

	return categories
}
