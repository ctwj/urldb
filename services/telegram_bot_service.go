package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
	"golang.org/x/net/proxy"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type TelegramBotService interface {
	Start() error
	Stop() error
	IsRunning() bool
	ReloadConfig() error
	GetRuntimeStatus() map[string]interface{}
	ValidateApiKey(apiKey string) (bool, map[string]interface{}, error)
	ValidateApiKeyWithProxy(apiKey string, proxyEnabled bool, proxyType, proxyHost string, proxyPort int, proxyUsername, proxyPassword string) (bool, map[string]interface{}, error)
	GetBotUsername() string
	SendMessage(chatID int64, text string) error
	SendMessageWithFormat(chatID int64, text string, parseMode string) error
	DeleteMessage(chatID int64, messageID int) error
	RegisterChannel(chatID int64, chatName, chatType string) error
	IsChannelRegistered(chatID int64) bool
	HandleWebhookUpdate(c interface{})
	CleanupDuplicateChannels() error
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
	ProxyEnabled       bool
	ProxyType          string // http, https, socks5
	ProxyHost          string
	ProxyPort          int
	ProxyUsername      string
	ProxyPassword      string
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

	utils.Info("[TELEGRAM] 从数据库加载到 %d 个配置项", len(configs))

	// 初始化默认值
	s.config.Enabled = false
	s.config.ApiKey = ""
	s.config.AutoReplyEnabled = false // 默认禁用自动回复
	s.config.AutoReplyTemplate = "您好！我可以帮您搜索网盘资源，请输入您要搜索的内容。"
	s.config.AutoDeleteEnabled = false
	s.config.AutoDeleteInterval = 60
	// 初始化代理默认值
	s.config.ProxyEnabled = false
	s.config.ProxyType = "http"
	s.config.ProxyHost = ""
	s.config.ProxyPort = 8080
	s.config.ProxyUsername = ""
	s.config.ProxyPassword = ""

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyTelegramBotEnabled:
			s.config.Enabled = config.Value == "true"
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (Enabled: %v)", config.Key, config.Value, s.config.Enabled)
		case entity.ConfigKeyTelegramBotApiKey:
			s.config.ApiKey = config.Value
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = [HIDDEN]", config.Key)
		case entity.ConfigKeyTelegramAutoReplyEnabled:
			s.config.AutoReplyEnabled = config.Value == "true"
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (AutoReplyEnabled: %v)", config.Key, config.Value, s.config.AutoReplyEnabled)
		case entity.ConfigKeyTelegramAutoReplyTemplate:
			if config.Value != "" {
				s.config.AutoReplyTemplate = config.Value
			}
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s", config.Key, config.Value)
		case entity.ConfigKeyTelegramAutoDeleteEnabled:
			s.config.AutoDeleteEnabled = config.Value == "true"
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (AutoDeleteEnabled: %v)", config.Key, config.Value, s.config.AutoDeleteEnabled)
		case entity.ConfigKeyTelegramAutoDeleteInterval:
			if config.Value != "" {
				fmt.Sscanf(config.Value, "%d", &s.config.AutoDeleteInterval)
			}
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (AutoDeleteInterval: %d)", config.Key, config.Value, s.config.AutoDeleteInterval)
		case entity.ConfigKeyTelegramProxyEnabled:
			s.config.ProxyEnabled = config.Value == "true"
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (ProxyEnabled: %v)", config.Key, config.Value, s.config.ProxyEnabled)
		case entity.ConfigKeyTelegramProxyType:
			s.config.ProxyType = config.Value
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (ProxyType: %s)", config.Key, config.Value, s.config.ProxyType)
		case entity.ConfigKeyTelegramProxyHost:
			s.config.ProxyHost = config.Value
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s", config.Key, "[HIDDEN]")
		case entity.ConfigKeyTelegramProxyPort:
			if config.Value != "" {
				fmt.Sscanf(config.Value, "%d", &s.config.ProxyPort)
			}
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s (ProxyPort: %d)", config.Key, config.Value, s.config.ProxyPort)
		case entity.ConfigKeyTelegramProxyUsername:
			s.config.ProxyUsername = config.Value
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s", config.Key, "[HIDDEN]")
		case entity.ConfigKeyTelegramProxyPassword:
			s.config.ProxyPassword = config.Value
			utils.Info("[TELEGRAM:CONFIG] 加载配置 %s = %s", config.Key, "[HIDDEN]")
		default:
			utils.Debug("未知配置: %s = %s", config.Key, config.Value)
		}
	}

	utils.Info("[TELEGRAM:SERVICE] Telegram Bot 配置加载完成: Enabled=%v, AutoReplyEnabled=%v, ApiKey长度=%d",
		s.config.Enabled, s.config.AutoReplyEnabled, len(s.config.ApiKey))
	return nil
}

// Start 启动机器人服务
func (s *TelegramBotServiceImpl) Start() error {
	if s.isRunning {
		utils.Info("[TELEGRAM:SERVICE] Telegram Bot 服务已经在运行中")
		return nil
	}

	// 加载配置
	if err := s.loadConfig(); err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if !s.config.Enabled || s.config.ApiKey == "" {
		utils.Info("[TELEGRAM:SERVICE] Telegram Bot 未启用或 API Key 未配置")
		return nil
	}

	// 创建 Bot 实例
	var bot *tgbotapi.BotAPI

	if s.config.ProxyEnabled && s.config.ProxyHost != "" {
		// 配置代理
		utils.Info("[TELEGRAM:PROXY] 配置代理: %s://%s:%d", s.config.ProxyType, s.config.ProxyHost, s.config.ProxyPort)

		var httpClient *http.Client

		if s.config.ProxyType == "socks5" {
			// SOCKS5 代理配置
			var auth *proxy.Auth
			if s.config.ProxyUsername != "" {
				auth = &proxy.Auth{
					User:     s.config.ProxyUsername,
					Password: s.config.ProxyPassword,
				}
			}

			dialer, proxyErr := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", s.config.ProxyHost, s.config.ProxyPort), auth, proxy.Direct)
			if proxyErr != nil {
				return fmt.Errorf("创建 SOCKS5 代理失败: %v", proxyErr)
			}

			httpClient = &http.Client{
				Transport: &http.Transport{
					Dial: dialer.Dial,
				},
				Timeout: 30 * time.Second,
			}
		} else {
			// HTTP/HTTPS 代理配置
			proxyURL := &url.URL{
				Scheme: s.config.ProxyType,
				Host:   fmt.Sprintf("%s:%d", s.config.ProxyHost, s.config.ProxyPort),
				User:   nil,
			}

			if s.config.ProxyUsername != "" {
				proxyURL.User = url.UserPassword(s.config.ProxyUsername, s.config.ProxyPassword)
			}

			httpClient = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
				Timeout: 30 * time.Second,
			}
		}

		botInstance, botErr := tgbotapi.NewBotAPIWithClient(s.config.ApiKey, tgbotapi.APIEndpoint, httpClient)
		if botErr != nil {
			return fmt.Errorf("创建 Telegram Bot (代理模式) 失败: %v", botErr)
		}
		bot = botInstance

		utils.Info("[TELEGRAM:PROXY] Telegram Bot 已配置代理连接")
	} else {
		// 直接连接（无代理）
		var err error
		bot, err = tgbotapi.NewBotAPI(s.config.ApiKey)
		if err != nil {
			return fmt.Errorf("创建 Telegram Bot 失败: %v", err)
		}

		utils.Info("[TELEGRAM:PROXY] Telegram Bot 使用直连模式")
	}

	s.bot = bot
	s.isRunning = true

	utils.Info("[TELEGRAM:SERVICE] Telegram Bot (@%s) 已启动", s.GetBotUsername())

	// 启动推送调度器
	s.startContentPusher()

	// 设置 webhook（在实际部署时配置）
	if err := s.setupWebhook(); err != nil {
		utils.Error("[TELEGRAM:SERVICE] 设置 Webhook 失败: %v", err)
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

	utils.Info("[TELEGRAM:SERVICE] Telegram Bot 服务已停止")
	return nil
}

