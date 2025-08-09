package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// TaskItemRepository 任务项仓库接口
type TaskItemRepository interface {
	GetByID(id uint) (*entity.TaskItem, error)
	Create(item *entity.TaskItem) error
	Delete(id uint) error
	DeleteByTaskID(taskID uint) error
	GetByTaskIDAndStatus(taskID uint, status string) ([]*entity.TaskItem, error)
	GetListByTaskID(taskID uint, page, pageSize int, status string) ([]*entity.TaskItem, int64, error)
	UpdateStatus(id uint, status string) error
	UpdateStatusAndOutput(id uint, status, outputData string) error
	GetStatsByTaskID(taskID uint) (map[string]int, error)
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
	var items []*entity.TaskItem
	err := r.db.Where("task_id = ? AND status = ?", taskID, status).Order("id ASC").Find(&items).Error
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
	err = query.Offset(offset).Limit(pageSize).Order("item_index ASC").Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UpdateStatus 更新任务项状态
func (r *TaskItemRepositoryImpl) UpdateStatus(id uint, status string) error {
	return r.db.Model(&entity.TaskItem{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateStatusAndOutput 更新任务项状态和输出数据
func (r *TaskItemRepositoryImpl) UpdateStatusAndOutput(id uint, status, outputData string) error {
	return r.db.Model(&entity.TaskItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"output_data": outputData,
	}).Error
}

// GetStatsByTaskID 获取任务项统计信息
func (r *TaskItemRepositoryImpl) GetStatsByTaskID(taskID uint) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := r.db.Model(&entity.TaskItem{}).
		Select("status, count(*) as count").
		Where("task_id = ?", taskID).
		Group("status").
		Find(&results).Error

	if err != nil {
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

	return stats, nil
}
