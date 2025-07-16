package entity

import (
	"time"

	"gorm.io/gorm"
)

// Cks 第三方平台账号cookie表
type Cks struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PanID       uint           `json:"pan_id" gorm:"not null;comment:平台ID"`
	Idx         int            `json:"idx" gorm:"comment:索引"`
	Ck          string         `json:"ck" gorm:"type:text;comment:cookie"`
	IsValid     bool           `json:"is_valid" gorm:"default:true;comment:是否有效"`
	Space       int64          `json:"space" gorm:"default:0;comment:总空间(字节)"`
	LeftSpace   int64          `json:"left_space" gorm:"default:0;comment:剩余空间(字节)"`
	UsedSpace   int64          `json:"used_space" gorm:"default:0;comment:已使用空间(字节)"`
	Username    string         `json:"username" gorm:"size:100;comment:用户名"`
	VipStatus   bool           `json:"vip_status" gorm:"default:false;comment:VIP状态"`
	ServiceType string         `json:"service_type" gorm:"size:20;comment:服务类型"`
	Remark      string         `json:"remark" gorm:"size:64;not null;comment:备注"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联关系
	Pan Pan `json:"pan" gorm:"foreignKey:PanID"`
}

// TableName 指定表名
func (Cks) TableName() string {
	return "cks"
}
