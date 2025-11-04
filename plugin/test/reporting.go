package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// TestReport represents a test report
type TestReport struct {
	SuiteName     string        `json:"suite_name"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	TotalTests    int           `json:"total_tests"`
	PassedTests   int           `json:"passed_tests"`
	FailedTests   int           `json:"failed_tests"`
	SkippedTests  int           `json:"skipped_tests"`
	TestCases     []TestCase    `json:"test_cases"`
	PluginReports []PluginReport `json:"plugin_reports"`
}

// TestCase represents a single test case
type TestCase struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // passed, failed, skipped
	Duration  float64   `json:"duration"` // in seconds
	Error     string    `json:"error,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// PluginReport represents a plugin-specific report
type PluginReport struct {
	PluginName    string            `json:"plugin_name"`
	PluginVersion string            `json:"plugin_version"`
	TestCases     []TestCase        `json:"test_cases"`
	Logs          []LogEntry        `json:"logs"`
	Performance   PerformanceReport `json:"performance"`
}

// PerformanceReport represents performance metrics
type PerformanceReport struct {
	AvgExecutionTime float64 `json:"avg_execution_time"` // in milliseconds
	TotalExecutions  int64   `json:"total_executions"`
	TotalErrors      int64   `json:"total_errors"`
	HealthScore      float64 `json:"health_score"`
}

// TestReporter is responsible for generating test reports
type TestReporter struct {
	report TestReport
}

// NewTestReporter creates a new test reporter
func NewTestReporter(suiteName string) *TestReporter {
	return &TestReporter{
		report: TestReport{
			SuiteName: suiteName,
			StartTime: time.Now(),
			TestCases: make([]TestCase, 0),
		},
	}
}

// StartTest marks the start of a test case
func (tr *TestReporter) StartTest(name string) {
	testCase := TestCase{
		Name:      name,
		StartTime: time.Now(),
	}
	tr.report.TestCases = append(tr.report.TestCases, testCase)
}

// EndTest marks the end of a test case
func (tr *TestReporter) EndTest(name string, passed bool, err error) {
	now := time.Now()

	// Find the test case
	for i := len(tr.report.TestCases) - 1; i >= 0; i-- {
		if tr.report.TestCases[i].Name == name && tr.report.TestCases[i].Status == "" {
			testCase := &tr.report.TestCases[i]
			testCase.EndTime = now
			testCase.Duration = now.Sub(testCase.StartTime).Seconds()

			if passed {
				testCase.Status = "passed"
				tr.report.PassedTests++
			} else {
				testCase.Status = "failed"
				tr.report.FailedTests++
				if err != nil {
					testCase.Error = err.Error()
				}
			}
			break
		}
	}

	tr.report.TotalTests++
}

// SkipTest marks a test as skipped
func (tr *TestReporter) SkipTest(name string) {
	now := time.Now()
	testCase := TestCase{
		Name:      name,
		Status:    "skipped",
		StartTime: now,
		EndTime:   now,
		Duration:  0,
	}
	tr.report.TestCases = append(tr.report.TestCases, testCase)
	tr.report.SkippedTests++
	tr.report.TotalTests++
}

// AddPluginReport adds a plugin-specific report
func (tr *TestReporter) AddPluginReport(plugin types.Plugin, logs []LogEntry, perf PerformanceReport) {
	pluginReport := PluginReport{
		PluginName:    plugin.Name(),
		PluginVersion: plugin.Version(),
		Logs:          logs,
		Performance:   perf,
	}
	tr.report.PluginReports = append(tr.report.PluginReports, pluginReport)
}

// GenerateReport generates the final test report
func (tr *TestReporter) GenerateReport() TestReport {
	tr.report.EndTime = time.Now()
	return tr.report
}

// GenerateTextReport generates a human-readable text report
func (tr *TestReporter) GenerateTextReport() string {
	var buf bytes.Buffer

	report := tr.GenerateReport()

	buf.WriteString(fmt.Sprintf("Test Report: %s\n", report.SuiteName))
	buf.WriteString(fmt.Sprintf("Start Time: %s\n", report.StartTime.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("End Time: %s\n", report.EndTime.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("Duration: %.2f seconds\n", report.EndTime.Sub(report.StartTime).Seconds()))
	buf.WriteString(fmt.Sprintf("Total Tests: %d\n", report.TotalTests))
	buf.WriteString(fmt.Sprintf("Passed: %d\n", report.PassedTests))
	buf.WriteString(fmt.Sprintf("Failed: %d\n", report.FailedTests))
	buf.WriteString(fmt.Sprintf("Skipped: %d\n", report.SkippedTests))
	buf.WriteString("\n")

	// Test cases
	buf.WriteString("Test Cases:\n")
	for _, tc := range report.TestCases {
		status := "✓"
		if tc.Status == "failed" {
			status = "✗"
		} else if tc.Status == "skipped" {
			status = "⊘"
		}

		buf.WriteString(fmt.Sprintf("  %s %s (%.3fs)\n", status, tc.Name, tc.Duration))
		if tc.Status == "failed" && tc.Error != "" {
			buf.WriteString(fmt.Sprintf("    Error: %s\n", tc.Error))
		}
	}

	// Plugin reports
	if len(report.PluginReports) > 0 {
		buf.WriteString("\nPlugin Reports:\n")
		for _, pr := range report.PluginReports {
			buf.WriteString(fmt.Sprintf("\n  Plugin: %s (v%s)\n", pr.PluginName, pr.PluginVersion))

			// Performance metrics
			buf.WriteString(fmt.Sprintf("    Avg Execution Time: %.2fms\n", pr.Performance.AvgExecutionTime))
			buf.WriteString(fmt.Sprintf("    Total Executions: %d\n", pr.Performance.TotalExecutions))
			buf.WriteString(fmt.Sprintf("    Total Errors: %d\n", pr.Performance.TotalErrors))
			buf.WriteString(fmt.Sprintf("    Health Score: %.2f\n", pr.Performance.HealthScore))

			// Logs summary
			if len(pr.Logs) > 0 {
				buf.WriteString("    Recent Logs:\n")
				// Show last 5 logs
				start := len(pr.Logs) - 5
				if start < 0 {
					start = 0
				}
				for i := start; i < len(pr.Logs); i++ {
					log := pr.Logs[i]
					buf.WriteString(fmt.Sprintf("      [%s] %s: %s\n", log.Level, log.Time.Format("15:04:05"), formatLogMessage(log.Message, log.Args)))
				}
			}
		}
	}

	return buf.String()
}

// GenerateJSONReport generates a JSON report
func (tr *TestReporter) GenerateJSONReport() (string, error) {
	report := tr.GenerateReport()

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// formatLogMessage formats a log message with its arguments
func formatLogMessage(message string, args []interface{}) string {
	if len(args) == 0 {
		return message
	}
	return fmt.Sprintf(message, args...)
}

// TestingTWrapper wraps *testing.T to integrate with our reporter
type TestingTWrapper struct {
	*testing.T
	reporter *TestReporter
}

// NewTestingTWrapper creates a new TestingTWrapper
func NewTestingTWrapper(t *testing.T, reporter *TestReporter) *TestingTWrapper {
	return &TestingTWrapper{
		T:        t,
		reporter: reporter,
	}
}

// Run wraps the testing.T.Run method to integrate with reporting
func (t *TestingTWrapper) Run(name string, f func(t *testing.T)) bool {
	t.reporter.StartTest(name)

	// Create a sub-test
	result := t.T.Run(name, func(subT *testing.T) {
		// Handle panic to ensure we report the test result
		defer func() {
			if r := recover(); r != nil {
				t.reporter.EndTest(name, false, fmt.Errorf("panic: %v", r))
				panic(r) // Re-panic
			}
		}()

		f(subT)

		// If we reach here, the test passed
		t.reporter.EndTest(name, !subT.Failed(), nil)
	})

	// If Run returned false, it means the test was skipped or failed to start
	if !result {
		t.reporter.SkipTest(name)
	}

	return result
}

// Error wraps the testing.T.Error method to integrate with reporting
func (t *TestingTWrapper) Error(args ...interface{}) {
	t.T.Error(args...)
}

// Errorf wraps the testing.T.Errorf method to integrate with reporting
func (t *TestingTWrapper) Errorf(format string, args ...interface{}) {
	t.T.Errorf(format, args...)
}

// FailNow wraps the testing.T.FailNow method to integrate with reporting
func (t *TestingTWrapper) FailNow() {
	t.T.FailNow()
}

// PluginTestHelper provides helper methods for testing plugins
type PluginTestHelper struct {
	reporter *TestReporter
}

// NewPluginTestHelper creates a new plugin test helper
func NewPluginTestHelper(reporter *TestReporter) *PluginTestHelper {
	return &PluginTestHelper{
		reporter: reporter,
	}
}

// ReportPluginTest generates a report for a plugin test
func (pth *PluginTestHelper) ReportPluginTest(plugin types.Plugin, context *TestPluginContext, instance *types.PluginInstance) {
	// Collect logs
	logs := context.GetLogs()

	// Collect performance metrics
	var perf PerformanceReport
	if instance != nil {
		perf = PerformanceReport{
			AvgExecutionTime: 0,
			TotalExecutions:  instance.TotalExecutions,
			TotalErrors:      instance.TotalErrors,
			HealthScore:      instance.HealthScore,
		}

		if instance.TotalExecutions > 0 {
			avgTime := float64(instance.TotalExecutionTime.Milliseconds()) / float64(instance.TotalExecutions)
			perf.AvgExecutionTime = avgTime
		}
	}

	// Add to reporter
	pth.reporter.AddPluginReport(plugin, logs, perf)
}

// AssertLogContains asserts that a log with the specified level and content exists
func (pth *PluginTestHelper) AssertLogContains(t *testing.T, context *TestPluginContext, level string, contains string) {
	if !context.AssertLogContains(t, level, contains) {
		t.Errorf("Expected to find log message containing '%s' with level '%s'", contains, level)
	}
}

// AssertNoErrorLogs asserts that there are no error logs
func (pth *PluginTestHelper) AssertNoErrorLogs(t *testing.T, context *TestPluginContext) {
	logs := context.GetLogs()
	for _, log := range logs {
		if log.Level == "ERROR" {
			t.Errorf("Unexpected error log: %s", formatLogMessage(log.Message, log.Args))
		}
	}
}

// AssertPluginStatus asserts that a plugin has the expected status
func (pth *PluginTestHelper) AssertPluginStatus(t *testing.T, instance *types.PluginInstance, expected types.PluginStatus) {
	if instance.Status != expected {
		t.Errorf("Expected plugin status to be %s, got %s", expected, instance.Status)
	}
}