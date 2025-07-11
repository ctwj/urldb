package dto

// HotDramaRequest 热播剧请求
type HotDramaRequest struct {
	Title     string  `json:"title" validate:"required"`
	Rating    float64 `json:"rating"`
	Year      string  `json:"year"`
	Directors string  `json:"directors"`
	Actors    string  `json:"actors"`
	Category  string  `json:"category"`
	SubType   string  `json:"sub_type"`
	Source    string  `json:"source"`
	DoubanID  string  `json:"douban_id"`
}

// HotDramaResponse 热播剧响应
type HotDramaResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Title     string  `json:"title"`
	Rating    float64 `json:"rating"`
	Year      string  `json:"year"`
	Directors string  `json:"directors"`
	Actors    string  `json:"actors"`
	Category  string  `json:"category"`
	SubType   string  `json:"sub_type"`
	Source    string  `json:"source"`
	DoubanID  string  `json:"douban_id"`
}

// HotDramaListResponse 热播剧列表响应
type HotDramaListResponse struct {
	Total int                `json:"total"`
	Items []HotDramaResponse `json:"items"`
}
