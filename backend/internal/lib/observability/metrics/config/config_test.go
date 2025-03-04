package metricscfg

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/otelgrpc"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ProvideTracerProvider(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		tracerProvider, err := cfg.ProvideMetricsProvider(context.Background(), logging.NewNoopLogger())

		assert.NoError(t, err)
		assert.NotNil(t, tracerProvider)
	})
}

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider: ProviderOtel,
			Otel: &otelgrpc.Config{
				CollectorEndpoint:  t.Name(),
				CollectionInterval: 1,
			},
		}

		assert.NoError(t, cfg.ValidateWithContext(context.Background()))
	})
}
