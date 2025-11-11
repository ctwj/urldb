package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
	"golang.org/x/net/proxy"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

// https://core.telegram.org/bots/api

type TelegramBotService interface {
	Start() error
	Stop() error
	IsRunning() bool
	ReloadConfig() error
	GetRuntimeStatus() map[string]interface{}
	ValidateApiKey(apiKey string) (bool, map[string]interface{}, error)
	ValidateApiKeyWithProxy(apiKey string, proxyEnabled bool, proxyType, proxyHost string, proxyPort int, proxyUsername, proxyPassword string) (bool, map[string]interface{}, error)
	GetBotUsername() string
	SendMessage(chatID int64, text string, img string) error
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
	readyRepo        repo.ReadyResourceRepository
	cronScheduler    *cron.Cron
	config           *TelegramBotConfig
	pushHistory      map[int64][]uint // 每个频道的推送历史记录，最多100条
	mu               sync.RWMutex     // 用于保护pushHistory的读写锁
	stopChan         chan struct{}    // 用于停止消息循环的channel
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
	readyResourceRepo repo.ReadyResourceRepository,
) TelegramBotService {
	return &TelegramBotServiceImpl{
		isRunning:        false,
		systemConfigRepo: systemConfigRepo,
		channelRepo:      channelRepo,
		resourceRepo:     resourceRepo,
		readyRepo:        readyResourceRepo,
		cronScheduler:    cron.New(),
		config:           &TelegramBotConfig{},
		pushHistory:      make(map[int64][]uint),
		stopChan:         make(chan struct{}),
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

	// 统计配置项数量，用于汇总日志
	configCount := 0

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyTelegramBotEnabled:
			s.config.Enabled = config.Value == "true"
		case entity.ConfigKeyTelegramBotApiKey:
			s.config.ApiKey = config.Value
		case entity.ConfigKeyTelegramAutoReplyEnabled:
			s.config.AutoReplyEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramAutoReplyTemplate:
			if config.Value != "" {
				s.config.AutoReplyTemplate = config.Value
			}
		case entity.ConfigKeyTelegramAutoDeleteEnabled:
			s.config.AutoDeleteEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramAutoDeleteInterval:
			if config.Value != "" {
				fmt.Sscanf(config.Value, "%d", &s.config.AutoDeleteInterval)
			}
		case entity.ConfigKeyTelegramProxyEnabled:
			s.config.ProxyEnabled = config.Value == "true"
		case entity.ConfigKeyTelegramProxyType:
			s.config.ProxyType = config.Value
		case entity.ConfigKeyTelegramProxyHost:
			s.config.ProxyHost = config.Value
		case entity.ConfigKeyTelegramProxyPort:
			if config.Value != "" {
				fmt.Sscanf(config.Value, "%d", &s.config.ProxyPort)
			}
		case entity.ConfigKeyTelegramProxyUsername:
			s.config.ProxyUsername = config.Value
		case entity.ConfigKeyTelegramProxyPassword:
			s.config.ProxyPassword = config.Value
		default:
			utils.Debug("未知Telegram配置: %s", config.Key)
		}
		configCount++
	}

	// 汇总输出配置加载结果，避免逐项日志
	proxyStatus := "禁用"
	if s.config.ProxyEnabled {
		proxyStatus = "启用"
	}

	utils.TelegramInfo("配置加载完成 - Bot启用: %v, 自动回复: %v, 代理: %s, 配置项数: %d",
		s.config.Enabled, s.config.AutoReplyEnabled, proxyStatus, configCount)
	return nil
}

