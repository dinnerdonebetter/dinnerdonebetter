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

func TestV1Client_BuildRecipeIterationStepExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		actual, err := c.BuildRecipeIterationStepExistsRequest(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeIterationStep.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_RecipeIterationStepExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeIterationStep.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps/%d", exampleRecipe.ID, exampleRecipeIterationStep.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.RecipeIterationStepExists(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeIterationStepExists(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetRecipeIterationStepRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeIterationStepRequest(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeIterationStep.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeIterationStep.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps/%d", exampleRecipe.ID, exampleRecipeIterationStep.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeIterationStep))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeIterationStep, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeIterationStep.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps/%d", exampleRecipe.ID, exampleRecipeIterationStep.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetRecipeIterationStepsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeIterationStepsRequest(ctx, exampleRecipe.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeIterationSteps(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)

		exampleRecipeIterationStepList := fakemodels.BuildFakeRecipeIterationStepList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps", exampleRecipe.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeIterationStepList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeIterationStepList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps", exampleRecipe.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateRecipeIterationStepRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateRecipeIterationStepRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps", exampleRecipe.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.RecipeIterationStepCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToRecipe = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeIterationStep))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateRecipeIterationStep(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeIterationStep, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeIterationStep(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateRecipeIterationStepRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateRecipeIterationStepRequest(ctx, exampleRecipeIterationStep)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps/%d", exampleRecipeIterationStep.BelongsToRecipe, exampleRecipeIterationStep.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleRecipeIterationStep))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateRecipeIterationStep(ctx, exampleRecipeIterationStep)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		err := buildTestClientWithInvalidURL(t).UpdateRecipeIterationStep(ctx, exampleRecipeIterationStep)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveRecipeIterationStepRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveRecipeIterationStepRequest(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeIterationStep.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iteration_steps/%d", exampleRecipe.ID, exampleRecipeIterationStep.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		err := buildTestClientWithInvalidURL(t).ArchiveRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.Error(t, err, "error should be returned")
	})
}
