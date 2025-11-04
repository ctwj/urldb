package config

import (
	"encoding/json"
	"fmt"
)

// ConfigField 定义配置字段的结构
type ConfigField struct {
	Key         string      `json:"key"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"` // string, int, bool, float, json
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Min         *float64    `json:"min,omitempty"`
	Max         *float64    `json:"max,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
	Encrypted   bool        `json:"encrypted,omitempty"`
}

// ConfigSchema 定义插件配置模式
type ConfigSchema struct {
	PluginName string        `json:"plugin_name"`
	Version    string        `json:"version"`
	Fields     []ConfigField `json:"fields"`
}

// NewConfigSchema 创建新的配置模式
func NewConfigSchema(pluginName, version string) *ConfigSchema {
	return &ConfigSchema{
		PluginName: pluginName,
		Version:    version,
		Fields:     make([]ConfigField, 0),
	}
}

// AddField 添加配置字段
func (s *ConfigSchema) AddField(field ConfigField) {
	s.Fields = append(s.Fields, field)
}

// GetField 获取配置字段
func (s *ConfigSchema) GetField(key string) (*ConfigField, bool) {
	for i := range s.Fields {
		if s.Fields[i].Key == key {
			return &s.Fields[i], true
		}
	}
	return nil, false
}

// ToJSON 将配置模式转换为JSON
func (s *ConfigSchema) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON 从JSON创建配置模式
func (s *ConfigSchema) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}

// Validate 验证配置是否符合模式
func (s *ConfigSchema) Validate(config map[string]interface{}) error {
	// 验证必需字段
	for _, field := range s.Fields {
		if field.Required {
			if _, exists := config[field.Key]; !exists {
				// 检查是否有默认值
				if field.Default == nil {
					return fmt.Errorf("required field '%s' is missing", field.Key)
				}
				// 设置默认值
				config[field.Key] = field.Default
			}
		}
	}

	// 验证字段类型和值
	for key, value := range config {
		field, exists := s.GetField(key)
		if !exists {
			// 允许额外字段存在，但不验证
			continue
		}

		if err := s.validateField(field, value); err != nil {
			return fmt.Errorf("field '%s': %v", key, err)
		}
	}

	return nil
}

// validateField 验证单个字段
func (s *ConfigSchema) validateField(field *ConfigField, value interface{}) error {
	switch field.Type {
	case "string":
		if str, ok := value.(string); ok {
			if field.Pattern != "" {
				// 这里可以添加正则表达式验证
			}
			if field.Enum != nil && len(field.Enum) > 0 {
				found := false
				for _, enumValue := range field.Enum {
					if str == enumValue {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("value '%s' is not in enum %v", str, field.Enum)
				}
			}
		} else {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "int":
		if num, ok := value.(int); ok {
			if field.Min != nil && float64(num) < *field.Min {
				return fmt.Errorf("value %d is less than minimum %f", num, *field.Min)
			}
			if field.Max != nil && float64(num) > *field.Max {
				return fmt.Errorf("value %d is greater than maximum %f", num, *field.Max)
			}
		} else if num, ok := value.(float64); ok {
			// 允许float64转换为int
			intVal := int(num)
			if field.Min != nil && float64(intVal) < *field.Min {
				return fmt.Errorf("value %d is less than minimum %f", intVal, *field.Min)
			}
			if field.Max != nil && float64(intVal) > *field.Max {
				return fmt.Errorf("value %d is greater than maximum %f", intVal, *field.Max)
			}
		} else {
			return fmt.Errorf("expected integer, got %T", value)
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "float":
		if num, ok := value.(float64); ok {
			if field.Min != nil && num < *field.Min {
				return fmt.Errorf("value %f is less than minimum %f", num, *field.Min)
			}
			if field.Max != nil && num > *field.Max {
				return fmt.Errorf("value %f is greater than maximum %f", num, *field.Max)
			}
		} else {
			return fmt.Errorf("expected float, got %T", value)
		}
	case "json":
		// JSON类型可以是任何结构，只需要能序列化为JSON
		if _, err := json.Marshal(value); err != nil {
			return fmt.Errorf("invalid JSON value: %v", err)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type)
	}

	return nil
}