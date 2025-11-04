package config

import (
	"encoding/json"
	"fmt"
	"time"
)

// ConfigVersion 配置版本
type ConfigVersion struct {
	Version     string                 `json:"version"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	Description string                 `json:"description,omitempty"`
	Author      string                 `json:"author,omitempty"`
}

// ConfigVersionManager 配置版本管理器
type ConfigVersionManager struct {
	versions map[string][]*ConfigVersion // plugin_name -> versions
	maxVersions int                     // 最大保留版本数
}

// NewConfigVersionManager 创建新的配置版本管理器
func NewConfigVersionManager(maxVersions int) *ConfigVersionManager {
	if maxVersions <= 0 {
		maxVersions = 10 // 默认保留10个版本
	}

	return &ConfigVersionManager{
		versions:    make(map[string][]*ConfigVersion),
		maxVersions: maxVersions,
	}
}

// SaveVersion 保存配置版本
func (m *ConfigVersionManager) SaveVersion(pluginName, version, description, author string, config map[string]interface{}) error {
	if pluginName == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	// 创建配置副本以避免引用问题
	configCopy := make(map[string]interface{})
	for k, v := range config {
		configCopy[k] = v
	}

	configVersion := &ConfigVersion{
		Version:     version,
		Config:      configCopy,
		CreatedAt:   time.Now(),
		Description: description,
		Author:      author,
	}

	// 添加到版本历史
	if _, exists := m.versions[pluginName]; !exists {
		m.versions[pluginName] = make([]*ConfigVersion, 0)
	}

	m.versions[pluginName] = append(m.versions[pluginName], configVersion)

	// 限制版本数量
	m.limitVersions(pluginName)

	return nil
}

// GetVersion 获取指定版本的配置
func (m *ConfigVersionManager) GetVersion(pluginName, version string) (*ConfigVersion, error) {
	versions, exists := m.versions[pluginName]
	if !exists {
		return nil, fmt.Errorf("no versions found for plugin '%s'", pluginName)
	}

	for _, v := range versions {
		if v.Version == version {
			return v, nil
		}
	}

	return nil, fmt.Errorf("version '%s' not found for plugin '%s'", version, pluginName)
}

// GetLatestVersion 获取最新版本的配置
func (m *ConfigVersionManager) GetLatestVersion(pluginName string) (*ConfigVersion, error) {
	versions, exists := m.versions[pluginName]
	if !exists || len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for plugin '%s'", pluginName)
	}

	// 返回最后一个（最新）版本
	return versions[len(versions)-1], nil
}

// ListVersions 列出插件的所有配置版本
func (m *ConfigVersionManager) ListVersions(pluginName string) ([]*ConfigVersion, error) {
	versions, exists := m.versions[pluginName]
	if !exists {
		return nil, fmt.Errorf("no versions found for plugin '%s'", pluginName)
	}

	// 返回副本以避免外部修改
	result := make([]*ConfigVersion, len(versions))
	copy(result, versions)

	return result, nil
}

// RevertToVersion 回滚到指定版本
func (m *ConfigVersionManager) RevertToVersion(pluginName, version string) (map[string]interface{}, error) {
	configVersion, err := m.GetVersion(pluginName, version)
	if err != nil {
		return nil, err
	}

	// 返回配置副本
	configCopy := make(map[string]interface{})
	for k, v := range configVersion.Config {
		configCopy[k] = v
	}

	return configCopy, nil
}

// DeleteVersions 删除插件的所有版本
func (m *ConfigVersionManager) DeleteVersions(pluginName string) {
	delete(m.versions, pluginName)
}

// limitVersions 限制版本数量
func (m *ConfigVersionManager) limitVersions(pluginName string) {
	versions := m.versions[pluginName]
	if len(versions) > m.maxVersions {
		// 保留最新的maxVersions个版本
		m.versions[pluginName] = versions[len(versions)-m.maxVersions:]
	}
}

// ToJSON 将配置版本转换为JSON
func (cv *ConfigVersion) ToJSON() ([]byte, error) {
	return json.Marshal(cv)
}

// FromJSON 从JSON创建配置版本
func (cv *ConfigVersion) FromJSON(data []byte) error {
	return json.Unmarshal(data, cv)
}