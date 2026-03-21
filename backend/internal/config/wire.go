package config

import (
	"github.com/google/wire"
	"github.com/verygoodsoftwarenotvirus/platform/server/http"
)

var (
	// ServiceConfigProviders represents this package's offering to the dependency injector.
	ServiceConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*APIServiceConfig),
			"Auth",
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
			new(*ServicesConfig),
			"Users",
			"DataPrivacy",
			"MealPlanning",
			"OAuth2Clients",
			"UploadedMedia",
			"Payments",
		),
		ProvideHTTPServerConfigFromAPIServiceConfig,
	)
)

func ProvideHTTPServerConfigFromAPIServiceConfig(cfg *APIServiceConfig) http.Config {
	return cfg.HTTPServer
}
