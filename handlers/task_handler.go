package handlers

import (
	"encoding/json"
	"fmt"
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
	PanID      uint   `json:"pan_id,omitempty"`
	Tags       []uint `json:"tags,omitempty"`
}

// CreateBatchTransferTask 创建批量转存任务
func (h *TaskHandler) CreateBatchTransferTask(c *gin.Context) {
	var req struct {
		Title            string                  `json:"title" binding:"required"`
		Description      string                  `json:"description"`
		Resources        []BatchTransferResource `json:"resources" binding:"required,min=1"`
		SelectedAccounts []uint                  `json:"selected_accounts,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("CreateBatchTransferTask - 用户创建批量转存任务 - 用户: %s, 任务标题: %s, 资源数量: %d, IP: %s", username, req.Title, len(req.Resources), clientIP)

	utils.Debug("创建批量转存任务: %s，资源数量: %d，选择账号数量: %d", req.Title, len(req.Resources), len(req.SelectedAccounts))

	// 构建任务配置
	taskConfig := map[string]interface{}{
		"selected_accounts": req.SelectedAccounts,
	}
	configJSON, _ := json.Marshal(taskConfig)

	// 创建任务
	newTask := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
		Type:        "transfer",
		Status:      "pending",
		TotalItems:  len(req.Resources),
		Config:      string(configJSON),
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
			PanID:      resource.PanID,
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

	utils.Debug("批量转存任务创建完成: %d, 共 %d 项", newTask.ID, len(req.Resources))

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

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("StartTask - 用户启动任务 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)

	err = h.taskManager.StartTask(uint(taskID))
	if err != nil {
		utils.Error("启动任务失败: %v", err)
		ErrorResponse(c, "启动任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Debug("启动任务: %d", taskID)

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

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("StopTask - 用户停止任务 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)

	err = h.taskManager.StopTask(uint(taskID))
	if err != nil {
		utils.Error("停止任务失败: %v", err)
		ErrorResponse(c, "停止任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Debug("停止任务: %d", taskID)

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

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("PauseTask - 用户暂停任务 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)

	err = h.taskManager.PauseTask(uint(taskID))
	if err != nil {
		utils.Error("暂停任务失败: %v", err)
		ErrorResponse(c, "暂停任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Debug("暂停任务: %d", taskID)

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
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	taskType := c.Query("taskType")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	
	// 获取任务列表
	tasks, total, err := h.repoMgr.TaskRepository.GetList(page, pageSize, taskType, status)
	if err != nil {
		utils.Error("获取任务列表失败: %v", err)
		ErrorResponse(c, "获取任务列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	
	// 获取任务运行状态
	var taskList []gin.H
	for _, task := range tasks {
		isRunning := h.taskManager.IsTaskRunning(task.ID)

		taskList = append(taskList, gin.H{
			"id":              task.ID,
			"title":           task.Title,
			"description":     task.Description,
			"type":            task.Type,
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
		"tasks":       taskList,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
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

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("DeleteTask - 用户删除任务 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)

	// 检查任务是否在运行
	if h.taskManager.IsTaskRunning(uint(taskID)) {
		utils.Warn("DeleteTask - 尝试删除正在运行的任务 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)
		ErrorResponse(c, "任务正在运行中，无法删除", http.StatusBadRequest)
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

	utils.Debug("任务删除成功: %d", taskID)
	utils.Info("DeleteTask - 任务删除成功 - 用户: %s, 任务ID: %d, IP: %s", username, taskID, clientIP)

	SuccessResponse(c, gin.H{
		"message": "任务删除成功",
	})
}

// CreateExpansionTask 创建扩容任务
func (h *TaskHandler) CreateExpansionTask(c *gin.Context) {
	var req struct {
		PanAccountID uint                   `json:"pan_account_id" binding:"required"`
		Description  string                 `json:"description"`
		DataSource   map[string]interface{} `json:"dataSource"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("CreateExpansionTask - 用户创建扩容任务 - 用户: %s, 账号ID: %d, IP: %s", username, req.PanAccountID, clientIP)

	utils.Debug("创建扩容任务: 账号ID %d", req.PanAccountID)

	// 获取账号信息，用于构建任务标题
	cks, err := h.repoMgr.CksRepository.FindByID(req.PanAccountID)
	if err != nil {
		utils.Error("获取账号信息失败: %v", err)
		ErrorResponse(c, "获取账号信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建账号名称
	accountName := cks.Username
	if accountName == "" {
		accountName = cks.Remark
	}
	if accountName == "" {
		accountName = fmt.Sprintf("账号%d", cks.ID)
	}

	// 构建任务配置（存储账号ID和数据源）
	taskConfig := map[string]interface{}{
		"pan_account_id": req.PanAccountID,
	}
	// 如果有数据源配置，添加到taskConfig中
	if req.DataSource != nil && len(req.DataSource) > 0 {
		taskConfig["data_source"] = req.DataSource
	}
	configJSON, _ := json.Marshal(taskConfig)

	// 创建任务标题，包含账号名称
	taskTitle := fmt.Sprintf("账号扩容 - %s", accountName)

	// 创建任务
	newTask := &entity.Task{
		Title:       taskTitle,
		Description: req.Description,
		Type:        "expansion",
		Status:      "pending",
		TotalItems:  1, // 扩容任务只有一个项目
		Config:      string(configJSON),
		CreatedAt:   utils.GetCurrentTime(),
		UpdatedAt:   utils.GetCurrentTime(),
	}

	if err := h.repoMgr.TaskRepository.Create(newTask); err != nil {
		utils.Error("创建扩容任务失败: %v", err)
		ErrorResponse(c, "创建任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建任务项
	expansionInput := task.ExpansionInput{
		PanAccountID: req.PanAccountID,
	}
	// 如果有数据源配置，添加到输入数据中
	if req.DataSource != nil && len(req.DataSource) > 0 {
		expansionInput.DataSource = req.DataSource
	}

	inputJSON, _ := json.Marshal(expansionInput)

	taskItem := &entity.TaskItem{
		TaskID:    newTask.ID,
		Status:    "pending",
		InputData: string(inputJSON),
		CreatedAt: utils.GetCurrentTime(),
		UpdatedAt: utils.GetCurrentTime(),
	}

	err = h.repoMgr.TaskItemRepository.Create(taskItem)
	if err != nil {
		utils.Error("创建扩容任务项失败: %v", err)
		// 继续处理，不返回错误
	}

	utils.Debug("扩容任务创建完成: %d", newTask.ID)

	SuccessResponse(c, gin.H{
		"task_id":     newTask.ID,
		"total_items": 1,
		"message":     "扩容任务创建成功",
	})
}

// GetExpansionAccounts 获取支持扩容的账号列表
func (h *TaskHandler) GetExpansionAccounts(c *gin.Context) {
	// 获取所有有效的账号
	cksList, err := h.repoMgr.CksRepository.FindByIsValid(false)
	if err != nil {
		utils.Error("获取账号列表失败: %v", err)
		ErrorResponse(c, "获取账号列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 过滤出 quark 账号
	var expansionAccounts []gin.H
	tasks, _, _ := h.repoMgr.TaskRepository.GetList(1, 1000, "expansion", "completed")
	for _, ck := range cksList {
		if ck.ServiceType == "quark" {
			// 使用 Username 作为账号名称，如果为空则使用 Remark
			accountName := ck.Username
			if accountName == "" {
				accountName = ck.Remark
			}
			if accountName == "" {
				accountName = "账号 " + fmt.Sprintf("%d", ck.ID)
			}

			// 检查是否已经扩容过
			expanded := false
			for _, task := range tasks {
				if task.Config != "" {
					var taskConfig map[string]interface{}
					if err := json.Unmarshal([]byte(task.Config), &taskConfig); err == nil {
						if configAccountID, ok := taskConfig["pan_account_id"].(float64); ok {
							if uint(configAccountID) == ck.ID {
								expanded = true
								break
							}
						}
					}
				}
			}

			expansionAccounts = append(expansionAccounts, gin.H{
				"id":           ck.ID,
				"name":         accountName,
				"service_type": ck.ServiceType,
				"expanded":     expanded,
				"total_space":  ck.Space,
				"used_space":   ck.UsedSpace,
				"created_at":   ck.CreatedAt,
				"updated_at":   ck.UpdatedAt,
			})
		}
	}

	SuccessResponse(c, gin.H{
		"accounts": expansionAccounts,
		"total":    len(expansionAccounts),
		"message":  "获取支持扩容账号列表成功",
	})
}

// GetExpansionOutput 获取账号扩容输出数据
func (h *TaskHandler) GetExpansionOutput(c *gin.Context) {
	accountIDStr := c.Param("accountId")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的账号ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.Debug("获取账号扩容输出数据: 账号ID %d", accountID)

	// 获取该账号的所有扩容任务
	tasks, _, err := h.repoMgr.TaskRepository.GetList(1, 1000, "expansion", "completed")
	if err != nil {
		utils.Error("获取扩容任务列表失败: %v", err)
		ErrorResponse(c, "获取扩容任务列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 查找该账号的扩容任务
	var targetTask *entity.Task
	for _, task := range tasks {
		if task.Config != "" {
			var taskConfig map[string]interface{}
			if err := json.Unmarshal([]byte(task.Config), &taskConfig); err == nil {
				if configAccountID, ok := taskConfig["pan_account_id"].(float64); ok {
					if uint(configAccountID) == uint(accountID) {
						targetTask = task
						break
					}
				}
			}
		}
	}

	if targetTask == nil {
		ErrorResponse(c, "该账号没有完成扩容任务", http.StatusNotFound)
		return
	}

	// 获取任务项，获取输出数据
	items, _, err := h.repoMgr.TaskItemRepository.GetListByTaskID(targetTask.ID, 1, 10, "completed")
	if err != nil {
		utils.Error("获取任务项失败: %v", err)
		ErrorResponse(c, "获取任务输出数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(items) == 0 {
		ErrorResponse(c, "任务项不存在", http.StatusNotFound)
		return
	}

	// 返回第一个完成的任务项的输出数据
	taskItem := items[0]
	var outputData map[string]interface{}
	if taskItem.OutputData != "" {
		if err := json.Unmarshal([]byte(taskItem.OutputData), &outputData); err != nil {
			utils.Error("解析输出数据失败: %v", err)
			ErrorResponse(c, "解析输出数据失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{
		"task_id":     targetTask.ID,
		"account_id":  accountID,
		"output_data": outputData,
		"message":     "获取扩容输出数据成功",
	})
}
