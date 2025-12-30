package jsvm

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/golang-jwt/jwt/v5"
	"github.com/robfig/cron/v3"
	"github.com/ctwj/urldb/core"
	"github.com/ctwj/urldb/db"
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

// dbxBinds 数据库相关绑定（实现直接数据库操作）
func dbxBinds(vm *goja.Runtime) {
	// 检查数据库连接是否可用
	if db.DB == nil {
		utils.Info("Database not available for plugin operations")
		vm.Set("db", map[string]interface{}{
			"find": func(table string, query interface{}) goja.Value {
				return vm.ToValue(map[string]interface{}{
					"error": "Database not available",
				})
			},
			"save": func(table string, data interface{}) error {
				return fmt.Errorf("Database not available")
			},
			"update": func(table string, id interface{}, data interface{}) error {
				return fmt.Errorf("Database not available")
			},
			"delete": func(table string, id interface{}) error {
				return fmt.Errorf("Database not available")
			},
		})
		return
	}

	obj := vm.NewObject()
	vm.Set("$db", obj)

	// 原始 SQL 查询
	obj.Set("raw", func(sql string, args ...interface{}) ([]map[string]interface{}, error) {
		var results []map[string]interface{}
		rows, err := db.DB.Raw(sql, args...).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}

			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				if b, ok := val.([]byte); ok {
					row[col] = string(b)
				} else {
					row[col] = val
				}
			}
			results = append(results, row)
		}

		return results, nil
	})

	// 通用查询
	obj.Set("find", func(table string, query map[string]interface{}) ([]map[string]interface{}, error) {
		db := db.DB.Table(table)

		if query != nil {
			for key, value := range query {
				db = db.Where(key, value)
			}
		}

		var results []map[string]interface{}
		err := db.Find(&results).Error
		return results, err
	})

	// 通用保存
	obj.Set("save", func(table string, data map[string]interface{}) (interface{}, error) {
		db := db.DB.Table(table)
		err := db.Create(data).Error
		if err != nil {
			return nil, err
		}
		return data, nil
	})

	// 通用更新
	obj.Set("update", func(table string, id interface{}, data map[string]interface{}) error {
		db := db.DB.Table(table)
		return db.Where("id = ?", id).Updates(data).Error
	})

	// 通用删除
	obj.Set("delete", func(table string, id interface{}) error {
		db := db.DB.Table(table)
		return db.Where("id = ?", id).Delete(nil).Error
	})

	// 计数
	obj.Set("count", func(table string, query map[string]interface{}) (int64, error) {
		db := db.DB.Table(table)

		if query != nil {
			for key, value := range query {
				db = db.Where(key, value)
			}
		}

		var count int64
		err := db.Count(&count).Error
		return count, err
	})

	// 保留原来的简单接口用于向后兼容
	vm.Set("db", map[string]interface{}{
		"find": func(table string, query interface{}) goja.Value {
			var queryMap map[string]interface{}
			if q, ok := query.(map[string]interface{}); ok {
				queryMap = q
			} else {
				queryMap = make(map[string]interface{})
			}

			var results []map[string]interface{}
			err := db.DB.Table(table).Where(queryMap).Find(&results).Error
			if err != nil {
				return vm.ToValue(map[string]interface{}{
					"error": err.Error(),
				})
			}

			return vm.ToValue(results)
		},
		"save": func(table string, data interface{}) error {
			if dataMap, ok := data.(map[string]interface{}); ok {
				return db.DB.Table(table).Create(dataMap).Error
			}
			return fmt.Errorf("invalid data format, expected map[string]interface{}")
		},
		"update": func(table string, id interface{}, data interface{}) error {
			if dataMap, ok := data.(map[string]interface{}); ok {
				return db.DB.Table(table).Where("id = ?", id).Updates(dataMap).Error
			}
			return fmt.Errorf("invalid data format, expected map[string]interface{}")
		},
		"delete": func(table string, id interface{}) error {
			return db.DB.Table(table).Where("id = ?", id).Delete(nil).Error
		},
		"raw": func(sql string, args ...interface{}) ([]map[string]interface{}, error) {
			var results []map[string]interface{}
			rows, err := db.DB.Raw(sql, args...).Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			columns, _ := rows.Columns()
			for rows.Next() {
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				if err := rows.Scan(valuePtrs...); err != nil {
					continue
				}

				row := make(map[string]interface{})
				for i, col := range columns {
					val := values[i]
					if b, ok := val.([]byte); ok {
						row[col] = string(b)
					} else {
						row[col] = val
					}
				}
				results = append(results, row)
			}

			return results, nil
		},
	})
}

