package tokenscfg

import (
	"context"
	"encoding/base64"
	"fmt"
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
		Provider                string `env:"PROVIDER"    json:"provider"`
		Audience                string `env:"AUDIENCE"    json:"audience"`
		Base64EncodedSigningKey string `env:"SIGNING_KEY" json:"base64EncodedSigningKey"`
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
		validation.Field(&cfg.Base64EncodedSigningKey, validation.Required),
	)
}

// ProvideTokenIssuer provides a token issuer.
func (cfg *Config) ProvideTokenIssuer(logger logging.Logger, tracerProvider tracing.TracerProvider) (tokens.Issuer, error) {
	decryptedSigningKey, err := base64.URLEncoding.DecodeString(cfg.Base64EncodedSigningKey)
	if err != nil {
		return nil, fmt.Errorf("decoding json web token signing key: %w", err)
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderJWT:
		return jwt.NewJWTSigner(logger, tracerProvider, cfg.Audience, decryptedSigningKey)
	case ProviderPASETO:
		return paseto.NewPASETOSigner(logger, tracerProvider, cfg.Audience, decryptedSigningKey)
	default:
		return tokens.NewNoopTokenIssuer(), nil
	}
}
