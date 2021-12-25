package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider: ProviderZerolog,
		}

		assert.NotNil(t, cfg.ProvideLogger())
	})

	T.Run("no provider", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}

		assert.NotNil(t, cfg.ProvideLogger())
	})
}
