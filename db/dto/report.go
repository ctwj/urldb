package dto

// ReportCreateRequest 举报创建请求
type ReportCreateRequest struct {
	ResourceKey string `json:"resource_key" validate:"required,max=255"`
	Reason      string `json:"reason" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=1000"`
	Contact     string `json:"contact" validate:"omitempty,max=255"`
	UserAgent   string `json:"user_agent" validate:"omitempty,max=1000"`
	IPAddress   string `json:"ip_address" validate:"omitempty,max=45"`
}

// ReportUpdateRequest 举报更新请求
type ReportUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=pending approved rejected"`
	Note   string `json:"note" validate:"omitempty,max=1000"`
}

// ResourceInfo 资源信息
type ResourceInfo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SaveURL     string `json:"save_url"`
	FileSize    string `json:"file_size"`
	Category    string `json:"category"`
	PanName     string `json:"pan_name"`
	ViewCount   int    `json:"view_count"`
	IsValid     bool   `json:"is_valid"`
	CreatedAt   string `json:"created_at"`
}

// ReportResponse 举报响应
type ReportResponse struct {
	ID          uint           `json:"id"`
	ResourceKey string         `json:"resource_key"`
	Reason      string         `json:"reason"`
	Description string         `json:"description"`
	Contact     string         `json:"contact"`
	UserAgent   string         `json:"user_agent"`
	IPAddress   string         `json:"ip_address"`
	Status      string         `json:"status"`
	Note        string         `json:"note"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	Resources   []ResourceInfo `json:"resources"` // 关联的资源列表
}

// ReportListRequest 举报列表请求
type ReportListRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Status   string `query:"status" validate:"omitempty,oneof=pending approved rejected"`
}