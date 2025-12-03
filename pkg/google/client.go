package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/searchconsole/v1"
	"google.golang.org/api/googleapi"
)

// Client Google Search Console API客户端
type Client struct {
	service       *searchconsole.Service
	indexingService *indexing.Service
	SiteURL       string
}

// Config 配置信息
type Config struct {
	CredentialsFile string `json:"credentials_file"`
	SiteURL         string `json:"website_url"`
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

	fmt.Printf("[GOOGLE-CLIENT] 初始化Google客户端，凭据文件: %s\n", config.CredentialsFile)
	fmt.Printf("[GOOGLE-CLIENT] 目标站点URL: %s\n", config.SiteURL)

	// 读取认证文件
	credentials, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("读取认证文件失败: %v", err)
	}
	fmt.Printf("[GOOGLE-CLIENT] 成功读取凭据文件，大小: %d bytes\n", len(credentials))

	// 检查凭据类型
	var credentialsMap map[string]interface{}
	if err := json.Unmarshal(credentials, &credentialsMap); err != nil {
		return nil, fmt.Errorf("解析凭据失败: %v", err)
	}

	// 根据凭据类型创建不同的配置
	credType, ok := credentialsMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("未知的凭据类型")
	}
	fmt.Printf("[GOOGLE-CLIENT] 凭据类型: %s\n", credType)

	// 提取服务账号邮箱（如果有的话）
	if credType == "service_account" {
		if email, exists := credentialsMap["client_email"]; exists {
			fmt.Printf("[GOOGLE-CLIENT] 服务账号邮箱: %s\n", email)
		}
	}

	// 组合作用域，包含Search Console和Indexing API
	scopes := []string{
		searchconsole.WebmastersScope,
		indexing.IndexingScope,
	}
	fmt.Printf("[GOOGLE-CLIENT] 使用的作用域: %v\n", scopes)

	var client *http.Client
	if credType == "service_account" {
		// 服务账号凭据
		fmt.Printf("[GOOGLE-CLIENT] 创建服务账号JWT配置...\n")
		jwtConfig, err := google.JWTConfigFromJSON(credentials, scopes...)
		if err != nil {
			return nil, fmt.Errorf("创建JWT配置失败: %v", err)
		}
		client = jwtConfig.Client(ctx)
		fmt.Printf("[GOOGLE-CLIENT] 服务账号客户端创建成功\n")
	} else {
		// OAuth2客户端凭据
		fmt.Printf("[GOOGLE-CLIENT] 创建OAuth2配置...\n")
		oauthConfig, err := google.ConfigFromJSON(credentials, scopes...)
		if err != nil {
			return nil, fmt.Errorf("创建OAuth配置失败: %v", err)
		}

		// 尝试从文件读取token
		token, err := tokenFromFile(config.TokenFile)
		if err != nil {
			// 如果没有token，启动web认证流程
			fmt.Printf("[GOOGLE-CLIENT] 未找到token文件，启动web认证流程...\n")
			token = getTokenFromWeb(oauthConfig)
			saveToken(config.TokenFile, token)
		}

		client = oauthConfig.Client(ctx, token)
		fmt.Printf("[GOOGLE-CLIENT] OAuth2客户端创建成功\n")
	}

	// 测试客户端连接
	fmt.Printf("[GOOGLE-CLIENT] 测试API客户端连接...\n")
	testURL := "https://www.googleapis.com/auth/webmasters"
	resp, err := client.Get(testURL)
	if err != nil {
		fmt.Printf("[GOOGLE-CLIENT] API客户端连接测试失败: %v\n", err)
		return nil, fmt.Errorf("API客户端连接测试失败: %v", err)
	}
	resp.Body.Close()
	fmt.Printf("[GOOGLE-CLIENT] API客户端连接测试成功\n")

	// 创建Search Console服务
	fmt.Printf("[GOOGLE-CLIENT] 创建Search Console服务...\n")
	service, err := searchconsole.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建Search Console服务失败: %v", err)
	}
	fmt.Printf("[GOOGLE-CLIENT] Search Console服务创建成功\n")

	// 创建Indexing API服务
	fmt.Printf("[GOOGLE-CLIENT] 创建Indexing API服务...\n")
	indexingService, err := indexing.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建Indexing服务失败: %v", err)
	}
	fmt.Printf("[GOOGLE-CLIENT] Indexing API服务创建成功\n")

	// 验证站点访问权限
	if config.SiteURL != "" {
		fmt.Printf("[GOOGLE-CLIENT] 验证站点访问权限: %s\n", config.SiteURL)
		sites, err := service.Sites.List().Do()
		if err != nil {
			fmt.Printf("[GOOGLE-CLIENT] 获取站点列表失败: %v\n", err)
			return nil, fmt.Errorf("无法访问Google Search Console API，请检查服务账号权限: %v", err)
		}

		fmt.Printf("[GOOGLE-CLIENT] 可访问的站点数量: %d\n", len(sites.SiteEntry))
		for i, site := range sites.SiteEntry {
			if i < 3 { // 只显示前3个站点
				fmt.Printf("[GOOGLE-CLIENT] 站点 %d: %s (权限级别: %s)\n", i+1, site.SiteUrl, site.PermissionLevel)
			}
		}

		// 检查目标站点是否在列表中
		targetSiteFound := false
		for _, site := range sites.SiteEntry {
			if strings.Contains(site.SiteUrl, strings.TrimPrefix(config.SiteURL, "https://")) ||
			   strings.Contains(strings.TrimPrefix(config.SiteURL, "https://"), site.SiteUrl) {
				targetSiteFound = true
				fmt.Printf("[GOOGLE-CLIENT] 目标站点权限验证成功: %s\n", site.SiteUrl)
				break
			}
		}

		if !targetSiteFound {
			fmt.Printf("[GOOGLE-CLIENT] 警告: 目标站点 %s 未在已验证站点列表中\n", config.SiteURL)
		}
	}

	fmt.Printf("[GOOGLE-CLIENT] Google客户端初始化完成\n")
	return &Client{
		service:         service,
		indexingService: indexingService,
		SiteURL:         config.SiteURL,
	}, nil
}

