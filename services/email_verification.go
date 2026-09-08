package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// 邮箱验证码错误（016-api-access-application FR-001）
var (
	// ErrEmailCodeTooFrequently 发码过于频繁（同邮箱 60 秒内重复请求）
	ErrEmailCodeTooFrequently = errors.New("EMAIL_CODE_TOO_FREQUENTLY")
)

const (
	emailCodeTTL            = 10 * time.Minute // 验证码有效期
	emailCodeResendInterval = 60 * time.Second // 同邮箱重发间隔
)

// emailCodeEntry 单个邮箱的验证码记录
type emailCodeEntry struct {
	code       string
	expiresAt  time.Time
	lastSentAt time.Time
}

var (
	emailCodeMutex sync.Mutex
	emailCodes     = map[string]emailCodeEntry{} // email -> entry（单实例内存态，重启失效可重发）
)

// generateEmailCode 生成 6 位数字验证码
func generateEmailCode() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// 极端情况退化：用时间戳取模
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// SendEmailVerificationCode 生成并向邮箱发送验证码（60 秒发码间隔）
func SendEmailVerificationCode(email string) error {
	now := utils.GetCurrentTime()

	emailCodeMutex.Lock()
	if entry, ok := emailCodes[email]; ok && entry.lastSentAt.Add(emailCodeResendInterval).After(now) {
		emailCodeMutex.Unlock()
		return ErrEmailCodeTooFrequently
	}
	code := generateEmailCode()
	emailCodes[email] = emailCodeEntry{
		code:       code,
		expiresAt:  now.Add(emailCodeTTL),
		lastSentAt: now,
	}
	emailCodeMutex.Unlock()

	subject := fmt.Sprintf("邮箱验证码：%s（%d 分钟内有效）", code, int(emailCodeTTL.Minutes()))
	body := fmt.Sprintf(
		`<div style="font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;">
<h2 style="color:#1f2937;">邮箱验证码</h2>
<p style="color:#4b5563;">您正在验证邮箱，请使用以下验证码（%d 分钟内有效）：</p>
<p style="font-size:28px;font-weight:bold;letter-spacing:6px;color:#2563eb;margin:20px 0;">%s</p>
<p style="color:#9ca3af;font-size:12px;">如果这不是您本人的操作，请忽略本邮件。</p>
</div>`, int(emailCodeTTL.Minutes()), code)

	if err := SendSystemMail(email, subject, body); err != nil {
		// 发送失败则清除记录，允许用户立即重试
		emailCodeMutex.Lock()
		delete(emailCodes, email)
		emailCodeMutex.Unlock()
		return err
	}
	return nil
}

// VerifyEmailCode 校验验证码（一次性，成功即销毁）
func VerifyEmailCode(email, code string) bool {
	emailCodeMutex.Lock()
	defer emailCodeMutex.Unlock()

	entry, ok := emailCodes[email]
	if !ok {
		return false
	}
	if utils.GetCurrentTime().After(entry.expiresAt) {
		delete(emailCodes, email)
		return false
	}
	if entry.code != code {
		return false
	}
	delete(emailCodes, email)
	return true
}

// CleanupExpiredEmailCodes 清理过期验证码（防内存缓慢增长；由调用方低频触发）
func CleanupExpiredEmailCodes() {
	now := utils.GetCurrentTime()
	emailCodeMutex.Lock()
	defer emailCodeMutex.Unlock()
	for email, entry := range emailCodes {
		if now.After(entry.expiresAt) {
			delete(emailCodes, email)
		}
	}
}