// IsRunning 检查机器人服务是否正在运行
func (s *TelegramBotServiceImpl) IsRunning() bool {
	return s.isRunning && s.bot != nil
}

// ReloadConfig 重新加载机器人配置
func (s *TelegramBotServiceImpl) ReloadConfig() error {
	utils.Info("[TELEGRAM:SERVICE] 开始重新加载配置...")

	// 重新加载配置
	if err := s.loadConfig(); err != nil {
		utils.Error("[TELEGRAM:SERVICE] 重新加载配置失败: %v", err)
		return fmt.Errorf("重新加载配置失败: %v", err)
	}

	utils.Info("[TELEGRAM:SERVICE] 配置重新加载完成: Enabled=%v, AutoReplyEnabled=%v",
		s.config.Enabled, s.config.AutoReplyEnabled)
	return nil
}

// GetRuntimeStatus 获取机器人运行时状态
func (s *TelegramBotServiceImpl) GetRuntimeStatus() map[string]interface{} {
	status := map[string]interface{}{
		"is_running":      s.IsRunning(),
		"bot_initialized": s.bot != nil,
		"config_loaded":   s.config != nil,
		"cron_running":    s.cronScheduler != nil,
		"username":        "",
		"uptime":          0,
	}

	if s.bot != nil {
		status["username"] = s.GetBotUsername()
	}

	return status
}

// ValidateApiKey 验证 API Key
func (s *TelegramBotServiceImpl) ValidateApiKey(apiKey string) (bool, map[string]interface{}, error) {
	if apiKey == "" {
		return false, nil, fmt.Errorf("API Key 不能为空")
	}

	var bot *tgbotapi.BotAPI
	var err error

	// 如果启用了代理，使用代理验证
	if s.config.ProxyEnabled && s.config.ProxyHost != "" {
		var httpClient *http.Client

		if s.config.ProxyType == "socks5" {
			var auth *proxy.Auth
			if s.config.ProxyUsername != "" {
				auth = &proxy.Auth{
					User:     s.config.ProxyUsername,
					Password: s.config.ProxyPassword,
				}
			}

			dialer, proxyErr := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", s.config.ProxyHost, s.config.ProxyPort), auth, proxy.Direct)
			if proxyErr != nil {
				// 如果代理失败，回退到直连
				utils.Warn("[TELEGRAM:PROXY] SOCKS5 代理验证失败，回退到直连: %v", proxyErr)
				bot, err = tgbotapi.NewBotAPI(apiKey)
			} else {
				httpClient = &http.Client{
					Transport: &http.Transport{
						Dial: dialer.Dial,
					},
					Timeout: 10 * time.Second,
				}
				bot, err = tgbotapi.NewBotAPIWithClient(apiKey, tgbotapi.APIEndpoint, httpClient)
			}
		} else {
			proxyURL := &url.URL{
				Scheme: s.config.ProxyType,
				Host:   fmt.Sprintf("%s:%d", s.config.ProxyHost, s.config.ProxyPort),
				User:   nil,
			}

			if s.config.ProxyUsername != "" {
				proxyURL.User = url.UserPassword(s.config.ProxyUsername, s.config.ProxyPassword)
			}

			httpClient = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
				Timeout: 10 * time.Second,
			}
			bot, err = tgbotapi.NewBotAPIWithClient(apiKey, tgbotapi.APIEndpoint, httpClient)
		}
	} else {
		// 直连验证
		bot, err = tgbotapi.NewBotAPI(apiKey)
	}

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

// ValidateApiKeyWithProxy 使用代理配置验证 API Key
func (s *TelegramBotServiceImpl) ValidateApiKeyWithProxy(apiKey string, proxyEnabled bool, proxyType, proxyHost string, proxyPort int, proxyUsername, proxyPassword string) (bool, map[string]interface{}, error) {
	if apiKey == "" {
		return false, nil, fmt.Errorf("API Key 不能为空")
	}

	var bot *tgbotapi.BotAPI
	var err error

	// 使用提供的代理配置进行校验
	if proxyEnabled && proxyHost != "" {
		var httpClient *http.Client

		if proxyType == "socks5" {
			var auth *proxy.Auth
			if proxyUsername != "" {
				auth = &proxy.Auth{
					User:     proxyUsername,
					Password: proxyPassword,
				}
			}

			dialer, proxyErr := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", proxyHost, proxyPort), auth, proxy.Direct)
			if proxyErr != nil {
				return false, nil, fmt.Errorf("创建 SOCKS5 代理失败: %v", proxyErr)
			}

			httpClient = &http.Client{
				Transport: &http.Transport{
					Dial: dialer.Dial,
				},
				Timeout: 10 * time.Second,
			}
		} else {
			proxyURL := &url.URL{
				Scheme: proxyType,
				Host:   fmt.Sprintf("%s:%d", proxyHost, proxyPort),
				User:   nil,
			}

			if proxyUsername != "" {
				proxyURL.User = url.UserPassword(proxyUsername, proxyPassword)
			}

			httpClient = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
				Timeout: 10 * time.Second,
			}
		}

		bot, err = tgbotapi.NewBotAPIWithClient(apiKey, tgbotapi.APIEndpoint, httpClient)
		if err != nil {
			utils.Error(fmt.Sprintf("[TELEGRAM:VALIDATE] 创建 Telegram Bot (代理校验) 失败 $v", err))
			return false, nil, fmt.Errorf("创建 Telegram Bot (代理校验) 失败: %v", err)
		}

		utils.Info("[TELEGRAM:VALIDATE] 使用代理配置校验 API Key")
	} else {
		// 直连校验
		bot, err = tgbotapi.NewBotAPI(apiKey)
		if err != nil {
			utils.Error(fmt.Sprintf("[TELEGRAM:VALIDATE] 创建 Telegram Bot 失败 $v", err))
			return false, nil, fmt.Errorf("无效的 API Key: %v", err)
		}

		utils.Info("[TELEGRAM:VALIDATE] 使用直连模式校验 API Key")
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
	utils.Info("[TELEGRAM:SERVICE] 使用长轮询模式处理消息")
	return nil
}

