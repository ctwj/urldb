package test

import (
	"testing"

	"github.com/ctwj/urldb/plugin/demo"
)

// TestSimpleFramework tests the basic functionality of the test framework
func TestSimpleFramework(t *testing.T) {
	// Test TestPluginContext
	ctx := NewTestPluginContext()

	// Test logging
	ctx.LogInfo("Test message")
	logs := ctx.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	// Test configuration
	ctx.SetConfig("test_key", "test_value")
	value, err := ctx.GetConfig("test_key")
	if err != nil {
		t.Errorf("Failed to get config: %v", err)
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value)
	}

	// Test data operations
	ctx.SetData("test_data_key", "test_data_value", "test_type")
	dataValue, err := ctx.GetData("test_data_key", "test_type")
	if err != nil {
		t.Errorf("Failed to get data: %v", err)
	}
	if dataValue != "test_data_value" {
		t.Errorf("Expected 'test_data_value', got '%v'", dataValue)
	}

	// Test plugin manager
	manager := NewTestPluginManager()
	plugin := demo.NewDemoPlugin()

	// Register plugin
	err = manager.RegisterPlugin(plugin)
	if err != nil {
		t.Errorf("Failed to register plugin: %v", err)
	}

	// Get plugin
	retrievedPlugin, err := manager.GetPlugin(plugin.Name())
	if err != nil {
		t.Errorf("Failed to get plugin: %v", err)
	}
	if retrievedPlugin.Name() != plugin.Name() {
		t.Errorf("Expected plugin name '%s', got '%s'", plugin.Name(), retrievedPlugin.Name())
	}
}