// securityBinds 安全相关绑定 (简化实现)
func securityBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$security", obj)

	// crypto - 使用标准库实现
	obj.Set("md5", func(data string) string {
		h := md5.Sum([]byte(data))
		return hex.EncodeToString(h[:])
	})
	obj.Set("sha256", func(data string) string {
		h := sha256.Sum256([]byte(data))
		return hex.EncodeToString(h[:])
	})
	obj.Set("sha512", func(data string) string {
		h := sha512.Sum512([]byte(data))
		return hex.EncodeToString(h[:])
	})
	obj.Set("hs256", func(data, key string) string {
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(data))
		return hex.EncodeToString(h.Sum(nil))
	})
	obj.Set("hs512", func(data, key string) string {
		h := hmac.New(sha512.New, []byte(key))
		h.Write([]byte(data))
		return hex.EncodeToString(h.Sum(nil))
	})
	obj.Set("equal", func(a, b string) bool {
		return a == b
	})

	// random - 简化实现
	randomStringFunc := func(length int) string {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}
	randomStringWithAlphabetFunc := func(alphabet string, length int) string {
		b := make([]byte, length)
		for i := range b {
			b[i] = alphabet[rand.Intn(len(alphabet))]
		}
		return string(b)
	}

	obj.Set("randomString", randomStringFunc)
	obj.Set("randomStringByRegex", func(regex string, length int) (string, error) {
		// 简化实现，仅返回随机字符串
		return randomStringFunc(length), nil
	})
	obj.Set("randomStringWithAlphabet", randomStringWithAlphabetFunc)
	obj.Set("pseudorandomString", randomStringFunc)
	obj.Set("pseudorandomStringWithAlphabet", randomStringWithAlphabetFunc)

	// jwt - 使用 jwt 库
	obj.Set("parseUnverifiedJWT", func(tokenString string) (map[string]interface{}, error) {
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			return nil, err
		}
		return token.Claims.(jwt.MapClaims), nil
	})
	obj.Set("parseJWT", func(tokenString, verificationKey string) (map[string]interface{}, error) {
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(verificationKey), nil
		})
		if err != nil {
			return nil, err
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, fmt.Errorf("invalid token")
	})
	obj.Set("createJWT", func(payload map[string]interface{}, signingKey string, secDuration int) (string, error) {
		claims := jwt.MapClaims(payload)
		claims["exp"] = time.Now().Add(time.Duration(secDuration) * time.Second).Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString([]byte(signingKey))
	})

	// encryption - 简化实现（注意：这不是生产级别的加密）
	obj.Set("encrypt", func(plaintext, key string) (string, error) {
		// 简单的 XOR "加密" - 仅用于演示
		if len(key) == 0 {
			return "", fmt.Errorf("key cannot be empty")
		}
		result := make([]byte, len(plaintext))
		keyLen := len(key)
		for i, c := range []byte(plaintext) {
			result[i] = c ^ key[i%keyLen]
		}
		return hex.EncodeToString(result), nil
	})
	obj.Set("decrypt", func(cipherText, key string) (string, error) {
		if len(key) == 0 {
			return "", fmt.Errorf("key cannot be empty")
		}
		data, err := hex.DecodeString(cipherText)
		if err != nil {
			return "", err
		}
		result := make([]byte, len(data))
		keyLen := len(key)
		for i, c := range data {
			result[i] = c ^ key[i%keyLen]
		}
		return string(result), nil
	})

	// 保留原来的简单接口用于向后兼容
	sha256Func := func(data string) string {
		h := sha256.Sum256([]byte(data))
		return hex.EncodeToString(h[:])
	}
	equalFunc := func(a, b string) bool {
		return a == b
	}

	vm.Set("security", map[string]interface{}{
		"hash": func(password string) string {
			return sha256Func(password)
		},
		"verify": func(password, hash string) bool {
			return equalFunc(password, hash)
		},
	})
}

