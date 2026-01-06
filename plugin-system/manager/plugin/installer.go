package plugin

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ctwj/urldb/utils"
)

// PluginPackage 插件包配置文件
type PluginPackage struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Main        string            `json:"main"`        // 主入口文件
	Hooks       []string          `json:"hooks"`       // 钩子文件列表
	Config      map[string]interface{} `json:"config"`  // 默认配置
	Dependencies []string         `json:"dependencies"` // 依赖插件
}

// PluginInstaller 插件安装器
type PluginInstaller struct {
	pluginsDir   string
	installedDir string
	tempDir      string
	db           *sql.DB // 数据库连接，用于执行迁移
}

// getMin 返回两个整数中的较小值
func getMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NewPluginInstaller 创建插件安装器
func NewPluginInstaller(baseDir string) *PluginInstaller {
	pluginsDir := filepath.Join(baseDir, "plugins")
	installedDir := filepath.Join(pluginsDir, "installed")
	tempDir := filepath.Join(pluginsDir, "temp")

	return &PluginInstaller{
		pluginsDir:   pluginsDir,
		installedDir: installedDir,
		tempDir:      tempDir,
	}
}

// NewPluginInstallerWithDB 创建带数据库连接的插件安装器
func NewPluginInstallerWithDB(baseDir string, db *sql.DB) *PluginInstaller {
	installer := NewPluginInstaller(baseDir)
	installer.db = db
	return installer
}

