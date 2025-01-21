package cookies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Base64EncodedHashKey:  t.Name(),
			Base64EncodedBlockKey: t.Name(),
		}

		assert.NoError(t, cfg.ValidateWithContext(context.Background()))
	})

	T.Run("with missing hash key", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Base64EncodedBlockKey: t.Name(),
		}

		assert.Error(t, cfg.ValidateWithContext(context.Background()))
	})

	T.Run("with missing block key", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Base64EncodedHashKey: t.Name(),
		}

		assert.Error(t, cfg.ValidateWithContext(context.Background()))
	})
}
