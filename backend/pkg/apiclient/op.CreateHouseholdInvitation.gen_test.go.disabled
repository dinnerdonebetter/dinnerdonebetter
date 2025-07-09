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

func TestClient_CreateAccountInvitation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/accounts/%s/invite"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()

		data := &AccountInvitation{}
		expected := &APIResponse[*AccountInvitation]{
			Data: data,
		}

		exampleInput := &AccountInvitationCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, accountID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateAccountInvitation(ctx, accountID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid account ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &AccountInvitationCreationRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateAccountInvitation(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()

		exampleInput := &AccountInvitationCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateAccountInvitation(ctx, accountID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()

		exampleInput := &AccountInvitationCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, accountID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateAccountInvitation(ctx, accountID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
