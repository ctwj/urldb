package repo

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// PluginConfigRepository 插件配置仓库
type PluginConfigRepository struct {
	db *gorm.DB
}

// NewPluginConfigRepository 创建插件配置仓库
func NewPluginConfigRepository(db *gorm.DB) *PluginConfigRepository {
	return &PluginConfigRepository{db: db}
}

// GetConfig 获取插件配置
func (r *PluginConfigRepository) GetConfig(pluginName string) (*entity.PluginConfig, error) {
	var config entity.PluginConfig
	err := r.db.Where("plugin_name = ?", pluginName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// SetConfig 设置插件配置
func (r *PluginConfigRepository) SetConfig(pluginName string, configData map[string]interface{}) error {
	// 序列化配置数据
	configJSON, err := json.Marshal(configData)
	if err != nil {
		return err
	}

	// 查找现有配置
	var existingConfig entity.PluginConfig
	err = r.db.Where("plugin_name = ?", pluginName).First(&existingConfig).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新配置
		newConfig := entity.PluginConfig{
			PluginName: pluginName,
			ConfigJSON: string(configJSON),
			Enabled:    true,
		}
		return r.db.Create(&newConfig).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有配置
		existingConfig.ConfigJSON = string(configJSON)
		return r.db.Save(&existingConfig).Error
	}
}

// SetEnabled 设置插件启用状态
func (r *PluginConfigRepository) SetEnabled(pluginName string, enabled bool) error {
	// 查找现有配置
	var existingConfig entity.PluginConfig
	err := r.db.Where("plugin_name = ?", pluginName).First(&existingConfig).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新配置记录
		newConfig := entity.PluginConfig{
			PluginName: pluginName,
			ConfigJSON: "{}",
			Enabled:    enabled,
		}
		return r.db.Create(&newConfig).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有配置
		return r.db.Model(&existingConfig).Update("enabled", enabled).Error
	}
}

