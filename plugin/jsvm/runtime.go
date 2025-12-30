package jsvm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/buffer"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

const typesFileName = "types.d.ts"

var defaultScriptPath = "plugin.js"

func init() {
	// For backward compatibility and consistency with the Go exposed
	// methods that operate with relative paths (e.g. `$os.writeFile`),
	// we define the "current JS module" as if it is a file in the current working directory
	// (the filename itself doesn't really matter and in our case the hook handlers are executed as separate "programs").
	//
	// This is necessary for `require(module)` to properly traverse parents node_modules (goja_nodejs#95).
	cwd, err := os.Getwd()
	if err != nil {
		// truly rare case, log just for debug purposes
		color.Yellow("Failed to retrieve the current working directory: %v", err)
	} else {
		defaultScriptPath = filepath.Join(cwd, defaultScriptPath)
	}
}

// Config defines the config options of the jsvm plugin.
type Config struct {
	// OnInit is an optional function that will be called
	// after a JS runtime is initialized, allowing you to
	// attach custom Go variables and functions.
	OnInit func(vm *goja.Runtime)

	// RouteRegister is a function to register custom routes from plugins
	RouteRegister func(method, path string, handler func() (interface{}, error)) error

	// HooksWatch enables auto app restarts when a JS app hook file changes.
	//
	// Note that currently the application cannot be automatically restarted on Windows
	// because the restart process relies on execve.
	HooksWatch bool

	// HooksDir specifies the JS app hooks directory.
	//
	// If not set it fallbacks to a relative "./hooks" directory.
	HooksDir string

	// HooksFilesPattern specifies a regular expression pattern that
	// identify which file to load by the hook vm(s).
	//
	// If not set it fallbacks to `^.*(\.plugin\.js|\.plugin\.ts)$`, aka. any
	// HooksDir file ending in ".plugin.js" or ".plugin.ts" (the last one is to enforce IDE linters).
	HooksFilesPattern string

	// HooksPoolSize specifies how many goja.Runtime instances to prewarm
	// and keep for the JS app hooks gorotines execution.
	//
	// Zero or negative value means that it will create a new goja.Runtime
	// on every fired goroutine.
	HooksPoolSize int

	// MigrationsDir specifies the JS migrations directory.
	//
	// If not set it fallbacks to a relative "./migrations" directory.
	MigrationsDir string

	// If not set it fallbacks to `^.*(\.js|\.ts)$`, aka. any MigrationDir file
	// ending in ".js" or ".ts" (the last one is to enforce IDE linters).
	MigrationsFilesPattern string

	// TypesDir specifies the directory where to store the embedded
	// TypeScript declarations file.
	//
	// If not set it fallbacks to ".".
	//
	// Note: Avoid using the same directory as the HooksDir when HooksWatch is enabled
	// to prevent unnecessary app restarts when the types file is initially created.
	TypesDir string
}

// MustRegister registers the jsvm plugin in the provided app instance
// and panics if it fails.
//
// Example usage:
//
//	jsvm.MustRegister(app, jsvm.Config{
//		OnInit: func(vm *goja.Runtime) {
//			// register custom bindings
//			vm.Set("myCustomVar", 123)
//		},
//	})
func MustRegister(app core.App, config Config) {
	if err := Register(app, config); err != nil {
		panic(err)
	}
}

// Register registers the jsvm plugin in the provided app instance.
func Register(app core.App, config Config) error {
	return RegisterWithRepo(app, config, nil)
}

// RegisterWithRepo registers the jsvm plugin in the provided app instance with repository manager.
func RegisterWithRepo(app core.App, config Config, repoManager *repo.RepositoryManager) error {
	// 检查 app 是否为 nil
	if app == nil {
		return fmt.Errorf("app instance is nil")
	}

	p := &plugin{app: app, config: config, repoManager: repoManager}

	if p.config.HooksDir == "" {
		p.config.HooksDir = filepath.Join(".", "hooks")
	}

	if p.config.MigrationsDir == "" {
		p.config.MigrationsDir = filepath.Join(".", "migrations")
	}

	if p.config.HooksFilesPattern == "" {
		p.config.HooksFilesPattern = `^.*(\.plugin\.js|\.plugin\.ts)$`
	}

	if p.config.MigrationsFilesPattern == "" {
		p.config.MigrationsFilesPattern = `^.*(\.js|\.ts)$`
	}

	if p.config.TypesDir == "" {
		p.config.TypesDir = "."
	}

	p.app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		err := e.Next()
		if err != nil {
			return err
		}

		// ensure that the user has the latest types declaration
		err = p.refreshTypesFile()
		if err != nil {
			color.Yellow("Unable to refresh app types file: %v", err)
		}

		return nil
	})

	if err := p.registerMigrations(); err != nil {
		return fmt.Errorf("registerMigrations: %w", err)
	}

	if err := p.registerHooks(); err != nil {
		return fmt.Errorf("registerHooks: %w", err)
	}

	return nil
}

