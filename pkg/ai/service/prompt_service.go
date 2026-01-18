package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// PromptService 提示词服务
type PromptService struct {
	db *gorm.DB
}

// NewPromptService 创建提示词服务实例
func NewPromptService(db *gorm.DB) *PromptService {
	return &PromptService{
		db: db,
	}
}

// GetPromptByType 根据类型获取提示词
func (s *PromptService) GetPromptByType(promptType string) (*entity.AIPrompt, error) {
	var prompt entity.AIPrompt
	err := s.db.Where("type = ? AND is_active = ?", promptType, true).First(&prompt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("未找到类型为 %s 的提示词", promptType)
		}
		return nil, fmt.Errorf("获取提示词失败: %v", err)
	}
	return &prompt, nil
}

// GetAllPrompts 获取所有提示词
func (s *PromptService) GetAllPrompts() ([]entity.AIPrompt, error) {
	var prompts []entity.AIPrompt
	err := s.db.Order("type ASC").Find(&prompts).Error
	if err != nil {
		return nil, fmt.Errorf("获取提示词列表失败: %v", err)
	}
	return prompts, nil
}

// UpdatePrompt 更新提示词
func (s *PromptService) UpdatePrompt(id uint, userContent string) error {
	result := s.db.Model(&entity.AIPrompt{}).Where("id = ?", id).Update("user_content", userContent)
	if result.Error != nil {
		return fmt.Errorf("更新提示词失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到ID为 %d 的提示词", id)
	}
	utils.Info("提示词 ID %d 已更新", id)
	return nil
}

// UpdatePromptWithDescription 更新提示词和描述
func (s *PromptService) UpdatePromptWithDescription(id uint, userContent, description string) error {
	result := s.db.Model(&entity.AIPrompt{}).Where("id = ?", id).Updates(map[string]interface{}{
		"user_content": userContent,
		"description":  description,
	})
	if result.Error != nil {
		return fmt.Errorf("更新提示词失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到ID为 %d 的提示词", id)
	}
	utils.Info("提示词 ID %d 已更新", id)
	return nil
}

// UpdateSystemPrompt 更新系统提示词
func (s *PromptService) UpdateSystemPrompt(id uint, systemContent string) error {
	result := s.db.Model(&entity.AIPrompt{}).Where("id = ?", id).Update("system_content", systemContent)
	if result.Error != nil {
		return fmt.Errorf("更新系统提示词失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到ID为 %d 的提示词", id)
	}
	utils.Info("系统提示词 ID %d 已更新", id)
	return nil
}

// UpdateUserPrompt 更新用户提示词
func (s *PromptService) UpdateUserPrompt(id uint, userContent string) error {
	result := s.db.Model(&entity.AIPrompt{}).Where("id = ?", id).Update("user_content", userContent)
	if result.Error != nil {
		return fmt.Errorf("更新用户提示词失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到ID为 %d 的提示词", id)
	}
	utils.Info("用户提示词 ID %d 已更新", id)
	return nil
}

// UpdateFullPrompt 完整更新提示词（系统+用户+描述）
func (s *PromptService) UpdateFullPrompt(id uint, systemContent, userContent, description string) error {
	updates := make(map[string]interface{})
	if systemContent != "" {
		updates["system_content"] = systemContent
	}
	if userContent != "" {
		updates["user_content"] = userContent
	}
	if description != "" {
		updates["description"] = description
	}

	if len(updates) == 0 {
		return fmt.Errorf("没有提供要更新的内容")
	}

	result := s.db.Model(&entity.AIPrompt{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新提示词失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到ID为 %d 的提示词", id)
	}
	utils.Info("提示词 ID %d 已完整更新", id)
	return nil
}

// TogglePromptStatus 切换提示词启用状态
func (s *PromptService) TogglePromptStatus(id uint) error {
	var prompt entity.AIPrompt
	if err := s.db.First(&prompt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("未找到ID为 %d 的提示词", id)
		}
		return fmt.Errorf("获取提示词失败: %v", err)
	}

	prompt.IsActive = !prompt.IsActive
	if err := s.db.Save(&prompt).Error; err != nil {
		return fmt.Errorf("更新提示词状态失败: %v", err)
	}

	status := "禁用"
	if prompt.IsActive {
		status = "启用"
	}
	utils.Info("提示词 %s 已%s", prompt.Name, status)
	return nil
}

// RenderPrompt 渲染提示词模板（保持向后兼容，渲染用户内容）
func (s *PromptService) RenderPrompt(prompt *entity.AIPrompt, data interface{}) (string, error) {
	if prompt == nil {
		return "", fmt.Errorf("提示词不能为空")
	}

	return s.RenderUserPrompt(prompt, data)
}

