package observability

import (
	"context"
	"testing"

	"github.com/prixfixeco/api_server/internal/observability/tracing/jaeger"

	"github.com/prixfixeco/api_server/internal/observability/metrics/prometheus"

	"github.com/prixfixeco/api_server/internal/observability/metrics/config"

	"github.com/stretchr/testify/assert"

	tracingcfg "github.com/prixfixeco/api_server/internal/observability/tracing/config"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderJaeger,
				Jaeger: &jaeger.Config{
					CollectorEndpoint:         "0.0.0.0",
					ServiceName:               t.Name(),
					SpanCollectionProbability: 1,
				},
			},
			Metrics: config.Config{
				Provider: config.ProviderPrometheus,
				Prometheus: &prometheus.Config{
					RuntimeMetricsCollectionInterval: config.DefaultMetricsCollectionInterval,
				},
			},
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
