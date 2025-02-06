package grpcapi

import (
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"

	"github.com/google/wire"
)

var (
	// ConfigProviders represents this package's offering to the dependency injector.
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.APIServiceConfig),
			"Auth",
			"Queues",
			"Email",
			"Analytics",
			"TextSearch",
			"FeatureFlags",
			"Encoding",
			"Events",
			"Observability",
			"Meta",
			"Routing",
			"GRPCServer",
			"HTTPServer",
			"Database",
			"Services",
		),
		wire.FieldsOf(
			new(*authentication.Config),
			"Tokens",
			"SSO",
			"OAuth2",
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
			"DataPrivacy",
			"OAuth2Clients",
			// Data
			"ValidEnumerations",
			"MealPlanning",
			"Recipes",
		),
	)
)
