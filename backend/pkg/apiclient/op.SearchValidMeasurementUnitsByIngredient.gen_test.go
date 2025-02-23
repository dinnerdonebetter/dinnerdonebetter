// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchValidMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_units/by_ingredient/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()
		validIngredientID := fake.BuildFakeID()

		list := []*ValidMeasurementUnit{}
		exampleResponse := &APIResponse[[]*ValidMeasurementUnit]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[ValidMeasurementUnit]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat, validIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.SearchValidMeasurementUnitsByIngredient(ctx, q, validIngredientID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty validIngredient ID", func(t *testing.T) {
		t.Parallel()

		q := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidMeasurementUnitsByIngredient(ctx, q, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()
		validIngredientID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchValidMeasurementUnitsByIngredient(ctx, q, validIngredientID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()
		validIngredientID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat, validIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidMeasurementUnitsByIngredient(ctx, q, validIngredientID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
