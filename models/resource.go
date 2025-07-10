package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Tags 自定义类型，用于处理 PostgreSQL 数组
type Tags []string

// Value 实现 driver.Valuer 接口
func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return pq.Array(t), nil
}

// Scan 实现 sql.Scanner 接口
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// 处理 PostgreSQL 数组格式: {tag1,tag2,tag3}
		if len(v) == 0 {
			*t = Tags{}
			return nil
		}

		// 移除 { } 并分割
		s := string(v)
		s = strings.Trim(s, "{}")
		if s == "" {
			*t = Tags{}
			return nil
		}

		tags := strings.Split(s, ",")
		*t = Tags(tags)
		return nil

	case string:
		// 处理字符串格式
		if v == "" {
			*t = Tags{}
			return nil
		}
		s := strings.Trim(v, "{}")
		if s == "" {
			*t = Tags{}
			return nil
		}
		tags := strings.Split(s, ",")
		*t = Tags(tags)
		return nil

	default:
		return fmt.Errorf("cannot scan %T into Tags", value)
	}
}

// Resource 资源模型 - 更新后的结构
type Resource struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`       // 单个URL字符串
	PanID        *int      `json:"pan_id"`    // 平台ID，标识链接类型
	QuarkURL     string    `json:"quark_url"` // 新增字段
	FileSize     string    `json:"file_size"` // 改为 string 类型
	CategoryID   *int      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	ViewCount    int       `json:"view_count"`
	IsValid      bool      `json:"is_valid"` // 新增字段
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Tag 标签模型
type Tag struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ResourceTag 资源标签关联表
type ResourceTag struct {
	ID         int       `json:"id"`
	ResourceID int       `json:"resource_id"`
	TagID      int       `json:"tag_id"`
	CreatedAt  time.Time `json:"created_at"`
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
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PanID       *int   `json:"pan_id"` // 平台ID
	QuarkURL    string `json:"quark_url"`
	FileSize    string `json:"file_size"`
	CategoryID  *int   `json:"category_id"`
	IsValid     bool   `json:"is_valid"`
	IsPublic    bool   `json:"is_public"`
	TagIDs      []int  `json:"tag_ids"` // 标签ID列表
}

// UpdateResourceRequest 更新资源请求
type UpdateResourceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PanID       *int   `json:"pan_id"` // 平台ID
	QuarkURL    string `json:"quark_url"`
	FileSize    string `json:"file_size"`
	CategoryID  *int   `json:"category_id"`
	IsValid     bool   `json:"is_valid"`
	IsPublic    bool   `json:"is_public"`
	TagIDs      []int  `json:"tag_ids"` // 标签ID列表
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

// ReadyResource 待处理资源模型
type ReadyResource struct {
	ID         int       `json:"id"`
	Title      *string   `json:"title"`
	URL        string    `json:"url"`
	CreateTime time.Time `json:"create_time"`
	IP         *string   `json:"ip"`
}

// CreateReadyResourceRequest 创建待处理资源请求
type CreateReadyResourceRequest struct {
	Title *string `json:"title"`
	URL   string  `json:"url" binding:"required"`
	IP    *string `json:"ip"`
}

// BatchCreateReadyResourceRequest 批量创建待处理资源请求
type BatchCreateReadyResourceRequest struct {
	Resources []CreateReadyResourceRequest `json:"resources" binding:"required"`
}

// Stats 统计信息
type Stats struct {
	TotalResources  int `json:"total_resources"`
	TotalCategories int `json:"total_categories"`
	TotalDownloads  int `json:"total_downloads"`
	TotalViews      int `json:"total_views"`
}