// messageLoop 消息处理循环（长轮询模式）
func (s *TelegramBotServiceImpl) messageLoop() {
	utils.Info("[TELEGRAM:MESSAGE] 开始监听 Telegram 消息更新...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	utils.Info("[TELEGRAM:MESSAGE] 消息监听循环已启动，等待消息...")

	for update := range updates {
		if update.Message != nil {
			utils.Info("[TELEGRAM:MESSAGE] 接收到新消息更新")
			s.handleMessage(update.Message)
		} else {
			utils.Debug("[TELEGRAM:MESSAGE] 接收到其他类型更新: %v", update)
		}
	}

	utils.Info("[TELEGRAM:MESSAGE] 消息监听循环已结束")
}

// handleMessage 处理接收到的消息
func (s *TelegramBotServiceImpl) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := strings.TrimSpace(message.Text)

	utils.Info("[TELEGRAM:MESSAGE] 收到消息: ChatID=%d, Text='%s', User=%s", chatID, text, message.From.UserName)

	if text == "" {
		return
	}

	// 处理 /register 命令（包括参数）
	if strings.HasPrefix(strings.ToLower(text), "/register") {
		utils.Info("[TELEGRAM:MESSAGE] 处理 /register 命令 from ChatID=%d", chatID)
		s.handleRegisterCommand(message)
		return
	}

	// 处理 /start 命令
	if strings.ToLower(text) == "/start" {
		utils.Info("[TELEGRAM:MESSAGE] 处理 /start 命令 from ChatID=%d", chatID)
		s.handleStartCommand(message)
		return
	}

	// 处理普通文本消息（搜索请求）
	if len(text) > 0 && !strings.HasPrefix(text, "/") {
		utils.Info("[TELEGRAM:MESSAGE] 处理搜索请求 from ChatID=%d: %s", chatID, text)
		s.handleSearchRequest(message)
		return
	}

	// 默认自动回复（只对正常消息，不对转发消息，且消息没有换行）
	if s.config.AutoReplyEnabled {
		// 检查是否是转发消息
		isForward := message.ForwardFrom != nil ||
			message.ForwardFromChat != nil ||
			message.ForwardDate != 0

		if isForward {
			utils.Info("[TELEGRAM:MESSAGE] 跳过自动回复，转发消息 from ChatID=%d", chatID)
		} else {
			// 检查消息是否包含换行符
			hasNewLine := strings.Contains(text, "\n") || strings.Contains(text, "\r")

			if hasNewLine {
				utils.Info("[TELEGRAM:MESSAGE] 跳过自动回复，消息包含换行 from ChatID=%d", chatID)
			} else {
				utils.Info("[TELEGRAM:MESSAGE] 发送自动回复 to ChatID=%d (AutoReplyEnabled=%v)", chatID, s.config.AutoReplyEnabled)
				s.sendReply(message, s.config.AutoReplyTemplate)
			}
		}
	} else {
		utils.Info("[TELEGRAM:MESSAGE] 跳过自动回复 to ChatID=%d (AutoReplyEnabled=%v)", chatID, s.config.AutoReplyEnabled)
	}
}

// handleRegisterCommand 处理注册命令
func (s *TelegramBotServiceImpl) handleRegisterCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := strings.TrimSpace(message.Text)

	// 检查是否是群组
	isGroup := message.Chat.IsGroup() || message.Chat.IsSuperGroup()

	if isGroup {
		// 群组中需要管理员权限
		if !s.isUserAdministrator(message.Chat.ID, message.From.ID) {
			errorMsg := "❌ *权限不足*\n\n只有群组管理员才能注册此群组用于推送。\n\n请联系管理员执行注册命令。"
			s.sendReply(message, errorMsg)
			return
		}

		// 检查是否已经注册了群组
		if s.hasActiveGroup() {
			errorMsg := "❌ *注册限制*\n\n系统最多只支持注册一个群组用于推送。\n\n请先注销现有群组，然后再注册新的群组。"
			s.sendReply(message, errorMsg)
			return
		}

		// 注册群组
		chatTitle := message.Chat.Title
		if chatTitle == "" {
			chatTitle = fmt.Sprintf("Group_%d", chatID)
		}

		err := s.RegisterChannel(chatID, chatTitle, "group")
		if err != nil {
			if strings.Contains(err.Error(), "该频道/群组已注册") {
				successMsg := fmt.Sprintf("⚠️ *群组已注册*\n\n群组: %s\n类型: 群组\n\n此群组已经注册，无需重复注册。", chatTitle)
				s.sendReply(message, successMsg)
			} else {
				errorMsg := fmt.Sprintf("❌ 注册失败: %v", err)
				s.sendReply(message, errorMsg)
			}
			return
		}

		successMsg := fmt.Sprintf("✅ *群组注册成功！*\n\n群组: %s\n类型: 群组\n\n现在可以向此群组推送资源内容了。", chatTitle)
		s.sendReply(message, successMsg)
		return
	}

	// 私聊处理
	parts := strings.Fields(text)

	if len(parts) == 1 {
		// 私聊中没有参数，显示注册帮助
		helpMsg := `🤖 *注册帮助*
*注册群组:*
* 添加机器人，为频道管理员
* 管理员发送 /register 命令

*注册频道:*
私聊机器人， 发送注册命令
支持两种格式：
• /register <频道ID> - 如: /register -1001234567890
• /register @用户名 - 如: /register @xypan

*获取频道ID的方法:*
1. 将机器人添加到频道并设为管理员
2. 向频道发送消息，查看机器人收到的消息
3. 频道ID通常是负数，如 -1001234567890

*示例:*
/register -1001234567890
/register @xypan

*注意:*
• 频道ID必须是纯数字（包括负号）
• 用户名格式必须以 @ 开头
• 机器人必须是频道的管理员才能注册
• 私聊不支持注册，只支持频道和群组注册`
		s.sendReply(message, helpMsg)
	} else if parts[1] == "help" || parts[1] == "-h" {
		// 显示注册帮助
		helpMsg := `🤖 *注册帮助*
*注册群组:*
* 添加机器人，为频道管理员
* 管理员发送 /register 命令

*注册频道:*
私聊机器人， 发送注册命令
支持两种格式：
• /register <频道ID> - 如: /register -1001234567890
• /register @用户名 - 如: /register @xypan

*获取频道ID的方法:*
1. 将机器人添加到频道并设为管理员
2. 向频道发送消息，查看机器人收到的消息
3. 频道ID通常是负数，如 -1001234567890

*示例:*
/register -1001234567890
/register @xypan

*注意:*
• 频道ID必须是纯数字（包括负号）
• 用户名格式必须以 @ 开头
• 机器人必须是频道的管理员才能注册`
		s.sendReply(message, helpMsg)
	} else {
		// 有参数，尝试注册频道
		channelIDStr := strings.TrimSpace(parts[1])
		s.handleChannelRegistration(message, channelIDStr)
	}
}

// handleStartCommand 处理开始命令
func (s *TelegramBotServiceImpl) handleStartCommand(message *tgbotapi.Message) {
	welcomeMsg := `🤖 欢迎使用老九网盘资源机器人！

• 发送 搜索 + 关键词 进行资源搜索
• 发送 /register 注册当前频道或群组，用于主动推送资源
• 私聊中使用 /register help 获取注册帮助
• 发送 /start 获取帮助信息
`

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

	utils.Info("[TELEGRAM:SEARCH] 处理搜索请求: %s", query)

	// 使用资源仓库进行搜索
	resources, total, err := s.resourceRepo.Search(query, nil, 1, 5) // 限制为5个结果
	if err != nil {
		utils.Error("[TELEGRAM:SEARCH] 搜索失败: %v", err)
		s.sendReply(message, "搜索服务暂时不可用，请稍后重试")
		return
	}

	if total == 0 {
		response := fmt.Sprintf("🔍 *搜索结果*\n\n关键词: `%s`\n\n❌ 未找到相关资源\n\n💡 建议:\n• 尝试使用更通用的关键词\n• 检查拼写是否正确\n• 减少关键词数量", query)
		// 没有找到资源，不使用资源自动删除
		s.sendReply(message, response)
		return
	}

	// 构建搜索结果消息
	resultText := fmt.Sprintf("🔍 *搜索结果*\n\n关键词: `%s`\n总共找到: %d 个资源\n\n", query, total)

	// 显示前5个结果
	for i, resource := range resources {
		if i >= 5 {
			break
		}

		// 清理资源标题和描述，确保UTF-8编码
		title := s.cleanResourceText(resource.Title)
		if len(title) > 50 {
			title = title[:47] + "..."
		}

		description := s.cleanResourceText(resource.Description)
		if len(description) > 100 {
			description = description[:97] + "..."
		}

		resultText += fmt.Sprintf("%d. *%s*\n%s\n\n", i+1, title, description)
	}

	// 如果有更多结果，添加提示
	if total > 5 {
		resultText += fmt.Sprintf("... 还有 %d 个结果\n\n", total-5)
		resultText += "💡 如需查看更多结果，请访问网站搜索"
	}

	// 使用包含资源的自动删除功能
	s.sendReplyWithResourceAutoDelete(message, resultText, len(resources))
}

