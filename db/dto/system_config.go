package dto

// SystemConfigRequest 系统配置请求
type SystemConfigRequest struct {
	// SEO 配置
	SiteTitle       *string `json:"site_title,omitempty"`
	SiteDescription *string `json:"site_description,omitempty"`
	Keywords        *string `json:"keywords,omitempty"`
	Author          *string `json:"author,omitempty"`
	Copyright       *string `json:"copyright,omitempty"`
	SiteLogo        *string `json:"site_logo,omitempty"`

	// 自动处理配置组
	AutoProcessReadyResources *bool `json:"auto_process_ready_resources,omitempty"` // 自动处理待处理资源
	AutoProcessInterval       *int  `json:"auto_process_interval,omitempty"`        // 自动处理间隔（分钟）
	AutoTransferEnabled       *bool `json:"auto_transfer_enabled,omitempty"`        // 开启自动转存
	AutoTransferLimitDays     *int  `json:"auto_transfer_limit_days,omitempty"`     // 自动转存限制天数（0表示不限制）
	AutoTransferMinSpace      *int  `json:"auto_transfer_min_space,omitempty"`      // 最小存储空间（GB）
	AutoFetchHotDramaEnabled  *bool `json:"auto_fetch_hot_drama_enabled,omitempty"` // 自动拉取热播剧名字

	// API配置
	ApiToken *string `json:"api_token,omitempty"` // 公开API访问令牌

	// 违禁词配置
	ForbiddenWords *string `json:"forbidden_words,omitempty"` // 违禁词列表，用逗号分隔

	// 广告配置
	AdKeywords   *string `json:"ad_keywords,omitempty"`    // 广告关键词列表，用逗号分隔
	AutoInsertAd *string `json:"auto_insert_ad,omitempty"` // 自动插入广告内容

	// 其他配置
	PageSize        *int  `json:"page_size,omitempty"`
	MaintenanceMode *bool `json:"maintenance_mode,omitempty"`
	EnableRegister  *bool `json:"enable_register,omitempty"` // 开启注册功能

	// 三方统计配置
	ThirdPartyStatsCode *string `json:"third_party_stats_code,omitempty"` // 三方统计代码

	// Meilisearch配置
	MeilisearchEnabled   *bool   `json:"meilisearch_enabled,omitempty"`
	MeilisearchHost      *string `json:"meilisearch_host,omitempty"`
	MeilisearchPort      *string `json:"meilisearch_port,omitempty"`
	MeilisearchMasterKey *string `json:"meilisearch_master_key,omitempty"`
	MeilisearchIndexName *string `json:"meilisearch_index_name,omitempty"`
}

// SystemConfigResponse 系统配置响应
type SystemConfigResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// SEO 配置
	SiteTitle       string `json:"site_title"`
	SiteDescription string `json:"site_description"`
	Keywords        string `json:"keywords"`
	Author          string `json:"author"`
	Copyright       string `json:"copyright"`
	SiteLogo        string `json:"site_logo"`

	// 自动处理配置组
	AutoProcessReadyResources bool `json:"auto_process_ready_resources"` // 自动处理待处理资源
	AutoProcessInterval       int  `json:"auto_process_interval"`        // 自动处理间隔（分钟）
	AutoTransferEnabled       bool `json:"auto_transfer_enabled"`        // 开启自动转存
	AutoTransferLimitDays     int  `json:"auto_transfer_limit_days"`     // 自动转存限制天数（0表示不限制）
	AutoTransferMinSpace      int  `json:"auto_transfer_min_space"`      // 最小存储空间（GB）
	AutoFetchHotDramaEnabled  bool `json:"auto_fetch_hot_drama_enabled"` // 自动拉取热播剧名字

	// API配置
	ApiToken string `json:"api_token"` // 公开API访问令牌

	// 违禁词配置
	ForbiddenWords string `json:"forbidden_words"` // 违禁词列表，用逗号分隔

	// 广告配置
	AdKeywords   string `json:"ad_keywords"`    // 广告关键词列表，用逗号分隔
	AutoInsertAd string `json:"auto_insert_ad"` // 自动插入广告内容

	// 其他配置
	PageSize        int  `json:"page_size"`
	MaintenanceMode bool `json:"maintenance_mode"`
	EnableRegister  bool `json:"enable_register"` // 开启注册功能

	// 三方统计配置
	ThirdPartyStatsCode string `json:"third_party_stats_code"` // 三方统计代码

	// Meilisearch配置
	MeilisearchEnabled   bool   `json:"meilisearch_enabled"`
	MeilisearchHost      string `json:"meilisearch_host"`
	MeilisearchPort      string `json:"meilisearch_port"`
	MeilisearchMasterKey string `json:"meilisearch_master_key"`
	MeilisearchIndexName string `json:"meilisearch_index_name"`
}

// SystemConfigItem 单个配置项
type SystemConfigItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// SystemConfigListResponse 配置列表响应
type SystemConfigListResponse struct {
	Configs []SystemConfigItem `json:"configs"`
}