// InspectURL 检查URL索引状态
func (c *Client) InspectURL(url string) (*URLInspectionResult, error) {
	fmt.Printf("[GOOGLE-CLIENT] 检查URL索引状态: %s\n", url)
	fmt.Printf("[GOOGLE-CLIENT] 使用站点URL: %s\n", c.SiteURL)

	// 验证URL格式
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("URL格式错误，必须以http://或https://开头: %s", url)
	}

	// 验证站点URL格式
	siteURL := c.SiteURL
	if !strings.HasSuffix(siteURL, "/") {
		siteURL = siteURL + "/"
	}
	if !strings.HasPrefix(siteURL, "sc-domain:") && !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}

	request := &searchconsole.InspectUrlIndexRequest{
		InspectionUrl: url,
		SiteUrl:       siteURL,
		LanguageCode:  "zh-CN",
	}

	fmt.Printf("[GOOGLE-CLIENT] 发送URL检查请求...\n")
	call := c.service.UrlInspection.Index.Inspect(request)
	response, err := call.Do()

	if err != nil {
		fmt.Printf("[GOOGLE-CLIENT] URL检查请求失败: %v\n", err)

		// 尝试解析Google API错误
		if googleErr, ok := err.(*googleapi.Error); ok {
			fmt.Printf("[GOOGLE-CLIENT] Google API错误详情: 代码=%d, 消息=%s\n", googleErr.Code, googleErr.Message)
			for _, e := range googleErr.Errors {
				fmt.Printf("[GOOGLE-CLIENT] 错误项: 原因=%s, 消息=%s\n", e.Reason, e.Message)
			}

			// 根据错误代码提供更具体的错误信息
			switch googleErr.Code {
			case 403:
				if strings.Contains(googleErr.Message, "forbidden") || strings.Contains(googleErr.Message, "permission") {
					return nil, fmt.Errorf("权限不足: 请确保服务账号已被授予该站点的访问权限。详情: %s", googleErr.Message)
				} else if strings.Contains(googleErr.Message, "indexing") {
					return nil, fmt.Errorf("Indexing API权限不足: 请确保服务账号已获得Indexing API访问权限。详情: %s", googleErr.Message)
				}
			case 404:
				return nil, fmt.Errorf("站点未找到: 请确保站点URL正确且已在Google Search Console中验证。详情: %s", googleErr.Message)
			case 400:
				return nil, fmt.Errorf("请求参数错误: %s", googleErr.Message)
			case 429:
				return nil, fmt.Errorf("请求过于频繁，请稍后重试: %s", googleErr.Message)
			default:
				return nil, fmt.Errorf("Google API错误 (代码: %d): %s", googleErr.Code, googleErr.Message)
			}
		}

		return nil, fmt.Errorf("检查URL失败: %v", err)
	}

	fmt.Printf("[GOOGLE-CLIENT] URL检查请求成功\n")

	// 转换响应
	result := &URLInspectionResult{}
	if response.InspectionResult != nil {
		fmt.Printf("[GOOGLE-CLIENT] 解析检查结果...\n")

		if response.InspectionResult.IndexStatusResult != nil {
			result.IndexStatusResult.IndexingState = string(response.InspectionResult.IndexStatusResult.IndexingState)
			if response.InspectionResult.IndexStatusResult.LastCrawlTime != "" {
				result.IndexStatusResult.LastCrawled = response.InspectionResult.IndexStatusResult.LastCrawlTime
				fmt.Printf("[GOOGLE-CLIENT] 最后抓取时间: %s\n", result.IndexStatusResult.LastCrawled)
			}
			fmt.Printf("[GOOGLE-CLIENT] 索引状态: %s\n", result.IndexStatusResult.IndexingState)
		}

		if response.InspectionResult.MobileUsabilityResult != nil {
			result.MobileUsabilityResult.MobileFriendly = response.InspectionResult.MobileUsabilityResult.Verdict == "MOBILE_USABILITY_VERdict_PASS"
			fmt.Printf("[GOOGLE-CLIENT] 移动友好性: %t\n", result.MobileUsabilityResult.MobileFriendly)
		}

		if response.InspectionResult.RichResultsResult != nil && response.InspectionResult.RichResultsResult.Verdict != "RICH_RESULTS_VERdict_PASS" {
			// 如果有富媒体结果检查信息
			result.RichResultsResult.Detected.Items = append(result.RichResultsResult.Detected.Items, struct {
				RichResultType string `json:"richResultType"`
			}{
				RichResultType: "UNKNOWN",
			})
		}
	} else {
		fmt.Printf("[GOOGLE-CLIENT] 警告: 响应中没有检查结果\n")
	}

	fmt.Printf("[GOOGLE-CLIENT] URL检查完成: %s\n", url)
	return result, nil
}

