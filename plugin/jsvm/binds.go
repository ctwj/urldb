package jsvm

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/robfig/cron/v3"
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// 全局cron调度器管理
type cronManager struct {
	scheduler *cron.Cron
	jobs      map[string]cron.EntryID
	jobsMux   sync.RWMutex
}

var globalCronManager = &cronManager{
	scheduler: cron.New(),
	jobs:      make(map[string]cron.EntryID),
}

// 初始化cron调度器
func init() {
	globalCronManager.scheduler.Start()
}

// baseBinds 基础API绑定
func baseBinds(vm *goja.Runtime) {

	// 工具函数
	vm.Set("jsonParse", func(str string) goja.Value {
		var result interface{}
		if err := json.Unmarshal([]byte(str), &result); err != nil {
			return vm.ToValue(nil)
		}
		return vm.ToValue(result)
	})

	vm.Set("jsonStringify", func(data goja.Value) string {
		jsonData, err := json.Marshal(data.Export())
		if err != nil {
			return ""
		}
		return string(jsonData)
	})

	vm.Set("sleep", func(ms int64) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	})

	vm.Set("timestamp", func() int64 {
		return time.Now().Unix()
	})
}

// dbxBinds 数据库相关绑定（简化版）
func dbxBinds(vm *goja.Runtime) {
	// 简化的数据库操作，实际需要适配到urldb的GORM
	vm.Set("db", map[string]interface{}{
		"find": func(table string, query interface{}) goja.Value {
			// TODO: 实现GORM查询
			utils.Info("DB find called for table: %s", table)
			return vm.ToValue([]interface{}{})
		},
		"save": func(table string, data interface{}) error {
			// TODO: 实现GORM保存
			utils.Info("DB save called for table: %s", table)
			return nil
		},
		"update": func(table string, id interface{}, data interface{}) error {
			// TODO: 实现GORM更新
			utils.Info("DB update called for table: %s, id: %v", table, id)
			return nil
		},
		"delete": func(table string, id interface{}) error {
			// TODO: 实现GORM删除
			utils.Info("DB delete called for table: %s, id: %v", table, id)
			return nil
		},
	})
}

// securityBinds 安全相关绑定
func securityBinds(vm *goja.Runtime) {
	vm.Set("security", map[string]interface{}{
		"hash": func(password string) string {
			// TODO: 实现密码哈希
			return "hashed_" + password
		},
		"verify": func(password, hash string) bool {
			// TODO: 实现密码验证
			return hash == "hashed_"+password
		},
	})
}

// osBinds 操作系统相关绑定
func osBinds(vm *goja.Runtime) {
	vm.Set("os", map[string]interface{}{
		"env": func(key string) string {
			// TODO: 实现环境变量获取
			return ""
		},
		"platform": func() string {
			return "linux"
		},
	})
}

// filepathBinds 文件路径相关绑定
func filepathBinds(vm *goja.Runtime) {
	vm.Set("filepath", map[string]interface{}{
		"join": func(parts ...string) string {
			result := ""
			for i, part := range parts {
				if i > 0 {
					result += "/"
				}
				result += part
			}
			return result
		},
		"base": func(path string) string {
			// 简化实现
			if idx := len(path) - 1; idx >= 0 && path[idx] == '/' {
				path = path[:idx]
			}
			for i := len(path) - 1; i >= 0; i-- {
				if path[i] == '/' {
					return path[i+1:]
				}
			}
			return path
		},
	})
}

// httpClientBinds HTTP客户端绑定
func httpClientBinds(vm *goja.Runtime) {
	vm.Set("http", map[string]interface{}{
		"get": func(url string, headers map[string]string) map[string]interface{} {
			// 简化的HTTP GET实现
			resp, err := http.Get(url)
			if err != nil {
				return map[string]interface{}{
					"status": 500,
					"error": err.Error(),
				}
			}
			defer resp.Body.Close()

			return map[string]interface{}{
				"status": resp.StatusCode,
				"headers": resp.Header,
				"body":   "Response body",
			}
		},
		"post": func(url string, data interface{}, headers map[string]string) map[string]interface{} {
			// 简化的HTTP POST实现
			return map[string]interface{}{
				"status": 200,
				"body":   "Mock POST response",
			}
		},
	})
}

// filesystemBinds 文件系统绑定
func filesystemBinds(vm *goja.Runtime) {
	vm.Set("fs", map[string]interface{}{
		"readFile": func(path string) string {
			// TODO: 实现文件读取
			return "File content from " + path
		},
		"writeFile": func(path string, content string) error {
			// TODO: 实现文件写入
			utils.Info("Write file called: %s", path)
			return nil
		},
	})
}

