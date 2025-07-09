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

func TestClient_GetAccountInvitationByID(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/accounts/%s/invitations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()
		accountInvitationID := fake.BuildFakeID()

		data := &AccountInvitation{}
		expected := &APIResponse[*AccountInvitation]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, accountID, accountInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetAccountInvitationByID(ctx, accountID, accountInvitationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid account ID", func(t *testing.T) {
		t.Parallel()

		accountInvitationID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetAccountInvitationByID(ctx, "", accountInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid accountInvitation ID", func(t *testing.T) {
		t.Parallel()

		accountID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetAccountInvitationByID(ctx, accountID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()
		accountInvitationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetAccountInvitationByID(ctx, accountID, accountInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()
		accountInvitationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, accountID, accountInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetAccountInvitationByID(ctx, accountID, accountInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
