package emailcfg

import (
	"context"
	"net/http"
	"strings"

	circuitbreaking2 "github.com/dinnerdonebetter/backend/internal/circuitbreaking"
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
		Sendgrid             *sendgrid.Config         `env:"init"     envPrefix:"SENDGRID_"         json:"sendgrid"`
		Mailgun              *mailgun.Config          `env:"init"     envPrefix:"MAILGUN_"          json:"mailgun"`
		Mailjet              *mailjet.Config          `env:"init"     envPrefix:"MAILJET_"          json:"mailjet"`
		CircuitBreakerConfig *circuitbreaking2.Config `env:"init"     envPrefix:"CIRCUIT_BREAKING_" json:"circuitBreakerConfig"`
		Provider             string                   `env:"PROVIDER" json:"provider"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Sendgrid, validation.When(cfg.Provider == ProviderSendgrid, validation.Required)),
		validation.Field(&cfg.Mailgun, validation.When(cfg.Provider == ProviderMailgun, validation.Required)),
		validation.Field(&cfg.Mailjet, validation.When(cfg.Provider == ProviderMailjet, validation.Required)),
	)
}

// ProvideEmailer provides an outbound_emailer.
func (cfg *Config) ProvideEmailer(logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client, circuitBreaker circuitbreaking2.CircuitBreaker) (email.Emailer, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderSendgrid:
		return sendgrid.NewSendGridEmailer(cfg.Sendgrid, logger, tracerProvider, client, circuitBreaker)
	case ProviderMailgun:
		return mailgun.NewMailgunEmailer(cfg.Mailgun, logger, tracerProvider, client, circuitBreaker)
	case ProviderMailjet:
		return mailjet.NewMailjetEmailer(cfg.Mailjet, logger, tracerProvider, client, circuitBreaker)
	default:
		logger.Debug("providing noop outbound_emailer")
		return email.NewNoopEmailer()
	}
}
