package featureflagscfg

import (
	"context"
	"net/http"
	"strings"

	circuitbreaking2 "github.com/dinnerdonebetter/backend/internal/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/featureflags"
	"github.com/dinnerdonebetter/backend/internal/featureflags/launchdarkly"
	"github.com/dinnerdonebetter/backend/internal/featureflags/posthog"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderLaunchDarkly is used to indicate the LaunchDarkly provider.
	ProviderLaunchDarkly = "launchdarkly"
	// ProviderPostHog is used to indicate the PostHog provider.
	ProviderPostHog = "posthog"
)

type (
	// Config configures our feature flag managers.
	Config struct {
		LaunchDarkly          *launchdarkly.Config     `env:"init"     envPrefix:"LAUNCH_DARKLY"     json:"launchDarkly"`
		PostHog               *posthog.Config          `env:"init"     envPrefix:"POSTHOG_"          json:"posthog"`
		CircuitBreakingConfig *circuitbreaking2.Config `env:"init"     envPrefix:"CIRCUIT_BREAKING_" json:"circuitBreakingConfig"`
		Provider              string                   `env:"PROVIDER" json:"provider"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Provider, validation.In(ProviderLaunchDarkly, ProviderPostHog)),
		validation.Field(&c.LaunchDarkly, validation.When(c.Provider == ProviderLaunchDarkly, validation.Required)),
		validation.Field(&c.PostHog, validation.When(c.Provider == ProviderPostHog, validation.Required)),
	)
}

func (c *Config) ProvideFeatureFlagManager(logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client, circuitBreaker circuitbreaking2.CircuitBreaker) (featureflags.FeatureFlagManager, error) {
	switch strings.TrimSpace(strings.ToLower(c.Provider)) {
	case ProviderLaunchDarkly:
		return launchdarkly.NewFeatureFlagManager(c.LaunchDarkly, logger, tracerProvider, httpClient, circuitBreaker)
	case ProviderPostHog:
		return posthog.NewFeatureFlagManager(c.PostHog, logger, tracerProvider, circuitBreaker)
	default:
		return featureflags.NewNoopFeatureFlagManager(), nil
	}
}
