package apiclients

import (
	authservice "github.com/prixfixeco/backend/internal/services/authentication"

	"github.com/google/wire"
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
		dataChangesTopicName:  cfg.DataChangesTopicName,
		minimumUsernameLength: cfg.MinimumUsernameLength,
		minimumPasswordLength: cfg.MinimumPasswordLength,
	}
}