// formsBinds 表单绑定（简化版）
func formsBinds(vm *goja.Runtime) {
	vm.Set("forms", map[string]interface{}{
		"validate": func(data interface{}, rules interface{}) bool {
			// TODO: 实现表单验证
			return true
		},
	})
}

// mailsBinds 邮件绑定（简化版）
func mailsBinds(vm *goja.Runtime) {
	vm.Set("mails", map[string]interface{}{
		"send": func(to, subject, body string) error {
			// TODO: 实现邮件发送
			utils.Info("Mail sent to %s: %s", to, subject)
			return nil
		},
	})
}

// apisBinds API绑定（urldb特定）
func apisBinds(vm *goja.Runtime) {
	vm.Set("apis", map[string]interface{}{
		"request": func(method, path string, data interface{}) map[string]interface{} {
			// TODO: 实现内部API调用
			return map[string]interface{}{
				"status": 200,
				"data":   "API response",
			}
		},
	})
}

// hooksBinds 钩子绑定
func hooksBinds(app core.App, vm *goja.Runtime, executors *vmsPool) {
	vm.Set("onURLAdd", func(handler goja.Value) {
		if _, ok := goja.AssertFunction(handler); ok {
			// 注册URL添加钩子
			app.OnURLAdd().BindFunc(func(e *core.URLEvent) error {
				// 从池中获取VM实例
				vm := executors.Get()
				defer executors.Put(vm)

				// 创建事件对象，包含 url 和 data 属性
				eventObj := vm.NewObject()
				if e.URL != nil {
					urlObj := vm.NewObject()
					urlObj.Set("id", e.URL.ID)
					urlObj.Set("key", e.URL.Key)
					urlObj.Set("title", e.URL.Title)
					urlObj.Set("url", e.URL.URL)
					urlObj.Set("description", e.URL.Description)
					urlObj.Set("category_id", e.URL.CategoryID)
					urlObj.Set("tags", e.URL.Tags)
					urlObj.Set("is_valid", e.URL.IsValid)
					urlObj.Set("is_public", e.URL.IsPublic)
					urlObj.Set("view_count", e.URL.ViewCount)
					urlObj.Set("created_at", e.URL.CreatedAt)
					urlObj.Set("updated_at", e.URL.UpdatedAt)
					eventObj.Set("url", urlObj)
				}

				if e.Data != nil {
					eventObj.Set("data", vm.ToValue(e.Data))
				}

				// 添加应用信息
				if e.App != nil {
					appObj := vm.NewObject()
					appObj.Set("name", "URLDB")
					appObj.Set("version", "1.0.0")
					eventObj.Set("app", appObj)
				}

				// 调用JavaScript处理器
				fn, _ := goja.AssertFunction(handler)
				_, err := fn(goja.Undefined(), eventObj)
				if err != nil {
					utils.Error("JavaScript hook error: %v", err)
				}

				return e.Next()
			})
		}
	})

	vm.Set("onUserLogin", func(handler goja.Value) {
		if _, ok := goja.AssertFunction(handler); ok {
			app.OnUserLogin().BindFunc(func(e *core.UserEvent) error {
				vm := executors.Get()
				defer executors.Put(vm)

				// 创建事件对象，包含 user 和 data 属性
				eventObj := vm.NewObject()
				if e.User != nil {
					userObj := vm.NewObject()
					userObj.Set("id", e.User.ID)
					userObj.Set("username", e.User.Username)
					userObj.Set("email", e.User.Email)
					userObj.Set("role", e.User.Role)
					userObj.Set("is_active", e.User.IsActive)
					userObj.Set("last_login", e.User.LastLogin)
					userObj.Set("created_at", e.User.CreatedAt)
					userObj.Set("updated_at", e.User.UpdatedAt)
					eventObj.Set("user", userObj)
				}

				if e.Data != nil {
					eventObj.Set("data", vm.ToValue(e.Data))
				}

				// 添加应用信息
				if e.App != nil {
					appObj := vm.NewObject()
					appObj.Set("name", "URLDB")
					appObj.Set("version", "1.0.0")
					eventObj.Set("app", appObj)
				}

				fn, _ := goja.AssertFunction(handler)
				_, err := fn(goja.Undefined(), eventObj)
				if err != nil {
					utils.Error("JavaScript hook error: %v", err)
				}

				return e.Next()
			})
		}
	})

	vm.Set("onURLAccess", func(handler goja.Value) {
		if _, ok := goja.AssertFunction(handler); ok {
			app.OnURLAccess().BindFunc(func(e *core.URLAccessEvent) error {
				vm := executors.Get()
				defer executors.Put(vm)

				// 创建事件对象，包含 url、access_log、request、response 属性
				eventObj := vm.NewObject()
				if e.URL != nil {
					urlObj := vm.NewObject()
					urlObj.Set("id", e.URL.ID)
					urlObj.Set("key", e.URL.Key)
					urlObj.Set("title", e.URL.Title)
					urlObj.Set("url", e.URL.URL)
					urlObj.Set("description", e.URL.Description)
					urlObj.Set("category_id", e.URL.CategoryID)
					urlObj.Set("tags", e.URL.Tags)
					urlObj.Set("is_valid", e.URL.IsValid)
					urlObj.Set("is_public", e.URL.IsPublic)
					urlObj.Set("view_count", e.URL.ViewCount)
					urlObj.Set("created_at", e.URL.CreatedAt)
					urlObj.Set("updated_at", e.URL.UpdatedAt)
					eventObj.Set("url", urlObj)
				}

				if e.AccessLog != nil {
					eventObj.Set("access_log", vm.ToValue(e.AccessLog))
				}

				if e.Request != nil {
					eventObj.Set("request", vm.ToValue(e.Request))
				}

				if e.Response != nil {
					eventObj.Set("response", vm.ToValue(e.Response))
				}

				// 添加应用信息
				if e.App != nil {
					appObj := vm.NewObject()
					appObj.Set("name", "URLDB")
					appObj.Set("version", "1.0.0")
					eventObj.Set("app", appObj)
				}

				fn, _ := goja.AssertFunction(handler)
				_, err := fn(goja.Undefined(), eventObj)
				if err != nil {
					utils.Error("JavaScript hook error: %v", err)
				}

				return e.Next()
			})
		}
	})
}

