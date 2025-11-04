package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin/cache"
	"github.com/ctwj/urldb/plugin/concurrency"
	"github.com/ctwj/urldb/plugin/security"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// PluginContext implements the PluginContext interface
type PluginContext struct {
	pluginName        string
	manager           *Manager
	config            map[string]interface{}
	taskManager       interface{}
	repoManager       *repo.RepositoryManager
	database          *gorm.DB
	securityMgr       *security.SecurityManager
	cacheManager      *cache.CacheManager
	concurrencyCtrl   *concurrency.ConcurrencyController
}

// NewPluginContext creates a new plugin context
func NewPluginContext(pluginName string, manager *Manager, config map[string]interface{}, repoManager *repo.RepositoryManager, database *gorm.DB) *PluginContext {
	// 创建插件专用的缓存管理器，TTL为10分钟
	cacheManager := cache.NewCacheManager(10 * time.Minute)

	// 创建并发控制器，全局限制为10个并发任务
	concurrencyCtrl := concurrency.NewConcurrencyController(10)

	return &PluginContext{
		pluginName:      pluginName,
		manager:         manager,
		config:          config,
		taskManager:     manager.taskManager,
		repoManager:     repoManager,
		database:        database,
		securityMgr:     manager.securityManager,
		cacheManager:    cacheManager,
		concurrencyCtrl: concurrencyCtrl,
	}
}

// LogDebug logs a debug message
func (pc *PluginContext) LogDebug(msg string, args ...interface{}) {
	utils.Debug("[%s] %s", pc.pluginName, fmt.Sprintf(msg, args...))
}

// LogInfo logs an info message
func (pc *PluginContext) LogInfo(msg string, args ...interface{}) {
	utils.Info("[%s] %s", pc.pluginName, fmt.Sprintf(msg, args...))
	// Log activity for monitoring
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "log_info", "", map[string]interface{}{
			"message": fmt.Sprintf(msg, args...),
		})
	}
}

// LogWarn logs a warning message
func (pc *PluginContext) LogWarn(msg string, args ...interface{}) {
	utils.Warn("[%s] %s", pc.pluginName, fmt.Sprintf(msg, args...))
	// Log activity for monitoring
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "log_warn", "", map[string]interface{}{
			"message": fmt.Sprintf(msg, args...),
		})
	}
}

// LogError logs an error message
func (pc *PluginContext) LogError(msg string, args ...interface{}) {
	utils.Error("[%s] %s", pc.pluginName, fmt.Sprintf(msg, args...))
	// Log activity for monitoring
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "log_error", "", map[string]interface{}{
			"message": fmt.Sprintf(msg, args...),
		})
	}
}

// GetConfig gets a configuration value
func (pc *PluginContext) GetConfig(key string) (interface{}, error) {
	// Check permission first
	if pc.securityMgr != nil {
		hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, security.PermissionConfigRead, pc.pluginName)
		if !hasPerm {
			pc.securityMgr.LogActivity(pc.pluginName, "permission_denied", "config_read", map[string]interface{}{
				"key": key,
			})
			return nil, fmt.Errorf("permission denied: plugin %s does not have config read permission", pc.pluginName)
		}
	}

	// 首先尝试从内存缓存中获取
	if pc.config != nil {
		if value, exists := pc.config[key]; exists {
			return value, nil
		}
	}

	// 然后尝试从插件缓存中获取
	cacheKey := fmt.Sprintf("config:%s", key)
	if value, err := pc.CacheGet(cacheKey); err == nil {
		return value, nil
	}

	// 如果缓存中没有，从数据库获取
	if pc.repoManager != nil {
		config, err := pc.repoManager.PluginConfigRepository.FindByPluginAndKey(pc.pluginName, key)
		if err != nil {
			return nil, fmt.Errorf("configuration key %s not found: %v", key, err)
		}

		// 将配置存入缓存，TTL为10分钟
		if err := pc.CacheSet(cacheKey, config.ConfigValue, 10*time.Minute); err != nil {
			pc.LogWarn("Failed to cache config: %v", err)
		}

		return config.ConfigValue, nil
	}

	return nil, fmt.Errorf("no configuration available")
}

