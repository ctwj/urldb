# urlDB Plugin Test Framework

This document describes the plugin test framework for urlDB.

## Overview

The plugin test framework provides a comprehensive set of tools for testing urlDB plugins. It includes:

1. **Unit Testing Framework** - For testing individual plugin components
2. **Integration Testing Environment** - For testing plugins in a complete system environment
3. **Test Reporting** - For generating detailed test reports
4. **Mock Objects** - For simulating system components during testing

## Components

### 1. Unit Testing Framework (`plugin/test/framework.go`)

The unit testing framework provides:

- `TestPluginContext` - A mock implementation of the PluginContext interface for testing plugin interactions with the system
- `TestPluginManager` - A test helper for managing plugin lifecycle in tests
- Logging and assertion utilities
- Configuration and data storage simulation
- Task scheduling simulation
- Cache system simulation
- Security permissions simulation
- Concurrency control simulation

### 2. Integration Testing Environment (`plugin/test/integration.go`)

The integration testing environment provides:

- `IntegrationTestSuite` - A complete integration test suite with database, repository manager, task manager, etc.
- `MockPlugin` - A mock plugin implementation for testing plugin manager functionality
- Various error scenario mock plugins
- Dependency relationship simulation
- Context operation simulation

### 3. Test Reporting (`plugin/test/reporting.go`)

The test reporting system provides:

- `TestReport` - Test report structure
- `TestReporter` - Test report generator
- `TestingTWrapper` - Wrapper for Go testing framework integration
- `PluginTestHelper` - Plugin test helper with specialized plugin testing functions

## Usage

### Writing Unit Tests

To write unit tests for plugins, follow this example:

```go
func TestMyPlugin(t *testing.T) {
    plugin := NewMyPlugin()

    // Create test context
    ctx := test.NewTestPluginContext()

    // Initialize plugin
    if err := plugin.Initialize(ctx); err != nil {
        t.Fatalf("Failed to initialize plugin: %v", err)
    }

    // Verify initialization logs
    if !ctx.AssertLogContains(t, "INFO", "Plugin initialized") {
        t.Error("Expected initialization log")
    }

    // Test other functionality...
}
```

### Writing Integration Tests

To write integration tests, follow this example:

```go
func TestMyPluginIntegration(t *testing.T) {
    // Create integration test suite
    suite := test.NewIntegrationTestSuite()
    suite.Setup(t)
    defer suite.Teardown()

    // Register plugin
    plugin := NewMyPlugin()
    if err := suite.RegisterPlugin(plugin); err != nil {
        t.Fatalf("Failed to register plugin: %v", err)
    }

    // Run integration test
    config := map[string]interface{}{
        "setting1": "value1",
    }

    suite.RunPluginIntegrationTest(t, plugin.Name(), config)
}
```

### Generating Test Reports

Test reports are automatically generated, but you can also create them manually:

```go
func TestWithReporting(t *testing.T) {
    // Create reporter
    reporter := test.NewTestReporter("MyTestSuite")
    wrapper := test.NewTestingTWrapper(t, reporter)

    // Use wrapper to run tests
    wrapper.Run("MyTest", func(t *testing.T) {
        // Test code...
    })

    // Generate report
    textReport := reporter.GenerateTextReport()
    t.Logf("Test Report:\n%s", textReport)
}
```

## Running Tests

### Run All Plugin Tests

```bash
go test ./plugin/...
```

### Run Specific Tests

```bash
go test ./plugin/demo/ -v
```

### Generate Test Coverage Report

```bash
go test ./plugin/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices

### 1. Test Plugin Lifecycle

Ensure you test the complete plugin lifecycle:

```go
func TestPluginLifecycle(t *testing.T) {
    manager := test.NewTestPluginManager()
    plugin := NewMyPlugin()

    // Register plugin
    manager.RegisterPlugin(plugin)

    // Test complete lifecycle
    config := map[string]interface{}{
        "config_key": "config_value",
    }

    if err := manager.RunPluginLifecycle(t, plugin.Name(), config); err != nil {
        t.Errorf("Plugin lifecycle failed: %v", err)
    }
}
```

### 2. Test Error Handling

Ensure you test plugin behavior under various error conditions:

```go
func TestPluginErrorHandling(t *testing.T) {
    // Test initialization error
    pluginWithInitError := test.NewIntegrationTestSuite().
        CreateMockPlugin("error-plugin", "1.0.0").
        WithErrorOnInitialize()

    ctx := test.NewTestPluginContext()
    if err := pluginWithInitError.Initialize(ctx); err == nil {
        t.Error("Expected initialize error")
    }
}
```

### 3. Test Dependencies

Test plugin dependency handling:

```go
func TestPluginDependencies(t *testing.T) {
    plugin := test.NewIntegrationTestSuite().
        CreateMockPlugin("dep-plugin", "1.0.0").
        WithDependencies([]string{"dep1", "dep2"})

    deps := plugin.Dependencies()
    if len(deps) != 2 {
        t.Errorf("Expected 2 dependencies, got %d", len(deps))
    }
}
```

### 4. Test Context Operations

Test plugin interactions with the system context:

```go
func TestPluginContextOperations(t *testing.T) {
    operations := []string{
        "log_info",
        "set_config",
        "get_data",
    }

    plugin := test.NewIntegrationTestSuite().
        CreateMockPlugin("context-plugin", "1.0.0").
        WithContextOperations(operations)

    ctx := test.NewTestPluginContext()
    plugin.Initialize(ctx)
    plugin.Start()

    // Verify operation results
    if !ctx.AssertLogContains(t, "INFO", "Info message") {
        t.Error("Expected info log")
    }
}
```

## Extending the Framework

### Adding New Test Features

To extend the test framework, you can:

1. Add new mock methods to `TestPluginContext`
2. Add new test helper methods to `TestPluginManager`
3. Add new reporting features to `TestReporter`

### Custom Report Formats

To create custom report formats, you can:

1. Extend the `TestReport` structure
2. Create new report generation methods
3. Implement specific report output formats

## Troubleshooting

### Common Issues

1. **Tests fail but no error message**
   - Check if test assertions are used correctly
   - Ensure test context is configured correctly

2. **Integration test environment setup fails**
   - Check database connection configuration
   - Ensure all dependent services are available

3. **Test reports are incomplete**
   - Ensure test reporter is used correctly
   - Check if tests complete normally