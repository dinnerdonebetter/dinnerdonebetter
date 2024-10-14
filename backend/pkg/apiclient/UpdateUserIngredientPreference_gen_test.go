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

func TestClient_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_ingredient_preferences/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userIngredientPreferenceID := fakes.BuildFakeID()

		data := fakes.BuildFakeUserIngredientPreference()
		expected := &types.APIResponse[*types.UserIngredientPreference]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeUserIngredientPreferenceUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, userIngredientPreferenceID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateUserIngredientPreference(ctx, userIngredientPreferenceID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid userIngredientPreference ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeUserIngredientPreferenceUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateUserIngredientPreference(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userIngredientPreferenceID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserIngredientPreferenceUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateUserIngredientPreference(ctx, userIngredientPreferenceID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userIngredientPreferenceID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserIngredientPreferenceUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, userIngredientPreferenceID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateUserIngredientPreference(ctx, userIngredientPreferenceID, exampleInput)

		assert.Error(t, err)
	})
}
