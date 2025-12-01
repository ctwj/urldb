package services

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

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
	utils.Debug("[WECHAT:MESSAGE] 处理文本消息 - AutoReplyEnabled: %v", s.config.AutoReplyEnabled)
	if !s.config.AutoReplyEnabled {
		utils.Info("[WECHAT:MESSAGE] 自动回复未启用")
		return nil, nil
	}

	keyword := strings.TrimSpace(msg.Content)
	utils.Info("[WECHAT:MESSAGE] 搜索关键词: '%s'", keyword)

	// 检查是否是分页命令
	if keyword == "上一页" || keyword == "prev" {
		return s.handlePrevPage(string(msg.FromUserName))
	}

	if keyword == "下一页" || keyword == "next" {
		return s.handleNextPage(string(msg.FromUserName))
	}

	// 检查是否是获取命令（例如：获取 1, 获取2等）
	if strings.HasPrefix(keyword, "获取") || strings.HasPrefix(keyword, "get") {
		return s.handleGetResource(string(msg.FromUserName), keyword)
	}

	// 检查是否为纯数字命令（例如：1, 2等），如果是，则将其作为获取资源命令处理
	if isPureNumber(keyword) {
		// 检查用户是否有搜索会话
		session := s.searchSessionManager.GetSession(string(msg.FromUserName))
		if session != nil {
			// 如果有搜索会话，则将数字作为获取资源命令处理
			return s.handleGetResource(string(msg.FromUserName), keyword)
		}
	}

	if keyword == "" {
		utils.Info("[WECHAT:MESSAGE] 关键词为空，返回提示消息")
		return message.NewText("请输入搜索关键词"), nil
	}

	// 检查搜索关键词是否包含违禁词
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return s.systemConfigRepo.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{} // 如果获取失败，使用空列表
	}

	// 检查关键词是否包含违禁词
	if len(cleanWords) > 0 {
		containsForbidden, matchedWords := utils.CheckContainsForbiddenWords(keyword, cleanWords)
		if containsForbidden {
			utils.Info("[WECHAT:MESSAGE] 搜索关键词包含违禁词: %v", matchedWords)
			return message.NewText("您的搜索关键词包含违禁内容，不予处理"), nil
		}
	}

	// 搜索资源
	utils.Debug("[WECHAT:MESSAGE] 开始搜索资源，限制数量: %d", s.config.SearchLimit)
	resources, err := s.SearchResources(keyword)
	if err != nil {
		utils.Error("[WECHAT:SEARCH] 搜索失败: %v", err)
		return message.NewText("搜索服务暂时不可用，请稍后重试"), nil
	}

	utils.Info("[WECHAT:MESSAGE] 搜索完成，找到 %d 个资源", len(resources))
	if len(resources) == 0 {
		utils.Info("[WECHAT:MESSAGE] 未找到相关资源，返回提示消息")
		return message.NewText(fmt.Sprintf("未找到关键词\"%s\"相关的资源，请尝试其他关键词", keyword)), nil
	}

	// 创建搜索会话并保存第一页结果
	s.searchSessionManager.CreateSession(string(msg.FromUserName), keyword, resources, 4)
	pageResources := s.searchSessionManager.GetCurrentPageResources(string(msg.FromUserName))

	// 格式化第一页搜索结果
	resultText := s.formatSearchResultsWithPagination(keyword, pageResources, string(msg.FromUserName))
	utils.Info("[WECHAT:MESSAGE] 格式化搜索结果，返回文本长度: %d", len(resultText))
	return message.NewText(resultText), nil
}

// handlePrevPage 处理上一页命令
func (s *WechatBotServiceImpl) handlePrevPage(userID string) (interface{}, error) {
	session := s.searchSessionManager.GetSession(userID)
	if session == nil {
		return message.NewText("没有找到搜索记录，请先进行搜索"), nil
	}

	if !s.searchSessionManager.HasPrevPage(userID) {
		return message.NewText("已经是第一页了"), nil
	}

	prevResources := s.searchSessionManager.PrevPage(userID)
	if prevResources == nil {
		return message.NewText("获取上一页失败"), nil
	}

	currentPage, totalPages, _, _ := s.searchSessionManager.GetPageInfo(userID)
	resultText := s.formatPageResources(session.Keyword, prevResources, currentPage, totalPages, userID)
	return message.NewText(resultText), nil
}

