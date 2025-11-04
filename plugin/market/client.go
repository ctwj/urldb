package market

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ctwj/urldb/utils"
)

// MarketClient 插件市场客户端
type MarketClient struct {
	config    MarketConfig
	httpClient *http.Client
}

// NewMarketClient 创建新的市场客户端
func NewMarketClient(config MarketConfig) *MarketClient {
	// 创建HTTP客户端
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
	}

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Timeout) * time.Second,
	}

	// 设置默认超时时间
	if config.Timeout <= 0 {
		config.Timeout = 30
	}

	return &MarketClient{
		config:     config,
		httpClient: httpClient,
	}
}

// Search 搜索插件
func (c *MarketClient) Search(req SearchRequest) (*SearchResponse, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s/plugins/search", c.config.APIEndpoint)

	// 序列化请求参数
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %v", err)
	}

	// 发送POST请求
	resp, err := c.doRequest("POST", apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to send search request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	return &searchResp, nil
}

// GetPlugin 获取插件详细信息
func (c *MarketClient) GetPlugin(pluginID string) (*PluginInfo, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s/plugins/%s", c.config.APIEndpoint, pluginID)

	// 发送GET请求
	resp, err := c.doRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send get plugin request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get plugin failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var plugin PluginInfo
	if err := json.NewDecoder(resp.Body).Decode(&plugin); err != nil {
		return nil, fmt.Errorf("failed to decode plugin response: %v", err)
	}

	return &plugin, nil
}

// DownloadPlugin 下载插件
func (c *MarketClient) DownloadPlugin(pluginID, version string) ([]byte, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s/plugins/%s/download", c.config.APIEndpoint, pluginID)
	if version != "" {
		apiURL = fmt.Sprintf("%s?version=%s", apiURL, version)
	}

	// 发送GET请求
	resp, err := c.doRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send download request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read download response: %v", err)
	}

	return data, nil
}

// GetPluginReviews 获取插件评价
func (c *MarketClient) GetPluginReviews(pluginID string, page, pageSize int) ([]PluginReview, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s/plugins/%s/reviews?page=%d&page_size=%d",
		c.config.APIEndpoint, pluginID, page, pageSize)

	// 发送GET请求
	resp, err := c.doRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send get reviews request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get reviews failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var reviews []PluginReview
	if err := json.NewDecoder(resp.Body).Decode(&reviews); err != nil {
		return nil, fmt.Errorf("failed to decode reviews response: %v", err)
	}

	return reviews, nil
}

// SubmitReview 提交评价
func (c *MarketClient) SubmitReview(review PluginReview) error {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s/plugins/%s/reviews", c.config.APIEndpoint, review.PluginID)

	// 序列化请求参数
	data, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("failed to marshal review: %v", err)
	}

	// 发送POST请求
	resp, err := c.doRequest("POST", apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to send submit review request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("submit review failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetCategories 获取插件分类
func (c *MarketClient) GetCategories() ([]PluginCategory, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s/categories", c.config.APIEndpoint)

	// 发送GET请求
	resp, err := c.doRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send get categories request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get categories failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var categories []PluginCategory
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to decode categories response: %v", err)
	}

	return categories, nil
}

// doRequest 执行HTTP请求
func (c *MarketClient) doRequest(method, url string, body []byte) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	// 添加自定义请求头
	for _, header := range c.config.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	// 添加User-Agent
	req.Header.Set("User-Agent", "urlDB-Plugin-Market-Client/1.0")

	// 发送请求
	utils.Debug("Sending %s request to %s", method, url)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetAuthToken 设置认证令牌
func (c *MarketClient) SetAuthToken(token string) {
	c.config.Headers = append(c.config.Headers, Header{
		Key:   "Authorization",
		Value: "Bearer " + token,
	})
}

// SetAPIKey 设置API密钥
func (c *MarketClient) SetAPIKey(apiKey string) {
	c.config.Headers = append(c.config.Headers, Header{
		Key:   "X-API-Key",
		Value: apiKey,
	})
}