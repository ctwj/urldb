package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// TelegramHandler Telegram 处理器
type TelegramHandler struct {
	telegramChannelRepo repo.TelegramChannelRepository
	systemConfigRepo    repo.SystemConfigRepository
	telegramBotService  services.TelegramBotService
}

// NewTelegramHandler 创建 Telegram 处理器
func NewTelegramHandler(
	telegramChannelRepo repo.TelegramChannelRepository,
	systemConfigRepo repo.SystemConfigRepository,
	telegramBotService services.TelegramBotService,
) *TelegramHandler {
	return &TelegramHandler{
		telegramChannelRepo: telegramChannelRepo,
		systemConfigRepo:    systemConfigRepo,
		telegramBotService:  telegramBotService,
	}
}

// GetBotConfig 获取机器人配置
func (h *TelegramHandler) GetBotConfig(c *gin.Context) {
	configs, err := h.systemConfigRepo.GetOrCreateDefault()
	if err != nil {
		ErrorResponse(c, "获取配置失败", http.StatusInternalServerError)
		return
	}

	botConfig := converter.SystemConfigToTelegramBotConfig(configs)
	SuccessResponse(c, botConfig)
}

// UpdateBotConfig 更新机器人配置
func (h *TelegramHandler) UpdateBotConfig(c *gin.Context) {
	var req dto.TelegramBotConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 转换为系统配置实体
	configs := converter.TelegramBotConfigRequestToSystemConfigs(req)

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
	if err := h.telegramBotService.ReloadConfig(); err != nil {
		ErrorResponse(c, "重新加载机器人配置失败", http.StatusInternalServerError)
		return
	}

	// 返回成功
	SuccessResponse(c, map[string]interface{}{
		"success": true,
		"message": "配置更新成功",
	})
}

// ValidateApiKey 校验 API Key
func (h *TelegramHandler) ValidateApiKey(c *gin.Context) {
	var req dto.ValidateTelegramApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	valid, botInfo, err := h.telegramBotService.ValidateApiKey(req.ApiKey)
	if err != nil {
		ErrorResponse(c, "校验失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ValidateTelegramApiKeyResponse{
		Valid:   valid,
		BotInfo: botInfo,
	}

	if !valid {
		response.Error = "无效的 API Key"
	}

	SuccessResponse(c, response)
}

// GetChannels 获取频道列表
func (h *TelegramHandler) GetChannels(c *gin.Context) {
	channels, err := h.telegramChannelRepo.FindAll()
	if err != nil {
		ErrorResponse(c, "获取频道列表失败", http.StatusInternalServerError)
		return
	}

	channelResponses := converter.TelegramChannelsToResponse(channels)
	SuccessResponse(c, channelResponses)
}

// CreateChannel 创建频道
func (h *TelegramHandler) CreateChannel(c *gin.Context) {
	var req dto.TelegramChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 检查频道是否已存在
	existing, err := h.telegramChannelRepo.FindByChatID(req.ChatID)
	if err == nil && existing != nil {
		ErrorResponse(c, "该频道/群组已注册", http.StatusBadRequest)
		return
	}

	// 获取当前用户信息作为注册者
	username := getCurrentUsername(c) // 需要实现获取用户信息的方法

	channel := converter.RequestToTelegramChannel(req, username)

	if err := h.telegramChannelRepo.Create(&channel); err != nil {
		ErrorResponse(c, "创建频道失败", http.StatusInternalServerError)
		return
	}

	response := converter.TelegramChannelToResponse(channel)
	SuccessResponse(c, response)
}

// UpdateChannel 更新频道
func (h *TelegramHandler) UpdateChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.TelegramChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 查找现有频道
	channel, err := h.telegramChannelRepo.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "频道不存在", http.StatusNotFound)
		return
	}

	// 更新频道信息
	channel.ChatName = req.ChatName
	channel.ChatType = req.ChatType
	channel.PushEnabled = req.PushEnabled
	channel.PushFrequency = req.PushFrequency
	channel.PushStartTime = req.PushStartTime
	channel.PushEndTime = req.PushEndTime
	channel.ContentCategories = req.ContentCategories
	channel.ContentTags = req.ContentTags
	channel.IsActive = req.IsActive

	if err := h.telegramChannelRepo.Update(channel); err != nil {
		ErrorResponse(c, "更新频道失败", http.StatusInternalServerError)
		return
	}

	response := converter.TelegramChannelToResponse(*channel)
	SuccessResponse(c, response)
}

// DeleteChannel 删除频道
func (h *TelegramHandler) DeleteChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	// 检查频道是否存在
	channel, err := h.telegramChannelRepo.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "频道不存在", http.StatusNotFound)
		return
	}

	// 删除频道
	if err := h.telegramChannelRepo.Delete(uint(id)); err != nil {
		ErrorResponse(c, "删除频道失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"success": true,
		"message": "频道 " + channel.ChatName + " 已成功移除",
	})
}

