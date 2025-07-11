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
	AutoProcessReadyResources bool `json:"auto_process_ready_resources"`                    // 自动处理待处理资源
	AutoProcessInterval       int  `json:"auto_process_interval" validate:"min=1,max=1440"` // 自动处理间隔（分钟）
	AutoTransferEnabled       bool `json:"auto_transfer_enabled"`                           // 开启自动转存
	AutoFetchHotDramaEnabled  bool `json:"auto_fetch_hot_drama_enabled"`                    // 自动拉取热播剧名字

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
	AutoFetchHotDramaEnabled  bool `json:"auto_fetch_hot_drama_enabled"` // 自动拉取热播剧名字

	// 其他配置
	PageSize        int  `json:"page_size"`
	MaintenanceMode bool `json:"maintenance_mode"`
}
