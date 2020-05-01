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

func TestV1Client_BuildRecipeTagExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		actual, err := c.BuildRecipeTagExistsRequest(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeTag.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_RecipeTagExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeTag.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags/%d", exampleRecipe.ID, exampleRecipeTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.RecipeTagExists(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeTagExists(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetRecipeTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeTagRequest(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeTag.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeTag.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags/%d", exampleRecipe.ID, exampleRecipeTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeTag))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeTag, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleRecipeTag.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags/%d", exampleRecipe.ID, exampleRecipeTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetRecipeTagsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetRecipeTagsRequest(ctx, exampleRecipe.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetRecipeTags(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)

		exampleRecipeTagList := fakemodels.BuildFakeRecipeTagList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags", exampleRecipe.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeTagList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeTags(ctx, exampleRecipe.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeTagList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeTags(ctx, exampleRecipe.ID, filter)

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
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags", exampleRecipe.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetRecipeTags(ctx, exampleRecipe.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateRecipeTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateRecipeTagRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags", exampleRecipe.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.RecipeTagCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToRecipe = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleRecipeTag))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateRecipeTag(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleRecipeTag, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipeTag(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateRecipeTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateRecipeTagRequest(ctx, exampleRecipeTag)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags/%d", exampleRecipeTag.BelongsToRecipe, exampleRecipeTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleRecipeTag))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateRecipeTag(ctx, exampleRecipeTag)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		err := buildTestClientWithInvalidURL(t).UpdateRecipeTag(ctx, exampleRecipeTag)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveRecipeTagRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveRecipeTagRequest(ctx, exampleRecipe.ID, exampleRecipeTag.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleRecipeTag.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_tags/%d", exampleRecipe.ID, exampleRecipeTag.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		err := buildTestClientWithInvalidURL(t).ArchiveRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.Error(t, err, "error should be returned")
	})
}
