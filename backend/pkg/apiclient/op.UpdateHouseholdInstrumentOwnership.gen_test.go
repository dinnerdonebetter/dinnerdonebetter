// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInstrumentOwnershipID := fake.BuildFakeID()

		data := &HouseholdInstrumentOwnership{}
		expected := &APIResponse[*HouseholdInstrumentOwnership]{
			Data: data,
		}

		exampleInput := &HouseholdInstrumentOwnershipUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdInstrumentOwnershipID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid householdInstrumentOwnership ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &HouseholdInstrumentOwnershipUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHouseholdInstrumentOwnership(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInstrumentOwnershipID := fake.BuildFakeID()

		exampleInput := &HouseholdInstrumentOwnershipUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInstrumentOwnershipID := fake.BuildFakeID()

		exampleInput := &HouseholdInstrumentOwnershipUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdInstrumentOwnershipID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, exampleInput)

		assert.Error(t, err)
	})
}
