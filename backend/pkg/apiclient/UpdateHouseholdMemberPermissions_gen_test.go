// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateHouseholdMemberPermissions(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s/members/%s/permissions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()
		userID := fakes.BuildFakeID()

		data := fakes.BuildFakeUserPermissionsResponse()
		expected := &types.APIResponse[*types.UserPermissionsResponse]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeModifyUserPermissionsInput()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, householdID, userID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, userID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		userID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeModifyUserPermissionsInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHouseholdMemberPermissions(ctx, "", userID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		householdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeModifyUserPermissionsInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()
		userID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeModifyUserPermissionsInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, userID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()
		userID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeModifyUserPermissionsInput()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, householdID, userID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, userID, exampleInput)

		assert.Error(t, err)
	})
}
