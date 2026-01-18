package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// ClassificationPreview 分类预览结构
type ClassificationPreview struct {
	SessionID              string    `json:"session_id"`
	ResourceID             uint      `json:"resource_id"`
	OriginalCategoryID     *uint     `json:"original_category_id"`
	OriginalCategoryName   string    `json:"original_category_name"`
	SuggestedCategoryID    uint      `json:"suggested_category_id"`
	SuggestedCategoryName  string    `json:"suggested_category_name"`
	Confidence             float64   `json:"confidence"`
	AIModelUsed            string    `json:"ai_model_used"`
	GeneratedAt            time.Time `json:"generated_at"`
}

// CategorySuggestion 分类建议结构
type CategorySuggestion struct {
	CategoryID uint     `json:"category_id"`
	CategoryName string `json:"category_name"`
	Confidence float64  `json:"confidence"`
	Reason     string   `json:"reason"`
}

// Classifier 分类器
type Classifier struct {
	client       *OpenAIClient
	promptService *PromptService  // 添加提示词服务
	repoManager  *repo.RepositoryManager
}

// NewClassifier 创建分类器
func NewClassifier(client *OpenAIClient, repoManager *repo.RepositoryManager) *Classifier {
	return &Classifier{
		client: client,
		promptService: NewPromptService(repoManager.GetDB()),
		repoManager: repoManager,
	}
}

// ClassifyResourcePreview 预分类资源 - AI 生成分类建议，但不立即应用到数据库
func (c *Classifier) ClassifyResourcePreview(resourceID uint) (*ClassificationPreview, error) {
	// 获取资源
	resource, err := c.repoManager.ResourceRepository.FindByID(resourceID)
	if err != nil {
		return nil, err
	}

	// 获取所有分类
	categories, err := c.repoManager.CategoryRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 从数据库获取分类提示词
	prompt, err := c.promptService.GetPromptByType("classification")
	if err != nil {
		return nil, fmt.Errorf("获取分类提示词失败: %v", err)
	}

	// 构建模板数据
	data := struct {
		Title       string
		Description string
		Type        string
		Categories  []entity.Category
	}{
		Title:       resource.Title,
		Description: resource.Description,
		Categories:  categories,
	}

	// 设置资源类型
	if resource.PanID != nil {
		pan, err := c.repoManager.PanRepository.FindByID(*resource.PanID)
		if err == nil && pan != nil {
			data.Type = pan.Name
		} else {
			data.Type = "未知"
		}
	} else {
		data.Type = "未知"
	}

	// 渲染提示词
	renderedPrompt, err := c.promptService.RenderPrompt(prompt, data)
	if err != nil {
		return nil, fmt.Errorf("渲染分类提示词失败: %v", err)
	}

	// 调用 AI API
	resp, err := c.client.CreateChatCompletion(renderedPrompt)
	if err != nil {
		return nil, fmt.Errorf("AI 分类请求失败: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI 未返回任何内容")
	}

	// 解析响应
	content := resp.Choices[0].Message.Content
	suggestion, err := c.parseCategorySuggestion(content)
	if err != nil {
		return nil, fmt.Errorf("解析分类结果失败: %v", err)
	}

	// 创建预览结果
	preview := &ClassificationPreview{
		ResourceID:             resourceID,
		OriginalCategoryID:     resource.CategoryID,
		OriginalCategoryName:   getCategoryName(resource.CategoryID, c.repoManager),
		SuggestedCategoryID:    suggestion.CategoryID,
		SuggestedCategoryName:  suggestion.CategoryName,
		Confidence:             suggestion.Confidence,
		AIModelUsed:            c.client.model,
		GeneratedAt:            time.Now(),
		SessionID:              generatePreviewSessionID(), // 生成会话ID用于确认操作
	}

	return preview, nil
}

// ApplyClassification 确认并应用分类建议 - 用户确认后才写入数据库
func (c *Classifier) ApplyClassification(preview *ClassificationPreview) error {
	// 获取最新资源数据
	resource, err := c.repoManager.ResourceRepository.FindByID(preview.ResourceID)
	if err != nil {
		return err
	}

	// 检查资源是否在预览期间被修改
	if resource.UpdatedAt.After(preview.GeneratedAt) {
		return fmt.Errorf("资源在预览期间已被修改，请重新分类")
	}

	// 验证目标分类是否存在
	_, err = c.repoManager.CategoryRepository.FindByID(preview.SuggestedCategoryID)
	if err != nil {
		return fmt.Errorf("目标分类不存在: %v", err)
	}

	// 应用分类建议
	resource.CategoryID = &preview.SuggestedCategoryID

	// 更新 AI 相关字段
	resource.AIModelUsed = &preview.AIModelUsed
	status := "completed"
	resource.AIGenerationStatus = &status
	timestamp := time.Now()
	resource.AIGenerationTimestamp = &timestamp
	resource.AILastRegeneration = &timestamp

	// 保存更新
	return c.repoManager.ResourceRepository.Update(resource)
}

