package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Username:       t.Name(),
			Password:       t.Name(),
			QueueAddresses: []string{t.Name()},
			DB:             1,
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}
