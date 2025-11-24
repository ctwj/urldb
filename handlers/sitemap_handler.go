package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/scheduler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	resourceRepo       repo.ResourceRepository
	systemConfigRepo   repo.SystemConfigRepository
	hotDramaRepo       repo.HotDramaRepository
	readyResourceRepo  repo.ReadyResourceRepository
	panRepo            repo.PanRepository
	cksRepo            repo.CksRepository
	tagRepo            repo.TagRepository
	categoryRepo       repo.CategoryRepository
)

// SetSitemapDependencies 注册Sitemap处理器依赖
func SetSitemapDependencies(
	resourceRepository repo.ResourceRepository,
	systemConfigRepository repo.SystemConfigRepository,
	hotDramaRepository repo.HotDramaRepository,
	readyResourceRepository repo.ReadyResourceRepository,
	panRepository repo.PanRepository,
	cksRepository repo.CksRepository,
	tagRepository repo.TagRepository,
	categoryRepository repo.CategoryRepository,
) {
	resourceRepo = resourceRepository
	systemConfigRepo = systemConfigRepository
	hotDramaRepo = hotDramaRepository
	readyResourceRepo = readyResourceRepository
	panRepo = panRepository
	cksRepo = cksRepository
	tagRepo = tagRepository
	categoryRepo = categoryRepository
}


const SITEMAP_MAX_URLS = 50000 // 每个sitemap最多5万个URL

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

// SitemapConfig sitemap配置
type SitemapConfig struct {
	AutoGenerate bool      `json:"auto_generate"`
	LastGenerate time.Time `json:"last_generate"`
	LastUpdate   time.Time `json:"last_update"`
}

// GetSitemapConfig 获取sitemap配置
func GetSitemapConfig(c *gin.Context) {
	// 从全局调度器获取配置
	enabled, err := scheduler.GetGlobalScheduler(
		hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
		panRepo, cksRepo, tagRepo, categoryRepo,
		repoManager.TaskItemRepository,
		repoManager.TaskRepository,
	).GetSitemapConfig()
	if err != nil && err != gorm.ErrRecordNotFound {
		// 如果获取失败，尝试从配置表中获取
		configStr, err := systemConfigRepo.GetConfigValue(entity.ConfigKeySitemapAutoGenerateEnabled)
		if err != nil && err != gorm.ErrRecordNotFound {
			ErrorResponse(c, "获取配置失败", http.StatusInternalServerError)
			return
		}
		enabled = configStr == "1" || configStr == "true"
	}

	// 获取最后生成时间（从配置中获取）
	configStr, err := systemConfigRepo.GetConfigValue(entity.ConfigKeySitemapLastGenerateTime)
	if err != nil && err != gorm.ErrRecordNotFound {
		// 如果获取失败，只返回启用状态
		config := SitemapConfig{
			AutoGenerate: enabled,
			LastGenerate: time.Time{}, // 空时间
			LastUpdate:   time.Now(),
		}
		SuccessResponse(c, config)
		return
	}

	var lastGenerateTime time.Time
	if configStr != "" {
		lastGenerateTime, _ = time.Parse("2006-01-02 15:04:05", configStr)
	}

	config := SitemapConfig{
		AutoGenerate: enabled,
		LastGenerate: lastGenerateTime,
		LastUpdate:   time.Now(),
	}

	SuccessResponse(c, config)
}