// cronBinds 定时任务绑定
func cronBinds(app core.App, vm *goja.Runtime, executors *vmsPool, repoManager *repo.RepositoryManager) {
	vm.Set("cron", map[string]interface{}{
		"add": func(name, schedule string, handler goja.Value) error {
			if fn, ok := goja.AssertFunction(handler); ok {
				// 实际添加到cron调度器
				globalCronManager.jobsMux.Lock()
				defer globalCronManager.jobsMux.Unlock()

				// 如果同名任务已存在，先移除
				if entryID, exists := globalCronManager.jobs[name]; exists {
					globalCronManager.scheduler.Remove(entryID)
					delete(globalCronManager.jobs, name)
					utils.Info("Removed existing cron job: %s", name)
				}

				// 创建包装函数，从池中获取VM实例
				wrappedFunc := func() {
					executor := executors.Get()
					defer executors.Put(executor)

					// 设置当前插件上下文
					pluginName := extractPluginNameFromCronJob(name)
					if pluginName != "" {
						executor.Set("_currentPluginName", pluginName)
					} else {
						executor.Set("_currentPluginName", "cron_job")
					}
					executor.Set("_repoManager", repoManager)

					_, err := fn(goja.Undefined())
					if err != nil {
						utils.Error("Cron job '%s' execution error: %v", name, err)
					} else {
						utils.Info("Cron job '%s' executed successfully", name)
					}
				}

				// 添加到调度器
				entryID, err := globalCronManager.scheduler.AddFunc(schedule, wrappedFunc)
				if err != nil {
					utils.Error("Failed to add cron job '%s': %v", name, err)
					return err
				}

				// 保存任务ID
				globalCronManager.jobs[name] = entryID
				utils.Info("Cron job registered and started: %s (%s)", name, schedule)
			}
			return nil
		},
	})

	// 为了兼容性，直接注册 cronAdd 函数
	vm.Set("cronAdd", func(name, schedule string, handler goja.Value) error {
		if fn, ok := goja.AssertFunction(handler); ok {
			// 实际添加到cron调度器
			globalCronManager.jobsMux.Lock()
			defer globalCronManager.jobsMux.Unlock()

			// 如果同名任务已存在，先移除
			if entryID, exists := globalCronManager.jobs[name]; exists {
				globalCronManager.scheduler.Remove(entryID)
				delete(globalCronManager.jobs, name)
				utils.Info("Removed existing cron job: %s", name)
			}

			// 创建包装函数，从池中获取VM实例
			wrappedFunc := func() {
				// 添加panic恢复机制，防止整个程序崩溃
				defer func() {
					if r := recover(); r != nil {
						utils.Error("Cron job '%s' panicked: %v", name, r)
						utils.Error("Stack trace: %v", r)
						// 不要重新panic，只是记录错误
					}
				}()

				// 检查插件是否启用
				pluginName := extractPluginNameFromCronJob(name)
				if repoManager != nil && pluginName != "" {
					if config, err := repoManager.PluginConfigRepository.GetConfig(pluginName); err == nil && config != nil && !config.Enabled {
						utils.Debug("Cron job '%s' skipped: plugin '%s' is disabled", name, pluginName)
						return
					}
				}

				executor := executors.Get()
				defer executors.Put(executor)

				// 设置当前插件上下文
				if pluginName != "" {
					executor.Set("_currentPluginName", pluginName)
					utils.Info("CRON: Set _currentPluginName to %s for VM %p", pluginName, executor)
				} else {
					executor.Set("_currentPluginName", "cron_job")
					utils.Info("CRON: Set _currentPluginName to cron_job for VM %p", executor)
				}
				executor.Set("_repoManager", repoManager)
				utils.Info("CRON: Set _repoManager for VM %p", executor)

				// 再次保护，防止VM调用时出错
				func() {
					defer func() {
						if r := recover(); r != nil {
							utils.Error("Cron job '%s' VM execution panicked: %v", name, r)
						}
					}()

					_, err := fn(goja.Undefined())
					if err != nil {
						utils.Error("Cron job '%s' execution error: %v", name, err)
					} else {
						utils.Info("Cron job '%s' executed successfully", name)
					}
				}()
			}

			// 添加到调度器
			entryID, err := globalCronManager.scheduler.AddFunc(schedule, wrappedFunc)
			if err != nil {
				utils.Error("Failed to add cron job '%s': %v", name, err)
				return err
			}

			// 保存任务ID
			globalCronManager.jobs[name] = entryID
			utils.Info("Cron job registered and started: %s (%s)", name, schedule)
		}
		return nil
	})
}

