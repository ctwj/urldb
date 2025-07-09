package models

import (
	"time"
)

// Resource 资源模型
type Resource struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	FileType      string    `json:"file_type"`
	CategoryID    *int      `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	Tags          []string  `json:"tags"`
	DownloadCount int       `json:"download_count"`
	ViewCount     int       `json:"view_count"`
	IsPublic      bool      `json:"is_public"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Category 分类模型
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateResourceRequest 创建资源请求
type CreateResourceRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	FilePath    string   `json:"file_path"`
	FileSize    int64    `json:"file_size"`
	FileType    string   `json:"file_type"`
	CategoryID  *int     `json:"category_id"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

// UpdateResourceRequest 更新资源请求
type UpdateResourceRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	FilePath    string   `json:"file_path"`
	FileSize    int64    `json:"file_size"`
	FileType    string   `json:"file_type"`
	CategoryID  *int     `json:"category_id"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query      string `json:"query"`
	CategoryID *int   `json:"category_id"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Resources []Resource `json:"resources"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

// Stats 统计信息
type Stats struct {
	TotalResources  int `json:"total_resources"`
	TotalCategories int `json:"total_categories"`
	TotalDownloads  int `json:"total_downloads"`
	TotalViews      int `json:"total_views"`
}
