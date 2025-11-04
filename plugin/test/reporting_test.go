package test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ctwj/urldb/plugin/demo"
	"github.com/ctwj/urldb/plugin/types"
)

// TestTestReporter tests the test reporter functionality
func TestTestReporter(t *testing.T) {
	reporter := NewTestReporter("TestSuite")

	// Test starting and ending a test
	reporter.StartTest("Test1")
	time.Sleep(10 * time.Millisecond) // Simulate test duration
	reporter.EndTest("Test1", true, nil)

	// Test failing test
	reporter.StartTest("Test2")
	reporter.EndTest("Test2", false, errors.New("test error"))

	// Test skipped test
	reporter.SkipTest("Test3")

	// Generate report
	report := reporter.GenerateReport()

	// Verify report contents
	if report.SuiteName != "TestSuite" {
		t.Errorf("Expected suite name 'TestSuite', got '%s'", report.SuiteName)
	}

	if report.TotalTests != 3 {
		t.Errorf("Expected 3 total tests, got %d", report.TotalTests)
	}

	if report.PassedTests != 1 {
		t.Errorf("Expected 1 passed test, got %d", report.PassedTests)
	}

	if report.FailedTests != 1 {
		t.Errorf("Expected 1 failed test, got %d", report.FailedTests)
	}

	if report.SkippedTests != 1 {
		t.Errorf("Expected 1 skipped test, got %d", report.SkippedTests)
	}

	// Verify test cases
	if len(report.TestCases) != 3 {
		t.Errorf("Expected 3 test cases, got %d", len(report.TestCases))
	}

	// Verify first test case
	tc1 := report.TestCases[0]
	if tc1.Name != "Test1" {
		t.Errorf("Expected test name 'Test1', got '%s'", tc1.Name)
	}
	if tc1.Status != "passed" {
		t.Errorf("Expected test status 'passed', got '%s'", tc1.Status)
	}
	if tc1.Duration <= 0 {
		t.Error("Expected positive duration")
	}

	// Verify second test case
	tc2 := report.TestCases[1]
	if tc2.Name != "Test2" {
		t.Errorf("Expected test name 'Test2', got '%s'", tc2.Name)
	}
	if tc2.Status != "failed" {
		t.Errorf("Expected test status 'failed', got '%s'", tc2.Status)
	}
	if tc2.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", tc2.Error)
	}

	// Verify third test case
	tc3 := report.TestCases[2]
	if tc3.Name != "Test3" {
		t.Errorf("Expected test name 'Test3', got '%s'", tc3.Name)
	}
	if tc3.Status != "skipped" {
		t.Errorf("Expected test status 'skipped', got '%s'", tc3.Status)
	}
}

// TestTextReportGeneration tests text report generation
func TestTextReportGeneration(t *testing.T) {
	reporter := NewTestReporter("TextReportTest")

	// Add some test data
	reporter.StartTest("SampleTest")
	reporter.EndTest("SampleTest", true, nil)

	// Generate text report
	textReport := reporter.GenerateTextReport()

	// Verify report contains expected elements
	if !contains(textReport, "Test Report: TextReportTest") {
		t.Error("Report should contain suite name")
	}

	if !contains(textReport, "SampleTest") {
		t.Error("Report should contain test name")
	}

	if !contains(textReport, "Passed: 1") {
		t.Error("Report should show 1 passed test")
	}

	t.Logf("Generated text report:\n%s", textReport)
}

// TestJSONReportGeneration tests JSON report generation
func TestJSONReportGeneration(t *testing.T) {
	reporter := NewTestReporter("JSONReportTest")

	// Add some test data
	reporter.StartTest("JSONTest")
	reporter.EndTest("JSONTest", false, errors.New("json test error"))

	// Generate JSON report
	jsonReport, err := reporter.GenerateJSONReport()
	if err != nil {
		t.Fatalf("Failed to generate JSON report: %v", err)
	}

	// Verify report contains expected elements
	if !contains(jsonReport, `"suite_name": "JSONReportTest"`) {
		t.Error("JSON report should contain suite name")
	}

	if !contains(jsonReport, `"name": "JSONTest"`) {
		t.Error("JSON report should contain test name")
	}

	if !contains(jsonReport, `"status": "failed"`) {
		t.Error("JSON report should show failed status")
	}

	if !contains(jsonReport, `"error": "json test error"`) {
		t.Error("JSON report should contain error message")
	}

	t.Logf("Generated JSON report:\n%s", jsonReport)
}

