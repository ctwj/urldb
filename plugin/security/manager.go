package security

import (
	"fmt"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// SecurityManager manages plugin security
type SecurityManager struct {
	db              *gorm.DB
	permissionSets  map[string]*PermissionSet
	behaviorMonitor *BehaviorMonitor
	mu              sync.RWMutex
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(database *gorm.DB) *SecurityManager {
	return &SecurityManager{
		db:              database,
		permissionSets:  make(map[string]*PermissionSet),
		behaviorMonitor: NewBehaviorMonitor(),
	}
}

// Initialize initializes the security manager
func (sm *SecurityManager) Initialize() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Load existing permissions from database
	if err := sm.loadPermissionsFromDB(); err != nil {
		return fmt.Errorf("failed to load permissions from database: %v", err)
	}

	utils.Info("Security manager initialized")
	return nil
}

// loadPermissionsFromDB loads permissions from the database
func (sm *SecurityManager) loadPermissionsFromDB() error {
	// This would load permissions from a database table
	// For now, we'll initialize with default permissions
	return nil
}

// GetPermissionSet returns the permission set for a plugin
func (sm *SecurityManager) GetPermissionSet(pluginName string) *PermissionSet {
	sm.mu.RLock()
	permSet, exists := sm.permissionSets[pluginName]
	sm.mu.RUnlock()

	if !exists {
		sm.mu.Lock()
		// Double-check after acquiring write lock
		permSet, exists = sm.permissionSets[pluginName]
		if !exists {
			// Create default permission set
			permSet = DefaultPluginPermissions(pluginName)
			sm.permissionSets[pluginName] = permSet
		}
		sm.mu.Unlock()
	}

	return permSet
}

// SetPermissionSet sets the permission set for a plugin
func (sm *SecurityManager) SetPermissionSet(pluginName string, permSet *PermissionSet) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.permissionSets[pluginName] = permSet
}

// CheckPermission checks if a plugin has a specific permission
func (sm *SecurityManager) CheckPermission(pluginName string, permissionType PermissionType, resource ...string) (bool, *Permission) {
	permSet := sm.GetPermissionSet(pluginName)
	if permSet == nil {
		return false, nil
	}

	return permSet.CheckPermission(permissionType, resource...)
}

// GrantPermission grants a permission to a plugin
func (sm *SecurityManager) GrantPermission(pluginName string, permission Permission) error {
	permSet := sm.GetPermissionSet(pluginName)
	if permSet == nil {
		permSet = NewPermissionSet()
		sm.SetPermissionSet(pluginName, permSet)
	}

	permSet.GrantPermission(permission)
	return nil
}

// RevokePermission revokes a permission from a plugin
func (sm *SecurityManager) RevokePermission(pluginName string, permissionType PermissionType, resource string) error {
	permSet := sm.GetPermissionSet(pluginName)
	if permSet == nil {
		return fmt.Errorf("no permission set found for plugin %s", pluginName)
	}

	return permSet.RevokePermission(permissionType, resource)
}

// ValidatePluginPermissions validates a plugin's permissions
func (sm *SecurityManager) ValidatePluginPermissions(pluginName string) error {
	permSet := sm.GetPermissionSet(pluginName)
	if permSet == nil {
		return fmt.Errorf("no permission set found for plugin %s", pluginName)
	}

	return permSet.ValidatePermissions()
}

// GetBehaviorMonitor returns the behavior monitor
func (sm *SecurityManager) GetBehaviorMonitor() *BehaviorMonitor {
	return sm.behaviorMonitor
}

// LogActivity logs a plugin activity
func (sm *SecurityManager) LogActivity(pluginName, activity, resource string, details interface{}) {
	sm.behaviorMonitor.LogActivity(pluginName, activity, resource, details)
}

// LogExecutionTime logs execution time for a plugin activity
func (sm *SecurityManager) LogExecutionTime(pluginName, activity, resource string, duration time.Duration) {
	sm.behaviorMonitor.LogExecutionTime(pluginName, activity, resource, duration)
}

