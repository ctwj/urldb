package common

import (
	"testing"
)

func TestExtractShareID_123pan(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantID   string
		wantType string
	}{
		{
			name:     "传统 /s/ 格式",
			url:      "https://www.123pan.com/s/i4uaTd-WHn0",
			wantID:   "i4uaTd-WHn0",
			wantType: Pan123Str,
		},
		{
			name:     "新 /123pan/ 格式",
			url:      "https://1856557151.share.123pan.cn/123pan/oJqrvd-KlG9ds",
			wantID:   "oJqrvd-KlG9ds",
			wantType: Pan123Str,
		},
		{
			name:     "share.123pan.cn /s/ 格式",
			url:      "https://abc123.share.123pan.cn/s/AbCdEf-12345",
			wantID:   "AbCdEf-12345",
			wantType: Pan123Str,
		},
		{
			name:     "123912.com 域名",
			url:      "https://www.123912.com/s/U8f2Td-ZeOX",
			wantID:   "U8f2Td-ZeOX",
			wantType: Pan123Str,
		},
		{
			name:     "123684.com 域名",
			url:      "https://www.123684.com/s/u9izjv-k3uWv",
			wantID:   "u9izjv-k3uWv",
			wantType: Pan123Str,
		},
		{
			name:     "带查询参数的链接",
			url:      "https://1856557151.share.123pan.cn/123pan/oJqrvd-KlG9ds?entry=2",
			wantID:   "oJqrvd-KlG9ds",
			wantType: Pan123Str,
		},
		{
			name:     "无效链接 - 不支持的路径",
			url:      "https://www.123pan.com/invalid/path",
			wantID:   "",
			wantType: NotFoundStr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotType := extractShareID(tt.url)
			if gotID != tt.wantID {
				t.Errorf("extractShareID() ID = %v, want %v", gotID, tt.wantID)
			}
			if gotType != tt.wantType {
				t.Errorf("extractShareID() Type = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}

func TestExtractShareIdString(t *testing.T) {
	// 测试导出的函数
	shareID, service := ExtractShareIdString("https://1856557151.share.123pan.cn/123pan/oJqrvd-KlG9ds")
	if shareID != "oJqrvd-KlG9ds" {
		t.Errorf("Expected shareID 'oJqrvd-KlG9ds', got '%s'", shareID)
	}
	if service != Pan123Str {
		t.Errorf("Expected service '%s', got '%s'", Pan123Str, service)
	}
}

func TestExtractShareID_OtherPans(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantID   string
		wantType string
	}{
		{
			name:     "夸克网盘",
			url:      "https://pan.quark.cn/s/9803af406f13",
			wantID:   "9803af406f13",
			wantType: QuarkStr,
		},
		{
			name:     "阿里云盘 alipan",
			url:      "https://www.alipan.com/s/QbaHJ71QjV1",
			wantID:   "QbaHJ71QjV1",
			wantType: AliyunStr,
		},
		{
			name:     "阿里云盘 aliyundrive",
			url:      "https://www.aliyundrive.com/s/hz1HHxhahsE",
			wantID:   "hz1HHxhahsE",
			wantType: AliyunStr,
		},
		{
			name:     "百度网盘",
			url:      "https://pan.baidu.com/s/1rIcc6X7D3rVzNSqivsRejw",
			wantID:   "1rIcc6X7D3rVzNSqivsRejw",
			wantType: BaiduStr,
		},
		{
			name:     "UC网盘",
			url:      "https://drive.uc.cn/s/e1ebe95d144c4",
			wantID:   "e1ebe95d144c4",
			wantType: UCStr,
		},
		{
			name:     "迅雷网盘",
			url:      "https://pan.xunlei.com/s/abc123-def456",
			wantID:   "abc123-def456",
			wantType: XunleiStr,
		},
		{
			name:     "天翼云盘 t/ 格式",
			url:      "https://cloud.189.cn/t/viy2quQzMBne",
			wantID:   "viy2quQzMBne",
			wantType: TianyiStr,
		},
		{
			name:     "天翼云盘 web/share 格式",
			url:      "https://cloud.189.cn/web/share?code=UfUjiiFRbymq",
			wantID:   "UfUjiiFRbymq",
			wantType: TianyiStr,
		},
		{
			name:     "115网盘",
			url:      "https://115.com/s/swhsaua36a1",
			wantID:   "swhsaua36a1",
			wantType: Pan115Str,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotType := extractShareID(tt.url)
			if gotID != tt.wantID {
				t.Errorf("extractShareID() ID = %v, want %v", gotID, tt.wantID)
			}
			if gotType != tt.wantType {
				t.Errorf("extractShareID() Type = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}
