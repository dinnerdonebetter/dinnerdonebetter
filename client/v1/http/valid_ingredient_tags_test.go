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

func TestV1Client_BuildValidIngredientTagExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		actual, err := c.BuildValidIngredientTagExistsRequest(ctx, exampleValidIngredientTag.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidIngredientTag.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ValidIngredientTagExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidIngredientTag.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredient_tags/%d", exampleValidIngredientTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.ValidIngredientTagExists(ctx, exampleValidIngredientTag.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidIngredientTagExists(ctx, exampleValidIngredientTag.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetValidIngredientTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetValidIngredientTagRequest(ctx, exampleValidIngredientTag.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidIngredientTag.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidIngredientTag.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredient_tags/%d", exampleValidIngredientTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientTag))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientTag(ctx, exampleValidIngredientTag.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidIngredientTag, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientTag(ctx, exampleValidIngredientTag.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidIngredientTag.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredient_tags/%d", exampleValidIngredientTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientTag(ctx, exampleValidIngredientTag.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetValidIngredientTagsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetValidIngredientTagsRequest(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetValidIngredientTags(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		exampleValidIngredientTagList := fakemodels.BuildFakeValidIngredientTagList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/valid_ingredient_tags", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientTagList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientTags(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidIngredientTagList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientTags(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/valid_ingredient_tags", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidIngredientTags(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateValidIngredientTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateValidIngredientTagRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/api/v1/valid_ingredient_tags", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.ValidIngredientTagCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientTag))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateValidIngredientTag(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidIngredientTag, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateValidIngredientTag(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateValidIngredientTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateValidIngredientTagRequest(ctx, exampleValidIngredientTag)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredient_tags/%d", exampleValidIngredientTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleValidIngredientTag))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateValidIngredientTag(ctx, exampleValidIngredientTag)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		err := buildTestClientWithInvalidURL(t).UpdateValidIngredientTag(ctx, exampleValidIngredientTag)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveValidIngredientTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveValidIngredientTagRequest(ctx, exampleValidIngredientTag.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidIngredientTag.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_ingredient_tags/%d", exampleValidIngredientTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		err := buildTestClientWithInvalidURL(t).ArchiveValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.Error(t, err, "error should be returned")
	})
}
