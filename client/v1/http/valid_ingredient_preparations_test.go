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

func TestV1Client_BuildValidIngredientPreparationExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID
		actual, err := c.BuildValidIngredientPreparationExistsRequest(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidIngredientPreparation.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidIngredientPreparation.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations/%d", exampleValidIngredient.ID, exampleValidIngredientPreparation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetValidIngredientPreparationRequest(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidIngredientPreparation.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidIngredientPreparation.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations/%d", exampleValidIngredient.ID, exampleValidIngredientPreparation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientPreparation))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidIngredientPreparation, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidIngredientPreparation.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations/%d", exampleValidIngredient.ID, exampleValidIngredientPreparation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetValidIngredientPreparationsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetValidIngredientPreparationsRequest(ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations", exampleValidIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientPreparationList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidIngredientPreparationList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations", exampleValidIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateValidIngredientPreparationRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations", exampleValidIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.ValidIngredientPreparationCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToValidIngredient = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientPreparation))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateValidIngredientPreparation(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidIngredientPreparation, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateValidIngredientPreparation(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateValidIngredientPreparationRequest(ctx, exampleValidIngredientPreparation)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations/%d", exampleValidIngredientPreparation.BelongsToValidIngredient, exampleValidIngredientPreparation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientPreparation))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		err := buildTestClientWithInvalidURL(t).UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveValidIngredientPreparationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveValidIngredientPreparationRequest(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidIngredientPreparation.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/valid_ingredient_preparations/%d", exampleValidIngredient.ID, exampleValidIngredientPreparation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		err := buildTestClientWithInvalidURL(t).ArchiveValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.Error(t, err, "error should be returned")
	})
}
