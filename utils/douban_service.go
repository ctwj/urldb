package utils

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

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
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Rating    Rating   `json:"rating"`
	Year      string   `json:"year"`
	Directors []string `json:"directors"`
	Actors    []string `json:"actors"`
}

// Rating 评分
type Rating struct {
	Value float64 `json:"value"`
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
		"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Referer":         "https://m.douban.com/",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
	})

	// 初始化电影榜单配置
	movieCategories := map[string]map[string]map[string]string{
		"热门电影": {
			"全部": {"category": "热门", "type": "全部"},
			"华语": {"category": "热门", "type": "华语"},
			"欧美": {"category": "热门", "type": "欧美"},
			"韩国": {"category": "热门", "type": "韩国"},
			"日本": {"category": "热门", "type": "日本"},
		},
		"最新电影": {
			"全部": {"category": "最新", "type": "全部"},
			"华语": {"category": "最新", "type": "华语"},
			"欧美": {"category": "最新", "type": "欧美"},
			"韩国": {"category": "最新", "type": "韩国"},
			"日本": {"category": "最新", "type": "日本"},
		},
		"豆瓣高分": {
			"全部": {"category": "豆瓣高分", "type": "全部"},
			"华语": {"category": "豆瓣高分", "type": "华语"},
			"欧美": {"category": "豆瓣高分", "type": "欧美"},
			"韩国": {"category": "豆瓣高分", "type": "韩国"},
			"日本": {"category": "豆瓣高分", "type": "日本"},
		},
		"冷门佳片": {
			"全部": {"category": "冷门佳片", "type": "全部"},
			"华语": {"category": "冷门佳片", "type": "华语"},
			"欧美": {"category": "冷门佳片", "type": "欧美"},
			"韩国": {"category": "冷门佳片", "type": "韩国"},
			"日本": {"category": "冷门佳片", "type": "日本"},
		},
	}

	// 初始化剧集榜单配置
	tvCategories := map[string]map[string]map[string]string{
		"最近热门剧集": {
			"综合":  {"category": "tv", "type": "tv"},
			"国产剧": {"category": "tv", "type": "tv_domestic"},
			"欧美剧": {"category": "tv", "type": "tv_american"},
			"日剧":  {"category": "tv", "type": "tv_japanese"},
			"韩剧":  {"category": "tv", "type": "tv_korean"},
			"动画":  {"category": "tv", "type": "tv_animation"},
			"纪录片": {"category": "tv", "type": "tv_documentary"},
		},
		"最近热门综艺": {
			"综合": {"category": "show", "type": "show"},
			"国内": {"category": "show", "type": "show_domestic"},
			"国外": {"category": "show", "type": "show_foreign"},
		},
	}

	return &DoubanService{
		baseURL:         "https://m.douban.com/rexxar/api/v2",
		client:          client,
		MovieCategories: movieCategories,
		TvCategories:    tvCategories,
	}
}

