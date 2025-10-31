package services

import (
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
	Enabled         bool
	AppID           string
	AppSecret       string
	Token           string
	EncodingAesKey  string
	WelcomeMessage  string
	AutoReplyEnabled bool
	SearchLimit     int
}

// WechatBotServiceImpl 微信公众号机器人服务实现
type WechatBotServiceImpl struct {
	isRunning        bool
	systemConfigRepo repo.SystemConfigRepository
	resourceRepo     repo.ResourceRepository
	readyRepo        repo.ReadyResourceRepository
	config           *WechatBotConfig
	wechatClient     *officialaccount.OfficialAccount
}

// NewWechatBotService 创建微信公众号机器人服务
func NewWechatBotService(
	systemConfigRepo repo.SystemConfigRepository,
	resourceRepo repo.ResourceRepository,
	readyResourceRepo repo.ReadyResourceRepository,
) WechatBotService {
	return &WechatBotServiceImpl{
		isRunning:        false,
		systemConfigRepo: systemConfigRepo,
		resourceRepo:     resourceRepo,
		readyRepo:        readyResourceRepo,
		config:           &WechatBotConfig{},
	}
}