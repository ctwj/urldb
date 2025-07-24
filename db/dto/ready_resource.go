package dto

// ReadyResourceRequest 待处理资源请求
type ReadyResourceRequest struct {
	Title       string   `json:"title" validate:"required" example:"示例资源标题"`
	Description string   `json:"description" example:"这是一个示例资源描述"`
	Url         []string `json:"url" validate:"required" example:"https://example.com/resource"`
	Category    string   `json:"category" example:"示例分类"`
	Tags        string   `json:"tags" example:"标签1,标签2"`
	Img         string   `json:"img" example:"https://example.com/image.jpg"`
	Source      string   `json:"source" example:"数据来源"`
	Extra       string   `json:"extra" example:"额外信息"`
}

// BatchReadyResourceRequest 批量待处理资源请求
type BatchReadyResourceRequest struct {
	Resources []ReadyResourceRequest `json:"resources" validate:"required"`
}
