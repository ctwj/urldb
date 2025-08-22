package entity

import (
	"time"

	"gorm.io/gorm"
)

// Resource 资源模型
type Resource struct {
	ID                  uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title               string         `json:"title" gorm:"size:255;not null;comment:资源标题"`
	Description         string         `json:"description" gorm:"type:text;comment:资源描述"`
	URL                 string         `json:"url" gorm:"size:128;comment:资源链接"`
	PanID               *uint          `json:"pan_id" gorm:"comment:平台ID"`
	SaveURL             string         `json:"save_url" gorm:"size:500;comment:转存后的链接"`
	FileSize            string         `json:"file_size" gorm:"size:100;comment:文件大小"`
	CategoryID          *uint          `json:"category_id" gorm:"comment:分类ID"`
	ViewCount           int            `json:"view_count" gorm:"default:0;comment:浏览次数"`
	IsValid             bool           `json:"is_valid" gorm:"default:true;comment:是否有效"`
	IsPublic            bool           `json:"is_public" gorm:"default:true;comment:是否公开"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Cover               string         `json:"cover" gorm:"size:500;comment:封面"`
	Author              string         `json:"author" gorm:"size:100;comment:作者"`
	ErrorMsg            string         `json:"error_msg" gorm:"size:255;comment:转存失败原因"`
	CkID                *uint          `json:"ck_id" gorm:"comment:账号ID"`
	Fid                 string         `json:"fid" gorm:"size:128;comment:网盘文件ID"`
	Key                 string         `json:"key" gorm:"size:64;index;comment:资源组标识，相同key表示同一组资源"`
	SyncedToMeilisearch bool           `json:"synced_to_meilisearch" gorm:"default:false;comment:是否已同步到Meilisearch"`
	SyncedAt            *time.Time     `json:"synced_at" gorm:"comment:同步时间"`

	// 关联关系
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	Pan      Pan      `json:"pan" gorm:"foreignKey:PanID"`
	Tags     []Tag    `json:"tags" gorm:"many2many:resource_tags;"`
}

// TableName 指定表名
func (Resource) TableName() string {
	return "resources"
}

// GetTitle 获取资源标题（实现utils.Resource接口）
func (r *Resource) GetTitle() string {
	return r.Title
}

// GetDescription 获取资源描述（实现utils.Resource接口）
func (r *Resource) GetDescription() string {
	return r.Description
}

// SetTitle 设置资源标题（实现utils.Resource接口）
func (r *Resource) SetTitle(title string) {
	r.Title = title
}

// SetDescription 设置资源描述（实现utils.Resource接口）
func (r *Resource) SetDescription(description string) {
	r.Description = description
}
