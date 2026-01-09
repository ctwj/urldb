package service

import (
	"strconv"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// ConfigManager 接口定义配置管理器的通用方法
type ConfigManager interface {
	GetConfigValue(key string) (string, error)
	GetConfigString(key string) (string, error)
	GetConfigBool(key string) (bool, error)
	GetConfigInt(key string) (int, error)
	GetConfigInt64(key string) (int64, error)
	GetConfigFloat64(key string) (float64, error)
	SetConfig(key, value string) error
	LoadAllConfigs() error
	RefreshConfigCache() error
}

// RepositoryConfigManager 是基于 RepositoryManager 的配置管理器实现
type RepositoryConfigManager struct {
	repoManager *repo.RepositoryManager
}

// NewRepositoryConfigManager 创建新的 RepositoryConfigManager
func NewRepositoryConfigManager(repoManager *repo.RepositoryManager) *RepositoryConfigManager {
	return &RepositoryConfigManager{
		repoManager: repoManager,
	}
}

// GetConfigValue 获取配置值
func (rcm *RepositoryConfigManager) GetConfigValue(key string) (string, error) {
	return rcm.repoManager.SystemConfigRepository.GetConfigValue(key)
}

// GetConfigBool 获取布尔值配置
func (rcm *RepositoryConfigManager) GetConfigBool(key string) (bool, error) {
	return rcm.repoManager.SystemConfigRepository.GetConfigBool(key)
}

// GetConfigInt 获取整数值配置
func (rcm *RepositoryConfigManager) GetConfigInt(key string) (int, error) {
	return rcm.repoManager.SystemConfigRepository.GetConfigInt(key)
}

// GetConfigInt64 获取64位整数值配置
func (rcm *RepositoryConfigManager) GetConfigInt64(key string) (int64, error) {
	value, err := rcm.GetConfigValue(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

// GetConfigFloat64 获取浮点数配置
func (rcm *RepositoryConfigManager) GetConfigFloat64(key string) (float64, error) {
	value, err := rcm.GetConfigValue(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(value, 64)
}

// GetConfigString 获取字符串配置值
func (rcm *RepositoryConfigManager) GetConfigString(key string) (string, error) {
	return rcm.GetConfigValue(key)
}

// SetConfig 设置配置
func (rcm *RepositoryConfigManager) SetConfig(key, value string) error {
	configs := []entity.SystemConfig{
		{
			Key:   key,
			Value: value,
			Type:  "string", // 默认类型，实际应该根据键来确定
		},
	}
	return rcm.repoManager.SystemConfigRepository.UpsertConfigs(configs)
}

// LoadAllConfigs 加载所有配置
func (rcm *RepositoryConfigManager) LoadAllConfigs() error {
	return nil // RepositoryManager本身不提供这个功能
}

// RefreshConfigCache 刷新配置缓存
func (rcm *RepositoryConfigManager) RefreshConfigCache() error {
	return nil // RepositoryManager本身不提供这个功能
}