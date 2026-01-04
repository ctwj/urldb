package bing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client Bing Webmaster API客户端
type Client struct {
	siteURL string
	apiKey  string
	client  *http.Client
}

// Config Bing配置
type Config struct {
	SiteURL string `json:"site_url"`
	APIKey  string `json:"api_key"`
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

	// 标准化API密钥
	apiKey := strings.TrimSpace(config.APIKey)

	fmt.Printf("[BING-CLIENT] 初始化Bing客户端，目标站点: %s, API密钥: %s\n",
		siteURL, apiKey)

	return &Client{
		siteURL: siteURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
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

	// 检查API密钥是否配置
	if c.apiKey == "" {
		fmt.Printf("[BING-CLIENT] API密钥未配置，使用ping API作为备选方案\n")
		return c.submitSitemapWithPingAPI(sitemapURL)
	}

	// 先验证sitemap是否可访问
	if err := c.VerifySitemap(sitemapURL); err != nil {
		fmt.Printf("[BING-CLIENT] sitemap验证失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("sitemap验证失败: %v", err),
			StatusCode: 0,
		}, nil
	}

	// 使用Bing Webmaster API提交sitemap
	return c.submitSitemapWithWebmasterAPI(sitemapURL)
}

// submitSitemapWithWebmasterAPI 使用Bing Webmaster API提交sitemap
func (c *Client) submitSitemapWithWebmasterAPI(sitemapURL string) (*SitemapSubmitResponse, error) {
	fmt.Printf("[BING-CLIENT] 使用Bing Webmaster API提交sitemap\n")

	// 构建Bing Webmaster API URL
	apiURL := "https://ssl.bing.com/webmaster/api.svc/json/SubmitUrl"

	// 构建请求体
	requestBody := map[string]interface{}{
		"siteUrl":    c.siteURL,
		"url":        sitemapURL,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("[BING-CLIENT] JSON编码失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("JSON编码失败: %v", err),
			StatusCode: 0,
		}, nil
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("[BING-CLIENT] 创建请求失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("创建请求失败: %v", err),
			StatusCode: 0,
		}, nil
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", "URLDB-Bing-Webmaster-API/1.0")

	fmt.Printf("[BING-CLIENT] 发送请求到: %s\n", apiURL)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("[BING-CLIENT] 请求失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("网络请求失败: %v", err),
			StatusCode: 0,
		}, nil
	}
	defer resp.Body.Close()

	fmt.Printf("[BING-CLIENT] 响应状态码: %d\n", resp.StatusCode)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[BING-CLIENT] 读取响应失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("读取响应失败: %v", err),
			StatusCode: resp.StatusCode,
		}, nil
	}

	// 解析响应
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("[BING-CLIENT] JSON解析失败: %v, 响应内容: %s\n", err, string(body))
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("JSON解析失败: %v", err),
			StatusCode: resp.StatusCode,
		}, nil
	}

	// 分析API响应
	success := resp.StatusCode == 200
	message := c.getWebmasterAPIMessage(resp.StatusCode, apiResponse)

	response := &SitemapSubmitResponse{
		Success:   success,
		Message:   message,
		StatusCode: resp.StatusCode,
	}

	if success {
		fmt.Printf("[BING-CLIENT] sitemap提交成功: %s\n", sitemapURL)
	} else {
		fmt.Printf("[BING-CLIENT] sitemap提交失败: %s (状态码: %d, 响应: %s)\n",
			sitemapURL, resp.StatusCode, string(body))
	}

	return response, nil
}

// submitSitemapWithPingAPI 使用传统的ping API作为备选方案
func (c *Client) submitSitemapWithPingAPI(sitemapURL string) (*SitemapSubmitResponse, error) {
	fmt.Printf("[BING-CLIENT] 使用ping API作为备选方案\n")

	// 先验证sitemap是否可访问
	if err := c.VerifySitemap(sitemapURL); err != nil {
		fmt.Printf("[BING-CLIENT] sitemap验证失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("sitemap验证失败: %v", err),
			StatusCode: 0,
		}, nil
	}

	// 构建Bing ping URL
	pingURL := fmt.Sprintf("https://www.bing.com/webmaster/ping.aspx?siteMap=%s", url.QueryEscape(sitemapURL))

	fmt.Printf("[BING-CLIENT] 发送ping请求到: %s\n", pingURL)

	// 发送请求
	resp, err := c.client.Get(pingURL)
	if err != nil {
		fmt.Printf("[BING-CLIENT] ping请求失败: %v\n", err)
		return &SitemapSubmitResponse{
			Success:   false,
			Message:   fmt.Sprintf("网络请求失败: %v", err),
			StatusCode: 0,
		}, nil
	}
	defer resp.Body.Close()

	fmt.Printf("[BING-CLIENT] ping响应状态码: %d\n", resp.StatusCode)

	// 解析响应
	success := resp.StatusCode == 200 || resp.StatusCode == 410
	message := c.getStatusMessage(resp.StatusCode)

	if resp.StatusCode == 410 {
		message = "Bing ping API已废弃，建议配置API密钥使用Webmaster API"
	}

	response := &SitemapSubmitResponse{
		Success:   success,
		Message:   message,
		StatusCode: resp.StatusCode,
	}

	return response, nil
}

// getWebmasterAPIMessage 解析Webmaster API响应消息
func (c *Client) getWebmasterAPIMessage(statusCode int, response map[string]interface{}) string {
	if statusCode == 200 {
		if d, ok := response["d"].(map[string]interface{}); ok {
			if success, ok := d["success"].(bool); ok && success {
				return "sitemap提交成功"
			} else if msg, ok := d["message"].(string); ok {
				return msg
			}
		}
		return "提交成功"
	}

	// 处理错误响应
	if errorInfo, ok := response["error"].(map[string]interface{}); ok {
		if msg, ok := errorInfo["message"].(string); ok {
			return msg
		}
		if code, ok := errorInfo["code"].(string); ok {
			return fmt.Sprintf("API错误: %s", code)
		}
	}

	return c.getStatusMessage(statusCode)
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
	case 410:
		return "Bing ping API已废弃，但sitemap可正常访问"
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