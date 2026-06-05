package authorization

const (
	// ReportAnalyticsEventsPermission allows reporting analytics events via the proxy (TrackEvent).
	ReportAnalyticsEventsPermission Permission = "report.analytics_events"
)

var (
	// AnalyticsPermissions contains all analytics-related permissions.
	AnalyticsPermissions = []Permission{
		ReportAnalyticsEventsPermission,
	}

	// AnalyticsAccountMemberPermissions contains analytics permissions for the account member role.
	// Pass to RegisterAccountMemberPermissions in the domain registration module.
	AnalyticsAccountMemberPermissions = []Permission{
		ReportAnalyticsEventsPermission,
	}
)
