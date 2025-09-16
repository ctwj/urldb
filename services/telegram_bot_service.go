package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type TelegramBotService interface {
	Start() error
	Stop() error
	ValidateApiKey(apiKey string) (bool, map[string]interface{}, error)
	GetBotUsername() string
	SendMessage(chatID int64, text string) error
	DeleteMessage(chatID int64, messageID int) error
	RegisterChannel(chatID int64, chatName, chatType string) error
	IsChannelRegistered(chatID int64) bool
	HandleWebhookUpdate(c interface{})
}

type TelegramBotServiceImpl struct {
	bot              *tgbotapi.BotAPI
	isRunning        bool
	systemConfigRepo repo.SystemConfigRepository
	channelRepo      repo.TelegramChannelRepository
	resourceRepo     repo.ResourceRepository // 添加资源仓库用于搜索
	cronScheduler    *cron.Cron
	config           *TelegramBotConfig
}

type TelegramBotConfig struct {
	Enabled            bool
	ApiKey             string
	AutoReplyEnabled   bool
	AutoReplyTemplate  string
	AutoDeleteEnabled  bool
	AutoDeleteInterval int // 分钟
}

func NewTelegramBotService(
	systemConfigRepo repo.SystemConfigRepository,
	channelRepo repo.TelegramChannelRepository,
	resourceRepo repo.ResourceRepository,
) TelegramBotService {
	return &TelegramBotServiceImpl{
		isRunning:        false,
		systemConfigRepo: systemConfigRepo,
		channelRepo:      channelRepo,
		resourceRepo:     resourceRepo,
		cronScheduler:    cron.New(),
		config:           &TelegramBotConfig{},
	}
}

// loadConfig 加载配置
func (s *TelegramBotServiceImpl) loadConfig() error {
	configs, err := s.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	utils.Info("从数据库加载到 %d 个配置项", len(configs))

	// 初始化默认值
	s.config.Enabled = false
	s.config.ApiKey = ""
	s.config.AutoReplyEnabled = false // 默认禁用自动回复
	s.config.AutoReplyTemplate = "您好！我可以帮您搜索网盘资源，请输入您要搜索的内容。"
	s.config.AutoDeleteEnabled = false
	s.config.AutoDeleteInterval = 60

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyTelegramBotEnabled:
			s.config.Enabled = config.Value == "true"
			utils.Info("加载配置 %s = %s (Enabled: %v)", config.Key, config.Value, s.config.Enabled)
		case entity.ConfigKeyTelegramBotApiKey:
			s.config.ApiKey = config.Value
			utils.Info("加载配置 %s = [HIDDEN]", config.Key)
		case entity.ConfigKeyTelegramAutoReplyEnabled:
			s.config.AutoReplyEnabled = config.Value == "true"
			utils.Info("加载配置 %s = %s (AutoReplyEnabled: %v)", config.Key, config.Value, s.config.AutoReplyEnabled)
		case entity.ConfigKeyTelegramAutoReplyTemplate:
			if config.Value != "" {
				s.config.AutoReplyTemplate = config.Value
			}
			utils.Info("加载配置 %s = %s", config.Key, config.Value)
		case entity.ConfigKeyTelegramAutoDeleteEnabled:
			s.config.AutoDeleteEnabled = config.Value == "true"
			utils.Info("加载配置 %s = %s (AutoDeleteEnabled: %v)", config.Key, config.Value, s.config.AutoDeleteEnabled)
		case entity.ConfigKeyTelegramAutoDeleteInterval:
			if config.Value != "" {
				fmt.Sscanf(config.Value, "%d", &s.config.AutoDeleteInterval)
			}
			utils.Info("加载配置 %s = %s (AutoDeleteInterval: %d)", config.Key, config.Value, s.config.AutoDeleteInterval)
		default:
			utils.Debug("未知配置: %s = %s", config.Key, config.Value)
		}
	}

	utils.Info("Telegram Bot 配置加载完成: Enabled=%v, AutoReplyEnabled=%v, ApiKey长度=%d",
		s.config.Enabled, s.config.AutoReplyEnabled, len(s.config.ApiKey))
	return nil
}

