package entity

import (
	"gorm.io/gorm"
	"time"
)

// CopyrightClaim 版权申述实体
type CopyrightClaim struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ResourceKey  string         `gorm:"type:varchar(255);not null;index" json:"resource_key"` // 资源唯一标识
	Identity     string         `gorm:"type:varchar(50);not null" json:"identity"`            // 申述人身份
	ProofType    string         `gorm:"type:varchar(50);not null" json:"proof_type"`          // 证明类型
	Reason       string         `gorm:"type:text;not null" json:"reason"`                     // 申述理由
	ContactInfo  string         `gorm:"type:varchar(255);not null" json:"contact_info"`       // 联系信息
	ClaimantName string         `gorm:"type:varchar(100);not null" json:"claimant_name"`      // 申述人姓名
	ProofFiles   string         `gorm:"type:text" json:"proof_files"`                         // 证明文件（JSON格式）
	UserAgent    string         `gorm:"type:text" json:"user_agent"`                          // 用户代理
	IPAddress    string         `gorm:"type:varchar(45)" json:"ip_address"`                   // IP地址
	Status       string         `gorm:"type:varchar(20);default:'pending'" json:"status"`     // 处理状态: pending, approved, rejected
	ProcessedAt  *time.Time     `json:"processed_at"`                                         // 处理时间
	ProcessedBy  *uint          `json:"processed_by"`                                         // 处理人ID
	Note         string         `gorm:"type:text" json:"note"`                                // 处理备注
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

// TableName 表名
func (CopyrightClaim) TableName() string {
	return "copyright_claims"
}