// MigrationEntry 存储迁移信息
type MigrationEntry struct {
	Name string
	Up   func(core.App) error
	Down func(core.App) error
}

type plugin struct {
	app         core.App
	config      Config
	repoManager *repo.RepositoryManager // RepositoryManager
	migrations  map[string]*MigrationEntry // 注册的迁移
}

// registerMigrations registers the JS migrations loader.
func (p *plugin) registerMigrations() error {
	// 初始化迁移存储
	p.migrations = make(map[string]*MigrationEntry)

	// fetch all js migrations sorted by their filename
	files, err := filesContent(p.config.MigrationsDir, p.config.MigrationsFilesPattern)
	if err != nil {
		return err
	}

	absHooksDir, err := filepath.Abs(p.config.HooksDir)
	if err != nil {
		return err
	}

	registry := new(require.Registry) // this can be shared by multiple runtimes

	for file, content := range files {
		vm := goja.New()

		registry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)
		buffer.Enable(vm)

		// 设置当前插件名称
		pluginName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		vm.Set("_currentPluginName", pluginName)
		vm.Set("_repoManager", p.repoManager)

		baseBinds(vm)
		dbxBinds(vm)
		securityBinds(vm)
		osBinds(vm)
		filepathBinds(vm)
		httpClientBinds(vm)
		filesystemBinds(vm)
		formsBinds(vm)
		mailsBinds(vm)

		vm.Set("__hooks", absHooksDir)

		vm.Set("migrate", func(up, down func(txApp core.App) error) {
			// 实现简化的迁移注册
			migrationName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

			// 注册迁移到内存中，实际执行可以通过 API 触发
			if !p.migrationRegistered(migrationName) {
				p.registerMigration(migrationName, up, down)
				utils.Info("Migration registered: %s (%s)", migrationName, file)
			} else {
				utils.Warn("Migration already registered: %s", migrationName)
			}
		})

		if p.config.OnInit != nil {
			p.config.OnInit(vm)
		}

		_, err := vm.RunScript(defaultScriptPath, string(content))
		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	return nil
}

