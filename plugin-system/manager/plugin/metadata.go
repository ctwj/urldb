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

// ConfigField 配置字段定义
type ConfigField struct {
	Type        string      `json:"type"`        // string, boolean, number, select, text
	Name        string      `json:"name"`        // 字段名
	Label       string      `json:"label"`       // 显示标签
	Description string      `json:"description"` // 描述
	Required    bool        `json:"required"`    // 是否必填
	Default     interface{} `json:"default"`     // 默认值
	Options     []string    `json:"options"`     // 选择项（仅select类型）
	Validation  string      `json:"validation"`  // 验证规则
}

// ScheduledTask 定时任务定义
type ScheduledTask struct {
	Name      string         `json:"name"`       // 任务名称
	Schedule  string         `json:"schedule"`   // 调度表达式
	Line      int            `json:"line"`       // 所在行号
	Frequency *CronFrequency `json:"frequency"`  // 执行频率信息
}

// PluginMetadata 插件元数据
type PluginMetadata struct {
	Name            string                 `json:"name"`
	DisplayName     string                 `json:"display_name"`     // 中文名字
	Version         string                 `json:"version"`
	Description     string                 `json:"description"`
	Author          string                 `json:"author"`
	License         string                 `json:"license"`
	Category        string                 `json:"category"`
	Dependencies    []string               `json:"dependencies"`
	Permissions     []string               `json:"permissions"`
	Hooks           []string               `json:"hooks"`
	ConfigFields    map[string]*ConfigField `json:"config_fields"`
	ConfigSchema    map[string]interface{} `json:"config_schema"`
	ScheduledTasks  []*ScheduledTask       `json:"scheduled_tasks"`  // 定时任务列表
	HasScheduledTask bool                  `json:"has_scheduled_task"` // 是否包含定时任务
	FilePath        string                 `json:"file_path"`
	FileSize        int64                  `json:"file_size"`
	FileHash        string                 `json:"file_hash"`
	Status          string                 `json:"status"`
	InstallTime     time.Time              `json:"install_time"`
	LastUpdated     time.Time              `json:"last_updated"`
}

// MetadataParser 元数据解析器
type MetadataParser struct {
	patterns map[string]*regexp.Regexp
}

