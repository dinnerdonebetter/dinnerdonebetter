package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			APIToken: t.Name(),
		}

		require.NoError(t, cfg.ValidateWithContext(ctx))
	})

	T.Run("with invalid token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Provider: "segment",
		}

		require.Error(t, cfg.ValidateWithContext(ctx))
	})
}
