package authentication

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
			Debug:                 false,
			EnableUserSignup:      false,
			MinimumUsernameLength: 123,
			MinimumPasswordLength: 123,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
