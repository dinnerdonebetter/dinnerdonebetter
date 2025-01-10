package jwt

import (
	"context"
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	exampleSigningKey = testutils.Example32ByteKey
	ed25519SigningKey = testutils.Example64ByteKey
	exampleJWT        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJUZXN0X2p3dFNpZ25lcl9Jc3N1ZUpXVC9zdGFuZGFyZCIsImV4cCI6MTcyNzU3MDU0OCwiaWF0IjoxNzI3NTY5OTQ4LCJpc3MiOiJkaW5uZXJkb25lYmV0dGVyIiwianRpIjoiY3JzYTA3NnRnM3FkdG1jY3E5MTAiLCJuYmYiOjE3Mjc1Njk4ODgsInN1YiI6ImNyc2EwNzZ0ZzNxZHRtY2NxOTBnIn0.tMASrQBoYAq4n1iwOElLqUQsYOARX5T1qxo8RKhvaAg"
	exampleExpiry     = time.Minute * 10
)

func Test_jwtSigner_IssueJWT(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		signer, err := NewJWTSigner(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), t.Name(), []byte(exampleSigningKey))
		require.NoError(t, err)

		ctx := context.Background()
		user := fakes.BuildFakeUser()

		actual, err := signer.IssueToken(ctx, user, exampleExpiry)
		assert.NoError(t, err)

		parsed, err := signer.ParseUserIDFromToken(ctx, actual)
		assert.NoError(t, err)
		assert.Equal(t, parsed, user.ID)
	})
}

func Test_jwtSigner_ParseJWT(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		signer, err := NewJWTSigner(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), t.Name(), []byte(exampleSigningKey))
		require.NoError(t, err)

		ctx := context.Background()
		user := fakes.BuildFakeUser()

		exampleToken, err := signer.IssueToken(ctx, user, exampleExpiry)
		assert.NoError(t, err)

		actual, err := signer.ParseUserIDFromToken(ctx, exampleToken)
		assert.NoError(t, err)
		assert.Equal(t, actual, user.ID)
	})

	T.Run("with invalid algo", func(t *testing.T) {
		t.Parallel()

		token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{})

		cryptoSigner := ed25519.PrivateKey(ed25519SigningKey)
		tokenString, err := token.SignedString(cryptoSigner)
		require.NoError(t, err)

		ctx := context.Background()

		signer, err := NewJWTSigner(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), t.Name(), []byte(exampleSigningKey))
		require.NoError(t, err)

		actual, err := signer.ParseUserIDFromToken(ctx, tokenString)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with invalid key", func(t *testing.T) {
		t.Parallel()

		exampleToken := exampleJWT
		signer, err := NewJWTSigner(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), t.Name(), []byte(exampleSigningKey))
		require.NoError(t, err)

		signer.(*jwtSigner).signingKey = nil

		ctx := context.Background()

		actual, err := signer.ParseUserIDFromToken(ctx, exampleToken)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})
}
