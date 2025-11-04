package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// ConcurrencyController 并发控制器
type ConcurrencyController struct {
	// 每个插件的最大并发数
	pluginLimits map[string]int
	// 当前每个插件的活跃任务数
	pluginActiveTasks map[string]int
	// 全局最大并发数
	globalLimit int
	// 全局活跃任务数
	globalActiveTasks int
	// 等待队列
	waitingTasks map[string][]*WaitingTask
	// 互斥锁
	mutex sync.Mutex
	// 条件变量用于任务等待和通知
	cond *sync.Cond
}

// WaitingTask 等待执行的任务
type WaitingTask struct {
	PluginName string
	TaskFunc   func() error
	Ctx        context.Context
	Cancel     context.CancelFunc
	ResultChan chan error
}

// NewConcurrencyController 创建新的并发控制器
func NewConcurrencyController(globalLimit int) *ConcurrencyController {
	cc := &ConcurrencyController{
		pluginLimits:      make(map[string]int),
		pluginActiveTasks: make(map[string]int),
		globalLimit:       globalLimit,
		waitingTasks:      make(map[string][]*WaitingTask),
	}

	cc.cond = sync.NewCond(&cc.mutex)
	return cc
}

// SetPluginLimit 设置插件的并发限制
func (cc *ConcurrencyController) SetPluginLimit(pluginName string, limit int) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	cc.pluginLimits[pluginName] = limit
	utils.Info("Set concurrency limit for plugin %s to %d", pluginName, limit)
}

// GetPluginLimit 获取插件的并发限制
func (cc *ConcurrencyController) GetPluginLimit(pluginName string) int {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	return cc.pluginLimits[pluginName]
}

// Execute 执行受并发控制的任务
func (cc *ConcurrencyController) Execute(ctx context.Context, pluginName string, taskFunc func() error) error {
	// 创建结果通道
	resultChan := make(chan error, 1)

	// 尝试获取执行许可
	permitted := cc.tryAcquire(pluginName)
	if !permitted {
		// 需要等待
		waitingTask := &WaitingTask{
			PluginName: pluginName,
			TaskFunc:   taskFunc,
			Ctx:        ctx,
			ResultChan: resultChan,
		}

		// 添加到等待队列
		cc.mutex.Lock()
		cc.waitingTasks[pluginName] = append(cc.waitingTasks[pluginName], waitingTask)
		cc.mutex.Unlock()

		utils.Debug("Task for plugin %s added to waiting queue", pluginName)

		// 等待结果或上下文取消
		select {
		case err := <-resultChan:
			return err
		case <-ctx.Done():
			// 从等待队列中移除任务
			cc.removeFromWaitingQueue(pluginName, waitingTask)
			return ctx.Err()
		}
	}

	// 可以立即执行
	defer cc.release(pluginName)

	// 在goroutine中执行任务，以支持超时和取消
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	taskResult := make(chan error, 1)
	go func() {
		defer close(taskResult)
		taskResult <- taskFunc()
	}()

	select {
	case err := <-taskResult:
		return err
	case <-taskCtx.Done():
		return taskCtx.Err()
	}
}

// tryAcquire 尝试获取执行许可
func (cc *ConcurrencyController) tryAcquire(pluginName string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// 检查全局限制
	if cc.globalActiveTasks >= cc.globalLimit {
		return false
	}

	// 检查插件限制
	pluginLimit := cc.pluginLimits[pluginName]
	if pluginLimit <= 0 {
		// 如果没有设置限制，使用默认值（全局限制的1/4，但至少为1）
		pluginLimit = cc.globalLimit / 4
		if pluginLimit < 1 {
			pluginLimit = 1
		}
	}

	if cc.pluginActiveTasks[pluginName] >= pluginLimit {
		return false
	}

	// 增加计数
	cc.globalActiveTasks++
	cc.pluginActiveTasks[pluginName]++
	return true
}

