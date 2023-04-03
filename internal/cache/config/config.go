package config

import (
	"context"
	"fmt"

	"github.com/prixfixeco/backend/internal/cache"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderMemory is the memory provider.
	ProviderMemory = "memory"
	// ProviderRedis is the redis provider.
	ProviderRedis = "redis"
)

type (
	// Config is the configuration for the cache.
	Config struct {
		Provider string   `json:"provider" mapstructure:"provider" toml:"provider"`
		Memory   struct{} `json:"memory" mapstructure:"memory" toml:"memory"`
		Redis    struct{} `json:"redis" mapstructure:"redis" toml:"redis"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ProviderMemory, ProviderRedis)),
		validation.Field(&cfg.Memory, validation.When(cfg.Provider == ProviderMemory, validation.Required)),
		validation.Field(&cfg.Redis, validation.When(cfg.Provider == ProviderRedis, validation.Required)),
	)
}

// ProvideCache validates a Config struct.
func ProvideCache[T cache.Cacheable](ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config) (cache.Cache[T], error) {
	switch cfg.Provider {
	case ProviderMemory:
		return nil, nil
	case ProviderRedis:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid cache provider: %q", cfg.Provider)
	}
}
