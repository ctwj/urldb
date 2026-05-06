package config

// PanRuleConfig 网盘规则配置
type PanRuleConfig struct {
	ID          int      `json:"id"`           // 规则ID
	Name        string   `json:"name"`         // 规则名称
	PanID       int      `json:"pan_id"`       // 网盘ID
	Domains     []string `json:"domains"`      // 域名列表
	URLPatterns []string `json:"url_patterns"` // URL正则模式列表
	Priority    int      `json:"priority"`     // 优先级
	Enabled     bool     `json:"enabled"`      // 是否启用
	Remark      string   `json:"remark"`       // 备注
}

// RemoteConfig 远程配置结构
type RemoteConfig struct {
	Version     int           `json:"version"`      // 配置版本号
	UpdatedAt   string        `json:"updated_at"`   // 更新时间
	Description string        `json:"description"`  // 更新说明
	Rules       []PanRuleConfig `json:"rules"`       // 网盘规则列表
}

// ConfigSource 配置来源
type ConfigSource int

const (
	SourceHardcoded ConfigSource = iota // 硬编码配置
	SourceRemote                        // 远程配置
	SourceMixed                         // 混合配置（远程+本地）
)

func (s ConfigSource) String() string {
	switch s {
	case SourceHardcoded:
		return "hardcoded"
	case SourceRemote:
		return "remote"
	case SourceMixed:
		return "mixed"
	default:
		return "unknown"
	}
}