// Start 启动机器人服务
func (s *TelegramBotServiceImpl) Start() error {
	if s.isRunning {
		utils.Info("Telegram Bot 服务已经在运行中")
		return nil
	}

	// 加载配置
	if err := s.loadConfig(); err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if !s.config.Enabled || s.config.ApiKey == "" {
		utils.Info("Telegram Bot 未启用或 API Key 未配置")
		return nil
	}

	// 创建 Bot 实例
	bot, err := tgbotapi.NewBotAPI(s.config.ApiKey)
	if err != nil {
		return fmt.Errorf("创建 Telegram Bot 失败: %v", err)
	}

	s.bot = bot
	s.isRunning = true

	utils.Info("Telegram Bot (@%s) 已启动", s.GetBotUsername())

	// 启动推送调度器
	s.startContentPusher()

	// 设置 webhook（在实际部署时配置）
	if err := s.setupWebhook(); err != nil {
		utils.Error("设置 Webhook 失败: %v", err)
	}

	// 启动消息处理循环（长轮询模式）
	go s.messageLoop()

	return nil
}

// Stop 停止机器人服务
func (s *TelegramBotServiceImpl) Stop() error {
	if !s.isRunning {
		return nil
	}

	s.isRunning = false

	if s.cronScheduler != nil {
		s.cronScheduler.Stop()
	}

	utils.Info("Telegram Bot 服务已停止")
	return nil
}

// ValidateApiKey 验证 API Key
func (s *TelegramBotServiceImpl) ValidateApiKey(apiKey string) (bool, map[string]interface{}, error) {
	if apiKey == "" {
		return false, nil, fmt.Errorf("API Key 不能为空")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return false, nil, fmt.Errorf("无效的 API Key: %v", err)
	}

	// 获取机器人信息
	botInfo, err := bot.GetMe()
	if err != nil {
		return false, nil, fmt.Errorf("获取机器人信息失败: %v", err)
	}

	botData := map[string]interface{}{
		"id":         botInfo.ID,
		"username":   strings.TrimPrefix(botInfo.UserName, "@"),
		"first_name": botInfo.FirstName,
		"last_name":  botInfo.LastName,
	}

	return true, botData, nil
}

// setupWebhook 设置 Webhook（可选）
func (s *TelegramBotServiceImpl) setupWebhook() error {
	// 在生产环境中，这里会设置 webhook URL
	// 暂时使用长轮询模式，不设置 webhook
	utils.Info("使用长轮询模式处理消息")
	return nil
}

// messageLoop 消息处理循环（长轮询模式）
func (s *TelegramBotServiceImpl) messageLoop() {
	utils.Info("开始监听 Telegram 消息更新...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	utils.Info("消息监听循环已启动，等待消息...")

	for update := range updates {
		if update.Message != nil {
			utils.Info("接收到新消息更新")
			s.handleMessage(update.Message)
		} else {
			utils.Debug("接收到其他类型更新: %v", update)
		}
	}

	utils.Info("消息监听循环已结束")
}

// handleMessage 处理接收到的消息
func (s *TelegramBotServiceImpl) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := strings.TrimSpace(message.Text)

	utils.Info("收到消息: ChatID=%d, Text='%s', User=%s", chatID, text, message.From.UserName)

	if text == "" {
		return
	}

	// 处理 /register 命令
	if strings.ToLower(text) == "/register" {
		utils.Info("处理 /register 命令 from ChatID=%d", chatID)
		s.handleRegisterCommand(message)
		return
	}

	// 处理 /start 命令
	if strings.ToLower(text) == "/start" {
		utils.Info("处理 /start 命令 from ChatID=%d", chatID)
		s.handleStartCommand(message)
		return
	}

	// 处理普通文本消息（搜索请求）
	if len(text) > 0 && !strings.HasPrefix(text, "/") {
		utils.Info("处理搜索请求 from ChatID=%d: %s", chatID, text)
		s.handleSearchRequest(message)
		return
	}

	// 默认自动回复
	if s.config.AutoReplyEnabled {
		utils.Info("发送自动回复 to ChatID=%d (AutoReplyEnabled=%v)", chatID, s.config.AutoReplyEnabled)
		s.sendReply(message, s.config.AutoReplyTemplate)
	} else {
		utils.Info("跳过自动回复 to ChatID=%d (AutoReplyEnabled=%v)", chatID, s.config.AutoReplyEnabled)
	}
}

