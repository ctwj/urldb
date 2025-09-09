package entity

import (
	"time"
)

// HotDrama 热播剧实体
type HotDrama struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 基本信息
	Title        string  `json:"title" gorm:"size:255;not null"` // 剧名
	CardSubtitle string  `json:"card_subtitle" gorm:"size:500"`  // 副标题
	EpisodesInfo string  `json:"episodes_info" gorm:"size:100"`  // 集数信息
	IsNew        bool    `json:"is_new" gorm:"default:false"`    // 是否新剧
	Rating       float64 `json:"rating" gorm:"default:0"`        // 评分
	RatingCount  int     `json:"rating_count" gorm:"default:0"`  // 评分人数
	Year         string  `json:"year" gorm:"size:10"`            // 年份
	Region       string  `json:"region" gorm:"size:100"`         // 地区
	Genres       string  `json:"genres" gorm:"size:500"`         // 类型（多个用逗号分隔）
	Directors    string  `json:"directors" gorm:"size:500"`      // 导演（多个用逗号分隔）
	Actors       string  `json:"actors" gorm:"size:1000"`        // 演员（多个用逗号分隔）
	PosterURL    string  `json:"poster_url" gorm:"size:500"`     // 海报URL

	// 分类信息
	Category string `json:"category" gorm:"size:50"` // 分类（电影/电视剧）
	SubType  string `json:"sub_type" gorm:"size:50"` // 子类型（华语/欧美/韩国/日本等）
	Rank     int    `json:"rank" gorm:"default:0"`   // 排序（豆瓣返回顺序）

	// 数据来源
	Source    string `json:"source" gorm:"size:50;default:'douban'"` // 数据来源
	DoubanID  string `json:"douban_id" gorm:"size:50"`               // 豆瓣ID
	DoubanURI string `json:"douban_uri" gorm:"size:200"`             // 豆瓣链接
}

// TableName 指定表名
func (HotDrama) TableName() string {
	return "hot_dramas"
}