// UpdateSitemapConfig 更新sitemap配置
func UpdateSitemapConfig(c *gin.Context) {
	var config SitemapConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	// 更新调度器配置
	if err := scheduler.GetGlobalScheduler(
		hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
		panRepo, cksRepo, tagRepo, categoryRepo,
		repoManager.TaskItemRepository,
		repoManager.TaskRepository,
	).UpdateSitemapConfig(config.AutoGenerate); err != nil {
		ErrorResponse(c, "更新调度器配置失败", http.StatusInternalServerError)
		return
	}

	// 保存自动生成功能状态
	autoGenerateStr := "0"
	if config.AutoGenerate {
		autoGenerateStr = "1"
	}
	autoGenerateConfig := entity.SystemConfig{
		Key:   entity.ConfigKeySitemapAutoGenerateEnabled,
		Value: autoGenerateStr,
		Type:  "bool",
	}

	// 保存最后生成时间
	lastGenerateStr := config.LastGenerate.Format("2006-01-02 15:04:05")
	lastGenerateConfig := entity.SystemConfig{
		Key:   entity.ConfigKeySitemapLastGenerateTime,
		Value: lastGenerateStr,
		Type:  "string",
	}

	configs := []entity.SystemConfig{autoGenerateConfig, lastGenerateConfig}
	if err := systemConfigRepo.UpsertConfigs(configs); err != nil {
		ErrorResponse(c, "保存配置失败", http.StatusInternalServerError)
		return
	}

	// 根据配置启动或停止调度器
	if config.AutoGenerate {
		scheduler.GetGlobalScheduler(
			hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
			panRepo, cksRepo, tagRepo, categoryRepo,
			repoManager.TaskItemRepository,
			repoManager.TaskRepository,
		).StartSitemapScheduler()
	} else {
		scheduler.GetGlobalScheduler(
			hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
			panRepo, cksRepo, tagRepo, categoryRepo,
			repoManager.TaskItemRepository,
			repoManager.TaskRepository,
		).StopSitemapScheduler()
	}

	SuccessResponse(c, config)
}

// GenerateSitemap 手动生成sitemap
func GenerateSitemap(c *gin.Context) {
	// 获取资源总数
	var total int64
	if err := resourceRepo.GetDB().Model(&entity.Resource{}).Count(&total).Error; err != nil {
		ErrorResponse(c, "获取资源总数失败", http.StatusInternalServerError)
		return
	}

	totalPages := int((total + SITEMAP_MAX_URLS - 1) / SITEMAP_MAX_URLS)

	// 获取全局调度器并立即执行sitemap生成
	globalScheduler := scheduler.GetGlobalScheduler(
		hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
		panRepo, cksRepo, tagRepo, categoryRepo,
		repoManager.TaskItemRepository,
		repoManager.TaskRepository,
	)

	// 手动触发sitemap生成
	globalScheduler.TriggerSitemapGeneration()

	// 记录最后生成时间为当前时间
	lastGenerateStr := time.Now().Format("2006-01-02 15:04:05")
	lastGenerateConfig := entity.SystemConfig{
		Key:   entity.ConfigKeySitemapLastGenerateTime,
		Value: lastGenerateStr,
		Type:  "string",
	}

	if err := systemConfigRepo.UpsertConfigs([]entity.SystemConfig{lastGenerateConfig}); err != nil {
		ErrorResponse(c, "更新最后生成时间失败", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"total_resources": total,
		"total_pages":     totalPages,
		"status":          "started",
		"message":         fmt.Sprintf("开始生成 %d 个sitemap文件", totalPages),
	}

	SuccessResponse(c, result)
}

// GetSitemapStatus 获取sitemap生成状态
func GetSitemapStatus(c *gin.Context) {
	// 获取资源总数
	var total int64
	if err := resourceRepo.GetDB().Model(&entity.Resource{}).Count(&total).Error; err != nil {
		ErrorResponse(c, "获取资源总数失败", http.StatusInternalServerError)
		return
	}

	// 计算需要生成的sitemap文件数量
	totalPages := int((total + SITEMAP_MAX_URLS - 1) / SITEMAP_MAX_URLS)

	// 获取最后生成时间
	lastGenerateStr, err := systemConfigRepo.GetConfigValue(entity.ConfigKeySitemapLastGenerateTime)
	if err != nil {
		// 如果没有记录，使用当前时间
		lastGenerateStr = time.Now().Format("2006-01-02 15:04:05")
	}

	lastGenerate, err := time.Parse("2006-01-02 15:04:05", lastGenerateStr)
	if err != nil {
		lastGenerate = time.Now()
	}

	// 检查调度器是否运行
	isRunning := scheduler.GetGlobalScheduler(
		hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
		panRepo, cksRepo, tagRepo, categoryRepo,
		repoManager.TaskItemRepository,
		repoManager.TaskRepository,
	).IsSitemapSchedulerRunning()

	// 获取自动生成功能状态
	autoGenerateEnabled, err := scheduler.GetGlobalScheduler(
		hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
		panRepo, cksRepo, tagRepo, categoryRepo,
		repoManager.TaskItemRepository,
		repoManager.TaskRepository,
	).GetSitemapConfig()
	if err != nil {
		// 如果调度器获取失败，从配置中获取
		configStr, err := systemConfigRepo.GetConfigValue(entity.ConfigKeySitemapAutoGenerateEnabled)
		if err != nil {
			autoGenerateEnabled = false
		} else {
			autoGenerateEnabled = configStr == "1" || configStr == "true"
		}
	}

	result := map[string]interface{}{
		"total_resources": total,
		"total_pages":     totalPages,
		"last_generate":   lastGenerate.Format("2006-01-02 15:04:05"),
		"status":          "ready",
		"is_running":      isRunning,
		"auto_generate":   autoGenerateEnabled,
	}

	SuccessResponse(c, result)
}

// SitemapIndexHandler sitemap索引文件处理器
func SitemapIndexHandler(c *gin.Context) {
	// 获取资源总数
	var total int64
	if err := resourceRepo.GetDB().Model(&entity.Resource{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资源总数失败"})
		return
	}

	totalPages := int((total + SITEMAP_MAX_URLS - 1) / SITEMAP_MAX_URLS)

	// 构建主机URL
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080" // 默认值
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	// 创建sitemap列表 - 现在文件保存在data/sitemap目录，通过/file/sitemap/路径访问
	var sitemaps []Sitemap
	for i := 0; i < totalPages; i++ {
		sitemapURL := fmt.Sprintf("%s/file/sitemap/sitemap-%d.xml", baseURL, i)
		sitemaps = append(sitemaps, Sitemap{
			Loc:     sitemapURL,
			LastMod: time.Now().Format("2006-01-02"),
		})
	}

	sitemapIndex := SitemapIndex{
		XMLNS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: sitemaps,
	}

	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, sitemapIndex)
}

// SitemapPageHandler sitemap页面处理器
func SitemapPageHandler(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页面参数"})
		return
	}

	offset := page * SITEMAP_MAX_URLS
	limit := SITEMAP_MAX_URLS

	var resources []entity.Resource
	if err := resourceRepo.GetDB().Offset(offset).Limit(limit).Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资源数据失败"})
		return
	}

	var urls []Url
	for _, resource := range resources {
		lastMod := resource.UpdatedAt
		if resource.CreatedAt.After(lastMod) {
			lastMod = resource.CreatedAt
		}

		urls = append(urls, Url{
			Loc:        fmt.Sprintf("/r/%s", resource.Key),
			LastMod:    lastMod.Format("2006-01-01"), // 只保留日期部分
			ChangeFreq: "weekly",
			Priority:   0.8,
		})
	}

	urlSet := UrlSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, urlSet)
}

