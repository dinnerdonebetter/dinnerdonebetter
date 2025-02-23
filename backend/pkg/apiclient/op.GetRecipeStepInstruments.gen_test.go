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

func TestClient_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		list := []*RecipeStepInstrument{}
		exampleResponse := &APIResponse[[]*RecipeStepInstrument]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[RecipeStepInstrument]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, recipeID, recipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty recipe ID", func(t *testing.T) {
		t.Parallel()

		recipeStepID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstruments(ctx, "", recipeStepID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with empty recipeStep ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstruments(ctx, recipeID, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipeStepID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, recipeID, recipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
