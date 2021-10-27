package websockets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Logging: logging.Config{},
		}

		require.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
