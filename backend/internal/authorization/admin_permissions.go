package authorization

const (
	// ReadUserDataPermission is a service admin permission.
	ReadUserDataPermission Permission = "admin.read_user_data"
	// PublishArbitraryQueueMessagePermission is a service admin permission.
	PublishArbitraryQueueMessagePermission Permission = "queues.publish.message"
)

var (
	// AdminPermissions contains all admin-specific permissions.
	AdminPermissions = []Permission{
		ReadUserDataPermission,
		PublishArbitraryQueueMessagePermission,
	}

	// AdminServiceAdminPermissions contains admin permissions for the service admin role.
	// Pass to RegisterServiceAdminPermissions in the domain registration module.
	AdminServiceAdminPermissions = []Permission{
		ReadUserDataPermission,
		PublishArbitraryQueueMessagePermission,
	}
)
