package secrets

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvSecretSource_GetSecret(t *testing.T) {
	t.Parallel()

	t.Run("returns set env var", func(t *testing.T) {
		t.Parallel()

		key := "TEST_SECRET_" + t.Name()
		value := "secret-value"
		require.NoError(t, os.Setenv(key, value))
		t.Cleanup(func() { _ = os.Unsetenv(key) })

		source := NewEnvSecretSource()
		ctx := context.Background()

		got, err := source.GetSecret(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	})

	t.Run("returns empty for unset env var", func(t *testing.T) {
		t.Parallel()

		key := "TEST_SECRET_UNSET_" + t.Name()
		require.NoError(t, os.Unsetenv(key))

		source := NewEnvSecretSource()
		ctx := context.Background()

		got, err := source.GetSecret(ctx, key)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("Close is safe", func(t *testing.T) {
		t.Parallel()

		source := NewEnvSecretSource()
		err := source.Close()
		require.NoError(t, err)
	})
}

func TestNoopSecretSource_GetSecret(t *testing.T) {
	t.Parallel()

	source := NewNoopSecretSource()
	ctx := context.Background()

	got, err := source.GetSecret(ctx, "any-key")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestNoopSecretSource_Close(t *testing.T) {
	t.Parallel()

	source := NewNoopSecretSource()
	err := source.Close()
	require.NoError(t, err)
}
