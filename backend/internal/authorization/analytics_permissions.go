package authorization

import (
	"github.com/mikespook/gorbac/v2"
)

const (
	// ReportAnalyticsEventsPermission allows reporting analytics events via the proxy (TrackEvent).
	ReportAnalyticsEventsPermission Permission = "report.analytics_events"
)

var (
	// AnalyticsPermissions contains all analytics-related permissions.
	AnalyticsPermissions = []gorbac.Permission{
		ReportAnalyticsEventsPermission,
	}
)