// GetPluginActivities returns activities for a plugin
func (sm *SecurityManager) GetPluginActivities(pluginName string, limit int) []PluginActivity {
	return sm.behaviorMonitor.GetActivities(pluginName, limit)
}

// GetSecurityAlerts returns security alerts
func (sm *SecurityManager) GetSecurityAlerts() []SecurityAlert {
	return sm.behaviorMonitor.GetAlerts()
}

// GetRecentAlerts returns recent security alerts
func (sm *SecurityManager) GetRecentAlerts(duration time.Duration) []SecurityAlert {
	return sm.behaviorMonitor.GetRecentAlerts(duration)
}

// SavePermissionsToDB saves permissions to the database
func (sm *SecurityManager) SavePermissionsToDB(pluginName string) error {
	// This would save permissions to a database table
	// Implementation depends on the database schema
	return nil
}

// LoadPermissionsFromDB loads permissions from the database for a specific plugin
func (sm *SecurityManager) LoadPermissionsFromDB(pluginName string) error {
	// This would load permissions from a database table for a specific plugin
	// Implementation depends on the database schema
	return nil
}

// CreateSecurityReport generates a security report for a plugin
func (sm *SecurityManager) CreateSecurityReport(pluginName string) *SecurityReport {
	report := &SecurityReport{
		PluginName:     pluginName,
		GeneratedAt:    time.Now(),
		Permissions:    make([]Permission, 0),
		Activities:     sm.GetPluginActivities(pluginName, 50),
		Alerts:         sm.GetRecentAlerts(24 * time.Hour),
		SecurityScore:  100.0, // Start with perfect score
		Issues:         make([]SecurityIssue, 0),
		Recommendations: make([]SecurityRecommendation, 0),
	}

	// Get permissions
	permSet := sm.GetPermissionSet(pluginName)
	if permSet != nil {
		allPerms := permSet.GetAllPermissions()
		for _, perms := range allPerms {
			report.Permissions = append(report.Permissions, perms...)
		}
	}

	// Calculate security score based on alerts and permissions
	score := 100.0
	for _, alert := range report.Alerts {
		switch alert.Severity {
		case "low":
			score -= 1
		case "medium":
			score -= 5
		case "high":
			score -= 15
		case "critical":
			score -= 30
		}
	}

	// Check for overly permissive permissions
	if permSet != nil {
		perms := permSet.GetPermissions(PermissionSystemWrite)
		if len(perms) > 0 {
			report.Issues = append(report.Issues, SecurityIssue{
				Type:        "overly_permissive",
				Description: "Plugin has system write permissions",
				Severity:    "high",
			})
			score -= 20
		}

		perms = permSet.GetPermissions(PermissionFileExec)
		if len(perms) > 0 {
			report.Issues = append(report.Issues, SecurityIssue{
				Type:        "overly_permissive",
				Description: "Plugin has file execution permissions",
				Severity:    "high",
			})
			score -= 20
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}
	report.SecurityScore = score

	// Generate recommendations
	if score < 80 {
		report.Recommendations = append(report.Recommendations, SecurityRecommendation{
			Type:        "security_improvement",
			Description: "Review and reduce plugin permissions to minimum required",
			Priority:    "high",
		})
	}

	if len(report.Alerts) > 5 {
		report.Recommendations = append(report.Recommendations, SecurityRecommendation{
			Type:        "monitoring",
			Description: "Increase monitoring for this plugin due to frequent alerts",
			Priority:    "medium",
		})
	}

	return report
}

// SecurityReport represents a security report for a plugin
type SecurityReport struct {
	PluginName      string                  `json:"plugin_name"`
	GeneratedAt     time.Time               `json:"generated_at"`
	Permissions     []Permission            `json:"permissions"`
	Activities      []PluginActivity        `json:"activities"`
	Alerts          []SecurityAlert         `json:"alerts"`
	SecurityScore   float64                 `json:"security_score"`
	Issues          []SecurityIssue         `json:"issues"`
	Recommendations []SecurityRecommendation `json:"recommendations"`
}

// SecurityIssue represents a security issue
type SecurityIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// SecurityRecommendation represents a security recommendation
type SecurityRecommendation struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}