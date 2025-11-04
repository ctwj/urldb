package test

import (
	"testing"

	"github.com/ctwj/urldb/plugin/demo"
	"github.com/ctwj/urldb/plugin/types"
)

// ExampleTest demonstrates how to use the plugin test framework
func ExampleTest(t *testing.T) {
	// Create a test reporter
	reporter := NewTestReporter("ExamplePluginTest")
	wrapper := NewTestingTWrapper(t, reporter)

	// Run tests using the wrapper
	wrapper.Run("ExamplePluginBasicTest", func(t *testing.T) {
		// Create plugin instance
		plugin := demo.NewDemoPlugin()

		// Test basic plugin properties
		if plugin.Name() != "demo-plugin" {
			t.Errorf("Expected plugin name 'demo-plugin', got '%s'", plugin.Name())
		}

		if plugin.Version() != "1.0.0" {
			t.Errorf("Expected plugin version '1.0.0', got '%s'", plugin.Version())
		}

		// Create test context
		ctx := NewTestPluginContext()

		// Test plugin initialization
		if err := plugin.Initialize(ctx); err != nil {
			t.Errorf("Failed to initialize plugin: %v", err)
		}

		// Verify initialization logs
		if !ctx.AssertLogContains(t, "INFO", "Demo plugin initialized") {
			t.Error("Expected initialization log message")
		}

		// Test plugin start
		if err := plugin.Start(); err != nil {
			t.Errorf("Failed to start plugin: %v", err)
		}

		// Verify start logs
		if !ctx.AssertLogContains(t, "INFO", "Demo plugin started") {
			t.Error("Expected start log message")
		}

		// Test plugin stop
		if err := plugin.Stop(); err != nil {
			t.Errorf("Failed to stop plugin: %v", err)
		}

		// Verify stop logs
		if !ctx.AssertLogContains(t, "INFO", "Demo plugin stopped") {
			t.Error("Expected stop log message")
		}

		// Test plugin cleanup
		if err := plugin.Cleanup(); err != nil {
			t.Errorf("Failed to cleanup plugin: %v", err)
		}

		// Verify cleanup logs
		if !ctx.AssertLogContains(t, "INFO", "Demo plugin cleaned up") {
			t.Error("Expected cleanup log message")
		}
	})

	// Generate and display reports
	textReport := reporter.GenerateTextReport()
	t.Logf("Text Report:\n%s", textReport)

	jsonReport, err := reporter.GenerateJSONReport()
	if err != nil {
		t.Errorf("Failed to generate JSON report: %v", err)
	} else {
		t.Logf("JSON Report: %s", jsonReport)
	}
}

// ExampleIntegrationTest demonstrates how to use the integration test suite
func ExampleIntegrationTest(t *testing.T) {
	// Create integration test suite
	suite := NewIntegrationTestSuite()

	// Setup the test environment
	suite.Setup(t)
	defer suite.Teardown()

	// Create and register a demo plugin
	demoPlugin := demo.NewDemoPlugin()
	if err := suite.RegisterPlugin(demoPlugin); err != nil {
		t.Fatalf("Failed to register demo plugin: %v", err)
	}

	// Run the plugin integration test
	config := map[string]interface{}{
		"interval": "1s", // Example configuration
	}

	suite.RunPluginIntegrationTest(t, demoPlugin.Name(), config)

	// You can also perform custom integration tests here
	// For example, testing database interactions:
	DB := suite.GetDB()
	if DB == nil {
		t.Error("Database should be available in integration test")
	}

	// Test repository operations
	repoManager := suite.GetRepositoryManager()
	if repoManager == nil {
		t.Error("Repository manager should be available in integration test")
	}

	// Test plugin manager
	pluginManager := suite.GetPluginManager()
	if pluginManager == nil {
		t.Error("Plugin manager should be available in integration test")
	}

	// Verify plugin is registered
	plugins := pluginManager.ListPlugins()
	found := false
	for _, plugin := range plugins {
		if plugin.Name == demoPlugin.Name() {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Demo plugin should be registered with plugin manager")
	}
}

// ExampleMockPluginTest demonstrates how to use mock plugins for testing
func ExampleMockPluginTest(t *testing.T) {
	// Create test reporter
	reporter := NewTestReporter("MockPluginTest")
	helper := NewPluginTestHelper(reporter)

	// Create a mock plugin
	mockPlugin := NewIntegrationTestSuite().CreateMockPlugin("mock-test", "1.0.0")
	ctx := NewTestPluginContext()

	// Initialize the plugin
	if err := mockPlugin.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize mock plugin: %v", err)
	}

	// Start the plugin
	if err := mockPlugin.Start(); err != nil {
		t.Fatalf("Failed to start mock plugin: %v", err)
	}

	// Create a plugin instance for reporting
	mockPluginWithCtx := mockPlugin.WithContext(ctx)
	instance := &types.PluginInstance{
		Plugin:            mockPluginWithCtx,
		Context:           ctx,
		Status:            types.StatusRunning,
		TotalExecutions:   10,
		TotalErrors:       1,
		HealthScore:       85.0,
		TotalExecutionTime: 100 * 1000 * 1000, // 100ms in nanoseconds
	}

	// Report the plugin test
	helper.ReportPluginTest(mockPluginWithCtx, ctx, instance)

	// Use helper methods to assert conditions
	helper.AssertLogContains(t, ctx, "INFO", "Mock plugin initialized")
	helper.AssertLogContains(t, ctx, "INFO", "Mock plugin started")

	// Generate report
	textReport := reporter.GenerateTextReport()
	t.Logf("Mock Plugin Test Report:\n%s", textReport)
}

