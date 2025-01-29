// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientGroupID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidIngredientGroup](t)

		expected := &APIResponse[*ValidIngredientGroup]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidIngredientGroupUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientGroupID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientGroup(ctx, validIngredientGroupID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientGroup ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidIngredientGroupUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientGroup(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientGroupID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidIngredientGroupUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientGroup(ctx, validIngredientGroupID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientGroupID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidIngredientGroupUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientGroupID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientGroup(ctx, validIngredientGroupID, exampleInput)

		assert.Error(t, err)
	})
}
