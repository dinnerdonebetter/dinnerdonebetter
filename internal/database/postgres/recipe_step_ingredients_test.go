package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_RecipeStepIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, "", exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, "", exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, "", exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, "", exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_getRecipeStepIngredientsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.getRecipeStepIngredientsForRecipe(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredients(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_createRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.createRecipeStepIngredient(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, "", exampleRecipeStepIngredient.ID))
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, ""))
	})
}
