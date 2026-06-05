package authorization

// RegisterCoreDomainPermissions registers the role-level permissions for all core service domains
// (waitlists, payments, webhooks, issuereports, uploadedmedia, notifications, comments) into the
// platform role sets. Call this once at application startup, before any permission checkers are
// constructed and before domain-specific registration (e.g. mealplanningregistration.RegisterForGRPCAPI).
//
// To add a new core domain: add per-role slices to its *_permissions.go file and append the
// corresponding Register* calls here.
func RegisterCoreDomainPermissions() {
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
