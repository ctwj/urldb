package repo

import (
	"gorm.io/gorm"
)

// PaginationResult 分页查询结果
type PaginationResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// PaginationOptions 分页查询选项
type PaginationOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	OrderBy  string                 `json:"order_by"`
	OrderDir string                 `json:"order_dir"` // asc or desc
	Preloads []string               `json:"preloads"`  // 需要预加载的关联
	Filters  map[string]interface{} `json:"filters"`   // 过滤条件
}

// DefaultPaginationOptions 默认分页选项
func DefaultPaginationOptions() *PaginationOptions {
	return &PaginationOptions{
		Page:     1,
		PageSize: 20,
		OrderBy:  "id",
		OrderDir: "desc",
		Preloads: []string{},
		Filters:  make(map[string]interface{}),
	}
}

// PaginatedQuery 通用分页查询函数
func PaginatedQuery[T any](db *gorm.DB, options *PaginationOptions) (*PaginationResult[T], error) {
	// 验证分页参数
	if options.Page < 1 {
		options.Page = 1
	}
	if options.PageSize < 1 || options.PageSize > 1000 {
		options.PageSize = 20
	}

	// 应用预加载
	query := db.Model(new(T))
	for _, preload := range options.Preloads {
		query = query.Preload(preload)
	}

	// 应用过滤条件
	for key, value := range options.Filters {
		// 处理特殊过滤条件
		switch key {
		case "search":
			// 搜索条件需要特殊处理
			if searchStr, ok := value.(string); ok && searchStr != "" {
				query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+searchStr+"%", "%"+searchStr+"%")
			}
		case "category_id":
			if categoryID, ok := value.(uint); ok {
				query = query.Where("category_id = ?", categoryID)
			}
		case "pan_id":
			if panID, ok := value.(uint); ok {
				query = query.Where("pan_id = ?", panID)
			}
		case "is_valid":
			if isValid, ok := value.(bool); ok {
				query = query.Where("is_valid = ?", isValid)
			}
		case "is_public":
			if isPublic, ok := value.(bool); ok {
				query = query.Where("is_public = ?", isPublic)
			}
		default:
			// 通用过滤条件
			query = query.Where(key+" = ?", value)
		}
	}

	// 应用排序
	orderClause := options.OrderBy + " " + options.OrderDir
	query = query.Order(orderClause)

	// 计算偏移量
	offset := (options.Page - 1) * options.PageSize

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询数据
	var data []T
	if err := query.Offset(offset).Limit(options.PageSize).Find(&data).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int((total + int64(options.PageSize) - 1) / int64(options.PageSize))

	return &PaginationResult[T]{
		Data:       data,
		Total:      total,
		Page:       options.Page,
		PageSize:   options.PageSize,
		TotalPages: totalPages,
	}, nil
}