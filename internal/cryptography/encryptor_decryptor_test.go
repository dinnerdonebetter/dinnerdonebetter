package cryptography

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/random"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardEncryptor(T *testing.T) {
	T.Parallel()

	T.Run("basic operation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expected := t.Name()
		secret, err := random.GenerateRawBytes(ctx, 32)
		require.NoError(t, err)

		encryptor := NewAESEncryptorDecryptor(tracing.NewNoopTracerProvider(), logging.NewNoopLogger())

		encrypted, err := encryptor.Encrypt(ctx, expected, string(secret))
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		actual, err := encryptor.Decrypt(ctx, encrypted, string(secret))
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
