package authorization

type (
	role int

	// Permission is a simple string alias.
	Permission string
)

var (
	// ServiceAdminPermissions is every service admin permission.
	// Domain-specific permissions are added at startup via RegisterServiceAdminPermissions.
	ServiceAdminPermissions = []Permission{
		ReadUserDataPermission,
		UpdateUserStatusPermission,
		ReadUserPermission,
		SearchUserPermission,
		ArchiveUserPermission,
		CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission,
		ArchiveServiceSettingsPermission,
		ImpersonateUserPermission,
		ManageUserSessionsPermission,
		PublishArbitraryQueueMessagePermission,
		// only admins can arbitrarily create these via the API, this is exclusively for integration test purposes.
		CreateServiceSettingsPermission,
	}

	// ServiceDataAdminPermissions is every service data admin permission.
	// Domain-specific permissions are added at startup via RegisterServiceDataAdminPermissions.
	ServiceDataAdminPermissions = []Permission{}

	// AccountAdminPermissions is every account admin permission.
	// Domain-specific permissions are added at startup via RegisterAccountAdminPermissions.
	AccountAdminPermissions = []Permission{
		UpdateAccountPermission,
		ArchiveAccountPermission,
		TransferAccountPermission,
		InviteUserToAccountPermission,
		ModifyMemberPermissionsForAccountPermission,
		RemoveMemberAccountPermission,
	}

	// AccountMemberPermissions is every account member permission.
	// Domain-specific permissions are added at startup via RegisterAccountMemberPermissions.
	AccountMemberPermissions = []Permission{
		ReportAnalyticsEventsPermission,
		ReadAuditLogEntriesPermission,
		ReadOAuth2ClientsPermission,
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission,
	}
)
