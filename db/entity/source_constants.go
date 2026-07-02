package entity

// 来源渠道常量（搜索 / 获取资源行为的统计归因）
// Feature: 009-statistics-enhancement
//
// 取值为 string，便于后续扩展新渠道而无需数据库迁移。
// SearchStat.Source / ResourceView.Source 共用本组常量。
const (
	SourceWeb    = "web"    // 网页前端
	SourceWechat = "wechat" // 微信公众号
	// 预留（首期不实现，仅占位以保证取值可扩展）：
	// SourceTelegram = "telegram"
	// SourceAPI      = "api"
)

// SourceDisplayName 返回来源渠道的中文展示名；未知来源原样返回。
func SourceDisplayName(source string) string {
	switch source {
	case SourceWeb:
		return "网页"
	case SourceWechat:
		return "公众号"
	default:
		return source
	}
}
