package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeStepIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("RecipeStepIngredientExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(true, nil)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetAllRecipeStepIngredientsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeStepIngredientsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := models.DefaultQueryFilter()
		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepIngredientList, nil)

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepIngredientList, nil)

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("CreateRecipeStepIngredient", mock.Anything, exampleInput).Return(exampleRecipeStepIngredient, nil)

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		c, mockDB := buildTestClient()

		mockDB.RecipeStepIngredientDataManager.On("UpdateRecipeStepIngredient", mock.Anything, exampleRecipeStepIngredient).Return(expected)

		err := c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("ArchiveRecipeStepIngredient", mock.Anything, exampleRecipeStepIngredient.BelongsToRecipeStep, exampleRecipeStepIngredient.ID).Return(expected)

		err := c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepIngredient.BelongsToRecipeStep, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