// 手动生成完整sitemap文件
func GenerateFullSitemap(c *gin.Context) {
	// 获取资源总数
	var total int64
	if err := resourceRepo.GetDB().Model(&entity.Resource{}).Count(&total).Error; err != nil {
		ErrorResponse(c, "获取资源总数失败", http.StatusInternalServerError)
		return
	}

	// 获取全局调度器并立即执行sitemap生成
	globalScheduler := scheduler.GetGlobalScheduler(
		hotDramaRepo, readyResourceRepo, resourceRepo, systemConfigRepo,
		panRepo, cksRepo, tagRepo, categoryRepo,
		repoManager.TaskItemRepository,
		repoManager.TaskRepository,
	)

	// 手动触发sitemap生成
	globalScheduler.TriggerSitemapGeneration()

	// 记录最后生成时间为当前时间
	lastGenerateStr := time.Now().Format("2006-01-02 15:04:05")
	lastGenerateConfig := entity.SystemConfig{
		Key:   entity.ConfigKeySitemapLastGenerateTime,
		Value: lastGenerateStr,
		Type:  "string",
	}

	if err := systemConfigRepo.UpsertConfigs([]entity.SystemConfig{lastGenerateConfig}); err != nil {
		ErrorResponse(c, "更新最后生成时间失败", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"message":         "Sitemap生成任务已启动",
		"total_resources": total,
		"status":          "processing",
		"estimated_time":  fmt.Sprintf("%d秒", total/1000), // 估算时间
	}

	SuccessResponse(c, result)
}