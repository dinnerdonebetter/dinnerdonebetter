// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_TransferHouseholdOwnership(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s/transfer"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdID := fake.BuildFakeID()

		data := &Household{}
		expected := &APIResponse[*Household]{
			Data: data,
		}

		exampleInput := &HouseholdOwnershipTransferInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, householdID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.TransferHouseholdOwnership(ctx, householdID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &HouseholdOwnershipTransferInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.TransferHouseholdOwnership(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdID := fake.BuildFakeID()

		exampleInput := &HouseholdOwnershipTransferInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.TransferHouseholdOwnership(ctx, householdID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		householdID := fake.BuildFakeID()

		exampleInput := &HouseholdOwnershipTransferInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, householdID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.TransferHouseholdOwnership(ctx, householdID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