// SetConfig sets a configuration value
func (pc *PluginContext) SetConfig(key string, value interface{}) error {
	// Check permission first
	if pc.securityMgr != nil {
		hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, security.PermissionConfigWrite, pc.pluginName)
		if !hasPerm {
			pc.securityMgr.LogActivity(pc.pluginName, "permission_denied", "config_write", map[string]interface{}{
				"key": key,
			})
			return fmt.Errorf("permission denied: plugin %s does not have config write permission", pc.pluginName)
		}
	}

	// 更新内存缓存
	if pc.config == nil {
		pc.config = make(map[string]interface{})
	}
	pc.config[key] = value

	// 同步到数据库
	if pc.repoManager != nil {
		// 将值转换为字符串
		var configValue string
		switch v := value.(type) {
		case string:
			configValue = v
		case int:
			configValue = fmt.Sprintf("%d", v)
		case bool:
			if v {
				configValue = "true"
			} else {
				configValue = "false"
			}
		default:
			configValue = fmt.Sprintf("%v", v)
		}

		// 确定配置类型
		var configType string
		switch value.(type) {
		case string:
			configType = "string"
		case int:
			configType = "int"
		case bool:
			configType = "bool"
		default:
			configType = "string"
		}

		// 更新或创建配置
		err := pc.repoManager.PluginConfigRepository.Upsert(pc.pluginName, key, configValue, configType, false, "")
		if err != nil {
			return fmt.Errorf("failed to save plugin config to database: %v", err)
		}
	}

	// 清除缓存中的配置
	cacheKey := fmt.Sprintf("config:%s", key)
	if err := pc.CacheDelete(cacheKey); err != nil {
		pc.LogWarn("Failed to delete cached config: %v", err)
	}

	// Log the activity
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "config_set", key, map[string]interface{}{
			"value": value,
		})
	}

	return nil
}

// GetData gets plugin data
func (pc *PluginContext) GetData(key string, dataType string) (interface{}, error) {
	// Check permission first
	if pc.securityMgr != nil {
		hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, security.PermissionDataRead, pc.pluginName)
		if !hasPerm {
			pc.securityMgr.LogActivity(pc.pluginName, "permission_denied", "data_read", map[string]interface{}{
				"key": key,
				"type": dataType,
			})
			return nil, fmt.Errorf("permission denied: plugin %s does not have data read permission", pc.pluginName)
		}
	}

	// 首先尝试从缓存获取
	cacheKey := fmt.Sprintf("data:%s:%s", dataType, key)
	if value, err := pc.CacheGet(cacheKey); err == nil {
		pc.LogDebug("Data retrieved from cache: key=%s, type=%s", key, dataType)
		return value, nil
	}

	if pc.repoManager == nil {
		return nil, fmt.Errorf("repository manager not available")
	}

	data, err := pc.repoManager.PluginDataRepository.FindByPluginAndKey(pc.pluginName, dataType, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin data: %v", err)
	}

	// 尝试解析JSON数据
	var value interface{}
	if err := json.Unmarshal([]byte(data.DataValue), &value); err != nil {
		// 如果解析失败，返回原始字符串
		value = data.DataValue
	}

	// 将数据存入缓存，TTL为5分钟
	if err := pc.CacheSet(cacheKey, value, 5*time.Minute); err != nil {
		pc.LogWarn("Failed to cache data: %v", err)
	}

	// Log the activity
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "data_read", key, map[string]interface{}{
			"type": dataType,
		})
	}

	return value, nil
}

// SetData sets plugin data
func (pc *PluginContext) SetData(key string, value interface{}, dataType string) error {
	// Check permission first
	if pc.securityMgr != nil {
		hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, security.PermissionDataWrite, pc.pluginName)
		if !hasPerm {
			pc.securityMgr.LogActivity(pc.pluginName, "permission_denied", "data_write", map[string]interface{}{
				"key": key,
				"type": dataType,
			})
			return fmt.Errorf("permission denied: plugin %s does not have data write permission", pc.pluginName)
		}
	}

	if pc.repoManager == nil {
		return fmt.Errorf("repository manager not available")
	}

	// 将值转换为JSON字符串
	var dataValue string
	switch v := value.(type) {
	case string:
		dataValue = v
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal data to JSON: %v", err)
		}
		dataValue = string(jsonData)
	}

	// 检查数据是否已存在
	existingData, err := pc.repoManager.PluginDataRepository.FindByPluginAndKey(pc.pluginName, dataType, key)
	if err != nil {
		// 数据不存在，创建新记录
		newData := &entity.PluginData{
			PluginName: pc.pluginName,
			DataType:   dataType,
			DataKey:    key,
			DataValue:  dataValue,
		}
		err = pc.repoManager.PluginDataRepository.Create(newData)
		if err != nil {
			return fmt.Errorf("failed to create plugin data: %v", err)
		}
	} else {
		// 数据已存在，更新记录
		existingData.DataValue = dataValue
		err = pc.repoManager.PluginDataRepository.Update(existingData)
		if err != nil {
			return fmt.Errorf("failed to update plugin data: %v", err)
		}
	}

	// 清除缓存中的旧数据
	cacheKey := fmt.Sprintf("data:%s:%s", dataType, key)
	if err := pc.CacheDelete(cacheKey); err != nil {
		pc.LogWarn("Failed to delete cached data: %v", err)
	}

	// Log the activity
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "data_write", key, map[string]interface{}{
			"type": dataType,
			"value": value,
		})
	}

	pc.LogInfo("Data set successfully: key=%s, type=%s", key, dataType)
	return nil
}

