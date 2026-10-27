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
	// ErrCredentialDisabled 凭证已被管理员停用（不可续期，需管理员恢复或重新申请）
	ErrCredentialDisabled = errors.New("凭证已被停用")
	// ErrCredentialPermanent 凭证永久有效（无到期时间，无需续期）
	ErrCredentialPermanent = errors.New("凭证永久有效，无需续期")
	// ErrInvalidRenewMonths 续期月数不在 1-6 范围内
	ErrInvalidRenewMonths = errors.New("续期月数需为 1-6 个月")
	// ErrApplicationPendingExists 已有待审核申请（开通/续期共用一个待审核位）
	ErrApplicationPendingExists = errors.New("已有申请在审核中")
)

// ApiRenewalMaxMonths 续期单次最大月数（1-6 可选）
const ApiRenewalMaxMonths = 6

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

// RenewApiCredential 用户续期 API 凭证（1-6 个月，可多次叠加）
// 与开通申请的区别：周期可选；凭证未过期时自动通过并立即顺延到期时间（autoApproved=true，
// 返回新到期时间）；已过期时转人工审核，同意后自审核时间起算延长（autoApproved=false）。
// 保留原 key，续期不换凭证。
func RenewApiCredential(userID uint, months int) (app *entity.ApiApplication, expiresAt *time.Time, autoApproved bool, err error) {
	if months < 1 || months > ApiRenewalMaxMonths {
		return nil, nil, false, ErrInvalidRenewMonths
	}

	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil {
		return nil, nil, false, err
	}
	if cred == nil {
		return nil, nil, false, ErrCredentialNotFound
	}
	if cred.Status != entity.ApiCredentialStatusActive {
		return nil, nil, false, ErrCredentialDisabled
	}
	if cred.ExpiresAt == nil {
		return nil, nil, false, ErrCredentialPermanent
	}

	pending, err := repoManager.ApiApplicationRepository.ExistsPendingByUser(userID)
	if err != nil {
		return nil, nil, false, err
	}
	if pending {
		return nil, nil, false, ErrApplicationPendingExists
	}

	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil {
		return nil, nil, false, err
	}
	if user == nil {
		return nil, nil, false, errors.New("用户不存在")
	}

	now := utils.GetCurrentTime()
	app = &entity.ApiApplication{
		UserID:         userID,
		Email:          user.Email,
		Type:           entity.ApiApplicationTypeRenewal,
		DurationMonths: months,
	}

	// 已过期：转人工审核，同意后自审核时间起算
	if cred.IsExpired(now) {
		app.Status = entity.ApiApplicationStatusPending
		if err := repoManager.ApiApplicationRepository.Create(app); err != nil {
			return nil, nil, false, err
		}
		utils.Info("RenewApiCredential - 凭证已过期，续期转人工审核 - userID: %d, months: %d, appID: %d",
			userID, months, app.ID)
		return app, nil, false, nil
	}

	// 未过期：自动通过，自当前到期时间顺延（提前续期不损失剩余时长）
	newExpiresAt := credentialRenewalExpiry(cred, months, now)
	cred.ExpiresAt = newExpiresAt
	if err := repoManager.ApiCredentialRepository.SaveForUser(cred); err != nil {
		return nil, nil, false, err
	}
	app.Status = entity.ApiApplicationStatusApproved
	app.ReviewedAt = &now
	app.ReviewerName = "系统自动通过"
	if err := repoManager.ApiApplicationRepository.Create(app); err != nil {
		// 到期时间已顺延，仅留痕失败，不影响续期结果
		utils.Error("RenewApiCredential - 自动续期留痕失败 - userID: %d, error: %v", userID, err)
	}
	utils.Info("RenewApiCredential - 未过期自动续期 - userID: %d, months: %d, expiresAt: %v",
		userID, months, newExpiresAt)

	go sendRenewalMail(app.Email, months, *newExpiresAt, true)

	return app, newExpiresAt, true, nil
}

