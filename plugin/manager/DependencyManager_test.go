package manager

import (
	"errors"
	"testing"

	"github.com/ctwj/urldb/plugin/types"
)

// MockPlugin is a mock plugin implementation for testing
type MockPlugin struct {
	name         string
	version      string
	description  string
	author       string
	dependencies []string
}

func NewMockPlugin(name string, dependencies []string) *MockPlugin {
	return &MockPlugin{
		name:         name,
		version:      "1.0.0",
		description:  "Mock plugin for testing",
		author:       "Test",
		dependencies: dependencies,
	}
}

func (m *MockPlugin) Name() string { return m.name }
func (m *MockPlugin) Version() string { return m.version }
func (m *MockPlugin) Description() string { return m.description }
func (m *MockPlugin) Author() string { return m.author }
func (m *MockPlugin) Initialize(ctx types.PluginContext) error { return nil }
func (m *MockPlugin) Start() error { return nil }
func (m *MockPlugin) Stop() error { return nil }
func (m *MockPlugin) Cleanup() error { return nil }
func (m *MockPlugin) Dependencies() []string { return m.dependencies }
func (m *MockPlugin) CheckDependencies() map[string]bool {
	result := make(map[string]bool)
	for _, dep := range m.dependencies {
		result[dep] = true // For testing, assume all dependencies are satisfied
	}
	return result
}

func TestDependencyManager(t *testing.T) {
	manager := NewManager(nil, nil, nil)
	depManager := NewDependencyManager(manager)

	// Test registering plugins with dependencies
	pluginA := NewMockPlugin("pluginA", []string{})
	pluginB := NewMockPlugin("pluginB", []string{"pluginA"})
	pluginC := NewMockPlugin("pluginC", []string{"pluginB"})

	manager.RegisterPlugin(pluginA)
	manager.RegisterPlugin(pluginB)
	manager.RegisterPlugin(pluginC)

	// Test dependency validation
	err := depManager.ValidateDependencies()
	if err != nil {
		t.Errorf("Expected no validation errors, got: %v", err)
	}

	// Test load order
	loadOrder, err := depManager.GetLoadOrder()
	if err != nil {
		t.Errorf("Expected no error when getting load order, got: %v", err)
	}

	// Verify that pluginA comes before pluginB, which comes before pluginC
	pluginAIndex, pluginBIndex, pluginCIndex := -1, -1, -1
	for i, name := range loadOrder {
		if name == "pluginA" {
			pluginAIndex = i
		} else if name == "pluginB" {
			pluginBIndex = i
		} else if name == "pluginC" {
			pluginCIndex = i
		}
	}

	if pluginAIndex >= pluginBIndex {
		t.Errorf("Expected pluginA to come before pluginB in load order")
	}
	if pluginBIndex >= pluginCIndex {
		t.Errorf("Expected pluginB to come before pluginC in load order")
	}

	// Test circular dependency detection
	manager2 := NewManager(nil, nil, nil)
	depManager2 := NewDependencyManager(manager2)

	// Create circular dependency: A -> B -> C -> A
	circularA := NewMockPlugin("circularA", []string{"circularC"}) // A depends on C
	circularB := NewMockPlugin("circularB", []string{"circularA"}) // B depends on A
	circularC := NewMockPlugin("circularC", []string{"circularB"}) // C depends on B

	manager2.RegisterPlugin(circularA)
	manager2.RegisterPlugin(circularB)
	manager2.RegisterPlugin(circularC)

	err = depManager2.ValidateDependencies()
	if err == nil {
		t.Errorf("Expected circular dependency error, got none")
	}
	if err != nil && err.Error()[:9] != "circular " {
		t.Errorf("Expected circular dependency error, got: %v", err)
	}

	// Test dependency checking for a single plugin
	manager3 := NewManager(nil, nil, nil)
	depManager3 := NewDependencyManager(manager3)

	// Create an instance for pluginA so we can check dependency status
	pluginA3 := NewMockPlugin("pluginA3", []string{})
	manager3.RegisterPlugin(pluginA3)

	// Simulate that pluginA3 is running
	manager3.instances["pluginA3"] = &PluginInstance{
		Plugin: pluginA3,
		Status: types.StatusRunning,
	}

	pluginD := NewMockPlugin("pluginD", []string{"pluginA3"})
	manager3.RegisterPlugin(pluginD)

	// Check dependencies for pluginD
	satisfied, unresolved, err := depManager3.CheckPluginDependencies("pluginD")
	if err != nil {
		t.Errorf("Expected no error when checking pluginD dependencies, got: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected pluginD dependencies to be satisfied, unresolved: %v", unresolved)
	}
}