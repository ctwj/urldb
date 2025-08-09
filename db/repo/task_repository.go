package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// TaskRepository 任务仓库接口
type TaskRepository interface {
	BaseRepository[entity.Task]
	FindWithItems(id uint) (*entity.Task, error)
	FindWithPagination(page, pageSize int) ([]entity.Task, int64, error)
	UpdateProgress(id uint, processed, success, failed int) error
	UpdateStatus(id uint, status entity.TaskStatus) error
	GetRunningTasks() ([]entity.Task, error)
}

// TaskRepositoryImpl 任务仓库实现
type TaskRepositoryImpl struct {
	BaseRepositoryImpl[entity.Task]
}

// NewTaskRepository 创建任务仓库
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.Task]{db: db},
	}
}

// FindWithItems 查找任务及其所有项目
func (r *TaskRepositoryImpl) FindWithItems(id uint) (*entity.Task, error) {
	var task entity.Task
	err := r.db.Preload("TaskItems").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// FindWithPagination 分页查询任务
func (r *TaskRepositoryImpl) FindWithPagination(page, pageSize int) ([]entity.Task, int64, error) {
	var tasks []entity.Task
	var total int64

	// 获取总数
	err := r.db.Model(&entity.Task{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateProgress 更新任务进度
func (r *TaskRepositoryImpl) UpdateProgress(id uint, processed, success, failed int) error {
	return r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"processed_items": processed,
		"success_items":   success,
		"failed_items":    failed,
	}).Error
}

// UpdateStatus 更新任务状态
func (r *TaskRepositoryImpl) UpdateStatus(id uint, status entity.TaskStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 如果状态为运行中，设置开始时间
	if status == entity.TaskStatusRunning {
		updates["started_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}

	// 如果状态为完成或失败，设置完成时间
	if status == entity.TaskStatusCompleted || status == entity.TaskStatusFailed {
		updates["completed_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}

	return r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(updates).Error
}

// GetRunningTasks 获取正在运行的任务
func (r *TaskRepositoryImpl) GetRunningTasks() ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.Where("status IN ?", []entity.TaskStatus{
		entity.TaskStatusRunning,
		entity.TaskStatusPending,
	}).Find(&tasks).Error
	return tasks, err
}
