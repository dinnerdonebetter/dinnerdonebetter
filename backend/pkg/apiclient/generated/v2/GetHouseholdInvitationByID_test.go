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

func TestClient_GetHouseholdInvitationByID(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s/invitations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()
		householdInvitationID := fakes.BuildFakeID()

		data := fakes.BuildFakeHouseholdInvitation()
		data.DestinationHousehold.WebhookEncryptionKey = ""
		data.FromUser.TwoFactorSecret = ""
		expected := &types.APIResponse[*types.HouseholdInvitation]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdID, householdInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, householdInvitationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected.Data, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		householdInvitationID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInvitationByID(ctx, "", householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid householdInvitation ID", func(t *testing.T) {
		t.Parallel()

		householdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()
		householdInvitationID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()
		householdInvitationID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdID, householdInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
