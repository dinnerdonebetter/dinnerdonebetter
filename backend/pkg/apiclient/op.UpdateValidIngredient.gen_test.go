// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredients/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientID := fake.BuildFakeID()

		data := &ValidIngredient{}
		expected := &APIResponse[*ValidIngredient]{
			Data: data,
		}

		exampleInput := &ValidIngredientUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidIngredient(ctx, validIngredientID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validIngredient ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &ValidIngredientUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidIngredient(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientID := fake.BuildFakeID()

		exampleInput := &ValidIngredientUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidIngredient(ctx, validIngredientID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validIngredientID := fake.BuildFakeID()

		exampleInput := &ValidIngredientUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validIngredientID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidIngredient(ctx, validIngredientID, exampleInput)

		assert.Error(t, err)
	})
}
