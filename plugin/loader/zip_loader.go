package loader

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// PluginMetadata 插件元数据
type PluginMetadata struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	Dependencies map[string]string `json:"dependencies"`
	EntryPoint   string            `json:"entry_point"`
	ConfigSchema map[string]interface{} `json:"config_schema"`
}

// ZipPluginLoader ZIP格式插件加载器
type ZipPluginLoader struct {
	pluginDir string
	tempDir   string
}

// NewZipPluginLoader 创建ZIP插件加载器
func NewZipPluginLoader(pluginDir string) *ZipPluginLoader {
	tempDir := filepath.Join(os.TempDir(), "urldb-plugins")
	return &ZipPluginLoader{
		pluginDir: pluginDir,
		tempDir:   tempDir,
	}
}

// InstallPlugin 安装ZIP格式的插件
func (z *ZipPluginLoader) InstallPlugin(zipData []byte) (*PluginMetadata, error) {
	// 创建临时目录
	if err := os.MkdirAll(z.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	// 创建临时文件
	tempFile, err := ioutil.TempFile(z.tempDir, "plugin-*.zip")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// 写入ZIP数据
	if _, err := tempFile.Write(zipData); err != nil {
		return nil, fmt.Errorf("failed to write zip data: %v", err)
	}
	tempFile.Close()

	// 读取ZIP文件
	reader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %v", err)
	}
	defer reader.Close()

	// 解析元数据
	metadata, err := z.extractMetadata(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %v", err)
	}

	// 验证插件
	if err := z.validatePlugin(reader, metadata); err != nil {
		return nil, fmt.Errorf("plugin validation failed: %v", err)
	}

	// 提取插件文件
	if err := z.extractPlugin(reader, metadata); err != nil {
		return nil, fmt.Errorf("failed to extract plugin: %v", err)
	}

	utils.Info("成功安装插件: %s v%s", metadata.Name, metadata.Version)
	return metadata, nil
}

// LoadPlugin 加载已安装的插件
func (z *ZipPluginLoader) LoadPlugin(name string) (types.Plugin, error) {
	pluginPath := filepath.Join(z.pluginDir, name, z.getPlatformBinaryName())

	// 检查插件文件是否存在
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	// 加载插件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %v", name, err)
	}

	// 查找插件符号
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin symbol not found in %s: %v", name, err)
	}

	// 类型断言
	pluginInstance, ok := sym.(types.Plugin)
	if !ok {
		return nil, fmt.Errorf("invalid plugin type in %s", name)
	}

	utils.Info("成功加载插件: %s", name)
	return pluginInstance, nil
}

// LoadAllPlugins 加载所有已安装的插件
func (z *ZipPluginLoader) LoadAllPlugins() ([]types.Plugin, error) {
	var plugins []types.Plugin

	// 确保插件目录存在
	if err := os.MkdirAll(z.pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// 读取插件目录
	dirs, err := ioutil.ReadDir(z.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		plugin, err := z.LoadPlugin(dir.Name())
		if err != nil {
			utils.Error("加载插件 %s 失败: %v", dir.Name(), err)
			continue
		}
		plugins = append(plugins, plugin)
	}

	utils.Info("总共加载了 %d 个插件", len(plugins))
	return plugins, nil
}

// UninstallPlugin 卸载插件
func (z *ZipPluginLoader) UninstallPlugin(name string) error {
	pluginDir := filepath.Join(z.pluginDir, name)

	// 检查插件是否存在
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s not found", name)
	}

	// 删除插件目录
	if err := os.RemoveAll(pluginDir); err != nil {
		return fmt.Errorf("failed to remove plugin directory: %v", err)
	}

	utils.Info("成功卸载插件: %s", name)
	return nil
}

// extractMetadata 从ZIP文件中提取元数据
func (z *ZipPluginLoader) extractMetadata(zipReader *zip.ReadCloser) (*PluginMetadata, error) {
	// 查找plugin.json文件
	for _, file := range zipReader.File {
		if file.Name == "plugin.json" {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			data, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			var metadata PluginMetadata
			if err := json.Unmarshal(data, &metadata); err != nil {
				return nil, err
			}

			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("plugin.json not found in zip file")
}

// validatePlugin 验证插件文件
func (z *ZipPluginLoader) validatePlugin(zipReader *zip.ReadCloser, metadata *PluginMetadata) error {
	// 检查必要的文件
	requiredFiles := []string{
		metadata.EntryPoint,
		"plugin.json",
	}

	for _, requiredFile := range requiredFiles {
		found := false
		for _, file := range zipReader.File {
			if file.Name == requiredFile {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("required file not found: %s", requiredFile)
		}
	}

	return nil
}

// extractPlugin 提取插件文件
func (z *ZipPluginLoader) extractPlugin(zipReader *zip.ReadCloser, metadata *PluginMetadata) error {
	pluginDir := filepath.Join(z.pluginDir, metadata.Name)

	// 创建插件目录
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// 提取所有文件
	for _, file := range zipReader.File {
		// 跳过目录
		if file.FileInfo().IsDir() {
			continue
		}

		// 打开源文件
		rc, err := file.Open()
		if err != nil {
			return err
		}

		// 创建目标文件
		targetPath := filepath.Join(pluginDir, file.Name)
		targetDir := filepath.Dir(targetPath)

		// 确保目录存在
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			rc.Close()
			return err
		}

		dst, err := os.Create(targetPath)
		if err != nil {
			rc.Close()
			return err
		}

		// 复制文件内容
		_, err = io.Copy(dst, rc)
		rc.Close()
		dst.Close()

		if err != nil {
			return err
		}

		// 设置文件权限
		if strings.HasSuffix(file.Name, z.getPlatformBinaryName()) {
			os.Chmod(targetPath, 0755)
		}
	}

	return nil
}

// getPlatformBinaryName 获取当前平台的二进制文件名
func (z *ZipPluginLoader) getPlatformBinaryName() string {
	switch runtime.GOOS {
	case "linux":
		return "plugin.so"
	case "windows":
		return "plugin.dll"
	case "darwin":
		return "plugin.dylib"
	default:
		return "plugin"
	}
}

// ListInstalledPlugins 列出已安装的插件
func (z *ZipPluginLoader) ListInstalledPlugins() ([]PluginMetadata, error) {
	var plugins []PluginMetadata

	// 确保插件目录存在
	if err := os.MkdirAll(z.pluginDir, 0755); err != nil {
		return nil, err
	}

	dirs, err := ioutil.ReadDir(z.pluginDir)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		metadataPath := filepath.Join(z.pluginDir, dir.Name(), "plugin.json")
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			continue
		}

		data, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			utils.Error("读取插件元数据失败: %v", err)
			continue
		}

		var metadata PluginMetadata
		if err := json.Unmarshal(data, &metadata); err != nil {
			utils.Error("解析插件元数据失败: %v", err)
			continue
		}

		plugins = append(plugins, metadata)
	}

	return plugins, nil
}