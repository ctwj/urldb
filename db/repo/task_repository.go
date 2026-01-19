package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
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
	UpdateTaskStats(id uint, processed, success, failed int) error
	UpdateStartedAt(id uint) error
	UpdateCompletedAt(id uint) error
	UpdateTotalItems(id uint, totalItems int) error
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
	startTime := utils.GetCurrentTime()
	var task entity.Task
	err := r.db.First(&task, id).Error
	queryDuration := time.Since(startTime)
	if err != nil {
		utils.Debug("GetByID失败: ID=%d, 错误=%v, 查询耗时=%v", id, err, queryDuration)
		return nil, err
	}
	utils.Debug("GetByID成功: ID=%d, 查询耗时=%v", id, queryDuration)
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
		query = query.Where("type = ?", taskType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		utils.Error("GetList获取总数失败: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		utils.Error("GetList查询失败: %v", err)
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateStatus 更新任务状态
func (r *TaskRepositoryImpl) UpdateStatus(id uint, status string) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.Task{}).Where("id = ?", id).Update("status", status).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateStatus失败: ID=%d, 状态=%s, 错误=%v, 更新耗时=%v", id, status, err, updateDuration)
		return err
	}
	utils.Debug("UpdateStatus成功: ID=%d, 状态=%s, 更新耗时=%v", id, status, updateDuration)
	return nil
}

// UpdateProgress 更新任务进度
func (r *TaskRepositoryImpl) UpdateProgress(id uint, progress float64, progressData string) error {
	startTime := utils.GetCurrentTime()
	// 检查progress和progress_data字段是否存在
	var count int64
	err := r.db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'progress'").Count(&count).Error
	if err != nil || count == 0 {
		// 如果检查失败或字段不存在，只更新processed_items等现有字段
		updateStart := utils.GetCurrentTime()
		err := r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
			"processed_items": progress, // 使用progress作为processed_items的近似值
		}).Error
		updateDuration := time.Since(updateStart)
		totalDuration := time.Since(startTime)
		if err != nil {
			utils.Error("UpdateProgress失败(字段不存在): ID=%d, 进度=%f, 错误=%v, 更新耗时=%v, 总耗时=%v", id, progress, err, updateDuration, totalDuration)
			return err
		}
		utils.Debug("UpdateProgress成功(字段不存在): ID=%d, 进度=%f, 更新耗时=%v, 总耗时=%v", id, progress, updateDuration, totalDuration)
		return nil
	}

	// 字段存在，正常更新
	updateStart := utils.GetCurrentTime()
	err = r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"progress":      progress,
		"progress_data": progressData,
	}).Error
	updateDuration := time.Since(updateStart)
	totalDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateProgress失败: ID=%d, 进度=%f, 错误=%v, 更新耗时=%v, 总耗时=%v", id, progress, err, updateDuration, totalDuration)
		return err
	}
	utils.Debug("UpdateProgress成功: ID=%d, 进度=%f, 更新耗时=%v, 总耗时=%v", id, progress, updateDuration, totalDuration)
	return nil
}

