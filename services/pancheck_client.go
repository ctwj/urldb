package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ctwj/urldb/utils"

	"github.com/go-resty/resty/v2"
)

// PanCheck 平台常量集合（PanCheck 服务端常量，固定发送全集，不依赖按平台过滤）
var pancheckAllPlatforms = []string{
	"quark", "uc", "baidu", "tianyi", "pan123", "pan115", "aliyun", "xunlei", "cmcc",
}

// LinkCheckStatus 链接检测结果状态（二态：有效/失效）
type LinkCheckStatus string

const (
	StatusValid    LinkCheckStatus = "valid"
	StatusInvalid  LinkCheckStatus = "invalid"
)

// LinkCheckOutcome 单个 URL 的检测结论
type LinkCheckOutcome struct {
	NormalizedURL string
	Status        LinkCheckStatus
	Platform      string
	FailReason    string
}

// pancheckRequest PanCheck 服务请求体
type pancheckRequest struct {
	Links             []string `json:"links"`
	SelectedPlatforms []string `json:"selected_platforms"`
}

// PanCheckClient PanCheck HTTP 客户端
type PanCheckClient struct {
	client *resty.Client
}

// NewPanCheckClient 创建 PanCheck 客户端
func NewPanCheckClient() *PanCheckClient {
	return &PanCheckClient{
		client: resty.New().
			SetTimeout(60 * time.Second).
			SetRetryCount(0),
	}
}

