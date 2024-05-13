package config

import (
	"context"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/email/mailgun"
	"github.com/dinnerdonebetter/backend/internal/email/mailjet"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderSendgrid represents SendGrid.
	ProviderSendgrid = "sendgrid"
	// ProviderMailgun represents Mailgun.
	ProviderMailgun = "mailgun"
	// ProviderMailjet represents Mailjet.
	ProviderMailjet = "mailjet"
)

type (
	// Config is the configuration structure.
	Config struct {
		Sendgrid *sendgrid.Config `json:"sendgrid" toml:"sendgrid,omitempty"`
		Mailgun  *mailgun.Config  `json:"mailgun"  toml:"mailgun,omitempty"`
		Mailjet  *mailjet.Config  `json:"mailjet"  toml:"mailjet,omitempty"`
		Provider string           `json:"provider" toml:"provider,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	cfg.EnsureDefaults()

	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Sendgrid, validation.When(cfg.Provider == ProviderSendgrid, validation.Required)),
		validation.Field(&cfg.Mailgun, validation.When(cfg.Provider == ProviderMailgun, validation.Required)),
		validation.Field(&cfg.Mailjet, validation.When(cfg.Provider == ProviderMailjet, validation.Required)),
	)
}

// EnsureDefaults ensures certain default values are set.
func (cfg *Config) EnsureDefaults() {
	switch cfg.Provider {
	case ProviderSendgrid:
		if cfg.Sendgrid != nil {
			cfg.Sendgrid.EnsureDefaults()
		}
	}
}

// ProvideEmailer provides an outbound_emailer.
func (cfg *Config) ProvideEmailer(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client) (email.Emailer, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderSendgrid:
		return sendgrid.NewSendGridEmailer(ctx, cfg.Sendgrid, logger, tracerProvider, client)
	case ProviderMailgun:
		return mailgun.NewMailgunEmailer(cfg.Mailgun, logger, tracerProvider, client)
	case ProviderMailjet:
		return mailjet.NewMailjetEmailer(cfg.Mailjet, logger, tracerProvider, client)
	default:
		logger.Debug("providing noop outbound_emailer")
		return email.NewNoopEmailer()
	}
}
