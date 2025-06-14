// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateAccountMemberPermissions(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/accounts/%s/members/%s/permissions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()
		userID := fake.BuildFakeID()

		data := &UserPermissionsResponse{}
		expected := &APIResponse[*UserPermissionsResponse]{
			Data: data,
		}

		exampleInput := &ModifyUserPermissionsInput{}

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, accountID, userID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateAccountMemberPermissions(ctx, accountID, userID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid account ID", func(t *testing.T) {
		t.Parallel()

		userID := fake.BuildFakeID()

		exampleInput := &ModifyUserPermissionsInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateAccountMemberPermissions(ctx, "", userID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		accountID := fake.BuildFakeID()

		exampleInput := &ModifyUserPermissionsInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateAccountMemberPermissions(ctx, accountID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()
		userID := fake.BuildFakeID()

		exampleInput := &ModifyUserPermissionsInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateAccountMemberPermissions(ctx, accountID, userID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()
		userID := fake.BuildFakeID()

		exampleInput := &ModifyUserPermissionsInput{}

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, accountID, userID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateAccountMemberPermissions(ctx, accountID, userID, exampleInput)

		assert.Error(t, err)
	})
}
