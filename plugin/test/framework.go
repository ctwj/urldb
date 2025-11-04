package test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// TestPluginContext is a mock implementation of PluginContext for testing
type TestPluginContext struct {
	mu           sync.RWMutex
	logs         []LogEntry
	config       map[string]interface{}
	data         map[string]map[string]interface{} // key: dataType, value: {key: value}
	cache        map[string]CacheEntry
	permissions  map[string]bool
	tasks        map[string]func()
	db           interface{}
	concurrencyLimit int
}

// LogEntry represents a log entry
type LogEntry struct {
	Level   string
	Message string
	Args    []interface{}
	Time    time.Time
}

// CacheEntry represents a cache entry
type CacheEntry struct {
	Value interface{}
	Expiry time.Time
}

// NewTestPluginContext creates a new test plugin context
func NewTestPluginContext() *TestPluginContext {
	return &TestPluginContext{
		logs:        make([]LogEntry, 0),
		config:      make(map[string]interface{}),
		data:        make(map[string]map[string]interface{}),
		cache:       make(map[string]CacheEntry),
		permissions: make(map[string]bool),
		tasks:       make(map[string]func()),
	}
}

// SetDB sets the database for testing
func (ctx *TestPluginContext) SetDB(db interface{}) {
	ctx.db = db
}

// Logging methods
func (ctx *TestPluginContext) LogDebug(msg string, args ...interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.logs = append(ctx.logs, LogEntry{
		Level:   "DEBUG",
		Message: msg,
		Args:    args,
		Time:    time.Now(),
	})
}

func (ctx *TestPluginContext) LogInfo(msg string, args ...interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.logs = append(ctx.logs, LogEntry{
		Level:   "INFO",
		Message: msg,
		Args:    args,
		Time:    time.Now(),
	})
}

func (ctx *TestPluginContext) LogWarn(msg string, args ...interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.logs = append(ctx.logs, LogEntry{
		Level:   "WARN",
		Message: msg,
		Args:    args,
		Time:    time.Now(),
	})
}

func (ctx *TestPluginContext) LogError(msg string, args ...interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.logs = append(ctx.logs, LogEntry{
		Level:   "ERROR",
		Message: msg,
		Args:    args,
		Time:    time.Now(),
	})
}

// GetLogs returns all logs
func (ctx *TestPluginContext) GetLogs() []LogEntry {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	logs := make([]LogEntry, len(ctx.logs))
	copy(logs, ctx.logs)
	return logs
}

// ClearLogs clears all logs
func (ctx *TestPluginContext) ClearLogs() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.logs = make([]LogEntry, 0)
}

// Configuration methods
func (ctx *TestPluginContext) GetConfig(key string) (interface{}, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if val, exists := ctx.config[key]; exists {
		return val, nil
	}
	return nil, fmt.Errorf("config key '%s' not found", key)
}

func (ctx *TestPluginContext) SetConfig(key string, value interface{}) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.config[key] = value
	return nil
}

// Data methods
func (ctx *TestPluginContext) GetData(key string, dataType string) (interface{}, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if dataMap, exists := ctx.data[dataType]; exists {
		if val, exists := dataMap[key]; exists {
			return val, nil
		}
	}
	return nil, fmt.Errorf("data key '%s' not found for type '%s'", key, dataType)
}

func (ctx *TestPluginContext) SetData(key string, value interface{}, dataType string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if _, exists := ctx.data[dataType]; !exists {
		ctx.data[dataType] = make(map[string]interface{})
	}
	ctx.data[dataType][key] = value
	return nil
}

func (ctx *TestPluginContext) DeleteData(key string, dataType string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if dataMap, exists := ctx.data[dataType]; exists {
		delete(dataMap, key)
	}
	return nil
}

// Task scheduling methods
func (ctx *TestPluginContext) RegisterTask(name string, task func()) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if _, exists := ctx.tasks[name]; exists {
		return fmt.Errorf("task '%s' already registered", name)
	}
	ctx.tasks[name] = task
	return nil
}

func (ctx *TestPluginContext) UnregisterTask(name string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if _, exists := ctx.tasks[name]; !exists {
		return fmt.Errorf("task '%s' not found", name)
	}
	delete(ctx.tasks, name)
	return nil
}

// GetTask returns a registered task
func (ctx *TestPluginContext) GetTask(name string) (func(), error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if task, exists := ctx.tasks[name]; exists {
		return task, nil
	}
	return nil, fmt.Errorf("task '%s' not found", name)
}

// GetTasks returns all registered tasks
func (ctx *TestPluginContext) GetTasks() map[string]func() {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	tasks := make(map[string]func())
	for name, task := range ctx.tasks {
		tasks[name] = task
	}
	return tasks
}

// Database access
func (ctx *TestPluginContext) GetDB() interface{} {
	return ctx.db
}

// Security methods
func (ctx *TestPluginContext) CheckPermission(permissionType string, resource ...string) (bool, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	key := permissionType
	if len(resource) > 0 {
		key = fmt.Sprintf("%s:%s", permissionType, resource[0])
	}
	if val, exists := ctx.permissions[key]; exists {
		return val, nil
	}
	return false, nil
}

func (ctx *TestPluginContext) RequestPermission(permissionType string, resource string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	key := fmt.Sprintf("%s:%s", permissionType, resource)
	ctx.permissions[key] = true
	return nil
}

func (ctx *TestPluginContext) GetSecurityReport() (interface{}, error) {
	return map[string]interface{}{
		"permissions": ctx.permissions,
	}, nil
}

