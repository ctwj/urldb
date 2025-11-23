package google

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/searchconsole/v1"
)

// Client Google Search Console API客户端
type Client struct {
	service *searchconsole.Service
	SiteURL string
}

// Config 配置信息
type Config struct {
	CredentialsFile string `json:"credentials_file"`
	SiteURL         string `json:"site_url"`
	TokenFile       string `json:"token_file"`
}

// URLInspectionRequest URL检查请求
type URLInspectionRequest struct {
	InspectionURL string `json:"inspectionUrl"`
	SiteURL       string `json:"siteUrl"`
	LanguageCode  string `json:"languageCode"`
}

// URLInspectionResult URL检查结果
type URLInspectionResult struct {
	IndexStatusResult struct {
		IndexingState string `json:"indexingState"`
		LastCrawled   string `json:"lastCrawled"`
		CrawlErrors   []struct {
			ErrorCode string `json:"errorCode"`
		} `json:"crawlErrors"`
	} `json:"indexStatusResult"`
	MobileUsabilityResult struct {
		MobileFriendly bool `json:"mobileFriendly"`
	} `json:"mobileUsabilityResult"`
	RichResultsResult struct {
		Detected struct {
			Items []struct {
				RichResultType string `json:"richResultType"`
			} `json:"items"`
		} `json:"detected"`
	} `json:"richResultsResult"`
}

// NewClient 创建新的客户端
func NewClient(config *Config) (*Client, error) {
	ctx := context.Background()

	// 读取认证文件
	credentials, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("读取认证文件失败: %v", err)
	}

	// 创建OAuth2配置
	oauthConfig, err := google.ConfigFromJSON(credentials, searchconsole.WebmastersScope)
	if err != nil {
		return nil, fmt.Errorf("创建OAuth配置失败: %v", err)
	}

	// 尝试从文件读取token
	token, err := tokenFromFile(config.TokenFile)
	if err != nil {
		// 如果没有token，启动web认证流程
		token = getTokenFromWeb(oauthConfig)
		saveToken(config.TokenFile, token)
	}

	// 创建HTTP客户端
	client := oauthConfig.Client(ctx, token)

	// 创建Search Console服务
	service, err := searchconsole.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建Search Console服务失败: %v", err)
	}

	return &Client{
		service: service,
		SiteURL: config.SiteURL,
	}, nil
}

// InspectURL 检查URL索引状态
func (c *Client) InspectURL(url string) (*URLInspectionResult, error) {
	request := &searchconsole.InspectUrlIndexRequest{
		InspectionUrl: url,
		SiteUrl:       c.SiteURL,
		LanguageCode:  "zh-CN",
	}

	call := c.service.UrlInspection.Index.Inspect(request)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("检查URL失败: %v", err)
	}

	// 转换响应
	result := &URLInspectionResult{}
	if response.InspectionResult != nil {
		if response.InspectionResult.IndexStatusResult != nil {
			result.IndexStatusResult.IndexingState = string(response.InspectionResult.IndexStatusResult.IndexingState)
			if response.InspectionResult.IndexStatusResult.LastCrawlTime != "" {
				result.IndexStatusResult.LastCrawled = response.InspectionResult.IndexStatusResult.LastCrawlTime
			}
		}

		if response.InspectionResult.MobileUsabilityResult != nil {
			result.MobileUsabilityResult.MobileFriendly = response.InspectionResult.MobileUsabilityResult.Verdict == "MOBILE_USABILITY_VERdict_PASS"
		}

		if response.InspectionResult.RichResultsResult != nil && response.InspectionResult.RichResultsResult.Verdict != "RICH_RESULTS_VERdict_PASS" {
			// 如果有富媒体结果检查信息
			result.RichResultsResult.Detected.Items = append(result.RichResultsResult.Detected.Items, struct {
				RichResultType string `json:"richResultType"`
			}{
				RichResultType: "UNKNOWN",
			})
		}
	}

	return result, nil
}

// SubmitSitemap 提交网站地图
func (c *Client) SubmitSitemap(sitemapURL string) error {
	call := c.service.Sitemaps.Submit(c.SiteURL, sitemapURL)
	err := call.Do()
	if err != nil {
		return fmt.Errorf("提交网站地图失败: %v", err)
	}

	return nil
}

// GetSites 获取已验证的网站列表
func (c *Client) GetSites() ([]*searchconsole.WmxSite, error) {
	call := c.service.Sites.List()
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("获取网站列表失败: %v", err)
	}

	return response.SiteEntry, nil
}

// GetSearchAnalytics 获取搜索分析数据
func (c *Client) GetSearchAnalytics(startDate, endDate string) (*searchconsole.SearchAnalyticsQueryResponse, error) {
	request := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      "web",
	}

	call := c.service.Searchanalytics.Query(c.SiteURL, request)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("获取搜索分析数据失败: %v", err)
	}

	return response, nil
}

// getTokenFromWeb 通过web流程获取token
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("请在浏览器中访问以下URL进行认证:\n%s\n", authURL)
	fmt.Printf("输入授权代码: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		panic(fmt.Sprintf("读取授权代码失败: %v", err))
	}

	token, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		panic(fmt.Sprintf("获取token失败: %v", err))
	}

	return token
}

// tokenFromFile 从文件读取token
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// saveToken 保存token到文件
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("保存凭证文件到: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(fmt.Sprintf("无法保存凭证文件: %v", err))
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

// BatchInspectURL 批量检查URL状态
func (c *Client) BatchInspectURL(urls []string, callback func(url string, result *URLInspectionResult, err error)) {
	semaphore := make(chan struct{}, 5) // 限制并发数

	for _, url := range urls {
		go func(u string) {
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			result, err := c.InspectURL(u)
			callback(u, result, err)
		}(url)

		// 避免请求过快
		time.Sleep(100 * time.Millisecond)
	}

	// 等待所有goroutine完成
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}
}