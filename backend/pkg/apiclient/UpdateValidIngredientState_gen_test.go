// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_states/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidIngredientState()
		expected := &types.APIResponse[*types.ValidIngredientState]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientStateID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientState(ctx, validIngredientStateID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientState ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientState(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientState(ctx, validIngredientStateID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientStateID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientState(ctx, validIngredientStateID, exampleInput)

		assert.Error(t, err)
	})
}
