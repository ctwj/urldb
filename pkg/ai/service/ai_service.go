package service

import (
	"fmt"

	"github.com/ctwj/urldb/db/repo"
	"github.com/sashabaranov/go-openai"
)

// AIService 主AI服务，提供通用AI能力供其他模块调用
type AIService struct {
	client        *OpenAIClient
	contentGen    *ContentGenerator
	classifier    *Classifier
	repoManager   *repo.RepositoryManager
}

// NewAIServiceWithConfig 创建AI服务
func NewAIServiceWithConfig(configManager ConfigManager, repoManager *repo.RepositoryManager) (*AIService, error) {
	client, err := NewOpenAIClientWithConfig(configManager)
	if err != nil {
		return nil, fmt.Errorf("创建OpenAI客户端失败: %v", err)
	}

	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)

	return &AIService{
		client:      client,
		contentGen:  contentGen,
		classifier:  classifier,
		repoManager: repoManager,
	}, nil
}

// NewAIService 创建AI服务
func NewAIService(client *OpenAIClient, repoManager *repo.RepositoryManager) (*AIService, error) {
	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)

	return &AIService{
		client:      client,
		contentGen:  contentGen,
		classifier:  classifier,
		repoManager: repoManager,
	}, nil
}

// GenerateText 通用文本生成 - 供其他模块调用
func (as *AIService) GenerateText(prompt string, options ...ChatOption) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	resp, err := as.client.Chat(messages, options...)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	return resp.Choices[0].Message.Content, nil
}

// AskQuestion 通用问答 - 供其他模块调用
func (as *AIService) AskQuestion(question string, context string) (string, error) {
	prompt := fmt.Sprintf("根据以下上下文回答问题：\n\n上下文：%s\n\n问题：%s\n\n请基于提供的上下文信息给出准确的回答。", context, question)
	return as.GenerateText(prompt, WithMaxTokens(500), WithTemperature(0.7))
}

// AnalyzeText 通用文本分析 - 供其他模块调用
func (as *AIService) AnalyzeText(text string, analysisType string) (string, error) {
	prompt := fmt.Sprintf("请对以下文本进行%s分析：\n\n%s", analysisType, text)
	return as.GenerateText(prompt, WithMaxTokens(300), WithTemperature(0.5))
}

// GenerateContentPreview 生成内容预览
func (as *AIService) GenerateContentPreview(resourceID uint) (*GeneratedContentPreview, error) {
	return as.contentGen.GenerateContentPreview(resourceID)
}

// ApplyGeneratedContent 应用生成的内容
func (as *AIService) ApplyGeneratedContent(preview *GeneratedContentPreview) error {
	return as.contentGen.ApplyGeneratedContent(preview)
}

// ClassifyResourcePreview 分类资源预览
func (as *AIService) ClassifyResourcePreview(resourceID uint) (*ClassificationPreview, error) {
	return as.classifier.ClassifyResourcePreview(resourceID)
}

// ApplyClassification 应用分类建议
func (as *AIService) ApplyClassification(preview *ClassificationPreview) error {
	return as.classifier.ApplyClassification(preview)
}

// TestConnection 测试AI连接
func (as *AIService) TestConnection() error {
	_, err := as.GenerateText("你好，请回复'连接正常'。", WithMaxTokens(10))
	return err
}

// ReloadClient 重新加载客户端配置
func (as *AIService) ReloadClient() error {
	return as.client.ReloadConfig()
}

// GetModel 获取当前使用的模型
func (as *AIService) GetModel() string {
	return as.client.model
}