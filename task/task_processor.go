package task

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// TaskProcessor 任务处理器接口
type TaskProcessor interface {
	Process(ctx context.Context, taskID uint, item *entity.TaskItem) error
	GetTaskType() string
}

// TaskManager 任务管理器
type TaskManager struct {
	processors map[string]TaskProcessor
	repoMgr    *repo.RepositoryManager
	mu         sync.RWMutex
	running    map[uint]context.CancelFunc // 正在运行的任务
}

// NewTaskManager 创建任务管理器
func NewTaskManager(repoMgr *repo.RepositoryManager) *TaskManager {
	return &TaskManager{
		processors: make(map[string]TaskProcessor),
		repoMgr:    repoMgr,
		running:    make(map[uint]context.CancelFunc),
	}
}

// RegisterProcessor 注册任务处理器
func (tm *TaskManager) RegisterProcessor(processor TaskProcessor) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.processors[processor.GetTaskType()] = processor
	utils.Debug("注册任务处理器: %s", processor.GetTaskType())
}

// getRegisteredProcessors 获取已注册的处理器列表（用于调试）
func (tm *TaskManager) getRegisteredProcessors() []string {
	var types []string
	for taskType := range tm.processors {
		types = append(types, taskType)
	}
	return types
}

// StartTask 启动任务
func (tm *TaskManager) StartTask(taskID uint) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	utils.Debug("StartTask: 尝试启动任务 %d", taskID)

	// 检查任务是否已在运行
	if _, exists := tm.running[taskID]; exists {
		utils.Debug("任务 %d 已在运行中", taskID)
		return fmt.Errorf("任务 %d 已在运行中", taskID)
	}

	// 获取任务信息
	task, err := tm.repoMgr.TaskRepository.GetByID(taskID)
	if err != nil {
		utils.Error("获取任务失败: %v", err)
		return fmt.Errorf("获取任务失败: %v", err)
	}

	utils.Debug("StartTask: 获取到任务 %d, 类型: %s, 状态: %s", task.ID, task.Type, task.Status)

	// 获取处理器
	processor, exists := tm.processors[string(task.Type)]
	if !exists {
		utils.Error("未找到任务类型 %s 的处理器, 已注册的处理器: %v", task.Type, tm.getRegisteredProcessors())
		return fmt.Errorf("未找到任务类型 %s 的处理器", task.Type)
	}

	utils.Debug("StartTask: 找到处理器 %s", task.Type)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	tm.running[taskID] = cancel

	utils.Debug("StartTask: 启动后台任务协程")
	// 启动后台任务
	go tm.processTask(ctx, task, processor)

	utils.Info("StartTask: 任务 %d 启动成功", taskID)
	return nil
}

// PauseTask 暂停任务
func (tm *TaskManager) PauseTask(taskID uint) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	utils.Info("PauseTask: 尝试暂停任务 %d", taskID)

	// 检查任务是否在运行
	cancel, exists := tm.running[taskID]
	if !exists {
		// 检查数据库中任务状态
		task, err := tm.repoMgr.TaskRepository.GetByID(taskID)
		if err != nil {
			utils.Error("获取任务信息失败: %v", err)
			return fmt.Errorf("获取任务信息失败: %v", err)
		}

		// 如果数据库中的状态是running，说明服务器重启了，直接更新状态
		if task.Status == "running" {
			utils.Info("任务 %d 在数据库中状态为running，但内存中不存在，可能是服务器重启，直接更新状态为paused", taskID)
			err := tm.repoMgr.TaskRepository.UpdateStatus(taskID, "paused")
			if err != nil {
				utils.Error("更新任务状态为暂停失败: %v", err)
				return fmt.Errorf("更新任务状态失败: %v", err)
			}
			utils.Info("任务 %d 暂停成功（服务器重启恢复）", taskID)
			return nil
		}

		utils.Info("任务 %d 未在运行，无法暂停", taskID)
		return fmt.Errorf("任务 %d 未在运行", taskID)
	}

	// 停止任务（类似stop，但状态标记为paused）
	cancel()
	delete(tm.running, taskID)

	// 更新任务状态为暂停
	err := tm.repoMgr.TaskRepository.UpdateStatus(taskID, "paused")
	if err != nil {
		utils.Error("更新任务状态为暂停失败: %v", err)
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	utils.Info("任务 %d 暂停成功", taskID)
	return nil
}

