package scheduler

import (
	"fmt"
	"sync"
	"time"

	panutils "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/utils"
)

// XunleiKeepaliveScheduler 迅雷 token 保活调度器。
//
// 迅雷没有"定时续期"机制，refresh_token 长期不被使用会自然过期（约 30 天），
// 且安卓方案无账号密码兜底，过期后无法自动恢复。本任务定期遍历所有有效迅雷账号，
// 调用 Keepalive 刷新 access_token（顺带轮转 refresh_token），给 refresh_token 续命，
// 避免闲置账号失效。
type XunleiKeepaliveScheduler struct {
	*BaseScheduler
	running bool
	mutex   sync.Mutex // 防止保活任务重叠执行
}

// NewXunleiKeepaliveScheduler 创建迅雷 token 保活调度器
func NewXunleiKeepaliveScheduler(base *BaseScheduler) *XunleiKeepaliveScheduler {
	return &XunleiKeepaliveScheduler{BaseScheduler: base}
}

// Start 启动迅雷 token 保活定时任务
func (s *XunleiKeepaliveScheduler) Start() {
	if s.running {
		utils.Debug("迅雷 token 保活任务已在运行中")
		return
	}
	s.running = true
	utils.Info("启动迅雷 token 保活定时任务")

	go func() {
		interval := 24 * time.Hour // refresh_token 有效期约 30 天，每天刷新一次足够续命
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		utils.Info(fmt.Sprintf("迅雷 token 保活任务已启动，间隔: %v", interval))

		// 启动后立即执行一次（续命 + 尽早暴露失效账号）
		s.keepalive()

		for {
			select {
			case <-ticker.C:
				if s.mutex.TryLock() {
					go func() {
						defer s.mutex.Unlock()
						s.keepalive()
					}()
				} else {
					utils.Debug("上一次迅雷 token 保活任务还在执行中，跳过本次")
				}
			case <-s.GetStopChan():
				utils.Info("停止迅雷 token 保活定时任务")
				return
			}
		}
	}()
}

// Stop 停止迅雷 token 保活定时任务
func (s *XunleiKeepaliveScheduler) Stop() {
	if !s.running {
		utils.Debug("迅雷 token 保活任务未在运行")
		return
	}
	s.GetStopChan() <- true
	s.running = false
	utils.Info("已发送停止信号给迅雷 token 保活任务")
}

// IsRunning 检查迅雷 token 保活任务是否正在运行
func (s *XunleiKeepaliveScheduler) IsRunning() bool {
	return s.running
}

// keepalive 遍历所有有效迅雷账号，刷新 token 续命
func (s *XunleiKeepaliveScheduler) keepalive() {
	utils.Debug("[迅雷保活] 开始刷新迅雷账号 token...")

	// 定位 xunlei 平台 ID
	pans, err := s.panRepo.FindAll()
	if err != nil {
		utils.Error(fmt.Sprintf("[迅雷保活] 获取平台列表失败: %v", err))
		return
	}
	var xunleiPanID uint
	found := false
	for _, p := range pans {
		if p.Name == "xunlei" {
			xunleiPanID = p.ID
			found = true
			break
		}
	}
	if !found {
		utils.Debug("[迅雷保活] 未找到 xunlei 平台，跳过")
		return
	}

	accounts, err := s.cksRepo.FindByPanID(xunleiPanID)
	if err != nil {
		utils.Error(fmt.Sprintf("[迅雷保活] 获取迅雷账号失败: %v", err))
		return
	}
	if len(accounts) == 0 {
		utils.Debug("[迅雷保活] 没有迅雷账号，跳过")
		return
	}

	factory := panutils.GetInstance()
	successCnt, failCnt, skipCnt := 0, 0, 0
	for i := range accounts {
		acc := accounts[i]
		if !acc.IsValid {
			skipCnt++
			continue
		}

		service, err := factory.CreatePanServiceByType(panutils.Xunlei, &panutils.PanConfig{})
		if err != nil {
			utils.Error(fmt.Sprintf("[迅雷保活] 账号 %d 创建服务失败: %v", acc.ID, err))
			failCnt++
			continue
		}

		xunlei := service.(*panutils.XunleiPanService)
		xunlei.SetCKSRepository(s.cksRepo, acc)
		if err := xunlei.Keepalive(); err != nil {
			utils.Error(fmt.Sprintf("[迅雷保活] 账号 %d (%s) 刷新失败: %v", acc.ID, acc.Username, err))
			failCnt++
			continue
		}
		successCnt++
		utils.Debug(fmt.Sprintf("[迅雷保活] 账号 %d (%s) 刷新成功", acc.ID, acc.Username))
	}

	utils.Info(fmt.Sprintf("[迅雷保活] 完成：成功 %d，失败 %d，跳过(无效) %d", successCnt, failCnt, skipCnt))
}