// DeleteData deletes plugin data
func (pc *PluginContext) DeleteData(key string, dataType string) error {
	// Check permission first
	if pc.securityMgr != nil {
		hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, security.PermissionDataWrite, pc.pluginName)
		if !hasPerm {
			pc.securityMgr.LogActivity(pc.pluginName, "permission_denied", "data_delete", map[string]interface{}{
				"key": key,
				"type": dataType,
			})
			return fmt.Errorf("permission denied: plugin %s does not have data delete permission", pc.pluginName)
		}
	}

	if pc.repoManager == nil {
		return fmt.Errorf("repository manager not available")
	}

	err := pc.repoManager.PluginDataRepository.DeleteByPluginAndKey(pc.pluginName, dataType, key)
	if err != nil {
		return fmt.Errorf("failed to delete plugin data: %v", err)
	}

	// 清除缓存中的数据
	cacheKey := fmt.Sprintf("data:%s:%s", dataType, key)
	if err := pc.CacheDelete(cacheKey); err != nil {
		pc.LogWarn("Failed to delete cached data: %v", err)
	}

	// Log the activity
	if pc.securityMgr != nil {
		pc.securityMgr.LogActivity(pc.pluginName, "data_delete", key, map[string]interface{}{
			"type": dataType,
		})
	}

	pc.LogInfo("Data deleted successfully: key=%s, type=%s", key, dataType)
	return nil
}

// RegisterTask registers a task with the task manager
func (pc *PluginContext) RegisterTask(name string, taskFunc func()) error {
	// Check permission first
	if pc.securityMgr != nil {
		hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, security.PermissionTaskSchedule, pc.pluginName)
		if !hasPerm {
			pc.securityMgr.LogActivity(pc.pluginName, "permission_denied", "task_schedule", map[string]interface{}{
				"name": name,
			})
			return fmt.Errorf("permission denied: plugin %s does not have task schedule permission", pc.pluginName)
		}
	}

	if pc.taskManager == nil {
		return fmt.Errorf("task manager not available")
	}

	// 创建任务处理器
	processor := task.NewPluginTaskProcessor(pc.pluginName, name, func(ctx context.Context, taskID uint, item *entity.TaskItem) error {
		// Log task execution start
		startTime := time.Now()
		if pc.securityMgr != nil {
			pc.securityMgr.LogActivity(pc.pluginName, "task_start", name, map[string]interface{}{
				"task_id": taskID,
			})
		}

		// 执行插件任务
		taskFunc()

		// Log task execution end
		duration := time.Since(startTime)
		if pc.securityMgr != nil {
			pc.securityMgr.LogExecutionTime(pc.pluginName, "task_execute", name, duration)
		}

		return nil
	})

	// 注册处理器到任务管理器
	if taskManager, ok := pc.taskManager.(*task.TaskManager); ok {
		taskManager.RegisterProcessor(processor)
		pc.LogInfo("Task registered: %s", name)

		// Log the activity
		if pc.securityMgr != nil {
			pc.securityMgr.LogActivity(pc.pluginName, "task_register", name, nil)
		}

		return nil
	}

	return fmt.Errorf("invalid task manager type")
}

// UnregisterTask unregisters a task from the task manager
func (pc *PluginContext) UnregisterTask(name string) error {
	// In a real implementation, this would unregister the task from the task manager
	// For now, we'll just log the operation
	pc.LogInfo("Unregistering task: %s", name)
	return nil
}

// GetDB returns the database connection
func (pc *PluginContext) GetDB() interface{} {
	return pc.database
}

