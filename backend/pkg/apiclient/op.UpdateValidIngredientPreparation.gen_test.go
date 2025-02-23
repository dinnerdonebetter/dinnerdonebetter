// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientPreparationID := fake.BuildFakeID()

		data := &ValidIngredientPreparation{}
		expected := &APIResponse[*ValidIngredientPreparation]{
			Data: data,
		}

		exampleInput := &ValidIngredientPreparationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientPreparation(ctx, validIngredientPreparationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientPreparation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &ValidIngredientPreparationUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientPreparation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientPreparationID := fake.BuildFakeID()

		exampleInput := &ValidIngredientPreparationUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientPreparation(ctx, validIngredientPreparationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientPreparationID := fake.BuildFakeID()

		exampleInput := &ValidIngredientPreparationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientPreparation(ctx, validIngredientPreparationID, exampleInput)

		assert.Error(t, err)
	})
}
