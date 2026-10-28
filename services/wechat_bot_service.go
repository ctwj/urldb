package services

import (
	"sync"

	"github.com/ctwj/urldb/db/repo"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// WechatBotService 微信公众号机器人服务接口
type WechatBotService interface {
	Start() error
	Stop() error
	IsRunning() bool
	ReloadConfig() error
	HandleMessage(msg *message.MixMessage) (interface{}, error)
	SendWelcomeMessage(openID string) error
	GetRuntimeStatus() map[string]interface{}
	GetConfig() *WechatBotConfig
}

// WechatBotConfig 微信公众号机器人配置
type WechatBotConfig struct {
	Enabled             bool
	AppID               string
	AppSecret           string
	Token               string
	EncodingAesKey      string
	WelcomeMessage      string
	AutoReplyEnabled    bool
	SearchLimit         int
	DynamicCodeEnabled  bool   // 「验证码」动态口令回复开关（TOTP，独立于自动回复开关）
	DynamicCodeSecret   string // 动态口令共享密钥，需与站点侧插件配置一致
	DynamicCodeTemplate string // 回复文案模板，{{code}} 占位符替换为当期口令
}

// WechatBotServiceImpl 微信公众号机器人服务实现
type WechatBotServiceImpl struct {
	isRunning        bool
	systemConfigRepo repo.SystemConfigRepository
	resourceRepo     repo.ResourceRepository
	readyRepo        repo.ReadyResourceRepository
	config           *WechatBotConfig
	wechatClient     *officialaccount.OfficialAccount
	searchSessionManager *SearchSessionManager
	// 012-wechat-bot-transfer：统一取链 + 去重锁
	linkService      ResourceLinkService // ResolveWithCheck 决策树
	linkInFlight      sync.Map           // 取链去重锁：resourceID -> struct{}（防并发重复转存）
}

// NewWechatBotService 创建微信公众号机器人服务
func NewWechatBotService(
	systemConfigRepo repo.SystemConfigRepository,
	resourceRepo repo.ResourceRepository,
	readyResourceRepo repo.ReadyResourceRepository,
	linkService ResourceLinkService,
) WechatBotService {
	return &WechatBotServiceImpl{
		isRunning:            false,
		systemConfigRepo:     systemConfigRepo,
		resourceRepo:         resourceRepo,
		readyRepo:            readyResourceRepo,
		config:               &WechatBotConfig{},
		searchSessionManager: GlobalSearchSessionManager,
		linkService:          linkService,
	}
}
