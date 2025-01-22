package api

import (
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

var (
	// ConfigProviders represents this package's offering to the dependency injector.
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.APIServiceConfig),
			"Queues",
			"Email",
			"Analytics",
			"Search",
			"FeatureFlags",
			"Encoding",
			"Events",
			"Observability",
			"Meta",
			"Routing",
			"HTTPServer",
			"Database",
			"Services",
		),
		wire.FieldsOf(
			new(*config.ServicesConfig),
			// Core
			"AuditLogEntries",
			"Auth",
			"Households",
			"HouseholdInvitations",
			"ServiceSettings",
			"ServiceSettingConfigurations",
			"Users",
			"UserNotifications",
			"Webhooks",
			"Workers",
			"DataPrivacy",
			"OAuth2Clients",
			// Data
			"ValidEnumerations",
			"MealPlanning",
			"Recipes",
		),
	)
)
