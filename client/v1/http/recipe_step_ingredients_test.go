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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV1Client_BuildGetRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		expectedID := uint64(1)
		actual, err := c.BuildGetRecipeStepIngredientRequest(ctx, expectedID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", expectedID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expected := &models.RecipeStepIngredient{
			ID: 1,
		}

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(expected.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipe_step_ingredients/%d", expected.ID), "expected and actual path don't match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(expected))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepIngredient(ctx, expected.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, expected, actual)
	})
}

func TestV1Client_BuildGetRecipeStepIngredientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeStepIngredientsRequest(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expected := &models.RecipeStepIngredientList{
			RecipeStepIngredients: []models.RecipeStepIngredient{
				{
					ID: 1,
				},
			},
		}

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/recipe_step_ingredients", "expected and actual path don't match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(expected))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepIngredients(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, expected, actual)
	})
}

func TestV1Client_BuildCreateRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		exampleInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    func(x uint64) *uint64 { return &x }(1),
			QuantityType:    "example",
			QuantityValue:   1.23,
			QuantityNotes:   "example",
			ProductOfRecipe: false,
			IngredientNotes: "example",
			RecipeStepID:    1,
		}
		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateRecipeStepIngredientRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expected := &models.RecipeStepIngredient{
			ID:              1,
			IngredientID:    func(x uint64) *uint64 { return &x }(1),
			QuantityType:    "example",
			QuantityValue:   1.23,
			QuantityNotes:   "example",
			ProductOfRecipe: false,
			IngredientNotes: "example",
			RecipeStepID:    1,
		}
		exampleInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/recipe_step_ingredients", "expected and actual path don't match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.RecipeStepIngredientCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(expected))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, expected, actual)
	})
}

func TestV1Client_BuildUpdateRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodPut
		exampleInput := &models.RecipeStepIngredient{
			ID: 1,
		}

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateRecipeStepIngredientRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expected := &models.RecipeStepIngredient{
			ID: 1,
		}

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipe_step_ingredients/%d", expected.ID), "expected and actual path don't match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(&models.RecipeStepIngredient{}))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateRecipeStepIngredient(ctx, expected)
		assert.NoError(t, err, "no error should be returned")
	})
}

func TestV1Client_BuildArchiveRecipeStepIngredientRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		expectedID := uint64(1)
		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveRecipeStepIngredientRequest(ctx, expectedID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", expectedID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expected := uint64(1)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipe_step_ingredients/%d", expected), "expected and actual path don't match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveRecipeStepIngredient(ctx, expected)
		assert.NoError(t, err, "no error should be returned")
	})
}
