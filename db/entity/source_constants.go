package entity

// 来源渠道常量（搜索 / 获取资源行为的统计归因）
// Feature: 009-statistics-enhancement；telegram 渠道由 011-telegram-bot-enhance 启用
//
// 取值为 string，便于后续扩展新渠道而无需数据库迁移。
// SearchStat.Source / ResourceView.Source 共用本组常量。
const (
	SourceWeb      = "web"      // 网页前端
	SourceWechat   = "wechat"   // 微信公众号
	SourceTelegram = "telegram" // 电报机器人（011-telegram-bot-enhance）
	// 预留（首期不实现，仅占位以保证取值可扩展）：
	// SourceAPI = "api"
)

// SourceDisplayName 返回来源渠道的中文展示名；未知来源原样返回。
func SourceDisplayName(source string) string {
	switch source {
	case SourceWeb:
		return "网页"
	case SourceWechat:
		return "公众号"
	case SourceTelegram:
		return "电报"
	default:
		return source
	}
}
