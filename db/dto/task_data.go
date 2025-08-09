package dto

import "fmt"

// BatchTransferInputData 批量转存任务的输入数据
type BatchTransferInputData struct {
	Title      string `json:"title"`       // 资源标题
	URL        string `json:"url"`         // 资源链接
	CategoryID *uint  `json:"category_id"` // 分类ID
	TagIDs     []uint `json:"tag_ids"`     // 标签ID列表
}

// BatchTransferOutputData 批量转存任务的输出数据
type BatchTransferOutputData struct {
	ResourceID uint   `json:"resource_id"` // 创建的资源ID
	SaveURL    string `json:"save_url"`    // 转存后的链接
	PlatformID uint   `json:"platform_id"` // 平台ID
}

// TaskItemData 通用任务项数据接口
type TaskItemData interface {
	// GetDisplayName 获取显示名称（用于前端显示）
	GetDisplayName() string
	// Validate 验证数据有效性
	Validate() error
}

// GetDisplayName 实现TaskItemData接口
func (data BatchTransferInputData) GetDisplayName() string {
	return data.Title
}

// Validate 验证批量转存输入数据
func (data BatchTransferInputData) Validate() error {
	if data.Title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if data.URL == "" {
		return fmt.Errorf("链接不能为空")
	}
	// 这里可以添加URL格式验证
	return nil
}

// GetDisplayName 实现TaskItemData接口
func (data BatchTransferOutputData) GetDisplayName() string {
	return fmt.Sprintf("ResourceID: %d", data.ResourceID)
}

// Validate 验证批量转存输出数据
func (data BatchTransferOutputData) Validate() error {
	if data.ResourceID == 0 {
		return fmt.Errorf("资源ID不能为空")
	}
	return nil
}