// GetAllConfigs 获取所有插件配置
func (r *PluginConfigRepository) GetAllConfigs() ([]entity.PluginConfig, error) {
	var configs []entity.PluginConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// GetEnabledPlugins 获取启用的插件列表
func (r *PluginConfigRepository) GetEnabledPlugins() ([]string, error) {
	var pluginNames []string
	err := r.db.Model(&entity.PluginConfig{}).
		Where("enabled = ?", true).
		Pluck("plugin_name", &pluginNames).Error
	return pluginNames, err
}

// DeleteConfig 删除插件配置
func (r *PluginConfigRepository) DeleteConfig(pluginName string) error {
	return r.db.Where("plugin_name = ?", pluginName).Delete(&entity.PluginConfig{}).Error
}

// PluginLogRepository 插件日志仓库
type PluginLogRepository struct {
	db *gorm.DB
}

// NewPluginLogRepository 创建插件日志仓库
func NewPluginLogRepository(db *gorm.DB) *PluginLogRepository {
	return &PluginLogRepository{db: db}
}

// CreateLog 创建插件日志
func (r *PluginLogRepository) CreateLog(log *entity.PluginLog) error {
	return r.db.Create(log).Error
}

// GetLogs 获取插件日志
func (r *PluginLogRepository) GetLogs(pluginName string, page, limit int) ([]entity.PluginLog, int64, error) {
	var logs []entity.PluginLog
	var total int64

	// 获取总数
	err := r.db.Model(&entity.PluginLog{}).
		Where("plugin_name = ?", pluginName).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * limit
	err = r.db.Where("plugin_name = ?", pluginName).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

// GetRecentLogs 获取最近的插件日志
func (r *PluginLogRepository) GetRecentLogs(pluginName string, limit int) ([]entity.PluginLog, error) {
	var logs []entity.PluginLog
	err := r.db.Where("plugin_name = ?", pluginName).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetErrorLogs 获取错误日志
func (r *PluginLogRepository) GetErrorLogs(pluginName string, page, limit int) ([]entity.PluginLog, int64, error) {
	var logs []entity.PluginLog
	var total int64

	// 获取总数
	err := r.db.Model(&entity.PluginLog{}).
		Where("plugin_name = ? AND success = ?", pluginName, false).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * limit
	err = r.db.Where("plugin_name = ? AND success = ?", pluginName, false).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

// DeleteOldLogs 删除旧日志
func (r *PluginLogRepository) DeleteOldLogs(olderThan time.Time) error {
	return r.db.Where("created_at < ?", olderThan).
		Delete(&entity.PluginLog{}).Error
}

// GetExecutionStats 获取执行统计
func (r *PluginLogRepository) GetExecutionStats(pluginName string, timeRange string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 根据时间范围计算时间过滤
	var timeFilter time.Time
	switch timeRange {
	case "1h":
		timeFilter = time.Now().Add(-1 * time.Hour)
	case "24h":
		timeFilter = time.Now().Add(-24 * time.Hour)
	case "7d":
		timeFilter = time.Now().Add(-7 * 24 * time.Hour)
	default:
		timeFilter = time.Now().Add(-24 * time.Hour)
	}

	// 总执行次数
	var totalExecutions int64
	r.db.Model(&entity.PluginLog{}).
		Where("plugin_name = ? AND created_at >= ?", pluginName, timeFilter).
		Count(&totalExecutions)

	// 成功执行次数
	var successExecutions int64
	r.db.Model(&entity.PluginLog{}).
		Where("plugin_name = ? AND success = ? AND created_at >= ?", pluginName, true, timeFilter).
		Count(&successExecutions)

	// 平均执行时间
	var avgExecutionTime sql.NullFloat64
	r.db.Model(&entity.PluginLog{}).
		Where("plugin_name = ? AND success = ? AND created_at >= ?", pluginName, true, timeFilter).
		Select("AVG(execution_time)").
		Scan(&avgExecutionTime)

	// 计算成功率
	successRate := float64(0)
	if totalExecutions > 0 {
		successRate = float64(successExecutions) / float64(totalExecutions) * 100
	}

	stats["total_executions"] = totalExecutions
	stats["success_executions"] = successExecutions
	stats["success_rate"] = successRate
	stats["average_time"] = avgExecutionTime.Float64

	return stats, nil
}

// CronJobRepository 定时任务仓库
type CronJobRepository struct {
	db *gorm.DB
}

// NewCronJobRepository 创建定时任务仓库
func NewCronJobRepository(db *gorm.DB) *CronJobRepository {
	return &CronJobRepository{db: db}
}

// CreateJob 创建定时任务
func (r *CronJobRepository) CreateJob(job *entity.CronJob) error {
	return r.db.Create(job).Error
}

// GetJob 获取定时任务
func (r *CronJobRepository) GetJob(jobName string) (*entity.CronJob, error) {
	var job entity.CronJob
	err := r.db.Where("job_name = ?", jobName).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateJob 更新定时任务
func (r *CronJobRepository) UpdateJob(job *entity.CronJob) error {
	return r.db.Save(job).Error
}

// GetEnabledJobs 获取启用的定时任务
func (r *CronJobRepository) GetEnabledJobs() ([]entity.CronJob, error) {
	var jobs []entity.CronJob
	err := r.db.Where("enabled = ?", true).Find(&jobs).Error
	return jobs, err
}

// DeleteJob 删除定时任务
func (r *CronJobRepository) DeleteJob(jobName string) error {
	return r.db.Where("job_name = ?", jobName).Delete(&entity.CronJob{}).Error
}

// UpdateJobRunTime 更新任务运行时间
func (r *CronJobRepository) UpdateJobRunTime(jobName string, lastRun, nextRun *time.Time) error {
	return r.db.Model(&entity.CronJob{}).
		Where("job_name = ?", jobName).
		Updates(map[string]interface{}{
			"last_run": lastRun,
			"next_run": nextRun,
		}).Error
}