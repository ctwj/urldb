package entity

import (
	"time"
)

// DownloadHistory 用户下载历史（记录登录用户获取资源链接的行为）
type DownloadHistory struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        uint      `json:"user_id" gorm:"not null;index;uniqueIndex:idx_download_history_user_resource,priority:1;comment:用户ID"`
	ResourceID    uint      `json:"resource_id" gorm:"not null;index;uniqueIndex:idx_download_history_user_resource,priority:2;comment:资源ID"`
	DownloadCount int       `json:"download_count" gorm:"default:1;comment:下载次数"`
	CreatedAt     time.Time `json:"created_at" gorm:"comment:首次下载时间"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"comment:最近下载时间"`
}

// TableName 指定表名
func (DownloadHistory) TableName() string {
	return "download_histories"
}
