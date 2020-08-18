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

func TestV1Client_BuildRecipeStepEventExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		actual, err := c.BuildRecipeStepEventExistsRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeStepEvent.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_RecipeStepEventExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeStepEvent.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.RecipeStepEventExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeStepEventExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetRecipeStepEventRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeStepEventRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeStepEvent.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeStepEvent.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepEvent))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeStepEvent, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeStepEvent.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetRecipeStepEventsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeStepEventsRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeStepEvents(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events", exampleRecipe.ID, exampleRecipeStep.ID)

		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepEventList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeStepEventList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events", exampleRecipe.ID, exampleRecipeStep.ID)

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
		actual, err := c.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateRecipeStepEventRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateRecipeStepEventRequest(ctx, exampleRecipe.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events", exampleRecipe.ID, exampleRecipeStep.ID)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.RecipeStepEventCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToRecipeStep = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepEvent))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateRecipeStepEvent(ctx, exampleRecipe.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeStepEvent, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeStepEvent(ctx, exampleRecipe.ID, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateRecipeStepEventRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateRecipeStepEventRequest(ctx, exampleRecipe.ID, exampleRecipeStepEvent)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events/%d", exampleRecipe.ID, exampleRecipeStepEvent.BelongsToRecipeStep, exampleRecipeStepEvent.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleRecipeStepEvent))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStepEvent)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		err := buildTestClientWithInvalidURL(t).UpdateRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStepEvent)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveRecipeStepEventRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveRecipeStepEventRequest(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeStepEvent.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_steps/%d/recipe_step_events/%d", exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		err := buildTestClientWithInvalidURL(t).ArchiveRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.Error(t, err, "error should be returned")
	})
}
