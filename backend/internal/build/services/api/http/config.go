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
			"HTTPServer",
			"Database",
			"Services",
			"AppleAppSiteAssociation",
		),
		wire.FieldsOf(
			new(*config.ServicesConfig),
			"Users",
			"DataPrivacy",
			"MealPlanning",
			"Auth",
			"OAuth2Clients",
		),
	)
)
