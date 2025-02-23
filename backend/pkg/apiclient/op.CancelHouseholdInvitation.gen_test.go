// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_CancelHouseholdInvitation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/household_invitations/%s/cancel"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdInvitationID := fake.BuildFakeID()

		data := &HouseholdInvitation{}
		expected := &APIResponse[*HouseholdInvitation]{
			Data: data,
		}

		exampleInput := &HouseholdInvitationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.CancelHouseholdInvitation(ctx, householdInvitationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid householdInvitation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &HouseholdInvitationUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.CancelHouseholdInvitation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdInvitationID := fake.BuildFakeID()

		exampleInput := &HouseholdInvitationUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.CancelHouseholdInvitation(ctx, householdInvitationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdInvitationID := fake.BuildFakeID()

		exampleInput := &HouseholdInvitationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.CancelHouseholdInvitation(ctx, householdInvitationID, exampleInput)

		assert.Error(t, err)
	})
}
