package utils

import "testing"

// 测试向量来自《微信公众号动态口令对接说明.md》第四节，用于自检算法实现正确性
func TestGenerateDynamicCode(t *testing.T) {
	const secret = "moss-wxpass-8f4c2a1e9b7d"
	cases := []struct {
		unixSeconds int64
		expected    string
	}{
		{1700000000, "534579"},
		{1700001800, "079502"},
		{1735689600, "428950"},
	}

	for _, c := range cases {
		if got := GenerateDynamicCode(secret, c.unixSeconds, DynamicCodePeriod); got != c.expected {
			t.Errorf("GenerateDynamicCode(%q, %d, %d) = %s, 期望 %s", secret, c.unixSeconds, DynamicCodePeriod, got, c.expected)
		}
	}
}
