// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_VerifyTOTPSecret(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/users/totp_secret/verify"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fakes.BuildFakeUser()
		// the hashed passwords is never transmitted over the wire.
		data.HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		data.TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		data.TwoFactorSecretVerifiedAt = nil

		expected := &types.APIResponse[*types.User]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeTOTPSecretVerificationInput(data)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.VerifyTOTPSecret(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeTOTPSecretVerificationInput(data)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.VerifyTOTPSecret(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeTOTPSecretVerificationInput(data)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.VerifyTOTPSecret(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
