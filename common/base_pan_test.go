package pan

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasePanService_GzipResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, _ = gz.Write([]byte(`{"errno":0,"result":{"bdstoken":"abc"}}`))
		_ = gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(w, &buf)
	}))
	defer srv.Close()

	svc := NewBasePanService(nil)
	// 模拟百度：手动写 Accept-Encoding，触发 Go transport 不自动解压的路径
	svc.SetHeader("Accept-Encoding", "gzip, deflate, br")

	data, err := svc.HTTPGet(srv.URL, nil)
	if err != nil {
		t.Fatalf("HTTPGet error: %v", err)
	}
	if !bytes.Contains(data, []byte(`"bdstoken":"abc"`)) {
		t.Fatalf("gzip 解压失败，got: %s", string(data))
	}
}

func TestBasePanService_HTTPPostForm(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		_, _ = io.Copy(buf, r.Body)
		gotBody = buf.String()
		_, _ = io.WriteString(w, `{"errno":0}`)
	}))
	defer srv.Close()

	svc := NewBasePanService(nil)
	data, err := svc.HTTPPostForm(srv.URL, "pwd=1234&vcode=", nil)
	if err != nil {
		t.Fatalf("HTTPPostForm error: %v", err)
	}
	if gotBody != "pwd=1234&vcode=" {
		t.Fatalf("body sent = %q, want raw form string (not JSON-wrapped)", gotBody)
	}
	if !bytes.Contains(data, []byte(`"errno":0`)) {
		t.Fatalf("unexpected response: %s", string(data))
	}
}
