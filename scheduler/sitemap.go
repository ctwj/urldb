package scheduler

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

const (
	SITEMAP_MAX_URLS = 50000 // 每个sitemap最多5万个URL
	SITEMAP_DIR      = "./data/sitemap" // sitemap文件目录
)

// SitemapScheduler Sitemap调度器
type SitemapScheduler struct {
	*BaseScheduler
	sitemapConfig entity.SystemConfig
	stopChan      chan bool
	isRunning     bool
}

// NewSitemapScheduler 创建Sitemap调度器
func NewSitemapScheduler(baseScheduler *BaseScheduler) *SitemapScheduler {
	return &SitemapScheduler{
		BaseScheduler: baseScheduler,
		stopChan:      make(chan bool),
		isRunning:     false,
	}
}

// Start 启动Sitemap调度任务
func (s *SitemapScheduler) Start() {
	if s.IsRunning() {
		utils.Debug("Sitemap定时任务已在运行中")
		return
	}

	s.SetRunning(true)
	utils.Info("开始启动Sitemap定时任务")

	go s.run()
}

// Stop 停止Sitemap调度任务
func (s *SitemapScheduler) Stop() {
	if !s.IsRunning() {
		utils.Debug("Sitemap定时任务未在运行")
		return
	}

	utils.Info("正在停止Sitemap定时任务")
	s.stopChan <- true
	s.SetRunning(false)
}

// IsRunning 检查Sitemap调度任务是否在运行
func (s *SitemapScheduler) IsRunning() bool {
	return s.isRunning
}

// SetRunning 设置运行状态
func (s *SitemapScheduler) SetRunning(running bool) {
	s.isRunning = running
}

// GetStopChan 获取停止通道
func (s *SitemapScheduler) GetStopChan() chan bool {
	return s.stopChan
}

// run 执行调度任务的主循环
func (s *SitemapScheduler) run() {
	utils.Info("Sitemap定时任务开始运行")

	// 立即执行一次
	s.generateSitemap()

	// 定时执行
	ticker := time.NewTicker(24 * time.Hour) // 每24小时执行一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			utils.Info("定时执行Sitemap生成任务")
			s.generateSitemap()
		case <-s.stopChan:
			utils.Info("收到停止信号，Sitemap调度任务退出")
			return
		}
	}
}

// generateSitemap 生成sitemap
func (s *SitemapScheduler) generateSitemap() {
	utils.Info("开始生成Sitemap...")

	startTime := time.Now()

	// 获取资源总数
	var total int64
	if err := s.BaseScheduler.resourceRepo.GetDB().Model(&entity.Resource{}).Count(&total).Error; err != nil {
		utils.Error("获取资源总数失败: %v", err)
		return
	}

	utils.Info("需要处理的资源总数: %d", total)

	if total == 0 {
		utils.Info("没有资源需要生成Sitemap")
		return
	}

	// 计算需要多少个sitemap文件
	totalPages := int((total + SITEMAP_MAX_URLS - 1) / SITEMAP_MAX_URLS)
	utils.Info("需要生成 %d 个sitemap文件", totalPages)

	// 确保目录存在
	if err := os.MkdirAll(SITEMAP_DIR, 0755); err != nil {
		utils.Error("创建sitemap目录失败: %v", err)
		return
	}

	// 生成每个sitemap文件
	for page := 0; page < totalPages; page++ {
		if s.SleepWithStopCheck(100 * time.Millisecond) { // 避免过于频繁的检查
			utils.Info("在生成sitemap过程中收到停止信号，退出生成")
			return
		}

		utils.Info("正在生成第 %d 个sitemap文件", page+1)

		if err := s.generateSitemapPage(page); err != nil {
			utils.Error("生成第 %d 个sitemap文件失败: %v", page, err)
		} else {
			utils.Info("成功生成第 %d 个sitemap文件", page+1)
		}
	}

	// 生成sitemap索引文件
	if err := s.generateSitemapIndex(totalPages); err != nil {
		utils.Error("生成sitemap索引文件失败: %v", err)
	} else {
		utils.Info("成功生成sitemap索引文件")
	}

	// 尝试获取网站基础URL
	baseURL, err := s.BaseScheduler.systemConfigRepo.GetConfigValue(entity.ConfigKeyWebsiteURL)
	if err != nil || baseURL == "" {
		baseURL = "https://yoursite.com" // 默认值
	}

	utils.Info("Sitemap生成完成，耗时: %v", time.Since(startTime))
	utils.Info("Sitemap地址: %s/sitemap.xml", baseURL)
}

