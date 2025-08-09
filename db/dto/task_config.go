package dto

import "fmt"

// BatchTransferTaskConfig 批量转存任务配置
type BatchTransferTaskConfig struct {
	CategoryID *uint  `json:"category_id"` // 默认分类ID
	TagIDs     []uint `json:"tag_ids"`     // 默认标签ID列表
}

// TaskConfig 通用任务配置接口
type TaskConfig interface {
	// Validate 验证配置有效性
	Validate() error
}

// Validate 验证批量转存任务配置
func (config BatchTransferTaskConfig) Validate() error {
	// 这里可以添加配置验证逻辑
	return nil
}

// 示例：未来可能的其他任务类型配置

// DataSyncTaskConfig 数据同步任务配置（示例）
type DataSyncTaskConfig struct {
	SourceType string `json:"source_type"` // 数据源类型
	TargetType string `json:"target_type"` // 目标类型
	SyncMode   string `json:"sync_mode"`   // 同步模式
}

// Validate 验证数据同步任务配置
func (config DataSyncTaskConfig) Validate() error {
	if config.SourceType == "" {
		return fmt.Errorf("数据源类型不能为空")
	}
	if config.TargetType == "" {
		return fmt.Errorf("目标类型不能为空")
	}
	return nil
}

// CleanupTaskConfig 清理任务配置（示例）
type CleanupTaskConfig struct {
	RetentionDays int    `json:"retention_days"` // 保留天数
	CleanupType   string `json:"cleanup_type"`   // 清理类型
}

// Validate 验证清理任务配置
func (config CleanupTaskConfig) Validate() error {
	if config.RetentionDays < 0 {
		return fmt.Errorf("保留天数不能为负数")
	}
	return nil
}
