// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMermaidDiagramForRecipe(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/mermaid"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		data := ""

		expected := &APIResponse[string]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMermaidDiagramForRecipe(ctx, recipeID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMermaidDiagramForRecipe(ctx, "")

		require.Empty(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMermaidDiagramForRecipe(ctx, recipeID)

		require.Empty(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		recipeID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMermaidDiagramForRecipe(ctx, recipeID)

		require.Empty(t, actual)
		assert.Error(t, err)
	})
}
