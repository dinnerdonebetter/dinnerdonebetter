package authorization

// RegisterCoreDomainPermissions registers the role-level permissions for all core service domains
// into the platform role sets. Call this once at application startup, before any permission checkers
// are constructed and before domain-specific registration (e.g. mealplanningregistration.RegisterForGRPCAPI).
//
// To add a new core domain: add per-role slices to its *_permissions.go file and append the
// corresponding Register* calls here.
func RegisterCoreDomainPermissions() {
	// admin (ReadUserData, PublishArbitraryQueueMessage)
	RegisterServiceAdminPermissions(AdminServiceAdminPermissions...)

	// auth (UpdateUserStatus, ReadUser, SearchUser, ArchiveUser, ImpersonateUser, ManageUserSessions)
	RegisterServiceAdminPermissions(AuthServiceAdminPermissions...)

	// identity (account membership management)
	RegisterAccountAdminPermissions(IdentityAccountAdminPermissions...)

	// oauth
	RegisterServiceAdminPermissions(OAuthServiceAdminPermissions...)
	RegisterAccountMemberPermissions(OAuthAccountMemberPermissions...)

	// settings
	RegisterServiceAdminPermissions(SettingsServiceAdminPermissions...)
	RegisterAccountMemberPermissions(SettingsAccountMemberPermissions...)

	// analytics
	RegisterAccountMemberPermissions(AnalyticsAccountMemberPermissions...)

	// audit
	RegisterAccountMemberPermissions(AuditAccountMemberPermissions...)

	// waitlists
	RegisterServiceAdminPermissions(WaitlistsServiceAdminPermissions...)
	RegisterAccountMemberPermissions(WaitlistsAccountMemberPermissions...)

	// payments
	RegisterServiceAdminPermissions(PaymentsServiceAdminPermissions...)
	RegisterAccountAdminPermissions(PaymentsAccountAdminPermissions...)
	RegisterAccountMemberPermissions(PaymentsAccountMemberPermissions...)

	// webhooks
	RegisterAccountAdminPermissions(WebhooksAccountAdminPermissions...)
	RegisterAccountMemberPermissions(WebhooksAccountMemberPermissions...)

	// issuereports
	RegisterAccountAdminPermissions(IssueReportsAccountAdminPermissions...)
	RegisterAccountMemberPermissions(IssueReportsAccountMemberPermissions...)

	// uploadedmedia
	RegisterAccountMemberPermissions(UploadedMediaAccountMemberPermissions...)

	// notifications
	RegisterServiceAdminPermissions(NotificationsServiceAdminPermissions...)
	RegisterAccountMemberPermissions(NotificationsAccountMemberPermissions...)

	// comments
	RegisterAccountMemberPermissions(CommentsAccountMemberPermissions...)
}
