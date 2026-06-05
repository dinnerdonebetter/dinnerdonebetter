package authorization

const (
	// CreateServiceSettingsPermission is an admin user permission.
	CreateServiceSettingsPermission Permission = "create.service_settings"
	// ReadServiceSettingsPermission is an admin user permission.
	ReadServiceSettingsPermission Permission = "read.service_settings"
	// SearchServiceSettingsPermission is an admin user permission.
	SearchServiceSettingsPermission Permission = "search.service_settings"
	// ArchiveServiceSettingsPermission is an admin user permission.
	ArchiveServiceSettingsPermission Permission = "archive.service_settings"

	// CreateServiceSettingConfigurationsPermission is an admin user permission.
	CreateServiceSettingConfigurationsPermission Permission = "create.service_setting_configurations"
	// ReadServiceSettingConfigurationsPermission is an admin user permission.
	ReadServiceSettingConfigurationsPermission Permission = "read.service_setting_configurations"
	// UpdateServiceSettingConfigurationsPermission is an admin user permission.
	UpdateServiceSettingConfigurationsPermission Permission = "update.service_setting_configurations"
	// ArchiveServiceSettingConfigurationsPermission is an admin user permission.
	ArchiveServiceSettingConfigurationsPermission Permission = "archive.service_setting_configurations"
)

var (
	// SettingsPermissions contains all settings-related permissions.
	SettingsPermissions = []Permission{
		CreateServiceSettingsPermission,
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		ArchiveServiceSettingsPermission,
		CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission,
	}

	// SettingsServiceAdminPermissions contains settings permissions for the service admin role.
	// Pass to RegisterServiceAdminPermissions in the domain registration module.
	SettingsServiceAdminPermissions = []Permission{
		CreateServiceSettingsPermission,
		ArchiveServiceSettingsPermission,
	}

	// SettingsAccountMemberPermissions contains settings permissions for the account member role.
	// Pass to RegisterAccountMemberPermissions in the domain registration module.
	SettingsAccountMemberPermissions = []Permission{
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission,
	}
)