// registerHooks registers the JS app hooks loader.
func (p *plugin) registerHooks() error {
	// fetch all js hooks sorted by their filename
	files, err := filesContent(p.config.HooksDir, p.config.HooksFilesPattern)
	if err != nil {
		return err
	}

	// prepend the types reference directive
	//
	// note: it is loaded during startup to handle conveniently also
	// the case when the HooksWatch option is enabled and the application
	// restart on newly created file
	for name, content := range files {
		if len(content) != 0 {
			// skip non-empty files for now to prevent accidental overwrite
			continue
		}
		path := filepath.Join(p.config.HooksDir, name)
		directive := `/// <reference path="` + p.relativeTypesPath(p.config.HooksDir) + `" />`
		if err := prependToEmptyFile(path, directive+"\n\n"); err != nil {
			color.Yellow("Unable to prepend the types reference: %v", err)
		}
	}

	// initialize the hooks dir watcher
	if p.config.HooksWatch {
		if err := p.watchHooks(); err != nil {
			color.Yellow("Unable to init hooks watcher: %v", err)
		}
	}

	if len(files) == 0 {
		// no need to register the vms since there are no entrypoint files anyway
		return nil
	}

	absHooksDir, err := filepath.Abs(p.config.HooksDir)
	if err != nil {
		return err
	}

	p.app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// 简化实现，暂时不绑定异常处理器
		return e.Next()
	})

	// safe to be shared across multiple vms
	requireRegistry := new(require.Registry)

	sharedBinds := func(vm *goja.Runtime) {
		requireRegistry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)
		buffer.Enable(vm)

		baseBinds(vm)
		dbxBinds(vm)
		filesystemBinds(vm)
		securityBinds(vm)
		osBinds(vm)
		filepathBinds(vm)
		httpClientBinds(vm)
		formsBinds(vm)
		apisBinds(vm)
		mailsBinds(vm)

		// 配置相关绑定（需要传递 repoManager）
		if p.repoManager != nil {
			configBinds(vm, p.repoManager)
		}

		vm.Set("$app", p.app)
		vm.Set("__hooks", absHooksDir)

		// 创建一个特殊的上下文获取函数，能够在执行时动态获取当前VM的上下文
		vm.Set("getCurrentPluginContext", func() map[string]interface{} {
			context := make(map[string]interface{})

			// 获取当前插件名称
			if cpn := vm.Get("_currentPluginName"); cpn != nil {
				if name, ok := cpn.Export().(string); ok && name != "" {
					context["pluginName"] = name
				} else {
					context["pluginName"] = "unknown"
				}
			} else {
				context["pluginName"] = "unknown"
			}

			// 获取 RepositoryManager
			if repoManager := vm.Get("_repoManager"); repoManager != nil {
				if rm, ok := repoManager.Export().(*repo.RepositoryManager); ok && rm != nil {
					context["repoManager"] = rm
				}
			}

			return context
		})

		// 创建一个使用全局 RepositoryManager 的 log 函数
		vm.Set("log", func(level, message, pluginName string) {
			// 输出到系统日志
			switch level {
			case "debug":
				utils.Debug(message)
			case "info":
				utils.Info(message)
			case "warn":
				utils.Warn(message)
			case "error":
				utils.Error(message)
			default:
				utils.Info(message)
			}

			// 如果没有提供插件名称，使用默认值
			if pluginName == "" {
				pluginName = "unknown"
			}

			// 使用全局的 RepositoryManager（从 p.repoManager 获取）
			if p.repoManager != nil {
				if pluginLogRepo := p.repoManager.PluginLogRepository; pluginLogRepo != nil {
					log := &entity.PluginLog{
						PluginName: pluginName,
						HookName:   "custom_log",
						Success:    level != "error", // error 级别的日志标记为不成功
						Message:    &message,        // 保存所有级别的日志消息
					}

					if level == "error" {
						log.ErrorMessage = &message // error 级别同时保存到错误消息字段
					}

					if err := pluginLogRepo.CreateLog(log); err != nil {
						utils.Error("FINAL-LOG: Failed to save plugin log to database: %v", err)
					} else {
						utils.Info("FINAL-LOG: Plugin log saved successfully to database: %s - %s", pluginName, message)
					}
				} else {
					utils.Error("FINAL-LOG: PLUGIN LOG REPO IS NIL")
				}
			} else {
				utils.Error("FINAL-LOG: GLOBAL REPO MANAGER NOT FOUND")
			}
		})

		if p.config.OnInit != nil {
			p.config.OnInit(vm)
		}
	}

	// initiliaze the executor vms
	executors := newPool(p.config.HooksPoolSize, func() *goja.Runtime {
		executor := goja.New()
		sharedBinds(executor)
		return executor
	})

	// initialize the loader vm
	loader := goja.New()
	sharedBinds(loader)
	hooksBinds(p.app, loader, executors)
	cronBinds(p.app, loader, executors, p.repoManager)
	routerBinds(p.app, loader, executors, p.repoManager, p.config.RouteRegister)

	for file, content := range files {
		func() {
			startTime := time.Now()
			var execErr error

			defer func() {
				if err := recover(); err != nil {
					execErr = fmt.Errorf("failed to execute %s:\n - %v", file, err)

					// 记录插件执行日志到数据库
					if logErr := p.recordPluginLog(file, "load", startTime, false, execErr.Error()); logErr != nil {
						utils.Error("Failed to record plugin execution log: %v", logErr)
					}

					if p.config.HooksWatch {
						color.Red("%v", execErr)
					} else {
						panic(execErr)
					}
				}
			}()

			_, err := loader.RunScript(defaultScriptPath, string(content))
			if err != nil {
				execErr = err
				panic(err)
			}

			// 记录插件成功执行日志到数据库
			if logErr := p.recordPluginLog(file, "load", startTime, true, ""); logErr != nil {
				utils.Error("Failed to record plugin execution log: %v", logErr)
			}
		}()
	}

	return nil
}