// RegisterChannelByCommand 通过命令注册频道（供内部调用）
func (h *TelegramHandler) RegisterChannelByCommand(chatID int64, chatName, chatType string) error {
	// 检查是否已注册
	existing, err := h.telegramChannelRepo.FindByChatID(chatID)
	if err == nil && existing != nil {
		// 已存在，返回成功
		return nil
	}

	// 创建新的频道记录
	channel := entity.TelegramChannel{
		ChatID:        chatID,
		ChatName:      chatName,
		ChatType:      chatType,
		PushEnabled:   true,
		PushFrequency: 24, // 默认24小时
		IsActive:      true,
		RegisteredBy:  "bot_command",
	}

	return h.telegramChannelRepo.Create(&channel)
}

// HandleWebhook 处理 Telegram Webhook
func (h *TelegramHandler) HandleWebhook(c *gin.Context) {
	// 将消息交给 bot 服务处理
	// 这里可以根据需要添加身份验证
	h.telegramBotService.HandleWebhookUpdate(c)
}

// GetBotStatus 获取机器人状态
func (h *TelegramHandler) GetBotStatus(c *gin.Context) {
	// 获取机器人运行时状态
	runtimeStatus := h.telegramBotService.GetRuntimeStatus()

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
		"api_key_configured": false,
	}

	for _, config := range configs {
		switch config.Key {
		case "telegram_bot_enabled":
			configStatus["enabled"] = config.Value == "true"
		case "telegram_auto_reply_enabled":
			configStatus["auto_reply_enabled"] = config.Value == "true"
		case "telegram_bot_api_key":
			configStatus["api_key_configured"] = config.Value != ""
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

// TestBotMessage 测试机器人消息发送
func (h *TelegramHandler) TestBotMessage(c *gin.Context) {
	var req struct {
		ChatID int64  `json:"chat_id" binding:"required"`
		Text   string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	err := h.telegramBotService.SendMessage(req.ChatID, req.Text)
	if err != nil {
		ErrorResponse(c, "发送消息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"success": true,
		"message": "测试消息已发送",
	})
}

// ReloadBotConfig 重新加载机器人配置
func (h *TelegramHandler) ReloadBotConfig(c *gin.Context) {
	// 这里可以实现重新加载配置的逻辑
	// 目前通过重启服务来实现配置重新加载

	SuccessResponse(c, map[string]interface{}{
		"success": true,
		"message": "请重启服务器以重新加载配置",
		"note":    "当前版本需要重启服务器才能重新加载机器人配置",
	})
}

// DebugBotConnection 调试机器人连接
func (h *TelegramHandler) DebugBotConnection(c *gin.Context) {
	// 获取机器人状态信息用于调试
	botUsername := h.telegramBotService.GetBotUsername()

	debugInfo := map[string]interface{}{
		"bot_username":      botUsername,
		"is_running":        botUsername != "",
		"timestamp":         "2024-01-01T12:00:00Z", // 当前时间
		"debugging_enabled": true,
		"expected_logs": []string{
			"[TELEGRAM:SERVICE] Telegram Bot (@username) 已启动",
			"[TELEGRAM:MESSAGE] 开始监听 Telegram 消息更新...",
			"[TELEGRAM:MESSAGE] 消息监听循环已启动，等待消息...",
			"[TELEGRAM:MESSAGE] 收到消息: ChatID=xxx, Text='/register'",
			"[TELEGRAM:MESSAGE] 处理 /register 命令 from ChatID=xxx",
		},
		"troubleshooting_steps": []string{
			"1. 检查服务器日志中是否有 TELEGRAM 相关日志",
			"2. 确认机器人已添加到群组并设为管理员",
			"3. 验证 API Key 是否正确",
			"4. 检查自动回复是否已启用",
			"5. 重启服务器重新加载配置",
		},
	}

	SuccessResponse(c, debugInfo)
}

// GetTelegramLogs 获取Telegram相关的日志
func (h *TelegramHandler) GetTelegramLogs(c *gin.Context) {
	// 解析查询参数
	hoursStr := c.DefaultQuery("hours", "24")
	limitStr := c.DefaultQuery("limit", "100")

	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 || hours > 720 { // 最多30天
		hours = 24
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 100
	}

	// 计算时间范围
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	// 获取日志
	logs, err := utils.GetTelegramLogs(&startTime, &endTime, limit)
	if err != nil {
		ErrorResponse(c, "获取日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
		"hours": hours,
		"limit": limit,
		"start": startTime.Format("2006-01-02 15:04:05"),
		"end":   endTime.Format("2006-01-02 15:04:05"),
	})
}

// GetTelegramLogStats 获取Telegram日志统计信息
func (h *TelegramHandler) GetTelegramLogStats(c *gin.Context) {
	hoursStr := c.DefaultQuery("hours", "24")

	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 || hours > 720 {
		hours = 24
	}

	stats, err := utils.GetTelegramLogStats(hours)
	if err != nil {
		ErrorResponse(c, "获取统计信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"stats": stats,
		"hours": hours,
	})
}

// ClearTelegramLogs 清理旧的Telegram日志
func (h *TelegramHandler) ClearTelegramLogs(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 30
	}

	err = utils.ClearOldTelegramLogs(days)
	if err != nil {
		ErrorResponse(c, "清理日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": fmt.Sprintf("已清理 %d 天前的日志文件", days),
		"days":    days,
	})
}

// getCurrentUsername 获取当前用户名（临时实现）
func getCurrentUsername(c *gin.Context) string {
	// 这里应该从中间件中获取用户信息
	// 暂时返回默认值
	return "admin"
}
