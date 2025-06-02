// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateIngredientID := fake.BuildFakeID()

		data := &ValidIngredientStateIngredient{}
		expected := &APIResponse[*ValidIngredientStateIngredient]{
			Data: data,
		}

		exampleInput := &ValidIngredientStateIngredientUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientStateIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredientStateIngredient(ctx, validIngredientStateIngredientID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredientStateIngredient ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &ValidIngredientStateIngredientUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredientStateIngredient(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateIngredientID := fake.BuildFakeID()

		exampleInput := &ValidIngredientStateIngredientUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredientStateIngredient(ctx, validIngredientStateIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validIngredientStateIngredientID := fake.BuildFakeID()

		exampleInput := &ValidIngredientStateIngredientUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientStateIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredientStateIngredient(ctx, validIngredientStateIngredientID, exampleInput)

		assert.Error(t, err)
	})
}
