package services

import (
	"context"
	"errors"
	"fmt"
	netmail "net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	mail "github.com/wneessen/go-mail"
)

// ErrSmtpNotConfigured SMTP 服务未配置（后台未填写 smtp_host/smtp_from）
var ErrSmtpNotConfigured = errors.New("SMTP_NOT_CONFIGURED")

// smtpConfig SMTP 发信配置（来自 SystemConfig，016-api-access-application）
type smtpConfig struct {
	host       string
	port       int
	username   string
	password   string
	from       string
	encryption string // ssl / starttls / plain
}

// loadSmtpConfig 从系统配置缓存读取 SMTP 配置
func loadSmtpConfig() (*smtpConfig, error) {
	repo := repoManager.SystemConfigRepository

	host, _ := repo.GetConfigValue(entity.ConfigKeySmtpHost)
	from, _ := repo.GetConfigValue(entity.ConfigKeySmtpFrom)
	if host == "" || from == "" {
		return nil, ErrSmtpNotConfigured
	}

	cfg := &smtpConfig{
		host:       host,
		port:       465,
		username:   "",
		password:   "",
		from:       from,
		encryption: "ssl",
	}

	if v, err := repo.GetConfigInt(entity.ConfigKeySmtpPort); err == nil && v > 0 && v <= 65535 {
		cfg.port = v
	}
	if v, _ := repo.GetConfigValue(entity.ConfigKeySmtpUsername); v != "" {
		cfg.username = v
	}
	if v, _ := repo.GetConfigValue(entity.ConfigKeySmtpPassword); v != "" {
		cfg.password = v
	}
	if v, _ := repo.GetConfigValue(entity.ConfigKeySmtpEncryption); v != "" {
		cfg.encryption = v
	}
	return cfg, nil
}

// buildSmtpClientOptions 构造 go-mail 客户端选项（发送与连接测试共用，保证测试行为与真实发送一致）
func buildSmtpClientOptions(cfg *smtpConfig) []mail.Option {
	var opts []mail.Option
	opts = append(opts, mail.WithPort(cfg.port), mail.WithTimeout(15*time.Second))

	switch cfg.encryption {
	case "starttls":
		opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	case "plain":
		opts = append(opts, mail.WithTLSPolicy(mail.NoTLS))
	default: // ssl
		opts = append(opts, mail.WithSSL())
	}

	if cfg.username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.username),
			mail.WithPassword(cfg.password),
		)
	}
	return opts
}

// SendSystemMail 通过后台 SMTP 配置发送 HTML 邮件（邮箱验证码、API 到期提醒共用通道）
func SendSystemMail(to, subject, htmlBody string) error {
	cfg, err := loadSmtpConfig()
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	if err := msg.From(cfg.from); err != nil {
		return fmt.Errorf("发件人地址无效: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("收件人地址无效: %w", err)
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, htmlBody)

	client, err := mail.NewClient(cfg.host, buildSmtpClientOptions(cfg)...)
	if err != nil {
		return fmt.Errorf("创建邮件客户端失败: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		utils.Error("SendSystemMail - 邮件发送失败 - to: %s, subject: %s, error: %v", to, subject, err)
		return fmt.Errorf("邮件发送失败: %w", err)
	}

	utils.Info("SendSystemMail - 邮件发送成功 - to: %s, subject: %s", to, subject)
	return nil
}

// SmtpTestRequest SMTP 连接测试请求（016-api-access-application）
// 字段为 nil 时回退到系统配置中已保存的值，便于校验未保存的表单或已保存的配置
type SmtpTestRequest struct {
	Host       *string
	Port       *string
	Username   *string
	Password   *string
	From       *string
	Encryption *string
}

// TestSmtpConnection 校验 SMTP 配置有效性：建立连接并完成认证（不实际发信）。
// 返回建连耗时，用于前端展示。
func TestSmtpConnection(req *SmtpTestRequest) (time.Duration, error) {
	// 先校验请求覆盖值（无需访问数据库）
	if req != nil {
		if req.Port != nil {
			p, err := strconv.Atoi(strings.TrimSpace(*req.Port))
			if err != nil || p <= 0 || p > 65535 {
				return 0, fmt.Errorf("SMTP 端口无效: %s", strings.TrimSpace(*req.Port))
			}
		}
		if req.Encryption != nil {
			switch strings.TrimSpace(*req.Encryption) {
			case "ssl", "starttls", "plain":
			default:
				return 0, fmt.Errorf("加密方式无效: %s（支持 ssl/starttls/plain）", strings.TrimSpace(*req.Encryption))
			}
		}
	}

	// 取已保存配置作为基线（可能未配置）
	var cfg *smtpConfig
	if c, err := loadSmtpConfig(); err == nil {
		cfg = c
	}

	// 用请求字段覆盖（校验前端表单当前值）
	if req != nil {
		if cfg == nil {
			cfg = &smtpConfig{port: 465, encryption: "ssl"}
		}
		if req.Host != nil {
			cfg.host = strings.TrimSpace(*req.Host)
		}
		if req.From != nil {
			cfg.from = strings.TrimSpace(*req.From)
		}
		if req.Username != nil {
			cfg.username = strings.TrimSpace(*req.Username)
		}
		if req.Password != nil {
			cfg.password = *req.Password
		}
		if req.Encryption != nil {
			cfg.encryption = strings.TrimSpace(*req.Encryption)
		}
		if req.Port != nil {
			p, _ := strconv.Atoi(strings.TrimSpace(*req.Port))
			cfg.port = p
		}
	}

	if cfg == nil || cfg.host == "" || cfg.from == "" {
		return 0, ErrSmtpNotConfigured
	}
	if _, err := netmail.ParseAddress(cfg.from); err != nil {
		return 0, fmt.Errorf("发件人地址无效: %s", cfg.from)
	}

	client, err := mail.NewClient(cfg.host, buildSmtpClientOptions(cfg)...)
	if err != nil {
		return 0, fmt.Errorf("创建邮件客户端失败: %w", err)
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := client.DialWithContext(ctx); err != nil {
		utils.Error("TestSmtpConnection - 连接失败 - host: %s, port: %d, error: %v", cfg.host, cfg.port, err)
		return 0, fmt.Errorf("连接失败: %w", err)
	}
	defer client.Close()

	elapsed := time.Since(start)
	utils.Info("TestSmtpConnection - 连接成功 - host: %s, port: %d, encryption: %s, 耗时: %s", cfg.host, cfg.port, cfg.encryption, elapsed)
	return elapsed, nil
}

// maskEmail 邮箱打码（a***@qq.com）
func maskEmail(email string) string {
	at := -1
	for i, r := range email {
		if r == '@' {
			at = i
			break
		}
	}
	if at <= 0 {
		return "***"
	}
	prefix := email[:1]
	return prefix + "***" + email[at:]
}
