package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

// DynamicCodePeriod 微信公众号动态口令时间窗（秒），与《微信公众号动态口令对接说明》一致：30 分钟
const DynamicCodePeriod int64 = 1800

// GenerateDynamicCode 按共享密钥与当前时间计算当期动态口令（TOTP 离线同步方案）。
// 算法规格见《微信公众号动态口令对接说明.md》：HMAC-SHA1（RFC 4226/6238 截断），
// 6 位数字（不足补 0），密钥为共享密钥字符串的 UTF-8 字节（不做 Base32 转换），
// counter = unix时间戳/period 向下取整，输入为 8 字节大端 counter。
func GenerateDynamicCode(secret string, unixSeconds, period int64) string {
	if period <= 0 {
		period = DynamicCodePeriod
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(unixSeconds/period))

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(buf)
	sum := mac.Sum(nil)

	off := sum[len(sum)-1] & 0x0f
	val := (uint32(sum[off])&0x7f)<<24 | uint32(sum[off+1])<<16 | uint32(sum[off+2])<<8 | uint32(sum[off+3])
	return fmt.Sprintf("%06d", val%1000000)
}
