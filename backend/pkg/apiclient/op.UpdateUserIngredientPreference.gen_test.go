// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/user_ingredient_preferences/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userIngredientPreferenceID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*UserIngredientPreference](t)

		expected := &APIResponse[*UserIngredientPreference]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*UserIngredientPreferenceUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, userIngredientPreferenceID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateUserIngredientPreference(ctx, userIngredientPreferenceID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid userIngredientPreference ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*UserIngredientPreferenceUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateUserIngredientPreference(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userIngredientPreferenceID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*UserIngredientPreferenceUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateUserIngredientPreference(ctx, userIngredientPreferenceID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userIngredientPreferenceID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*UserIngredientPreferenceUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, userIngredientPreferenceID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateUserIngredientPreference(ctx, userIngredientPreferenceID, exampleInput)

		assert.Error(t, err)
	})
}
