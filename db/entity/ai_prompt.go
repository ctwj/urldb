package entity

import (
	"time"
)

// AIPrompt AI提示词实体
type AIPrompt struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:100;not null;comment:提示词名称"`
	Type          string    `json:"type" gorm:"size:50;not null;uniqueIndex;comment:提示词类型"`
	SystemContent string    `json:"system_content" gorm:"type:text;comment:系统提示词内容"`
	UserContent   string    `json:"user_content" gorm:"type:text;comment:用户提示词模板"`
	Description   string    `json:"description" gorm:"size:500;comment:提示词描述"`
	Variables     string    `json:"variables" gorm:"type:text;comment:支持的变量列表(JSON格式)"`
	IsActive      bool      `json:"is_active" gorm:"default:true;comment:是否启用"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// PromptType 提示词类型常量
const (
	PromptTypeContentGeneration = "content_generation" // 内容生成
	PromptTypeClassification     = "classification"     // 分类推荐
	PromptTypeToolSystem         = "tool_system"         // 工具调用系统
	PromptTypeQATemplate         = "qa_template"         // 问答模板
	PromptTypeAnalysisTemplate   = "analysis_template"   // 文本分析模板
)

// GetPromptTypeName 获取提示词类型的显示名称
func GetPromptTypeName(promptType string) string {
	switch promptType {
	case PromptTypeContentGeneration:
		return "内容生成"
	case PromptTypeClassification:
		return "分类推荐"
	case PromptTypeToolSystem:
		return "工具调用系统"
	case PromptTypeQATemplate:
		return "问答模板"
	case PromptTypeAnalysisTemplate:
		return "文本分析模板"
	default:
		return "未知类型"
	}
}

// GetPromptTypeDescription 获取提示词类型的使用位置
func GetPromptTypeDescription(promptType string) string {
	switch promptType {
	case PromptTypeContentGeneration:
		return "数据管理 → 资源管理 → 智能处理/优化标题按钮"
	case PromptTypeClassification:
		return "数据管理 → 资源管理 → 智能分类按钮"
	case PromptTypeToolSystem:
		return "系统配置 → AI配置 → AI聊天/验证/测试按钮"
	case PromptTypeQATemplate:
		return "后端API接口 → /api/ai/ask（前端未直接使用）"
	case PromptTypeAnalysisTemplate:
		return "后端API接口 → /api/ai/analyze（前端未直接使用）"
	default:
		return "未定义的使用位置"
	}
}

// GetAllPromptTypes 获取所有提示词类型
func GetAllPromptTypes() []string {
	return []string{
		PromptTypeContentGeneration,
		PromptTypeClassification,
		PromptTypeToolSystem,
		PromptTypeQATemplate,
		PromptTypeAnalysisTemplate,
	}
}

// TableName 指定表名
func (AIPrompt) TableName() string {
	return "ai_prompts"
}