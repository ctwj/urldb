package service

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// GeneratedContentPreview 内容生成预览结构
type GeneratedContentPreview struct {
	SessionID             string     `json:"session_id"`
	ResourceID            uint       `json:"resource_id"`
	OriginalTitle         string     `json:"original_title"`
	OriginalDescription   string     `json:"original_description"`
	GeneratedTitle        string     `json:"generated_title"`
	GeneratedDescription  string     `json:"generated_description"`
	GeneratedSEOTitle     string     `json:"generated_seo_title"`
	GeneratedSEODescription string   `json:"generated_seo_description"`
	GeneratedSEOKeywords  []string   `json:"generated_seo_keywords"`
	AIModelUsed           string     `json:"ai_model_used"`
	GeneratedAt           time.Time  `json:"generated_at"`
}

// ContentGenerator 内容生成器
type ContentGenerator struct {
	client       *OpenAIClient
	promptService *PromptService  // 添加提示词服务
	repoManager  *repo.RepositoryManager
}

// NewContentGenerator 创建内容生成器
func NewContentGenerator(client *OpenAIClient, repoManager *repo.RepositoryManager) *ContentGenerator {
	return &ContentGenerator{
		client: client,
		promptService: NewPromptService(repoManager.GetDB()),
		repoManager: repoManager,
	}
}

// GenerateContentPreview 预生成内容 - 用户发起后 AI 生成内容，但不立即写入数据库
func (cg *ContentGenerator) GenerateContentPreview(resourceID uint) (*GeneratedContentPreview, error) {
	// 获取资源
	resource, err := cg.repoManager.ResourceRepository.FindByID(resourceID)
	if err != nil {
		return nil, err
	}

	// 从数据库获取内容生成提示词
	prompt, err := cg.promptService.GetPromptByType("content_generation")
	if err != nil {
		return nil, fmt.Errorf("获取内容生成提示词失败: %v", err)
	}

	// 构建模板数据
	data := struct {
		Title       string
		Description string
		Type        string
	}{
		Title:       resource.Title,
		Description: resource.Description,
	}

	// 设置资源类型
	if resource.PanID != nil {
		pan, err := cg.repoManager.PanRepository.FindByID(*resource.PanID)
		if err == nil && pan != nil {
			data.Type = pan.Name
		} else {
			data.Type = "未知"
		}
	} else {
		data.Type = "未知"
	}

	// 渲染提示词
	renderedPrompt, err := cg.promptService.RenderPrompt(prompt, data)
	if err != nil {
		return nil, fmt.Errorf("渲染内容生成提示词失败: %v", err)
	}

	// 调用 AI API
	resp, err := cg.client.CreateChatCompletion(renderedPrompt)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI 未返回任何内容")
	}

	// 解析响应
	content := resp.Choices[0].Message.Content
	generatedContent, err := cg.parseGeneratedContent(content)
	if err != nil {
		return nil, fmt.Errorf("解析生成内容失败: %v", err)
	}

	// 创建预览结果
	preview := &GeneratedContentPreview{
		ResourceID:            resourceID,
		OriginalTitle:         resource.Title,
		OriginalDescription:   resource.Description,
		GeneratedTitle:        generatedContent.Title,
		GeneratedDescription:  generatedContent.Description,
		GeneratedSEOTitle:     generatedContent.SEOTitle,
		GeneratedSEODescription: generatedContent.SEODescription,
		GeneratedSEOKeywords:  generatedContent.SEOKeywords,
		AIModelUsed:           cg.client.model,
		GeneratedAt:           time.Now(),
		SessionID:             generatePreviewSessionID(), // 生成会话ID用于确认操作
	}

	return preview, nil
}

// ApplyGeneratedContent 确认并应用生成的内容 - 用户确认后才写入数据库
func (cg *ContentGenerator) ApplyGeneratedContent(preview *GeneratedContentPreview) error {
	// 获取最新资源数据
	resource, err := cg.repoManager.ResourceRepository.FindByID(preview.ResourceID)
	if err != nil {
		return err
	}

	// 检查资源是否在预览期间被修改
	if resource.UpdatedAt.After(preview.GeneratedAt) {
		return fmt.Errorf("资源在预览期间已被修改，请重新生成")
	}

	// 应用生成的内容
	if preview.GeneratedTitle != "" {
		resource.Title = preview.GeneratedTitle
	}
	if preview.GeneratedDescription != "" {
		resource.Description = preview.GeneratedDescription
	}

	// 更新 SEO 相关字段
	resource.SEOTitle = &preview.GeneratedSEOTitle
	resource.SEODescription = &preview.GeneratedSEODescription
	resource.SEOKeywords = &preview.GeneratedSEOKeywords

	// 更新 AI 相关字段
	resource.AIModelUsed = &preview.AIModelUsed
	status := "completed"
	resource.AIGenerationStatus = &status
	timestamp := time.Now()
	resource.AIGenerationTimestamp = &timestamp
	resource.AILastRegeneration = &timestamp

	// 保存更新
	return cg.repoManager.ResourceRepository.Update(resource)
}

