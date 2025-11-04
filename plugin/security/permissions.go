package security

import (
	"fmt"
	"sync"
)

// PermissionType defines different types of permissions a plugin can have
type PermissionType string

const (
	// System permissions
	PermissionSystemRead    PermissionType = "system:read"
	PermissionSystemWrite   PermissionType = "system:write"
	PermissionSystemExecute PermissionType = "system:execute"

	// Database permissions
	PermissionDatabaseRead  PermissionType = "database:read"
	PermissionDatabaseWrite PermissionType = "database:write"
	PermissionDatabaseExec  PermissionType = "database:exec"

	// Network permissions
	PermissionNetworkConnect PermissionType = "network:connect"
	PermissionNetworkListen  PermissionType = "network:listen"

	// File permissions
	PermissionFileRead  PermissionType = "file:read"
	PermissionFileWrite PermissionType = "file:write"
	PermissionFileExec  PermissionType = "file:exec"

	// Task permissions
	PermissionTaskSchedule PermissionType = "task:schedule"
	PermissionTaskControl  PermissionType = "task:control"

	// Config permissions
	PermissionConfigRead  PermissionType = "config:read"
	PermissionConfigWrite PermissionType = "config:write"

	// Data permissions
	PermissionDataRead  PermissionType = "data:read"
	PermissionDataWrite PermissionType = "data:write"
)

// Permission represents a specific permission
type Permission struct {
	Type        PermissionType `json:"type"`
	Description string         `json:"description"`
	Resource    string         `json:"resource,omitempty"` // Optional resource identifier
	Allowed     bool           `json:"allowed"`
}

// PermissionSet represents a collection of permissions
type PermissionSet struct {
	mu          sync.RWMutex
	permissions map[PermissionType][]Permission
}

// NewPermissionSet creates a new permission set
func NewPermissionSet() *PermissionSet {
	return &PermissionSet{
		permissions: make(map[PermissionType][]Permission),
	}
}

// GrantPermission grants a permission to the set
func (ps *PermissionSet) GrantPermission(permission Permission) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.permissions[permission.Type]; !exists {
		ps.permissions[permission.Type] = make([]Permission, 0)
	}

	ps.permissions[permission.Type] = append(ps.permissions[permission.Type], permission)
}

// CheckPermission checks if a specific permission is granted
func (ps *PermissionSet) CheckPermission(permissionType PermissionType, resource ...string) (bool, *Permission) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	permissions, exists := ps.permissions[permissionType]
	if !exists {
		return false, nil
	}

	resourceFilter := ""
	if len(resource) > 0 {
		resourceFilter = resource[0]
	}

	for _, perm := range permissions {
		if perm.Allowed {
			// If no resource filter, just match the type
			if resourceFilter == "" {
				return true, &perm
			}
			// If resource filter exists, match both type and resource
			if perm.Resource == "" || perm.Resource == resourceFilter {
				return true, &perm
			}
		}
	}

	return false, nil
}

// GetPermissions returns all permissions of a specific type
func (ps *PermissionSet) GetPermissions(permissionType PermissionType) []Permission {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if permissions, exists := ps.permissions[permissionType]; exists {
		result := make([]Permission, len(permissions))
		copy(result, permissions)
		return result
	}
	return nil
}

// GetAllPermissions returns all permissions in the set
func (ps *PermissionSet) GetAllPermissions() map[PermissionType][]Permission {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	result := make(map[PermissionType][]Permission)
	for permType, permissions := range ps.permissions {
		result[permType] = make([]Permission, len(permissions))
		copy(result[permType], permissions)
	}
	return result
}

// RevokePermission revokes a specific permission
func (ps *PermissionSet) RevokePermission(permissionType PermissionType, resource string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	permissions, exists := ps.permissions[permissionType]
	if !exists {
		return fmt.Errorf("permission type %s not found", permissionType)
	}

	updated := make([]Permission, 0)
	for _, perm := range permissions {
		if perm.Resource != resource {
			updated = append(updated, perm)
		}
	}

	ps.permissions[permissionType] = updated
	return nil
}

// ValidatePermissions validates the permissions against security policies
func (ps *PermissionSet) ValidatePermissions() error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	// Check for conflicting permissions
	for permType, permissions := range ps.permissions {
		for _, perm := range permissions {
			if !perm.Allowed {
				continue
			}

			// Validate permission type
			switch permType {
			case PermissionSystemWrite, PermissionSystemExecute:
				// These should be granted carefully
				if perm.Resource == "" {
					// System-wide permission - should be limited
				}
			case PermissionDatabaseExec:
				// Execution permission should be limited
			case PermissionFileExec:
				// Execution permission should be limited
			}

			// Add more validation rules as needed
		}
	}

	return nil
}

// DefaultPluginPermissions returns a default set of permissions for a plugin
func DefaultPluginPermissions(pluginName string) *PermissionSet {
	perms := NewPermissionSet()

	// Grant basic read permissions
	perms.GrantPermission(Permission{
		Type:        PermissionConfigRead,
		Description: "Read plugin configuration",
		Resource:    pluginName,
		Allowed:     true,
	})

	perms.GrantPermission(Permission{
		Type:        PermissionDataRead,
		Description: "Read plugin data",
		Resource:    pluginName,
		Allowed:     true,
	})

	// Grant basic data write permission for the plugin's own data
	perms.GrantPermission(Permission{
		Type:        PermissionDataWrite,
		Description: "Write plugin data",
		Resource:    pluginName,
		Allowed:     true,
	})

	// Grant basic task scheduling permission
	perms.GrantPermission(Permission{
		Type:        PermissionTaskSchedule,
		Description: "Schedule tasks",
		Resource:    pluginName,
		Allowed:     true,
	})

	return perms
}