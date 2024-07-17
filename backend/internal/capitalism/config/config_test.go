package config

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/capitalism/stripe"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Enabled:  true,
			Provider: StripeProvider,
			Stripe:   &stripe.Config{APIKey: t.Name()},
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})

	T.Run("returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Enabled: false,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})

	T.Run("with invalid config", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Enabled:  true,
			Provider: StripeProvider,
		}

		assert.Error(t, cfg.ValidateWithContext(ctx))
	})
}
