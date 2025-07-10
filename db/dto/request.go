package dto

// CreatePanRequest 创建平台请求
type CreatePanRequest struct {
	Name      string `json:"name" binding:"required"`
	Key       int    `json:"key"`
	Ck        string `json:"ck"`
	IsValid   bool   `json:"is_valid"`
	Space     int64  `json:"space"`
	LeftSpace int64  `json:"left_space"`
	Remark    string `json:"remark"`
}

// UpdatePanRequest 更新平台请求
type UpdatePanRequest struct {
	Name      string `json:"name"`
	Key       int    `json:"key"`
	Ck        string `json:"ck"`
	IsValid   bool   `json:"is_valid"`
	Space     int64  `json:"space"`
	LeftSpace int64  `json:"left_space"`
	Remark    string `json:"remark"`
}

// CreateCksRequest 创建cookie请求
type CreateCksRequest struct {
	PanID  uint   `json:"pan_id" binding:"required"`
	T      string `json:"t"`
	Idx    int    `json:"idx"`
	Ck     string `json:"ck"`
	Remark string `json:"remark"`
}

// UpdateCksRequest 更新cookie请求
type UpdateCksRequest struct {
	PanID  uint   `json:"pan_id"`
	T      string `json:"t"`
	Idx    int    `json:"idx"`
	Ck     string `json:"ck"`
	Remark string `json:"remark"`
}

// CreateResourceRequest 创建资源请求
type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PanID       *uint  `json:"pan_id"`
	QuarkURL    string `json:"quark_url"`
	FileSize    string `json:"file_size"`
	CategoryID  *uint  `json:"category_id"`
	IsValid     bool   `json:"is_valid"`
	IsPublic    bool   `json:"is_public"`
	TagIDs      []uint `json:"tag_ids"`
}

// UpdateResourceRequest 更新资源请求
type UpdateResourceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PanID       *uint  `json:"pan_id"`
	QuarkURL    string `json:"quark_url"`
	FileSize    string `json:"file_size"`
	CategoryID  *uint  `json:"category_id"`
	IsValid     bool   `json:"is_valid"`
	IsPublic    bool   `json:"is_public"`
	TagIDs      []uint `json:"tag_ids"`
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

// SearchRequest 搜索请求
type SearchRequest struct {
	Query      string `json:"query"`
	CategoryID *uint  `json:"category_id"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}