// ensureDirectories 确保目录存在
func (pi *PluginInstaller) ensureDirectories() error {
	dirs := []string{pi.pluginsDir, pi.installedDir, pi.tempDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// InstallFromFile 从文件安装插件
func (pi *PluginInstaller) InstallFromFile(zipPath string) error {
	if err := pi.ensureDirectories(); err != nil {
		return err
	}

	// 解压到临时目录
	tempPluginDir, err := pi.extractToTemp(zipPath)
	if err != nil {
		return fmt.Errorf("failed to extract plugin: %w", err)
	}

	// 读取插件配置
	pkg, err := pi.readPluginConfig(tempPluginDir)
	if err != nil {
		pi.cleanupTemp(tempPluginDir)
		return fmt.Errorf("failed to read plugin config: %w", err)
	}

	// 验证插件
	if err := pi.validatePlugin(pkg, tempPluginDir); err != nil {
		pi.cleanupTemp(tempPluginDir)
		return fmt.Errorf("plugin validation failed: %w", err)
	}

	// 安装到目标目录
	installDir := filepath.Join(pi.installedDir, pkg.Name)
	if err := pi.installToDirectory(tempPluginDir, installDir); err != nil {
		pi.cleanupTemp(tempPluginDir)
		return fmt.Errorf("failed to install plugin: %w", err)
	}

	// 执行安装迁移
	if err := pi.executeInstallMigration(installDir); err != nil {
		pi.cleanupTemp(tempPluginDir)
		return fmt.Errorf("failed to execute install migration: %w", err)
	}

	// 清理临时文件
	pi.cleanupTemp(tempPluginDir)

	utils.Info("Plugin '%s' v%s installed successfully", pkg.Name, pkg.Version)
	return nil
}

// InstallFromURL 从URL安装插件
func (pi *PluginInstaller) InstallFromURL(url string) error {
	// 检查URL是否指向JS文件
	isJSFile := strings.HasSuffix(strings.ToLower(url), ".plugin.js")
	isZipFile := strings.HasSuffix(strings.ToLower(url), ".zip")

	// 下载插件文件
	downloadedPath, err := pi.downloadPlugin(url)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}
	defer os.Remove(downloadedPath)

	// 检查下载的文件类型并进行相应处理
	if isJSFile || strings.HasSuffix(strings.ToLower(downloadedPath), ".plugin.js") {
		// 这是一个JS单文件插件，为了与您提到的手动放置在hooks目录的行为一致，
		// 我们直接将其复制到hooks目录下
		return pi.installJSSingleFileToHooks(downloadedPath)
	} else if isZipFile || strings.HasSuffix(strings.ToLower(downloadedPath), ".zip") {
		// 这是一个ZIP压缩包，使用InstallFromFile方法
		return pi.InstallFromFile(downloadedPath)
	} else {
		// 尝试检测文件类型
		content, err := os.ReadFile(downloadedPath)
		if err != nil {
			return fmt.Errorf("failed to read downloaded file: %w", err)
		}
		contentStr := string(content)

		// 检查是否是JavaScript插件文件（通常包含plugin.js标识符或JSDoc注释）
		if strings.Contains(contentStr, "@name") || strings.Contains(contentStr, "onURLAdd") || strings.Contains(contentStr, "onUserLogin") || strings.Contains(contentStr, "onURLAccess") || strings.Contains(contentStr, "routerAdd") || strings.Contains(contentStr, "cronAdd") {
			// 假设是插件JS文件，直接复制到hooks目录
			return pi.installJSSingleFileToHooks(downloadedPath)
		} else {
			// 默认作为ZIP文件处理
			return pi.InstallFromFile(downloadedPath)
		}
	}
}

// installJSSingleFileToHooks 安装JS单文件到hooks目录
func (pi *PluginInstaller) installJSSingleFileToHooks(filePath string) error {
	utils.Info("Installing JS single file to hooks: %s", filePath)

	// 读取插件文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		utils.Error("Failed to read plugin file: %v", err)
		return fmt.Errorf("failed to read plugin file: %w", err)
	}

	utils.Info("Downloaded file size: %d bytes", len(content))

	// 解析插件元数据
	metadataParser := NewMetadataParser()
	metadata, err := metadataParser.ParseContent(content)
	if err != nil {
		utils.Error("Failed to parse plugin metadata: %v", err)
		utils.Info("File content preview: %s", string(content[:getMin(len(content), 200)]))
		return fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	utils.Info("Parsed plugin metadata - Name: %s, Version: %s, Description: %s",
		metadata.Name, metadata.Version, metadata.Description)

	// 确保hooks目录存在
	hooksDir := "./plugin-system/hooks"
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		utils.Error("Failed to create hooks directory: %v", err)
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// 从元数据中获取插件名称，如果获取不到则从文件名推断
	pluginName := metadata.Name
	if pluginName == "" {
		// 从文件名推断插件名
		baseName := filepath.Base(filePath)
		pluginName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		utils.Info("Using filename-based plugin name: %s", pluginName)
	} else {
		utils.Info("Using metadata-based plugin name: %s", pluginName)
	}

	// 确保插件名称是有效的文件名
	origPluginName := pluginName
	pluginName = strings.ReplaceAll(pluginName, " ", "_")  // 替换空格为下划线
	pluginName = strings.ReplaceAll(pluginName, "/", "_")  // 替换斜杠为下划线
	pluginName = strings.ReplaceAll(pluginName, "\\", "_") // 替换反斜杠为下划线
	pluginName = strings.ReplaceAll(pluginName, ":", "_")  // 替换冒号为下划线
	pluginName = strings.ReplaceAll(pluginName, "*", "_")  // 替换星号为下划线
	pluginName = strings.ReplaceAll(pluginName, "?", "_")  // 替换问号为下划线
	pluginName = strings.ReplaceAll(pluginName, "\"", "_") // 替换引号为下划线
	pluginName = strings.ReplaceAll(pluginName, "<", "_")  // 替换小于号为下划线
	pluginName = strings.ReplaceAll(pluginName, ">", "_")  // 替换大于号为下划线
	pluginName = strings.ReplaceAll(pluginName, "|", "_")  // 替换竖线为下划线

	if origPluginName != pluginName {
		utils.Info("Plugin name sanitized from '%s' to '%s'", origPluginName, pluginName)
	}

	// 创建目标文件路径
	destFile := filepath.Join(hooksDir, pluginName+".plugin.js")
	utils.Info("Target file path: %s", destFile)

	// 检查文件是否已存在
	if _, err := os.Stat(destFile); err == nil {
		utils.Warn("Plugin file already exists: %s", destFile)
		return fmt.Errorf("plugin file already exists: %s", destFile)
	}

	// 复制插件文件到hooks目录
	if err := os.WriteFile(destFile, content, 0644); err != nil {
		utils.Error("Failed to copy plugin file to hooks directory: %v", err)
		return fmt.Errorf("failed to copy plugin file to hooks directory: %w", err)
	}

	utils.Info("JS file plugin '%s' installed to hooks directory successfully", pluginName)
	return nil
}

// extractToTemp 解压插件包到临时目录
func (pi *PluginInstaller) extractToTemp(zipPath string) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// 创建临时目录
	tempDir, err := os.MkdirTemp(pi.tempDir, "plugin-*")
	if err != nil {
		return "", err
	}

	// 解压文件
	for _, file := range reader.File {
		path := filepath.Join(tempDir, file.Name)

		// 防止路径遍历攻击
		if !strings.HasPrefix(path, tempDir) {
			continue
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return tempDir, err
		}

		// 解压文件
		fileReader, err := file.Open()
		if err != nil {
			return tempDir, err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			fileReader.Close()
			return tempDir, err
		}

		_, err = io.Copy(targetFile, fileReader)
		fileReader.Close()
		targetFile.Close()

		if err != nil {
			return tempDir, err
		}
	}

	return tempDir, nil
}

