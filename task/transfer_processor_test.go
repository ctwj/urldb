package task

import "testing"

func TestTransferProcessor_isValidURL(t *testing.T) {
	tp := &TransferProcessor{} // isValidURL has no deps; nil repoMgr is fine

	tests := []struct {
		name string
		url  string
		want bool
	}{
		// existing platforms (regression)
		{"quark", "https://pan.quark.cn/s/abc123", true},
		{"xunlei", "https://pan.xunlei.com/s/abcdef", true},
		{"invalid random", "https://example.com/s/x", false},
		// baidu (new)
		{"baidu /s/ short", "https://pan.baidu.com/s/1abc_def-ghi", true},
		{"baidu /s/ with pwd", "https://pan.baidu.com/s/1abcdefg?pwd=1234", true},
		{"baidu http /s/", "http://pan.baidu.com/s/1abc123XYZ", true},
		{"baidu share/init", "https://pan.baidu.com/share/init?surl=abcdefg", true},
		{"baidu share/init with pwd", "https://pan.baidu.com/share/init?surl=abc_def&pwd=wxyz", true},
		{"not baidu (bare host)", "https://baidu.com/s/1abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tp.isValidURL(tt.url)
			if got != tt.want {
				t.Fatalf("isValidURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}
