package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// 用户 API 开放访问处理器（016-api-access-application）
// 覆盖：邮箱认证（A1/A2）、开通申请（A3）、状态查询（A4）、凭证重置（A5）

// maskEmailForDisplay 邮箱打码（a***@qq.com）
func maskEmailForDisplay(email string) string {
	at := strings.Index(email, "@")
	if at <= 0 {
		return "***"
	}
	return email[:1] + "***" + email[at:]
}

// SendEmailCode 发送邮箱验证码（A1，FR-001 前置能力）
func SendEmailCode(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}
	if strings.TrimSpace(user.Email) == "" {
		ErrorResponse(c, "请先在个人信息中填写邮箱", http.StatusBadRequest)
		return
	}

	if err := services.SendEmailVerificationCode(user.Email); err != nil {
		if errors.Is(err, services.ErrEmailCodeTooFrequently) {
			ErrorResponse(c, "发送过于频繁，请稍后再试", http.StatusTooManyRequests)
			return
		}
		if errors.Is(err, services.ErrSmtpNotConfigured) {
			ErrorResponse(c, "邮件服务未配置，请联系管理员", http.StatusBadRequest)
			return
		}
		utils.Error("SendEmailCode - 验证码发送失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "验证码发送失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"sent":  true,
		"email": maskEmailForDisplay(user.Email),
	})
}

// VerifyEmailCode 校验邮箱验证码（A2，通过后置 email_verified=true）
func VerifyEmailCode(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "验证码不能为空", http.StatusBadRequest)
		return
	}

	userID := c.GetUint("user_id")
	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	if !services.VerifyEmailCode(user.Email, strings.TrimSpace(req.Code)) {
		ErrorResponse(c, "验证码错误或已过期", http.StatusBadRequest)
		return
	}

	user.EmailVerified = true
	if err := repoManager.UserRepository.Update(user); err != nil {
		utils.Error("VerifyEmailCode - 认证状态写入失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"email_verified": true})
}

// ApplyApiAccess 提交 API 开通申请（A3，FR-001/FR-002/FR-003）
func ApplyApiAccess(c *gin.Context) {
	var req struct {
		Purpose string `json:"purpose" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "用途说明不能为空", http.StatusBadRequest)
		return
	}
	purpose := strings.TrimSpace(req.Purpose)
	if len([]rune(purpose)) > 500 {
		ErrorResponse(c, "用途说明不能超过 500 字", http.StatusBadRequest)
		return
	}

	userID := c.GetUint("user_id")
	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	// 顺序校验：邮箱认证 → 待审核查重 → 已开通查重（FR-003）
	if !user.EmailVerified {
		ErrorResponse(c, "请先完成邮箱认证", http.StatusForbidden)
		return
	}
	pending, err := repoManager.ApiApplicationRepository.ExistsPendingByUser(userID)
	if err != nil {
		utils.Error("ApplyApiAccess - 待审核查重失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}
	if pending {
		ErrorResponse(c, "已有申请在审核中，请耐心等待", http.StatusBadRequest)
		return
	}
	now := utils.GetCurrentTime()
	granted, err := repoManager.ApiCredentialRepository.ExistsActiveUnexpiredByUser(userID, now)
	if err != nil {
		utils.Error("ApplyApiAccess - 开通状态查重失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}
	if granted {
		ErrorResponse(c, "API 权限已开通，无需重复申请", http.StatusBadRequest)
		return
	}

	app := &entity.ApiApplication{
		UserID:  userID,
		Email:   user.Email,
		Purpose: purpose,
		Status:  entity.ApiApplicationStatusPending,
	}
	if err := repoManager.ApiApplicationRepository.Create(app); err != nil {
		utils.Error("ApplyApiAccess - 申请创建失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "提交失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	utils.Info("ApplyApiAccess - 申请提交成功 - userID: %d, appID: %d", userID, app.ID)
	SuccessResponse(c, app)
}

// GetMyApiStatus 查询我的 API 状态（A4，FR-007）
func GetMyApiStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := repoManager.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		ErrorResponse(c, "用户不存在", http.StatusNotFound)
		return
	}

	app, err := repoManager.ApiApplicationRepository.FindLatestByUser(userID)
	if err != nil {
		utils.Error("GetMyApiStatus - 申请查询失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil {
		utils.Error("GetMyApiStatus - 凭证查询失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	var credential interface{}
	if cred == nil {
		credential = gin.H{"exists": false}
	} else {
		now := utils.GetCurrentTime()
		credential = gin.H{
			"exists":       true,
			"key_prefix":   cred.KeyPrefix,
			"status":       cred.Status,
			"expires_at":   cred.ExpiresAt,
			"expired":      cred.IsExpired(now),
			"last_used_at": cred.LastUsedAt,
			"created_at":   cred.CreatedAt,
		}
	}

	SuccessResponse(c, gin.H{
		"email_verified": user.EmailVerified,
		"application":    app,
		"credential":     credential,
	})
}

// ResetMyApiKey 重置 API 凭证（A5，FR-008：旧凭证立即失效；expires_at 不变）
func ResetMyApiKey(c *gin.Context) {
	userID := c.GetUint("user_id")

	cred, err := repoManager.ApiCredentialRepository.FindByUserID(userID)
	if err != nil {
		utils.Error("ResetMyApiKey - 凭证查询失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}
	if cred == nil {
		ErrorResponse(c, "尚未开通 API 权限", http.StatusBadRequest)
		return
	}

	now := utils.GetCurrentTime()
	if cred.Status != entity.ApiCredentialStatusActive {
		ErrorResponse(c, "权限已停用，如需使用请重新申请", http.StatusForbidden)
		return
	}
	if cred.IsExpired(now) {
		ErrorResponse(c, "凭证已过期，请重新申请", http.StatusForbidden)
		return
	}

	plain, prefix, hash, err := services.NewApiKey()
	if err != nil {
		utils.Error("ResetMyApiKey - 凭证生成失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}
	cred.KeyHash = hash
	cred.KeyPrefix = prefix
	cred.ReminderSentAt = nil
	if err := repoManager.ApiCredentialRepository.SaveForUser(cred); err != nil {
		utils.Error("ResetMyApiKey - 凭证保存失败 - userID: %d, error: %v", userID, err)
		ErrorResponse(c, "操作失败，请稍后重试", http.StatusInternalServerError)
		return
	}

	utils.Info("ResetMyApiKey - 凭证已重置 - userID: %d", userID)
	SuccessResponse(c, gin.H{
		"key":        plain,
		"key_prefix": prefix,
		"expires_at": cred.ExpiresAt,
	})
}