// Start 启动机器人服务
func (s *TelegramBotServiceImpl) Start() error {
	// 确保机器人完全停止状态
	if s.isRunning && s.bot != nil {
		utils.Info("[TELEGRAM:SERVICE] Telegram Bot 服务已经在运行中")
		return nil
	}

	// 如果isRunning为true但bot为nil，说明状态不一致，需要清理
	if s.isRunning && s.bot == nil {
		utils.Info("[TELEGRAM:SERVICE] 检测到不一致状态，清理残留资源")
		s.isRunning = false
	}

	// 加载配置
	if err := s.loadConfig(); err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 加载推送历史记录
	if err := s.loadPushHistory(); err != nil {
		utils.Error("[TELEGRAM:SERVICE] 加载推送历史记录失败: %v", err)
		// 不返回错误，继续启动服务
	}

	if !s.config.Enabled || s.config.ApiKey == "" {
		utils.Info("[TELEGRAM:SERVICE] Telegram Bot 未启用或 API Key 未配置")
		// 如果机器人当前正在运行，需要停止它
		if s.isRunning {
			utils.Info("[TELEGRAM:SERVICE] 机器人已被禁用，停止正在运行的服务")
			s.Stop()
		}
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

	// 重置停止信号channel
	s.stopChan = make(chan struct{})

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

	utils.Info("[TELEGRAM:SERVICE] 开始停止 Telegram Bot 服务")

	s.isRunning = false

	// 安全地发送停止信号给消息循环
	select {
	case <-s.stopChan:
		// channel 已经关闭
	default:
		// channel 未关闭，安全关闭
		close(s.stopChan)
	}

	if s.cronScheduler != nil {
		s.cronScheduler.Stop()
	}

	// 清理机器人实例以避免冲突
	s.bot = nil

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
			utils.Error(fmt.Sprintf("[TELEGRAM:VALIDATE] 创建 Telegram Bot (代理校验) 失败: %v", err))
			return false, nil, fmt.Errorf("创建 Telegram Bot (代理校验) 失败: %v", err)
		}

		utils.Info("[TELEGRAM:VALIDATE] 使用代理配置校验 API Key")
	} else {
		// 直连校验
		bot, err = tgbotapi.NewBotAPI(apiKey)
		if err != nil {
			utils.Error(fmt.Sprintf("[TELEGRAM:VALIDATE] 创建 Telegram Bot 失败: %v", err))
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

	// 确保机器人实例存在
	if s.bot == nil {
		utils.Error("[TELEGRAM:MESSAGE] 机器人实例为空，无法启动消息监听循环")
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	utils.Info("[TELEGRAM:MESSAGE] 消息监听循环已启动，等待消息...")

	for {
		select {
		case <-s.stopChan:
			utils.Info("[TELEGRAM:MESSAGE] 收到停止信号，退出消息监听循环")
			return
		case update, ok := <-updates:
			if !ok {
				utils.Info("[TELEGRAM:MESSAGE] updates channel 已关闭，退出消息监听循环")
				return
			}
			// 在处理消息前检查机器人是否仍在运行
			if !s.isRunning || s.bot == nil {
				utils.Info("[TELEGRAM:MESSAGE] 机器人已停止，忽略接收到的消息")
				return
			}
			if update.Message != nil {
				utils.Info("[TELEGRAM:MESSAGE] 接收到新消息更新")
				s.handleMessage(update.Message)
			} else {
				utils.Debug("[TELEGRAM:MESSAGE] 接收到其他类型更新: %v", update)
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (s *TelegramBotServiceImpl) handleMessage(message *tgbotapi.Message) {
	// 检查机器人是否正在运行且已启用
	if !s.isRunning || !s.config.Enabled {
		utils.Info("[TELEGRAM:MESSAGE] 机器人已停止或禁用，跳过消息处理: ChatID=%d", message.Chat.ID)
		return
	}

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

	// 处理 /s 命令
	if strings.HasPrefix(strings.ToLower(text), "/s ") {
		utils.Info("[TELEGRAM:MESSAGE] 处理 /s 命令 from ChatID=%d", chatID)
		// 提取搜索关键词
		keyword := strings.TrimSpace(text[3:]) // 去掉 "/s " 前缀
		if keyword != "" {
			utils.Info("[TELEGRAM:MESSAGE] 处理搜索请求 from ChatID=%d: %s", chatID, keyword)
			s.handleSearchRequest(message, keyword)
			return
		}
	}

	if len(text) == 0 {
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
				// 处理普通文本消息（搜索请求）
				re := regexp.MustCompile(`^【(\d+)】.*?`)
				matches := re.FindStringSubmatch(text)
				if len(matches) >= 2 {
					utils.Info("[TELEGRAM:MESSAGE] 处理搜索请求 from ChatID=%d: %s", chatID, text)
					num, _ := strconv.Atoi(matches[1])
					s.handleResourceRequest(message, uint(num))
					return
				}
				sre := regexp.MustCompile(`^搜索(.*?)$`)
				smatches := sre.FindStringSubmatch(text)
				if len(smatches) >= 2 {
					utils.Info("[TELEGRAM:MESSAGE] 处理搜索请求 from ChatID=%d: %s", chatID, text)
					keyword := strings.TrimSpace(smatches[1])
					s.handleSearchRequest(message, keyword)
					return
				}
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

		// 检查当前活跃的Telegram项目总数（频道+群组）
		activeItemCount := s.hasActiveTelegramItems()
		if activeItemCount >= 3 {
			errorMsg := "❌ *注册限制*\n\n系统最多支持注册3个频道/群组用于推送。\n\n当前已注册: %d个，请先注销现有频道/群组，然后再注册新的。"
			s.sendReply(message, fmt.Sprintf(errorMsg, activeItemCount))
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
• 发送 /s 关键词 进行资源搜索（命令形式）
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
func (s *TelegramBotServiceImpl) handleResourceRequest(message *tgbotapi.Message, id uint) {

	// 使用资源仓库进行搜索
	resources, err := s.resourceRepo.FindByIDs([]uint{uint(id)}) // 限制为5个结果
	if err != nil {
		utils.Error("[TELEGRAM:SEARCH] 搜索失败: %v", err)
		s.sendReply(message, "搜索服务暂时不可用，请稍后重试")
		return
	}

	if len(resources) == 0 {
		s.sendReply(message, "未找到该资源")
		return
	}

	// 构建搜索结果消息
	resultText := ""

	// 显示前5个结果
	for i, resource := range resources {
		if i >= 1 {
			break
		}
		title := s.cleanMessageTextForHTML(resource.Title)

		if resource.SaveURL != "" {
			resultText += fmt.Sprintf("<b>%d. %s</>\n<i>%s</i>\n", i+1, title, resource.SaveURL)
		} else {
			resultText += fmt.Sprintf("<b>%d. %s</>\n<i>%s</i>\n", i+1, title, resource.URL)
		}
	}

	// 使用包含资源的自动删除功能
	s.sendReplyWithResourceAutoDelete(message, resultText)
	s.sendReply(message, "资源已发送，请注意查收")
}

// handleSearchRequest 处理搜索请求
func (s *TelegramBotServiceImpl) handleSearchRequest(message *tgbotapi.Message, keyword string) {

	// 使用资源仓库进行搜索
	resources, total, err := s.resourceRepo.Search(keyword, nil, 1, 5) // 限制为5个结果
	if err != nil {
		utils.Error("[TELEGRAM:SEARCH] 搜索失败: %v", err)
		s.sendReply(message, "搜索服务暂时不可用，请稍后重试")
		return
	}

	if total == 0 {
		response := fmt.Sprintf("🔍 *搜索结果*\n\n关键词: `%s`\n\n❌ 未找到相关资源\n\n💡 建议:\n• 尝试使用更通用的关键词\n• 检查拼写是否正确\n• 减少关键词数量", keyword)
		// 没有找到资源，不使用资源自动删除
		s.sendReply(message, response)
		return
	}

	// 构建搜索结果消息
	resultText := fmt.Sprintf("🔍 *搜索结果* 总共找到: %d 个资源\n\n", total)

	// 显示前5个结果
	for i, resource := range resources {
		if i >= 5 {
			break
		}
		title := s.cleanMessageTextForHTML(resource.Title)
		// description := s.cleanMessageTextForHTML(resource.Description)
		if resource.SaveURL != "" {
			resultText += fmt.Sprintf("<b>%s</b>\n<a href=\"%s\">%s</a>\n", title, resource.SaveURL, resource.SaveURL)
		} else {
			resultText += fmt.Sprintf("<b>%s</b>\n<a href=\"%s\">%s</a>\n", title, resource.URL, resource.URL)
		}
	}

	// 如果有更多结果，添加提示
	if total > 5 {
		resultText += fmt.Sprintf("... 还有 %d 个结果\n\n", total-5)
	}

	resultText += "<i>如果资源失效请访问，发送搜索 + 关键字，可以搜索资源</i>"

	// 使用包含资源的自动删除功能
	s.sendReplyWithResourceAutoDelete(message, resultText)
}

// sendReply 发送回复消息
func (s *TelegramBotServiceImpl) sendReply(message *tgbotapi.Message, text string) {
	s.sendReplyWithAutoDelete(message, text, false)
}

// sendReplyWithAutoDelete 发送回复消息，支持指定是否自动删除
func (s *TelegramBotServiceImpl) sendReplyWithAutoDelete(message *tgbotapi.Message, text string, autoDelete bool) {
	// 清理消息文本，确保UTF-8编码
	originalText := text
	utils.Info("[TELEGRAM:MESSAGE] 尝试发送回复消息到 ChatID=%d, 原始长度=%d, 清理后长度=%d", message.Chat.ID, len(originalText), len(text))

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.MessageID

	utils.Debug("[TELEGRAM:MESSAGE] 发送Markdown版本消息: %s", text[:min(100, len(text))])

	sentMsg, err := s.bot.Send(msg)
	if err != nil {
		utils.Error("[TELEGRAM:MESSAGE:ERROR] 发送Markdown消息失败: %v", err)
		return
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

// cleanMessageTextForHTML 清理消息文本为HTML格式
func (s *TelegramBotServiceImpl) cleanMessageTextForHTML(text string) string {
	if text == "" {
		return text
	}
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

// sendReplyWithResourceAutoDelete 发送包含资源的回复消息，自动添加删除提醒
func (s *TelegramBotServiceImpl) sendReplyWithResourceAutoDelete(message *tgbotapi.Message, text string) {
	// 如果启用了自动删除且有资源，在消息中添加删除提醒
	if s.config.AutoDeleteEnabled && s.config.AutoDeleteInterval > 0 {
		deleteNotice := fmt.Sprintf("\n\n⏰ <b>此消息将在 %d 分钟后自动删除</b>", s.config.AutoDeleteInterval)
		text += deleteNotice
	}

	// 使用资源消息的特殊删除逻辑
	s.sendReplyWithAutoDelete(message, text, true)
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

	// 过滤出在允许推送时间段内的频道
	validChannels := s.filterChannelsByTimeRange(channels)
	if len(validChannels) == 0 {
		utils.Info("[TELEGRAM:PUSH] 所有频道都不在推送时间段内")
		return
	}

	utils.Info("[TELEGRAM:PUSH] 开始推送内容到 %d 个频道（过滤前: %d 个频道）", len(validChannels), len(channels))

	for _, channel := range validChannels {
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
	message, img := s.buildPushMessage(channel, resources)

	// 3. 发送消息（推送消息不自动删除，使用 HTML 格式）
	err := s.SendMessage(channel.ChatID, message, img)
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

	// 5. 记录推送的资源ID到历史记录，避免重复推送
	for _, resource := range resources {
		var resourceID uint
		switch r := resource.(type) {
		case *entity.Resource:
			resourceID = r.ID
		case entity.Resource:
			resourceID = r.ID
		default:
			utils.Error("[TELEGRAM:PUSH] 无效的资源类型: %T", resource)
			continue
		}
		s.addPushedResourceID(channel.ChatID, resourceID)
	}

	utils.Info("[TELEGRAM:PUSH:SUCCESS] 成功推送内容到频道: %s (%d 条资源)", channel.ChatName, len(resources))
}

// findResourcesForChannel 查找适合频道的资源
func (s *TelegramBotServiceImpl) findResourcesForChannel(channel entity.TelegramChannel) []interface{} {
	utils.Info("[TELEGRAM:PUSH] 开始为频道 %s (%d) 查找资源", channel.ChatName, channel.ChatID)

	// 获取最近推送的历史资源ID，避免重复推送
	excludeResourceIDs := s.getRecentlyPushedResourceIDs(channel.ChatID)

	// 解析资源策略
	strategy := channel.ResourceStrategy
	if strategy == "" {
		strategy = "random" // 默认纯随机
	}

	utils.Info("[TELEGRAM:PUSH] 使用策略: %s, 时间限制: %s, 排除最近推送资源数: %d",
		strategy, channel.TimeLimit, len(excludeResourceIDs))

	// 根据策略获取资源
	switch strategy {
	case "latest":
		// 最新优先策略 - 获取最近的资源
		return s.findLatestResources(channel, excludeResourceIDs)
	case "transferred":
		// 已转存优先策略 - 优先获取有转存链接的资源
		return s.findTransferredResources(channel, excludeResourceIDs)
	case "random":
		// 纯随机策略（原逻辑）
		return s.findRandomResources(channel, excludeResourceIDs)
	default:
		// 默认随机策略
		return s.findRandomResources(channel, excludeResourceIDs)
	}
}

// findLatestResources 查找最新资源
func (s *TelegramBotServiceImpl) findLatestResources(channel entity.TelegramChannel, excludeResourceIDs []uint) []interface{} {
	params := s.buildFilterParams(channel)

	// 在数据库查询中排除已推送的资源
	if len(excludeResourceIDs) > 0 {
		params["exclude_ids"] = excludeResourceIDs
	}

	// 使用现有的搜索功能，按更新时间倒序获取最新资源
	resources, _, err := s.resourceRepo.SearchWithFilters(params)
	if err != nil {
		utils.Error("[TELEGRAM:PUSH] 获取最新资源失败: %v", err)
		return s.findRandomResources(channel, excludeResourceIDs) // 回退到随机策略
	}

	// 应用时间限制
	if channel.TimeLimit != "none" && len(resources) > 0 {
		resources = s.applyTimeFilter(resources, channel.TimeLimit)
	}

	if len(resources) == 0 {
		utils.Info("[TELEGRAM:PUSH] 没有找到符合条件的最新资源，尝试获取随机资源")
		return s.findRandomResources(channel, excludeResourceIDs) // 回退到随机策略
	}

	// 返回最新资源（第一条）
	utils.Info("[TELEGRAM:PUSH] 成功获取最新资源: %s", resources[0].Title)
	return []interface{}{&resources[0]}
}

// findTransferredResources 查找已转存资源
func (s *TelegramBotServiceImpl) findTransferredResources(channel entity.TelegramChannel, excludeResourceIDs []uint) []interface{} {
	params := s.buildFilterParams(channel)

	// 添加转存链接条件
	params["has_save_url"] = true

	// 在数据库查询中排除已推送的资源
	if len(excludeResourceIDs) > 0 {
		params["exclude_ids"] = excludeResourceIDs
	}

	// 优先获取有转存链接的资源
	resources, _, err := s.resourceRepo.SearchWithFilters(params)
	if err != nil {
		utils.Error("[TELEGRAM:PUSH] 获取已转存资源失败: %v", err)
		return []interface{}{}
	}

	// 应用时间限制
	if channel.TimeLimit != "none" && len(resources) > 0 {
		resources = s.applyTimeFilter(resources, channel.TimeLimit)
	}

	if len(resources) == 0 {
		utils.Info("[TELEGRAM:PUSH] 没有找到符合条件的已转存资源，尝试获取随机资源")
		// 如果没有已转存资源，回退到随机策略
		return s.findRandomResources(channel, excludeResourceIDs)
	}

	// 返回第一个有转存链接的资源
	utils.Info("[TELEGRAM:PUSH] 成功获取已转存资源: %s", resources[0].Title)
	return []interface{}{&resources[0]}
}

// findRandomResources 查找随机资源（原有逻辑）
func (s *TelegramBotServiceImpl) findRandomResources(channel entity.TelegramChannel, excludeResourceIDs []uint) []interface{} {
	params := s.buildFilterParams(channel)

	// 如果是已转存优先策略但没有找到转存资源，这里会回退到随机策略
	// 此时不需要额外的转存链接条件，让随机函数处理

	// 在数据库查询中排除已推送的资源
	if len(excludeResourceIDs) > 0 {
		params["exclude_ids"] = excludeResourceIDs
	}

	// 使用搜索功能获取候选资源，然后过滤
	params["limit"] = 100 // 获取更多候选资源
	candidateResources, _, err := s.resourceRepo.SearchWithFilters(params)
	if err != nil {
		utils.Error("[TELEGRAM:PUSH] 获取候选资源失败: %v", err)
		return []interface{}{}
	}

	// 应用时间限制
	if channel.TimeLimit != "none" && len(candidateResources) > 0 {
		candidateResources = s.applyTimeFilter(candidateResources, channel.TimeLimit)
	}

	// 如果还有候选资源，随机选择一个
	if len(candidateResources) > 0 {
		// 简单随机选择（未来可以考虑使用更好的随机算法）
		randomIndex := time.Now().Nanosecond() % len(candidateResources)
		selectedResource := candidateResources[randomIndex]

		utils.Info("[TELEGRAM:PUSH] 成功获取随机资源: %s (从 %d 个候选资源中选择)",
			selectedResource.Title, len(candidateResources))
		return []interface{}{&selectedResource}
	}

	// 如果候选资源不足，回退到数据库随机函数
	defer func() {
		if r := recover(); r != nil {
			utils.Warn("[TELEGRAM:PUSH] 随机查询失败，回退到传统方法: %v", r)
		}
	}()

	randomResource, err := s.resourceRepo.GetRandomResourceWithFilters(params["category"].(string), params["tag"].(string), channel.IsPushSavedInfo)
	if err == nil && randomResource != nil {
		utils.Info("[TELEGRAM:PUSH] 使用数据库随机函数获取资源: %s", randomResource.Title)
		return []interface{}{randomResource}
	}

	return []interface{}{}
}

// applyTimeFilter 应用时间限制过滤
func (s *TelegramBotServiceImpl) applyTimeFilter(resources []entity.Resource, timeLimit string) []entity.Resource {
	now := time.Now()
	var filtered []entity.Resource

	for _, resource := range resources {
		include := false

		switch timeLimit {
		case "week":
			// 一周内
			if resource.CreatedAt.After(now.AddDate(0, 0, -7)) {
				include = true
			}
		case "month":
			// 一月内
			if resource.CreatedAt.After(now.AddDate(0, -1, 0)) {
				include = true
			}
		case "none":
			// 无限制，包含所有
			include = true
		}

		if include {
			filtered = append(filtered, resource)
		}
	}

	return filtered
}

// buildFilterParams 构建过滤参数
func (s *TelegramBotServiceImpl) buildFilterParams(channel entity.TelegramChannel) map[string]interface{} {
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

	return params
}

// buildPushMessage 构建推送消息
func (s *TelegramBotServiceImpl) buildPushMessage(channel entity.TelegramChannel, resources []interface{}) (string, string) {
	var resource *entity.Resource

	// 处理两种可能的类型：*entity.Resource 或 entity.Resource
	switch r := resources[0].(type) {
	case *entity.Resource:
		resource = r
	case entity.Resource:
		resource = &r
	default:
		utils.Error("[TELEGRAM:PUSH] 无效的资源类型: %T", resources[0])
		return "", ""
	}

	message := fmt.Sprintf("🆕 <b>%s</b>\n", s.cleanMessageTextForHTML(resource.Title))

	if resource.Description != "" {
		message += fmt.Sprintf("<blockquote>%s</blockquote>\n", s.cleanMessageTextForHTML(resource.Description))
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
	message += fmt.Sprintf("\n💡 评论区评论 (<code>【%v】%s</code>) 即可获取资源，括号内名称点击可复制📋\n", resource.ID, resource.Title)

	img := ""
	if resource.Cover != "" {
		img = resource.Cover
	} else {
		// 从 readyRepo 中取出 extra 字段，解析 JSON 获取 fid，用于构造图片URL
		// readyResources, err := s.readyRepo.FindByKey(resource.Key)
		// if err == nil && len(readyResources) > 0 {
		// 	readyResource := readyResources[0]
		// 	if readyResource.Extra != "" {
		// 		var extraData map[string]interface{}
		// 		if err := json.Unmarshal([]byte(readyResource.Extra), &extraData); err == nil {
		// 			if fid, ok := extraData["fid"].(string); ok && fid != "" {
		// 				img = fid
		// 			}
		// 		}
		// 	}
		// }
	}

	return message, img
}

func (s *TelegramBotServiceImpl) GetImgUrl(fid string) string {
	return fmt.Sprintf("http://tg.9book.top:3000/api/tool/file/%s", fid)
}

// GetBotUsername 获取机器人用户名
func (s *TelegramBotServiceImpl) GetBotUsername() string {
	if s.bot != nil {
		return s.bot.Self.UserName
	}
	return ""
}

// SendMessage 发送消息（默认使用 HTML 格式）
func (s *TelegramBotServiceImpl) SendMessage(chatID int64, text string, img string) error {
	// 检查机器人是否正在运行且已启用
	if !s.isRunning || !s.config.Enabled {
		utils.Info("[TELEGRAM:MESSAGE] 机器人已停止或禁用，跳过发送消息: ChatID=%d", chatID)
		return fmt.Errorf("机器人已停止或禁用")
	}

	if img == "" {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		_, err := s.bot.Send(msg)
		if err != nil {
			utils.Error("[TELEGRAM:MESSAGE:ERROR] 发送消息失败: %v", err)
		}
		return err
	} else {
		// 如果 img 以 http 开头，则为图片URL，否则为文件remote_id
		if strings.HasPrefix(img, "http") {
			// 发送图片URL
			photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(img))
			photoMsg.Caption = text
			photoMsg.ParseMode = "HTML"
			_, err := s.bot.Send(photoMsg)
			if err != nil {
				utils.Error("[TELEGRAM:MESSAGE:ERROR] 发送图片消息失败: %v", err)
			}
			return err
		} else {
			// imgUrl := s.GetImgUrl(img)
			//todo  判断 imgUrl 是否可用
			// 发送文件ID
			photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(img))
			photoMsg.Caption = text
			photoMsg.ParseMode = "HTML"
			_, err := s.bot.Send(photoMsg)
			if err != nil {
				utils.Error("[TELEGRAM:MESSAGE:ERROR] 发送图片消息失败: %v", err)
			}
			return err
		}
	}
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
		PushFrequency:     5,       // 默认5分钟
		PushStartTime:     "08:30", // 默认开始时间8:30
		PushEndTime:       "11:30", // 默认结束时间11:30
		IsActive:          true,
		RegisteredBy:      "bot_command",
		RegisteredAt:      time.Now(),
		ContentCategories: "",
		ContentTags:       "",
		API:               "",       // 后续可配置
		Token:             "",       // 后续可配置
		ApiType:           "l9",     // 默认l9类型
		IsPushSavedInfo:   false,    // 默认推送所有资源
		ResourceStrategy:  "random", // 默认纯随机
		TimeLimit:         "none",   // 默认无限制
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

// hasActiveGroup 检查当前活跃的群组数量
func (s *TelegramBotServiceImpl) hasActiveGroup() int {
	channels, err := s.channelRepo.FindByChatType("group")
	if err != nil {
		utils.Error("[TELEGRAM:LIMIT] 检查活跃群组失败: %v", err)
		return 0
	}

	// 统计活跃的群组数量
	activeCount := 0
	for _, channel := range channels {
		if channel.IsActive {
			activeCount++
		}
	}
	return activeCount
}

// hasActiveChannel 检查当前活跃的频道数量
func (s *TelegramBotServiceImpl) hasActiveChannel() int {
	channels, err := s.channelRepo.FindByChatType("channel")
	if err != nil {
		utils.Error("[TELEGRAM:LIMIT] 检查活跃频道失败: %v", err)
		return 0
	}

	// 统计活跃的频道数量
	activeCount := 0
	for _, channel := range channels {
		if channel.IsActive {
			activeCount++
		}
	}
	return activeCount
}

// hasActiveTelegramItems 检查当前活跃的Telegram项目（频道+群组）总数
func (s *TelegramBotServiceImpl) hasActiveTelegramItems() int {
	chatTypes := []string{"channel", "group"}
	channels, err := s.channelRepo.FindActiveChannelsByTypes(chatTypes)
	if err != nil {
		utils.Error("[TELEGRAM:LIMIT] 检查活跃Telegram项目失败: %v", err)
		return 0
	}
	return len(channels)
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

		// 检查当前活跃的Telegram项目总数（频道+群组）
		activeItemCount := s.hasActiveTelegramItems()
		if activeItemCount >= 3 {
			errorMsg := "❌ *注册限制*\n\n系统最多支持注册3个频道/群组用于推送。\n\n当前已注册: %d个，请先注销现有频道/群组，然后再注册新的。"
			s.sendReply(message, fmt.Sprintf(errorMsg, activeItemCount))
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
		if existingChannel.ResourceStrategy == "" {
			existingChannel.ResourceStrategy = "random"
		}
		if existingChannel.TimeLimit == "" {
			existingChannel.TimeLimit = "none"
		}
		if existingChannel.PushFrequency == 0 {
			existingChannel.PushFrequency = 5
		}
		if existingChannel.PushStartTime == "" {
			existingChannel.PushStartTime = "08:30"
		}
		if existingChannel.PushEndTime == "" {
			existingChannel.PushEndTime = "11:30"
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
		PushFrequency:     1, // 默认1分钟
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

// savePushHistory 保存指定频道的推送历史记录到文件（每行一个消息ID）
func (s *TelegramBotServiceImpl) savePushHistory(chatID int64) {
	// 获取指定频道的历史记录
	history, exists := s.pushHistory[chatID]
	if !exists {
		history = []uint{}
	}

	// 确保目录存在
	dir := "./data/telegram_push_history"
	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.Error("[TELEGRAM:PUSH] 创建数据目录失败: %v", err)
		return
	}

	// 写入文件，每个频道一个文件，每行一个消息ID
	filename := filepath.Join(dir, fmt.Sprintf("%d.txt", chatID))

	// 构建文件内容（每行一个消息ID）
	var content strings.Builder
	for _, resourceID := range history {
		content.WriteString(fmt.Sprintf("%d\n", resourceID))
	}

	if err := os.WriteFile(filename, []byte(content.String()), 0644); err != nil {
		utils.Error("[TELEGRAM:PUSH] 保存推送历史记录到文件失败: %v", err)
		return
	}

	utils.Debug("[TELEGRAM:PUSH] 成功保存频道 %d 的推送历史记录到文件: %s", chatID, filename)
}

// addPushedResourceID 添加已推送的资源ID到历史记录
func (s *TelegramBotServiceImpl) addPushedResourceID(chatID int64, resourceID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前频道的历史记录
	history := s.pushHistory[chatID]
	if history == nil {
		history = []uint{}
	}

	// 检查是否已经超过5000条记录
	if len(history) >= 5000 {
		// 移除旧的2500条记录，保留最新的2500条记录
		startIndex := len(history) - 2500
		history = history[startIndex:]
		utils.Info("[TELEGRAM:PUSH] 频道 %d 推送历史记录已满(5000条)，移除旧的2500条记录，保留最新的2500条", chatID)
	}

	// 添加新的资源ID到历史记录
	history = append(history, resourceID)
	s.pushHistory[chatID] = history

	utils.Debug("[TELEGRAM:PUSH] 添加推送历史，ChatID: %d, ResourceID: %d, 当前历史记录数: %d",
		chatID, resourceID, len(history))

	// 保存到文件（只保存当前频道）
	s.savePushHistory(chatID)
}

// loadPushHistory 从文件加载推送历史记录（每行一个消息ID）
func (s *TelegramBotServiceImpl) loadPushHistory() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查目录是否存在
	dir := "./data/telegram_push_history"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		utils.Info("[TELEGRAM:PUSH] 推送历史记录目录不存在，使用空的历史记录")
		return nil
	}

	// 读取目录中的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		utils.Error("[TELEGRAM:PUSH] 读取推送历史记录目录失败: %v", err)
		return err
	}

	// 初始化推送历史记录映射
	s.pushHistory = make(map[int64][]uint)

	// 遍历所有文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 检查文件名格式是否为 *.txt
		filename := file.Name()
		if !strings.HasSuffix(filename, ".txt") {
			continue
		}

		// 提取chatID
		chatIDStr := strings.TrimSuffix(filename, ".txt")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			utils.Warn("[TELEGRAM:PUSH] 无法解析频道ID文件名: %s", filename)
			continue
		}

		// 读取文件内容
		fullPath := filepath.Join(dir, filename)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			utils.Error("[TELEGRAM:PUSH] 读取推送历史记录文件失败: %s, %v", fullPath, err)
			continue
		}

		// 解析每行的消息ID
		lines := strings.Split(string(data), "\n")
		var resourceIDs []uint

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			resourceID, err := strconv.ParseUint(line, 10, 32)
			if err != nil {
				utils.Warn("[TELEGRAM:PUSH] 无法解析消息ID: %s in file %s", line, filename)
				continue
			}

			resourceIDs = append(resourceIDs, uint(resourceID))
		}

		// 只保留最多5000条记录
		if len(resourceIDs) > 5000 {
			// 保留最新的5000条记录
			startIndex := len(resourceIDs) - 5000
			resourceIDs = resourceIDs[startIndex:]
		}

		s.pushHistory[chatID] = resourceIDs
		utils.Debug("[TELEGRAM:PUSH] 加载频道 %d 的历史记录，共 %d 条", chatID, len(resourceIDs))
	}

	utils.Info("[TELEGRAM:PUSH] 成功从文件加载推送历史记录，共 %d 个频道", len(s.pushHistory))
	return nil
}

// getRecentlyPushedResourceIDs 获取最近推送过的资源ID列表
func (s *TelegramBotServiceImpl) getRecentlyPushedResourceIDs(chatID int64) []uint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回该频道的推送历史记录
	if history, exists := s.pushHistory[chatID]; exists {
		utils.Debug("[TELEGRAM:PUSH] 获取推送历史，ChatID: %d, 历史记录数: %d", chatID, len(history))
		// 返回副本，避免外部修改
		result := make([]uint, len(history))
		copy(result, history)
		return result
	}

	utils.Debug("[TELEGRAM:PUSH] 获取推送历史，ChatID: %d, 无历史记录", chatID)
	return []uint{}
}

// excludePushedResources 从候选资源中排除已推送过的资源
func (s *TelegramBotServiceImpl) excludePushedResources(resources []entity.Resource, excludeIDs []uint) []entity.Resource {
	if len(excludeIDs) == 0 {
		return resources
	}

	utils.Debug("[TELEGRAM:PUSH] 排除 %d 个已推送资源", len(excludeIDs))

	// 创建排除ID的映射，提高查找效率
	excludeMap := make(map[uint]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	// 过滤资源列表
	var filtered []entity.Resource
	for _, resource := range resources {
		if !excludeMap[resource.ID] {
			filtered = append(filtered, resource)
		}
	}

	utils.Debug("[TELEGRAM:PUSH] 过滤后剩余 %d 个资源", len(filtered))
	return filtered
}

// filterChannelsByTimeRange 过滤出在允许推送时间段内的频道
func (s *TelegramBotServiceImpl) filterChannelsByTimeRange(channels []entity.TelegramChannel) []entity.TelegramChannel {
	now := time.Now()
	currentTime := now.Format("15:04") // HH:MM 格式

	var filteredChannels []entity.TelegramChannel

	for _, channel := range channels {
		// 检查是否在推送时间段内
		if !s.isChannelInPushTimeRange(channel, currentTime) {
			utils.Info("[TELEGRAM:PUSH] 频道 %s 不在推送时间段内 (当前: %s, 允许: %s-%s)",
				channel.ChatName, currentTime, channel.PushStartTime, channel.PushEndTime)
			continue
		}

		filteredChannels = append(filteredChannels, channel)
	}

	utils.Info("[TELEGRAM:PUSH] 时间段过滤结果: %d/%d 个频道在允许推送时间段内",
		len(filteredChannels), len(channels))
	return filteredChannels
}

// isChannelInPushTimeRange 检查频道是否在推送时间段内
func (s *TelegramBotServiceImpl) isChannelInPushTimeRange(channel entity.TelegramChannel, currentTime string) bool {
	// 如果开始时间或结束时间为空，允许推送
	if channel.PushStartTime == "" || channel.PushEndTime == "" {
		return true
	}

	startTime := channel.PushStartTime
	endTime := channel.PushEndTime

	// 比较时间（假设时间格式为 HH:MM）
	if startTime <= endTime {
		// 同一天时间段，例如 08:30 - 11:30
		return currentTime >= startTime && currentTime <= endTime
	} else {
		// 跨天时间段，例如 22:00 - 06:00
		return currentTime >= startTime || currentTime <= endTime
	}
}