// StopTask 停止任务
func (tm *TaskManager) StopTask(taskID uint) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	cancel, exists := tm.running[taskID]
	if !exists {
		// 检查数据库中任务状态
		task, err := tm.repoMgr.TaskRepository.GetByID(taskID)
		if err != nil {
			utils.Error("获取任务信息失败: %v", err)
			return fmt.Errorf("获取任务信息失败: %v", err)
		}

		// 如果数据库中的状态是running，说明服务器重启了，直接更新状态
		if task.Status == "running" {
			utils.Info("任务 %d 在数据库中状态为running，但内存中不存在，可能是服务器重启，直接更新状态为paused", taskID)
			err := tm.repoMgr.TaskRepository.UpdateStatus(taskID, "paused")
			if err != nil {
				utils.Error("更新任务状态失败: %v", err)
				return fmt.Errorf("更新任务状态失败: %v", err)
			}
			utils.Info("任务 %d 停止成功（服务器重启恢复）", taskID)
			return nil
		}

		return fmt.Errorf("任务 %d 未在运行", taskID)
	}

	cancel()
	delete(tm.running, taskID)

	// 更新任务状态为暂停
	err := tm.repoMgr.TaskRepository.UpdateStatus(taskID, "paused")
	if err != nil {
		utils.Error("更新任务状态失败: %v", err)
	}

	return nil
}

// processTask 处理任务
func (tm *TaskManager) processTask(ctx context.Context, task *entity.Task, processor TaskProcessor) {
	defer func() {
		tm.mu.Lock()
		delete(tm.running, task.ID)
		tm.mu.Unlock()
		utils.Debug("processTask: 任务 %d 处理完成，清理资源", task.ID)
	}()

	utils.Debug("processTask: 开始处理任务: %d, 类型: %s", task.ID, task.Type)

	// 更新任务状态为运行中
	err := tm.repoMgr.TaskRepository.UpdateStatus(task.ID, "running")
	if err != nil {
		utils.Error("更新任务状态失败: %v", err)
		return
	}

	// 获取任务项统计信息，用于计算正确的进度
	stats, err := tm.repoMgr.TaskItemRepository.GetStatsByTaskID(task.ID)
	if err != nil {
		utils.Error("获取任务项统计失败: %v", err)
		stats = map[string]int{
			"total":      0,
			"pending":    0,
			"processing": 0,
			"completed":  0,
			"failed":     0,
		}
	}

	// 获取待处理的任务项
	items, err := tm.repoMgr.TaskItemRepository.GetByTaskIDAndStatus(task.ID, "pending")
	if err != nil {
		utils.Error("获取任务项失败: %v", err)
		tm.markTaskFailed(task.ID, fmt.Sprintf("获取任务项失败: %v", err))
		return
	}

	// 计算总任务项数和已完成的项数
	totalItems := stats["total"]
	completedItems := stats["completed"]
	initialFailedItems := stats["failed"]
	processingItems := stats["processing"]

	// 如果当前批次有处理中的任务项，重置它们为pending状态（服务器重启恢复）
	if processingItems > 0 {
		utils.Debug("任务 %d 发现 %d 个处理中的任务项，重置为pending状态", task.ID, processingItems)
		err = tm.repoMgr.TaskItemRepository.ResetProcessingItems(task.ID)
		if err != nil {
			utils.Error("重置处理中任务项失败: %v", err)
		}
		// 重新获取待处理的任务项
		items, err = tm.repoMgr.TaskItemRepository.GetByTaskIDAndStatus(task.ID, "pending")
		if err != nil {
			utils.Error("重新获取任务项失败: %v", err)
			tm.markTaskFailed(task.ID, fmt.Sprintf("重新获取任务项失败: %v", err))
			return
		}
	}

	currentBatchItems := len(items)
	processedItems := completedItems + initialFailedItems // 已经处理的项目数
	successItems := completedItems
	failedItems := initialFailedItems

	utils.Debug("任务 %d 统计信息: 总计=%d, 已完成=%d, 已失败=%d, 待处理=%d",
		task.ID, totalItems, completedItems, failedItems, currentBatchItems)

	for _, item := range items {
		select {
		case <-ctx.Done():
			utils.Debug("任务 %d 被取消", task.ID)
			return
		default:
			// 处理单个任务项
			err := tm.processTaskItem(ctx, task.ID, item, processor)
			processedItems++

			if err != nil {
				failedItems++
				utils.Error("处理任务项 %d 失败: %v", item.ID, err)
			} else {
				successItems++
			}

			// 更新任务进度（基于总任务项数）
			if totalItems > 0 {
				progress := float64(processedItems) / float64(totalItems) * 100
				tm.updateTaskProgress(task.ID, progress, processedItems, successItems, failedItems)
			}
		}
	}

	// 任务完成
	status := "completed"
	message := fmt.Sprintf("任务完成，共处理 %d 项，成功 %d 项，失败 %d 项", processedItems, successItems, failedItems)

	if failedItems > 0 && successItems == 0 {
		status = "failed"
		message = fmt.Sprintf("任务失败，共处理 %d 项，全部失败", processedItems)
	} else if failedItems > 0 {
		status = "partial_success"
		message = fmt.Sprintf("任务部分成功，共处理 %d 项，成功 %d 项，失败 %d 项", processedItems, successItems, failedItems)
	}

	err = tm.repoMgr.TaskRepository.UpdateStatusAndMessage(task.ID, status, message)
	if err != nil {
		utils.Error("更新任务状态失败: %v", err)
	}

	utils.Info("任务 %d 处理完成: %s", task.ID, message)
}