// configBinds 配置相关绑定
func configBinds(vm *goja.Runtime, repoManager *repo.RepositoryManager) {
	// 获取插件配置函数
	vm.Set("getPluginConfig", func(pluginName string) goja.Value {
		// 从数据库查询插件配置
		config, err := repoManager.PluginConfigRepository.GetConfig(pluginName)
		if err != nil {
			utils.Error("Failed to get plugin config for %s: %v", pluginName, err)
			return vm.ToValue(nil)
		}

		// 解析配置 JSON
		var configData interface{}
		if err := json.Unmarshal([]byte(config.ConfigJSON), &configData); err != nil {
			utils.Error("Failed to parse config JSON for %s: %v", pluginName, err)
			return vm.ToValue(nil)
		}

		utils.Info("Plugin config loaded for %s: %v", pluginName, configData)
		return vm.ToValue(configData)
	})

	// 设置插件配置函数
	vm.Set("setPluginConfig", func(pluginName string, configData goja.Value) error {
		// 保存到数据库
		err := repoManager.PluginConfigRepository.SetConfig(pluginName, configData.Export().(map[string]interface{}))
		if err != nil {
			utils.Error("Failed to save config for %s: %v", pluginName, err)
			return err
		}

		utils.Info("Plugin config saved for %s", pluginName)
		return nil
	})

	// 获取插件启用状态
	vm.Set("isPluginEnabled", func(pluginName string) bool {
		config, err := repoManager.PluginConfigRepository.GetConfig(pluginName)
		if err != nil {
			utils.Error("Failed to get plugin status for %s: %v", pluginName, err)
			return false
		}
		return config.Enabled
	})

	// 设置插件启用状态
	vm.Set("setPluginEnabled", func(pluginName string, enabled bool) error {
		err := repoManager.PluginConfigRepository.SetEnabled(pluginName, enabled)
		if err != nil {
			utils.Error("Failed to set plugin status for %s: %v", pluginName, err)
			return err
		}

		utils.Info("Plugin %s enabled: %v", pluginName, enabled)
		return nil
	})
}

