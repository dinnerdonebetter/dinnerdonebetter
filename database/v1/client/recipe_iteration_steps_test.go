package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeIterationStepExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("RecipeIterationStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(true, nil)

		actual, err := c.RecipeIterationStepExists(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(exampleRecipeIterationStep, nil)

		actual, err := c.GetRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStep, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeIterationStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("GetAllRecipeIterationStepsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeIterationStepsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeIterationSteps(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := models.DefaultQueryFilter()
		exampleRecipeIterationStepList := fakemodels.BuildFakeRecipeIterationStepList()

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("GetRecipeIterationSteps", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeIterationStepList, nil)

		actual, err := c.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStepList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeIterationStepList := fakemodels.BuildFakeRecipeIterationStepList()

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("GetRecipeIterationSteps", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeIterationStepList, nil)

		actual, err := c.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStepList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("CreateRecipeIterationStep", mock.Anything, exampleInput).Return(exampleRecipeIterationStep, nil)

		actual, err := c.CreateRecipeIterationStep(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStep, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		c, mockDB := buildTestClient()

		mockDB.RecipeIterationStepDataManager.On("UpdateRecipeIterationStep", mock.Anything, exampleRecipeIterationStep).Return(expected)

		err := c.UpdateRecipeIterationStep(ctx, exampleRecipeIterationStep)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeIterationStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationStepDataManager.On("ArchiveRecipeIterationStep", mock.Anything, exampleRecipeIterationStep.BelongsToRecipe, exampleRecipeIterationStep.ID).Return(expected)

		err := c.ArchiveRecipeIterationStep(ctx, exampleRecipeIterationStep.BelongsToRecipe, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
