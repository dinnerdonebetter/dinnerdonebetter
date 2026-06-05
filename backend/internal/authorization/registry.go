package authorization

import "sync"

var registryMu sync.Mutex

// RegisterServiceAdminPermissions appends domain-specific permissions to ServiceAdminPermissions.
// Call this once during application startup, before any permission checkers are constructed.
func RegisterServiceAdminPermissions(perms ...Permission) {
	registryMu.Lock()
	ServiceAdminPermissions = append(ServiceAdminPermissions, perms...)
	registryMu.Unlock()
}

// RegisterServiceDataAdminPermissions appends domain-specific permissions to ServiceDataAdminPermissions.
// Call this once during application startup, before any permission checkers are constructed.
func RegisterServiceDataAdminPermissions(perms ...Permission) {
	registryMu.Lock()
	ServiceDataAdminPermissions = append(ServiceDataAdminPermissions, perms...)
	registryMu.Unlock()
}

// RegisterAccountAdminPermissions appends domain-specific permissions to AccountAdminPermissions.
// Call this once during application startup, before any permission checkers are constructed.
func RegisterAccountAdminPermissions(perms ...Permission) {
	registryMu.Lock()
	AccountAdminPermissions = append(AccountAdminPermissions, perms...)
	registryMu.Unlock()
}

// RegisterAccountMemberPermissions appends domain-specific permissions to AccountMemberPermissions.
// Call this once during application startup, before any permission checkers are constructed.
func RegisterAccountMemberPermissions(perms ...Permission) {
	registryMu.Lock()
	AccountMemberPermissions = append(AccountMemberPermissions, perms...)
	registryMu.Unlock()
}
