package utils

import (
	"testing"
)

func TestExtractServiceType(t *testing.T) {
	tests := []struct {
		url         string
		expected    ServiceType
		description string
	}{
		{
			url:         "https://pan.quark.cn/s/123456",
			expected:    Quark,
			description: "夸克网盘",
		},
		{
			url:         "https://www.alipan.com/s/123456",
			expected:    Alipan,
			description: "阿里云盘",
		},
		{
			url:         "https://www.aliyundrive.com/s/123456",
			expected:    Alipan,
			description: "阿里云盘别名",
		},
		{
			url:         "https://pan.baidu.com/s/123456",
			expected:    BaiduPan,
			description: "百度网盘",
		},
		{
			url:         "https://drive.uc.cn/s/123456",
			expected:    UC,
			description: "UC网盘",
		},
		{
			url:         "https://fast.uc.cn/s/123456",
			expected:    UC,
			description: "UC网盘别名",
		},
		{
			url:         "https://example.com/s/123456",
			expected:    NotFound,
			description: "不支持的链接",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := ExtractServiceType(test.url)
			if result != test.expected {
				t.Errorf("期望 %v, 实际 %v", test.expected, result)
			}
		})
	}
}

func TestExtractShareId(t *testing.T) {
	tests := []struct {
		url          string
		expectedID   string
		expectedType ServiceType
		description  string
	}{
		{
			url:          "https://pan.quark.cn/s/123456",
			expectedID:   "123456",
			expectedType: Quark,
			description:  "夸克网盘",
		},
		{
			url:          "https://www.alipan.com/s/123456?entry=abc",
			expectedID:   "123456",
			expectedType: Alipan,
			description:  "阿里云盘带参数",
		},
		{
			url:          "https://pan.baidu.com/s/123456#section",
			expectedID:   "123456",
			expectedType: BaiduPan,
			description:  "百度网盘带锚点",
		},
		{
			url:          "https://example.com/s/123456",
			expectedID:   "123456",
			expectedType: NotFound,
			description:  "不支持的链接",
		},
		{
			url:          "https://pan.quark.cn/other/123456",
			expectedID:   "",
			expectedType: NotFound,
			description:  "无效格式",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			shareID, serviceType := ExtractShareId(test.url)
			if shareID != test.expectedID {
				t.Errorf("期望分享ID %s, 实际 %s", test.expectedID, shareID)
			}
			if serviceType != test.expectedType {
				t.Errorf("期望服务类型 %v, 实际 %v", test.expectedType, serviceType)
			}
		})
	}
}

func TestPanFactory(t *testing.T) {
	factory := NewPanFactory()

	tests := []struct {
		url         string
		config      *PanConfig
		shouldError bool
		description string
	}{
		{
			url: "https://pan.quark.cn/s/123456",
			config: &PanConfig{
				URL:         "https://pan.quark.cn/s/123456",
				IsType:      0,
				ExpiredType: 1,
			},
			shouldError: false,
			description: "夸克网盘",
		},
		{
			url: "https://www.alipan.com/s/123456",
			config: &PanConfig{
				URL:         "https://www.alipan.com/s/123456",
				IsType:      0,
				ExpiredType: 1,
			},
			shouldError: false,
			description: "阿里云盘",
		},
		{
			url: "https://example.com/s/123456",
			config: &PanConfig{
				URL:         "https://example.com/s/123456",
				IsType:      0,
				ExpiredType: 1,
			},
			shouldError: true,
			description: "不支持的链接",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			service, err := factory.CreatePanService(test.url, test.config)

			if test.shouldError {
				if err == nil {
					t.Error("期望错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误，但得到错误: %v", err)
				return
			}

			if service == nil {
				t.Error("服务不应该为nil")
				return
			}

			// 测试服务类型
			serviceType := service.GetServiceType()
			expectedType := ExtractServiceType(test.url)
			if serviceType != expectedType {
				t.Errorf("期望服务类型 %v, 实际 %v", expectedType, serviceType)
			}
		})
	}
}

func TestServiceTypeString(t *testing.T) {
	tests := []struct {
		serviceType ServiceType
		expected    string
	}{
		{Quark, "quark"},
		{Alipan, "alipan"},
		{BaiduPan, "baidu"},
		{UC, "uc"},
		{NotFound, "unknown"},
		{ServiceType(999), "unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.serviceType.String()
			if result != test.expected {
				t.Errorf("期望 %s, 实际 %s", test.expected, result)
			}
		})
	}
}
