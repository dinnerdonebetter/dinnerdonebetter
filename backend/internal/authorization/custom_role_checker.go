package authorization

// CustomRolePermissionChecker holds a precomputed set of permission IDs
// derived from one or more custom roles assigned to a user.
type CustomRolePermissionChecker struct {
	permissions map[string]struct{}
}

// NewCustomRolePermissionChecker creates a checker from a set of permission IDs.
func NewCustomRolePermissionChecker(permissionIDs map[string]struct{}) *CustomRolePermissionChecker {
	return &CustomRolePermissionChecker{
		permissions: permissionIDs,
	}
}

// HasPermission returns whether the custom role set includes the given permission.
func (c *CustomRolePermissionChecker) HasPermission(p Permission) bool {
	if c == nil || c.permissions == nil {
		return false
	}
	_, ok := c.permissions[p.ID()]
	return ok
}
