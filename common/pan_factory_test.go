package pan

import (
	"testing"
)

func TestExtractServiceType_123pan(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want ServiceType
	}{
		{
			name: "传统 123pan.com 域名",
			url:  "https://www.123pan.com/s/i4uaTd-WHn0",
			want: Pan123,
		},
		{
			name: "share.123pan.cn 子域名",
			url:  "https://1856557151.share.123pan.cn/123pan/oJqrvd-KlG9ds",
			want: Pan123,
		},
		{
			name: "123912.com 域名",
			url:  "https://www.123912.com/s/U8f2Td-ZeOX",
			want: Pan123,
		},
		{
			name: "123684.com 域名",
			url:  "https://www.123684.com/s/u9izjv-k3uWv",
			want: Pan123,
		},
		{
			name: "123865.com 域名",
			url:  "https://www.123865.com/s/test123",
			want: Pan123,
		},
		{
			name: "123685.com 域名",
			url:  "https://www.123685.com/s/test456",
			want: Pan123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractServiceType(tt.url)
			if got != tt.want {
				t.Errorf("ExtractServiceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractShareId_123pan(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		wantID string
		want   ServiceType
	}{
		{
			name:   "传统 /s/ 格式",
			url:    "https://www.123pan.com/s/i4uaTd-WHn0",
			wantID: "i4uaTd-WHn0",
			want:   Pan123,
		},
		{
			name:   "新 /123pan/ 格式",
			url:    "https://1856557151.share.123pan.cn/123pan/oJqrvd-KlG9ds",
			wantID: "oJqrvd-KlG9ds",
			want:   Pan123,
		},
		{
			name:   "带查询参数",
			url:    "https://1856557151.share.123pan.cn/123pan/oJqrvd-KlG9ds?entry=2",
			wantID: "oJqrvd-KlG9ds",
			want:   Pan123,
		},
		{
			name:   "带锚点",
			url:    "https://www.123pan.com/s/i4uaTd-WHn0#folder",
			wantID: "i4uaTd-WHn0",
			want:   Pan123,
		},
		{
			name:   "带 entry 参数",
			url:    "https://www.123pan.com/s/i4uaTd-WHn0?entry=123",
			wantID: "i4uaTd-WHn0",
			want:   Pan123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, got := ExtractShareId(tt.url)
			if gotID != tt.wantID {
				t.Errorf("ExtractShareId() ID = %v, want %v", gotID, tt.wantID)
			}
			if got != tt.want {
				t.Errorf("ExtractShareId() Type = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractServiceType_AllPans(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want ServiceType
	}{
		{"夸克网盘", "https://pan.quark.cn/s/abc123", Quark},
		{"阿里云盘 alipan", "https://www.alipan.com/s/abc123", Alipan},
		{"阿里云盘 aliyundrive", "https://www.aliyundrive.com/s/abc123", Alipan},
		{"百度网盘", "https://pan.baidu.com/s/abc123", BaiduPan},
		{"UC网盘", "https://drive.uc.cn/s/abc123", UC},
		{"迅雷网盘", "https://pan.xunlei.com/s/abc123", Xunlei},
		{"天翼云盘", "https://cloud.189.cn/t/abc123", Tianyi},
		{"115网盘 115.com", "https://115.com/s/abc123", Pan115},
		{"115网盘 115cdn.com", "https://115cdn.com/s/abc123", Pan115},
		{"115网盘 anxia.com", "https://anxia.com/s/abc123", Pan115},
		{"不支持的链接", "https://example.com/s/abc123", NotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractServiceType(tt.url)
			if got != tt.want {
				t.Errorf("ExtractServiceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceType_String(t *testing.T) {
	tests := []struct {
		st   ServiceType
		want string
	}{
		{Quark, "quark"},
		{Alipan, "alipan"},
		{BaiduPan, "baidu"},
		{UC, "uc"},
		{Xunlei, "xunlei"},
		{Tianyi, "tianyi"},
		{Pan123, "123pan"},
		{Pan115, "115"},
		{NotFound, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.st.String(); got != tt.want {
				t.Errorf("ServiceType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
