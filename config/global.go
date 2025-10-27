package config

import (
	"sync"
)

var (
	globalConfigManager *ConfigManager
	once                sync.Once
)

// SetGlobalConfigManager 设置全局配置管理器
func SetGlobalConfigManager(cm *ConfigManager) {
	globalConfigManager = cm
}

// GetGlobalConfigManager 获取全局配置管理器
func GetGlobalConfigManager() *ConfigManager {
	return globalConfigManager
}

// GetConfig 获取配置值（全局函数）
func GetConfig(key string) (*ConfigItem, error) {
	if globalConfigManager == nil {
		return nil, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfig(key)
}

// GetConfigValue 获取配置值（全局函数）
func GetConfigValue(key string) (string, error) {
	if globalConfigManager == nil {
		return "", ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigValue(key)
}

// GetConfigBool 获取布尔配置值（全局函数）
func GetConfigBool(key string) (bool, error) {
	if globalConfigManager == nil {
		return false, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigBool(key)
}

// GetConfigInt 获取整数配置值（全局函数）
func GetConfigInt(key string) (int, error) {
	if globalConfigManager == nil {
		return 0, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigInt(key)
}

// GetConfigInt64 获取64位整数配置值（全局函数）
func GetConfigInt64(key string) (int64, error) {
	if globalConfigManager == nil {
		return 0, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigInt64(key)
}

// GetConfigFloat64 获取浮点数配置值（全局函数）
func GetConfigFloat64(key string) (float64, error) {
	if globalConfigManager == nil {
		return 0, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigFloat64(key)
}

// SetConfig 设置配置值（全局函数）
func SetConfig(key, value string) error {
	if globalConfigManager == nil {
		return ErrConfigManagerNotInitialized
	}
	return globalConfigManager.SetConfig(key, value)
}

// SetConfigWithType 设置配置值（指定类型，全局函数）
func SetConfigWithType(key, value, configType string) error {
	if globalConfigManager == nil {
		return ErrConfigManagerNotInitialized
	}
	return globalConfigManager.SetConfigWithType(key, value, configType)
}

// GetConfigWithEnvFallback 获取配置值（环境变量优先，全局函数）
func GetConfigWithEnvFallback(configKey, envKey string) (string, error) {
	if globalConfigManager == nil {
		return "", ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigWithEnvFallback(configKey, envKey)
}

// GetConfigIntWithEnvFallback 获取整数配置值（环境变量优先，全局函数）
func GetConfigIntWithEnvFallback(configKey, envKey string) (int, error) {
	if globalConfigManager == nil {
		return 0, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigIntWithEnvFallback(configKey, envKey)
}

// GetConfigBoolWithEnvFallback 获取布尔配置值（环境变量优先，全局函数）
func GetConfigBoolWithEnvFallback(configKey, envKey string) (bool, error) {
	if globalConfigManager == nil {
		return false, ErrConfigManagerNotInitialized
	}
	return globalConfigManager.GetConfigBoolWithEnvFallback(configKey, envKey)
}

// ErrConfigManagerNotInitialized 配置管理器未初始化错误
var ErrConfigManagerNotInitialized = &ConfigError{
	Code:    "CONFIG_MANAGER_NOT_INITIALIZED",
	Message: "配置管理器未初始化",
}

// ConfigError 配置错误
type ConfigError struct {
	Code    string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}