// readPluginConfig 读取插件配置
func (pi *PluginInstaller) readPluginConfig(pluginDir string) (*PluginPackage, error) {
	configPath := filepath.Join(pluginDir, "package.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var pkg PluginPackage
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	// 设置默认主入口文件
	if pkg.Main == "" {
		pkg.Main = "index.js"
	}

	return &pkg, nil
}

// validatePlugin 验证插件
func (pi *PluginInstaller) validatePlugin(pkg *PluginPackage, pluginDir string) error {
	// 检查插件名称格式
	if !regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(pkg.Name) {
		return fmt.Errorf("invalid plugin name: %s (only lowercase letters, numbers, hyphens and underscores allowed)", pkg.Name)
	}

	// 检查主入口文件是否存在
	mainPath := filepath.Join(pluginDir, pkg.Main)
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		return fmt.Errorf("main entry file not found: %s", pkg.Main)
	}

	// 检查钩子文件是否存在
	for _, hook := range pkg.Hooks {
		hookPath := filepath.Join(pluginDir, hook)
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			return fmt.Errorf("hook file not found: %s", hook)
		}
	}

	// 检查是否已安装
	installDir := filepath.Join(pi.installedDir, pkg.Name)
	if _, err := os.Stat(installDir); !os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' is already installed", pkg.Name)
	}

	return nil
}

// installToDirectory 安装插件到目标目录
func (pi *PluginInstaller) installToDirectory(srcDir, destDir string) error {
	// 复制整个插件目录
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// 复制文件
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

// downloadPlugin 下载插件包
func (pi *PluginInstaller) downloadPlugin(url string) (string, error) {
	utils.Info("Starting download from URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		utils.Error("Failed to make HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.Error("Download failed with status: %d", resp.StatusCode)
		return "", fmt.Errorf("failed to download plugin: HTTP %d", resp.StatusCode)
	}

	// 检查内容类型
	contentType := resp.Header.Get("Content-Type")
	utils.Info("Downloaded content type: %s", contentType)

	// 根据URL后缀创建相应类型的临时文件
	var tempFile *os.File
	if strings.HasSuffix(strings.ToLower(url), ".plugin.js") {
		tempFile, err = os.CreateTemp(pi.tempDir, "download-*.js")
	} else {
		tempFile, err = os.CreateTemp(pi.tempDir, "download-*")
	}

	if err != nil {
		utils.Error("Failed to create temp file: %v", err)
		return "", err
	}
	defer tempFile.Close()

	// 保存文件
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		utils.Error("Failed to save downloaded file: %v", err)
		os.Remove(tempFile.Name())
		return "", err
	}

	utils.Info("File downloaded successfully to: %s", tempFile.Name())
	return tempFile.Name(), nil
}

// cleanupTemp 清理临时目录
func (pi *PluginInstaller) cleanupTemp(tempDir string) {
	os.RemoveAll(tempDir)
}

// Uninstall 卸载插件
func (pi *PluginInstaller) Uninstall(pluginName string) error {
	installDir := filepath.Join(pi.installedDir, pluginName)

	// 检查插件是否存在
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' is not installed", pluginName)
	}

	// 执行卸载迁移
	if err := pi.executeUninstallMigration(installDir); err != nil {
		return fmt.Errorf("failed to execute uninstall migration: %w", err)
	}

	// 删除插件目录
	if err := os.RemoveAll(installDir); err != nil {
		return fmt.Errorf("failed to uninstall plugin '%s': %w", pluginName, err)
	}

	utils.Info("Plugin '%s' uninstalled successfully", pluginName)
	return nil
}