// UpdateStatusAndMessage 更新任务状态和消息
func (r *TaskRepositoryImpl) UpdateStatusAndMessage(id uint, status, message string) error {
	startTime := utils.GetCurrentTime()
	// 检查message字段是否存在
	var count int64
	err := r.db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'message'").Count(&count).Error
	if err != nil {
		// 如果检查失败，只更新状态
		updateStart := utils.GetCurrentTime()
		err := r.db.Model(&entity.Task{}).Where("id = ?", id).Update("status", status).Error
		updateDuration := time.Since(updateStart)
		totalDuration := time.Since(startTime)
		if err != nil {
			utils.Error("UpdateStatusAndMessage失败(检查失败): ID=%d, 状态=%s, 错误=%v, 更新耗时=%v, 总耗时=%v", id, status, err, updateDuration, totalDuration)
			return err
		}
		utils.Debug("UpdateStatusAndMessage成功(检查失败): ID=%d, 状态=%s, 更新耗时=%v, 总耗时=%v", id, status, updateDuration, totalDuration)
		return nil
	}

	if count > 0 {
		// message字段存在，更新状态和消息
		updateStart := utils.GetCurrentTime()
		err := r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":  status,
			"message": message,
		}).Error
		updateDuration := time.Since(updateStart)
		totalDuration := time.Since(startTime)
		if err != nil {
			utils.Error("UpdateStatusAndMessage失败(字段存在): ID=%d, 状态=%s, 错误=%v, 更新耗时=%v, 总耗时=%v", id, status, err, updateDuration, totalDuration)
			return err
		}
		utils.Debug("UpdateStatusAndMessage成功(字段存在): ID=%d, 状态=%s, 更新耗时=%v, 总耗时=%v", id, status, updateDuration, totalDuration)
		return nil
	} else {
		// message字段不存在，只更新状态
		updateStart := utils.GetCurrentTime()
		err := r.db.Model(&entity.Task{}).Where("id = ?", id).Update("status", status).Error
		updateDuration := time.Since(updateStart)
		totalDuration := time.Since(startTime)
		if err != nil {
			utils.Error("UpdateStatusAndMessage失败(字段不存在): ID=%d, 状态=%s, 错误=%v, 更新耗时=%v, 总耗时=%v", id, status, err, updateDuration, totalDuration)
			return err
		}
		utils.Debug("UpdateStatusAndMessage成功(字段不存在): ID=%d, 状态=%s, 更新耗时=%v, 总耗时=%v", id, status, updateDuration, totalDuration)
		return nil
	}
}

// UpdateTaskStats 更新任务统计信息
func (r *TaskRepositoryImpl) UpdateTaskStats(id uint, processed, success, failed int) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"processed_items": processed,
		"success_items":   success,
		"failed_items":    failed,
	}).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateTaskStats失败: ID=%d, 处理数=%d, 成功数=%d, 失败数=%d, 错误=%v, 更新耗时=%v", id, processed, success, failed, err, updateDuration)
		return err
	}
	utils.Debug("UpdateTaskStats成功: ID=%d, 处理数=%d, 成功数=%d, 失败数=%d, 更新耗时=%v", id, processed, success, failed, updateDuration)
	return nil
}

// UpdateStartedAt 更新任务开始时间
func (r *TaskRepositoryImpl) UpdateStartedAt(id uint) error {
	startTime := utils.GetCurrentTime()
	now := time.Now()
	err := r.db.Model(&entity.Task{}).Where("id = ?", id).Update("started_at", now).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateStartedAt失败: ID=%d, 错误=%v, 更新耗时=%v", id, err, updateDuration)
		return err
	}
	utils.Debug("UpdateStartedAt成功: ID=%d, 更新耗时=%v", id, updateDuration)
	return nil
}

// UpdateCompletedAt 更新任务完成时间
func (r *TaskRepositoryImpl) UpdateCompletedAt(id uint) error {
	startTime := utils.GetCurrentTime()
	now := time.Now()
	err := r.db.Model(&entity.Task{}).Where("id = ?", id).Update("completed_at", now).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateCompletedAt失败: ID=%d, 错误=%v, 更新耗时=%v", id, err, updateDuration)
		return err
	}
	utils.Debug("UpdateCompletedAt成功: ID=%d, 更新耗时=%v", id, updateDuration)
	return nil
}

// UpdateTotalItems 更新任务总项目数
func (r *TaskRepositoryImpl) UpdateTotalItems(id uint, totalItems int) error {
	startTime := utils.GetCurrentTime()
	err := r.db.Model(&entity.Task{}).Where("id = ?", id).Update("total_items", totalItems).Error
	updateDuration := time.Since(startTime)
	if err != nil {
		utils.Error("UpdateTotalItems失败: ID=%d, 总项目数=%d, 错误=%v, 更新耗时=%v", id, totalItems, err, updateDuration)
		return err
	}
	utils.Debug("UpdateTotalItems成功: ID=%d, 总项目数=%d, 更新耗时=%v", id, totalItems, updateDuration)
	return nil
}
