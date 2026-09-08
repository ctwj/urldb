package entity

import (
	"time"

	"gorm.io/gorm"
)

// API 凭证状态（016-api-access-application）
const (
	ApiCredentialStatusActive   = "active"   // 有效（是否过期为派生态，见 IsExpired）
	ApiCredentialStatusDisabled = "disabled" // 已被管理员停用
)

// ApiCredential 用户 API 调用凭证（一用户一行，重置=同行换 key）
type ApiCredential struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint           `json:"user_id" gorm:"not null;uniqueIndex;comment:归属用户ID"`
	KeyHash        string         `json:"-" gorm:"size:64;not null;uniqueIndex;comment:明文key的SHA-256"`
	KeyPrefix      string         `json:"key_prefix" gorm:"size:16;not null;comment:明文前12字符(展示辨认)"`
	Status         string         `json:"status" gorm:"size:20;not null;default:active;comment:状态:active/disabled"`
	ExpiresAt      *time.Time     `json:"expires_at" gorm:"comment:到期时间;NULL=永久"`
	ReminderSentAt *time.Time     `json:"reminder_sent_at" gorm:"comment:到期前7天提醒已发送时间"`
	LastUsedAt     *time.Time     `json:"last_used_at" gorm:"comment:最近调用时间"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (ApiCredential) TableName() string {
	return "api_credentials"
}

// IsExpired 派生态：到期时间非空且已过（FR-016/FR-017，过期不落库，调用时惰性判定）
func (c *ApiCredential) IsExpired(now time.Time) bool {
	return c.ExpiresAt != nil && now.After(*c.ExpiresAt)
}
