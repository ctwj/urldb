package config

import (
	"fmt"
	"sync"
)

// ConfigManager 插件配置管理器
type ConfigManager struct {
	schemas   map[string]*ConfigSchema
	templates *ConfigTemplateManager
	versions  *ConfigVersionManager
	validator *ConfigValidator
	mutex     sync.RWMutex
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		schemas:   make(map[string]*ConfigSchema),
		templates: NewConfigTemplateManager(),
		versions:  NewConfigVersionManager(10),
	}
}

// RegisterSchema 注册配置模式
func (m *ConfigManager) RegisterSchema(schema *ConfigSchema) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if schema.PluginName == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	m.schemas[schema.PluginName] = schema
	return nil
}

// GetSchema 获取配置模式
func (m *ConfigManager) GetSchema(pluginName string) (*ConfigSchema, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	schema, exists := m.schemas[pluginName]
	if !exists {
		return nil, fmt.Errorf("schema not found for plugin '%s'", pluginName)
	}

	return schema, nil
}

// ValidateConfig 验证插件配置
func (m *ConfigManager) ValidateConfig(pluginName string, config map[string]interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	schema, exists := m.schemas[pluginName]
	if !exists {
		return fmt.Errorf("schema not found for plugin '%s'", pluginName)
	}

	validator := NewConfigValidator(schema)
	return validator.Validate(config)
}

// ApplyTemplate 应用配置模板
func (m *ConfigManager) ApplyTemplate(pluginName, templateName string, config map[string]interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.templates.ApplyTemplate(templateName, config)
}

// SaveVersion 保存配置版本
func (m *ConfigManager) SaveVersion(pluginName, version, description, author string, config map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.versions.SaveVersion(pluginName, version, description, author, config)
}

// GetLatestVersion 获取最新配置版本
func (m *ConfigManager) GetLatestVersion(pluginName string) (map[string]interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	version, err := m.versions.GetLatestVersion(pluginName)
	if err != nil {
		return nil, err
	}

	// 返回配置副本
	configCopy := make(map[string]interface{})
	for k, v := range version.Config {
		configCopy[k] = v
	}

	return configCopy, nil
}

// RevertToVersion 回滚到指定版本
func (m *ConfigManager) RevertToVersion(pluginName, version string) (map[string]interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.versions.RevertToVersion(pluginName, version)
}

// ListVersions 列出配置版本
func (m *ConfigManager) ListVersions(pluginName string) ([]*ConfigVersion, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.versions.ListVersions(pluginName)
}

// RegisterTemplate 注册配置模板
func (m *ConfigManager) RegisterTemplate(template *ConfigTemplate) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.templates.RegisterTemplate(template)
}

// GetTemplate 获取配置模板
func (m *ConfigManager) GetTemplate(name string) (*ConfigTemplate, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.templates.GetTemplate(name)
}

// ListTemplates 列出所有模板
func (m *ConfigManager) ListTemplates() []*ConfigTemplate {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.templates.ListTemplates()
}

// ApplyDefaults 应用默认值
func (m *ConfigManager) ApplyDefaults(pluginName string, config map[string]interface{}) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	schema, exists := m.schemas[pluginName]
	if !exists {
		return fmt.Errorf("schema not found for plugin '%s'", pluginName)
	}

	validator := NewConfigValidator(schema)
	validator.ApplyDefaults(config)
	return nil
}