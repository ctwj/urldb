package models

import (
	"time"
)

// Pan 第三方平台表
type Pan struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`       // Qurak
	Key       int       `json:"key"`        // quark
	Ck        string    `json:"ck"`         // cookie
	IsValid   bool      `json:"is_valid"`   // 是否有效
	Space     int64     `json:"space"`      // 空间
	LeftSpace int64     `json:"left_space"` // 剩余空间
	Remark    string    `json:"remark"`     // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Cks 第三方平台账号cookie表
type Cks struct {
	ID     int    `json:"id"`
	PanID  int    `json:"pan_id"` // pan ID
	T      string `json:"t"`      // cookie类型
	Idx    int    `json:"idx"`    // index
	Ck     string `json:"ck"`     // cookie
	Remark string `json:"remark"` // 备注
}

// CreatePanRequest 创建平台请求
type CreatePanRequest struct {
	Name      string `json:"name" binding:"required"`
	Key       int    `json:"key"`
	Ck        string `json:"ck"`
	IsValid   bool   `json:"is_valid"`
	Space     int64  `json:"space"`
	LeftSpace int64  `json:"left_space"`
	Remark    string `json:"remark"`
}

// UpdatePanRequest 更新平台请求
type UpdatePanRequest struct {
	Name      string `json:"name"`
	Key       int    `json:"key"`
	Ck        string `json:"ck"`
	IsValid   bool   `json:"is_valid"`
	Space     int64  `json:"space"`
	LeftSpace int64  `json:"left_space"`
	Remark    string `json:"remark"`
}

// CreateCksRequest 创建cookie请求
type CreateCksRequest struct {
	PanID  int    `json:"pan_id" binding:"required"`
	T      string `json:"t"`
	Idx    int    `json:"idx"`
	Ck     string `json:"ck"`
	Remark string `json:"remark"`
}

// UpdateCksRequest 更新cookie请求
type UpdateCksRequest struct {
	PanID  int    `json:"pan_id"`
	T      string `json:"t"`
	Idx    int    `json:"idx"`
	Ck     string `json:"ck"`
	Remark string `json:"remark"`
}
