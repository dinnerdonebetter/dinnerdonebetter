package config

import (
	"context"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/email/sendgrid"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	// ProviderSendgrid represents SendGrid.
	ProviderSendgrid = "sendgrid"
)

type (
	// Config is the configuration structure.
	Config struct {
		Provider string `json:"provider" mapstructure:"provider" toml:"provider,omitempty"`
		APIToken string `json:"apiToken" mapstructure:"api_token" toml:"api_token,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.APIToken, validation.When(strings.EqualFold(strings.TrimSpace(cfg.Provider), ProviderSendgrid), validation.Required)),
	)
}

// ProvideEmailer provides an emailer.
func (cfg *Config) ProvideEmailer(logger logging.Logger, client *http.Client) (email.Emailer, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderSendgrid:
		return sendgrid.NewSendGridEmailer(cfg.APIToken, logger, tracing.NewNoopTracerProvider(), client)
	default:
		return email.NewNoopEmailer()
	}
}