// CheckPermission checks if the plugin has the specified permission
func (pc *PluginContext) CheckPermission(permissionType string, resource ...string) (bool, error) {
	if pc.securityMgr == nil {
		return false, fmt.Errorf("security manager not available")
	}

	permType := security.PermissionType(permissionType)
	hasPerm, _ := pc.securityMgr.CheckPermission(pc.pluginName, permType, resource...)
	return hasPerm, nil
}

// RequestPermission requests a permission for the plugin
func (pc *PluginContext) RequestPermission(permissionType string, resource string) error {
	if pc.securityMgr == nil {
		return fmt.Errorf("security manager not available")
	}

	permType := security.PermissionType(permissionType)
	permission := security.Permission{
		Type:     permType,
		Resource: resource,
		Allowed:  false, // Default to not allowed, requires manual approval
	}

	// Log the request for review
	pc.LogInfo("Permission request: type=%s, resource=%s", permissionType, resource)
	pc.securityMgr.LogActivity(pc.pluginName, "permission_request", resource, map[string]interface{}{
		"type": permissionType,
		"resource": resource,
	})

	return pc.securityMgr.GrantPermission(pc.pluginName, permission)
}

// GetSecurityReport returns a security report for the plugin
func (pc *PluginContext) GetSecurityReport() (interface{}, error) {
	if pc.securityMgr == nil {
		return nil, fmt.Errorf("security manager not available")
	}

	report := pc.securityMgr.CreateSecurityReport(pc.pluginName)
	return report, nil
}

// CacheSet 设置缓存项
func (pc *PluginContext) CacheSet(key string, value interface{}, ttl time.Duration) error {
	if pc.cacheManager == nil {
		return fmt.Errorf("cache manager not available")
	}

	// 构造带插件名称的缓存键，避免键冲突
	cacheKey := fmt.Sprintf("plugin:%s:%s", pc.pluginName, key)
	pc.cacheManager.Set(cacheKey, value, ttl)
	return nil
}

// CacheGet 获取缓存项
func (pc *PluginContext) CacheGet(key string) (interface{}, error) {
	if pc.cacheManager == nil {
		return nil, fmt.Errorf("cache manager not available")
	}

	// 构造带插件名称的缓存键
	cacheKey := fmt.Sprintf("plugin:%s:%s", pc.pluginName, key)
	value, exists := pc.cacheManager.Get(cacheKey)
	if !exists {
		return nil, fmt.Errorf("cache key %s not found", key)
	}

	return value, nil
}

// CacheDelete 删除缓存项
func (pc *PluginContext) CacheDelete(key string) error {
	if pc.cacheManager == nil {
		return fmt.Errorf("cache manager not available")
	}

	// 构造带插件名称的缓存键
	cacheKey := fmt.Sprintf("plugin:%s:%s", pc.pluginName, key)
	pc.cacheManager.Delete(cacheKey)
	return nil
}

// CacheClear 清空插件的所有缓存
func (pc *PluginContext) CacheClear() error {
	if pc.cacheManager == nil {
		return fmt.Errorf("cache manager not available")
	}

	// 注意：这里我们不能直接清空整个缓存管理器，因为它是所有插件共享的
	// 在实际实现中，可能需要更复杂的机制来清空特定插件的缓存
	pc.LogWarn("CacheClear is not implemented due to shared cache manager")
	return nil
}

// ConcurrencyExecute 在并发控制下执行任务
func (pc *PluginContext) ConcurrencyExecute(ctx context.Context, taskFunc func() error) error {
	if pc.concurrencyCtrl == nil {
		return fmt.Errorf("concurrency controller not available")
	}

	return pc.concurrencyCtrl.Execute(ctx, pc.pluginName, taskFunc)
}

// SetConcurrencyLimit 设置插件的并发限制
func (pc *PluginContext) SetConcurrencyLimit(limit int) error {
	if pc.concurrencyCtrl == nil {
		return fmt.Errorf("concurrency controller not available")
	}

	// 通过管理器设置插件限制
	if pc.manager != nil {
		pc.manager.SetPluginConcurrencyLimit(pc.pluginName, limit)
	} else {
		// 如果没有管理器引用，直接在本地控制器上设置
		pc.concurrencyCtrl.SetPluginLimit(pc.pluginName, limit)
	}

	return nil
}

// GetConcurrencyStats 获取并发控制统计信息
func (pc *PluginContext) GetConcurrencyStats() (map[string]interface{}, error) {
	if pc.concurrencyCtrl == nil {
		return nil, fmt.Errorf("concurrency controller not available")
	}

	stats := pc.concurrencyCtrl.GetStats()
	return stats, nil
}