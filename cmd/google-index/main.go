package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ctwj/urldb/pkg/google"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	config := &google.Config{
		CredentialsFile: "credentials.json",
		SiteURL:         "https://your-site.com/", // 替换为你的网站
		TokenFile:       "token.json",
	}

	client, err := google.NewClient(config)
	if err != nil {
		log.Fatalf("创建Google客户端失败: %v", err)
	}

	switch command {
	case "inspect":
		if len(os.Args) < 3 {
			fmt.Println("请提供要检查的URL")
			return
		}
		inspectSingleURL(client, os.Args[2])

	case "batch":
		if len(os.Args) < 3 {
			fmt.Println("请提供包含URL列表的文件")
			return
		}
		batchInspectURLs(client, os.Args[2])

	case "sites":
		listSites(client)

	case "analytics":
		getAnalytics(client)

	case "sitemap":
		if len(os.Args) < 3 {
			fmt.Println("请提供网站地图URL")
			return
		}
		submitSitemap(client, os.Args[2])

	default:
		fmt.Printf("未知命令: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Google Search Console API 工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  google-index inspect <url>        - 检查单个URL的索引状态")
	fmt.Println("  google-index batch <file>         - 批量检查URL状态")
	fmt.Println("  google-index sites                - 列出已验证的网站")
	fmt.Println("  google-index analytics            - 获取搜索分析数据")
	fmt.Println("  google-index sitemap <url>        - 提交网站地图")
	fmt.Println()
	fmt.Println("配置:")
	fmt.Println("  - 创建 credentials.json 文件 (Google Cloud Console 下载)")
	fmt.Println("  - 修改 config.SiteURL 为你的网站URL")
}

func inspectSingleURL(client *google.Client, url string) {
	fmt.Printf("正在检查URL: %s\n", url)

	result, err := client.InspectURL(url)
	if err != nil {
		log.Printf("检查失败: %v", err)
		return
	}

	printResult(url, result)
}

func batchInspectURLs(client *google.Client, filename string) {
	urls, err := readURLsFromFile(filename)
	if err != nil {
		log.Fatalf("读取URL文件失败: %v", err)
	}

	fmt.Printf("开始批量检查 %d 个URL...\n", len(urls))

	results := make(chan struct {
		url    string
		result *google.URLInspectionResult
		err    error
	}, len(urls))

	client.BatchInspectURL(urls, func(url string, result *google.URLInspectionResult, err error) {
		results <- struct {
			url    string
			result *google.URLInspectionResult
			err    error
		}{url, result, err}
	})

	// 收集并打印结果
	fmt.Println("\n检查结果:")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-50s %-15s %-15s %-20s\n", "URL", "索引状态", "移动友好", "最后抓取")
	fmt.Println(strings.Repeat("-", 100))

	for i := 0; i < len(urls); i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("%-50s %-15s\n", truncate(res.url, 47), "ERROR")
			continue
		}

		indexStatus := res.result.IndexStatusResult.IndexingState
		mobileFriendly := "否"
		if res.result.MobileUsabilityResult.MobileFriendly {
			mobileFriendly = "是"
		}
		lastCrawl := res.result.IndexStatusResult.LastCrawled
		if lastCrawl == "" {
			lastCrawl = "未知"
		}

		fmt.Printf("%-50s %-15s %-15s %-20s\n",
			truncate(res.url, 47), indexStatus, mobileFriendly, lastCrawl)
	}

	fmt.Println(strings.Repeat("-", 100))
}

func listSites(client *google.Client) {
	sites, err := client.GetSites()
	if err != nil {
		log.Printf("获取网站列表失败: %v", err)
		return
	}

	fmt.Println("已验证的网站:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-50s %-15s %-15s\n", "网站URL", "权限级别", "验证状态")
	fmt.Println(strings.Repeat("-", 80))

	for _, site := range sites {
		permissionLevel := string(site.PermissionLevel)
		verified := "否"
		if site.SiteUrl == client.SiteURL {
			verified = "是"
		}

		fmt.Printf("%-50s %-15s %-15s\n",
			truncate(site.SiteUrl, 47), permissionLevel, verified)
	}

	fmt.Println(strings.Repeat("-", 80))
}

func getAnalytics(client *google.Client) {
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02") // 最近30天

	fmt.Printf("获取搜索分析数据 (%s 到 %s)...\n", startDate, endDate)

	analytics, err := client.GetSearchAnalytics(startDate, endDate)
	if err != nil {
		log.Printf("获取分析数据失败: %v", err)
		return
	}

	// 计算总计数据
	var totalClicks, totalImpressions float64
	var totalPosition float64

	for _, row := range analytics.Rows {
		totalClicks += row.Clicks
		totalImpressions += row.Impressions
		totalPosition += row.Position
	}

	avgCTR := float64(0)
	if totalImpressions > 0 {
		avgCTR = float64(totalClicks) / float64(totalImpressions) * 100
	}

	avgPosition := float64(0)
	if len(analytics.Rows) > 0 {
		avgPosition = totalPosition / float64(len(analytics.Rows))
	}

	fmt.Println("\n搜索分析摘要:")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("总点击数: %.0f\n", totalClicks)
	fmt.Printf("总展示次数: %.0f\n", totalImpressions)
	fmt.Printf("平均点击率: %.2f%%\n", avgCTR)
	fmt.Printf("平均排名: %.1f\n", avgPosition)
	fmt.Println(strings.Repeat("-", 60))

	// 显示前10个页面
	if len(analytics.Rows) > 0 {
		fmt.Println("\n热门页面 (前10):")
		fmt.Printf("%-50s %-10s %-10s %-10s\n", "页面", "点击", "展示", "排名")
		fmt.Println(strings.Repeat("-", 80))

		maxRows := len(analytics.Rows)
		if maxRows > 10 {
			maxRows = 10
		}

		for i := 0; i < maxRows; i++ {
			row := analytics.Rows[i]
			fmt.Printf("%-50s %-10d %-10d %-10.1f\n",
				truncate(row.Keys[0], 47), row.Clicks, row.Impressions, row.Position)
		}
	}
}

func submitSitemap(client *google.Client, sitemapURL string) {
	fmt.Printf("正在提交网站地图: %s\n", sitemapURL)

	err := client.SubmitSitemap(sitemapURL)
	if err != nil {
		log.Printf("提交网站地图失败: %v", err)
		return
	}

	fmt.Println("网站地图提交成功!")
}

func readURLsFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var urls []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	return urls, nil
}

func printResult(url string, result *google.URLInspectionResult) {
	fmt.Println("\nURL检查结果:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("URL: %s\n", url)
	fmt.Println(strings.Repeat("-", 80))

	fmt.Printf("索引状态: %s\n", result.IndexStatusResult.IndexingState)
	if result.IndexStatusResult.LastCrawled != "" {
		fmt.Printf("最后抓取: %s\n", result.IndexStatusResult.LastCrawled)
	}

	fmt.Printf("移动友好: %t\n", result.MobileUsabilityResult.MobileFriendly)

	if len(result.RichResultsResult.Detected.Items) > 0 {
		fmt.Println("富媒体结果:")
		for _, item := range result.RichResultsResult.Detected.Items {
			fmt.Printf("  - %s\n", item.RichResultType)
		}
	}

	if len(result.IndexStatusResult.CrawlErrors) > 0 {
		fmt.Println("抓取错误:")
		for _, err := range result.IndexStatusResult.CrawlErrors {
			fmt.Printf("  - %s\n", err.ErrorCode)
		}
	}

	fmt.Println(strings.Repeat("=", 80))
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}