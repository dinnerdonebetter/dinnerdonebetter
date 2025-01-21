package cookies

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKey = "HEREISA32CHARSECRETWHICHISMADEUP"
)

func buildConfigForTest() *Config {
	return &Config{
		Base64EncodedHashKey:  base64.RawURLEncoding.EncodeToString([]byte(testKey)),
		Base64EncodedBlockKey: base64.RawURLEncoding.EncodeToString([]byte(testKey)),
	}
}

func TestNewCookieManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		m, err := NewCookieManager(buildConfigForTest(), tracing.NewNoopTracerProvider())
		assert.NoError(t, err)
		assert.NotNil(t, m)
	})
}

type example struct {
	Name string
}

func Test_manager_Encode(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, err := NewCookieManager(buildConfigForTest(), tracing.NewNoopTracerProvider())
		require.NoError(t, err)
		require.NotNil(t, m)

		actual, err := m.Encode(ctx, "test", &example{Name: t.Name()})
		require.NoError(t, err)
		assert.NotEmpty(t, actual)
	})
}

func Test_manager_Decode(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, err := NewCookieManager(buildConfigForTest(), tracing.NewNoopTracerProvider())
		require.NoError(t, err)
		require.NotNil(t, m)

		encoded, err := m.Encode(ctx, "test", &example{Name: t.Name()})
		require.NoError(t, err)
		assert.NotEmpty(t, encoded)

		var actual example
		err = m.Decode(ctx, "test", encoded, &actual)
		require.NoError(t, m.Decode(ctx, "test", encoded, &actual))
		assert.Equal(t, actual.Name, t.Name())
	})
}