// ExampleWithReporting demonstrates how to use the test framework with comprehensive reporting
func ExampleWithReporting(t *testing.T) {
	// Create test reporter
	reporter := NewTestReporter("ComprehensivePluginTest")
	helper := NewPluginTestHelper(reporter)

	// Create plugin and context
	plugin := demo.NewDemoPlugin()
	ctx := NewTestPluginContext()

	// Initialize plugin and record instance
	if err := plugin.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Start plugin
	if err := plugin.Start(); err != nil {
		t.Fatalf("Failed to start plugin: %v", err)
	}

	// Execute plugin functionality
	// demoPlugin := plugin.(*demo.DemoPlugin)
	// demoPlugin.FetchAndLogResource()

	// Stop plugin
	if err := plugin.Stop(); err != nil {
		t.Fatalf("Failed to stop plugin: %v", err)
	}

	// Create instance for reporting
	instance := &types.PluginInstance{
		Plugin:            plugin,
		Context:           ctx,
		Status:            types.StatusStopped,
		TotalExecutions:   1,
		TotalErrors:       0,
		HealthScore:       100.0,
		TotalExecutionTime: 50 * 1000 * 1000, // 50ms in nanoseconds
	}

	// Report plugin test
	helper.ReportPluginTest(plugin, ctx, instance)

	// Use helper assertions
	helper.AssertNoErrorLogs(t, ctx)

	// Generate comprehensive reports
	textReport := reporter.GenerateTextReport()
	jsonReport, err := reporter.GenerateJSONReport()
	if err != nil {
		t.Errorf("Failed to generate JSON report: %v", err)
	}

	// Log reports
	t.Logf("Comprehensive Test Report:\n%s", textReport)
	t.Logf("JSON Report: %s", jsonReport)
}

// ExampleErrorHandlingTest demonstrates how to test error handling scenarios
func ExampleErrorHandlingTest(t *testing.T) {
	// Create test reporter
	reporter := NewTestReporter("ErrorHandlingTest")

	// Test plugins with intentional errors
	t.Run("InitializeErrorTest", func(t *testing.T) {
		plugin := NewIntegrationTestSuite().
			CreateMockPlugin("init-error", "1.0.0").
			WithErrorOnInitialize()

		ctx := NewTestPluginContext()
		if err := plugin.Initialize(ctx); err == nil {
			t.Error("Expected initialization error")
		} else {
			t.Logf("Expected error occurred: %v", err)
		}
	})

	t.Run("StartErrorTest", func(t *testing.T) {
		plugin := NewIntegrationTestSuite().
			CreateMockPlugin("start-error", "1.0.0").
			WithErrorOnInitialize()

		ctx := NewTestPluginContext()
		if err := plugin.Initialize(ctx); err == nil {
			if err := plugin.Start(); err == nil {
				t.Error("Expected start error")
			} else {
				t.Logf("Expected error occurred: %v", err)
			}
		}
	})

	// Generate report
	textReport := reporter.GenerateTextReport()
	t.Logf("Error Handling Test Report:\n%s", textReport)
}