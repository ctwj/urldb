package bing

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client Bing Webmaster API客户端
type Client struct {
	siteURL string
	client   *http.Client
}

// Config Bing配置
type Config struct {
	SiteURL string `json:"site_url"`
}

// SitemapSubmitResponse sitemap提交响应
type SitemapSubmitResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	StatusCode int   `json:"status_code"`
}


// NewClient 创建新的Bing客户端
func NewClient(config *Config) (*Client, error) {
	// 标准化站点URL
	siteURL := strings.TrimSpace(config.SiteURL)
	if siteURL != "" && !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}

	fmt.Printf("[BING-CLIENT] 初始化Bing客户端，目标站点: %s\n", siteURL)

	return &Client{
		siteURL: siteURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// SubmitSitemap 提交网站地图到Bing
func (c *Client) SubmitSitemap(sitemapURL string) (*SitemapSubmitResponse, error) {
	fmt.Printf("[BING-CLIENT] 提交网站地图到Bing: %s\n", sitemapURL)

	// 验证sitemap URL
	if !strings.HasPrefix(sitemapURL, "http://") && !strings.HasPrefix(sitemapURL, "https://") {
		return nil, fmt.Errorf("sitemap URL格式错误，必须以http://或https://开头")
	}

	// 构建Bing ping URL
	pingURL := fmt.Sprintf("https://www.bing.com/webmaster/ping.aspx?siteMap=%s", url.QueryEscape(sitemapURL))

	fmt.Printf("[BING-CLIENT] 发送请求到: %s\n", pingURL)

	// 发送请求
	resp, err := c.client.Get(pingURL)
	if err != nil {
		fmt.Printf("[BING-CLIENT] 请求失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("网络请求失败: %v", err),
			StatusCode: 0,
		}, err
	}
	defer resp.Body.Close()

	fmt.Printf("[BING-CLIENT] 响应状态码: %d\n", resp.StatusCode)

	// 解析响应
	success := resp.StatusCode == 200
	message := c.getStatusMessage(resp.StatusCode)

	response := &SitemapSubmitResponse{
		Success:   success,
		Message:   message,
		StatusCode: resp.StatusCode,
	}

	if success {
		fmt.Printf("[BING-CLIENT] 网站地图提交成功: %s\n", sitemapURL)
	} else {
		fmt.Printf("[BING-CLIENT] 网站地图提交失败: %s (状态码: %d)\n", sitemapURL, resp.StatusCode)
	}

	return response, nil
}

// getStatusMessage 根据状态码获取消息
func (c *Client) getStatusMessage(statusCode int) string {
	switch statusCode {
	case 200:
		return "提交成功"
	case 400:
		return "请求参数错误"
	case 404:
		return "网站地图未找到或无法访问"
	case 429:
		return "请求过于频繁，请稍后重试"
	case 500:
		return "Bing服务器内部错误"
	default:
		return fmt.Sprintf("未知错误 (状态码: %d)", statusCode)
	}
}


// BatchSubmitSitemaps 批量提交网站地图
func (c *Client) BatchSubmitSitemaps(sitemapURLs []string) ([]*SitemapSubmitResponse, error) {
	fmt.Printf("[BING-CLIENT] 批量提交 %d 个网站地图\n", len(sitemapURLs))

	responses := make([]*SitemapSubmitResponse, len(sitemapURLs))

	for i, sitemapURL := range sitemapURLs {
		response, err := c.SubmitSitemap(sitemapURL)
		if err != nil {
			response = &SitemapSubmitResponse{
				Success:   false,
				Message:   fmt.Sprintf("提交失败: %v", err),
				StatusCode: 0,
			}
		}
		responses[i] = response

		// Bing建议间隔1秒以上
		if i < len(sitemapURLs)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	successCount := 0
	for _, resp := range responses {
		if resp.Success {
			successCount++
		}
	}

	fmt.Printf("[BING-CLIENT] 批量提交完成: 成功 %d/%d\n", successCount, len(responses))
	return responses, nil
}

// VerifySitemap 验证网站地图可访问性
func (c *Client) VerifySitemap(sitemapURL string) error {
	fmt.Printf("[BING-CLIENT] 验证网站地图可访问性: %s\n", sitemapURL)

	resp, err := c.client.Get(sitemapURL)
	if err != nil {
		return fmt.Errorf("无法访问网站地图: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("网站地图返回错误状态码: %d", resp.StatusCode)
	}

	// 检查Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "xml") && !strings.Contains(contentType, "text/xml") {
		fmt.Printf("[BING-CLIENT] 警告: Content-Type不是XML格式: %s\n", contentType)
	}

	fmt.Printf("[BING-CLIENT] 网站地图验证成功: %s\n", sitemapURL)
	return nil
}