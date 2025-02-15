// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetHouseholdInstrumentOwnership(T *testing.T) {
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

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdInstrumentOwnershipID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid householdInstrumentOwnership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetHouseholdInstrumentOwnership(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInstrumentOwnershipID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		householdInstrumentOwnershipID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, householdInstrumentOwnershipID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