// SubmitSitemap 提交网站地图
func (c *Client) SubmitSitemap(sitemapURL string) error {
	fmt.Printf("[GOOGLE-CLIENT] 提交网站地图: %s\n", sitemapURL)
	fmt.Printf("[GOOGLE-CLIENT] 目标站点: %s\n", c.SiteURL)

	// 验证站点URL格式
	siteURL := c.SiteURL
	if !strings.HasSuffix(siteURL, "/") {
		siteURL = siteURL + "/"
	}
	if !strings.HasPrefix(siteURL, "sc-domain:") && !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}

	fmt.Printf("[GOOGLE-CLIENT] 格式化后的站点URL: %s\n", siteURL)

	call := c.service.Sitemaps.Submit(siteURL, sitemapURL)
	fmt.Printf("[GOOGLE-CLIENT] 发送网站地图提交请求...\n")
	err := call.Do()

	if err != nil {
		fmt.Printf("[GOOGLE-CLIENT] 网站地图提交失败: %v\n", err)

		// 尝试解析Google API错误
		if googleErr, ok := err.(*googleapi.Error); ok {
			fmt.Printf("[GOOGLE-CLIENT] Google API错误详情: 代码=%d, 消息=%s\n", googleErr.Code, googleErr.Message)
			for _, e := range googleErr.Errors {
				fmt.Printf("[GOOGLE-CLIENT] 错误项: 原因=%s, 消息=%s\n", e.Reason, e.Message)
			}

			// 根据错误代码提供更具体的错误信息
			switch googleErr.Code {
			case 403:
				return fmt.Errorf("网站地图提交权限不足: 请确保服务账号已被授予该站点的完整权限。详情: %s", googleErr.Message)
			case 404:
				return fmt.Errorf("站点或网站地图未找到: 请确保站点URL正确且网站地图可访问。详情: %s", googleErr.Message)
			case 400:
				if strings.Contains(googleErr.Message, "sitemap") {
					return fmt.Errorf("网站地图格式错误或无法访问: %s", googleErr.Message)
				}
				return fmt.Errorf("请求参数错误: %s", googleErr.Message)
			case 429:
				return fmt.Errorf("请求过于频繁，请稍后重试: %s", googleErr.Message)
			default:
				return fmt.Errorf("Google API错误 (代码: %d): %s", googleErr.Code, googleErr.Message)
			}
		}

		return fmt.Errorf("提交网站地图失败: %v", err)
	}

	fmt.Printf("[GOOGLE-CLIENT] 网站地图提交成功: %s\n", sitemapURL)
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

