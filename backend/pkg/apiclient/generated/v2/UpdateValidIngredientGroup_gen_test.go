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

func TestClient_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_groups/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientGroupID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidIngredientGroup()
		expected := &types.APIResponse[*types.ValidIngredientGroup]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientGroupID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientGroup(ctx, validIngredientGroupID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientGroup ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientGroup(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientGroupID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientGroup(ctx, validIngredientGroupID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientGroupID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientGroupID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientGroup(ctx, validIngredientGroupID, exampleInput)

		assert.Error(t, err)
	})
}