// Check 调用 PanCheck 服务检测一批 URL。
// 返回的 map 以规范化 URL 为键，仅包含得出有效/失效结论的链接；
// 未得出结论（pending/既不在 valid 也不在 invalid）的链接不出现在 map 中。
// 当发生网络超时、非 2xx 响应、JSON 解析失败时返回 error，调用方应将整批视为未得出结论。
func (c *PanCheckClient) Check(ctx context.Context, host string, timeoutSeconds int, urls []string) (map[string]LinkCheckOutcome, error) {
	if len(urls) == 0 {
		return map[string]LinkCheckOutcome{}, nil
	}
	if host == "" {
		return nil, fmt.Errorf("pancheck host 未配置")
	}

	endpoint := strings.TrimRight(host, "/") + "/api/v1/links/check"
	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	// 用 context 控制单次请求超时（resty Request 无 SetTimeout）
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 规范化所有请求链接，建立 规范化URL -> 原始URL 映射（PanCheck 侧也做规范化，这里用于匹配）
	normalizedSet := make(map[string]string, len(urls))
	requestLinks := make([]string, 0, len(urls))
	for _, raw := range urls {
		n := NormalizeURL(raw)
		if n == "" {
			continue
		}
		if _, ok := normalizedSet[n]; !ok {
			normalizedSet[n] = raw
			requestLinks = append(requestLinks, n)
		}
	}
	if len(requestLinks) == 0 {
		return map[string]LinkCheckOutcome{}, nil
	}

	reqBody := pancheckRequest{
		Links:             requestLinks,
		SelectedPlatforms: pancheckAllPlatforms,
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(endpoint)
	if err != nil {
		utils.Warn("PanCheck 请求失败: %v", err)
		return nil, fmt.Errorf("pancheck 请求失败: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		utils.Warn("PanCheck 返回非 2xx 状态码: %d", resp.StatusCode())
		return nil, fmt.Errorf("pancheck 返回状态码 %d", resp.StatusCode())
	}

	outcomes, err := parsePanCheckResponse(resp.Body(), requestLinks)
	if err != nil {
		utils.Warn("PanCheck 响应解析失败: %v", err)
		return nil, err
	}
	return outcomes, nil
}

// parsePanCheckResponse 容错解析 PanCheck 响应。
// - 同时兼容字符串数组与对象数组形式 [{url,platform,reason}]
// - 字段名容错（valid/invalid/pending 各组）
// - 按 normalized URL 匹配；invalid 优先于 valid
// - 既不在 valid 也不在 invalid 的链接视为未得出结论（不写入 map）
func parsePanCheckResponse(body []byte, requestedNormalized []string) (map[string]LinkCheckOutcome, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("响应非合法 JSON: %w", err)
	}

	validLinks := extractLinkGroup(raw, []string{"valid_links", "available", "ok", "valid"})
	invalidLinks := extractLinkGroup(raw, []string{"invalid_links", "unavailable", "expired", "dead", "invalid"})

	requested := make(map[string]struct{}, len(requestedNormalized))
	for _, n := range requestedNormalized {
		requested[n] = struct{}{}
	}

	outcomes := make(map[string]LinkCheckOutcome, len(requestedNormalized))

	// invalid 优先于 valid：先收集 valid，再用 invalid 覆盖
	for n, info := range validLinks {
		if _, ok := requested[n]; !ok {
			// 服务端返回的链接不在请求集合中，按规范化匹配失败则跳过
			continue
		}
		outcomes[n] = LinkCheckOutcome{
			NormalizedURL: n,
			Status:        StatusValid,
			Platform:      info.platform,
			FailReason:    "",
		}
	}
	for n, info := range invalidLinks {
		if _, ok := requested[n]; !ok {
			continue
		}
		reason := info.reason
		if reason == "" {
			reason = "链接已失效"
		}
		outcomes[n] = LinkCheckOutcome{
			NormalizedURL: n,
			Status:        StatusInvalid,
			Platform:      info.platform,
			FailReason:    reason,
		}
	}

	return outcomes, nil
}

// linkInfo 从对象数组项中提取的附加信息
type linkInfo struct {
	platform string
	reason   string
}

// extractLinkGroup 从响应原始 map 中按候选字段名提取一个链接分组，返回 规范化URL -> 信息。
func extractLinkGroup(raw map[string]json.RawMessage, candidates []string) map[string]linkInfo {
	result := make(map[string]linkInfo)
	for _, key := range candidates {
		chunk, ok := raw[key]
		if !ok {
			continue
		}
		// 尝试解析为字符串数组
		var asStrings []string
		if err := json.Unmarshal(chunk, &asStrings); err == nil {
			for _, s := range asStrings {
				n := NormalizeURL(s)
				if n == "" {
					n = strings.TrimSpace(s)
				}
				if n != "" {
					if _, exists := result[n]; !exists {
						result[n] = linkInfo{}
					}
				}
			}
			return result
		}
		// 尝试解析为对象数组
		var asObjects []map[string]any
		if err := json.Unmarshal(chunk, &asObjects); err == nil {
			for _, obj := range asObjects {
				rawURL := readStringField(obj, []string{"url", "link", "normalized_url"})
				if rawURL == "" {
					continue
				}
				n := NormalizeURL(rawURL)
				if n == "" {
					n = strings.TrimSpace(rawURL)
				}
				if n == "" {
					continue
				}
				if _, exists := result[n]; !exists {
					result[n] = linkInfo{
						platform: readStringField(obj, []string{"platform", "type"}),
						reason:   readStringField(obj, []string{"reason", "fail_reason", "message", "error"}),
					}
				}
			}
			return result
		}
	}
	return result
}

// readStringField 从对象中按候选键读取首个非空字符串值
func readStringField(obj map[string]any, candidates []string) string {
	for _, key := range candidates {
		if v, ok := obj[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// NormalizeURL 规范化 URL（对齐 tg_tool normalize_url）：
// 1. TrimSpace
// 2. 去掉 fragment（#...）
// 3. scheme + host 转小写（保留 path/query 原始大小写）
// 4. 去掉末尾 /
func NormalizeURL(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	// 去掉 fragment
	if idx := strings.Index(s, "#"); idx >= 0 {
		s = s[:idx]
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		// 非标准 URL：退化为 trim 后的字符串去尾斜杠
		return strings.TrimRight(s, "/")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	// 去掉末尾斜杠（仅对 path 处理）
	u.Path = strings.TrimRight(u.Path, "/")
	// 保留 query 原始大小写
	return u.String()
}

// URLHash 计算规范化 URL 的 SHA-256 十六进制摘要，作为缓存键
func URLHash(normalized string) string {
	h := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(h[:])
}