// handleRegisterCommand 处理注册命令
func (s *TelegramBotServiceImpl) handleRegisterCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	chatTitle := message.Chat.Title
	if chatTitle == "" {
		// 如果没有标题，使用用户名作为名称
		if message.Chat.UserName != "" {
			chatTitle = message.Chat.UserName
		} else {
			chatTitle = fmt.Sprintf("Chat_%d", chatID)
		}
	}

	chatType := "private"
	if message.Chat.IsChannel() {
		chatType = "channel"
	} else if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		chatType = "group"
	}

	err := s.RegisterChannel(chatID, chatTitle, chatType)
	if err != nil {
		errorMsg := fmt.Sprintf("注册失败: %v", err)
		s.sendReply(message, errorMsg)
		return
	}

	successMsg := fmt.Sprintf("✅ 注册成功！\n\n频道/群组: %s\n类型: %s\n\n现在可以向此频道推送资源内容了。", chatTitle, chatType)
	s.sendReply(message, successMsg)
}

// handleStartCommand 处理开始命令
func (s *TelegramBotServiceImpl) handleStartCommand(message *tgbotapi.Message) {
	welcomeMsg := `🤖 欢迎使用网盘资源机器人！

我会帮您搜索网盘资源。使用方法：
• 直接发送关键词搜索资源
• 发送 /register 注册当前频道用于推送

享受使用吧！`

	if s.config.AutoReplyEnabled && s.config.AutoReplyTemplate != "" {
		welcomeMsg += "\n\n" + s.config.AutoReplyTemplate
	}

	s.sendReply(message, welcomeMsg)
}

// handleSearchRequest 处理搜索请求
func (s *TelegramBotServiceImpl) handleSearchRequest(message *tgbotapi.Message) {
	query := strings.TrimSpace(message.Text)
	if query == "" {
		s.sendReply(message, "请输入搜索关键词")
		return
	}

	// 这里使用简单的资源搜索，实际项目中需要完善搜索逻辑
	// resources, err := s.resourceRepo.Search(query, nil, 0, 10)
	// 暂时模拟一个搜索结果
	results := []string{
		fmt.Sprintf("🔍 搜索关键词: %s", query),
		"暂无相关资源，请尝试其他关键词。",
		"",
		fmt.Sprintf("💡 提示：如需精确搜索，请使用更具体的关键词。"),
	}

	resultText := strings.Join(results, "\n")
	s.sendReply(message, resultText)
}

// sendReply 发送回复消息
func (s *TelegramBotServiceImpl) sendReply(message *tgbotapi.Message, text string) {
	utils.Info("尝试发送回复消息到 ChatID=%d: %s", message.Chat.ID, text)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID

	sentMsg, err := s.bot.Send(msg)
	if err != nil {
		utils.Error("发送消息失败: %v", err)
		return
	}

	utils.Info("消息发送成功 to ChatID=%d, MessageID=%d", sentMsg.Chat.ID, sentMsg.MessageID)

	// 如果启用了自动删除，启动删除定时器
	if s.config.AutoDeleteEnabled && s.config.AutoDeleteInterval > 0 {
		time.AfterFunc(time.Duration(s.config.AutoDeleteInterval)*time.Minute, func() {
			deleteConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    sentMsg.Chat.ID,
				MessageID: sentMsg.MessageID,
			}
			_, err := s.bot.Request(deleteConfig)
			if err != nil {
				utils.Error("删除消息失败: %v", err)
			}
		})
	}
}

// startContentPusher 启动内容推送器
func (s *TelegramBotServiceImpl) startContentPusher() {
	// 每小时检查一次需要推送的频道
	s.cronScheduler.AddFunc("@every 1h", func() {
		s.pushContentToChannels()
	})

	s.cronScheduler.Start()
	utils.Info("内容推送调度器已启动")
}

// pushContentToChannels 推送内容到频道
func (s *TelegramBotServiceImpl) pushContentToChannels() {
	// 获取需要推送的频道
	channels, err := s.channelRepo.FindDueForPush()
	if err != nil {
		utils.Error("获取推送频道失败: %v", err)
		return
	}

	if len(channels) == 0 {
		utils.Debug("没有需要推送的频道")
		return
	}

	utils.Info("开始推送内容到 %d 个频道", len(channels))

	for _, channel := range channels {
		go s.pushToChannel(channel)
	}
}

