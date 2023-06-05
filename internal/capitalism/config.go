package capitalism

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// StripeProvider is the key that indicates Stripe should be used for payments.
	StripeProvider = "stripe"
)

type (
	// Config allows for the configuration of this package and its subpackages.
	Config struct {
		Stripe   *StripeConfig `json:"stripe"   toml:"stripe"`
		Provider string        `json:"provider" toml:"provider"`
		Enabled  bool          `json:"enabled"  toml:"enabled"`
	}

	// StripeConfig configures our Stripe interface.
	StripeConfig struct {
		APIKey        string `json:"apiKey"        toml:"api_key"`
		SuccessURL    string `json:"successURL"    toml:"success_url"`
		CancelURL     string `json:"cancelURL"     toml:"cancel_url"`
		WebhookSecret string `json:"webhookSecret" toml:"webhook_secret"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a StripeConfig struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	if !cfg.Enabled {
		return nil
	}

	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(StripeProvider)),
		validation.Field(&cfg.Stripe, validation.When(cfg.Provider == StripeProvider, validation.Required)),
	)
}

var _ validation.ValidatableWithContext = (*StripeConfig)(nil)

// ValidateWithContext validates a StripeConfig struct.
func (cfg *StripeConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.APIKey, validation.Required),
	)
}
