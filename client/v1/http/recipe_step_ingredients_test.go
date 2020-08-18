package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV1Client_BuildRecipeStepIngredientExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		actual, err := c.BuildRecipeStepIngredientExistsRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeStepIngredient.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_RecipeStepIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeStepIngredient.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeStepIngredientRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeStepIngredient.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeStepIngredient.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepIngredient))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeStepIngredient, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeStepIngredient.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetRecipeStepIngredientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeStepIngredientsRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients", exampleRecipe.ID, exampleRecipeStep.ID)

		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepIngredientList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeStepIngredientList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients", exampleRecipe.ID, exampleRecipeStep.ID)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateRecipeStepIngredientRequest(ctx, exampleRecipe.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients", exampleRecipe.ID, exampleRecipeStep.ID)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.RecipeStepIngredientCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToRecipeStep = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepIngredient))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateRecipeStepIngredient(ctx, exampleRecipe.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeStepIngredient, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeStepIngredient(ctx, exampleRecipe.ID, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateRecipeStepIngredientRequest(ctx, exampleRecipe.ID, exampleRecipeStepIngredient)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d", exampleRecipe.ID, exampleRecipeStepIngredient.BelongsToRecipeStep, exampleRecipeStepIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepIngredient))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStepIngredient)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		err := buildTestClientWithInvalidURL(t).UpdateRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStepIngredient)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveRecipeStepIngredientRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeStepIngredient.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		err := buildTestClientWithInvalidURL(t).ArchiveRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err, "error should be returned")
	})
}
