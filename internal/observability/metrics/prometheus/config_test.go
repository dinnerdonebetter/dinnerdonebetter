package prometheus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func TestConfig_ProvideInstrumentationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			RuntimeMetricsCollectionInterval: minimumRuntimeCollectionInterval,
		}

		actual, err := cfg.ProvideInstrumentationHandler(logging.NewNoopLogger())
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			RuntimeMetricsCollectionInterval: minimumRuntimeCollectionInterval,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}

func Test_initiatePrometheusExporter(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		initiatePrometheusExporter()
	})
}