// GetMovieRanking 获取电影榜单数据
func (ds *DoubanService) GetMovieRanking(category, rankingType string, start, limit int) (*DoubanResult, error) {
	log.Printf("获取电影榜单: %s - %s, start: %d, limit: %d", category, rankingType, start, limit)

	// 构建请求参数
	params := map[string]string{
		"start": strconv.Itoa(start),
		"limit": strconv.Itoa(limit),
	}

	// 根据不同的category和type添加特定参数
	if category != "热门" || rankingType != "全部" {
		if rankingType != "全部" {
			params["type"] = rankingType
		}
		if category != "热门" {
			params["category"] = category
		}
	}

	var response *resty.Response
	var err error

	// 尝试调用豆瓣API
	response, err = ds.client.R().
		SetQueryParams(params).
		Get(ds.baseURL + "/subject/recent_hot/movie")

	if err != nil {
		log.Printf("豆瓣API调用失败，使用模拟数据: %v", err)
		// 如果豆瓣API调用失败，使用模拟数据
		mockData := ds.getMockMovieData()
		mockData.IsMockData = true
		mockData.MockReason = "API调用失败"

		return &DoubanResult{
			Success: true,
			Data:    mockData,
		}, nil
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(response.Body(), &apiResponse); err != nil {
		log.Printf("解析API响应失败: %v", err)
		mockData := ds.getMockMovieData()
		mockData.IsMockData = true
		mockData.MockReason = "解析API响应失败"

		return &DoubanResult{
			Success: true,
			Data:    mockData,
		}, nil
	}

	// 处理豆瓣移动端API的响应格式
	items := ds.extractItems(apiResponse)
	categories := ds.extractCategories(apiResponse)

	// 如果没有获取到真实数据，使用模拟数据
	isMockData := false
	mockReason := ""

	if len(items) == 0 {
		log.Println("API返回空数据，使用模拟数据")
		mockData := ds.getMockMovieData()
		items = mockData.Items
		isMockData = true
		mockReason = "API返回空数据"
	}

	// 如果没有获取到categories，使用默认的电影分类
	if len(categories) == 0 {
		categories = []DoubanCategory{
			{Category: "热门", Selected: true, Type: "全部", Title: "热门"},
			{Category: "最新", Selected: false, Type: "全部", Title: "最新"},
			{Category: "豆瓣高分", Selected: false, Type: "全部", Title: "豆瓣高分"},
			{Category: "冷门佳片", Selected: false, Type: "全部", Title: "冷门佳片"},
			{Category: "热门", Selected: false, Type: "华语", Title: "华语"},
			{Category: "热门", Selected: false, Type: "欧美", Title: "欧美"},
			{Category: "热门", Selected: false, Type: "韩国", Title: "韩国"},
			{Category: "热门", Selected: false, Type: "日本", Title: "日本"},
		}
	}

	// 根据请求的category和type更新selected状态
	for i := range categories {
		categories[i].Selected = categories[i].Category == category && categories[i].Type == rankingType
	}

	// 限制返回数量
	if len(items) > limit {
		items = items[:limit]
	}

	result := &DoubanResponse{
		Items:      items,
		Total:      len(items),
		Categories: categories,
		IsMockData: isMockData,
		MockReason: mockReason,
	}

	if isMockData {
		result.Notice = "⚠️ 这是模拟数据，非豆瓣实时数据"
	}

	return &DoubanResult{
		Success: true,
		Data:    result,
	}, nil
}

// GetTvRanking 获取电视剧榜单数据
func (ds *DoubanService) GetTvRanking(category, rankingType string, start, limit int) (*DoubanResult, error) {
	log.Printf("获取电视剧榜单: %s - %s, start: %d, limit: %d", category, rankingType, start, limit)

	// 构建请求参数
	params := map[string]string{
		"start": strconv.Itoa(start),
		"limit": strconv.Itoa(limit),
	}

	// 根据不同的category和type添加特定参数
	if category != "tv" || rankingType != "tv" {
		if rankingType != "tv" {
			params["type"] = rankingType
		}
		if category != "tv" {
			params["category"] = category
		}
	}

	var response *resty.Response
	var err error

	// 尝试调用豆瓣API
	response, err = ds.client.R().
		SetQueryParams(params).
		Get(ds.baseURL + "/subject/recent_hot/tv")

	if err != nil {
		log.Printf("豆瓣TV API调用失败，使用模拟数据: %v", err)
		mockData := ds.getMockTvData()
		mockData.IsMockData = true
		mockData.MockReason = "API调用失败"

		return &DoubanResult{
			Success: true,
			Data:    mockData,
		}, nil
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(response.Body(), &apiResponse); err != nil {
		log.Printf("解析TV API响应失败: %v", err)
		mockData := ds.getMockTvData()
		mockData.IsMockData = true
		mockData.MockReason = "解析API响应失败"

		return &DoubanResult{
			Success: true,
			Data:    mockData,
		}, nil
	}

	// 处理豆瓣移动端API的响应格式
	items := ds.extractItems(apiResponse)
	categories := ds.extractCategories(apiResponse)

	// 如果没有获取到真实数据，使用模拟数据
	isMockData := false
	mockReason := ""

	if len(items) == 0 {
		log.Println("TV API返回空数据，使用模拟数据")
		mockData := ds.getMockTvData()
		items = mockData.Items
		isMockData = true
		mockReason = "API返回空数据"
	}

	// 如果没有获取到categories，使用默认的电视剧分类
	if len(categories) == 0 {
		categories = []DoubanCategory{
			{Category: "tv", Selected: true, Type: "tv", Title: "综合"},
			{Category: "tv", Selected: false, Type: "tv_domestic", Title: "国产剧"},
			{Category: "show", Selected: false, Type: "show", Title: "综艺"},
			{Category: "tv", Selected: false, Type: "tv_american", Title: "欧美剧"},
			{Category: "tv", Selected: false, Type: "tv_japanese", Title: "日剧"},
			{Category: "tv", Selected: false, Type: "tv_korean", Title: "韩剧"},
			{Category: "tv", Selected: false, Type: "tv_animation", Title: "动画"},
			{Category: "tv", Selected: false, Type: "tv_documentary", Title: "纪录片"},
		}
	}

	// 根据请求的category和type更新selected状态
	for i := range categories {
		categories[i].Selected = categories[i].Category == category && categories[i].Type == rankingType
	}

	// 限制返回数量
	if len(items) > limit {
		items = items[:limit]
	}

	result := &DoubanResponse{
		Items:      items,
		Total:      len(items),
		Categories: categories,
		IsMockData: isMockData,
		MockReason: mockReason,
	}

	if isMockData {
		result.Notice = "⚠️ 这是模拟数据，非豆瓣实时数据"
	}

	return &DoubanResult{
		Success: true,
		Data:    result,
	}, nil
}

// GetMovieCategories 获取支持的电影类别
func (ds *DoubanService) GetMovieCategories() map[string]map[string]map[string]string {
	return ds.MovieCategories
}

// GetTvCategories 获取支持的电视剧类别
func (ds *DoubanService) GetTvCategories() map[string]map[string]map[string]string {
	return ds.TvCategories
}

// GetAllCategories 获取所有支持的类别
func (ds *DoubanService) GetAllCategories() map[string]interface{} {
	return map[string]interface{}{
		"movie": ds.GetMovieCategories(),
		"tv":    ds.GetTvCategories(),
	}
}

// GetMovieSubCategories 获取电影特定大类下的小类
func (ds *DoubanService) GetMovieSubCategories(mainCategory string) map[string]map[string]string {
	return ds.MovieCategories[mainCategory]
}

// GetTvSubCategories 获取剧集特定大类下的小类
func (ds *DoubanService) GetTvSubCategories(mainCategory string) map[string]map[string]string {
	return ds.TvCategories[mainCategory]
}

// getMockMovieData 获取模拟电影数据
func (ds *DoubanService) getMockMovieData() *DoubanResponse {
	return &DoubanResponse{
		Notice: "⚠️ 这是模拟数据，非豆瓣实时数据",
		Items: []DoubanItem{
			{
				ID:        "1292052",
				Title:     "肖申克的救赎",
				Rating:    Rating{Value: 9.7},
				Year:      "1994",
				Directors: []string{"弗兰克·德拉邦特"},
				Actors:    []string{"蒂姆·罗宾斯", "摩根·弗里曼"},
			},
			{
				ID:        "1291546",
				Title:     "霸王别姬",
				Rating:    Rating{Value: 9.6},
				Year:      "1993",
				Directors: []string{"陈凯歌"},
				Actors:    []string{"张国荣", "张丰毅", "巩俐"},
			},
			{
				ID:        "1295644",
				Title:     "阿甘正传",
				Rating:    Rating{Value: 9.5},
				Year:      "1994",
				Directors: []string{"罗伯特·泽米吉斯"},
				Actors:    []string{"汤姆·汉克斯", "罗宾·怀特"},
			},
		},
		Total: 3,
	}
}

// getMockTvData 获取模拟电视剧数据
func (ds *DoubanService) getMockTvData() *DoubanResponse {
	return &DoubanResponse{
		Notice: "⚠️ 这是模拟数据，非豆瓣实时数据",
		Items: []DoubanItem{
			{
				ID:        "26794435",
				Title:     "请回答1988",
				Rating:    Rating{Value: 9.7},
				Year:      "2015",
				Directors: []string{"申元浩"},
				Actors:    []string{"李惠利", "朴宝剑", "柳俊烈"},
			},
			{
				ID:        "1309163",
				Title:     "大明王朝1566",
				Rating:    Rating{Value: 9.7},
				Year:      "2007",
				Directors: []string{"张黎"},
				Actors:    []string{"陈宝国", "黄志忠", "王庆祥"},
			},
			{
				ID:        "1309169",
				Title:     "亮剑",
				Rating:    Rating{Value: 9.3},
				Year:      "2005",
				Directors: []string{"陈健", "张前"},
				Actors:    []string{"李幼斌", "何政军", "张光北"},
			},
		},
		Total: 3,
	}
}

// extractItems 从API响应中提取项目列表
func (ds *DoubanService) extractItems(response map[string]interface{}) []DoubanItem {
	var items []DoubanItem

	// 尝试从不同的字段获取items
	if itemsData, ok := response["items"]; ok {
		if itemsBytes, err := json.Marshal(itemsData); err == nil {
			json.Unmarshal(itemsBytes, &items)
		}
	} else if subjectsData, ok := response["subjects"]; ok {
		if subjectsBytes, err := json.Marshal(subjectsData); err == nil {
			json.Unmarshal(subjectsBytes, &items)
		}
	}

	return items
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

// FetchHotDramaNames 获取热播剧名字（用于定时任务）
func (ds *DoubanService) FetchHotDramaNames() ([]string, error) {
	var dramaNames []string

	// 获取电影热门榜单
	movieResult, err := ds.GetMovieRanking("热门", "全部", 0, 10)
	if err != nil {
		log.Printf("获取电影榜单失败: %v", err)
	} else if movieResult.Success && movieResult.Data != nil {
		for _, item := range movieResult.Data.Items {
			dramaNames = append(dramaNames, item.Title)
		}
	}

	// 获取电视剧热门榜单
	tvResult, err := ds.GetTvRanking("tv", "tv", 0, 10)
	if err != nil {
		log.Printf("获取电视剧榜单失败: %v", err)
	} else if tvResult.Success && tvResult.Data != nil {
		for _, item := range tvResult.Data.Items {
			dramaNames = append(dramaNames, item.Title)
		}
	}

	log.Printf("获取到 %d 个热播剧名字", len(dramaNames))
	return dramaNames, nil
}