// pushToChannel 推送内容到一个频道
func (s *TelegramBotServiceImpl) pushToChannel(channel entity.TelegramChannel) {
	// 这里实现推送逻辑
	// 1. 根据频道设置过滤资源
	resources := s.findResourcesForChannel(channel)
	if len(resources) == 0 {
		utils.Debug("频道 %s 没有可推送的内容", channel.ChatName)
		return
	}

	// 2. 构建推送消息
	message := s.buildPushMessage(channel, resources)

	// 3. 发送消息
	err := s.SendMessage(channel.ChatID, message)
	if err != nil {
		utils.Error("推送失败到频道 %s (%d): %v", channel.ChatName, channel.ChatID, err)
		// 可以考虑将频道标记为非活跃或记录错误
		return
	}

	// 4. 更新最后推送时间
	err = s.channelRepo.UpdateLastPushAt(channel.ID, time.Now())
	if err != nil {
		utils.Error("更新推送时间失败: %v", err)
	}

	utils.Info("成功推送内容到频道: %s", channel.ChatName)
}

// findResourcesForChannel 查找适合频道的资源
func (s *TelegramBotServiceImpl) findResourcesForChannel(channel entity.TelegramChannel) []interface{} {
	// 这里需要实现根据频道配置过滤资源
	// 暂时返回空数组，实际实现中需要查询资源数据库
	return []interface{}{}
}

// buildPushMessage 构建推送消息
func (s *TelegramBotServiceImpl) buildPushMessage(channel entity.TelegramChannel, resources []interface{}) string {
	message := fmt.Sprintf("📢 **%s**\n\n", channel.ChatName)

	if len(resources) == 0 {
		message += "暂无新内容推送"
	} else {
		message += fmt.Sprintf("🆕 发现 %d 个新资源:\n\n", len(resources))
		// 这里需要格式化资源列表
		message += "*详细资源列表请查看网站*"
	}

	message += fmt.Sprintf("\n\n⏰ 下次推送: %d 小时后", channel.PushFrequency)

	return message
}

// GetBotUsername 获取机器人用户名
func (s *TelegramBotServiceImpl) GetBotUsername() string {
	if s.bot != nil {
		return s.bot.Self.UserName
	}
	return ""
}

// SendMessage 发送消息
func (s *TelegramBotServiceImpl) SendMessage(chatID int64, text string) error {
	if s.bot == nil {
		return fmt.Errorf("Bot 未初始化")
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := s.bot.Send(msg)
	return err
}

// DeleteMessage 删除消息
func (s *TelegramBotServiceImpl) DeleteMessage(chatID int64, messageID int) error {
	if s.bot == nil {
		return fmt.Errorf("Bot 未初始化")
	}

	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := s.bot.Request(deleteConfig)
	return err
}

// RegisterChannel 注册频道
func (s *TelegramBotServiceImpl) RegisterChannel(chatID int64, chatName, chatType string) error {
	// 检查是否已注册
	if s.IsChannelRegistered(chatID) {
		return fmt.Errorf("该频道/群组已注册")
	}

	channel := entity.TelegramChannel{
		ChatID:            chatID,
		ChatName:          chatName,
		ChatType:          chatType,
		PushEnabled:       true,
		PushFrequency:     24, // 默认24小时
		IsActive:          true,
		RegisteredBy:      "bot_command",
		RegisteredAt:      time.Now(),
		ContentCategories: "",
		ContentTags:       "",
	}

	return s.channelRepo.Create(&channel)
}

// IsChannelRegistered 检查频道是否已注册
func (s *TelegramBotServiceImpl) IsChannelRegistered(chatID int64) bool {
	channel, err := s.channelRepo.FindByChatID(chatID)
	return err == nil && channel != nil
}

// HandleWebhookUpdate 处理 Webhook 更新（预留接口，目前使用长轮询）
func (s *TelegramBotServiceImpl) HandleWebhookUpdate(c interface{}) {
	// 目前使用长轮询模式，webhook 接口预留
	// 将来可以实现从 webhook 接收消息的处理逻辑
	// 如果需要实现 webhook 模式，可以在这里添加处理逻辑
}
