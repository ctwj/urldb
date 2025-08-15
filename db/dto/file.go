package dto

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	IsPublic bool `json:"is_public" form:"is_public"` // 是否公开
}

// FileResponse 文件响应
type FileResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// 文件信息
	OriginalName string `json:"original_name"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
	MimeType     string `json:"mime_type"`

	// 访问信息
	AccessURL string `json:"access_url"`

	// 用户信息
	UserID uint   `json:"user_id"`
	User   string `json:"user"` // 用户名

	// 状态信息
	Status    string `json:"status"`
	IsPublic  bool   `json:"is_public"`
	IsDeleted bool   `json:"is_deleted"`
}

// FileListRequest 文件列表请求
type FileListRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Search   string `json:"search" form:"search"`
	FileType string `json:"file_type" form:"file_type"`
	Status   string `json:"status" form:"status"`
	UserID   uint   `json:"user_id" form:"user_id"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Files []FileResponse `json:"files"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	File    FileResponse `json:"file"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
}

// FileDeleteRequest 文件删除请求
type FileDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// FileUpdateRequest 文件更新请求
type FileUpdateRequest struct {
	ID       uint   `json:"id" binding:"required"`
	IsPublic *bool  `json:"is_public"`
	Status   string `json:"status"`
}
