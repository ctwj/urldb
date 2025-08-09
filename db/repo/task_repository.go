package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// TaskRepository 任务仓库接口
type TaskRepository interface {
	GetByID(id uint) (*entity.Task, error)
	Create(task *entity.Task) error
	Delete(id uint) error
	GetList(page, pageSize int, taskType, status string) ([]*entity.Task, int64, error)
	UpdateStatus(id uint, status string) error
	UpdateProgress(id uint, progress float64, progressData string) error
	UpdateStatusAndMessage(id uint, status, message string) error
}

// TaskRepositoryImpl 任务仓库实现
type TaskRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓库
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{
		db: db,
	}
}

// GetByID 根据ID获取任务
func (r *TaskRepositoryImpl) GetByID(id uint) (*entity.Task, error) {
	var task entity.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Create 创建任务
func (r *TaskRepositoryImpl) Create(task *entity.Task) error {
	return r.db.Create(task).Error
}

// Delete 删除任务
func (r *TaskRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Task{}, id).Error
}

// GetList 获取任务列表
func (r *TaskRepositoryImpl) GetList(page, pageSize int, taskType, status string) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	query := r.db.Model(&entity.Task{})

	// 添加过滤条件
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}
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
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateStatus 更新任务状态
func (r *TaskRepositoryImpl) UpdateStatus(id uint, status string) error {
	return r.db.Model(&entity.Task{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateProgress 更新任务进度
func (r *TaskRepositoryImpl) UpdateProgress(id uint, progress float64, progressData string) error {
	return r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"progress":      progress,
		"progress_data": progressData,
	}).Error
}

// UpdateStatusAndMessage 更新任务状态和消息
func (r *TaskRepositoryImpl) UpdateStatusAndMessage(id uint, status, message string) error {
	return r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  status,
		"message": message,
	}).Error
}
