package config

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/featureflags"
	"github.com/dinnerdonebetter/backend/internal/featureflags/launchdarkly"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderLaunchDarkly is used to indicate the LaunchDarkly provider.
	ProviderLaunchDarkly = "launchdarkly"
)

type (
	// Config configures our feature flag managers.
	Config struct {
		LaunchDarkly *launchdarkly.Config
		Provider     string
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Provider, validation.In(ProviderLaunchDarkly)),
		validation.Field(&c.LaunchDarkly, validation.When(c.Provider == ProviderLaunchDarkly, validation.Required)),
	)
}

func (c *Config) ProvideFeatureFlagManager(logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client) (featureflags.FeatureFlagManager, error) {
	switch c.Provider {
	case ProviderLaunchDarkly:
		return launchdarkly.NewFeatureFlagManager(c.LaunchDarkly, logger, tracerProvider, httpClient)
	default:
		return featureflags.NewNoopFeatureFlagManager(), nil
	}
}