// generateSitemapPage 生成单个sitemap页面
func (s *SitemapScheduler) generateSitemapPage(page int) error {
	offset := page * SITEMAP_MAX_URLS
	limit := SITEMAP_MAX_URLS

	var resources []entity.Resource
	if err := s.BaseScheduler.resourceRepo.GetDB().Offset(offset).Limit(limit).Find(&resources).Error; err != nil {
		return fmt.Errorf("获取资源数据失败: %w", err)
	}

	var urls []Url
	for _, resource := range resources {
		lastMod := resource.UpdatedAt
		if resource.CreatedAt.After(lastMod) {
			lastMod = resource.CreatedAt
		}

		urls = append(urls, Url{
			Loc:        fmt.Sprintf("/r/%s", resource.Key),
			LastMod:    lastMod.Format("2006-01-02"), // 只保留日期部分
			ChangeFreq: "weekly",
			Priority:   0.8,
		})
	}

	urlSet := UrlSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	filename := filepath.Join(SITEMAP_DIR, fmt.Sprintf("sitemap-%d.xml", page))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	file.WriteString(xml.Header)
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlSet); err != nil {
		return fmt.Errorf("写入XML失败: %w", err)
	}

	return nil
}

// generateSitemapIndex 生成sitemap索引文件
func (s *SitemapScheduler) generateSitemapIndex(totalPages int) error {
	// 构建主机URL - 这里使用默认URL，实际应用中应从配置获取
	baseURL, err := s.BaseScheduler.systemConfigRepo.GetConfigValue(entity.ConfigKeyWebsiteURL)
	if err != nil || baseURL == "" {
		baseURL = "https://yoursite.com" // 默认值
	}

	// 移除URL末尾的斜杠
	baseURL = strings.TrimSuffix(baseURL, "/")

	var sitemaps []Sitemap
	for i := 0; i < totalPages; i++ {
		sitemapURL := fmt.Sprintf("%s/sitemap-%d.xml", baseURL, i)
		sitemaps = append(sitemaps, Sitemap{
			Loc:     sitemapURL,
			LastMod: time.Now().Format("2006-01-02"),
		})
	}

	sitemapIndex := SitemapIndex{
		XMLNS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: sitemaps,
	}

	filename := filepath.Join(SITEMAP_DIR, "sitemap.xml")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建索引文件失败: %w", err)
	}
	defer file.Close()

	file.WriteString(xml.Header)
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(sitemapIndex); err != nil {
		return fmt.Errorf("写入索引XML失败: %w", err)
	}

	return nil
}

// GetSitemapConfig 获取Sitemap配置
func (s *SitemapScheduler) GetSitemapConfig() (bool, error) {
	configStr, err := s.BaseScheduler.systemConfigRepo.GetConfigValue(entity.ConfigKeySitemapConfig)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	// 解析配置字符串，这里简化处理
	return configStr == "1" || configStr == "true", nil
}

// UpdateSitemapConfig 更新Sitemap配置
func (s *SitemapScheduler) UpdateSitemapConfig(enabled bool) error {
	configStr := "0"
	if enabled {
		configStr = "1"
	}

	config := entity.SystemConfig{
		Key:   entity.ConfigKeySitemapConfig,
		Value: configStr,
		Type:  "bool",
	}

	// 由于repository没有直接的SetConfig方法，我们使用UpsertConfigs
	configs := []entity.SystemConfig{config}
	return s.BaseScheduler.systemConfigRepo.UpsertConfigs(configs)
}

// UrlSet sitemap内容
type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []Url    `xml:"url"`
}

// Url 单个URL信息
type Url struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

// SitemapIndex sitemap索引结构
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	XMLNS    string    `xml:"xmlns,attr"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// Sitemap 单个sitemap信息
type Sitemap struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}