// ListInstalled 列出已安装的插件
func (pi *PluginInstaller) ListInstalled() ([]*PluginPackage, error) {
	entries, err := os.ReadDir(pi.installedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*PluginPackage{}, nil
		}
		return nil, err
	}

	var plugins []*PluginPackage

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginDir := filepath.Join(pi.installedDir, entry.Name())
		pkg, err := pi.readPluginConfig(pluginDir)
		if err != nil {
			utils.Error("Failed to read plugin config for '%s': %v", entry.Name(), err)
			continue
		}

		plugins = append(plugins, pkg)
	}

	return plugins, nil
}

// IsInstalled 检查插件是否已安装
func (pi *PluginInstaller) IsInstalled(pluginName string) bool {
	installDir := filepath.Join(pi.installedDir, pluginName)
	_, err := os.Stat(installDir)
	return !os.IsNotExist(err)
}

// InstallSingleFile 安装单文件插件
func (pi *PluginInstaller) InstallSingleFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("plugin file not found: %s", filePath)
	}

	// 读取插件文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read plugin file: %w", err)
	}

	// 解析插件元数据
	metadataParser := NewMetadataParser()
	metadata, err := metadataParser.ParseContent(content)
	if err != nil {
		return fmt.Errorf("failed to parse plugin metadata: %w", err)
	}

	// 创建插件包对象
	pkg := &PluginPackage{
		Name:        metadata.Name,
		Version:     metadata.Version,
		Description: metadata.Description,
		Author:      metadata.Author,
		Main:        filepath.Base(filePath),
		Hooks:       []string{filepath.Base(filePath)},
	}

	// 创建安装目录
	installDir := filepath.Join(pi.installedDir, pkg.Name)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// 创建hooks子目录
	hooksDir := filepath.Join(installDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// 复制插件文件到hooks目录
	destFile := filepath.Join(hooksDir, filepath.Base(filePath))
	if err := os.WriteFile(destFile, content, 0644); err != nil {
		return fmt.Errorf("failed to copy plugin file: %w", err)
	}

	// 创建package.json
	packageJSON := map[string]interface{}{
		"name":        pkg.Name,
		"version":     pkg.Version,
		"description": pkg.Description,
		"author":      pkg.Author,
		"main":        pkg.Main,
		"hooks":       pkg.Hooks,
		"config":      map[string]interface{}{"enabled": true},
	}

	packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create package.json: %w", err)
	}

	packageJSONPath := filepath.Join(installDir, "package.json")
	if err := os.WriteFile(packageJSONPath, packageJSONBytes, 0644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	utils.Info("Single-file plugin '%s' v%s installed successfully", pkg.Name, pkg.Version)
	return nil
}

// executeMigration 执行迁移 SQL 文件
func (pi *PluginInstaller) executeMigration(pluginDir, migrationType string) error {
	utils.Info("Starting %s migration execution for plugin: %s", migrationType, pluginDir)

	if pi.db == nil {
		utils.Warn("Database connection not available, skipping migration")
		return nil
	}

	migrationFile := filepath.Join(pluginDir, "migrate", migrationType+".sql")
	utils.Info("Checking migration file: %s", migrationFile)

	if _, err := os.Stat(migrationFile); os.IsNotExist(err) {
		// 迁移文件不存在，不是错误
		utils.Info("Migration file does not exist: %s", migrationFile)
		return nil
	}

	utils.Info("Executing %s migration for plugin: %s", migrationType, pluginDir)

	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
	}

	utils.Info("Migration SQL content: %s", string(content))

	// 执行 SQL
	_, err = pi.db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute %s migration: %w", migrationType, err)
	}

	utils.Info("Successfully executed %s migration for plugin", migrationType)
	return nil
}

// executeInstallMigration 执行安装迁移
func (pi *PluginInstaller) executeInstallMigration(pluginDir string) error {
	return pi.executeMigration(pluginDir, "install")
}

// executeUninstallMigration 执行卸载迁移
func (pi *PluginInstaller) executeUninstallMigration(pluginDir string) error {
	return pi.executeMigration(pluginDir, "uninstall")
}