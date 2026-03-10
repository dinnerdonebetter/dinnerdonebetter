package secrets

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ProvideSecretSource(t *testing.T) {
	t.Parallel()

	t.Run("nil config returns env source", func(t *testing.T) {
		t.Parallel()

		var cfg *Config
		source, err := cfg.ProvideSecretSource()
		require.NoError(t, err)
		require.NotNil(t, source)

		key := "TEST_NIL_CONFIG_" + t.Name()
		value := "from-env"
		require.NoError(t, os.Setenv(key, value))
		t.Cleanup(func() { _ = os.Unsetenv(key) })

		got, err := source.GetSecret(context.Background(), key)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	})

	t.Run("empty provider returns env source", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{Provider: ""}
		source, err := cfg.ProvideSecretSource()
		require.NoError(t, err)
		require.NotNil(t, source)

		key := "TEST_EMPTY_PROVIDER_" + t.Name()
		value := "from-env"
		require.NoError(t, os.Setenv(key, value))
		t.Cleanup(func() { _ = os.Unsetenv(key) })

		got, err := source.GetSecret(context.Background(), key)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	})

	t.Run("env provider returns env source", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{Provider: ProviderEnv}
		source, err := cfg.ProvideSecretSource()
		require.NoError(t, err)
		require.NotNil(t, source)

		key := "TEST_ENV_PROVIDER_" + t.Name()
		value := "from-env"
		require.NoError(t, os.Setenv(key, value))
		t.Cleanup(func() { _ = os.Unsetenv(key) })

		got, err := source.GetSecret(context.Background(), key)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	})

	t.Run("noop provider returns noop source", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{Provider: ProviderNoop}
		source, err := cfg.ProvideSecretSource()
		require.NoError(t, err)
		require.NotNil(t, source)

		got, err := source.GetSecret(context.Background(), "any")
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("unknown provider returns error", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{Provider: "vault"}
		source, err := cfg.ProvideSecretSource()
		require.Error(t, err)
		assert.Nil(t, source)
		assert.Contains(t, err.Error(), "unknown")
	})
}
