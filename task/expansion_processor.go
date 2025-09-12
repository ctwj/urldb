package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// ExpansionProcessor 扩容任务处理器
type ExpansionProcessor struct {
	repoMgr *repo.RepositoryManager
}

// NewExpansionProcessor 创建扩容任务处理器
func NewExpansionProcessor(repoMgr *repo.RepositoryManager) *ExpansionProcessor {
	return &ExpansionProcessor{
		repoMgr: repoMgr,
	}
}

// GetTaskType 获取任务类型
func (ep *ExpansionProcessor) GetTaskType() string {
	return "expansion"
}

// ExpansionInput 扩容任务输入数据结构
type ExpansionInput struct {
	PanAccountID uint `json:"pan_account_id"`
}

// ExpansionOutput 扩容任务输出数据结构
type ExpansionOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Time    string `json:"time"`
}

// Process 处理扩容任务项
func (ep *ExpansionProcessor) Process(ctx context.Context, taskID uint, item *entity.TaskItem) error {
	utils.Info("开始处理扩容任务项: %d", item.ID)

	// 解析输入数据
	var input ExpansionInput
	if err := json.Unmarshal([]byte(item.InputData), &input); err != nil {
		return fmt.Errorf("解析输入数据失败: %v", err)
	}

	// 验证输入数据
	if err := ep.validateInput(&input); err != nil {
		return fmt.Errorf("输入数据验证失败: %v", err)
	}

	// 检查账号是否已经扩容过
	exists, err := ep.checkExpansionExists(input.PanAccountID)
	if err != nil {
		utils.Error("检查扩容记录失败: %v", err)
		return fmt.Errorf("检查扩容记录失败: %v", err)
	}

	if exists {
		output := ExpansionOutput{
			Success: false,
			Message: "账号已扩容过",
			Error:   "每个账号只能扩容一次",
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Info("账号已扩容过，跳过扩容: 账号ID %d", input.PanAccountID)
		return fmt.Errorf("账号已扩容过")
	}

	// 检查账号类型（只支持quark账号）
	if err := ep.checkAccountType(input.PanAccountID); err != nil {
		output := ExpansionOutput{
			Success: false,
			Message: "账号类型不支持扩容",
			Error:   err.Error(),
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("账号类型不支持扩容: %v", err)
		return err
	}

	// 执行扩容操作（这里留空，直接返回成功）
	if err := ep.performExpansion(ctx, input.PanAccountID); err != nil {
		output := ExpansionOutput{
			Success: false,
			Message: "扩容失败",
			Error:   err.Error(),
			Time:    utils.GetCurrentTimeString(),
		}

		outputJSON, _ := json.Marshal(output)
		item.OutputData = string(outputJSON)

		utils.Error("扩容任务项处理失败: %d, 错误: %v", item.ID, err)
		return fmt.Errorf("扩容失败: %v", err)
	}

	// 扩容成功
	output := ExpansionOutput{
		Success: true,
		Message: "扩容成功",
		Time:    utils.GetCurrentTimeString(),
	}

	outputJSON, _ := json.Marshal(output)
	item.OutputData = string(outputJSON)

	utils.Info("扩容任务项处理完成: %d, 账号ID: %d", item.ID, input.PanAccountID)
	return nil
}

// validateInput 验证输入数据
func (ep *ExpansionProcessor) validateInput(input *ExpansionInput) error {
	if input.PanAccountID == 0 {
		return fmt.Errorf("账号ID不能为空")
	}
	return nil
}

// checkExpansionExists 检查账号是否已经扩容过
func (ep *ExpansionProcessor) checkExpansionExists(panAccountID uint) (bool, error) {
	// 查询所有expansion类型的任务
	tasks, _, err := ep.repoMgr.TaskRepository.GetList(1, 1000, "expansion", "completed")
	if err != nil {
		return false, fmt.Errorf("获取扩容任务列表失败: %v", err)
	}

	// 检查每个任务的配置中是否包含该账号ID
	for _, task := range tasks {
		if task.Config != "" {
			var taskConfig map[string]interface{}
			if err := json.Unmarshal([]byte(task.Config), &taskConfig); err == nil {
				if configAccountID, ok := taskConfig["pan_account_id"].(float64); ok {
					if uint(configAccountID) == panAccountID {
						// 找到了该账号的扩容任务，检查任务状态
						if task.Status == "completed" {
							// 如果任务已完成，说明已经扩容过
							return true, nil
						}
					}
				}
			}
		}
	}

	return false, nil
}

// checkAccountType 检查账号类型（只支持quark账号）
func (ep *ExpansionProcessor) checkAccountType(panAccountID uint) error {
	// 获取账号信息
	cks, err := ep.repoMgr.CksRepository.FindByID(panAccountID)
	if err != nil {
		return fmt.Errorf("获取账号信息失败: %v", err)
	}

	// 检查是否为quark账号
	if cks.ServiceType != "quark" {
		return fmt.Errorf("只支持quark账号扩容，当前账号类型: %s", cks.ServiceType)
	}

	return nil
}

// performExpansion 执行扩容操作
func (ep *ExpansionProcessor) performExpansion(ctx context.Context, panAccountID uint) error {
	// 扩容逻辑暂时留空，直接返回成功
	// TODO: 实现具体的扩容逻辑

	utils.Info("执行扩容操作，账号ID: %d", panAccountID)

	// 模拟扩容操作延迟
	// time.Sleep(2 * time.Second)

	return nil
}
