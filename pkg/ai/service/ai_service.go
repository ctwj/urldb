package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ctwj/urldb/db/repo"
	"github.com/sashabaranov/go-openai"
)

// AIConfig AI配置结构
type AIConfig struct {
	APIKey      *string
	APIURL      *string
	Model       *string
	MaxTokens   *int
	Temperature *float32
	Timeout     *int
	RetryCount  *int
}

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
	_, err := as.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return err
}

// TestConnectionWithResponse 测试AI连接并返回响应
func (as *AIService) TestConnectionWithResponse() (string, error) {
	response, err := as.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return response, err
}

// TestConnectionWithConfig 使用临时配置测试AI连接
func (as *AIService) TestConnectionWithConfig(config *AIConfig) error {
	// 创建临时客户端
	tempClient, err := as.createTempClient(config)
	if err != nil {
		return fmt.Errorf("创建临时客户端失败: %v", err)
	}

	// 创建临时 AIService
	tempAIService := &AIService{
		client:      tempClient,
		contentGen:  NewContentGenerator(tempClient, as.repoManager),
		classifier:  NewClassifier(tempClient, as.repoManager),
		repoManager: as.repoManager,
	}

	// 使用临时 AIService 询问模型信息
	_, err = tempAIService.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return err
}

// TestConnectionWithConfigAndResponse 使用临时配置测试AI连接并返回响应
func (as *AIService) TestConnectionWithConfigAndResponse(config *AIConfig) (string, error) {
	// 创建临时客户端
	tempClient, err := as.createTempClient(config)
	if err != nil {
		return "", fmt.Errorf("创建临时客户端失败: %v", err)
	}

	// 创建临时 AIService
	tempAIService := &AIService{
		client:      tempClient,
		contentGen:  NewContentGenerator(tempClient, as.repoManager),
		classifier:  NewClassifier(tempClient, as.repoManager),
		repoManager: as.repoManager,
	}

	// 使用临时 AIService 询问模型信息
	response, err := tempAIService.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return response, err
}

// createTempClient 创建临时客户端
func (as *AIService) createTempClient(config *AIConfig) (*OpenAIClient, error) {
	if config.APIKey == nil || *config.APIKey == "" {
		return nil, fmt.Errorf("API Key 不能为空")
	}

	// 设置默认值
	baseURL := "https://api.openai.com/v1"
	if config.APIURL != nil && *config.APIURL != "" {
		baseURL = *config.APIURL
	}

	model := "gpt-3.5-turbo"
	if config.Model != nil && *config.Model != "" {
		model = *config.Model
	}

	timeout := 30 * time.Second
	if config.Timeout != nil {
		timeout = time.Duration(*config.Timeout) * time.Second
	}

	retryCount := 3
	if config.RetryCount != nil {
		retryCount = *config.RetryCount
	}

	// 创建 OpenAI 客户端配置
	clientConfig := openai.DefaultConfig(*config.APIKey)
	clientConfig.BaseURL = baseURL

	// 设置超时
	clientConfig.HTTPClient = &http.Client{
		Timeout: timeout,
	}

	// 创建 OpenAI 客户端
	openaiClient := openai.NewClientWithConfig(clientConfig)

	// 创建临时 OpenAI 客户端包装器
	tempOpenAIClient := &OpenAIClient{
		apiKey:       *config.APIKey,
		baseURL:      baseURL,
		model:        model,
		organization: "",
		proxy:        "",
		timeout:      timeout,
		retryCount:   retryCount,
		client:       openaiClient,
		config:       nil,
	}

	return tempOpenAIClient, nil
}

// ReloadClient 重新加载客户端配置
func (as *AIService) ReloadClient() error {
	return as.client.ReloadConfig()
}

// GetModel 获取当前使用的模型
func (as *AIService) GetModel() string {
	return as.client.model
}