// routerBinds 路由绑定
func routerBinds(app core.App, vm *goja.Runtime, executors *vmsPool, repoManager *repo.RepositoryManager, routeRegister func(method, path string, handler func() (interface{}, error)) error) {
	vm.Set("router", map[string]interface{}{
		"add": func(method, path string, handler goja.Value) error {
			if _, ok := goja.AssertFunction(handler); ok {
				if routeRegister != nil {
					// 将 JavaScript handler 转换为 Go handler
					goHandler := func() (interface{}, error) {
						vm := executors.Get()
						defer executors.Put(vm)

						// 创建一个模拟的事件对象，提供 json 方法
						event := map[string]interface{}{
							"json": func(status int, data interface{}) interface{} {
								return map[string]interface{}{
									"status": status,
									"data":   data,
								}
							},
						}

						fn, _ := goja.AssertFunction(handler)
						result, err := fn(goja.Undefined(), vm.ToValue(event))
						if err != nil {
							return nil, err
						}

						// 导出结果
						exported := result.Export()
						if resultMap, ok := exported.(map[string]interface{}); ok {
							if _, hasStatus := resultMap["status"]; hasStatus {
								if data, hasData := resultMap["data"]; hasData {
									return data, nil
								}
							}
						}
						return exported, nil
					}
					return routeRegister(method, path, goHandler)
				}
				// 如果没有注册路由器，只记录日志
				utils.Info("Route registered (no router bind): %s %s", method, path)
			}
			return nil
		},
	})

	// 为了兼容性，直接注册 routerAdd 函数
	vm.Set("routerAdd", func(method, path string, handler goja.Value) error {
		if _, ok := goja.AssertFunction(handler); ok {
			if routeRegister != nil {
				// 将 JavaScript handler 转换为 Go handler
				goHandler := func() (interface{}, error) {
					// 添加panic恢复机制，防止整个程序崩溃
					defer func() {
						if r := recover(); r != nil {
							utils.Error("Route handler panicked: %v", r)
						}
					}()

					vm := executors.Get()
					defer executors.Put(vm)

					// 设置当前插件上下文（从路由注册时推断）
					vm.Set("_currentPluginName", "route_handler")
					vm.Set("_repoManager", repoManager)

					// 创建一个模拟的事件对象，提供 json 方法
					event := map[string]interface{}{
						"json": func(status int, data interface{}) interface{} {
							return map[string]interface{}{
								"status": status,
								"data":   data,
							}
						},
					}

					// 保护JavaScript执行
					var finalResult interface{}
					func() {
						defer func() {
							if r := recover(); r != nil {
								utils.Error("Route VM execution panicked: %v", r)
							}
						}()

						fn, _ := goja.AssertFunction(handler)
						result, err := fn(goja.Undefined(), vm.ToValue(event))
						if err != nil {
							utils.Error("Route execution error: %v", err)
							return
						}

						// 导出结果
						exported := result.Export()
						if resultMap, ok := exported.(map[string]interface{}); ok {
							if _, hasStatus := resultMap["status"]; hasStatus {
								if data, hasData := resultMap["data"]; hasData {
									utils.Info("Plugin route handler success: %v", data)
									finalResult = data
									return
								}
							}
						}
						utils.Info("Plugin route handler success: %v", exported)
						finalResult = exported
					}()

					// 返回处理结果
					if finalResult != nil {
						return finalResult, nil
					}

					// 返回默认响应，避免nil
					return map[string]interface{}{
						"message": "Plugin route executed",
						"success": true,
					}, nil
				}
				return routeRegister(method, path, goHandler)
			}
			// 如果没有注册路由器，只记录日志
			utils.Info("Route registered (no router bind): %s %s", method, path)
		}
		return nil
	})
}

// extractPluginNameFromCronJob 从cron任务名称中提取插件名称
func extractPluginNameFromCronJob(cronJobName string) string {
	// 常见的任务名称模式：
	// - config_demo_task -> config_demo
	// - analytics_engine_task -> analytics_engine
	// - test-job -> test
	// - daily_report -> daily_report (可能是独立的)

	// 如果以 _task 结尾，去掉后缀
	if strings.HasSuffix(cronJobName, "_task") {
		return strings.TrimSuffix(cronJobName, "_task")
	}

	// 如果包含连字符，取第一部分
	if strings.Contains(cronJobName, "-") {
		parts := strings.Split(cronJobName, "-")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// 如果包含下划线，尝试推断插件名
	if strings.Contains(cronJobName, "_") {
		// 对于像 daily_report 这样的名称，可能本身就是插件名
		// 但对于像 config_demo_task 这样的，我们已经在上面处理了
		parts := strings.Split(cronJobName, "_")
		if len(parts) >= 2 {
			// 检查是否是常见的任务后缀
			lastPart := parts[len(parts)-1]
			if lastPart == "task" || lastPart == "job" || lastPart == "report" {
				return strings.Join(parts[:len(parts)-1], "_")
			}
		}
	}

	// 默认返回原名称
	return cronJobName
}