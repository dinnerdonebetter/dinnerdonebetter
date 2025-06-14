// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_AcceptAccountInvitation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/account_invitations/%s/accept"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInvitationID := fake.BuildFakeID()

		data := &AccountInvitation{}
		expected := &APIResponse[*AccountInvitation]{
			Data: data,
		}

		exampleInput := &AccountInvitationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, accountInvitationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.AcceptAccountInvitation(ctx, accountInvitationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid accountInvitation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &AccountInvitationUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.AcceptAccountInvitation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInvitationID := fake.BuildFakeID()

		exampleInput := &AccountInvitationUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.AcceptAccountInvitation(ctx, accountInvitationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInvitationID := fake.BuildFakeID()

		exampleInput := &AccountInvitationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, accountInvitationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.AcceptAccountInvitation(ctx, accountInvitationID, exampleInput)

		assert.Error(t, err)
	})
}
