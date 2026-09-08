package dto

import "time"

// API 开放访问共享契约结构（016-api-access-application，contracts/api.md）

// ApiCredentialInfo 用户侧凭证概要（A4 credential 字段）
type ApiCredentialInfo struct {
	Exists     bool       `json:"exists"`
	KeyPrefix  string     `json:"key_prefix,omitempty"`
	Status     string     `json:"status,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Expired    bool       `json:"expired,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

// ApiApplicationAdminItem 管理端申请列表项（B1）
type ApiApplicationAdminItem struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Purpose      string     `json:"purpose"`
	Status       string     `json:"status"`
	RejectReason string     `json:"reject_reason"`
	ReviewerName string     `json:"reviewer_name"`
	ReviewedAt   *time.Time `json:"reviewed_at"`
	CreatedAt    time.Time  `json:"created_at"`

	// 已同意申请关联的凭证状态（approved 行填充，其余为空值）
	CredentialStatus  string     `json:"credential_status,omitempty"`              // active/disabled
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`                     // 到期时间，NULL=永久
	CredentialExpired bool       `json:"credential_expired,omitempty"`             // 凭证是否已过期（派生态）
}
