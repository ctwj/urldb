package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// loadConfig 加载微信配置
func (s *WechatBotServiceImpl) loadConfig() error {
	configs, err := s.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	utils.Info("[WECHAT] 从数据库加载到 %d 个配置项", len(configs))

	// 初始化默认值
	s.config.Enabled = false
	s.config.AppID = ""
	s.config.AppSecret = ""
	s.config.Token = ""
	s.config.EncodingAesKey = ""
	s.config.WelcomeMessage = "欢迎关注老九网盘资源库！发送关键词即可搜索资源。"
	s.config.AutoReplyEnabled = true
	s.config.SearchLimit = 5

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyWechatBotEnabled:
			s.config.Enabled = config.Value == "true"
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = %s (Enabled: %v)", config.Key, config.Value, s.config.Enabled)
		case entity.ConfigKeyWechatAppId:
			s.config.AppID = config.Value
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = [HIDDEN]", config.Key)
		case entity.ConfigKeyWechatAppSecret:
			s.config.AppSecret = config.Value
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = [HIDDEN]", config.Key)
		case entity.ConfigKeyWechatToken:
			s.config.Token = config.Value
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = [HIDDEN]", config.Key)
		case entity.ConfigKeyWechatEncodingAesKey:
			s.config.EncodingAesKey = config.Value
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = [HIDDEN]", config.Key)
		case entity.ConfigKeyWechatWelcomeMessage:
			if config.Value != "" {
				s.config.WelcomeMessage = config.Value
			}
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = %s", config.Key, config.Value)
		case entity.ConfigKeyWechatAutoReplyEnabled:
			s.config.AutoReplyEnabled = config.Value == "true"
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = %s (AutoReplyEnabled: %v)", config.Key, config.Value, s.config.AutoReplyEnabled)
		case entity.ConfigKeyWechatSearchLimit:
			if config.Value != "" {
				limit, err := strconv.Atoi(config.Value)
				if err == nil && limit > 0 {
					s.config.SearchLimit = limit
				}
			}
			utils.Info("[WECHAT:CONFIG] 加载配置 %s = %s (SearchLimit: %d)", config.Key, config.Value, s.config.SearchLimit)
		}
	}

	utils.Info("[WECHAT:SERVICE] 微信公众号机器人配置加载完成: Enabled=%v, AutoReplyEnabled=%v",
		s.config.Enabled, s.config.AutoReplyEnabled)
	return nil
}

// Start 启动微信公众号机器人服务
func (s *WechatBotServiceImpl) Start() error {
	if s.isRunning {
		utils.Info("[WECHAT:SERVICE] 微信公众号机器人服务已经在运行中")
		return nil
	}

	// 加载配置
	if err := s.loadConfig(); err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if !s.config.Enabled || s.config.AppID == "" || s.config.AppSecret == "" {
		utils.Info("[WECHAT:SERVICE] 微信公众号机器人未启用或配置不完整")
		return nil
	}

	// 创建微信客户端
	cfg := &config.Config{
		AppID:          s.config.AppID,
		AppSecret:      s.config.AppSecret,
		Token:          s.config.Token,
		EncodingAESKey: s.config.EncodingAesKey,
		Cache:          cache.NewMemory(),
	}
	s.wechatClient = officialaccount.NewOfficialAccount(cfg)

	s.isRunning = true
	utils.Info("[WECHAT:SERVICE] 微信公众号机器人服务已启动")
	return nil
}

// Stop 停止微信公众号机器人服务
func (s *WechatBotServiceImpl) Stop() error {
	if !s.isRunning {
		return nil
	}

	s.isRunning = false
	utils.Info("[WECHAT:SERVICE] 微信公众号机器人服务已停止")
	return nil
}

// IsRunning 检查微信公众号机器人服务是否正在运行
func (s *WechatBotServiceImpl) IsRunning() bool {
	return s.isRunning
}

// ReloadConfig 重新加载微信公众号机器人配置
func (s *WechatBotServiceImpl) ReloadConfig() error {
	utils.Info("[WECHAT:SERVICE] 开始重新加载配置...")

	// 重新加载配置
	if err := s.loadConfig(); err != nil {
		utils.Error("[WECHAT:SERVICE] 重新加载配置失败: %v", err)
		return fmt.Errorf("重新加载配置失败: %v", err)
	}

	utils.Info("[WECHAT:SERVICE] 配置重新加载完成: Enabled=%v, AutoReplyEnabled=%v",
		s.config.Enabled, s.config.AutoReplyEnabled)
	return nil
}

