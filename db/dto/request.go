package dto

// CreatePanRequest 创建平台请求
type CreatePanRequest struct {
	Name   string `json:"name" binding:"required"`
	Key    int    `json:"key"`
	Icon   string `json:"icon"`
	Remark string `json:"remark"`
}

// UpdatePanRequest 更新平台请求
type UpdatePanRequest struct {
	Name   string `json:"name"`
	Key    int    `json:"key"`
	Icon   string `json:"icon"`
	Remark string `json:"remark"`
}

// CreateCksRequest 创建cookie请求
type CreateCksRequest struct {
	PanID       uint   `json:"pan_id" binding:"required"`
	Idx         int    `json:"idx"`
	Ck          string `json:"ck"`
	IsValid     bool   `json:"is_valid"`
	Space       int64  `json:"space"`
	LeftSpace   int64  `json:"left_space"`
	UsedSpace   int64  `json:"used_space"`
	Username    string `json:"username"`
	VipStatus   bool   `json:"vip_status"`
	ServiceType string `json:"service_type"`
	Remark      string `json:"remark"`
}

// UpdateCksRequest 更新cookie请求
type UpdateCksRequest struct {
	PanID       uint   `json:"pan_id"`
	Idx         int    `json:"idx"`
	Ck          string `json:"ck"`
	IsValid     bool   `json:"is_valid"`
	Space       int64  `json:"space"`
	LeftSpace   int64  `json:"left_space"`
	UsedSpace   int64  `json:"used_space"`
	Username    string `json:"username"`
	VipStatus   bool   `json:"vip_status"`
	ServiceType string `json:"service_type"`
	Remark      string `json:"remark"`
}

// CreateResourceRequest 创建资源请求
type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PanID       *uint  `json:"pan_id"`
	SaveURL     string `json:"save_url"`
	FileSize    string `json:"file_size"`
	CategoryID  *uint  `json:"category_id"`
	IsValid     bool   `json:"is_valid"`
	IsPublic    bool   `json:"is_public"`
	TagIDs      []uint `json:"tag_ids"`
	Cover       string `json:"cover"`
	Author      string `json:"author"`
	ErrorMsg    string `json:"error_msg"`
}

// UpdateResourceRequest 更新资源请求
type UpdateResourceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PanID       *uint  `json:"pan_id"`
	SaveURL     string `json:"save_url"`
	FileSize    string `json:"file_size"`
	CategoryID  *uint  `json:"category_id"`
	IsValid     bool   `json:"is_valid"`
	IsPublic    bool   `json:"is_public"`
	TagIDs      []uint `json:"tag_ids"`
	Cover       string `json:"cover"`
	Author      string `json:"author"`
	ErrorMsg    string `json:"error_msg"`
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
	CategoryID  *uint  `json:"category_id"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  *uint  `json:"category_id"`
}

// CreateReadyResourceRequest 创建待处理资源请求
type CreateReadyResourceRequest struct {
	Title       *string `json:"title"`
	Description string  `json:"description"`
	URL         string  `json:"url" binding:"required"`
	Category    string  `json:"category"`
	Tags        string  `json:"tags"`
	Img         string  `json:"img"`
	Source      string  `json:"source"`
	Extra       string  `json:"extra"`
	IP          *string `json:"ip"`
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
