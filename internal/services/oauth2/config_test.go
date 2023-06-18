package oauth2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		cfg := &Config{
			AccessTokenLifespan:  1,
			RefreshTokenLifespan: 1,
			Domain:               t.Name(),
			DataChangesTopicName: t.Name(),
			Debug:                false,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
