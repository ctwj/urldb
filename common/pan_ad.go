package pan

import (
	"log"
	"regexp"
	"strings"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// 本文件提供网盘"广告清理 / 自定义广告插入"的共享工具函数。
//
// 背景：夸克网盘（quark_pan.go）已自带一套广告处理方法实现并稳定运行；为满足
// FR-012（不回归夸克），不改动夸克的既有方法。新增网盘（如 UC）统一复用以下
// 包级函数，避免在每个 *_pan.go 中重复实现。
//
// 相关系统配置项（db/entity）：
//   - ConfigKeyAdKeywords：广告关键词（命中则删除来源广告文件）
//   - ConfigKeyAutoInsertAd：自动插入广告（转存后随机插入一条自定义广告）
//
// 依赖的包级缓存变量（configRefreshChan / systemConfigOnce / systemConfigRepo）
// 声明于 quark_pan.go，同包内直接引用。

// getAdSystemConfigValue 读取系统配置值，并在收到缓存刷新信号时清空缓存。
func getAdSystemConfigValue(key string) (string, error) {
	// 检查是否需要刷新缓存
	select {
	case <-configRefreshChan:
		systemConfigOnce.Do(func() {
			systemConfigRepo = repo.NewSystemConfigRepository(db.DB)
		})
		systemConfigRepo.ClearConfigCache()
	default:
		// 没有刷新信号，继续使用缓存
	}

	systemConfigOnce.Do(func() {
		systemConfigRepo = repo.NewSystemConfigRepository(db.DB)
	})
	return systemConfigRepo.GetConfigValue(key)
}

// splitAdKeywords 按中英文逗号分割关键词。
func splitAdKeywords(keywordsStr string) []string {
	if keywordsStr == "" {
		return []string{}
	}
	re := regexp.MustCompile(`[,，]`)
	parts := re.Split(keywordsStr, -1)

	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitAdURLs 按换行符分割广告 URL 列表。
func splitAdURLs(autoInsertAdStr string) []string {
	if autoInsertAdStr == "" {
		return []string{}
	}
	lines := strings.Split(autoInsertAdStr, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// checkAdKeywordsInFilename 检查文件名是否包含任一关键词（大小写不敏感）。
func checkAdKeywordsInFilename(filename string, keywords []string) bool {
	lowercaseFilename := strings.ToLower(filename)
	for _, keyword := range keywords {
		if strings.Contains(lowercaseFilename, strings.ToLower(keyword)) {
			log.Printf("文件 %s 包含广告关键词: %s", filename, keyword)
			return true
		}
	}
	return false
}

// containsAdKeywords 检查文件名是否命中系统配置的广告关键词。
func containsAdKeywords(filename string) bool {
	adKeywordsStr, err := getAdSystemConfigValue(entity.ConfigKeyAdKeywords)
	if err != nil {
		log.Printf("获取广告关键词配置失败: %v", err)
		return false
	}
	if adKeywordsStr == "" {
		return false
	}
	adKeywords := splitAdKeywords(adKeywordsStr)
	return checkAdKeywordsInFilename(filename, adKeywords)
}

// extractAdFileIDs 从广告 URL 列表中提取分享 ID（复用同包 ExtractShareId）。
func extractAdFileIDs(adURLs []string) []string {
	var result []string
	for _, url := range adURLs {
		shareID, _ := ExtractShareId(url)
		if shareID != "" {
			result = append(result, shareID)
		}
	}
	return result
}
