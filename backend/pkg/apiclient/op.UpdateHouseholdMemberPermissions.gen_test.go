// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateHouseholdMemberPermissions(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s/members/%s/permissions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fake.BuildFakeID()
		userID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*UserPermissionsResponse](t)

		expected := &APIResponse[*UserPermissionsResponse]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ModifyUserPermissionsInput](t)

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, householdID, userID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, userID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		userID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ModifyUserPermissionsInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHouseholdMemberPermissions(ctx, "", userID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		householdID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ModifyUserPermissionsInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fake.BuildFakeID()
		userID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ModifyUserPermissionsInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, userID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fake.BuildFakeID()
		userID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ModifyUserPermissionsInput](t)

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, householdID, userID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateHouseholdMemberPermissions(ctx, householdID, userID, exampleInput)

		assert.Error(t, err)
	})
}