// TestPrompt 测试提示词
func (s *PromptService) TestPrompt(promptType string, testData interface{}) (string, error) {
	prompt, err := s.GetPromptByType(promptType)
	if err != nil {
		return "", err
	}

	rendered, err := s.RenderPrompt(prompt, testData)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

// RenderPromptByType 根据类型直接渲染提示词
func (s *PromptService) RenderPromptByType(promptType string, data interface{}) (string, error) {
	prompt, err := s.GetPromptByType(promptType)
	if err != nil {
		return "", err
	}

	return s.RenderPrompt(prompt, data)
}

// RenderSystemPrompt 渲染系统提示词
func (s *PromptService) RenderSystemPrompt(prompt *entity.AIPrompt, data interface{}) (string, error) {
	if prompt == nil {
		return "", fmt.Errorf("提示词不能为空")
	}

	tmpl, err := template.New("system_prompt").Parse(prompt.SystemContent)
	if err != nil {
		return "", fmt.Errorf("解析系统提示词模板失败: %v", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("渲染系统提示词模板失败: %v", err)
	}

	return result.String(), nil
}

// RenderUserPrompt 渲染用户提示词
func (s *PromptService) RenderUserPrompt(prompt *entity.AIPrompt, data interface{}) (string, error) {
	if prompt == nil {
		return "", fmt.Errorf("提示词不能为空")
	}

	tmpl, err := template.New("user_prompt").Parse(prompt.UserContent)
	if err != nil {
		return "", fmt.Errorf("解析用户提示词模板失败: %v", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("渲染用户提示词模板失败: %v", err)
	}

	return result.String(), nil
}

// RenderSystemPromptByType 根据类型渲染系统提示词
func (s *PromptService) RenderSystemPromptByType(promptType string, data interface{}) (string, error) {
	prompt, err := s.GetPromptByType(promptType)
	if err != nil {
		return "", err
	}

	return s.RenderSystemPrompt(prompt, data)
}

// RenderUserPromptByType 根据类型渲染用户提示词
func (s *PromptService) RenderUserPromptByType(promptType string, data interface{}) (string, error) {
	prompt, err := s.GetPromptByType(promptType)
	if err != nil {
		return "", err
	}

	return s.RenderUserPrompt(prompt, data)
}

// GetPromptVariables 获取提示词变量列表
func (s *PromptService) GetPromptVariables(prompt *entity.AIPrompt) ([]string, error) {
	if prompt.Variables == "" {
		return []string{}, nil
	}

	var variables []string
	if err := json.Unmarshal([]byte(prompt.Variables), &variables); err != nil {
		return nil, fmt.Errorf("解析变量列表失败: %v", err)
	}

	return variables, nil
}

// CreateDefaultPrompts 创建默认提示词
func (s *PromptService) CreateDefaultPrompts() error {
	defaultPrompts := []entity.AIPrompt{
		// 内容生成提示词
		{
			Name:      "内容生成提示词",
			Type:      entity.PromptTypeContentGeneration,
			SystemContent: "你是一个专业的内容优化专家，擅长为各类资源创建吸引人的标题、描述和SEO内容。\n\n你需要根据用户提供的资源信息，生成更优的内容并严格按照JSON格式返回结果。\n\n返回格式要求：\n{\n  \"title\": \"优化后的标题\",\n  \"description\": \"优化后的描述\",\n  \"seo_title\": \"SEO标题\",\n  \"seo_description\": \"SEO描述\",\n  \"seo_keywords\": [\"关键词1\", \"关键词2\"]\n}",
			UserContent:   "请根据以下资源信息生成更优的标题、描述和SEO内容：\n\n原始标题: {{.Title}}\n原始描述: {{.Description}}\n资源类型: {{.Type}}\n\n请提供：\n1. 优化后的标题（更吸引人，更准确）\n2. 详细的资源描述（更全面，更有说服力）\n3. SEO友好的标题（包含关键词，适合搜索引擎）\n4. SEO友好的描述（简洁明了，突出重点）\n5. 相关的SEO关键词（用逗号分隔，便于搜索）",
			Description:   "用于生成资源标题、描述和SEO内容的提示词",
			Variables:     `["Title", "Description", "Type"]`,
			IsActive:      true,
		},
		// 分类推荐提示词
		{
			Name:      "分类推荐提示词",
			Type:      entity.PromptTypeClassification,
			SystemContent: "你是一个专业的分类推荐专家，擅长分析资源内容并为其推荐最合适的分类。\n\n你需要根据资源信息和现有分类列表，推荐最适合的分类并提供详细的推荐理由。\n\n返回格式要求：\n{\n  \"category_id\": 1,\n  \"category_name\": \"分类名称\",\n  \"confidence\": 0.9,\n  \"reason\": \"推荐理由\"\n}",
			UserContent:   "请根据以下资源信息为其推荐最合适的分类：\n\n资源标题: {{.Title}}\n资源描述: {{.Description}}\n资源类型: {{.Type}}\n\n现有分类列表：\n{{range .Categories}}- {{.ID}}: {{.Name}}\n{{end}}\n\n请分析资源内容并推荐最适合的分类ID和分类名称，同时提供置信度（0-1之间的数值）和详细的推荐理由。",
			Description:   "用于资源分类推荐的提示词",
			Variables:     `["Title", "Description", "Type", "Categories"]`,
			IsActive:      true,
		},
		// 工具调用系统提示词
		{
			Name:      "工具调用系统提示词",
			Type:      entity.PromptTypeToolSystem,
			SystemContent: "你叫 老九助手，你是一个充满智慧的辅助专家，可以回答用户的各种问题，并且可以调用各种mcp工具为用户获取更加专业的回答。\n\n【核心规则】\n1. 如果用户询问时间、日期、搜索信息或其他需要实时数据的问题，你必须使用相应的工具\n2. 不要猜测或编造信息，必须使用工具获取准确的数据\n3. 调用工具后，根据工具返回的结果给用户准确的回答\n4. 【最重要】调用工具时，必须提供所有必需的参数，不要省略任何 required 参数\n5. 根据工具的参数定义和用户的问题，智能选择合适的参数值\n\n【工具调用格式 - 仅使用JSON格式】\n格式：<工具名称: {\"参数名\": \"参数值\"}>\n示例：<search: {\"query\": \"人工智能最新进展\"}>\n示例：<current_time: {\"format\": \"YYYY-MM-DD HH:mm:ss\"}>\n\n【严格约束 - 违反将导致工具调用失败】\n⚠️  绝对不要使用空对象 {}\n⚠️  必须提供所有标记为 required 的参数\n⚠️  确保工具名称与可用工具列表中的名称完全一致\n⚠️  所有参数值都必须用双引号包裹\n⚠️  根据用户问题的具体需求，选择最合适的参数值\n\n【关键工具参数要求】\n🔹 current_time: 【必须提供 format 参数】\n   - 正确示例：<current_time: {\"format\": \"YYYY-MM-DD HH:mm:ss\"}>\n   - 时间格式选择：\n     * 用户问\"今天几号\" → 用 \"YYYY-MM-DD\"\n     * 用户问\"现在几点\" → 用 \"HH:mm:ss\" \n     * 用户问\"现在时间\" → 用 \"YYYY-MM-DD HH:mm:ss\"\n   - 可选参数：timezone (如 \"Asia/Shanghai\")\n\n🔹 relative_time: 【必须提供 time 参数】\n   - 格式：YYYY-MM-DD HH:mm:ss\n   - 示例：<relative_time: {\"time\": \"2025-01-01 12:00:00\"}>\n\n🔹 search: 【必须提供 query 参数】\n   - 示例：<search: {\"query\": \"人工智能最新进展\"}>\n   - 可选参数：max_results (默认25)\n\n🔹 fetch_content: 【必须提供 url 参数】\n   - 示例：<fetch_content: {\"url\": \"https://example.com\"}>\n\n【警告】如果工具调用失败，检查是否遗漏了 required 参数！",
			UserContent:   "【可用工具列表及参数要求】\n{{.ToolListWithParams}}\n\n【重要提醒】\n⚠️ 调用工具时必须提供所有必需参数（标记为【必需】的参数）\n⚠️ 使用格式：<工具名称: {\"参数名\": \"参数值\"}>\n⚠️ 绝对不要使用空对象 {}\n\n用户的请求会在这里提供，请根据用户的问题使用相应的工具来获取准确信息并回答。",
			Description:   "AI助手的系统指令和工具调用规则",
			Variables:     `["ToolListWithParams"]`,
			IsActive:      true,
		},
		// 问答模板提示词
		{
			Name:      "问答模板提示词",
			Type:      entity.PromptTypeQATemplate,
			SystemContent: "你是一个专业的问答助手，擅长基于提供的上下文信息给出准确的回答。\n\n你需要严格根据上下文信息回答问题，不要编造或推测信息。如果上下文中没有相关信息，请明确说明。",
			UserContent:   "根据以下上下文回答问题：\n\n上下文：{{.Context}}\n\n问题：{{.Question}}\n\n请基于提供的上下文信息给出准确的回答。",
			Description:   "基于上下文的问答模板",
			Variables:     `["Context", "Question"]`,
			IsActive:      true,
		},
		// 文本分析模板提示词
		{
			Name:      "文本分析模板提示词",
			Type:      entity.PromptTypeAnalysisTemplate,
			SystemContent: "你是一个专业的文本分析专家，擅长对各类文本进行深入分析。\n\n你需要根据用户指定的分析类型，对提供的文本进行全面、准确的分析，并提供有价值的见解。",
			UserContent:   "请对以下文本进行{{.AnalysisType}}分析：\n\n{{.Text}}",
			Description:   "文本分析指令模板",
			Variables:     `["Text", "AnalysisType"]`,
			IsActive:      true,
		},
	}

	for _, prompt := range defaultPrompts {
		// 检查是否已存在相同类型的提示词
		var existing entity.AIPrompt
		err := s.db.Where("type = ?", prompt.Type).First(&existing).Error
		if err == nil {
			utils.Debug("提示词类型 %s 已存在，跳过创建", prompt.Type)
			continue
		}
		if err != gorm.ErrRecordNotFound {
			utils.Error("检查提示词是否存在时出错: %v", err)
			continue
		}

		// 创建新的提示词
		if err := s.db.Create(&prompt).Error; err != nil {
			utils.Error("创建默认提示词 %s 失败: %v", prompt.Name, err)
		} else {
			utils.Info("创建默认提示词: %s", prompt.Name)
		}
	}

	return nil
}