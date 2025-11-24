package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// TaskItemRepository 任务项仓库接口
type TaskItemRepository interface {
	GetByID(id uint) (*entity.TaskItem, error)
	Create(item *entity.TaskItem) error
	Update(item *entity.TaskItem) error
	Delete(id uint) error
	DeleteByTaskID(taskID uint) error
	GetByTaskIDAndStatus(taskID uint, status string) ([]*entity.TaskItem, error)
	GetListByTaskID(taskID uint, page, pageSize int, status string) ([]*entity.TaskItem, int64, error)
	UpdateStatus(id uint, status string) error
	UpdateStatusAndOutput(id uint, status, outputData string) error
	GetStatsByTaskID(taskID uint) (map[string]int, error)
	GetIndexStats() (map[string]int, error)
	ResetProcessingItems(taskID uint) error
}

// TaskItemRepositoryImpl 任务项仓库实现
type TaskItemRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskItemRepository 创建任务项仓库
func NewTaskItemRepository(db *gorm.DB) TaskItemRepository {
	return &TaskItemRepositoryImpl{
		db: db,
	}
}

// GetByID 根据ID获取任务项
func (r *TaskItemRepositoryImpl) GetByID(id uint) (*entity.TaskItem, error) {
	var item entity.TaskItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建任务项
func (r *TaskItemRepositoryImpl) Create(item *entity.TaskItem) error {
	return r.db.Create(item).Error
}

// Update 更新任务项
func (r *TaskItemRepositoryImpl) Update(item *entity.TaskItem) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.TaskItem{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"status":           item.Status,
		"error_message":    item.ErrorMessage,
		"index_status":     item.IndexStatus,
		"mobile_friendly":  item.MobileFriendly,
		"last_crawled":     item.LastCrawled,
		"status_code":      item.StatusCode,
		"input_data":       item.InputData,
		"output_data":      item.OutputData,
		"process_log":      item.ProcessLog,
		"url":              item.URL,
		"inspect_result":   item.InspectResult,
		"processed_at":     item.ProcessedAt,
		"updated_at":       time.Now(),
	}).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("Update任务项失败: ID=%d, 错误=%v, 更新耗时=%v", item.ID, err, updateDuration)
		return err
	}
	utils.Debug("Update任务项成功: ID=%d, 更新耗时=%v", item.ID, updateDuration)
	return nil
}

// Delete 删除任务项
func (r *TaskItemRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.TaskItem{}, id).Error
}

// DeleteByTaskID 根据任务ID删除所有任务项
func (r *TaskItemRepositoryImpl) DeleteByTaskID(taskID uint) error {
	return r.db.Where("task_id = ?", taskID).Delete(&entity.TaskItem{}).Error
}

// GetByTaskIDAndStatus 根据任务ID和状态获取任务项
func (r *TaskItemRepositoryImpl) GetByTaskIDAndStatus(taskID uint, status string) ([]*entity.TaskItem, error) {
	startTime := utils.GetCurrentTime()
	var items []*entity.TaskItem
	err := r.db.Where("task_id = ? AND status = ?", taskID, status).Order("id ASC").Find(&items).Error
	queryDuration := time.Since(startTime)
	if err != nil {
		utils.Error("GetByTaskIDAndStatus失败: 任务ID=%d, 状态=%s, 错误=%v, 查询耗时=%v", taskID, status, err, queryDuration)
		return nil, err
	}
	utils.Debug("GetByTaskIDAndStatus成功: 任务ID=%d, 状态=%s, 数量=%d, 查询耗时=%v", taskID, status, len(items), queryDuration)
	return items, err
}

// GetListByTaskID 根据任务ID分页获取任务项
func (r *TaskItemRepositoryImpl) GetListByTaskID(taskID uint, page, pageSize int, status string) ([]*entity.TaskItem, int64, error) {
	var items []*entity.TaskItem
	var total int64

	query := r.db.Model(&entity.TaskItem{}).Where("task_id = ?", taskID)

	// 添加状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UpdateStatus 更新任务项状态
func (r *TaskItemRepositoryImpl) UpdateStatus(id uint, status string) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.TaskItem{}).Where("id = ?", id).Update("status", status).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateStatus失败: ID=%d, 状态=%s, 错误=%v, 更新耗时=%v", id, status, err, updateDuration)
		return err
	}
	utils.Debug("UpdateStatus成功: ID=%d, 状态=%s, 更新耗时=%v", id, status, updateDuration)
	return nil
}

// UpdateStatusAndOutput 更新任务项状态和输出数据
func (r *TaskItemRepositoryImpl) UpdateStatusAndOutput(id uint, status, outputData string) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.TaskItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"output_data": outputData,
	}).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateStatusAndOutput失败: ID=%d, 状态=%s, 错误=%v, 更新耗时=%v", id, status, err, updateDuration)
		return err
	}
	utils.Debug("UpdateStatusAndOutput成功: ID=%d, 状态=%s, 更新耗时=%v", id, status, updateDuration)
	return nil
}

// GetStatsByTaskID 获取任务项统计信息
func (r *TaskItemRepositoryImpl) GetStatsByTaskID(taskID uint) (map[string]int, error) {
	startTime := utils.GetCurrentTime()
	var results []struct {
		Status string
		Count  int
	}

	err := r.db.Model(&entity.TaskItem{}).
		Select("status, count(*) as count").
		Where("task_id = ?", taskID).
		Group("status").
		Find(&results).Error

	queryDuration := time.Since(startTime)
	if err != nil {
		utils.Error("GetStatsByTaskID失败: 任务ID=%d, 错误=%v, 查询耗时=%v", taskID, err, queryDuration)
		return nil, err
	}

	stats := map[string]int{
		"total":      0,
		"pending":    0,
		"processing": 0,
		"completed":  0,
		"failed":     0,
	}

	for _, result := range results {
		stats[result.Status] = result.Count
		stats["total"] += result.Count
	}

	totalDuration := time.Since(startTime)
	utils.Debug("GetStatsByTaskID成功: 任务ID=%d, 统计信息=%v, 查询耗时=%v, 总耗时=%v", taskID, stats, queryDuration, totalDuration)
	return stats, nil
}

// ResetProcessingItems 重置处理中的任务项为pending状态
func (r *TaskItemRepositoryImpl) ResetProcessingItems(taskID uint) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.TaskItem{}).
		Where("task_id = ? AND status = ?", taskID, "processing").
		Update("status", "pending").Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("ResetProcessingItems失败: 任务ID=%d, 错误=%v, 更新耗时=%v", taskID, err, updateDuration)
		return err
	}
	utils.Debug("ResetProcessingItems成功: 任务ID=%d, 更新耗时=%v", taskID, updateDuration)
	return nil
}

// GetIndexStats 获取索引统计信息
func (r *TaskItemRepositoryImpl) GetIndexStats() (map[string]int, error) {
	stats := make(map[string]int)

	// 统计各种状态的数量
	statuses := []string{"completed", "failed", "pending"}

	for _, status := range statuses {
		var count int64
		err := r.db.Model(&entity.TaskItem{}).Where("status = ?", status).Count(&count).Error
		if err != nil {
			return nil, err
		}

		switch status {
		case "completed":
			stats["indexed"] = int(count)
		case "failed":
			stats["error"] = int(count)
		case "pending":
			stats["not_indexed"] = int(count)
		}
	}

	return stats, nil
}