// GeneratedContent 解析后的生成内容
type GeneratedContent struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	SEOTitle        string   `json:"seo_title"`
	SEODescription  string   `json:"seo_description"`
	SEOKeywords     []string `json:"seo_keywords"`
}

// parseGeneratedContent 解析生成的内容
func (cg *ContentGenerator) parseGeneratedContent(content string) (*GeneratedContent, error) {
	// 提取JSON部分
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	jsonMatch := re.FindString(content)

	if jsonMatch == "" {
		// 如果没有找到JSON格式，尝试从文本中提取信息
		return cg.extractContentFromText(content)
	}

	// 这里我们简单地解析JSON，实际实现中可能需要使用json包进行解析
	// 为简化处理，我们直接从文本中提取信息
	return cg.extractContentFromText(content)
}

// extractContentFromText 从文本中提取内容
func (cg *ContentGenerator) extractContentFromText(content string) (*GeneratedContent, error) {
	// 简化处理：直接从文本中提取相关信息
	// 实际实现中可能需要更复杂的解析逻辑
	lines := strings.Split(content, "\n")

	genContent := &GeneratedContent{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"title\":") || strings.Contains(line, "优化后的标题") || strings.Contains(line, "标题:") {
			// 提取标题
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				title := strings.TrimSpace(line[colonIndex+1:])
				title = strings.Trim(title, "\"")
				genContent.Title = title
			}
		} else if strings.HasPrefix(line, "\"description\":") || strings.Contains(line, "优化后的描述") || strings.Contains(line, "描述:") {
			// 提取描述
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				desc := strings.TrimSpace(line[colonIndex+1:])
				desc = strings.Trim(desc, "\"")
				genContent.Description = desc
			}
		} else if strings.Contains(line, "seo_title") || strings.Contains(line, "SEO标题") {
			// 提取SEO标题
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				seoTitle := strings.TrimSpace(line[colonIndex+1:])
				seoTitle = strings.Trim(seoTitle, "\"")
				genContent.SEOTitle = seoTitle
			}
		} else if strings.Contains(line, "seo_description") || strings.Contains(line, "SEO描述") {
			// 提取SEO描述
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				seoDesc := strings.TrimSpace(line[colonIndex+1:])
				seoDesc = strings.Trim(seoDesc, "\"")
				genContent.SEODescription = seoDesc
			}
		} else if strings.Contains(line, "seo_keywords") || strings.Contains(line, "关键词") {
			// 提取SEO关键词
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 {
				keywordsStr := strings.TrimSpace(line[colonIndex+1:])
				// 简单解析关键词列表
				keywordsStr = strings.Trim(keywordsStr, "[]\"")
				keywords := strings.Split(keywordsStr, ",")
				for i, keyword := range keywords {
					keywords[i] = strings.TrimSpace(keyword)
				}
				genContent.SEOKeywords = keywords
			}
		}
	}

	return genContent, nil
}

// generatePreviewSessionID 生成预览会话ID
func generatePreviewSessionID() string {
	return fmt.Sprintf("preview_%d_%x", time.Now().Unix(),
		sha256.Sum256([]byte(fmt.Sprintf("%d_%s", time.Now().UnixNano(),
			strconv.Itoa(rand.Int())))))
}

// buildPrompt 构建提示词
func (cg *ContentGenerator) buildPrompt(resource *entity.Resource) string {
	var prompt strings.Builder
	prompt.WriteString("请根据以下资源信息生成更优的标题、描述和SEO内容：\n\n")
	prompt.WriteString(fmt.Sprintf("原始标题: %s\n", resource.Title))
	prompt.WriteString(fmt.Sprintf("原始描述: %s\n", resource.Description))
	if resource.PanID != nil {
		pan, err := cg.repoManager.PanRepository.FindByID(*resource.PanID)
		if err == nil && pan != nil {
			prompt.WriteString(fmt.Sprintf("资源类型: %s\n", pan.Name))
		} else {
			prompt.WriteString("资源类型: 未知\n")
		}
	} else {
		prompt.WriteString("资源类型: 未知\n")
	}
	prompt.WriteString("\n请提供：\n")
	prompt.WriteString("1. 优化后的标题\n")
	prompt.WriteString("2. 详细的资源描述\n")
	prompt.WriteString("3. SEO友好的标题\n")
	prompt.WriteString("4. SEO友好的描述\n")
	prompt.WriteString("5. 相关的SEO关键词（用逗号分隔）\n\n")
	prompt.WriteString("请以JSON格式返回结果，格式如下：\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"title\": \"优化后的标题\",\n")
	prompt.WriteString("  \"description\": \"优化后的描述\",\n")
	prompt.WriteString("  \"seo_title\": \"SEO标题\",\n")
	prompt.WriteString("  \"seo_description\": \"SEO描述\",\n")
	prompt.WriteString("  \"seo_keywords\": [\"关键词1\", \"关键词2\"]\n")
	prompt.WriteString("}")

	return prompt.String()
}