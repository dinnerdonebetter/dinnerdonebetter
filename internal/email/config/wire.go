package config

import (
	"net/http"

	"github.com/google/wire"

	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/observability/logging"
)

var (
	// Providers are what we provide to dependency injection.
	Providers = wire.NewSet(
		ProvideEmailer,
	)
)

// ProvideEmailer provides an email.Emailer from a config.
func ProvideEmailer(cfg *Config, logger logging.Logger, client *http.Client) (email.Emailer, error) {
	return cfg.ProvideEmailer(logger, client)
}
