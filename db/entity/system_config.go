package entity

import (
	"time"
)

// SystemConfig 系统配置实体
type SystemConfig struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// SEO 配置
	SiteTitle       string `json:"site_title" gorm:"size:200;not null;default:'网盘资源管理系统'"`
	SiteDescription string `json:"site_description" gorm:"size:500"`
	Keywords        string `json:"keywords" gorm:"size:500"`
	Author          string `json:"author" gorm:"size:100"`
	Copyright       string `json:"copyright" gorm:"size:200"`

	// 自动处理配置组
	AutoProcessReadyResources bool `json:"auto_process_ready_resources" gorm:"default:false"` // 自动处理待处理资源
	AutoProcessInterval       int  `json:"auto_process_interval" gorm:"default:30"`           // 自动处理间隔（分钟）
	AutoTransferEnabled       bool `json:"auto_transfer_enabled" gorm:"default:false"`        // 开启自动转存
	AutoFetchHotDramaEnabled  bool `json:"auto_fetch_hot_drama_enabled" gorm:"default:false"` // 自动拉取热播剧名字

	// API配置
	ApiToken string `json:"api_token" gorm:"size:100;uniqueIndex"` // 公开API访问令牌

	// 其他配置
	PageSize        int  `json:"page_size" gorm:"default:100"`
	MaintenanceMode bool `json:"maintenance_mode" gorm:"default:false"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}
