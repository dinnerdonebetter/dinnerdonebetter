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

func TestClient_GetHouseholdInvitationByID(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s/invitations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fake.BuildFakeID()
		householdInvitationID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*HouseholdInvitation](t)
		expected := &APIResponse[*HouseholdInvitation]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdID, householdInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, householdInvitationID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		householdInvitationID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInvitationByID(ctx, "", householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid householdInvitation ID", func(t *testing.T) {
		t.Parallel()

		householdID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fake.BuildFakeID()
		householdInvitationID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fake.BuildFakeID()
		householdInvitationID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdID, householdInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholdInvitationByID(ctx, householdID, householdInvitationID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
