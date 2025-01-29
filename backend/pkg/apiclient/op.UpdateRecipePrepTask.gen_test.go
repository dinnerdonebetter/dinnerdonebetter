// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateRecipePrepTask(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipePrepTaskID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*RecipePrepTask](t)

		expected := &APIResponse[*RecipePrepTask]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*RecipePrepTaskUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipePrepTaskID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateRecipePrepTask(ctx, recipeID, recipePrepTaskID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipePrepTaskID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipePrepTaskUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipePrepTask(ctx, "", recipePrepTaskID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid recipePrepTask ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipePrepTaskUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateRecipePrepTask(ctx, recipeID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipePrepTaskID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipePrepTaskUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateRecipePrepTask(ctx, recipeID, recipePrepTaskID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipePrepTaskID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*RecipePrepTaskUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, recipeID, recipePrepTaskID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateRecipePrepTask(ctx, recipeID, recipePrepTaskID, exampleInput)

		assert.Error(t, err)
	})
}
