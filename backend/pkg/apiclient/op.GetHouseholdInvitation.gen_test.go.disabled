// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAccountInvitation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/account_invitations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInvitationID := fake.BuildFakeID()

		data := &AccountInvitation{}
		expected := &APIResponse[*AccountInvitation]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, accountInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetAccountInvitation(ctx, accountInvitationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid accountInvitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetAccountInvitation(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInvitationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetAccountInvitation(ctx, accountInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInvitationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, accountInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetAccountInvitation(ctx, accountInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
