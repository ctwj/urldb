package pan

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"
)

// TestXLSign_Structure 验证签名结构：必须为 "1." + 32 位十六进制。
func TestXLSign_Structure(t *testing.T) {
	sign := xlSignWithTimestamp(xlProfileAndroid, "test-device-id", "1700000000000")
	if !strings.HasPrefix(sign, "1.") {
		t.Fatalf("签名应以 '1.' 开头，实际: %s", sign)
	}
	hash := strings.TrimPrefix(sign, "1.")
	if len(hash) != 32 {
		t.Fatalf("签名哈希长度应为 32，实际: %d (%s)", len(hash), hash)
	}
	if _, err := hex.DecodeString(hash); err != nil {
		t.Fatalf("签名哈希应为合法十六进制: %v", err)
	}
}

// TestXLSign_Deterministic 验证确定性：相同输入相同签名；不同输入/profile 不同签名。
func TestXLSign_Deterministic(t *testing.T) {
	s1 := xlSignWithTimestamp(xlProfileAndroid, "device-abc", "1700000000000")
	if xlSignWithTimestamp(xlProfileAndroid, "device-abc", "1700000000000") != s1 {
		t.Fatal("相同输入应产生相同签名")
	}
	if xlSignWithTimestamp(xlProfileAndroid, "device-xyz", "1700000000000") == s1 {
		t.Fatal("不同 deviceID 应产生不同签名")
	}
	if xlSignWithTimestamp(xlProfileAndroid, "device-abc", "1700000000999") == s1 {
		t.Fatal("不同 timestamp 应产生不同签名")
	}
	// 不同 profile（浏览器）应产生不同签名（不同 client_id/盐值）
	if xlSignWithTimestamp(xlProfileBrowser, "device-abc", "1700000000000") == s1 {
		t.Fatal("不同 profile 应产生不同签名")
	}
}

// TestXLSign_AlgorithmMatchesReference 验证算法与参考实现等价（固化算法防误改）。
func TestXLSign_AlgorithmMatchesReference(t *testing.T) {
	p := xlProfileAndroid
	deviceID := "925b7631473a13716b791d7f28289cad"
	timestamp := "1645241033384"
	want := p.ClientID + p.ClientVersion + p.PackageName + deviceID + timestamp
	for _, a := range p.Algorithms {
		sum := md5.Sum([]byte(want + a))
		want = hex.EncodeToString(sum[:])
	}
	want = "1." + want
	got := xlSignWithTimestamp(p, deviceID, timestamp)
	if got != want {
		t.Fatalf("签名算法与参考不一致:\n want=%s\n got =%s", want, got)
	}
}

// TestXLProfileByType 验证 profile 选择（默认 android）。
func TestXLProfileByType(t *testing.T) {
	if xlProfileByType("browser").ClientID != xlProfileBrowser.ClientID {
		t.Fatal("browser 类型应返回浏览器 profile")
	}
	if xlProfileByType("android").ClientID != xlProfileAndroid.ClientID {
		t.Fatal("android 类型应返回下载管家 profile")
	}
	if xlProfileByType("").ClientID != xlProfileAndroid.ClientID {
		t.Fatal("空类型应默认返回下载管家 profile（向后兼容）")
	}
}

// TestDeriveDeviceID 验证设备标识派生：稳定、独立、32 位 hex（R-05）。
func TestDeriveDeviceID(t *testing.T) {
	d1 := deriveDeviceID("user1", "pass1")
	if d1 != deriveDeviceID("user1", "pass1") {
		t.Fatal("相同账号应派生相同设备标识")
	}
	if len(d1) != 32 {
		t.Fatalf("设备标识应为 32 位，实际: %d", len(d1))
	}
	if _, err := hex.DecodeString(d1); err != nil {
		t.Fatalf("设备标识应为合法十六进制: %v", err)
	}
	if deriveDeviceID("user2", "pass2") == d1 {
		t.Fatal("不同账号应派生不同设备标识")
	}
}
