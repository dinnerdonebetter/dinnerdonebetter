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

func TestClient_UpdateHousehold(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()

		data := fakes.BuildFakeHousehold()
		data.WebhookEncryptionKey = ""

		expected := &types.APIResponse[*types.Household]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeHouseholdUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateHousehold(ctx, householdID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeHouseholdUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateHousehold(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeHouseholdUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateHousehold(ctx, householdID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeHouseholdUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, householdID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateHousehold(ctx, householdID, exampleInput)

		assert.Error(t, err)
	})
}
