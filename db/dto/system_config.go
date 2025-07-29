package dto

// SystemConfigRequest 系统配置请求
type SystemConfigRequest struct {
	// SEO 配置
	SiteTitle       string `json:"site_title" validate:"required"`
	SiteDescription string `json:"site_description"`
	Keywords        string `json:"keywords"`
	Author          string `json:"author"`
	Copyright       string `json:"copyright"`

	// 自动处理配置组
	AutoProcessReadyResources bool `json:"auto_process_ready_resources"`                        // 自动处理待处理资源
	AutoProcessInterval       int  `json:"auto_process_interval" validate:"min=1,max=1440"`     // 自动处理间隔（分钟）
	AutoTransferEnabled       bool `json:"auto_transfer_enabled"`                               // 开启自动转存
	AutoTransferLimitDays     int  `json:"auto_transfer_limit_days" validate:"min=0,max=365"`   // 自动转存限制天数（0表示不限制）
	AutoTransferMinSpace      int  `json:"auto_transfer_min_space" validate:"min=100,max=1024"` // 最小存储空间（GB）
	AutoFetchHotDramaEnabled  bool `json:"auto_fetch_hot_drama_enabled"`                        // 自动拉取热播剧名字

	// API配置
	ApiToken string `json:"api_token"` // 公开API访问令牌

	// 其他配置
	PageSize        int  `json:"page_size" validate:"min=10,max=500"`
	MaintenanceMode bool `json:"maintenance_mode"`
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

	// 自动处理配置组
	AutoProcessReadyResources bool `json:"auto_process_ready_resources"` // 自动处理待处理资源
	AutoProcessInterval       int  `json:"auto_process_interval"`        // 自动处理间隔（分钟）
	AutoTransferEnabled       bool `json:"auto_transfer_enabled"`        // 开启自动转存
	AutoTransferLimitDays     int  `json:"auto_transfer_limit_days"`     // 自动转存限制天数（0表示不限制）
	AutoTransferMinSpace      int  `json:"auto_transfer_min_space"`      // 最小存储空间（GB）
	AutoFetchHotDramaEnabled  bool `json:"auto_fetch_hot_drama_enabled"` // 自动拉取热播剧名字

	// API配置
	ApiToken string `json:"api_token"` // 公开API访问令牌

	// 其他配置
	PageSize        int  `json:"page_size"`
	MaintenanceMode bool `json:"maintenance_mode"`
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
