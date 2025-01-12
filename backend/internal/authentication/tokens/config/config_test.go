package tokenscfg

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Provider:   ProviderJWT,
			Audience:   t.Name(),
			SigningKey: random.MustGenerateRawBytes(ctx, 32),
		}

		require.NoError(t, cfg.ValidateWithContext(ctx))
	})

	T.Run("with missing key", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Provider: ProviderJWT,
		}

		require.Error(t, cfg.ValidateWithContext(ctx))
	})
}

func TestConfig_ProvideTokenIssuer(T *testing.T) {
	T.Parallel()

	T.Run("with SendGrid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider:   ProviderJWT,
			Audience:   t.Name(),
			SigningKey: random.MustGenerateRawBytes(ctx, 32),
		}

		actual, err := cfg.ProvideTokenIssuer(logger, tracing.NewNoopTracerProvider())
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider: "",
		}

		actual, err := cfg.ProvideTokenIssuer(logger, tracing.NewNoopTracerProvider())
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}
