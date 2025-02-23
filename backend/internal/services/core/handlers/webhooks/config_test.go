package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		cfg := &Config{
			Debug: false,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