// sendReply 发送回复消息
func (s *TelegramBotServiceImpl) sendReply(message *tgbotapi.Message, text string) {
	s.sendReplyWithAutoDelete(message, text, s.config.AutoDeleteEnabled)
}

// sendReplyWithAutoDelete 发送回复消息，支持指定是否自动删除
func (s *TelegramBotServiceImpl) sendReplyWithAutoDelete(message *tgbotapi.Message, text string, autoDelete bool) {
	// 清理消息文本，确保UTF-8编码
	originalText := text
	text = s.cleanMessageText(text)
	utils.Info("[TELEGRAM:MESSAGE] 尝试发送回复消息到 ChatID=%d, 原始长度=%d, 清理后长度=%d", message.Chat.ID, len(originalText), len(text))

	// 检查清理后的文本是否有效
	if len(text) == 0 {
		utils.Error("[TELEGRAM:MESSAGE:ERROR] 清理后消息为空，无法发送")
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyToMessageID = message.MessageID

	utils.Debug("[TELEGRAM:MESSAGE] 发送Markdown版本消息: %s", text[:min(100, len(text))])

	sentMsg, err := s.bot.Send(msg)
	if err != nil {
		utils.Error("[TELEGRAM:MESSAGE:ERROR] 发送Markdown消息失败: %v", err)
		// 如果是UTF-8编码错误或Markdown错误，尝试发送纯文本版本
		if strings.Contains(err.Error(), "UTF-8") || strings.Contains(err.Error(), "Bad Request") || strings.Contains(err.Error(), "strings must be encoded") {
			utils.Info("[TELEGRAM:MESSAGE] 尝试发送纯文本版本...")
			plainText := s.cleanMessageTextForPlain(originalText)
			utils.Debug("[TELEGRAM:MESSAGE] 发送纯文本版本消息: %s", plainText[:min(100, len(plainText))])

			msg.ParseMode = ""
			msg.Text = plainText
			sentMsg, err = s.bot.Send(msg)
			if err != nil {
				utils.Error("[TELEGRAM:MESSAGE:ERROR] 纯文本发送也失败: %v", err)
				return
			}
		} else {
			return
		}
	}

	utils.Info("[TELEGRAM:MESSAGE:SUCCESS] 消息发送成功 to ChatID=%d, MessageID=%d", sentMsg.Chat.ID, sentMsg.MessageID)

	// 如果启用了自动删除，启动删除定时器
	if autoDelete && s.config.AutoDeleteInterval > 0 {
		utils.Info("[TELEGRAM:MESSAGE] 设置自动删除定时器: %d 分钟后删除消息", s.config.AutoDeleteInterval)
		time.AfterFunc(time.Duration(s.config.AutoDeleteInterval)*time.Minute, func() {
			deleteConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    sentMsg.Chat.ID,
				MessageID: sentMsg.MessageID,
			}
			_, err := s.bot.Request(deleteConfig)
			if err != nil {
				utils.Error("[TELEGRAM:MESSAGE:ERROR] 删除消息失败: %v", err)
			} else {
				utils.Info("[TELEGRAM:MESSAGE] 消息已自动删除: ChatID=%d, MessageID=%d", sentMsg.Chat.ID, sentMsg.MessageID)
			}
		})
	}
}

// 辅助函数：返回两个数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// cleanMessageText 清理消息文本，确保UTF-8编码和Markdown格式兼容
func (s *TelegramBotServiceImpl) cleanMessageText(text string) string {
	if text == "" {
		return text
	}

	// 记录原始消息用于调试
	utils.Debug("[TELEGRAM:CLEAN] 原始消息长度: %d", len(text))

	// 清理Markdown特殊字符
	text = strings.ReplaceAll(text, "\\", "\\\\") // 转义反斜杠
	text = strings.ReplaceAll(text, "*", "\\*")   // 转义星号
	text = strings.ReplaceAll(text, "_", "\\_")   // 转义下划线
	text = strings.ReplaceAll(text, "`", "\\`")   // 转义反引号
	text = strings.ReplaceAll(text, "[", "\\[")   // 转义方括号
	text = strings.ReplaceAll(text, "]", "\\]")   // 转义方括号

	// 移除可能的控制字符
	text = strings.Map(func(r rune) rune {
		if r < 32 && r != 9 && r != 10 && r != 13 { // 保留tab、换行、回车
			return -1 // 删除控制字符
		}
		return r
	}, text)

	// 限制消息长度（Telegram单条消息最大4096字符）
	if len(text) > 4000 {
		text = text[:4000] + "..."
		utils.Debug("[TELEGRAM:CLEAN] 消息已截断，长于4000字符")
	}

	utils.Debug("[TELEGRAM:CLEAN] 清理后消息长度: %d", len(text))
	return text
}

// cleanMessageTextForPlain 清理消息文本为纯文本格式
func (s *TelegramBotServiceImpl) cleanMessageTextForPlain(text string) string {
	if text == "" {
		return "空消息"
	}

	utils.Debug("[TELEGRAM:CLEAN:PLAIN] 原始纯文本消息长度: %d", len(text))

	// 移除Markdown格式字符
	text = strings.ReplaceAll(text, "*", "")  // 移除粗体
	text = strings.ReplaceAll(text, "_", "")  // 移除斜体
	text = strings.ReplaceAll(text, "`", "")  // 移除代码
	text = strings.ReplaceAll(text, "[", "(") // 替换链接开始
	text = strings.ReplaceAll(text, "]", ")") // 替换链接结束
	text = strings.ReplaceAll(text, "\\", "") // 移除转义符

	// 移除可能的控制字符
	text = strings.Map(func(r rune) rune {
		if r < 32 && r != 9 && r != 10 && r != 13 { // 保留tab、换行、回车
			return -1 // 删除控制字符
		}
		return r
	}, text)

	// 如果清理后消息为空，返回默认消息
	if strings.TrimSpace(text) == "" {
		text = "消息内容无法显示"
	}

	// 限制消息长度
	if len(text) > 4000 {
		text = text[:4000] + "..."
		utils.Debug("[TELEGRAM:CLEAN:PLAIN] 纯文本消息已截断，长于4000字符")
	}

	utils.Debug("[TELEGRAM:CLEAN:PLAIN] 清理后纯文本消息长度: %d", len(text))
	return text
}

