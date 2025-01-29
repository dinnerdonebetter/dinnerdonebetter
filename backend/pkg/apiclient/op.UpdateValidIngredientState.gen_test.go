// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_states/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidIngredientState](t)

		expected := &APIResponse[*ValidIngredientState]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidIngredientStateUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientStateID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientState(ctx, validIngredientStateID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientState ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidIngredientStateUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientState(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidIngredientStateUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientState(ctx, validIngredientStateID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidIngredientStateUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientStateID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientState(ctx, validIngredientStateID, exampleInput)

		assert.Error(t, err)
	})
}
