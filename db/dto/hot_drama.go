package dto

// HotDramaRequest 热播剧请求
type HotDramaRequest struct {
	Title        string  `json:"title" validate:"required"`
	CardSubtitle string  `json:"card_subtitle"`
	EpisodesInfo string  `json:"episodes_info"`
	IsNew        bool    `json:"is_new"`
	Rating       float64 `json:"rating"`
	RatingCount  int     `json:"rating_count"`
	Year         string  `json:"year"`
	Region       string  `json:"region"`
	Genres       string  `json:"genres"`
	Directors    string  `json:"directors"`
	Actors       string  `json:"actors"`
	PosterURL    string  `json:"poster_url"`
	Category     string  `json:"category"`
	SubType      string  `json:"sub_type"`
	Source       string  `json:"source"`
	DoubanID     string  `json:"douban_id"`
	DoubanURI    string  `json:"douban_uri"`
}

// HotDramaResponse 热播剧响应
type HotDramaResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Title        string  `json:"title"`
	CardSubtitle string  `json:"card_subtitle"`
	EpisodesInfo string  `json:"episodes_info"`
	IsNew        bool    `json:"is_new"`
	Rating       float64 `json:"rating"`
	RatingCount  int     `json:"rating_count"`
	Year         string  `json:"year"`
	Region       string  `json:"region"`
	Genres       string  `json:"genres"`
	Directors    string  `json:"directors"`
	Actors       string  `json:"actors"`
	PosterURL    string  `json:"poster_url"`
	Category     string  `json:"category"`
	SubType      string  `json:"sub_type"`
	Source       string  `json:"source"`
	DoubanID     string  `json:"douban_id"`
	DoubanURI    string  `json:"douban_uri"`
}

// HotDramaListResponse 热播剧列表响应
type HotDramaListResponse struct {
	Total int                `json:"total"`
	Items []HotDramaResponse `json:"items"`
}
