package test

import (
	"testing"

	"github.com/ctwj/urldb/plugin/demo"
)

// Test framework tests
func TestPluginFramework(t *testing.T) {
	t.Run("TestPluginContext", func(t *testing.T) {
		ctx := NewTestPluginContext()

		// Test logging
		ctx.LogInfo("Test message with args: %d", 42)
		logs := ctx.GetLogs()
		if len(logs) != 1 {
			t.Errorf("Expected 1 log, got %d", len(logs))
		}

		// Test configuration
		ctx.SetConfig("test_key", "test_value")
		if val, err := ctx.GetConfig("test_key"); err != nil || val != "test_value" {
			t.Errorf("Config test failed: %v", err)
		}

		// Test data operations
		ctx.SetData("data_key", "data_value", "test_type")
		if val, err := ctx.GetData("data_key", "test_type"); err != nil || val != "data_value" {
			t.Errorf("Data test failed: %v", err)
		}

		// Test task registration
		taskCalled := false
		taskFunc := func() { taskCalled = true }
		ctx.RegisterTask("test_task", taskFunc)

		if task, err := ctx.GetTask("test_task"); err != nil {
			t.Errorf("Task registration failed: %v", err)
		} else {
			task()
			if !taskCalled {
				t.Error("Task was not called")
			}
		}
	})

	t.Run("TestPluginManager", func(t *testing.T) {
		manager := NewTestPluginManager()
		plugin := demo.NewDemoPlugin()

		// Register plugin
		if err := manager.RegisterPlugin(plugin); err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Get plugin
		if _, err := manager.GetPlugin(plugin.Name()); err != nil {
			t.Fatalf("Failed to get plugin: %v", err)
		}

		// Test lifecycle
		config := map[string]interface{}{
			"test": "value",
		}

		if err := manager.RunPluginLifecycle(t, plugin.Name(), config); err != nil {
			t.Errorf("Plugin lifecycle failed: %v", err)
		}
	})

	t.Run("TestLogAssertion", func(t *testing.T) {
		ctx := NewTestPluginContext()
		ctx.LogInfo("This is a test message")

		if !ctx.AssertLogContains(t, "INFO", "test message") {
			t.Error("Log assertion failed")
		}

		if ctx.AssertLogContains(t, "ERROR", "test message") {
			t.Error("Log assertion should have failed for wrong level")
		}
	})
}