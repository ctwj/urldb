package services

import (
	"testing"
	"time"

	"github.com/ctwj/urldb/db/entity"
)

// TestCredentialRenewalExpiry 续期到期时间计算：未过期自当前到期时间顺延（叠加不损失剩余时长），
// 已过期/无到期时间自 now 起算（月份加法沿用 time.AddDate 的自然归一，如 11-30 +3月 → 次年 03-02）
func TestCredentialRenewalExpiry(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Date(2026, 10, 27, 10, 0, 0, 0, loc)
	expAt := time.Date(2026, 11, 5, 0, 0, 0, 0, loc)

	cases := []struct {
		name   string
		cred   *entity.ApiCredential
		months int
		want   time.Time
	}{
		{"未过期自当前到期时间顺延", &entity.ApiCredential{ExpiresAt: &expAt}, 3,
			time.Date(2027, 2, 5, 0, 0, 0, 0, loc)},
		{"未过期可多次叠加", &entity.ApiCredential{ExpiresAt: &expAt}, 6,
			time.Date(2027, 5, 5, 0, 0, 0, 0, loc)},
		{"已过期自当前时间起算", &entity.ApiCredential{ExpiresAt: &time.Time{}}, 1,
			time.Date(2026, 11, 27, 10, 0, 0, 0, loc)},
		{"无到期时间(永久,防御)自当前时间起算", &entity.ApiCredential{}, 2,
			time.Date(2026, 12, 27, 10, 0, 0, 0, loc)},
		{"月末溢出按自然归一", &entity.ApiCredential{ExpiresAt: ptrTime(time.Date(2026, 11, 30, 0, 0, 0, 0, loc))}, 3,
			time.Date(2027, 3, 2, 0, 0, 0, 0, loc)},
	}
	for _, c := range cases {
		got := credentialRenewalExpiry(c.cred, c.months, now)
		if !got.Equal(c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func ptrTime(t time.Time) *time.Time { return &t }
