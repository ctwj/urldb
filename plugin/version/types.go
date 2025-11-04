package version

import (
	"time"
)

// PluginVersion 插件版本信息
type PluginVersion struct {
	ID          string            `json:"id"`
	PluginName  string            `json:"plugin_name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ReleaseDate time.Time         `json:"release_date"`
	Size        int64             `json:"size"`
	DownloadURL string            `json:"download_url"`
	Checksum    string            `json:"checksum"`
	Dependencies []string         `json:"dependencies"`
	Compatibility []string        `json:"compatibility"`
	Changelog   string            `json:"changelog"`
	Status      VersionStatus     `json:"status"` // 版本状态
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// VersionStatus 版本状态
type VersionStatus string

const (
	VersionStatusActive   VersionStatus = "active"   // 激活状态
	VersionStatusInactive VersionStatus = "inactive" // 未激活状态
	VersionStatusDeprecated VersionStatus = "deprecated" // 已废弃
	VersionStatusBeta     VersionStatus = "beta"     // 测试版本
	VersionStatusStable   VersionStatus = "stable"   // 稳定版本
)

// VersionHistory 版本历史
type VersionHistory struct {
	PluginName string          `json:"plugin_name"`
	Versions   []PluginVersion `json:"versions"`
	Current    string          `json:"current"` // 当前版本
}

// VersionComparison 版本比较结果
type VersionComparison struct {
	Current   string `json:"current"`
	Available string `json:"available"`
	IsNewer   bool   `json:"is_newer"`
	IsOlder   bool   `json:"is_older"`
	IsSame    bool   `json:"is_same"`
}

// VersionRequirement 版本要求
type VersionRequirement struct {
	PluginName string `json:"plugin_name"`
	MinVersion string `json:"min_version"` // 最小版本要求
	MaxVersion string `json:"max_version"` // 最大版本要求
	ExactVersion string `json:"exact_version"` // 精确版本要求
}

// CompatibilityInfo 兼容性信息
type CompatibilityInfo struct {
	PluginName string   `json:"plugin_name"`
	Version    string   `json:"version"`
	Compatible bool     `json:"compatible"`
	Reason     string   `json:"reason,omitempty"` // 不兼容原因
	SystemInfo SystemInfo `json:"system_info"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
}