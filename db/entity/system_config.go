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

	// 自动处理配置
	AutoProcessReadyResources bool `json:"auto_process_ready_resources" gorm:"default:false"`
	AutoProcessInterval       int  `json:"auto_process_interval" gorm:"default:30"` // 分钟

	// 其他配置
	PageSize        int  `json:"page_size" gorm:"default:100"`
	MaintenanceMode bool `json:"maintenance_mode" gorm:"default:false"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}