// cleanResourceText 清理从数据库读取的资源文本
func (s *TelegramBotServiceImpl) cleanResourceText(text string) string {
	if text == "" {
		return text
	}

	// 记录原始文本用于调试（只记录前50字符避免日志过长）
	debugText := text
	if len(text) > 50 {
		debugText = text[:47] + "..."
	}
	utils.Debug("[TELEGRAM:CLEAN:RESOURCE] 原始资源文本: %s", debugText)

	// 移除可能的控制字符，但保留中文字符
	text = strings.Map(func(r rune) rune {
		if r < 32 && r != 9 && r != 10 && r != 13 { // 保留tab、换行、回车
			return -1 // 删除控制字符
		}
		// 注意：不再移除超出BMP的字符，因为中文字符可能需要这些码点
		return r
	}, text)

	// 移除零宽度字符和其他不可见字符，但保留中文字符
	text = strings.ReplaceAll(text, "\u200B", "") // 零宽度空格
	text = strings.ReplaceAll(text, "\u200C", "") // 零宽度非连接符
	text = strings.ReplaceAll(text, "\u200D", "") // 零宽度连接符
	text = strings.ReplaceAll(text, "\uFEFF", "") // 字节顺序标记

	// 移除其他可能的垃圾字符，但非常保守
	text = strings.ReplaceAll(text, "\u0000", "") // 空字符
	text = strings.ReplaceAll(text, "\uFFFD", "") // 替换字符

	// 如果清理后为空，返回默认文本
	if strings.TrimSpace(text) == "" {
		text = "无标题"
	}

	utils.Debug("[TELEGRAM:CLEAN:RESOURCE] 清理后资源文本长度: %d", len(text))
	return text
}

// sendReplyWithResourceAutoDelete 发送包含资源的回复消息，自动添加删除提醒
func (s *TelegramBotServiceImpl) sendReplyWithResourceAutoDelete(message *tgbotapi.Message, text string, resourceCount int) {
	// 如果启用了自动删除且有资源，在消息中添加删除提醒
	if s.config.AutoDeleteEnabled && s.config.AutoDeleteInterval > 0 && resourceCount > 0 {
		deleteNotice := fmt.Sprintf("\n\n⏰ *此消息将在 %d 分钟后自动删除*", s.config.AutoDeleteInterval)
		text += deleteNotice
		utils.Info("[TELEGRAM:MESSAGE] 添加删除提醒到包含资源的回复消息")
	}

	// 使用资源消息的特殊删除逻辑
	s.sendReplyWithAutoDelete(message, text, s.config.AutoDeleteEnabled && resourceCount > 0)
}

// startContentPusher 启动内容推送器
func (s *TelegramBotServiceImpl) startContentPusher() {
	// 每分钟检查一次需要推送的频道
	s.cronScheduler.AddFunc("@every 1m", func() {
		s.pushContentToChannels()
	})

	s.cronScheduler.Start()
	utils.Info("[TELEGRAM:PUSH] 内容推送调度器已启动")
}

// pushContentToChannels 推送内容到频道
func (s *TelegramBotServiceImpl) pushContentToChannels() {
	// 获取需要推送的频道
	channels, err := s.channelRepo.FindDueForPush()
	if err != nil {
		utils.Error("[TELEGRAM:PUSH:ERROR] 获取推送频道失败: %v", err)
		return
	}

	if len(channels) == 0 {
		utils.Debug("[TELEGRAM:PUSH] 没有需要推送的频道")
		return
	}

	utils.Info("[TELEGRAM:PUSH] 开始推送内容到 %d 个频道", len(channels))

	for _, channel := range channels {
		go s.pushToChannel(channel)
	}
}

// pushToChannel 推送内容到一个频道
func (s *TelegramBotServiceImpl) pushToChannel(channel entity.TelegramChannel) {
	utils.Info("[TELEGRAM:PUSH] 开始推送到频道: %s (ID: %d)", channel.ChatName, channel.ChatID)

	// 1. 根据频道设置过滤资源
	resources := s.findResourcesForChannel(channel)
	if len(resources) == 0 {
		utils.Info("[TELEGRAM:PUSH] 频道 %s 没有可推送的内容", channel.ChatName)
		return
	}

	// 2. 构建推送消息
	message := s.buildPushMessage(channel, resources)

	// 3. 发送消息（推送消息不自动删除，使用 Markdown 格式）
	err := s.SendMessageWithFormat(channel.ChatID, message, "MarkdownV2")
	if err != nil {
		utils.Error("[TELEGRAM:PUSH:ERROR] 推送失败到频道 %s (%d): %v", channel.ChatName, channel.ChatID, err)
		return
	}

	// 4. 更新最后推送时间
	err = s.channelRepo.UpdateLastPushAt(channel.ID, time.Now())
	if err != nil {
		utils.Error("[TELEGRAM:PUSH:ERROR] 更新推送时间失败: %v", err)
		return
	}

	utils.Info("[TELEGRAM:PUSH:SUCCESS] 成功推送内容到频道: %s (%d 条资源)", channel.ChatName, len(resources))
}

