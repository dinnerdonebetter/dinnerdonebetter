package loggingcfg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Provider: ProviderZerolog,
		}

		l, err := cfg.ProvideLogger(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})

	T.Run("no provider", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}
		l, err := cfg.ProvideLogger(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})
}
