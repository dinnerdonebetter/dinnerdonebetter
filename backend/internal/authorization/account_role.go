package authorization

import (
	"encoding/gob"

	"github.com/mikespook/gorbac/v2"
)

type (
	// AccountRole describes a role a user has for an account context.
	AccountRole role

	// AccountRolePermissionsChecker checks permissions for one or more account Roles.
	AccountRolePermissionsChecker interface {
		HasPermission(Permission) bool
	}
)

const (
	// AccountMemberRole is a role for a plain account participant.
	AccountMemberRole AccountRole = iota
	// AccountAdminRole is a role for someone who can manipulate the specifics of an account.
	AccountAdminRole AccountRole = iota
	// AccountRestrictedRole is a role with zero built-in permissions, used when custom roles are the sole source.
	AccountRestrictedRole AccountRole = iota

	AccountAdminRoleName      = "account_admin"
	AccountMemberRoleName     = "account_member"
	AccountRestrictedRoleName = "account_restricted"
)

var (
	accountAdmin      = gorbac.NewStdRole(AccountAdminRoleName)
	accountMember     = gorbac.NewStdRole(AccountMemberRoleName)
	accountRestricted = gorbac.NewStdRole(AccountRestrictedRoleName)
)

type accountRoleCollection struct {
	customPerms *CustomRolePermissionChecker
	Roles       []string
}

func init() {
	gob.Register(accountRoleCollection{})
}

// NewAccountRolePermissionChecker returns a new checker for a set of Roles.
func NewAccountRolePermissionChecker(roles ...string) AccountRolePermissionsChecker {
	return &accountRoleCollection{
		Roles: roles,
	}
}

// NewAccountRolePermissionCheckerWithCustomRoles returns a new checker that combines built-in roles with custom role permissions.
func NewAccountRolePermissionCheckerWithCustomRoles(customPerms *CustomRolePermissionChecker, roles ...string) AccountRolePermissionsChecker {
	return &accountRoleCollection{
		Roles:       roles,
		customPerms: customPerms,
	}
}

func (r AccountRole) String() string {
	switch r {
	case AccountMemberRole:
		return AccountMemberRoleName
	case AccountAdminRole:
		return AccountAdminRoleName
	case AccountRestrictedRole:
		return AccountRestrictedRoleName
	default:
		return ""
	}
}

// HasPermission returns whether a user can do something or not.
func (r accountRoleCollection) HasPermission(p Permission) bool {
	if hasPermission(p, r.Roles...) {
		return true
	}
	if r.customPerms != nil {
		return r.customPerms.HasPermission(p)
	}
	return false
}

func hasPermission(p Permission, roles ...string) bool {
	if len(roles) == 0 {
		return false
	}

	for _, r := range roles {
		if globalAuthorizer.IsGranted(r, p, nil) {
			return true
		}
	}

	return false
}