// NewMetadataParser 创建元数据解析器
func NewMetadataParser() *MetadataParser {
	return &MetadataParser{
		patterns: map[string]*regexp.Regexp{
			"name":         regexp.MustCompile(`@name\s+([^\s\n]+)`),
			"display_name": regexp.MustCompile(`@display_name\s+(.+)`),
			"version":      regexp.MustCompile(`@version\s+([\d\.]+)`),
			"description":  regexp.MustCompile(`@description\s+(.+)`),
			"author":       regexp.MustCompile(`@author\s+(.+)`),
			"license":      regexp.MustCompile(`@license\s+(.+)`),
			"category":     regexp.MustCompile(`@category\s+([^\s\n]+)`),
			"dependencies": regexp.MustCompile(`@dependencies\s+\[(.+)\]`),
			"permissions":  regexp.MustCompile(`@permissions\s+\[(.+)\]`),
			"hooks":        regexp.MustCompile(`@hooks\s+\[(.+)\]`),
			"config_start": regexp.MustCompile(`@config`),
			"config_end":   regexp.MustCompile(`@config`),
			"field":        regexp.MustCompile(`@field\s+\{([^}]+)\}\s+(\w+)\s+([^"]+)\s+"([^"]*)"`),
			"field_options": regexp.MustCompile(`\[(.*?)\]`),
			"cron_add":     regexp.MustCompile(`cron\.add\s*\(\s*["']([^"']+)["']\s*,\s*["']([^"']+)["']`),
			"cronadd":      regexp.MustCompile(`cronAdd\s*\(\s*["']([^"']+)["']\s*,\s*["']([^"']+)["']`),
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
		FilePath:         filePath,
		FileSize:         info.Size(),
		InstallTime:      info.ModTime(),
		LastUpdated:      info.ModTime(),
		Status:           "installed",
		ConfigFields:     make(map[string]*ConfigField),
		ScheduledTasks:   []*ScheduledTask{},
		HasScheduledTask: false,
	}

	scanner := bufio.NewScanner(file)
	inConfigBlock := false
	var configLines []string
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "//") && !strings.Contains(line, "@") {
			// 即使是注释行也要检查定时任务
			p.checkForScheduledTasks(line, lineNumber, metadata)
			continue
		}

		// 检查定时任务（非注释行）
		p.checkForScheduledTasks(line, lineNumber, metadata)

		// 处理配置块开始
		if strings.Contains(line, "@config") && !inConfigBlock {
			inConfigBlock = true
			continue
		}

		// 处理配置块结束
		if strings.Contains(line, "@config") && inConfigBlock {
			inConfigBlock = false
			// 解析配置块
			metadata.ConfigFields = parseConfigFields(configLines)
			configLines = []string{}
			continue
		}

		// 在配置块内收集行
		if inConfigBlock {
			configLines = append(configLines, line)
			continue
		}

		// 处理旧的配置格式（向后兼容）
		if strings.Contains(line, "@config_schema") {
			if match := p.patterns["config"].FindStringSubmatch(line); len(match) > 1 {
				metadata.ConfigSchema = parseJSONSchema(match[1])
			}
			continue
		}

		// 解析其他元数据
		for key, pattern := range p.patterns {
			if key == "config" || key == "config_start" || key == "config_end" || key == "field" || key == "field_options" {
				continue
			}
			if match := pattern.FindStringSubmatch(line); len(match) > 1 {
				switch key {
				case "name":
					metadata.Name = match[1]
				case "display_name":
					metadata.DisplayName = match[1]
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

	// 设置默认值 - 优先使用文件名（去掉.plugin.js后缀）
	if metadata.Name == "" {
		baseName := filepath.Base(filePath)
		// 去掉 .plugin.js 或 .plugin.ts 后缀
		if strings.HasSuffix(baseName, ".plugin.js") {
			metadata.Name = strings.TrimSuffix(baseName, ".plugin.js")
		} else if strings.HasSuffix(baseName, ".plugin.ts") {
			metadata.Name = strings.TrimSuffix(baseName, ".plugin.ts")
		} else {
			// 去掉普通扩展名作为后备
			metadata.Name = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		}
	}
	if metadata.DisplayName == "" {
		metadata.DisplayName = metadata.Name
		fmt.Printf("DEBUG: 使用默认 display_name: %s (文件: %s)\n", metadata.DisplayName, filepath.Base(filePath))
	} else {
		fmt.Printf("DEBUG: 解析到的 display_name: %s (文件: %s)\n", metadata.DisplayName, filepath.Base(filePath))
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

// parseConfigFields 解析配置字段
func parseConfigFields(lines []string) map[string]*ConfigField {
	fields := make(map[string]*ConfigField)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "@field") {
			continue
		}

		// 解析 @field 行
		// 格式: @field {type} name label "description" ["option1", "option2", ...] @default "value"
		re := regexp.MustCompile(`@field\s+\{([^}]+)\}\s+(\w+)\s+([^"]+)\s+"([^"]*)"(.*)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) < 5 {
			continue
		}

		fieldType := strings.TrimSpace(matches[1])
		fieldName := matches[2]
		fieldLabel := matches[3]
		fieldDesc := matches[4]
		remaining := matches[5]

		field := &ConfigField{
			Type:        fieldType,
			Name:        fieldName,
			Label:       fieldLabel,
			Description: fieldDesc,
			Required:    true,  // 默认必填
		}

		// 检查是否有 @optional 标记
		if strings.Contains(remaining, "@optional") {
			field.Required = false
		}

		// 解析默认值
		defaultValue := parseDefaultValue(remaining, fieldType)
		if defaultValue != nil {
			field.Default = defaultValue
		} else {
			// 设置默认值（如果没有 @default 注释）
			switch fieldType {
			case "boolean":
				field.Default = false
			case "number":
				field.Default = 0
			case "string", "text":
				field.Default = ""
			case "select":
				field.Default = ""
				// 解析选项
				if options := parseFieldOptions(remaining); len(options) > 0 {
					field.Options = options
					field.Default = options[0] // 默认第一个选项
				}
			}
		}

		// 对于 select 类型，无论如何都要解析选项
		if fieldType == "select" {
			if options := parseFieldOptions(remaining); len(options) > 0 {
				field.Options = options
				// 如果没有设置默认值，使用第一个选项
				if defaultValue == nil {
					field.Default = options[0]
				}
			}
		}

		fields[fieldName] = field
	}

	return fields
}

// parseDefaultValue 解析默认值
func parseDefaultValue(input string, fieldType string) interface{} {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var defaultValue string

	// 先尝试匹配带引号的值 @default "value"
	reQuoted := regexp.MustCompile(`@default\s+"([^"]*)"`)
	matches := reQuoted.FindStringSubmatch(input)
	if len(matches) >= 2 {
		defaultValue = strings.TrimSpace(matches[1])
	} else {
		// 尝试匹配不带引号的值 @default value
		reUnquoted := regexp.MustCompile(`@default\s+(\S+)`)
		matches := reUnquoted.FindStringSubmatch(input)
		if len(matches) >= 2 {
			defaultValue = strings.TrimSpace(matches[1])
		} else {
			return nil
		}
	}

	if defaultValue == "" {
		return nil
	}

	// 根据字段类型转换默认值
	switch fieldType {
	case "boolean":
		// 处理布尔值
		lowerValue := strings.ToLower(defaultValue)
		if lowerValue == "true" || lowerValue == "1" || lowerValue == "yes" || lowerValue == "on" {
			return true
		}
		return false
	case "number":
		// 处理数字
		if intValue, err := strconv.Atoi(defaultValue); err == nil {
			return intValue
		}
		if floatValue, err := strconv.ParseFloat(defaultValue, 64); err == nil {
			return floatValue
		}
		return 0
	default:
		// 字符串、文本、选择框等直接返回字符串
		return defaultValue
	}
}

// parseFieldOptions 解析字段选项
func parseFieldOptions(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	// 查找方括号内的内容
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return []string{}
	}

	optionsStr := strings.TrimSpace(matches[1])
	if optionsStr == "" {
		return []string{}
	}

	// 解析选项
	options := strings.Split(optionsStr, ",")
	result := make([]string, 0, len(options))

	for _, option := range options {
		option = strings.TrimSpace(option)
		option = strings.Trim(option, `"`)
		if option != "" {
			result = append(result, option)
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
	nameMap := make(map[string]string) // 插件名称到文件路径的映射

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

			// 检查名称冲突
			if existingPath, exists := nameMap[metadata.Name]; exists {
				fmt.Printf("ERROR: Plugin name conflict detected!\n")
				fmt.Printf("  Plugin name '%s' is already used by: %s\n", metadata.Name, existingPath)
				fmt.Printf("  Conflicting file: %s\n", path)
				fmt.Printf("  Second plugin will be skipped to prevent conflicts.\n")
				return nil // 跳过这个插件
			}

			nameMap[metadata.Name] = path
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

// checkForScheduledTasks 检查定时任务
func (p *MetadataParser) checkForScheduledTasks(line string, lineNumber int, metadata *PluginMetadata) {
	// 检查 cron.add 调用
	if match := p.patterns["cron_add"].FindStringSubmatch(line); len(match) > 2 {
		taskName := match[1]
		schedule := match[2]
		frequency := ParseCronExpression(schedule)
		task := &ScheduledTask{
			Name:      taskName,
			Schedule:  schedule,
			Line:      lineNumber,
			Frequency: frequency,
		}
		metadata.ScheduledTasks = append(metadata.ScheduledTasks, task)
		metadata.HasScheduledTask = true
	}

	// 检查 cronAdd 调用
	if match := p.patterns["cronadd"].FindStringSubmatch(line); len(match) > 2 {
		taskName := match[1]
		schedule := match[2]
		frequency := ParseCronExpression(schedule)
		task := &ScheduledTask{
			Name:      taskName,
			Schedule:  schedule,
			Line:      lineNumber,
			Frequency: frequency,
		}
		metadata.ScheduledTasks = append(metadata.ScheduledTasks, task)
		metadata.HasScheduledTask = true
	}
}

// ParseContent 从内容解析插件元数据
func (p *MetadataParser) ParseContent(content []byte) (*PluginMetadata, error) {
	// 创建基本的元数据对象
	metadata := &PluginMetadata{
		FilePath:         "",
		FileSize:         int64(len(content)),
		InstallTime:      time.Now(),
		LastUpdated:      time.Now(),
		Status:           "installed",
		ConfigFields:     make(map[string]*ConfigField),
		ScheduledTasks:   []*ScheduledTask{},
		HasScheduledTask: false,
	}

	// 将内容转换为字符串并按行分割
	lines := strings.Split(string(content), "\n")
	lineNumber := 0

	for _, line := range lines {
		lineNumber++
		line = strings.TrimSpace(line)

		// 解析元数据注释
		if strings.HasPrefix(line, "/**") {
			// 开始解析元数据块
			continue
		}

		if strings.HasPrefix(line, "* @name ") {
			metadata.Name = strings.TrimSpace(strings.TrimPrefix(line, "* @name "))
		} else if strings.HasPrefix(line, "* @description ") {
			metadata.Description = strings.TrimSpace(strings.TrimPrefix(line, "* @description "))
		} else if strings.HasPrefix(line, "* @version ") {
			metadata.Version = strings.TrimSpace(strings.TrimPrefix(line, "* @version "))
		} else if strings.HasPrefix(line, "* @author ") {
			metadata.Author = strings.TrimSpace(strings.TrimPrefix(line, "* @author "))
		} else if strings.HasPrefix(line, "* @license ") {
			metadata.License = strings.TrimSpace(strings.TrimPrefix(line, "* @license "))
		} else if strings.HasPrefix(line, "* @category ") {
			metadata.Category = strings.TrimSpace(strings.TrimPrefix(line, "* @category "))
		}

		// 检查定时任务
		p.checkForScheduledTasks(line, lineNumber, metadata)
	}

	// 设置默认值
	if metadata.Name == "" {
		metadata.Name = "unknown"
	}
	if metadata.Version == "" {
		metadata.Version = "1.0.0"
	}
	if metadata.Description == "" {
		metadata.Description = "No description available"
	}

	return metadata, nil
}