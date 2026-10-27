package scheduler

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"
)

// apiReminderDays 到期提醒档位：到期前 5/2/1 天各提醒一次（FR-018）
var apiReminderDays = []int{5, 2, 1}

// ApiExpirationScheduler API 凭证到期提醒调度器（016-api-access-application，FR-018）
// 每日一次，按自然日实时计算，无需落库去重标记——每日任务固定 24h 间隔保证同一自然日
// 只跑一轮，每个档位恰好命中一次：
//   - 剩余自然日恰为 5/2/1 天 → 到期提醒邮件
//   - 已过期且过期后第一个自然日（到期次日）→ 过期通知邮件
//
// 发送失败当日不重试（次日已不在档位）；同日手动触发 RunOnceManual 会重复发送同一档。
// 永久/已停用凭证不提醒。
type ApiExpirationScheduler struct {
	*BaseScheduler
	running bool
	mutex   sync.Mutex
}

// NewApiExpirationScheduler 创建 API 凭证到期提醒调度器
func NewApiExpirationScheduler(base *BaseScheduler) *ApiExpirationScheduler {
	return &ApiExpirationScheduler{BaseScheduler: base}
}

// Start 启动每日提醒任务
func (s *ApiExpirationScheduler) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		utils.Debug("API 凭证到期提醒任务已在运行中")
		return
	}
	s.running = true
	utils.Info("启动 API 凭证到期提醒调度任务（每日）")

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runOnce(context.Background())
			case <-s.GetStopChan():
				utils.Info("停止 API 凭证到期提醒调度任务")
				return
			}
		}
	}()
}

// RunOnceManual 手动触发一轮（运维/测试用；与当日定时任务同跑会对同档位重复发送）
func (s *ApiExpirationScheduler) RunOnceManual() {
	s.runOnce(context.Background())
}

// calendarDaysUntil 按自然日计算 from 到 to 还隔几天（同一自然日为 0，to 早于 from 时为负；两侧均为 UTC 零点，差值恒为 24h 整数倍）
// 先把 to 归一到 from 的时区，保证两侧日期在同一时区下比较
func calendarDaysUntil(from, to time.Time) int {
	to = to.In(from.Location())
	f := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	t := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	return int(t.Sub(f).Hours() / 24)
}

// runOnce 执行一轮提醒
func (s *ApiExpirationScheduler) runOnce(ctx context.Context) {
	credRepo := GetGlobalApiCredentialRepo()
	if credRepo == nil {
		utils.Warn("API 凭证到期提醒 - 凭证仓储未注入，跳过本轮")
		return
	}

	now := utils.GetCurrentTime()
	// 候选集：到期时间落在 过去2天 ~ 未来6天（到期日差5天的最迟约6×24h后到期，
	// 到期次日的最早约2×24h前到期），档位匹配由各分支按自然日判定
	list, err := credRepo.FindExpiringBetween(now.AddDate(0, 0, -2), now.AddDate(0, 0, 6))
	if err != nil {
		utils.Error("API 凭证到期提醒 - 待提醒查询失败: %v", err)
		return
	}
	if len(list) == 0 {
		return
	}

	utils.Info("API 凭证到期提醒 - 本轮候选 %d 条", len(list))
	sent := 0
	for _, cred := range list {
		if cred.ExpiresAt == nil {
			continue
		}
		if cred.IsExpired(now) {
			// 过期通知：过期后第一个自然日（到期次日）发送一次
			if calendarDaysUntil(*cred.ExpiresAt, now) != 1 {
				continue
			}
			if err := services.SendApiExpiredMail(cred.UserID, *cred.ExpiresAt); err != nil {
				utils.Error("API 凭证过期通知 - 发送失败 - userID: %d, error: %v", cred.UserID, err)
				continue
			}
			sent++
			continue
		}

		daysLeft := calendarDaysUntil(now, *cred.ExpiresAt)
		if !slices.Contains(apiReminderDays, daysLeft) {
			continue
		}
		if err := services.SendApiExpirationReminderMail(cred.UserID, *cred.ExpiresAt, daysLeft); err != nil {
			utils.Error("API 凭证到期提醒 - 发送失败 - userID: %d, daysLeft: %d, error: %v", cred.UserID, daysLeft, err)
			continue
		}
		sent++
	}
	utils.Info("API 凭证到期提醒 - 本轮完成：成功 %d / 候选 %d", sent, len(list))
}
