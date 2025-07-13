package pan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BasePanService 基础网盘服务
type BasePanService struct {
	config     *PanConfig
	httpClient *http.Client
	headers    map[string]string
}

// NewBasePanService 创建基础网盘服务
func NewBasePanService(config *PanConfig) *BasePanService {
	return &BasePanService{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// SetHeader 设置请求头
func (b *BasePanService) SetHeader(key, value string) {
	b.headers[key] = value
}

// SetHeaders 批量设置请求头
func (b *BasePanService) SetHeaders(headers map[string]string) {
	for key, value := range headers {
		b.headers[key] = value
	}
}

// GetHeader 获取请求头
func (b *BasePanService) GetHeader(key string) string {
	return b.headers[key]
}

// GetConfig 获取配置
func (b *BasePanService) GetConfig() *PanConfig {
	return b.config
}

// HTTPGet 发送GET请求
func (b *BasePanService) HTTPGet(requestURL string, queryParams map[string]string) ([]byte, error) {
	// 构建查询参数
	if len(queryParams) > 0 {
		u, err := url.Parse(requestURL)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		requestURL = u.String()
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %d, %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HTTPPost 发送POST请求
func (b *BasePanService) HTTPPost(requestURL string, data interface{}, queryParams map[string]string) ([]byte, error) {
	var body io.Reader

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	// 构建查询参数
	if len(queryParams) > 0 {
		u, err := url.Parse(requestURL)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		requestURL = u.String()
	}

	req, err := http.NewRequest("POST", requestURL, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	// 如果没有设置Content-Type，默认设置为application/json
	if _, exists := b.headers["Content-Type"]; !exists && data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %d, %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// HTTPPut 发送PUT请求
func (b *BasePanService) HTTPPut(requestURL string, data interface{}) ([]byte, error) {
	var body io.Reader

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("PUT", requestURL, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	// 如果没有设置Content-Type，默认设置为application/json
	if _, exists := b.headers["Content-Type"]; !exists && data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %d, %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// HTTPDelete 发送DELETE请求
func (b *BasePanService) HTTPDelete(requestURL string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %d, %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ExecuteWithRetry 带重试的请求执行
func (b *BasePanService) ExecuteWithRetry(executeFunc func() ([]byte, error), maxRetries int, retryDelay time.Duration) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		data, err := executeFunc()
		if err == nil {
			return data, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("重试%d次后仍然失败: %v", maxRetries, lastErr)
}

// ParseJSONResponse 解析JSON响应
func (b *BasePanService) ParseJSONResponse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// IsValidURL 验证URL格式
func (b *BasePanService) IsValidURL(urlStr string) bool {
	_, err := url.Parse(urlStr)
	return err == nil
}

// ExtractFileName 从URL中提取文件名
func (b *BasePanService) ExtractFileName(urlStr string) string {
	parts := strings.Split(urlStr, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// SanitizeFileName 清理文件名
func (b *BasePanService) SanitizeFileName(fileName string) string {
	// 移除或替换不合法的字符
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := fileName

	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	return result
}
