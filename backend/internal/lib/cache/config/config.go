package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/cache"
	"github.com/dinnerdonebetter/backend/internal/lib/cache/memory"
	"github.com/dinnerdonebetter/backend/internal/lib/cache/redis"

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
		Redis    *redis.Config `env:"init"     envPrefix:"REDIS_" json:"redis"`
		Provider string        `env:"PROVIDER" json:"provider"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ProviderMemory, ProviderRedis)),
		validation.Field(&cfg.Redis, validation.When(cfg.Provider == ProviderRedis, validation.Required)),
	)
}

// ProvideCache provides a Cache.
func ProvideCache[T any](cfg *Config) (cache.Cache[T], error) {
	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case ProviderMemory:
		return memory.NewInMemoryCache[T](), nil
	case ProviderRedis:
		return redis.NewRedisCache[T](cfg.Redis, time.Hour), nil
	default:
		return nil, fmt.Errorf("invalid cache provider: %q", cfg.Provider)
	}
}
