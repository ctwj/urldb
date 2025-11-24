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

	// Google索引专用方法
	GetDistinctProcessedURLs() ([]string, error)
	GetLatestURLStatus(url string) (*entity.TaskItem, error)
	UpsertURLStatusRecords(taskID uint, urlResults []*URLStatusResult) error
	CleanupOldRecords() error
}

// URLStatusResult 用于批量处理的结果
type URLStatusResult struct {
	URL            string
	IndexStatus    string
	InspectResult  string
	MobileFriendly bool
	StatusCode     int
	LastCrawled    *time.Time
	ErrorMessage   string
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

// GetDistinctProcessedURLs 获取所有已处理的URL（去重）
func (r *TaskItemRepositoryImpl) GetDistinctProcessedURLs() ([]string, error) {
	startTime := utils.GetCurrentTime()
	var urls []string

	// 只返回成功处理的URL，避免处理失败的URL重复尝试
	err := r.db.Model(&entity.TaskItem{}).
		Where("status = ? AND url != ?", "completed", "").
		Distinct("url").
		Pluck("url", &urls).Error

	queryDuration := time.Since(startTime)
	if err != nil {
		utils.Error("GetDistinctProcessedURLs失败: 错误=%v, 查询耗时=%v", err, queryDuration)
		return nil, err
	}

	utils.Debug("GetDistinctProcessedURLs成功: URL数量=%d, 查询耗时=%v", len(urls), queryDuration)
	return urls, nil
}

// GetLatestURLStatus 获取URL的最新处理状态
func (r *TaskItemRepositoryImpl) GetLatestURLStatus(url string) (*entity.TaskItem, error) {
	startTime := utils.GetCurrentTime()
	var item entity.TaskItem

	err := r.db.Where("url = ?", url).
		Order("created_at DESC").
		First(&item).Error

	queryDuration := time.Since(startTime)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Debug("GetLatestURLStatus: URL未找到=%s, 查询耗时=%v", url, queryDuration)
			return nil, nil
		}
		utils.Error("GetLatestURLStatus失败: URL=%s, 错误=%v, 查询耗时=%v", url, err, queryDuration)
		return nil, err
	}

	utils.Debug("GetLatestURLStatus成功: URL=%s, 状态=%s, 查询耗时=%v", url, item.Status, queryDuration)
	return &item, nil
}

// UpsertURLStatusRecords 批量创建或更新URL状态
func (r *TaskItemRepositoryImpl) UpsertURLStatusRecords(taskID uint, urlResults []*URLStatusResult) error {
	startTime := utils.GetCurrentTime()

	if len(urlResults) == 0 {
		return nil
	}

	// 批量操作，减少数据库查询次数
	for _, result := range urlResults {
		// 查找现有记录
		existing, err := r.GetLatestURLStatus(result.URL)

		if err != nil && err != gorm.ErrRecordNotFound {
			utils.Error("UpsertURLStatusRecords查询失败: URL=%s, 错误=%v", result.URL, err)
			continue
		}

		now := time.Now()

		if existing != nil && existing.ID > 0 {
			// 更新现有记录（只更新状态变化的）
			if existing.IndexStatus != result.IndexStatus || existing.StatusCode != result.StatusCode {
				existing.IndexStatus = result.IndexStatus
				existing.InspectResult = result.InspectResult
				existing.MobileFriendly = result.MobileFriendly
				existing.StatusCode = result.StatusCode
				existing.LastCrawled = result.LastCrawled
				existing.ErrorMessage = result.ErrorMessage
				existing.ProcessedAt = &now

				if err := r.Update(existing); err != nil {
					utils.Error("UpsertURLStatusRecords更新失败: URL=%s, 错误=%v", result.URL, err)
					continue
				}
			}
		} else {
			// 创建新记录
			newItem := &entity.TaskItem{
				TaskID:         taskID,
				URL:            result.URL,
				Status:         "completed",
				IndexStatus:    result.IndexStatus,
				InspectResult:  result.InspectResult,
				MobileFriendly: result.MobileFriendly,
				StatusCode:     result.StatusCode,
				LastCrawled:    result.LastCrawled,
				ErrorMessage:   result.ErrorMessage,
				ProcessedAt:    &now,
			}

			if err := r.Create(newItem); err != nil {
				utils.Error("UpsertURLStatusRecords创建失败: URL=%s, 错误=%v", result.URL, err)
				continue
			}
		}
	}

	totalDuration := time.Since(startTime)
	utils.Info("UpsertURLStatusRecords完成: 数量=%d, 耗时=%v", len(urlResults), totalDuration)
	return nil
}

// CleanupOldRecords 清理旧记录，保留每个URL的最新记录
func (r *TaskItemRepositoryImpl) CleanupOldRecords() error {
	startTime := utils.GetCurrentTime()

	// 1. 找出每个URL的最新记录ID
	var latestIDs []uint
	err := r.db.Table("task_items").
		Select("MAX(id) as id").
		Where("url != '' AND status = ?", "completed").
		Group("url").
		Pluck("id", &latestIDs).Error

	if err != nil {
		utils.Error("CleanupOldRecords获取最新ID失败: 错误=%v", err)
		return err
	}

	// 2. 删除所有非最新的已完成记录
	deleteResult := r.db.Where("status = ? AND id NOT IN (?)", "completed", latestIDs).
		Delete(&entity.TaskItem{})

	if deleteResult.Error != nil {
		utils.Error("CleanupOldRecords删除旧记录失败: 错误=%v", deleteResult.Error)
		return deleteResult.Error
	}

	// 3. 清理失败的旧记录（保留1周）
	failureCutoff := time.Now().AddDate(0, 0, -7)
	failureDeleteResult := r.db.Where("status = ? AND created_at < ?", "failed", failureCutoff).
		Delete(&entity.TaskItem{})

	if failureDeleteResult.Error != nil {
		utils.Error("CleanupOldRecords删除失败记录失败: 错误=%v", failureDeleteResult.Error)
		return failureDeleteResult.Error
	}

	totalDuration := time.Since(startTime)
	utils.Info("CleanupOldRecords完成: 删除完成记录=%d, 删除失败记录=%d, 耗时=%v",
		deleteResult.RowsAffected, failureDeleteResult.RowsAffected, totalDuration)

	return nil
}
