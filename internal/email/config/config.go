package config

import (
	"context"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/email/noop"
	"github.com/prixfixeco/api_server/internal/email/sendgrid"
	"github.com/prixfixeco/api_server/internal/observability/logging"
)

type (
	// Config is the configuration structure.
	Config struct {
		Provider string
		APIToken string
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.APIToken, validation.When(strings.EqualFold(strings.TrimSpace(cfg.Provider), "sendgrid"), validation.Required)),
	)
}

// ProvideEmailer provides an emailer.
func (cfg *Config) ProvideEmailer(logger logging.Logger, client *http.Client) (email.Emailer, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "sendgrid":
		return sendgrid.NewSendGridEmailer(cfg.APIToken, logger, client)
	default:
		return noop.NewNoopEmailer(), nil
	}
}