// findResourcesForChannel 查找适合频道的资源
func (s *TelegramBotServiceImpl) findResourcesForChannel(channel entity.TelegramChannel) []interface{} {
	utils.Info("[TELEGRAM:PUSH] 开始为频道 %s (%d) 查找资源", channel.ChatName, channel.ChatID)

	params := map[string]interface{}{"category": "", "tag": ""}

	if channel.ContentCategories != "" {
		categories := strings.Split(channel.ContentCategories, ",")
		for i, category := range categories {
			categories[i] = strings.TrimSpace(category)
		}
		params["category"] = categories[0]
	}

	if channel.ContentTags != "" {
		tags := strings.Split(channel.ContentTags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		params["tag"] = tags[0]
	}

	// 尝试使用 PostgreSQL 的随机功能
	defer func() {
		if r := recover(); r != nil {
			utils.Warn("[TELEGRAM:PUSH] 随机查询失败，回退到传统方法: %v", r)
		}
	}()

	randomResource, err := s.resourceRepo.GetRandomResourceWithFilters(params["category"].(string), params["tag"].(string), channel.IsPushSavedInfo)
	if err == nil && randomResource != nil {
		utils.Info("[TELEGRAM:PUSH] 成功获取随机资源: %s", randomResource.Title)
		return []interface{}{randomResource}
	}

	return []interface{}{}
}

// buildPushMessage 构建推送消息
func (s *TelegramBotServiceImpl) buildPushMessage(channel entity.TelegramChannel, resources []interface{}) string {
	resource := resources[0].(*entity.Resource)

	message := fmt.Sprintf("🆕 %s\n\n", s.cleanResourceText(resource.Title))

	if resource.Description != "" {
		message += fmt.Sprintf("📝 %s\n\n", s.cleanResourceText(resource.Description))
	}

	// 添加标签
	if len(resource.Tags) > 0 {
		message += "\n🏷️ "
		for i, tag := range resource.Tags {
			if i > 0 {
				message += " "
			}
			message += fmt.Sprintf("#%s", tag.Name)
		}
		message += "\n"
	}

	// 添加资源信息
	message += fmt.Sprintf("\n💡 评论区评论 (【%s】%s) 即可获取资源，括号内名称点击可复制📋\n", resource.Key, resource.Title)

	return message
}

// GetBotUsername 获取机器人用户名
func (s *TelegramBotServiceImpl) GetBotUsername() string {
	if s.bot != nil {
		return s.bot.Self.UserName
	}
	return ""
}

// SendMessage 发送消息（默认使用 MarkdownV2 格式）
func (s *TelegramBotServiceImpl) SendMessage(chatID int64, text string) error {
	return s.SendMessageWithFormat(chatID, text, "MarkdownV2")
}

// SendMessageWithFormat 发送消息，支持指定格式
func (s *TelegramBotServiceImpl) SendMessageWithFormat(chatID int64, text string, parseMode string) error {
	if s.bot == nil {
		return fmt.Errorf("Bot 未初始化")
	}

	// 根据格式选择不同的文本清理方法
	var cleanedText string
	switch parseMode {
	case "Markdown", "MarkdownV2":
		cleanedText = s.cleanMessageText(text)
	case "HTML":
		cleanedText = s.cleanMessageTextForPlain(text) // HTML 格式暂时使用纯文本清理
	default: // 纯文本或其他格式
		cleanedText = s.cleanMessageTextForPlain(text)
		parseMode = "" // Telegram API 中空字符串表示纯文本
	}

	msg := tgbotapi.NewMessage(chatID, cleanedText)
	msg.ParseMode = parseMode

	// 检测并添加代码实体（只在 Markdown 格式下）
	if parseMode == "Markdown" || parseMode == "MarkdownV2" {
		entities := s.parseCodeEntities(text, cleanedText)
		if len(entities) > 0 {
			msg.Entities = entities
			utils.Info("[TELEGRAM:MESSAGE] 为消息添加了 %d 个代码实体", len(entities))
		}
	}

	msg1 := tgbotapi.NewMessage(chatID, "*bold text*\n"+
		"_italic \n"+
		"__underline__\n"+
		"~strikethrough~\n"+
		"||spoiler||\n"+
		"*bold _italic bold ~italic bold strikethrough ||italic bold strikethrough spoiler||~ __underline italic bold___ bold*\n"+
		"[inline URL](http://www.example.com/)\n"+
		"[inline mention of a user](tg://user?id=123456789)\n"+
		"![👍](tg://emoji?id=5368324170671202286)\n"+
		"`inline fixed-width code`\n"+
		"```\n"+
		"pre-formatted fixed-width code block\n"+
		"```\n"+
		"```python\n"+
		"pre-formatted fixed-width code block written in the Python programming language\n"+
		"```\n"+
		">Block quotation started\n"+
		">Block quotation continued\n"+
		">Block quotation continued\n"+
		">Block quotation continued\n"+
		">The last line of the block quotation\n"+
		"**>The expandable block quotation started right after the previous block quotation\n"+
		">It is separated from the previous block quotation by an empty bold entity\n"+
		">Expandable block quotation continued\n"+
		">Hidden by default part of the expandable block quotation started\n"+
		">Expandable block quotation continued\n"+
		">The last line of the expandable block quotation with the expandability mark||")
	s.bot.Send(msg1)

	_, err := s.bot.Send(msg)
	if err != nil {
		utils.Error("[TELEGRAM:MESSAGE:ERROR] 发送消息失败 (格式: %s): %v", parseMode, err)
		// 如果是格式错误，尝试发送纯文本版本
		if strings.Contains(err.Error(), "parse") || strings.Contains(err.Error(), "Bad Request") {
			utils.Info("[TELEGRAM:MESSAGE] 尝试发送纯文本版本...")
			msg.ParseMode = ""
			msg.Text = s.cleanMessageTextForPlain(text)
			msg.Entities = nil // 纯文本模式下不使用实体
			_, err = s.bot.Send(msg)
		}
	}
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
		PushFrequency:     5, // 默认5分钟
		IsActive:          true,
		RegisteredBy:      "bot_command",
		RegisteredAt:      time.Now(),
		ContentCategories: "",
		ContentTags:       "",
		API:               "",    // 后续可配置
		Token:             "",    // 后续可配置
		ApiType:           "l9",  // 默认l9类型
		IsPushSavedInfo:   false, // 默认推送所有资源
	}

	return s.channelRepo.Create(&channel)
}

// IsChannelRegistered 检查频道是否已注册
func (s *TelegramBotServiceImpl) IsChannelRegistered(chatID int64) bool {
	channel, err := s.channelRepo.FindByChatID(chatID)
	return err == nil && channel != nil
}

// isUserAdministrator 检查用户是否为群组管理员
func (s *TelegramBotServiceImpl) isUserAdministrator(chatID int64, userID int64) bool {
	if s.bot == nil {
		return false
	}

	// 获取用户在群组中的信息
	memberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	}

	member, err := s.bot.GetChatMember(memberConfig)
	if err != nil {
		utils.Error("[TELEGRAM:ADMIN] 获取用户群组成员信息失败: %v", err)
		return false
	}

	// 检查用户是否为管理员或创建者
	userStatus := string(member.Status)
	return userStatus == "administrator" || userStatus == "creator"
}

// isBotAdministrator 检查机器人是否为频道管理员
func (s *TelegramBotServiceImpl) isBotAdministrator(chatID int64) bool {
	if s.bot == nil {
		return false
	}

	// 获取机器人自己的信息
	botInfo, err := s.bot.GetMe()
	if err != nil {
		utils.Error("[TELEGRAM:ADMIN:BOT] 获取机器人信息失败: %v", err)
		return false
	}

	// 获取机器人作为频道成员的信息
	memberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: botInfo.ID,
		},
	}

	member, err := s.bot.GetChatMember(memberConfig)
	if err != nil {
		utils.Error("[TELEGRAM:ADMIN:BOT] 获取机器人频道成员信息失败: %v", err)
		return false
	}

	// 检查机器人是否为管理员或创建者
	botStatus := string(member.Status)
	utils.Info("[TELEGRAM:ADMIN:BOT] 机器人状态: %s (ChatID: %d)", botStatus, chatID)
	return botStatus == "administrator" || botStatus == "creator"
}

// hasActiveGroup 检查是否已经注册了活跃的群组
func (s *TelegramBotServiceImpl) hasActiveGroup() bool {
	channels, err := s.channelRepo.FindByChatType("group")
	if err != nil {
		utils.Error("[TELEGRAM:LIMIT] 检查活跃群组失败: %v", err)
		return false
	}

	// 检查是否有活跃的群组
	for _, channel := range channels {
		if channel.IsActive {
			return true
		}
	}
	return false
}

// hasActiveChannel 检查是否已经注册了活跃的频道
func (s *TelegramBotServiceImpl) hasActiveChannel() bool {
	channels, err := s.channelRepo.FindByChatType("channel")
	if err != nil {
		utils.Error("[TELEGRAM:LIMIT] 检查活跃频道失败: %v", err)
		return false
	}

	// 检查是否有活跃的频道
	for _, channel := range channels {
		if channel.IsActive {
			return true
		}
	}
	return false
}

