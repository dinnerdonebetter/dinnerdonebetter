package authorization

const (
	// CreateCustomRolesPermission is a permission for creating custom roles.
	CreateCustomRolesPermission Permission = "create.custom_roles"
	// ReadCustomRolesPermission is a permission for reading custom roles.
	ReadCustomRolesPermission Permission = "read.custom_roles"
	// UpdateCustomRolesPermission is a permission for updating custom roles.
	UpdateCustomRolesPermission Permission = "update.custom_roles"
	// ArchiveCustomRolesPermission is a permission for archiving custom roles.
	ArchiveCustomRolesPermission Permission = "archive.custom_roles"
	// AssignCustomRolesPermission is a permission for assigning custom roles to users.
	AssignCustomRolesPermission Permission = "assign.custom_roles"
)
