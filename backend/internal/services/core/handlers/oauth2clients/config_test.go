package oauth2clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cfg := &Config{
			OAuth2ClientCreationDisabled: true,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
