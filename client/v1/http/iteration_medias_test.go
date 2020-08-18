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

func TestV1Client_BuildIterationMediaExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		actual, err := c.BuildIterationMediaExistsRequest(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleIterationMedia.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_IterationMediaExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleIterationMedia.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias/%d", exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.IterationMediaExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.IterationMediaExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetIterationMediaRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetIterationMediaRequest(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleIterationMedia.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleIterationMedia.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias/%d", exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleIterationMedia))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleIterationMedia, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleIterationMedia.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias/%d", exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetIterationMediasRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetIterationMediasRequest(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetIterationMedias(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		filter := (*models.QueryFilter)(nil)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias", exampleRecipe.ID, exampleRecipeIteration.ID)

		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleIterationMediaList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleIterationMediaList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		filter := (*models.QueryFilter)(nil)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias", exampleRecipe.ID, exampleRecipeIteration.ID)

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
		actual, err := c.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateIterationMediaRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateIterationMediaRequest(ctx, exampleRecipe.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		expectedPath := fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias", exampleRecipe.ID, exampleRecipeIteration.ID)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.IterationMediaCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToRecipeIteration = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleIterationMedia))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateIterationMedia(ctx, exampleRecipe.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleIterationMedia, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateIterationMedia(ctx, exampleRecipe.ID, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateIterationMediaRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateIterationMediaRequest(ctx, exampleRecipe.ID, exampleIterationMedia)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias/%d", exampleRecipe.ID, exampleIterationMedia.BelongsToRecipeIteration, exampleIterationMedia.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleIterationMedia))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateIterationMedia(ctx, exampleRecipe.ID, exampleIterationMedia)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		err := buildTestClientWithInvalidURL(t).UpdateIterationMedia(ctx, exampleRecipe.ID, exampleIterationMedia)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveIterationMediaRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveIterationMediaRequest(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleIterationMedia.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/recipes/%d/recipe_iterations/%d/iteration_medias/%d", exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		err := buildTestClientWithInvalidURL(t).ArchiveIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.Error(t, err, "error should be returned")
	})
}
