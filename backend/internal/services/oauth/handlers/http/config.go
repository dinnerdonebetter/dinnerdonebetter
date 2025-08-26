package http

import (
	"context"
	"time"

	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config is our configuration.
	Config struct {
		_ struct{} `json:"-"`

		Tokens               tokenscfg.Config `envPrefix:"TOKENS_"          json:"tokens"`
		Domain               string           `env:"DOMAIN"                 json:"domain"`
		AccessTokenLifespan  time.Duration    `env:"ACCESS_TOKEN_LIFESPAN"  json:"accessTokenLifespan"`
		RefreshTokenLifespan time.Duration    `env:"REFRESH_TOKEN_LIFESPAN" json:"refreshTokenLifespan"`
		Debug                bool             `env:"DEBUG"                  json:"debug"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.AccessTokenLifespan, validation.Required),
		validation.Field(&cfg.RefreshTokenLifespan, validation.Required),
		validation.Field(&cfg.Domain, validation.Required),
	)
}