// processTaskItem 处理单个任务项
func (tm *TaskManager) processTaskItem(ctx context.Context, taskID uint, item *entity.TaskItem, processor TaskProcessor) error {
	// 更新任务项状态为处理中
	err := tm.repoMgr.TaskItemRepository.UpdateStatus(item.ID, "processing")
	if err != nil {
		return fmt.Errorf("更新任务项状态失败: %v", err)
	}

	// 处理任务项
	err = processor.Process(ctx, taskID, item)

	if err != nil {
		// 处理失败
		outputData := map[string]interface{}{
			"error": err.Error(),
			"time":  utils.GetCurrentTime(),
		}
		outputJSON, _ := json.Marshal(outputData)

		updateErr := tm.repoMgr.TaskItemRepository.UpdateStatusAndOutput(item.ID, "failed", string(outputJSON))
		if updateErr != nil {
			utils.Error("更新失败任务项状态失败: %v", updateErr)
		}
		return err
	}

	// 处理成功
	outputData := map[string]interface{}{
		"success": true,
		"time":    utils.GetCurrentTime(),
	}
	outputJSON, _ := json.Marshal(outputData)

	err = tm.repoMgr.TaskItemRepository.UpdateStatusAndOutput(item.ID, "completed", string(outputJSON))
	if err != nil {
		utils.Error("更新成功任务项状态失败: %v", err)
	}

	return nil
}

// updateTaskProgress 更新任务进度
func (tm *TaskManager) updateTaskProgress(taskID uint, progress float64, processed, success, failed int) {
	// 更新任务统计信息
	err := tm.repoMgr.TaskRepository.UpdateTaskStats(taskID, processed, success, failed)
	if err != nil {
		utils.Error("更新任务统计信息失败: %v", err)
	}

	// 更新进度数据（用于兼容性）
	progressData := map[string]interface{}{
		"progress":  progress,
		"processed": processed,
		"success":   success,
		"failed":    failed,
		"time":      utils.GetCurrentTime(),
	}

	progressJSON, _ := json.Marshal(progressData)

	err = tm.repoMgr.TaskRepository.UpdateProgress(taskID, progress, string(progressJSON))
	if err != nil {
		utils.Error("更新任务进度数据失败: %v", err)
	}
}

// markTaskFailed 标记任务失败
func (tm *TaskManager) markTaskFailed(taskID uint, message string) {
	err := tm.repoMgr.TaskRepository.UpdateStatusAndMessage(taskID, "failed", message)
	if err != nil {
		utils.Error("标记任务失败状态失败: %v", err)
	}
}

// GetTaskStatus 获取任务状态
func (tm *TaskManager) GetTaskStatus(taskID uint) (string, error) {
	task, err := tm.repoMgr.TaskRepository.GetByID(taskID)
	if err != nil {
		return "", err
	}
	return string(task.Status), nil
}

// IsTaskRunning 检查任务是否在运行
func (tm *TaskManager) IsTaskRunning(taskID uint) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	_, exists := tm.running[taskID]
	return exists
}

// RecoverRunningTasks 恢复运行中的任务（服务器重启后调用）
func (tm *TaskManager) RecoverRunningTasks() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	utils.Info("开始恢复运行中的任务...")

	// 获取数据库中状态为running的任务
	tasks, _, err := tm.repoMgr.TaskRepository.GetList(1, 1000, "", "running")
	if err != nil {
		utils.Error("获取运行中任务失败: %v", err)
		return fmt.Errorf("获取运行中任务失败: %v", err)
	}

	recoveredCount := 0
	for _, task := range tasks {
		// 检查任务是否已在内存中运行
		if _, exists := tm.running[task.ID]; exists {
			utils.Info("任务 %d 已在内存中运行，跳过恢复", task.ID)
			continue
		}

		// 获取处理器
		processor, exists := tm.processors[string(task.Type)]
		if !exists {
			utils.Error("未找到任务类型 %s 的处理器，跳过恢复任务 %d", task.Type, task.ID)
			// 将任务状态重置为pending，避免卡在running状态
			tm.repoMgr.TaskRepository.UpdateStatus(task.ID, "pending")
			continue
		}

		// 创建上下文并恢复任务
		ctx, cancel := context.WithCancel(context.Background())
		tm.running[task.ID] = cancel

		utils.Info("恢复任务 %d (类型: %s)", task.ID, task.Type)
		go tm.processTask(ctx, task, processor)
		recoveredCount++
	}

	utils.Info("任务恢复完成，共恢复 %d 个任务", recoveredCount)
	return nil
}
