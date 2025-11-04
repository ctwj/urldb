package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/plugin"
	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/monitor"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// IntegrationTestSuite provides a complete integration test environment
type IntegrationTestSuite struct {
	DB              *gorm.DB
	RepoManager     *repo.RepositoryManager
	PluginManager   *manager.Manager
	TaskManager     *task.TaskManager
	PluginMonitor   *monitor.PluginMonitor
	TestPlugins     map[string]types.Plugin
	CleanupFuncs    []func()
	tempDBFile      string
}

// NewIntegrationTestSuite creates a new integration test suite
func NewIntegrationTestSuite() *IntegrationTestSuite {
	return &IntegrationTestSuite{
		TestPlugins: make(map[string]types.Plugin),
		CleanupFuncs: make([]func(), 0),
	}
}

// Setup initializes the integration test environment
func (its *IntegrationTestSuite) Setup(t *testing.T) {
	// Initialize logger
	utils.InitLogger(nil)

	// Setup test database
	its.setupTestDatabase(t)

	// Setup repository manager
	its.setupRepositoryManager()

	// Setup task manager
	its.setupTaskManager()

	// Setup plugin monitor
	its.setupPluginMonitor()

	// Setup plugin manager
	its.setupPluginManager()

	// Setup test data
	its.setupTestData(t)
}

// setupTestDatabase sets up a test database
func (its *IntegrationTestSuite) setupTestDatabase(t *testing.T) {
	// For integration tests, we'll use an in-memory SQLite database
	// In a real scenario, you might want to use a separate test database

	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "urldb_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database file: %v", err)
	}
	its.tempDBFile = tempFile.Name()
	tempFile.Close()

	// Set environment variables for test database
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "url_db_test")

	// Initialize database connection
	err = db.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	its.DB = db.DB

	// Add cleanup function
	its.CleanupFuncs = append(its.CleanupFuncs, func() {
		if its.DB != nil {
			sqlDB, _ := its.DB.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
		os.Remove(its.tempDBFile)
	})
}

// setupRepositoryManager sets up the repository manager
func (its *IntegrationTestSuite) setupRepositoryManager() {
	its.RepoManager = repo.NewRepositoryManager(its.DB)
}

// setupTaskManager sets up the task manager
func (its *IntegrationTestSuite) setupTaskManager() {
	// Create a simple task manager for testing
	its.TaskManager = task.NewTaskManager(its.RepoManager)
}

// setupPluginMonitor sets up the plugin monitor
func (its *IntegrationTestSuite) setupPluginMonitor() {
	its.PluginMonitor = monitor.NewPluginMonitor()
}

// setupPluginManager sets up the plugin manager
func (its *IntegrationTestSuite) setupPluginManager() {
	its.PluginManager = manager.NewManager(its.TaskManager, its.RepoManager, its.DB, its.PluginMonitor)
	plugin.GlobalManager = its.PluginManager
}

// setupTestData sets up initial test data
func (its *IntegrationTestSuite) setupTestData(t *testing.T) {
	// Create tables
	if err := its.DB.AutoMigrate(
		&entity.Resource{},
		&entity.ReadyResource{},
		&entity.Pan{},
		&entity.Cks{},
		&entity.Category{},
		&entity.Tag{},
		&entity.User{},
		&entity.SystemConfig{},
		&entity.Task{},
		&entity.TaskItem{},
		&entity.PluginConfig{},
		&entity.PluginData{},
	); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	// Insert some test data
	testUser := &entity.User{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
		Role:     "admin",
	}

	if err := its.RepoManager.UserRepository.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Add test pan data
	testPan := &entity.Pan{
		Name: "testpan",
		Key:  99,
		Icon: "<i class=\"fas fa-cloud text-blue-500\"></i>",
		Remark: "Test Pan",
	}

	if err := its.DB.Create(testPan).Error; err != nil {
		t.Fatalf("Failed to create test pan: %v", err)
	}

	// Add test cks data
	testCks := &entity.Cks{
		PanID:   testPan.ID,
		Ck:      "test_cookie",
		IsValid: true,
		Remark:  "Test Ck",
	}

	if err := its.DB.Create(testCks).Error; err != nil {
		t.Fatalf("Failed to create test cks: %v", err)
	}

	// Add more test data as needed
}

// Teardown cleans up the integration test environment
func (its *IntegrationTestSuite) Teardown() {
	// Run cleanup functions in reverse order
	for i := len(its.CleanupFuncs) - 1; i >= 0; i-- {
		its.CleanupFuncs[i]()
	}
}

