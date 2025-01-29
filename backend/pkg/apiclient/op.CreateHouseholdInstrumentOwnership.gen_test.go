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

func TestClient_CreateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fake.BuildFakeForTest[*HouseholdInstrumentOwnership](t)

		expected := &APIResponse[*HouseholdInstrumentOwnership]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*HouseholdInstrumentOwnershipCreationRequestInput](t)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fake.BuildFakeForTest[*HouseholdInstrumentOwnershipCreationRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fake.BuildFakeForTest[*HouseholdInstrumentOwnershipCreationRequestInput](t)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