// handleNextPage 处理下一页命令
func (s *WechatBotServiceImpl) handleNextPage(userID string) (interface{}, error) {
	session := s.searchSessionManager.GetSession(userID)
	if session == nil {
		return message.NewText("没有找到搜索记录，请先进行搜索"), nil
	}

	if !s.searchSessionManager.HasNextPage(userID) {
		return message.NewText("已经是最后一页了"), nil
	}

	nextResources := s.searchSessionManager.NextPage(userID)
	if nextResources == nil {
		return message.NewText("获取下一页失败"), nil
	}

	currentPage, totalPages, _, _ := s.searchSessionManager.GetPageInfo(userID)
	resultText := s.formatPageResources(session.Keyword, nextResources, currentPage, totalPages, userID)
	return message.NewText(resultText), nil
}

// handleGetResource 处理获取资源命令
func (s *WechatBotServiceImpl) handleGetResource(userID, command string) (interface{}, error) {
	session := s.searchSessionManager.GetSession(userID)
	if session == nil {
		return message.NewText("没有找到搜索记录，请先进行搜索"), nil
	}

	// 检查是否只输入了"获取"或"get"，没有指定编号
	if command == "获取" || command == "get" {
		return message.NewText("📌 请输入要获取的资源编号\n\n💡 提示：回复\"获取 1\"或\"get 1\"获取第一个资源的详细信息"), nil
	}

	// 检查是否为纯数字命令（如 "1", "2" 等），如果是则转换为 "获取X" 格式
	if isPureNumber(command) {
		command = "获取" + command
	}

	// 解析命令，例如："获取 1" 或 "get 2"
	// 支持"获取4"这种没有空格的格式
	var index int
	var err error
	patterns := []string{"获取%d", "获取 %d", "get%d", "get %d"}

	parsed := false
	for _, pattern := range patterns {
		_, err = fmt.Sscanf(command, pattern, &index)
		if err == nil {
			parsed = true
			break
		}
	}

	if !parsed {
		return message.NewText("❌ 命令格式错误\n\n📌 正确格式：\n   • 获取 1\n   • get 1\n   • 获取1\n   • get1\n   • 直接输入数字 1"), nil
	}

	if index < 1 || index > len(session.Resources) {
		return message.NewText(fmt.Sprintf("❌ 资源编号超出范围\n\n📌 请输入 1-%d 之间的数字\n💡 提示：回复\"获取 %d\"获取第%d个资源", len(session.Resources), index, index)), nil
	}

	// 获取指定资源
	resource := session.Resources[index-1]

	// 格式化资源详细信息（美化输出）
	var result strings.Builder
	// result.WriteString(fmt.Sprintf("📌 资源详情\n\n"))

	// 标题
	result.WriteString(fmt.Sprintf("📌 标题: %s\n", resource.Title))

	// 描述
	if resource.Description != "" {
		result.WriteString(fmt.Sprintf("\n📝 描述:\n   %s\n", resource.Description))
	}

	// 文件大小
	if resource.FileSize != "" {
		result.WriteString(fmt.Sprintf("\n📊 大小: %s\n", resource.FileSize))
	}

	// 作者
	if resource.Author != "" {
		result.WriteString(fmt.Sprintf("\n👤 作者: %s\n", resource.Author))
	}

	// 分类
	if resource.Category.Name != "" {
		result.WriteString(fmt.Sprintf("\n📂 分类: %s\n", resource.Category.Name))
	}

	// 标签
	if len(resource.Tags) > 0 {
		result.WriteString("\n🏷️ 标签: ")
		var tags []string
		for _, tag := range resource.Tags {
			tags = append(tags, tag.Name)
		}
		result.WriteString(fmt.Sprintf("%s\n", strings.Join(tags, " ")))
	}

	// 链接（美化）
	if resource.SaveURL != "" {
		result.WriteString(fmt.Sprintf("\n📥 转存链接:\n   %s", resource.SaveURL))
	} else if resource.URL != "" {
		result.WriteString(fmt.Sprintf("\n🔗 资源链接:\n   %s", resource.URL))
	}

	// 添加操作提示
	result.WriteString(fmt.Sprintf("\n\n💡 提示：回复\"获取 %d\"可再次查看此资源", index))

	return message.NewText(result.String()), nil
}

