// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/lib/internal/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetHouseholdInvitation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/household_invitations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInvitationID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*HouseholdInvitation](t)
		expected := &APIResponse[*HouseholdInvitation]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetHouseholdInvitation(ctx, householdInvitationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid householdInvitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInvitation(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInvitationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholdInvitation(ctx, householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInvitationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholdInvitation(ctx, householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