// RegisterPlugin registers a plugin for integration testing
func (its *IntegrationTestSuite) RegisterPlugin(plugin types.Plugin) error {
	its.TestPlugins[plugin.Name()] = plugin
	return its.PluginManager.RegisterPlugin(plugin)
}

// GetPluginManager returns the plugin manager
func (its *IntegrationTestSuite) GetPluginManager() *manager.Manager {
	return its.PluginManager
}

// GetRepositoryManager returns the repository manager
func (its *IntegrationTestSuite) GetRepositoryManager() *repo.RepositoryManager {
	return its.RepoManager
}

// GetDB returns the database connection
func (its *IntegrationTestSuite) GetDB() *gorm.DB {
	return its.DB
}

// RunPluginIntegrationTest runs a complete integration test for a plugin
func (its *IntegrationTestSuite) RunPluginIntegrationTest(t *testing.T, pluginName string, config map[string]interface{}) {
	// Initialize plugin
	t.Run("Initialize", func(t *testing.T) {
		err := its.PluginManager.InitializePlugin(pluginName, config)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}
	})

	// Start plugin
	t.Run("Start", func(t *testing.T) {
		err := its.PluginManager.StartPlugin(pluginName)
		if err != nil {
			t.Fatalf("Failed to start plugin: %v", err)
		}

		// Check plugin status
		status := its.PluginManager.GetPluginStatus(pluginName)
		if status != types.StatusRunning {
			t.Errorf("Expected plugin status to be running, got %s", status)
		}
	})

	// Test plugin functionality
	t.Run("Functionality", func(t *testing.T) {
		// This would depend on the specific plugin being tested
		// For now, we'll just wait a bit to allow any scheduled tasks to run
		time.Sleep(100 * time.Millisecond)
	})

	// Stop plugin
	t.Run("Stop", func(t *testing.T) {
		err := its.PluginManager.StopPlugin(pluginName)
		if err != nil {
			t.Fatalf("Failed to stop plugin: %v", err)
		}

		// Check plugin status
		status := its.PluginManager.GetPluginStatus(pluginName)
		if status != types.StatusStopped {
			t.Errorf("Expected plugin status to be stopped, got %s", status)
		}
	})

	// Cleanup plugin
	t.Run("Cleanup", func(t *testing.T) {
		// Get plugin instance to access context
		plugin, err := its.PluginManager.GetPlugin(pluginName)
		if err != nil {
			t.Fatalf("Failed to get plugin: %v", err)
		}

		err = plugin.Cleanup()
		if err != nil {
			t.Fatalf("Failed to cleanup plugin: %v", err)
		}
	})
}

// CreateMockPlugin creates a mock plugin for testing
func (its *IntegrationTestSuite) CreateMockPlugin(name, version string) *MockPlugin {
	return &MockPlugin{
		name:        name,
		version:     version,
		description: fmt.Sprintf("Mock plugin for testing: %s", name),
		author:      "Test Suite",
	}
}

// MockPlugin is a mock plugin implementation for testing
type MockPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
	started     bool
	stopped     bool
	cleaned     bool
}

// Name returns the plugin name
func (p *MockPlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *MockPlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *MockPlugin) Description() string {
	return p.description
}

// Author returns the plugin author
func (p *MockPlugin) Author() string {
	return p.author
}

// Initialize initializes the plugin
func (p *MockPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Mock plugin initialized: %s", p.name)
	return nil
}

// Start starts the plugin
func (p *MockPlugin) Start() error {
	p.started = true
	if p.context != nil {
		p.context.LogInfo("Mock plugin started: %s", p.name)
	}
	return nil
}

// Stop stops the plugin
func (p *MockPlugin) Stop() error {
	p.stopped = true
	if p.context != nil {
		p.context.LogInfo("Mock plugin stopped: %s", p.name)
	}
	return nil
}

// Cleanup cleans up the plugin
func (p *MockPlugin) Cleanup() error {
	p.cleaned = true
	if p.context != nil {
		p.context.LogInfo("Mock plugin cleaned up: %s", p.name)
	}
	return nil
}

// Dependencies returns the plugin dependencies
func (p *MockPlugin) Dependencies() []string {
	return []string{}
}

// CheckDependencies checks the plugin dependencies
func (p *MockPlugin) CheckDependencies() map[string]bool {
	return make(map[string]bool)
}

// IsStarted returns whether the plugin has been started
func (p *MockPlugin) IsStarted() bool {
	return p.started
}

