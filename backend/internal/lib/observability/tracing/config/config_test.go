package tracingcfg

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/oteltrace"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ProvideTracerProvider(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}

		tracerProvider, err := cfg.ProvideTracerProvider(
			context.Background(),
			logging.NewNoopLogger(),
		)

		assert.NoError(t, err)
		assert.NotNil(t, tracerProvider)
	})
}

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:                  ProviderOtel,
			ServiceName:               t.Name(),
			SpanCollectionProbability: 1,
			Otel: &oteltrace.Config{
				CollectorEndpoint: t.Name(),
			},
		}

		assert.NoError(t, cfg.ValidateWithContext(context.Background()))
	})
}
