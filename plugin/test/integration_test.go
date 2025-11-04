package test

import (
	"testing"

	"github.com/ctwj/urldb/plugin/demo"
)

// TestIntegrationSuite tests the integration test suite
func TestIntegrationSuite(t *testing.T) {
	// Create integration test suite
	suite := NewIntegrationTestSuite()

	// Setup
	suite.Setup(t)
	defer suite.Teardown()

	// Test suite setup
	if suite.DB == nil {
		t.Error("Database not initialized")
	}

	if suite.RepoManager == nil {
		t.Error("Repository manager not initialized")
	}

	if suite.PluginManager == nil {
		t.Error("Plugin manager not initialized")
	}

	// Test registering a plugin
	mockPlugin := suite.CreateMockPlugin("test-plugin", "1.0.0")
	if err := suite.RegisterPlugin(mockPlugin); err != nil {
		t.Errorf("Failed to register mock plugin: %v", err)
	}

	// Test getting plugin manager
	if suite.GetPluginManager() == nil {
		t.Error("Failed to get plugin manager")
	}

	// Test getting repository manager
	if suite.GetRepositoryManager() == nil {
		t.Error("Failed to get repository manager")
	}

	// Test getting database
	if suite.GetDB() == nil {
		t.Error("Failed to get database")
	}
}

// TestPluginIntegration tests a complete plugin integration
func TestPluginIntegration(t *testing.T) {
	// Create integration test suite
	suite := NewIntegrationTestSuite()

	// Setup
	suite.Setup(t)
	defer suite.Teardown()

	// Create and register demo plugin
	demoPlugin := demo.NewDemoPlugin()
	if err := suite.RegisterPlugin(demoPlugin); err != nil {
		t.Fatalf("Failed to register demo plugin: %v", err)
	}

	// Run integration test
	config := map[string]interface{}{
		"interval": "1s",
	}

	suite.RunPluginIntegrationTest(t, demoPlugin.Name(), config)
}

// TestMockPlugin tests the mock plugin functionality
func TestMockPlugin(t *testing.T) {
	// Test basic mock plugin
	plugin := NewIntegrationTestSuite().CreateMockPlugin("test-mock", "1.0.0")

	if plugin.Name() != "test-mock" {
		t.Errorf("Expected name 'test-mock', got '%s'", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}

	// Test plugin lifecycle
	ctx := NewTestPluginContext()
	if err := plugin.Initialize(ctx); err != nil {
		t.Errorf("Failed to initialize mock plugin: %v", err)
	}

	if err := plugin.Start(); err != nil {
		t.Errorf("Failed to start mock plugin: %v", err)
	}

	if !plugin.IsStarted() {
		t.Error("Plugin should be marked as started")
	}

	if err := plugin.Stop(); err != nil {
		t.Errorf("Failed to stop mock plugin: %v", err)
	}

	if !plugin.IsStopped() {
		t.Error("Plugin should be marked as stopped")
	}

	if err := plugin.Cleanup(); err != nil {
		t.Errorf("Failed to cleanup mock plugin: %v", err)
	}

	if !plugin.IsCleaned() {
		t.Error("Plugin should be marked as cleaned")
	}
}

// TestMockPluginWithError tests mock plugin with errors
func TestMockPluginWithError(t *testing.T) {
	// Test plugin with initialize error
	pluginWithError := NewIntegrationTestSuite().CreateMockPlugin("error-plugin", "1.0.0").WithErrorOnInitialize()
	ctx := NewTestPluginContext()

	if err := pluginWithError.Initialize(ctx); err == nil {
		t.Error("Expected initialize error")
	}

	// Test plugin with start error
	pluginWithStartError := NewIntegrationTestSuite().CreateMockPlugin("start-error-plugin", "1.0.0").WithErrorOnStart()

	if err := pluginWithStartError.Initialize(ctx); err != nil {
		t.Errorf("Failed to initialize plugin: %v", err)
	}

	if err := pluginWithStartError.Start(); err == nil {
		t.Error("Expected start error")
	}

	// Test plugin with stop error
	pluginWithStopError := NewIntegrationTestSuite().CreateMockPlugin("stop-error-plugin", "1.0.0").WithErrorOnStop()

	if err := pluginWithStopError.Initialize(ctx); err != nil {
		t.Errorf("Failed to initialize plugin: %v", err)
	}

	if err := pluginWithStopError.Start(); err != nil {
		t.Errorf("Failed to start plugin: %v", err)
	}

	if err := pluginWithStopError.Stop(); err == nil {
		t.Error("Expected stop error")
	}
}

// TestMockPluginWithDependencies tests mock plugin with dependencies
func TestMockPluginWithDependencies(t *testing.T) {
	// Create plugin with dependencies
	plugin := NewIntegrationTestSuite().CreateMockPlugin("dep-plugin", "1.0.0").WithDependencies([]string{"dep1", "dep2"})

	deps := plugin.Dependencies()
	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	checkDeps := plugin.CheckDependencies()
	if len(checkDeps) != 2 {
		t.Errorf("Expected 2 dependency checks, got %d", len(checkDeps))
	}

	// Verify all dependencies are marked as satisfied
	for _, dep := range []string{"dep1", "dep2"} {
		if satisfied, exists := checkDeps[dep]; !exists || !satisfied {
			t.Errorf("Dependency %s should be satisfied", dep)
		}
	}
}

// TestMockPluginWithContextOperations tests mock plugin with context operations
func TestMockPluginWithContextOperations(t *testing.T) {
	// Create plugin with context operations
	operations := []string{
		"log_info",
		"set_config",
		"set_data",
		"register_task",
		"cache_set",
	}

	plugin := NewIntegrationTestSuite().CreateMockPlugin("context-op-plugin", "1.0.0").WithContextOperations(operations)
	ctx := NewTestPluginContext()

	// Initialize plugin
	if err := plugin.Initialize(ctx); err != nil {
		t.Errorf("Failed to initialize plugin: %v", err)
	}

	// Start plugin (this will execute context operations)
	if err := plugin.Start(); err != nil {
		t.Errorf("Failed to start plugin: %v", err)
	}

	// Verify context operations were performed
	// Check logs
	logs := ctx.GetLogs()
	if len(logs) == 0 {
		t.Error("Expected logs from context operations")
	}

	// Check config was set
	if val, err := ctx.GetConfig("test_key"); err != nil || val != "test_value" {
		t.Error("Expected config to be set")
	}

	// Check data was set
	if val, err := ctx.GetData("test_key", "test_type"); err != nil || val != "test_value" {
		t.Error("Expected data to be set")
	}

	// Check task was registered
	if _, err := ctx.GetTask("test_task"); err != nil {
		t.Error("Expected task to be registered")
	}

	// Check cache was set
	if _, err := ctx.CacheGet("test_key"); err != nil {
		t.Error("Expected cache item to be set")
	}
}