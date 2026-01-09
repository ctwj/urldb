package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient 包装 OpenAI 客户端
type OpenAIClient struct {
	apiKey       string
	baseURL      string
	model        string
	organization string
	proxy        string
	timeout      time.Duration
	retryCount   int
	client       *openai.Client
	config       ConfigManager // 添加配置管理器
}

// NewOpenAIClientWithConfig 创建新的 OpenAI 客户端
func NewOpenAIClientWithConfig(configManager ConfigManager) (*OpenAIClient, error) {
	// 从数据库配置系统读取参数
	apiKey, err := configManager.GetConfigString("ai_api_key")
	if err != nil || apiKey == "" {
		// 尝试从环境变量读取
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("AI API Key 未配置")
		}
	}

	baseURL, _ := configManager.GetConfigString("ai_api_url")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model, _ := configManager.GetConfigString("ai_model")
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	organization, _ := configManager.GetConfigString("ai_organization")
	proxy, _ := configManager.GetConfigString("ai_proxy")

	timeoutSec, _ := configManager.GetConfigInt("ai_timeout")
	if timeoutSec == 0 {
		timeoutSec = 30
	}

	retryCount, _ := configManager.GetConfigInt("ai_retry_count")
	if retryCount == 0 {
		retryCount = 3
	}

	// 构建 OpenAI 客户端配置
	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = baseURL
	if organization != "" {
		clientConfig.OrgID = organization
	}

	// 设置代理
	if proxy != "" {
		clientConfig.HTTPClient = &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(proxy)
				},
			},
		}
	} else {
		clientConfig.HTTPClient = &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		}
	}

	return &OpenAIClient{
		apiKey:       apiKey,
		baseURL:      baseURL,
		model:        model,
		organization: organization,
		proxy:        proxy,
		timeout:      time.Duration(timeoutSec) * time.Second,
		retryCount:   retryCount,
		client:       openai.NewClientWithConfig(clientConfig),
		config:       configManager, // 设置配置管理器
	}, nil
}

// NewOpenAIClient 创建新的 OpenAI 客户端 (兼容旧方法，使用默认配置)
func NewOpenAIClient() (*OpenAIClient, error) {
	// 使用空的配置管理器，实际使用时应使用 NewOpenAIClientWithConfig
	apiKey := os.Getenv("OPENAI_API_KEY")
	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = "https://api.openai.com/v1"

	return &OpenAIClient{
		apiKey:       apiKey,
		baseURL:      "https://api.openai.com/v1",
		model:        "gpt-3.5-turbo",
		timeout:      30 * time.Second,
		retryCount:   3,
		client:       openai.NewClientWithConfig(clientConfig),
	}, nil
}

// 从数据库重新加载配置
func (c *OpenAIClient) ReloadConfig() error {
	if c.config == nil {
		return fmt.Errorf("配置管理器未初始化")
	}

	newClient, err := NewOpenAIClientWithConfig(c.config)
	if err != nil {
		return err
	}

	c.client = newClient.client
	c.apiKey = newClient.apiKey
	c.baseURL = newClient.baseURL
	c.model = newClient.model
	c.organization = newClient.organization
	c.proxy = newClient.proxy
	c.timeout = newClient.timeout
	c.retryCount = newClient.retryCount

	return nil
}

// Chat 调用对话接口
func (c *OpenAIClient) Chat(messages []openai.ChatCompletionMessage, options ...ChatOption) (*openai.ChatCompletionResponse, error) {
	request := openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	// 应用选项
	for _, option := range options {
		option(&request)
	}

	resp, err := c.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChatOption 对话选项类型
type ChatOption func(*openai.ChatCompletionRequest)

// WithMaxTokens 设置最大令牌数
func WithMaxTokens(tokens int) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		req.MaxTokens = tokens
	}
}

// WithTemperature 设置温度
func WithTemperature(temp float32) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		req.Temperature = temp
	}
}

// WithSystemPrompt 设置系统提示
func WithSystemPrompt(prompt string) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		// 确保消息数组以系统消息开始
		if len(req.Messages) == 0 {
			req.Messages = append([]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			}, req.Messages...)
		} else {
			// 如果已有消息，替换或添加系统消息
			req.Messages[0] = openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			}
		}
	}
}

// CreateChatCompletion 创建对话完成
func (c *OpenAIClient) CreateChatCompletion(prompt string) (*openai.ChatCompletionResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	return c.Chat(messages)
}