package google

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Sitemap 网站地图结构
type Sitemap struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL 网站地图URL项
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// SitemapIndex 网站地图索引
type SitemapIndex struct {
	XMLName xml.Name     `xml:"sitemapindex"`
	Xmlns   string       `xml:"xmlns,attr"`
	Sitemaps []SitemapRef `xml:"sitemap"`
}

// SitemapRef 网站地图引用
type SitemapRef struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// GenerateSitemap 生成网站地图
func GenerateSitemap(urls []string, filename string) error {
	sitemap := Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, url := range urls {
		sitemapURL := SitemapURL{
			Loc:        strings.TrimSpace(url),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.8",
		}
		sitemap.URLs = append(sitemap.URLs, sitemapURL)
	}

	data, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return fmt.Errorf("生成XML失败: %v", err)
	}

	// 添加XML头部
	xmlData := []byte(xml.Header + string(data))

	return os.WriteFile(filename, xmlData, 0644)
}

// GenerateSitemapIndex 生成网站地图索引
func GenerateSitemapIndex(sitemaps []string, filename string) error {
	index := SitemapIndex{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, sitemap := range sitemaps {
		ref := SitemapRef{
			Loc:     strings.TrimSpace(sitemap),
			LastMod: time.Now().Format("2006-01-02"),
		}
		index.Sitemaps = append(index.Sitemaps, ref)
	}

	data, err := xml.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("生成XML失败: %v", err)
	}

	// 添加XML头部
	xmlData := []byte(xml.Header + string(data))

	return os.WriteFile(filename, xmlData, 0644)
}

// SplitSitemap 将大量URL分割成多个网站地图
func SplitSitemap(urls []string, maxURLsPerSitemap int, baseURL string) ([]string, error) {
	if maxURLsPerSitemap <= 0 {
		maxURLsPerSitemap = 50000 // Google限制
	}

	var sitemapFiles []string
	totalURLs := len(urls)
	sitemapCount := (totalURLs + maxURLsPerSitemap - 1) / maxURLsPerSitemap

	for i := 0; i < sitemapCount; i++ {
		start := i * maxURLsPerSitemap
		end := start + maxURLsPerSitemap
		if end > totalURLs {
			end = totalURLs
		}

		sitemapURLs := urls[start:end]
		filename := fmt.Sprintf("sitemap_part_%d.xml", i+1)

		err := GenerateSitemap(sitemapURLs, filename)
		if err != nil {
			return nil, fmt.Errorf("生成网站地图 %s 失败: %v", filename, err)
		}

		sitemapFiles = append(sitemapFiles, baseURL+filename)
		fmt.Printf("生成网站地图: %s (%d URLs)\n", filename, len(sitemapURLs))
	}

	// 生成网站地图索引
	if len(sitemapFiles) > 1 {
		indexFiles := make([]string, len(sitemapFiles))
		for i, file := range sitemapFiles {
			indexFiles[i] = file
		}

		err := GenerateSitemapIndex(indexFiles, "sitemap_index.xml")
		if err != nil {
			return nil, fmt.Errorf("生成网站地图索引失败: %v", err)
		}

		fmt.Printf("生成网站地图索引: sitemap_index.xml\n")
	}

	return sitemapFiles, nil
}

// PingSearchEngines 通知搜索引擎网站地图更新
func PingSearchEngines(sitemapURL string) error {
	searchEngines := []string{
		"http://www.google.com/webmasters/sitemaps/ping?sitemap=",
		"http://www.bing.com/webmaster/ping.aspx?siteMap=",
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, engine := range searchEngines {
		fullURL := engine + sitemapURL

		resp, err := client.Get(fullURL)
		if err != nil {
			fmt.Printf("通知搜索引擎失败 %s: %v\n", engine, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == 200 {
			fmt.Printf("成功通知: %s\n", engine)
		} else {
			fmt.Printf("通知失败: %s (状态码: %d)\n", engine, resp.StatusCode)
		}
	}

	return nil
}