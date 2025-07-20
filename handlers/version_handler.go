package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
)

// VersionResponse 版本响应结构
type VersionResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Time    time.Time   `json:"time"`
}

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	versionInfo := utils.GetVersionInfo()

	response := VersionResponse{
		Success: true,
		Data:    versionInfo,
		Message: "版本信息获取成功",
		Time:    time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// GetVersionString 获取版本字符串
func GetVersionString(c *gin.Context) {
	versionString := utils.GetVersionString()

	response := VersionResponse{
		Success: true,
		Data: map[string]string{
			"version": versionString,
		},
		Message: "版本字符串获取成功",
		Time:    time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// GetFullVersionInfo 获取完整版本信息
func GetFullVersionInfo(c *gin.Context) {
	fullInfo := utils.GetFullVersionInfo()

	response := VersionResponse{
		Success: true,
		Data: map[string]string{
			"version_info": fullInfo,
		},
		Message: "完整版本信息获取成功",
		Time:    time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// CheckUpdate 检查更新
func CheckUpdate(c *gin.Context) {
	currentVersion := utils.GetVersionInfo().Version

	// 从GitHub API获取最新版本信息
	latestVersion, err := getLatestVersionFromGitHub()
	if err != nil {
		// 如果GitHub API失败，使用模拟数据
		latestVersion = "1.0.0"
	}

	hasUpdate := utils.IsVersionNewer(latestVersion, currentVersion)

	response := VersionResponse{
		Success: true,
		Data: map[string]interface{}{
			"current_version":  currentVersion,
			"latest_version":   latestVersion,
			"has_update":       hasUpdate,
			"update_available": hasUpdate,
			"update_url":       "https://github.com/ctwj/urldb/releases/latest",
		},
		Message: "更新检查完成",
		Time:    time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// getLatestVersionFromGitHub 从GitHub获取最新版本
func getLatestVersionFromGitHub() (string, error) {
	// 使用GitHub API获取最新Release
	url := "https://api.github.com/repos/ctwj/urldb/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API返回状态码: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// 移除版本号前的 'v' 前缀
	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}