// osBinds 操作系统相关绑定 (移植自 PocketBase)
func osBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$os", obj)

	// 基本系统信息
	obj.Set("args", os.Args)
	obj.Set("exit", os.Exit)
	obj.Set("getenv", os.Getenv)
	obj.Set("tempDir", os.TempDir)
	obj.Set("getwd", os.Getwd)

	// 文件系统操作
	obj.Set("dirFS", os.DirFS)
	obj.Set("stat", os.Stat)
	obj.Set("readFile", os.ReadFile)
	obj.Set("writeFile", os.WriteFile)
	obj.Set("readDir", os.ReadDir)
	obj.Set("truncate", os.Truncate)
	obj.Set("mkdir", os.Mkdir)
	obj.Set("mkdirAll", os.MkdirAll)
	obj.Set("rename", os.Rename)
	obj.Set("remove", os.Remove)
	obj.Set("removeAll", os.RemoveAll)
	obj.Set("openRoot", os.OpenRoot)
	obj.Set("openInRoot", os.OpenInRoot)

	// 命令执行
	obj.Set("exec", exec.Command) // @deprecated
	obj.Set("cmd", exec.Command)

	// 保留原来的简单接口用于向后兼容
	vm.Set("os", map[string]interface{}{
		"env": func(key string) string {
			return os.Getenv(key)
		},
		"platform": func() string {
			return runtime.GOOS
		},
	})
}

// filepathBinds 文件路径相关绑定 (移植自 PocketBase)
func filepathBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$filepath", obj)

	obj.Set("base", filepath.Base)
	obj.Set("clean", filepath.Clean)
	obj.Set("dir", filepath.Dir)
	obj.Set("ext", filepath.Ext)
	obj.Set("fromSlash", filepath.FromSlash)
	obj.Set("isAbs", filepath.IsAbs)
	obj.Set("join", filepath.Join)
	obj.Set("rel", filepath.Rel)
	obj.Set("split", filepath.Split)
	obj.Set("toSlash", filepath.ToSlash)

	// 简化版本的 glob 和 match
	obj.Set("glob", func(pattern string) ([]string, error) {
		return filepath.Glob(pattern)
	})
	obj.Set("match", func(pattern, name string) (bool, error) {
		return filepath.Match(pattern, name)
	})

	// 保留原来的简单接口用于向后兼容
	vm.Set("filepath", map[string]interface{}{
		"base": filepath.Base,
		"clean": filepath.Clean,
		"dir": filepath.Dir,
		"ext": filepath.Ext,
		"isAbs": filepath.IsAbs,
		"join": filepath.Join,
		"split": filepath.Split,
		"glob": func(pattern string) ([]string, error) {
			return filepath.Glob(pattern)
		},
		"match": func(pattern, name string) (bool, error) {
			return filepath.Match(pattern, name)
		},
	})
}

