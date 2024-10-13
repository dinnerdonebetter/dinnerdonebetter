// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_TransferHouseholdOwnership(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%s/transfer"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()

		data := fakes.BuildFakeHousehold()
		data.WebhookEncryptionKey = ""

		expected := &types.APIResponse[*types.Household]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeHouseholdOwnershipTransferInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, householdID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.TransferHouseholdOwnership(ctx, householdID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeHouseholdOwnershipTransferInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.TransferHouseholdOwnership(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeHouseholdOwnershipTransferInput()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.TransferHouseholdOwnership(ctx, householdID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeHouseholdOwnershipTransferInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, householdID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.TransferHouseholdOwnership(ctx, householdID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
