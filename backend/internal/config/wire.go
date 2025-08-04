package config

import (
	"github.com/dinnerdonebetter/backend/internal/platform/server/http"

	"github.com/google/wire"
)

var (
	// ServiceConfigProviders represents this package's offering to the dependency injector.
	ServiceConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*APIServiceConfig),
			"Observability",
			"Email",
			"Analytics",
			"FeatureFlags",
			"Encoding",
			"Routing",
			"Database",
			"Meta",
			"Events",
			"Queues",
			"TextSearch",
			"Services",
		),
		wire.FieldsOf(
			new(*AdminWebappConfig),
			"Cookies",
		),
		wire.FieldsOf(
			new(*ServicesConfig),
			// Core
			"AuditLogEntries",
			"Auth",
			"Accounts",
			"AccountInvitations",
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
			"MealPlanning",
		),
		ProvideHTTPServerConfigFromAPIServiceConfig,
		ProvideHTTPServerConfigFromAdminWebappConfig,
	)
)

func ProvideHTTPServerConfigFromAPIServiceConfig(cfg *APIServiceConfig) http.Config {
	return cfg.HTTPServer
}

func ProvideHTTPServerConfigFromAdminWebappConfig(cfg *AdminWebappConfig) http.Config {
	return cfg.HTTPServer
}
