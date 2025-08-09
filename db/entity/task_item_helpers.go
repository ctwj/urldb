package entity

import (
	"encoding/json"
	"fmt"

	"github.com/ctwj/urldb/db/dto"
)

// SetInputData 设置输入数据（将结构体转换为JSON字符串）
func (item *TaskItem) SetInputData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化输入数据失败: %v", err)
	}
	item.InputData = string(jsonData)
	return nil
}

// GetInputData 获取输入数据（根据任务类型解析JSON）
func (item *TaskItem) GetInputData(taskType TaskType) (interface{}, error) {
	if item.InputData == "" {
		return nil, fmt.Errorf("输入数据为空")
	}

	switch taskType {
	case TaskTypeBatchTransfer:
		var data dto.BatchTransferInputData
		err := json.Unmarshal([]byte(item.InputData), &data)
		if err != nil {
			return nil, fmt.Errorf("解析批量转存输入数据失败: %v", err)
		}
		return data, nil
	default:
		// 对于未知任务类型，返回原始JSON数据
		var data map[string]interface{}
		err := json.Unmarshal([]byte(item.InputData), &data)
		if err != nil {
			return nil, fmt.Errorf("解析输入数据失败: %v", err)
		}
		return data, nil
	}
}

// SetOutputData 设置输出数据（将结构体转换为JSON字符串）
func (item *TaskItem) SetOutputData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化输出数据失败: %v", err)
	}
	item.OutputData = string(jsonData)
	return nil
}

// GetOutputData 获取输出数据（根据任务类型解析JSON）
func (item *TaskItem) GetOutputData(taskType TaskType) (interface{}, error) {
	if item.OutputData == "" {
		return nil, fmt.Errorf("输出数据为空")
	}

	switch taskType {
	case TaskTypeBatchTransfer:
		var data dto.BatchTransferOutputData
		err := json.Unmarshal([]byte(item.OutputData), &data)
		if err != nil {
			return nil, fmt.Errorf("解析批量转存输出数据失败: %v", err)
		}
		return data, nil
	default:
		// 对于未知任务类型，返回原始JSON数据
		var data map[string]interface{}
		err := json.Unmarshal([]byte(item.OutputData), &data)
		if err != nil {
			return nil, fmt.Errorf("解析输出数据失败: %v", err)
		}
		return data, nil
	}
}

// GetDisplayName 获取显示名称（用于前端显示）
func (item *TaskItem) GetDisplayName(taskType TaskType) string {
	inputData, err := item.GetInputData(taskType)
	if err != nil {
		return fmt.Sprintf("TaskItem#%d", item.ID)
	}

	switch taskType {
	case TaskTypeBatchTransfer:
		if data, ok := inputData.(dto.BatchTransferInputData); ok {
			return data.Title
		}
	}

	return fmt.Sprintf("TaskItem#%d", item.ID)
}

// AddProcessLog 添加处理日志
func (item *TaskItem) AddProcessLog(message string) {
	if item.ProcessLog == "" {
		item.ProcessLog = message
	} else {
		item.ProcessLog += "\n" + message
	}
}
