package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ctwj/urldb/utils"
)

// RemoteFetcher 远程配置获取器
type RemoteFetcher struct {
	repoURL      string        // GitHub仓库URL
	filePath     string        // 配置文件路径
	branch       string        // 分支名
	timeout      time.Duration // 请求超时
	lastFetch    time.Time     // 上次获取时间
	lastVersion  int           // 上次获取的版本
	lastHash     string        // 上次获取的配置哈希
}

// NewRemoteFetcher 创建远程配置获取器
func NewRemoteFetcher(repoURL, filePath, branch string) *RemoteFetcher {
	return &RemoteFetcher{
		repoURL:  repoURL,
		filePath: filePath,
		branch:   branch,
		timeout:  10 * time.Second,
	}
}

// SetTimeout 设置请求超时
func (f *RemoteFetcher) SetTimeout(timeout time.Duration) {
	f.timeout = timeout
}

// Fetch 获取远程配置
func (f *RemoteFetcher) Fetch() (*RemoteConfig, error) {
	url := f.buildURL()
	utils.Info("Fetching remote config from: %s", url)

	client := &http.Client{
		Timeout: f.timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch config: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var config RemoteConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	hash := sha256.Sum256(body)
	f.lastHash = hex.EncodeToString(hash[:])
	f.lastFetch = time.Now()
	f.lastVersion = config.Version

	utils.Info("Successfully fetched remote config version %d", config.Version)

	return &config, nil
}

// CheckUpdate 检查是否有更新
func (f *RemoteFetcher) CheckUpdate() (bool, int, error) {
	url := f.buildURL()

	client := &http.Client{
		Timeout: f.timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, 0, fmt.Errorf("failed to check update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, 0, fmt.Errorf("failed to check update: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var config RemoteConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return false, 0, fmt.Errorf("failed to parse config: %w", err)
	}

	hasUpdate := config.Version > f.lastVersion
	return hasUpdate, config.Version, nil
}

// GetLastVersion 获取上次获取的版本
func (f *RemoteFetcher) GetLastVersion() int {
	return f.lastVersion
}

// GetLastFetchTime 获取上次获取时间
func (f *RemoteFetcher) GetLastFetchTime() time.Time {
	return f.lastFetch
}

// GetLastHash 获取上次配置哈希
func (f *RemoteFetcher) GetLastHash() string {
	return f.lastHash
}

// buildURL 构建GitHub Raw URL
func (f *RemoteFetcher) buildURL() string {
	return fmt.Sprintf("%s/raw/%s/%s", f.repoURL, f.branch, f.filePath)
}