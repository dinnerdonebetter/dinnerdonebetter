package tokenscfg

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/authentication/tokens/jwt"
	"github.com/dinnerdonebetter/backend/internal/authentication/tokens/paseto"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderJWT represents JWT.
	ProviderJWT = "jwt"
	// ProviderPASETO represents PASETO.
	ProviderPASETO = "paseto"
)

type (
	// Config is the configuration structure.
	Config struct {
		Provider   string `env:"PROVIDER"    json:"provider"`
		Audience   string `env:"AUDIENCE"    json:"audience"`
		SigningKey []byte `env:"SIGNING_KEY" json:"signing_key"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Provider, validation.In(ProviderJWT, ProviderPASETO)),
		validation.Field(&cfg.Audience, validation.Required),
		validation.Field(&cfg.SigningKey, validation.Required),
	)
}

// ProvideTokenIssuer provides a token issuer.
func (cfg *Config) ProvideTokenIssuer(logger logging.Logger, tracerProvider tracing.TracerProvider) (tokens.Issuer, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderJWT:
		return jwt.NewJWTSigner(logger, tracerProvider, cfg.Audience, cfg.SigningKey)
	case ProviderPASETO:
		return paseto.NewPASETOSigner(logger, tracerProvider, cfg.Audience, cfg.SigningKey)
	default:
		return tokens.NewNoopTokenIssuer(), nil
	}
}