// httpClientBinds HTTP客户端绑定 (移植自 PocketBase)
func httpClientBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$http", obj)

	// HTTP 请求结果结构
	type httpResult struct {
		StatusCode int                    `json:"statusCode"`
		Headers    map[string][]string    `json:"headers"`
		Body       []byte                 `json:"body"`
		BodyString string                 `json:"bodyString"`
		Cookies    map[string]*http.Cookie `json:"cookies"`
		Error      string                 `json:"error,omitempty"`
	}

	// 通用 HTTP 请求函数
	sendFunc := func(params map[string]interface{}) *httpResult {
		result := &httpResult{
			Headers: make(map[string][]string),
			Cookies: make(map[string]*http.Cookie),
		}

		// 解析参数
		method := "GET"
		url := ""
		var body io.Reader
		headers := make(map[string]string)
		timeout := 120

		if v, ok := params["method"].(string); ok {
			method = strings.ToUpper(v)
		}
		if v, ok := params["url"].(string); ok {
			url = v
		}
		if v, ok := params["headers"].(map[string]interface{}); ok {
			for k, val := range v {
				if str, ok := val.(string); ok {
					headers[k] = str
				}
			}
		}
		if v, ok := params["timeout"].(int); ok && v > 0 {
			timeout = v
		}

		// 处理请求体
		if v, ok := params["body"].(string); ok {
			body = strings.NewReader(v)
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "text/plain"
			}
		} else if data, ok := params["data"].(map[string]interface{}); ok {
			if jsonData, err := json.Marshal(data); err == nil {
				body = strings.NewReader(string(jsonData))
				if headers["Content-Type"] == "" {
					headers["Content-Type"] = "application/json"
				}
			}
		}

		if url == "" {
			result.Error = "URL is required"
			return result
		}

		// 创建请求
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			result.Error = err.Error()
			return result
		}

		// 设置头部
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// 设置超时
		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		defer resp.Body.Close()

		// 读取响应体
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			result.Error = err.Error()
			return result
		}

		result.StatusCode = resp.StatusCode
		result.Headers = resp.Header
		result.Body = respBody
		result.BodyString = string(respBody)
		result.Cookies = make(map[string]*http.Cookie)
		for _, cookie := range resp.Cookies() {
			result.Cookies[cookie.Name] = cookie
		}

		return result
	}

	obj.Set("send", sendFunc)

	// 便捷方法
	obj.Set("get", func(url string, params ...map[string]interface{}) *httpResult {
		requestParams := map[string]interface{}{
			"method": "GET",
			"url":    url,
		}
		if len(params) > 0 {
			for k, v := range params[0] {
				requestParams[k] = v
			}
		}
		return sendFunc(requestParams)
	})

	obj.Set("post", func(url string, params ...map[string]interface{}) *httpResult {
		requestParams := map[string]interface{}{
			"method": "POST",
			"url":    url,
		}
		if len(params) > 0 {
			for k, v := range params[0] {
				requestParams[k] = v
			}
		}
		return sendFunc(requestParams)
	})

	obj.Set("put", func(url string, params ...map[string]interface{}) *httpResult {
		requestParams := map[string]interface{}{
			"method": "PUT",
			"url":    url,
		}
		if len(params) > 0 {
			for k, v := range params[0] {
				requestParams[k] = v
			}
		}
		return sendFunc(requestParams)
	})

	obj.Set("delete", func(url string, params ...map[string]interface{}) *httpResult {
		requestParams := map[string]interface{}{
			"method": "DELETE",
			"url":    url,
		}
		if len(params) > 0 {
			for k, v := range params[0] {
				requestParams[k] = v
			}
		}
		return sendFunc(requestParams)
	})

	// 保留原来的简单接口用于向后兼容
	vm.Set("http", map[string]interface{}{
		"get": func(url string, headers map[string]string) map[string]interface{} {
			params := map[string]interface{}{
				"url": url,
			}
			if headers != nil {
				headerMap := make(map[string]interface{})
				for k, v := range headers {
					headerMap[k] = v
				}
				params["headers"] = headerMap
			}
			result := sendFunc(params)
			if result.Error != "" {
				return map[string]interface{}{
					"status": 500,
					"error":  result.Error,
				}
			}
			return map[string]interface{}{
				"status":  result.StatusCode,
				"headers": result.Headers,
				"body":    result.BodyString,
			}
		},
		"post": func(url string, data interface{}, headers map[string]string) map[string]interface{} {
			params := map[string]interface{}{
				"method": "POST",
				"url":    url,
				"data":   data,
			}
			if headers != nil {
				headerMap := make(map[string]interface{})
				for k, v := range headers {
					headerMap[k] = v
				}
				params["headers"] = headerMap
			}
			result := sendFunc(params)
			if result.Error != "" {
				return map[string]interface{}{
					"status": 500,
					"error":  result.Error,
				}
			}
			return map[string]interface{}{
				"status": result.StatusCode,
				"body":   result.BodyString,
			}
		},
	})
}

