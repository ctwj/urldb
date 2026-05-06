package config

import (
	"encoding/json"
	"os"

	"github.com/ctwj/urldb/utils"
)

// FileLoader 本地文件加载器
type FileLoader struct {
	filePath string
}

// NewFileLoader 创建本地文件加载器
func NewFileLoader(filePath string) *FileLoader {
	return &FileLoader{
		filePath: filePath,
	}
}

// Load 从本地文件加载配置
func (l *FileLoader) Load() (*RemoteConfig, error) {
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		return nil, err
	}

	var config RemoteConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	utils.Info("Loaded config from local file: %s, version: %d", l.filePath, config.Version)
	return &config, nil
}

// Exists 检查文件是否存在
func (l *FileLoader) Exists() bool {
	_, err := os.Stat(l.filePath)
	return err == nil
}