// handleChannelRegistration 处理频道注册（支持频道ID和用户名）
func (s *TelegramBotServiceImpl) handleChannelRegistration(message *tgbotapi.Message, channelParam string) {
	channelParam = strings.TrimSpace(channelParam)

	var chat tgbotapi.Chat
	var err error
	var identifier string

	// 首先获取频道信息，然后检查机器人权限
	// 这一步会在后面的逻辑中完成，获取chat对象后再检查权限

	// 判断是频道ID还是用户名格式
	if strings.HasPrefix(channelParam, "@") {
		// 用户名格式：@username
		username := strings.TrimPrefix(channelParam, "@")
		if username == "" {
			errorMsg := "❌ *用户名格式错误*\n\n用户名不能为空，如 @mychannel"
			s.sendReply(message, errorMsg)
			return
		}

		// 尝试通过用户名获取频道信息
		// 手动构造请求URL并发送
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getChat", s.config.ApiKey)
		data := url.Values{}
		data.Set("chat_id", "@"+username)

		client := &http.Client{Timeout: 10 * time.Second}

		// 如果有代理，配置代理
		if s.config.ProxyEnabled && s.config.ProxyHost != "" {
			var proxyClient *http.Client
			if s.config.ProxyType == "socks5" {
				// SOCKS5代理配置
				auth := &proxy.Auth{}
				if s.config.ProxyUsername != "" {
					auth.User = s.config.ProxyUsername
					auth.Password = s.config.ProxyPassword
				}
				dialer, proxyErr := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", s.config.ProxyHost, s.config.ProxyPort), auth, proxy.Direct)
				if proxyErr != nil {
					errorMsg := fmt.Sprintf("❌ *代理配置错误*\n\n无法连接到代理服务器: %v", proxyErr)
					s.sendReply(message, errorMsg)
					return
				}
				proxyClient = &http.Client{
					Transport: &http.Transport{
						Dial: dialer.Dial,
					},
					Timeout: 10 * time.Second,
				}
			} else {
				// HTTP/HTTPS代理配置
				proxyURL := &url.URL{
					Scheme: s.config.ProxyType,
					Host:   fmt.Sprintf("%s:%d", s.config.ProxyHost, s.config.ProxyPort),
				}
				if s.config.ProxyUsername != "" {
					proxyURL.User = url.UserPassword(s.config.ProxyUsername, s.config.ProxyPassword)
				}
				proxyClient = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyURL),
					},
					Timeout: 10 * time.Second,
				}
			}
			client = proxyClient
		}

		resp, httpErr := client.PostForm(apiURL, data)
		if httpErr != nil {
			errorMsg := fmt.Sprintf("❌ *无法访问频道*\n\n请确保:\n• 机器人已被添加到频道 @%s\n• 机器人已被设为频道管理员\n• 用户名正确\n\n错误详情: %v", username, httpErr)
			s.sendReply(message, errorMsg)
			return
		}
		defer resp.Body.Close()

		// 解析响应
		var apiResponse struct {
			OK     bool `json:"ok"`
			Result struct {
				ID       int64  `json:"id"`
				Title    string `json:"title"`
				Username string `json:"username"`
				Type     string `json:"type"`
			} `json:"result"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			errorMsg := "❌ *解析服务器响应失败*\n\n请稍后重试"
			s.sendReply(message, errorMsg)
			return
		}

		if !apiResponse.OK {
			errorMsg := fmt.Sprintf("❌ *获取频道信息失败*\n\n错误: %s", apiResponse.Description)
			s.sendReply(message, errorMsg)
			return
		}

		// 检查是否是频道
		if apiResponse.Result.Type != "channel" {
			errorMsg := "❌ *这不是一个频道*\n\n请提供有效的频道用户名。"
			s.sendReply(message, errorMsg)
			return
		}

		// 构造Chat对象
		chat = tgbotapi.Chat{
			ID:       apiResponse.Result.ID,
			Title:    apiResponse.Result.Title,
			UserName: apiResponse.Result.Username,
			Type:     apiResponse.Result.Type,
		}

		identifier = fmt.Sprintf("@%s", username)

		// 检查机器人是否是频道管理员
		if !s.isBotAdministrator(chat.ID) {
			errorMsg := "❌ *权限不足*\n\n机器人必须是频道的管理员才能注册此频道用于推送。\n\n请先将机器人添加为频道管理员，然后重试注册命令。"
			s.sendReply(message, errorMsg)
			return
		}

	} else if strings.HasPrefix(channelParam, "-") && len(channelParam) > 10 {
		// 频道ID格式：-1001234567890
		channelID, parseErr := strconv.ParseInt(channelParam, 10, 64)
		if parseErr != nil {
			errorMsg := fmt.Sprintf("❌ *频道ID格式错误*\n\n频道ID必须是数字，如 -1001234567890\n\n您输入的: %s", channelParam)
			s.sendReply(message, errorMsg)
			return
		}

		// 通过频道ID获取频道信息
		chat, err = s.bot.GetChat(tgbotapi.ChatInfoConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: channelID,
			},
		})

		if err != nil {
			errorMsg := fmt.Sprintf("❌ *无法访问频道*\n\n请确保:\n• 机器人已被添加到频道\n• 机器人已被设为频道管理员\n• 频道ID正确\n\n错误详情: %v", err)
			s.sendReply(message, errorMsg)
			return
		}

		// 检查是否已经是频道
		if !chat.IsChannel() {
			errorMsg := "❌ *这不是一个频道*\n\n请提供有效的频道ID。"
			s.sendReply(message, errorMsg)
			return
		}

		// 检查机器人是否是频道管理员
		if !s.isBotAdministrator(chat.ID) {
			errorMsg := "❌ *权限不足*\n\n机器人必须是频道的管理员才能注册此频道用于推送。\n\n请先将机器人添加为频道管理员，然后重试注册命令。"
			s.sendReply(message, errorMsg)
			return
		}

		// 检查是否已经注册了频道
		if s.hasActiveChannel() {
			errorMsg := "❌ *注册限制*\n\n系统最多只支持注册一个频道用于推送。\n\n请先注销现有频道，然后再注册新的频道。"
			s.sendReply(message, errorMsg)
			return
		}

		identifier = fmt.Sprintf("ID: %d", chat.ID)

	} else {
		// 无效格式
		errorMsg := fmt.Sprintf("❌ *格式错误*\n\n支持的格式:\n• 频道ID: -1001234567890\n• 用户名: @mychannel\n\n您输入的: %s", channelParam)
		s.sendReply(message, errorMsg)
		return
	}

	// 尝试查找现有频道
	existingChannel, findErr := s.channelRepo.FindByChatID(chat.ID)

	if findErr == nil && existingChannel != nil {
		// 频道已存在，更新信息
		existingChannel.ChatName = chat.Title
		existingChannel.RegisteredBy = message.From.UserName
		existingChannel.RegisteredAt = time.Now()
		existingChannel.IsActive = true
		existingChannel.PushEnabled = true
		// 为现有频道设置默认值
		if existingChannel.ApiType == "" {
			existingChannel.ApiType = "telegram"
		}

		err := s.channelRepo.Update(existingChannel)
		if err != nil {
			errorMsg := fmt.Sprintf("❌ 频道更新失败: %v", err)
			s.sendReply(message, errorMsg)
			return
		}

		successMsg := fmt.Sprintf("✅ *频道更新成功！*\n\n频道: %s\n%s\n类型: 频道\n\n频道信息已更新，现在可以正常推送内容。", chat.Title, identifier)
		s.sendReply(message, successMsg)
		return
	}

	// 频道不存在，创建新记录
	channel := entity.TelegramChannel{
		ChatID:            chat.ID,
		ChatName:          chat.Title,
		ChatType:          "channel",
		PushEnabled:       true,
		PushFrequency:     60, // 默认1小时
		IsActive:          true,
		RegisteredBy:      message.From.UserName,
		RegisteredAt:      time.Now(),
		ContentCategories: "",
		ContentTags:       "",
		API:               "",         // 后续可配置
		Token:             "",         // 后续可配置
		ApiType:           "telegram", // 默认telegram类型
		IsPushSavedInfo:   false,      // 默认推送所有资源
	}

	createErr := s.channelRepo.Create(&channel)
	if createErr != nil {
		// 如果创建失败，可能是因为并发或其他问题，再次尝试查找
		if existing, retryErr := s.channelRepo.FindByChatID(chat.ID); retryErr == nil && existing != nil {
			successMsg := fmt.Sprintf("⚠️ *频道已注册*\n\n频道: %s\n%s\n类型: 频道\n\n此频道已经注册，无需重复注册。", chat.Title, identifier)
			s.sendReply(message, successMsg)
		} else {
			errorMsg := fmt.Sprintf("❌ 频道注册失败: %v", createErr)
			s.sendReply(message, errorMsg)
		}
		return
	}

	successMsg := fmt.Sprintf("✅ *频道注册成功！*\n\n频道: %s\n%s\n类型: 频道\n\n现在可以向此频道推送资源内容了。\n\n可以通过管理界面调整推送设置。", chat.Title, identifier)
	s.sendReply(message, successMsg)
}

// HandleWebhookUpdate 处理 Webhook 更新（预留接口，目前使用长轮询）
func (s *TelegramBotServiceImpl) HandleWebhookUpdate(c interface{}) {
	// 目前使用长轮询模式，webhook 接口预留
	// 将来可以实现从 webhook 接收消息的处理逻辑
	// 如果需要实现 webhook 模式，可以在这里添加处理逻辑
}

// CleanupDuplicateChannels 清理数据库中的重复频道记录
func (s *TelegramBotServiceImpl) CleanupDuplicateChannels() error {
	utils.Info("[TELEGRAM:CLEANUP] 开始清理重复的频道记录...")

	err := s.channelRepo.CleanupDuplicateChannels()
	if err != nil {
		utils.Error("[TELEGRAM:CLEANUP:ERROR] 清理重复频道记录失败: %v", err)
		return fmt.Errorf("清理重复频道记录失败: %v", err)
	}

	utils.Info("[TELEGRAM:CLEANUP:SUCCESS] 成功清理重复的频道记录")
	return nil
}

// parseCodeEntities 解析消息中的代码实体
func (s *TelegramBotServiceImpl) parseCodeEntities(originalText string, cleanedText string) []tgbotapi.MessageEntity {
	var entities []tgbotapi.MessageEntity

	// 定义开始和结束标记
	startMarker := "评论区评论 ("
	endMarker := ") 即可获取资源"

	// 在原始文本中查找标记
	start := strings.Index(originalText, startMarker)
	if start == -1 {
		return entities
	}

	// 计算代码块的开始位置（在开始标记之后）
	codeStart := start + len(startMarker)

	// 查找结束标记
	end := strings.Index(originalText[codeStart:], endMarker)
	if end == -1 {
		return entities
	}

	// 计算代码块的结束位置
	codeEnd := codeStart + end

	// 确保代码内容不为空
	if codeEnd <= codeStart {
		return entities
	}

	// 获取原始代码内容
	originalCodeContent := originalText[codeStart:codeEnd]

	// 在清理后的文本中查找相同的代码内容，计算新的偏移量
	cleanedStart := strings.Index(cleanedText, originalCodeContent)
	if cleanedStart == -1 {
		// 如果找不到完全匹配的内容，使用精确偏移计算
		cleanedStart = s.findPreciseOffset(originalText, cleanedText, codeStart)
	}

	// 验证清理后偏移量是否有效
	if cleanedStart < 0 || cleanedStart >= len(cleanedText) {
		utils.Warn("[TELEGRAM:MESSAGE] 无法计算有效的实体偏移量")
		return entities
	}

	// 安全地获取清理后的代码内容（确保不超出字符串边界）
	cleanedEnd := cleanedStart + len(originalCodeContent)
	if cleanedEnd > len(cleanedText) {
		cleanedEnd = len(cleanedText)
	}
	cleanedCodeContent := cleanedText[cleanedStart:cleanedEnd]

	// 确保清理后的代码内容不为空
	if strings.TrimSpace(cleanedCodeContent) == "" {
		return entities
	}

	// 创建代码实体，使用 UTF-8 字符计数
	codeEntity := tgbotapi.MessageEntity{
		Type:   "code",
		Offset: utf8.RuneCountInString(cleanedText[:cleanedStart]), // 使用 UTF-8 字符计数
		Length: utf8.RuneCountInString(cleanedCodeContent),         // 使用 UTF-8 字符计数
	}

	entities = append(entities, codeEntity)

	utils.Info("[TELEGRAM:MESSAGE] 检测到代码实体: 原始位置=%d-%d, 清理后位置=%d-%d",
		codeStart, codeEnd, cleanedStart, cleanedEnd)
	utils.Info("[TELEGRAM:MESSAGE] 原始代码内容: %s", originalCodeContent)
	utils.Info("[TELEGRAM:MESSAGE] 清理后代码内容: %s", cleanedCodeContent)
	utils.Info("[TELEGRAM:MESSAGE] 实体偏移量: %d, 长度: %d", codeEntity.Offset, codeEntity.Length)

	return entities
}

// findPreciseOffset 通过字符级别的精确匹配计算清理后文本中的偏移量
func (s *TelegramBotServiceImpl) findPreciseOffset(originalText string, cleanedText string, originalOffset int) int {
	// 获取原始文本中指定位置前后的上下文
	contextSize := 50
	originalContext := originalText[max(0, originalOffset-contextSize):min(len(originalText), originalOffset+contextSize)]

	// 在清理后的文本中查找相似的上下文
	bestMatch := -1
	maxSimilarity := 0.0

	for i := 0; i <= len(cleanedText)-len(originalContext); i++ {
		candidate := cleanedText[i:min(len(cleanedText), i+len(originalContext))]
		similarity := s.calculateSimilarity(originalContext, candidate)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			bestMatch = i + (originalOffset - max(0, originalOffset-contextSize))
		}
	}

	// 如果相似度足够高，返回最佳匹配
	if maxSimilarity > 0.7 {
		return max(0, min(len(cleanedText)-1, bestMatch))
	}

	// 回退到比例估算
	return s.calculateCleanedOffset(originalText, cleanedText, originalOffset)
}

// calculateSimilarity 计算两个字符串的相似度
func (s *TelegramBotServiceImpl) calculateSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	// 简单字符匹配相似度
	matches := 0
	minLen := min(len(s1), len(s2))

	for i := 0; i < minLen; i++ {
		if s1[i] == s2[i] {
			matches++
		}
	}

	return float64(matches) / float64(minLen)
}

// calculateCleanedOffset 计算清理后文本中的偏移量（比例估算）
func (s *TelegramBotServiceImpl) calculateCleanedOffset(originalText string, cleanedText string, originalOffset int) int {
	// 计算清理后文本中对应位置的近似偏移量
	// 这种方法通过比较字符比例来估算位置
	if len(originalText) == 0 {
		return 0
	}

	originalRatio := float64(originalOffset) / float64(len(originalText))
	estimatedOffset := int(float64(len(cleanedText)) * originalRatio)

	// 确保偏移量在有效范围内
	if estimatedOffset < 0 {
		estimatedOffset = 0
	}
	if estimatedOffset >= len(cleanedText) {
		estimatedOffset = len(cleanedText) - 1
	}

	return estimatedOffset
}

// 辅助函数：返回两个数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
