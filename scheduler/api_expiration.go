package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"
)

// ApiExpirationScheduler API 凭证到期提醒调度器（016-api-access-application，FR-018）
// 每日一次：到期时间落在未来 6~7 天窗口、active、未发过提醒的凭证 → 发提醒邮件并打标记。
// 发送失败不打标记，次日窗口仍在，天然重试；永久/已停用/已过期凭证不提醒。
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

// RunOnceManual 手动触发一轮（运维/测试用）
func (s *ApiExpirationScheduler) RunOnceManual() {
	s.runOnce(context.Background())
}

// runOnce 执行一轮提醒
func (s *ApiExpirationScheduler) runOnce(ctx context.Context) {
	credRepo := GetGlobalApiCredentialRepo()
	if credRepo == nil {
		utils.Warn("API 凭证到期提醒 - 凭证仓储未注入，跳过本轮")
		return
	}

	now := utils.GetCurrentTime()
	// 提醒窗口：到期前 7 天（取 [now+6d, now+7d]，每日任务下每套凭证恰好命中一次窗口）
	start := now.AddDate(0, 0, 6)
	end := now.AddDate(0, 0, 7)

	list, err := credRepo.FindExpiringBetween(start, end)
	if err != nil {
		utils.Error("API 凭证到期提醒 - 待提醒查询失败: %v", err)
		return
	}
	if len(list) == 0 {
		return
	}

	utils.Info("API 凭证到期提醒 - 本轮待提醒 %d 条", len(list))
	sent := 0
	for _, cred := range list {
		if cred.ExpiresAt == nil {
			continue
		}
		ok, err := services.SendApiExpirationReminderMail(cred.UserID, *cred.ExpiresAt)
		if err != nil {
			utils.Error("API 凭证到期提醒 - 发送失败 - userID: %d, error: %v", cred.UserID, err)
			continue
		}
		if !ok {
			// 用户不存在或未填邮箱：打标记避免每日空转
			nowMark := utils.GetCurrentTime()
			cred.ReminderSentAt = &nowMark
			if err := credRepo.Update(cred); err != nil {
				utils.Error("API 凭证到期提醒 - 标记写入失败 - userID: %d, error: %v", cred.UserID, err)
			}
			continue
		}

		nowMark := utils.GetCurrentTime()
		cred.ReminderSentAt = &nowMark
		if err := credRepo.Update(cred); err != nil {
			utils.Error("API 凭证到期提醒 - 标记写入失败 - userID: %d, error: %v", cred.UserID, err)
			continue
		}
		sent++
	}
	utils.Info("API 凭证到期提醒 - 本轮完成：成功 %d / 共 %d", sent, len(list))
}
