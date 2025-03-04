package analyticscfg

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/analytics/posthog"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics/rudderstack"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics/segment"
	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/stretchr/testify/require"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cfg := &Config{
			Provider: ProviderSegment,
			Segment:  &segment.Config{APIToken: t.Name()},
		}

		require.NoError(t, cfg.ValidateWithContext(ctx))
	})

	T.Run("with invalid token", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cfg := &Config{
			Provider: ProviderSegment,
		}

		require.Error(t, cfg.ValidateWithContext(ctx))
	})
}

func TestConfig_ProvideCollector(T *testing.T) {
	T.Parallel()

	allProviders := []string{
		ProviderSegment,
		ProviderRudderstack,
		ProviderPostHog,
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		for _, provider := range allProviders {
			cfg := &Config{
				Provider:       provider,
				Segment:        &segment.Config{APIToken: t.Name()},
				Rudderstack:    &rudderstack.Config{DataPlaneURL: t.Name(), APIKey: t.Name()},
				Posthog:        &posthog.Config{APIKey: t.Name()},
				CircuitBreaker: circuitbreaking.Config{},
			}

			_, err := cfg.ProvideCollector(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider())
			require.NoError(t, err)
		}
	})

	T.Run("with invalid values", func(t *testing.T) {
		t.Parallel()

		for _, provider := range allProviders {
			cfg := &Config{
				Provider:    provider,
				Segment:     &segment.Config{},
				Rudderstack: &rudderstack.Config{},
				Posthog:     &posthog.Config{},
			}

			_, err := cfg.ProvideCollector(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), metrics.NewNoopMetricsProvider())
			require.Error(t, err)
		}
	})
}
