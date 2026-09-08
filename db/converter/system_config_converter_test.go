package converter

import (
	"testing"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// 回归测试（016-api-access-application）：
// 管理员仅修改邮件服务配置并保存时，请求体只含 smtp_* 字段，
// 此前 RequestToSystemConfig 对这类请求返回空切片，导致 handler 误报
// "配置数据转换失败" 并返回 500，且 SMTP 配置从未入库。
func TestRequestToSystemConfigSmtpOnly(t *testing.T) {
	host := "smtp.qq.com"
	port := "465"
	username := "noreply@qq.com"
	password := "auth-code"
	from := "noreply@qq.com"
	encryption := "ssl"

	req := &dto.SystemConfigRequest{
		SmtpHost:       &host,
		SmtpPort:       &port,
		SmtpUsername:   &username,
		SmtpPassword:   &password,
		SmtpFrom:       &from,
		SmtpEncryption: &encryption,
	}

	configs := RequestToSystemConfig(req)
	if len(configs) != 6 {
		t.Fatalf("期望生成 6 条配置，实际 %d 条", len(configs))
	}

	want := map[string]string{
		entity.ConfigKeySmtpHost:       host,
		entity.ConfigKeySmtpPort:       port,
		entity.ConfigKeySmtpUsername:   username,
		entity.ConfigKeySmtpPassword:   password,
		entity.ConfigKeySmtpFrom:       from,
		entity.ConfigKeySmtpEncryption: encryption,
	}
	got := make(map[string]string, len(configs))
	for _, c := range configs {
		if c.Type != entity.ConfigTypeString {
			t.Errorf("配置 %s 的类型应为 %s，实际 %s", c.Key, entity.ConfigTypeString, c.Type)
		}
		got[c.Key] = c.Value
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("配置 %s = %q，期望 %q", k, got[k], v)
		}
	}
}

// API 开放配置仅提交时（前端 feature-config.vue「API 开放配置」组）应正确产出 4 条 int 配置
func TestRequestToSystemConfigApiOpenOnly(t *testing.T) {
	days := 90
	minute := 60
	hour := 1200
	day := 6000

	req := &dto.SystemConfigRequest{
		ApiDefaultValidityDays: &days,
		ApiRateLimitMinute:     &minute,
		ApiRateLimitHour:       &hour,
		ApiRateLimitDay:        &day,
	}

	configs := RequestToSystemConfig(req)
	if len(configs) != 4 {
		t.Fatalf("期望生成 4 条配置，实际 %d 条", len(configs))
	}

	want := map[string]string{
		entity.ConfigKeyApiDefaultValidityDays: "90",
		entity.ConfigKeyApiRateLimitMinute:     "60",
		entity.ConfigKeyApiRateLimitHour:       "1200",
		entity.ConfigKeyApiRateLimitDay:        "6000",
	}
	got := make(map[string]string, len(configs))
	for _, c := range configs {
		if c.Type != entity.ConfigTypeInt {
			t.Errorf("配置 %s 的类型应为 %s，实际 %s", c.Key, entity.ConfigTypeInt, c.Type)
		}
		got[c.Key] = c.Value
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("配置 %s = %q，期望 %q", k, got[k], v)
		}
	}
}

// 响应转换应回填 api_* 字段（此前缺失导致管理后台表单回显为默认值）
func TestSystemConfigToResponseApiOpen(t *testing.T) {
	configs := []entity.SystemConfig{
		{Key: entity.ConfigKeyApiDefaultValidityDays, Value: "90"},
		{Key: entity.ConfigKeyApiRateLimitMinute, Value: "60"},
		{Key: entity.ConfigKeyApiRateLimitHour, Value: "1200"},
		{Key: entity.ConfigKeyApiRateLimitDay, Value: "6000"},
	}

	resp := SystemConfigToResponse(configs)
	if resp.ApiDefaultValidityDays != 90 {
		t.Errorf("ApiDefaultValidityDays = %d，期望 90", resp.ApiDefaultValidityDays)
	}
	if resp.ApiRateLimitMinute != 60 || resp.ApiRateLimitHour != 1200 || resp.ApiRateLimitDay != 6000 {
		t.Errorf("限流回显错误: %d/%d/%d", resp.ApiRateLimitMinute, resp.ApiRateLimitHour, resp.ApiRateLimitDay)
	}
}

// 空请求应返回空切片（handler 据此返回 400 而非 500）
func TestRequestToSystemConfigEmpty(t *testing.T) {
	if configs := RequestToSystemConfig(&dto.SystemConfigRequest{}); len(configs) != 0 {
		t.Fatalf("空请求应返回空切片，实际 %d 条", len(configs))
	}
}

// 响应转换应回填 smtp_* 字段（此前缺失导致后台表单永远回显为空）
func TestSystemConfigToResponseSmtp(t *testing.T) {
	configs := []entity.SystemConfig{
		{Key: entity.ConfigKeySmtpHost, Value: "smtp.qq.com"},
		{Key: entity.ConfigKeySmtpPort, Value: "465"},
		{Key: entity.ConfigKeySmtpFrom, Value: "noreply@qq.com"},
		{Key: entity.ConfigKeySmtpEncryption, Value: "ssl"},
	}

	resp := SystemConfigToResponse(configs)
	if resp.SmtpHost != "smtp.qq.com" {
		t.Errorf("SmtpHost = %q，期望 %q", resp.SmtpHost, "smtp.qq.com")
	}
	if resp.SmtpPort != "465" {
		t.Errorf("SmtpPort = %q，期望 %q", resp.SmtpPort, "465")
	}
	if resp.SmtpFrom != "noreply@qq.com" {
		t.Errorf("SmtpFrom = %q，期望 %q", resp.SmtpFrom, "noreply@qq.com")
	}
	if resp.SmtpEncryption != "ssl" {
		t.Errorf("SmtpEncryption = %q，期望 %q", resp.SmtpEncryption, "ssl")
	}
}