// GetRuntimeStatus 获取微信公众号机器人运行时状态
func (s *WechatBotServiceImpl) GetRuntimeStatus() map[string]interface{} {
	status := map[string]interface{}{
		"is_running":    s.IsRunning(),
		"config_loaded": s.config != nil,
		"app_id":        s.config.AppID,
	}

	return status
}

// GetConfig 获取当前配置
func (s *WechatBotServiceImpl) GetConfig() *WechatBotConfig {
	return s.config
}

// HandleMessage 处理微信消息
func (s *WechatBotServiceImpl) HandleMessage(msg *message.MixMessage) (interface{}, error) {
	utils.Info("[WECHAT:MESSAGE] 收到消息: FromUserName=%s, MsgType=%s, Event=%s, Content=%s",
		msg.FromUserName, msg.MsgType, msg.Event, msg.Content)

	switch msg.MsgType {
	case message.MsgTypeText:
		return s.handleTextMessage(msg)
	case message.MsgTypeEvent:
		return s.handleEventMessage(msg)
	default:
		return nil, nil // 不处理其他类型消息
	}
}

// handleTextMessage 处理文本消息
func (s *WechatBotServiceImpl) handleTextMessage(msg *message.MixMessage) (interface{}, error) {
	if !s.config.AutoReplyEnabled {
		return nil, nil
	}

	keyword := strings.TrimSpace(msg.Content)
	if keyword == "" {
		return message.NewText("请输入搜索关键词"), nil
	}

	// 搜索资源
	resources, err := s.SearchResources(keyword)
	if err != nil {
		utils.Error("[WECHAT:SEARCH] 搜索失败: %v", err)
		return message.NewText("搜索服务暂时不可用，请稍后重试"), nil
	}

	if len(resources) == 0 {
		return message.NewText(fmt.Sprintf("未找到关键词\"%s\"相关的资源，请尝试其他关键词", keyword)), nil
	}

	// 格式化搜索结果
	resultText := s.formatSearchResults(keyword, resources)
	return message.NewText(resultText), nil
}

// handleEventMessage 处理事件消息
func (s *WechatBotServiceImpl) handleEventMessage(msg *message.MixMessage) (interface{}, error) {
	if msg.Event == message.EventSubscribe {
		// 新用户关注
		return message.NewText(s.config.WelcomeMessage), nil
	}
	return nil, nil
}

// SearchResources 搜索资源
func (s *WechatBotServiceImpl) SearchResources(keyword string) ([]entity.Resource, error) {
	// 使用现有的资源搜索功能
	resources, total, err := s.resourceRepo.Search(keyword, nil, 1, s.config.SearchLimit)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return []entity.Resource{}, nil
	}

	return resources, nil
}

// formatSearchResults 格式化搜索结果
func (s *WechatBotServiceImpl) formatSearchResults(keyword string, resources []entity.Resource) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("🔍 搜索\"%s\"的结果（共%d条）：\n\n", keyword, len(resources)))

	for i, resource := range resources {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, resource.Title))
		if resource.Description != "" {
			desc := resource.Description
			if len(desc) > 50 {
				desc = desc[:50] + "..."
			}
			result.WriteString(fmt.Sprintf("   %s\n", desc))
		}
		if resource.SaveURL != "" {
			result.WriteString(fmt.Sprintf("   转存链接：%s\n", resource.SaveURL))
		} else if resource.URL != "" {
			result.WriteString(fmt.Sprintf("   资源链接：%s\n", resource.URL))
		}
		result.WriteString("\n")
	}

	result.WriteString("💡 提示：回复资源编号可获取详细信息")
	return result.String()
}

// SendWelcomeMessage 发送欢迎消息（预留接口，实际通过事件处理）
func (s *WechatBotServiceImpl) SendWelcomeMessage(openID string) error {
	// 实际上欢迎消息是通过关注事件自动发送的
	// 这里提供一个手动发送的接口
	if !s.isRunning || s.wechatClient == nil {
		return fmt.Errorf("微信客户端未初始化")
	}

	// 注意：Customer API 需要额外的权限，这里仅作示例
	// 实际应用中可能需要使用模板消息或其他方式
	return nil
}