// PublishURL 提交URL到Google索引
func (c *Client) PublishURL(url string, urlType string) error {
	fmt.Printf("[GOOGLE-CLIENT] 提交URL到Google索引: %s (类型: %s)\n", url, urlType)
	fmt.Printf("[GOOGLE-CLIENT] 目标站点: %s\n", c.SiteURL)

	if c.indexingService == nil {
		return fmt.Errorf("Indexing服务未初始化")
	}

	// 验证URL格式
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL格式错误，必须以http://或https://开头: %s", url)
	}

	// 验证URL类型
	if urlType != "URL_UPDATED" && urlType != "URL_DELETED" {
		return fmt.Errorf("URL类型必须是 'URL_UPDATED' 或 'URL_DELETED'，当前: %s", urlType)
	}

	// 检查URL和站点域名的匹配关系（放宽验证，只要域名相同即可）
	siteDomain := extractDomain(c.SiteURL)
	urlDomain := extractDomain(url)

	if siteDomain != urlDomain {
		fmt.Printf("[GOOGLE-CLIENT] 警告: URL域名 (%s) 与站点域名 (%s) 不匹配\n", urlDomain, siteDomain)
		// 不返回错误，只是警告，因为有时子域名也可以提交
	}

	// 创建发布请求
	publication := &indexing.UrlNotification{
		Url:        url,
		Type:       urlType, // "URL_UPDATED" 或 "URL_DELETED"
	}

	fmt.Printf("[GOOGLE-CLIENT] 发送URL索引提交请求...\n")
	call := c.indexingService.UrlNotifications.Publish(publication)
	response, err := call.Do()

	if err != nil {
		fmt.Printf("[GOOGLE-CLIENT] URL索引提交失败: %v\n", err)

		// 尝试解析Google API错误
		if googleErr, ok := err.(*googleapi.Error); ok {
			fmt.Printf("[GOOGLE-CLIENT] Google API错误详情: 代码=%d, 消息=%s\n", googleErr.Code, googleErr.Message)
			for _, e := range googleErr.Errors {
				fmt.Printf("[GOOGLE-CLIENT] 错误项: 原因=%s, 消息=%s\n", e.Reason, e.Message)
			}

			// 根据错误代码提供更具体的错误信息
			switch googleErr.Code {
			case 403:
				if strings.Contains(googleErr.Message, "Indexing API") {
					return fmt.Errorf("Indexing API权限不足: 请确保服务账号已获得Indexing API访问权限。这需要在Google Cloud Console中启用Indexing API并授权服务账号。详情: %s", googleErr.Message)
				}
				return fmt.Errorf("URL索引提交权限不足: 请确保服务账号已被授予该站点的索引提交权限。详情: %s", googleErr.Message)
			case 429:
				return fmt.Errorf("URL索引提交请求过于频繁，Indexing API有严格的速率限制。请稍后重试。详情: %s", googleErr.Message)
			case 400:
				if strings.Contains(googleErr.Message, "url") {
					return fmt.Errorf("URL格式错误或无法访问: %s", googleErr.Message)
				}
				return fmt.Errorf("请求参数错误: %s", googleErr.Message)
			case 404:
				return fmt.Errorf("URL未找到或站点未验证: 请确保URL可访问且站点已在Google Search Console中验证。详情: %s", googleErr.Message)
			default:
				return fmt.Errorf("Google Indexing API错误 (代码: %d): %s", googleErr.Code, googleErr.Message)
			}
		}

		return fmt.Errorf("提交URL到索引失败: %v", err)
	}

	if response != nil {
		fmt.Printf("[GOOGLE-CLIENT] URL索引提交响应: %+v\n", response)
		if response.UrlNotificationMetadata != nil {
			fmt.Printf("[GOOGLE-CLIENT] 提交状态: %s\n", response.UrlNotificationMetadata.LatestUpdate.Type)
			if response.UrlNotificationMetadata.LatestUpdate.Url != "" {
				fmt.Printf("[GOOGLE-CLIENT] 提交的URL: %s\n", response.UrlNotificationMetadata.LatestUpdate.Url)
			}
		}
	}

	fmt.Printf("[GOOGLE-CLIENT] URL索引提交成功: %s\n", url)
	return nil
}

// extractDomain 从URL中提取域名
func extractDomain(url string) string {
	// 移除协议
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}

	// 移除路径
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	// 移除sc-domain:前缀
	if strings.HasPrefix(url, "sc-domain:") {
		url = url[11:]
	}

	return url
}

// BatchPublishURLs 批量提交URL到Google索引
func (c *Client) BatchPublishURLs(urls []string, urlType string, callback func(url string, success bool, err error)) {
	semaphore := make(chan struct{}, 3) // Indexing API限制更严格，降低并发数

	for _, url := range urls {
		go func(u string) {
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			err := c.PublishURL(u, urlType)
			callback(u, err == nil, err)

			// Indexing API有更严格的频率限制
			time.Sleep(1 * time.Second)
		}(url)
	}

	// 等待所有goroutine完成
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}
}

// GetURLNotificationStatus 获取URL通知状态
func (c *Client) GetURLNotificationStatus(url string) (*indexing.UrlNotificationMetadata, error) {
	if c.indexingService == nil {
		return nil, fmt.Errorf("Indexing服务未初始化")
	}

	call := c.indexingService.UrlNotifications.GetMetadata()
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("获取URL通知状态失败: %v", err)
	}

	return response, nil
}