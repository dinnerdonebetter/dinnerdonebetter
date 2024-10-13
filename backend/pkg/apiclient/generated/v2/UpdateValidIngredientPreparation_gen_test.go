// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientPreparationID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidIngredientPreparation()
		expected := &types.APIResponse[*types.ValidIngredientPreparation]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientPreparation(ctx, validIngredientPreparationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientPreparation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientPreparation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientPreparationID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientPreparation(ctx, validIngredientPreparationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientPreparationID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientPreparation(ctx, validIngredientPreparationID, exampleInput)

		assert.Error(t, err)
	})
}
