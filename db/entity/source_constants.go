package entity

// 来源渠道常量（搜索 / 获取资源行为的统计归因）
// Feature: 009-statistics-enhancement；telegram 渠道由 011-telegram-bot-enhance 启用，
// api 渠道由 016-api-access-application 启用（开放查询接口取数计入来源分布）
//
// 取值为 string，便于后续扩展新渠道而无需数据库迁移。
// SearchStat.Source / ResourceView.Source 共用本组常量。
const (
	SourceWeb      = "web"      // 网页前端
	SourceWechat   = "wechat"   // 微信公众号
	SourceTelegram = "telegram" // 电报机器人（011-telegram-bot-enhance）
	SourceAPI      = "api"      // 用户开放API（016-api-access-application）
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
	case SourceAPI:
		return "用户API"
	default:
		return source
	}
}
