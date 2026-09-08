package services

import (
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// 三窗口频率限制器（016-api-access-application US4，FR-012/FR-013）
// 每用户独立三窗口固定计数：自然分钟 / 自然小时 / Asia-Shanghai 自然日（与 015 限日口径一致）。
// 进程内存态：单实例部署适用，重启清零（research.md R5 已记录接受）。

// apiWindowCounters 单用户三窗口计数
type apiWindowCounters struct {
	minuteKey int64
	minute    int64
	hourKey   int64
	hour      int64
	dayKey    int64
	day       int64
}

var (
	apiRateMutex    sync.Mutex
	apiRateCounters = map[uint]*apiWindowCounters{} // userID -> counters
)

// getApiRateLimit 读某档配置（缺省回退 30/600/3000；0=不限；负值按缺省）
func getApiRateLimit(key string, def int) int {
	v, err := repoManager.SystemConfigRepository.GetConfigInt(key)
	if err != nil || v < 0 {
		return def
	}
	return v
}

// windowKeys 计算三窗口标识（窗口起始时间戳）
func windowKeys(now time.Time) (minuteKey, hourKey, dayKey int64) {
	minuteKey = now.Truncate(time.Minute).Unix()
	hourKey = now.Truncate(time.Hour).Unix()
	// 自然日按系统时区（默认 Asia/Shanghai，见 utils.GetCurrentTime）
	loc := now.Location()
	y, m, d := now.In(loc).Date()
	dayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)
	dayKey = dayStart.Unix()
	return
}

// AllowApiRequest 判定并计数：三档均未超限则各窗口 +1 并放行；任一超限返回对应维度
// 返回（是否放行, 超限维度 minute/hour/day, 该维度上限）
func AllowApiRequest(userID uint) (bool, string, int) {
	now := utils.GetCurrentTime()
	minuteKey, hourKey, dayKey := windowKeys(now)

	limits := [3]int{
		getApiRateLimit(entity.ConfigKeyApiRateLimitMinute, 30),
		getApiRateLimit(entity.ConfigKeyApiRateLimitHour, 600),
		getApiRateLimit(entity.ConfigKeyApiRateLimitDay, 3000),
	}

	apiRateMutex.Lock()
	defer apiRateMutex.Unlock()

	c, ok := apiRateCounters[userID]
	if !ok {
		c = &apiWindowCounters{}
		apiRateCounters[userID] = c
	}

	// 窗口切换清零
	if c.minuteKey != minuteKey {
		c.minuteKey, c.minute = minuteKey, 0
	}
	if c.hourKey != hourKey {
		c.hourKey, c.hour = hourKey, 0
	}
	if c.dayKey != dayKey {
		c.dayKey, c.day = dayKey, 0
	}

	// 依次判定：任一超限即拒绝（FR-013）
	if limits[0] > 0 && c.minute >= int64(limits[0]) {
		return false, "minute", limits[0]
	}
	if limits[1] > 0 && c.hour >= int64(limits[1]) {
		return false, "hour", limits[1]
	}
	if limits[2] > 0 && c.day >= int64(limits[2]) {
		return false, "day", limits[2]
	}

	// 放行并计数
	c.minute++
	c.hour++
	c.day++
	return true, "", 0
}
