package aes

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardEncryptor(T *testing.T) {
	T.Parallel()

	T.Run("basic operation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expected := t.Name()
		secret, err := random.GenerateHexEncodedString(ctx, 16)
		require.NoError(t, err)

		encryptor, err := NewEncryptorDecryptor(tracing.NewNoopTracerProvider(), logging.NewNoopLogger(), []byte(secret))
		require.NotNil(t, encryptor)
		require.NoError(t, err)

		encrypted, err := encryptor.Encrypt(ctx, expected)
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		actual, err := encryptor.Decrypt(ctx, encrypted)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
