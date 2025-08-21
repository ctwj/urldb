package utils

import (
	"regexp"
	"strings"

	"github.com/ctwj/urldb/db/entity"
)

// ForbiddenWordsProcessor 违禁词处理器
type ForbiddenWordsProcessor struct{}

// NewForbiddenWordsProcessor 创建违禁词处理器实例
func NewForbiddenWordsProcessor() *ForbiddenWordsProcessor {
	return &ForbiddenWordsProcessor{}
}

// CheckContainsForbiddenWords 检查字符串是否包含违禁词
// 参数：
//   - text: 要检查的文本
//   - forbiddenWords: 违禁词列表
//
// 返回：
//   - bool: 是否包含违禁词
//   - []string: 匹配到的违禁词列表
func (p *ForbiddenWordsProcessor) CheckContainsForbiddenWords(text string, forbiddenWords []string) (bool, []string) {
	if len(forbiddenWords) == 0 {
		return false, nil
	}

	var matchedWords []string
	textLower := strings.ToLower(text)

	for _, word := range forbiddenWords {
		wordLower := strings.ToLower(word)
		if strings.Contains(textLower, wordLower) {
			matchedWords = append(matchedWords, word)
		}
	}

	return len(matchedWords) > 0, matchedWords
}

// ReplaceForbiddenWords 替换字符串中的违禁词为 *
// 参数：
//   - text: 要处理的文本
//   - forbiddenWords: 违禁词列表
//
// 返回：
//   - string: 替换后的文本
func (p *ForbiddenWordsProcessor) ReplaceForbiddenWords(text string, forbiddenWords []string) string {
	if len(forbiddenWords) == 0 {
		return text
	}

	result := text
	// 按长度降序排序，避免短词替换后影响长词的匹配
	sortedWords := make([]string, len(forbiddenWords))
	copy(sortedWords, forbiddenWords)

	// 简单的长度排序（这里可以优化为更复杂的排序）
	for i := 0; i < len(sortedWords)-1; i++ {
		for j := i + 1; j < len(sortedWords); j++ {
			if len(sortedWords[i]) < len(sortedWords[j]) {
				sortedWords[i], sortedWords[j] = sortedWords[j], sortedWords[i]
			}
		}
	}

	for _, word := range sortedWords {
		// 使用正则表达式进行不区分大小写的替换
		// 对于中文，不使用单词边界，直接替换
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
		// 使用字符长度而不是字节长度
		charCount := len([]rune(word))
		result = re.ReplaceAllString(result, strings.Repeat("*", charCount))
	}

	return result
}

// ReplaceForbiddenWordsWithHighlight 替换字符串中的违禁词为 *（处理高亮标记）
// 参数：
//   - text: 要处理的文本（可能包含高亮标记）
//   - forbiddenWords: 违禁词列表
//
// 返回：
//   - string: 替换后的文本
func (p *ForbiddenWordsProcessor) ReplaceForbiddenWordsWithHighlight(text string, forbiddenWords []string) string {
	if len(forbiddenWords) == 0 {
		return text
	}

	// 1. 先移除所有高亮标记，获取纯文本
	cleanText := regexp.MustCompile(`<mark>(.*?)</mark>`).ReplaceAllString(text, "$1")

	// 2. 检查纯文本中是否包含违禁词
	hasForbidden := false
	for _, word := range forbiddenWords {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
		if re.MatchString(cleanText) {
			hasForbidden = true
			break
		}
	}

	// 3. 如果包含违禁词，则替换非高亮文本
	if hasForbidden {
		return p.ReplaceForbiddenWords(text, forbiddenWords)
	}

	// 4. 如果不包含违禁词，直接返回原文本
	return text
}

// ProcessForbiddenWords 处理违禁词：检查并替换
// 参数：
//   - text: 要处理的文本
//   - forbiddenWords: 违禁词列表
//
// 返回：
//   - bool: 是否包含违禁词
//   - []string: 匹配到的违禁词列表
//   - string: 替换后的文本
func (p *ForbiddenWordsProcessor) ProcessForbiddenWords(text string, forbiddenWords []string) (bool, []string, string) {
	contains, matchedWords := p.CheckContainsForbiddenWords(text, forbiddenWords)
	replacedText := p.ReplaceForbiddenWords(text, forbiddenWords)
	return contains, matchedWords, replacedText
}

// ParseForbiddenWordsConfig 解析违禁词配置字符串
// 参数：
//   - config: 违禁词配置字符串，多个词用逗号分隔
//
// 返回：
//   - []string: 处理后的违禁词列表
func (p *ForbiddenWordsProcessor) ParseForbiddenWordsConfig(config string) []string {
	if config == "" {
		return nil
	}

	words := strings.Split(config, ",")
	var cleanWords []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			cleanWords = append(cleanWords, word)
		}
	}

	return cleanWords
}

// 全局实例，方便直接调用
var DefaultForbiddenWordsProcessor = NewForbiddenWordsProcessor()