// Cache methods
func (ctx *TestPluginContext) CacheSet(key string, value interface{}, ttl time.Duration) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	expiry := time.Now().Add(ttl)
	ctx.cache[key] = CacheEntry{
		Value:  value,
		Expiry: expiry,
	}
	return nil
}

func (ctx *TestPluginContext) CacheGet(key string) (interface{}, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if entry, exists := ctx.cache[key]; exists {
		if time.Now().Before(entry.Expiry) {
			return entry.Value, nil
		}
		// Expired, remove it
		ctx.mu.RUnlock()
		ctx.mu.Lock()
		delete(ctx.cache, key)
		ctx.mu.Unlock()
		ctx.mu.RLock()
	}
	return nil, fmt.Errorf("cache key '%s' not found or expired", key)
}

func (ctx *TestPluginContext) CacheDelete(key string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	delete(ctx.cache, key)
	return nil
}

func (ctx *TestPluginContext) CacheClear() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.cache = make(map[string]CacheEntry)
	return nil
}

// Concurrency methods
func (ctx *TestPluginContext) ConcurrencyExecute(ctx2 context.Context, taskFunc func() error) error {
	// For testing, we execute directly without concurrency control
	return taskFunc()
}

func (ctx *TestPluginContext) SetConcurrencyLimit(limit int) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.concurrencyLimit = limit
	return nil
}

func (ctx *TestPluginContext) GetConcurrencyStats() (map[string]interface{}, error) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return map[string]interface{}{
		"limit": ctx.concurrencyLimit,
	}, nil
}

// TestPluginManager is a test helper for managing plugin lifecycle in tests
type TestPluginManager struct {
	plugins map[string]types.Plugin
	contexts map[string]*TestPluginContext
}

// NewTestPluginManager creates a new test plugin manager
func NewTestPluginManager() *TestPluginManager {
	return &TestPluginManager{
		plugins:  make(map[string]types.Plugin),
		contexts: make(map[string]*TestPluginContext),
	}
}

// RegisterPlugin registers a plugin for testing
func (tpm *TestPluginManager) RegisterPlugin(plugin types.Plugin) error {
	name := plugin.Name()
	if _, exists := tpm.plugins[name]; exists {
		return fmt.Errorf("plugin '%s' already registered", name)
	}
	tpm.plugins[name] = plugin
	return nil
}

// GetPlugin returns a registered plugin
func (tpm *TestPluginManager) GetPlugin(name string) (types.Plugin, error) {
	if plugin, exists := tpm.plugins[name]; exists {
		return plugin, nil
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}

// GetContext returns a test context for a plugin
func (tpm *TestPluginManager) GetContext(pluginName string) *TestPluginContext {
	if ctx, exists := tpm.contexts[pluginName]; exists {
		return ctx
	}
	ctx := NewTestPluginContext()
	tpm.contexts[pluginName] = ctx
	return ctx
}

// InitializePlugin initializes a plugin for testing
func (tpm *TestPluginManager) InitializePlugin(t *testing.T, pluginName string, config map[string]interface{}) error {
	plugin, err := tpm.GetPlugin(pluginName)
	if err != nil {
		return err
	}

	ctx := tpm.GetContext(pluginName)
	for key, value := range config {
		ctx.SetConfig(key, value)
	}

	if err := plugin.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize plugin '%s': %v", pluginName, err)
	}

	return nil
}

// StartPlugin starts a plugin for testing
func (tpm *TestPluginManager) StartPlugin(t *testing.T, pluginName string) error {
	plugin, err := tpm.GetPlugin(pluginName)
	if err != nil {
		return err
	}

	if err := plugin.Start(); err != nil {
		return fmt.Errorf("failed to start plugin '%s': %v", pluginName, err)
	}

	return nil
}

// StopPlugin stops a plugin for testing
func (tpm *TestPluginManager) StopPlugin(t *testing.T, pluginName string) error {
	plugin, err := tpm.GetPlugin(pluginName)
	if err != nil {
		return err
	}

	if err := plugin.Stop(); err != nil {
		return fmt.Errorf("failed to stop plugin '%s': %v", pluginName, err)
	}

	return nil
}

// CleanupPlugin cleans up a plugin for testing
func (tpm *TestPluginManager) CleanupPlugin(t *testing.T, pluginName string) error {
	plugin, err := tpm.GetPlugin(pluginName)
	if err != nil {
		return err
	}

	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup plugin '%s': %v", pluginName, err)
	}

	return nil
}

// RunPluginLifecycle runs the complete plugin lifecycle for testing
func (tpm *TestPluginManager) RunPluginLifecycle(t *testing.T, pluginName string, config map[string]interface{}) error {
	// Initialize
	if err := tpm.InitializePlugin(t, pluginName, config); err != nil {
		return err
	}

	// Start
	if err := tpm.StartPlugin(t, pluginName); err != nil {
		return err
	}

	// Stop
	if err := tpm.StopPlugin(t, pluginName); err != nil {
		return err
	}

	// Cleanup
	if err := tpm.CleanupPlugin(t, pluginName); err != nil {
		return err
	}

	return nil
}

// AssertLogContains checks if a log message contains the specified text
func (ctx *TestPluginContext) AssertLogContains(t *testing.T, level string, contains string) bool {
	logs := ctx.GetLogs()
	for _, log := range logs {
		if log.Level == level && containsInMessage(log.Message, log.Args, contains) {
			return true
		}
	}
	return false
}

// containsInMessage checks if the formatted message contains the specified text
func containsInMessage(message string, args []interface{}, contains string) bool {
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	return strings.Contains(formatted, contains)
}