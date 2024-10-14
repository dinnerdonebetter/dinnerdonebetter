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

func TestClient_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidIngredient()
		expected := &types.APIResponse[*types.ValidIngredient]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidIngredientUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredient(ctx, validIngredientID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredient ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidIngredientUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredient(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredient(ctx, validIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredient(ctx, validIngredientID, exampleInput)

		assert.Error(t, err)
	})
}
