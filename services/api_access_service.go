package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// API 开通/凭证业务错误（016-api-access-application）
var (
	// ErrApplicationNotPending 申请不在待审核状态（非法状态流转）
	ErrApplicationNotPending = errors.New("申请不在待审核状态")
	// ErrApplicationNotFound 申请不存在
	ErrApplicationNotFound = errors.New("申请不存在")
	// ErrCredentialNotFound 用户没有可操作的凭证
	ErrCredentialNotFound = errors.New("凭证不存在")
)

// apiKeyAlphabet API key 随机段字符集（Base62）
const apiKeyAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// apiKeyRandomLength 随机段长度（43 位，含前缀共 49 字符）
const apiKeyRandomLength = 43

// NewApiKey 生成新 API key：返回（明文、展示前缀、SHA-256 哈希）。明文仅返回一次。
func NewApiKey() (plain, prefix, hash string, err error) {
	random, err := gonanoid.Generate(apiKeyAlphabet, apiKeyRandomLength)
	if err != nil {
		return "", "", "", err
	}
	plain = "urldb_" + random
	sum := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(sum[:])
	if len(plain) > 12 {
		prefix = plain[:12]
	} else {
		prefix = plain
	}
	return plain, prefix, hash, nil
}

// GetApiDefaultValidityDays 默认有效天数（缺省 30；0=永久；负值按 30 处理）
func GetApiDefaultValidityDays() int {
	v, err := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyApiDefaultValidityDays)
	if err != nil || v < 0 {
		return 30
	}
	return v
}

// ApproveApplication 同意申请并签发/重建凭证（FR-005/FR-016）
// 返回更新后的申请单与新凭证到期时间（永久时为 nil）
func ApproveApplication(appID, reviewerID uint, reviewerName string) (*entity.ApiApplication, *time.Time, error) {
	app, err := repoManager.ApiApplicationRepository.FindByID(appID)
	if err != nil {
		return nil, nil, ErrApplicationNotFound
	}
	if app.Status != entity.ApiApplicationStatusPending {
		return nil, nil, ErrApplicationNotPending
	}

	now := utils.GetCurrentTime()
	app.Status = entity.ApiApplicationStatusApproved
	app.ReviewerID = &reviewerID
	app.ReviewerName = reviewerName
	app.ReviewedAt = &now
	if err := repoManager.ApiApplicationRepository.Update(app); err != nil {
		return nil, nil, err
	}

	// 签发/重建凭证：到期时间一次性固化（不追溯，spec 澄清决议 2）
	_, prefix, hash, err := NewApiKey()
	if err != nil {
		return nil, nil, err
	}
	var expiresAt *time.Time
	if days := GetApiDefaultValidityDays(); days > 0 {
		t := now.Add(time.Duration(days) * 24 * time.Hour)
		expiresAt = &t
	}
	cred := &entity.ApiCredential{
		UserID:    app.UserID,
		KeyHash:   hash,
		KeyPrefix: prefix,
		Status:    entity.ApiCredentialStatusActive,
		ExpiresAt: expiresAt,
	}
	if err := repoManager.ApiCredentialRepository.SaveForUser(cred); err != nil {
		return nil, nil, err
	}

	utils.Info("ApproveApplication - 申请已同意并签发凭证 - appID: %d, userID: %d, expiresAt: %v",
		appID, app.UserID, expiresAt)
	return app, expiresAt, nil
}

// RejectApplication 拒绝申请（理由可选，FR-006）
func RejectApplication(appID, reviewerID uint, reviewerName, reason string) (*entity.ApiApplication, error) {
	app, err := repoManager.ApiApplicationRepository.FindByID(appID)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	if app.Status != entity.ApiApplicationStatusPending {
		return nil, ErrApplicationNotPending
	}

	now := utils.GetCurrentTime()
	app.Status = entity.ApiApplicationStatusRejected
	app.RejectReason = reason
	app.ReviewerID = &reviewerID
	app.ReviewerName = reviewerName
	app.ReviewedAt = &now
	if err := repoManager.ApiApplicationRepository.Update(app); err != nil {
		return nil, err
	}
	return app, nil
}

// DisableUserCredential 停用用户 API 权限（FR-014，立即失效）
func DisableUserCredential(userID uint) error {
	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil || cred == nil {
		return ErrCredentialNotFound
	}
	if cred.Status != entity.ApiCredentialStatusActive {
		return ErrCredentialNotFound
	}
	cred.Status = entity.ApiCredentialStatusDisabled
	return repoManager.ApiCredentialRepository.SaveForUser(cred)
}

// EnableUserCredential 启用用户 API 权限（B6，恢复被停用的凭证）
// 仅恢复 status=active；到期时间为派生态，若已过期需另行调整到期时间
func EnableUserCredential(userID uint) error {
	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil || cred == nil {
		return ErrCredentialNotFound
	}
	if cred.Status != entity.ApiCredentialStatusDisabled {
		return ErrCredentialNotFound
	}
	cred.Status = entity.ApiCredentialStatusActive
	return repoManager.ApiCredentialRepository.SaveForUser(cred)
}

// UpdateUserCredentialExpiry 管理员调整凭证到期时间（B5）
// expiresAt 为 nil 表示永久；同时清空 reminder_sent_at，使新到期日临近时可再次提醒
func UpdateUserCredentialExpiry(userID uint, expiresAt *time.Time) (*entity.ApiCredential, error) {
	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil || cred == nil {
		return nil, ErrCredentialNotFound
	}
	cred.ExpiresAt = expiresAt
	cred.ReminderSentAt = nil
	if err := repoManager.ApiCredentialRepository.SaveForUser(cred); err != nil {
		return nil, err
	}
	utils.Info("UpdateUserCredentialExpiry - 到期时间已调整 - userID: %d, expiresAt: %v", userID, expiresAt)
	return cred, nil
}

// SendApiExpirationReminderMail 发送凭证到期提醒邮件（FR-018，供 scheduler 调用）
// 返回是否实际发送（用户邮箱缺失/发信失败均不视为成功，reminder_sent_at 不落值）
func SendApiExpirationReminderMail(userID uint, expiresAt time.Time) (bool, error) {
	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil || user == nil || user.Email == "" {
		return false, err
	}

	subject := "您的 API 访问凭证即将到期"
	body := "<div style=\"font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;\">" +
		"<h2 style=\"color:#1f2937;\">API 访问凭证即将到期</h2>" +
		"<p style=\"color:#4b5563;\">您的 API 访问凭证将于 <b>" + expiresAt.Format("2006-01-02 15:04") +
		"</b> 到期。</p>" +
		"<p style=\"color:#4b5563;\">到期后 API 将无法继续调用。如需继续使用，请在到期后重新提交开通申请。</p>" +
		"<p style=\"color:#9ca3af;font-size:12px;\">本邮件为系统自动发送，请勿回复。</p></div>"

	if err := SendSystemMail(user.Email, subject, body); err != nil {
		return false, err
	}
	return true, nil
}
