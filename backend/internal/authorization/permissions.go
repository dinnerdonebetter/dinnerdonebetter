package authorization

type (
	role int

	// Permission is a simple string alias.
	Permission string
)

var (
	// ServiceAdminPermissions is every service admin permission.
	// All permissions are registered at startup via RegisterServiceAdminPermissions.
	ServiceAdminPermissions = []Permission{}

	// ServiceDataAdminPermissions is every service data admin permission.
	// All permissions are registered at startup via RegisterServiceDataAdminPermissions.
	ServiceDataAdminPermissions = []Permission{}

	// AccountAdminPermissions is every account admin permission.
	// All permissions are registered at startup via RegisterAccountAdminPermissions.
	AccountAdminPermissions = []Permission{}

	// AccountMemberPermissions is every account member permission.
	// All permissions are registered at startup via RegisterAccountMemberPermissions.
	AccountMemberPermissions = []Permission{}
)