// release 释放执行许可
func (cc *ConcurrencyController) release(pluginName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// 减少计数
	cc.globalActiveTasks--
	if cc.globalActiveTasks < 0 {
		cc.globalActiveTasks = 0
	}

	cc.pluginActiveTasks[pluginName]--
	if cc.pluginActiveTasks[pluginName] < 0 {
		cc.pluginActiveTasks[pluginName] = 0
	}

	// 检查是否有等待的任务可以执行
	cc.checkWaitingTasks()

	// 通知等待的goroutine
	cc.cond.Broadcast()
}

// checkWaitingTasks 检查等待队列中的任务是否可以执行
func (cc *ConcurrencyController) checkWaitingTasks() {
	for pluginName, tasks := range cc.waitingTasks {
		if len(tasks) > 0 {
			// 检查是否可以获得执行许可
			if cc.tryAcquire(pluginName) {
				// 获取第一个等待的任务
				task := tasks[0]
				cc.waitingTasks[pluginName] = tasks[1:]

				// 在新的goroutine中执行任务
				go cc.executeWaitingTask(task, pluginName)
			}
		}
	}
}

// executeWaitingTask 执行等待的任务
func (cc *ConcurrencyController) executeWaitingTask(task *WaitingTask, pluginName string) {
	defer cc.release(pluginName)

	// 检查上下文是否已取消
	select {
	case <-task.Ctx.Done():
		task.ResultChan <- task.Ctx.Err()
		return
	default:
	}

	// 执行任务
	err := task.TaskFunc()

	// 发送结果
	select {
	case task.ResultChan <- err:
	default:
		// 结果通道已关闭或已满
		utils.Warn("Failed to send result for waiting task of plugin %s", pluginName)
	}
}

// removeFromWaitingQueue 从等待队列中移除任务
func (cc *ConcurrencyController) removeFromWaitingQueue(pluginName string, task *WaitingTask) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if tasks, exists := cc.waitingTasks[pluginName]; exists {
		for i, t := range tasks {
			if t == task {
				// 从队列中移除
				cc.waitingTasks[pluginName] = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
	}
}

// GetStats 获取并发控制器的统计信息
func (cc *ConcurrencyController) GetStats() map[string]interface{} {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	stats := make(map[string]interface{})
	stats["global_limit"] = cc.globalLimit
	stats["global_active"] = cc.globalActiveTasks

	pluginStats := make(map[string]interface{})
	for pluginName, limit := range cc.pluginLimits {
		pluginStat := make(map[string]interface{})
		pluginStat["limit"] = limit
		pluginStat["active"] = cc.pluginActiveTasks[pluginName]
		if tasks, exists := cc.waitingTasks[pluginName]; exists {
			pluginStat["waiting"] = len(tasks)
		} else {
			pluginStat["waiting"] = 0
		}
		pluginStats[pluginName] = pluginStat
	}
	stats["plugins"] = pluginStats

	return stats
}

// WaitForAvailable 等待直到有可用的并发槽位
func (cc *ConcurrencyController) WaitForAvailable(ctx context.Context, pluginName string) error {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	for {
		// 检查是否有可用的槽位
		pluginLimit := cc.pluginLimits[pluginName]
		if pluginLimit <= 0 {
			// 如果没有设置限制，使用默认值
			pluginLimit = cc.globalLimit / 4
			if pluginLimit < 1 {
				pluginLimit = 1
			}
		}

		if cc.globalActiveTasks < cc.globalLimit && cc.pluginActiveTasks[pluginName] < pluginLimit {
			return nil // 有可用槽位
		}

		// 等待或超时
		waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// 使用goroutine来处理条件等待
		waitDone := make(chan struct{})
		go func() {
			defer close(waitDone)
			cc.cond.Wait()
		}()

		select {
		case <-waitDone:
			// 条件满足，继续检查
		case <-waitCtx.Done():
			// 超时或取消
			return fmt.Errorf("timeout waiting for available slot: %v", waitCtx.Err())
		}
	}
}