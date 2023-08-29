package stripe

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config configures our Stripe interface.
	Config struct {
		APIKey        string `json:"apiKey"        toml:"api_key"`
		SuccessURL    string `json:"successURL"    toml:"success_url"`
		CancelURL     string `json:"cancelURL"     toml:"cancel_url"`
		WebhookSecret string `json:"webhookSecret" toml:"webhook_secret"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.APIKey, validation.Required),
	)
}
