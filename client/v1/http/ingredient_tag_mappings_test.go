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

func TestV1Client_BuildIngredientTagMappingExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		actual, err := c.BuildIngredientTagMappingExistsRequest(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleIngredientTagMapping.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_IngredientTagMappingExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleIngredientTagMapping.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings/%d", exampleValidIngredient.ID, exampleIngredientTagMapping.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.IngredientTagMappingExists(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.IngredientTagMappingExists(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetIngredientTagMappingRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetIngredientTagMappingRequest(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleIngredientTagMapping.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleIngredientTagMapping.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings/%d", exampleValidIngredient.ID, exampleIngredientTagMapping.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleIngredientTagMapping))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleIngredientTagMapping, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleIngredientTagMapping.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings/%d", exampleValidIngredient.ID, exampleIngredientTagMapping.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetIngredientTagMappingsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetIngredientTagMappingsRequest(ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetIngredientTagMappings(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)

		exampleIngredientTagMappingList := fakemodels.BuildFakeIngredientTagMappingList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings", exampleValidIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleIngredientTagMappingList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleIngredientTagMappingList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)

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
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings", exampleValidIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateIngredientTagMappingRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateIngredientTagMappingRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings", exampleValidIngredient.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.IngredientTagMappingCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToValidIngredient = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleIngredientTagMapping))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateIngredientTagMapping(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleIngredientTagMapping, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateIngredientTagMapping(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateIngredientTagMappingRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateIngredientTagMappingRequest(ctx, exampleIngredientTagMapping)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings/%d", exampleIngredientTagMapping.BelongsToValidIngredient, exampleIngredientTagMapping.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleIngredientTagMapping))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateIngredientTagMapping(ctx, exampleIngredientTagMapping)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		err := buildTestClientWithInvalidURL(t).UpdateIngredientTagMapping(ctx, exampleIngredientTagMapping)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveIngredientTagMappingRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveIngredientTagMappingRequest(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleIngredientTagMapping.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredients/%d/ingredient_tag_mappings/%d", exampleValidIngredient.ID, exampleIngredientTagMapping.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		err := buildTestClientWithInvalidURL(t).ArchiveIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.Error(t, err, "error should be returned")
	})
}
