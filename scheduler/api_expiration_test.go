package scheduler

import (
	"testing"
	"time"
)

// TestCalendarDaysUntil 自然日剩余天数计算（含时区归一与午夜边界）
func TestCalendarDaysUntil(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Date(2026, 10, 27, 23, 30, 0, 0, loc)
	cases := []struct {
		name string
		exp  time.Time
		want int
	}{
		{"常规5天", time.Date(2026, 11, 1, 10, 0, 0, 0, loc), 5},
		{"到期在当天稍晚也算整日", time.Date(2026, 11, 1, 0, 30, 0, 0, loc), 5},
		{"2天档", time.Date(2026, 10, 29, 8, 0, 0, 0, loc), 2},
		{"1天档", time.Date(2026, 10, 28, 23, 59, 0, 0, loc), 1},
		{"到期当天为0", time.Date(2026, 10, 27, 23, 59, 0, 0, loc), 0},
		// 到期时间以 UTC 返回时先归一到应用时区：UTC 11-01 17:00 = 上海 11-02 01:00
		{"UTC午夜后归一时区", time.Date(2026, 11, 1, 17, 0, 0, 0, time.UTC), 6},
		{"UTC午夜前归一时区", time.Date(2026, 11, 1, 5, 0, 0, 0, time.UTC), 5},
		// 过期方向：到期日在昨天 → 过期后第 1 个自然日
		// 到期日在昨天：now→exp 方向为负（调度器过期分支用 exp→now 方向，见下方单独断言）
		{"到期日在昨天", time.Date(2026, 10, 26, 18, 0, 0, 0, loc), -1},
	}
	for _, c := range cases {
		if got := calendarDaysUntil(now, c.exp); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

// 过期分支实际方向：calendarDaysUntil(到期时间, now) == 1 即过期后第一个自然日
func TestCalendarDaysUntilExpiredDirection(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Date(2026, 10, 27, 8, 0, 0, 0, loc)
	expiredYesterday := time.Date(2026, 10, 26, 18, 0, 0, 0, loc)
	if got := calendarDaysUntil(expiredYesterday, now); got != 1 {
		t.Errorf("过期次日方向: got %d, want 1", got)
	}
	expiredTwoDaysAgo := time.Date(2026, 10, 25, 18, 0, 0, 0, loc)
	if got := calendarDaysUntil(expiredTwoDaysAgo, now); got != 2 {
		t.Errorf("过期两天方向: got %d, want 2", got)
	}
}
