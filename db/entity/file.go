package entity

import (
	"time"
)

// File 文件实体
type File struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 文件信息
	OriginalName string `json:"original_name" gorm:"size:255;not null;comment:原始文件名"`
	FileName     string `json:"file_name" gorm:"size:255;not null;unique;comment:存储文件名"`
	FilePath     string `json:"file_path" gorm:"size:500;not null;comment:文件路径"`
	FileSize     int64  `json:"file_size" gorm:"not null;comment:文件大小(字节)"`
	FileType     string `json:"file_type" gorm:"size:100;not null;comment:文件类型"`
	MimeType     string `json:"mime_type" gorm:"size:100;comment:MIME类型"`

	// 访问信息
	AccessURL string `json:"access_url" gorm:"size:500;comment:访问URL"`

	// 用户信息
	UserID uint `json:"user_id" gorm:"comment:上传用户ID"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// 状态信息
	Status    string `json:"status" gorm:"size:20;default:'active';comment:文件状态"`
	IsPublic  bool   `json:"is_public" gorm:"default:true;comment:是否公开"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false;comment:是否已删除"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// FileStatus 文件状态常量
const (
	FileStatusActive   = "active"   // 正常
	FileStatusInactive = "inactive" // 禁用
	FileStatusDeleted  = "deleted"  // 已删除
)
