package dto

import "time"

// SearchResponse 搜索响应
type SearchResponse struct {
	Resources []ResourceResponse `json:"resources"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

// ResourceResponse 资源响应
type ResourceResponse struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	URL          string        `json:"url"`
	PanID        *uint         `json:"pan_id"`
	QuarkURL     string        `json:"quark_url"`
	FileSize     string        `json:"file_size"`
	CategoryID   *uint         `json:"category_id"`
	CategoryName string        `json:"category_name"`
	ViewCount    int           `json:"view_count"`
	IsValid      bool          `json:"is_valid"`
	IsPublic     bool          `json:"is_public"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Tags         []TagResponse `json:"tags"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ResourceCount int64    `json:"resource_count"`
	TagNames      []string `json:"tag_names"`
}

// TagResponse 标签响应
type TagResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CategoryID    *uint  `json:"category_id"`
	CategoryName  string `json:"category_name"`
	ResourceCount int64  `json:"resource_count"`
}

// PanResponse 平台响应
type PanResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Key    int    `json:"key"`
	Icon   string `json:"icon"`
	Remark string `json:"remark"`
}

// CksResponse Cookie响应
type CksResponse struct {
	ID        uint   `json:"id"`
	PanID     uint   `json:"pan_id"`
	Idx       int    `json:"idx"`
	Ck        string `json:"ck"`
	IsValid   bool   `json:"is_valid"`
	Space     int64  `json:"space"`
	LeftSpace int64  `json:"left_space"`
	Remark    string `json:"remark"`
}

// ReadyResourceResponse 待处理资源响应
type ReadyResourceResponse struct {
	ID          uint      `json:"id"`
	Title       *string   `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Category    string    `json:"category"`
	Tags        string    `json:"tags"`
	Img         string    `json:"img"`
	Source      string    `json:"source"`
	Extra       string    `json:"extra"`
	CreateTime  time.Time `json:"create_time"`
	IP          *string   `json:"ip"`
}

// Stats 统计信息
type Stats struct {
	TotalResources  int64 `json:"total_resources"`
	TotalCategories int64 `json:"total_categories"`
	TotalTags       int64 `json:"total_tags"`
	TotalViews      int64 `json:"total_views"`
}
