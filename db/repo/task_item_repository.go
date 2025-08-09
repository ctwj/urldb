package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// TaskItemRepository 任务项仓库接口
type TaskItemRepository interface {
	BaseRepository[entity.TaskItem]
	FindByTaskID(taskID uint) ([]entity.TaskItem, error)
	FindByTaskIDWithPagination(taskID uint, page, pageSize int) ([]entity.TaskItem, int64, error)
	UpdateStatus(id uint, status entity.TaskItemStatus, errorMsg string) error
	UpdateSuccess(id uint, outputData string) error
	GetPendingItemsByTaskID(taskID uint) ([]entity.TaskItem, error)
	BatchCreate(items []entity.TaskItem) error
}

// TaskItemRepositoryImpl 任务项仓库实现
type TaskItemRepositoryImpl struct {
	BaseRepositoryImpl[entity.TaskItem]
}

// NewTaskItemRepository 创建任务项仓库
func NewTaskItemRepository(db *gorm.DB) TaskItemRepository {
	return &TaskItemRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.TaskItem]{db: db},
	}
}

// FindByTaskID 根据任务ID查找所有任务项
func (r *TaskItemRepositoryImpl) FindByTaskID(taskID uint) ([]entity.TaskItem, error) {
	var items []entity.TaskItem
	err := r.db.Where("task_id = ?", taskID).Order("created_at ASC").Find(&items).Error
	return items, err
}

// FindByTaskIDWithPagination 根据任务ID分页查找任务项
func (r *TaskItemRepositoryImpl) FindByTaskIDWithPagination(taskID uint, page, pageSize int) ([]entity.TaskItem, int64, error) {
	var items []entity.TaskItem
	var total int64

	// 获取总数
	err := r.db.Model(&entity.TaskItem{}).Where("task_id = ?", taskID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Where("task_id = ?", taskID).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UpdateStatus 更新任务项状态
func (r *TaskItemRepositoryImpl) UpdateStatus(id uint, status entity.TaskItemStatus, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	if status != entity.TaskItemStatusPending {
		updates["processed_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}

	return r.db.Model(&entity.TaskItem{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateSuccess 更新任务项为成功状态
func (r *TaskItemRepositoryImpl) UpdateSuccess(id uint, outputData string) error {
	updates := map[string]interface{}{
		"status":       entity.TaskItemStatusSuccess,
		"output_data":  outputData,
		"processed_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}

	return r.db.Model(&entity.TaskItem{}).Where("id = ?", id).Updates(updates).Error
}

// GetPendingItemsByTaskID 获取任务的待处理项目
func (r *TaskItemRepositoryImpl) GetPendingItemsByTaskID(taskID uint) ([]entity.TaskItem, error) {
	var items []entity.TaskItem
	err := r.db.Where("task_id = ? AND status = ?", taskID, entity.TaskItemStatusPending).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

// BatchCreate 批量创建任务项
func (r *TaskItemRepositoryImpl) BatchCreate(items []entity.TaskItem) error {
	return r.db.CreateInBatches(items, 100).Error
}
