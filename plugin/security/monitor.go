package security

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// BehaviorMonitor monitors plugin behavior
type BehaviorMonitor struct {
	mu               sync.RWMutex
	pluginActivities map[string][]PluginActivity
	maxActivities    int
	thresholds       BehaviorThresholds
	alerts           []SecurityAlert
}

// PluginActivity represents a plugin activity
type PluginActivity struct {
	Timestamp   time.Time     `json:"timestamp"`
	Activity    string        `json:"activity"`
	Resource    string        `json:"resource,omitempty"`
	Details     interface{}   `json:"details,omitempty"`
	ExecutionTime time.Duration `json:"execution_time,omitempty"`
}

// BehaviorThresholds defines thresholds for behavior monitoring
type BehaviorThresholds struct {
	MaxDatabaseQueries     int           `json:"max_database_queries"`
	MaxNetworkConnections  int           `json:"max_network_connections"`
	MaxFileOperations      int           `json:"max_file_operations"`
	MaxExecutionTime       time.Duration `json:"max_execution_time"`
	MaxMemoryUsage         int64         `json:"max_memory_usage"`
	AlertThreshold         int           `json:"alert_threshold"`
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	Timestamp time.Time `json:"timestamp"`
	Plugin    string    `json:"plugin"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"` // low, medium, high, critical
	Details   string    `json:"details,omitempty"`
}

// NewBehaviorMonitor creates a new behavior monitor
func NewBehaviorMonitor() *BehaviorMonitor {
	return &BehaviorMonitor{
		pluginActivities: make(map[string][]PluginActivity),
		maxActivities:    1000, // Keep last 1000 activities per plugin
		thresholds: BehaviorThresholds{
			MaxDatabaseQueries:    1000,
			MaxNetworkConnections: 100,
			MaxFileOperations:     500,
			MaxExecutionTime:      30 * time.Second,
			MaxMemoryUsage:        100 * 1024 * 1024, // 100MB
			AlertThreshold:        5,
		},
		alerts: make([]SecurityAlert, 0),
	}
}

// LogActivity logs a plugin activity
func (bm *BehaviorMonitor) LogActivity(pluginName, activity, resource string, details interface{}) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	activityRecord := PluginActivity{
		Timestamp: time.Now(),
		Activity:  activity,
		Resource:  resource,
		Details:   details,
	}

	if _, exists := bm.pluginActivities[pluginName]; !exists {
		bm.pluginActivities[pluginName] = make([]PluginActivity, 0, bm.maxActivities)
	}

	activities := bm.pluginActivities[pluginName]
	if len(activities) >= bm.maxActivities {
		// Remove oldest activity
		activities = activities[1:]
	}

	activities = append(activities, activityRecord)
	bm.pluginActivities[pluginName] = activities

	// Check for suspicious behavior
	bm.checkBehavior(pluginName, activityRecord)
}

// LogExecutionTime logs execution time for an activity
func (bm *BehaviorMonitor) LogExecutionTime(pluginName, activity, resource string, duration time.Duration) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	activityRecord := PluginActivity{
		Timestamp:     time.Now(),
		Activity:      activity,
		Resource:      resource,
		ExecutionTime: duration,
	}

	if _, exists := bm.pluginActivities[pluginName]; !exists {
		bm.pluginActivities[pluginName] = make([]PluginActivity, 0, bm.maxActivities)
	}

	activities := bm.pluginActivities[pluginName]
	if len(activities) >= bm.maxActivities {
		activities = activities[1:]
	}

	activities = append(activities, activityRecord)
	bm.pluginActivities[pluginName] = activities

	// Check for suspicious behavior
	bm.checkBehavior(pluginName, activityRecord)
}

// GetActivities returns activities for a plugin
func (bm *BehaviorMonitor) GetActivities(pluginName string, limit int) []PluginActivity {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	activities, exists := bm.pluginActivities[pluginName]
	if !exists {
		return nil
	}

	if limit <= 0 || limit >= len(activities) {
		result := make([]PluginActivity, len(activities))
		copy(result, activities)
		return result
	}

	// Return last N activities
	start := len(activities) - limit
	result := make([]PluginActivity, limit)
	copy(result, activities[start:])
	return result
}

// GetAllActivities returns all activities
func (bm *BehaviorMonitor) GetAllActivities() map[string][]PluginActivity {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	result := make(map[string][]PluginActivity)
	for plugin, activities := range bm.pluginActivities {
		result[plugin] = make([]PluginActivity, len(activities))
		copy(result[plugin], activities)
	}
	return result
}

// GetAlerts returns security alerts
func (bm *BehaviorMonitor) GetAlerts() []SecurityAlert {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	result := make([]SecurityAlert, len(bm.alerts))
	copy(result, bm.alerts)
	return result
}

