package jsvm

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dop251/goja"
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/utils"
)

// baseBinds 基础API绑定
func baseBinds(vm *goja.Runtime) {
	// 日志函数 - 记录到插件日志表和系统日志
	vm.Set("log", func(level, message string) {
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

		// TODO: 记录到插件日志表，需要获取当前插件名称
		// 这里需要传递插件名称和钩子名称信息
	})

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

				// 调用JavaScript处理器
				fn, _ := goja.AssertFunction(handler)
				_, err := fn(goja.Undefined(), vm.ToValue(e.URL), vm.ToValue(e.Data))
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

				fn, _ := goja.AssertFunction(handler)
				_, err := fn(goja.Undefined(), vm.ToValue(e.User), vm.ToValue(e.Data))
				if err != nil {
					utils.Error("JavaScript hook error: %v", err)
				}

				return e.Next()
			})
		}
	})

	vm.Set("onAPIRequest", func(handler goja.Value) {
		if _, ok := goja.AssertFunction(handler); ok {
			// TODO: 实现API请求钩子
			utils.Info("API request handler registered")
		}
	})

	vm.Set("onURLAccess", func(handler goja.Value) {
		if _, ok := goja.AssertFunction(handler); ok {
			// TODO: 实现URL访问钩子
			utils.Info("URL access handler registered")
		}
	})
}

// cronBinds 定时任务绑定
func cronBinds(app core.App, vm *goja.Runtime, executors *vmsPool) {
	vm.Set("cron", map[string]interface{}{
		"add": func(name, schedule string, handler goja.Value) error {
			if _, ok := goja.AssertFunction(handler); ok {
				// TODO: 实现定时任务注册
				utils.Info("Cron job registered: %s (%s)", name, schedule)
			}
			return nil
		},
	})

	// 为了兼容性，直接注册 cronAdd 函数
	vm.Set("cronAdd", func(name, schedule string, handler goja.Value) error {
		if _, ok := goja.AssertFunction(handler); ok {
			// TODO: 实现定时任务注册
			utils.Info("Cron job registered: %s (%s)", name, schedule)
		}
		return nil
	})
}

// routerBinds 路由绑定
func routerBinds(app core.App, vm *goja.Runtime, executors *vmsPool, routeRegister func(method, path string, handler func() (interface{}, error)) error) {
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
	})
}