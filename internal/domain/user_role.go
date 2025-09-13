package domain

import (
	"errors"
	"fmt"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin   UserRole = "ADMIN"
	UserRoleManager UserRole = "MANAGER"
	UserRoleViewer  UserRole = "VIEWER"
)

// AllUserRoles returns all valid user roles
func AllUserRoles() []UserRole {
	return []UserRole{
		UserRoleAdmin,
		UserRoleManager,
		UserRoleViewer,
	}
}

// IsValid checks if the user role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin, UserRoleManager, UserRoleViewer:
		return true
	default:
		return false
	}
}

// String returns the string representation of the user role
func (r UserRole) String() string {
	return string(r)
}

// HasPermission checks if the role has a specific permission
func (r UserRole) HasPermission(permission string) bool {
	permissions := r.GetPermissions()
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetPermissions returns all permissions for this role
func (r UserRole) GetPermissions() []string {
	switch r {
	case UserRoleAdmin:
		return []string{
			"employee:read",
			"employee:write",
			"employee:delete",
			"user:read",
			"user:write",
			"user:delete",
			"audit:read",
			"system:admin",
		}
	case UserRoleManager:
		return []string{
			"employee:read",
			"employee:write",
			"user:read",
			"audit:read",
		}
	case UserRoleViewer:
		return []string{
			"employee:read",
			"audit:read",
		}
	default:
		return []string{}
	}
}

// CanAccessSalary checks if the role can access salary information
func (r UserRole) CanAccessSalary() bool {
	return r == UserRoleAdmin || r == UserRoleManager
}

// CanManageUsers checks if the role can manage users
func (r UserRole) CanManageUsers() bool {
	return r == UserRoleAdmin
}

// CanDeleteEmployees checks if the role can delete employees
func (r UserRole) CanDeleteEmployees() bool {
	return r == UserRoleAdmin
}

// CanViewAuditLogs checks if the role can view audit logs
func (r UserRole) CanViewAuditLogs() bool {
	return true // All roles can view audit logs
}

// ParseUserRole parses a string into a UserRole
func ParseUserRole(role string) (UserRole, error) {
	ur := UserRole(role)
	if !ur.IsValid() {
		return "", fmt.Errorf("invalid user role: %s", role)
	}
	return ur, nil
}

// MustParseUserRole parses a string into a UserRole or panics
func MustParseUserRole(role string) UserRole {
	ur, err := ParseUserRole(role)
	if err != nil {
		panic(err)
	}
	return ur
}

var (
	ErrInvalidUserRole = errors.New("invalid user role")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)