// extendUserCredentialExpiry 延长用户凭证有效期并落库，返回新到期时间
func extendUserCredentialExpiry(userID uint, months int, now time.Time) (*time.Time, error) {
	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cred == nil {
		return nil, ErrCredentialNotFound
	}
	if cred.Status != entity.ApiCredentialStatusActive {
		return nil, ErrCredentialDisabled
	}
	expiresAt := credentialRenewalExpiry(cred, months, now)
	cred.ExpiresAt = expiresAt
	if err := repoManager.ApiCredentialRepository.SaveForUser(cred); err != nil {
		return nil, err
	}
	return expiresAt, nil
}

// credentialRenewalExpiry 计算续期后的到期时间：
// 未过期自当前到期时间顺延（叠加不损失剩余时长），已过期自 now 起算
func credentialRenewalExpiry(cred *entity.ApiCredential, months int, now time.Time) *time.Time {
	base := now
	if cred.ExpiresAt != nil && !cred.IsExpired(now) {
		base = *cred.ExpiresAt
	}
	t := base.AddDate(0, months, 0)
	return &t
}

// ApproveApplication 同意申请并签发/重建凭证（FR-005/FR-016）
// 开通申请：签发/重建凭证；续期申请：保留原 key 仅延长有效期。
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

	// 续期申请（过期后转审）：保留原 key，仅延长有效期（未过期自当前到期时间顺延，已过期自审核时间起算）
	if app.Type == entity.ApiApplicationTypeRenewal {
		expiresAt, err := extendUserCredentialExpiry(app.UserID, app.DurationMonths, now)
		if err != nil {
			return nil, nil, err
		}
		app.Status = entity.ApiApplicationStatusApproved
		app.ReviewerID = &reviewerID
		app.ReviewerName = reviewerName
		app.ReviewedAt = &now
		if err := repoManager.ApiApplicationRepository.Update(app); err != nil {
			return nil, nil, err
		}
		utils.Info("ApproveApplication - 续期申请已同意 - appID: %d, userID: %d, months: %d, expiresAt: %v",
			appID, app.UserID, app.DurationMonths, expiresAt)
		go sendRenewalMail(app.Email, app.DurationMonths, *expiresAt, false)
		return app, expiresAt, nil
	}

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
		"<p style=\"color:#4b5563;\">到期后凭证将自动失效，API 将无法继续调用。请提前前往用户中心「API 访问」页面续期（未过期续期自动通过）。</p>" +
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
		"<p style=\"color:#4b5563;\">如需继续使用，请前往用户中心「API 访问」页面提交续期申请（审核通过后恢复，API Key 保持不变）。</p>" +
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

// sendRenewalMail 异步发送续期结果邮件（失败仅记日志，不影响续期结果）
func sendRenewalMail(to string, months int, newExpiresAt time.Time, auto bool) {
	if err := SendApiRenewalMail(to, months, newExpiresAt, auto); err != nil {
		if errors.Is(err, ErrSmtpNotConfigured) {
			utils.Warn("sendRenewalMail - SMTP 未配置，跳过续期邮件 - to: %s", to)
			return
		}
		utils.Error("sendRenewalMail - 续期邮件发送失败 - to: %s, error: %v", to, err)
	}
}

// SendApiRenewalMail 发送续期结果邮件（auto=true 未过期自动续期成功；false 续期申请审核通过）
func SendApiRenewalMail(to string, months int, newExpiresAt time.Time, auto bool) error {
	if to == "" {
		return nil
	}
	title := "您的 API 续期申请已通过"
	lead := "您的 API 续期申请已通过审核。"
	if auto {
		title = "您的 API 凭证续期成功"
		lead = "您的 API 凭证未过期，续期已自动通过并即时生效。"
	}
	body := "<div style=\"font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;\">" +
		"<h2 style=\"color:#1f2937;\">" + title + "</h2>" +
		"<p style=\"color:#4b5563;\">" + lead + "</p>" +
		"<p style=\"color:#4b5563;\">本次续期 <b>" + strconv.Itoa(months) + "</b> 个月，凭证有效期至：<b>" +
		newExpiresAt.Format("2006-01-02 15:04") + "</b>，API Key 保持不变。</p>" +
		"<p style=\"color:#9ca3af;font-size:12px;\">本邮件为系统自动发送，请勿回复。</p></div>"
	return SendSystemMail(to, title, body)
}