// parseCategorySuggestion 解析分类建议
func (c *Classifier) parseCategorySuggestion(content string) (*CategorySuggestion, error) {
	// 提取JSON部分
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	jsonMatch := re.FindString(content)

	if jsonMatch == "" {
		// 如果没有找到JSON格式，尝试从文本中提取信息
		return c.extractCategoryFromText(content)
	}

	// 这里我们简单地解析JSON，实际实现中可能需要使用json包进行解析
	// 为简化处理，我们直接从文本中提取信息
	return c.extractCategoryFromText(content)
}

// extractCategoryFromText 从文本中提取分类信息
func (c *Classifier) extractCategoryFromText(content string) (*CategorySuggestion, error) {
	// 简化处理：直接从文本中提取相关信息
	// 实际实现中可能需要更复杂的解析逻辑
	lines := strings.Split(content, "\n")

	suggestion := &CategorySuggestion{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "category_id") || strings.Contains(line, "分类ID") {
			// 提取分类ID
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				idStr := strings.TrimSpace(line[colonIndex+1:])
				idStr = strings.Trim(idStr, "\",")
				// 简单解析数字
				if id, err := parseNumber(idStr); err == nil {
					suggestion.CategoryID = uint(id)
				}
			}
		} else if strings.Contains(line, "category_name") || strings.Contains(line, "分类名称") {
			// 提取分类名称
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				name := strings.TrimSpace(line[colonIndex+1:])
				name = strings.Trim(name, "\",")
				suggestion.CategoryName = name
			}
		} else if strings.Contains(line, "confidence") || strings.Contains(line, "置信度") {
			// 提取置信度
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				confStr := strings.TrimSpace(line[colonIndex+1:])
				confStr = strings.Trim(confStr, "\",")
				// 简单解析数字
				if conf, err := parseFloat(confStr); err == nil {
					suggestion.Confidence = conf
				}
			}
		} else if strings.Contains(line, "reason") || strings.Contains(line, "推荐理由") {
			// 提取推荐理由
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				reason := strings.TrimSpace(line[colonIndex+1:])
				reason = strings.Trim(reason, "\",")
				suggestion.Reason = reason
			}
		}
	}

	// 如果没有提取到分类ID，尝试从文本中查找数字
	if suggestion.CategoryID == 0 {
		id := extractFirstNumber(content)
		if id > 0 {
			suggestion.CategoryID = uint(id)
		}
	}

	// 如果没有提取到置信度，设置默认值
	if suggestion.Confidence <= 0 || suggestion.Confidence > 1 {
		suggestion.Confidence = 0.5 // 默认置信度
	}

	return suggestion, nil
}

// getCategoryName 获取分类名称的辅助函数
func getCategoryName(categoryID *uint, repoManager *repo.RepositoryManager) string {
	if categoryID == nil {
		return "未分类"
	}
	category, err := repoManager.CategoryRepository.FindByID(*categoryID)
	if err != nil {
		return "未知分类"
	}
	return category.Name
}

// parseNumber 简单解析数字
func parseNumber(s string) (int, error) {
	// 提取第一个数字
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("未找到数字")
	}

	// 尝试转换为整数
	id, err := strconv.Atoi(match)
	if err != nil {
		return 0, fmt.Errorf("无法将字符串转换为整数: %v", err)
	}
	return id, nil
}

// parseFloat 简单解析浮点数
func parseFloat(s string) (float64, error) {
	// 提取第一个浮点数
	re := regexp.MustCompile(`\d+\.?\d*`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("未找到数字")
	}

	// 尝试转换为浮点数
	conf, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, fmt.Errorf("无法将字符串转换为浮点数: %v", err)
	}

	// 确保置信度在0-1之间
	if conf < 0 {
		conf = 0
	} else if conf > 1 {
		conf = 1
	}

	return conf, nil
}

// extractFirstNumber 提取文本中的第一个数字
func extractFirstNumber(text string) uint {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(text)
	if match == "" {
		return 0
	}

	// 尝试转换为整数
	id, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}

	// 确保返回非负数
	if id < 0 {
		return 0
	}

	return uint(id)
}