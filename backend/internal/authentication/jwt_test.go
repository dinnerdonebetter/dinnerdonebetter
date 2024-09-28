package authentication

import (
	"crypto/ed25519"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	exampleSigningKey = "HEREISA32CHARSECRETWHICHISMADEUP"
	ed25519SigningKey = "HEREISA64CHARSECRETWHICHISMADEUPHEREISA64CHARSECRETWHICHISMADEUP"
)

func Test_jwtSigner_IssueJWT(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.StZO4Shf1x3YHTqX8V3V03kSA83qepATCaisem6hnNI"
		signer, err := NewJWTSigner([]byte(exampleSigningKey))
		require.NoError(t, err)

		actual, err := signer.IssueJWT()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_jwtSigner_ParseJWT(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.StZO4Shf1x3YHTqX8V3V03kSA83qepATCaisem6hnNI"
		signer, err := NewJWTSigner([]byte(exampleSigningKey))
		require.NoError(t, err)

		actual, err := signer.ParseJWT(exampleToken)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	T.Run("with invalid algo", func(t *testing.T) {
		t.Parallel()

		token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{})

		cryptoSigner := ed25519.PrivateKey(ed25519SigningKey)
		tokenString, err := token.SignedString(cryptoSigner)
		require.NoError(t, err)

		signer, err := NewJWTSigner([]byte(exampleSigningKey))
		require.NoError(t, err)

		actual, err := signer.ParseJWT(tokenString)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid key", func(t *testing.T) {
		t.Parallel()

		exampleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.StZO4Shf1x3YHTqX8V3V03kSA83qepATCaisem6hnNI"
		signer, err := NewJWTSigner([]byte(exampleSigningKey))
		require.NoError(t, err)

		signer.(*jwtSigner).signingKey = nil

		actual, err := signer.ParseJWT(exampleToken)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}
