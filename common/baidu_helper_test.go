package pan

import (
	"strings"
	"testing"
)

func TestExtractSurl(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"/s/ 带前导1和pwd", "https://pan.baidu.com/s/1abcDef?pwd=1234", "abcDef"},
		{"/share/init?surl= 格式", "https://pan.baidu.com/share/init?surl=abcDef", "abcDef"},
		{"/s/ 仅前导1无pwd", "https://pan.baidu.com/s/1xxx", "xxx"},
		{"/s/ 无前导1", "https://pan.baidu.com/s/abcDef", "abcDef"},
		{"非法URL", "https://example.com/whatever", ""},
		{"空字符串", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractSurl(tt.url)
			if got != tt.want {
				t.Errorf("ExtractSurl(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestParseSharePageHTML(t *testing.T) {
	page := `{"shareid":1234567890,"share_uk":"9876543","fs_id":1001,"server_filename":"电影A.mp4","isdir":0,"fs_id":1002,"server_filename":"剧集B","isdir":1}`

	files, err := ParseSharePageHTML(page)
	if err != nil {
		t.Fatalf("ParseSharePageHTML() unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("ParseSharePageHTML() got %d files, want 2", len(files))
	}
	if files[0].FsID != 1001 || files[0].ServerName != "电影A.mp4" || files[0].IsDir {
		t.Errorf("files[0] = %+v, want FsID=1001 ServerName=电影A.mp4 IsDir=false", files[0])
	}
	if files[1].FsID != 1002 || files[1].ServerName != "剧集B" || !files[1].IsDir {
		t.Errorf("files[1] = %+v, want FsID=1002 ServerName=剧集B IsDir=true", files[1])
	}
	if files[0].ShareID != 1234567890 || files[0].UK != 9876543 {
		t.Errorf("files[0] ShareID=%d UK=%d, want 1234567890/9876543", files[0].ShareID, files[0].UK)
	}
}

func TestParseSharePageHTML_Invalid(t *testing.T) {
	_, err := ParseSharePageHTML("nothing useful here")
	if err == nil {
		t.Fatal("ParseSharePageHTML() expected error for invalid page, got nil")
	}
	if !strings.Contains(err.Error(), "失效") {
		t.Errorf("error should mention 失效, got: %v", err)
	}
}

func TestErrorCodeMap(t *testing.T) {
	cases := map[int]string{
		0:      "转存成功",
		-6:     "无痕模式获取 Cookie",
		-10:    "容量不足",
		20:     "容量不足",
		4:      "同名文件",
		200025: "提取码输入错误",
	}
	for code, want := range cases {
		got, ok := ErrorCodeMap[code]
		if !ok {
			t.Errorf("ErrorCodeMap[%d] missing", code)
			continue
		}
		if !strings.Contains(got, want) {
			t.Errorf("ErrorCodeMap[%d] = %q, want substring %q", code, got, want)
		}
	}
}

func TestErrnoMessage(t *testing.T) {
	// 已知码返回映射消息
	if got := ErrnoMessage(-10); !strings.Contains(got, "容量不足") {
		t.Errorf("ErrnoMessage(-10) = %q, want 容量不足", got)
	}
	// 未知码返回兜底文案
	if got := ErrnoMessage(999999); !strings.Contains(got, "未知错误") {
		t.Errorf("ErrnoMessage(999999) = %q, want 未知错误", got)
	}
}
