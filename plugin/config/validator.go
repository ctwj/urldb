package config

import (
	"fmt"
	"regexp"
)

// ConfigValidator 配置验证器
type ConfigValidator struct {
	schema *ConfigSchema
}

// NewConfigValidator 创建新的配置验证器
func NewConfigValidator(schema *ConfigSchema) *ConfigValidator {
	return &ConfigValidator{
		schema: schema,
	}
}

// Validate 验证配置
func (v *ConfigValidator) Validate(config map[string]interface{}) error {
	if v.schema == nil {
		return fmt.Errorf("no schema provided for validation")
	}

	return v.schema.Validate(config)
}

// ValidateField 验证单个字段
func (v *ConfigValidator) ValidateField(key string, value interface{}) error {
	if v.schema == nil {
		return fmt.Errorf("no schema provided for validation")
	}

	field, exists := v.schema.GetField(key)
	if !exists {
		return fmt.Errorf("field '%s' not found in schema", key)
	}

	return v.schema.validateField(field, value)
}

// ApplyDefaults 应用默认值
func (v *ConfigValidator) ApplyDefaults(config map[string]interface{}) {
	if v.schema == nil {
		return
	}

	for _, field := range v.schema.Fields {
		if field.Required && field.Default != nil {
			if _, exists := config[field.Key]; !exists {
				config[field.Key] = field.Default
			}
		}
	}
}

// GetSchema 获取配置模式
func (v *ConfigValidator) GetSchema() *ConfigSchema {
	return v.schema
}

// ValidatePattern 验证字符串模式
func (v *ConfigValidator) ValidatePattern(pattern, value string) error {
	if pattern == "" {
		return nil
	}

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}

	if !matched {
		return fmt.Errorf("value '%s' does not match pattern '%s'", value, pattern)
	}

	return nil
}