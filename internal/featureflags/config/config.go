package config

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderLaunchDarkly is used to indicate the LaunchDarkly provider.
	ProviderLaunchDarkly = "launchdarkly"
)

type (
	// LaunchDarklyConfig configures our launch darkly implementation.
	LaunchDarklyConfig struct {
		SDKKey      string
		InitTimeout time.Duration
	}

	// Config configures our feature flag managers.
	Config struct {
		LaunchDarkly *LaunchDarklyConfig
		Provider     string
	}
)

// ValidateWithContext validates the config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Provider, validation.In(ProviderLaunchDarkly)),
		validation.Field(&c.LaunchDarkly, validation.When(c.Provider == ProviderLaunchDarkly, validation.Required)),
	)
}
