package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

// VersionInfo 版本信息结构
type VersionInfo struct {
	Version     string    `json:"version"`
	BuildTime   time.Time `json:"build_time"`
	GitCommit   string    `json:"git_commit"`
	GitBranch   string    `json:"git_branch"`
	GoVersion   string    `json:"go_version"`
	NodeVersion string    `json:"node_version"`
	Platform    string    `json:"platform"`
	Arch        string    `json:"arch"`
}

// 编译时注入的版本信息
var (
	Version   = "1.0.0"
	BuildTime = time.Now().Format("2006-01-02 15:04:05")
	GitCommit = "unknown"
	GitBranch = "unknown"
)

// GetVersionInfo 获取版本信息
func GetVersionInfo() *VersionInfo {
	buildTime, _ := time.Parse("2006-01-02 15:04:05", BuildTime)

	return &VersionInfo{
		Version:     Version,
		BuildTime:   buildTime,
		GitCommit:   GitCommit,
		GitBranch:   GitBranch,
		GoVersion:   runtime.Version(),
		NodeVersion: getNodeVersion(),
		Platform:    runtime.GOOS,
		Arch:        runtime.GOARCH,
	}
}

// GetVersionString 获取版本字符串
func GetVersionString() string {
	info := GetVersionInfo()
	return fmt.Sprintf("v%s (%s)", info.Version, info.GitCommit)
}

// GetFullVersionInfo 获取完整版本信息
func GetFullVersionInfo() string {
	info := GetVersionInfo()
	return fmt.Sprintf(`版本信息:
  版本号: v%s
  构建时间: %s
  Git提交: %s
  Git分支: %s
  Go版本: %s
  Node版本: %s
  平台: %s/%s`,
		info.Version,
		info.BuildTime.Format("2006-01-02 15:04:05"),
		info.GitCommit,
		info.GitBranch,
		info.GoVersion,
		info.NodeVersion,
		info.Platform,
		info.Arch,
	)
}

// LoadVersionFromFile 从文件加载版本信息
func LoadVersionFromFile(filename string) (*VersionInfo, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var info VersionInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// SaveVersionToFile 保存版本信息到文件
func SaveVersionToFile(filename string, info *VersionInfo) error {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// getNodeVersion 获取Node.js版本
func getNodeVersion() string {
	// 这里可以通过执行 node --version 来获取
	// 为了简化，返回一个默认值
	return "N/A"
}

// IsVersionNewer 比较版本号
func IsVersionNewer(version1, version2 string) bool {
	// 简单的版本比较，可以根据需要实现更复杂的逻辑
	return version1 > version2
}

// GetVersionComponents 获取版本号组件
func GetVersionComponents(version string) (major, minor, patch int, err error) {
	var majorStr, minorStr, patchStr string
	_, err = fmt.Sscanf(version, "%s.%s.%s", &majorStr, &minorStr, &patchStr)
	if err != nil {
		return 0, 0, 0, err
	}

	// 这里可以添加更复杂的版本号解析逻辑
	return 1, 0, 0, nil
}