// GetRecentAlerts returns recent security alerts
func (bm *BehaviorMonitor) GetRecentAlerts(duration time.Duration) []SecurityAlert {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	recent := make([]SecurityAlert, 0)
	for _, alert := range bm.alerts {
		if alert.Timestamp.After(cutoff) {
			recent = append(recent, alert)
		}
	}
	return recent
}

// checkBehavior checks for suspicious behavior and generates alerts
func (bm *BehaviorMonitor) checkBehavior(pluginName string, activity PluginActivity) {
	// Check execution time
	if activity.ExecutionTime > bm.thresholds.MaxExecutionTime {
		bm.generateAlert(pluginName, "slow_execution", fmt.Sprintf("Plugin took %v to execute", activity.ExecutionTime), "medium", activity)
	}

	// Check for excessive database queries
	if activity.Activity == "database_query" {
		// Count recent database queries for this plugin
		recentQueries := 0
		for _, act := range bm.pluginActivities[pluginName] {
			if act.Activity == "database_query" &&
				act.Timestamp.After(time.Now().Add(-1*time.Minute)) {
				recentQueries++
			}
		}
		if recentQueries > bm.thresholds.MaxDatabaseQueries {
			bm.generateAlert(pluginName, "excessive_db_queries",
				fmt.Sprintf("Plugin executed %d database queries in last minute", recentQueries),
				"high", map[string]interface{}{
					"count": recentQueries,
					"limit": bm.thresholds.MaxDatabaseQueries,
				})
		}
	}

	// Check for excessive file operations
	if activity.Activity == "file_operation" {
		// Count recent file operations for this plugin
		recentFiles := 0
		for _, act := range bm.pluginActivities[pluginName] {
			if act.Activity == "file_operation" &&
				act.Timestamp.After(time.Now().Add(-1*time.Minute)) {
				recentFiles++
			}
		}
		if recentFiles > bm.thresholds.MaxFileOperations {
			bm.generateAlert(pluginName, "excessive_file_ops",
				fmt.Sprintf("Plugin performed %d file operations in last minute", recentFiles),
				"high", map[string]interface{}{
					"count": recentFiles,
					"limit": bm.thresholds.MaxFileOperations,
				})
		}
	}

	// Check for suspicious network activity
	if activity.Activity == "network_connect" {
		details, ok := activity.Details.(map[string]interface{})
		if ok {
			host, _ := details["host"].(string)
			// Check if connecting to suspicious hosts
			if isSuspiciousHost(host) {
				bm.generateAlert(pluginName, "suspicious_network",
					fmt.Sprintf("Plugin connecting to suspicious host: %s", host),
					"high", details)
			}
		}
	}
}

// generateAlert generates a security alert
func (bm *BehaviorMonitor) generateAlert(pluginName, alertType, message, severity string, details interface{}) {
	alert := SecurityAlert{
		Timestamp: time.Now(),
		Plugin:    pluginName,
		Type:      alertType,
		Message:   message,
		Severity:  severity,
	}

	if details != nil {
		detailBytes, err := json.Marshal(details)
		if err == nil {
			alert.Details = string(detailBytes)
		}
	}

	bm.alerts = append(bm.alerts, alert)

	// Log the alert
	utils.Warn("Security Alert [%s] for plugin %s: %s", severity, pluginName, message)

	// Keep only recent alerts (last 1000)
	if len(bm.alerts) > 1000 {
		bm.alerts = bm.alerts[len(bm.alerts)-1000:]
	}
}

// isSuspiciousHost checks if a host is suspicious
func isSuspiciousHost(host string) bool {
	suspiciousHosts := []string{
		"127.0.0.1",
		"localhost",
		"0.0.0.0",
	}

	for _, suspicious := range suspiciousHosts {
		if host == suspicious {
			return true
		}
	}

	return false
}

// SetThresholds sets behavior monitoring thresholds
func (bm *BehaviorMonitor) SetThresholds(thresholds BehaviorThresholds) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.thresholds = thresholds
}

// GetThresholds returns current behavior monitoring thresholds
func (bm *BehaviorMonitor) GetThresholds() BehaviorThresholds {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.thresholds
}

// ClearActivities clears activities for a plugin
func (bm *BehaviorMonitor) ClearActivities(pluginName string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.pluginActivities, pluginName)
}

// ClearAllActivities clears all activities
func (bm *BehaviorMonitor) ClearAllActivities() {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.pluginActivities = make(map[string][]PluginActivity)
}

// ClearAlerts clears all alerts
func (bm *BehaviorMonitor) ClearAlerts() {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.alerts = make([]SecurityAlert, 0)
}

// ExportActivities exports activities for a plugin
func (bm *BehaviorMonitor) ExportActivities(pluginName string) ([]byte, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	activities, exists := bm.pluginActivities[pluginName]
	if !exists {
		return nil, fmt.Errorf("no activities found for plugin %s", pluginName)
	}

	return json.Marshal(activities)
}