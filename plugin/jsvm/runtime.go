package jsvm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ctwj/urldb/core"
)

const defaultScriptPath = "pb.js"

// Config JSVM 配置（简化版，不依赖 goja）
type Config struct {
	// OnInit VM 初始化回调
	OnInit func(vm interface{})

	// HooksWatch 启用文件监控热重载
	HooksWatch bool

	// HooksDir 钩子文件目录
	HooksDir string

	// HooksFilesPattern 钩子文件匹配模式
	HooksFilesPattern string

	// HooksPoolSize VM 池大小
	HooksPoolSize int

	// MigrationsDir 迁移文件目录
	MigrationsDir string

	// MigrationsFilesPattern 迁移文件匹配模式
	MigrationsFilesPattern string

	// TypesDir TypeScript 类型定义目录
	TypesDir string
}

// plugin JSVM 插件实例（简化版）
type plugin struct {
	app    core.App
	config Config
}

// Register 注册 JSVM 插件（简化版实现）
func Register(app core.App, config Config) error {
	p := &plugin{
		app:    app,
		config: config,
	}

	// 设置默认值
	if p.config.HooksDir == "" {
		p.config.HooksDir = filepath.Join(app.DataDir(), "../hooks")
	}

	if p.config.MigrationsDir == "" {
		p.config.MigrationsDir = filepath.Join(app.DataDir(), "../migrations")
	}

	if p.config.HooksFilesPattern == "" {
		p.config.HooksFilesPattern = `^.*(\.pb\.js|\.pb\.ts)$`
	}

	if p.config.MigrationsFilesPattern == "" {
		p.config.MigrationsFilesPattern = `^.*(\.js|\.ts)$`
	}

	if p.config.TypesDir == "" {
		p.config.TypesDir = app.DataDir()
	}

	// 注册钩子
	p.app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		err := e.Next()
		if err != nil {
			return err
		}

		// 刷新类型文件
		err = p.refreshTypesFile()
		if err != nil {
			fmt.Printf("Unable to refresh app types file: %v\n", err)
		}

		return nil
	})

	// 注册迁移和钩子（简化版）
	if err := p.registerMigrations(); err != nil {
		return fmt.Errorf("registerMigrations: %w", err)
	}

	if err := p.registerHooks(); err != nil {
		return fmt.Errorf("registerHooks: %w", err)
	}

	return nil
}

// MustRegister 注册插件，失败时 panic
func MustRegister(app core.App, config Config) {
	if err := Register(app, config); err != nil {
		panic(err)
	}
}

// registerMigrations 注册 JavaScript 迁移（简化版）
func (p *plugin) registerMigrations() error {
	files, err := p.filesContent(p.config.MigrationsDir, p.config.MigrationsFilesPattern)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf("Found %d migration files\n", len(files))

	// 简化版：只记录文件，不执行 JavaScript
	for file := range files {
		fmt.Printf("Migration file found: %s\n", file)
	}

	return nil
}

// registerHooks 注册 JavaScript 钩子（简化版）
func (p *plugin) registerHooks() error {
	files, err := p.filesContent(p.config.HooksDir, p.config.HooksFilesPattern)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf("Found %d hook files\n", len(files))

	// 简化版：只记录文件，不执行 JavaScript
	for file := range files {
		fmt.Printf("Hook file found: %s\n", file)
	}

	return nil
}

// filesContent 获取目录下匹配模式的文件内容
func (p *plugin) filesContent(dirPath, pattern string) (map[string][]byte, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string][]byte{}, nil
		}
		return nil, err
	}

	var exp *regexp.Regexp
	if pattern != "" {
		exp, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
	}

	result := map[string][]byte{}
	for _, f := range files {
		if f.IsDir() || (exp != nil && !exp.MatchString(f.Name())) {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, err
		}

		result[f.Name()] = raw
	}

	return result, nil
}

// refreshTypesFile 刷新 TypeScript 类型文件（简化版）
func (p *plugin) refreshTypesFile() error {
	typesFile := filepath.Join(p.config.TypesDir, "types.d.ts")

	// 基础类型定义
	typesContent := `// URLDB Plugin System TypeScript Definitions

declare global {
  // 应用接口
  interface App {
  }

  // URL 模型
  interface URL {
    id: string;
    url: string;
    title: string;
    category: string;
    tags: string[];
    createdAt: Date;
    updatedAt: Date;
  }

  // 用户模型
  interface User {
    id: string;
    username: string;
    email: string;
    createdAt: Date;
  }

  // 钩子事件
  interface URLEvent {
    app: App;
    url: URL;
    data: Record<string, any>;
    next(): void;
  }

  interface UserEvent {
    app: App;
    user: User;
    data: Record<string, any>;
    next(): void;
  }

  interface APIEvent {
    app: App;
    request: any;
    path: string;
    method: string;
    headers: Record<string, string>;
    body: any;
    next(): void;
  }
}

// 钩子函数声明
declare function onURLAdd(handler: (e: URLEvent) => void): void;
declare function onURLAccess(handler: (e: URLEvent) => void): void;
declare function onUserLogin(handler: (e: UserEvent) => void): void;
declare function onAPIRequest(handler: (e: APIEvent) => void): void;

// 路由函数声明
declare function routerAdd(method: string, path: string, handler: (ctx: any) => void): void;

// 定时任务函数声明
declare function cronAdd(name: string, schedule: string, handler: () => void): void;

export {};
`

	// 确保目录存在
	dir := filepath.Dir(typesFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(typesFile, []byte(typesContent), 0644)
}