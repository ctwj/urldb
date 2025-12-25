package plugin

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// PluginMetadata 插件元数据
type PluginMetadata struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	License      string            `json:"license"`
	Category     string            `json:"category"`
	Dependencies []string          `json:"dependencies"`
	Permissions  []string          `json:"permissions"`
	Hooks        []string          `json:"hooks"`
	ConfigSchema map[string]interface{} `json:"config_schema"`
	FilePath     string            `json:"file_path"`
	FileSize     int64             `json:"file_size"`
	FileHash     string            `json:"file_hash"`
	Status       string            `json:"status"`
	InstallTime  time.Time         `json:"install_time"`
	LastUpdated  time.Time         `json:"last_updated"`
}

// MetadataParser 元数据解析器
type MetadataParser struct {
	patterns map[string]*regexp.Regexp
}

// NewMetadataParser 创建元数据解析器
func NewMetadataParser() *MetadataParser {
	return &MetadataParser{
		patterns: map[string]*regexp.Regexp{
			"name":         regexp.MustCompile(`@name\s+(\w+)`),
			"version":      regexp.MustCompile(`@version\s+([\d\.]+)`),
			"description":  regexp.MustCompile(`@description\s+(.+)`),
			"author":       regexp.MustCompile(`@author\s+(.+)`),
			"license":      regexp.MustCompile(`@license\s+(.+)`),
			"category":     regexp.MustCompile(`@category\s+(\w+)`),
			"dependencies": regexp.MustCompile(`@dependencies\s+\[(.+)\]`),
			"permissions":  regexp.MustCompile(`@permissions\s+\[(.+)\]`),
			"hooks":        regexp.MustCompile(`@hooks\s+\[(.+)\]`),
			"config":       regexp.MustCompile(`@config_schema\s+\{([\s\S]*?)\}`),
		},
	}
}

// ParseFile 解析插件文件的元数据
func (p *MetadataParser) ParseFile(filePath string) (*PluginMetadata, error) {
	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// 读取文件内容
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 解析元数据
	metadata := &PluginMetadata{
		FilePath:    filePath,
		FileSize:    info.Size(),
		InstallTime: info.ModTime(),
		LastUpdated: info.ModTime(),
		Status:      "installed",
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "//") && !strings.Contains(line, "@") {
			continue
		}

		// 处理配置块
		if strings.Contains(line, "@config_schema") {
			if match := p.patterns["config"].FindStringSubmatch(line); len(match) > 1 {
				metadata.ConfigSchema = parseJSONSchema(match[1])
			}
			continue
		}

		// 解析其他元数据
		for key, pattern := range p.patterns {
			if key == "config" {
				continue
			}
			if match := pattern.FindStringSubmatch(line); len(match) > 1 {
				switch key {
				case "name":
					metadata.Name = match[1]
				case "version":
					metadata.Version = match[1]
				case "description":
					metadata.Description = match[1]
				case "author":
					metadata.Author = match[1]
				case "license":
					metadata.License = match[1]
				case "category":
					metadata.Category = match[1]
				case "dependencies":
					metadata.Dependencies = parseStringArray(match[1])
				case "permissions":
					metadata.Permissions = parseStringArray(match[1])
				case "hooks":
					metadata.Hooks = parseStringArray(match[1])
				}
			}
		}
	}

	// 设置默认值
	if metadata.Name == "" {
		metadata.Name = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}
	if metadata.Version == "" {
		metadata.Version = "1.0.0"
	}
	if metadata.Description == "" {
		metadata.Description = "No description available"
	}
	if metadata.Category == "" {
		metadata.Category = "utility"
	}

	// 计算文件哈希
	hash, err := calculateFileHash(filePath)
	if err == nil {
		metadata.FileHash = hash
	}

	return metadata, scanner.Err()
}

// parseStringArray 解析字符串数组
func parseStringArray(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	// 移除引号和空格
	items := strings.Split(input, ",")
	result := make([]string, 0, len(items))

	for _, item := range items {
		item = strings.TrimSpace(item)
		item = strings.Trim(item, `"'`)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}

// parseJSONSchema 简单的JSON Schema解析
func parseJSONSchema(input string) map[string]interface{} {
	schema := make(map[string]interface{})

	// 简化解析，实际项目中应该使用JSON解析库
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"type"`) {
			if strings.Contains(line, `"object"`) {
				schema["type"] = "object"
			}
		}
	}

	return schema
}

// calculateFileHash 计算文件哈希（简化版本）
func calculateFileHash(filePath string) (string, error) {
	// 这里应该使用实际的哈希算法，如SHA256
	// 为了简化，我们返回文件大小和修改时间作为伪哈希
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("hash_%d_%d", info.Size(), info.ModTime().Unix()), nil
}

// ScanDirectory 扫描目录中的所有插件
func (p *MetadataParser) ScanDirectory(dirPath string) ([]*PluginMetadata, error) {
	var plugins []*PluginMetadata

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .js 文件
		if !info.IsDir() && strings.HasSuffix(path, ".js") {
			metadata, err := p.ParseFile(path)
			if err != nil {
				// 记录错误但继续处理其他文件
				fmt.Printf("Warning: failed to parse %s: %v\n", path, err)
				return nil
			}
			plugins = append(plugins, metadata)
		}

		return nil
	})

	return plugins, err
}

// GetPluginStatus 获取插件状态
func GetPluginStatus(pluginName string) string {
	// 这里应该检查数据库中的插件状态
	// 暂时返回默认状态
	return "installed"
}

// UpdatePluginStatus 更新插件状态
func UpdatePluginStatus(pluginName, status string) error {
	// 这里应该更新数据库中的插件状态
	// 暂时只打印日志
	fmt.Printf("Updating plugin %s status to %s\n", pluginName, status)
	return nil
}

// ValidateMetadata 验证元数据完整性
func (p *MetadataParser) ValidateMetadata(metadata *PluginMetadata) error {
	if metadata.Name == "" {
		return fmt.Errorf("plugin name is required")
	}

	if metadata.Version == "" {
		return fmt.Errorf("plugin version is required")
	}

	if !isValidVersion(metadata.Version) {
		return fmt.Errorf("invalid version format: %s", metadata.Version)
	}

	if metadata.FilePath == "" {
		return fmt.Errorf("file path is required")
	}

	return nil
}

// isValidVersion 验证版本格式
func isValidVersion(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}

	return true
}