// IsStopped returns whether the plugin has been stopped
func (p *MockPlugin) IsStopped() bool {
	return p.stopped
}

// IsCleaned returns whether the plugin has been cleaned up
func (p *MockPlugin) IsCleaned() bool {
	return p.cleaned
}

// GetContext returns the plugin context
func (p *MockPlugin) GetContext() types.PluginContext {
	return p.context
}

// WithContext sets the plugin context
func (p *MockPlugin) WithContext(ctx types.PluginContext) *MockPlugin {
	p.context = ctx
	return p
}

// WithErrorOnStart configures the plugin to return an error on Start
func (p *MockPlugin) WithErrorOnStart() *MockPluginWithError {
	return &MockPluginWithError{
		MockPlugin: p,
		errorOn:    "start",
	}
}

// WithErrorOnStop configures the plugin to return an error on Stop
func (p *MockPlugin) WithErrorOnStop() *MockPluginWithError {
	return &MockPluginWithError{
		MockPlugin: p,
		errorOn:    "stop",
	}
}

// WithErrorOnInitialize configures the plugin to return an error on Initialize
func (p *MockPlugin) WithErrorOnInitialize() *MockPluginWithError {
	return &MockPluginWithError{
		MockPlugin: p,
		errorOn:    "initialize",
	}
}

// MockPluginWithError is a mock plugin that returns errors
type MockPluginWithError struct {
	*MockPlugin
	errorOn string
}

// Initialize initializes the plugin with error
func (p *MockPluginWithError) Initialize(ctx types.PluginContext) error {
	if p.errorOn == "initialize" {
		return fmt.Errorf("intentional initialize error")
	}
	return p.MockPlugin.Initialize(ctx)
}

// Start starts the plugin with error
func (p *MockPluginWithError) Start() error {
	if p.errorOn == "start" {
		return fmt.Errorf("intentional start error")
	}
	return p.MockPlugin.Start()
}

// Stop stops the plugin with error
func (p *MockPluginWithError) Stop() error {
	if p.errorOn == "stop" {
		return fmt.Errorf("intentional stop error")
	}
	return p.MockPlugin.Stop()
}

// MockPluginWithDependencies is a mock plugin with dependencies
type MockPluginWithDependencies struct {
	*MockPlugin
	dependencies []string
}

// WithDependencies creates a mock plugin with dependencies
func (p *MockPlugin) WithDependencies(deps []string) *MockPluginWithDependencies {
	return &MockPluginWithDependencies{
		MockPlugin:   p,
		dependencies: deps,
	}
}

// Dependencies returns the plugin dependencies
func (p *MockPluginWithDependencies) Dependencies() []string {
	return p.dependencies
}

// CheckDependencies checks the plugin dependencies
func (p *MockPluginWithDependencies) CheckDependencies() map[string]bool {
	result := make(map[string]bool)
	for _, dep := range p.dependencies {
		// In a real implementation, this would check if the dependency is satisfied
		// For testing, we'll assume all dependencies are satisfied
		result[dep] = true
	}
	return result
}

// MockPluginWithContextOperations is a mock plugin that performs context operations
type MockPluginWithContextOperations struct {
	*MockPlugin
	operations []string
}

// WithContextOperations creates a mock plugin that performs context operations
func (p *MockPlugin) WithContextOperations(ops []string) *MockPluginWithContextOperations {
	return &MockPluginWithContextOperations{
		MockPlugin: p,
		operations: ops,
	}
}

// Start starts the plugin and performs context operations
func (p *MockPluginWithContextOperations) Start() error {
	if err := p.MockPlugin.Start(); err != nil {
		return err
	}

	if p.context != nil {
		for _, op := range p.operations {
			switch op {
			case "log_debug":
				p.context.LogDebug("Debug message")
			case "log_info":
				p.context.LogInfo("Info message")
			case "log_warn":
				p.context.LogWarn("Warning message")
			case "log_error":
				p.context.LogError("Error message")
			case "set_config":
				p.context.SetConfig("test_key", "test_value")
			case "get_config":
				p.context.GetConfig("test_key")
			case "set_data":
				p.context.SetData("test_key", "test_value", "test_type")
			case "get_data":
				p.context.GetData("test_key", "test_type")
			case "register_task":
				p.context.RegisterTask("test_task", func() {})
			case "cache_set":
				p.context.CacheSet("test_key", "test_value", time.Minute)
			case "cache_get":
				p.context.CacheGet("test_key")
			}
		}
	}

	return nil
}