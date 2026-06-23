package pan

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"
)

// TestXLCaptchaSign_Structure 验证签名结构：必须为 "1." + 32 位十六进制。
func TestXLCaptchaSign_Structure(t *testing.T) {
	timestamp, sign := xlCaptchaSign("test-device-id")
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
	if timestamp == "" {
		t.Fatal("时间戳不应为空")
	}
}

// TestXLCaptchaSign_Deterministic 验证确定性：相同输入产生相同签名；不同输入产生不同签名。
func TestXLCaptchaSign_Deterministic(t *testing.T) {
	s1 := xlCaptchaSignWithTimestamp("device-abc", "1700000000000")
	s2 := xlCaptchaSignWithTimestamp("device-abc", "1700000000000")
	if s1 != s2 {
		t.Fatalf("相同输入应产生相同签名: %s vs %s", s1, s2)
	}
	if xlCaptchaSignWithTimestamp("device-xyz", "1700000000000") == s1 {
		t.Fatal("不同 deviceID 应产生不同签名")
	}
	if xlCaptchaSignWithTimestamp("device-abc", "1700000000999") == s1 {
		t.Fatal("不同 timestamp 应产生不同签名")
	}
}

// TestXLCaptchaSign_AlgorithmMatchesReference 验证算法与参考实现（OpenList thunder）等价。
// 手动按相同拼接顺序与盐值顺序做多轮 MD5，与函数输出比对，固化算法防止未来误改。
func TestXLCaptchaSign_AlgorithmMatchesReference(t *testing.T) {
	deviceID := "925b7631473a13716b791d7f28289cad"
	timestamp := "1645241033384"

	want := XLClientID + XLClientVersion + XLPackageName + deviceID + timestamp
	for _, a := range Algorithms {
		sum := md5.Sum([]byte(want + a))
		want = hex.EncodeToString(sum[:])
	}
	want = "1." + want

	got := xlCaptchaSignWithTimestamp(deviceID, timestamp)
	if got != want {
		t.Fatalf("签名算法与参考不一致:\n want=%s\n got =%s", want, got)
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
