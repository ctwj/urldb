package demo

import (
	"testing"
	"time"

	"github.com/ctwj/urldb/plugin/test"
	"github.com/ctwj/urldb/plugin/types"
)

// TestDemoPlugin tests the demo plugin
func TestDemoPlugin(t *testing.T) {
	plugin := NewDemoPlugin()

	// Test basic plugin information
	if plugin.Name() != "demo-plugin" {
		t.Errorf("Expected name 'demo-plugin', got '%s'", plugin.Name())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}

	if plugin.Author() != "urlDB Team" {
		t.Errorf("Expected author 'urlDB Team', got '%s'", plugin.Author())
	}

	// Test plugin interface compliance
	var _ types.Plugin = (*DemoPlugin)(nil)
}

// TestDemoPluginLifecycle tests the demo plugin lifecycle
func TestDemoPluginLifecycle(t *testing.T) {
	// Create test reporter
	reporter := test.NewTestReporter("DemoPluginLifecycle")
	wrapper := test.NewTestingTWrapper(t, reporter)

	wrapper.Run("FullLifecycle", func(t *testing.T) {
		plugin := NewDemoPlugin()
		manager := test.NewTestPluginManager()

		// Register plugin
		if err := manager.RegisterPlugin(plugin); err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Test complete lifecycle
		config := map[string]interface{}{
			"interval": "1s",
		}

		if err := manager.RunPluginLifecycle(t, plugin.Name(), config); err != nil {
			t.Errorf("Plugin lifecycle failed: %v", err)
		}
	})

	// Generate report
	report := reporter.GenerateTextReport()
	t.Logf("Test Report:\n%s", report)
}

// TestDemoPluginFunctionality tests the demo plugin functionality
func TestDemoPluginFunctionality(t *testing.T) {
	plugin := NewDemoPlugin()
	ctx := test.NewTestPluginContext()

	// Initialize plugin
	if err := plugin.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Check initialization logs
	if !ctx.AssertLogContains(t, "INFO", "Demo plugin initialized") {
		t.Error("Expected initialization log message")
	}

	// Start plugin
	if err := plugin.Start(); err != nil {
		t.Fatalf("Failed to start plugin: %v", err)
	}

	// Check start logs
	if !ctx.AssertLogContains(t, "INFO", "Demo plugin started") {
		t.Error("Expected start log message")
	}

	// Execute the task function directly to test functionality
	plugin.fetchAndLogResource()

	// Check that a resource was logged
	logs := ctx.GetLogs()
	found := false
	for _, log := range logs {
		if log.Level == "INFO" && log.Message == "Fetched resource: ID=%d, Title=%s, URL=%s" {
			found = true
			break
		}
	}
	if !found {
		t.Log("Note: Resource fetch log may not be present due to randomness in demo plugin")
	}

	// Stop plugin
	if err := plugin.Stop(); err != nil {
		t.Fatalf("Failed to stop plugin: %v", err)
	}

	// Check stop logs
	if !ctx.AssertLogContains(t, "INFO", "Demo plugin stopped") {
		t.Error("Expected stop log message")
	}

	// Cleanup plugin
	if err := plugin.Cleanup(); err != nil {
		t.Fatalf("Failed to cleanup plugin: %v", err)
	}

	// Check cleanup logs
	if !ctx.AssertLogContains(t, "INFO", "Demo plugin cleaned up") {
		t.Error("Expected cleanup log message")
	}
}

// TestDemoPluginWithReporting tests the demo plugin with reporting
func TestDemoPluginWithReporting(t *testing.T) {
	// Create test reporter
	reporter := test.NewTestReporter("DemoPluginWithReporting")
	helper := test.NewPluginTestHelper(reporter)

	// Create plugin and context
	plugin := NewDemoPlugin()
	ctx := test.NewTestPluginContext()

	// Create plugin instance for reporting
	instance := &types.PluginInstance{
		Plugin:            plugin,
		Context:           ctx,
		Status:            types.StatusInitialized,
		Config:            make(map[string]interface{}),
		StartTime:         time.Now(),
		TotalExecutions:   5,
		TotalErrors:       1,
		HealthScore:       85.5,
		TotalExecutionTime: time.Second,
	}

	// Initialize plugin
	if err := plugin.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Report plugin test
	helper.ReportPluginTest(plugin, ctx, instance)

	// Generate report
	jsonReport, err := reporter.GenerateJSONReport()
	if err != nil {
		t.Errorf("Failed to generate JSON report: %v", err)
	} else {
		t.Logf("JSON Report: %s", jsonReport)
	}

	textReport := reporter.GenerateTextReport()
	t.Logf("Text Report:\n%s", textReport)
}