// formatSearchResultsWithPagination 格式化带分页的搜索结果
func (s *WechatBotServiceImpl) formatSearchResultsWithPagination(keyword string, resources []entity.Resource, userID string) string {
	currentPage, totalPages, _, _ := s.searchSessionManager.GetPageInfo(userID)
	return s.formatPageResources(keyword, resources, currentPage, totalPages, userID)
}

// formatPageResources 格式化页面资源
// 根据用户需求，搜索结果中不显示资源链接，只显示标题和描述
func (s *WechatBotServiceImpl) formatPageResources(keyword string, resources []entity.Resource, currentPage, totalPages int, userID string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("🔍 搜索\"%s\"的结果（第%d/%d页）：\n\n", keyword, currentPage, totalPages))

	for i, resource := range resources {
		// 构建当前资源的文本表示
		var resourceText strings.Builder

		// 计算全局索引（当前页的第i个资源在整个结果中的位置）
		globalIndex := (currentPage-1)*4 + i + 1
		resourceText.WriteString(fmt.Sprintf("%d. 📌 %s\n", globalIndex, resource.Title))

		if resource.Description != "" {
			// 限制描述长度以避免消息过长（正确处理中文字符）
			desc := resource.Description
			// 将字符串转换为 rune 切片以正确处理中文字符
			runes := []rune(desc)
			if len(runes) > 50 {
				desc = string(runes[:50]) + "..."
			}
			resourceText.WriteString(fmt.Sprintf("   📝 %s\n", desc))
		}

		// 添加标签显示（格式：🏷️标签,空格,再接别的标签）
		if len(resource.Tags) > 0 {
			var tags []string
			for _, tag := range resource.Tags {
				tags = append(tags, "🏷️"+tag.Name)
			}
			// 限制标签数量以避免消息过长
			if len(tags) > 5 {
				tags = tags[:5]
			}
			resourceText.WriteString(fmt.Sprintf("   %s\n", strings.Join(tags, " ")))
		}

		resourceText.WriteString(fmt.Sprintf("   👉 回复\"获取 %d\"查看详细信息\n", globalIndex))
		resourceText.WriteString("\n")

		// 预计算添加当前资源后的消息长度
		tempMessage := result.String() + resourceText.String()

		// 添加分页提示和预留空间
		if currentPage > 1 || currentPage < totalPages {
			tempMessage += "💡 提示：回复\""
			if currentPage > 1 && currentPage < totalPages {
				tempMessage += "上一页\"或\"下一页"
			} else if currentPage > 1 {
				tempMessage += "上一页"
			} else {
				tempMessage += "下一页"
			}
			tempMessage += "\"翻页\n"
		}

		// 检查添加当前资源后是否会超过微信限制
		tempRunes := []rune(tempMessage)
		if len(tempRunes) > 550 {
			result.WriteString("💡 内容较多，请翻页查看更多\n")
			break
		}

		// 如果不会超过限制，则添加当前资源到结果中
		result.WriteString(resourceText.String())
	}

	// 添加分页提示
	var pageTips []string
	if currentPage > 1 {
		pageTips = append(pageTips, "上一页")
	}
	if currentPage < totalPages {
		pageTips = append(pageTips, "下一页")
	}

	if len(pageTips) > 0 {
		result.WriteString(fmt.Sprintf("💡 提示：回复\"%s\"翻页\n", strings.Join(pageTips, "\"或\"")))
	}

	// 确保消息不超过微信限制（正确处理中文字符）
	message := result.String()
	// 将字符串转换为 rune 切片以正确处理中文字符
	runes := []rune(message)
	if len(runes) > 600 {
		// 如果还是超过限制，截断消息（微信建议不超过600个字符）
		message = string(runes[:597]) + "..."
	}

	return message
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
	// 使用统一搜索函数（包含Meilisearch优先搜索和违禁词处理）
	return UnifiedSearchResources(keyword, s.config.SearchLimit, s.systemConfigRepo, s.resourceRepo)
}

// formatSearchResults 格式化搜索结果
func (s *WechatBotServiceImpl) formatSearchResults(keyword string, resources []entity.Resource) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("🔍 搜索\"%s\"的结果（共%d条）：\n\n", keyword, len(resources)))

	for i, resource := range resources {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, resource.Title))
		if resource.Cover != "" {
			result.WriteString(fmt.Sprintf("   ![封面](%s)\n", resource.Cover))
		}
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

// isPureNumber 检查字符串是否为纯数字
func isPureNumber(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