// normalizeExceptions registers a global error handler that
// wraps the extracted goja exception error value for consistency
// when throwing or returning errors.
func (p *plugin) normalizeServeExceptions(e interface{}) error {
	// 简化实现，直接返回
	return nil
}

// watchHooks initializes a hooks file watcher that will restart the
// application (*if possible) in case of a change in the hooks directory.
//
// This method does nothing if the hooks directory is missing.
func (p *plugin) watchHooks() error {
	watchDir := p.config.HooksDir

	hooksDirInfo, err := os.Lstat(p.config.HooksDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // no hooks dir to watch
		}
		return err
	}

	if hooksDirInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		watchDir, err = filepath.EvalSymlinks(p.config.HooksDir)
		if err != nil {
			return fmt.Errorf("failed to resolve hooksDir symlink: %w", err)
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	var debounceTimer *time.Timer

	stopDebounceTimer := func() {
		if debounceTimer != nil {
			debounceTimer.Stop()
			debounceTimer = nil
		}
	}

	p.app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		watcher.Close()

		stopDebounceTimer()

		return e.Next()
	})

	// start listening for events.
	go func() {
		defer stopDebounceTimer()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				stopDebounceTimer()

				debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
					// app restart is currently not supported on Windows
					if runtime.GOOS == "windows" {
						color.Yellow("File %s changed, please restart the app manually", event.Name)
					} else {
						color.Yellow("File %s changed, restarting...", event.Name)
						if err := p.app.Restart(); err != nil {
							color.Red("Failed to restart the app:", err)
						}
					}
				})
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				color.Red("Watch error:", err)
			}
		}
	}()

	// add directories to watch
	//
	// @todo replace once recursive watcher is added (https://github.com/fsnotify/fsnotify/issues/18)
	dirsErr := filepath.WalkDir(watchDir, func(path string, entry os.DirEntry, err error) error {
		// ignore hidden directories, node_modules, symlinks, sockets, etc.
		if !entry.IsDir() || entry.Name() == "node_modules" || strings.HasPrefix(entry.Name(), ".") {
			return nil
		}

		return watcher.Add(path)
	})
	if dirsErr != nil {
		watcher.Close()
	}

	return dirsErr
}

// fullTypesPathReturns returns the full path to the generated TS file.
func (p *plugin) fullTypesPath() string {
	return filepath.Join(p.config.TypesDir, typesFileName)
}

// relativeTypesPath returns a path to the generated TS file relative
// to the specified basepath.
//
// It fallbacks to the full path if generating the relative path fails.
func (p *plugin) relativeTypesPath(basepath string) string {
	fullPath := p.fullTypesPath()

	rel, err := filepath.Rel(basepath, fullPath)
	if err != nil {
		// fallback to the full path
		rel = fullPath
	}

	return rel
}

