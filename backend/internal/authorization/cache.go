package authorization

import (
	"context"
	"sync"
)

// RolePermissionCacheLoader loads all custom role permission mappings from the database.
type RolePermissionCacheLoader func(ctx context.Context) (map[string][]string, error)

// RolePermissionCache provides fast in-memory lookups of custom role permissions.
type RolePermissionCache struct {
	rolePermissions map[string]map[string]struct{}
	mu              sync.RWMutex
}

// NewRolePermissionCache creates a new empty cache.
func NewRolePermissionCache() *RolePermissionCache {
	return &RolePermissionCache{
		rolePermissions: make(map[string]map[string]struct{}),
	}
}

// Refresh reloads all role→permission mappings from the database.
func (c *RolePermissionCache) Refresh(ctx context.Context, loader RolePermissionCacheLoader) error {
	data, err := loader(ctx)
	if err != nil {
		return err
	}

	built := make(map[string]map[string]struct{}, len(data))
	for roleID, perms := range data {
		set := make(map[string]struct{}, len(perms))
		for _, p := range perms {
			set[p] = struct{}{}
		}
		built[roleID] = set
	}

	c.mu.Lock()
	c.rolePermissions = built
	c.mu.Unlock()

	return nil
}

// GetPermissionsForRoles returns the union of all permissions for the given role IDs.
func (c *RolePermissionCache) GetPermissionsForRoles(roleIDs []string) map[string]struct{} {
	if len(roleIDs) == 0 {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]struct{})
	for _, roleID := range roleIDs {
		if perms, ok := c.rolePermissions[roleID]; ok {
			for p := range perms {
				result[p] = struct{}{}
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
