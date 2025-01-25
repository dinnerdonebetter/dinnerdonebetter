package featureflagscfg

import (
	"context"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags/launchdarkly"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags/posthog"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

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
		LaunchDarkly   *launchdarkly.Config   `env:"init"     envPrefix:"LAUNCH_DARKLY"     json:"launchDarkly"`
		PostHog        *posthog.Config        `env:"init"     envPrefix:"POSTHOG_"          json:"posthog"`
		Provider       string                 `env:"PROVIDER" json:"provider"`
		CircuitBreaker circuitbreaking.Config `env:"init"     envPrefix:"CIRCUIT_BREAKING_" json:"circuitBreakingConfig"`
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

func (c *Config) ProvideFeatureFlagManager(logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client, circuitBreaker circuitbreaking.CircuitBreaker) (featureflags.FeatureFlagManager, error) {
	switch strings.TrimSpace(strings.ToLower(c.Provider)) {
	case ProviderLaunchDarkly:
		return launchdarkly.NewFeatureFlagManager(c.LaunchDarkly, logger, tracerProvider, httpClient, circuitBreaker)
	case ProviderPostHog:
		return posthog.NewFeatureFlagManager(c.PostHog, logger, tracerProvider, circuitBreaker)
	default:
		return featureflags.NewNoopFeatureFlagManager(), nil
	}
}
