package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

func TestProvideCollector(T *testing.T) {
	T.Parallel()

	T.Run("noop", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		logger := logging.NewNoopLogger()

		actual, err := ProvideCollector(cfg, logger, tracing.NewNoopTracerProvider())
		require.NoError(t, err)
		require.NotNil(t, actual)
	})

	T.Run("with segment", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider: ProviderSegment,
			APIToken: t.Name(),
		}
		logger := logging.NewNoopLogger()

		actual, err := ProvideCollector(cfg, logger, tracing.NewNoopTracerProvider())
		require.NoError(t, err)
		require.NotNil(t, actual)
	})
}
