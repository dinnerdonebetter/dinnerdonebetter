// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_AcceptHouseholdInvitation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/household_invitations/%s/accept"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInvitationID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*HouseholdInvitation](t)

		expected := &APIResponse[*HouseholdInvitation]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*HouseholdInvitationUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.AcceptHouseholdInvitation(ctx, householdInvitationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid householdInvitation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*HouseholdInvitationUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.AcceptHouseholdInvitation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInvitationID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*HouseholdInvitationUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.AcceptHouseholdInvitation(ctx, householdInvitationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInvitationID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*HouseholdInvitationUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.AcceptHouseholdInvitation(ctx, householdInvitationID, exampleInput)

		assert.Error(t, err)
	})
}
