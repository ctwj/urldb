package handlers

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// WechatHandler 微信公众号处理器
type WechatHandler struct {
	wechatService    services.WechatBotService
	systemConfigRepo repo.SystemConfigRepository
}

// NewWechatHandler 创建微信公众号处理器
func NewWechatHandler(
	wechatService services.WechatBotService,
	systemConfigRepo repo.SystemConfigRepository,
) *WechatHandler {
	return &WechatHandler{
		wechatService:    wechatService,
		systemConfigRepo: systemConfigRepo,
	}
}

// HandleWechatMessage 处理微信消息推送
func (h *WechatHandler) HandleWechatMessage(c *gin.Context) {
	// 验证微信消息签名
	if !h.validateSignature(c) {
		utils.Error("[WECHAT:VALIDATE] 签名验证失败")
		c.String(http.StatusForbidden, "签名验证失败")
		return
	}

	// 处理微信验证请求
	if c.Request.Method == "GET" {
		echostr := c.Query("echostr")
		utils.Info("[WECHAT:VERIFY] 微信服务器验证成功, echostr=%s", echostr)
		c.String(http.StatusOK, echostr)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Error("[WECHAT:MESSAGE] 读取请求体失败: %v", err)
		c.String(http.StatusBadRequest, "读取请求体失败")
		return
	}

	// 解析微信消息
	var msg message.MixMessage
	if err := xml.Unmarshal(body, &msg); err != nil {
		utils.Error("[WECHAT:MESSAGE] 解析微信消息失败: %v", err)
		c.String(http.StatusBadRequest, "消息格式错误")
		return
	}

	// 处理消息
	reply, err := h.wechatService.HandleMessage(&msg)
	if err != nil {
		utils.Error("[WECHAT:MESSAGE] 处理微信消息失败: %v", err)
		c.String(http.StatusInternalServerError, "处理失败")
		return
	}

	// 如果有回复内容，发送回复
	if reply != nil {
		responseXML, err := xml.Marshal(reply)
		if err != nil {
			utils.Error("[WECHAT:MESSAGE] 序列化回复消息失败: %v", err)
			c.String(http.StatusInternalServerError, "回复失败")
			return
		}
		c.Data(http.StatusOK, "application/xml", responseXML)
	} else {
		c.String(http.StatusOK, "success")
	}
}

// GetBotConfig 获取微信机器人配置
func (h *WechatHandler) GetBotConfig(c *gin.Context) {
	configs, err := h.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取配置失败", http.StatusInternalServerError)
		return
	}

	botConfig := converter.SystemConfigToWechatBotConfig(configs)
	SuccessResponse(c, botConfig)
}

// UpdateBotConfig 更新微信机器人配置
func (h *WechatHandler) UpdateBotConfig(c *gin.Context) {
	var req dto.WechatBotConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 转换为系统配置实体
	configs := converter.WechatBotConfigRequestToSystemConfigs(req)

	// 保存配置
	if len(configs) > 0 {
		err := h.systemConfigRepo.UpsertConfigs(configs)
		if err != nil {
			ErrorResponse(c, "保存配置失败", http.StatusInternalServerError)
			return
		}
	}

	// 重新加载配置缓存
	if err := h.systemConfigRepo.SafeRefreshConfigCache(); err != nil {
		ErrorResponse(c, "刷新配置缓存失败", http.StatusInternalServerError)
		return
	}

	// 重新加载机器人服务配置
	if err := h.wechatService.ReloadConfig(); err != nil {
		ErrorResponse(c, "重新加载机器人配置失败", http.StatusInternalServerError)
		return
	}

	// 配置更新完成后，尝试启动机器人（如果未运行且配置有效）
	if startErr := h.wechatService.Start(); startErr != nil {
		utils.Warn("[WECHAT:HANDLER] 配置更新后尝试启动机器人失败: %v", startErr)
		// 启动失败不影响配置保存，只记录警告
	}

	// 返回成功
	SuccessResponse(c, map[string]interface{}{
		"success": true,
		"message": "配置更新成功，机器人已尝试启动",
	})
}

// GetBotStatus 获取机器人状态
func (h *WechatHandler) GetBotStatus(c *gin.Context) {
	// 获取机器人运行时状态
	runtimeStatus := h.wechatService.GetRuntimeStatus()

	// 获取配置状态
	configs, err := h.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取配置失败", http.StatusInternalServerError)
		return
	}

	// 解析配置状态
	configStatus := map[string]interface{}{
		"enabled":            false,
		"auto_reply_enabled": false,
		"app_id_configured":  false,
		"token_configured":   false,
	}

	for _, config := range configs {
		switch config.Key {
		case entity.ConfigKeyWechatBotEnabled:
			configStatus["enabled"] = config.Value == "true"
		case entity.ConfigKeyWechatAutoReplyEnabled:
			configStatus["auto_reply_enabled"] = config.Value == "true"
		case entity.ConfigKeyWechatAppId:
			configStatus["app_id_configured"] = config.Value != ""
		case entity.ConfigKeyWechatToken:
			configStatus["token_configured"] = config.Value != ""
		}
	}

	// 合并状态信息
	status := map[string]interface{}{
		"config":         configStatus,
		"runtime":        runtimeStatus,
		"overall_status": runtimeStatus["is_running"].(bool),
		"status_text": func() string {
			if runtimeStatus["is_running"].(bool) {
				return "运行中"
			} else if configStatus["enabled"].(bool) {
				return "已启用但未运行"
			} else {
				return "已停止"
			}
		}(),
	}

	SuccessResponse(c, status)
}

// validateSignature 验证微信消息签名
func (h *WechatHandler) validateSignature(c *gin.Context) bool {
	// 获取配置中的Token
	configs, err := h.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		return false
	}

	var token string
	for _, config := range configs {
		if config.Key == entity.ConfigKeyWechatToken {
			token = config.Value
			break
		}
	}

	if token == "" {
		// 如果没有配置Token，跳过签名验证（开发模式）
		return true
	}

	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")

	// 验证签名
	tmpArr := []string{token, timestamp, nonce}
	sort.Strings(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	tmpStr = fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))

	return tmpStr == signature
}