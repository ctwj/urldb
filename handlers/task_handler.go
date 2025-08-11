package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	repoMgr     *repo.RepositoryManager
	taskManager *task.TaskManager
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(repoMgr *repo.RepositoryManager, taskManager *task.TaskManager) *TaskHandler {
	return &TaskHandler{
		repoMgr:     repoMgr,
		taskManager: taskManager,
	}
}

// 批量转存任务资源项
type BatchTransferResource struct {
	Title      string `json:"title" binding:"required"`
	URL        string `json:"url" binding:"required"`
	CategoryID uint   `json:"category_id,omitempty"`
	Tags       []uint `json:"tags,omitempty"`
}

// CreateBatchTransferTask 创建批量转存任务
func (h *TaskHandler) CreateBatchTransferTask(c *gin.Context) {
	var req struct {
		Title       string                  `json:"title" binding:"required"`
		Description string                  `json:"description"`
		Resources   []BatchTransferResource `json:"resources" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.Info("创建批量转存任务: %s，资源数量: %d", req.Title, len(req.Resources))

	// 创建任务
	newTask := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
		Type:        "transfer",
		Status:      "pending",
		TotalItems:  len(req.Resources),
		CreatedAt:   utils.GetCurrentTime(),
		UpdatedAt:   utils.GetCurrentTime(),
	}

	err := h.repoMgr.TaskRepository.Create(newTask)
	if err != nil {
		utils.Error("创建任务失败: %v", err)
		ErrorResponse(c, "创建任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建任务项
	for _, resource := range req.Resources {
		// 构建转存输入数据
		transferInput := task.TransferInput{
			Title:      resource.Title,
			URL:        resource.URL,
			CategoryID: resource.CategoryID,
			Tags:       resource.Tags,
		}

		inputJSON, _ := json.Marshal(transferInput)

		taskItem := &entity.TaskItem{
			TaskID:    newTask.ID,
			Status:    "pending",
			InputData: string(inputJSON),
			CreatedAt: utils.GetCurrentTime(),
			UpdatedAt: utils.GetCurrentTime(),
		}

		err = h.repoMgr.TaskItemRepository.Create(taskItem)
		if err != nil {
			utils.Error("创建任务项失败: %v", err)
			// 继续创建其他任务项
		}
	}

	utils.Info("批量转存任务创建完成: %d, 共 %d 项", newTask.ID, len(req.Resources))

	SuccessResponse(c, gin.H{
		"task_id":     newTask.ID,
		"total_items": len(req.Resources),
		"message":     "任务创建成功",
	})
}

// StartTask 启动任务
func (h *TaskHandler) StartTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.Info("启动任务: %d", taskID)

	err = h.taskManager.StartTask(uint(taskID))
	if err != nil {
		utils.Error("启动任务失败: %v", err)
		ErrorResponse(c, "启动任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "任务启动成功",
	})
}

// StopTask 停止任务
func (h *TaskHandler) StopTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.Info("停止任务: %d", taskID)

	err = h.taskManager.StopTask(uint(taskID))
	if err != nil {
		utils.Error("停止任务失败: %v", err)
		ErrorResponse(c, "停止任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "任务停止成功",
	})
}

// PauseTask 暂停任务
func (h *TaskHandler) PauseTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.Info("暂停任务: %d", taskID)

	err = h.taskManager.PauseTask(uint(taskID))
	if err != nil {
		utils.Error("暂停任务失败: %v", err)
		ErrorResponse(c, "暂停任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "任务暂停成功",
	})
}

// GetTaskStatus 获取任务状态
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 获取任务详情
	task, err := h.repoMgr.TaskRepository.GetByID(uint(taskID))
	if err != nil {
		ErrorResponse(c, "任务不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	// 获取任务项统计
	stats, err := h.repoMgr.TaskItemRepository.GetStatsByTaskID(uint(taskID))
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

	// 检查任务是否在运行
	isRunning := h.taskManager.IsTaskRunning(uint(taskID))

	SuccessResponse(c, gin.H{
		"id":              task.ID,
		"title":           task.Title,
		"description":     task.Description,
		"task_type":       task.Type,
		"status":          task.Status,
		"total_items":     task.TotalItems,
		"processed_items": task.ProcessedItems,
		"success_items":   task.SuccessItems,
		"failed_items":    task.FailedItems,
		"is_running":      isRunning,
		"stats":           stats,
		"created_at":      task.CreatedAt,
		"updated_at":      task.UpdatedAt,
	})
}

// GetTasks 获取任务列表
func (h *TaskHandler) GetTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	taskType := c.Query("task_type")
	status := c.Query("status")

	utils.Info("GetTasks: 获取任务列表 page=%d, pageSize=%d, taskType=%s, status=%s", page, pageSize, taskType, status)

	tasks, total, err := h.repoMgr.TaskRepository.GetList(page, pageSize, taskType, status)
	if err != nil {
		utils.Error("获取任务列表失败: %v", err)
		ErrorResponse(c, "获取任务列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("GetTasks: 从数据库获取到 %d 个任务", len(tasks))

	// 为每个任务添加运行状态
	var result []gin.H
	for _, task := range tasks {
		isRunning := h.taskManager.IsTaskRunning(task.ID)
		utils.Info("GetTasks: 任务 %d (%s) 数据库状态: %s, TaskManager运行状态: %v", task.ID, task.Title, task.Status, isRunning)

		result = append(result, gin.H{
			"id":              task.ID,
			"title":           task.Title,
			"description":     task.Description,
			"task_type":       task.Type,
			"status":          task.Status,
			"total_items":     task.TotalItems,
			"processed_items": task.ProcessedItems,
			"success_items":   task.SuccessItems,
			"failed_items":    task.FailedItems,
			"is_running":      isRunning,
			"created_at":      task.CreatedAt,
			"updated_at":      task.UpdatedAt,
		})
	}

	SuccessResponse(c, gin.H{
		"items": result,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetTaskItems 获取任务项列表
func (h *TaskHandler) GetTaskItems(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10000"))
	status := c.Query("status")

	items, total, err := h.repoMgr.TaskItemRepository.GetListByTaskID(uint(taskID), page, pageSize, status)
	if err != nil {
		utils.Error("获取任务项列表失败: %v", err)
		ErrorResponse(c, "获取任务项列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 解析输入和输出数据
	var result []gin.H
	for _, item := range items {
		itemData := gin.H{
			"id":         item.ID,
			"status":     item.Status,
			"created_at": item.CreatedAt,
			"updated_at": item.UpdatedAt,
		}

		// 解析输入数据
		if item.InputData != "" {
			var inputData map[string]interface{}
			if err := json.Unmarshal([]byte(item.InputData), &inputData); err == nil {
				itemData["input"] = inputData
			}
		}

		// 解析输出数据
		if item.OutputData != "" {
			var outputData map[string]interface{}
			if err := json.Unmarshal([]byte(item.OutputData), &outputData); err == nil {
				itemData["output"] = outputData
			}
		}

		result = append(result, itemData)
	}

	SuccessResponse(c, gin.H{
		"items": result,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的任务ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 检查任务是否在运行
	if h.taskManager.IsTaskRunning(uint(taskID)) {
		ErrorResponse(c, "任务正在运行，请先停止任务", http.StatusBadRequest)
		return
	}

	// 删除任务项
	err = h.repoMgr.TaskItemRepository.DeleteByTaskID(uint(taskID))
	if err != nil {
		utils.Error("删除任务项失败: %v", err)
		ErrorResponse(c, "删除任务项失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 删除任务
	err = h.repoMgr.TaskRepository.Delete(uint(taskID))
	if err != nil {
		utils.Error("删除任务失败: %v", err)
		ErrorResponse(c, "删除任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("任务删除成功: %d", taskID)
	SuccessResponse(c, gin.H{
		"message": "任务删除成功",
	})
}