// refreshTypesFile saves the embedded TS declarations as a file on the disk.
func (p *plugin) refreshTypesFile() error {
	fullPath := p.fullTypesPath()

	// ensure that the types directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Basic TypeScript definitions for urldb
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
declare function onURLAccess(handler: (e: URLAccessEvent) => void): void;
declare function onUserLogin(handler: (e: UserEvent) => void): void;

// 路由函数声明
declare function routerAdd(method: string, path: string, handler: (ctx: any) => void): void;

// 定时任务函数声明
declare function cronAdd(name: string, schedule: string, handler: () => void): void;

// 全局变量
declare const $app: App;
declare const __hooks: string;

export {};
`

	// read the first timestamp line of the old file (if exists) and compare it to the embedded one
	// (note: ignore errors to allow always overwriting the file if it is invalid)
	existingFile, err := os.Open(fullPath)
	if err == nil {
		// For simplicity, always overwrite the file
		existingFile.Close()
	}

	return os.WriteFile(fullPath, []byte(typesContent), 0644)
}

// prependToEmptyFile prepends the specified text to an empty file.
//
// If the file is not empty this method does nothing.
func prependToEmptyFile(path, text string) error {
	info, err := os.Stat(path)

	if err == nil && info.Size() == 0 {
		return os.WriteFile(path, []byte(text), 0644)
	}

	return err
}

// filesContent returns a map with all direct files within the specified dir and their content.
//
// If directory with dirPath is missing or no files matching the pattern were found,
// it returns an empty map and no error.
//
// If pattern is empty string it matches all root files.
func filesContent(dirPath string, pattern string) (map[string][]byte, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string][]byte{}, nil
		}
		return nil, err
	}

	var exp *regexp.Regexp
	if pattern != "" {
		var err error
		if exp, err = regexp.Compile(pattern); err != nil {
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

// recordPluginLog 记录插件执行日志到数据库
func (p *plugin) recordPluginLog(pluginName, hookName string, startTime time.Time, success bool, errorMessage string) error {
	// 计算执行时间
	executionTime := time.Since(startTime).Milliseconds()

	// 使用插件的真实名称（应该已经是解析后的名称）
	name := pluginName
	// 如果传入的还是文件名格式，做兼容处理
	if strings.HasSuffix(name, ".plugin.js") {
		name = strings.TrimSuffix(name, ".plugin.js")
	} else if strings.HasSuffix(name, ".plugin.ts") {
		name = strings.TrimSuffix(name, ".plugin.ts")
	}

	// 创建插件日志记录
	log := &entity.PluginLog{
		PluginName:    name,
		HookName:      hookName,
		ExecutionTime: int(executionTime),
		Success:       success,
	}

	if !success && errorMessage != "" {
		log.ErrorMessage = &errorMessage
	}

	// 尝试记录到数据库
	if p.repoManager != nil {
		// 直接使用 RepositoryManager
		if pluginLogRepo := p.repoManager.PluginLogRepository; pluginLogRepo != nil {
			if err := pluginLogRepo.CreateLog(log); err != nil {
				utils.Error("Failed to save plugin log to database: %v", err)
				// 如果数据库记录失败，仍然记录到系统日志
			} else {
				// 数据库记录成功，也记录到系统日志用于调试
				if success {
					utils.Info("Plugin '%s' executed successfully (%s) in %dms (logged to db)", name, hookName, executionTime)
				} else {
					utils.Error("Plugin '%s' execution failed (%s): %s (logged to db)", name, hookName, errorMessage)
				}
				return nil
			}
		}
	}

	// 如果无法访问数据库或记录失败，记录到系统日志
	if success {
		utils.Info("Plugin '%s' executed successfully (%s) in %dms", name, hookName, executionTime)
	} else {
		utils.Error("Plugin '%s' execution failed (%s): %s", name, hookName, errorMessage)
	}

	return nil
}

// normalizeException normalizes goja exceptions for consistent error handling.
func normalizeException(err error) error {
	// Simple implementation - in real PocketBase this would be more sophisticated
	return err
}

// migrationRegistered 检查迁移是否已注册
func (p *plugin) migrationRegistered(name string) bool {
	_, exists := p.migrations[name]
	return exists
}

// registerMigration 注册迁移到内存中
func (p *plugin) registerMigration(name string, up, down func(core.App) error) {
	p.migrations[name] = &MigrationEntry{
		Name: name,
		Up:   up,
		Down: down,
	}
}

// GetRegisteredMigrations 返回已注册的迁移列表（可用于API端点）
func (p *plugin) GetRegisteredMigrations() map[string]*MigrationEntry {
	if p.migrations == nil {
		return make(map[string]*MigrationEntry)
	}

	// 返回副本以避免外部修改
	result := make(map[string]*MigrationEntry)
	for k, v := range p.migrations {
		result[k] = v
	}
	return result
}

// ExecuteMigration 执行指定的迁移
func (p *plugin) ExecuteMigration(name string) error {
	if p.migrations == nil {
		return fmt.Errorf("no migrations registered")
	}

	migration, exists := p.migrations[name]
	if !exists {
		return fmt.Errorf("migration '%s' not found", name)
	}

	utils.Info("Executing migration: %s", name)
	err := migration.Up(p.app)
	if err != nil {
		utils.Error("Migration '%s' failed: %v", name, err)
		return fmt.Errorf("migration '%s' execution failed: %w", name, err)
	}

	utils.Info("Migration '%s' executed successfully", name)
	return nil
}