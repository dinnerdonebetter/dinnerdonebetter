package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Tracing: tracing.Config{
				Provider: tracing.Jaeger,
			},
			Metrics: metrics.Config{
				Provider:                         metrics.Prometheus,
				RuntimeMetricsCollectionInterval: metrics.DefaultMetricsCollectionInterval,
			},
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
