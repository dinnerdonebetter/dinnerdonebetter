package apiclients

import (
	"github.com/google/wire"

	authservice "github.com/prixfixeco/backend/internal/services/authentication"
)

var (
	// Providers are what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvideConfig,
		ProvideAPIClientsService,
	)
)

// ProvideConfig converts an auth Config to a local Config.
func ProvideConfig(cfg *authservice.Config) *Config {
	return &Config{
		minimumUsernameLength: cfg.MinimumUsernameLength,
		minimumPasswordLength: cfg.MinimumPasswordLength,
	}
}
