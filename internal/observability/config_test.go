package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/metrics"
	tracingcfg "github.com/prixfixeco/api_server/internal/observability/tracing/config"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.Jaeger,
			},
			Metrics: metrics.Config{
				Provider:                         metrics.Prometheus,
				RuntimeMetricsCollectionInterval: metrics.DefaultMetricsCollectionInterval,
			},
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
