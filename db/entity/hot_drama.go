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
	Title     string  `json:"title" gorm:"size:255;not null"` // 剧名
	Rating    float64 `json:"rating" gorm:"default:0"`        // 评分
	Year      string  `json:"year" gorm:"size:10"`            // 年份
	Directors string  `json:"directors" gorm:"size:500"`      // 导演（多个用逗号分隔）
	Actors    string  `json:"actors" gorm:"size:1000"`        // 演员（多个用逗号分隔）

	// 分类信息
	Category string `json:"category" gorm:"size:50"` // 分类（电影/电视剧）
	SubType  string `json:"sub_type" gorm:"size:50"` // 子类型（华语/欧美/韩国/日本等）

	// 数据来源
	Source   string `json:"source" gorm:"size:50;default:'douban'"` // 数据来源
	DoubanID string `json:"douban_id" gorm:"size:50"`               // 豆瓣ID
}

// TableName 指定表名
func (HotDrama) TableName() string {
	return "hot_dramas"
}