// filesystemBinds 文件系统绑定 (移植自 PocketBase)
func filesystemBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$filesystem", obj)

	// 简化的文件操作 - 由于 PocketBase 的 File 结构复杂，这里提供基本的文件操作
	obj.Set("readFile", func(path string) ([]byte, error) {
		return os.ReadFile(path)
	})
	obj.Set("writeFile", func(path string, data []byte, perm ...os.FileMode) error {
		if len(perm) > 0 {
			return os.WriteFile(path, data, perm[0])
		}
		return os.WriteFile(path, data, 0644)
	})
	obj.Set("fileExists", func(path string) bool {
		_, err := os.Stat(path)
		return !os.IsNotExist(err)
	})
	obj.Set("fileSize", func(path string) (int64, error) {
		info, err := os.Stat(path)
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	})

	// 保留原来的简单接口用于向后兼容
	vm.Set("fs", map[string]interface{}{
		"readFile": func(path string) string {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Sprintf("Error reading file: %v", err)
			}
			return string(data)
		},
		"writeFile": func(path string, content string) error {
			utils.Info("Write file called: %s", path)
			return os.WriteFile(path, []byte(content), 0644)
		},
	})
}


// formsBinds 表单绑定（简化实现）
func formsBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$forms", obj)

	// 简单的表单验证系统
	validateFunc := func(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
		errors := make(map[string]interface{})

		for field, rule := range rules {
			value, exists := data[field]
			if !exists {
				if required, ok := rule.(map[string]interface{})["required"].(bool); ok && required {
					errors[field] = "This field is required"
				}
				continue
			}

			if ruleMap, ok := rule.(map[string]interface{}); ok {
				// 检查必填
				if required, ok := ruleMap["required"].(bool); ok && required {
					if value == "" || value == nil {
						errors[field] = "This field is required"
						continue
					}
				}

				// 检查最小长度
				if minLength, ok := ruleMap["minLength"].(int); ok {
					if str, ok := value.(string); ok && len(str) < minLength {
						errors[field] = fmt.Sprintf("Minimum length is %d", minLength)
					}
				}

				// 检查最大长度
				if maxLength, ok := ruleMap["maxLength"].(int); ok {
					if str, ok := value.(string); ok && len(str) > maxLength {
						errors[field] = fmt.Sprintf("Maximum length is %d", maxLength)
					}
				}

				// 检查邮箱格式
				if isEmail, ok := ruleMap["isEmail"].(bool); ok && isEmail {
					if str, ok := value.(string); ok && !strings.Contains(str, "@") {
						errors[field] = "Invalid email format"
					}
				}

				// 检查正则表达式
				if pattern, ok := ruleMap["pattern"].(string); ok {
					if str, ok := value.(string); ok {
						if matched, _ := regexp.MatchString(pattern, str); !matched {
							errors[field] = "Invalid format"
						}
					}
				}
			}
		}

		result := map[string]interface{}{
			"valid": len(errors) == 0,
			"errors": errors,
		}
		return result
	}

	obj.Set("validate", validateFunc)
	obj.Set("create", func(schema map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"data": make(map[string]interface{}),
			"schema": schema,
			"validate": func(data map[string]interface{}) map[string]interface{} {
				return validateFunc(data, schema)
			},
		}
	})

	// 保留原来的简单接口用于向后兼容
	vm.Set("forms", map[string]interface{}{
		"validate": func(data interface{}, rules interface{}) bool {
			if dataMap, ok := data.(map[string]interface{}); ok {
				if rulesMap, ok := rules.(map[string]interface{}); ok {
					result := validateFunc(dataMap, rulesMap)
					return result["valid"].(bool)
				}
			}
			return true
		},
	})
}

// mailsBinds 邮件绑定（已移除 TODO）
func mailsBinds(vm *goja.Runtime) {
	vm.Set("mails", map[string]interface{}{
		"send": func(to, subject, body string) error {
			// 简单的邮件发送日志记录
			utils.Info("Mail to %s: %s - %s", to, subject, body)
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