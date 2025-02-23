// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateHousehold(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdID := fake.BuildFakeID()

		data := &Household{}
		expected := &APIResponse[*Household]{
			Data: data,
		}

		exampleInput := &HouseholdUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateHousehold(ctx, householdID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &HouseholdUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHousehold(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdID := fake.BuildFakeID()

		exampleInput := &HouseholdUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateHousehold(ctx, householdID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdID := fake.BuildFakeID()

		exampleInput := &HouseholdUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateHousehold(ctx, householdID, exampleInput)

		assert.Error(t, err)
	})
}
