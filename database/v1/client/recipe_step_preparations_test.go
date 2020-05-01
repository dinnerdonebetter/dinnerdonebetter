package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeStepPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("RecipeStepPreparationExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(true, nil)

		actual, err := c.RecipeStepPreparationExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(exampleRecipeStepPreparation, nil)

		actual, err := c.GetRecipeStepPreparation(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("GetAllRecipeStepPreparationsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeStepPreparationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepPreparations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := models.DefaultQueryFilter()
		exampleRecipeStepPreparationList := fakemodels.BuildFakeRecipeStepPreparationList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("GetRecipeStepPreparations", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepPreparationList, nil)

		actual, err := c.GetRecipeStepPreparations(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeStepPreparationList := fakemodels.BuildFakeRecipeStepPreparationList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("GetRecipeStepPreparations", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepPreparationList, nil)

		actual, err := c.GetRecipeStepPreparations(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("CreateRecipeStepPreparation", mock.Anything, exampleInput).Return(exampleRecipeStepPreparation, nil)

		actual, err := c.CreateRecipeStepPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()

		c, mockDB := buildTestClient()

		mockDB.RecipeStepPreparationDataManager.On("UpdateRecipeStepPreparation", mock.Anything, exampleRecipeStepPreparation).Return(expected)

		err := c.UpdateRecipeStepPreparation(ctx, exampleRecipeStepPreparation)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepPreparationDataManager.On("ArchiveRecipeStepPreparation", mock.Anything, exampleRecipeStepPreparation.BelongsToRecipeStep, exampleRecipeStepPreparation.ID).Return(expected)

		err := c.ArchiveRecipeStepPreparation(ctx, exampleRecipeStepPreparation.BelongsToRecipeStep, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
