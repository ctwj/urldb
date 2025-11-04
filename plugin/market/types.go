package market

import (
	"time"
)

// PluginInfo 插件市场中的插件信息
type PluginInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Homepage    string    `json:"homepage"`
	Repository  string    `json:"repository"`
	License     string    `json:"license"`
	Tags        []string  `json:"tags"`
	Category    string    `json:"category"`
	Downloads   int64     `json:"downloads"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Size        int64     `json:"size"`
	DownloadURL string    `json:"download_url"`
	Checksum    string    `json:"checksum"`
	Dependencies []string `json:"dependencies"`
	Compatibility []string `json:"compatibility"` // 兼容的系统/架构
	Changelog   string    `json:"changelog"`
	Screenshots []string  `json:"screenshots"`
	Documentation string  `json:"documentation"`
}

// PluginCategory 插件分类
type PluginCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// PluginReview 插件评价
type PluginReview struct {
	ID         string    `json:"id"`
	PluginID   string    `json:"plugin_id"`
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	Rating     int       `json:"rating"` // 1-5星
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query       string   `json:"query"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	MinRating   float64  `json:"min_rating"`
	OrderBy     string   `json:"order_by"` // downloads, rating, updated_at
	OrderDir    string   `json:"order_dir"` // asc, desc
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Plugins    []PluginInfo `json:"plugins"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

// InstallRequest 安装请求
type InstallRequest struct {
	PluginID   string `json:"plugin_id"`
	Version    string `json:"version"` // 如果为空则安装最新版本
	Force      bool   `json:"force"`   // 是否强制安装
}

// InstallResponse 安装响应
type InstallResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	PluginName string `json:"plugin_name"`
	Version    string `json:"version"`
	Error      string `json:"error,omitempty"`
}

// MarketConfig 插件市场配置
type MarketConfig struct {
	APIEndpoint string   `json:"api_endpoint"`
	RegistryURL string   `json:"registry_url"`
	Timeout     int      `json:"timeout"` // 超时时间（秒）
	Proxy       string   `json:"proxy,omitempty"`
	Insecure    bool     `json:"insecure"` // 是否跳过SSL验证
	Headers     []Header `json:"headers,omitempty"`
}

// Header HTTP请求头
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}