// 便捷函数，直接调用全局实例

// CheckContainsForbiddenWords 检查字符串是否包含违禁词（便捷函数）
func CheckContainsForbiddenWords(text string, forbiddenWords []string) (bool, []string) {
	return DefaultForbiddenWordsProcessor.CheckContainsForbiddenWords(text, forbiddenWords)
}

// ReplaceForbiddenWords 替换字符串中的违禁词为 *（便捷函数）
func ReplaceForbiddenWords(text string, forbiddenWords []string) string {
	return DefaultForbiddenWordsProcessor.ReplaceForbiddenWords(text, forbiddenWords)
}

// ReplaceForbiddenWordsWithHighlight 替换字符串中的违禁词为 *（处理高亮标记，便捷函数）
func ReplaceForbiddenWordsWithHighlight(text string, forbiddenWords []string) string {
	return DefaultForbiddenWordsProcessor.ReplaceForbiddenWordsWithHighlight(text, forbiddenWords)
}

// ProcessForbiddenWords 处理违禁词：检查并替换（便捷函数）
func ProcessForbiddenWords(text string, forbiddenWords []string) (bool, []string, string) {
	return DefaultForbiddenWordsProcessor.ProcessForbiddenWords(text, forbiddenWords)
}

// ParseForbiddenWordsConfig 解析违禁词配置字符串（便捷函数）
func ParseForbiddenWordsConfig(config string) []string {
	return DefaultForbiddenWordsProcessor.ParseForbiddenWordsConfig(config)
}

// RemoveDuplicates 去除字符串切片中的重复项
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// ResourceForbiddenInfo 资源违禁词信息
type ResourceForbiddenInfo struct {
	HasForbiddenWords bool     `json:"has_forbidden_words"`
	ForbiddenWords    []string `json:"forbidden_words"`
	ProcessedTitle    string   `json:"-"` // 不序列化，仅内部使用
	ProcessedDesc     string   `json:"-"` // 不序列化，仅内部使用
}

// CheckResourceForbiddenWords 检查资源是否包含违禁词（检查标题和描述）
// 参数：
//   - title: 资源标题
//   - description: 资源描述
//   - forbiddenWords: 违禁词列表
//
// 返回：
//   - ResourceForbiddenInfo: 包含检查结果和处理后的文本
func CheckResourceForbiddenWords(title, description string, forbiddenWords []string) ResourceForbiddenInfo {

	if len(forbiddenWords) == 0 {
		return ResourceForbiddenInfo{
			HasForbiddenWords: false,
			ForbiddenWords:    []string{},
			ProcessedTitle:    title,
			ProcessedDesc:     description,
		}
	}

	// 分别检查标题和描述
	titleHasForbidden, titleMatchedWords := CheckContainsForbiddenWords(title, forbiddenWords)
	descHasForbidden, descMatchedWords := CheckContainsForbiddenWords(description, forbiddenWords)

	// 合并结果
	hasForbiddenWords := titleHasForbidden || descHasForbidden
	var matchedWords []string
	if titleHasForbidden {
		matchedWords = append(matchedWords, titleMatchedWords...)
	}
	if descHasForbidden {
		matchedWords = append(matchedWords, descMatchedWords...)
	}
	// 去重
	matchedWords = RemoveDuplicates(matchedWords)

	// 处理文本（替换违禁词）
	processedTitle := ReplaceForbiddenWords(title, forbiddenWords)
	processedDesc := ReplaceForbiddenWords(description, forbiddenWords)

	return ResourceForbiddenInfo{
		HasForbiddenWords: hasForbiddenWords,
		ForbiddenWords:    matchedWords,
		ProcessedTitle:    processedTitle,
		ProcessedDesc:     processedDesc,
	}
}

// GetForbiddenWordsFromConfig 从系统配置获取违禁词列表
// 参数：
//   - getConfigFunc: 获取配置的函数
//
// 返回：
//   - []string: 解析后的违禁词列表
//   - error: 获取配置时的错误
func GetForbiddenWordsFromConfig(getConfigFunc func() (string, error)) ([]string, error) {
	forbiddenWords, err := getConfigFunc()
	if err != nil {
		return nil, err
	}
	return ParseForbiddenWordsConfig(forbiddenWords), nil
}

// ProcessResourcesForbiddenWords 批量处理资源的违禁词
// 参数：
//   - resources: 资源切片
//   - forbiddenWords: 违禁词列表
//
// 返回：
//   - 处理后的资源切片
func ProcessResourcesForbiddenWords(resources []entity.Resource, forbiddenWords []string) []entity.Resource {
	if len(forbiddenWords) == 0 {
		return resources
	}

	for i := range resources {
		// 处理标题中的违禁词
		resources[i].Title = ReplaceForbiddenWords(resources[i].Title, forbiddenWords)
		// 处理描述中的违禁词
		resources[i].Description = ReplaceForbiddenWords(resources[i].Description, forbiddenWords)
	}

	return resources
}
