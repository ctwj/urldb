package entity

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Password  string         `json:"-" gorm:"size:255;not null;comment:密码"`
	Email     string         `json:"email" gorm:"size:100;comment:邮箱"`
	Role      string         `json:"role" gorm:"size:20;default:'user';comment:角色"`
	IsActive  bool           `json:"is_active" gorm:"default:true;comment:是否激活"`
	LastLogin *time.Time     `json:"last_login" gorm:"comment:最后登录时间"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