// TestPluginReportGeneration tests plugin report generation
func TestPluginReportGeneration(t *testing.T) {
	reporter := NewTestReporter("PluginReportTest")

	// Create a demo plugin
	plugin := demo.NewDemoPlugin()

	// Create test logs
	logs := []LogEntry{
		{Level: "INFO", Message: "Plugin initialized", Time: time.Now()},
		{Level: "DEBUG", Message: "Processing resource", Args: []interface{}{123}, Time: time.Now()},
		{Level: "WARN", Message: "Resource not found", Time: time.Now()},
	}

	// Create performance report
	perf := PerformanceReport{
		AvgExecutionTime: 15.5,
		TotalExecutions:  100,
		TotalErrors:      5,
		HealthScore:      92.5,
	}

	// Add plugin report
	reporter.AddPluginReport(plugin, logs, perf)

	// Generate text report
	textReport := reporter.GenerateTextReport()

	// Verify plugin report is included
	if !contains(textReport, "Plugin: demo-plugin") {
		t.Error("Report should contain plugin name")
	}

	if !contains(textReport, "Avg Execution Time: 15.50ms") {
		t.Error("Report should contain average execution time")
	}

	if !contains(textReport, "Health Score: 92.50") {
		t.Error("Report should contain health score")
	}

	if !contains(textReport, "Recent Logs:") {
		t.Error("Report should contain recent logs section")
	}

	t.Logf("Generated plugin text report:\n%s", textReport)
}

// TestTestingTWrapper tests the TestingTWrapper
func TestTestingTWrapper(t *testing.T) {
	reporter := NewTestReporter("WrapperTest")
	wrapper := NewTestingTWrapper(t, reporter)

	// Test running a sub-test
	result := wrapper.Run("SubTest1", func(t *testing.T) {
		// This test should pass
	})

	if !result {
		t.Error("SubTest1 should have passed")
	}

	// Test running a failing sub-test
	result = wrapper.Run("SubTest2", func(t *testing.T) {
		t.Error("This test is designed to fail")
	})

	if result {
		t.Error("SubTest2 should have failed")
	}

	// Generate report
	report := reporter.GenerateReport()
	if report.TotalTests != 2 {
		t.Errorf("Expected 2 total tests, got %d", report.TotalTests)
	}

	if report.PassedTests != 1 {
		t.Errorf("Expected 1 passed test, got %d", report.PassedTests)
	}

	if report.FailedTests != 1 {
		t.Errorf("Expected 1 failed test, got %d", report.FailedTests)
	}
}

// TestPluginTestHelper tests the PluginTestHelper
func TestPluginTestHelper(t *testing.T) {
	reporter := NewTestReporter("HelperTest")
	helper := NewPluginTestHelper(reporter)

	// Create test context and plugin
	ctx := NewTestPluginContext()
	plugin := demo.NewDemoPlugin()

	// Add some logs
	ctx.LogInfo("Plugin started successfully")
	ctx.LogError("Failed to process resource: %s", "resource_id")

	// Create plugin instance
	instance := &types.PluginInstance{
		Plugin:            plugin,
		Context:           ctx,
		Status:            types.StatusRunning,
		TotalExecutions:   50,
		TotalErrors:       2,
		HealthScore:       95.0,
		TotalExecutionTime: 500 * time.Millisecond,
	}

	// Report plugin test
	helper.ReportPluginTest(plugin, ctx, instance)

	// Test assertions
	helper.AssertLogContains(t, ctx, "INFO", "Plugin started")
	helper.AssertLogContains(t, ctx, "ERROR", "Failed to process")

	// Generate report
	textReport := reporter.GenerateTextReport()
	if !contains(textReport, "Plugin: demo-plugin") {
		t.Error("Report should contain plugin information")
	}

	if !contains(textReport, "Avg Execution Time: 10.00ms") {
		t.Error("Report should contain calculated average execution time")
	}

	t.Logf("Generated helper test report:\n%s", textReport)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && strings.Index(s, substr) != -1
}