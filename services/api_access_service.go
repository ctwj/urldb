package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
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

	// 审批通过邮件通知（异步发送，失败仅记日志，不影响审批结果）
	go func() {
		if err := SendApiApprovalMail(app.Email, expiresAt); err != nil {
			if errors.Is(err, ErrSmtpNotConfigured) {
				utils.Warn("ApproveApplication - SMTP 未配置，跳过审批通过邮件 - appID: %d", appID)
				return
			}
			utils.Error("ApproveApplication - 审批通过邮件发送失败 - appID: %d, userID: %d, error: %v", appID, app.UserID, err)
		}
	}()

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
// expiresAt 为 nil 表示永久；到期提醒按自然日实时计算，新到期日自动生效
func UpdateUserCredentialExpiry(userID uint, expiresAt *time.Time) (*entity.ApiCredential, error) {
	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil || cred == nil {
		return nil, ErrCredentialNotFound
	}
	cred.ExpiresAt = expiresAt
	if err := repoManager.ApiCredentialRepository.SaveForUser(cred); err != nil {
		return nil, err
	}
	utils.Info("UpdateUserCredentialExpiry - 到期时间已调整 - userID: %d, expiresAt: %v", userID, expiresAt)
	return cred, nil
}

// SendApiExpirationReminderMail 发送凭证到期提醒邮件（FR-018，供 scheduler 调用）
// daysLeft 为按自然日计算的剩余天数；用户不存在或未填邮箱时静默跳过
func SendApiExpirationReminderMail(userID uint, expiresAt time.Time, daysLeft int) error {
	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil || user.Email == "" {
		utils.Debug("SendApiExpirationReminderMail - 用户不存在或未填邮箱，跳过 - userID: %d", userID)
		return nil
	}

	subject := "您的 API 访问凭证即将到期"
	body := "<div style=\"font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;\">" +
		"<h2 style=\"color:#1f2937;\">API 访问凭证即将到期</h2>" +
		"<p style=\"color:#4b5563;\">您的 API 访问凭证将于 <b>" + expiresAt.Format("2006-01-02 15:04") +
		"</b> 到期，还剩约 <b>" + strconv.Itoa(daysLeft) + "</b> 天。</p>" +
		"<p style=\"color:#4b5563;\">到期后凭证将自动失效，API 将无法继续调用。如需继续使用，请在到期后重新提交开通申请。</p>" +
		"<p style=\"color:#9ca3af;font-size:12px;\">本邮件为系统自动发送，请勿回复。</p></div>"
	return SendSystemMail(user.Email, subject, body)
}

// SendApiApprovalMail 发送审批通过通知邮件（附带凭证到期时间）

// SendApiExpiredMail 发送凭证过期通知邮件（FR-018，供 scheduler 调用，过期后第一个自然日触发）
// 用户不存在或未填邮箱时静默跳过
func SendApiExpiredMail(userID uint, expiredAt time.Time) error {
	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil || user.Email == "" {
		utils.Debug("SendApiExpiredMail - 用户不存在或未填邮箱，跳过 - userID: %d", userID)
		return nil
	}

	subject := "您的 API 访问凭证已过期"
	body := "<div style=\"font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;\">" +
		"<h2 style=\"color:#dc2626;\">API 访问凭证已过期</h2>" +
		"<p style=\"color:#4b5563;\">您的 API 访问凭证已于 <b>" + expiredAt.Format("2006-01-02 15:04") +
		"</b> 过期，现已自动失效，API 调用将返回 403。</p>" +
		"<p style=\"color:#4b5563;\">如需继续使用，请前往用户中心「API 访问」页面重新提交开通申请。</p>" +
		"<p style=\"color:#9ca3af;font-size:12px;\">本邮件为系统自动发送，请勿回复。</p></div>"
	return SendSystemMail(user.Email, subject, body)
}
// to 为申请时认证的邮箱；expiresAt 为 nil 表示永久有效
func SendApiApprovalMail(to string, expiresAt *time.Time) error {
	if to == "" {
		return nil
	}
	expiryText := "永久有效"
	if expiresAt != nil {
		expiryText = "<b>" + expiresAt.Format("2006-01-02 15:04") + "</b>"
	}
	subject := "您的 API 访问申请已通过"
	body := "<div style=\"font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;\">" +
		"<h2 style=\"color:#1f2937;\">API 访问申请已通过</h2>" +
		"<p style=\"color:#4b5563;\">您的 API 访问申请已通过审核，API 凭证已签发。</p>" +
		"<p style=\"color:#4b5563;\">凭证有效期至：" + expiryText + "，到期后凭证将自动失效。</p>" +
		"<p style=\"color:#4b5563;\">请前往用户中心「API 访问」页面查看凭证并获取 API Key。</p>" +
		"<p style=\"color:#9ca3af;font-size:12px;\">本邮件为系统自动发送，请勿回复。</p></div>"
	return SendSystemMail(to, subject, body)
}
