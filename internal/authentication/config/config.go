package config

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/argon2"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	ProviderArgon2 = "argon2"
)

type (
	Config struct {
		Argon2   argon2.Config `json:"argon2"   toml:"argon2"`
		Provider string        `json:"provider" toml:"provider"`
	}
)

func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ProviderArgon2)),
	)
}

// ProvideAuthenticator provides an authentication.Authenticator implementation based on the config.
func ProvideAuthenticator(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (authentication.Authenticator, error) {
	switch cfg.Provider {
	case ProviderArgon2:
		return argon2.ProvideAuthenticator(logger, tracerProvider), nil
	default:
		return nil, fmt.Errorf("unknown provider: %q", cfg.Provider)
	}
}
