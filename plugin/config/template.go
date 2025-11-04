package config

import (
	"encoding/json"
	"fmt"
)

// ConfigTemplate 配置模板
type ConfigTemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	SchemaRef   string                 `json:"schema_ref,omitempty"` // 引用的模式ID
	Version     string                 `json:"version"`
}

// ConfigTemplateManager 配置模板管理器
type ConfigTemplateManager struct {
	templates map[string]*ConfigTemplate
}

// NewConfigTemplateManager 创建新的配置模板管理器
func NewConfigTemplateManager() *ConfigTemplateManager {
	return &ConfigTemplateManager{
		templates: make(map[string]*ConfigTemplate),
	}
}

// RegisterTemplate 注册配置模板
func (m *ConfigTemplateManager) RegisterTemplate(template *ConfigTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	m.templates[template.Name] = template
	return nil
}

// GetTemplate 获取配置模板
func (m *ConfigTemplateManager) GetTemplate(name string) (*ConfigTemplate, error) {
	template, exists := m.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found", name)
	}

	return template, nil
}

// ListTemplates 列出所有模板
func (m *ConfigTemplateManager) ListTemplates() []*ConfigTemplate {
	templates := make([]*ConfigTemplate, 0, len(m.templates))
	for _, template := range m.templates {
		templates = append(templates, template)
	}
	return templates
}

// ApplyTemplate 应用模板到配置
func (m *ConfigTemplateManager) ApplyTemplate(templateName string, config map[string]interface{}) error {
	template, err := m.GetTemplate(templateName)
	if err != nil {
		return err
	}

	// 将模板配置合并到目标配置中
	for key, value := range template.Config {
		if _, exists := config[key]; !exists {
			config[key] = value
		}
	}

	return nil
}

// CreateTemplateFromConfig 从配置创建模板
func (m *ConfigTemplateManager) CreateTemplateFromConfig(name, description string, config map[string]interface{}) *ConfigTemplate {
	return &ConfigTemplate{
		Name:        name,
		Description: description,
		Config:      config,
		Version:     "1.0.0",
	}
}

// ToJSON 将模板转换为JSON
func (t *ConfigTemplate) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// FromJSON 从JSON创建模板
func (t *ConfigTemplate) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}